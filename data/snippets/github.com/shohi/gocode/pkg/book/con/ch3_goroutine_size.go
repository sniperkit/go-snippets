package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c
	}

	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	before := memConsumed()

	for i := numGoroutines; i > 0; i-- {
		go noop()
	}

	wg.Wait()
	after := memConsumed()
	fmt.Printf("before: %.3f Kb\n", float64(before/1000))
	fmt.Printf("before: %.3f Kb\n", float64(after/1000))
	fmt.Printf("per size: %.3f kb\n", float64((after-before)/numGoroutines/1000))
	fmt.Printf("active goroutines: %v\n", runtime.NumGoroutine())
}
