// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/breadtubetv/bake/util"
	"github.com/spf13/cobra"
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Import a video by ID",
	Long:  `Import a YouTube video by ID and assign it to a creator.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := util.ImportVideo(id, creator)
		if err != nil {
			log.Fatalf("could not import video: %v", err)
		}
	},
}

var (
	id       string
	creator  string
	provider string
)

func init() {
	importRootCmd.AddCommand(videoCmd)

	videoCmd.Flags().StringVar(&id, "id", "", "ID of the video, e.g. xspEtjnSfQA is the ID for https://www.youtube.com/watch?v=xspEtjnSfQA")
	videoCmd.Flags().StringVarP(&creator, "creator", "c", "", "Creator slug for the imported video")
	videoCmd.Flags().StringVarP(&provider, "provider", "p", "", "Video provider to import from - e.g. youtube")

	videoCmd.MarkFlagRequired("id")
	videoCmd.MarkFlagRequired("creator")
	videoCmd.MarkFlagRequired("provider")
}
