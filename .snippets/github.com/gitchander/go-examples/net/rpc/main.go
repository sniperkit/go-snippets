package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Server struct{}

func (this *Server) Negate(i int64, reply *int64) error {
	*reply = -i
	return nil
}

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}
func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	param := 781
	var result int64
	err = c.Call("Server.Negate", int64(param), &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Server.Negate(%d) = %d", param, result)
	}
}

func main() {
	go server()
	go client()

	var input string
	fmt.Scanln(&input)
}
