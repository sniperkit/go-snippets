package main

import (
	"fmt"
	"time"
)

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}()
	return c
}

func main() {
	joe := boring("Joe")
	ann := boring("Ann")
	for i := 0; i < 10; i++ {
		fmt.Printf("You say: %q\n", <-joe)
		fmt.Printf("You say: %q\n", <-ann)
	}
	fmt.Println("You're boring; I'm leaving.")
}
