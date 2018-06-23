package main

import (
	"fmt"
)

type MyError struct {
	name string
	val  int
}

func (e MyError) Error() string {
	return fmt.Sprintf("Error: name: %s, value: %d", e.name, e.val)
}

func printF() error {
	return MyError{
		name: "workshop",
		val:  100,
	}
}

func main() {
	err := printF()
	fmt.Println(err)
}
