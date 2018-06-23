package config

import (
	"fmt"
	"encoding/json"
	"os"
)

type Configuration struct {
	Port string
	Database string
}

var fileName = "./config.json"

func LoadConfig () *Configuration {
	configuration := Configuration{}

	if _, err := os.Stat(fileName); err == nil {
		file, _ := os.Open(fileName)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	return &configuration
}
