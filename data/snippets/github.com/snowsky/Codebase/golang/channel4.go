package main

import "fmt"

func main() {
	c := make(chan int)
	// done := make(chan bool)
	n := 10
	for i := 0; i < n; i++ {
		go func(x int) {
			fmt.Println("GOGOGO!")
			c <- x
			close(c)
		}(i)
	}

	go func() {
		for v := range c {
			fmt.Println(v)
		}
	}()

	// close(c)
}
