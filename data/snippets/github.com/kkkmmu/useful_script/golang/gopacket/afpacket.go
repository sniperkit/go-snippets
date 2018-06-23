package main

import (
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

func main() {
	c, err := net.ListenPacket("ip4:ospf", "0.0.0.0") // OSPF for IPv4
	if err != nil {
		log.Fatal(err)
	}

	en0, err := net.InterfaceByName("ens33")
	if err != nil {
		log.Fatal(err)
	}
	r, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.SetMulticastInterface(en0); err != nil {
		log.Fatal(err)
	}
	if _, err := r.MulticastInterface(); err != nil {
		log.Fatal(err)
	}

	if err := r.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err)
	}
	allSPFRouters := net.IPAddr{IP: net.IPv4(224, 0, 0, 5)}
	if err := r.JoinGroup(en0, &allSPFRouters); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 10000)
	_, _, err = c.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	log.Println("Received a packet")
}
