package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

const (
	channelYAMLOldFormat = `name: "Angie Speaks"
slug: "angiespeaks"
url: "https://www.youtube.com/channel/UCUtloyZ_Iu4BJekIqPLc_fQ"
description: "Anarchist Leftist channel with a creative and mystical flair!"
subscribers: 8367
tags: ["breadtube"]`
	channelYAMLNewFormat = `name: "anarchopac"
slug: "anarchopac"
description: "I'm a disabled pan-sexual trans woman who talks about anarchism, feminism and marxism."
subscribers:
providers:
  youtube:
    name: "anarchopac"
    url: "https://www.youtube.com/user/anarchopac"
    subscribers: 15438
    description: ""
  patreon:
    name: "anarchopac"
    url: "https://www.patreon.com/anarchopac"
    description: "left-wing youtube videos"
tags: ["breadtube"]`
)

func TestChannelUnmarshalYAML_OldFormat(t *testing.T) {
	channel := Channel{}

	err := yaml.Unmarshal([]byte(channelYAMLOldFormat), &channel)
	assert.NoError(t, err)
	assert.Equal(t, "Angie Speaks", channel.Name)
	assert.Equal(t, "angiespeaks", channel.Slug)
	assert.Len(t, channel.Providers, 0)

	require.Len(t, channel.remnant, 3)
	url := channel.remnant["url"]
	assert.Equal(t, url, "https://www.youtube.com/channel/UCUtloyZ_Iu4BJekIqPLc_fQ")

	description := channel.remnant["description"]
	assert.Equal(t, "Anarchist Leftist channel with a creative and mystical flair!", description)

	subscribers := channel.remnant["subscribers"]
	assert.Equal(t, 8367, subscribers)
}

func TestChannelUnmarshalYAML_NewFormat(t *testing.T) {
	channel := Channel{}

	err := yaml.Unmarshal([]byte(channelYAMLNewFormat), &channel)
	assert.NoError(t, err)
	assert.Equal(t, "anarchopac", channel.Slug)
	require.Len(t, channel.Providers, 2)

	youtubeProvider := channel.Providers["youtube"]
	assert.Equal(t, "anarchopac", youtubeProvider.Name)
	assert.Equal(t, MustParseURL("https://www.youtube.com/user/anarchopac"), youtubeProvider.URL)
	assert.Equal(t, uint64(15438), youtubeProvider.Subscribers)

	patreonProvider := channel.Providers["patreon"]
	assert.Equal(t, MustParseURL("https://www.patreon.com/anarchopac"), patreonProvider.URL)

	// subscribers is omitted because it's nil
	assert.Len(t, channel.remnant, 1)
}

func TestChannelUnmarshalYAML_Nils(t *testing.T) {
	channel := Channel{}

	err := yaml.Unmarshal([]byte(`name: "anarchopac"
slug: nil
description:
subscribers:
providers:
  youtube:
    name:
    url: null
    subscribers:
    description:
`), &channel)
	assert.NoError(t, err)
}

func TestChannelYouTubeURL(t *testing.T) {
	channel := Channel{}

	err := yaml.Unmarshal([]byte(channelYAMLNewFormat), &channel)
	assert.NoError(t, err)
	assert.Equal(t, MustParseURL("https://www.youtube.com/user/anarchopac"), channel.YouTubeURL())

	err = yaml.Unmarshal([]byte(channelYAMLOldFormat), &channel)
	assert.NoError(t, err)
	assert.Equal(t, MustParseURL("https://www.youtube.com/channel/UCUtloyZ_Iu4BJekIqPLc_fQ"), channel.YouTubeURL())
}

func TestChannelUnmarshalYAML_NilsURLFallback(t *testing.T) {
	channel := Channel{}

	err := yaml.Unmarshal([]byte(`name: "anarchopac"
slug: nil
description:
subscribers:
url: http://youtube.com/channel/cancelled
providers:
  youtube:
    name:
    url: null
    subscribers:
    description:
`), &channel)
	assert.NoError(t, err)
	assert.Equal(t, MustParseURL("http://youtube.com/channel/cancelled"), channel.YouTubeURL())
}
