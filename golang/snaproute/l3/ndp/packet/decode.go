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
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/ndp/config"
	"net"
)

func getEthLayer(pkt gopacket.Packet, eth *layers.Ethernet) error {
	ethLayer := pkt.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return errors.New("Decoding ethernet layer failed")
	}
	*eth = *ethLayer.(*layers.Ethernet)
	return nil
}

func getDot1QLayer(pkt gopacket.Packet, dot1q *layers.Dot1Q) int32 {
	dot1qLayer := pkt.Layer(layers.LayerTypeDot1Q)
	if dot1qLayer == nil {
		return -1
	}
	*dot1q = *dot1qLayer.(*layers.Dot1Q)
	if dot1q.VLANIdentifier == 0 {
		return config.INTERNAL_VLAN
	}

	return int32(dot1q.VLANIdentifier)
}

/*
 *			ICMPv6 MESSAGE FORMAT
 *
 *    0                   1                   2                   3
 *    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 *   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *   |     Type      |     Code      |          Checksum             |
 *   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *   |                           Reserved                            |
 *   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *   |                                                               |
 *   +                                                               +
 *   |                                                               |
 *   +                       Target Address                          +
 *   |                                                               |
 *   +                                                               +
 *   |                                                               |
 *   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *   |   Options ...
 *   +-+-+-+-+-+-+-+-+-+-+-+-
 *
 *  API: given a packet it will fill in ip header and icmpv6
 */
func getIpAndICMPv6Hdr(pkt gopacket.Packet, ipv6Hdr *layers.IPv6, icmpv6Hdr *layers.ICMPv6) error {
	ipLayer := pkt.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil {
		return errors.New("Invalid IPv6 layer")
	}
	*ipv6Hdr = *ipLayer.(*layers.IPv6)
	ipPayload := ipLayer.LayerPayload()
	icmpv6Hdr.DecodeFromBytes(ipPayload, nil)
	return nil
}

func (p *Packet) decodeNS(hdr *layers.ICMPv6, srcIP, dstIP net.IP) (*NDInfo, error) {
	ndInfo := &NDInfo{}
	ndInfo.PktType = layers.ICMPv6TypeNeighborSolicitation
	ndInfo.DecodeNDInfo(hdr.TypeBytes, hdr.LayerPayload())
	if ndInfo.IsTargetMulticast() {
		return nil, errors.New(fmt.Sprintln("Targent Address specified", ndInfo.TargetAddress,
			"is a multicast address"))
	}
	err := ndInfo.ValidateNDSInfo(srcIP, dstIP)
	if err != nil {
		return nil, err
	}
	return ndInfo, nil
}

func (p *Packet) decodeNA(hdr *layers.ICMPv6, srcIP, dstIP net.IP) (*NDInfo, error) {
	ndInfo := &NDInfo{}
	ndInfo.PktType = layers.ICMPv6TypeNeighborAdvertisement
	ndInfo.DecodeNDInfo(hdr.TypeBytes, hdr.LayerPayload())
	if ndInfo.IsTargetMulticast() {
		return nil, errors.New(fmt.Sprintln("Targent Address specified", ndInfo.TargetAddress,
			"is a multicast address"))
	}
	err := ndInfo.ValidateNDAInfo(hdr.TypeBytes, dstIP)
	if err != nil {
		return nil, err
	}
	return ndInfo, nil
}

func (p *Packet) decodeRA(hdr *layers.ICMPv6, srcIP, dstIP net.IP) (*NDInfo, error) {
	//prefixFound := false
	ndInfo := &NDInfo{}
	ndInfo.PktType = layers.ICMPv6TypeRouterAdvertisement
	ndInfo.DecodeRAInfo(hdr.TypeBytes, hdr.LayerPayload())
	err := ndInfo.ValidateRAInfo()
	if err != nil {
		return ndInfo, err
	}
	return ndInfo, nil
}

func (p *Packet) decodeICMPv6Hdr(hdr *layers.ICMPv6, srcIP net.IP, dstIP net.IP) (*NDInfo, error) {
	ndInfo := &NDInfo{}
	var err error
	// Validating checksum received, if success then only start parsing icmp payload
	err = validateChecksum(srcIP, dstIP, hdr)
	if err != nil {
		return nil, err
	}
	typeCode := hdr.TypeCode
	if typeCode.Code() != ICMPV6_CODE {
		return nil, errors.New(fmt.Sprintln("Invalid Code", typeCode.Code()))
	}
	switch typeCode.Type() {
	case layers.ICMPv6TypeNeighborSolicitation:
		ndInfo, err = p.decodeNS(hdr, srcIP, dstIP)

	case layers.ICMPv6TypeNeighborAdvertisement:
		ndInfo, err = p.decodeNA(hdr, srcIP, dstIP)

	case layers.ICMPv6TypeRouterSolicitation:
		return nil, errors.New("Router Solicitation is not yet supported")

	case layers.ICMPv6TypeRouterAdvertisement:
		ndInfo, err = p.decodeRA(hdr, srcIP, dstIP)
	default:
		return nil, errors.New(fmt.Sprintln("Not Supported ICMPv6 Type:", typeCode.Type()))
	}
	if err != nil {
		return nil, err
	}
	return ndInfo, nil
}

/* API: Get IPv6 & ICMPv6 Header
 *      Does Validation of IPv6
 *      Does Validation of ICMPv6
 * Validation Conditions are defined below, if anyone of them do not satisfy discard the packet:
 *  - The IP Hop Limit field has a value of 255, i.e., the packet
 *   could not possibly have been forwarded by a router. <- done
 *
 *  - ICMP Checksum is valid. <- done
 *
 *  - ICMP Code is 0. <- done
 *
 *  - ICMP length (derived from the IP length) is 24 or more octets. <- done
 *
 *  - Target Address is not a multicast address. <- done
 *
 *  - All included options have a length that is greater than zero. <- @TODO: need to add this later
 *
 *  - If the IP source address is the unspecified address, the IP
 *    destination address is a solicited-node multicast address. <- done
 *
 *  - If the IP source address is the unspecified address, there is no
 *    source link-layer address option in the message. <- @TODO: need to be done later
 */
func (p *Packet) DecodeND(pkt gopacket.Packet) (*NDInfo, error) {
	// first decode all the layers
	icmpv6Hdr := &layers.ICMPv6{}
	ipv6Hdr := &layers.IPv6{}
	eth := &layers.Ethernet{}
	dot1q := &layers.Dot1Q{}
	var err error

	// Get Ethernet Layer
	err = getEthLayer(pkt, eth)
	if err != nil {
		return nil, err
	}

	// Get Dot1Q Layer
	dot1qTag := getDot1QLayer(pkt, dot1q)

	// First get ipv6 and icmp6 information
	err = getIpAndICMPv6Hdr(pkt, ipv6Hdr, icmpv6Hdr)
	if err != nil {
		return nil, err
	}

	// Validating ipv6 header
	err = validateIPv6Hdr(ipv6Hdr, icmpv6Hdr.TypeCode.Type())
	if err != nil {
		return nil, err
	}

	// Validating icmpv6 header
	ndInfo, err := p.decodeICMPv6Hdr(icmpv6Hdr, ipv6Hdr.SrcIP, ipv6Hdr.DstIP)
	if err != nil {
		return ndInfo, err
	}

	ndInfo.SrcMac = eth.SrcMAC.String()   // Update SRC MAC From ethernet
	ndInfo.DstMac = eth.DstMAC.String()   // Update DST MAC from ethernet
	ndInfo.SrcIp = ipv6Hdr.SrcIP.String() // copy sender ip address to this
	ndInfo.DstIp = ipv6Hdr.DstIP.String() // copy destination ip
	ndInfo.Dot1Q = dot1qTag
	return ndInfo, nil
}
