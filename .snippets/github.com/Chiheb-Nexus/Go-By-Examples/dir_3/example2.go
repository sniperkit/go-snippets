package main

import (
	"fmt"
)

func main() {
	tab := make([]int, 10) // Create a slice of int

	for i := range tab {

		if i%2 == 0 {
			continue // Skip iteration
		}

		tab[i] = i << uint(i)

		if tab[i] > 2000 {
			break // Break iteration
		}
	}

	for _, value := range tab {
		fmt.Printf("%d\n", value)
	}
}
