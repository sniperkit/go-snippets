package main

import (
	"fmt"
)

func main() {
	var a [6]int
	fmt.Println("empty:", a)

	a[3] = 100
	fmt.Println("set: ", a)
	fmt.Println("len: ", len(a))

	//b := [5]int{1, 2, 3, 4, 5}
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
