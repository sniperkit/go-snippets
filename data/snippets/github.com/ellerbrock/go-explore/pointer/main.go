package main

import "fmt"

func main() {

	// a := 42
	var a int = 42

	fmt.Println("value of a: ", a)
	fmt.Println("memory address of a: ", &a)

	// makes b a point of the memory address of a
	// b is a type pointer to a integer
	// * is an operator
	// b := &a
	var b *int = &a

	fmt.Println("value of *b: ", *b)
	fmt.Println("memory address of b: ", b)

	*b = 100

	fmt.Println("value of a: ", a)
	fmt.Println("value of *b: ", *b)
}
