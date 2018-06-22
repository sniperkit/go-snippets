package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := counter("AUS counter ...")
	d := counter("SFO counter ...")
	// NOTE : these two functions need to wait for each other and count in lockstep
	// Write a fan-in function to avoid this
	for i := 0; i < 10; i++ {
		fmt.Printf("%q\n", <-c)
		fmt.Printf("%q\n", <-d)
	}

	fmt.Println("Program terminated...")
}

// Returns receive-only channel of strings
func counter(msg string) <-chan string {
	c := make(chan string)

	// anynomous function
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // return the channel to the caller
}
