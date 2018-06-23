package main

import "fmt"

func main() {

	p := fmt.Println
	//gomp
	for i := 0; i < 10; i -= -1 {
		p(10 - i)
	}

}
