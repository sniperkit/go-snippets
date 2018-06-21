package main

import (
	"fmt"
	"net"
	"log"
)

func main() {
	protocol := "icmp"
	netaddr, err := net.ResolveIPAddr("ip4", "10.71.1.122")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenIP("ip4:"+protocol, netaddr)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 1024)
	for {
		numRead, _,  err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", buf[:numRead])
	}
}
