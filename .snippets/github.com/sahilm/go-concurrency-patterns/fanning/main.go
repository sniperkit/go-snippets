/*
An example of fanning-out work and then fanning it in by merging channels. The merge fn takes n channels
and returns a single out channel with n channels multiplexed onto it.
*/
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func gen(nums ...int) chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in chan int) chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	rand.Seed(time.Now().Unix())
	var nums []int
	for i := 0; i < 100; i++ {
		nums = append(nums, rand.Intn(100))
	}

	in := gen(nums...)

	// Distribute the sq work across runtime.NumCPU() val of goroutines
	var cs []chan int
	for i := 0; i < runtime.NumCPU(); i++ {
		fmt.Printf("starting worker: %v\n", i+1)
		cs = append(cs, sq(in))
	}

	// Consume the merged output from all channels
	for n := range merge(cs...) {
		fmt.Println(n)
	}
}

func merge(cs ...chan int) chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
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
