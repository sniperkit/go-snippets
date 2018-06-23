// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config takes care of the configuration file parsing.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DevMine/srctool/log"
)

// Configuration constants
const (
	ConfigFolder     = "srctool"      // Configuration folder name
	DataFolder       = "srctool"      // Data folder name
	ParsersFolder    = "parsers"      // Parsers folder name
	ConfigFileName   = "srctool.conf" // Configuration file name
	ChecksumFileName = "MD5SUM"       // Checksum file name

	// DefaultConfigDir is the default configuration directoy when
	// $XDG_CONFIG_HOME is not set.
	DefaultConfigDir = ".config"

	// DefaultDataDir is the default data directoy when $XDG_DATA_HOME
	// is not set.
	DefaultDataDir = ".local/share"
)

const (
	archExt = ".zip" // archive extension
)

// default config file
const defaultConfigFile = `{
	"download_server_url": "http://dl.devmine.ch/parsers"
}`

// Config holds the configuration of srctool.
type Config struct {
	DownloadServerURL string `json:"download_server_url"`
}

// New creates a new Config initialized with the values defined in the
// configuration file located in $XDG_CONFIG_HOME/srctool/srctool.conf.
// If $XDG_CONFIG_HOME is not set, it uses the directory "$HOME/.config/" as
// config home. If some files or directories do not already exist, it creates
// them automatically.
func New() (*Config, error) {
	if err := createConfigDir(); err != nil {
		return nil, err
	}

	if err := createDataDir(); err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadFile(filepath.Join(ConfigDir(), ConfigFileName))
	if err != nil {
		log.Debug("config:", err)
		return nil, errors.New("cannot read configuration file")
	}

	cfg := new(Config)
	if err = json.Unmarshal(bs, &cfg); err != nil {
		log.Debug("config:", err)
		return nil, errors.New("malformed configuration file")
	}

	if err = cfg.verify(); err != nil {
		return nil, fmt.Errorf("config: %v", err)
	}

	return cfg, nil
}

// createConfigDir creates the configuration directory if it does not already
// exists.
func createConfigDir() error {
	var err error

	if _, err = os.Stat(ConfigDir()); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("config dir '%s' does not exist", ConfigDir()))
		log.Info("creating the data dir...")

		if err = os.MkdirAll(ConfigDir(), 0755); err != nil {
			return err
		}

		log.Success("config dir successfully created")
	} else if err != nil {
		return err
	}

	if err = createConfigFile(); err != nil {
		return err
	}

	return nil
}

// createDataDir creates the data directory if it does not already exists.
func createDataDir() error {
	var err error
	path := filepath.Join(DataDir(), ParsersFolder)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("data dir '%s' does not exist", path))
		log.Info("creating the data dir...")

		if err = os.MkdirAll(path, 0755); err != nil {
			return err
		}

		log.Success("data dir successfully created")
	} else if err != nil {
		return err
	}

	return nil
}

// createConfigFile creates a default configuration file if it does not already
// exists.
func createConfigFile() error {
	filePath := filepath.Join(ConfigDir(), ConfigFileName)

	var err error

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		if err = ioutil.WriteFile(filePath, []byte(defaultConfigFile), 0644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// Save saves the Config on the disk.
func (c Config) Save() error {
	bs, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Debug(err)
		return errors.New("unable to save config file")
	}

	configPath := filepath.Join(ConfigDir(), ConfigFileName)

	fi, err := os.Stat(configPath)
	if err != nil {
		log.Debug(err)
		return errors.New("unable to save config file")
	}

	if err = ioutil.WriteFile(configPath, bs, fi.Mode()); err != nil {
		log.Debug(err)
		return errors.New("unable to save config file")
	}

	return nil
}

// verify the correctness of the config values
func (c Config) verify() error {
	if _, err := url.Parse(c.DownloadServerURL); err != nil {
		return errors.New(fmt.Sprint("config:", err))
	}

	return nil
}

// ConfigDir returns the configuration directory of srctool.
func ConfigDir() string {
	configHome := filepath.Join(os.Getenv("HOME"), DefaultConfigDir)
	if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		configHome = xdg
	}

	return filepath.Join(configHome, ConfigFolder)
}

// DataDir returns the data directory of srctool.
func DataDir() string {
	dataHome := filepath.Join(os.Getenv("HOME"), DefaultDataDir)
	if xdg := os.Getenv("XDG_DATA_HOME"); len(xdg) > 0 {
		dataHome = xdg
	}

	return filepath.Join(dataHome, DataFolder)
}

// ConfigFilePath returns the path of the configuration file.
func ConfigFilePath() string {
	return filepath.Join(DataDir(), ConfigFileName)
}

// ParsersDir returns the path of the parsers directory.
func ParsersDir() string {
	return filepath.Join(DataDir(), ParsersFolder)
}

// ParserPath return the local path of a parser.
func ParserPath(parserName string) string {
	return filepath.Join(ParsersDir(), parserName)
}

// RemoteChecksumsPath returns the path of the remotes checksums file.
func RemoteChecksumsPath(serverURL string) string {
	url := serverURL
	if url[len(url)-1] != '/' {
		url += "/"
	}

	return url + "MD5SUMS"
}

// LocalChecksumPath returns the path of the checksum file for a given parser.
func LocalChecksumPath(parserName string) string {
	return filepath.Join(ParserPath(parserName), "MD5SUM")
}

// TempPath returns the temporary path of the compressed parser.
func TempPath(parserName string) string {
	return filepath.Join(os.TempDir(), parserName+archExt)
}

// RemoteParserPath returns the remote parser path.
func RemoteParserPath(parserName string) string {
	return filepath.Join(runtime.GOOS, runtime.GOARCH, parserName+archExt)
}

// ParserURI returns the full URI of a parser on the download server.
func ParserURI(serverURL, parserName string) string {
	url := serverURL
	if url[len(url)-1] != '/' {
		url += "/"
	}

	return url + RemoteParserPath(parserName)
}
