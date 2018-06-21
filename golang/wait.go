package main

import (
	"os"
	"os/exec"
)

func wait(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func main() {
	wait("ls", "-al")
	wait("ps", "aelf")
}
