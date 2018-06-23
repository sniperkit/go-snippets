package main

import "fmt"

func main() {
	fmt.Println(c())
}

func c() (i int) {
	defer func() { i++ }() // defer function changes return value i to 2
	return 1
}
