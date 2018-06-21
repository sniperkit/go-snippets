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
	"ribd"
	"testing"
	//	"time"
)

var ipv6AddrList []testIpInfo
var ipv6RouteList []*ribd.IPv6Route

func InitIpv6AddrInfoList() {
	ipv6AddrList = make([]testIpInfo, 0)
	ipv6AddrList = append(ipv6AddrList, testIpInfo{ipAddr: "2002::1.2.10.2", cidr: "2002::/64"})
	ipv6AddrList = append(ipv6AddrList, testIpInfo{ipAddr: "2002::10.1.2.10.2", cidr: "2002::/64"})
	ipv6AddrList = append(ipv6AddrList, testIpInfo{ipAddr: "2003::1.2.10.2", cidr: "2003::/64"})
	ipv6AddrList = append(ipv6AddrList, testIpInfo{ipAddr: "2004::1.2.10.2", cidr: "2003::/64"})
	ipv6AddrList = append(ipv6AddrList, testIpInfo{ipAddr: "20013::1.2.10.2", cidr: "20013::/64"})

}
func InitIpv6RouteList() {
	ipv6RouteList = make([]*ribd.IPv6Route, 0)
	ipv6RouteList = append(ipv6RouteList, &ribd.IPv6Route{
		DestinationNw: "2000:1::/64",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2002::1234:5678:9abc:1234"}},
		Protocol:      "EBGP",
	})
	ipv6RouteList = append(ipv6RouteList, &ribd.IPv6Route{
		DestinationNw: "2000:3::/64",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2002::1234:5678:9abc:1234"}},
		Protocol:      "STATIC",
	})
}
func TestInitv6RtProcessApiTestServer(t *testing.T) {
	fmt.Println("Initv6RtProcessApiTestServer")
	StartTestServer()
	TestProcessLogicalIntfCreateEvent(t)
	TestIPv6IntfCreateEvent(t)
	TestPolicyConditionConfigCreate(t)
	TestPolicyStmtConfigCreate(t)
	TestPolicyDefinitionConfigCreate(t)
	TestUpdateApplyPolicy(t)
	fmt.Println("****************")
}

func TestProcessV6RouteCreateConfig(t *testing.T) {
	fmt.Println("****TestProcessV6RouteCreateConfig****")
	for _, v6route := range ipv6RouteList {
		val_err := server.IPv6RouteConfigValidationCheck(ipv6RouteList[0], "add")
		if val_err != nil {
			fmt.Println("Validation failed for route:", ipv6RouteList[0], " with error:", val_err)
			continue
		}
		val, err := server.ProcessV6RouteCreateConfig(v6route, FIBAndRIB, ribd.Int(len(destNetSlice)))
		fmt.Println("val = ", val, " err: ", err, " for route:", v6route)
	}
	val, err := server.ProcessV6RouteCreateConfig(ipv6RouteList[0], FIBAndRIB, ribd.Int(len(destNetSlice)))
	fmt.Println("val = ", val, " err: ", err, " for route:", ipv6RouteList[0])
	TestGetRouteReachability(t)
	TestResolveNextHop(t)
	TestGetRoute(t)
	TestProcessIPv6IntfStateChangeEvents(t)
	fmt.Println("************************************")
}
func TestProcessv6RoutePatchUpdateConfig(t *testing.T) {
	fmt.Println("****TestProcessv6RoutePatchUpdateConfig****")
	for _, v6Route := range ipv6RouteList {
		for _, op := range patchOpList {
			//	fmt.Println("Applying patch:", op, " to route:", v6Route)
			testRoute := *v6Route
			val_err := server.IPv6RouteConfigValidationCheckForPatchUpdate(&testRoute, &testRoute, []*ribd.PatchOpInfo{op})
			if val_err != nil {
				fmt.Println("Validaion for Patch Update for route:", testRoute, "and patch op: ", op, " failed with err:", val_err)
				continue
			}
			val, err := server.Processv6RoutePatchUpdateConfig(&testRoute, &testRoute, []*ribd.PatchOpInfo{op})
			fmt.Println("val = ", val, " err: ", err, " for testRoute:", testRoute)
			TestGetRoute(t)
		}
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
func TestProcessv6RouteUpdateConfig(t *testing.T) {
	fmt.Println("****TestProcessv6RouteUpdateConfig****")
	for _, v6Route := range ipv6RouteList {
		var newRoute ribd.IPv6Route
		newRoute = *v6Route
		newRoute.Cost = 80
		newRoute.NextHop = []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2002::1234:5678:9abc:1234", Weight: 20}}
		attrSet := make([]bool, 6) //number of fields in ribd.IPv4Route
		attrSet[3] = true          //set cost attr to true
		attrSet[4] = true          //NUll route attr to true
		attrSet[5] = true          //set next hop ip attr to true
		val_err := server.IPv6RouteConfigValidationCheckForUpdate(v6Route, &newRoute, attrSet)
		if val_err != nil {
			fmt.Println("val_err:", val_err, " for v6Route:", v6Route, " newRoute:", newRoute, " attrSet:", attrSet)
			continue
		}
		val, err := server.Processv6RouteUpdateConfig(v6Route, &newRoute, attrSet)
		fmt.Println("val = ", val, " err: ", err)
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
func TestProcessv6RouteDeleteConfig(t *testing.T) {
	fmt.Println("****TestProcessv6RouteDeleteConfig****")
	for _, v6Route := range ipv6RouteList {
		val_err := server.IPv6RouteConfigValidationCheck(v6Route, "del")
		if val_err != nil {
			fmt.Println("Validation failed for route:", v6Route, " with error:", val_err)
			continue
		}
		val, err := server.ProcessV6RouteDeleteConfig(v6Route, FIBAndRIB)
		fmt.Println("val = ", val, " err: ", err)
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
