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
	"net"
	"strconv"
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

func convertIPv4ToUint32(ip []byte) uint32 {
	var val uint32 = 0

	val = val + uint32(ip[0])
	val = (val << 8) + uint32(ip[1])
	val = (val << 8) + uint32(ip[2])
	val = (val << 8) + uint32(ip[3])

	return val
}

func computeChkSum(pkt []byte) uint16 {
	var csum uint32

	for i := 0; i < len(pkt); i += 2 {
		csum += uint32(pkt[i]) << 8
		csum += uint32(pkt[i+1])
	}
	chkSum := ^uint16((csum >> 16) + csum)
	return chkSum
}

func computeTcpIPChkSum(data []byte, csum uint32) uint16 {
	// to handle odd lengths, we loop to length - 1, incrementing by 2, then
	// handle the last byte specifically by checking against the original
	// length.
	length := len(data) - 1
	for i := 0; i < length; i += 2 {
		// For our test packet, doing this manually is about 25% faster
		// (740 ns vs. 1000ns) than doing it by calling binary.BigEndian.Uint16.
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum + (csum >> 16))
}

func getHWAddr(macAddr string) (mac net.HardwareAddr) {
	mac, err := net.ParseMAC(macAddr)
	if mac == nil || err != nil {
		return nil
	}

	return mac
}

func getIP(ipAddr string) (ip net.IP) {
	ip = net.ParseIP(ipAddr)
	if ip == nil {
		return ip
	}
	ip = ip.To4()
	return ip
}

func convertUint32ToIPv4(val uint32) string {
	p0 := int(val & 0xFF)
	p1 := int((val >> 8) & 0xFF)
	p2 := int((val >> 16) & 0xFF)
	p3 := int((val >> 24) & 0xFF)
	str := strconv.Itoa(p3) + "." + strconv.Itoa(p2) + "." +
		strconv.Itoa(p1) + "." + strconv.Itoa(p0)

	return str
}

func convertUint32ToNetIPv4(val uint32) net.IP {
	return getIP(convertUint32ToIPv4(val))
}
