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
package server

import (
	"fmt"
	"l3/ndp/debug"
	"net"
	"strings"
)

func createNeighborKey(mac, ip, intfName string) string {
	return mac + "_" + ip + "_" + intfName
}

func splitNeighborKey(nbrKey string) []string {
	return strings.Split(nbrKey, "_")
}

func isLinkLocal(ipAddr string) bool {
	ip, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		ip = net.ParseIP(ipAddr)
	}
	return ip.IsLinkLocalUnicast() && (ip.To4() == nil)
}

/*
func baseFilter(macAddr string) (filter string) {
	filter = fmt.Sprintf("%s%s%s", NDP_PCAP_FILTER, NDP_ETHER_SRC)
	debug.Logger.Info("new filter is:", filter)
	return filter
}
*/

func getNewFilter(macAddr string) (filter string) {
	filter = fmt.Sprintf("%s%s%s", NDP_PCAP_FILTER, NDP_ETHER_SRC, macAddr)
	debug.Logger.Info("updating filter to:", filter)
	return filter
}

func (svr *NDPServer) IsIPv6Addr(ipAddr string) bool {
	ip, _, _ := net.ParseCIDR(ipAddr)
	if ip.To4() == nil {
		return true
	}

	return false
}
