package main

import "github.com/bblfsh/tools"

type CyclomaticComp struct {
	Common
}

func (c *CyclomaticComp) Execute(args []string) error {
	return c.execute(args, tools.CyclomaticComplexity{})
}
