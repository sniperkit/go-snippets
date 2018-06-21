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

type PolicyASPathSetConfig struct {
	Name       string
	ASPathList []string
}
type PolicyASPathSet struct {
	Name                string
	ASPathList          []string
	PolicyConditionList []string
	MatchInfoList       []MatchASPathConditionInfo
	LocalDBSliceIdx     int
}
type PolicyMatchASPathSetCondition struct {
	ASPath    string
	ASPathSet string
}
type MatchASPathConditionInfo struct {
	UseSet bool
	Set    string
	ASPath interface{}
}

func (db *PolicyEngineDB) UpdateASPathSet(condition PolicyCondition, setName string, op int) (err error) {
	db.Logger.Info("UpdateASPathSet for communityset ", setName)
	var i int
	item := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(setName))
	if item == nil {
		db.Logger.Info("ASPath set ", setName, " not defined")
		err = errors.New("ASPath set not defined")
		return err
	}
	set := item.(PolicyASPathSet)
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
	db.PolicyASPathSetDB.Set(patriciaDB.Prefix(set.Name), set)
	return err
}

func (db *PolicyEngineDB) ValidatePolicyASPathSetCreate(inCfg PolicyASPathSetConfig) (err error) {
	db.Logger.Info("ValidatePolicyASPathSetCreate")
	set := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(inCfg.Name))
	if set != nil {
		db.Logger.Err("Duplicate Condition name")
		err = errors.New("Duplicate policy ASPathSet definition")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyASPathSet(cfg PolicyASPathSetConfig) (val bool, err error) {
	db.Logger.Info("PolicyEngineDB CreatePolicyASPathSet :", cfg.Name)
	policyASPathSet := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyASPathSet == nil {
		db.Logger.Info("Defining a new policy aspath set with name ", cfg.Name)
		list := make([]string, 0)
		matchInfoList := make([]MatchASPathConditionInfo, 0)
		db.Logger.Info("cfg.ASPathList:", cfg.ASPathList)
		for _, v := range cfg.ASPathList {
			db.Logger.Info("range over cfg.ASPathList, current value:", v)
			list = append(list, v)
			var conditionInfo MatchASPathConditionInfo
			conditionInfo.UseSet = false
			val, err := bgpUtils.GetAsPathRegex(v)
			if err != nil {
				db.Logger.Err("GetAsPathRegex failed with err:", err)
				return false, err
			}
			conditionInfo.ASPath = val
			matchInfoList = append(matchInfoList, conditionInfo)
		}
		db.Logger.Info("insert ASPath set with ASPathList:", list, " matchInfoList:", matchInfoList)
		if ok := db.PolicyASPathSetDB.Insert(patriciaDB.Prefix(cfg.Name), PolicyASPathSet{Name: cfg.Name, ASPathList: list, MatchInfoList: matchInfoList}); ok != true {
			db.Logger.Info(" return value not ok")
			err = errors.New("Error creating policy aspath set in the DB")
			return false, err
		}
		db.LocalPolicyASPathSetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate policy aspath set"))
		err = errors.New("Duplicate policy policy aspath set definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) ValidatePolicyASPathSetDelete(cfg PolicyASPathSetConfig) (err error) {
	item := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("ASPath Set ", cfg.Name, "not found in the DB")
		err = errors.New("ASPath Set not found")
		return err
	}
	set := item.(PolicyASPathSet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This aspath set is currently being used by a policy condition. Try deleting the condition before deleting the aspath set")
		err = errors.New("This aspath set is currently being used by a policy condition. Try deleting the condition before deleting the aspath set")
		return err
	}
	return nil
}
func (db *PolicyEngineDB) DeletePolicyASPathSet(cfg PolicyASPathSetConfig) (val bool, err error) {
	db.Logger.Info("DeletePolicyASPathSet")
	err = db.ValidatePolicyASPathSetDelete(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy aspath set deletion with err:", err)
		return false, err
	}
	item := db.PolicyASPathSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("aspath set ", cfg.Name, "not found in the DB")
		err = errors.New("ASPath set not found")
		return false, err
	}
	set := item.(PolicyASPathSet)
	if len(set.PolicyConditionList) != 0 {
		db.Logger.Err("This aspath set is currently being used by a policy condition. Try deleting the condition before deleting the aspath set")
		err = errors.New("This aspath set is currently being used by a policy condition. Try deleting the condition before deleting the aspath set")
		return false, err
	}
	deleted := db.PolicyASPathSetDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info("Found and deleted aspath set ", cfg.Name)
		db.LocalPolicyASPathSetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyMatchASPathSetCondition(inCfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchASPathSetCondition"))
	cfg := inCfg.MatchASPathConditionInfo
	var conditionInfo MatchASPathConditionInfo
	var conditionGetBulkInfo string
	if len(cfg.ASPathSet) == 0 && len(cfg.ASPath) == 0 {
		db.Logger.Err(fmt.Sprintln("Empty aspath set/nil aspath"))
		err = errors.New("Empty aspath set/nil aspath")
		return false, err
	}
	if len(cfg.ASPathSet) != 0 && len(cfg.ASPath) != 0 {
		db.Logger.Err(fmt.Sprintln("Cannot provide both aspath set and individual ASPath"))
		err = errors.New("Cannot provide both aspath set and individual aspath")
		return false, err
	}
	if len(cfg.ASPath) != 0 {
		conditionGetBulkInfo = "match ASPath " + cfg.ASPath
		conditionInfo.UseSet = false
		//check if community is a well-known community
		val, err := bgpUtils.GetAsPathRegex(cfg.ASPath)
		if err != nil {
			db.Logger.Err("GetAsPathRegex failed with err:", err)
			return false, err
		}
		conditionInfo.ASPath = val
	} else if len(cfg.ASPathSet) != 0 {
		conditionInfo.UseSet = true
		conditionInfo.Set = cfg.ASPathSet
		conditionGetBulkInfo = "match ASPath set " + cfg.ASPathSet
	}
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", inCfg.Name))
		newPolicyCondition := PolicyCondition{Name: inCfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeASPathMatch, ConditionInfo: conditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo
		if len(cfg.ASPathSet) != 0 {
			db.Logger.Info("Policy Condition has ", cfg.ASPathSet, " aspath set")
			err = db.UpdateASPathSet(newPolicyCondition, cfg.ASPathSet, add)
			if err != nil {
				db.Logger.Info("UpdateASPathSet returned err ", err)
				err = errors.New("Error with UpdateASPathSet")
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
