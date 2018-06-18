package main

import "fmt"

func swap(x, y string) (string, string) {
	return y, x
}

func main() {
	x, y := swap("world", "hello")
	fmt.Printf("This function returns two strings -> %s %s", x, y)
}
