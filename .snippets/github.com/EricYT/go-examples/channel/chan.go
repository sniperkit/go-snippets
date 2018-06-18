package main

import "fmt"
import "time"
import "sync"

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	message := make(chan string)

	go func() {
		//time.Sleep(time.Millisecond * 3000)
		fmt.Println("Goroutine Block until a receiver is ready")
		message <- "ping"
		wg.Done()
	}()

	time.Sleep(time.Millisecond * 3000)
	mes := <-message
	fmt.Println("Block until receive a message:", mes)

	wg.Wait()

}
