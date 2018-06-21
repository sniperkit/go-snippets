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

// This file defines all interfaces provided for the Vlan service
package rpc

import (
	"asicdInt"
	"asicdServices"
)

func (svcHdlr AsicDaemonServiceHandler) CreateVlan(vlanObj *asicdServices.Vlan) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateVlan(oldVlanObj, newVlanObj *asicdServices.Vlan, attrset []bool, op []*asicdServices.PatchOpInfo) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteVlan(vlanObj *asicdServices.Vlan) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkVlan(currMarker, count asicdInt.Int) (*asicdInt.VlanGetInfo, error) {
	bulkObj := asicdInt.NewVlanGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkVlanState(currMarker, count asicdServices.Int) (*asicdServices.VlanStateGetInfo, error) {
	bulkObj := asicdServices.NewVlanStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetVlanState(vlanId int32) (*asicdServices.VlanState, error) {
	vlanState := asicdServices.NewVlanState()
	return vlanState, nil
}
