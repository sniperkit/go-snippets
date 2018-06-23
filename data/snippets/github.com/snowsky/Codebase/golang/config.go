package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ednsConfig struct {
	DNSServer   string `json:"dns_server" yaml:"dns_server"`
	RedisServer string `json:"redis_server" yaml:"redis_server"`
}

// type ednsConfig struct {
// 	Firewall_network_rules map[string][]string
// }

func main() {
	filename, _ := filepath.Abs("./edns.yml.bak")
	yamlFile, err := ioutil.ReadFile(filename)
	// fmt.Println(yamlFile)
	if err != nil {
		fmt.Println(err)
	}

	var config ednsConfig
	err = yaml.Unmarshal(yamlFile, &config)
	fmt.Println(config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Value: %#v\n", config.DNSServer)
}
