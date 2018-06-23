package main

import "fmt"

func main() {

	p := fmt.Println
	//gomp
	for i := 0; i < 134; i += 123 {
		p(10 - i)
	}

}
