package main

import (
	"fmt"
	"os"
)

func main() {
	dir, err := os.Open(".")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range fileInfo {
		fmt.Println(f.Name())
	}
}
