package main

import "fmt"

func main() {
	a, b := doIt(8)

	fmt.Println("square: ", a)
	fmt.Println("is even: ", b)
}

func doIt(num int) (int, bool) {
	return num * num, num%2 == 0
}
