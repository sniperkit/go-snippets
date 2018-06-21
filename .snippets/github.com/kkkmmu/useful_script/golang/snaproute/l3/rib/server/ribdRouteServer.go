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

// ribdRouteServer.go
package server

import (
	defs "l3/rib/ribdCommonDefs"
	"ribd"
	"strconv"
)

type RouteConfigInfo struct {
	OrigRoute *ribd.IPv4Route
	NewRoute  *ribd.IPv4Route
	Attrset   []bool
	Op        string //"add"/"del"/"update"
}
type TrackReachabilityInfo struct {
	IpAddr   string
	Protocol string
	Op       string
}
type NextHopInfoKey struct {
	nextHopIp string
}
type NextHopInfo struct {
	refCount int //number of routes using this as a next hop
}
type RouteCountInfo struct {
	totalcount int
	ecmpcount  int
}
type PerProtocolRouteInfo struct {
	v4routeMap map[string]RouteCountInfo
	v6routeMap map[string]RouteCountInfo
	totalcount RouteCountInfo
}

var ProtocolRouteMap map[string]PerProtocolRouteInfo  //map[string]int
var InterfaceRouteMap map[string]PerProtocolRouteInfo //map[string]int

func UpdateV4ProtocolRouteMap(protocol string, op string, value string, ecmp bool) {
	var info PerProtocolRouteInfo

	if ProtocolRouteMap == nil {
		if op == "del" {
			return
		}
		ProtocolRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	}
	info, ok := ProtocolRouteMap[protocol]
	if !ok || info.v4routeMap == nil {
		if op == "del" {
			return
		}
		if info.v4routeMap == nil {
			info.v4routeMap = make(map[string]RouteCountInfo)
		}
	}
	protocolroutemap := info.v4routeMap
	if protocolroutemap == nil {
		return
	}
	totalcount := info.totalcount
	count, ok := protocolroutemap[value]
	if !ok {
		if op == "del" {
			return
		}
	}
	if op == "add" {
		count.totalcount++
		totalcount.totalcount++
		if ecmp {
			if count.ecmpcount == 0 {
				//first time ecmp route is added, add an additional count
				count.ecmpcount++
				totalcount.ecmpcount++
			}
			count.ecmpcount++
			totalcount.ecmpcount++
		}
	} else if op == "del" {
		count.totalcount--
		totalcount.totalcount--
		if ecmp {
			count.ecmpcount--
			totalcount.ecmpcount--
			if count.totalcount <= 1 {
				count.ecmpcount--
				totalcount.ecmpcount--
			}
		}
	}
	protocolroutemap[value] = count
	info.v4routeMap = protocolroutemap
	info.totalcount = totalcount
	ProtocolRouteMap[protocol] = info
}
func UpdateV6ProtocolRouteMap(protocol string, op string, value string, ecmp bool) {
	var info PerProtocolRouteInfo

	if ProtocolRouteMap == nil {
		if op == "del" {
			return
		}
		ProtocolRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	}
	info, ok := ProtocolRouteMap[protocol]
	if !ok || info.v6routeMap == nil {
		if op == "del" {
			return
		}
		if info.v6routeMap == nil {
			info.v6routeMap = make(map[string]RouteCountInfo)
		}
	}
	protocolroutemap := info.v6routeMap
	if protocolroutemap == nil {
		return
	}
	totalcount := info.totalcount
	count, ok := protocolroutemap[value]
	if !ok {
		if op == "del" {
			return
		}
	}
	if op == "add" {
		count.totalcount++
		totalcount.totalcount++
		if ecmp {
			if count.ecmpcount == 0 {
				//first time ecmp route is added, add an additional count
				count.ecmpcount++
				totalcount.ecmpcount++
			}
			count.ecmpcount++
			totalcount.ecmpcount++
		}
	} else if op == "del" {
		count.totalcount--
		totalcount.totalcount--
		if ecmp {
			count.ecmpcount--
			totalcount.ecmpcount--
			if count.totalcount <= 1 {
				count.ecmpcount--
				totalcount.ecmpcount--
			}
		}
	}
	protocolroutemap[value] = count
	info.v6routeMap = protocolroutemap
	info.totalcount = totalcount
	ProtocolRouteMap[protocol] = info
}
func UpdateProtocolRouteMap(protocol string, op string, ipType defs.IPType, value string, ecmp bool) {
	//logger.Debug("UpdateProtocolRouteMap,protocol:", protocol, " iptype:", ipType)
	if ipType == defs.IPv4 {
		UpdateV4ProtocolRouteMap(protocol, op, value, ecmp)
	} else {
		UpdateV6ProtocolRouteMap(protocol, op, value, ecmp)
	}

}

func UpdateV4InterfaceRouteMap(intfref string, op string, value string, ecmp bool) {
	var info PerProtocolRouteInfo

	if InterfaceRouteMap == nil {
		if op == "del" {
			return
		}
		InterfaceRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	}
	info, ok := InterfaceRouteMap[intfref]
	if !ok || info.v4routeMap == nil {
		if op == "del" {
			return
		}
		if info.v4routeMap == nil {
			info.v4routeMap = make(map[string]RouteCountInfo)
		}
	}
	interfaceroutemap := info.v4routeMap
	if interfaceroutemap == nil {
		return
	}
	totalcount := info.totalcount
	count, ok := interfaceroutemap[value]
	if !ok {
		if op == "del" {
			return
		}
	}
	if op == "add" {
		count.totalcount++
		totalcount.totalcount++
		if ecmp {
			if count.ecmpcount == 0 {
				//first time ecmp route is added, add an additional count
				count.ecmpcount++
				totalcount.ecmpcount++
			}
			count.ecmpcount++
			totalcount.ecmpcount++
		}
	} else if op == "del" {
		count.totalcount--
		totalcount.totalcount--
		if ecmp {
			count.ecmpcount--
			totalcount.ecmpcount--
			if count.totalcount <= 1 {
				count.ecmpcount--
				totalcount.ecmpcount--
			}
		}
	}
	interfaceroutemap[value] = count
	info.v4routeMap = interfaceroutemap
	info.totalcount = totalcount
	InterfaceRouteMap[intfref] = info
}
func UpdateV6InterfaceRouteMap(intfref string, op string, value string, ecmp bool) {
	var info PerProtocolRouteInfo

	if InterfaceRouteMap == nil {
		if op == "del" {
			return
		}
		InterfaceRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	}
	info, ok := InterfaceRouteMap[intfref]
	if !ok || info.v6routeMap == nil {
		if op == "del" {
			return
		}
		if info.v6routeMap == nil {
			info.v6routeMap = make(map[string]RouteCountInfo)
		}
	}
	interfaceroutemap := info.v6routeMap
	if interfaceroutemap == nil {
		return
	}
	totalcount := info.totalcount
	count, ok := interfaceroutemap[value]
	if !ok {
		if op == "del" {
			return
		}
	}
	if op == "add" {
		count.totalcount++
		totalcount.totalcount++
		if ecmp {
			if count.ecmpcount == 0 {
				//first time ecmp route is added, add an additional count
				count.ecmpcount++
				totalcount.ecmpcount++
			}
			count.ecmpcount++
			totalcount.ecmpcount++
		}
	} else if op == "del" {
		count.totalcount--
		totalcount.totalcount--
		if ecmp {
			count.ecmpcount--
			totalcount.ecmpcount--
			if count.totalcount <= 1 {
				count.ecmpcount--
				totalcount.ecmpcount--
			}
		}
	}
	interfaceroutemap[value] = count
	info.v6routeMap = interfaceroutemap
	info.totalcount = totalcount
	InterfaceRouteMap[intfref] = info
}
func UpdateInterfaceRouteMap(intf int, op string, ipType defs.IPType, value string, ecmp bool) {
	intfref := strconv.Itoa(int(intf))
	intfEntry, ok := IntfIdNameMap[int32(intf)]
	if ok {
		//logger.Debug("Map foud for ifndex : ", routeInfoList[sel].nextHopIfIndex, "Name = ", intfEntry.name)
		intfref = intfEntry.name
	}
	if ipType == defs.IPv4 {
		UpdateV4InterfaceRouteMap(intfref, op, value, ecmp)
	} else {
		UpdateV6InterfaceRouteMap(intfref, op, value, ecmp)
	}

}
func (ribdServiceHandler *RIBDServer) StartRouteProcessServer() {
	logger.Info("Starting the routeserver loop")
	ProtocolRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	for {
		select {
		case routeConf := <-ribdServiceHandler.RouteConfCh:
			//logger.Debug(fmt.Sprintln("received message on RouteConfCh channel, op: ", routeConf.Op)
			if routeConf.Op == defs.Add {
				ribdServiceHandler.ProcessV4RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), FIBAndRIB, ribd.Int(len(destNetSlice)))
			} else if routeConf.Op == defs.AddFIBOnly {
				ribdServiceHandler.ProcessV4RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), FIBOnly, routeConf.AdditionalParams.(ribd.Int))
			} else if routeConf.Op == defs.AddBulk {
				ribdServiceHandler.ProcessBulkRouteCreateConfig(routeConf.OrigBulkRouteConfigObject) //.([]*ribd.IPv4Route))
			} else if routeConf.Op == defs.Del {
				ribdServiceHandler.ProcessV4RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), FIBAndRIB)
			} else if routeConf.Op == defs.DelFIBOnly {
				ribdServiceHandler.ProcessV4RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), FIBOnly)
			} else if routeConf.Op == defs.Update {
				if routeConf.PatchOp == nil || len(routeConf.PatchOp) == 0 {
					ribdServiceHandler.Processv4RouteUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), routeConf.NewConfigObject.(*ribd.IPv4Route), routeConf.AttrSet)
				} else {
					ribdServiceHandler.Processv4RoutePatchUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), routeConf.NewConfigObject.(*ribd.IPv4Route), routeConf.PatchOp)
				}
			} else if routeConf.Op == defs.Addv6 {
				//create ipv6 route
				ribdServiceHandler.ProcessV6RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), FIBAndRIB, ribd.Int(len(destNetSlice)))
			} else if routeConf.Op == defs.Addv6FIBOnly {
				//create ipv6 route
				ribdServiceHandler.ProcessV6RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), FIBOnly, routeConf.AdditionalParams.(ribd.Int))
			} else if routeConf.Op == defs.Delv6 {
				//delete ipv6 route
				ribdServiceHandler.ProcessV6RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), FIBAndRIB)
			} else if routeConf.Op == defs.Delv6FIBOnly {
				//delete ipv6 route
				ribdServiceHandler.ProcessV6RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), FIBOnly)
			} else if routeConf.Op == defs.Updatev6 {
				//update ipv6 route
				if routeConf.PatchOp == nil || len(routeConf.PatchOp) == 0 {
					ribdServiceHandler.Processv6RouteUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), routeConf.NewConfigObject.(*ribd.IPv6Route), routeConf.AttrSet)
				} else {
					ribdServiceHandler.Processv6RoutePatchUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), routeConf.NewConfigObject.(*ribd.IPv6Route), routeConf.PatchOp)
				}
			}
		}
	}
}
