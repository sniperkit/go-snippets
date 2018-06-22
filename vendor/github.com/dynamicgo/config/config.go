package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Jeffail/gabs"
)

// Value config item value
type Value interface{}

// Config golang json config object
type Config struct {
	container *gabs.Container
}

// New load config from json bytes
func New(source []byte) (*Config, error) {
	parsed, err := gabs.ParseJSON(source)

	if err != nil {
		return nil, err
	}

	return &Config{
		container: parsed,
	}, nil
}

// NewFromFile load config from json file
func NewFromFile(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	return New(data)
}

// Reload reload config from source bytes
func (config *Config) Reload(source []byte) error {
	parsed, err := gabs.ParseJSON(source)

	if err != nil {
		return err
	}

	config.container = parsed

	return nil
}

// Get get config value
func (config *Config) Get(path string) Value {
	return config.container.Path(path).Data()
}

// GetObject get config value as object
func (config *Config) GetObject(path string, v interface{}) error {

	value, ok := config.tryGet(path)

	if !ok {
		return fmt.Errorf("config %s not found", path)
	}

	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func (config *Config) tryGet(path string) (Value, bool) {

	if config.Has(path) {
		return config.container.Path(path).Data(), true
	}

	return nil, false
}

// GetInt64 get config value as Int
func (config *Config) GetInt64(path string, defaultval int64) int64 {

	value, ok := config.tryGet(path)

	if !ok {
		return defaultval
	}

	if val, ok := value.(float64); ok {
		return int64(val)
	}

	return defaultval
}

// GetBool get config value as bool
func (config *Config) GetBool(path string, defaultval bool) bool {

	value, ok := config.tryGet(path)

	if !ok {
		return defaultval
	}

	if val, ok := value.(bool); ok {
		return val
	}

	return defaultval
}

// GetDuration fetch config value as time.Duration
func (config *Config) GetDuration(path string, defaultval time.Duration) time.Duration {
	return time.Duration(config.GetInt64(path, int64(defaultval)))
}

// GetString get config value as Int
func (config *Config) GetString(path string, defaultval string) string {
	value, ok := config.tryGet(path)

	if !ok {
		return defaultval
	}

	if val, ok := value.(string); ok {
		return val
	}

	return defaultval
}

// Has check if has config item indicate by path
func (config *Config) Has(path string) bool {
	return config.container.ExistsP(path)
}

// String print config as string
func (config *Config) String() string {
	return config.container.String()
}

var config = &Config{}

// GlobalConfig global builtin config
var GlobalConfig = config

// Load load config from source bytes
func Load(source []byte) {
	config.Reload(source)
}

// LoadFromFile load config from config json file
func LoadFromFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	config.Reload(data)

	return nil
}

// Get global method, get config value from global config object
func Get(path string) Value {
	return config.Get(path)
}

// GetBool get config value as bool
func GetBool(path string, defaultval bool) bool {
	return config.GetBool(path, defaultval)
}

// Has global method, check if global config object has the config item
func Has(path string) bool {
	return config.Has(path)
}

// GetInt64 get config value as Int
func GetInt64(path string, defaultval int64) int64 {
	return config.GetInt64(path, defaultval)
}

// GetDuration get config value as Int
func GetDuration(path string, defaultval time.Duration) time.Duration {
	return config.GetDuration(path, defaultval)
}

// GetString get config value as String
func GetString(path string, defaultval string) string {
	return config.GetString(path, defaultval)
}

// GetObject get config value as object
func GetObject(path string, v interface{}) error {
	return config.GetObject(path, v)
}
