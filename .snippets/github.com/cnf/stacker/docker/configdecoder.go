package docker

import (
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// ConfigDecoder object.
type ConfigDecoder struct{}

// RegisterConfigDecoder with the engine.
func RegisterConfigDecoder() {
	engine.RegisterConfigDecoder("docker", &ConfigDecoder{})
}

// DecodeConfig returns a parsed config object
func (cd *ConfigDecoder) DecodeConfig(data interface{}) (engine.Config, error) {
	var md mapstructure.Metadata
	cfg := dockerwrapper.NewConfig()
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   &cfg,
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
