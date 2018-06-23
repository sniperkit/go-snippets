package main

import (
	"fmt"
	"runtime"
	"sync"
)

var a int = 1

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg := sync.WaitGroup{}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go GO(&wg, i)
	}
	wg.Wait()
}

func GO(wg *sync.WaitGroup, index int) {
	// a := 1
	for i := 0; i < 100000000; i++ {
		a += i
	}
	fmt.Println(index, a)
	// if index == 9 {
	// 	c <- true
	// }
	wg.Done()
}
