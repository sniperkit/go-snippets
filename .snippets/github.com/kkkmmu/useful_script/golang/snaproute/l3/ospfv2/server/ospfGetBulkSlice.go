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

package server

import (
	"time"
)

func (server *OSPFV2Server) InitGetBulkSliceRefresh() {
	server.GetBulkData.SliceRefreshCh = make(chan bool)
	server.GetBulkData.SliceRefreshDoneCh = make(chan bool)
	server.GetBulkData.SliceRefreshDuration = time.Duration(10) * time.Minute
}

func (server *OSPFV2Server) GetBulkSliceRefresh() {
	refreshGetBulkSliceFunc := func() {
		server.GetBulkData.SliceRefreshCh <- true
		server.logger.Info("Slice Refresh is in progress")
		<-server.GetBulkData.SliceRefreshDoneCh
		server.GetBulkData.SliceRefreshTimer.Reset(server.GetBulkData.SliceRefreshDuration)
	}
	server.GetBulkData.SliceRefreshTimer = time.AfterFunc(server.GetBulkData.SliceRefreshDuration, refreshGetBulkSliceFunc)
}
