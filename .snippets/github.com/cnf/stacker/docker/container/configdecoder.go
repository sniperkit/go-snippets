package container

import (
	"github.com/cnf/structhash"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/engine"
)

// ConfigDecoder object
type ConfigDecoder struct{}

// Version returns the version of the config
// TODO: needs actual code
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

// RegisterConfigDecoder registers this config decoder with the engine.
func RegisterConfigDecoder() {
	engine.RegisterConfigDecoder("container", &ConfigDecoder{})
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
	if !isValidName(cfg.Name) {
		glog.Warningf("not a valid name: '%s'", cfg.Name)
		return nil, engine.ErrNotAValidConfig
	}
	for i := range cfg.Dependencies {
		if !isValidName(cfg.Dependencies[i]) {
			glog.Warningf("not a valid name: '%s'", cfg.Dependencies[i])
			return nil, engine.ErrNotAValidConfig
		}
	}
	if !isValidState(cfg.State) {
		glog.Warningf("state can not be '%s'", cfg.State)
		return nil, engine.ErrNotAValidConfig
	}
	if !cfg.hasValidRestart() {
		glog.Warningf("restart can not be '%s'", cfg.State)
		return nil, engine.ErrNotAValidConfig
	}
	return cfg, nil
}
