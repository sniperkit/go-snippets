package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10)
	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
			time.Sleep(time.Second) // 这行与上一行交换之后只打印出0-3
		}
		done <- true
	}()

	go func() {
		for {
			fmt.Println(<-c)
		}
	}()

	time.Sleep(time.Second)

	close(c)
	if <-done {
		return
	}
}
