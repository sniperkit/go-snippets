package main

// change versionString by cmd:
// go run -ldflags "-X main.versionString=1.0" ..

import "fmt"

var versionString = "unset"

func main() {
	fmt.Println("Version: ", versionString)
}
