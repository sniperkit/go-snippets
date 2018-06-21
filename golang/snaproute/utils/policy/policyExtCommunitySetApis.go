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

type PolicyExtendedCommunityInfo struct {
	Type  string
	Value string
}

type PolicyExtendedCommunitySetConfig struct {
	Name                  string
	ExtendedCommunityList []PolicyExtendedCommunityInfo
}
type PolicyExtendedCommunitySet struct {
	Name                  string
	ExtendedCommunityList []PolicyExtendedCommunityInfo
	PolicyConditionList   []string
	MatchInfoList         []MatchExtendedCommunityConditionInfo
	LocalDBSliceIdx       int
}
type PolicyMatchExtendedCommunitySetCondition struct {
	ExtendedCommunity    PolicyExtendedCommunityInfo
	ExtendedCommunitySet string
}
type MatchExtendedCommunityConditionInfo struct {
	UseSet            bool
	Set               string
	ExtendedCommunity uint64
}

func (db *PolicyEngineDB) UpdateExtendedCommunitySet(condition PolicyCondition, setName string, op int) (err error) {
	db.Logger.Info("UpdateExtendedCommunitySet for communityset ", setName)
	var i int
	item := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(setName))
	if item == nil {
		db.Logger.Info("ExtendedCommunity set ", setName, " not defined")
		err = errors.New("ExtendedCommunity set not defined")
		return err
	}
	set := item.(PolicyExtendedCommunitySet)
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
				db.Logger.Info("Found the condition in the policy extended community set table, deleting it")
				found = true
				break
			}
		}
		if found {
			set.PolicyConditionList = append(set.PolicyConditionList[:i], set.PolicyConditionList[i+1:]...)
		}
	}
	db.PolicyExtendedCommunitySetDB.Set(patriciaDB.Prefix(set.Name), set)
	return err
}

func (db *PolicyEngineDB) ValidatePolicyExtendedCommunitySetCreate(inCfg PolicyExtendedCommunitySetConfig) (err error) {
	db.Logger.Info("ValidatePolicyExtendedCommunitySetCreate")
	set := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(inCfg.Name))
	if set != nil {
		db.Logger.Err("Duplicate Condition name")
		err = errors.New("Duplicate policy ExtendedCommunitySet definition")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyExtendedCommunitySet(cfg PolicyExtendedCommunitySetConfig) (val bool, err error) {
	db.Logger.Info("PolicyEngineDB CreatePolicyExtendedCommunitySet :", cfg.Name)
	policyExtendedCommunitySet := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyExtendedCommunitySet == nil {
		db.Logger.Info("Defining a new policy extended communuty set with name ", cfg.Name)
		list := make([]PolicyExtendedCommunityInfo, 0)
		matchInfoList := make([]MatchExtendedCommunityConditionInfo, 0)
		db.Logger.Info("cfg.ExtendedCommunityList:", cfg.ExtendedCommunityList)
		for _, v := range cfg.ExtendedCommunityList {
			db.Logger.Info("range over cfg.ExtendedCommunityList, current value:", v)
			list = append(list, v)
			var conditionInfo MatchExtendedCommunityConditionInfo
			conditionInfo.UseSet = false
			match, err := bgpUtils.EncodeExtCommunity(bgpUtils.ExtCommunity{v.Type, v.Value})
			if err != nil {
				db.Logger.Err(fmt.Sprintln("EncodeExtCommunity returned err:", err))
				return false, err
			}
			conditionInfo.ExtendedCommunity = match
			matchInfoList = append(matchInfoList, conditionInfo)
		}
		db.Logger.Info("insert ExtendedCommunity set with ExtendedCommunityList:", list, " matchInfoList:", matchInfoList)
		if ok := db.PolicyExtendedCommunitySetDB.Insert(patriciaDB.Prefix(cfg.Name), PolicyExtendedCommunitySet{Name: cfg.Name, ExtendedCommunityList: list, MatchInfoList: matchInfoList}); ok != true {
			db.Logger.Info(" return value not ok")
			err = errors.New("Error creating policy extended community set in the DB")
			return false, err
		}
		db.LocalPolicyExtendedCommunitySetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate policy extended community set"))
		err = errors.New("Duplicate policy policy extended community set definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) ValidatePolicyExtendedCommunitySetDelete(cfg PolicyExtendedCommunitySetConfig) (err error) {
	item := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("ExtendedCommunity Set ", cfg.Name, "not found in the DB")
		err = errors.New("ExtendedCommunity Set not found")
		return err
	}
	set := item.(PolicyExtendedCommunitySet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This extended community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		err = errors.New("This extended community set is currently being used by a policy condition. Try deleting the condition before deleting the community set")
		return err
	}
	return nil
}
func (db *PolicyEngineDB) DeletePolicyExtendedCommunitySet(cfg PolicyExtendedCommunitySetConfig) (val bool, err error) {
	db.Logger.Info("DeletePolicyExtendedCommunitySet")
	err = db.ValidatePolicyExtendedCommunitySetDelete(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy extended community set deletion with err:", err)
		return false, err
	}
	item := db.PolicyExtendedCommunitySetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("extended community set ", cfg.Name, "not found in the DB")
		err = errors.New("ExtendedCommunity set not found")
		return false, err
	}
	set := item.(PolicyExtendedCommunitySet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This extended community set is currently being used by a policy condition. Try deleting the condition before deleting the extended community set")
		err = errors.New("This extended community set is currently being used by a policy condition. Try deleting the condition before deleting the extended community set")
		return false, err
	}
	deleted := db.PolicyExtendedCommunitySetDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info("Found and deleted extended community set ", cfg.Name)
		db.LocalPolicyExtendedCommunitySetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyMatchExtendedCommunitySetCondition(inCfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchExtendedCommunitySetCondition"))
	cfg := inCfg.MatchExtendedCommunityConditionInfo
	var conditionInfo MatchExtendedCommunityConditionInfo
	var conditionGetBulkInfo string
	if len(cfg.ExtendedCommunitySet) == 0 && len(cfg.ExtendedCommunity.Value) == 0 {
		db.Logger.Err(fmt.Sprintln("Empty extended community set/nil extended community"))
		err = errors.New("Empty extended community set/nil extended community")
		return false, err
	}
	if len(cfg.ExtendedCommunitySet) != 0 && len(cfg.ExtendedCommunity.Value) != 0 {
		db.Logger.Err(fmt.Sprintln("Cannot provide both extended community set and individual ExtendedCommunity"))
		err = errors.New("Cannot provide both extended community set and individual extended community")
		return false, err
	}
	if len(cfg.ExtendedCommunity.Value) != 0 {
		conditionGetBulkInfo = "match ExtendedCommunity " + cfg.ExtendedCommunity.Type + ":" + cfg.ExtendedCommunity.Value
		conditionInfo.UseSet = false
		match, err := bgpUtils.EncodeExtCommunity(bgpUtils.ExtCommunity{cfg.ExtendedCommunity.Type, cfg.ExtendedCommunity.Value})
		if err != nil {
			db.Logger.Err(fmt.Sprintln("EncodeExtCommunity returned err:", err))
			return false, err
		}
		conditionInfo.ExtendedCommunity = match
	} else if len(cfg.ExtendedCommunitySet) != 0 {
		conditionInfo.UseSet = true
		conditionInfo.Set = cfg.ExtendedCommunitySet
		conditionGetBulkInfo = "match ExtendedCommunity set " + cfg.ExtendedCommunitySet
	}
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", inCfg.Name))
		newPolicyCondition := PolicyCondition{Name: inCfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch, ConditionInfo: conditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo
		if len(cfg.ExtendedCommunitySet) != 0 {
			db.Logger.Info("Policy Condition has ", cfg.ExtendedCommunitySet, " extended community set")
			err = db.UpdateExtendedCommunitySet(newPolicyCondition, cfg.ExtendedCommunitySet, add)
			if err != nil {
				db.Logger.Info("UpdateExtendedCommunitySet returned err ", err)
				err = errors.New("Error with UpdateExtendedCommunitySet")
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
