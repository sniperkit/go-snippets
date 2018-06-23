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
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"net"
)

const (
	// ip/vrrp header Check Defines
	VRRP_TTL                             = uint8(255)          // ttl for ipv4 packet
	VRRP_HOP_LIMIT                       = 255                 // hop limit for ipv6 packet
	VRRP_PKT_TYPE_ADVERTISEMENT          = 1                   // Only one type is supported which is advertisement
	VRRP_PROTO_ID                        = 112                 // vrrp protocol id
	VRRP_RSVD                            = 0                   // vrrp rsvd bits
	VRRP_HDR_CREATE_CHECKSUM             = 0                   // inital hdr checksum
	VRRP_HEADER_SIZE_EXCLUDING_IPVX      = 8                   // 8 bytes...
	VRRP_IPV4_VERSION               byte = 4                   // ip version
	VRRP_IPV6_VERSION               byte = 6                   // ip version
	VRRP_IPV4_HEADER_MIN_SIZE            = 20                  // min ipv4 header size
	VRRP_PROTOCOL_MAC                    = "01:00:5e:00:00:12" // protocol mac used in the encoding packet
	VRRP_V6_PROTOCOL_MAC                 = "33:33:00:00:00:12" // protocol mac used for ipv6 address
	VRRP_V4_GROUP_IP                     = "224.0.0.18"        // ipv4 group address
	VRRP_V6_GROUP_IP                     = "ff02::12"          // ipv6 group address
	VRRP_HEADER_MIN_SIZE                 = 20                  // min vrrp header size for version2

	// error message from Packet
	VRRP_CHECKSUM_ERR      = "VRRP checksum failure"
	VRRP_INCORRECT_VERSION = "Version is not correct for received VRRP Packet"
	VRRP_INCORRECT_FIELDS  = "Field like type/count ip addr/Advertisement Interval are not valid"
)

type Header struct {
	Version      uint8
	Type         uint8
	VirtualRtrId uint8
	Priority     uint8
	CountIPAddr  uint8
	Rsvd         uint8
	MaxAdverInt  uint16
	CheckSum     uint16
	IpAddr       []net.IP
}

type PacketInfo struct {
	Hdr          *Header
	DstMac       string
	SrcMac       string
	Version      uint8
	Vrid         uint8
	Priority     uint8
	AdvertiseInt uint16
	VirutalMac   string
	IpAddr       string // this is IP Header SRC IP
	Vip          string // this is the ip used in vrrp header
	DstIp        string
	IpType       int
}

func Init() *PacketInfo {
	pktInfo := &PacketInfo{}
	return pktInfo
}

func computeChecksum(content []byte) uint16 {
	var csum uint32
	var rv uint16
	for i := 0; i < len(content); i += 2 {
		csum += uint32(content[i]) << 8
		csum += uint32(content[i+1])
	}
	rv = ^uint16((csum >> 16) + csum)

	return rv
}

func (p *PacketInfo) validateHeader(hdr *Header, layerContent []byte) error {
	if hdr.Version != common.VERSION2 && hdr.Version != common.VERSION3 {
		return errors.New(VRRP_INCORRECT_VERSION)
	}
	// Set Checksum to 0 for verifying checksum
	binary.BigEndian.PutUint16(layerContent[6:8], 0)
	// Verify checksum
	chksum := computeChecksum(layerContent)
	if chksum != hdr.CheckSum {
		debug.Logger.Err(chksum, "!=", hdr.CheckSum)
		return errors.New(fmt.Sprintln(VRRP_CHECKSUM_ERR, "computed CheckSum:", chksum, "received chksum in pkt:", hdr.CheckSum))
	}

	// Verify VRRP fields
	if hdr.CountIPAddr == 0 || hdr.MaxAdverInt == 0 || hdr.Type == 0 {
		return errors.New(VRRP_INCORRECT_FIELDS)
	}
	return nil
}

func pseudoheaderChecksum(ip *layers.IPv6) (csum uint32) {
	for i := 0; i < 16; i += 2 {
		csum += uint32(ip.SrcIP[i]) << 8
		csum += uint32(ip.SrcIP[i+1])
		csum += uint32(ip.DstIP[i]) << 8
		csum += uint32(ip.DstIP[i+1])
	}
	return csum
}

func computeV6Checksum(ipv6Hdr *layers.IPv6, layerContent []byte) uint16 {
	length := uint32(len(layerContent))
	// get ipv6 pseudo header
	csum := pseudoheaderChecksum(ipv6Hdr)
	csum += uint32(ipv6Hdr.NextHeader)
	csum += length & 0xffff
	csum += length >> 16
	//csum := uint32(0)
	cLen := len(layerContent) - 1
	for i := 0; i < cLen; i += 2 {
		csum += uint32(layerContent[i]) << 8
		csum += uint32(layerContent[i+1])
	}
	if len(layerContent)%2 == 1 {
		csum += uint32(layerContent[cLen]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum + (csum >> 16))
}

func (p *PacketInfo) validateV6Header(hdr *Header, ipv6Hdr *layers.IPv6, layerContent []byte) error {
	if hdr.Version != common.VERSION2 && hdr.Version != common.VERSION3 {
		return errors.New(VRRP_INCORRECT_VERSION)
	}
	binary.BigEndian.PutUint16(layerContent[6:8], 0)
	chksum := computeV6Checksum(ipv6Hdr, layerContent)
	if chksum != hdr.CheckSum {
		debug.Logger.Err(chksum, "!=", hdr.CheckSum)
		return errors.New(fmt.Sprintln(VRRP_CHECKSUM_ERR, "computed CheckSum:", chksum, "received chksum in pkt:", hdr.CheckSum))
	}
	return nil
}
