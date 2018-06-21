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
	"net"
)

var ProtocolFamilyDefaultRouteMap = map[uint32]NLRI{
	GetProtocolFamily(AfiIP, SafiUnicast):  NewIPPrefix(net.ParseIP("0.0.0.0"), uint8(0)),
	GetProtocolFamily(AfiIP6, SafiUnicast): NewIPPrefix(net.ParseIP("::"), uint8(0)),
}

func ConstructNLRIForDefaultRoutes(update map[uint32]bool) (valid, invalid map[uint32][]NLRI) {
	for protoFamily, add := range update {
		ipPrefix, ok := ProtocolFamilyDefaultRouteMap[protoFamily]
		if !ok {
			continue
		}
		if add {
			if valid == nil {
				valid = make(map[uint32][]NLRI)
			}
			valid[protoFamily] = make([]NLRI, 1)
			valid[protoFamily][0] = ipPrefix
		} else {
			if invalid == nil {
				invalid = make(map[uint32][]NLRI)
			}
			invalid[protoFamily] = make([]NLRI, 1)
			invalid[protoFamily][0] = ipPrefix
		}
	}

	return valid, invalid
}
