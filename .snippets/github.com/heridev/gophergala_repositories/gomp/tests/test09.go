package main

import "fmt"

func main() {

	p := fmt.Println
	//gomp
	for i := 0; i < 10; i += 2 {
		for j := 0; j < 20; j++ {
			p(10 - i + j)
		}
	}

}
