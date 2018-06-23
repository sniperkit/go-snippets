package main

import "fmt"

type T int

func IsClosed(c <-chan T) bool {
	select {
	case <-c:
		return true
	default:
	}
	return false
}

func main() {
	c := make(chan T)
	fmt.Println("Channel closed: ", IsClosed(c))
	close(c)
	fmt.Println("Channel closed: ", IsClosed(c))
}
