package managers

import (
	"os"
	"os/exec"
)

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd
}
