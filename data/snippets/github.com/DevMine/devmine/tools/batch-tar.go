//usr/bin/env go run $0 $@; exit
// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package batch-tar.go is able to tar VCS repositories located in several
// folders.
// Note that tar is in-place, ie: the initial directory is removed and only
// the directory in tar format remains.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [(OPTIONS)] [REPOSITORIES ROOT FOLDER]\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}
	depthflag := flag.Uint("d", 0, "depth level where to find repositories")
	maxGoroutines := flag.Uint("g", uint(runtime.NumCPU()), "max number of goroutines to spawn")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "invalid # of arguments")
		flag.Usage()
	}

	reposDir := flag.Arg(0)

	tarBin, err := exec.LookPath("tar")
	if err != nil {
		fatal(err)
	}

	tasks := make(chan *exec.Cmd, 0)
	var wg sync.WaitGroup
	for w := uint(0); w < *maxGoroutines; w++ {
		wg.Add(1)
		go func() {
			for cmd := range tasks {
				out, err := cmd.CombinedOutput()
				fmt.Print(string(out))
				if err != nil {
					fmt.Println(err)
				}
			}
			wg.Done()
		}()
	}

	iterateRepos(tasks, tarBin, reposDir, *depthflag)

	close(tasks)
	wg.Wait()
}

func iterateRepos(tasks chan *exec.Cmd, tarBin, dirPath string, depth uint) {
	fis, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fatal(err)
	}

	if depth == 0 {
		for _, fi := range fis {
			if !fi.IsDir() {
				continue
			}

			fmt.Println("adding repository: ", fi.Name(), " to the tasks pool")
			repoPath := filepath.Join(dirPath, fi.Name())

			tasks <- exec.Command(tarBin, "--remove-files", "-cf", repoPath+".tar",
				"-C", path.Dir(repoPath), "--", path.Base(repoPath))
		}
		return
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		iterateRepos(tasks, tarBin, filepath.Join(dirPath, fi.Name()), depth-1)
	}
}

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
