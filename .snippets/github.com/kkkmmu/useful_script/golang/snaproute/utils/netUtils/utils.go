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

// netUtils.go
package netUtils

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
	"utils/patriciaDB"
)

func GetNetowrkPrefixFromStrings(ipAddr string, mask string) (prefix patriciaDB.Prefix, err error) {
	destNetIpAddr, err := GetIP(ipAddr)
	if err != nil {
		fmt.Println("destNetIpAddr invalid")
		return prefix, err
	}
	networkMaskAddr, err := GetIP(mask)
	if err != nil {
		fmt.Println("networkMaskAddr invalid")
		return prefix, err
	}
	prefix, err = GetNetworkPrefix(destNetIpAddr, networkMaskAddr)
	if err != nil {
		fmt.Println("err=", err)
		return prefix, err
	}
	return prefix, err
}
func GetNetworkPrefixFromCIDR(ipAddr string) (ipPrefix patriciaDB.Prefix, err error) {
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return ipPrefix, err
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	//ipMaskStr := net.IP(ipMask).String()
	ipPrefix, err = GetNetowrkPrefixFromStrings(ipAddrStr, (net.IP(ipNet.Mask)).String()) //ipMaskStr)

	return ipPrefix, err
}
func GetIPInt(ip net.IP) (ipInt int, err error) {
	if ip == nil {
		fmt.Printf("ip address %v invalid\n", ip)
		return ipInt, errors.New("Invalid destination network IP Address")
	}
	ip = ip.To4()
	parsedPrefixIP := int(ip[3]) | int(ip[2])<<8 | int(ip[1])<<16 | int(ip[0])<<24
	ipInt = parsedPrefixIP
	return ipInt, nil
}

func GetIP(ipAddr string) (ip net.IP, err error) {
	ip = net.ParseIP(ipAddr)
	if ip == nil {
		return ip, errors.New("Invalid destination network IP Address")
	}
	return ip, nil
}
func IsZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}
func IsZerosIPString(ipAddr string) (bool, error) {
	ip, err := GetIP(ipAddr)
	if err != nil {
		fmt.Println("invalid IP address")
		return false, errors.New("Invalid IP address")
	}
	if IsIPv4Mask(ip) {
		return IsZeros(ip[12:15]), nil
	} else {
		return IsZeros(ip), nil
	}
	//fmt.Println("ip:", ip, "len(ip):", len(ip), "ip[12:15]:", ip[12:15], " net.IP(ipAddr):", net.IP(ipAddr))
	return IsZeros(ip), nil
}
func IsIPv4Mask(mask net.IP) bool {
	if IsZeros(mask[0:10]) &&
		mask[10] == 0xff &&
		mask[11] == 0xff {
		return true
	}
	return false
}

func GetPrefixLen(networkMask net.IP) (prefixLen int, err error) {
	//fmt.Println("GetPrefixLen() for mask:", networkMask)
	mask := net.IPMask(networkMask)
	//fmt.Println("mask:", mask)
	if IsIPv4Mask(net.IP(mask)) {
		prefixLen, _ = mask[12:16].Size()
	} else {
		prefixLen, _ = mask.Size()
	}
	//fmt.Println("prefixLen = ", prefixLen, " err:", err)
	return prefixLen, err
}
func GetNetworkPrefix(destNetIp net.IP, networkMask net.IP) (destNet patriciaDB.Prefix, err error) {
	prefixLen, err := GetPrefixLen(networkMask)
	if err != nil {
		fmt.Println("err when getting prefixLen, err= ", err)
		return destNet, errors.New(fmt.Sprintln("Invalid networkmask ", networkMask))
	}
	//	fmt.Println("prefixLen for :", networkMask, ":", prefixLen)
	var netIp net.IP
	vdestMask := net.IPMask(networkMask)
	if IsIPv4Mask(net.IP(vdestMask)) {
		//	fmt.Println("v4,vdestMask[12:16]:", vdestMask[12:16])
		netIp = destNetIp.Mask(vdestMask[12:16])
	} else {
		netIp = destNetIp.Mask(vdestMask)
	}
	if netIp == nil {
		fmt.Println("netip nil ")
		return destNet, errors.New("netIp nil")
	}
	numbytes := prefixLen / 8
	if (prefixLen % 8) != 0 {
		numbytes++
	}
	destNet = make([]byte, numbytes)
	//	fmt.Println("numbytes:", numbytes, " netIp:", netIp, " len(netIp):", len(netIp))
	for i := 0; i < numbytes && i < len(netIp); i++ {
		destNet[i] = netIp[i]
	}
	return destNet, err
}
func GetCIDR(ipAddr string, mask string) (addr string, err error) {
	destNetIpAddr, err := GetIP(ipAddr)
	if err != nil {
		fmt.Println("destNetIpAddr invalid")
		return addr, err
	}
	maskIP, err := GetIP(mask)
	if err != nil {
		fmt.Println("err in getting mask IP for mask string", mask)
		return addr, err
	}
	prefixLen, err := GetPrefixLen(maskIP)
	if err != nil {
		fmt.Println("err in getting prefix len for mask string", mask)
		return addr, err
	}
	addr = (destNetIpAddr.Mask(net.IPMask(maskIP))).String() + "/" + strconv.Itoa(prefixLen)
	return addr, err
}
func CheckIfInRange(testIPAddr, ipAddr string, lowPrefixLen int, highPrefixLen int) bool {
	//fmt.Println("testIPAddr:", testIPAddr, " ipAddr:", ipAddr, " lowPrefixLen:", lowPrefixLen, " highPrefixLen:", highPrefixLen)
	testIPStrings := strings.SplitAfter(testIPAddr, "/")
	if len(testIPStrings) != 2 {
		fmt.Println("Invalid cidr :", testIPAddr)
		return false
	}
	testPrefixLenStr := testIPStrings[1]
	testPrefixLen, _ := strconv.Atoi(testPrefixLenStr)

	baseIPStrings := strings.SplitAfter(ipAddr, "/")
	if len(baseIPStrings) != 2 {
		fmt.Println("Invalid cidr :", ipAddr)
		return false
	}
	basePrefixLenStr := baseIPStrings[1]
	basePrefixLen, _ := strconv.Atoi(basePrefixLenStr)

	testAddr, _, err := net.ParseCIDR(testIPAddr)
	if err != nil {
		fmt.Println("error parsing test address:", testIPAddr)
		return false
	}
	if lowPrefixLen == -1 && highPrefixLen == -1 {
		//exact case
		if testPrefixLen != basePrefixLen {
			//fmt.Println("Prefix len ", testPrefixLen, "is not the exact prefix as", basePrefixLen, " for ", testIPAddr, " and ", ipAddr)
			return false
		}
	} else {
		//for a range
		if testPrefixLen > highPrefixLen || testPrefixLen < lowPrefixLen {
			//fmt.Println("Prefix len ", testPrefixLen, " not with range for ", lowPrefixLen, " and", highPrefixLen, " for ", testIPAddr, " and ", ipAddr)
			return false
		}
	}
	testCidr := testAddr.String() + "/" + basePrefixLenStr
	testIpPrefix, err := GetNetworkPrefixFromCIDR(testCidr)
	if err != nil {
		fmt.Println("Invalid ipPrefix for the route ", testAddr)
		return false
	}
	ipPrefix, err := GetNetworkPrefixFromCIDR(ipAddr)
	if err != nil {
		fmt.Println("Invalid ipPrefix for the route ", ipAddr)
		return false
	}
	if bytes.Equal(testIpPrefix, ipPrefix) {
		return true
	}

	return false
}

func IsIPv6Addr(ipAddr string) bool {
	ip, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		ip = net.ParseIP(ipAddr)
		if ip == nil {
			return false
		}
	}
	ip1 := ip.To4()
	if ip1 == nil {
		return true
	}
	ip2 := ip.To16()
	if len(ip1) == len(ip2) {
		return true
	}
	return false
}

func IsIPv4Addr(ipAddr string) bool {
	ip, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		ip = net.ParseIP(ipAddr)
		if ip == nil {
			return false
		}
	}
	ip1 := ip.To4()
	if ip1 == nil {
		return false
	}
	ip2 := ip.To16()
	if len(ip1) != len(ip2) {
		return true
	}
	return false
}

func IsIpv6LinkLocal(ipAddr string) bool {
	ip, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		ip = net.ParseIP(ipAddr)
	}
	return ip.IsLinkLocalUnicast() && (ip.To4() == nil)
}

const (
	SUBNET_MASK_31_BIT  = 31
	SUBNET_MASK_32_BIT  = 32
	SUBNET_MASK_128_BIT = 128
)

// accepts input only in CIDR format
func IsIpAddrValid(ipAddr string, ipType int) bool {
	ip, ipnet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return false
	}
	ipStr := ipnet.String()
	size, _ := ipnet.Mask.Size()
	if ipType == syscall.AF_INET {
		if size == SUBNET_MASK_31_BIT || size == SUBNET_MASK_32_BIT {
			return true
		}
		ip = ip.To4()
	} else if ipType == syscall.AF_INET6 {
		if size == SUBNET_MASK_128_BIT {
			return true
		}
		ip = ip.To16()
	} else {
		return false // invalid ip type
	}
	subnet := ipnet.Mask
	for idx, _ := range subnet {
		subnet[idx] = byte(subnet[idx]) ^ 0xff
	}
	broadcast := make(net.IP, 0)
	for idx, _ := range ip {
		broadcast = append(broadcast, ip[idx]|subnet[idx])
	}
	broadcastIp := broadcast.String() + "/" + strconv.Itoa(size)
	if strings.Compare(ipAddr, ipStr) == 0 || strings.Compare(ipAddr, broadcastIp) == 0 {
		return false
	}
	return true
}
