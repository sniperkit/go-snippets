package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig("test_config/launcher.toml")
	if err != nil {
		t.Error(err)
	}
	if config.Hosts["localhost"].Name != "127.0.0.1" {
		t.Errorf("Expected host name to be 127.0.0.1, but got %v", config.Hosts["localhost"].Name)
	}
	if config.Hosts["localhost"].User != "user" {
		t.Errorf("Expected host user to be user, but got %v", config.Hosts["localhost"].User)
	}
	if config.Scripts["ls_lh"].Content != "ls -lh && sleep 1 && ls -lh" {
		t.Errorf("Expected ls_lh script content to be 'ls -lh && sleep 1 && ls -lh', but got %v",
			config.Scripts["ls_lh"].Content)
	}
}
