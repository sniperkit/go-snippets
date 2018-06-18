package main

import "fmt"

func main() {
	doIt(1, 2, 3, 4, 5)
}

func doIt(numbers ...int) {
	for _, v := range numbers {
		fmt.Println(v)
	}
}
