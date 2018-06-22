package main

import (
	"fmt"
)

func calc(vals ...int) {
	fmt.Println("array count:", len(vals))
	for index, val := range vals {
		fmt.Println(index, ":", val)
	}
}

func main() {
	// calc()
	// calc(1, 2)
	// calc(1, 2, 3)
	calc([]int{1, 2, 3, 4, 5}...)
}
