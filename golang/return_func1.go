package main

import "fmt"

func main() {
	p2 := plusX(3)
	fmt.Printf("%v\n", p2(2))
}
func plusTwo() func(int) int {
	return func(x int) int { return x + 2 }
}

func plusX(x int) func(int) int {
	return func(y int) int { return x + y }
}
