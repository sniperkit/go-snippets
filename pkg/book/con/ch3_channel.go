package main

import (
	"fmt"
	"sync"
)

func simpleExample() {
	stringCh := make(chan string)

	go func() {
		stringCh <- "Hello Channel"
	}()

	fmt.Println(<-stringCh)
}

func begionExample() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func bufferExample() {
	var ch chan interface{}
	ch = make(chan interface{}, 4)

	fmt.Println(<-ch)

	for val := range ch {
		fmt.Printf("%v\n", val)
	}

}

func main() {
	// begionExample()
	bufferExample()
}
