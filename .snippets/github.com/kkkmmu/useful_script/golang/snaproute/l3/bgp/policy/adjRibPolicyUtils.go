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

// adjRibPolicyUtils.go
package policy

import (
	bgprib "l3/bgp/rib"
	"l3/bgp/utils"
	"utils/patriciaDB"
	utilspolicy "utils/policy"
)

func (eng *AdjRibPPolicyEngine) AdjRIBDeleteRoutePolicyState(route *bgprib.AdjRIBPathIdRoute, policyName string) {
	utils.Logger.Info("deleteRoutePolicyState")
	found := false
	idx := 0
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

func deleteAdjRIBRoutePolicyStateAll(route *bgprib.AdjRIBPathIdRoute) {
	utils.Logger.Debug("deleteAdjRIBRoutePolicyStateAll")
	route.PolicyList = nil
	route.PolicyStmt = ""
	return
}

func addAdjRIBRoutePolicyState(route *bgprib.AdjRIBPathIdRoute, policy string, policyStmt string) {
	utils.Logger.Debug("addAdjRIBRoutePolicyState")
	route.PolicyList = append(route.PolicyList, policy)
	route.PolicyStmt = policyStmt
	return
}

func UpdateAdjRIBRoutePolicyState(route *bgprib.AdjRIBPathIdRoute, op int, policy string, policyStmt string) {
	utils.Logger.Debug("UpdateAdjRIBRoutePolicyState - op=%d", op)
	if op == DelAll {
		deleteAdjRIBRoutePolicyStateAll(route)
		//deletePolicyRouteMapEntry(route, policy)
	} else if op == Add {
		addAdjRIBRoutePolicyState(route, policy, policyStmt)
	}
}

func (eng *AdjRibPPolicyEngine) addAdjRIBPolicyRouteMap(route *bgprib.AdjRIBPathIdRoute, policy string) {
	utils.Logger.Debugf("addAdjRIBPolicyRouteMap - route=%+v, policy=%s", route, policy)
	var newRoute string
	newRoute = route.NLRI.GetCIDR()
	ipPrefix, err := GetNetworkPrefixFromCIDR(newRoute)
	if err != nil {
		utils.Logger.Info("Invalid ip prefix")
		return
	}
	utils.Logger.Info("Adding ip prefix %s %v ", newRoute, ipPrefix)
	policyInfo := eng.PolicyEngine.PolicyDB.Get(patriciaDB.Prefix(policy))
	if policyInfo == nil {
		utils.Logger.Info("Unexpected:policyInfo nil for policy ", policy)
		return
	}
	tempPolicy := policyInfo.(utilspolicy.Policy)
	policyExtensions := tempPolicy.Extensions.(AdjRibPolicyExtensions)
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
		utils.Logger.Info("IP: ", policyExtensions.RouteInfoList[i].NLRI.GetCIDR(), " neighbor: ",
			policyExtensions.RouteInfoList[i].Neighbor)
		if policyExtensions.RouteInfoList[i].NLRI.GetPrefix().String() == route.NLRI.GetPrefix().String() &&
			policyExtensions.RouteInfoList[i].NLRI.GetLength() == route.NLRI.GetLength() &&
			policyExtensions.RouteInfoList[i].Neighbor.String() == route.Neighbor.String() {
			utils.Logger.Info("route already is a part of ", policy, "'s routeInfolist")
			found = true
		}
	}
	if found == false {
		policyExtensions.RouteInfoList = append(policyExtensions.RouteInfoList, route)
	}
	eng.PolicyEngine.PolicyDB.Set(patriciaDB.Prefix(policy), tempPolicy)
}

func deleteAdjRIBPolicyRouteMap(route *bgprib.AdjRIBPathIdRoute, policy string) {
	//fmt.Println("deletePolicyRouteMap")
}

func (eng *AdjRibPPolicyEngine) UpdateAdjRIBPolicyRouteMap(route *bgprib.AdjRIBPathIdRoute, policy string, op int) {
	utils.Logger.Debugf("UpdateAdjRIBPolicyRouteMap - op=%d", op)
	if op == Add {
		eng.addAdjRIBPolicyRouteMap(route, policy)
	} else if op == Del {
		deleteAdjRIBPolicyRouteMap(route, policy)
	}

}
