package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	parser = flags.NewNamedParser("srclib-docker", flags.Default)
	cwd    = getCWD()
)

// GlobalOpt contains global options.
var GlobalOpt struct {
	Verbose bool `short:"v" description:"show verbose output"`
}

func init() {
	parser.LongDescription = "srclib-docker is a srclib toolchain that scans and analyzes Dockerfiles in a repository or tree."
	parser.AddGroup("Global options", "", &GlobalOpt)
}

func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

func main() {
	log.SetFlags(0)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
