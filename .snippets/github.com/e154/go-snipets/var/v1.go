package main

import (
	"fmt"
)

func main() {
	a := 1

	if a > 0 {
		a := 2

		a += 1
	}

	fmt.Println( a ) // Что будет выведено на экран? И почему?
}