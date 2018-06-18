//Package cmd provides helper functions for working with the os/exec library.
package cmd

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

//Run executes the command & throws away stdout.
//Stderr and error exit codes are reported as the returned error.
func Run(cmd *exec.Cmd) error {
	return PipeWith(cmd, func(stdout io.Reader) error { return nil })
}

//PipeWith executes the command & sends the provided event handler the stdout io.Reader.
//Errors returned from the event handler are forwarded on.
//Stderr and error exit codes are reported as the returned error.
func PipeWith(cmd *exec.Cmd, action func(io.Reader) error) error {
	stdout, stdOutErr := cmd.StdoutPipe()
	if stdOutErr != nil {
		return stdOutErr
	}

	stderr, stderrErr := cmd.StderrPipe()
	if stderrErr != nil {
		return stderrErr
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := action(stdout); err != nil {
		return err
	}

	stderrLines, stderrLinesErr := scanReader(stderr)
	if stderrLinesErr != nil {
		return stderrLinesErr
	}

	err := newCmdError(cmd.Wait(), stderrLines)
	return err
}

//Output executes the command & returns stdout.
//Stderr and error exit codes are reported as the returned error.
func Output(cmd *exec.Cmd) ([]string, error) {
	var results []string
	err := PipeWith(cmd, func(stdout io.Reader) error {
		var scanerr error
		results, scanerr = scanReader(stdout)
		return scanerr
	})

	return results, err
}

func scanReader(reader io.Reader) ([]string, error) {
	var results []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	return results, scanner.Err()
}

//Error is returned when an exit error is encountered.
//It also captures the stderr lines from the command.
type Error struct {
	Exit   *exec.ExitError
	StdErr []string
}

func (s *Error) Error() string {
	stderr := strings.Join(s.StdErr, "\n")
	if stderr != "" {
		return stderr
	}

	return s.Exit.Error()
}

func newCmdError(err error, stderr []string) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return &Error{exitErr, stderr}
	}

	return err
}
