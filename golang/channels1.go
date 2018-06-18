package main

import (
	"fmt"
	"math/rand"
	"time"
)

func readIn(c1, c2 <-chan int) <-chan int {
	return nil
}

var done chan bool = make(chan bool)

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; i < 3; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
		done <- true
	}()

	<-done
	return c
}

func main() {
	// c1, c2 := make(chan int)
	// go readIn(c1, c2 <-chan int)
	c := boring("boring!")
	for i := 0; ; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	// go func ()  {
	//   if <-done {
	// 		fmt.Println("You're boring; I'm leaving...")
	// 	}
	// }
	// close(c)
}
