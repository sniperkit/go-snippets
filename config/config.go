package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

var configFile string

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General GeneralConfig
	AliOSS  AliOSSConfig
}

type GeneralConfig struct {
	SnippetFile string `toml:"snippetfile"`
	Editor      string `toml:"editor"`
	Column      int    `toml:"column"`
	SelectCmd   string `toml:"selectcmd"`
}

type AliOSSConfig struct {
	AccessID   string `toml:"access_id"`
	AccessKey  string `toml:"access_key"`
	BucketName string `toml:"bucket_name"`
	Endpoint   string `toml:"endpoint"`
	AutoSync   bool   `toml:"auto_sync"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Query     string
	Delimiter string
	Force     bool
	OneLine   bool
	Color     bool
	Tag       bool
}

func (cfg *Config) Load(file string) error {
	configFile = file
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		cfg.General.SnippetFile = expandPath(cfg.General.SnippetFile)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir, _ := GetDefaultConfigDir()
	cfg.General.SnippetFile = filepath.Join(dir, "data", "snippet.toml")
	_, err = os.Create(cfg.General.SnippetFile)
	if err != nil {
		return err
	}

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		cfg.General.Editor = "vim"
	}
	cfg.General.Column = 40
	cfg.General.SelectCmd = "peco"

	return toml.NewEncoder(f).Encode(cfg)
}

func (cfg *Config) Switch(snippetFile string) error {
	dir, _ := GetDefaultConfigDir()
	cfg.General.SnippetFile = filepath.Join(dir, "data", snippetFile)

	_, err := os.Stat(cfg.General.SnippetFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		_, err := os.Create(cfg.General.SnippetFile)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	return toml.NewEncoder(f).Encode(cfg)
}

func GetDefaultConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "pet")
		}
		dir = filepath.Join(dir, "pet")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "pet")
	}
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}
