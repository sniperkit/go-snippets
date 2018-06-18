package cmd

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func ExamplePipeWith() {
	tempdir, _ := ioutil.TempDir("", "example-pipewith")
	defer os.RemoveAll(tempdir)

	os.MkdirAll(filepath.Join(tempdir, "sub1"), 0744)
	os.MkdirAll(filepath.Join(tempdir, "sub2"), 0744)

	PipeWith(exec.Command("ls", tempdir), func(stdout io.Reader) error {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		return nil
	})

	//Output:
	//sub1
	//sub2
}

func ExampleOutput() {
	tempdir, _ := ioutil.TempDir("", "example-output")
	defer os.RemoveAll(tempdir)

	os.MkdirAll(filepath.Join(tempdir, "sub1"), 0744)
	os.MkdirAll(filepath.Join(tempdir, "sub2"), 0744)

	out, _ := Output(exec.Command("ls", tempdir))

	for _, k := range out {
		fmt.Println(k)
	}

	//Output:
	//sub1
	//sub2
}

func ExampleRun() {
	tempdir, _ := ioutil.TempDir("", "example-run")
	defer os.RemoveAll(tempdir)

	os.MkdirAll(filepath.Join(tempdir, "sub1"), 0744)
	os.MkdirAll(filepath.Join(tempdir, "sub2"), 0744)

	if runErr := Run(exec.Command("ls", tempdir)); runErr != nil {
		fmt.Printf("Got error %v", runErr)
	}

	//Output:
}
