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
	Data     DataConfig      `json:"data"`
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

	// CommitsPerTransaction is useful to set how many commits to insert into
	// the database per transaction. Note that this option is only useful if
	// commit deltas are NOT inserted as well. Defaults to 1000000.
	CommitsPerTransaction uint `json:"commits_per_transaction"`
}

// DataConfig is used to specify some data to retrieve or not.
type DataConfig struct {
	// TmpDir can be used to specify a temporary working directory. If
	// left unspecified, the default system temporary directory will be used.
	// If you have a ramdisk, you are advised to use it here.
	TmpDir string `json:"tmp_dir"`

	// TmpDirFileSizeLimit can be used to specify the maximum size in GB of an
	// object to be temporarily placed in TmpDir for processing. Files of size
	// larger than this value will not be processed in TmpDir.
	TmpDirFileSizeLimit float64 `json:"tmp_dir_file_size_limit"`

	CommitDeltas  bool `json:"commit_deltas"`
	CommitPatches bool `json:"commit_patches"`
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

	if cfg.Data.TmpDirFileSizeLimit < 0.01 {
		cfg.Data.TmpDirFileSizeLimit = 0.01
	}

	if cfg.Database.CommitsPerTransaction == 0 {
		cfg.Database.CommitsPerTransaction = 1000000
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

	err = c.Data.verify()
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

	if dc.CommitsPerTransaction == 0 {
		return errors.New("commits per transaction must be positive")
	}

	return nil
}

func (dc DataConfig) verify() error {
	if dc.CommitPatches && !dc.CommitDeltas {
		return errors.New("commit patches may only be specified along with commit deltas")
	}

	return nil
}
