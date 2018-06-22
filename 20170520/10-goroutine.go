package main

import (
	"fmt"
	// "time"
	"sync"
)

// func hello() {
// 	fmt.Println("Hi, workshop")
// }

func main() {
	wg := sync.WaitGroup{}
	name := "workshop"
	wg.Add(1)
	go func(name string) {
		fmt.Println("Hi,", name)
		wg.Done()
	}(name)
	wg.Wait()
	// time.Sleep(100 * time.Millisecond)

	name = "golang"
	ch := make(chan bool)
	go func(name string) {
		fmt.Println("Hi,", name)
		ch <- true
	}(name)
	<-ch
}
