package main

import "fmt"

func main() {
	for i := 1; i <= 20; i++ {
		if i < 10 {
			fmt.Printf("0%v\n", i)
		} else {
			fmt.Printf("%v\n", i)
		}
	}
}
