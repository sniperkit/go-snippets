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
	"l3/ospfv2/api"
	"models/objects"
	"ospfv2d"
)

func (rpcHdl *rpcServiceHandler) restoreOspfv2GlobalConfFromDB() (bool, error) {
	rpcHdl.logger.Info("Restoring Ospfv2 Global Config From DB")
	var ospfv2Gbl objects.Ospfv2Global

	ospfGblList, err := rpcHdl.dbHdl.GetAllObjFromDb(ospfv2Gbl)
	if err != nil {
		return false, errors.New("Failed to retireve Ospfv2Global object info from DB")
	}
	rpcHdl.logger.Info("ospfGblList:", ospfGblList)
	for idx := 0; idx < len(ospfGblList); idx++ {
		dbObj := ospfGblList[idx].(objects.Ospfv2Global)
		obj := new(ospfv2d.Ospfv2Global)
		objects.Convertospfv2dOspfv2GlobalObjToThrift(&dbObj, obj)
		convObj, err := convertFromRPCFmtOspfv2Global(obj)
		if err != nil {
			return false, err
		}
		ok, err := api.CreateOspfv2Global(convObj)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (rpcHdl *rpcServiceHandler) CreateOspfv2Global(config *ospfv2d.Ospfv2Global) (bool, error) {
	cfg, err := convertFromRPCFmtOspfv2Global(config)
	if err != nil {
		return false, err
	}
	rv, err := api.CreateOspfv2Global(cfg)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) UpdateOspfv2Global(oldConfig, newConfig *ospfv2d.Ospfv2Global, attrset []bool, op []*ospfv2d.PatchOpInfo) (bool, error) {
	convOldCfg, err := convertFromRPCFmtOspfv2Global(oldConfig)
	if err != nil {
		return false, err
	}
	convNewCfg, err := convertFromRPCFmtOspfv2Global(newConfig)
	if err != nil {
		return false, err
	}
	rv, err := api.UpdateOspfv2Global(convOldCfg, convNewCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) DeleteOspfv2Global(config *ospfv2d.Ospfv2Global) (bool, error) {
	cfg, err := convertFromRPCFmtOspfv2Global(config)
	if err != nil {
		return false, err
	}
	rv, err := api.DeleteOspfv2Global(cfg)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetOspfv2GlobalState(Vrf string) (*ospfv2d.Ospfv2GlobalState, error) {
	var convObj *ospfv2d.Ospfv2GlobalState
	// Need to be updated when we support Vrf
	if Vrf != "default" {
		return nil, errors.New("Unsupported Vrf")
	}
	obj, err := api.GetOspfv2GlobalState(Vrf)
	if err == nil {
		convObj = convertToRPCFmtOspfv2GlobalState(obj)
	}
	return convObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkOspfv2GlobalState(fromIdx, count ospfv2d.Int) (*ospfv2d.Ospfv2GlobalStateGetInfo, error) {
	var getBulkInfo ospfv2d.Ospfv2GlobalStateGetInfo
	info, err := api.GetBulkOspfv2GlobalState(int(fromIdx), int(count))
	if info == nil || err != nil {
		return &getBulkInfo, err
	}
	getBulkInfo.StartIdx = fromIdx
	getBulkInfo.EndIdx = ospfv2d.Int(info.EndIdx)
	getBulkInfo.More = info.More
	getBulkInfo.Count = ospfv2d.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkInfo.Ospfv2GlobalStateList = append(getBulkInfo.Ospfv2GlobalStateList,
			convertToRPCFmtOspfv2GlobalState(info.List[idx]))
	}
	return &getBulkInfo, err
}
