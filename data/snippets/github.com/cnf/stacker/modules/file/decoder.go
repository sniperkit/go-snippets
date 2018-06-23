package file

import (
	"path/filepath"

	"github.com/cnf/structhash"
	"github.com/mitchellh/mapstructure"

	"github.com/cnf/stacker/engine"
)

// Config ...
type Config struct {
	LogOutput string `mapstructure:"log_output" version:"1"`
	ConfigDir string `mapstructure:"config_dir" version:"1"`
	Watch     bool   `mapstructure:"watch" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "module",
		Module: "file",
	}
}

// NewConfig returns a stacker Config with defaults
func NewConfig() *Config {
	return &Config{
		LogOutput: "stderr",
		ConfigDir: "/etc/stacker/conf.d",
		Watch:     false,
	}
}

// Decoder ...
type Decoder struct{}

// Version ...
func (sc *Config) Version() int {
	return 1
}

// Hash returns the unique hash for this config object
func (sc *Config) Hash() string {
	hash, err := structhash.Hash(sc, engine.LatestHashVersion)
	if err != nil {
		return ""
	}
	return hash
}

// RegisterConfigDecoder with the engine.
func RegisterConfigDecoder() {
	engine.RegisterConfigDecoder("stacker", &Decoder{})
}

// DecodeConfig returns a parsed config object
func (d *Decoder) DecodeConfig(data interface{}) (engine.Config, error) {
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
	cdp, err := filepath.Abs(cfg.ConfigDir)
	if err != nil {
		return nil, err
	}
	cfg.ConfigDir = filepath.Clean(cdp)
	return cfg, nil
}
