package util

import (
	"fmt"
	"io/ioutil"
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
