package server

import (
	"encoding/json"
	"github.com/go-fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type LConfig struct {
	Data *Config
	path string
}

type Config struct {
	// Urls to block.
	Block []string
	// File System Dir -> url mapping.
	// NOTE: '~' is not supported as special path!
	Dirs []Dir
}

type Dir struct {
	FPath string
	Url   string
}

func LiveConfig(path string) (*LConfig, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(path)
	if err != nil {
		watcher.Close()
		return nil, err
	}
	config := &LConfig{nil, path}
	if err := config.update(); err != nil {
		log.Println("Config update error:", err)
	}
	go watch(config, watcher)
	return config, nil
}

func watch(c *LConfig, watcher *fsnotify.Watcher) {
	for {
		select {
		case evn := <-watcher.Events:
			//log.Println("fsnotify:event:", evn)
			if evn.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				if err := c.update(); err != nil {
					log.Println("Config update error:", err)
				}
			}
		case err := <-watcher.Errors:
			log.Println("fsnotify:error:", err)
		}
	}
}

func (c *LConfig) update() error {
	buffer, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}
	config := Config{}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		return err
	}
	for _, dir := range config.Dirs {
		path, err := filepath.Rel(".", dir.FPath)
		if err != nil {
			log.Println(err)
			// TODO: remove invalid entry
			continue
		}
		dir.FPath = path
		if !strings.HasSuffix(dir.Url, "/") {
			dir.Url += "/"
		}
	}
	// TODO: make the update thread safe!
	c.Data = &config
	return nil
}
