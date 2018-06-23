package main

import (
	"fmt"
	"time"
)

func main() {
	c := counter("SFO")
	timeout := time.After(2 * time.Second)

	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-timeout:
			fmt.Println("Was timed out...")
			return
		}
	}

}

// Returns receive-only channel of strings
func counter(msg string) <-chan string {
	c := make(chan string)

	// anynomous function
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			//time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // return the channel to the caller
}
