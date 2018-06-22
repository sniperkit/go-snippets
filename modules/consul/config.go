package consul

import (
	"github.com/cnf/structhash"
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/engine"
)

// Config object for consul
type Config struct {
	Address    string `mapstructure:"address" version:"1"`
	Datacenter string `mapstructure:"datacenter" version:"1"`
	WaitTime   string `mapstructure:"wait_time" version:"1"`
	Token      string `mapstructure:"token" version:"1"`

	Discovery bool `mapstructure:"discovery" version:"1"`
	Exposed   bool `mapstructure:"exposed" version:"1"`
	Published bool `mapstructure:"published" version:"1"`

	Confstore bool `mapstructure:"confstore" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "module",
		Module: "consul",
	}
}

// NewConfig returns a Consul config with defaults
func NewConfig() *Config {
	return &Config{
		Address:   "http://172.17.42.1:8500/",
		Discovery: false,
		Exposed:   false,
		Published: false,
		Confstore: false,
	}
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
	engine.RegisterConfigDecoder("consul", &ConfigDecoder{})
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

// ValidateConfig returns a validated object
// TODO: remove
func (c *Config) ValidateConfig(cf *Config) (*Config, error) {
	return cf, nil
}
