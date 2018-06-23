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

// policyEngine.go
package policy

import (
	//"reflect"
	"sort"
	//	"strconv"
	//	"strings"
	bgpUtils "utils/bgpUtils"
	"utils/netUtils"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
	//	"utils/commonDefs"
	//"net"
	//	"asicdServices"
	//	"asicd/asicdConstDefs"
	//"bytes"
	//  "database/sql"
)

func (db *PolicyEngineDB) ActionListHasAction(actionList []PolicyAction, actionType int, action string) (match bool) {
	db.Logger.Info("ActionListHasAction for action ", action)
	return match
}

func (db *PolicyEngineDB) ActionNameListHasAction(actionList []string, actionType int, action string) (match bool) {
	db.Logger.Info("ActionListHasAction for action ", action)
	return match
}

func (db *PolicyEngineDB) PolicyEngineCheckActionsForEntity(entity PolicyEngineFilterEntityParams, policyConditionType int) (actionList []string) {
	db.Logger.Info("PolicyEngineTest to see if there are any policies for condition ", policyConditionType)
	var policyStmtList []string
	switch policyConditionType {
	case policyCommonDefs.PolicyConditionTypeProtocolMatch:
		policyStmtList = db.ProtocolPolicyListDB[entity.RouteProtocol]
		break
		/* Uncomment if needed in future
		case policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch:
			break
		case policyCommonDefs.PolicyConditionTypeNeighborMatch:
			policyStmtList = db.NeighborPolicyListDB[entity.Neighbor]
			break
		case policyCommonDefs.PolicyConditionTypeCommunityMatch:
			policyStmtList = db.CommunityPolicyListDB[entity.Community]
			break
		case policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch:
			policyStmtList = db.ExtCommunityPolicyListDB[entity.ExtendedCommunity]
			break
		case policyCommonDefs.PolicyConditionTypeASPathMatch:
			policyStmtList = db.ASPathPolicyListDB[entity.ASPath]
			break
		case policyCommonDefs.PolicyConditionTypeLocalPrefMatch:
			policyStmtList = db.LocalPrefPolicyListDB[entity.LocalPref]
			break*/
	default:
		db.Logger.Info("policyList not maintained for conditonType:", policyConditionType)
		return nil
	}
	if policyStmtList == nil || len(policyStmtList) == 0 {
		db.Logger.Info("no policy statements configured for this protocol")
		return nil
	}
	for i := 0; i < len(policyStmtList); i++ {
		db.Logger.Info("Found policy stmt ", policyStmtList[i], " for this entity")
		policyList := db.PolicyStmtPolicyMapDB[policyStmtList[i]]
		if policyList == nil || len(policyList) == 0 {
			db.Logger.Info("No policies configured for this entity")
			return nil
		}
		for j := 0; j < len(policyList); j++ {
			db.Logger.Info("Found policy ", policyList[j], "for this statement")
			policyStmtInfo := db.PolicyStmtDB.Get(patriciaDB.Prefix(policyStmtList[i]))
			if policyStmtInfo == nil {
				db.Logger.Info("Did not find this stmt in the DB")
				return nil
			}
			policyStmt := policyStmtInfo.(PolicyStmt)
			if db.ConditionCheckValid(entity, policyStmt.Conditions, policyStmt) {
				db.Logger.Info("All conditions valid for this route, so this policy will be potentially applied to this route")
				return policyStmt.Actions
			}
		}
	}
	return actionList
}

func (db *PolicyEngineDB) PolicyEngineUndoActionsPolicyStmt(policy Policy, policyStmt PolicyStmt, params interface{}, conditionsAndActionsList ConditionsAndActionsList) {
	db.Logger.Info("policyEngineUndoActionsPolicyStmt")
	if conditionsAndActionsList.ActionList == nil {
		db.Logger.Info("No actions")
		return
	}
	var i int
	conditionInfoList := make([]interface{}, 0)
	for j := 0; j < len(conditionsAndActionsList.ConditionList); j++ {
		conditionInfoList = append(conditionInfoList, conditionsAndActionsList.ConditionList[j].ConditionInfo)
	}

	for i = 0; i < len(conditionsAndActionsList.ActionList); i++ {
		db.Logger.Info("Find policy action number ", i, " name ", conditionsAndActionsList.ActionList[i], " in the action database")
		/*
			actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[i]))
			if actionItem == nil {
				db.Logger.Info("Did not find action ", conditionsAndActionsList.ActionList[i], " in the action database")
				continue
			}
			actionInfo := actionItem.(PolicyAction)
		*/
		policyAction := conditionsAndActionsList.ActionList[i]
		if db.UndoActionfuncMap[policyAction.ActionType] != nil {
			db.UndoActionfuncMap[policyAction.ActionType](policyAction.ActionInfo, conditionInfoList, policy, params, policyStmt)
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineUndoPolicyForEntity(entity PolicyEngineFilterEntityParams, policy Policy, policyEngineApplyInfo PolicyEngineApplyInfo, params interface{}) bool {
	db.Logger.Info("policyEngineUndoPolicyForRoute - policy name ", policy.Name, "  route: ", entity.DestNetIp, " type:", entity.RouteProtocol, " policyEngineApplyInfo:", policyEngineApplyInfo)
	if db.GetPolicyEntityMapIndex == nil {
		return false
	}
	policyEntityIndex := db.GetPolicyEntityMapIndex(entity, policy.Name)
	if policyEntityIndex == nil {
		db.Logger.Info("policy entity map index nil")
		return false
	}
	policyStmtMap := db.PolicyEntityMap[policyEntityIndex]
	if policyStmtMap.PolicyStmtMap == nil {
		db.Logger.Info("Unexpected:None of the policy statements of this policy have been applied on this route")
		return false
	}
	undoStmtMap := make(map[string]bool)
	for _, undoStmt := range policyEngineApplyInfo.StmtList {
		undoStmtMap[undoStmt] = true
	}
	undoConditionsMap := make(map[string]bool)
	for _, undoCondition := range policyEngineApplyInfo.ConditionsList {
		undoConditionsMap[undoCondition] = true
	}
	ret := true
	for stmt, conditionsAndActionsList := range policyStmtMap.PolicyStmtMap {
		db.Logger.Info("Applied policyStmtName ", stmt)
		//if the undo stmt list is non zero, then this is not the case for policy delete but for policy update
		if policyEngineApplyInfo.StmtList != nil && len(policyEngineApplyInfo.StmtList) > 0 {
			_, ok := undoStmtMap[stmt]
			if !ok {
				db.Logger.Info("this statement ", stmt, " is not the one to be removed from the policy")
				//return value should be false, so the policy is not deleted from the entity
				ret = false
				continue
			}
		}
		//if the undo conditions list is non zero, then this is not the case for policy delete but for policy update
		db.Logger.Debug("PolicyEngineUndoPolicyForEntity - policyStmtMap:", policyStmtMap)
		for _, condition := range conditionsAndActionsList.ConditionList {
			if policyEngineApplyInfo.ConditionsList != nil && len(policyEngineApplyInfo.ConditionsList) > 0 {
				db.Logger.Debug("PolicyEngineUndoPolicyForEntity:checking if condition ", condition.Name, " present in the undoConditionsMap:", undoConditionsMap)
				_, ok := undoConditionsMap[condition.Name]
				if !ok {
					db.Logger.Info("this condition ", condition.Name, " is not the one to be removed from the policy stmt")
					//return value should be false, so the policy is not deleted from the entity
					ret = false
					continue
				}
			}
		}
		if ret == false {
			continue
		}
		policyStmt := db.PolicyStmtDB.Get(patriciaDB.Prefix(stmt))
		if policyStmt == nil {
			db.Logger.Info("Invalid policyStmt")
			continue
		}
		db.PolicyEngineUndoActionsPolicyStmt(policy, policyStmt.(PolicyStmt), params, conditionsAndActionsList)
		//check if the route still exists - it may have been deleted by the previous statement action
		if db.IsEntityPresentFunc != nil {
			if !(db.IsEntityPresentFunc(params)) {
				db.Logger.Info("This entity no longer exists")
				break
			}
		}
	}
	return ret
}
func (db *PolicyEngineDB) PolicyEngineUndoApplyPolicyForEntity(entity PolicyEngineFilterEntityParams, updateInfo PolicyEngineApplyInfo, params interface{}) bool {
	info := updateInfo.ApplyPolicy
	match, _ := db.PolicyEngineMatchConditions(entity, info.Conditions, "all")
	db.Logger.Info("In PolicyEngineUndoApplyPolicyForEntity, match = ", match)
	if !match {
		db.Logger.Info("Apply Policy conditions dont match for this entity")
		return false
	}
	return db.PolicyEngineUndoPolicyForEntity(entity, info.ApplyPolicy, updateInfo, params)
}
func (db *PolicyEngineDB) PolicyEngineImplementActions(entity PolicyEngineFilterEntityParams, action PolicyAction,
	conditionInfoList []interface{}, policy Policy, params interface{}, policyStmt PolicyStmt) (
	policyActionList []PolicyAction) {
	db.Logger.Info("policyEngineImplementActions")
	policyActionList = make([]PolicyAction, 0)
	addActionToList := false
	switch action.ActionType {
	case policyCommonDefs.PolicyActionTypeRouteDisposition, policyCommonDefs.PolicyActionTypeRouteRedistribute,
		policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, policyCommonDefs.PolicyActionTypeAggregate,
		policyCommonDefs.PolicyActionTypeRIBIn, policyCommonDefs.PolicyActionTypeRIBOut, policyCommonDefs.PolicyActionTypeSetCommunity,
		policyCommonDefs.PolicyActionTypeSetExtendedCommunity, policyCommonDefs.PolicyActionTypeSetLocalPref:
		if entity.DeletePath == true {
			db.Logger.Info("action to be reversed", action.ActionType)
			if db.UndoActionfuncMap[action.ActionType] != nil {
				db.UndoActionfuncMap[action.ActionType](action.ActionInfo, conditionInfoList, policy, params, policyStmt)
			}
			addActionToList = true
		} else { //if entity.CreatePath == true or neither create/delete is valid - in case this function is called a a part of policy create{
			db.Logger.Info("action to be applied", action.ActionType)
			if db.ActionfuncMap[action.ActionType] != nil {
				db.ActionfuncMap[action.ActionType](action.ActionInfo, conditionInfoList, policy, params, policyStmt)
			}
			addActionToList = true
		}
	default:
		db.Logger.Err("UnknownInvalid type of action")
		break
	}
	if addActionToList == true {
		policyActionList = append(policyActionList, action)
	}
	return policyActionList
}

/*
func (db *PolicyEngineDB) FindPrefixMatch(ipAddr string, ipPrefix patriciaDB.Prefix, condition PolicyCondition) (match bool) {
	db.Logger.Info("Prefix match policy ", policyName)
	policyListItem := db.PrefixPolicyListDB.GetLongestPrefixNode(ipPrefix)
	if policyListItem == nil {
		db.Logger.Info("intf stored at prefix ", ipPrefix, " is nil")
		return false
	}
	if policyListItem != nil && reflect.TypeOf(policyListItem).Kind() != reflect.Slice {
		db.Logger.Err("Incorrect data type for this prefix ")
		return false
	}
	policyListSlice := reflect.ValueOf(policyListItem)
	for idx := 0; idx < policyListSlice.Len(); idx++ {
		prefixPolicyListInfo := policyListSlice.Index(idx).Interface().(PrefixPolicyListInfo)
		if prefixPolicyListInfo.policyName != policyName {
			db.Logger.Info("Found a potential match for this prefix but the policy ", policyName, " is not what we are looking for")
			continue
		}
		if prefixPolicyListInfo.lowRange == -1 && prefixPolicyListInfo.highRange == -1 {
			db.Logger.Info("Looking for exact match condition for prefix ", prefixPolicyListInfo.ipPrefix)
			if bytes.Equal(ipPrefix, prefixPolicyListInfo.ipPrefix) {
				db.Logger.Info(" Matched the prefix")
				return true
			} else {
				db.Logger.Info(" Did not match the exact prefix")
				return false
			}
		}
		tempSlice := strings.Split(ipAddr, "/")
		maskLen, err := strconv.Atoi(tempSlice[1])
		if err != nil {
			db.Logger.Err("err getting maskLen")
			return false
		}
		db.Logger.Info("Mask len = ", maskLen)
		if maskLen < prefixPolicyListInfo.lowRange || maskLen > prefixPolicyListInfo.highRange {
			db.Logger.Info("Mask range of the route ", maskLen, " not within the required mask range:", prefixPolicyListInfo.lowRange, "..", prefixPolicyListInfo.highRange)
			return false
		} else {
			db.Logger.Info("Mask range of the route ", maskLen, " within the required mask range:", prefixPolicyListInfo.lowRange, "..", prefixPolicyListInfo.highRange)
			return true
		}
	}
	return match
}
*/
func (db *PolicyEngineDB) FindPrefixMatch(ipAddr string, condition PolicyCondition) (match bool) {
	db.Logger.Info("ipAddr : ", ipAddr, " condition.IpPrefix: ", condition.ConditionInfo.(MatchPrefixConditionInfo).IpPrefix, " conditionInfo,MaskLengthRange: ", condition.ConditionInfo.(MatchPrefixConditionInfo).Prefix.IpPrefix, condition.ConditionInfo.(MatchPrefixConditionInfo).UsePrefixSet)
	conditionInfo := condition.ConditionInfo.(MatchPrefixConditionInfo)
	if conditionInfo.UsePrefixSet {
		db.Logger.Info("FindPrefixMatch:use prefixset")
		item := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(conditionInfo.PrefixSet))
		if item != nil {
			prefixSet := item.(PolicyPrefixSet)
			db.Logger.Info("FindPrefixMatch, prefixSet matchInfoList :", prefixSet.MatchInfoList)
			for _, matchInfo := range prefixSet.MatchInfoList {
				match = netUtils.CheckIfInRange(ipAddr, matchInfo.Prefix.IpPrefix, matchInfo.LowRange, matchInfo.HighRange)
				if match {
					db.Logger.Info("Matched prefix for:", matchInfo.Prefix.IpPrefix, matchInfo.LowRange, matchInfo.HighRange)
					break
				}
			}
		} else {
			db.Logger.Err("prefix set ", conditionInfo.PrefixSet, " not found")
		}
	} else {
		match = netUtils.CheckIfInRange(ipAddr, conditionInfo.Prefix.IpPrefix, conditionInfo.LowRange, conditionInfo.HighRange)
	}
	/*	if conditionInfo.LowRange == -1 && conditionInfo.HighRange == -1 {
		_, ipNet, err := net.ParseCIDR(condition.ConditionInfo.(MatchPrefixConditionInfo).Prefix.IpPrefix)
		if err != nil {
			return false
		}
		if bytes.Equal(ipPrefix, conditionInfo.IpPrefix) {
			db.Logger.Info(" Matched the prefix")
			return true
		}
		networkMask := ipNet.Mask
		vdestMask := net.IPv4Mask(networkMask[0], networkMask[1], networkMask[2], networkMask[3])
		destIp := (net.IP(ipPrefix)).Mask(vdestMask)
		db.Logger.Info("networkMask: ", networkMask, " vdestMask: ", vdestMask, " destIp: ", destIp, "Looking for exact match condition for prefix ", conditionInfo.IpPrefix, " and ", destIp)
		if bytes.Equal(destIp, conditionInfo.IpPrefix) {
			db.Logger.Info(" Matched the prefix")
			return true
		} else {
			db.Logger.Info(" Did not match the exact prefix")
			return false
		}
	}*/
	/*	tempSlice := strings.Split(ipAddr, "/")
		maskLen, err := strconv.Atoi(tempSlice[1])
		if err != nil {
			db.Logger.Err("err getting maskLen")
			return false
		}
		db.Logger.Info("Mask len = ", maskLen)
		if maskLen < conditionInfo.LowRange || maskLen > conditionInfo.HighRange {
			db.Logger.Info("Mask range of the route ", maskLen, " not within the required mask range:", conditionInfo.LowRange, "-", conditionInfo.HighRange)
			return false
		} else {
			db.Logger.Info("Mask range of the route ", maskLen, " within the required mask range:", conditionInfo.LowRange, "-", conditionInfo.HighRange)
			return true
		}*/
	/*baseIp, _, err := net.ParseCIDR(condition.ConditionInfo.(MatchPrefixConditionInfo).Prefix.IpPrefix)
	if err != nil {
		db.Logger.Info("Invalid condition ip:", condition.ConditionInfo.(MatchPrefixConditionInfo).Prefix.IpPrefix)
		return false
	}
	testIp, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		db.Logger.Err("Invalid ipAddr:", ipAddr)
		return false
	}
	match = netUtils.CheckIfInRange(testIp.String(), baseIp.String(), conditionInfo.LowRange, conditionInfo.HighRange)*/
	return match
}

func (db *PolicyEngineDB) DstIpPrefixMatchConditionfunc(entity PolicyEngineFilterEntityParams, condition PolicyCondition) (match bool) {
	db.Logger.Info("dstIpPrefixMatchConditionfunc")
	/*ipPrefix, err := netUtils.GetNetworkPrefixFromCIDR(entity.DestNetIp)
	if err != nil {
		db.Logger.Info("Invalid ipPrefix for the route ", entity.DestNetIp)
		return false
	}*/
	match = db.FindPrefixMatch(entity.DestNetIp, condition)
	if match {
		db.Logger.Info("Found a match for this prefix")
	}
	return match
}
func (db *PolicyEngineDB) ProtocolMatchConditionfunc(entity PolicyEngineFilterEntityParams, condition PolicyCondition) (match bool) {
	db.Logger.Info("protocolMatchConditionfunc: check if policy protocol: ", condition.ConditionInfo.(string), " matches entity protocol: ", entity.RouteProtocol)
	matchProto := condition.ConditionInfo.(string)
	if matchProto == entity.RouteProtocol {
		db.Logger.Info("Protocol condition matches")
		match = true
	}
	return match
}

func (db *PolicyEngineDB) NeighborMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	db.Logger.Info("NeighborMatchConditionfunc: check if policy neighbor:", condition.ConditionInfo.(string),
		"matches entity neighbor: ", entity.Neighbor)
	matchNeighbor := condition.ConditionInfo.(string)
	if matchNeighbor == entity.Neighbor {
		db.Logger.Info("Neighbor matches")
		match = true
	}
	return match
}

func (db *PolicyEngineDB) CommunityMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	conditionInfo := condition.ConditionInfo.(MatchCommunityConditionInfo)
	if conditionInfo.UseSet {
		db.Logger.Info("FindCommunityMatch:useset")
		item := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(conditionInfo.Set))
		if item != nil {
			set := item.(PolicyCommunitySet)
			db.Logger.Info("FindCommunityMatch, communitySet matchInfoList :", set.MatchInfoList)
			for _, community := range entity.Community {
				for _, matchInfo := range set.MatchInfoList {
					if matchInfo.Community == community {
						db.Logger.Info("Community", community, "matches")
						match = true
						break
					}
				}
				if match {
					break
				}
			}
		} else {
			db.Logger.Err("community set ", conditionInfo.Set, " not found")
		}
	} else {
		for _, community := range entity.Community {
			if conditionInfo.Community == community {
				db.Logger.Info("Community", community, "matches")
				match = true
				break
			}
		}
	}
	return match
}

func (db *PolicyEngineDB) ASPathMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	db.Logger.Info("ASPathMatchConditionfunc: check if policy as path matches entity ASPath: ", entity.ASPath)
	conditionInfo := condition.ConditionInfo.(MatchASPathConditionInfo)
	if conditionInfo.UseSet {
		db.Logger.Info("FindASPathMatch:useset")
		item := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(conditionInfo.Set))
		if item != nil {
			set := item.(PolicyASPathSet)
			db.Logger.Info("FindASPathMatch, ASPathSet matchInfoList :", set.MatchInfoList)
			for _, matchInfo := range set.MatchInfoList {
				if bgpUtils.MatchASPath(matchInfo.ASPath, entity.ASPath) {
					db.Logger.Info("ASPath matches")
					match = true
					break
				}
			}
		} else {
			db.Logger.Err("aspath set ", conditionInfo.Set, " not found")
		}
	} else {
		if bgpUtils.MatchASPath(conditionInfo.ASPath, entity.ASPath) {
			db.Logger.Info("ASPath matches")
			match = true
		}
	}
	return match
}

func (db *PolicyEngineDB) LocalPrefMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	db.Logger.Info("LocalPrefMatchConditionfunc: check if policy LocalPref:", condition.ConditionInfo.(uint32),
		"matches entity LocalPref: ", entity.LocalPref)
	matchEntity := condition.ConditionInfo.(uint32)
	if matchEntity == entity.LocalPref {
		db.Logger.Info("LocalPref matches")
		match = true
	}
	return match
}

func (db *PolicyEngineDB) MEDMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	db.Logger.Info("MEDMatchConditionfunc: check if policy MED:", condition.ConditionInfo.(uint32),
		"matches entity MED: ", entity.MED)
	matchEntity := condition.ConditionInfo.(uint32)
	if matchEntity == entity.MED {
		db.Logger.Info("MED matches")
		match = true
	}
	return match
}

func (db *PolicyEngineDB) ExtendedCommunityMatchConditionfunc(entity PolicyEngineFilterEntityParams,
	condition PolicyCondition) (match bool) {
	db.Logger.Info("ExtendedCommunityMatchConditionfunc: check if policy extended community matches entity extended community: ", entity.ExtendedCommunity)
	conditionInfo := condition.ConditionInfo.(MatchExtendedCommunityConditionInfo)
	if conditionInfo.UseSet {
		db.Logger.Info("FindExtendedCommunityMatch:useset")
		item := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(conditionInfo.Set))
		if item != nil {
			set := item.(PolicyExtendedCommunitySet)
			db.Logger.Info("FindExtendedCommunityMatch, extendedcommunitySet matchInfoList :", set.MatchInfoList)
			for _, extComm := range entity.ExtendedCommunity {
				for _, matchInfo := range set.MatchInfoList {
					if matchInfo.ExtendedCommunity == extComm {
						db.Logger.Info("Extended Community", extComm, "matches")
						match = true
						break
					}
				}
				if match {
					break
				}
			}
		} else {
			db.Logger.Err("Extended community set ", conditionInfo.Set, " not found")
		}
	} else {
		for _, extComm := range entity.ExtendedCommunity {
			if conditionInfo.ExtendedCommunity == extComm {
				db.Logger.Info("Extended Community", extComm, "matches")
				match = true
				break
			}
		}
	}
	return match
}

func (db *PolicyEngineDB) ConditionCheckValid(entity PolicyEngineFilterEntityParams, conditionsList []string, policyStmt PolicyStmt) (valid bool) {
	db.Logger.Info("conditionCheckValid")
	valid = true
	if conditionsList == nil {
		db.Logger.Info("No conditions to match, so valid")
		return true
	}
	for i := 0; i < len(conditionsList); i++ {
		db.Logger.Info("Find policy condition number ", i, " name ", policyStmt.Conditions[i], " in the condition database")
		conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(conditionsList[i]))
		if conditionItem == nil {
			db.Logger.Info("Did not find condition ", conditionsList[i], " in the condition database")
			continue
		}
		condition := conditionItem.(PolicyCondition)
		db.Logger.Info("policy condition number ", i, " type ", condition.ConditionType)
		if db.ConditionCheckfuncMap[condition.ConditionType] != nil {
			match := db.ConditionCheckfuncMap[condition.ConditionType](entity, condition)
			if !match {
				db.Logger.Info("Condition does not match")
				return false
			}
		}
	}
	db.Logger.Info("returning valid= ", valid)
	return valid
}
func (db *PolicyEngineDB) PolicyEngineMatchConditions(entity PolicyEngineFilterEntityParams, conditions []string, matchConditions string) (match bool, conditionsList []PolicyCondition) {
	db.Logger.Info("policyEngineMatchConditions")
	var i int
	allConditionsMatch := true
	anyConditionsMatch := false
	addConditiontoList := false
	conditionsList = make([]PolicyCondition, 0)
	if len(conditions) == 0 {
		allConditionsMatch = true
		anyConditionsMatch = true
	} else {
		for i = 0; i < len(conditions); i++ {
			addConditiontoList = false
			db.Logger.Debug("Find policy condition number ", i, " name ", conditions[i], " in the condition database")
			conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(conditions[i]))
			if conditionItem == nil {
				db.Logger.Info("Did not find condition ", conditions[i], " in the condition database")
				continue
			}
			condition := conditionItem.(PolicyCondition)
			db.Logger.Debug("policy condition number ", i, "  type ", condition.ConditionType)
			if db.ConditionCheckfuncMap[condition.ConditionType] != nil {
				match = db.ConditionCheckfuncMap[condition.ConditionType](entity, condition)
				db.Logger.Debug("match for condition:", condition.Name, " : ", match)
				if match {
					db.Logger.Info("Condition match found")
					anyConditionsMatch = true
					addConditiontoList = true
				} else {
					db.Logger.Info("Condition:", condition.Name, " not matched, set allConditionsMatch to false")
					allConditionsMatch = false
				}
			}
			if addConditiontoList == true {
				conditionsList = append(conditionsList, condition)
			}
		}
	}
	if matchConditions == "all" && allConditionsMatch == true {
		db.Logger.Info("retuning true because matchConditions:", matchConditions, " and allConditionsMatch:", allConditionsMatch)
		return true, conditionsList
	}
	if matchConditions == "any" && anyConditionsMatch == true {
		db.Logger.Info("returning true because matchConditions:", matchConditions, " and anyConditionsMatch:", anyConditionsMatch)
		return true, conditionsList
	}
	return false, conditionsList
}
func (db *PolicyEngineDB) PolicyEngineApplyPolicyStmt(entity *PolicyEngineFilterEntityParams, info ApplyPolicyInfo,
	policyStmt PolicyStmt, policyPath int, params interface{}, hit *bool, deleted *bool) {
	policy := info.ApplyPolicy
	db.Logger.Info("policyEngineApplyPolicyStmt - ", policyStmt.Name)
	var policyConditionList []PolicyCondition
	var conditionList []PolicyCondition
	conditionInfoList := make([]interface{}, 0)
	var match bool
	if policyStmt.Conditions == nil && info.Conditions == nil {
		db.Logger.Info("No policy conditions")
		*hit = true
	} else {
		//match, ret_conditionList := db.PolicyEngineMatchConditions(*entity, policyStmt)
		match, policyConditionList = db.PolicyEngineMatchConditions(*entity, policyStmt.Conditions, policyStmt.MatchConditions)
		db.Logger.Info("match = ", match)
		*hit = match
		if !match {
			db.Logger.Info("Stmt Conditions do not match")
			return
		}
		db.Logger.Debug("PolicyEngineApplyStmt policyConditionList after checking with the policystmt:", policyConditionList)
		for j := 0; j < len(policyConditionList); j++ {
			conditionInfoList = append(conditionInfoList, policyConditionList[j].ConditionInfo)
		}
		db.Logger.Debug("PolicyEngineApplyStmt conditionInfoList after adding policyConditionList:", conditionInfoList)
		match, conditionList = db.PolicyEngineMatchConditions(*entity, info.Conditions, "all")
		db.Logger.Debug("PolicyEngineApplyStmt conditionList after checking with the applyInfo:", conditionList)
		db.Logger.Info("match = ", match)
		*hit = match
		if !match {
			db.Logger.Info("Extra Conditions do not match")
			return
		}
		for j := 0; j < len(conditionList); j++ {
			conditionInfoList = append(conditionInfoList, conditionList[j].ConditionInfo)
			policyConditionList = append(policyConditionList, conditionList[j])
		}
		db.Logger.Debug("PolicyEngineApplyStmt conditionInfoList after checking with the applyInfo:", conditionInfoList)
	}
	actionList := db.PolicyEngineImplementActions(*entity, info.Action, conditionInfoList, info.ApplyPolicy, params, policyStmt)
	if db.ActionListHasAction(actionList, policyCommonDefs.PolicyActionTypeRouteDisposition, "Reject") {
		db.Logger.Info("Reject action was applied for this entity")
		*deleted = true
	}
	//check if the route still exists - it may have been deleted by the previous statement action
	if db.IsEntityPresentFunc != nil {
		*deleted = !(db.IsEntityPresentFunc(params))
	}
	policyInfoGet := db.PolicyDB.Get(patriciaDB.Prefix(policy.Name))
	if policyInfoGet != nil {
		policyInfo := policyInfoGet.(Policy)
		db.Logger.Info("PolicyEngineApplyPolicyStmt, before updateEntityDB db.Global:", db.Global, " policyInfo:", policyInfo)
	}
	db.AddPolicyEntityMapEntry(*entity, policy.Name, policyStmt.Name, policyConditionList, actionList)
	if db.UpdateEntityDB != nil {
		policyDetails := PolicyDetails{Policy: policy.Name, PolicyStmt: policyStmt.Name, ConditionList: conditionList, ActionList: actionList, EntityDeleted: *deleted}
		db.UpdateEntityDB(policyDetails, params)
	}
	policyInfoGet = db.PolicyDB.Get(patriciaDB.Prefix(policy.Name))
	if policyInfoGet != nil {
		policyInfo := policyInfoGet.(Policy)
		db.Logger.Info("PolicyEngineApplyPolicyStmt, after updateEntityDB db.Global:", db.Global, " policyInfo:", policyInfo)
	}
}

func (db *PolicyEngineDB) PolicyEngineApplyPolicy(entity *PolicyEngineFilterEntityParams, info ApplyPolicyInfo, policyEngineApplyInfo PolicyEngineApplyInfo, policyPath int, params interface{}, hit *bool) {
	db.Logger.Info("policyEngineApplyPolicy - ", info.ApplyPolicy.Name)
	policy := info.ApplyPolicy

	stmtMap := make(map[string]bool)
	for _, stmt := range policyEngineApplyInfo.StmtList {
		stmtMap[stmt] = true
	}
	var policyStmtKeys []int
	deleted := false
	for k := range policy.PolicyStmtPrecedenceMap {
		db.Logger.Info("key k = ", k)
		policyStmtKeys = append(policyStmtKeys, k)
	}
	sort.Ints(policyStmtKeys)
	for i := 0; i < len(policyStmtKeys); i++ {
		db.Logger.Info("Key: ", policyStmtKeys[i], " policyStmtName ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
		//if the undo stmt list is non zero, then this is not the case for policy delete but for policy update
		if policyEngineApplyInfo.StmtList != nil && len(policyEngineApplyInfo.StmtList) > 0 {
			_, ok := stmtMap[policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]]]
			if !ok {
				db.Logger.Info("this statement ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]], " is not the one to be added to the policy")
				//return value should be false, so the policy is not deleted from the entity
				continue
			}
		}
		policyStmt := db.PolicyStmtDB.Get((patriciaDB.Prefix(policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])))
		if policyStmt == nil {
			db.Logger.Info("Invalid policyStmt")
			continue
		}
		db.PolicyEngineApplyPolicyStmt(entity, info, policyStmt.(PolicyStmt), policyPath, params, hit, &deleted)
		if deleted == true {
			db.Logger.Info("Entity was deleted as a part of the policyStmt ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
			break
		}
		if *hit == true {
			if policy.MatchType == "any" {
				db.Logger.Info("Match type for policy ", policy.Name, " is any and the policy stmt ", (policyStmt.(PolicyStmt)).Name, " is a hit, no more policy statements will be executed")
				break
			}
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineApplyForEntity(entity PolicyEngineFilterEntityParams, policyData interface{}, params interface{}) {
	db.Logger.Info("policyEngineApplyForEntity")
	updateInfo := policyData.(PolicyEngineApplyInfo)
	info := updateInfo.ApplyPolicy
	policy := info.ApplyPolicy
	policyHit := false
	if len(entity.PolicyList) == 0 {
		db.Logger.Info("This route has no policy applied to it so far, just apply the new policy")
		db.PolicyEngineApplyPolicy(&entity, info, updateInfo, policyCommonDefs.PolicyPath_All, params, &policyHit)
	} else {
		db.Logger.Info("This route already has policy applied to it - len(route.PolicyList) - ", len(entity.PolicyList))

		for i := 0; i < len(entity.PolicyList); i++ {
			db.Logger.Info("policy at index ", i)
			policyInfo := db.PolicyDB.Get(patriciaDB.Prefix(entity.PolicyList[i]))
			if policyInfo == nil {
				db.Logger.Info("Unexpected: Invalid policy in the route policy list")
			} else {
				oldPolicy := policyInfo.(Policy)
				if !isPolicyTypeSame(oldPolicy, policy) {
					db.Logger.Info("The policy type applied currently is not the same as new policy, so apply new policy")
					db.PolicyEngineApplyPolicy(&entity, info, updateInfo, policyCommonDefs.PolicyPath_All, params, &policyHit)
				} else if oldPolicy.Precedence < policy.Precedence {
					db.Logger.Info("The policy types are same and precedence of the policy applied currently is lower than the new policy, so do nothing")
					return
				} else {
					db.Logger.Info("The new policy's precedence is lower, so undo old policy's actions and apply the new policy")
					//db.PolicyEngineUndoPolicyForEntity(entity, oldPolicy, params)
					db.PolicyEngineApplyPolicy(&entity, info, updateInfo, policyCommonDefs.PolicyPath_All, params, &policyHit)
				}
			}
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineReverseGlobalPolicyStmt(policy Policy, policyStmt PolicyStmt) {
	db.Logger.Info("policyEngineApplyGlobalPolicyStmt - ", policyStmt.Name)
	var conditionItem interface{} = nil
	//global policies can only have statements with 1 condition and 1 action
	if policyStmt.Actions == nil {
		db.Logger.Info("No policy actions defined")
		return
	}
	if policyStmt.Conditions == nil {
		db.Logger.Info("No policy conditions")
	} else {
		if len(policyStmt.Conditions) > 1 {
			db.Logger.Info("only 1 condition allowed for global policy stmt")
			return
		}
		conditionItem = db.PolicyConditionsDB.Get(patriciaDB.Prefix(policyStmt.Conditions[0]))
		if conditionItem == nil {
			db.Logger.Info("Condition ", policyStmt.Conditions[0], " not found")
			return
		}
		actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[0]))
		if actionItem == nil {
			db.Logger.Info("Action ", policyStmt.Actions[0], " not found")
			return
		}
		actionInfo := actionItem.(PolicyAction)
		if db.UndoActionfuncMap[actionInfo.ActionType] != nil {
			//since global policies have just 1 condition, we can pass that as the params to the undo call
			db.UndoActionfuncMap[actionInfo.ActionType](actionItem, nil, policy, conditionItem, policyStmt)
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineApplyGlobalPolicyStmt(policy Policy, policyStmt PolicyStmt) {
	db.Logger.Info("policyEngineApplyGlobalPolicyStmt - ", policyStmt.Name)
	var conditionItem interface{} = nil
	//global policies can only have statements with 1 condition and 1 action
	if policyStmt.Actions == nil {
		db.Logger.Info("No policy actions defined")
		return
	}
	if policyStmt.Conditions == nil {
		db.Logger.Info("No policy conditions")
	} else {
		if len(policyStmt.Conditions) > 1 {
			db.Logger.Info("only 1 condition allowed for global policy stmt")
			return
		}
		conditionItem = db.PolicyConditionsDB.Get(patriciaDB.Prefix(policyStmt.Conditions[0]))
		if conditionItem == nil {
			db.Logger.Info("Condition ", policyStmt.Conditions[0], " not found")
			return
		}
		policyCondition := conditionItem.(PolicyCondition)
		conditionInfoList := make([]interface{}, 0)
		conditionInfoList = append(conditionInfoList, policyCondition.ConditionInfo)

		actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[0]))
		if actionItem == nil {
			db.Logger.Info("Action ", policyStmt.Actions[0], " not found")
			return
		}
		actionInfo := actionItem.(PolicyAction)
		if db.ActionfuncMap[actionInfo.ActionType] != nil {
			db.ActionfuncMap[actionInfo.ActionType](actionInfo.ActionInfo, conditionInfoList, policy, nil, policyStmt)
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineReverseGlobalPolicy(policy Policy) {
	db.Logger.Info("policyEngineReverseGlobalPolicy")
	var policyStmtKeys []int
	for k := range policy.PolicyStmtPrecedenceMap {
		db.Logger.Info("key k = ", k)
		policyStmtKeys = append(policyStmtKeys, k)
	}
	sort.Ints(policyStmtKeys)
	for i := 0; i < len(policyStmtKeys); i++ {
		db.Logger.Info("Key: ", policyStmtKeys[i], " policyStmtName ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
		policyStmt := db.PolicyStmtDB.Get((patriciaDB.Prefix(policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])))
		if policyStmt == nil {
			db.Logger.Info("Invalid policyStmt")
			continue
		}
		db.PolicyEngineReverseGlobalPolicyStmt(policy, policyStmt.(PolicyStmt))
	}
}
func (db *PolicyEngineDB) PolicyEngineApplyGlobalPolicy(policy Policy) {
	db.Logger.Info("policyEngineApplyGlobalPolicy")
	var policyStmtKeys []int
	for k := range policy.PolicyStmtPrecedenceMap {
		db.Logger.Info("key k = ", k)
		policyStmtKeys = append(policyStmtKeys, k)
	}
	sort.Ints(policyStmtKeys)
	for i := 0; i < len(policyStmtKeys); i++ {
		db.Logger.Info("Key: ", policyStmtKeys[i], " policyStmtName ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
		policyStmt := db.PolicyStmtDB.Get((patriciaDB.Prefix(policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])))
		if policyStmt == nil {
			db.Logger.Info("Invalid policyStmt")
			continue
		}
		db.PolicyEngineApplyGlobalPolicyStmt(policy, policyStmt.(PolicyStmt))
	}
}

func (db *PolicyEngineDB) PolicyEngineTraverseAndApplyPolicy(info ApplyPolicyInfo, stmtList []string, conditionsList []string) {
	db.Logger.Info("PolicyEngineTraverseAndApplyPolicy -  apply policy ", info.ApplyPolicy.Name)
	if db.TraverseAndApplyPolicyFunc != nil {
		db.Logger.Info("Calling TraverseAndApplyPolicyFunc function")
		updateStmtInfo := PolicyEngineApplyInfo{
			ApplyPolicy:    info,
			StmtList:       stmtList,
			ConditionsList: conditionsList,
		}
		db.TraverseAndApplyPolicyFunc(updateStmtInfo, db.PolicyEngineApplyForEntity)
	}
	/*	if policy.ExportPolicy || policy.ImportPolicy {
			db.Logger.Info("Applying import/export policy to all routes"))
			if db.TraverseAndApplyPolicyFunc != nil {
				db.Logger.Info("Calling TraverseAndApplyPolicyFunc function"))
				db.TraverseAndApplyPolicyFunc(policy, db.PolicyEngineApplyForEntity)
			}
		} else if policy.GlobalPolicy {
			db.Logger.Info("Need to apply global policy"))
			db.PolicyEngineApplyGlobalPolicy(policy)
		}*/
}

/*func (db *PolicyEngineDB) PolicyEngineTraverseAndReversePolicyStmts(info ApplyPolicyInfo, stmtList []string, conditionsList []string) {
	db.Logger.Info("PolicyEngineTraverseAndReversePolicyStmts -policy:", info.ApplyPolicy.Name)
	updateStmtInfo := PolicyEngineApplyInfo{
		ApplyPolicy:    info,
		StmtList:       stmtList,
		ConditionsList: conditionsList,
	}
	if db.TraverseAndReversePolicyFunc != nil {
		db.TraverseAndReversePolicyFunc(updateStmtInfo)
	}
}
*/
func (db *PolicyEngineDB) PolicyEngineTraverseAndReversePolicy(info ApplyPolicyInfo, stmtList []string, conditionsList []string) {
	db.Logger.Info("PolicyEngineTraverseAndReversePolicy -  reverse policy ", info.ApplyPolicy.Name)
	if db.TraverseAndReversePolicyFunc != nil {
		db.Logger.Info("Calling PolicyEngineTraverseAndReversePolicy function")
		updateStmtInfo := PolicyEngineApplyInfo{
			ApplyPolicy:    info,
			StmtList:       stmtList,
			ConditionsList: conditionsList,
		}
		db.TraverseAndReversePolicyFunc(updateStmtInfo)
	}
	/*	if policy.ExportPolicy || policy.ImportPolicy {
			db.Logger.Info("Reversing import/export policy ")
			db.TraverseAndReversePolicyFunc(info)
		} else if policy.GlobalPolicy {
			db.Logger.Info("Need to reverse global policy")
			db.PolicyEngineReverseGlobalPolicy(info)
		}*/
}

func (db *PolicyEngineDB) PolicyEngineFilter(entity PolicyEngineFilterEntityParams, policyPath int, params interface{}) {
	/*db.Logger.Info("PolicyEngineFilter")
	var policyPath_Str string
	if policyPath == policyCommonDefs.PolicyPath_Import {
		policyPath_Str = "Import"
	} else if policyPath == policyCommonDefs.PolicyPath_Export {
		policyPath_Str = "Export"
	} else if policyPath == policyCommonDefs.PolicyPath_All {
		policyPath_Str = "ALL"
		db.Logger.Info("policy path ", policyPath_Str, " unexpected in this function")
		return
	}
	db.Logger.Info("PolicyEngineFilter for policypath ", policyPath_Str, "create = ", entity.CreatePath, " delete = ", entity.DeletePath, " route: ", entity.DestNetIp, " protocol type: ", entity.RouteProtocol)*/
	var policyKeys []int
	var policyHit bool
	idx := 0
	var policyInfo interface{}
	if policyPath == policyCommonDefs.PolicyPath_Import {
		for k := range db.ImportPolicyPrecedenceMap {
			policyKeys = append(policyKeys, k)
		}
	} else if policyPath == policyCommonDefs.PolicyPath_Export {
		for k := range db.ExportPolicyPrecedenceMap {
			policyKeys = append(policyKeys, k)
		}
	}
	sort.Ints(policyKeys)
	for {
		if entity.DeletePath == true { //policyEngineFilter called during delete
			if entity.PolicyList != nil {
				if idx >= len(entity.PolicyList) {
					break
				}
				//db.Logger.Info("getting policy ", idx, " from entity.PolicyList")
				policyInfo = db.PolicyDB.Get(patriciaDB.Prefix(entity.PolicyList[idx]))
				if policyInfo == nil {
					db.Logger.Info("policy nil for ", entity.PolicyList[idx], " during delete path of policyengin filter")
					continue
				}
				idx++
				if policyInfo.(Policy).ExportPolicy && policyPath == policyCommonDefs.PolicyPath_Import || policyInfo.(Policy).ImportPolicy && policyPath == policyCommonDefs.PolicyPath_Export {
					//		db.Logger.Info("policy ", policyInfo.(Policy).Name, " not the same type as the policypath -", policyPath_Str)
					continue
				}
			} else {
				//db.Logger.Info("PolicyList empty and this is a delete operation, so break")
				break
			}
		} else if entity.CreatePath == true { //policyEngine filter called during create
			//db.Logger.Info("idx = ", idx, " len(policyKeys):", len(policyKeys))
			if idx >= len(policyKeys) {
				break
			}
			policyName := ""
			if policyPath == policyCommonDefs.PolicyPath_Import {
				policyName = db.ImportPolicyPrecedenceMap[policyKeys[idx]]
			} else if policyPath == policyCommonDefs.PolicyPath_Export {
				policyName = db.ExportPolicyPrecedenceMap[policyKeys[idx]]
			}
			//db.Logger.Info("getting policy  ", idx, " policyKeys[idx] = ", policyKeys[idx], " ", policyName, " from PolicyDB")
			policyInfo = db.PolicyDB.Get((patriciaDB.Prefix(policyName)))
			idx++
		}
		if policyInfo == nil {
			db.Logger.Info("Nil policy")
			break
		}
		policy := policyInfo.(Policy)
		localPolicyDB := *db.LocalPolicyDB
		if localPolicyDB != nil {
			if len(localPolicyDB) > int(policy.LocalDBSliceIdx) && localPolicyDB[policy.LocalDBSliceIdx].IsValid == false {
				//db.Logger.Info("Invalid policy at localDB slice idx ", policy.LocalDBSliceIdx)
				continue
			}
		}
		info, ok := db.ApplyPolicyMap[policy.Name]
		if !ok || info.Count == 0 {
			db.Logger.Info("no application for this policy ", policy.Name)
			continue
		}
		applyList := info.InfoList
		for j := 0; j < len(applyList); j++ {
			db.PolicyEngineApplyPolicy(&entity, applyList[j], PolicyEngineApplyInfo{}, policyPath, params, &policyHit)
			if policyHit {
				//db.Logger.Info("Policy ", policy.Name, " applied to the route")
				break
			}
		}
	}
	if entity.PolicyHitCounter == 0 {
		//db.Logger.Info("Need to apply default policy, policyPath = ", policyPath, "policyPath_Str= ", policyPath_Str)
		if policyPath == policyCommonDefs.PolicyPath_Import {
			//db.Logger.Info("Applying default import policy")
			if db.DefaultImportPolicyActionFunc != nil {
				db.DefaultImportPolicyActionFunc(nil, nil, Policy{}, params, PolicyStmt{})
			}
		} else if policyPath == policyCommonDefs.PolicyPath_Export {
			//db.Logger.Info("Applying default export policy")
			if db.DefaultExportPolicyActionFunc != nil {
				db.DefaultExportPolicyActionFunc(nil, nil, Policy{}, params, PolicyStmt{})
			}
		}
	}
	if entity.DeletePath == true {
		db.DeletePolicyEntityMapEntry(entity, "")
	}
}
