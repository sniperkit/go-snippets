//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package netUtils

import (
	"fmt"
	"net"
	"syscall"
	"testing"
)

type IPRange struct {
	testAddr       string
	baseAddr       string
	lowPrefixLen   int
	highPrefixLen  int
	expectedresult bool
}

var IPRangeData []IPRange

func IPAddrStringToU8List(ipAddr string) []uint8 {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return ip
	}
	return ip
}
func TestInitNetUtils(t *testing.T) {
	IPRangeData = make([]IPRange, 0)
	IPRangeData = append(IPRangeData, IPRange{"192.168.1.1/31", "192.168.1.0/24", 24, 32, true})
	IPRangeData = append(IPRangeData, IPRange{"192.168.2.1/31", "192.168.1.0/24", 26, 32, false})
	IPRangeData = append(IPRangeData, IPRange{"192.168.2.1/31", "192.168.0.0/16", 16, 32, true})
	IPRangeData = append(IPRangeData, IPRange{"192.167.2.1/31", "192.168.0.0/16", 16, 32, false})
	IPRangeData = append(IPRangeData, IPRange{"192.168.1.1/31", "200.1.1.0", 24, 32, false})
	IPRangeData = append(IPRangeData, IPRange{"192.168.0.7/31", "192.168.0.0/26", -1, -1, false})
	IPRangeData = append(IPRangeData, IPRange{"2003::11:1:10:1/127", "5001:6000:7000::0/48", 48, 128, false})
	IPRangeData = append(IPRangeData, IPRange{"2003::11:1:10:1/127", "2003:11:1::0/64", 64, 128, false})
	IPRangeData = append(IPRangeData, IPRange{"2003::11:1:10:1/127", "2003::0/64", 64, 128, true})
	IPRangeData = append(IPRangeData, IPRange{"5001::11:1:10:1/127", "5001:6000:7000::0", 48, 128, false})
	IPRangeData = append(IPRangeData, IPRange{"5001:6000:7000::11:1:10:1/127", "5001:6000:7000::0/48", 48, 128, true})
	IPRangeData = append(IPRangeData, IPRange{"2001::172:16:0:29/127", "2001::172:16:0:0/96", -1, -1, false})
	IPRangeData = append(IPRangeData, IPRange{"2000::192:16:0:29/127", "2000::192:16:0:0/96", -1, -1, false})
	IPRangeData = append(IPRangeData, IPRange{"fe80::e0:ecff:fe26:a7f0/128", "fe80::/10", -1, -1, false})
	IPRangeData = append(IPRangeData, IPRange{"fe80::e0:ecff:fe26:a7f0/128", "fe80::/10", 10, 128, true})
	IPRangeData = append(IPRangeData, IPRange{"1000:192:168::3/128", "1000::/16", 16, 128, true})
	IPRangeData = append(IPRangeData, IPRange{"192.168.0.2/31", "192.168.0.3/31", -1, -1, true})
	IPRangeData = append(IPRangeData, IPRange{"192.168.0.90/31", "192.168.0.0/26", 26, 32, false})
	IPRangeData = append(IPRangeData, IPRange{"192.168.0.10/32", "192.168.0.0/24", 31, 32, true})
	IPRangeData = append(IPRangeData, IPRange{"172.32.3.39/32", "172.31.14.148/20", 20, 32, true})
}
func TestGetNetworkPrefix(t *testing.T) {
	fmt.Println("****TestGetNetworkPrefix****")
	ip := "10.1.10.1"
	mask := "255.255.255.0"
	prefix, err := GetNetworkPrefix(net.IP(ip), net.IP(mask))
	fmt.Println("prefix:", prefix, " err:", err)
	fmt.Println("****************")
}

func TestGetNetworkPrefixFromStrings(t *testing.T) {
	fmt.Println("****TestGetNetworkPrefixFromString****")
	ip := "10.1.10.1"
	mask := "255.255.255.0"
	prefix, err := GetNetowrkPrefixFromStrings(ip, mask)
	fmt.Println("prefix:", prefix, " err:", err)
	ip = "10.1.10.1"
	mask = "255.255.255.254"
	prefix, err = GetNetowrkPrefixFromStrings(ip, mask)
	ip = "2000::192:16:0:19"
	mask = "255.255.255.254"
	prefix, err = GetNetowrkPrefixFromStrings(ip, mask)
	fmt.Println("prefix:", prefix, " err:", err)
	ip = "2222:22:1:1::"
	mask = "ffff:ffff:ffff:ffff::"
	prefix, err = GetNetowrkPrefixFromStrings(ip, mask)
	fmt.Println("prefix:", prefix, " err:", err)
	ip = "2222:22:1:2::"
	mask = "ffff:ffff:ffff:ffff::"
	prefix, err = GetNetowrkPrefixFromStrings(ip, mask)
	fmt.Println("prefix:", prefix, " err:", err)
	fmt.Println("****************")
}

func TestGetPrefixLen(t *testing.T) {
	fmt.Println("****TestGetPrefixLen()****")
	ip := "255.255.255.0"

	netIP, err := GetIP(ip)
	if err != nil {
		fmt.Println("netIP invalid")
	}
	prefixLen, err := GetPrefixLen(netIP)
	fmt.Println("netIp:", netIP)
	fmt.Println("prefixLen:", prefixLen, " err:", err, " for ip:", ip)

	ip = "0.0.0.0"
	netIP, err = GetIP(ip)
	if err != nil {
		fmt.Println("netIP invalid")
	}
	prefixLen, err = GetPrefixLen(netIP)
	fmt.Println("prefixLen:", prefixLen, " err:", err, " for ip:", ip)

	ip = "255.254.0.0"
	parsedIP := IPAddrStringToU8List(ip)
	fmt.Println("parsedIP:", parsedIP, " for ip:", ip)
	prefixLen, err = GetPrefixLen(parsedIP)
	fmt.Println("prefixLen:", prefixLen, " err:", err, " for ip:", ip)

	ip = "255.254.0.0"
	prefixLen, err = GetPrefixLen(net.IP(ip))
	fmt.Println("prefixLen:", prefixLen, " err:", err, " for ip:", ip)

	ip = "11.1.10.2"
	netIP, err = GetIP(ip)
	if err != nil {
		fmt.Println("netIP invalid")
	}
	fmt.Println("netIp:", netIP)
	prefixLen, err = GetPrefixLen(netIP)
	fmt.Println("prefixLen:", prefixLen, " err:", err, " for ip:", ip)
}

func TestGetNetworkPrefixFromCIDR(t *testing.T) {
	fmt.Println("****TestGetNetworkPrefixFromCIDR****")
	ip := "10.1.10.1/24"
	prefix, err := GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "10.1.10.0/24"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "192.168.11.1/31"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "fe80::/64"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "80.16.16.16/32"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "5010:1010::/32"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	ip = "2000::192:16:0:18/31"
	prefix, err = GetNetworkPrefixFromCIDR(ip)
	fmt.Println("prefix:", prefix, " err:", err, " for ip:", ip)
	fmt.Println("****************")
}

func TestCheckIfInRange(t *testing.T) {
	fmt.Println("****TestCheckIfInRange()****")
	for _, data := range IPRangeData {
		result := CheckIfInRange(data.testAddr, data.baseAddr, data.lowPrefixLen, data.highPrefixLen)
		if result != data.expectedresult {
			t.Error("match result for ", data, "is ", result, " expected result is:", data.expectedresult)
		}
		//fmt.Println("match result for ", data, "is ", result, " expected result is:", data.expectedresult)
	}
	fmt.Println("****************************")
}
func TestIsZeroesIPString(t *testing.T) {
	fmt.Println("****TestIsZeroesIPString()****")
	ipAddr := "0.0.0.0"
	isZeroes, err := IsZerosIPString(ipAddr)
	fmt.Println("isZeroes:", isZeroes, " err:", err, " for ipAddr:", ipAddr)
	ipAddr = "10.10.10.10"
	isZeroes, err = IsZerosIPString(ipAddr)
	fmt.Println("isZeroes:", isZeroes, " err:", err, " for ipAddr:", ipAddr)
	ipAddr = "0:0:0:0:0:0:0:0"
	isZeroes, err = IsZerosIPString(ipAddr)
	fmt.Println("isZeroes:", isZeroes, " err:", err, " for ipAddr:", ipAddr)
	fmt.Println("*******************************")
}

func TestCheckIPv4Address(t *testing.T) {
	ip := "10.1.10.1/24"
	fmt.Println("****TestCheckIPv4Address()****")
	if IsIPv4Addr(ip) == false {
		t.Error(ip, "is ipv4 address")
	}
	if IsIPv6Addr(ip) == true {
		t.Error(ip, "is not ipv6 address")
	}

	if IsIPv4Addr("10.1.10.20") != true {
		t.Error("10.1.10.20 should fail the check")
	}

	if IsIPv6Addr("10.1.10.20") == true {
		t.Error("10.1.10.20 is ipv4 address not ipv6")
	}
	fmt.Println("****************************")
}

func TestCheckIPv6Address(t *testing.T) {
	ip := "2003::2/64"
	fmt.Println("****TestCheckIPv6Address()****")
	if IsIPv6Addr(ip) == false {
		t.Error(ip, "is ipv6 address")
	}

	if IsIPv4Addr(ip) == true {
		t.Error(ip, "is not ipv4 address")
	}
	if IsIPv6Addr("2003::2") != true {
		t.Error("2003::2 should fail the check")
	}
	if IsIPv4Addr("2003::2") == true {
		t.Error("2003::2 is ipv6 address not ipv4")
	}
	fmt.Println("****************************")
}

func TestGetCIDR(t *testing.T) {
	fmt.Println("****TestGetCIDR()****")
	ipAddr := "40.1.10.0"
	mask := "255.255.255.0"
	cidrAddr, err := GetCIDR(ipAddr, mask)
	fmt.Println("cidrAddr:", cidrAddr, " err:", err, " for ipAddr:", ipAddr, " mask:", mask)
	fmt.Println("****************************")
}

func TestIpIpAddrValid(t *testing.T) {
	fmt.Println("****TestIsIpAddrValid()****")
	ipAddr := "192.16.0.0/31"
	valid := IsIpAddrValid(ipAddr, syscall.AF_INET)
	if !valid {
		t.Error("ip Addr should be valid:", ipAddr)
		return
	}
	ipAddr = "192.16.0.0/32"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET)
	if !valid {
		t.Error("ip Addr should be valid:", ipAddr)
		return
	}
	ipAddr = "192.16.0.0/24"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET)
	if valid {
		t.Error("ip Addr should not be valid:", ipAddr)
		return
	}

	ipAddr = "1000::192:16:0:1/64"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET6)
	if !valid {
		t.Error("ip Addr should be valid:", ipAddr)
		return
	}
	ipAddr = "1000::192:16:0:1/128"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET6)
	if !valid {
		t.Error("ip Addr should be valid:", ipAddr)
		return
	}
	ipAddr = "1000::192:16:0:1/127"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET6)
	if valid {
		t.Error("ip Addr should not be valid:", ipAddr)
		return
	}
	ipAddr = "3000::e/126"
	valid = IsIpAddrValid(ipAddr, syscall.AF_INET6)
	if !valid {
		t.Error("ip Addr should be valid:", ipAddr)
		return
	}

	fmt.Println("****************************")
}
