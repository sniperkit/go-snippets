package main

import (
	"fmt"
	"log"
	"net"
	"golang.org/x/net/ipv4"
)

func main() {
	conn, err := net.ListenPacket("ip4:icmp", "10.71.1.122")
	if err != nil {
		log.Fatal(err)
	}

	r, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 1024)
	for {
		log.Println("Try to receive packet!\n")
		_, _, _, err := r.ReadFrom(buf)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		fmt.Printf("%v\n", buf[:])
	}
}
