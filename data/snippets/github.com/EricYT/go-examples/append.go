package main

import (
	"log"
)

func append(slice, data []byte) []byte {
	l := len(slice)
	if l+len(data) > cap(slice) { //relocate
		newSlice := make([]byte, ((l + len(data)) * 2))
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : 1+len(data)]
	for i, c := range data {
		slice[l+i] = c
	}
	return slice
}

func main() {
	log.Println("log test message")
}
