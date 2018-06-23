package commands

import "github.com/codegangsta/cli"

// AddCommands adds the available commands to the command line interface.
func AddCommands(app *cli.App) {
	RegisterDownload(app)
}
