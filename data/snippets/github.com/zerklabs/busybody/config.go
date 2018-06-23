package busybody

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

const DefaultSwimInterval = "2m0s"
const DefaultSwimTimeout = "1m0s"

type BusyConfig struct {
	DeflateCompression      bool          `toml:"deflate_compression"`
	DeflateCompressionLevel int           `toml:"deflate_compression_level"`
	LogLevel                int           `toml:"log_level"`
	Peers                   []string      `toml:"peers"`
	SharedKey               string        `toml:"shared_key"`
	SnappyCompression       bool          `toml:"snappy_compression"`
	SwimIntervalStr         string        `toml:"swim_interval"`
	SwimTimeoutStr          string        `toml:"swim_timeout"`
	SwimInterval            time.Duration `toml:"-"`
	SwimTimeout             time.Duration `toml:"-"`
	Uri                     string        `toml:"uri"`
	ZlibCompression         bool          `toml:"zlib_compression"`
}

func ParseConfig(config []byte) (*BusyConfig, error) {
	var conf BusyConfig
	if _, err := toml.Decode(string(config), &conf); err != nil {
		return nil, err
	}

	if conf.Uri == "" {
		return nil, fmt.Errorf("uri required in config")
	}

	if conf.SwimIntervalStr == "" {
		conf.SwimIntervalStr = DefaultSwimInterval
	}

	if conf.SwimTimeoutStr == "" {
		conf.SwimTimeoutStr = DefaultSwimTimeout
	}

	conf.SwimTimeout, _ = time.ParseDuration(DefaultSwimTimeout)
	conf.SwimInterval, _ = time.ParseDuration(DefaultSwimInterval)

	if conf.SnappyCompression && conf.DeflateCompression && conf.ZlibCompression {
		return nil, fmt.Errorf("only one of snappy, deflate or zlib can be used")
	}

	return &conf, nil
}
