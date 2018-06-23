package main

import "github.com/bblfsh/tools"

type NPath struct {
	Common
}

func (c *NPath) Execute(args []string) error {
	return c.execute(args, tools.NPath{})
}
