//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package rpc

import (
	"errors"
	"fmt"
	"l3/ospfv2/api"
	"models/objects"
	"ospfv2d"
)

func (rpcHdl *rpcServiceHandler) restoreOspfv2IntfConfFromDB() (bool, error) {
	rpcHdl.logger.Info("Restoring Ospfv2 Intf Config From DB")
	var ospfv2Intf objects.Ospfv2Intf

	ospfIntfList, err := rpcHdl.dbHdl.GetAllObjFromDb(ospfv2Intf)
	if err != nil {
		return false, errors.New("Failed to retireve Ospfv2Intf object info from DB")
	}
	for idx := 0; idx < len(ospfIntfList); idx++ {
		dbObj := ospfIntfList[idx].(objects.Ospfv2Intf)
		obj := new(ospfv2d.Ospfv2Intf)
		objects.Convertospfv2dOspfv2IntfObjToThrift(&dbObj, obj)
		convObj, err := convertFromRPCFmtOspfv2Intf(obj)
		if err != nil {
			return false, err
		}
		ok, err := api.CreateOspfv2Intf(convObj)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (rpcHdl *rpcServiceHandler) CreateOspfv2Intf(config *ospfv2d.Ospfv2Intf) (bool, error) {
	cfg, err := convertFromRPCFmtOspfv2Intf(config)
	if err != nil {
		return false, err
	}
	rv, err := api.CreateOspfv2Intf(cfg)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) UpdateOspfv2Intf(oldConfig, newConfig *ospfv2d.Ospfv2Intf, attrset []bool, op []*ospfv2d.PatchOpInfo) (bool, error) {
	convOldCfg, err := convertFromRPCFmtOspfv2Intf(oldConfig)
	if err != nil {
		return false, err
	}
	convNewCfg, err := convertFromRPCFmtOspfv2Intf(newConfig)
	if err != nil {
		return false, err
	}
	rv, err := api.UpdateOspfv2Intf(convOldCfg, convNewCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) DeleteOspfv2Intf(config *ospfv2d.Ospfv2Intf) (bool, error) {
	cfg, err := convertFromRPCFmtOspfv2Intf(config)
	if err != nil {
		return false, err
	}
	rv, err := api.DeleteOspfv2Intf(cfg)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetOspfv2IntfState(IpAddress string, AddressLessIfIdx int32) (*ospfv2d.Ospfv2IntfState, error) {
	var convObj *ospfv2d.Ospfv2IntfState
	ipAddr, err := convertDotNotationToUint32(IpAddress)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Invalid IpAddress", err))
	}
	addrLessIfIdx := uint32(AddressLessIfIdx)
	obj, err := api.GetOspfv2IntfState(ipAddr, addrLessIfIdx)
	if err == nil {
		convObj = convertToRPCFmtOspfv2IntfState(obj)
	}
	return convObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkOspfv2IntfState(fromIdx, count ospfv2d.Int) (*ospfv2d.Ospfv2IntfStateGetInfo, error) {
	var getBulkInfo ospfv2d.Ospfv2IntfStateGetInfo
	info, err := api.GetBulkOspfv2IntfState(int(fromIdx), int(count))
	if info == nil || err != nil {
		return &getBulkInfo, err
	}
	getBulkInfo.StartIdx = fromIdx
	getBulkInfo.EndIdx = ospfv2d.Int(info.EndIdx)
	getBulkInfo.More = info.More
	getBulkInfo.Count = ospfv2d.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkInfo.Ospfv2IntfStateList = append(getBulkInfo.Ospfv2IntfStateList,
			convertToRPCFmtOspfv2IntfState(info.List[idx]))
	}
	return &getBulkInfo, err
}
