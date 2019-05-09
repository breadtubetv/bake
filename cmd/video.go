package cmd

import (
	"log"

	"github.com/breadtubetv/bake/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Import a video by ID",
	Long:  `Import a YouTube video by ID and assign it to a creator.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := util.ImportVideo(id, creator, viper.GetString("projectRoot"))
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
