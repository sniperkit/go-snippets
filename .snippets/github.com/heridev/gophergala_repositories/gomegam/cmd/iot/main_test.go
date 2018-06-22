package main

import (
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestCommandsFromBaseManagerAreRegistered(c *gocheck.C) {
	baseManager := cmd.BuildBaseManager("megam", version, header)
	manager := buildManager("megam")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, gocheck.Equals, true)
		c.Assert(command, gocheck.FitsTypeOf, instance)
	}
}