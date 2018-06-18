package main

import (
	"fmt"
)

func main() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two" // chan 在同一个GOroutine中
	close(queue)

	for elem := range queue {
		fmt.Println(elem)
	}
}
