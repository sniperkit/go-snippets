package main

import "fmt"

func main() {
	sum := 0
	for i := 0; i<10; i++ {
		sum += i
	}
	fmt.Printf("For %d\n", sum)

	//For continued
	sum = 1
	for ; sum < 1000; {
		sum += sum
	}
	fmt.Printf("For continued %d\n", sum)

	//For is Go's "while"
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Printf("For is Go's \"while\" %d\n", sum)

	//Forever
//	for {
//
//	}
}
