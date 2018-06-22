package netrpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"testing"

	"github.com/shohi/gopattern/rpc/server"
)

func TestRPC(t *testing.T) {
	serverAddress := "0.0.0.0"
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	var wg sync.WaitGroup
	wg.Add(1)

	// start server
	go func() {
		l, e := net.Listen("tcp", ":1234")
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
	}()

	// start client
	go func() {
		defer wg.Done()
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}

		// sync way
		args := &server.Args{A: 7, B: 8}
		var reply int
		err = client.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		log.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

		// async way
		quotient := new(server.Quotient)
		divCall := client.Go("Arith.Divide", args, quotient, nil)
		replyCall := <-divCall.Done
		res, _ := replyCall.Reply.(*server.Quotient)
		log.Printf("Arith: %d/%d = %v\n", args.A, args.B, *res)
	}()

	wg.Wait()
}
