package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {
	ch := make(chan int, 100)
	finished := make(chan bool, 1)

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(num chan<- int, val int) {
			num <- val
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			wg.Done()
		}(ch, i)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

loop:
	for {
		select {
		case num := <-ch:
			fmt.Println(num)
		case <-finished:
			fmt.Println("finished")
			break loop
		case <-time.After(100 * time.Millisecond):
			fmt.Println("timeout")
			break loop
		}
	}
}
