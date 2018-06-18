package cmd

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestCmdError(t *testing.T) {
	cmd, _, gErr := getFailingCmd()
	if gErr != nil {
		t.Errorf("%v", gErr)
	}

	err := Run(cmd)
	if err == nil {
		t.Errorf("Should have error")
	}

	if _, ok := err.(*Error); !ok {
		t.Errorf("Should be CmdError: %v", err)
	}
}

func TestCmdErrorDoesntTrapNils(t *testing.T) {
	err := newCmdError(nil, []string{"foo"})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorDoesntTrapNonExitErrors(t *testing.T) {
	in := fmt.Errorf("foo")
	err := newCmdError(in, []string{"foo"})
	if err != in {
		t.Errorf("%v", err)
	}
}

//to get an exit err I need to run a process
//and get that error
func getRealExitError() error {
	cmd, _, gErr := getFailingCmd()
	if gErr != nil {
		return gErr
	}
	return cmd.Run()
}

func getFailingCmd() (*exec.Cmd, string, error) {
	tempDir, err := ioutil.TempDir("", "failingcmd")
	if err != nil {
		return nil, "", err
	}

	os.RemoveAll(tempDir)

	//ls on a non-existent directory should fail
	//not portable to systems without ls
	return exec.Command("ls", tempDir), tempDir, nil
}

func TestCmdErrorTrapsExitErrors(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, []string{"foo"})
	if err == in {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorUsesStdErr(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, []string{"foo", "bar"})
	if err == in {
		t.Errorf("%v", err)
	}

	if err.Error() != "foo\nbar" {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorUsesInErrorInFaceOfNilStdErr(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, nil)

	if err.Error() != in.Error() {
		t.Errorf("%v", err)
	}
}
