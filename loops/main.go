package main

import "fmt"

func main() {

	for i := 0; i <= 3; i++ {
		fmt.Println(i)
	}

	fmt.Println("----")

	x := 0
	for x <= 3 {
		fmt.Println(x)
		x++
	}

	fmt.Println("----")

	x = 0
	for {
		fmt.Println(x)
		if x >= 3 {
			break
		}
		x++
	}
}
