package main

import "fmt"

func main() {
	// The defer is LIFO order.
	for i := 0; i < 10; i++ {
		defer fmt.Printf("%d ", i)
	}
}
