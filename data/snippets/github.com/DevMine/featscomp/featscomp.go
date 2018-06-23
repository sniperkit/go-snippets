// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package featscomp computes features.
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/DevMine/featscomp/config"
	"github.com/DevMine/featscomp/features"
	"github.com/DevMine/featscomp/util"
)

const (
	averageForksName       = "average_forks"
	averageStarsName       = "average_stars"
	commitsCountName       = "commits_count"
	contributionsCountName = "contributions_count"
	followersCountName     = "followers_count"
	hireableName           = "hireable"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [CONFIGURATION_FILE]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	configPath := flag.Arg(0)
	if len(configPath) == 0 {
		flag.Usage()
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		fatal(err)
	}

	db, err := util.OpenDBSession(cfg.Database)
	if err != nil {
		fatal(err)
	}
	defer db.Close()

	var feats []features.Feature

	cf := cfg.Features

	// insert features in the database
	// TODO refactor this since all features implement the same interface.
	if cf.AverageForks {
		averageForks, err := features.NewAverageForks(averageForksName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, averageForks)
	}

	if cf.AverageStars {
		averageStars, err := features.NewAverageStars(averageStarsName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, averageStars)
	}

	if cf.CommitsCount {
		commitsCount, err := features.NewCommitsCount(commitsCountName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, commitsCount)
	}

	if cf.ContributionsCount {
		contributionsCount, err := features.NewContributionsCount(contributionsCountName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, contributionsCount)
	}

	if cf.FollowersCount {
		followersCount, err := features.NewFollowersCount(followersCountName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, followersCount)
	}

	if cf.Hireable {
		hireable, err := features.NewHireable(hireableName, db)
		if err != nil {
			fatal(err)
		}
		feats = append(feats, hireable)
	}

	run(feats)
}

func run(feats []features.Feature) {
	var wg sync.WaitGroup
	wg.Add(len(feats))

	tic := time.Now()
	for _, f := range feats {
		fmt.Printf("computing feature %s\n", f.Name())
		go func(f features.Feature) {
			defer wg.Done()
			if err := f.Score(); err != nil {
				util.FPrintlnErr(err)
			}
		}(f)
	}
	wg.Wait()

	toc := time.Now()
	fmt.Println("total elapsed time: ", toc.Sub(tic))
}

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
