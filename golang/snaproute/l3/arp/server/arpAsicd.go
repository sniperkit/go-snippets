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

package server

import (
	"errors"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/clntUtils/clntIntfs"
)

type AsicdMsgType uint8

const (
	CreateAsicdEntry AsicdMsgType = 1
	DeleteAsicdEntry AsicdMsgType = 2
)

type AsicdMsg struct {
	MsgType AsicdMsgType
	IpAddr  string
	MacAddr string
	VlanId  int32
	IfIdx   int32
}

func (server *ARPServer) processAsicdNotification(msg clntIntfs.NotifyMsg) {
	switch msg.(type) {
	case asicdClntDefs.L2IntfStateNotifyMsg:
		l2Msg := msg.(asicdClntDefs.L2IntfStateNotifyMsg)
		server.logger.Debug("L2IntfStateNotifyMsg:", l2Msg)
		server.processL2StateChange(l2Msg)
		server.dumpInfra()
	case asicdClntDefs.IPv4L3IntfStateNotifyMsg:
		l3Msg := msg.(asicdClntDefs.IPv4L3IntfStateNotifyMsg)
		server.dumpInfra()
		server.logger.Debug("IPv4L3IntfStateNotifyMsg:", l3Msg)
		server.processIPv4L3StateChange(l3Msg)
		server.dumpInfra()
	case asicdClntDefs.VlanNotifyMsg:
		vlanMsg := msg.(asicdClntDefs.VlanNotifyMsg)
		server.logger.Debug("VlanNotifyMsg:", vlanMsg)
		server.updateVlanInfra(vlanMsg)
		server.dumpInfra()
	case asicdClntDefs.LagNotifyMsg:
		lagMsg := msg.(asicdClntDefs.LagNotifyMsg)
		server.logger.Debug("LagNotifyMsg:", lagMsg)
		server.updateLagInfra(lagMsg)
		server.dumpInfra()
	case asicdClntDefs.IPv4IntfNotifyMsg:
		ipv4Msg := msg.(asicdClntDefs.IPv4IntfNotifyMsg)
		server.logger.Debug("IPv4IntfNotifyMsg:", ipv4Msg)
		server.updateIPv4Infra(ipv4Msg)
		server.dumpInfra()
	case asicdClntDefs.IPv4NbrMacMoveNotifyMsg:
		macMoveMsg := msg.(asicdClntDefs.IPv4NbrMacMoveNotifyMsg)
		server.processIPv4NbrMacMove(macMoveMsg)
		server.dumpInfra()
	case asicdClntDefs.IPv4VirtualIntfNotifyMsg:
		virIntfMsg := msg.(asicdClntDefs.IPv4VirtualIntfNotifyMsg)
		server.logger.Info("Msg Virtual Intf:", virIntfMsg)
		server.processVirtualIntfEvent(virIntfMsg)
	case asicdClntDefs.IPv4VirtualIntfStateNotifyMsg:
		virStMsg := msg.(asicdClntDefs.IPv4VirtualIntfStateNotifyMsg)
		server.logger.Info("Virtual Intf State Change Message:", virStMsg)
		server.processVirtualIntfStateEvent(virStMsg)
	}
}

func (server *ARPServer) processAsicdMsg(msg AsicdMsg) error {
	if server.AsicdPlugin == nil {
		return errors.New("No asicd plugin")
	}
	switch msg.MsgType {
	case CreateAsicdEntry:
		server.logger.Debug("CreateAsicdEntry:", msg)
		_, err := server.AsicdPlugin.CreateIPv4Neighbor(msg.IpAddr, msg.MacAddr, msg.VlanId, msg.IfIdx)
		if err != nil {
			server.logger.Err("Asicd Create IPv4 Neighbor failed for IpAddr:", msg.IpAddr, "VlanId:", msg.VlanId, "IfIdx:", msg.IfIdx, "err:", err)
			return err
		}
	case DeleteAsicdEntry:
		server.logger.Debug("DeleteAsicdEntry:", msg)
		_, err := server.AsicdPlugin.DeleteIPv4Neighbor(msg.IpAddr, "00:00:00:00:00:00", 0, 0)
		if err != nil {
			server.logger.Err("Asicd was unable to delete neigbhor entry for", msg.IpAddr, "err:", err)
			return err
		}
	default:
		err := errors.New("Invalid Asicd Msg Type")
		return err
	}
	return nil
}
