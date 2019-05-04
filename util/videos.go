package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
