package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

func server() {
	// listen on a port
	ln, err := net.Listen("tcp", ":9527")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// accept
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	// receive a message
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Receive message ", msg)
	}
	c.Close()
}

func client() {
	// connect to server
	c, err := net.Dial("tcp", "127.0.0.1:9527")
	if err != nil {
		fmt.Println(err)
		return
	}

	// send a message
	msg := "hello world"
	fmt.Println("send a message:", msg)
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
}

func main() {
	go server()
	go client()
	go client()

	var input string
	fmt.Scanln(&input)
}
