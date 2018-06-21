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

// This file defines all interfaces provided for L2 service
package rpc

import (
	"asicdInt"
	"asicdServices"
)

//LAG related APIs
func (svcHdlr AsicDaemonServiceHandler) CreateLag(ifName string, hashType int32, ports string) (rval int32, err error) {
	return (int32)(0), nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteLag(ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateLag(ifIndex, hashType int32, ports string) (rval int32, err error) {
	return (int32)(0), nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateLagCfgIntfList(ifName string, intfList []int32) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteLagCfgIntfList(ifName string, intfList []int32) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateLagCfgIntfList(ifName string, intfList []int32) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkLag(currMarker, count asicdInt.Int) (*asicdInt.LagGetInfo, error) {
	bulkObj := asicdInt.NewLagGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateStg(vlanList []int32) (stgId int32, err error) {
	return int32(0), nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteStg(stgId int32) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) SetPortStpState(stgId, port, stpState int32) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetPortStpState(stgId, port int32) (stpState int32, err error) {
	return stpState, err
}

func (svcHdlr AsicDaemonServiceHandler) UpdateStgVlanList(stgId int32, vlanList []int32) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) FlushFdbStgGroup(stgId, port int32) error {
	return nil
}

//Protocol Mac Addr related APIs
func (svcHdlr AsicDaemonServiceHandler) EnablePacketReception(
	macObj *asicdInt.RsvdProtocolMacConfig) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DisablePacketReception(
	macObj *asicdInt.RsvdProtocolMacConfig) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetMacTableEntryState(macAddr string) (*asicdServices.MacTableEntryState, error) {
	macTableEntryState := asicdServices.NewMacTableEntryState()
	return macTableEntryState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkMacTableEntryState(currMarker, count asicdServices.Int) (*asicdServices.MacTableEntryStateGetInfo, error) {
	bulkObj := asicdServices.NewMacTableEntryStateGetInfo()
	return bulkObj, nil
}
