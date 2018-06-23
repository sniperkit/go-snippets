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
	"asicd/asicdCommonDefs"
	"asicdInt"
	"asicdServices"
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"l3/ospf/config"
	//"utils/commonDefs"
)

type AsicdClient struct {
	OspfClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

func (server *OSPFServer) createASICdSubscriber() {
	for {
		server.logger.Info("Read on ASICd subscriber socket...")
		asicdrxBuf, err := server.asicdSubSocket.Recv(0)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Recv on ASICd subscriber socket failed with error:", err))
			server.asicdSubSocketErrCh <- err
			continue
		}
		server.logger.Info(fmt.Sprintln("ASIC subscriber recv returned:", asicdrxBuf))
		server.asicdSubSocketCh <- asicdrxBuf
	}
}

func (server *OSPFServer) listenForASICdUpdates(address string) error {
	var err error
	if server.asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		server.logger.Err(fmt.Sprintln("Failed to create ASICd subscribe socket, error:", err))
		return err
	}

	if err = server.asicdSubSocket.Subscribe(""); err != nil {
		server.logger.Err(fmt.Sprintln("Failed to subscribe to \"\" on ASICd subscribe socket, error:", err))
		return err
	}

	if _, err = server.asicdSubSocket.Connect(address); err != nil {
		server.logger.Err(fmt.Sprintln("Failed to connect to ASICd publisher socket, address:", address, "error:", err))
		return err
	}

	server.logger.Info(fmt.Sprintln("Connected to ASICd publisher at address:", address))
	if err = server.asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		server.logger.Err(fmt.Sprintln("Failed to set the buffer size for ASICd publisher socket, error:", err))
		return err
	}
	return nil
}

func (server *OSPFServer) processAsicdNotification(asicdrxBuf []byte) {
	var msg asicdCommonDefs.AsicdNotification
	err := json.Unmarshal(asicdrxBuf, &msg)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Unable to unmarshal asicdrxBuf:", asicdrxBuf))
		return
	}
	/*
		if msg.MsgType == asicdCommonDefs.NOTIFY_PORT_CONFIG_MTU_CHANGE {
			var mtuChangeMsg asicdCommonDefs.PortConfigMtuChgNotigyMsg
			err = json.Unmarshal(msg.Msg, &mtuChangeMsg)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Mtu change :Unable to unmarshal msg :", msg.Msg))
				return
			}
			server.UpdateMtu(mtuChangeMsg.IfIndex, mtuChangeMsg.Mtu)
		}
	*/

	if msg.MsgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE ||
		msg.MsgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_DELETE {
		var newLogicalIntfMgs asicdCommonDefs.LogicalIntfNotifyMsg
		err = json.Unmarshal(msg.Msg, &newLogicalIntfMgs)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to unmarshal msg : ", msg.Msg))
			return
		}
		server.UpdateLogicalIntfInfra(newLogicalIntfMgs, msg.MsgType)
	}

	if msg.MsgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE ||
		msg.MsgType == asicdCommonDefs.NOTIFY_IPV4INTF_DELETE {
		var NewIpv4IntfMsg asicdCommonDefs.IPv4IntfNotifyMsg
		err = json.Unmarshal(msg.Msg, &NewIpv4IntfMsg)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to unmarshal msg:", msg.Msg))
			return
		}
		server.UpdateIPv4Infra(NewIpv4IntfMsg, msg.MsgType)
	} else if msg.MsgType == asicdCommonDefs.NOTIFY_VLAN_CREATE ||
		msg.MsgType == asicdCommonDefs.NOTIFY_VLAN_DELETE {
		var vlanNotifyMsg asicdCommonDefs.VlanNotifyMsg
		err = json.Unmarshal(msg.Msg, &vlanNotifyMsg)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to unmashal vlanNotifyMsg:", msg.Msg))
			return
		}
		server.UpdateVlanInfra(vlanNotifyMsg, msg.MsgType)
	}

	if msg.MsgType == asicdCommonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE {
		var newIpv4StateMsg asicdCommonDefs.IPv4L3IntfStateNotifyMsg
		err = json.Unmarshal(msg.Msg, &newIpv4StateMsg)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to unmarshal msg :", msg))
			return
		}
		server.processIntfStateChange(newIpv4StateMsg.IpAddr,
			newIpv4StateMsg.IfIndex,
			config.Status(newIpv4StateMsg.IfState),
			msg.MsgType)
	}
}

func (server *OSPFServer) initAsicdForRxMulticastPkt() (err error) {
	// All SPF Router
	allSPFRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLSPFROUTERMAC,
		MacAddrMask: MASKMAC,
	}
	if server.asicdClient.ClientHdl == nil {
		server.logger.Err("Null asicd client handle")
		return nil
	}
	ret, err := server.asicdClient.ClientHdl.EnablePacketReception(&allSPFRtrMacConf)
	if !ret {
		server.logger.Info(fmt.Sprintln("Adding reserved mac failed", ALLSPFROUTERMAC))
		return err
	}

	// All D Router
	allDRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLDROUTERMAC,
		MacAddrMask: MASKMAC,
	}

	ret, err = server.asicdClient.ClientHdl.EnablePacketReception(&allDRtrMacConf)
	if !ret {
		server.logger.Info(fmt.Sprintln("Adding reserved mac failed", ALLDROUTERMAC))
		return err
	}
	return nil
}

const (
	numberedP2P   int = 0
	unnumberedP2P int = 1
	broadcast     int = 2
)
