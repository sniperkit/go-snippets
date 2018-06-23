package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	exit1 := make(chan bool)
	var wg sync.WaitGroup

	f := func(exit <-chan bool, wg *sync.WaitGroup) {

		ever := true
		for ever {
			select {
			case <-exit:
				ever = false
				fmt.Println("exit")

			case <-time.After(time.Second * 2):
				fmt.Println("timeout")
			}
		}

		wg.Done()
	}

	wg.Add(1)
	go f(exit1, &wg)

	time.Sleep(time.Second * 10)

	exit1 <- true
	wg.Wait()
}
