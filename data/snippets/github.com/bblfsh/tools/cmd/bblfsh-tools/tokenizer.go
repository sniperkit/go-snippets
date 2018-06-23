package main

import "github.com/bblfsh/tools"

type Tokenizer struct {
	Common
}

func (c *Tokenizer) Execute(args []string) error {
	return c.execute(args, tools.Tokenizer{})
}
