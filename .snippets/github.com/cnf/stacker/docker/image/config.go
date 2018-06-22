package image

import (
	"fmt"
	"strings"

	"github.com/cnf/stacker/engine"
)

// Config struct for an Image
type Config struct {
	Name     string `mapstructure:"name" version:"1"`
	Registry string `mapstructure:"registry" version:"1"`
	State    string `mapstructure:"state" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "crt",
		Module: "image",
		ID:     c.Name,
	}
}

// NewConfig returns an ImageConfig with defaults
func NewConfig() *Config {
	return &Config{
		State: "present",
	}
}

func parseRepoName(name string) (repo, tag string, err error) {
	spos := strings.IndexRune(name, '/')
	cpos := strings.IndexRune(name, ':')
	if cpos != -1 && cpos > spos {
		tag = name[cpos+1:]
		repo = name[:cpos]
		return repo, tag, nil
	}
	// TODO: don't hardcode this
	return name, "latest", nil
}

func sanitizeName(cfg *Config) string {
	// TODO: ...
	repo, tag, err := parseRepoName(cfg.Name)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", repo, tag)

}
