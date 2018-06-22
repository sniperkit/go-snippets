// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/DevMine/srcanlzr/anlzr"
	"github.com/DevMine/srcanlzr/src"
)

const (
	version = "0.0.0"
)

func formatOutput(r *anlzr.Result) ([]byte, error) {
	var bs []byte
	var err error

	switch *format {
	case "JSON":
		bs, err = json.Marshal(r)
	case "XML":
		bs, err = xml.Marshal(r)
	case "protobuf":
		err = errors.New("protobuf is not yet implemented")
	default:
		err = errors.New("unsupported output format")
	}

	return bs, err
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// program flags
var (
	format         = flag.String("f", "JSON", "Output format. Possible values are: JSON, XML, protobuf")
	outputFileName = flag.String("o", "", "Output file name. By default, the output is set to stdout")
	cpuprofile     = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile     = flag.String("memprofile", "", "write memory profile to this file")
	vflag          = flag.Bool("v", false, "Print version.")
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [JSON PATH]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	if *vflag {
		fmt.Printf("%s - v%s\n", filepath.Base(os.Args[0]), version)
		return
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	var reader io.Reader

	if l := len(flag.Args()); l > 1 {
		// more that one file are passed as arguments
		fmt.Fprint(os.Stderr, "too many arguments\n\n")
		flag.Usage()
	} else if l == 0 {
		// no argument, we read from stdin
		reader = os.Stdin
	} else {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		reader = f
	}

	out := os.Stdout
	if len(*outputFileName) > 0 {
		var err error
		out, err = os.Open(*outputFileName)
		if err != nil {
			fatal(err)
		}
		defer out.Close()
	}

	p, err := src.Decode(reader)
	if err != nil {
		fatal(err)
	}

	res, err := anlzr.RunAnalyzers(p, anlzr.LoC{}, anlzr.Complexity{}, anlzr.LocPerLang{}, anlzr.CommentRatios{})
	if err != nil {
		fatal(err)
	}

	bs, err := formatOutput(res)
	if err != nil {
		fatal(err)
	}

	fmt.Fprintln(out, string(bs))
}
