package main

import (
	"net"
	"fmt"
	"net/http"
	"time"
	_ "net/http/pprof"
)

func metric() error {
	sock, err := net.Listen("tcp", "localhost:8123")
	if err != nil {
		return err
	}
	go func() {
		fmt.Println("HTTP now available at port 8123")
		http.Serve(sock, nil)
	}()

	return nil
}

func main() {
	metric()
	for {
		time.Sleep(time.Second)
	}
}