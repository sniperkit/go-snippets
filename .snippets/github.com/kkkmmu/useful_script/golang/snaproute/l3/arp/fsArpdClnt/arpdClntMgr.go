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
	"fmt"
	"utils/clntUtils/clntIntfs"
	"utils/ipcutils"
	"utils/logging"
)

type FSArpdClntMgr struct {
	ipcutils.IPCClientBase
	NCtrlCh   chan bool
	ClientHdl *arpd.ARPDServicesClient
}

var Logger logging.LoggerIntf

func NewArpdClntInit(clntInitParams *clntIntfs.BaseClntInitParams) (*FSArpdClntMgr, error) {
	var err error
	fsArpdClntMgr := new(FSArpdClntMgr)
	Logger = clntInitParams.Logger

	err = fsArpdClntMgr.GetArpdThriftClientHdl(clntInitParams)
	if fsArpdClntMgr.ClientHdl == nil || err != nil {
		Logger.Err("Unable Initialize Arpd Client", err)
		return nil, errors.New(fmt.Sprintln("Unable Initialize Arpd Client", err))
	}

	if clntInitParams.NHdl != nil {
		err = fsArpdClntMgr.InitFSArpdSubscriber(clntInitParams)
		if err != nil {
			Logger.Err("Unable Initialize Arpd Client", err)
			fsArpdClntMgr.DeinitArpdThriftClientHdl()
			return nil, errors.New(fmt.Sprintln("Unable Initialize Arpd Client", err))
		}
	}
	return fsArpdClntMgr, nil
}

func (arpdClientMgr *FSArpdClntMgr) ArpdClntDeinit() {
	if arpdClientMgr.ClientHdl != nil {
		arpdClientMgr.DeinitArpdThriftClientHdl()
		arpdClientMgr.ClientHdl = nil
	}
	if arpdClientMgr.NCtrlCh != nil {
		arpdClientMgr.DeinitFSArpdSubscriber()
		arpdClientMgr.NCtrlCh = nil
	}
}
