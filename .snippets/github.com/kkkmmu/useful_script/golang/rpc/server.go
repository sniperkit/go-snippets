package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Echo int

func (t *Echo) Hi(args string, reply *string) error {
	*reply = "echo:" + args
	return nil
}

func main() {
	rpc.Register(new(Echo))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
