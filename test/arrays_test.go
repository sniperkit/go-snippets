package test

import (
	"fmt"
	"testing"
)

func TestArrays(t *testing.T) {

	// 1. How to declare?
	var a [3]int
	fmt.Printf("1. a = %v\n", a)

	// 2. How to declare and initialize?
	var b = [3]int{1, 2, 3}
	fmt.Printf("2. b = %v\n", b)

	// 3. How to get the size?
	fmt.Printf("3. b size = %d\n", len(b))

	// 4. How to get an element?
	fmt.Printf("4. b[2] = %d\n", b[2])

	// 5. How to traverse?
	fmt.Printf("5. traversing an array:\n")
	for index, value := range b {
		fmt.Printf("\tb[%d]=%d\n", index, value)
	}

}
