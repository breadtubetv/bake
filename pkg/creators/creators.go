package creators

import (
	"fmt"

	"github.com/breadtubetv/bake/pkg/providers"
	"github.com/spf13/viper"
)

// Creator refers to a single BreadTube creator.
// A creator can have many providers.
type Creator struct {
	name      string
	permalink string
	slug      string
	tags      []string
	Providers providers.ProviderI
	Provider  map[string]providers.Provider
}

var config *viper.Viper

// AddConfig allows package user to pass a viper configuation to the pkg
func AddConfig(conf *viper.Viper) error {
	config = conf
	return nil
}

// FindCreatorBySlug retreives the creator based on the provided
// slug, loads the creator, and returns its creator object.
func FindCreatorBySlug(slug string) *Creator {
	fmt.Print(slug)
	return &Creator{}
}
