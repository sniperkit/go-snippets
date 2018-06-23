package main

import (
	"github.com/Ken2mer/knmr/logger"
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	agentCommand,
	githubCommand,
	godocCommand,
	linebotCommand,
	twitterCommand,
}

var flags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug, d",
		Usage: "enable debug output in logs",
	},
}

var before = func(context *cli.Context) error {
	debugEnabled := context.GlobalBool("debug")
	if debugEnabled {
		logger.Verbose = true
	}
	return nil
}
