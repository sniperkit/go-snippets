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
func (m RIBDServer) ProcessPolicyPrefixSetConfigCreate(cfg *ribd.PolicyPrefixSet, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyPrefixSetConfigCreate:", cfg.Name)
	prefixList := make([]policy.PolicyPrefix, 0)
	for _, prefix := range cfg.PrefixList {
		prefixList = append(prefixList, policy.PolicyPrefix{IpPrefix: prefix.Prefix, MasklengthRange: prefix.MaskLengthRange})
	}
	newCfg := policy.PolicyPrefixSetConfig{Name: cfg.Name, PrefixList: prefixList}
	val, err = db.CreatePolicyPrefixSet(newCfg)
	return val, err
}

/*
   Function to delete policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyPrefixSetConfigDelete(cfg *ribd.PolicyPrefixSet, db *policy.PolicyEngineDB) (val bool, err error) {
	logger.Debug("ProcessPolicyPrefixSetConfigDelete: ", cfg.Name)
	newCfg := policy.PolicyPrefixSetConfig{Name: cfg.Name}
	val, err = db.DeletePolicyPrefixSet(newCfg)
	return val, err
}

/*
   Function to patch update policy prefix set in the policyEngineDB
*/
func (m RIBDServer) ProcessPolicyPrefixSetConfigPatchUpdate(origCfg *ribd.PolicyPrefixSet, newCfg *ribd.PolicyPrefixSet, op []*ribd.PatchOpInfo, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyPrefixSetConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy prefix set")
		return errors.New("Policy prefix set to be updated is different than the original one")
	}
	for idx := 0; idx < len(op); idx++ {
		switch op[idx].Path {
		case "PrefixList":
			logger.Debug("Patch update for PrefixList")
			newPolicyObj := policy.PolicyPrefixSetConfig{
				Name: origCfg.Name,
			}
			newPolicyObj.PrefixList = make([]policy.PolicyPrefix, 0)
			valueObjArr := []ribd.PolicyPrefix{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				//logger.Debug("error unmarshaling value:", err))
				return errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			logger.Debug("Number of prefixes:", len(valueObjArr))
			for _, val := range valueObjArr {
				logger.Debug("ipPrefix - ", val.Prefix, " masklengthrange:", val.MaskLengthRange)
				newPolicyObj.PrefixList = append(newPolicyObj.PrefixList, policy.PolicyPrefix{
					IpPrefix:        val.Prefix,
					MasklengthRange: val.MaskLengthRange,
				})
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
func (m RIBDServer) ProcessPolicyPrefixSetConfigUpdate(origCfg *ribd.PolicyPrefixSet, newCfg *ribd.PolicyPrefixSet, attrset []bool, db *policy.PolicyEngineDB) (err error) {
	logger.Debug("ProcessPolicyPrefixSetConfigUpdate:", origCfg.Name)
	if origCfg.Name != newCfg.Name {
		logger.Err("Update for a different policy prefix set statement")
		return errors.New("Policy prefix set statement to be updated is different than the original one")
	}
	return err
}

func (m RIBDServer) GetBulkPolicyPrefixSetState(fromIndex ribd.Int, rcount ribd.Int, db *policy.PolicyEngineDB) (policyPrefixSets *ribd.PolicyPrefixSetStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyPrefixSetState")
	PolicyPrefixSetDB := db.PolicyPrefixSetDB
	localPolicyPrefixSetDB := *db.LocalPolicyPrefixSetDB
	var i, validCount, toIndex ribd.Int
	var tempNode []ribd.PolicyPrefixSetState = make([]ribd.PolicyPrefixSetState, rcount)
	var nextNode *ribd.PolicyPrefixSetState
	var returnNodes []*ribd.PolicyPrefixSetState
	var returnGetInfo ribd.PolicyPrefixSetStateGetInfo
	i = 0
	policyPrefixSets = &returnGetInfo
	more := true
	if localPolicyPrefixSetDB == nil {
		logger.Debug("localPolicyPrefixSetDB not initialized")
		return policyPrefixSets, err
	}
	for ; ; i++ {
		if i+fromIndex >= ribd.Int(len(localPolicyPrefixSetDB)) {
			logger.Debug("All the policy prefix sets fetched")
			more = false
			break
		}
		if localPolicyPrefixSetDB[i+fromIndex].IsValid == false {
			logger.Debug("Invalid policy prefix set")
			continue
		}
		if validCount == rcount {
			logger.Debug("Enough policy prefix sets fetched")
			break
		}
		prefixNodeGet := PolicyPrefixSetDB.Get(localPolicyPrefixSetDB[i+fromIndex].Prefix)
		if prefixNodeGet != nil {
			prefixNode := prefixNodeGet.(policy.PolicyPrefixSet)
			nextNode = &tempNode[validCount]
			nextNode.Name = prefixNode.Name
			nextNode.PrefixList = make([]*ribd.PolicyPrefix, 0)
			for _, prefix := range prefixNode.PrefixList {
				nextNode.PrefixList = append(nextNode.PrefixList, &ribd.PolicyPrefix{prefix.IpPrefix, prefix.MasklengthRange})
			}
			logger.Info("len(nextNode.PrefixList):", len(nextNode.PrefixList), " len(prefixNode.PrefixList:", len(prefixNode.PrefixList))
			if prefixNode.PolicyConditionList != nil {
				nextNode.PolicyConditionList = make([]string, 0)
			}
			for idx := 0; idx < len(prefixNode.PolicyConditionList); idx++ {
				nextNode.PolicyConditionList = append(nextNode.PolicyConditionList, prefixNode.PolicyConditionList[idx])
			}
			toIndex = ribd.Int(prefixNode.LocalDBSliceIdx)
			if len(returnNodes) == 0 {
				returnNodes = make([]*ribd.PolicyPrefixSetState, 0)
			}
			returnNodes = append(returnNodes, nextNode)
			validCount++
		}
	}
	policyPrefixSets.PolicyPrefixSetStateList = returnNodes
	policyPrefixSets.StartIdx = fromIndex
	policyPrefixSets.EndIdx = toIndex + 1
	policyPrefixSets.More = more
	policyPrefixSets.Count = validCount
	return policyPrefixSets, err
}
