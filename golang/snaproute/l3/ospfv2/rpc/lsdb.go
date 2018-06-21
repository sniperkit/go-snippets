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
	"l3/ospfv2/api"
	"ospfv2d"
)

func (rpcHdl *rpcServiceHandler) GetOspfv2LsdbState(LSType string, LsId, AreaId, AdvRouterId string) (*ospfv2d.Ospfv2LsdbState, error) {
	var convObj *ospfv2d.Ospfv2LsdbState
	lsType, err := convertFromRPCFmtLSType(LSType)
	if err != nil {
		return nil, err
	}
	lsId, err := convertDotNotationToUint32(LsId)
	if err != nil {
		return nil, err
	}
	areaId, err := convertDotNotationToUint32(AreaId)
	if err != nil {
		return nil, err
	}
	advRouterId, err := convertDotNotationToUint32(AdvRouterId)
	if err != nil {
		return nil, err
	}
	obj, err := api.GetOspfv2LsdbState(lsType, lsId, areaId, advRouterId)
	if err == nil {
		convObj = convertToRPCFmtOspfv2LsdbState(obj)
	}
	return convObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkOspfv2LsdbState(fromIdx, count ospfv2d.Int) (*ospfv2d.Ospfv2LsdbStateGetInfo, error) {
	var getBulkInfo ospfv2d.Ospfv2LsdbStateGetInfo
	info, err := api.GetBulkOspfv2LsdbState(int(fromIdx), int(count))
	if info == nil || err != nil {
		return &getBulkInfo, err
	}
	getBulkInfo.StartIdx = fromIdx
	getBulkInfo.EndIdx = ospfv2d.Int(info.EndIdx)
	getBulkInfo.More = info.More
	getBulkInfo.Count = ospfv2d.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkInfo.Ospfv2LsdbStateList = append(getBulkInfo.Ospfv2LsdbStateList,
			convertToRPCFmtOspfv2LsdbState(info.List[idx]))
	}
	return &getBulkInfo, err
}
