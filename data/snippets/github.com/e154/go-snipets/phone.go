package main

import (
	"regexp"
	"log"
	"fmt"
)

func main() {

	example := "+7(913) 318-95-46"

	// Make a Regex to say we only want
	reg, err := regexp.Compile("[a-zA-Z-+/(/) ]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(example, "")

	fmt.Printf("A string of %s becomes %s \n", example, processedString)

}
