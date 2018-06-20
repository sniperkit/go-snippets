package config

import (
	"testing"
)

func TestConfiguration_NewConnectionString(t *testing.T) {
	conf := Configuration{}
	conf.Load("application.dist.yaml")
	DSN, err := conf.NewConnectionString()
	if err != nil {
		t.Errorf("%s", err)
	}
	test_DSN := "username:password@tcp(hostname:3306)/database?options"
	if DSN != test_DSN {
		t.Errorf("Provided DSN (%s) doesn't match test DSN (%s)", DSN, test_DSN)
	}
}

func TestConfiguration_NewConnectionString_NoConfLoaded(t *testing.T) {
	conf := Configuration{}
	if connStr, err := conf.NewConnectionString(); err == nil {
		t.Errorf("Error not set to nil, DSN is: %s", connStr)
	}
}
