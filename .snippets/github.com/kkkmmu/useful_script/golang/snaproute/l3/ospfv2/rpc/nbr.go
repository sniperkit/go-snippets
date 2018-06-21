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
	"ospfv2d"
)

func (rpcHdl *rpcServiceHandler) GetOspfv2NbrState(IpAddr string, AddressLessIfIdx int32) (*ospfv2d.Ospfv2NbrState, error) {
	var convObj *ospfv2d.Ospfv2NbrState
	ipAddr, err := convertDotNotationToUint32(IpAddr)
	if err != nil {
		return nil, errors.New("Invalid IP Address")
	}
	addrLessIfIdx := uint32(AddressLessIfIdx)
	obj, err := api.GetOspfv2NbrState(ipAddr, addrLessIfIdx)
	if err == nil {
		convObj = convertToRPCFmtOspfv2NbrState(obj)
	}
	return convObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkOspfv2NbrState(fromIdx, count ospfv2d.Int) (*ospfv2d.Ospfv2NbrStateGetInfo, error) {
	var getBulkInfo ospfv2d.Ospfv2NbrStateGetInfo
	info, err := api.GetBulkOspfv2NbrState(int(fromIdx), int(count))
	if info == nil || err != nil {
		return &getBulkInfo, err
	}
	getBulkInfo.StartIdx = fromIdx
	getBulkInfo.EndIdx = ospfv2d.Int(info.EndIdx)
	getBulkInfo.More = info.More
	getBulkInfo.Count = ospfv2d.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkInfo.Ospfv2NbrStateList = append(getBulkInfo.Ospfv2NbrStateList,
			convertToRPCFmtOspfv2NbrState(info.List[idx]))
	}
	return &getBulkInfo, err
}
