package main

import (
	"fmt"
)

func plus(a int, b int) int {
	return a + b
}
func main() {
	res := plus(1, 2)
	fmt.Println("1+2 = ", res)
}
