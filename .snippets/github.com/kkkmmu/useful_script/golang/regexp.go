package main

import (
	"fmt"
	"regexp"
)

var test = `window.synccheck={retcode:"0",selector:"2"}`
var re = regexp.MustCompile(`window\.synccheck=\{retcode\:\"(?P<retcode>[0-9]+)\"\,selector\:\"(?P<selector>[0-9]+)\"\}`)

func main() {
	matches := re.FindStringSubmatch(test)
	if len(matches) > 1 {
		fmt.Println(matches[1])
		fmt.Println(matches[2])
	}
}
