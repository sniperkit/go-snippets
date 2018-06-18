package main

import "fmt"

var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func main() {
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	//Range continued
	fmt.Println("Range continued")
	p := make([]int, 10)
	for i := range p {
		p[i] = 1<<uint(i)
	}
	for _, value := range p{
		fmt.Printf("%d\n", value)
	}
}
