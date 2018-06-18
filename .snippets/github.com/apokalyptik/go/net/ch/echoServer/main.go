package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/apokalyptik/net/ch"
)

func echoServer(i, o, e chan string) {
	for {
		select {
		case in := <-i:
			o <- fmt.Sprintf("%s", in)
		case <-e:
			return
		}
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	var handler = ch.New(c)
	echoServer(handler.In, handler.Out, handler.Err)
}

var listen = "127.0.0.1:3333"

func init() {
	flag.StringVar(&listen, "listen", listen, "address to listen on")
}

func main() {
	flag.Parse()
	if addr, err := net.ResolveTCPAddr("tcp", listen); err != nil {
		log.Fatal(err)
	} else {
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Fatal(err)
		} else {
			for {
				if conn, err := listener.Accept(); err != nil {
					log.Fatal(err)
				} else {
					go handleConnection(conn)
				}
			}
		}
	}
}
