package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// CreateChannelVideoFolder creates an empty folder within the videos data folder
func CreateChannelVideoFolder(channel *Channel, projectRoot string) error {
	folder := path.Join(projectRoot, fmt.Sprintf("/data/videos/%s", channel.Slug))

	// Make a video directory with a .gitignore
	videoFolder := os.Mkdir(folder, os.ModePerm)
	// There was a suggestion to gitignore, but when importing I like to immediately start committing videos
	os.OpenFile(path.Join(folder, "TODO.yml"), os.O_RDONLY|os.O_CREATE, 0666)

	return videoFolder
}

// GetCreatorVideos returns a list of video IDs for a given creator
func GetCreatorVideos(slug string, projectRoot string) ([]string, error) {
	folder := path.Join(projectRoot, fmt.Sprintf("/data/videos/%s", slug))

	// Using "channels" for consistency, but will change to creators in refactor.
	videoDir, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var videoIds []string

	for _, file := range videoDir {
		if !file.IsDir() {
			fileName := file.Name()
			videoIds = append(videoIds, strings.TrimSuffix(fileName, filepath.Ext(fileName)))
		}
	}

	return videoIds, nil
}

// ImportVideo will import a YouTube video based on an ID and create
// a new file in the videos data folder for the specified creator
func ImportVideo(id, creator, projectRoot string) error {
	channel, ok := LoadChannels(projectRoot).Find(creator)
	if !ok {
		log.Fatalf("creator %v not found", creator)
	}

	// TODO: Get video description via YouTube API
	// TODO: Add file to data/videos/<creator> dir
	// TODO: Add ID to channel page under 'videos'
	channelDir := fmt.Sprintf("%s/data/videos/%s", projectRoot, creator)
	if _, err := os.Stat(channelDir); os.IsNotExist(err) {
		err := CreateChannelVideoFolder(channel, projectRoot)
		if err != nil {
			log.Fatalf("unable to create folder for %v: %v", creator, err)
		}
	}

	return nil
}
