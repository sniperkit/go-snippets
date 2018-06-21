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

// System
package lacp

import (
	"fmt"
	"net"
)

// 6.4.5 Variables associated with the System
type LacpSystem struct {
	// System Priority
	Actor_System_priority uint16
	// MAC address component of the System Id
	Actor_System [6]uint8
}

func (s *LacpSystem) LacpSystemActorSystemIdSet(Actor_System net.HardwareAddr) {
	s.Actor_System = convertNetHwAddressToSysIdKey(Actor_System)
}

func (s *LacpSystem) LacpSystemActorSystemPrioritySet(Actor_System_priority uint16) {
	s.Actor_System_priority = Actor_System_priority
}

func (s *LacpSystem) LacpSystemConvertSystemIdToString() string {
	sysId := LacpSystemIdGet(*s)
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
		sysId[0],
		sysId[1],
		sysId[2],
		sysId[3],
		sysId[4],
		sysId[5],
		sysId[6],
		sysId[7],
	)
}

//6.3.2 System identification
//The globally unique identifier used to identify a System shall be the concatenation of a globally
//administered individual MAC address and the System Priority. The MAC address chosen may be the
//individual MAC address associated with one of the Aggregation Ports of the System. In the case of DRNI
//(Clause 9), all Portal Systems in a Portal have the same System Identifier, which is provided by the
//concatenation of the Portal’s administrated MAC address (7.4.1.1.4) and the Portal’s System Priority
//(7.4.1.1.5).
//
//Where it is necessary to perform numerical comparisons between System Identifiers, each System Identifier
//is considered to be an eight octet unsigned binary number, constructed as follows:
//
// a) The two most significant octets of the System Identifier comprise the System Priority. The System
//    Priority value is taken to be an unsigned binary number; the most significant octet of the System
//    Priority forms the most significant octet of the System Identifier.
//
// b) The third most significant octet of the System Identifier is derived from the initial octet of the MAC
//    address; the least significant bit of the octet is assigned the value of the first bit of the MAC address,
//    the next most significant bit of the octet is assigned the value of the next bit of the MAC address,
//    and so on. The fourth through eighth octets are similarly assigned the second through sixth octets of
//    the MAC address.
func LacpSystemIdGet(s LacpSystem) [8]uint8 {

	var sysId [8]uint8

	mac := s.Actor_System

	sysId[0] = uint8(s.Actor_System_priority >> 8 & 0xff)
	sysId[1] = uint8(s.Actor_System_priority & 0xff)
	sysId[2] = mac[0]
	sysId[3] = mac[1]
	sysId[4] = mac[2]
	sysId[5] = mac[3]
	sysId[6] = mac[4]
	sysId[7] = mac[5]
	return sysId
}
