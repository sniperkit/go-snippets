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
	"net"
)

/*
 *  Range for Solicited Node Multicast Address from RFC 4291 FF02:0:0:0:0:1:FF00:0000 to FF02:0:0:0:0:1:FFFF:FFFF
 *  if srcIp == "::", i.e Unspecified address then dstIP should be solicited-node address FF02:0:0:0:0:1:FFXX:XXXX
 *  if srcIP == "::", then there should not be any source link-layer option in message
 */
func (nd *NDInfo) ValidateNDSInfo(srcIP net.IP, dstIP net.IP) error {
	if srcIP.IsUnspecified() {
		if !(dstIP[0] == IPV6_MULTICAST_BYTE && dstIP[1]&0x0f == 0x02 &&
			dstIP[11]&0x0f == 0x01 && dstIP[12] == IPV6_MULTICAST_BYTE) {
			return errors.New(fmt.Sprintln("Destination IP address",
				dstIP.String(), "is not Solicited-Node Multicast Address"))
		}
		options := nd.Options
		if len(options) > 0 {
			for _, option := range options {
				if option.Type == NDOptionTypeSourceLinkLayerAddress {
					return errors.New(fmt.Sprintln("During ND Solicitation with Unspecified",
						"address Source Link Layer Option should not be set"))
				}
			}
		}
	}
	return nil
}

/*
 * If the IP Destination Address is a multicast address the
 *       Solicited flag is zero.
 * All included options have a length that is greater than zero.
 */
func (nd *NDInfo) ValidateNDAInfo(icmpFlags []byte, dstIP net.IP) error {
	if dstIP.IsMulticast() {
		flags := binary.BigEndian.Uint16(icmpFlags[0:2])
		if (flags & 0x4000) == 0x4000 {
			return errors.New(fmt.Sprintln("Check for If Destination Address is a multicast",
				"address then the Solicited flag is zero, Failed"))
		}
	}
	// @TODO: need to add support for options length
	return nil
}

func validateIPv6Hdr(hdr *layers.IPv6, layerType uint8) error {
	if hdr.HopLimit != HOP_LIMIT {
		return errors.New(fmt.Sprintln("Invalid Hop Limit", hdr.HopLimit))
	}
	switch layerType {
	case layers.ICMPv6TypeNeighborSolicitation, layers.ICMPv6TypeNeighborAdvertisement:
		if hdr.Length < ICMPV6_MIN_LENGTH {
			return errors.New(fmt.Sprintf("Invalid ICMP length %d", hdr.Length))
		}
	case layers.ICMPv6TypeRouterAdvertisement:
		if hdr.Length < ICMPV6_MIN_LENGTH_RA {
			return errors.New(fmt.Sprintf("Invalid ICMP length %d", hdr.Length))
		}
	}
	return nil
}

/*
 * Validate
 *	- All included options have a length that is greater than zero.
 *
 * Cache below information during validation
 *	- Source Link-Layer Address
 *	- Prefix Information
 *	- MTU options
 */
func (nd *NDInfo) ValidateRAInfo() error {
	options := nd.Options
	if len(options) > 0 {
		for _, option := range options {
			switch option.Type {
			case NDOptionTypeSourceLinkLayerAddress:
				if option.Length == 0 {
					return errors.New(fmt.Sprintln("During Router Advertisement",
						"Source Link Layer Option has length as zero"))
				}
			case NDOptionTypeMTU:
				if option.Length == 0 {
					return errors.New(fmt.Sprintln("During Router Advertisement",
						"MTU Option has length as zero"))
				}
			}
		}
	}
	return nil
}
