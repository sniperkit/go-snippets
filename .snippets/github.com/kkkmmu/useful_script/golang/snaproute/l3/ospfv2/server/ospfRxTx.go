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
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import ()

func (server *OSPFV2Server) StopIntfRxTxPkt(intfKey IntfConfKey) {
	server.StopOspfRecvPkts(intfKey)
	//Nothing to stop for Tx
	server.DeinitRxPkt(intfKey)
	server.DeinitTxPkt(intfKey)
	server.logger.Info("StopIntfRxTxPkt() successfully")
}

func (server *OSPFV2Server) StartIntfRxTxPkt(intfKey IntfConfKey) {
	err := server.InitRxPkt(intfKey)
	if err != nil {
		server.logger.Err("Error: InitRxPkt()", err)
		return
	}
	err = server.InitTxPkt(intfKey)
	if err != nil {
		server.logger.Err("Error: InitTxPkt()", err)
		return
	}
	go server.StartOspfRecvPkts(intfKey)
	server.logger.Info("StartIntfRxTxPkt() successfully")
}
