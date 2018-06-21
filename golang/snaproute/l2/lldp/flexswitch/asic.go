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

package flexswitch

import (
	"asicd/asicdCommonDefs"
	"asicdServices"
	"encoding/json"
	"errors"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"l2/lldp/api"
	"l2/lldp/config"
	"l2/lldp/utils"
	"strconv"
	"time"
	"utils/commonDefs"
	"utils/ipcutils"
)

type AsicPlugin struct {
	asicdClient    *asicdServices.ASICDServicesClient
	asicdSubSocket *nanomsg.SubSocket
}

func connectAsicd(filePath string, asicdClient chan *asicdServices.ASICDServicesClient) {
	fileName := filePath + CLIENTS_FILE_NAME

	clientJson, err := getClient(fileName, "asicd")
	if err != nil || clientJson == nil {
		asicdClient <- nil
		return
	}

	clientTransport, protocolFactory, err := ipcutils.CreateIPCHandles("localhost:" +
		strconv.Itoa(clientJson.Port))
	if err != nil {
		debug.Logger.Info("Failed to connect to ASICd, retrying until success")
		count := 0
		ticker := time.NewTicker(time.Duration(250) * time.Millisecond)
		for _ = range ticker.C {
			clientTransport, protocolFactory, err =
				ipcutils.CreateIPCHandles("localhost:" +
					strconv.Itoa(clientJson.Port))
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				debug.Logger.Info("Still waiting to connect to ASICd")
			}
		}
	}
	client := asicdServices.NewASICDServicesClientFactory(clientTransport,
		protocolFactory)
	asicdClient <- client
}

func NewAsicPlugin(fileName string) (*AsicPlugin, error) {
	var asicdClient *asicdServices.ASICDServicesClient = nil
	asicdClientCh := make(chan *asicdServices.ASICDServicesClient)

	debug.Logger.Info("Connecting to ASICd")
	go connectAsicd(fileName, asicdClientCh)
	asicdClient = <-asicdClientCh
	if asicdClient == nil {
		debug.Logger.Err("Failed to connecto to ASICd")
		return nil, errors.New("Failed to connect to ASICd")
	}

	mgr := &AsicPlugin{
		asicdClient: asicdClient,
	}
	return mgr, nil

}

/*  Helper function to get bulk port state information from asicd
 */
func (p *AsicPlugin) getPortStates() []*config.PortInfo {
	debug.Logger.Info("Get Port State List")
	currMarker := int64(asicdCommonDefs.MIN_SYS_PORTS)
	more := false
	objCount := 0
	count := 500
	portStates := make([]*config.PortInfo, 0)
	for {
		bulkInfo, err := p.asicdClient.GetBulkPortState(asicdServices.Int(currMarker), asicdServices.Int(count))
		if err != nil {
			debug.Logger.Err(fmt.Sprintln(": getting bulk port config from asicd failed with reason", err))
			break
		}
		objCount = int(bulkInfo.Count)
		more = bool(bulkInfo.More)
		currMarker = int64(bulkInfo.EndIdx)
		for i := 0; i < objCount; i++ {
			obj := bulkInfo.PortStateList[i]
			port := &config.PortInfo{
				IfIndex:   obj.IfIndex,
				OperState: obj.OperState,
				Name:      obj.IntfRef,
				Pvid:      obj.Pvid,
			}
			pObj, err := p.asicdClient.GetPort(obj.IntfRef)
			if err != nil {
				debug.Logger.Err(fmt.Sprintln("Getting mac address for", obj.Name, "failed, error:", err))
			} else {
				port.MacAddr = pObj.MacAddr
				port.Description = pObj.Description
			}
			debug.Logger.Debug("Adding port Name, OperState, IfIndex:", port, "to portStates")
			portStates = append(portStates, port)
		}
		if more == false {
			break
		}
	}
	debug.Logger.Info("Done with Port State list")
	return portStates
}

func (p *AsicPlugin) GetPortsInfo() []*config.PortInfo {
	portStates := p.getPortStates()
	return portStates
}

func (p *AsicPlugin) connectSubSocket() error {
	var err error
	address := asicdCommonDefs.PUB_SOCKET_ADDR
	debug.Logger.Info(" setting up asicd update listener")
	if p.asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		debug.Logger.Err(fmt.Sprintln("Failed to create ASIC subscribe socket, error:",
			err))
		return err
	}

	if err = p.asicdSubSocket.Subscribe(""); err != nil {
		debug.Logger.Err(fmt.Sprintln("Failed to subscribe to ASIC subscribe socket",
			"error:", err))
		return err
	}

	if _, err = p.asicdSubSocket.Connect(address); err != nil {
		debug.Logger.Err(fmt.Sprintln("Failed to connect to ASIC publisher socket",
			"address:", address, "error:", err))
		return err
	}

	debug.Logger.Info(fmt.Sprintln(" Connected to ASIC publisher at address:", address))
	if err = p.asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		debug.Logger.Err(fmt.Sprintln(" Failed to set the buffer size for ASIC publisher",
			"socket, error:", err))
		return err
	}
	debug.Logger.Info("asicd update listener is set")
	return nil
}

func (p *AsicPlugin) listenAsicdUpdates() {
	for {
		rxBuf, err := p.asicdSubSocket.Recv(0)
		if err != nil {
			debug.Logger.Err(fmt.Sprintln(
				"Recv on asicd Subscriber socket failed with error:", err))
			continue
		}
		var msg asicdCommonDefs.AsicdNotification
		err = json.Unmarshal(rxBuf, &msg)
		if err != nil {
			debug.Logger.Err(fmt.Sprintln("Unable to Unmarshal asicd msg:", msg.Msg))
			continue
		}
		switch msg.MsgType {
		case asicdCommonDefs.NOTIFY_L2INTF_STATE_CHANGE:
			var l2IntfStateNotifyMsg asicdCommonDefs.L2IntfStateNotifyMsg
			err = json.Unmarshal(msg.Msg, &l2IntfStateNotifyMsg)
			if err != nil {
				debug.Logger.Err("Unable to Unmarshal l2 intf state change:", msg.Msg)
				continue
			}
			debug.Logger.Debug("Got Notification from Asicd Subscriber socket for ifIndex:", l2IntfStateNotifyMsg.IfIndex,
				"State:", l2IntfStateNotifyMsg.IfState)
			if l2IntfStateNotifyMsg.IfState == asicdCommonDefs.INTF_STATE_UP {
				api.SendPortStateChange(l2IntfStateNotifyMsg.IfIndex, "UP")
			} else {
				api.SendPortStateChange(l2IntfStateNotifyMsg.IfIndex, "DOWN")
			}
		case asicdCommonDefs.NOTIFY_PORT_ATTR_CHANGE:
			var portAttrMsg asicdCommonDefs.PortAttrChangeNotifyMsg
			err = json.Unmarshal(msg.Msg, &portAttrMsg)
			if err != nil {
				debug.Logger.Err("Unable to Unmarshal Port Attr Change Notify message err:", err)
				continue
			}

			debug.Logger.Debug("Received port attribute change message:", portAttrMsg)
			if (portAttrMsg.AttrMask&commonDefs.PORT_ATTR_DESCRIPTION) == commonDefs.PORT_ATTR_DESCRIPTION ||
				(portAttrMsg.AttrMask&commonDefs.PORT_ATTR_PVID) == commonDefs.PORT_ATTR_PVID {
				debug.Logger.Debug("Received description attribute change")
				api.SendPortAttrChange(portAttrMsg.IfIndex, portAttrMsg.Description, portAttrMsg.Pvid)
			}
		}
	}

}
func (p *AsicPlugin) Start() {
	err := p.connectSubSocket()
	if err != nil {
		return
	}
	go p.listenAsicdUpdates()
}
