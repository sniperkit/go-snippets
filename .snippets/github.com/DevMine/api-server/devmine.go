// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package apiserver starts the DevMine projects API server.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"github.com/DevMine/api-server/cache"
	"github.com/DevMine/api-server/config"
	"github.com/DevMine/api-server/srv"
)

func fatal(a ...interface{}) {
	glog.Error(a)
	os.Exit(1)
}

func main() {
	configPath := flag.String("c", "", "configuration file")
	flag.Parse()

	// Make sure we finish writing logs before exiting.
	defer glog.Flush()

	glog.Info("starting the API server...")

	if len(*configPath) == 0 {
		fatal("no configuration specified")
	}

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		fatal(err)
	}

	db, err := srv.OpenDBSession(cfg.Database)
	if err != nil {
		fatal(err)
	}
	defer db.Close()

	glog.Info("caching data...")
	tic := time.Now()
	err = cache.LoadCache(db)
	if err != nil {
		fatal(err)
	}
	toc := time.Now()
	glog.Info("done in ", toc.Sub(tic))

	router := srv.SetupRouter(db, cfg.Server.EnableCors)
	addr := fmt.Sprintf("%s:%d", cfg.Server.HostName, cfg.Server.Port)
	glog.Infof("listening on %s...\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		fatal(err)
	}
}
