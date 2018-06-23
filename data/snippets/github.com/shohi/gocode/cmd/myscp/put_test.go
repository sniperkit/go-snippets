package main

import (
	"log"
	"testing"
)

func TestLoadSSHConfig(t *testing.T) {
	cfg, err := loadSSHConfig()
	log.Println(err)
	log.Println(cfg)
}
