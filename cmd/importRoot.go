package cmd

import (
	"github.com/spf13/cobra"
)

// importRootCmd represents the importRoot command
var importRootCmd = &cobra.Command{
	Use:   "import",
	Short: "Import resources, currently works with video",
	Long: `Command to import resources. Not to be confused with the channel import command.

	Currently the only valid resource for import is the video resource. However,
	this command will later form the basis for importing all resources.`,
	ValidArgs: []string{"video"},
	Args:      cobra.ExactArgs(1),
	Run:       func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(importRootCmd)
}
