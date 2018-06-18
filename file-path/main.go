package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// example 1
	path1, _ := filepath.Abs(".")
	fmt.Println(path1)

	// example 2
	path2, _ := os.Getwd()
	fmt.Println(path2)
}
