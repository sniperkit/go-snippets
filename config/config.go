// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config takes care of the configuration file parsing.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

// sslModes corresponds to the SSL modes available for the connection to the
// PostgreSQL database.
// See http://www.postgresql.org/docs/9.4/static/libpq-ssl.html for details.
var sslModes = map[string]bool{
	"disable":     true,
	"require":     true,
	"verify-ca":   true,
	"verify-full": true,
}

// Config is the main configuration structure.
type Config struct {
	Database *DatabaseConfig `json:"database"`
}

// DatabaseConfig is a configuration for PostgreSQL database connection
// information
type DatabaseConfig struct {
	HostName string `json:"hostname"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`

	// Can take values: disable, require, verify-ca or verify-full
	SSLMode string `json:"ssl_mode"`
}

// ReadConfig reads a JSON formatted configuration file, verifies the values
// of the configuration parameters and fills the Config structure.
func ReadConfig(path string) (*Config, error) {
	if len(path) == 0 {
		return &Config{}, nil
	}

	// TODO maybe use a safer function like io.Copy
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := json.Unmarshal(bs, cfg); err != nil {
		return nil, err
	}

	if err := cfg.verify(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) verify() error {

	err := c.Database.verify()
	if err != nil {
		return err
	}

	return nil
}

func (dc DatabaseConfig) verify() error {
	if len(strings.Trim(dc.HostName, " ")) == 0 {
		return errors.New("database hostname cannot be empty")
	}

	if dc.Port <= 0 {
		return errors.New("database port must be greater than 0")
	}

	if len(strings.Trim(dc.UserName, " ")) == 0 {
		return errors.New("database username cannot be empty")
	}

	if len(strings.Trim(dc.DBName, " ")) == 0 {
		return errors.New("database name cannot be empty")
	}

	if _, ok := sslModes[dc.SSLMode]; !ok {
		return errors.New("database can only be disable, require, verify-ca or verify-full")
	}

	return nil
}
