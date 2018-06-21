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

package rpc

import (
	"net"
	"strings"
)

func convertIPStrToUint32(ipStr string) (uint32, bool) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, false
	}

	IP := ip.To4()
	ipAddr := uint32(IP[0])<<24 | uint32(IP[1])<<16 | uint32(IP[2])<<8 | uint32(IP[3])
	return ipAddr, true
}

func parseIPRangeStr(ipAddrRangeStr string) (uint32, uint32, bool) {
	ret := strings.Contains(ipAddrRangeStr, "-")
	if !ret {
		return 0, 0, false
	}
	ipStr := strings.Split(ipAddrRangeStr, "-")
	if len(ipStr) != 2 {
		return 0, 0, false
	}
	lIP := strings.TrimSpace(ipStr[0])
	hIP := strings.TrimSpace(ipStr[1])
	lowerIP, ret := convertIPStrToUint32(lIP)
	if !ret {
		return 0, 0, false
	}
	higherIP, ret := convertIPStrToUint32(hIP)
	if !ret {
		return 0, 0, false
	}
	return lowerIP, higherIP, true
}
