package main

import (
	"fmt"
	"time"
)

func main() {
	// make a channel that capital is 2
	message := make(chan string, 2)

	go func() {
		message <- "foo"
		message <- "bar"
		fmt.Println("goroutine send message without synchronization")
	}()

	fmt.Println("Sleep a time to test whether the goroutine is blocked")
	time.Sleep(3 * time.Second)

	fmt.Println(<-message)
	fmt.Println(<-message)
}
