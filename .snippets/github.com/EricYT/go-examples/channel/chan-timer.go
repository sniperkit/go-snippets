package main

import (
	"time"

	"github.com/Unknwon/log"
)

func timer(d time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		time.Sleep(d)
		c <- 1
		close(c)
	}()
	return c
}

func main() {
	for i := 0; i < 10; i++ {
		c := timer(1 * time.Second)
		res := <-c
		log.Debug("main receive res:", res)
	}
}
