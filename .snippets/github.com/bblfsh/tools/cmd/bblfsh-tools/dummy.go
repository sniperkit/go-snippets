package main

import "github.com/bblfsh/tools"

type Dummy struct {
	Common
}

func (c *Dummy) Execute(args []string) error {
	return c.execute(args, tools.Dummy{})
}
