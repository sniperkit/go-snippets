package main

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"ospf"
)

var OSPF = ospf.Header{
	Version:        2,
	Type:           1,
	PacketLength:   24,
	RouterID:       net.IPv4(10, 1, 1, 1),
	AreaID:         net.IPv4(0, 0, 0, 1),
	CheckSum:       0,
	AuType:         0,
	Authentication: 0,
}

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

	var ospf bytes.Buffer
	binary.Write(&ospf, binary.BigEndian, &OSPF)

	iph := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0xc0, // DSCP CS6
		TotalLen: ipv4.HeaderLen + ospf.Len(),
		TTL:      1,
		Protocol: 89,
		Dst:      allSPFRouters.IP.To4(),
	}

	var buffer bytes.Buffer
	os, err := OSPF.Marshal()
	if err != nil {
		panic(err)
	}
	binary.Write(&buffer, binary.BigEndian, os)
	OSPF.CheckSum = OSPFCheckSum(buffer.Bytes())
	os, err = OSPF.Marshal()
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, os)
	var cm *ipv4.ControlMessage
	cm = &ipv4.ControlMessage{IfIndex: en0.Index}

	if err := r.WriteTo(iph, buffer.Bytes(), cm); err != nil {
		log.Fatal(err)
	}
}

func OSPFCheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}
