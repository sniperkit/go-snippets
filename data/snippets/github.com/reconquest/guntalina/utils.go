package main

import (
	"fmt"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

func execute(cmdline string) (string, error) {
	parser := shellwords.NewParser()
	parser.ParseBacktick = false
	parser.ParseEnv = true

	args, err := parser.Parse(cmdline)
	if err != nil {
		return "", fmt.Errorf("can't parse command: %s", err)
	}

	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()

	return string(output), err
}
