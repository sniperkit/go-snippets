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

// policyUtils.go
package policy

import (
	"errors"
	"fmt"
	bgprib "l3/bgp/rib"
	"l3/bgp/utils"
	"net"
	"strconv"
	"strings"
	"utils/patriciaDB"
	utilspolicy "utils/policy"
)

const (
	Add = iota
	Del
	DelAll
	Invalidate
)

const (
	Invalid = -1
	Valid   = 0
)

type RouteParams struct {
	DestNetIp  string
	PrefixLen  uint16
	NextHopIp  string
	CreateType int
	DeleteType int
}

type PolicyRouteIndex struct {
	DestNetIP string //CIDR format
	Policy    string
}

type localDB struct {
	prefix     patriciaDB.Prefix
	isValid    bool
	precedence int
	nextHopIp  string
}

type ConditionsAndActionsList struct {
	conditionList []string
	actionList    []string
}
type PolicyStmtMap struct {
	policyStmtMap map[string]ConditionsAndActionsList
}

var PolicyRouteMap map[PolicyRouteIndex]PolicyStmtMap

func getPolicyEnityKey(entity utilspolicy.PolicyEngineFilterEntityParams, policy string) (
	policyEntityKey utilspolicy.PolicyEntityMapIndex) {
	utils.Logger.Info("getPolicyEnityKey entity =", entity, "policy =", policy)
	policyEntityKey = PolicyRouteIndex{DestNetIP: entity.DestNetIp, Policy: policy}
	return policyEntityKey
}

func getIPInt(ip net.IP) (ipInt int, err error) {
	if ip == nil {
		fmt.Printf("ip address %v invalid\n", ip)
		return ipInt, errors.New("Invalid destination network IP Address")
	}
	ip = ip.To4()
	parsedPrefixIP := int(ip[3]) | int(ip[2])<<8 | int(ip[1])<<16 | int(ip[0])<<24
	ipInt = parsedPrefixIP
	return ipInt, nil
}

func getIP(ipAddr string) (ip net.IP, err error) {
	ip = net.ParseIP(ipAddr)
	if ip == nil {
		return ip, errors.New("Invalid destination network IP Address")
	}
	ip = ip.To4()
	return ip, nil
}

func getPrefixLen(networkMask net.IP) (prefixLen int, err error) {
	ipInt, err := getIPInt(networkMask)
	if err != nil {
		return -1, err
	}
	for prefixLen = 0; ipInt != 0; ipInt >>= 1 {
		prefixLen += ipInt & 1
	}
	return prefixLen, nil
}
func getNetworkPrefix(destNetIp net.IP, networkMask net.IP) (destNet patriciaDB.Prefix, err error) {
	prefixLen, err := getPrefixLen(networkMask)
	if err != nil {
		utils.Logger.Info("err when getting prefixLen, err= ", err)
		return destNet, err
	}
	vdestMask := net.IPv4Mask(networkMask[0], networkMask[1], networkMask[2], networkMask[3])
	netIp := destNetIp.Mask(vdestMask)
	numbytes := prefixLen / 8
	if (prefixLen % 8) != 0 {
		numbytes++
	}
	destNet = make([]byte, numbytes)
	for i := 0; i < numbytes; i++ {
		destNet[i] = netIp[i]
	}
	return destNet, err
}
func getNetowrkPrefixFromStrings(ipAddr string, mask string) (prefix patriciaDB.Prefix, err error) {
	destNetIpAddr, err := getIP(ipAddr)
	if err != nil {
		utils.Logger.Info("destNetIpAddr invalid")
		return prefix, err
	}
	networkMaskAddr, err := getIP(mask)
	if err != nil {
		utils.Logger.Info("networkMaskAddr invalid")
		return prefix, err
	}
	prefix, err = getNetworkPrefix(destNetIpAddr, networkMaskAddr)
	if err != nil {
		utils.Logger.Info("err=", err)
		return prefix, err
	}
	return prefix, err
}

func GetNetworkPrefixFromCIDR(ipAddr string) (ipPrefix patriciaDB.Prefix, err error) {
	//var ipMask net.IP
	_, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return ipPrefix, err
	}
	/*
		ipMask = make(net.IP, 4)
		copy(ipMask, ipNet.Mask)
		ipAddrStr := ip.String()
		ipMaskStr := net.IP(ipMask).String()
		ipPrefix, err = getNetowrkPrefixFromStrings(ipAddrStr, ipMaskStr)
	*/
	i := strings.IndexByte(ipAddr, '/')
	prefixLen, _ := strconv.Atoi(ipAddr[i+1:])
	numbytes := (prefixLen + 7) / 8
	destNet := make([]byte, numbytes)
	for i := 0; i < numbytes; i++ {
		destNet[i] = ipNet.IP[i]
	}

	return patriciaDB.Prefix(destNet), err
}

func (eng *LocRibPolicyEngine) DeleteRoutePolicyState(route *bgprib.Route, policyName string) {
	utils.Logger.Info("deleteRoutePolicyState")
	found := false
	idx := 0
	/*    if routeInfoRecordList.policyList[policyName] != nil {
		delete(routeInfoRecordList.policyList, policyName)
	}*/
	for idx = 0; idx < len(route.PolicyList); idx++ {
		if route.PolicyList[idx] == policyName {
			found = true
			break
		}
	}

	if !found {
		utils.Logger.Info("Policy ", policyName, "not found in policyList of route", route)
		return
	}

	route.PolicyList = append(route.PolicyList[:idx], route.PolicyList[idx+1:]...)
}

func deleteRoutePolicyStateAll(route *bgprib.Route) {
	utils.Logger.Info("deleteRoutePolicyStateAll")
	route.PolicyList = nil
	return
}

func deletePolicyRouteMapEntry(route *bgprib.Route, policy string) {
	utils.Logger.Info("deletePolicyRouteMapEntry for policy ", policy, "route ", route.Dest.BGPRouteState.GetNetwork(), "/",
		route.Dest.BGPRouteState.GetCIDRLen())
	if PolicyRouteMap == nil {
		utils.Logger.Info("PolicyRouteMap empty")
		return
	}
	destNetIP := route.Dest.BGPRouteState.GetNetwork() + "/" + strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen()))
	policyRouteIndex := PolicyRouteIndex{DestNetIP: destNetIP, Policy: policy}
	//PolicyRouteMap[policyRouteIndex].policyStmtMap=nil
	delete(PolicyRouteMap, policyRouteIndex)
}

func addRoutePolicyState(route *bgprib.Route, policy string, policyStmt string) {
	utils.Logger.Info("addRoutePolicyState")
	route.PolicyList = append(route.PolicyList, policy)
	return
}

func UpdateRoutePolicyState(route *bgprib.Route, op int, policy string, policyStmt string) {
	utils.Logger.Info("updateRoutePolicyState")
	if op == DelAll {
		deleteRoutePolicyStateAll(route)
		//deletePolicyRouteMapEntry(route, policy)
	} else if op == Add {
		addRoutePolicyState(route, policy, policyStmt)
	}
}

func (eng *LocRibPolicyEngine) addPolicyRouteMap(route *bgprib.Route, policy string) {
	utils.Logger.Info("addPolicyRouteMap")
	//policy.hitCounter++
	//ipPrefix, err := getNetowrkPrefixFromStrings(route.Network, route.Mask)
	var newRoute string
	newRoute = route.Dest.BGPRouteState.GetNetwork() + "/" + strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen()))
	ipPrefix, err := GetNetworkPrefixFromCIDR(newRoute)
	if err != nil {
		utils.Logger.Info("Invalid ip prefix")
		return
	}
	//  newRoute := string(ipPrefix[:])
	utils.Logger.Info("Adding ip prefix %s %v ", newRoute, ipPrefix)
	policyInfo := eng.PolicyEngine.PolicyDB.Get(patriciaDB.Prefix(policy))
	if policyInfo == nil {
		utils.Logger.Info("Unexpected:policyInfo nil for policy ", policy)
		return
	}
	tempPolicy := policyInfo.(utilspolicy.Policy)
	policyExtensions := tempPolicy.Extensions.(PolicyExtensions)
	policyExtensions.HitCounter++

	utils.Logger.Info("routelist len= ", len(policyExtensions.RouteList), " prefix list so far")
	found := false
	for i := 0; i < len(policyExtensions.RouteList); i++ {
		utils.Logger.Info(policyExtensions.RouteList[i])
		if policyExtensions.RouteList[i] == newRoute {
			utils.Logger.Info(newRoute, " already is a part of ", policy, "'s routelist")
			found = true
		}
	}
	if !found {
		policyExtensions.RouteList = append(policyExtensions.RouteList, newRoute)
	}

	found = false
	utils.Logger.Info("routeInfoList details")
	for i := 0; i < len(policyExtensions.RouteInfoList); i++ {
		utils.Logger.Info("IP: ", policyExtensions.RouteInfoList[i].Dest.BGPRouteState.GetNetwork(), "/",
			policyExtensions.RouteInfoList[i].Dest.BGPRouteState.GetCIDRLen(), " nextHop: ",
			policyExtensions.RouteInfoList[i].PathInfo.NextHop)
		if policyExtensions.RouteInfoList[i].Dest.BGPRouteState.GetNetwork() == route.Dest.BGPRouteState.GetNetwork() &&
			policyExtensions.RouteInfoList[i].Dest.BGPRouteState.GetCIDRLen() == route.Dest.BGPRouteState.GetCIDRLen() &&
			policyExtensions.RouteInfoList[i].PathInfo.NextHop == route.PathInfo.NextHop {
			utils.Logger.Info("route already is a part of ", policy, "'s routeInfolist")
			found = true
		}
	}
	if found == false {
		policyExtensions.RouteInfoList = append(policyExtensions.RouteInfoList, route)
	}
	eng.PolicyEngine.PolicyDB.Set(patriciaDB.Prefix(policy), tempPolicy)
}

func deletePolicyRouteMap(route *bgprib.Route, policy string) {
	//fmt.Println("deletePolicyRouteMap")
}

func (eng *LocRibPolicyEngine) UpdatePolicyRouteMap(route *bgprib.Route, policy string, op int) {
	utils.Logger.Info("updatePolicyRouteMap")
	if op == Add {
		eng.addPolicyRouteMap(route, policy)
	} else if op == Del {
		deletePolicyRouteMap(route, policy)
	}

}
