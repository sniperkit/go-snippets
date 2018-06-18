package main

import (
	"fmt"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Something went wrong: %v\n", err)
		}
	}()

	// trigger panic with division by zero
	for i := 3; i >= 0; i-- {
		fmt.Printf("%d divided by %d is %d\n", 3, i, 3/i)
	}
}
