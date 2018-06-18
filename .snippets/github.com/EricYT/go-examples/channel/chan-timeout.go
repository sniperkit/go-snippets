package main

import "fmt"
import "time"

type T struct {
	name  string // name of the object
	value int    // its value
}

func main() {
	c1 := make(chan string, 1)

	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "message 1"
	}()

	select {
	case message := <-c1:
		fmt.Println("message from c1:", message)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
	}

	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "message 2"
	}()

	select {
	case message := <-c1:
		fmt.Println("message from c1:", message)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout 1")
	}
}
