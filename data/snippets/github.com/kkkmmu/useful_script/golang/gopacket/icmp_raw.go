package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.ListenPacket("ip4:icmp", "10.71.1.122")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 1024)
	for {
		numRead, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", buf[:numRead])
	}
}
