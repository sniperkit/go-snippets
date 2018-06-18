package main

import (
	"bufio"
	"fmt"
	"os"
)

// use of "bufio.NewScanner" solves the problem with whitespaces via "fmt.Scanln"
func main() {
	// example user input: "Maik Ellerbrock" works well with bufio
	// with fmt.Scanln it would only return "Maik"
	fullName := bscan("first and last name")
	city := bscan("city")

	fmt.Printf("Hello %s from %s!", fullName, city)
}

func bscan(msg string) string {
	var res string

	fmt.Printf("Please enter your %s: ", msg)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		res = scanner.Text()
		break
	}
	return res
}
