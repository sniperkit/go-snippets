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

type PolicyConditionConfig struct {
	Name                                string
	ConditionType                       string
	MatchProtocolConditionInfo          string
	MatchDstIpPrefixConditionInfo       PolicyDstIpMatchPrefixSetCondition
	MatchCommunityConditionInfo         PolicyMatchCommunitySetCondition
	MatchExtendedCommunityConditionInfo PolicyMatchExtendedCommunitySetCondition
	MatchNeighborConditionInfo          string
	MatchASPathConditionInfo            PolicyMatchASPathSetCondition
	MatchLocalPrefConditionInfo         uint32
	MatchMEDConditionInfo               uint32
	//MatchNeighborConditionInfo   PolicyMatchNeighborSetCondition
	//MatchTagConditionInfo   PolicyMatchTagSetCondition
}

type PolicyCondition struct {
	Name                 string
	ConditionType        int
	ConditionInfo        interface{}
	PolicyStmtList       []string
	ConditionGetBulkInfo string
	LocalDBSliceIdx      int
}

func (db *PolicyEngineDB) CreatePolicyMatchProtocolCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchProtocolCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on protocol ", cfg.MatchProtocolConditionInfo))
		matchProto := cfg.MatchProtocolConditionInfo
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeProtocolMatch, ConditionInfo: matchProto, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match Protocol " + matchProto
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}

/*
func (db *PolicyEngineDB) CreatePolicyMatchCommunityCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchCommunityCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		var val uint32
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on community ", cfg.MatchCommunityConditionInfo))
		//check if community is a well-known community
		val, err := bgpUtils.GetCommunityValue(cfg.MatchCommunityConditionInfo)
		if err != nil {
			db.Logger.Err("GetCommunityValue return error:", err, " for ", cfg.MatchCommunityConditionInfo)
			return false, err
		}
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeCommunityMatch, ConditionInfo: val, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match Community " + cfg.MatchCommunityConditionInfo
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) CreatePolicyMatchExtendedCommunityCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchExtendedCommunityCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on extended community ", cfg.MatchExtendedCommunityConditionInfo))
		match, err := bgpUtils.EncodeExtCommunity(bgpUtils.ExtCommunity{cfg.MatchExtendedCommunityConditionInfo.Type, cfg.MatchExtendedCommunityConditionInfo.Value})
		if err != nil {
			db.Logger.Err(fmt.Sprintln("EncodeExtCommunity returned err:", err))
			return false, err
		}
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch, ConditionInfo: match, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match Extended Community " + match
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) CreatePolicyMatchASPathCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchASPathCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on AS PATH ", cfg.MatchASPathConditionInfo))
		val, err := bgpUtils.GetAsPathRegex(cfg.MatchASPathConditionInfo)
		if err != nil {
			db.Logger.Err("GetAsPathRegex failed with err:", err)
			return false, err
		}
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeASPathMatch, ConditionInfo: val, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match ASPath " + cfg.MatchASPathConditionInfo
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
*/
func (db *PolicyEngineDB) CreatePolicyMatchLocalPrefCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchLocalPrefCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on local pref ", cfg.MatchLocalPrefConditionInfo))
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeLocalPrefMatch, ConditionInfo: cfg.MatchLocalPrefConditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match LocalPref " + strconv.Itoa(int(cfg.MatchLocalPrefConditionInfo))
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) CreatePolicyMatchMEDCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchMEDCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on MED ", cfg.MatchMEDConditionInfo))
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeMEDMatch, ConditionInfo: cfg.MatchMEDConditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match MED " + strconv.Itoa(int(cfg.MatchMEDConditionInfo))
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyMatchNeighborCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchNeighborCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on neighbor ", cfg.MatchNeighborConditionInfo))
		matchNeighbor := cfg.MatchNeighborConditionInfo
		newPolicyCondition := PolicyCondition{Name: cfg.Name,
			ConditionType: policyCommonDefs.PolicyConditionTypeNeighborMatch, ConditionInfo: matchNeighbor,
			LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match Neighbor " + matchNeighbor
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) ValidateConditionConfigCreate(inCfg PolicyConditionConfig) (err error) {
	db.Logger.Info(fmt.Sprintln("ValidateConditionConfigCreate"))
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition != nil {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return err
	}
	switch inCfg.ConditionType {
	case "MatchProtocol":
		break
	case "MatchDstIpPrefix":
		cfg := inCfg.MatchDstIpPrefixConditionInfo
		db.Logger.Info("ValidateConditionConfigCreate:cfg.PrefixSet:", cfg.PrefixSet, " len(cfg.PrefixSet):", len(cfg.PrefixSet))
		if len(cfg.PrefixSet) == 0 && len(cfg.Prefix.IpPrefix) == 0 {
			db.Logger.Err(fmt.Sprintln("Empty prefix set/nil prefix"))
			err = errors.New("Empty prefix set/nil prefix")
			return err
		}
		if len(cfg.PrefixSet) != 0 && len(cfg.Prefix.IpPrefix) != 0 {
			db.Logger.Err(fmt.Sprintln("Cannot provide both prefix set and individual prefix"))
			err = errors.New("Cannot provide both prefix set and individual prefix")
			return err
		}
		if len(cfg.Prefix.IpPrefix) != 0 {
			_, err = netUtils.GetNetworkPrefixFromCIDR(cfg.Prefix.IpPrefix)
			if err != nil {
				db.Logger.Err(fmt.Sprintln("ipPrefix invalid "))
				return errors.New("ipPrefix invalid")
			}
			if cfg.Prefix.MasklengthRange == "exact" {
			} else {
				maskList := strings.Split(cfg.Prefix.MasklengthRange, "-")
				if len(maskList) != 2 {
					db.Logger.Err(fmt.Sprintln("Invalid masklength range"))
					return errors.New("Invalid masklength range")
				}
				_, err = strconv.Atoi(maskList[0])
				if err != nil {
					db.Logger.Err(fmt.Sprintln("lowRange mask not valid"))
					return errors.New("lowRange mask not valid")
				}
				_, err = strconv.Atoi(maskList[1])
				if err != nil {
					db.Logger.Err(fmt.Sprintln("highRange mask not valid"))
					return errors.New("highRange mask not valid")
				}
			}
		}
	case "MatchNeighbor":
		break
	case "MatchCommunity":
	case "MatchExtendedCommunity":
	case "MatchLocalPref":
	case "MatchMED":
	case "MatchASPath":
		break
	default:
		db.Logger.Err(fmt.Sprintln("Unknown condition type ", inCfg.ConditionType))
		err = errors.New("Unknown condition type")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyCondition"))
	err = db.ValidateConditionConfigCreate(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy condition creation with err:", err)
		return false, err
	}
	switch cfg.ConditionType {
	case "MatchDstIpPrefix":
		val, err = db.CreatePolicyDstIpMatchPrefixSetCondition(cfg)
		break
	case "MatchProtocol":
		val, err = db.CreatePolicyMatchProtocolCondition(cfg)
		break
	case "MatchCommunity":
		val, err = db.CreatePolicyMatchCommunitySetCondition(cfg)
		break
	case "MatchLocalPref":
		val, err = db.CreatePolicyMatchLocalPrefCondition(cfg)
		break
	case "MatchMED":
		val, err = db.CreatePolicyMatchMEDCondition(cfg)
		break
	case "MatchASPath":
		val, err = db.CreatePolicyMatchASPathSetCondition(cfg)
		break
	case "MatchExtendedCommunity":
		val, err = db.CreatePolicyMatchExtendedCommunitySetCondition(cfg)
		break
	case "MatchNeighbor":
		val, err = db.CreatePolicyMatchNeighborCondition(cfg)
		break

	default:
		db.Logger.Err(fmt.Sprintln("Unknown condition type ", cfg.ConditionType))
		err = errors.New("Unknown condition type")
		return false, err
	}
	return val, err
}
func (db *PolicyEngineDB) ValidateConditionConfigDelete(cfg PolicyConditionConfig) (err error) {
	conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if conditionItem == nil {
		db.Logger.Err(fmt.Sprintln("Condition ", cfg.Name, "not found in the DB"))
		err = errors.New("Condition not found")
		return err
	}
	condition := conditionItem.(PolicyCondition)
	if len(condition.PolicyStmtList) != 0 {
		db.Logger.Err(fmt.Sprintln("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition"))
		err = errors.New("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition")
		return err
	}
	return nil
}
func (db *PolicyEngineDB) DeletePolicyCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("DeletePolicyCondition"))
	err = db.ValidateConditionConfigDelete(cfg)
	if err != nil {
		db.Logger.Err("Validation failed for policy condition deletion with err:", err)
		return false, err
	}
	conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if conditionItem == nil {
		db.Logger.Err(fmt.Sprintln("Condition ", cfg.Name, "not found in the DB"))
		err = errors.New("Condition not found")
		return false, err
	}
	condition := conditionItem.(PolicyCondition)
	if len(condition.PolicyStmtList) != 0 {
		db.Logger.Err(fmt.Sprintln("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition"))
		err = errors.New("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition")
		return false, err
	}
	deleted := db.PolicyConditionsDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info(fmt.Sprintln("Found and deleted condition ", cfg.Name))
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
		if condition.ConditionType == policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch {
			conditionInfo := condition.ConditionInfo.(MatchPrefixConditionInfo)
			if len(conditionInfo.PrefixSet) != 0 {
				err = db.UpdatePrefixSet(condition, conditionInfo.PrefixSet, del)
				if err != nil {
					db.Logger.Info("UpdatePrefixSet returned err ", err)
					err = errors.New("Error with UpdatePrefixSet")
					return false, err
				}
			}
		}
	}
	return true, err
}

func (db *PolicyEngineDB) UpdatePolicyCondition(cfg PolicyConditionConfig, attr string) (err error) {
	func_msg := "UpdatePolicyCondition for " + cfg.Name
	db.Logger.Debug(func_msg)
	conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if conditionItem == nil {
		db.Logger.Err(fmt.Sprintln("Condition ", cfg.Name, "not found in the DB"))
		err = errors.New("Condition not found")
		return err
	}
	//	condition := conditionItem.(PolicyCondition)
	switch attr {
	case "ConditionType":
	case "Protocol":
	case "IpPrefix":
	case "MaskLengthRange":
	case "PrefixSet":
	}
	return err
}
