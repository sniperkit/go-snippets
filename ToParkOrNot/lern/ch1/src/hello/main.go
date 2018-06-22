package main

import (
	"fmt"
)

func main() {
	message := "The quick brown fox jumped over the lazy dog.\n"

	for i, c := range message {
		fmt.Printf("%d %c\n", i, c)
	}

	fmt.Printf(message)
}
