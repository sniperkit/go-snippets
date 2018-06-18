package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	fmt.Println("boom-go")

	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("command context:%+v\n", c)
		fmt.Println("boom! I say!")
		return nil
	}

	app.Run(os.Args)
}
