package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/example/tasks"
)

var (
	app        *cli.App
	configPath string
)

func init()  {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "machinery"
	app.Usage = "machinery worker and send example tasks with machinery send"
	app.Author = "Richard Knop"
	app.Email = "risoknop@gmail.com"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "",
			Destination: &configPath,
			Usage:       "Path to a configuration file",
		},
	}
}

func main() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "sms-worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				return smsWorker()
			},
		},
		{
			Name:  "send-sms",
			Usage: "send example tasks ",
			Action: func(c *cli.Context) error {
				return sendSms()
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func startServer() (server *machinery.Server, err error) {
	// Create server instance
	server, err = machinery.NewServer(loadConfig())
	if err != nil {
		return
	}

	// Register tasks
	tasks := map[string]interface{}{
		"add":               exampletasks.Add,
		"multiply":          exampletasks.Multiply,
		"panic_task":        exampletasks.PanicTask,
		"single_sms": SendSingleSms,
		"multiple_sms": SendMultipleSms,
		"long_running_task": LongRunningTask,
		"delayed_task": DelayedTask,
	}

	err = server.RegisterTasks(tasks)
	return
}
