package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan int)
	n := 10
	go func() {
		for i := 0; i < n; i++ {
			fmt.Println("GOGOGO!", i)
			c <- i
		}
		close(c)
	}()

	for v := range c {
		fmt.Println(v)
	}
}
