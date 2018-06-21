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

// networkStmtPolicyUtils.go
package policy

import (
	"l3/bgp/config"
	"l3/bgp/utils"
	"utils/patriciaDB"
	utilspolicy "utils/policy"
)

func (eng *NetworkStmtPolicyEngine) DeleteNSPolicyState(ns *config.BGPNetworkStatement, policyName string) {
	utils.Logger.Info("NSDeleteRoutePolicyState")

	if ns.Policy != policyName {
		utils.Logger.Info("Policy ", policyName, "is not applied to ns", ns)
		return
	}

	ns.PolicyStmtList = nil
	ns.PolicyStmtList = make([]string, 0)
}

func deleteNSRoutePolicyStateAll(ns *config.BGPNetworkStatement) {
	utils.Logger.Info("deleteNSRoutePolicyStateAll")
	ns.PolicyList = nil
	ns.PolicyList = make([]string, 0)
	ns.PolicyStmtList = nil
	ns.PolicyStmtList = make([]string, 0)
}

func addNSRoutePolicyState(ns *config.BGPNetworkStatement, policy string, policyStmt string) {
	utils.Logger.Info("addNSRoutePolicyState")
	ns.PolicyList = append(ns.PolicyList, policy)
	ns.PolicyStmtList = append(ns.PolicyStmtList, policyStmt)
}

func UpdateNSPolicyState(ns *config.BGPNetworkStatement, op int, policy string, policyStmt string) {
	utils.Logger.Infof("UpdateNSRoutePolicyState - op=%d, policy=%s stmt=%s, ns=%+v", op, policy, policyStmt, ns)
	if op == DelAll {
		deleteNSRoutePolicyStateAll(ns)
		//deletePolicyRouteMapEntry(ns, policy)
	} else if op == Add {
		addNSRoutePolicyState(ns, policy, policyStmt)
	}
}

func (eng *NetworkStmtPolicyEngine) addPolicyNSMap(ns *config.BGPNetworkStatement, policy string) {
	utils.Logger.Infof("addPolicyNSMap - ns=%+v, policy=%s", ns, policy)
	var newPrefix string
	newPrefix = ns.IPPrefix
	utils.Logger.Info("Adding ip prefix", newPrefix)
	policyInfo := eng.PolicyEngine.PolicyDB.Get(patriciaDB.Prefix(policy))
	if policyInfo == nil {
		utils.Logger.Info("Unexpected:policyInfo nil for policy ", policy)
		return
	}
	tempPolicy := policyInfo.(utilspolicy.Policy)
	policyExtensions := tempPolicy.Extensions.(NetworkStmtPolicyExtensions)
	policyExtensions.HitCounter++

	utils.Logger.Info("routelist len= ", len(policyExtensions.PrefixList), " prefix list so far")
	found := false
	for i := 0; i < len(policyExtensions.PrefixList); i++ {
		utils.Logger.Info(policyExtensions.PrefixList[i])
		if policyExtensions.PrefixList[i] == newPrefix {
			utils.Logger.Info(newPrefix, " already is a part of ", policy, "'s routelist")
			found = true
		}
	}
	if !found {
		utils.Logger.Info("addPolicyNSMap: add prefix", newPrefix, "to PrefixList")
		policyExtensions.PrefixList = append(policyExtensions.PrefixList, newPrefix)
	}

	found = false
	utils.Logger.Info("routeInfoList details")
	for i := 0; i < len(policyExtensions.PrefixInfoList); i++ {
		utils.Logger.Info("IP: ", policyExtensions.PrefixInfoList[i].IPPrefix, " neighbor: ",
			policyExtensions.PrefixInfoList[i].AddressFamily)
		if policyExtensions.PrefixInfoList[i].IPPrefix == newPrefix &&
			policyExtensions.PrefixInfoList[i].AddressFamily == ns.AddressFamily {
			utils.Logger.Info("ns", ns, "already is a part of ", policy, "'s routeInfolist")
			found = true
		}
	}
	if !found {
		utils.Logger.Info("addPolicyNSMap: add network statement", *ns, "to PrefixInfoList")
		policyExtensions.PrefixInfoList = append(policyExtensions.PrefixInfoList, ns)
	}
	tempPolicy.Extensions = policyExtensions
	eng.PolicyEngine.PolicyDB.Set(patriciaDB.Prefix(policy), tempPolicy)
}

func (eng *NetworkStmtPolicyEngine) deletePolicyNSMap(ns *config.BGPNetworkStatement, policy string) {
	utils.Logger.Infof("deleteNSPolicyRouteMap - ns=%+v, policy=%s", ns, policy)
	var newPrefix string
	newPrefix = ns.IPPrefix
	utils.Logger.Info("Adding ip prefix %s", newPrefix)
	policyInfo := eng.PolicyEngine.PolicyDB.Get(patriciaDB.Prefix(policy))
	if policyInfo == nil {
		utils.Logger.Info("deleteNSPolicyRouteMap:policy %s not found", policy)
		return
	}
	tempPolicy := policyInfo.(utilspolicy.Policy)
	policyExtensions := tempPolicy.Extensions.(NetworkStmtPolicyExtensions)
	policyExtensions.HitCounter--

	utils.Logger.Info("routelist len= ", len(policyExtensions.PrefixList), " prefix list so far")
	found := false
	for i := 0; i < len(policyExtensions.PrefixList); i++ {
		utils.Logger.Info(policyExtensions.PrefixList[i])
		if policyExtensions.PrefixList[i] == newPrefix {
			utils.Logger.Info(newPrefix, " already is a part of ", policy, "'s routelist")
			found = true
		}
	}
	if !found {
		policyExtensions.PrefixList = append(policyExtensions.PrefixList, newPrefix)
	}

	found = false
	utils.Logger.Info("routeInfoList details")
	for i := 0; i < len(policyExtensions.PrefixInfoList); i++ {
		utils.Logger.Info("IP: ", policyExtensions.PrefixInfoList[i].IPPrefix, " neighbor: ",
			policyExtensions.PrefixInfoList[i].AddressFamily)
		if policyExtensions.PrefixInfoList[i].IPPrefix == newPrefix &&
			policyExtensions.PrefixInfoList[i].AddressFamily == ns.AddressFamily {
			utils.Logger.Info("ns", ns, "already is a part of ", policy, "'s routeInfolist")
			found = true
		}
	}
	if found == false {
		policyExtensions.PrefixInfoList = append(policyExtensions.PrefixInfoList, ns)
	}

	tempPolicy.Extensions = policyExtensions
	eng.PolicyEngine.PolicyDB.Set(patriciaDB.Prefix(policy), tempPolicy)
}

func (eng *NetworkStmtPolicyEngine) UpdatePolicyNSMap(ns *config.BGPNetworkStatement, policy string, op int) {
	utils.Logger.Infof("UpdateNSPolicyRouteMap - op=%d", op)
	if op == Add {
		eng.addPolicyNSMap(ns, policy)
	} else if op == Del {
		eng.deletePolicyNSMap(ns, policy)
	}

}
