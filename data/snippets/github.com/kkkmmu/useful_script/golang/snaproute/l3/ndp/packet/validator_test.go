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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func TestValidateIpv6UnicastNSHdr(t *testing.T) {
	initPacketTestBasics()
	p := gopacket.NewPacket(nsBaseTestPkt, layers.LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	icmpv6Hdr := &layers.ICMPv6{}
	ipv6Hdr := &layers.IPv6{}
	err := getIpAndICMPv6Hdr(p, ipv6Hdr, icmpv6Hdr)
	if err != nil {
		t.Error("Decoding ipv6 and icmpv6 header failed", err)
	}
	validateIPv6Hdr(ipv6Hdr, layers.ICMPv6TypeNeighborSolicitation)
	if err != nil {
		t.Error("Validating IPv6 Hdr failed", err)
	}
}

func TestValidateICMPv6UnicastNSChecksum(t *testing.T) {
	initPacketTestBasics()
	p := gopacket.NewPacket(nsBaseTestPkt, layers.LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	icmpv6Hdr := &layers.ICMPv6{}
	ipv6Hdr := &layers.IPv6{}
	err := getIpAndICMPv6Hdr(p, ipv6Hdr, icmpv6Hdr)
	if err != nil {
		t.Error("Decoding ipv6 and icmpv6 header failed", err)
	}
	err = validateChecksum(ipv6Hdr.SrcIP, ipv6Hdr.DstIP, icmpv6Hdr)
	if err != nil {
		t.Error("Validating Checksum failed", err)
	}
}

func TestPseudoChecksumBuf(t *testing.T) {
	initPacketTestBasics()
	p := gopacket.NewPacket(naBaseTestPkt, layers.LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	icmpv6Hdr := &layers.ICMPv6{}
	ipv6Hdr := &layers.IPv6{}

	err := getIpAndICMPv6Hdr(p, ipv6Hdr, icmpv6Hdr)
	if err != nil {
		t.Error("Decoding ipv6 and icmpv6 header failed", err)
	}
	buf := createPseudoHeader(ipv6Hdr.SrcIP, ipv6Hdr.DstIP, icmpv6Hdr)
	if buf[39] != ICMPV6_NEXT_HEADER {
		t.Error("creating pseudo header failed")
	}
	if len(buf) != 40 {
		t.Error("invalid pseudo header for checksum calculation")
	}
}
