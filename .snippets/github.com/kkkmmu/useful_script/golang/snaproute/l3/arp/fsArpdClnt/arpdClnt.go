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
	"strconv"
	"sync"
	"time"
	"utils/cfgParser"
	"utils/clntUtils/clntIntfs"
	"utils/ipcutils"
)

var arpdMutex *sync.Mutex = &sync.Mutex{}

func (arpdClientMgr *FSArpdClntMgr) GetArpdThriftClientHdl(clntInitParams *clntIntfs.BaseClntInitParams) error {
	arpdClientMgr.Name = "arpd"
	port, err := cfgParser.GetDmnPortFromClientJson(arpdClientMgr.Name, clntInitParams.ParamsFile)
	if err != nil {
		Logger.Err("Error opening client connection for arpd", err)
		return err
	}
	Logger.Debug("found arpd at port", port)
	arpdClientMgr.Address = "localhost:" + strconv.Itoa(port)
	arpdClientMgr.TTransport, arpdClientMgr.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(arpdClientMgr.Address)
	if err != nil {
		Logger.Err("Failed to connect to Arpd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			arpdClientMgr.TTransport, arpdClientMgr.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(arpdClientMgr.Address)
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				Logger.Err("Still can't connect to Arpd, retrying..")
			}
		}

	}
	Logger.Info("Connected to Arpd")
	arpdClientMgr.ClientHdl = arpd.NewARPDServicesClientFactory(arpdClientMgr.TTransport, arpdClientMgr.PtrProtocolFactory)
	return nil
}

func (arpdClientMgr *FSArpdClntMgr) DeinitArpdThriftClientHdl() {
	arpdClientMgr.CloseIPCHandles()
}
