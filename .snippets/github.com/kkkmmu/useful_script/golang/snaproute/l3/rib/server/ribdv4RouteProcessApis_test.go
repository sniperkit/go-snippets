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
	"encoding/json"
	"fmt"
	"ribd"
	"ribdInt"
	"strconv"
	"testing"
	//	"time"
)

type testIpInfo struct {
	ipAddr string
	mask   string
	cidr   string
}

var ipv4AddrList []testIpInfo
var ipv4RouteList []*ribd.IPv4Route
var patchOpList []*ribd.PatchOpInfo

func InitIpv4AddrInfoList() {
	ipv4AddrList = make([]testIpInfo, 0)
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"11.1.10.2", "255.255.255.0", "11.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"21.1.10.2", "255.255.255.0", "21.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"12.1.10.2/24", "255.255.255.0", "12.1.10.2/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"12.1.10.20/24", "255.255.255.0", "12.1.10.20/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"13.1.10.2", "255.255.255.0", "13.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"22.1.10.2", "255.255.255.0", "22.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"33.1.10.2", "255.255.255.0", "33.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"40.0.1.2", "255.255.255.0", "40.0.1.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"40.1.10.2", "255.255.255.0", "40.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"50.1.10.2", "255.255.255.0", "50.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"41.1.10.2", "255.255.255.0", "41.1.10.0/24"})
	ipv4AddrList = append(ipv4AddrList, testIpInfo{"60.1.10.2", "255.255.255.0", "60.1.10.0/24"})

}
func InitIpv4RouteList() {
	ipv4RouteList = make([]*ribd.IPv4Route, 0)
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "40.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "11.1.10.2", NextHopIntRef: "lo1"}},
		Protocol:      "EBGP",
	})
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "41.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "11.1.10.2", NextHopIntRef: "lo1"}},
		Protocol:      "EBGP",
	})
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "50.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "12.1.10.20", NextHopIntRef: "lo2"}},
		Protocol:      "STATIC",
	})
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "40.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "22.1.10.2", NextHopIntRef: "22"}},
		Protocol:      "CONNECTED",
	})
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "60.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "12.1.10.20", NextHopIntRef: "lo2"}},
		Protocol:      "STATIC",
		Cost:          20,
	})
	ipv4RouteList = append(ipv4RouteList, &ribd.IPv4Route{
		DestinationNw: "40.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "33.1.10.2", NextHopIntRef: "33"}},
		Protocol:      "STATIC",
	})
}
func InitPatchOpList() {
	patchOpList = make([]*ribd.PatchOpInfo, 0)
	nhbytes, _ := json.Marshal([]*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "13.1.10.2", NextHopIntRef: "lo3"}})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "add",
		Path:  "NextHop",
		Value: string(nhbytes),
	})
	nh6bytes, _ := json.Marshal([]*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2003::1234:5678:9abc:1234", NextHopIntRef: "lo5"}})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "add",
		Path:  "NextHop",
		Value: string(nh6bytes),
	})
	costbytes, _ := json.Marshal(10)
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "add",
		Path:  "Cost",
		Value: string(costbytes),
	})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "remove",
		Path:  "NextHop",
		Value: string(nhbytes),
	})
	nhbytes1, _ := json.Marshal([]*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "14.1.10.2", NextHopIntRef: "lo4"}})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "remove",
		Path:  "NextHop",
		Value: string(nhbytes1),
	})
	nh6bytes, _ = json.Marshal([]*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2003::1234:5678:9abc:1234", NextHopIntRef: "lo5"}})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "remove",
		Path:  "NextHop",
		Value: string(nh6bytes),
	})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "remove",
		Path:  "Cost",
		Value: string(costbytes),
	})
	patchOpList = append(patchOpList, &ribd.PatchOpInfo{
		Op:    "test",
		Path:  "Cost",
		Value: string(costbytes),
	})
}
func TestInitv4RtProcessApiTestServer(t *testing.T) {
	fmt.Println("Initv4RtProcessApiTestServer")
	StartTestServer()
	TestProcessLogicalIntfCreateEvent(t)
	TestIPv4IntfCreateEvent(t)
	TestPolicyConditionConfigCreate(t)
	TestPolicyStmtConfigCreate(t)
	TestPolicyDefinitionConfigCreate(t)
	TestUpdateApplyPolicy(t)
	fmt.Println("****************")
}

func TestGetRouteReachability(t *testing.T) {
	//fmt.Println("**** Test GetRouteReachability****")
	for _, ipAddr := range ipv4AddrList {
		//fmt.Println("check route reachability of ipv4Addr:", ipAddr.ipAddr)
		nh, err := server.GetRouteReachabilityInfo(ipAddr.ipAddr, -1)
		if err != nil {
			fmt.Println("TestGetRouteReachability:error ", err, " getting route reachability for ip:", ipAddr, " nh:", nh)
			continue
		}
		//fmt.Println("TestGetRouteReachability: for ip:", ipAddr.ipAddr, ": nh:", nh)
	}
	for _, ipAddr := range ipv6AddrList {
		nh, err := server.GetRouteReachabilityInfo(ipAddr.ipAddr, -1)
		if err != nil {
			fmt.Println("TestGetRouteReachability:error ", err, " getting route reachability for ip:", ipAddr, " nh:", nh)
			continue
		}
		//fmt.Println("TestGetRouteReachability: for ip:", ipAddr.ipAddr, ": nh:", nh)
	}
	//fmt.Println("*************************************")
}
func TestResolveNextHop(t *testing.T) {
	//fmt.Println("****TestResolveNextHop****")
	for _, ipAddr := range ipv4AddrList {
		nh, rnh, err := ResolveNextHop(ipAddr.ipAddr)
		fmt.Println("TestResolveNextHop:nh:", nh, " rnh:", rnh, " err:", err, " for ipAddr:", ipAddr.ipAddr)
	}
	for _, ipAddr := range ipv6AddrList {
		nh, rnh, err := ResolveNextHop(ipAddr.ipAddr)
		fmt.Println("TestResolveNextHop;nh:", nh, " rnh:", rnh, " err:", err, " for ipAddr:", ipAddr.ipAddr)
	}
	//fmt.Println("****************************")
}
func TestGetRoute(t *testing.T) {
	//fmt.Println("**** TestGetRoute****")
	for _, ipInfo := range ipv4AddrList {
		rt, err := server.Getv4Route(ipInfo.cidr)
		if err != nil {
			fmt.Println("TestGetRoute:error getting ip info for ip:", ipInfo.cidr, " err:", err, " routeInfo:", rt)
			continue
		}
		//fmt.Println("TestGetRoute:rt info:", rt)
	}
	for _, ipInfo := range ipv6AddrList {
		rt, err := server.Getv6Route(ipInfo.cidr)
		if err != nil {
			fmt.Println("TestGetRoute;error getting ip info for ip:", ipInfo.cidr, "err:", err, " routeInfo:", rt)
			continue
		}
		//fmt.Println("TestGetRoute;rt info:", rt)
	}
	fmt.Println("TestGetRoute:Routes per protocol**************")
	stats, _ := server.GetBulkRouteStatsPerProtocolState(0, 10)
	fmt.Println(stats)
	//fmt.Println("*********************************")
}
func TestProcessV4RouteCreateConfig(t *testing.T) {
	fmt.Println("****TestProcessRouteCreateConfig****")
	for _, v4route := range ipv4RouteList {
		val_err := server.RouteConfigValidationCheck(ipv4RouteList[0], "add")
		if val_err != nil {
			fmt.Println("Validation failed for route:", ipv4RouteList[0], " with error:", val_err)
			continue
		}
		val, err := server.ProcessV4RouteCreateConfig(v4route, FIBAndRIB, ribd.Int(len(destNetSlice)))
		fmt.Println("val = ", val, " err: ", err, " for route:", v4route)
	}
	val, err := server.ProcessV4RouteCreateConfig(ipv4RouteList[0], FIBAndRIB, ribd.Int(len(destNetSlice)))
	fmt.Println("val = ", val, " err: ", err, " for route:", ipv4RouteList[0])
	TestGetRouteReachability(t)
	TestResolveNextHop(t)
	TestGetRoute(t)
	TestProcessIPv4IntfStateChangeEvents(t)
	fmt.Println("************************************")
}
func TestScaleRouteCreate(t *testing.T) {
	fmt.Println("****TestScaleRouteCreate****")
	//timeFmt := "2006-01-02 15:04:05.999999999 -0700 PDT"
	var count int = 0
	var maxCount int = 30000
	intByt2 := 1
	intByt3 := 1
	byte1 := "22"
	byte4 := "0"
	var routes []*ribdInt.IPv4RouteConfig
	var route []ribdInt.IPv4RouteConfig
	//var scaleTestStartTime string
	//var scaleTestEndTime string
	//var err error
	//var startTime time.Time
	//var endTime time.Time
	routeCount, _ := server.GetTotalv4RouteCount()
	fmt.Println("Route count before scale test start:", routeCount)
	routes = make([]*ribdInt.IPv4RouteConfig, 0)
	route = make([]ribdInt.IPv4RouteConfig, maxCount)
	for {
		if intByt3 > 254 {
			intByt3 = 1
			intByt2++
		} else {
			intByt3++
		}
		if intByt2 > 254 {
			intByt2 = 1
		} //else {
		//intByt2++
		//}

		byte2 := strconv.Itoa(intByt2)
		byte3 := strconv.Itoa(intByt3)
		rtNet := byte1 + "." + byte2 + "." + byte3 + "." + byte4
		route[count].DestinationNw = rtNet
		route[count].NetworkMask = "255.255.255.0"
		route[count].NextHop = make([]*ribdInt.RouteNextHopInfo, 0)
		nh := ribdInt.RouteNextHopInfo{
			NextHopIp: "11.1.10.2",
		}
		route[count].NextHop = append(route[count].NextHop, &nh)
		route[count].Protocol = "STATIC"
		routes = append(routes, &route[count])
		count++
		if maxCount == count {
			fmt.Println("Done. Total route configs added ", count)
			break
		}
		//fmt.Println("Creating Route ", route)
		/*		_, err := server.ProcessRouteCreateConfig(&route)
				if err == nil {
					if count == 0 {
						fmt.Println("recording starttime as ", routeCreatedTime)
						scaleTestStartTime = routeCreatedTime
					}
					count++
				} else {
					fmt.Println("Call failed", err, "count: ", count)
					return
				}
				if maxCount == count {
					fmt.Println("Done. Total calls executed", count)
					fmt.Println("recording endtime as ", routeCreatedTime)
					scaleTestEndTime = routeCreatedTime
					break
				}*/
	}
	/*	fmt.Println("startTime:", scaleTestStartTime)
		startTime, err := time.Parse(timeFmt, scaleTestStartTime)
		if err != nil {
			fmt.Println("err parsing obj time:", scaleTestStartTime, " into timeFmt:", timeFmt, " err:", err)
			return
		}
		fmt.Println("endTime:", scaleTestEndTime)
		endTime, err := time.Parse(timeFmt, scaleTestEndTime)
		if err != nil {
			fmt.Println("err parsing obj time:", scaleTestEndTime, " into timeFmt:", timeFmt, " err:", err)
			return
		}
		fmt.Println("Time to install ", maxCount, " number of routes is:", "duration:", endTime.Sub(startTime))
	*/
	/*	server.ProcessBulkRouteCreateConfig(routes)
		scaleTestStartTime, err = server.Getv4RouteCreatedTime(routeCount + 1)
		if err != nil {
			fmt.Println("err ", err, " getting routecreated time for route #", routeCount+1)
			return
		}
		fmt.Println("startTime:", scaleTestStartTime, " for the ", routeCount+1, " route")
		startTime, err = time.Parse(timeFmt, scaleTestStartTime)
		if err != nil {
			fmt.Println("err parsing obj time:", scaleTestStartTime, " into timeFmt:", timeFmt, " err:", err)
			return
		}
		scaleTestEndTime, err = server.Getv4RouteCreatedTime(routeCount + maxCount)
		if err != nil {
			fmt.Println("err ", err, " getting routecreated time for route #", routeCount+maxCount)
			for {
				scaleTestEndTime, err = server.Getv4RouteCreatedTime(routeCount + maxCount)
				if err == nil {
					break
				}
			}
			//return
		}
		fmt.Println("endTime:", scaleTestEndTime, " after the ", routeCount+maxCount, " route")
		endTime, err = time.Parse(timeFmt, scaleTestEndTime)
		if err != nil {
			fmt.Println("err parsing obj time:", scaleTestEndTime, " into timeFmt:", timeFmt, " err:", err)
			return
		}
		fmt.Println("Getv4RouteCreatedTime() method Time to install ", maxCount, " number of routes is:", "duration:", endTime.Sub(startTime))
	*/
}
func TestProcessv4RoutePatchUpdateConfig(t *testing.T) {
	fmt.Println("****TestProcessRoutePatchUpdateConfig****")
	for _, v4Route := range ipv4RouteList {
		for _, op := range patchOpList {
			//fmt.Println("Applying patch:", op, " to route:", v4Route)
			testRoute := *v4Route
			val_err := server.RouteConfigValidationCheckForPatchUpdate(&testRoute, &testRoute, []*ribd.PatchOpInfo{op})
			if val_err != nil {
				fmt.Println("Validaion for Patch Update for route:", testRoute, "and patch op: ", op, " failed with err:", val_err)
				continue
			}
			val, err := server.Processv4RoutePatchUpdateConfig(&testRoute, &testRoute, []*ribd.PatchOpInfo{op})
			fmt.Println("val = ", val, " err: ", err, " for testRoute:", testRoute)
			TestGetRoute(t)
		}
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
func TestProcessv4RouteUpdateConfig(t *testing.T) {
	fmt.Println("****TestProcessRouteUpdateConfig****")
	for _, v4Route := range ipv4RouteList {
		var newRoute ribd.IPv4Route
		newRoute = *v4Route
		newRoute.Cost = 80
		newRoute.NextHop = []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "12.1.10.20", Weight: 20}}
		attrSet := make([]bool, 6) //number of fields in ribd.IPv4Route
		attrSet[3] = true          //set cost attr to true
		attrSet[4] = true          //NUll route attr to true
		attrSet[5] = true          //set next hop ip attr to true
		val_err := server.RouteConfigValidationCheckForUpdate(v4Route, &newRoute, attrSet)
		if val_err != nil {
			fmt.Println("val_err:", val_err, " for v4Route:", v4Route, " newRoute:", newRoute, " attrSet:", attrSet)
			continue
		}
		val, err := server.Processv4RouteUpdateConfig(v4Route, &newRoute, attrSet)
		fmt.Println("val = ", val, " err: ", err)
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
func TestProcessv4RouteDeleteConfig(t *testing.T) {
	fmt.Println("****TestProcessRouteDeleteConfig****")
	for _, v4Route := range ipv4RouteList {
		val_err := server.RouteConfigValidationCheck(v4Route, "del")
		if val_err != nil {
			fmt.Println("Validation failed for route:", v4Route, " with error:", val_err)
			continue
		}
		val, err := server.ProcessV4RouteDeleteConfig(v4Route, FIBAndRIB)
		fmt.Println("val = ", val, " err: ", err, " for v4Route:", v4Route)
	}
	TestGetRouteReachability(t)
	TestGetRoute(t)
	TestResolveNextHop(t)
	fmt.Println("************************************")
}
