//usr/bin/env go run $0 $@; exit
// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package genmd5 generates MD5 sums from parsers archives, suitable to
// generate the MD5SUM file for the parsers repository.
// This tool is useful because the output of standard tools to generate MD5
// sums may vary especially md5sum(1) from Linux and md5i(1) from the BSDs.
package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	extflag := flag.String("x", ".zip", "file extension to restrict to")
	flag.Usage = func() {
		fmt.Printf("usage: %s [(DIRECTORY)]\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	var path string
	args := len(flag.Args())
	if args < 1 {
		path = "./"
	} else if args == 1 {
		path = flag.Arg(0)
	} else {
		flag.Usage()
	}

	if err := genMD5Sums(path, *extflag); err != nil {
		log.Fatal(err)
	}
}

func genMD5Sums(dirPath, extension string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// all parsers are zip archives
		if filepath.Ext(info.Name()) != extension {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		sum := md5.Sum(file)
		fmt.Printf("%s %s\n", hex.EncodeToString(sum[:]), path)

		return nil
	})
}
