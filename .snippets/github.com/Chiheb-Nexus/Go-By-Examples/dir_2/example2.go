package main

import (
	"fmt"
)

func main() {

	var a [2][3]string

	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			a[i][j] = fmt.Sprintf("| Row: %d ; Column: %d |", i, j)
		}
	}

	for i := range a {
		fmt.Printf("%d- %s\n", i, a[i])
	}

	fmt.Printf("%q\n", a)
}
