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
	//"l3/ospfv2/api"
	"ospfv2d"
)

func (rpcHdl *rpcServiceHandler) GetOspfv2RouteState(DestId string, AddrMask string, DestType string) (*ospfv2d.Ospfv2RouteState, error) {
	/*
		var convObj *ospfv2d.Ospfv2RouteState
		//TODO
		destId := uint32(0)
		addrMask := uint32(0)
		destType := uint32(0)
		obj, err := api.GetOspfv2RouteState(destId, addrMask, destType)
		if err == nil {
			convObj = convertToRPCFmtOspfv2RouteState(obj)
		}
		return convObj, err
	*/
	return nil, errors.New("This call should not come to Ospfv2 Daemon")
}

func (rpcHdl *rpcServiceHandler) GetBulkOspfv2RouteState(fromIdx, count ospfv2d.Int) (*ospfv2d.Ospfv2RouteStateGetInfo, error) {
	/*
		var getBulkInfo ospfv2d.Ospfv2RouteStateGetInfo
		info, err := api.GetBulkOspfv2RouteState(int(fromIdx), int(count))
		getBulkInfo.StartIdx = fromIdx
		getBulkInfo.EndIdx = ospfv2d.Int(info.EndIdx)
		getBulkInfo.More = info.More
		getBulkInfo.Count = ospfv2d.Int(len(info.List))
		for idx := 0; idx < len(info.List); idx++ {
			getBulkInfo.Ospfv2RouteStateList = append(getBulkInfo.Ospfv2RouteStateList,
				convertToRPCFmtOspfv2RouteState(info.List[idx]))
		}
		return &getBulkInfo, err
	*/
	return nil, errors.New("This call should not come to Ospfv2 Daemon")
}
