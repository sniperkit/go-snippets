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
	"fmt"
	"github.com/google/gopacket/layers"
	"l3/ndp/debug"
	"net"
)

type NDOption struct {
	Type   NDOptionType
	Length byte
	Value  []byte
}

/*
 *  Struct is super set of NS/NA and RS/RA
 *  Depending on the packet type fill in the necessary information and use it
 */
type NDInfo struct {
	// NS/NA Information
	TargetAddress  net.IP
	PktType        layers.ICMPv6TypeCode
	Dot1Q          int32
	SrcMac         string
	SrcIp          string
	DstIp          string
	DstMac         string
	LearnedIntfRef string
	LearnedIfIndex int32
	// RA Information
	CurHopLimit    uint8
	ReservedFlags  uint8
	RouterLifetime uint16
	ReachableTime  uint32
	RetransTime    uint32

	// For All Types
	Options []*NDOption
}

/*		ND Solicitation Packet Format Rcvd From ICPMv6
 *    0                   1                   2                   3
 *    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
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
 */

func DecodeOptionLayer(payload []byte) *NDOption {
	ndOpt := &NDOption{}
	ndOpt.Type = NDOptionType(payload[0])
	ndOpt.Length = payload[1]
	ndOpt.Value = append(ndOpt.Value, payload[2:]...)
	return ndOpt
}

func (nd *NDInfo) DecodeNDInfo(typeByte, payload []byte) {
	if nd.TargetAddress == nil {
		nd.TargetAddress = make(net.IP, IPV6_ADDRESS_BYTES, IPV6_ADDRESS_BYTES)
	}
	debug.Logger.Debug(fmt.Sprintf("typeByte is 0x%x, 0x%x, 0x%x, 0x%x", typeByte[0], typeByte[1], typeByte[2], typeByte[3]))
	nd.ReservedFlags = typeByte[0] // taking Router:Solicited:Override flags from the message
	debug.Logger.Debug(fmt.Sprintf("reserved flags are 0x%x", nd.ReservedFlags))
	copy(nd.TargetAddress, payload[0:IPV6_ADDRESS_BYTES])
	if len(payload) > IPV6_ADDRESS_BYTES {
		//decode option layer also
		ndOpt := DecodeOptionLayer(payload[IPV6_ADDRESS_BYTES:])
		nd.Options = append(nd.Options, ndOpt)
	}
}

/*
 *  0                   1                   2                   3
 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |     Type      |     Code      |          Checksum             |     <------ ICMPV6 Info Start - 0
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  | Cur Hop Limit |M|O|  Reserved |       Router Lifetime         |     <------ ICMPV6 Ends here  - 4
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |                         Reachable Time                        |	  <------ MIN RA Info Start - 8, 0 - 4
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |                          Retrans Timer                        |     <------ MIN RA Info Ends  - 12, 4 - 8
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |   Options ...
 *  +-+-+-+-+-+-+-+-+-+-+-+-
 *  @TODO: Handle Prefix Information Option Type
 */
func (nd *NDInfo) DecodeRAInfo(typeByte, payload []byte) {
	nd.CurHopLimit = typeByte[0]
	nd.ReservedFlags = typeByte[1]
	nd.RouterLifetime = binary.BigEndian.Uint16(typeByte[2:4])
	nd.ReachableTime = binary.BigEndian.Uint32(payload[:4])
	nd.RetransTime = binary.BigEndian.Uint32(payload[4:8])
	// if more than min payload length then it means that we have got options
	if len(payload) > ICMPV6_MIN_PAYLOAD_LENGTH_RA {
		for base := ICMPV6_MIN_PAYLOAD_LENGTH_RA; base < len(payload); base = base + 8 {
			if base+8 > len(payload) {
				break
			}
			ndOpt := DecodeOptionLayer(payload[base:(base + 8)])
			nd.Options = append(nd.Options, ndOpt)
		}
	}
}

/*
 *  According to RFC 2375 https://tools.ietf.org/html/rfc2375 all ipv6 multicast address have first byte as
 *  FF or 0xff, so compare that with the Target address first byte.
 */
func (nd *NDInfo) IsTargetMulticast() bool {
	if nd.TargetAddress.IsMulticast() {
		return true
	}
	return false
}
