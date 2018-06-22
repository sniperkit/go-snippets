package main

import (
	"os"
	"path"
	"sort"

	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

var version string

var app = &cli.App{
	Name:    "swm",
	Version: version,
	Usage:   "swm <command>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "profile", Usage: "The profile for the TMUX session", Value: getDefaultProfile()},
		&cli.StringFlag{Name: "story", Usage: "The story for the TMUX session", Value: getDefaultStory()},
		&cli.StringFlag{Name: "code-path", Usage: "The absolute path to the code path", Value: path.Join(os.Getenv("HOME"), "code")},
		&cli.StringFlag{Name: "ignore-pattern", Usage: "The Regex pattern to ignore", Value: "^.snapshots$"},
		&cli.BoolFlag{Name: "debug", Usage: "enable debug mode"},
	},
	Before: createLogger,
	Authors: []*cli.Author{
		{
			Name:  "Wael Nasreddine",
			Email: "wael.nasreddine@gmail.com",
		},
	},
	Commands: []*cli.Command{
		// coder for code management
		coderCmd,
		// tmux for switch client
		tmuxCmd,
	},
}

func init() {
	// sort the commands/flags
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	sortCommands(app.Commands)
}

func main() {
	// run the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("error occurred")
	}
}
