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

// ribdPolicyActionApis.go
package policy

import (
	"errors"
	"fmt"
	"strconv"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)

type RedistributeActionInfo struct {
	Redistribute               bool
	RedistributeTargetProtocol string
}

type PolicyAggregateActionInfo struct {
	GenerateASSet   bool
	SendSummaryOnly bool
}

type PolicyAction struct {
	Name              string
	ActionType        int
	ActionInfo        interface{}
	PolicyStmtList    []string
	ActionGetBulkInfo string
	LocalDBSliceIdx   int
}

type PolicyActionConfig struct {
	Name                           string
	ActionType                     string
	SetAdminDistanceValue          int
	Accept                         bool
	Reject                         bool
	RedistributeAction             string
	RedistributeTargetProtocol     string
	NetworkStatementTargetProtocol string
	GenerateASSet                  bool
	SendSummaryOnly                bool
}

func (db *PolicyEngineDB) CreatePolicyRouteDispositionAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreateRouteDispositionAction"))
	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name))
		routeDispositionAction := ""
		if cfg.Accept == true {
			routeDispositionAction = "permit"
		} else if cfg.Reject == true {
			routeDispositionAction = "deny"
		} else {
			db.Logger.Err(fmt.Sprintln("User should set either one of permit/deny to true for this action type"))
			err = errors.New("User should set either one of permit/deny to true for this action type")
			return false, err
		}
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: policyCommonDefs.PolicyActionTypeRouteDisposition, ActionInfo: routeDispositionAction, LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = routeDispositionAction
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyAdminDistanceAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyAdminDistanceAction"))
	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name, "Setting admin distance value to ", cfg.SetAdminDistanceValue))
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: policyCommonDefs.PoilcyActionTypeSetAdminDistance, ActionInfo: cfg.SetAdminDistanceValue, LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = "Set admin distance to value " + strconv.Itoa(int(cfg.SetAdminDistanceValue))
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) CreatePolicyNetworkStatementAdvertiseAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyNetworkStatementAdvertiseAction"))
	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name))
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, ActionInfo: cfg.NetworkStatementTargetProtocol, LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = "Advertise network statement to " + cfg.NetworkStatementTargetProtocol
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) CreatePolicyRedistributionAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyRedistributionAction"))

	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name))
		redistributeActionInfo := RedistributeActionInfo{RedistributeTargetProtocol: cfg.RedistributeTargetProtocol}
		if cfg.RedistributeAction == "Allow" {
			redistributeActionInfo.Redistribute = true
		} else if cfg.RedistributeAction == "Block" {
			redistributeActionInfo.Redistribute = false
		} else {
			db.Logger.Err(fmt.Sprintln("Invalid redistribute option ", cfg.RedistributeAction, " - should be either Allow/Block"))
			err = errors.New("Invalid redistribute option")
			return false, err
		}
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: policyCommonDefs.PolicyActionTypeRouteRedistribute, ActionInfo: redistributeActionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = cfg.RedistributeAction + " Redistribute to Target Protocol " + cfg.RedistributeTargetProtocol
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyAggregateAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyAggregateAction"))

	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name))
		aggregateActionInfo := PolicyAggregateActionInfo{GenerateASSet: cfg.GenerateASSet, SendSummaryOnly: cfg.SendSummaryOnly}
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: policyCommonDefs.PolicyActionTypeAggregate, ActionInfo: aggregateActionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = "Aggregate action set GenerateASSet to " +
			strconv.FormatBool(cfg.GenerateASSet) + " set SendSummaryOnly to " + strconv.FormatBool(cfg.SendSummaryOnly)
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyRIBInOutAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyRIBInAction"))

	var actionType int
	switch cfg.ActionType {
	case "RIBIn":
		actionType = policyCommonDefs.PolicyActionTypeRIBIn
	case "RIBOut":
		actionType = policyCommonDefs.PolicyActionTypeRIBOut
	default:
		db.Logger.Err(fmt.Sprintln("Unknown action type ", cfg.ActionType))
		err = errors.New("Unknown action type")
		return false, err
	}

	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyAction == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy action with name ", cfg.Name))
		ribInOutAction := ""
		if cfg.Accept == true {
			ribInOutAction = "permit"
		} else if cfg.Reject == true {
			ribInOutAction = "deny"
		} else {
			db.Logger.Err(fmt.Sprintln("User should set either one of permit/deny to true for this action type"))
			err = errors.New("User should set either one of permit/deny to true for this action type")
			return false, err
		}
		newPolicyAction := PolicyAction{Name: cfg.Name, ActionType: actionType, ActionInfo: ribInOutAction,
			LocalDBSliceIdx: (len(*db.LocalPolicyActionsDB))}
		newPolicyAction.ActionGetBulkInfo = ribInOutAction
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error inserting action in DB")
			return false, err
		}
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate action name"))
		err = errors.New("Duplicate policy action definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyAction"))
	switch cfg.ActionType {
	case "RouteDisposition":
		val, err = db.CreatePolicyRouteDispositionAction(cfg)
		break
	case "Redistribution":
		val, err = db.CreatePolicyRedistributionAction(cfg)
		break
	case "SetAdminDistance":
		val, err = db.CreatePolicyAdminDistanceAction(cfg)
		break
	case "NetworkStatementAdvertise":
		val, err = db.CreatePolicyNetworkStatementAdvertiseAction(cfg)
		break
	case "Aggregate":
		val, err = db.CreatePolicyAggregateAction(cfg)
		break
	case "RIBIn":
		val, err = db.CreatePolicyRIBInOutAction(cfg)
		break
	case "RIBOut":
		val, err = db.CreatePolicyRIBInOutAction(cfg)
		break
	default:
		db.Logger.Err(fmt.Sprintln("Unknown action type ", cfg.ActionType))
		err = errors.New("Unknown action type")
		return false, err
	}
	return val, err
}

func (db *PolicyEngineDB) DeletePolicyAction(cfg PolicyActionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("DeletePolicyAction"))
	actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if actionItem == nil {
		db.Logger.Err(fmt.Sprintln("action ", cfg.Name, "not found in the DB"))
		err = errors.New("action not found")
		return false, err
	}
	action := actionItem.(PolicyAction)
	if len(action.PolicyStmtList) != 0 {
		db.Logger.Err(fmt.Sprintln("This action is currently being used by one or more policy statements. Try deleting the stmt before deleting the action"))
		err = errors.New("This action is currently being used by one or more policy statements. Try deleting the stmt before deleting the action")
		return false, err
	}
	deleted := db.PolicyActionsDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info(fmt.Sprintln("Found and deleted actions ", cfg.Name))
		db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}
