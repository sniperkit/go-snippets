package main

import (
	"fmt"
	"os"
)

func main() {
	// a behavior test of using relative file path

	// test 1: `cd path-to/lab/go && go run gofile.go` → file exists
	// test 2: `cd path-to/lab && go run go/gofile.go` → file doesn't exists
	if _, err := os.Stat("./gofile.go"); err == nil {
		fmt.Println("file exist")
	} else {
		fmt.Println("file doesn't exist")
	}
}
