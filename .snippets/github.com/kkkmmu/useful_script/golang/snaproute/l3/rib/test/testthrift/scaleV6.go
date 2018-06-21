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

package routeThriftTest

import (
	"fmt"
	"ribd"
	"strconv"
	"time"
)

func v6Add(client *ribd.RIBDServicesClient, nextHopIpStr string, maxCount int64) (err error) {
	fmt.Println("v6Add")
	var count int64 = 0
	//var maxCount int = 30000
	intByt4 := 1
	intByt3 := 1
	intByt2 := 22
	intByt1 := 2222
	start := time.Now()
	var route ribd.IPv6Route
	routeCount, _ := client.GetTotalv6RouteCount()
	fmt.Println("Route count before scale test start:", routeCount)
	for {
		if intByt4 > 253 && intByt3 > 254 {
			intByt2++
			intByt3 = 1
			intByt4 = 1
		}
		if intByt4 > 254 {
			intByt4 = 1
			intByt3++
		} else {
			intByt4++
		}
		if intByt3 > 254 {
			intByt2 = 1
		} //else {
		//intByt2++
		//}

		route = ribd.IPv6Route{}
		byte1 := strconv.Itoa(intByt1)
		byte2 := strconv.Itoa(intByt2)
		byte3 := strconv.Itoa(intByt3)
		byte4 := strconv.Itoa(intByt4)
		rtNet := byte1 + ":" + byte2 + ":" + byte3 + ":" + byte4 + "::/64"
		route.DestinationNw = rtNet
		route.NextHop = make([]*ribd.NextHopInfo, 0)
		nh := ribd.NextHopInfo{
			NextHopIp: nextHopIpStr, //"2002::1234:5678:9abc:1234",
		}
		route.NextHop = append(route.NextHop, &nh)
		route.Protocol = "STATIC"
		rv := client.OnewayCreateIPv6Route(&route)
		if rv == nil {
			count++
		} else {
			fmt.Println("Call failed", rv, "count: ", count)
			elapsed := time.Since(start)
			fmt.Println(" ## Elapsed time is ", elapsed)
			return nil
		}
		if maxCount == count {
			fmt.Println("Done. Total calls executed", count)
			break
		}

	}
	//elapsed := time.Since(start)
	//	fmt.Println(" ## Elapsed time is ", elapsed)
	return nil
}
func v6Del(client *ribd.RIBDServicesClient, nextHopIpStr string) (err error) {
	fmt.Println("v6Del")
	count, _ := client.GetTotalv6RouteCount()
	startCount := count
	fmt.Println("Deleting ", count, " number of v6 routes")
	intByt4 := 1
	intByt3 := 1
	intByt2 := 22
	intByt1 := 2222
	start := time.Now()
	var route ribd.IPv6Route
	for {
		if intByt4 > 253 && intByt3 > 254 {
			intByt2++
			intByt3 = 1
			intByt4 = 1
		}
		if intByt4 > 254 {
			intByt4 = 1
			intByt3++
		} else {
			intByt4++
		}
		if intByt3 > 254 {
			intByt2 = 1
		} //else {
		//intByt2++
		//}

		route = ribd.IPv6Route{}
		byte1 := strconv.Itoa(intByt1)
		byte2 := strconv.Itoa(intByt2)
		byte3 := strconv.Itoa(intByt3)
		byte4 := strconv.Itoa(intByt4)
		rtNet := byte1 + ":" + byte2 + ":" + byte3 + ":" + byte4 + "::/64"
		route.DestinationNw = rtNet
		route.NextHop = make([]*ribd.NextHopInfo, 0)
		nh := ribd.NextHopInfo{
			NextHopIp: nextHopIpStr, //"2002::1234:5678:9abc:1234",
		}
		route.NextHop = append(route.NextHop, &nh)
		route.Protocol = "STATIC"
		rv := client.OnewayDeleteIPv6Route(&route)
		if rv == nil {
			count--
		} else {
			fmt.Println("Call failed", rv, "count: ", count)
			elapsed := time.Since(start)
			fmt.Println(" ## Elapsed time is ", elapsed)
			return nil
		}
		if count == 2 {
			fmt.Println("Done. Total calls executed", startCount-count)
			break
		}

	}
	//elapsed := time.Since(start)
	//	fmt.Println(" ## Elapsed time is ", elapsed)
	return nil
}

//func main() {
func ScaleV6Add(ribdClient *ribd.RIBDServicesClient, nextHopIpStr string, number int64) {
	v6Add(ribdClient, nextHopIpStr, number) //ribd.NewRIBDServicesClientFactory(transport, protocolFactory))
}

func ScaleV6Del(ribdClient *ribd.RIBDServicesClient, nextHopIpStr string) {
	v6Del(ribdClient, nextHopIpStr)
}
