package cmd

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const DASHFILE_JSON = "dashfiles.json"

type Dashfile struct {
	Dotfiles string `json:"dotfiles"`
}

func DefaultDashfile() *Dashfile {
	return &Dashfile{
		Dotfiles: "dotfiles",
	}
}

func (dashFile *Dashfile) readFromWorkspace(workspace string) error {
	file := filepath.Join(workspace, DASHFILE_JSON)
	if fileExists(file) {
		if data, err := ioutil.ReadFile(file); err != nil {
			return err
		} else if err := json.Unmarshal(data, dashFile); err != nil {
			return err
		}
	}
	return nil
}
