package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("nonblocking-go")

	var wg sync.WaitGroup
	wg.Add(3)

	var signal chan struct{} = make(chan struct{})
	var pipe chan int = make(chan int)

	// goroutine 1
	go func() {
		defer wg.Done()
		var in chan<- struct{}
		for {
			select {
			case now := <-time.After(time.Second * 3):
				fmt.Println("goroutine 1 signal time: ", now)
				in = signal
			case in <- struct{}{}:
				in = nil
			}
		}
	}()

	// goroutine 2
	go func() {
		defer wg.Done()

		var index int
		var out chan<- int
		for {
			select {
			case now := <-time.After(time.Second * 1):
				fmt.Println("goroutine 2 ticket: ", now)
			case <-signal:
				fmt.Println("goroutine 2 receive signal now: ", time.Now())
				index++
				//pipe <- index
				out = pipe
			case out <- index:
				fmt.Println("goroutine 2 send index now: ", time.Now())
				out = nil
			}
		}
	}()

	// goroutine 3
	go func() {
		defer wg.Done()
		for {
			select {
			case index := <-pipe:
				fmt.Println("goroutine 3 receive index: ", index, " now: ", time.Now())
				time.Sleep(time.Second * 7)
			}
		}
	}()

	wg.Wait()
}
