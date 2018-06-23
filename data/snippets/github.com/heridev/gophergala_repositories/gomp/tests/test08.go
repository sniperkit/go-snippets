package main

import "fmt"

func main() {

	p := fmt.Println
	var a [10]int
	//gomp
	for i := 0; i < 10; i++ {
		a[i] = 1
	}
	//gomp
	for i := 0; i < 10; i += a[0] {
		p(10 - i)
	}

}
