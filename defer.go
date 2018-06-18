package main

import "fmt"

func main() {
	x := testReturn(true)
	fmt.Println(x)
}

func testReturn(x bool) uint {
	defer second()
	if x {
		first()
		return 1
	} else {
		first()
		return 2
	}
}

func first() {
	fmt.Println("first")
}

func second() {
	fmt.Println("second")
}
