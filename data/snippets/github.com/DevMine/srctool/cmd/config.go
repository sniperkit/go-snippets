// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Config command provides options for creating a default config file, getting
// values and setting configuration values.
func Config(c *cli.Context) {
	// This will create the configuration directory and file if it does not
	// already exist.
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if c.Bool("server-url") {
		if len(c.Args()) > 1 {
			log.Fatal("invalid number of argument")
		}

		if !c.Args().Present() {
			fmt.Println("server-url = ", cfg.DownloadServerURL)
			return
		}

		cfg.DownloadServerURL = c.Args().First()
		if err = cfg.Save(); err != nil {
			log.Fatal(err)
		}

		log.Success("download server URL successfully updated")
		return
	}
}
