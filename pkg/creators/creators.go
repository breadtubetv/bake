package creators

import (
	"fmt"

	"github.com/breadtubetv/bake/pkg/providers"
	"github.com/spf13/viper"
)

// Creator refers to a single BreadTube creator.
// A creator can have many providers.
type Creator struct {
	Name      string
	Permalink string
	Slug      string
	Tags      []string
	Providers providers.Providers
}

// FindCreatorBySlug retreives the creator based on the provided
// slug, loads the creator, and returns its creator object.
func FindCreatorBySlug(slug string) *Creator {
	fmt.Print(slug)
	fmt.Println(viper.GetString("projectRoot"))
	return &Creator{}
}

func NewCreator(name, slug string) (*Creator, error) {
	// Check if creator already exists
	// Generate Permalink
	return nil, nil
}

// func (c *Creator) LoadProviders() *Creator
// func (c *Creator) AddTag(tag string) *Creator
