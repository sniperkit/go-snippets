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
	"utils/clntUtils/clntDefs/asicdClntDefs"
	//"l3/rib/ribdCommonDefs"
	"testing"
)

var logicalIntfList []asicdClntDefs.LogicalIntfNotifyMsg
var vlanList []asicdClntDefs.VlanNotifyMsg
var ipv4IntfList []asicdClntDefs.IPv4IntfNotifyMsg
var ipv6IntfList []asicdClntDefs.IPv6IntfNotifyMsg

func InitLogicalIntfList() {
	logicalIntfList = make([]asicdClntDefs.LogicalIntfNotifyMsg, 0)
	logicalIntfList = append(logicalIntfList, asicdClntDefs.LogicalIntfNotifyMsg{
		IfIndex:         1,
		LogicalIntfName: "lo1",
	})
	logicalIntfList = append(logicalIntfList, asicdClntDefs.LogicalIntfNotifyMsg{
		IfIndex:         2,
		LogicalIntfName: "lo2",
	})
	logicalIntfList = append(logicalIntfList, asicdClntDefs.LogicalIntfNotifyMsg{
		IfIndex:         3,
		LogicalIntfName: "lo3",
	})
	logicalIntfList = append(logicalIntfList, asicdClntDefs.LogicalIntfNotifyMsg{
		IfIndex:         4,
		LogicalIntfName: "lo4",
	})
	logicalIntfList = append(logicalIntfList, asicdClntDefs.LogicalIntfNotifyMsg{
		IfIndex:         5,
		LogicalIntfName: "lo5",
	})
}
func InitVlanList() {
	vlanList = make([]asicdClntDefs.VlanNotifyMsg, 0)
	vlanList = append(vlanList, asicdClntDefs.VlanNotifyMsg{
		VlanId:   100,
		VlanName: "vlan100",
	})
	vlanList = append(vlanList, asicdClntDefs.VlanNotifyMsg{
		VlanId:   200,
		VlanName: "vlan200",
	})
}
func InitIPv4IntfList() {
	ipv4IntfList = make([]asicdClntDefs.IPv4IntfNotifyMsg, 0)
	ipv4IntfList = append(ipv4IntfList, asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "11.1.10.1/24",
		IfIndex: 1,
	})
	ipv4IntfList = append(ipv4IntfList, asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "21.1.10.1/24",
		IfIndex: 2,
	})
	ipv4IntfList = append(ipv4IntfList, asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "31.1.10.1/24",
		IfIndex: 3,
	})
	ipv4IntfList = append(ipv4IntfList, asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "35.1.10.1/24",
		IfIndex: 35,
	})
}
func InitIPv6IntfList() {
	ipv6IntfList = make([]asicdClntDefs.IPv6IntfNotifyMsg, 0)
	ipv6IntfList = append(ipv6IntfList, asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2002::1/64",
		IfIndex: 4,
	})
	ipv6IntfList = append(ipv6IntfList, asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2002::/64",
		IfIndex: 3,
	})
	ipv6IntfList = append(ipv6IntfList, asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2003::1/64",
		IfIndex: 40,
	})
	ipv6IntfList = append(ipv6IntfList, asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2003::1/64",
		IfIndex: 5,
	})
}
func TestInitRtEventHdlrTestServer(t *testing.T) {
	fmt.Println("****Init Route event handler Server****")
	StartTestServer()
	fmt.Println("****************")
}
func TestProcessLogicalIntfCreateEvent(t *testing.T) {
	fmt.Println("**** Test LogicalIntfCreate event ****")
	fmt.Println("IntfIdNameMap before:")
	fmt.Println(IntfIdNameMap)
	fmt.Println("IfNameToIfIndex before:")
	fmt.Println(IfNameToIfIndex)
	for _, lo := range logicalIntfList {
		server.ProcessLogicalIntfCreateEvent(lo)
	}
	fmt.Println("IntfIdNameMap after:")
	fmt.Println(IntfIdNameMap)
	fmt.Println("IfNameToIfIndex after:")
	fmt.Println(IfNameToIfIndex)
	fmt.Println("***************************************")
}
func TestVlanCreateEvent(t *testing.T) {
	fmt.Println("**** TestVlanCreateEvent event ****")
	fmt.Println("IntfIdNameMap before:")
	fmt.Println(IntfIdNameMap)
	fmt.Println("IfNameToIfIndex before:")
	fmt.Println(IfNameToIfIndex)
	for _, vlan := range vlanList {
		server.ProcessVlanCreateEvent(vlan)
	}
	fmt.Println("IntfIdNameMap after:")
	fmt.Println(IntfIdNameMap)
	fmt.Println("IfNameToIfIndex after:")
	fmt.Println(IfNameToIfIndex)
	fmt.Println("***************************************")
}
func TestIPv4IntfCreateEvent(t *testing.T) {
	fmt.Println("**** TestIPv4IntfCreateEvent event ****")
	for _, v4Intf := range ipv4IntfList {
		server.ProcessIPv4IntfCreateEvent(v4Intf)
	}
	TestGetRouteReachability(t)
	fmt.Println("***************************************")
}
func TestIPv6IntfCreateEvent(t *testing.T) {
	fmt.Println("**** TestIPv6IntfCreateEvent event ****")
	for _, v6Intf := range ipv6IntfList {
		server.ProcessIPv6IntfCreateEvent(v6Intf)
	}
	TestGetRouteReachability(t)
	fmt.Println("***************************************")
}
func TestProcessIPv4IntfStateChangeEvents(t *testing.T) {
	fmt.Println("****TestProcessIPv4IntfStateChangeEvents()****")
	//TestGetRouteReachability(t)
	for _, ipInfo := range ipv4AddrList {
		server.ProcessIPv4IntfDownEvent(ipInfo.ipAddr, -1)
	}
	TestGetRouteReachability(t)
	for _, ipInfo := range ipv4AddrList {
		server.ProcessIPv4IntfUpEvent(ipInfo.ipAddr, -1)
	}
	TestGetRouteReachability(t)
	fmt.Println("********************************************")
}
func TestProcessIPv6IntfStateChangeEvents(t *testing.T) {
	fmt.Println("****TestProcessIPv6IntfStateChangeEvents()****")
	//TestGetRouteReachability(t)
	for _, ipInfo := range ipv6AddrList {
		server.ProcessIPv6IntfDownEvent(ipInfo.ipAddr, -1)
	}
	TestGetRouteReachability(t)
	for _, ipInfo := range ipv6AddrList {
		server.ProcessIPv6IntfUpEvent(ipInfo.ipAddr, -1)
	}
	TestGetRouteReachability(t)
	fmt.Println("********************************************")
}
func TestIPv4IntfDeleteEvent(t *testing.T) {
	fmt.Println("**** TestIPv4IntfDeleteEvent event ****")
	v4Intf := asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "31.1.10.2/24",
		IfIndex: 3,
	}
	server.ProcessIPv4IntfDeleteEvent(v4Intf)
	v4Intf = asicdClntDefs.IPv4IntfNotifyMsg{
		IpAddr:  "61.1.10.2/24",
		IfIndex: 6,
	}
	server.ProcessIPv4IntfDeleteEvent(v4Intf)
	fmt.Println("***************************************")
}
func TestIPv6IntfDeleteEvent(t *testing.T) {
	fmt.Println("**** TestIPv6IntfDeleteEvent event ****")
	v6Intf := asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2002::1/64",
		IfIndex: 4,
	}
	server.ProcessIPv6IntfDeleteEvent(v6Intf)
	v6Intf = asicdClntDefs.IPv6IntfNotifyMsg{
		IpAddr:  "2006::1/64",
		IfIndex: 6,
	}
	server.ProcessIPv6IntfDeleteEvent(v6Intf)
	fmt.Println("***************************************")
}
