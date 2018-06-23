package main

import (
	"log"
	"log/syslog"
)

func main() {
	exampleLogger()
}

func exampleSimple() {

	sl, err := syslog.New(syslog.LOG_LOCAL5|syslog.LOG_NOTICE, "simple")
	if err != nil {
		log.Fatal(err)
	}
	defer sl.Close()

	sl.Notice("Hello, Syslog!")
}

func exampleLogger() {

	sl, err := syslog.New(syslog.LOG_LOCAL5|syslog.LOG_NOTICE, "logger")
	if err != nil {
		log.Fatal(err)
	}
	defer sl.Close()

	log.SetOutput(sl)

	log.Print("Hello, Syslog! Example Logger.")
}
