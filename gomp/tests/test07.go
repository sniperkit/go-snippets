package main

import "fmt"

func main() {

	p := fmt.Println
	a := 123
	b := 123
	//gomp
	for i := a; i < a+2*b; i += b / 32 {
		p(10 - i)
	}

}
