package main

import "fmt"
import "path/filepath"

func main() {
	var name string = "CenterOS7_2.ovf"
	var pattern string = "CenterOS7_2*.ovf"

	math, err := filepath.Match(pattern, name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if math {
		fmt.Println("math")
	}
}
