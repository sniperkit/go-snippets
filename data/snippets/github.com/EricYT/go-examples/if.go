package main

import (
	"fmt"
)

func main() {
	var x string
	x = "hello world"

	if x := 1; x > 0 {
		fmt.Println("if local variable")
	}

	fmt.Println(x)

	//
	for pos, char := range "我们都是好同志" {
		fmt.Printf("character %c starts at byte position %d\n", char, pos)
	}
}
