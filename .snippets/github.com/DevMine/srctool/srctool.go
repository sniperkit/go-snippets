// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// srctool is a command line tool to manage source code parsers. It is able to
// download parsers from a web server, install them and run them. In short, it is a
// manager for source code parsers.
package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/cmd"
	"github.com/DevMine/srctool/log"
)

func main() {
	app := cli.NewApp()
	app.Name = "srctool"
	app.Usage = "tool for parsing source code"
	app.Version = "1.0.0"
	app.Author = "The DevMine authors"
	app.Email = "contact@devmine.ch"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "d",
			Usage: "enable debug mode",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "install one or all language parser(s)",
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.Install(c)
			},
		},
		{
			Name:      "delete",
			ShortName: "d",
			Usage:     "delete one or all language parser(s)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dry",
					Usage: "dry mode",
				},
			},
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.Delete(c)
			},
		},
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "update one or all language parser(s)",
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.Update(c)
			},
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list installed or available parsers",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "r",
					Usage: "list remote parsers",
				},
			},
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.List(c)
			},
		},
		{
			Name:      "parse",
			ShortName: "p",
			Usage:     "parse a project",
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.Parse(c)
			},
		},
		{
			Name:      "config",
			ShortName: "c",
			Usage:     "create config file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "server-url",
					Usage: "get/set the download server URL",
				},
			},
			Action: func(c *cli.Context) {
				log.SetDebugMode(c.GlobalBool("d"))
				cmd.Config(c)
			},
		},
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "print program version",
			Action: func(c *cli.Context) {
				fmt.Printf("%s - v%s\n", app.Name, app.Version)
			},
		},
	}

	app.Run(os.Args)
}
