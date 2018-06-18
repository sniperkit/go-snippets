package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	go func() {
		messages <- "result"
	}()

	time.Sleep(time.Second * 1)

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}
	select {
	case msg := <-messages:
		fmt.Println("received messages", msg)
	default:
		fmt.Println("no messages received")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

// Basic sends and receives on channels are blocking.
// However, we can use select with a default clause to
// implement non-blocking sends, receives, and even non-blocking multi-way selects.
