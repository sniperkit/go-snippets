package main

import (
	"flag"
	"fmt"

	"github.com/EricYT/go-examples/flags/bar"
)

var foo string

func init() {
	flag.StringVar(&foo, "foo", "default foo", "foo option")
}

func main() {
	// Parse
	flag.Parse()
	// Generate a random value between 0 and max
	fmt.Println("main function:", foo)
	fmt.Println("bar is ", bar.Display())
}
