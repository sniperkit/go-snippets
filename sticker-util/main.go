package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/hiro4bbh/sticker/plugin"
)

func init() {
	gob.Register([]interface{}(nil))
}

func main() {
	plugin.InitializePlugin()
	opts := NewOptions(os.Args[0], os.Stdout, os.Stderr)
	if err := opts.Parse(os.Args[1:]); err != nil {
		fmt.Printf("sticker-util: %s\n", err)
		os.Exit(1)
	}
	if err := opts.Run(); err != nil {
		fmt.Printf("sticker-util: %s\n", err)
		os.Exit(1)
	}
}
