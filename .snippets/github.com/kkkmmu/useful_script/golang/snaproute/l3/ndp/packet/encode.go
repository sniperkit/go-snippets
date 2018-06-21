//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//
package packet

import (
	"encoding/binary"
	_ "fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/ndp/debug"
	"net"
)

func (pkt *Packet) constructEthLayer() *layers.Ethernet {
	// Ethernet Layer Information
	srcMAC, _ := net.ParseMAC(pkt.SrcMac)
	dstMAC, _ := net.ParseMAC(pkt.DstMac)
	eth := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv6,
	}

	return eth
}

func (pkt *Packet) constructIPv6Layer() *layers.IPv6 {
	// IPv6 Layer Information
	sip := net.ParseIP(pkt.SrcIp)
	dip := net.ParseIP(pkt.DstIp)

	ipv6 := &layers.IPv6{
		Version:      IPV6_VERSION,
		TrafficClass: 0,
		NextHeader:   layers.IPProtocolICMPv6,
		SrcIP:        sip,
		DstIP:        dip,
		HopLimit:     HOP_LIMIT,
	}

	return ipv6
}

func constructICMPv6NS(srcMac net.HardwareAddr, ipv6 *layers.IPv6) []byte {
	// ICMPV6 Layer Information
	payload := make([]byte, ICMPV6_MIN_LENGTH)
	payload[0] = byte(layers.ICMPv6TypeNeighborSolicitation)
	payload[1] = byte(0)
	binary.BigEndian.PutUint16(payload[2:4], 0) // Putting zero for checksum before calculating checksum
	binary.BigEndian.PutUint32(payload[4:], 0)  // RESERVED FLAG...
	copy(payload[8:], ipv6.DstIP.To16())

	// Append Source Link Layer Option here
	srcOption := NDOption{
		Type:   NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
		Value:  srcMac,
	}
	payload = append(payload, byte(srcOption.Type))
	payload = append(payload, srcOption.Length)
	payload = append(payload, srcOption.Value...)
	binary.BigEndian.PutUint16(payload[2:4], getCheckSum(ipv6, payload))
	return payload
}

func constructICMPv6RA(srcMac net.HardwareAddr, ipv6 *layers.IPv6) []byte {
	// ICMPV6 Layer Information
	payload := make([]byte, ICMPV6_MIN_LENGTH_RA)
	payload[0] = byte(layers.ICMPv6TypeRouterAdvertisement)
	payload[1] = byte(0)
	binary.BigEndian.PutUint16(payload[2:4], 0) // Putting zero for checksum before calculating checksum
	payload[4] = byte(64)
	payload[5] = byte(0)
	binary.BigEndian.PutUint16(payload[6:8], 1800) // Router Lifetime
	binary.BigEndian.PutUint32(payload[8:12], 0)   // reachable time
	binary.BigEndian.PutUint32(payload[12:16], 0)  // retrans time

	// Append Source Link Layer Option here
	srcOption := NDOption{
		Type:   NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
		Value:  srcMac,
	}

	mtuOption := NDOption{
		Type:   NDOptionTypeMTU,
		Length: 1,
		// Reserved is added as first 2 bytes in value
		Value: []byte{0x00, 0x00, 0x00, 0x00, 0x05, 0xdc}, // 1500 hardcoded need to change this
	}
	payload = append(payload, byte(srcOption.Type))
	payload = append(payload, srcOption.Length)
	payload = append(payload, srcOption.Value...)

	payload = append(payload, byte(mtuOption.Type))
	payload = append(payload, mtuOption.Length)
	payload = append(payload, mtuOption.Value...)
	binary.BigEndian.PutUint16(payload[2:4], getCheckSum(ipv6, payload))
	return payload
}

func (pkt *Packet) Encode() []byte {
	eth := pkt.constructEthLayer()
	ipv6 := pkt.constructIPv6Layer()

	var icmpv6Payload []byte
	switch pkt.PType {
	case layers.ICMPv6TypeNeighborSolicitation:
		icmpv6Payload = constructICMPv6NS(eth.SrcMAC, ipv6)
	case layers.ICMPv6TypeRouterAdvertisement:
		icmpv6Payload = constructICMPv6RA(eth.SrcMAC, ipv6)
	}

	ipv6.Length = uint16(len(icmpv6Payload))
	// GoPacket serialized buffer that will be used to send out raw bytes
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	//debug.Logger.Debug("Sending pkt (DMAC, SMAC):(", eth.DstMAC.String(), ",", eth.SrcMAC.String(), ") and (SIP,DIP):(", ipv6.SrcIP.String(), ",", ipv6.DstIP.String(), ")")

	err := gopacket.SerializeLayers(buffer, options, eth, ipv6, gopacket.Payload(icmpv6Payload))
	if err != nil {
		//fmt.Println("serialize layers failed, err:", err)
		debug.Logger.Err("serialize layers failed, err:", err)
	}
	return buffer.Bytes()
}
