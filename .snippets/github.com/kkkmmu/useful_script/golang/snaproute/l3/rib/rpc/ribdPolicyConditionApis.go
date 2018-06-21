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

// ribdPolicyConditionApis.go
package rpc

import (
	"errors"
	defs "l3/rib/ribdCommonDefs"
	"l3/rib/server"
	"ribd"
	//"utils/policy"
)

func (m RIBDServicesHandler) CreatePolicyPrefixSet(cfg *ribd.PolicyPrefixSet) (val bool, err error) {
	logger.Debug("CreatePolicyPrefixSet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyPrefixSet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) UpdatePolicyPrefixSet(origconfig *ribd.PolicyPrefixSet, newconfig *ribd.PolicyPrefixSet, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) {
	if op == nil || len(op) == 0 {
		//update op
		logger.Info("Update op for policy prefix set definition")

	} else {
		//patch op
		logger.Info("patch op:", op, " for policy prefix set definition")
	}
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		PatchOp:          op,
		Op:               defs.UpdatePolicyPrefixSet,
	}
	return val, err
}
func (m RIBDServicesHandler) DeletePolicyPrefixSet(cfg *ribd.PolicyPrefixSet) (val bool, err error) {
	logger.Debug("DeletePolicyPrefixSet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyPrefixSet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) GetBulkPolicyPrefixSetState(fromIndex ribd.Int, count ribd.Int) (state *ribd.PolicyPrefixSetStateGetInfo, err error) {
	logger.Debug("GetBulkPolicyPrefixSetState")
	ret, err := m.server.GetBulkPolicyPrefixSetState(fromIndex, count, m.server.GlobalPolicyEngineDB)
	return ret, err
}
func (m RIBDServicesHandler) GetPolicyPrefixSetState(name string) (state *ribd.PolicyPrefixSetState, err error) {
	state = ribd.NewPolicyPrefixSetState()
	return state, err
}

func (m RIBDServicesHandler) CreatePolicyASPathSet(cfg *ribd.PolicyASPathSet) (val bool, err error) {
	logger.Debug("CreatePolicyASPathSet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyASPathSet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) UpdatePolicyASPathSet(origconfig *ribd.PolicyASPathSet, newconfig *ribd.PolicyASPathSet, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) {
	if op == nil || len(op) == 0 {
		//update op
		logger.Info("Update op for policy ASPath set definition")

	} else {
		//patch op
		logger.Info("patch op:", op, " for policy ASPath set definition")
	}
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		PatchOp:          op,
		Op:               defs.UpdatePolicyASPathSet,
	}
	return val, err
}
func (m RIBDServicesHandler) DeletePolicyASPathSet(cfg *ribd.PolicyASPathSet) (val bool, err error) {
	logger.Debug("DeletePolicyASPathSet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyASPathSet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) GetBulkPolicyASPathSetState(fromIndex ribd.Int, count ribd.Int) (state *ribd.PolicyASPathSetStateGetInfo, err error) {
	logger.Debug("GetBulkPolicyASPathSetState")
	ret, err := m.server.GetBulkPolicyASPathSetState(fromIndex, count, m.server.GlobalPolicyEngineDB)
	return ret, err
}
func (m RIBDServicesHandler) GetPolicyASPathSetState(name string) (state *ribd.PolicyASPathSetState, err error) {
	state = ribd.NewPolicyASPathSetState()
	return state, err
}

func (m RIBDServicesHandler) CreatePolicyCommunitySet(cfg *ribd.PolicyCommunitySet) (val bool, err error) {
	logger.Debug("CreatePolicyCommunitySet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyCommunitySet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) UpdatePolicyCommunitySet(origconfig *ribd.PolicyCommunitySet, newconfig *ribd.PolicyCommunitySet, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) {
	if op == nil || len(op) == 0 {
		//update op
		logger.Info("Update op for policy Community set definition")

	} else {
		//patch op
		logger.Info("patch op:", op, " for policy community set definition")
	}
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		PatchOp:          op,
		Op:               defs.UpdatePolicyCommunitySet,
	}
	return val, err
}
func (m RIBDServicesHandler) DeletePolicyCommunitySet(cfg *ribd.PolicyCommunitySet) (val bool, err error) {
	logger.Debug("DeletePolicyCommunitySet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyCommunitySet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) GetBulkPolicyCommunitySetState(fromIndex ribd.Int, count ribd.Int) (state *ribd.PolicyCommunitySetStateGetInfo, err error) {
	logger.Debug("GetBulkPolicyCommunitySetState")
	ret, err := m.server.GetBulkPolicyCommunitySetState(fromIndex, count, m.server.GlobalPolicyEngineDB)
	return ret, err
}
func (m RIBDServicesHandler) GetPolicyCommunitySetState(name string) (state *ribd.PolicyCommunitySetState, err error) {
	state = ribd.NewPolicyCommunitySetState()
	return state, err
}

func (m RIBDServicesHandler) CreatePolicyExtendedCommunitySet(cfg *ribd.PolicyExtendedCommunitySet) (val bool, err error) {
	logger.Debug("CreatePolicyExtendedCommunitySet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyExtendedCommunitySet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) UpdatePolicyExtendedCommunitySet(origconfig *ribd.PolicyExtendedCommunitySet, newconfig *ribd.PolicyExtendedCommunitySet, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) {
	if op == nil || len(op) == 0 {
		//update op
		logger.Info("Update op for policy Extended Community set definition")

	} else {
		//patch op
		logger.Info("patch op:", op, " for policy Extended community set definition")
	}
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		PatchOp:          op,
		Op:               defs.UpdatePolicyExtendedCommunitySet,
	}
	return val, err
}
func (m RIBDServicesHandler) DeletePolicyExtendedCommunitySet(cfg *ribd.PolicyExtendedCommunitySet) (val bool, err error) {
	logger.Debug("DeletePolicyExtendedCommunitySet: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyExtendedCommunitySet,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) GetBulkPolicyExtendedCommunitySetState(fromIndex ribd.Int, count ribd.Int) (state *ribd.PolicyExtendedCommunitySetStateGetInfo, err error) {
	logger.Debug("GetBulkPolicyExtendedCommunitySetState")
	ret, err := m.server.GetBulkPolicyExtendedCommunitySetState(fromIndex, count, m.server.GlobalPolicyEngineDB)
	return ret, err
}
func (m RIBDServicesHandler) GetPolicyExtendedCommunitySetState(name string) (state *ribd.PolicyExtendedCommunitySetState, err error) {
	state = ribd.NewPolicyExtendedCommunitySetState()
	return state, err
}

func (m RIBDServicesHandler) CreatePolicyCondition(cfg *ribd.PolicyCondition) (val bool, err error) {
	logger.Debug("CreatePolicyConditioncfg: ", cfg.Name, " cfg:", *cfg)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyCondition,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) DeletePolicyCondition(cfg *ribd.PolicyCondition) (val bool, err error) {
	logger.Debug("DeletePolicyConditionConfig: ", cfg.Name)
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyCondition,
	}
	err = <-m.server.PolicyConfDone
	if err == nil {
		val = true
	}
	return val, err
}
func (m RIBDServicesHandler) UpdatePolicyCondition(origconfig *ribd.PolicyCondition, newconfig *ribd.PolicyCondition, attrset []bool, op []*ribd.PatchOpInfo) (val bool, err error) {
	logger.Debug("UpdatePolicyConditionConfig:UpdatePolicyCondition: ", newconfig.Name)
	if op == nil || len(op) == 0 {
		//update op
		logger.Info("Update op for policy condition definition")

	} else {
		//patch op
		logger.Err("patch op:", op, " for policy condition definition not supported")
		return false, errors.New("Patch update for policy condition not supported")
	}
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: origconfig,
		NewConfigObject:  newconfig,
		AttrSet:          attrset,
		PatchOp:          op,
		Op:               defs.UpdatePolicyCondition,
	}
	return true, err
}
func (m RIBDServicesHandler) GetPolicyConditionState(name string) (*ribd.PolicyConditionState, error) {
	logger.Debug("Get state for Policy Condition")
	retState := ribd.NewPolicyConditionState()
	return retState, nil
}
func (m RIBDServicesHandler) GetBulkPolicyConditionState(fromIndex ribd.Int, rcount ribd.Int) (policyConditions *ribd.PolicyConditionStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Debug("GetBulkPolicyConditionState")
	ret, err := m.server.GetBulkPolicyConditionState(fromIndex, rcount, m.server.GlobalPolicyEngineDB)
	return ret, err
}
