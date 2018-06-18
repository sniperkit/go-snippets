package main

import (
	"fmt"
)

func main() {
	words := []string{"foo", "bar", "baz"}
	done := make(chan bool)
	defer close(done)
	for _, word := range words {
		go func(word string) {
			// time.Sleep(1 * time.Second)
			fmt.Println(word)
			done <- true
		}(word)
	}

	// Do concurrent things here

	// This blocks and waits for signal from channel
	for i := 0; i < len(words); i++ {
		<-done
	}
}
