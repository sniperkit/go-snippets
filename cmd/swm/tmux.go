package main

import (
	"github.com/kalbasit/swm/tmux"

	cli "gopkg.in/urfave/cli.v2"
)

var tmuxCmd = &cli.Command{
	Name: "tmux",
	Subcommands: []*cli.Command{
		// switch client switches tmux client
		{
			Name:   "switch-client",
			Usage:  "Switch the client within the session for this profile and story",
			Action: tmuxSwitchClient,
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "kill-pane", Usage: "kill the TMUX pane after switch client"},
			},
		},
		// vim exit will save/exit any open vim
		{
			Name:   "vim-exit",
			Usage:  "Close all of open Vim within the session for this profile and story",
			Action: tmuxVimExit,
		},
		// kill-server will kill the server
		{
			Name:   "kill-server",
			Usage:  "Kill the server closes the tmux session for this profile and story",
			Action: tmuxKillServer,
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "vim-exit", Usage: "if vim is found running, kill it"},
			},
		},
	},
}

func tmuxSwitchClient(ctx *cli.Context) error {
	tmuxManager, err := newTmuxManager(ctx)
	if err != nil {
		return err
	}

	return tmuxManager.SwitchClient(ctx.Bool("kill-pane"))
}

func tmuxVimExit(ctx *cli.Context) error {
	tmuxManager, err := newTmuxManager(ctx)
	if err != nil {
		return err
	}

	return tmuxManager.VimExit()
}

func tmuxKillServer(ctx *cli.Context) error {
	tmuxManager, err := newTmuxManager(ctx)
	if err != nil {
		return err
	}

	return tmuxManager.KillServer(ctx.Bool("vim-exit"))
}

func newTmuxManager(ctx *cli.Context) (tmux.Manager, error) {
	// create a new coder
	c, err := newCoder(ctx)
	if err != nil {
		return nil, err
	}
	// scan the code path
	if err := c.Scan(); err != nil {
		return nil, err
	}
	// create a new TMUX manager
	tmuxManager := tmux.New(&tmux.Options{
		Coder:   c,
		Profile: ctx.String("profile"),
		Story:   ctx.String("story"),
	})

	return tmuxManager, nil
}
