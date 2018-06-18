package main

import (
	"encoding/gob"
	"fmt"
	"net"
	//"time"
)

func server() {
	// listen on a port
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	// receive the message
	var msg string

	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println("Received", msg)
		}
		//time.Sleep(300)
	}

	c.Close()
}

func client() {
	// connect to the server
	c, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// send the message
	messages := []string{
		"Hello World",
		"0123456789",
		"Test",
	}

	//fmt.Println("Sending", msg)

	e := gob.NewEncoder(c)
	for _, m := range messages {
		if err = e.Encode(m); err != nil {
			fmt.Println("+", err)
		}
	}
}

func main() {

	go server()
	go client()

	var input string
	fmt.Scanln(&input)
}
