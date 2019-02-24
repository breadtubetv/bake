package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [provider] [channel_url]",
	Short: "Import a Channel into BreadtubeTV",
	Long: fmt.Sprintf(`Add the supplied channel into BreadtubeTV, without having to edit JSON.
	
	Available providers: %s`, strings.Join(ProviderNames(), ", ")),
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var provider = args[0]
		var channelURL, err = url.Parse(args[1])

		if err != nil {
			log.Fatalf("Improperly formatted URL provided '%s': %v", args[1], err)
		}

		if _, ok := Providers[provider]; !ok {
			log.Fatalf("No provider exists called %s", provider)
		}

		log.Printf("Importing %s...\n", channelURL)
		Providers[provider]["channel_import"].(func(*url.URL))(channelURL)
	},
}

func init() {
	channelCmd.AddCommand(importCmd)
}

func importChannel(channelURL string) {
	fmt.Printf("called %s\n", channelURL)
}
