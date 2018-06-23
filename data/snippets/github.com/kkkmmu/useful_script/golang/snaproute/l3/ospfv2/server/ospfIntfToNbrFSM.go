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

import ()

func (server *OSPFV2Server) SendHelloEventMsg(msg NbrHelloEventMsg) {
	server.logger.Debug("Sending msg to Nbr State Machine", msg)
	server.MessagingChData.IntfToNbrFSMChData.NbrHelloEventCh <- msg
}

func (server *OSPFV2Server) SendDeleteNbrsMsg(key IntfConfKey) {
	msg := DeleteNbrMsg{
		IntfKey: key,
	}
	server.logger.Info("Send message to Nbr state machine to delete nbrs", msg)
	server.MessagingChData.IntfToNbrFSMChData.DeleteNbrCh <- msg
}

func (server *OSPFV2Server) SendNetworkDRChangeMsg(key IntfConfKey, oldState, newState uint8) {
	msg := NetworkDRChangeMsg{
		IntfKey:         key,
		OldIntfFSMState: oldState,
		NewIntfFSMState: newState,
	}
	server.logger.Info("Sending Network DR change message", msg)
	server.MessagingChData.IntfToNbrFSMChData.NetworkDRChangeCh <- msg
}
