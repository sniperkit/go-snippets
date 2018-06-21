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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"ribd"
	"ribdInt"
	"strings"
	"utils/patriciaDB"
	"utils/policy"
	"utils/policy/policyCommonDefs"
)

/*
   This structure can be used along with policyDefinitionConfig object to pass on any application specific
   info to policy engine
*/
type PolicyExtensions struct {
	hitCounter    int
	routeList     []string
	routeInfoList []ribdInt.Routes
}
type Policy struct {
	*policy.Policy
	hitCounter    int
	routeList     []string
	routeInfoList []ribdInt.Routes
}

/*
   Function to create policy condition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyConditionConfigCreate(cfg *ribd.PolicyCondition, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyConditionConfigCreate:CreatePolicyConditioncfg: ", cfg.Name, " cfg:", *cfg)
	newPolicy := policy.PolicyConditionConfig{Name: cfg.Name,
		ConditionType:               cfg.ConditionType,
		MatchProtocolConditionInfo:  cfg.Protocol,
		MatchLocalPrefConditionInfo: uint32(cfg.LocalPref),
		MatchMEDConditionInfo:       uint32(cfg.MED),
	}
	matchPrefix := policy.PolicyPrefix{IpPrefix: cfg.IpPrefix, MasklengthRange: cfg.MaskLengthRange}
	newPolicy.MatchDstIpPrefixConditionInfo = policy.PolicyDstIpMatchPrefixSetCondition{Prefix: matchPrefix, PrefixSet: cfg.PrefixSet}

	newPolicy.MatchCommunityConditionInfo = policy.PolicyMatchCommunitySetCondition{Community: cfg.Community, CommunitySet: cfg.CommunitySet}
	logger.Debug("cfg.CommunitySet:", cfg.CommunitySet, " newPolicy.MatchCommunityConditionInfo:", newPolicy.MatchCommunityConditionInfo)
	matchExtendedCommunityInfo := policy.PolicyExtendedCommunityInfo{cfg.ExtendedCommunityType, cfg.ExtendedCommunityValue}
	newPolicy.MatchExtendedCommunityConditionInfo = policy.PolicyMatchExtendedCommunitySetCondition{ExtendedCommunity: matchExtendedCommunityInfo, ExtendedCommunitySet: cfg.ExtendedCommunitySet}

	newPolicy.MatchASPathConditionInfo = policy.PolicyMatchASPathSetCondition{ASPath: cfg.ASPath, ASPathSet: cfg.ASPathSet}
	val, err = db.CreatePolicyCondition(newPolicy)
	return val, err
}

/*
   Function to delete policy condition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyConditionConfigDelete(cfg *ribd.PolicyCondition, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyConditionConfigDelete:DeletePolicyCondition: ", cfg.Name)
	newPolicy := policy.PolicyConditionConfig{Name: cfg.Name}
	val, err = db.DeletePolicyCondition(newPolicy)
	return val, err
}

/*
   Function to update policy condition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyConditionConfigUpdate(origCfg *ribd.PolicyCondition, newCfg *ribd.PolicyCondition, attrset []bool, db *policy.PolicyEngineDB) (err error) {
	func_msg := "ProcessPolicyConditionConfigUpdate for condition " + origCfg.Name
	logger.Debug(func_msg)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy condition statement")
		return errors.New("Policy prefix condition to be updated is different than the original one")
	}
	if attrset != nil {
		objTyp := reflect.TypeOf(*origCfg)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				if objName == "Protocol" {
					logger.Debug(func_msg, " Attr to be updated is Protocol")
					newPolicyCfg := policy.PolicyConditionConfig{
						Name: origCfg.Name,
						MatchProtocolConditionInfo: newCfg.Protocol,
					}
					err = db.UpdatePolicyCondition(newPolicyCfg, "Protocol")
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for Protocol attribute")
						return err
					}
				} else if objName == "IpPrefix" {
					logger.Debug(func_msg, " Attr to be updated is IpPrefix")
					matchPrefix := policy.PolicyPrefix{IpPrefix: newCfg.IpPrefix}
					newPolicyCfg := policy.PolicyConditionConfig{
						Name: origCfg.Name,
					}
					newPolicyCfg.MatchDstIpPrefixConditionInfo = policy.PolicyDstIpMatchPrefixSetCondition{Prefix: matchPrefix}
					err = db.UpdatePolicyCondition(newPolicyCfg, "IpPrefix")
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for Protocol ipprefix")
						return err
					}
				} else if objName == "MaskLengthRange" {
					logger.Debug(func_msg, " Attr to be updated is MaskLengthRange")
					matchPrefix := policy.PolicyPrefix{MasklengthRange: newCfg.MaskLengthRange}
					newPolicyCfg := policy.PolicyConditionConfig{
						Name: origCfg.Name,
					}
					newPolicyCfg.MatchDstIpPrefixConditionInfo = policy.PolicyDstIpMatchPrefixSetCondition{Prefix: matchPrefix}
					err = db.UpdatePolicyCondition(newPolicyCfg, "MaskLengthRange")
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for MaskLengthRange  attribute")
						return err
					}
				} else if objName == "PrefixSet" {
					logger.Debug(func_msg, " Attr to be updated is PrefixSet")
					newPolicyCfg := policy.PolicyConditionConfig{
						Name: origCfg.Name,
					}
					newPolicyCfg.MatchDstIpPrefixConditionInfo = policy.PolicyDstIpMatchPrefixSetCondition{PrefixSet: newCfg.PrefixSet}
					err = db.UpdatePolicyCondition(newPolicyCfg, "PrefixSet")
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for PrefixSet attribute")
						return err
					}
				} else if objName == "ConditionType" {
					logger.Debug(func_msg, " Attr to be updated is ConditionType")
					newPolicyCfg := policy.PolicyConditionConfig{
						Name:          origCfg.Name,
						ConditionType: newCfg.ConditionType,
					}
					err = db.UpdatePolicyCondition(newPolicyCfg, "ConditionType")
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for ConditionType attribute")
						return err
					}
				} else {
					logger.Err(fmt.Sprintln("Update of ", objName, " not supported"))
					return errors.New(fmt.Sprintln("PolicyCondition update for attribute ", objName, " not supported"))
				}
			}
		}
	}
	return err
}

/*
   Function to create policy action in the policyEngineDB
*/
/*func (m RIBDServer) ProcessPolicyActionConfigCreate(cfg *ribdInt.PolicyAction, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyActionConfigCreate:CreatePolicyAction"))
	newAction := policy.PolicyActionConfig{Name: cfg.Name, ActionType: cfg.ActionType, SetAdminDistanceValue: int(cfg.SetAdminDistanceValue), Accept: cfg.Accept, Reject: cfg.Reject, RedistributeAction: cfg.RedistributeAction, RedistributeTargetProtocol: cfg.RedistributeTargetProtocol, NetworkStatementTargetProtocol: cfg.NetworkStatementTargetProtocol}
	val, err = db.CreatePolicyAction(newAction)
	return val, err
}
*/
/*
   Function to delete policy action in the policyEngineDB
*/
/*func (m RIBDServer) ProcessPolicyActionConfigDelete(cfg *ribdInt.PolicyAction, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyActionConfigDelete:CreatePolicyAction"))
	newAction := policy.PolicyActionConfig{Name: cfg.Name}
	val, err = db.DeletePolicyAction(newAction)
	return val, err
}
*/
/*
   Function to create policy stmt in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyStmtConfigCreate(cfg *ribd.PolicyStmt, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyStatementCreate:CreatePolicyStatement")
	newPolicyStmt := policy.PolicyStmtConfig{Name: cfg.Name, MatchConditions: cfg.MatchConditions}
	if len(cfg.Conditions) != 0 {
		newPolicyStmt.Conditions = make([]string, 0)
		for i := 0; i < len(cfg.Conditions); i++ {
			newPolicyStmt.Conditions = append(newPolicyStmt.Conditions, cfg.Conditions[i])
		}
	}
	newPolicyStmt.Actions = make([]string, 0)
	newPolicyStmt.Actions = append(newPolicyStmt.Actions, cfg.Action)
	newPolicyStmt.SetActions = make([]policy.PolicyActionCfg, 0)
	for _, setAction := range cfg.SetActions {
		newPolicyStmt.SetActions = append(newPolicyStmt.SetActions, policy.PolicyActionCfg{
			Attr:              setAction.Attr,
			Community:         setAction.Community,
			ExtendedCommunity: policy.PolicyExtendedCommunityInfo{setAction.ExtendedCommunityType, setAction.ExtendedCommunityValue},
			LocalPref:         uint32(setAction.LocalPref),
			MED:               uint32(setAction.MED),
			PrependASPath:     setAction.PrependASPath,
		})
	}
	err = db.CreatePolicyStatement(newPolicyStmt)
	return err
}

/*
   Function to delete policy stmt in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyStmtConfigDelete(cfg *ribd.PolicyStmt, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyStatementDelete:DeletePolicyStatement for name ", cfg.Name)
	stmt := policy.PolicyStmtConfig{Name: cfg.Name}
	err = db.DeletePolicyStatement(stmt)
	return err
}

/*
   Function to patch update policy stmt in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyStmtConfigPatchUpdate(origCfg *ribd.PolicyStmt, newCfg *ribd.PolicyStmt, op []*ribd.PatchOpInfo, db *policy.PolicyEngineDB) (err error) {
	func_msg := "ProcessPolicyStmtConfigPatchUpdate for stmt " + origCfg.Name + ":"
	logger.Debug(func_msg)
	if origCfg.Name != newCfg.Name {
		logger.Err(func_msg, " Update for a different policy stmt")
		return errors.New("Policy stmt to be updated is different than the original one")
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "Conditions":
			logger.Debug(func_msg, " Patch update for Conditions")
			newPolicyStmt := policy.PolicyStmtConfig{
				Name:            origCfg.Name,
				MatchConditions: origCfg.MatchConditions,
			}
			newPolicyStmt.Actions = make([]string, 0)
			newPolicyStmt.Actions = append(newPolicyStmt.Actions, origCfg.Action)
			newPolicyStmt.Conditions = make([]string, 0)
			var valueObjArr []string
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				logger.Debug(func_msg, "error unmarshaling value:", err)
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug(func_msg, " Number of conditions:", len(valueObjArr))
			for _, val := range valueObjArr {
				logger.Debug(func_msg, " condition to be removed - ", val)
				newPolicyStmt.Conditions = append(newPolicyStmt.Conditions, val)
			}
			switch op[idx].Op {
			case "add":
				db.UpdateAddPolicyStmtConditions(newPolicyStmt)
			case "remove":
				db.UpdateRemovePolicyStmtConditions(newPolicyStmt)
			default:
				logger.Err("Operation ", op[idx].Op, " not supported")
			}
		default:
			logger.Err(func_msg, " Patch update for attribute:", op[idx].Path, " not supported")
			err = errors.New(fmt.Sprintln("Operation ", op[idx].Op, " not supported"))
		}
	}
	return err
}

/*
   Function to update policy stmt in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyStmtConfigUpdate(origCfg *ribd.PolicyStmt, newCfg *ribd.PolicyStmt, attrset []bool, db *policy.PolicyEngineDB) (err error) {
	func_msg := "ProcessPolicyStmtConfigUpdate for statement " + origCfg.Name
	logger.Debug(func_msg)
	if origCfg.Name != newCfg.Name {
		logger.Err(func_msg, "Update for a different policy statement")
		return errors.New("Policy statement to be updated is different than the original one")
	}
	if attrset != nil {
		objTyp := reflect.TypeOf(*origCfg)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrset[i] {
				if objName == "MatchConditions" {
					logger.Debug(func_msg, " Attr to be updated is MatchType")
					newPolicyStmt := policy.PolicyStmtConfig{
						Name:            origCfg.Name,
						MatchConditions: newCfg.MatchConditions,
					}
					err = db.UpdatePolicyStmtMatchTypeAttr(newPolicyStmt)
					if err != nil {
						db.Logger.Err(func_msg, " policylib returned err:", err, " for matchtype attribute")
						return err
					}
				} else {
					logger.Err(fmt.Sprintln("Update of ", objName, " not supported"))
					return errors.New(fmt.Sprintln("PolicyStmt update for attribute ", objName, " not supported"))
				}
			}
		}
	}
	return err
}

/*
   Function to create policy definition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyDefinitionConfigCreate(cfg *ribd.PolicyDefinition, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyDefinitionCreate:CreatePolicyDefinition")
	newPolicy := policy.PolicyDefinitionConfig{Name: cfg.Name, Precedence: int(cfg.Priority), MatchType: cfg.MatchType, PolicyType: cfg.PolicyType}
	newPolicy.PolicyDefinitionStatements = make([]policy.PolicyDefinitionStmtPrecedence, 0)
	var policyDefinitionStatement policy.PolicyDefinitionStmtPrecedence
	for i := 0; i < len(cfg.StatementList); i++ {
		policyDefinitionStatement.Precedence = int(cfg.StatementList[i].Priority)
		policyDefinitionStatement.Statement = cfg.StatementList[i].Statement
		newPolicy.PolicyDefinitionStatements = append(newPolicy.PolicyDefinitionStatements, policyDefinitionStatement)
	}
	newPolicy.Extensions = PolicyExtensions{}
	err = db.CreatePolicyDefinition(newPolicy)
	return err
}

/*
   Function to delete policy definition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyDefinitionConfigDelete(cfg *ribd.PolicyDefinition, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyDefinitionDelete:DeletePolicyDefinition for name ", cfg.Name)
	policy := policy.PolicyDefinitionConfig{Name: cfg.Name}
	err = db.DeletePolicyDefinition(policy)
	return err
}

/*
   Function to patch update policy definition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyDefinitionConfigPatchUpdate(origCfg *ribd.PolicyDefinition, newCfg *ribd.PolicyDefinition, op []*ribd.PatchOpInfo, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyDefinitionConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy")
		return errors.New("Policy to be updated is different than the original one")
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "StatementList":
			logger.Debug("Patch update for StatementList")
			/*newconfig should only have the next hops that have to be added or deleted*/
			newPolicy := policy.PolicyDefinitionConfig{
				Name:       origCfg.Name,
				Precedence: int(origCfg.Priority),
				MatchType:  origCfg.MatchType,
				PolicyType: origCfg.PolicyType,
			}
			newPolicy.PolicyDefinitionStatements = make([]policy.PolicyDefinitionStmtPrecedence, 0)
			newPolicy.Extensions = PolicyExtensions{}
			//logger.Debug("value = ", op[idx].Value)
			valueObjArr := []ribd.PolicyDefinitionStmtPriority{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				logger.Debug("error unmarshaling value:", err)
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug("Number of statements:", len(valueObjArr))
			for _, val := range valueObjArr {
				logger.Debug("stmtInfo: pri - ", val.Priority, " stmt: ", val.Statement)
				newPolicy.PolicyDefinitionStatements = append(newPolicy.PolicyDefinitionStatements, policy.PolicyDefinitionStmtPrecedence{
					Precedence: int(val.Priority),
					Statement:  val.Statement,
				})
			}
			switch op[idx].Op {
			case "add":
				db.UpdateAddPolicyDefinitionStmts(newPolicy)
			case "remove":
				db.UpdateRemovePolicyDefinitionStmts(newPolicy)
			default:
				logger.Err("Operation ", op[idx].Op, " not supported")
			}
		default:
			logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
			err = errors.New(fmt.Sprintln("Operation ", op[idx].Op, " not supported"))
		}
	}
	return err
}

/*
   Function to update policy definition in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyDefinitionConfigUpdate(origCfg *ribd.PolicyDefinition, newCfg *ribd.PolicyDefinition, attrset []bool, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyDefinitionConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy")
		return errors.New("Policy to be updated is different than the original one")
	}
	return err
}

func (m RIBDServer) GetBulkPolicyConditionState(fromIndex ribd.Int, rcount ribd.Int, db *policy.PolicyEngineDB) (policyConditions *ribd.PolicyConditionStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyConditionState")
	PolicyConditionsDB := db.PolicyConditionsDB
	localPolicyConditionsDB := *db.LocalPolicyConditionsDB
	var i, validCount, toIndex ribd.Int
	var tempNode []ribd.PolicyConditionState = make([]ribd.PolicyConditionState, rcount)
	var nextNode *ribd.PolicyConditionState
	var returnNodes []*ribd.PolicyConditionState
	var returnGetInfo ribd.PolicyConditionStateGetInfo
	i = 0
	policyConditions = &returnGetInfo
	more := true
	if localPolicyConditionsDB == nil {
		logger.Debug("PolicyDefinitionStmtMatchProtocolConditionGetInfo not initialized")
		return policyConditions, err
	}
	for ; ; i++ {
		if i+fromIndex >= ribd.Int(len(localPolicyConditionsDB)) {
			logger.Debug("All the policy conditions fetched")
			more = false
			break
		}
		if localPolicyConditionsDB[i+fromIndex].IsValid == false {
			logger.Debug("Invalid policy condition statement")
			continue
		}
		if validCount == rcount {
			logger.Debug("Enough policy conditions fetched")
			break
		}
		prefixNodeGet := PolicyConditionsDB.Get(localPolicyConditionsDB[i+fromIndex].Prefix)
		if prefixNodeGet != nil {
			prefixNode := prefixNodeGet.(policy.PolicyCondition)
			if strings.HasPrefix(prefixNode.Name, "__Internal") && strings.HasSuffix(prefixNode.Name, "__") {
				//this is implicitly created condition as a part of apply config
				continue
			}
			nextNode = &tempNode[validCount]
			nextNode.Name = prefixNode.Name
			nextNode.ConditionInfo = prefixNode.ConditionGetBulkInfo
			if prefixNode.PolicyStmtList != nil {
				nextNode.PolicyStmtList = make([]string, 0)
			}
			for idx := 0; idx < len(prefixNode.PolicyStmtList); idx++ {
				nextNode.PolicyStmtList = append(nextNode.PolicyStmtList, prefixNode.PolicyStmtList[idx])
			}
			toIndex = ribd.Int(prefixNode.LocalDBSliceIdx)
			if len(returnNodes) == 0 {
				returnNodes = make([]*ribd.PolicyConditionState, 0)
			}
			returnNodes = append(returnNodes, nextNode)
			validCount++
		}
	}
	policyConditions.PolicyConditionStateList = returnNodes
	policyConditions.StartIdx = fromIndex
	policyConditions.EndIdx = toIndex + 1
	policyConditions.More = more
	policyConditions.Count = validCount
	return policyConditions, err
}

func (m RIBDServer) GetBulkPolicyStmtState(fromIndex ribd.Int, rcount ribd.Int, db *policy.PolicyEngineDB) (policyStmts *ribd.PolicyStmtStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyStmtState")
	PolicyStmtDB := db.PolicyStmtDB
	localPolicyStmtDB := *db.LocalPolicyStmtDB
	var i, validCount, toIndex ribd.Int
	var tempNode []ribd.PolicyStmtState = make([]ribd.PolicyStmtState, rcount)
	var nextNode *ribd.PolicyStmtState
	var returnNodes []*ribd.PolicyStmtState
	var returnGetInfo ribd.PolicyStmtStateGetInfo
	i = 0
	policyStmts = &returnGetInfo
	more := true
	if localPolicyStmtDB == nil {
		logger.Debug("destNetSlice not initialized")
		return policyStmts, err
	}
	for ; ; i++ {
		if i+fromIndex >= ribd.Int(len(localPolicyStmtDB)) {
			logger.Debug("All the policy statements fetched")
			more = false
			break
		}
		if localPolicyStmtDB[i+fromIndex].IsValid == false {
			logger.Debug("Invalid policy statement")
			continue
		}
		if validCount == rcount {
			logger.Debug("Enough policy statements fetched")
			break
		}
		logger.Debug("Fetching trie record for index ", i+fromIndex, " and prefix ", (localPolicyStmtDB[i+fromIndex].Prefix))
		prefixNodeGet := PolicyStmtDB.Get(localPolicyStmtDB[i+fromIndex].Prefix)
		if prefixNodeGet != nil {
			prefixNode := prefixNodeGet.(policy.PolicyStmt)
			nextNode = &tempNode[validCount]
			nextNode.Name = prefixNode.Name
			nextNode.Conditions = prefixNode.Conditions
			nextNode.Action = prefixNode.Actions[0]
			if prefixNode.SetActionsCfg != nil {
				nextNode.SetActions = make([]*ribd.PolicyAction, 0)
			}
			for _, setAction := range prefixNode.SetActionsCfg {
				nextNode.SetActions = append(nextNode.SetActions, &ribd.PolicyAction{
					Attr:                   setAction.Attr,
					Community:              setAction.Community,
					ExtendedCommunityType:  setAction.ExtendedCommunity.Type,
					ExtendedCommunityValue: setAction.ExtendedCommunity.Value,
					LocalPref:              int32(setAction.LocalPref)})
			}
			if prefixNode.PolicyList != nil {
				nextNode.PolicyList = make([]string, 0)
			}
			for idx := 0; idx < len(prefixNode.PolicyList); idx++ {
				nextNode.PolicyList = append(nextNode.PolicyList, prefixNode.PolicyList[idx])
			}
			toIndex = ribd.Int(prefixNode.LocalDBSliceIdx)
			if len(returnNodes) == 0 {
				returnNodes = make([]*ribd.PolicyStmtState, 0)
			}
			returnNodes = append(returnNodes, nextNode)
			validCount++
		}
	}
	policyStmts.PolicyStmtStateList = returnNodes
	policyStmts.StartIdx = fromIndex
	policyStmts.EndIdx = toIndex + 1
	policyStmts.More = more
	policyStmts.Count = validCount
	return policyStmts, err
}

func (m RIBDServer) GetBulkPolicyDefinitionState(fromIndex ribd.Int, rcount ribd.Int, db *policy.PolicyEngineDB) (policyStmts *ribd.PolicyDefinitionStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyDefinitionState")
	PolicyDB := db.PolicyDB
	localPolicyDB := *db.LocalPolicyDB
	var i, validCount, toIndex ribd.Int
	var tempNode []ribd.PolicyDefinitionState = make([]ribd.PolicyDefinitionState, rcount)
	var nextNode *ribd.PolicyDefinitionState
	var returnNodes []*ribd.PolicyDefinitionState
	var returnGetInfo ribd.PolicyDefinitionStateGetInfo
	i = 0
	policyStmts = &returnGetInfo
	more := true
	if localPolicyDB == nil {
		logger.Debug("LocalPolicyDB not initialized")
		return policyStmts, err
	}
	for ; ; i++ {
		logger.Debug("Fetching trie record for index %d\n", i+fromIndex)
		if i+fromIndex >= ribd.Int(len(localPolicyDB)) {
			logger.Debug("All the policies fetched")
			more = false
			break
		}
		if localPolicyDB[i+fromIndex].IsValid == false {
			logger.Debug("Invalid policy")
			continue
		}
		if validCount == rcount {
			logger.Debug("Enough policies fetched")
			break
		}
		logger.Debug("Fetching trie record for index %d and prefix %v\n", i+fromIndex, (localPolicyDB[i+fromIndex].Prefix))
		prefixNodeGet := PolicyDB.Get(localPolicyDB[i+fromIndex].Prefix)
		if prefixNodeGet != nil {
			prefixNode := prefixNodeGet.(policy.Policy)
			nextNode = &tempNode[validCount]
			nextNode.Name = prefixNode.Name
			extensions := prefixNode.Extensions.(PolicyExtensions)
			nextNode.IpPrefixList = make([]string, 0)
			for k := 0; k < len(extensions.routeList); k++ {
				nextNode.IpPrefixList = append(nextNode.IpPrefixList, extensions.routeList[k])
			}
			toIndex = ribd.Int(prefixNode.LocalDBSliceIdx)
			if len(returnNodes) == 0 {
				returnNodes = make([]*ribd.PolicyDefinitionState, 0)
			}
			returnNodes = append(returnNodes, nextNode)
			validCount++
		}
	}
	policyStmts.PolicyDefinitionStateList = returnNodes
	policyStmts.StartIdx = fromIndex
	policyStmts.EndIdx = toIndex + 1
	policyStmts.More = more
	policyStmts.Count = validCount
	return policyStmts, err
}

/*
    Function called when apply policy is called by an application
	Inputs:
	        info - type ApplyPolicyInfo - specifies the policy,
			                             Source protocol which is applying the policy,
                                          Conditions when to apply the policy,
                                          Action - what needs to be done on a hit
            apply - type bool - whether to apply the policy
*/
func (m *RIBDServer) UpdateApplyPolicyList(applyList []*ribdInt.ApplyPolicyInfo, undoList []*ribdInt.ApplyPolicyInfo, apply bool, db *policy.PolicyEngineDB) {
	logger.Debug("UpdateApplyPolicyList")
	for _, applyListInfo := range applyList {
		m.UpdateApplyPolicy(applyListInfo, apply, db)
	}
	for _, undoListInfo := range undoList {
		m.UndoApplyPolicy(undoListInfo, apply, db)
	}
}
func (m *RIBDServer) UndoApplyPolicy(info *ribdInt.ApplyPolicyInfo, apply bool, db *policy.PolicyEngineDB) {
	logger.Debug("UndoApplyPolicy with apply set to ", apply)
	var err error
	source := info.Source
	conditionName := ""
	policyName := info.Policy
	action := info.Action
	var policyAction policy.PolicyAction
	conditionNameList := make([]string, 0)

	policyDB := db.PolicyDB
	policyConditionsDB := db.PolicyConditionsDB

	nodeGet := policyDB.Get(patriciaDB.Prefix(policyName))
	if nodeGet == nil {
		logger.Err("Policy ", policyName, " not defined")
		return
	}
	node := nodeGet.(policy.Policy)
	conditions := make([]ribdInt.ConditionInfo, 0)
	for i := 0; i < len(info.Conditions); i++ {
		conditions = append(conditions, *info.Conditions[i])
	}
	logger.Debug("RIB handler UndoApplyPolicy source:", source, " policy:", policyName, " action:", action, " apply:", apply, "conditions: ")
	for j := 0; j < len(conditions); j++ {
		logger.Debug("ConditionType =  :", conditions[j].ConditionType)
		switch conditions[j].ConditionType {
		case "MatchProtocol":
			logger.Debug(conditions[j].Protocol)
			conditionName = "__InternalMatch" + conditions[j].Protocol + "__"
			ok := policyConditionsDB.Match(patriciaDB.Prefix(conditionName))
			if !ok {
				logger.Debug("condition ", conditionName, " not found")
				return
			}
		case "MatchDstIpPrefix":
		case "MatchSrcIpPrefix":
			logger.Debug("IpPrefix:", conditions[j].IpPrefix, "MasklengthRange:", conditions[j].MasklengthRange)
		default:
			logger.Err("Invalid condition type:", conditions[j].ConditionType)
			return
		}
		if err == nil {
			logger.Debug("Adding condition ", conditionName, " to conditionNameList")
			conditionNameList = append(conditionNameList, conditionName)
		}
	}
	switch action {
	case "Redistribution":
		logger.Debug("Setting up Redistribution action map")
		redistributeActionInfo := policy.RedistributeActionInfo{false, source}
		policyAction = policy.PolicyAction{Name: action, ActionType: policyCommonDefs.PolicyActionTypeRouteRedistribute, ActionInfo: redistributeActionInfo}
		break
	default:
		logger.Debug("Action ", action, "currently a no-op")
		return
	}
	/*
	   Call the policy library updateApplyPolicy function
	*/
	logger.Debug("Calling undo applypolicy with conditionNameList: ", conditionNameList)
	db.UpdateUndoApplyPolicy(policy.ApplyPolicyInfo{node, policyAction, conditionNameList}, apply)
	return
}
func (m RIBDServer) UpdateApplyPolicy(info *ribdInt.ApplyPolicyInfo, apply bool, db *policy.PolicyEngineDB) {
	logger.Debug("UpdateApplyPolicy with apply set to ", apply)
	var err error
	conditionName := ""
	source := info.Source
	policyName := info.Policy
	action := info.Action
	var policyAction policy.PolicyAction
	conditionNameList := make([]string, 0)

	policyDB := db.PolicyDB
	policyConditionsDB := db.PolicyConditionsDB

	var node policy.Policy
	node.Name = policyName
	nodeGet := policyDB.Get(patriciaDB.Prefix(policyName))
	if nodeGet == nil {
		logger.Err("Policy ", policyName, " not defined")
		//return
	} else {
		node = nodeGet.(policy.Policy)
	}
	//if apply {
	conditions := make([]ribdInt.ConditionInfo, 0)
	for i := 0; i < len(info.Conditions); i++ {
		conditions = append(conditions, *info.Conditions[i])
	}
	logger.Debug("RIB handler UpdateApplyPolicy source:", source, " policy:", policyName, " action:", action, " apply:", apply, "conditions: ")
	for j := 0; j < len(conditions); j++ {
		logger.Debug("ConditionType =  :", conditions[j].ConditionType)
		switch conditions[j].ConditionType {
		case "MatchProtocol":
			logger.Debug(conditions[j].Protocol)
			conditionName = "__InternalMatch" + conditions[j].Protocol + "__"
			ok := policyConditionsDB.Match(patriciaDB.Prefix(conditionName))
			if !ok {
				logger.Debug("Define condition ", conditionName)
				policyCondition := ribd.PolicyCondition{Name: conditionName, ConditionType: conditions[j].ConditionType, Protocol: conditions[j].Protocol}
				_, err = m.ProcessPolicyConditionConfigCreate(&policyCondition, db)
			}
		case "MatchDstIpPrefix":
		case "MatchSrcIpPrefix":
			logger.Debug("IpPrefix:", conditions[j].IpPrefix, "MasklengthRange:", conditions[j].MasklengthRange)
		default:
			logger.Err("Invalid condition type:", conditions[j].ConditionType)
			return
		}
		if err == nil {
			logger.Debug("Adding condition ", conditionName, " to conditionNameList")
			conditionNameList = append(conditionNameList, conditionName)
		}
	}
	//}
	//define Action
	switch action {
	case "Redistribution":
		logger.Debug("Setting up Redistribution action map")
		redistributeActionInfo := policy.RedistributeActionInfo{true, source}
		policyAction = policy.PolicyAction{Name: action, ActionType: policyCommonDefs.PolicyActionTypeRouteRedistribute, ActionInfo: redistributeActionInfo}
		break
	default:
		logger.Debug("Action ", action, "currently a no-op")
		return
	}
	/*
	   Call the policy library updateApplyPolicy function
	*/
	logger.Debug("Calling applypolicy with conditionNameList: ", conditionNameList)
	db.UpdateApplyPolicy(policy.ApplyPolicyInfo{node, policyAction, conditionNameList}, apply)
	return
}
