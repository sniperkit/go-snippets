package main

import "fmt"

func main() {
	x := uint(2)
	double(&x)
	fmt.Println(x)
}

/* Pointer */
func double(xPtr *uint) {
	*xPtr *= 2
}
