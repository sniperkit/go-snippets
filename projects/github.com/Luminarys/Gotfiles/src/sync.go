package main

import (
	"os/exec"
)

func sync(dir string) error {
	cmd := "rsync"
	args := []string{"-r", "--existing", homeDir + "/", dir}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return err
	}
	return nil
}
