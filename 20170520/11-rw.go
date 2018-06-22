package main

import (
	"fmt"
	"runtime"
	"sync"
)

var count = 0
var wg = sync.WaitGroup{}

// var m = sync.RWMutex{}

func main() {
	// runtime.GOMAXPROCS(-1)
	// fmt.Println(runtime.GOMAXPROCS(-1))
	for i := 0; i < 10; i++ {
		wg.Add(2)
		// m.RLock()
		go sayHello()
		// m.Lock()
		go addCount()
	}
	wg.Wait()
}

func sayHello() {
	// m.RLock()
	fmt.Printf("Hello %#v\n", count)
	// m.RUnlock()
	wg.Done()
}

func addCount() {
	// m.Lock()
	count++
	// m.Unlock()
	wg.Done()
}
