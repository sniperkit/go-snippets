package main

import "fmt"
import "time"

func main() {

	chan1 := make(chan string)
	chan2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		chan1 <- "message 1"
	}()

	go func() {
		time.Sleep(time.Second * 1)
		chan2 <- "message 2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case mess1 := <-chan1:
			fmt.Println("message from :", mess1)
		case mess2 := <-chan2:
			fmt.Println("message from :", mess2)
		}
	}
}
