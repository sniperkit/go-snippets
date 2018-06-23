package main

import (
	"github.com/RichardKnop/machinery/v1/config"
)

func loadConfig() *config.Config {
	return config.NewFromYaml("config.yml", true, true)
}
