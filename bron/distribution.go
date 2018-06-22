package main

import (

	"fmt"
	"io/ioutil"

)

func filterDistribution(templates map[string]*Template, contentFile string, verbosePtr bool, quietPtr bool, statusPtr bool) map[string]int {

	counts := map[string]int{}

	language, t := determineLanguage(templates, contentFile, verbosePtr, quietPtr)

	if (!quietPtr && verbosePtr) {
		if language != "unknown" {
			fmt.Println("language determined as", language)

			// XXX multi-line is not accurate
			// XXX counting lines that have comments but also code, as code only
			commentMarkers := t.Comments
			fmt.Printf("whole thing, commentMarkers, %v", commentMarkers)
		} else {
			fmt.Println("skipping loc filtering for", contentFile, " due to not being able to determine language.")
		}
	}

// determine language
// strip comments & count lines
// strip white space and count
// count remainder

	//coder, _ := regexp.Compile("^[^\\s]

	file, err := ioutil.ReadFile(contentFile)
	check(err)

	if (!quietPtr && verbosePtr) {
		fmt.Println("found", len(file), "characters in file", contentFile)
	}

	counts["loc"] = len(file)

	return counts

}
