package videos

import (
	"net/url"
)

type Video struct {
	ID     string
	Title  string
	Source string
	URL    url.URL
}
