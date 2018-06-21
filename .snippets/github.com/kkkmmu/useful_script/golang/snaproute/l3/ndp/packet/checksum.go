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
	"github.com/google/gopacket/layers"
	"net"
)

func calculateChecksum(content []byte) uint16 {
	var csum uint32
	for i := 0; i < len(content); i += 2 {
		csum += uint32(content[i]) << 8
		csum += uint32(content[i+1])
	}
	return ^uint16((csum >> 16) + csum)
}

/*
 *	          ICMPv6 PSEUDO-HDR MESSAGE FORMAT
 *
 *    0                   1                   2                   3
 *    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |                                                               |
 * +                                                               +
 * |                                                               |
 * +                         Source Address                        +
 * |                                                               |
 * +                                                               +
 * |                                                               |
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |                                                               |
 * +                                                               +
 * |                                                               |
 * +                      Destination Address                      +
 * |                                                               |
 * +                                                               +
 * |                                                               |
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |                   Upper-Layer Packet Length                   |
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |                      zero                     |  Next Header  |
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 */
//func createPseudoHeader(srcIP, dstIP net.IP, icmpv6Hdr *layers.ICMPv6) []byte {
func createPseudoHeader(srcIP, dstIP net.IP, info interface{}) []byte {
	var buf []byte
	/*
	 *   PSEUDO HEADER BYTE START
	 */
	buf = append(buf, srcIP...)
	buf = append(buf, dstIP...)
	buf = append(buf, 0)
	buf = append(buf, 0)
	switch info.(type) {
	case *layers.ICMPv6:
		icmpv6Hdr := info.(*layers.ICMPv6)
		buf = append(buf, byte((ICMP_HDR_LENGTH+len(icmpv6Hdr.LayerPayload()))/256))
		buf = append(buf, byte((ICMP_HDR_LENGTH+len(icmpv6Hdr.LayerPayload()))%256))
	case []byte:
		payload := info.([]byte)
		buf = append(buf, byte((ICMP_HDR_LENGTH+len(payload))/256))
		buf = append(buf, byte((ICMP_HDR_LENGTH+len(payload))%256))
	}
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, ICMPV6_NEXT_HEADER)
	/*
	 *   PSEUDO HEADER BYTE END
	 */
	return buf
}

/*
 * This is used during validating checksum received in packet
 */
func validateChecksum(srcIP, dstIP net.IP, icmpv6Hdr *layers.ICMPv6) error {
	var buf []byte
	buf = append(buf, createPseudoHeader(srcIP, dstIP, icmpv6Hdr)...)
	/*
	 *   ICMPv6 HEADER BYTE START
	 */
	buf = append(buf, icmpv6Hdr.TypeCode.Type())
	buf = append(buf, icmpv6Hdr.TypeCode.Code())
	// add 2 bytes of Checksum..
	for idx := 0; idx < 2; idx++ {
		buf = append(buf, 0)
	}
	// add typebytes which is [4]bytes
	buf = append(buf, icmpv6Hdr.TypeBytes...)
	// Copy the payload which includes TargetAddress & Options..
	buf = append(buf, icmpv6Hdr.LayerPayload()...)
	// Pad to the next 32-bit boundary
	for idx := 0; idx < 4-(len(icmpv6Hdr.LayerPayload())/4); idx++ {
		buf = append(buf, 0)
	}
	/*
	 *   ICMPv6 HEADER BYTE END
	 */
	rv := calculateChecksum(buf)
	if rv != icmpv6Hdr.Checksum {
		return errors.New(fmt.Sprintf("Calculated Checksum 0x%x and wanted checksum is 0x%x",
			rv, icmpv6Hdr.Checksum))
	}
	return nil
}

/*
 * This is used during sending solicitation packet
 */
func getCheckSum(ipv6 *layers.IPv6, payload []byte) uint16 {
	var buf []byte
	buf = append(buf, createPseudoHeader(ipv6.SrcIP, ipv6.DstIP, payload[8:])...)

	buf = append(buf, payload...)
	// Pad to the next 32-bit boundary
	for idx := 0; idx < 4-(len(payload)/4); idx++ {
		buf = append(buf, 0)
	}
	return calculateChecksum(buf)
}
