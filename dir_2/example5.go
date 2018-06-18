package main

import (
	"fmt"
)

func main() {

	names := []string{
		"Chiheb",
		"Nexus",
		"Alex",
	}
	
	fmt.Printf("Length: %d\n", len(names))

	cities := make([]string, 42)
	fmt.Printf("Length: %d\n", len(cities))

	countries := []string{}
	fmt.Printf("Length: %d\n", len(countries))
}
