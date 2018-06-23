package main

import (
	"encoding/binary"
	//	"net"
	"fmt"
)

type Message struct {
	Version  uint8
	Length   uint32
	Sequence uint64
	Message  []byte
}

func main() {
	msg := &Message{
		Version:  100,
		Length:   256,
		Sequence: 1234,
		Message:  []byte("Hello world"),
	}

	var a uint8
	var b int32
	var c float64

	fmt.Println(binary.Size(a))
	fmt.Println(binary.Size(b))
	fmt.Println(binary.Size(c))
	fmt.Println(binary.Size(msg))
}
