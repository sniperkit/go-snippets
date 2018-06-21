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
	"encoding/json"
	nanomsg "github.com/op/go-nanomsg"
	"l3/rib/ribdCommonDefs"
)

func (server *BFDServer) CreateRIBdSubscriber() {
	server.logger.Info("Listen for RIBd updates")
	server.listenForRIBdUpdates(ribdCommonDefs.PUB_SOCKET_BFDD_ADDR)
	for {
		server.logger.Info("Read on RIBd subscriber socket...")
		rxBuf, err := server.ribdSubSocket.Recv(0)
		if err != nil {
			server.logger.Err("Recv on RIBd subscriber socket failed with error:", err)
			server.ribdSubSocketErrCh <- err
			continue
		}
		server.logger.Info("RIB subscriber recv returned:", rxBuf)
		server.ribdSubSocketCh <- rxBuf
	}
}

func (server *BFDServer) listenForRIBdUpdates(address string) error {
	var err error
	if server.ribdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		server.logger.Err("Failed to create RIBd subscribe socket, error:", err)
		return err
	}

	if _, err = server.ribdSubSocket.Connect(address); err != nil {
		server.logger.Err("Failed to connect to RIBd publisher socket, address:", address, "error:", err)
		return err
	}

	if err = server.ribdSubSocket.Subscribe(""); err != nil {
		server.logger.Err("Failed to subscribe to \"\" on RIBd subscribe socket, error:", err)
		return err
	}

	server.logger.Info("Connected to RIBd publisher at address:", address)
	if err = server.ribdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		server.logger.Err("Failed to set the buffer size for RIBd publisher socket, error:", err)
		return err
	}
	return nil
}

func (server *BFDServer) processRibdNotification(rxBuf []byte) error {
	var msg ribdCommonDefs.RibdNotifyMsg
	err := json.Unmarshal(rxBuf, &msg)
	if err != nil {
		server.logger.Err("Unable to unmarshal rxBuf:", rxBuf)
		return err
	}
	switch msg.MsgType {
	case ribdCommonDefs.NOTIFY_ROUTE_REACHABILITY_STATUS_UPDATE:
		server.logger.Info("Received NOTIFY_ROUTE_REACHABILITY_STATUS_UPDATE")
		var msgInfo ribdCommonDefs.RouteReachabilityStatusMsgInfo
		err = json.Unmarshal(msg.MsgBuf, &msgInfo)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg:", msg.MsgBuf)
			return err
		}
		server.logger.Info(" IP ", msgInfo.Network, " reachabilityStatus: ", msgInfo.IsReachable)
		if msgInfo.IsReachable {
			server.logger.Info(" NextHop IP:", msgInfo.NextHopIntf.NextHopIp, " IfIndex ", msgInfo.NextHopIntf.NextHopIfIndex)
			server.HandleNextHopChange(msgInfo.Network, int32(msgInfo.NextHopIntf.NextHopIfIndex), true)
		} else {
			server.logger.Info(" NextHop IP:", msgInfo.NextHopIntf.NextHopIp, " is not reachable ")
			server.HandleNextHopChange(msgInfo.Network, 0, false)
		}
		break
	default:
		break
	}
	return nil
}
