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

// ribdUtils.go
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	defs "l3/rib/ribdCommonDefs"
	"net"
	"ribd"
	"ribdInt"
	"sort"
	"strconv"
	"strings"
	"utils/patriciaDB"
	"utils/policy"

	"github.com/op/go-nanomsg"
)

type RouteDistanceConfig struct {
	defaultDistance    int
	configuredDistance int
}
type AdminDistanceSlice []ribd.RouteDistanceState
type RedistributeRouteInfo struct {
	route ribdInt.Routes
}
type RedistributionPolicyInfo struct {
	policy     string
	policyStmt string
}
type PublisherMapInfo struct {
	pub_ipc    string
	pub_socket *nanomsg.PubSocket
}

var RedistributeRouteMap map[string][]RedistributeRouteInfo
var RedistributionPolicyMap map[string]RedistributionPolicyInfo
var TrackReachabilityMap map[string][]string //map[ipAddr][]protocols
var RouteProtocolTypeMapDB map[string]int
var ReverseRouteProtoTypeMapDB map[int]string
var ProtocolAdminDistanceMapDB map[string]RouteDistanceConfig
var ProtocolAdminDistanceSlice AdminDistanceSlice
var PublisherInfoMap map[string]PublisherMapInfo
var RIBD_PUB *nanomsg.PubSocket
var RIBD_POLICY_PUB *nanomsg.PubSocket

func InitPublisher(pub_str string) (pub *nanomsg.PubSocket) {
	logger.Info("Setting up %s", pub_str, "publisher")
	pub, err := nanomsg.NewPubSocket()
	if err != nil {
		logger.Err("Failed to open pub socket")
		return nil
	}
	ep, err := pub.Bind(pub_str)
	if err != nil {
		logger.Info("Failed to bind pub socket - ", ep)
		return nil
	}
	err = pub.SetSendBuffer(1024 * 1024)
	if err != nil {
		logger.Err("Failed to set send buffer size")
		return nil
	}
	return pub
}

func BuildPublisherMap() {
	RIBD_PUB = InitPublisher(defs.PUB_SOCKET_ADDR)
	RIBD_POLICY_PUB = InitPublisher(defs.PUB_SOCKET_POLICY_ADDR)
	for k, _ := range RouteProtocolTypeMapDB {
		logger.Info("Building publisher map for protocol ", k)
		if k == "CONNECTED" || k == "STATIC" {
			logger.Info("Publisher info for protocol ", k, " not required")
			continue
		}
		if k == "IBGP" || k == "EBGP" {
			continue
		}
		pub_ipc := "ipc:///tmp/ribd_" + strings.ToLower(k) + "d.ipc"
		logger.Info("pub_ipc:", pub_ipc)
		pub := InitPublisher(pub_ipc)
		PublisherInfoMap[k] = PublisherMapInfo{pub_ipc, pub}
	}
	PublisherInfoMap["EBGP"] = PublisherInfoMap["BGP"]
	PublisherInfoMap["IBGP"] = PublisherInfoMap["BGP"]
	PublisherInfoMap["BFD"] = PublisherMapInfo{defs.PUB_SOCKET_BFDD_ADDR, InitPublisher(defs.PUB_SOCKET_BFDD_ADDR)}
	PublisherInfoMap["VXLAN"] = PublisherMapInfo{defs.PUB_SOCKET_VXLAND_ADDR, InitPublisher(defs.PUB_SOCKET_VXLAND_ADDR)}
}
func BuildRouteProtocolTypeMapDB() {
	RouteProtocolTypeMapDB["CONNECTED"] = defs.CONNECTED
	RouteProtocolTypeMapDB["EBGP"] = defs.EBGP
	RouteProtocolTypeMapDB["IBGP"] = defs.IBGP
	RouteProtocolTypeMapDB["BGP"] = defs.BGP
	RouteProtocolTypeMapDB["OSPF"] = defs.OSPF
	RouteProtocolTypeMapDB["STATIC"] = defs.STATIC

	//reverse
	ReverseRouteProtoTypeMapDB[defs.CONNECTED] = "CONNECTED"
	ReverseRouteProtoTypeMapDB[defs.IBGP] = "IBGP"
	ReverseRouteProtoTypeMapDB[defs.EBGP] = "EBGP"
	ReverseRouteProtoTypeMapDB[defs.BGP] = "BGP"
	ReverseRouteProtoTypeMapDB[defs.STATIC] = "STATIC"
	ReverseRouteProtoTypeMapDB[defs.OSPF] = "OSPF"
}
func BuildProtocolAdminDistanceMapDB() {
	ProtocolAdminDistanceMapDB["CONNECTED"] = RouteDistanceConfig{defaultDistance: 0, configuredDistance: -1}
	ProtocolAdminDistanceMapDB["STATIC"] = RouteDistanceConfig{defaultDistance: 1, configuredDistance: -1}
	ProtocolAdminDistanceMapDB["EBGP"] = RouteDistanceConfig{defaultDistance: 20, configuredDistance: -1}
	ProtocolAdminDistanceMapDB["IBGP"] = RouteDistanceConfig{defaultDistance: 200, configuredDistance: -1}
	ProtocolAdminDistanceMapDB["OSPF"] = RouteDistanceConfig{defaultDistance: 110, configuredDistance: -1}
}
func (slice AdminDistanceSlice) Len() int {
	return len(slice)
}
func (slice AdminDistanceSlice) Less(i, j int) bool {
	return slice[i].Distance < slice[j].Distance
}
func (slice AdminDistanceSlice) Swap(i, j int) {
	slice[i].Protocol, slice[j].Protocol = slice[j].Protocol, slice[i].Protocol
	slice[i].Distance, slice[j].Distance = slice[j].Distance, slice[i].Distance
}
func BuildProtocolAdminDistanceSlice(force bool) {
	distance := 0
	protocol := ""
	if ProtocolAdminDistanceSlice != nil && force == false {
		//dont build it if it is already built
		return
	}
	ProtocolAdminDistanceSlice = nil
	ProtocolAdminDistanceSlice = make([]ribd.RouteDistanceState, 0)
	if ProtocolAdminDistanceMapDB == nil {
		ProtocolAdminDistanceMapDB = make(map[string]RouteDistanceConfig)
		BuildProtocolAdminDistanceMapDB()
	}
	for k, v := range ProtocolAdminDistanceMapDB {
		protocol = k
		distance = v.defaultDistance
		if v.configuredDistance != -1 {
			distance = v.configuredDistance
		}
		routeDistance := ribd.RouteDistanceState{Protocol: protocol, Distance: int32(distance)}
		ProtocolAdminDistanceSlice = append(ProtocolAdminDistanceSlice, routeDistance)
	}
	sort.Sort(ProtocolAdminDistanceSlice)
}
func (m RIBDServer) ConvertIntfStrToIfIndexStr(intfString string) (ifIndex string, err error) {
	if val, err := strconv.Atoi(intfString); err == nil {
		//Verify ifIndex is valid
		//logger.Info("IfIndex = ", val)
		_, ok := IntfIdNameMap[int32(val)]
		if !ok {
			logger.Err("Cannot create ip route on a unknown L3 interface")
			return ifIndex, errors.New("Cannot create ip route on a unknown L3 interface")
		}
		ifIndex = intfString
	} else {
		//Verify ifName is valid
		if _, ok := IfNameToIfIndex[intfString]; !ok {
			return ifIndex, errors.New("Invalid ifName value")
		}
		ifIndex = strconv.Itoa(int(IfNameToIfIndex[intfString]))
	}
	return ifIndex, nil
}

func arpResolveCalled(key NextHopInfoKey) bool {
	if RouteServiceHandler.NextHopInfoMap == nil {
		return false
	}
	info, ok := RouteServiceHandler.NextHopInfoMap[key]
	if !ok || info.refCount == 0 {
		logger.Info("Arp resolve not called for ", key.nextHopIp)
		return false
	}
	return true
}
func updateNextHopMap(key NextHopInfoKey, op int) (count int) {
	/*opStr := ""
	if op == add {
		opStr = "incrementing"
	} else if op == del {
		opStr = "decrementing"
	}
	logger.Info(opStr, " nextHop Map for ", key.nextHopIp)
	*/
	if RouteServiceHandler.NextHopInfoMap == nil {
		return -1
	}
	info, ok := RouteServiceHandler.NextHopInfoMap[key]
	if !ok {
		RouteServiceHandler.NextHopInfoMap[key] = NextHopInfo{1}
		count = 1
	} else {
		if op == add {
			info.refCount++
		} else if op == del {
			info.refCount--
		}
		RouteServiceHandler.NextHopInfoMap[key] = info
		count = info.refCount
	}
	//logger.Info("Updated refcount = ", count)
	return count
}
func findElement(list []string, element string) int {
	index := -1
	for i := 0; i < len(list); i++ {
		if list[i] == element {
			logger.Info("Found element ", element, " at index ", i)
			return i
		}
	}
	//logger.Info("Element ", element, " not added to the list")
	return index
}
func buildPolicyEntityFromRoute(route ribdInt.Routes, params interface{}) (entity policy.PolicyEngineFilterEntityParams, err error) {
	routeInfo := params.(RouteParams)
	//logger.Info("buildPolicyEntityFromRoute: createType: ", routeInfo.createType, " delete type: ", routeInfo.deleteType)
	destNetIp, err := getCIDR(route.Ipaddr, route.Mask)
	if err != nil {
		logger.Info("error getting CIDR address for ", route.Ipaddr, ":", route.Mask)
		return entity, err
	}
	entity.DestNetIp = destNetIp
	//logger.Info("buildPolicyEntityFromRoute: destNetIp:", entity.DestNetIp)
	entity.NextHopIp = route.NextHopIp
	entity.RouteProtocol = ReverseRouteProtoTypeMapDB[int(route.Prototype)]
	if routeInfo.createType != Invalid {
		entity.CreatePath = true
	}
	if routeInfo.deleteType != Invalid {
		entity.DeletePath = true
	}
	return entity, err
}
func BuildRouteParamsFromRouteInoRecord(routeInfoRecord RouteInfoRecord) RouteParams {
	var params RouteParams
	params.ipType = routeInfoRecord.ipType
	params.routeType = ribd.Int(routeInfoRecord.protocol)
	params.destNetIp = routeInfoRecord.destNetIp.String()
	params.sliceIdx = ribd.Int(routeInfoRecord.sliceIdx)
	params.networkMask = routeInfoRecord.networkMask.String()
	params.metric = routeInfoRecord.metric
	params.nextHopIp = routeInfoRecord.nextHopIp.String()
	params.nextHopIfIndex = routeInfoRecord.nextHopIfIndex
	return params
}
func BuildRouteParamsFromribdIPv4Route(cfg *ribd.IPv4Route, createType int, deleteType int, sliceIdx ribd.Int) RouteParams {
	nextHopIp := cfg.NextHop[0].NextHopIp
	if cfg.NullRoute == true { //commonDefs.IfTypeNull {
		logger.Info("null route create request")
		nextHopIp = "255.255.255.255"
	}
	nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[0].NextHopIntRef)
	params := RouteParams{destNetIp: cfg.DestinationNw,
		ipType:         defs.IPv4,
		networkMask:    cfg.NetworkMask,
		nextHopIp:      nextHopIp,
		nextHopIfIndex: ribd.Int(nextHopIntRef),
		weight:         ribd.Int(cfg.NextHop[0].Weight),
		metric:         ribd.Int(cfg.Cost),
		routeType:      ribd.Int(RouteProtocolTypeMapDB[cfg.Protocol]),
		sliceIdx:       ribd.Int(sliceIdx),
		createType:     ribd.Int(createType),
		deleteType:     ribd.Int(deleteType),
	}
	return params
}
func BuildRouteParamsFromribdIPv6Route(cfg *ribd.IPv6Route, createType int, deleteType int, sliceIdx ribd.Int) RouteParams {
	nextHopIp := cfg.NextHop[0].NextHopIp
	if cfg.NullRoute == true { //commonDefs.IfTypeNull {
		logger.Info("null route create request")
		nextHopIp = "255.255.255.255" //TBD: mask IP for null route next hop
	}
	nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[0].NextHopIntRef)
	params := RouteParams{destNetIp: cfg.DestinationNw,
		ipType:         defs.IPv6,
		networkMask:    cfg.NetworkMask,
		nextHopIp:      nextHopIp,
		nextHopIfIndex: ribd.Int(nextHopIntRef),
		weight:         ribd.Int(cfg.NextHop[0].Weight),
		metric:         ribd.Int(cfg.Cost),
		routeType:      ribd.Int(RouteProtocolTypeMapDB[cfg.Protocol]),
		sliceIdx:       ribd.Int(sliceIdx),
		createType:     ribd.Int(createType),
		deleteType:     ribd.Int(deleteType),
	}
	return params
}
func BuildPolicyRouteFromribdIPv4Route(cfg *ribd.IPv4Route) (policyRoute ribdInt.Routes) {
	nextHopIp := cfg.NextHop[0].NextHopIp
	if cfg.NullRoute == true { //commonDefs.IfTypeNull {
		logger.Info("null route create request")
		nextHopIp = "255.255.255.255"
	}
	nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[0].NextHopIntRef)
	policyRoute = ribdInt.Routes{Ipaddr: cfg.DestinationNw,
		IPAddrType: ribdInt.Int(defs.IPv4),
		Mask:       cfg.NetworkMask,
		NextHopIp:  nextHopIp,
		IfIndex:    ribdInt.Int(nextHopIntRef), //cfg.NextHopInfp[0].NextHopIntRef,
		Weight:     ribdInt.Int(cfg.NextHop[0].Weight),
		Metric:     ribdInt.Int(cfg.Cost),
		Prototype:  ribdInt.Int(RouteProtocolTypeMapDB[cfg.Protocol]),
	}
	return policyRoute
}
func BuildPolicyRouteFromribdIPv6Route(cfg *ribd.IPv6Route) (policyRoute ribdInt.Routes) {
	nextHopIp := cfg.NextHop[0].NextHopIp
	if cfg.NullRoute == true { //commonDefs.IfTypeNull {
		logger.Info("null route create request")
		//nextHopIp = "255.255.255.255"  //TBD: mask IP for null route next hop
	}
	nextHopIntRef, _ := strconv.Atoi(cfg.NextHop[0].NextHopIntRef)
	policyRoute = ribdInt.Routes{Ipaddr: cfg.DestinationNw,
		IPAddrType: ribdInt.Int(defs.IPv6),
		Mask:       cfg.NetworkMask,
		NextHopIp:  nextHopIp,
		IfIndex:    ribdInt.Int(nextHopIntRef), //cfg.NextHopInfp[0].NextHopIntRef,
		Weight:     ribdInt.Int(cfg.NextHop[0].Weight),
		Metric:     ribdInt.Int(cfg.Cost),
		Prototype:  ribdInt.Int(RouteProtocolTypeMapDB[cfg.Protocol]),
	}
	return policyRoute
}
func findRouteWithNextHop(routeInfoList []RouteInfoRecord, nextHopIpType defs.IPType, nextHopIP string, nextHopIfIndex ribd.Int) (found bool, routeInfoRecord RouteInfoRecord, index int) {
	//logger.Info("findRouteWithNextHop ", nextHopIP, " and type:", nextHopIpType, " and ifIndex:", nextHopIfIndex)
	index = -1
	for i := 0; i < len(routeInfoList); i++ {
		//	logger.Info("findRouteWithNextHop: current route type:", routeInfoList[i].nextHopIpType)
		if routeInfoList[i].nextHopIpType == nextHopIpType {
			//logger.Info("findRouteWithNextHop():same ip type,routeInfoList[i]:", routeInfoList[i])
			//logger.Info("Next hop IP present")
			if nextHopIP != "" && routeInfoList[i].nextHopIp.String() != nextHopIP {
				//logger.Info("findRouteWithNextHop(),nextHopIP ", nextHopIP, " not the same as route next hop ip:", routeInfoList[i].nextHopIp.String())
				continue
			}
			if nextHopIfIndex != -1 && routeInfoList[i].nextHopIfIndex != nextHopIfIndex {
				//logger.Info("nextHopIfIndex:", nextHopIfIndex, " routeInfoList[i].nextHopIfIndex:", routeInfoList[i].nextHopIfIndex, " do not match")
				//fmt.Println("nextHopIfIndex:", nextHopIfIndex, " routeInfoList[i].nextHopIfIndex:", routeInfoList[i].nextHopIfIndex, " do not match")
				continue
			}
			found = true
			routeInfoRecord = routeInfoList[i]
			index = i
			break
		}
	}
	return found, routeInfoRecord, index
}
func newNextHop(ipType defs.IPType, ip string, nextHopIfIndex ribd.Int, routeInfoList []RouteInfoRecord) (isNewNextHop bool) {
	//logger.Info("newNextHop")
	isNewNextHop = true
	for i := 0; i < len(routeInfoList); i++ {
		if routeInfoList[i].nextHopIp.String() == ip && routeInfoList[i].nextHopIfIndex == nextHopIfIndex && routeInfoList[i].nextHopIpType == ipType {
			//logger.Info("Next hop already present for nexthopIP:", ip, " ipType:", ipType, " ifIndex:", nextHopIfIndex)
			isNewNextHop = false
		}
	}
	return isNewNextHop
}
func isSameRoute(selectedRoute ribdInt.Routes, route ribdInt.Routes) (same bool) {
	//logger.Info("isSameRoute")
	if selectedRoute.IPAddrType == route.IPAddrType && selectedRoute.Ipaddr == route.Ipaddr && selectedRoute.Mask == route.Mask && selectedRoute.Prototype == route.Prototype {
		same = true
	}
	return same
}
func getPolicyRouteMapIndex(entity policy.PolicyEngineFilterEntityParams, policy string) (policyRouteIndex policy.PolicyEntityMapIndex) {
	//logger.Info("getPolicyRouteMapIndex")
	policyRouteIndex = PolicyRouteIndex{destNetIP: entity.DestNetIp, policy: policy}
	logger.Info("getPolicyRouteMapIndex:Returning policyRouteIndex as : ", policyRouteIndex)
	return policyRouteIndex
}

/*
   Update routelist for policy
*/
func addPolicyRouteMap(route ribdInt.Routes, policyName string) {
	//logger.Info("addPolicyRouteMap for route ", route, " policy:", policyName)
	ipPrefix, err := getNetowrkPrefixFromStrings(route.Ipaddr, route.Mask)
	if err != nil {
		logger.Err("Invalid ip prefix")
		return
	}
	maskIp, err := getIP(route.Mask)
	if err != nil {
		return
	}
	prefixLen, err := getPrefixLen(maskIp)
	if err != nil {
		return
	}
	var newRoute string
	found := false
	newRoute = route.Ipaddr + "/" + strconv.Itoa(prefixLen)
	//	newRoute := string(ipPrefix[:])
	logger.Info("addPolicyRouteMap for route ", route, " policy:", policyName, " Adding ip prefix ", newRoute, ipPrefix)
	policyInfo := PolicyEngineDB.PolicyDB.Get(patriciaDB.Prefix(policyName))
	if policyInfo == nil {
		logger.Info("Unexpected:policyInfo nil for policy ", policyName)
		return
	}
	tempPolicyInfo := policyInfo.(policy.Policy)
	tempPolicy := tempPolicyInfo.Extensions.(PolicyExtensions)
	tempPolicy.hitCounter++
	if tempPolicy.routeList == nil {
		logger.Info("routeList nil")
		tempPolicy.routeList = make([]string, 0)
	}
	//logger.Info("routelist len= ", len(tempPolicy.routeList))
	for i := 0; i < len(tempPolicy.routeList); i++ {
		//	logger.Info(" policy ", policyName, " routeList contains ", tempPolicy.routeList[i])
		if tempPolicy.routeList[i] == newRoute {
			//logger.Info(newRoute, " already is a part of ", policyName, "'s routelist")
			found = true
		}
	}
	if !found {
		tempPolicy.routeList = append(tempPolicy.routeList, newRoute)
	}
	found = false
	//logger.Info("routeInfoList details")
	for i := 0; i < len(tempPolicy.routeInfoList); i++ {
		//logger.Info("IP: ", tempPolicy.routeInfoList[i].Ipaddr, ":", tempPolicy.routeInfoList[i].Mask, " protocolType: ", ReverseRouteProtoTypeMapDB[int(tempPolicy.routeInfoList[i].Prototype)])
		if tempPolicy.routeInfoList[i].Ipaddr == route.Ipaddr && tempPolicy.routeInfoList[i].Mask == route.Mask && tempPolicy.routeInfoList[i].Prototype == route.Prototype {
			//		logger.Info("route already is a part of ", policyName, "'s routeInfolist")
			found = true
		}
	}
	if tempPolicy.routeInfoList == nil {
		tempPolicy.routeInfoList = make([]ribdInt.Routes, 0)
	}
	if found == false {
		tempPolicy.routeInfoList = append(tempPolicy.routeInfoList, route)
	}
	tempPolicyInfo.Extensions = tempPolicy
	PolicyEngineDB.PolicyDB.Set(patriciaDB.Prefix(policyName), tempPolicyInfo)
}
func deletePolicyRouteMap(route ribdInt.Routes, policyName string) {
	//logger.Info("deletePolicyRouteMap")
}
func updatePolicyRouteMap(route ribdInt.Routes, policy string, op int) {
	//logger.Info("updatePolicyRouteMap")
	if op == add {
		addPolicyRouteMap(route, policy)
	} else if op == del {
		deletePolicyRouteMap(route, policy)
	}

}

func deleteRoutePolicyStateAll(route ribdInt.Routes) {
	//logger.Info("deleteRoutePolicyStateAll")
	destNet, err := getNetowrkPrefixFromStrings(route.Ipaddr, route.Mask)
	if err != nil {
		return
	}

	routeInfoRecordListItem := RouteInfoMapGet(defs.IPType(route.IPAddrType), destNet)
	if routeInfoRecordListItem == nil {
		logger.Info(" entry not found for prefix %v", destNet)
		return
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	routeInfoRecordList.policyHitCounter = ribd.Int(route.PolicyHitCounter)
	routeInfoRecordList.policyList = nil //append(routeInfoRecordList.policyList[:0])
	RouteInfoMapSet(defs.IPType(route.IPAddrType), destNet, routeInfoRecordList)
	return
}
func addRoutePolicyState(route ribdInt.Routes, policy string, policyStmt string) {
	//logger.Info("addRoutePolicyState for ", route.Ipaddr, ":", route.Mask, " ipType:", route.IPAddrType)
	destNet, err := getNetowrkPrefixFromStrings(route.Ipaddr, route.Mask)
	if err != nil {
		return
	}

	routeInfoRecordListItem := RouteInfoMapGet(defs.IPType(route.IPAddrType), destNet)
	if routeInfoRecordListItem == nil {
		logger.Info("Unexpected - entry not found for prefix ", destNet)
		return
	}
	//logger.Info("Adding policy ", policy, " to route ", destNet)
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	found := false
	idx := 0
	for idx = 0; idx < len(routeInfoRecordList.policyList); idx++ {
		if routeInfoRecordList.policyList[idx] == policy {
			found = true
			break
		}
	}
	if found {
		logger.Info("Policy ", policy, "already a part of policyList of route ", destNet)
		return
	}
	routeInfoRecordList.policyHitCounter = ribd.Int(route.PolicyHitCounter)
	if routeInfoRecordList.policyList == nil {
		routeInfoRecordList.policyList = make([]string, 0)
	}
	/*	policyStmtList := routeInfoRecordList.policyList[policy]
		if policyStmtList == nil {
		   policyStmtList = make([]string,0)
		}
		policyStmtList = append(policyStmtList,policyStmt)
	    routeInfoRecordList.policyList[policy] = policyStmtList*/
	routeInfoRecordList.policyList = append(routeInfoRecordList.policyList, policy)
	RouteInfoMapSet(defs.IPType(route.IPAddrType), destNet, routeInfoRecordList)
	//logger.Debug("Adding to DBRouteCh from addRoutePolicyState")
	RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
		OrigConfigObject: RouteDBInfo{routeInfoRecordList.routeInfoProtocolMap[routeInfoRecordList.selectedRouteProtocol][0], routeInfoRecordList},
		Op:               defs.Add,
	}
	//RouteServiceHandler.WriteIPv4RouteStateEntryToDB(RouteDBInfo{routeInfoRecordList.routeInfoProtocolMap[routeInfoRecordList.selectedRouteProtocol][0], routeInfoRecordList})
	return
}
func deleteRoutePolicyState(ipType defs.IPType, ipPrefix patriciaDB.Prefix, policyName string) {
	//logger.Info("deleteRoutePolicyState")
	found := false
	idx := 0
	routeInfoRecordListItem := RouteInfoMapGet(ipType, ipPrefix)
	if routeInfoRecordListItem == nil {
		logger.Info("routeInfoRecordListItem nil for prefix ", ipPrefix)
		return
	}
	routeInfoRecordList := routeInfoRecordListItem.(RouteInfoRecordList)
	/*    if routeInfoRecordList.policyList[policyName] != nil {
		delete(routeInfoRecordList.policyList, policyName)
	}*/
	for idx = 0; idx < len(routeInfoRecordList.policyList); idx++ {
		if routeInfoRecordList.policyList[idx] == policyName {
			found = true
			break
		}
	}
	if !found {
		logger.Info("Policy ", policyName, "not found in policyList of route ", ipPrefix)
		return
	}
	if len(routeInfoRecordList.policyList) <= idx+1 {
		//logger.Info("last element, routeInfoRecordList.policyList:", routeInfoRecordList.policyList)
		routeInfoRecordList.policyList = routeInfoRecordList.policyList[:idx]
		//logger.Info("routeInfoRecordList.policyList after deleting:", routeInfoRecordList.policyList)
	} else {
		routeInfoRecordList.policyList = append(routeInfoRecordList.policyList[:idx], routeInfoRecordList.policyList[idx+1:]...)
	}
	RouteInfoMapSet(ipType, ipPrefix, routeInfoRecordList)
	//logger.Debug("Adding to DBRouteCh from deleteRoutePolicyState")
	RouteServiceHandler.DBRouteCh <- RIBdServerConfig{
		OrigConfigObject: RouteDBInfo{routeInfoRecordList.routeInfoProtocolMap[routeInfoRecordList.selectedRouteProtocol][0], routeInfoRecordList},
		Op:               defs.Add,
	}
}

func updateRoutePolicyState(route ribdInt.Routes, op int, policy string, policyStmt string) {
	//logger.Info("updateRoutePolicyState")
	if op == delAll {
		deleteRoutePolicyStateAll(route)
	} else if op == add {
		addRoutePolicyState(route, policy, policyStmt)
	}
}
func UpdateRedistributeTargetMap(evt int, protocol string, route ribdInt.Routes) {
	//logger.Info("UpdateRedistributeTargetMap")
	if evt == defs.NOTIFY_ROUTE_CREATED {
		redistributeMapInfo := RedistributeRouteMap[protocol]
		if redistributeMapInfo == nil {
			redistributeMapInfo = make([]RedistributeRouteInfo, 0)
		}
		redistributeRouteInfo := RedistributeRouteInfo{route: route}
		redistributeMapInfo = append(redistributeMapInfo, redistributeRouteInfo)
		RedistributeRouteMap[protocol] = redistributeMapInfo
	} else if evt == defs.NOTIFY_ROUTE_DELETED {
		redistributeMapInfo := RedistributeRouteMap[protocol]
		if redistributeMapInfo != nil {
			found := false
			i := 0
			for i = 0; i < len(redistributeMapInfo); i++ {
				if isSameRoute((redistributeMapInfo[i].route), route) {
					logger.Info("Found the route that is to be taken off the redistribution list for ", protocol)
					found = true
					break
				}
			}
			if found {
				if len(redistributeMapInfo) <= i+1 {
					redistributeMapInfo = redistributeMapInfo[:i]
				} else {
					redistributeMapInfo = append(redistributeMapInfo[:i], redistributeMapInfo[i+1:]...)
				}
			}
			RedistributeRouteMap[protocol] = redistributeMapInfo
		}
	}
}
func RedistributionNotificationSend(PUB *nanomsg.PubSocket, route ribdInt.Routes, evt int, targetProtocol string) {
	//logger.Info("RedistributionNotificationSend")
	msgBuf := defs.RoutelistInfo{RouteInfo: route}
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_ROUTE_CREATED {
		evtStr = " NOTIFY_ROUTE_CREATED "
	} else if evt == defs.NOTIFY_ROUTE_DELETED {
		evtStr = " NOTIFY_ROUTE_DELETED "
	}
	eventInfo := "Redistribute "
	if route.NetworkStatement == true {
		eventInfo = " Advertise Network Statement "
	}
	eventInfo = eventInfo + evtStr + " for route " + route.Ipaddr + " " + route.Mask + " type " + ReverseRouteProtoTypeMapDB[int(route.Prototype)] + " to " + targetProtocol
	//logger.Info("Adding ", evtStr, " for route ", route.Ipaddr, " ", route.Mask, " to notification channel")
	RouteServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}
func RouteReachabilityStatusNotificationSend(targetProtocol string, info RouteReachabilityStatusInfo) {
	logger.Info("RouteReachabilityStatusNotificationSend for protocol ", targetProtocol)
	publisherInfo, ok := PublisherInfoMap[targetProtocol]
	if !ok {
		logger.Info("Publisher not found for protocol ", targetProtocol)
		return
	}
	evt := defs.NOTIFY_ROUTE_REACHABILITY_STATUS_UPDATE
	PUB := publisherInfo.pub_socket
	msgInfo := defs.RouteReachabilityStatusMsgInfo{}
	msgInfo.Network = info.destNet
	if info.status == "Up" || info.status == "Updated" {
		msgInfo.IsReachable = true
	}
	msgInfo.NextHopIntf = info.nextHopIntf
	msgBuf := msgInfo
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	eventInfo := "Update Route Reachability status " + info.status + " for network " + info.destNet + " for protocol " + targetProtocol
	if info.status == "Up" {
		eventInfo = eventInfo + " NextHop IP: " + info.nextHopIntf.NextHopIp + " Index: " + strconv.Itoa(int(info.nextHopIntf.NextHopIfIndex))
	}
	//logger.Info("Adding  NOTIFY_ROUTE_REACHABILITY_STATUS_UPDATE with status ", info.status, " for network ", info.destNet, " to notification channel")
	RouteServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}
func RouteReachabilityStatusUpdate(targetProtocol string, info RouteReachabilityStatusInfo) {
	logger.Info("RouteReachabilityStatusUpdate targetProtocol ", targetProtocol, " info:", info)
	if targetProtocol != "NONE" {
		RouteReachabilityStatusNotificationSend(targetProtocol, info)
	}
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(info.destNet)
	if err != nil {
		logger.Err("Error getting IP from cidr: ", info.destNet)
		return
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	destIpPrefix, err := getNetowrkPrefixFromStrings(ipAddrStr, ipMaskStr)
	if err != nil {
		logger.Err("Error getting ip prefix for ip:", ipAddrStr, " mask:", ipMaskStr)
		return
	}
	logger.Info("RouteReachabilityStatusUpdate(), destIpPrefix:", destIpPrefix)
	//check the TrackReachabilityMap to see if any other protocols are interested in receiving updates for this network
	for k, list := range TrackReachabilityMap {
		prefix, err := getNetowrkPrefixFromStrings(k, ipMaskStr)
		if err != nil {
			logger.Err("Error getting ip prefix for ip:", k, " mask:", ipMaskStr)
			continue
		}
		logger.Info("RouteReachabilityStatusUpdate(), prefix:", prefix, " for k:", k, " list:", list)
		if bytes.Equal(destIpPrefix, prefix) {
			for idx := 0; idx < len(list); idx++ {
				logger.Info(" protocol ", list[idx], " interested in receving reachability updates for ipAddr ", info.destNet)
				info.destNet = k
				RouteReachabilityStatusNotificationSend(list[idx], info)
			}
		}
	}
	return
}

/*func getIPInt(ip net.IP) (ipInt int, err error) {
	if ip == nil {
		//logger.Info(fmt.Sprintln("ip address invalid", ip))
		return ipInt, errors.New("Invalid destination network IP Address")
	}
	ip = ip.To4()
	if ip == nil {
		//logger.Err("ip.To4 nil")
		return ipInt, errors.New("ip.To4() nil")
	}
	parsedPrefixIP := int(ip[3]) | int(ip[2])<<8 | int(ip[1])<<16 | int(ip[0])<<24
	ipInt = parsedPrefixIP
	return ipInt, nil
}
*/
func getIP(ipAddr string) (ip net.IP, err error) {
	//logger.Debug("getIP for ipAddr:", ipAddr)
	ip = net.ParseIP(ipAddr)
	if ip == nil {
		return ip, errors.New("Invalid destination network IP Address")
	}
	//ip = ip.To4()
	////logger.Debug(fmt.Sprintln("ip after ip.to4():", ip))
	return ip, nil
}
func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}
func isIPv4Mask(mask net.IP) bool {
	/*if len(mask) < 5 {
		return false
	}*/
	if isZeros(mask[0:10]) &&
		mask[10] == 0xff &&
		mask[11] == 0xff {
		return true
	}
	return false
}

func getPrefixLen(networkMask net.IP) (prefixLen int, err error) {
	//logger.Debug("getPrefixLen for networkMask: ", networkMask)
	/*	ipInt, err := getIPInt(networkMask)
		if err != nil {
			return -1, err
		}
		for prefixLen = 0; ipInt != 0; ipInt >>= 1 {
			prefixLen += ipInt & 1
		}*/
	mask := net.IPMask(networkMask)
	if isIPv4Mask(net.IP(mask)) {
		prefixLen, _ = mask[12:16].Size()
	} else {
		prefixLen, _ = mask.Size()
	}
	//	prefixLen,bits := mask.Size()
	//logger.Debug(fmt.Sprintln("prefixLen = ", prefixLen))
	return prefixLen, err
}
func validateNetworkPrefix(ipAddr string, mask string) (destNet patriciaDB.Prefix, err error) {
	//logger.Debug("validateNetworkPrefix for ip ", ipAddr, " mask: ", mask)
	destNetIp, err := getIP(ipAddr)
	if err != nil {
		logger.Err("destNetIpAddr ", ipAddr, " invalid")
		return destNet, err
	}
	networkMask, err := getIP(mask)
	if err != nil {
		logger.Err("networkMaskAddr invalid")
		return destNet, err
	}
	prefixLen, err := getPrefixLen(networkMask)
	if err != nil {
		logger.Err("err when getting prefixLen, err= ", err)
		return destNet, errors.New(fmt.Sprintln("Invalid networkmask ", networkMask))
	}
	vdestMask := net.IPMask(networkMask) //net.IPv4Mask(networkMask[0], networkMask[1], networkMask[2], networkMask[3])
	netIp := destNetIp.Mask(vdestMask)
	if netIp == nil {
		logger.Err("netIp nil for ipAddr:", ipAddr, " mask:", mask)
		return destNet, errors.New("netIp nil")
	}
	//logger.Debug("netIP: ", netIp, " destNetIp ", destNetIp)
	if !(bytes.Equal(destNetIp, netIp)) {
		logger.Err("Cannot have ip : ", destNetIp, " more specific than mask ")
		return destNet, errors.New(fmt.Sprintln("IP address ", destNetIp, " more specific than mask ", networkMask))
	}
	numbytes := prefixLen / 8
	if (prefixLen % 8) != 0 {
		numbytes++
	}
	destNet = make([]byte, numbytes)
	for i := 0; i < numbytes && i < len(netIp); i++ {
		destNet[i] = netIp[i]
	}
	return destNet, err
}
func getNetworkPrefix(destNetIp net.IP, networkMask net.IP) (destNet patriciaDB.Prefix, nwAddr string, err error) {
	//logger.Debug("getNetworkPrefix for ip: ", destNetIp, "  networkMask: ", networkMask)
	prefixLen, err := getPrefixLen(networkMask)
	if err != nil {
		logger.Err("getNetworkPrefix for ip: ", destNetIp, "  networkMask: ", networkMask, " err when getting prefixLen, err= ", err)
		return destNet, nwAddr, errors.New(fmt.Sprintln("Invalid networkmask ", networkMask))
	}
	numbytes := prefixLen / 8
	if (prefixLen % 8) != 0 {
		numbytes++
	}
	var netIp net.IP
	vdestMask := net.IPMask(networkMask)
	//logger.Debug("vdestMask:", vdestMask)
	if isIPv4Mask(net.IP(vdestMask)) {
		netIp = destNetIp.Mask(vdestMask[12:16])
		//	logger.Debug("getNetworkPrefix for ip: ", destNetIp, "  networkMask: ", networkMask, " ipv4 case, netIp = ", netIp, " vdestMask:", vdestMask[12:16], " nwAddr:", nwAddr)
	} else {
		netIp = destNetIp.Mask(vdestMask)
		//	logger.Debug("getNetworkPrefix for ip: ", destNetIp, "  networkMask: ", networkMask, " ipv6 case, netIp = ", netIp, " vdestMask:", vdestMask, " nwAddr:", nwAddr)
	}
	if netIp == nil {
		//	logger.Debug("getNetworkPrefix for ip ", destNetIp, " networkMask: ", networkMask, " netIP nil : prefixLen  = ", prefixLen, "  numbytes:", numbytes)
		return destNet, nwAddr, errors.New("netIp nil")
	}
	destNet = make([]byte, numbytes)
	for i := 0; i < numbytes && i < len(netIp); i++ {
		destNet[i] = netIp[i]
		//	logger.Debug("destnet[", i, "]:", destNet[i], " netIp[", i, "]:", netIp[i])
	}
	nwAddr = (destNetIp.Mask(net.IPMask(networkMask))).String() + "/" + strconv.Itoa(prefixLen)
	//logger.Debug("destNet:", destNet, "nwAddr:", nwAddr)
	return destNet, nwAddr, err
}
func getNetowrkPrefixFromStrings(ipAddr string, mask string) (prefix patriciaDB.Prefix, err error) {
	//logger.Debug("getNetowrkPrefixFromStrings for ip ", ipAddr, " mask: ", mask)
	destNetIpAddr, err := getIP(ipAddr)
	if err != nil {
		logger.Info("getNetowrkPrefixFromStrings for ip ", ipAddr, " mask: ", mask, " destNetIpAddr ", ipAddr, " invalid")
		return prefix, err
	}
	//logger.Debug("getNetworkPrefixFrmStrings:destNetIpAddr:", destNetIpAddr)
	networkMaskAddr, err := getIP(mask)
	if err != nil {
		logger.Err("getNetowrkPrefixFromStrings for ip ", ipAddr, " mask: ", mask, " networkMaskAddr invalid")
		return prefix, err
	}
	prefix, _, err = getNetworkPrefix(destNetIpAddr, networkMaskAddr)
	if err != nil {
		logger.Info("getNetowrkPrefixFromStrings for ip ", ipAddr, " mask: ", mask, " err=", err)
		return prefix, err
	}
	return prefix, err
}
func getNetworkPrefixFromCIDR(ipAddr string) (ipPrefix patriciaDB.Prefix, err error) {
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return ipPrefix, err
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	//ipMaskStr := net.IP(ipMask).String()
	//logger.Debug("getNetowrkPrefixFromCIDR for ip ", ipAddr, " calling getNetowrkPrefixFromStrings(", ipAddrStr, ",", (net.IP(ipNet.Mask)).String(), ")")
	ipPrefix, err = getNetowrkPrefixFromStrings(ipAddrStr, (net.IP(ipNet.Mask)).String()) //ipMaskStr)
	return ipPrefix, err
}
func getCIDR(ipAddr string, mask string) (addr string, err error) {
	destNetIpAddr, err := getIP(ipAddr)
	if err != nil {
		logger.Err("destNetIpAddr invalid")
		return addr, err
	}
	maskIP, err := getIP(mask)
	if err != nil {
		logger.Err("err in getting mask IP for mask string", mask)
		return addr, err
	}
	prefixLen, err := getPrefixLen(maskIP)
	if err != nil {
		logger.Err("err in getting prefix len for mask string", mask)
		return addr, err
	}
	addr = (destNetIpAddr.Mask(net.IPMask(maskIP))).String() + "/" + strconv.Itoa(prefixLen)
	if isIPv4Mask(maskIP) {
		addr = (destNetIpAddr.Mask(net.IPMask(maskIP[12:16]))).String() + "/" + strconv.Itoa(prefixLen)
		//logger.Debug("ipv4 case, addr = ", addr, " maskIP:", maskIP[12:16])
	}
	return addr, err
}
