// Package youtube implements the functions used in combination with Youtube as a provider
package youtube

import (
	"net/url"

	"github.com/gophergala/yps/provider"
)

type (
	yt struct{}
)

var (
	ytb *yt
)

func (y *yt) IsValidURL(u string) bool {
	ur, err := url.Parse(u)
	if err != nil {
		return false
	}

	if ur.Scheme != "https" || ur.Host != "www.youtube.com" {
		return false
	}

	if ur.Path == "/watch" {
		return ur.Query().Get("v") != ""
	} else if ur.Path == "/playlist" {
		return ur.Query().Get("list") != ""
	}

	return false
}

func (y *yt) IsVideo(u string) bool {
	ur, err := url.Parse(u)
	if err != nil {
		return false
	}

	if ur.Scheme != "https" || ur.Host != "www.youtube.com" {
		return false
	}

	if ur.Path == "/watch" {
		return ur.Query().Get("v") != ""
	}

	return false
}

func (y *yt) IsPlaylist(u string) bool {
	ur, err := url.Parse(u)
	if err != nil {
		return false
	}

	if ur.Scheme != "https" || ur.Host != "www.youtube.com" {
		return false
	}

	if ur.Path == "/playlist" {
		return ur.Query().Get("list") != ""
	}

	return false
}

func (y *yt) URLToFile(url string, downloadFunc provider.DownloadFunc) (string, error) {
	_, _ = downloadFunc(url)

	return url, nil
}

// NewYoutube returns the Youtube implementation of the provider.Provider functionality
func NewYoutube() provider.Provider {
	if ytb == nil {
		ytb = &yt{}
	}

	return ytb
}
