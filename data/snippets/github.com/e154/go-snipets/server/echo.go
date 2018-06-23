package main

import (
	"net"
	"fmt"
	"log"
	"io"
)

const (
	address = "localhost"
	port = 9999
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		p, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go io.Copy(p, p)
	}
}
