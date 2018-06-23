package main

import (
	"github.com/spf13/viper"
)

var (
	configuration *Config
)

type Config struct {
	MongoUrl        string
}

func readConfiguration() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	configuration = &Config{
		MongoUrl:                viper.GetString("mongo.url"),
	}
}
