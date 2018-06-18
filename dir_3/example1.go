package main

import (
	"fmt"
)

func main() {
	var tab = []int{1, 2, 3, 4, 5, 6}

	for key, value := range tab {
		fmt.Printf("Key: %d | Value: %d\n", key, value)
	}

	cities := map[string]int{
		"New York":    8520,
		"Los Angelos": 3650,
		"Tunis":       4520,
		"Paris":       3560,
	}
	for key, value := range cities {
		fmt.Printf("%s has %d startups\n", key, value)
	}
}
