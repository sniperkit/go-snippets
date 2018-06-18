package main

import (
	"fmt"
	"os"

	"github.com/uoregon-libraries/gopkg/bagit"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <path to bag directory>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	var b = bagit.New(os.Args[1])
	var err = b.WriteTagFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating tag files for %q: %s\n", os.Args[1], err)
	}
}
