package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

/*
 * Constants
 */

var (
	configFileLocation string
	fileLock           sync.Mutex
)

type Config struct {
	Port     string   `json:"port"`
	Keywords []string `json:"keywords"`
}

var appConfig Config

func GetConfig() *Config {
	return &appConfig
}

func Load(filename string) error {

	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return errors.New("Cannot open datafile")
	}

	defer file.Close()

	configFileLocation = filename

	decoder := json.NewDecoder(file)
	decErr := decoder.Decode(&appConfig)
	if decErr != nil {
		return errors.New("Error while decoding items: " + decErr.Error())
	}

	log.Println("Loaded config from file")
	return nil
}

func Save() error {

	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Create(configFileLocation)
	if err != nil {
		return errors.New("Cannot create config file: " + err.Error())
	}
	defer file.Close()

	b, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		return errors.New("Cannot marshal config json: " + err.Error())
	}

	_, err = file.Write(b)
	if err != nil {
		return errors.New("Cannot write config file: " + err.Error())
	}

	return nil
}
