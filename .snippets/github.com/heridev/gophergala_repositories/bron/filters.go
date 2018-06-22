package main

import (

	"fmt"

)

func parse(files []string, templates map[string]*Template, verbosePtr bool, quietPtr bool, statusPtr bool) {

	for _, file := range files {
		d := "going to parse..." + file
		if (!quietPtr && verbosePtr) {
			fmt.Println(colorize(d))
		}
		filterDistribution(templates, file, verbosePtr, quietPtr, statusPtr)
	}
	if (!quietPtr && verbosePtr) {
		fmt.Println("all done")
	}

}
