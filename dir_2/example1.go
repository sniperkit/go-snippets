package main

import (
	"fmt"
)

func main() {
	// First manner
	/*
		var a [2]string
		a[0] = "Hello"
		a[1] = "World!"
	*/

	// Second manner
	/*
		a := [2]string{0: "Hello", 1: "World!"}
	*/

	// Third manner
	/*
		a := [2]string{}
		a[0] = "Hello"
		a[1] = "World!"
	*/

	// Fourth manner with ellipsis
	a := [...]string{"Hello", "World!"}

	fmt.Printf("%s %s\n", a[0], a[1])
	// See array entries
	fmt.Printf("%q\n", a)
}
