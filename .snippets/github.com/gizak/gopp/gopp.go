package main

import (
	"fmt"
	"github.com/gizak/gopp/cmd"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "gopp: %s", err.Error())
		os.Exit(1)
	}
}
