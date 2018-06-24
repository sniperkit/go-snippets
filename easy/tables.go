package main

import "fmt"

func main() {
	for i := 1; i < 13; i++ {
		for j := 1; j < 13; j++ {
			fmt.Printf("%4d", i*j)
		}
		fmt.Println()
	}
}
