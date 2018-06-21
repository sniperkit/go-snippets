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

package server

import (
	"fmt"
	//"net"
	"testing"
)

func TestInitRtUtilsTestServer(t *testing.T) {
	fmt.Println("****Init Route Utils Server****")
	StartTestServer()
	TestProcessLogicalIntfCreateEvent(t)
	fmt.Println("****************")
}
func TestConvertIntfStrToIfIndexStr(t *testing.T) {
	fmt.Println("**** TestConvertIntfStrToIfIndexStr event ****")
	var ifIndex string
	var err error
	ifIndex, err = server.ConvertIntfStrToIfIndexStr("lo1")
	fmt.Println("ifIndex,err:", ifIndex, err, " for lo1")
	ifIndex, err = server.ConvertIntfStrToIfIndexStr("lo10")
	fmt.Println("ifIndex,err:", ifIndex, err, " for lo10")
	fmt.Println("***************************************")
}
func TestGetNetworkPrefixFromStrings(t *testing.T) {
	fmt.Println("****TestGetNetworkPrefixFromStrings****")
	ipAddr := "11.1.10.2"
	mask := "255.255.255.0"
	ipPrefix, err := getNetowrkPrefixFromStrings(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	ipAddr = "11.1.10.2/24"
	ipPrefix, err = getNetowrkPrefixFromStrings(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	ipAddr = "11.1.10.35"
	mask = "255.0.0.0"
	ipPrefix, err = getNetowrkPrefixFromStrings(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	fmt.Println("***********************************")
}
func TestGetNetworkPrefixFromCIDR(t *testing.T) {
	fmt.Println("****TestGetNetworkPrefixFromCIDR****")
	ipAddr := "11.1.10.2"
	ipPrefix, err := getNetworkPrefixFromCIDR(ipAddr)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr)
	ipAddr = "11.1.10.2/24"
	ipPrefix, err = getNetworkPrefixFromCIDR(ipAddr)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr)
	ipAddr = "11.1.10.2/33"
	ipPrefix, err = getNetworkPrefixFromCIDR(ipAddr)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr)
	fmt.Println("***********************************")
}
func TestValidateNetworkPrefix(t *testing.T) {
	fmt.Println("****TestValidateNetworkPrefix****")
	ipAddr := "11.1.10.2"
	mask := "255.255.255.0"
	ipPrefix, err := validateNetworkPrefix(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	ipAddr = "11.1.10.2/24"
	ipPrefix, err = validateNetworkPrefix(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	ipAddr = "11.1.10.35"
	mask = "255.255.255.255"
	ipPrefix, err = validateNetworkPrefix(ipAddr, mask)
	fmt.Println("ipPfrefix,err:", ipPrefix, " ", err, " for ipAddr:", ipAddr, " mask:", mask)
	fmt.Println("*********************************")
}
func TestGetPrefixLen(t *testing.T) {
	fmt.Println("****TestGetPrefixLen()****")
	ip := "10.1.10.1"
	netIP, err := getIP(ip)
	if err != nil {
		fmt.Println("netIP invalid")
	}
	prefixLen, err := getPrefixLen(netIP)
	fmt.Println("prefixLen,err:", prefixLen, ",", err, " for ip:", ip)
	ip = "255.255.0.0"
	netIP, err = getIP(ip)
	if err != nil {
		fmt.Println("netIP invalid")
	}
	prefixLen, err = getPrefixLen(netIP)
	fmt.Println("prefixLen,err:", prefixLen, ",", err, " for ip:", ip)
	fmt.Println("**************************")
}
