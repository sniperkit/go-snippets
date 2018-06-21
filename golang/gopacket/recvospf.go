package main

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

var allSPFRouters = net.IPAddr{IP: net.IPv4(224, 0, 0, 5)}
var allDOtherRouters = net.IPAddr{IP: net.IPv4(224, 0, 0, 5)}

func main() {
	conn, err := net.ListenPacket("ip4:ospf", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	r, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatal(err)
	}

	err = r.JoinGroup(nil, &allSPFRouters)
	if err != nil {
		log.Fatal(err)
	}
	defer r.LeaveGroup(nil, &allSPFRouters)

	buf := make([]byte, 10240)
	for {
		log.Println("Try to receive packet!\n")
		_, _, _, err := r.ReadFrom(buf)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		fmt.Printf("%s\n", string(buf[:]))
	}
}
