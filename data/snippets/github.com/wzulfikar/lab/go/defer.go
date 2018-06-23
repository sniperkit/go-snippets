package main

import "fmt"

func main() {
	defer fmt.Println("deferred code")

	if true {
		fmt.Println("if block")
		return
	}
	fmt.Println("end")
}
