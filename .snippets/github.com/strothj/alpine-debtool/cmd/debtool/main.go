package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/strothj/alpine-debtool/commands"
)

func main() {
	app := cli.NewApp()
	commands.AddCommands(app)
	app.Run(os.Args)
}
