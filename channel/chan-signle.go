package main

import "fmt"

func ping(to chan<- string, message string) {
	to <- message
}

func pong(pong chan<- string, ping <-chan string) {
	mes := <-ping
	pong <- mes
}

func main() {
	pingc := make(chan string, 1)
	pongc := make(chan string, 1)

	ping(pingc, "passed message")
	pong(pongc, pingc)

	fmt.Println(<-pongc)
}
