package container

import (
	"regexp"

	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

var (
	validContainerNameChars   = `[a-zA-Z0-9][a-zA-Z0-9_.-]`
	validContainerNamePattern = regexp.MustCompile(`^/?` + validContainerNameChars + `+$`)
)

// Config object for a container
type Config struct {
	// Create Options
	Name       string   `mapstructure:"name" version:"1"`
	Hostname   string   `mapstructure:"hostname" version:"1"`
	User       string   `mapstructure:"user" version:"1"`
	Memory     string   `mapstructure:"memory" version:"1"`
	CPUShares  int64    `mapstructure:"cpu_shares" version:"1"`
	CPUSet     string   `mapstructure:"cpu_set" version:"1"`
	Attach     []string `mapstructure:"attach" version:"1"`
	Expose     []string `mapstructure:"expose" version:"1"`
	TTY        bool     `mapstructure:"tty" version:"1"`
	Env        []string `mapstructure:"env" version:"1"`
	Cmd        []string `mapstructure:"cmd" version:"1"`
	Image      string   `mapstructure:"image" version:"1"`
	Volumes    []string `mapstructure:"volumes" version:"1"`
	Workdir    string   `mapstructure:"workdir" version:"1"`
	Entrypoint []string `mapstructure:"entrypoint" version:"1"`

	// Start Options
	CapAdd      []string `mapstructure:"cap_add" version:"1"`
	CapDrop     []string `mapstructure:"cap_drop" version:"1"`
	CIDFile     string   `mapstructure:"cid_file" version:"1" lastversion:"1"`
	LXCConf     []string `mapstructure:"lxc_conf" version:"1" lastversion:"1"`
	Privileged  bool     `mapstructure:"privileged" version:"1"`
	Publish     []string `mapstructure:"publish" version:"1"`
	PublishAll  bool     `mapstructure:"publish_all" version:"1"`
	Links       []string `mapstructure:"link" version:"1"`
	DNS         []string `mapstructure:"dns" version:"1"`
	DNSSearch   []string `mapstructure:"dns_search" version:"1"`
	VolumesFrom []string `mapstructure:"volumes_from" version:"1"`
	Net         string   `mapstructure:"net" version:"1"`

	// We set these, no override allowed
	Detach   bool `mapstructure:"-" version:"0" lastversion:"0"`
	SigProxy bool `mapstructure:"-" version:"0" lastversion:"0"`

	// No support in the API?
	Device []string `mapstructure:"device" version:"1"`

	// Stacker options
	Remove       bool     `mapstructure:"remove" version:"1"`
	Restart      string   `mapstructure:"restart" version:"1"`
	Dependencies []string `mapstructure:"dependencies" version:"1"`
	State        string   `mapstructure:"state" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "crt",
		Module: "container",
		ID:     c.Name,
	}
}

// NewConfig returns a ContainerConfig with defaults
func NewConfig() *Config {
	return &Config{
		SigProxy: true,
		Detach:   true,
		Restart:  "always",
		Remove:   true,
	}
}

// validates a container name
func isValidName(name string) bool {
	if name == "" {
		return false
	}
	if !validContainerNamePattern.MatchString(name) {
		glog.Warningf("invalid container name (%s), only %s allowed", name, validContainerNameChars)
		return false
	}
	return true
}

// validates container states
func isValidState(state string) bool {
	validStates := map[string]bool{
		"running": true,
		"created": true,
		"stopped": true,
		"removed": true,
	}
	if validStates[state] {
		return true
	}
	return false
}

func (c *Config) hasValidRestart() bool {
	validRestart := map[string]bool{
		"always":   true,
		"never":    true,
		"on-error": true,
	}
	if validRestart[c.Restart] {
		return true
	}
	return false
}
