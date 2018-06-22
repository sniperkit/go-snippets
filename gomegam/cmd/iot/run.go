package main

import (
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gnuflag"
)

type IOTStart struct {
	fs *gnuflag.FlagSet
}

func (g *IOTStart) Info() *cmd.Info {
	desc := `starts the iot base web daemon.
	
	
	`
	return &cmd.Info{
		Name:    "start",
		Usage:   `start`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *IOTStart) Run(context *cmd.Context, client *cmd.Client) error {
	RunWeb()
	return nil
}

func (c *IOTStart) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("iot", gnuflag.ExitOnError)
	}
	return c.fs
}

