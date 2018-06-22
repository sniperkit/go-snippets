// cli.go defines the toolchain's command line interface.
package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	// TODO: Replace 'quickstart' with the name of the language this toolchain is for.
	flagParser = flags.NewNamedParser("srclib-quickstart", flags.Default)
	cwd        = getCWD()
)

// init is called before main.
func init() {
	// TODO: Replace this with a description of the toolchain.
	flagParser.LongDescription = "srclib-quickstart will guide you in building a srclib backend."
}

func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// main is the entry point of the executable.
func main() {
	log.SetFlags(0)
	if _, err := flagParser.Parse(); err != nil {
		os.Exit(1)
	}
}
