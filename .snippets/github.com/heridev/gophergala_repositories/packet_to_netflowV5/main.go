package main

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"fmt"
	"net"
)

func main() {
	if handle, err := pcap.OpenLive("lo", 1600, true, 0); err != nil {
		panic(err)
	} else {
		var eth layers.Ethernet
		var ip4 layers.IPv4
		var ip6 layers.IPv6
		var tcp layers.TCP
		var udp layers.UDP

		// flow fields
		var sourceIP net.IP
		var destinationIP net.IP
		var ipProtocol layers.IPProtocol
		var sourcePort uint16
		var destinationPort uint16
		var tOS uint8

		// flowFields := make(map[string])

		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp)
		decoded := []gopacket.LayerType{}

		for {
			packetData, _, _ := handle.ReadPacketData()
			parser.DecodeLayers(packetData, &decoded)
			for _, layerType := range decoded {
				switch layerType {
				case layers.LayerTypeIPv6:
					sourceIP = ip6.SrcIP
					destinationIP = ip6.DstIP
					ipProtocol = ip6.NextHeader
				case layers.LayerTypeIPv4:
					sourceIP = ip4.SrcIP
					destinationIP = ip4.DstIP
					tOS = ip4.TOS
					ipProtocol = ip4.Protocol
				case layers.LayerTypeTCP:
					sourcePort = uint16(tcp.SrcPort)
					destinationPort = uint16(tcp.DstPort)
				case layers.LayerTypeUDP:
					sourcePort = uint16(udp.SrcPort)
					destinationPort = uint16(udp.DstPort)

				}
				fmt.Println("Results: ", sourceIP, destinationIP, ipProtocol, sourcePort, destinationPort, tOS)
			}
		}

	}
}

type FlowPayload struct {
	SourceIP        [4]byte
	DestinationIp   [4]byte
	NextHop         [4]byte
	Input           [2]byte
	Output          [2]byte
	dPkts           [4]byte
	dOctets         [4]byte
	First           [4]byte
	Last            [4]byte
	SourcePort      [2]byte
	DestinationPort [2]byte
	Pad1            byte
	TCPFlags        byte
	Prot            byte
	Tos             byte
	SrcAS           [2]byte
	DstAs           [2]byte
	SrcMask         byte
	DstMask         byte
	Pad2            [2]byte
}
