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
	"errors"
	"models/actions"
)

func (arpdClientMgr *FSArpdClntMgr) ExecuteActionArpDeleteByIfName(cfg *actions.ArpDeleteByIfName) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpDeleteByIfName(cfg)
		retVal, err := arpdClientMgr.ClientHdl.ExecuteActionArpDeleteByIfName(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) ExecuteActionArpDeleteByIPv4Addr(cfg *actions.ArpDeleteByIPv4Addr) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpDeleteByIPv4Addr(cfg)
		retVal, err := arpdClientMgr.ClientHdl.ExecuteActionArpDeleteByIPv4Addr(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) ExecuteActionArpRefreshByIfName(cfg *actions.ArpRefreshByIfName) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpRefreshByIfName(cfg)
		retVal, err := arpdClientMgr.ClientHdl.ExecuteActionArpRefreshByIfName(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) ExecuteActionArpRefreshByIPv4Addr(cfg *actions.ArpRefreshByIPv4Addr) (bool, error) {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		convCfg := convertToThriftFromArpRefreshByIPv4Addr(cfg)
		retVal, err := arpdClientMgr.ClientHdl.ExecuteActionArpRefreshByIPv4Addr(convCfg)
		arpdMutex.Unlock()
		return retVal, err
	}
	return false, errors.New("Arpd Client Handle is nil")
}
