package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/gerred/publikey/publikey"
)

func main() {
	app := cli.NewApp()
	app.Name = "publikey"
	app.Author = Author
	app.Email = Email
	app.Version = Version
	app.Usage = "interact with the publikey API for public SSH key management"
	app.Commands = []cli.Command{
		publikey.NewListCommand(),
		publikey.NewServerCommand(),
		publikey.NewRegisterCommand(),
		publikey.NewAddCommand(),
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host, H",
			Value: "api.publikey.io",
			Usage: "publikey host url (including port)",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "publikey user, if not logged in",
		},
	}

	app.Run(os.Args)
}
