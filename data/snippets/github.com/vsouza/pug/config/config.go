package config

import (
	"github.com/spf13/viper"
)

// Config is an exported struct to store viper environment configuration
type Config struct {
	envConf *viper.Viper
}

// NewConfig is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func NewConfig(env string) (*Config, error) {
	var err error
	c := &Config{}
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName(env)
	v.AddConfigPath("config/")
	err = v.ReadInConfig()
	c.envConf = v
	return c, err
}

// GetString is an exported method that receives a string key and get a
// string configuration in config struct
func (c *Config) GetString(key string) string {
	return c.envConf.GetString(key)
}

// GetInt is an exported method that receives a string key and get a
// configuration int item in config struct
func (c *Config) GetInt(key string) int {
	return c.envConf.GetInt(key)
}
