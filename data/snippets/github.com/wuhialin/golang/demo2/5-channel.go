package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const N = 26

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.Gosched()

	var wait sync.WaitGroup
	wait.Add(2 * N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wait.Done()
			time.Sleep(time.Duration(i * 50000000))
			fmt.Printf("%c", 'a'+i)
		}(i)
		go func(i int) {
			defer wait.Done()
			fmt.Printf("%c", 'A'+i)
		}(i)
	}
	wait.Wait()
}
