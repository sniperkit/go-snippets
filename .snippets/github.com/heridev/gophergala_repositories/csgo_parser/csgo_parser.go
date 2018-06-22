package main

import (
	"flag"
	"fmt"
)

func main() {
	filename := flag.String("f", "", "Name of the demo file")

	flag.Parse()

	fmt.Println(*filename)

	// use DemoFileDump to open
	var demo DemoFileDump
	if demo.Open(*filename) == true {
		fmt.Println("Demo file is good")

		demo.DoDump()
	}
}
