// standalone "updater" of the xmlstats 
package main

import (
	"fmt"
    "github.com/gophergala/NextMatch/updater/xmlstats"
	"log"
	"os"
)

//TODO: add parameter for game and date
//TODO: convert BySport in a method instead of a function and have xmlstats.Token as instance attribute
func main() {

	xmlstats.Token = os.Getenv("XMLSTATS_TOKEN")
	if len(xmlstats.Token) == 0 {
		log.Fatal("Specify XMLSTATS_TOKEN environment variable")
	}
	events, err := xmlstats.BySport("nba", "20150124")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", events)
}

