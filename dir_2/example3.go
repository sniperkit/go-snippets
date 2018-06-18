package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5}
	
	// Print a[1] .. a[3]
	fmt.Println(a[1:4])
	
	// Print a[0] .. a[2]
	fmt.Println(a[:3])
	
	// Print a[2] .. len(a)
	fmt.Println(a[2:])
}
