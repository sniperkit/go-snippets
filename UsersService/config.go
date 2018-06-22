package main

import (
	"github.com/spf13/viper"
	"fmt"
	"strconv"
)

type Config interface {
	GetValue(string) string
	GetIntValue(string) int
}

type configuration map[string]interface{}

var config configuration

type BaseConfig struct {
}

func (self *BaseConfig) Load() {
	self.LoadWithOptions(map[string]interface{}{})
}
func (self *BaseConfig) LoadWithOptions(options map[string]interface{}) {
	viper.SetDefault("port", "3000")
	viper.SetDefault("log_level", "warn")
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	if options["configPath"] != nil {
		viper.AddConfigPath(options["configPath"].(string))
	} else {
		viper.AddConfigPath("./")
		viper.AddConfigPath("../")
	}
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	config = configuration{}
}

func (self *BaseConfig) GetValue(key string) string {
	if _, ok := config[key]; !ok {
		config[key] = getStringOrPanic(key)
	}
	return config[key].(string)
}

func (self BaseConfig) GetIntValue(key string) int {
	if _, ok := config[key]; !ok {
		config[key] = getIntOrPanic(key)
	}
	return config[key].(int)
}

func (self *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	if _, ok := config[key]; !ok {
		var value string
		if value = viper.GetString(key); !viper.IsSet(key) {
			value = defaultValue
		}
		config[key] = value
	}
	return config[key].(string)
}


func getStringOrPanic(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func getIntOrPanic(key string) int {
	checkKey(key)
	v, err := strconv.Atoi(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return v
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("%s key is not set", key))
	}
}


func panicIfErrorForKey(err error, key string) {
	if err != nil {
		panic(fmt.Errorf("Could not parse key: %s. Error: %v", key, err))
	}
}
