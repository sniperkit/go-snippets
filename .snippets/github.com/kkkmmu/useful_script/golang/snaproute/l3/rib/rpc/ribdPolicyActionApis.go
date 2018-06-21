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
package rpc

import (
	defs "l3/rib/ribdCommonDefs"
	"l3/rib/server"
	"ribdInt"
)

func (m RIBDServicesHandler) CreatePolicyAction(cfg *ribdInt.PolicyAction) (val bool, err error) {
	logger.Info("CreatePolicyAction")
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.AddPolicyAction,
	}
	return true, err
}

func (m RIBDServicesHandler) DeletePolicyAction(cfg *ribdInt.PolicyAction) (val bool, err error) {
	logger.Info("CreatePolicyAction")
	m.server.PolicyConfCh <- server.RIBdServerConfig{
		OrigConfigObject: cfg,
		Op:               defs.DelPolicyAction,
	}
	return true, err
}

func (m RIBDServicesHandler) UpdatePolicyAction(origconfig *ribdInt.PolicyAction, newconfig *ribdInt.PolicyAction, attrset []bool, op []*ribdInt.PatchOpInfo) (val bool, err error) {
	logger.Info("UpdatePolicyAction")
	return true, err
}

/*func (m RIBDServicesHandler) GetPolicyActionState(name string) (*ribdInt.PolicyActionState, error) {
	logger.Info("Get state for Policy Action")
	retState := ribd.NewPolicyActionState()
	return retState, nil
}
func (m RIBDServicesHandler) GetBulkPolicyActionState(fromIndex ribd.Int, rcount ribd.Int) (policyActions *ribdInt.PolicyActionStateGetInfo, err error) { //(routes []*ribd.Routes, err error) {
	logger.Info(fmt.Sprintln("GetBulkPolicyActionState"))
	policyActions,err = m.server.GetBulkPolicyActionState(fromIndex,rcount)
	return policyActions, err
}*/
