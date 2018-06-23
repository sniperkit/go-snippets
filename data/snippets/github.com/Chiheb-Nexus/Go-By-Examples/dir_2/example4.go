package main

import (
	"fmt"
)

func main() {
	// First manner
	/*
		names := make([]string, 3)
		names[0] = "Nexus"
		names[1] = "Chiheb"
		names[2] = "Alex"
	*/

	// Second manner
	/*
		names := []string{}
		names = append(names, "Nexus")
		names = append(names, "Chiheb")
		names = append(names, "Alex")
	*/

	// Third manner
	/*
		names := []string{}
		names = append(names, "Nexus", "Chiheb", "Alex")
	*/

	// Fourth manner using ellipsis
	// Note: slice's types must match. If not go will throw an error !
	names := []string{"Chiheb", "Nexus"}
	names_2 := []string{"Alex", "Paul"}
	names = append(names, names_2...)

	fmt.Printf("%q\n", names)
}
