package main

import (
	"path/filepath"
	"fmt"
)

func main() {
	file := "20170121_004649_nodes.go"
	pattern := "2017*.go"
	match, err := filepath.Match(pattern, file)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}

	if match {
		fmt.Println("ok")
	} else  {
		fmt.Println("not found")
	}
}
