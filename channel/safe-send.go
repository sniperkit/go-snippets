package main

import "fmt"

type T int

func SafeSend(c chan<- T, value T) (closed bool) {
	defer func() {
		if recover() != nil {
			// the function result can be altered
			// in a defer function calling
			closed = true
		}
	}()
	c <- value   // panic if c is closed
	return false // <<=>> closed = false; return
}

func main() {
	c := make(chan T, 1)
	fmt.Println("send a message result: ", SafeSend(c, T(1)))
	close(c)
	fmt.Println("send a message again result: ", SafeSend(c, T(2)))
}
