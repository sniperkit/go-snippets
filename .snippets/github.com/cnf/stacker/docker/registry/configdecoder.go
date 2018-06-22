package registry

import (
	"github.com/cnf/structhash"
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/engine"
)

// ConfigDecoder object
type ConfigDecoder struct{}

// Version of our config object
// TODO: automate
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
	engine.RegisterConfigDecoder("registry", &ConfigDecoder{})
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
	// TODO: check valid name
	// if cfg.Name == "" || !checkValidName(cfg.Name) {
	// 	glog.Warningf("not a valid name: '%s'", cfg.Name)
	// 	return nil, engine.ErrNotAValidConfig
	// }
	return cfg, nil
	// TODO: actually do stuff
}
