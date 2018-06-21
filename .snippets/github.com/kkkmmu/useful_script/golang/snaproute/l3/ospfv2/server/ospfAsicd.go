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
	"asicd/asicdCommonDefs"
	"asicdInt"
	"asicdServices"
	"encoding/json"
	//"git.apache.org/thrift.git/lib/go/thrift"
	nanomsg "github.com/op/go-nanomsg"
	"strconv"
	"time"
	"utils/commonDefs"
	"utils/ipcutils"
)

type AsicdClient struct {
	OspfClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

type AsicdCommStruct struct {
	asicdSubSocketCh    chan []byte
	asicdClient         AsicdClient
	asicdSubSocket      *nanomsg.SubSocket
	asicdSubSocketErrCh chan error
}

func (server *OSPFV2Server) initAsicdComm() error {
	server.asicdComm.asicdSubSocketCh = make(chan []byte)
	server.asicdComm.asicdSubSocketErrCh = make(chan error)
	return nil
}

func (server *OSPFV2Server) ConnectToAsicdServer(port int) {
	var err error
	server.logger.Info("found asicd at port", port)
	server.asicdComm.asicdClient.Address = "localhost:" + strconv.Itoa(port)
	server.asicdComm.asicdClient.Transport, server.asicdComm.asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.asicdComm.asicdClient.Address)
	if err != nil {
		server.logger.Info("Failed to connect to Asicd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			server.asicdComm.asicdClient.Transport, server.asicdComm.asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.asicdComm.asicdClient.Address)
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				server.logger.Info("Still can't connect to asicd, retrying..")
			}
		}
	}
	server.logger.Info("Ospfd is connected to asicd")
	server.asicdComm.asicdClient.ClientHdl = asicdServices.NewASICDServicesClientFactory(server.asicdComm.asicdClient.Transport, server.asicdComm.asicdClient.PtrProtocolFactory)
	server.asicdComm.asicdClient.IsConnected = true
}

func (server *OSPFV2Server) StartAsicdSubscriber() {
	server.logger.Info("Listen for asicd updates")
	server.listenForAsicdUpdates(asicdCommonDefs.PUB_SOCKET_ADDR)
	go server.createAsicdSubscriber()
}

func (server *OSPFV2Server) createAsicdSubscriber() {
	for {
		server.logger.Info("Read on asicd subscriber socket...")
		asicdrxBuf, err := server.asicdComm.asicdSubSocket.Recv(0)
		if err != nil {
			server.logger.Err("Recv on ASICd subscriber socket failed with error:", err)
			server.asicdComm.asicdSubSocketErrCh <- err
			continue
		}
		server.logger.Info("asic subscriber recv returned:", asicdrxBuf)
		server.asicdComm.asicdSubSocketCh <- asicdrxBuf
	}
}

func (server *OSPFV2Server) listenForAsicdUpdates(address string) error {
	var err error
	if server.asicdComm.asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		server.logger.Err("Failed to create ASICd subscribe socket, error:", err)
		return err
	}

	if err = server.asicdComm.asicdSubSocket.Subscribe(""); err != nil {
		server.logger.Err("Failed to subscribe to \"\" on ASICd subscribe socket, error:", err)
		return err
	}

	if _, err = server.asicdComm.asicdSubSocket.Connect(address); err != nil {
		server.logger.Err("Failed to connect to ASICd publisher socket, address:", address, "error:", err)
		return err
	}

	server.logger.Info("Connected to ASICd publisher at address:", address)
	if err = server.asicdComm.asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		server.logger.Err("Failed to set the buffer size for ASICd publisher socket, error:", err)
		return err
	}
	return nil
}

func (server *OSPFV2Server) processAsicdNotification(asicdRxBuf []byte) {
	var asicdMsg asicdCommonDefs.AsicdNotification
	err := json.Unmarshal(asicdRxBuf, &asicdMsg)
	if err != nil {
		server.logger.Err("Unable to unmarshal asicdrxBuf:", asicdRxBuf)
		return
	}
	if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_PORT_ATTR_CHANGE {
		var msg asicdCommonDefs.PortAttrChangeNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Mtu change :Unable to unmarshal msg :", asicdMsg.Msg)
			return
		}
		if (msg.AttrMask & commonDefs.PORT_ATTR_MTU) == commonDefs.PORT_ATTR_MTU {
			server.UpdateMtu(msg)
		}
	} else if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_DELETE {
		var msg asicdCommonDefs.LogicalIntfNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg : ", asicdMsg.Msg)
			return
		}
		server.UpdateLogicalIntfInfra(msg, asicdMsg.MsgType)
	} else if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_IPV4INTF_DELETE {
		//server.logger.Info("Recv NOTIFY_IPV4INTF_CREATE/NOTIFY_IPV4INTF_CREATE Msg")
		var msg asicdCommonDefs.IPv4IntfNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg:", asicdMsg.Msg)
			return
		}
		server.UpdateIPv4Infra(msg, asicdMsg.MsgType)
	} else if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE {
		//server.logger.Info("Recv NOTIFY_IPV4_L3INTF_STATE_CHANGE Msg")
		var msg asicdCommonDefs.IPv4L3IntfStateNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg :", asicdMsg.Msg)
			return
		}
		server.ProcessIPv4StateChange(msg)
	} else if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_VLAN_CREATE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_VLAN_DELETE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_VLAN_UPDATE {
		var msg asicdCommonDefs.VlanNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg :", asicdMsg.Msg)
			return
		}
		server.ProcessVlanNotify(msg, asicdMsg.MsgType)
	} else if asicdMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_CREATE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_DELETE ||
		asicdMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_UPDATE {
		var msg asicdCommonDefs.LagNotifyMsg
		err = json.Unmarshal(asicdMsg.Msg, &msg)
		if err != nil {
			server.logger.Err("Unable to unmarshal msg :", asicdMsg.Msg)
			return
		}
		server.ProcessLagNotify(msg, asicdMsg.MsgType)
	}
}

func (server *OSPFV2Server) initAsicdForRxMulticastPkt() (err error) {
	// All SPF Router
	allSPFRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLSPFROUTERMAC,
		MacAddrMask: MASKMAC,
	}
	if server.asicdComm.asicdClient.ClientHdl == nil {
		server.logger.Err("Null asicd client handle")
		return nil
	}
	ret, err := server.asicdComm.asicdClient.ClientHdl.EnablePacketReception(&allSPFRtrMacConf)
	if !ret {
		server.logger.Info("Adding reserved mac failed", ALLSPFROUTERMAC)
		return err
	}

	// All D Router
	allDRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLDROUTERMAC,
		MacAddrMask: MASKMAC,
	}

	ret, err = server.asicdComm.asicdClient.ClientHdl.EnablePacketReception(&allDRtrMacConf)
	if !ret {
		server.logger.Info("Adding reserved mac failed", ALLDROUTERMAC)
		return err
	}
	return nil
}

func (server *OSPFV2Server) deinitAsicdForRxMulticastPkt() (err error) {
	// All SPF Router
	allSPFRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLSPFROUTERMAC,
		MacAddrMask: MASKMAC,
	}
	if server.asicdComm.asicdClient.ClientHdl == nil {
		server.logger.Err("Null asicd client handle")
		return nil
	}
	ret, err := server.asicdComm.asicdClient.ClientHdl.DisablePacketReception(&allSPFRtrMacConf)
	if !ret {
		server.logger.Info("Deleting reserved mac failed", ALLSPFROUTERMAC)
		return err
	}

	// All D Router
	allDRtrMacConf := asicdInt.RsvdProtocolMacConfig{
		MacAddr:     ALLDROUTERMAC,
		MacAddrMask: MASKMAC,
	}

	ret, err = server.asicdComm.asicdClient.ClientHdl.DisablePacketReception(&allDRtrMacConf)
	if !ret {
		server.logger.Info("Deleting reserved mac failed", ALLDROUTERMAC)
		return err
	}
	return nil
}
