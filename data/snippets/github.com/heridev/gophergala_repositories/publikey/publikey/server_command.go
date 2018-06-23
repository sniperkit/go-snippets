package publikey

import (
	"github.com/codegangsta/cli"
	"github.com/gerred/publikey/server"
)

func NewServerCommand() cli.Command {
	return cli.Command{
		Name:      "server",
		ShortName: "s",
		Usage:     "start publikey server for hosting public keys",
		Action:    handleServerCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "port, p",
				Value: "7890",
				Usage: "port to run the publikey server on",
			},
			cli.StringFlag{
				Name:  "data-file",
				Value: "publikey.db",
				Usage: "location for the BoltDB database",
			},
		},
	}
}

func handleServerCommand(c *cli.Context) {
	server.Serve(c.String("port"), c.String("data-file"))
}
