package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

func config() map[string]map[string]string {
	content, err := ioutil.ReadFile(filepath.Join("..", "config.json"))
	if err != nil {
		log.Fatalln(err)
	}
	config := make(map[string]map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
