package main

import "fmt"

func main() {
	ch := make(chan int)
	n := 10
	for i := 0; i < n; i++ {
		go func(c chan int, x int) {
			fmt.Println("GOGOGO!", x)
			c <- x
		}(ch, i)
	}

	for v := range ch {
		fmt.Println(v)
	}
	// close(c)
}
