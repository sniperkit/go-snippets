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

// bgp.go
package packet

import (
	"errors"
	"fmt"
	"l3/bgp/config"
	"l3/rib/ribdCommonDefs"
	"net"
)

type AFI uint16
type SAFI uint8

const (
	AfiIP AFI = iota + 1
	AfiIP6
)

const (
	SafiUnicast SAFI = iota + 1
	SafiMulticast
)

var ProtocolFamilyMap = map[string]uint32{
	"ipv4-unicast": GetProtocolFamily(AfiIP, SafiUnicast),
	"ipv6-unicast": GetProtocolFamily(AfiIP6, SafiUnicast),
	//"ipv4-multicast": GetProtocolFamily(AfiIP, SafiMulticast),
	//"ipv6-multicast": GetProtocolFamily(AfiIP6, SafiMulticast),
}

var AFINextHopLenMap = map[AFI]int{
	AfiIP:  4,
	AfiIP6: 16,
}

var AFINextHop = map[AFI]net.IP{
	AfiIP:  net.IPv4zero,
	AfiIP6: net.IPv6zero,
}

var RIBdAddressTypeToAFI = map[ribdCommonDefs.IPType]AFI{
	ribdCommonDefs.IPv4: AfiIP,
	ribdCommonDefs.IPv6: AfiIP6,
}

var PeerAddrTypeProtoFamilyMap = map[config.PeerAddressType]uint32{
	config.PeerAddressV4: GetProtocolFamily(AfiIP, SafiUnicast),
	config.PeerAddressV6: GetProtocolFamily(AfiIP6, SafiUnicast),
}

func GetProtocolFromConfig(afiSafis *[]config.AfiSafiConfig, neighborAddress net.IP) (map[uint32]bool, bool) {
	afiSafiMap := make(map[uint32]bool)
	rv := true
	for _, afiSafi := range *afiSafis {
		if afiSafiVal, ok := ProtocolFamilyMap[afiSafi.AfiSafiName]; ok {
			afiSafiMap[afiSafiVal] = true
		} else {
			rv = false
			break
		}
	}

	if len(afiSafiMap) == 0 {
		if neighborAddress.To4() == nil {
			afiSafiMap[ProtocolFamilyMap["ipv6-unicast"]] = true
		} else {
			afiSafiMap[ProtocolFamilyMap["ipv4-unicast"]] = true
		}
	}
	return afiSafiMap, rv
}

func GetProtocolFamily(afi AFI, safi SAFI) uint32 {
	return uint32(afi<<8) | uint32(safi)
}

func GetAfiSafi(protocolFamily uint32) (AFI, SAFI) {
	return AFI(protocolFamily >> 8), SAFI(protocolFamily & 0xFF)
}

func GetAddressLengthForFamily(protoFamily uint32) int {
	afi, _ := GetAfiSafi(protoFamily)
	if addrLen, ok := AFINextHopLenMap[afi]; ok {
		return addrLen
	}
	return -1
}

func GetProtocolFamilyFromAddressLen(bytes int) uint32 {
	for afi, length := range AFINextHopLenMap {
		if length == bytes {
			return GetProtocolFamily(afi, SafiUnicast)
		}
	}
	return 0
}

func GetZeroNextHopForFamily(protoFamily uint32) net.IP {
	afi, _ := GetAfiSafi(protoFamily)
	if nh, ok := AFINextHop[afi]; ok {
		return nh
	}
	return nil
}

func GetProtocolFromOpenMsg(openMsg *BGPOpen) map[uint32]bool {
	afiSafiMap := make(map[uint32]bool)
	for _, optParam := range openMsg.OptParams {
		if capabilities, ok := optParam.(*BGPOptParamCapability); ok {
			for _, capability := range capabilities.Value {
				if val, ok := capability.(*BGPCapMPExt); ok {
					afiSafiMap[GetProtocolFamily(val.AFI, val.SAFI)] = true
				}
			}
		}
	}

	return afiSafiMap
}

func GetProtocolFamilyFromAddrType(addrType ribdCommonDefs.IPType) (uint32, error) {
	if afi, ok := RIBdAddressTypeToAFI[addrType]; ok {
		return GetProtocolFamily(afi, SafiUnicast), nil
	}

	return 0, errors.New(fmt.Sprintf("Address family not found for address type %d", addrType))
}

func GetProtocolFamilyFromPeerAddrType(addrType config.PeerAddressType) (uint32, error) {
	if pf, ok := PeerAddrTypeProtoFamilyMap[addrType]; ok {
		return pf, nil
	}

	return 0, errors.New(fmt.Sprintf("Address family not found for peer address type %d", addrType))
}
