package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/MediaMath/cove"
	"github.com/pkg/browser"
)

func main() {
	outputDir := flag.String("o", "", "If provided output files will be written to this dir and not opened.")
	mode := flag.String("mode", "set", "Which cover mode to use.")
	short := flag.Bool("short", false, "Run only short tests.")
	keepProfile := flag.Bool("keep", false, "Will not remove coverage profile files if set.")
	flag.Parse()

	openReport := false
	reportPath := *outputDir
	if reportPath == "" {
		openReport = true
		reportPath, _ = ioutil.TempDir("", "cvr")
	}

	paths := flag.Args()

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	packs, pathErr := cove.Packages(cove.PackagePatternsFromStrings(paths)...)
	logError(pathErr)

	anyCoverage := false
	for _, pack := range packs {
		profile, coverErr := cove.CoverageProfile(*short, *mode, reportPath, pack)
		logError(coverErr)

		if profile != "" {
			if !*keepProfile {
				defer os.RemoveAll(profile)
			}

			report, reportErr := cove.CoverageReport(profile, reportPath)
			logError(reportErr)

			if report != "" {
				anyCoverage = true
				if openReport {
					browser.OpenFile(report)
				}
			}
		}
	}

	if !anyCoverage {
		fmt.Printf("No coverage for %s\n", strings.Join(paths, ","))
	}
}

func logError(err error) {
	if err != nil {
		log.Printf("%s\n", err)
	}
}
