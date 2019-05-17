package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/breadtubetv/bake/util"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/yaml.v2"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const accessJSON = `{
	"installed": {
		"client_id": "660935947237-ajqve9kv3n0nnhonhnc5j638fsfan31o.apps.googleusercontent.com",
		"project_id": "sacred-dahlia-229511",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_secret": "piloG0CTWO7PVlbwX8drp6xU",
		"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
  }
}`

// LoadYoutube initalises the Youtube service
func LoadYoutube() map[string]interface{} {
	return map[string]interface{}{
		"config":         config,
		"channel_import": importChannel,
		"video_import":   ImportVideo,
	}
}

func config() {
	client := getClient(youtube.YoutubeReadonlyScope)

	_, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	log.Printf("Successfully authenticated and cached credentials.")
}

// FetchDetails returns the YouTube details for a channel
func FetchDetails(channelURL *util.URL) (util.Provider, error) {
	id := path.Base(channelURL.Path)
	category := path.Base(path.Dir(channelURL.Path))

	client := getClient(youtube.YoutubeReadonlyScope)
	service, err := youtube.New(client)
	if err != nil {
		return util.Provider{}, fmt.Errorf("error creating YouTube client: %v", err)
	}

	call := service.Channels.List("snippet,statistics")
	if category == "channel" {
		call = call.Id(id)
	} else {
		call = call.ForUsername(id)
	}
	response, err := call.Do()
	handleError(err, "")

	if len(response.Items) == 0 {
		return util.Provider{}, fmt.Errorf("could not find channel from URL")
	}

	channelName := ""
	channelDescription := ""
	channelSubscriberCount := uint64(0)
	for _, item := range response.Items {
		channelName = item.Snippet.Title
		channelDescription = item.Snippet.Description
		channelSubscriberCount = item.Statistics.SubscriberCount
		break
	}

	return util.Provider{
		Description: channelDescription,
		Name:        channelName,
		URL:         channelURL,
		Slug:        id,
		Subscribers: channelSubscriberCount,
	}, nil
}

func fetchProfileImageURL(url *util.URL) (string, error) {
	id := path.Base(url.Path)
	category := path.Base(path.Dir(url.Path))
	client := getClient(youtube.YoutubeReadonlyScope)

	youtubeSvc, err := youtube.New(client)
	if err != nil {
		return "", fmt.Errorf("fetchProfileImage: Error while creating YouTube service: %v", err)
	}

	call := youtubeSvc.Channels.List("snippet").Fields("items(snippet/thumbnails)")
	if category == "channel" {
		call = call.Id(id)
	} else {
		call = call.ForUsername(id)
	}
	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("Error retrieving channel profile picture, please download manually.\nErr: %v", err.Error())
	}

	imgURL := response.Items[0].Snippet.Thumbnails.Default.Url
	return imgURL, nil
}

func saveImage(imgURL string, slug string, projectRoot string) error {
	resp, err := http.Get(imgURL)
	if err != nil {
		return fmt.Errorf("couldn't retreive image: %v", err)
	}
	defer resp.Body.Close()

	filePath := fmt.Sprintf("%s/static/img/channels/%s.jpg", projectRoot, slug)
	img, _ := os.Create(filePath)
	defer img.Close()

	_, err = io.Copy(img, resp.Body)
	if err != nil {
		return fmt.Errorf("Error saving channel profile picture, please download manually.\nErr: %v", err.Error())
	}
	log.Printf("Saving %s", filePath)
	return nil
}

func formatChannelDetails(slug string, channelURL *util.URL) (util.Channel, error) {
	provider, err := FetchDetails(channelURL)
	if err != nil {
		return util.Channel{}, err
	}

	return util.Channel{
		Name:      provider.Name,
		Slug:      slug,
		Providers: map[string]util.Provider{"youtube": provider},
	}, nil
}

func importChannel(slug string, channelURL *util.URL, projectRoot string) {
	dataDir := path.Join(projectRoot, "/data/channels")
	channelList := util.LoadChannels(dataDir)

	importedChannel, err := formatChannelDetails(slug, channelURL)
	if err != nil {
		log.Fatalf("Error obtaining channel info: %v", err)
	}

	channel, ok := channelList.Find(slug)
	if ok {
		log.Fatalf("Channel with slug '%s' already exists, use update instead.", slug)
	}
	channel.Name = importedChannel.Name
	channel.Slug = importedChannel.Slug
	channel.Permalink = importedChannel.Slug
	channel.Providers = importedChannel.Providers

	log.Printf("Title: %s, Count: %d\n", channel.Name, channel.Providers["youtube"].Subscribers)
	imgURL, err := fetchProfileImageURL(channelURL)
	if err == nil {
		// TrimSuffix will need to be switched to projectRoot when we merge #185
		err = saveImage(imgURL, slug, strings.TrimSuffix(dataDir, "/data/channels"))
	}

	if err != nil {
		log.Println(err.Error())
	}

	err = util.SaveChannel(channel, dataDir)
	if err != nil {
		log.Fatalf("Error saving channel '%s': %v", slug, err)
	}

	_ = util.CreateChannelVideoFolder(channel, projectRoot)

	err = util.CreateChannelPage(channel, projectRoot)
	if err != nil {
		log.Printf("Unable to create channel page for %s, please create manually.", slug)
	}
}

// ImportVideo will import a YouTube video based on an ID and create
// a new file in the videos data folder for the specified creator
func ImportVideo(id, creator, projectRoot string) error {
	channel, ok := util.LoadChannels(projectRoot + "/data/channels").Find(creator)
	if !ok {
		log.Fatalf("creator %v not found", creator)
	}

	creatorDir := fmt.Sprintf("%s/data/videos/%s", projectRoot, creator)
	if _, err := os.Stat(creatorDir); os.IsNotExist(err) {
		err := util.CreateChannelVideoFolder(channel, projectRoot)
		if err != nil {
			log.Fatalf("unable to create folder for %v: %v", creator, err)
		}
	}

	vid, err := getVideo(id)
	vid.Channel = creator
	if err != nil {
		return err
	}

	videoFile := fmt.Sprintf("%s/%s.yml", creatorDir, vid.ID)
	f, err := os.Create(videoFile)
	if err != nil {
		return fmt.Errorf("could not create file for video '%s': %v", id, err)
	}
	defer f.Close()

	data, err := yaml.Marshal(vid)
	if err != nil {
		return fmt.Errorf("couldn't marshal video data: %v", err)
	}
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("couldn't write to file: %s: %v", f.Name(), err)
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("couldn't sync file: %s: %v", f.Name(), err)
	}
	log.Printf("created video file %v", videoFile)

	err = channel.GetChannelPage(projectRoot).AddVideo(id, projectRoot)
	if err != nil {
		return fmt.Errorf("couldn't update channel page: %v", err)
	}

	return nil
}

// Video represents the a YouTube video
type Video struct {
	ID          string `yaml:"id"`
	Title       string
	Description string
	Source      string
	Channel     string
}

// GetVideo retreives video details from YouTube
func getVideo(videoID string) (*Video, error) {
	client := getClient(youtube.YoutubeReadonlyScope)
	yt, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %v", err)
	}

	call := yt.Videos.List("snippet").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error calling the YouTube API: %v", err)
	}

	video := &Video{
		ID:          resp.Items[0].Id,
		Title:       resp.Items[0].Snippet.Title,
		Description: resp.Items[0].Snippet.Description,
		Source:      "youtube",
	}

	return video, nil
}

const launchWebServer = true

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(scope string) *http.Client {
	if apiKey := os.Getenv("YOUTUBE_API"); apiKey != "" {
		return &http.Client{
			Transport: &transport.APIKey{Key: apiKey},
		}
	}

	ctx := context.Background()

	b := []byte(accessJSON)

	// If modifying the scope, delete your previously saved credentials
	// at ~/.credentials/youtube-go.json
	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Use a redirect URI like this for a web app. The redirect URI must be a
	// valid one for your OAuth2 credentials.
	config.RedirectURL = "http://localhost:8090"
	// Use the following redirect URI if launchWebServer=false in oauth2.go
	// config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		if launchWebServer {
			fmt.Println("Trying to get token from web")
			tok, err = getTokenFromWeb(config, authURL)
		} else {
			fmt.Println("Trying to get token from prompt")
			tok, err = getTokenFromPrompt(config, authURL)
		}
		if err == nil {
			saveToken(cacheFile, tok)
		}
	}
	return config.Client(ctx, tok)
}

// startWebServer starts a web server that listens on http://localhost:8080.
// The webserver waits for an oauth code in the three-legged auth flow.
func startWebServer() (codeCh chan string, err error) {
	listener, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		return nil, err
	}
	codeCh = make(chan string)

	serve := func() {
		err := http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.FormValue("code")
			codeCh <- code // send code to OAuth flow
			listener.Close()
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
		}))
		if err != nil {
			panic(fmt.Sprintf("error starting server for OAuth callback: %v", err))
		}
	}

	go serve()

	return codeCh, nil
}

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

// Exchange the authorization code for an access token
func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	return tok, nil
}

// getTokenFromPrompt uses Config to request a Token and prompts the user
// to enter the token on the command line. It returns the retrieved Token.
func getTokenFromPrompt(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	var code string
	fmt.Printf("Go to the following link in your browser. After completing "+
		"the authorization flow, enter the authorization code on the command "+
		"line: \n%v\n", authURL)

	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	fmt.Println(authURL)
	return exchangeToken(config, code)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	codeCh, err := startWebServer()
	if err != nil {
		fmt.Printf("Unable to start a web server.")
		return nil, err
	}

	err = openURL(authURL)
	if err != nil {
		log.Fatalf("Unable to open authorization URL in web server: %v", err)
	} else {
		fmt.Println("Your browser has been opened to an authorization URL.",
			" This program will resume once authorization has been provided.")
		fmt.Println(authURL)
	}

	// Wait for the web server to get the code.
	code := <-codeCh
	return exchangeToken(config, code)
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	err = os.MkdirAll(tokenCacheDir, 0700)
	if err != nil {
		return "", err
	}

	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Println("trying to save token")
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Printf("couldn't store token: %v", err)
	}
}
