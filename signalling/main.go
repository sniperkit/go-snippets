/*
An example of cancelling upstream producers when downstream consumers have finished.
*/
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func gen(done <-chan struct{}, nums ...int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}

		}
	}()
	return out
}

func sq(done <-chan struct{}, in chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}

		}
	}()
	return out
}

func main() {
	rand.Seed(time.Now().Unix())
	var nums []int
	for i := 0; i < 100; i++ {
		nums = append(nums, rand.Intn(100))
	}

	done := make(chan struct{})
	defer close(done)

	in := gen(done, nums...)

	// Distribute the sq work across runtime.NumCPU() val of goroutines
	var cs []chan int
	for i := 0; i < runtime.NumCPU(); i++ {
		fmt.Printf("starting worker: %v\n", i+1)
		cs = append(cs, sq(done, in))
	}

	// Consume the merged output from all channels
	for n := range merge(done, cs...) {
		fmt.Println(n)
	}
}

func merge(done chan struct{}, cs ...chan int) chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
