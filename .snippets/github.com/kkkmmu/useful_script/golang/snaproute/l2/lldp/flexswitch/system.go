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
	"encoding/json"
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/sysdCommonDefs"
	"l2/lldp/api"
	"l2/lldp/config"
	"l2/lldp/utils"
	"models/objects"
	"time"
	"utils/dbutils"
	"utils/eventUtils"
)

type SystemPlugin struct {
	sysdSubSocket *nanomsg.SubSocket
}

func NewSystemPlugin(fileName string, db *dbutils.DBUtil) (*SystemPlugin, error) {
	mgr := &SystemPlugin{}
	err := eventUtils.InitEvents("LLDPD", db, db, debug.Logger, 1000)
	if err != nil {
		debug.Logger.Info("unable to initialize event utils", err)
	}
	return mgr, nil
}

func (p *SystemPlugin) connectSubSocket() error {
	var err error
	address := sysdCommonDefs.PUB_SOCKET_ADDR
	debug.Logger.Info(" setting up sysd update listener")
	if p.sysdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		debug.Logger.Err("Failed to create SYS subscribe socket, error:", err)
		return err
	}

	if err = p.sysdSubSocket.Subscribe(""); err != nil {
		debug.Logger.Err("Failed to subscribe to SYS subscribe socket error:", err)
		return err
	}

	if _, err = p.sysdSubSocket.Connect(address); err != nil {
		debug.Logger.Err("Failed to connect to SYS publisher socket",
			"address:", address, "error:", err)
		return err
	}

	debug.Logger.Info("Connected to SYS publisher at address:", address)
	if err = p.sysdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		debug.Logger.Err(" Failed to set the buffer size for SYS publisher",
			"socket, error:", err)
		return err
	}
	debug.Logger.Info("sysd update listener is set")
	return nil
}

func (p *SystemPlugin) listenSystemdUpdates() {
	for {
		debug.Logger.Debug("Read on System Subscriber socket....")
		rxBuf, err := p.sysdSubSocket.Recv(0)
		if err != nil {
			debug.Logger.Err("Recv on sysd Subscriber socket failed with error:", err)
			continue
		}
		var msg sysdCommonDefs.Notification
		err = json.Unmarshal(rxBuf, &msg)
		if err != nil {
			debug.Logger.Err("Unable to Unmarshal sysd err:", err)
			continue
		}
		debug.Logger.Debug("LLDP recv msg type:", msg.Type)
		switch msg.Type {
		case sysdCommonDefs.SYSTEM_Info:
			var systemInfo objects.SystemParam
			err = json.Unmarshal(msg.Payload, &systemInfo)
			sysInfo := config.SystemInfo{
				Vrf:         systemInfo.Vrf,
				MgmtIp:      systemInfo.MgmtIp,
				Hostname:    systemInfo.Hostname,
				SwitchMac:   systemInfo.SwitchMac,
				SwVersion:   systemInfo.SwVersion,
				Description: systemInfo.Description,
			}
			debug.Logger.Debug("LLDP received system update:", sysInfo)
			api.UpdateCache(&sysInfo)
		}
	}
}

func (p *SystemPlugin) Start() {
	err := p.connectSubSocket()
	if err != nil {
		return
	}
	go p.listenSystemdUpdates()
}

func (p *SystemPlugin) GetSystemInfo(dbHdl *dbutils.DBUtil) *config.SystemInfo {
	sysInfo := &config.SystemInfo{}
	count := 0
	for {
		if count%10000 == 0 {
			debug.Logger.Info("Still trying to get db info for System Param")
		}
		count++
		if dbHdl == nil {
			time.Sleep(250 * time.Millisecond)
			continue
		}

		var dbObj objects.SystemParam
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			debug.Logger.Err("DB query failed for System Info, retry in next 250 Millisecond")
			time.Sleep(250 * time.Millisecond)
			continue
		}
		for idx := 0; idx < len(objList); idx++ {
			dbObject := objList[idx].(objects.SystemParam)
			sysInfo.SwitchMac = dbObject.SwitchMac
			sysInfo.MgmtIp = dbObject.MgmtIp
			sysInfo.SwVersion = dbObject.SwVersion
			sysInfo.Description = dbObject.Description
			sysInfo.Hostname = dbObject.Hostname
			sysInfo.Vrf = dbObject.Vrf
			return sysInfo
		}
	}
}
