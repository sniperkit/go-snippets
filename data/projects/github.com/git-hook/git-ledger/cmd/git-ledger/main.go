package main

import (
	"os"

	"github.com/Masterminds/semver"
	"github.com/urfave/cli"
)

var VERSION = semver.MustParse("0.1.1")

func main() {

	app := cli.NewApp()

	app.Name = "git-ledger"
	app.Usage = "Index your git clones."
	app.Version = VERSION.String()

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Usage:     "Start tracking an existing repository.",
			ArgsUsage: "[path]",
			Action:    add,
		},
		{
			Name:      "find",
			Usage:     "Print the location of a tracked repository.",
			ArgsUsage: "[path]",
			Action:    find,
		},
		{
			Name:   "ls",
			Usage:  "Print all tracked repositories",
			Action: ls,
		},
		{
			Name:      "rm",
			Usage:     "Stop tracking an existing repository.",
			ArgsUsage: "[path]",
			Action:    rm,
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
