package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	envVars, err := parseJSONFile("config.dev.json")
	if err != nil {
		log.Fatal("Failed to parse json config", err)
	}
	for varKey := range envVars {
		log.Printf("Set key %s and value %s", varKey, envVars[varKey])
		err := os.Setenv(varKey, envVars[varKey])
		if err != nil {
			log.Println("Failed to set env var")
		}
	}
}

// all parsed data will converted into map[string]string so it can be processed fruther using os.SetEnv
func parseJSONFile(filepath string) (map[string]string, error) {
	jsonContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	vars := make(map[string]interface{})
	err = json.Unmarshal(jsonContent, &vars)
	if err != nil {
		return nil, err
	}

	sanitzedVars := make(map[string]string)
	for key := range vars {
		switch vars[key].(type) {
		case string:
			// do something with string
			sanitzedVars[key] = vars[key].(string)
		default:
			// do something if not string
			return nil, fmt.Errorf("Key %s have non 'string' value data type. Value is: %v", key, vars[key])
		}
	}
	return sanitzedVars, nil
}
