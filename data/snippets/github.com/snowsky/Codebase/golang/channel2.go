package main

import "fmt"

func main() {
	c := make(chan int)
	n := 10
	for i := 0; i < n; i++ {
		go func(x int) {
			fmt.Println("GOGOGO!")
			c <- x
		}(i)
	}

	for v := range c {
		fmt.Println(v)
	}
	// close(c)
}
