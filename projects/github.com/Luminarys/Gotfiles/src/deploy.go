package main

import (
	"os/exec"
)

func deploy(path string, force bool) error {
	cmd := "rsync"
	var args []string
	if force {
		args = []string{"-r", "--exclude", ".git/", path, homeDir + "/"}
	} else {
		args = []string{"-r", "--exclude", ".git/", "--ignore-existing", path, homeDir + "/"}
	}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return err
	}
	return nil
}
