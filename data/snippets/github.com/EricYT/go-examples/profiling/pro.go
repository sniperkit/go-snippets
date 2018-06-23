package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	var count int = 10
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			wg.Done()
			var sum int = 1
			for j := 0; j < 100000; j++ {
				sum = sum + j
			}
			result <- sum
		}()
	}
	wg.Wait()

	for {
		fmt.Println("select start")
		select {
		case s := <-result:
			fmt.Println("get result:", s)
		case <-time.After(time.Second * 10):
			fmt.Println("main process over")
			return
		}
	}
}
