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

type NDOptionType byte

const (
	NDOptionTypeSourceLinkLayerAddress NDOptionType = 1
	NDOptionTypeTargetLinkLayerAddress NDOptionType = 2
	NDOptionTypePrefixInfo             NDOptionType = 3
	NDOptionTypeRedirectHeader         NDOptionType = 4
	NDOptionTypeMTU                    NDOptionType = 5
)

const (
	HOP_LIMIT                              = 255
	ICMPV6_CODE                            = 0
	ICMP_HDR_LENGTH                        = 8
	UNSPECIFIED_IP_ADDRESS                 = "::"
	IPV6_ICMPV6_MULTICAST_DST_MAC          = "33:33:00:00:00:00"
	IPV6_ADDRESS_BYTES                     = 16
	IPV6_MULTICAST_BYTE             byte   = 0xff
	IPV6_VERSION                    byte   = 6
	ICMPV6_MIN_LENGTH               uint16 = 24
	ICMPV6_NEXT_HEADER              byte   = 58
	ICMPV6_SOURCE_LINK_LAYER_LENGTH uint16 = 8
	SOLICITATED_NODE_ADDRESS               = "ff02::1:ff00:0000"
	SOLICITATED_SRC_IP                     = "::"

	// Router Advertisement Specific Constants
	ICMPV6_MIN_LENGTH_RA         uint16 = 16
	ICMPV6_MIN_PAYLOAD_LENGTH_RA        = 8
)
