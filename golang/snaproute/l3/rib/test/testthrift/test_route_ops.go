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

// test_route_ops
package routeThriftTest

import (
	"fmt"
	"ribd"
	"ribdInt"
	"sync"
)

var route ribd.IPv4Route
var v6route ribd.IPv6Route
var ipv4RouteList []ribd.IPv4Route
var ipv6RouteList []ribd.IPv6Route
var reachabilityTestList []string
var Wg sync.WaitGroup

func GetTotalRouteCount(client *ribd.RIBDServicesClient) {
	fmt.Println("GetTotalRouteCount")
	number, _ := client.GetTotalv4RouteCount()
	fmt.Println("Number of v4 routes:", number)
	number, _ = client.GetTotalv6RouteCount()
	fmt.Println("Number of v6 routes:", number)
}
func GetRouteCreatedTime(client *ribd.RIBDServicesClient, number int) {
	fmt.Println("GetRouteCreatedTime")
	time, err := client.Getv4RouteCreatedTime(ribdInt.Int(number))
	fmt.Println("v4 route created time err: ", err, " time:", time)
	time, err = client.Getv6RouteCreatedTime(ribdInt.Int(number))
	fmt.Println("v6 route created time err: ", err, " time:", time)
}
func Createv4Routes(client *ribd.RIBDServicesClient) {
	fmt.Println("Createv4Routes")
	for _, route := range ipv4RouteList {
		fmt.Println("creating routes for :", route)
		client.CreateIPv4Route(&route)
	}
}
func Deletev4Routes(client *ribd.RIBDServicesClient) {
	fmt.Println("Deletev4Routes")
	for _, route := range ipv4RouteList {
		v, e := client.DeleteIPv4Route(&route)
		fmt.Println("v:", v, "e:", e, " for deleteIPRoute for ", route)
	}
}
func Createv6Routes(client *ribd.RIBDServicesClient) {
	fmt.Println("Createv6Routes")
	for _, v6route := range ipv6RouteList {
		fmt.Println("creating v6 routes for :", v6route)
		client.CreateIPv6Route(&v6route)
	}
}
func Deletev6Routes(client *ribd.RIBDServicesClient) {
	fmt.Println("Deletev6Routes")
	for _, v6route := range ipv6RouteList {
		v, e := client.DeleteIPv6Route(&v6route)
		fmt.Println("v:", v, "e:", e, " for deleteIPRoute for ", v6route)
	}
}
func CheckRouteReachability(client *ribd.RIBDServicesClient) {
	fmt.Println("CheckRouteReachability")
	for _, dest := range reachabilityTestList {
		nhIntf, err := client.GetRouteReachabilityInfo(dest, -1)
		fmt.Println("nhIntf info for ", dest, ":", nhIntf, " err:", err)
	}
}
func Createv4RouteList() {
	ipv4RouteList = make([]ribd.IPv4Route, 0)
	/*ipv4RouteList = append(ipv4RouteList, ribd.IPv4Route{
		DestinationNw: "40.1.10.0",
		NetworkMask:   "255.255.255.0",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.1.2"}},
		Protocol:      "STATIC",
	})*/
	ipv4RouteList = append(ipv4RouteList, ribd.IPv4Route{
		DestinationNw: "40.10.0.0/16",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.1.2"}},
		Protocol:      "EBGP",
		NullRoute:     true,
	})

	ipv4RouteList = append(ipv4RouteList, ribd.IPv4Route{
		DestinationNw: "50.10.0.0/16",
		Protocol:      "EBGP",
		NullRoute:     true,
	})
	ipv4RouteList = append(ipv4RouteList, ribd.IPv4Route{
		DestinationNw: "40.10.0.0/24",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.2.2"}},
		Protocol:      "STATIC",
	})
	/*	ipv4RouteList = append(ipv4RouteList, ribd.IPv4Route{
		DestinationNw: "40.1.1.1/24",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.2.2"}},
		Protocol:      "EBGP",
		NullRoute:     true,
	})*/

	//reachability test list
	reachabilityTestList = make([]string, 0)
	reachabilityTestList = append(reachabilityTestList, "40.0.1.2")
	reachabilityTestList = append(reachabilityTestList, "40.1.1.2")
	reachabilityTestList = append(reachabilityTestList, "40.1.10.2")
}
func Createv6RouteList() {
	ipv6RouteList = make([]ribd.IPv6Route, 0)
	ipv6RouteList = append(ipv6RouteList, ribd.IPv6Route{
		DestinationNw: "2001::/64",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2002::1234:5678:9abc:1234"}},
		Protocol:      "EBGP",
		NullRoute:     true,
	})

	ipv6RouteList = append(ipv6RouteList, ribd.IPv6Route{
		DestinationNw: "2003::/64",
		Protocol:      "EBGP",
		NullRoute:     true,
	})
	ipv6RouteList = append(ipv6RouteList, ribd.IPv6Route{
		DestinationNw: "2003::5/127",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "2002::1234:5678:9abc:1234"}},
		Protocol:      "STATIC",
	})

}
func LoopTest(client *ribd.RIBDServicesClient) {
	route := ribd.IPv4Route{
		DestinationNw: "40.10.0.0/16",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.1.2"}},
		Protocol:      "EBGP",
		NullRoute:     true,
	}
	route1 := ribd.IPv4Route{
		DestinationNw: "40.10.0.0/24",
		NextHop:       []*ribd.NextHopInfo{&ribd.NextHopInfo{NextHopIp: "40.1.2.2"}},
		Protocol:      "STATIC",
	}
	for {
		client.CreateIPv4Route(&route)
		client.DeleteIPv4Route(&route)
		client.CreateIPv4Route(&route1)
		client.DeleteIPv4Route(&route1)
	}
}
