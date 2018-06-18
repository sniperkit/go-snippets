package main

import "fmt"

func main() {

	// Go has pointers. A pointer holds the memory address of a value.
	// The type *T is a pointer to a T value. Its zero value is nil.
	var p *int

	// The & operator generates a pointer to its operand.
	i := 42
	p = &i

	// The * operator denotes the pointer's underlying value.
	fmt.Println(*p) // read i through the pointer p
	*p = 21         // set i through the pointer p

	// This is known as "dereferencing" or "indirecting".
	// Unlike C, Go has no pointer arithmetic.

	var a1 int = 10
	fmt.Printf("Address of a variable: %x\n", &a1)

	// Actual variable declaration.
	var a int = 20

	// Pointer variable declaration.
	var ip *int

	// Store address of a in pointer variable.
	ip = &a
	fmt.Printf("Address of a variable: %x\n", &a)

	// address stored in pointer variable.
	fmt.Printf("Address stored in ip variable: %x\n", ip)

	// access the value using the pointer.
	fmt.Printf("Value of *ip variable: %d\n", *ip)

}
