package main

import "fmt"

func main() {
	queue := make(chan int, 2)

	queue <- 1
	queue <- 2
	close(queue)

	for value := range queue {
		fmt.Println("Message ", value)
	}
}
