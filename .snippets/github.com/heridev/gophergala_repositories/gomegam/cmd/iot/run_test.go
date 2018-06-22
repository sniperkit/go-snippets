package main

import (
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestIOTStartInfo(c *gocheck.C) {
	desc := `starts the iot base web daemon.
	
	
	`
	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := IOTStart{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}
