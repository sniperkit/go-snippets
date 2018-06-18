package main

import "fmt"

func main() {
	name := "Maik Ellerbrock"
	char := ""

	for x := 0; x < len(name); x++ {
		char = string(name[x])
		fmt.Println("char: ", char, " |  byte: ", []byte(char))
	}
}
