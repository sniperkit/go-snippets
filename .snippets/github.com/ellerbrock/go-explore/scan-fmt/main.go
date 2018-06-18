package main

import "fmt"

func main() {

	// Problem:
	// -------
	// Lets say the User enters: "Maik Ellerbrock"
	// name would only be "Maik"
	//
	// Workaround:
	// -----------
	// use of the "bufio" Library
	// https://golang.org/pkg/bufio/#NewScanner

	var name string

	fmt.Print("Please enter your Name: ")
	fmt.Scanln(&name)

	fmt.Println("Hello ", name)
}
