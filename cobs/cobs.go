package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gophergala/cobs/backend"
	"github.com/gophergala/cobs/builder"
)

func main() {
	app := cli.NewApp()
	app.Name = "CoBS"
	app.Usage = "Container Build Service"
	app.Commands = []cli.Command{
		{
			Name: "backend",
			Action: func(c *cli.Context) {
				fmt.Println("run backend")
				backend.Run()
			},
		},
		{
			Name: "builder",
			Action: func(c *cli.Context) {
				fmt.Println("run builder")
				builder.Run()
			},
		},
		//{
		//	Name: "hunter",
		//	Action: func(c *cli.Context) {
		//		fmt.Println("run hunter")
		//		hunter.Run()
		//	},
		//},
		//		{
		//			Name: "instrumenter",
		//			Action: func(c *cli.Context) {
		//				fmt.Println("run instrumenter")
		//				instrumenter.Run()
		//			},
		//		},
		{
			Name: "all",
			Action: func(c *cli.Context) {
				fmt.Println("run all")
				go backend.Run()
				builder.Run()
			},
		},
	}
	app.Run(os.Args)
}
