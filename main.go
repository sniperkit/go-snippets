// dt is a command line tool for storing data together on the distributed web
package main

import (
	"os"

	"github.com/datatogether/dt/cmd"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}

	mode := os.Getenv("GOLANG_ENV")
	if mode != "PRODUCTION" {
		log.Out = os.Stdout
	} else {
		log.Out = os.Stderr
	}

	cmd.SetLogger(log)
}

func main() {
	// Catch errors & pretty-print.
	// comment this out to get stack traces back.
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				cmd.PrintErr(err)
			} else {
				log.Info(r)
			}
		}
	}()

	cmd.Execute()
}
