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

// PolicyConditionApis.go
package policy

import (
	"errors"
	"fmt"
	bgpUtils "utils/bgpUtils"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)

type PolicyCommunitySetConfig struct {
	Name          string
	CommunityList []string
}
type PolicyCommunitySet struct {
	Name                string
	CommunityList       []string
	PolicyConditionList []string
	MatchInfoList       []MatchCommunityConditionInfo
	LocalDBSliceIdx     int
}
type PolicyMatchCommunitySetCondition struct {
	Community    string
	CommunitySet string
}
type MatchCommunityConditionInfo struct {
	UseSet    bool
	Set       string
	Community uint32
}

func (db *PolicyEngineDB) UpdateCommunitySet(condition PolicyCondition, setName string, op int) (err error) {
	db.Logger.Info("UpdateCommunitySet for communityset ", setName)
	var i int
	item := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(setName))
	if item == nil {
		db.Logger.Info("Community set ", setName, " not defined")
		err = errors.New("Community set not defined")
		return err
	}
	set := item.(PolicyCommunitySet)
	if set.PolicyConditionList == nil {
		if op == del {
			db.Logger.Info("set.PolicyConditionList nil")
			return err
		}
		set.PolicyConditionList = make([]string, 0)
	}
	if op == add {
		set.PolicyConditionList = append(set.PolicyConditionList, condition.Name)
	}
	found := false
	if op == del {
		for i = 0; i < len(set.PolicyConditionList); i++ {
			if set.PolicyConditionList[i] == condition.Name {
				db.Logger.Info("Found the condition in the policy community set table, deleting it")
				found = true
				break
			}
		}
		if found {
			set.PolicyConditionList = append(set.PolicyConditionList[:i], set.PolicyConditionList[i+1:]...)
		}
	}
	db.PolicyCommunitySetDB.Set(patriciaDB.Prefix(set.Name), set)
	return err
}

func (db *PolicyEngineDB) ValidatePolicyCommunitySetCreate(inCfg PolicyCommunitySetConfig) (err error) {
	db.Logger.Info("ValidatePolicyCommunitySetCreate")
	set := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(inCfg.Name))
	if set != nil {
		db.Logger.Err("Duplicate Condition name")
		err = errors.New("Duplicate policy CommunitySet definition")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyCommunitySet(cfg PolicyCommunitySetConfig) (val bool, err error) {
	db.Logger.Info("PolicyEngineDB CreatePolicyCommunitySet :", cfg.Name)
	policyCommunitySet := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCommunitySet == nil {
		db.Logger.Info("Defining a new policy communuty set with name ", cfg.Name)
		list := make([]string, 0)
		matchInfoList := make([]MatchCommunityConditionInfo, 0)
		db.Logger.Info("cfg.CommunityList:", cfg.CommunityList)
		for _, v := range cfg.CommunityList {
			db.Logger.Info("range over cfg.CommunityList, current value:", v)
			list = append(list, v)
			var conditionInfo MatchCommunityConditionInfo
			conditionInfo.UseSet = false
			var val uint32
			//check if community is a well-known community
			val, err := bgpUtils.GetCommunityValue(v)
			if err != nil {
				db.Logger.Err("GetCommunityValue return error:", err, " for ", v)
				return false, err
			}
			conditionInfo.Community = val
			matchInfoList = append(matchInfoList, conditionInfo)
		}
		db.Logger.Info("insert Community set with CommunityList:", list, " matchInfoList:", matchInfoList)
		if ok := db.PolicyCommunitySetDB.Insert(patriciaDB.Prefix(cfg.Name), PolicyCommunitySet{Name: cfg.Name, CommunityList: list, MatchInfoList: matchInfoList}); ok != true {
			db.Logger.Info(" return value not ok")
			err = errors.New("Error creating policy community set in the DB")
			return false, err
		}
		db.LocalPolicyCommunitySetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate policy community set"))
		err = errors.New("Duplicate policy policy community set definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) ValidatePolicyCommunitySetDelete(cfg PolicyCommunitySetConfig) (err error) {
	item := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("Community Set ", cfg.Name, "not found in the DB")
		err = errors.New("Community Set not found")
		return err
	}
	set := item.(PolicyCommunitySet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		err = errors.New("This community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		return err
	}
	return nil
}
func (db *PolicyEngineDB) DeletePolicyCommunitySet(cfg PolicyCommunitySetConfig) (val bool, err error) {
	db.Logger.Info("DeletePolicyCommunitySet")
	err = db.ValidatePolicyCommunitySetDelete(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy community set deletion with err:", err)
		return false, err
	}
	item := db.PolicyCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("community set ", cfg.Name, "not found in the DB")
		err = errors.New("Community set not found")
		return false, err
	}
	set := item.(PolicyCommunitySet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		err = errors.New("This community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		return false, err
	}
	deleted := db.PolicyCommunitySetDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info("Found and deleted community set ", cfg.Name)
		db.LocalPolicyCommunitySetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyMatchCommunitySetCondition(inCfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchCommunitySetCondition inCfg:", inCfg))
	cfg := inCfg.MatchCommunityConditionInfo
	var conditionInfo MatchCommunityConditionInfo
	var conditionGetBulkInfo string
	db.Logger.Info("CreatePolicyMatchCommunitySetCondition:cfg.CommunitySet:", cfg.CommunitySet, " cfg.Community:", cfg.Community)
	if len(cfg.CommunitySet) == 0 && len(cfg.Community) == 0 {
		db.Logger.Err(fmt.Sprintln("Empty community set/nil community"))
		err = errors.New("Empty community set/nil community")
		return false, err
	}
	if len(cfg.CommunitySet) != 0 && len(cfg.Community) != 0 {
		db.Logger.Err(fmt.Sprintln("Cannot provide both community set and individual Community"))
		err = errors.New("Cannot provide both community set and individual community")
		return false, err
	}
	if len(cfg.Community) != 0 {
		conditionGetBulkInfo = "match Community " + cfg.Community
		conditionInfo.UseSet = false
		var val uint32
		//check if community is a well-known community
		val, err := bgpUtils.GetCommunityValue(cfg.Community)
		if err != nil {
			db.Logger.Err("GetCommunityValue return error:", err, " for ", cfg.Community)
			return false, err
		}
		conditionInfo.Community = val
	} else if len(cfg.CommunitySet) != 0 {
		conditionInfo.UseSet = true
		conditionInfo.Set = cfg.CommunitySet
		conditionGetBulkInfo = "match Community set " + cfg.CommunitySet
	}
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", inCfg.Name))
		newPolicyCondition := PolicyCondition{Name: inCfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeCommunityMatch, ConditionInfo: conditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo
		if len(cfg.CommunitySet) != 0 {
			db.Logger.Info("Policy Condition has ", cfg.CommunitySet, " community set")
			err = db.UpdateCommunitySet(newPolicyCondition, cfg.CommunitySet, add)
			if err != nil {
				db.Logger.Info("UpdateCommunitySet returned err ", err)
				err = errors.New("Error with UpdateCommunitySet")
				return false, err
			}
		}
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(inCfg.Name), newPolicyCondition); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(inCfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
