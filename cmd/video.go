package cmd

import (
	"log"
	"os"
	"regexp"

	"github.com/breadtubetv/bake/pkg/providers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Import a video by ID",
	Long:  `Import a YouTube video by ID and assign it to a creator.`,
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" && url == "" {
			log.Fatal("command must include either video ID or URL")
		}
		if id != "" && url != "" {
			log.Fatal("both video ID and URL provided, expected only one")
		}

		if url != "" {
			re := regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`)
			if match := re.MatchString(url); match {
				subs := re.FindStringSubmatch(url)
				id = subs[1]
			} else {
				log.Fatal("the given URL is not a valid YouTube URL")
			}
		}

		err := providers.ImportVideo(id, creator, os.ExpandEnv(viper.GetString("projectRoot")))
		if err != nil {
			log.Fatalf("could not import video: %v", err)
		}
	},
}

var (
	id       string
	url      string
	creator  string
	provider string
)

func init() {
	importRootCmd.AddCommand(videoCmd)

	videoCmd.Flags().StringVar(&id, "id", "", "ID of the video, e.g. xspEtjnSfQA is the ID for https://www.youtube.com/watch?v=xspEtjnSfQA")
	videoCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the video, e.g. https://www.youtube.com/watch?v=xspEtjnSfQA. Use instead of --id.")
	videoCmd.Flags().StringVarP(&creator, "creator", "c", "", "Creator slug for the imported video")
	videoCmd.Flags().StringVarP(&provider, "provider", "p", "", "Video provider to import from - e.g. youtube")

	videoCmd.MarkFlagRequired("creator")
	videoCmd.MarkFlagRequired("provider")
}
