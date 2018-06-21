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

package fsArpdClnt

import (
	"arpd"
	"errors"
	"models/objects"
	"utils/clntUtils/clntDefs"
	"utils/clntUtils/clntDefs/arpdClntDefs"
)

func (arpdClientMgr *FSArpdClntMgr) GetArpEntryState(IpAddr string) (*objects.ArpEntryState, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		obj, err := arpdClientMgr.ClientHdl.GetArpEntryState(IpAddr)
		arpdMutex.Unlock()
		if err != nil {
			return nil, err
		}
		retObj := convertFromThriftToArpEntryState(obj)
		return retObj, nil
	}
	return nil, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) GetBulkArpEntryState(fromIdx, count int) (*arpdClntDefs.ArpEntryStateGetInfo, error) {
	if arpdClientMgr.ClientHdl == nil {
		return nil, errors.New("Arpd Client Handle is nil")
	}
	arpdMutex.Lock()
	bulkInfo, err := arpdClientMgr.ClientHdl.GetBulkArpEntryState(arpd.Int(fromIdx), arpd.Int(count))
	arpdMutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var retObj arpdClntDefs.ArpEntryStateGetInfo
	retObj.StartIdx = int32(bulkInfo.StartIdx)
	retObj.EndIdx = int32(bulkInfo.EndIdx)
	retObj.Count = int32(bulkInfo.Count)
	retObj.More = bulkInfo.More
	retObj.ArpEntryStateList = make([]*objects.ArpEntryState, int(retObj.Count))
	for idx := 0; idx < int(retObj.Count); idx++ {
		retObj.ArpEntryStateList[idx] = convertFromThriftToArpEntryState(bulkInfo.ArpEntryStateList[idx])
	}
	return &retObj, nil
}

func (arpdClientMgr *FSArpdClntMgr) GetArpLinuxEntryState(IpAddr string) (*objects.ArpLinuxEntryState, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		obj, err := arpdClientMgr.ClientHdl.GetArpLinuxEntryState(IpAddr)
		arpdMutex.Unlock()
		if err != nil {
			return nil, err
		}
		retObj := convertFromThriftToArpLinuxEntryState(obj)
		return retObj, nil
	}
	return nil, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) GetBulkArpLinuxEntryState(fromIdx, count int) (*arpdClntDefs.ArpLinuxEntryStateGetInfo, error) {
	if arpdClientMgr.ClientHdl == nil {
		return nil, errors.New("Arpd Client Handle is nil")
	}
	arpdMutex.Lock()
	bulkInfo, err := arpdClientMgr.ClientHdl.GetBulkArpLinuxEntryState(arpd.Int(fromIdx), arpd.Int(count))
	arpdMutex.Unlock()
	if bulkInfo == nil {
		return nil, err
	}
	var retObj arpdClntDefs.ArpLinuxEntryStateGetInfo
	retObj.StartIdx = int32(bulkInfo.StartIdx)
	retObj.EndIdx = int32(bulkInfo.EndIdx)
	retObj.Count = int32(bulkInfo.Count)
	retObj.More = bulkInfo.More
	retObj.ArpLinuxEntryStateList = make([]*objects.ArpLinuxEntryState, int(retObj.Count))
	for idx := 0; idx < int(retObj.Count); idx++ {
		retObj.ArpLinuxEntryStateList[idx] = convertFromThriftToArpLinuxEntryState(bulkInfo.ArpLinuxEntryStateList[idx])
	}
	return &retObj, nil
}

func (arpdClientMgr *FSArpdClntMgr) CreateArpGlobal(cfg *objects.ArpGlobal) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpGlobal(cfg)
		retVal, err := arpdClientMgr.ClientHdl.CreateArpGlobal(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) UpdateArpGlobal(origCfg, newCfg *objects.ArpGlobal, attrset []bool, op []*clntDefs.PatchOpInfo) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convOrigCfg := convertToThriftFromArpGlobal(origCfg)
		convNewCfg := convertToThriftFromArpGlobal(newCfg)
		convOp := convertToThriftPatchOpInfo(op)
		retVal, err := arpdClientMgr.ClientHdl.UpdateArpGlobal(convOrigCfg, convNewCfg, attrset, convOp)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) DeleteArpGlobal(cfg *objects.ArpGlobal) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpGlobal(cfg)
		retVal, err := arpdClientMgr.ClientHdl.DeleteArpGlobal(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}
