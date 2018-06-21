package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"l3/rib/test/testthrift"
	"l3/rib/testutils"
	"os"
	"strconv"
	//"time"
)

func main() {
	ribdClient := testutils.GetRIBdClient()
	if ribdClient == nil {
		fmt.Println("RIBd client nil")
		return
	}

	routeThriftTest.Createv4RouteList()
	routeThriftTest.Createv6RouteList()
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to dial out to Redis server")
		return
	}
	defer conn.Close()

	route_ops := os.Args[1:]
	fmt.Println("op:", route_ops)
	for i := 0; i < len(route_ops); i++ {
		op := route_ops[i]
		switch op {
		case "createv4":
			fmt.Println("Create v4 route test")
			routeThriftTest.Createv4Routes(ribdClient)
		case "createv6":
			fmt.Println("Create v6 route test")
			routeThriftTest.Createv6Routes(ribdClient)
		case "verify":
			fmt.Println("Verify reachability info")
			routeThriftTest.CheckRouteReachability(ribdClient)
		case "deletev4":
			fmt.Println("Delete v4 route test")
			routeThriftTest.Deletev4Routes(ribdClient)
		case "deletev6":
			fmt.Println("Delete v6 route test")
			routeThriftTest.Deletev6Routes(ribdClient)
		case "scalev4Add":
			if (i + 2) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main scalev4Add <nexthopIpStr> <number>")
				break
			}
			nextHopIpStr := route_ops[i+1]
			number, _ := strconv.Atoi(route_ops[i+2])
			i++
			fmt.Println("Scale test for ", number, " v4 routes", " with nextHopIpStr:", nextHopIpStr)
			routeThriftTest.ScaleV4Add(ribdClient, nextHopIpStr, int64(number))
		case "scalev6Add":
			if (i + 2) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main scalev6Add <nextHopIpStr> <number>")
				break
			}
			nextHopIpStr := route_ops[i+1]
			number, _ := strconv.Atoi(route_ops[i+2])
			i++
			fmt.Println("Scale test for ", number, " v6 routes", " with nextHopIpStr:", nextHopIpStr)
			routeThriftTest.ScaleV6Add(ribdClient, nextHopIpStr, int64(number))
		case "scalev4Del":
			if (i + 1) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main scalev4Del <nextHopIpStr>")
				break
			}
			nextHopIpStr := route_ops[i+1]
			fmt.Println("Scale test for deleting v4 routes with nextHopIpStr", nextHopIpStr)
			routeThriftTest.ScaleV4Del(ribdClient, nextHopIpStr)
		case "scalev6Del":
			if (i + 1) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main scalev6Del <nextHopIpStr>")
				break
			}
			nextHopIpStr := route_ops[i+1]
			fmt.Println("Scale test for deleting v6 routes with nextHopIpStr", nextHopIpStr)
			routeThriftTest.ScaleV6Del(ribdClient, nextHopIpStr)
		case "ECMPScalev4Add":
			if (i + 2) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main ECMPScalev4Add <num of nextHops> <num of scale routes>")
				break
			}
			ecmpCount, _ := strconv.Atoi(route_ops[i+1])
			routeThriftTest.Wg.Add(ecmpCount)
			scaleCount, _ := strconv.Atoi(route_ops[i+2])
			fmt.Println("Scale test for adding ", ecmpCount, " -way ", scaleCount, " number of v4 routes")
			count := 1
			for {
				nextHopIpStr := "40.1." + strconv.Itoa(count) + ".2"
				fmt.Println("Adding ", scaleCount, " routes with nextHop:", nextHopIpStr)
				ribdClientNew := testutils.GetRIBdClient()
				if ribdClientNew == nil {
					fmt.Println("RIBd client New nil")
					return
				}

				routeThriftTest.EcmpScalev4Add(ribdClientNew, nextHopIpStr, int64(scaleCount))
				if count == ecmpCount {
					break
				}
				count++
			}
			routeThriftTest.Wg.Wait()
			//time.Sleep(time.Second * 60)
			//fmt.Println("After sleep")
		case "ECMPScalev4Del":
			if (i + 1) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main ECMPScalev4Del <num of nextHops>")
				break
			}
			ecmpCount, _ := strconv.Atoi(route_ops[i+1])
			routeThriftTest.Wg.Add(ecmpCount)
			fmt.Println("Scale test for deleting ", ecmpCount, " -way v4 routes")
			count := 1
			for {
				nextHopIpStr := "40.1." + strconv.Itoa(count) + ".2"
				fmt.Println("Deleting routes with nextHop:", nextHopIpStr)
				ribdClientNew := testutils.GetRIBdClient()
				if ribdClientNew == nil {
					fmt.Println("RIBd client New nil")
					return
				}

				routeThriftTest.EcmpScalev4Del(ribdClientNew, nextHopIpStr)
				if count == ecmpCount {
					break
				}
				count++
			}
			routeThriftTest.Wg.Wait()
			//time.Sleep(time.Second * 60)
			//fmt.Println("After sleep")
		case "RouteCount":
			fmt.Println("RouteCount")
			routeThriftTest.GetTotalRouteCount(ribdClient)
		case "GetV4RoutesFromDB":
			fmt.Println("GetV4RoutesFromDB")
			if (i + 2) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main GetV4RoutesFromDB <startIndex> <count>")
				break
			}
			start, _ := strconv.Atoi(route_ops[1])
			count, _ := strconv.Atoi(route_ops[2])
			routeThriftTest.GetBulkV4RoutesFromDb(ribdClient, int64(start), int64(count), conn)
		case "LoopTest":
			fmt.Println("LoopTest")
			routeThriftTest.LoopTest(ribdClient)
		case "Time":
			if (i + 1) == len(route_ops) {
				fmt.Println("Incorrect usage: should be ./main Time <number>")
				break
			}
			number, _ := strconv.Atoi(route_ops[i+1])
			i++
			fmt.Println("Time for ", number, " route creation")
			routeThriftTest.GetRouteCreatedTime(ribdClient, number)
		}
	}
}
