package main

import (
	"os"

	"github.com/Ken2mer/knmr/logger"
	"github.com/urfave/cli"
)

var (
	name      = "knmr"
	version   = "0.0.0"
	gitcommit string
)

func main() {
	logger.Infof("%s version: %s-%s\n", name, version, gitcommit)

	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = "knmr"

	app.Commands = commands
	app.Flags = flags
	app.Before = before

	if err := app.Run(os.Args); err != nil {
		logger.Errorf("error: %v\n", err)
		os.Exit(1)
	}
}
