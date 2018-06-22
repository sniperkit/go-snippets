package main

import "fmt"

func test(a, b int) bool {
	if a > b {
		return true
	}

	return false
}

func main() {
	a := 1

	if a > 2 {
		fmt.Println("a > 2")
	}

	// b := 100
	if b := 1; b >= 1 {
		fmt.Println("b:", b)
		fmt.Println("b > 1")
	}

	if check := test(2, 1); check {
		fmt.Println("a > b")
	}

}
