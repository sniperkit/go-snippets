package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/scpm"
	// "log"
	"os"
	"runtime"
)

var (
	version     string
	globalFlags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 22,
		},
		cli.StringFlag{
			Name:  "identity, i",
			Usage: "ssh identity to use for connecting to the host",
		},
		cli.StringFlag{
			Name:  "in",
			Value: "",
			Usage: "/path/to/file or /path/to/folder",
		},
		cli.StringSliceFlag{
			Name:  "path",
			Value: &cli.StringSlice{},
			Usage: "user@example.com:/path/to",
		},
	}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "scpm"
	// app.EnableBashCompletion = true
	app.Author = "gronpipmaster"
	app.Email = "gronpipmaster@gmail.com"
	app.Version = version
	app.Usage = "Copy files over ssh protocol to multiple servers."
	app.Flags = globalFlags
	app.Action = action
	app.Run(os.Args)
}

func action(ctx *cli.Context) {
	hosts := []*scpm.Host{}
	for _, host := range ctx.GlobalStringSlice("path") {
		h, err := scpm.NewHost(host, ctx.GlobalString("identity"), ctx.GlobalInt("port"))
		if err != nil {
			fatalln(err)
		}
		hosts = append(hosts, h)
	}
	if len(ctx.GlobalString("in")) == 0 {
		fatalln("Field --in required.")
	}
	scp, err := scpm.New(
		hosts,
		ctx.GlobalString("in"),
	)
	if err != nil {
		fatalln(err)
	}
	//Create chanels for wait quit signal
	quit := make(chan bool)
	//Init and run
	go scp.Run(quit)
	for {
		select {
		case <-quit:
			os.Exit(0)
		}
	}
}

func fatalln(i ...interface{}) {
	fmt.Println(i...)
	os.Exit(1)
}
