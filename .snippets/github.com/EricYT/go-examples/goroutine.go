package main

import (
	"fmt"
	"time"
)

var counter int

func main() {
	countc := make(chan int, 1)

	for i := 0; i < 3; i++ {
		go count(countc)
	}

	for {
		value := <-countc
		fmt.Println("channel count ", value)
	}

	time.Sleep(time.Millisecond * 10)
}

func count(c chan int) {
	counter++
	fmt.Println(counter)
	c <- counter
}
