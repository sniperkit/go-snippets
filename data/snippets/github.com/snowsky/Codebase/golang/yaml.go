package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Firewall_network_rules map[string][]string
}

type Options struct {
	Src string
	Dst string
}

func main() {
	filename, _ := filepath.Abs("./edns.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	fmt.Println(yamlFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", config.Firewall_network_rules)
}
