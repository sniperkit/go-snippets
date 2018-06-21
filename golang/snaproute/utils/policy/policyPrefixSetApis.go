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
	"strconv"
	"strings"
	"utils/netUtils"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)

type PolicyPrefix struct {
	IpPrefix        string //CIDR eg: 1.1.1.2/24
	MasklengthRange string //exact or a specific range 21..24
}
type PolicyPrefixSetConfig struct {
	Name       string
	PrefixList []PolicyPrefix
}
type PolicyPrefixSet struct {
	Name                string
	PrefixList          []PolicyPrefix
	PolicyConditionList []string
	MatchInfoList       []MatchPrefixConditionInfo
	LocalDBSliceIdx     int
}
type PolicyDstIpMatchPrefixSetCondition struct {
	PrefixSet string
	Prefix    PolicyPrefix
}
type MatchPrefixConditionInfo struct {
	UsePrefixSet bool
	PrefixSet    string
	DstIpMatch   bool
	SrcIpMatch   bool
	Prefix       PolicyPrefix
	IpPrefix     patriciaDB.Prefix //network prefix
	LowRange     int
	HighRange    int
}

func (db *PolicyEngineDB) UpdatePrefixSet(condition PolicyCondition, prefixSetName string, op int) (err error) {
	db.Logger.Info("UpdatePrefixSet for prefixset ", prefixSetName)
	var i int
	item := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(prefixSetName))
	if item == nil {
		db.Logger.Info("prefix set ", prefixSetName, " not defined")
		err = errors.New("prefix set not defined")
		return err
	}
	prefixSet := item.(PolicyPrefixSet)
	if prefixSet.PolicyConditionList == nil {
		if op == del {
			db.Logger.Info("prefixSet.PolicyConditionList nil")
			return err
		}
		prefixSet.PolicyConditionList = make([]string, 0)
	}
	if op == add {
		prefixSet.PolicyConditionList = append(prefixSet.PolicyConditionList, condition.Name)
	}
	found := false
	if op == del {
		for i = 0; i < len(prefixSet.PolicyConditionList); i++ {
			if prefixSet.PolicyConditionList[i] == condition.Name {
				db.Logger.Info("Found the condition in the policy prefix set table, deleting it")
				found = true
				break
			}
		}
		if found {
			prefixSet.PolicyConditionList = append(prefixSet.PolicyConditionList[:i], prefixSet.PolicyConditionList[i+1:]...)
		}
	}
	db.PolicyPrefixSetDB.Set(patriciaDB.Prefix(prefixSet.Name), prefixSet)
	return err
}

func (db *PolicyEngineDB) ValidatePolicyPrefixSetCreate(inCfg PolicyPrefixSetConfig) (err error) {
	db.Logger.Info("ValidatePolicyPrefixSetCreate")
	policyPrefixSet := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyPrefixSet != nil {
		db.Logger.Err("Duplicate Condition name")
		err = errors.New("Duplicate policy prefixSet definition")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyPrefixSet(cfg PolicyPrefixSetConfig) (val bool, err error) {
	db.Logger.Info("PolicyEngineDB CreatePolicyPrefixSet :", cfg.Name)
	policyPrefixSet := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyPrefixSet == nil {
		db.Logger.Info("Defining a new policy prefix set with name ", cfg.Name)
		prefixList := make([]PolicyPrefix, 0)
		matchInfoList := make([]MatchPrefixConditionInfo, 0)
		db.Logger.Info("cfg.PrefixList:", cfg.PrefixList)
		for _, prefix := range cfg.PrefixList {
			db.Logger.Info("range over cfg.PrefixList, current prefix:", prefix)
			prefixList = append(prefixList, prefix)
			var conditionInfo MatchPrefixConditionInfo
			conditionInfo.HighRange = -1
			conditionInfo.LowRange = -1
			if len(prefix.IpPrefix) != 0 {
				conditionInfo.UsePrefixSet = false
				conditionInfo.Prefix.IpPrefix = prefix.IpPrefix
				conditionInfo.Prefix.MasklengthRange = prefix.MasklengthRange
				conditionInfo.IpPrefix, err = netUtils.GetNetworkPrefixFromCIDR(conditionInfo.Prefix.IpPrefix)
				if err != nil {
					db.Logger.Err(fmt.Sprintln("ipPrefix invalid "))
					return
				}
				if prefix.MasklengthRange == "exact" {
				} else {
					maskList := strings.Split(conditionInfo.Prefix.MasklengthRange, "-")
					if len(maskList) != 2 {
						db.Logger.Err(fmt.Sprintln("Invalid masklength range"))
						return
					}
					conditionInfo.LowRange, err = strconv.Atoi(maskList[0])
					if err != nil {
						db.Logger.Err(fmt.Sprintln("lowRange mask not valid"))
						return
					}
					conditionInfo.HighRange, err = strconv.Atoi(maskList[1])
					if err != nil {
						db.Logger.Err(fmt.Sprintln("highRange mask not valid"))
						return
					}
					db.Logger.Info(fmt.Sprintln("lowRange = ", conditionInfo.LowRange, " highrange = ", conditionInfo.HighRange))
				}
				matchInfoList = append(matchInfoList, conditionInfo)
			}
		}
		db.Logger.Info("insert prefix set with prefixList:", prefixList, " matchInfoList:", matchInfoList)
		if ok := db.PolicyPrefixSetDB.Insert(patriciaDB.Prefix(cfg.Name), PolicyPrefixSet{Name: cfg.Name, PrefixList: prefixList, MatchInfoList: matchInfoList}); ok != true {
			db.Logger.Info(" return value not ok")
			err = errors.New("Error creating policy prefix set in the DB")
			return false, err
		}
		db.LocalPolicyPrefixSetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate policy prefix set"))
		err = errors.New("Duplicate policy policy prefix set definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) ValidatePolicyPrefixSetDelete(cfg PolicyPrefixSetConfig) (err error) {
	item := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("Prefix Set ", cfg.Name, "not found in the DB")
		err = errors.New("Prefix Set not found")
		return err
	}
	prefixSet := item.(PolicyPrefixSet)
	if len(prefixSet.PolicyConditionList) != 0 {
		db.Logger.Err("This prefix set is currently being used by a policy condition. Try deleting the condition before deleting the prefix set")
		err = errors.New("This prefix set is currently being used by a policy condition. Try deleting the condition before deleting the prefixset")
		return err
	}
	return nil
}
func (db *PolicyEngineDB) DeletePolicyPrefixSet(cfg PolicyPrefixSetConfig) (val bool, err error) {
	db.Logger.Info("DeletePolicyPrefixSet")
	err = db.ValidatePolicyPrefixSetDelete(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy prefix set deletion with err:", err)
		return false, err
	}
	item := db.PolicyPrefixSetDB.Get(patriciaDB.Prefix(cfg.Name))
	if item == nil {
		db.Logger.Err("Prefix set ", cfg.Name, "not found in the DB")
		err = errors.New("Prefix set not found")
		return false, err
	}
	prefixSet := item.(PolicyPrefixSet)
	if len(prefixSet.PolicyConditionList) != 0 {
		db.Logger.Err("This prefix set is currently being used by a policy condition. Try deleting the condition before deleting the prefix set")
		err = errors.New("This prefix set is currently being used by a policy condition. Try deleting the condition before deleting the prefixset")
		return false, err
	}
	deleted := db.PolicyPrefixSetDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info("Found and deleted prefix set ", cfg.Name)
		db.LocalPolicyPrefixSetDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyDstIpMatchPrefixSetCondition(inCfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyDstIpMatchPrefixSetCondition"))
	cfg := inCfg.MatchDstIpPrefixConditionInfo
	var conditionInfo MatchPrefixConditionInfo
	conditionInfo.HighRange = -1
	conditionInfo.LowRange = -1
	var conditionGetBulkInfo string
	if len(cfg.PrefixSet) == 0 && len(cfg.Prefix.IpPrefix) == 0 {
		db.Logger.Err(fmt.Sprintln("Empty prefix set/nil prefix"))
		err = errors.New("Empty prefix set/nil prefix")
		return false, err
	}
	if len(cfg.PrefixSet) != 0 && len(cfg.Prefix.IpPrefix) != 0 {
		db.Logger.Err(fmt.Sprintln("Cannot provide both prefix set and individual prefix"))
		err = errors.New("Cannot provide both prefix set and individual prefix")
		return false, err
	}
	if len(cfg.Prefix.IpPrefix) != 0 {
		conditionGetBulkInfo = "match destination Prefix " + cfg.Prefix.IpPrefix + "MasklengthRange " + cfg.Prefix.MasklengthRange
		conditionInfo.UsePrefixSet = false
		conditionInfo.Prefix.IpPrefix = cfg.Prefix.IpPrefix
		conditionInfo.Prefix.MasklengthRange = cfg.Prefix.MasklengthRange
		conditionInfo.IpPrefix, err = netUtils.GetNetworkPrefixFromCIDR(conditionInfo.Prefix.IpPrefix)
		if err != nil {
			db.Logger.Err(fmt.Sprintln("ipPrefix invalid "))
			return
		}
		if cfg.Prefix.MasklengthRange == "exact" {
		} else {
			maskList := strings.Split(conditionInfo.Prefix.MasklengthRange, "-")
			if len(maskList) != 2 {
				db.Logger.Err(fmt.Sprintln("Invalid masklength range"))
				return
			}
			conditionInfo.LowRange, err = strconv.Atoi(maskList[0])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("lowRange mask not valid"))
				return
			}
			conditionInfo.HighRange, err = strconv.Atoi(maskList[1])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("highRange mask not valid"))
				return
			}
			db.Logger.Info(fmt.Sprintln("lowRange = ", conditionInfo.LowRange, " highrange = ", conditionInfo.HighRange))
		}
	} else if len(cfg.PrefixSet) != 0 {
		conditionInfo.UsePrefixSet = true
		conditionInfo.PrefixSet = cfg.PrefixSet
		conditionGetBulkInfo = "match destination Prefix " + cfg.PrefixSet
	}
	conditionInfo.DstIpMatch = true
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", inCfg.Name))
		newPolicyCondition := PolicyCondition{Name: inCfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch, ConditionInfo: conditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo
		if len(cfg.PrefixSet) != 0 {
			db.Logger.Info("Policy Condition has ", cfg.PrefixSet, " prefix set")
			err = db.UpdatePrefixSet(newPolicyCondition, cfg.PrefixSet, add)
			if err != nil {
				db.Logger.Info("UpdatePrefixSet returned err ", err)
				err = errors.New("Error with UpdatePrefixSet")
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
