package main

import (
	"fmt"
)

func main() {

	s := "Go编程"
	fmt.Println(len(s))
	fmt.Println(len([]rune(s)))
}
