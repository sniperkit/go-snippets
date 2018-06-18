package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)
	for i := 0; i < 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(time.Millisecond * 200)

	for req := range requests {
		<-limiter
		fmt.Println("Request", req, time.Now())
	}

	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Millisecond * 2000) {
			burstyLimiter <- t
		}
	}()

	burstyRequest := make(chan int, 5)
	for i := 0; i < 5; i++ {
		burstyRequest <- i
	}
	close(burstyRequest)
	for req := range burstyRequest {
		<-burstyLimiter
		fmt.Println("Request", req, time.Now())
	}
}
