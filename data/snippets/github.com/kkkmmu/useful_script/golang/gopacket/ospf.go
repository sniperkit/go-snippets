package main

import (
	"log"
	"net"
	"runtime"

	"golang.org/x/net/ipv4"
)

func main() {
	c, err := net.ListenPacket("ip4:89", "0.0.0.0") // OSPF for IPv4
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	r, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}

	en0, err := net.InterfaceByName("ens33")
	if err != nil {
		log.Fatal(err)
	}
	allSPFRouters := net.IPAddr{IP: net.IPv4(224, 0, 0, 5)}
	if err := r.JoinGroup(en0, &allSPFRouters); err != nil {
		log.Fatal(err)
	}
	defer r.LeaveGroup(en0, &allSPFRouters)

	hello := make([]byte, 24) // fake hello data, you need to implement this
	ospf := make([]byte, 24)  // fake ospf header, you need to implement this
	ospf[0] = 2               // version 2
	ospf[1] = 2               // hello packet
	ospf = append(ospf, hello...)
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0xc0, // DSCP CS6
		TotalLen: ipv4.HeaderLen + len(ospf),
		TTL:      1,
		Protocol: 89,
		Dst:      allSPFRouters.IP.To4(),
	}

	var cm *ipv4.ControlMessage
	switch runtime.GOOS {
	case "darwin", "linux":
		cm = &ipv4.ControlMessage{IfIndex: en0.Index}
	default:
		if err := r.SetMulticastInterface(en0); err != nil {
			log.Fatal(err)
		}
	}
	if err := r.WriteTo(iph, ospf, cm); err != nil {
		log.Fatal(err)
	}
}
