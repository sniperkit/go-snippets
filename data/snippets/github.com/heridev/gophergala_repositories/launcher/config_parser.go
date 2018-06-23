package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Hosts   map[string]*Host
	Scripts map[string]*Script
}

func ParseConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	conf := &Config{}
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
