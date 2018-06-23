// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repotool is able to fetch information from a source code repository.
// Typically, it can get all commits, their authors and commiters and so on
// and return this information in a JSON object.
// Currently, on the Git VCS is supported.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/DevMine/repotool/config"
	"github.com/DevMine/srcanlzr/src"

	"github.com/DevMine/repotool/repo"
)

const version = "1.0.0"

// program flags
var (
	versionflag       = flag.Bool("version", false, "print version.")
	srctoolflag       = flag.String("srctool", "", "read json file produced by srctool (give stdin to read from stdin)")
	cpuprofileflag    = flag.String("cpuprofile", "", "write cpu profile to file")
	tmpDirflag        = flag.String("tmpdir", "", "temporary directory location")
	fileSizeLimitflag = flag.Float64("filesizelimit", 0.1, "maximum size, in GB, for a file to be processed in the temporary directory location")
	deltasflag        = flag.Bool("deltas", false, "fetch commit deltas")
	patchesflag       = flag.Bool("patches", false, "fetch commit patches")
)

func main() {
	var err error

	flag.Usage = func() {
		fmt.Printf("usage: %s [OPTION(S)] [REPOSITORY PATH]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	if *versionflag {
		fmt.Printf("%s - %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	if *cpuprofileflag != "" {
		f, err := os.Create(*cpuprofileflag)
		if err != nil {
			fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "invalid # of arguments")
		flag.Usage()
	}

	cfg := new(config.Config)
	cfg.Data.TmpDir = *tmpDirflag
	cfg.Data.TmpDirFileSizeLimit = *fileSizeLimitflag
	cfg.Data.CommitDeltas = *deltasflag
	cfg.Data.CommitPatches = *patchesflag

	repoPath := flag.Arg(0)
	var repository repo.Repo
	repository, err = repo.New(cfg.Data, repoPath)
	if err != nil {
		fatal(err)
	}
	defer func() {
		repository.Cleanup()
		if err != nil {
			fatal(err)
		}
	}()

	fmt.Fprintln(os.Stderr, "fetching repository commits...")
	tic := time.Now()
	err = repository.FetchCommits()
	if err != nil {
		return
	}
	toc := time.Now()
	fmt.Fprintln(os.Stderr, "done in ", toc.Sub(tic))

	if *srctoolflag == "" {
		var bs []byte
		bs, err = json.Marshal(repository)
		if err != nil {
			return
		}
		fmt.Println(string(bs))
	} else {
		var r *bufio.Reader
		if *srctoolflag == strings.ToLower("stdin") {
			// read from stdin
			r = bufio.NewReader(os.Stdin)
		} else {
			// read from srctool json file
			var f *os.File
			if f, err = os.Open(*srctoolflag); err != nil {
				return
			}
			r = bufio.NewReader(f)
		}

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, r); err != nil {
			return
		}

		var p *src.Project
		p, err = src.Decode(buf)
		if err != nil {
			return
		}

		p.Repo = repository.GetRepository()
		if err := p.Encode(os.Stdout); err != nil {
			return
		}
	}
}

// fatal prints an error on standard error stream and exits.
func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
