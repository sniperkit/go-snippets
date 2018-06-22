package main

import (
	"fmt"
	"time"
)

func main() {
	c := fanIn(sleepCounter("SFO", 2), sleepCounter("AUS", 4))
	for i := 1; i < 15; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("Excellent example of an fan-in function")
}

// Combines the output of both the channels
func fanIn(input1 <-chan string, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			// Send to channel c, recv from input1
			c <- <-input1
		}
	}()
	go func() {
		for {
			// Send to channel c, recv from input2
			c <- <-input2
		}
	}()

	return c
}

// Returns receive-only channel of strings
func sleepCounter(name string, wait int) <-chan string {
	c := make(chan string)

	// goroutine inside a function
	go func() {
		for i := 1; ; i++ {
			c <- fmt.Sprintf("%s %d", name, i)
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}()
	return c // return the channel to the caller
}
