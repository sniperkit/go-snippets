package main

import (
	"log"
	"os"
)

var stdOut, stdErr = log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0)

func main() {
	stdOut.Println("This is good output")
	stdErr.Fatal("this will log to stdErr and exit with status code 1")
}
