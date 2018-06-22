package logger

import (
	"github.com/cnf/structhash"
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/engine"
)

// Config object for logger
type Config struct {
	Location string `toml:"location" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "module",
		Module: "logger",
	}
}

// NewConfig returns a logger config with defaults
func NewConfig() *Config {
	return &Config{}
}

// ConfigDecoder ...
type ConfigDecoder struct{}

// Version ...
func (c *Config) Version() int {
	return 1
}

// Hash returns the unique hash for this config object
func (c *Config) Hash() string {
	hash, err := structhash.Hash(c, engine.LatestHashVersion)
	if err != nil {
		return ""
	}
	return hash
}

// RegisterConfigDecoder with the engine.
func RegisterConfigDecoder() {
	engine.RegisterConfigDecoder("logger", &ConfigDecoder{})
}

// DecodeConfig returns a parsed config object
func (cd *ConfigDecoder) DecodeConfig(data interface{}) (engine.Config, error) {
	var md mapstructure.Metadata
	cfg := NewConfig()
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   cfg,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(data); err != nil {
		return nil, err
	}
	return cfg, nil
}
