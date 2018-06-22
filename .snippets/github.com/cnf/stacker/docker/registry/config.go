package registry

import (
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

const (
	// INDEXSERVER is the default docker index server
	INDEXSERVER = "https://index.docker.io/v1/"
	// REGISTRYSERVER ...
	REGISTRYSERVER = "https://registry-1.docker.io/v1/"
)

// Config structure for a Registry
type Config struct {
	Name          string `mapstructure:"-"`
	Username      string `mapstructure:"username" version:"1"`
	Password      string `mapstructure:"password" version:"1"`
	Email         string `mapstructure:"email" version:"1"`
	ServerAddress string `mapstructure:"address" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "crt",
		Module: "registry",
		ID:     c.Name,
	}
}

// NewConfig returns a RegistryConfig with defaults
func NewConfig() *Config {
	return &Config{}
}

// DecodeConfig returns a parsed config object
func DecodeConfig(cont []toml.Primitive) ([]*Config, error) {
	cfg := []*Config{}
	for _, v := range cont {
		config := NewConfig()
		if err := toml.PrimitiveDecode(v, &config); err != nil {
			glog.Error(err)
			continue
		}

		name, err := nameFromServerAddress(config.ServerAddress)
		if err != nil {
			continue
		}
		config.Name = name

		cfg = append(cfg, config)
	}
	return cfg, nil
}
