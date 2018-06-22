package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRemoteURL(t *testing.T) {
	tests := map[string]*remoteURL{
		"git@github.com:kalbasit/swm.git": &remoteURL{
			protocol:      "",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: ":",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"git@github.com:kalbasit/swm": &remoteURL{
			protocol:      "",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: ":",
			path:          "kalbasit/swm",
			extension:     "",
		},

		"http://git@github.com/kalbasit/swm.git": &remoteURL{
			protocol:      "http",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"http://git@github.com/kalbasit/swm": &remoteURL{
			protocol:      "http",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     "",
		},
		"https://git@github.com/kalbasit/swm.git": &remoteURL{
			protocol:      "https",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"https://git@github.com/kalbasit/swm": &remoteURL{
			protocol:      "https",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     "",
		},

		"http://github.com/kalbasit/swm.git": &remoteURL{
			protocol:      "http",
			username:      "",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"http://github.com/kalbasit/swm": &remoteURL{
			protocol:      "http",
			username:      "",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     "",
		},
		"https://github.com/kalbasit/swm.git": &remoteURL{
			protocol:      "https",
			username:      "",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"https://github.com/kalbasit/swm": &remoteURL{
			protocol:      "https",
			username:      "",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     "",
		},

		"ssh://git@github.com/kalbasit/swm.git": &remoteURL{
			protocol:      "ssh",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     ".git",
		},
		"ssh://git@github.com/kalbasit/swm": &remoteURL{
			protocol:      "ssh",
			username:      "git",
			hostname:      "github.com",
			pathSeparator: "/",
			path:          "kalbasit/swm",
			extension:     "",
		},
	}

	for url, expected := range tests {
		assert.Equal(t, expected, parseRemoteURL(url))
	}
}

func TestRemoteURLString(t *testing.T) {
	tests := []string{
		"git@github.com:kalbasit/swm.git",
		"git@github.com:kalbasit/swm",
		"http://git@github.com/kalbasit/swm.git",
		"http://git@github.com/kalbasit/swm",
		"https://git@github.com/kalbasit/swm.git",
		"https://git@github.com/kalbasit/swm",
		"http://github.com/kalbasit/swm.git",
		"http://github.com/kalbasit/swm",
		"https://github.com/kalbasit/swm.git",
		"https://github.com/kalbasit/swm",
		"ssh://git@github.com/kalbasit/swm.git",
		"ssh://git@github.com/kalbasit/swm",
	}

	for _, url := range tests {
		assert.Equal(t, url, parseRemoteURL(url).String())
	}
}
