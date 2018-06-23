// Copyright 2014 Gabriel Guzman <gabe@lifewaza.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Heavily based on the gofmt command:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala/gomua"
	"io"
	"net/mail"
	"os"
	"path/filepath"
)

var exitCode = 0

// Header Parser (hp) takes an email message as input and returns
//  the message's headers.

// For now, just take input on standard in, assuming one message
//  at a time.
// Later, allow user to specify a filename, or directory as args.

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	msg, err := mail.ReadMessage(in)
	if err != nil {
		return err
	}

	//subject := msg.Header.Get("Subject")
	for k, v := range msg.Header {
		fmt.Printf("%v: %v\n", k, v)
	}
	return err
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil {
		err = processFile(path, nil, false)
	}
	if err != nil {
		fmt.Printf("%s", err)
		exitCode = 2
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		if err := processFile("<standard input>", os.Stdin, true); err != nil {
			fmt.Printf("%v", err)
			exitCode = 2
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		gomua.Scan(path)
	}
	os.Exit(exitCode)
}
