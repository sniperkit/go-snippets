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
	"ribd"
	"utils/policy"
)

/*
   Function to create policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyCommunitySetConfigCreate(cfg *ribd.PolicyCommunitySet, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyCommunitySetConfigCreate:", cfg.Name)
	list := make([]string, 0)
	for _, v := range cfg.CommunityList {
		list = append(list, v)
	}
	newCfg := policy.PolicyCommunitySetConfig{Name: cfg.Name, CommunityList: list}
	val, err = db.CreatePolicyCommunitySet(newCfg)
	return val, err
}

/*
   Function to delete policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyCommunitySetConfigDelete(cfg *ribd.PolicyCommunitySet, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyCommunitySetConfigDelete: ", cfg.Name)
	newCfg := policy.PolicyCommunitySetConfig{Name: cfg.Name}
	val, err = db.DeletePolicyCommunitySet(newCfg)
	return val, err
}

/*
   Function to patch update policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyCommunitySetConfigPatchUpdate(origCfg *ribd.PolicyCommunitySet, newCfg *ribd.PolicyCommunitySet, op []*ribd.PatchOpInfo, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyCommunitySetConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy community set")
		return errors.New("Policy community set to be updated is different than the original one")
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "CommunityList":
			logger.Debug("Patch update for CommunityList")
			newPolicyObj := policy.PolicyCommunitySetConfig{
				Name: origCfg.Name,
			}
			newPolicyObj.CommunityList = make([]string, 0)
			var valueObjArr []string
			valueObjArr = make([]string, 0)
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				//logger.Debug("error unmarshaling value:", err))
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug("Number of communities:", len(valueObjArr))
			for _, val := range valueObjArr {
				newPolicyObj.CommunityList = append(newPolicyObj.CommunityList, val)
			}
			switch op[idx].Op {
			case "add":
				//db.UpdateAddPolicyDefinition(newPolicy)
			case "remove":
				//db.UpdateRemovePolicyDefinition(newconfig)
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
   Function to update policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyCommunitySetConfigUpdate(origCfg *ribd.PolicyCommunitySet, newCfg *ribd.PolicyCommunitySet, attrset []bool, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyCommunitySetConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy community set statement")
		return errors.New("Policy community set statement to be updated is different than the original one")
	}
	return err
}

func (m RIBDServer) GetBulkPolicyCommunitySetState(fromIndex ribd.Int, rcount ribd.Int, db *policy.PolicyEngineDB) (policyCommunitySets *ribd.PolicyCommunitySetStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyCommunitySetState")
	PolicyCommunitySetDB := db.PolicyCommunitySetDB
	localPolicyCommunitySetDB := *db.LocalPolicyCommunitySetDB
	var i, validCount, toIndex ribd.Int
	var tempNode []ribd.PolicyCommunitySetState = make([]ribd.PolicyCommunitySetState, rcount)
	var nextNode *ribd.PolicyCommunitySetState
	var returnNodes []*ribd.PolicyCommunitySetState
	var returnGetInfo ribd.PolicyCommunitySetStateGetInfo
	i = 0
	policyCommunitySets = &returnGetInfo
	more := true
	if localPolicyCommunitySetDB == nil {
		logger.Debug("localPolicyCommunitySetDB not initialized")
		return policyCommunitySets, err
	}
	for ; ; i++ {
		if i+fromIndex >= ribd.Int(len(localPolicyCommunitySetDB)) {
			logger.Debug("All the policy Community sets fetched")
			more = false
			break
		}
		if localPolicyCommunitySetDB[i+fromIndex].IsValid == false {
			logger.Debug("Invalid policy Community set")
			continue
		}
		if validCount == rcount {
			logger.Debug("Enough policy Community sets fetched")
			break
		}
		prefixNodeGet := PolicyCommunitySetDB.Get(localPolicyCommunitySetDB[i+fromIndex].Prefix)
		if prefixNodeGet != nil {
			prefixNode := prefixNodeGet.(policy.PolicyCommunitySet)
			nextNode = &tempNode[validCount]
			nextNode.Name = prefixNode.Name
			nextNode.CommunityList = make([]string, 0)
			for _, prefix := range prefixNode.CommunityList {
				nextNode.CommunityList = append(nextNode.CommunityList, prefix)
			}
			if prefixNode.PolicyConditionList != nil {
				nextNode.PolicyConditionList = make([]string, 0)
			}
			for idx := 0; idx < len(prefixNode.PolicyConditionList); idx++ {
				nextNode.PolicyConditionList = append(nextNode.PolicyConditionList, prefixNode.PolicyConditionList[idx])
			}
			toIndex = ribd.Int(prefixNode.LocalDBSliceIdx)
			if len(returnNodes) == 0 {
				returnNodes = make([]*ribd.PolicyCommunitySetState, 0)
			}
			returnNodes = append(returnNodes, nextNode)
			validCount++
		}
	}
	policyCommunitySets.PolicyCommunitySetStateList = returnNodes
	policyCommunitySets.StartIdx = fromIndex
	policyCommunitySets.EndIdx = toIndex + 1
	policyCommunitySets.More = more
	policyCommunitySets.Count = validCount
	return policyCommunitySets, err
}
