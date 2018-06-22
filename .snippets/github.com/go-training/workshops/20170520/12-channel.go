package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	wg.Add(2)
	ch := make(chan int, 50)
	go func(num <-chan int) {
		defer wg.Done()
		for {
			if i, ok := <-num; ok {
				fmt.Println(i)
			} else {
				break
			}
		}
	}(ch)
	go func(ch chan<- int) {
		defer wg.Done()
		ch <- 100
		ch <- 200
		close(ch)
	}(ch)
	wg.Wait()
}
