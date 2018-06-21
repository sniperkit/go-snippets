package main

import (
	"log"
	"regexp"
)

func main() {
	//<span id="gamesPage" style="display:none;"> 11 </span>
	reg := regexp.MustCompile(`<span\s+id="gamesPage".*?>\s*(\d+)\s*</span>`)
	q := reg.FindStringSubmatch(`<span id="gamesPage" style="display:none;"> 11 </span>`)
	if q != nil {
		log.Println(q[1])
	}
}
