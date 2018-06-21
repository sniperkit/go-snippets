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
	"errors"
	"net"
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

type PortProperty struct {
	Name     string
	VlanName string
	VlanId   uint16
	IpAddr   string
}

type VlanProperty struct {
	Name       string
	UntagPorts []int32
	IpAddr     string
}

type LagProperty struct {
	Links []int32
}

func (server *BFDServer) updateVlanPropertyMap(msg asicdClntDefs.VlanNotifyMsg) {
	if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_CREATE {
		ent := server.vlanPropertyMap[int32(msg.VlanId)]
		ent.Name = msg.VlanName
		ent.UntagPorts = msg.UntagPorts
		server.vlanPropertyMap[int32(msg.VlanId)] = ent
	} else if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_DELETE {
		delete(server.vlanPropertyMap, int32(msg.VlanId))
	}
}

func (server *BFDServer) updatePortPropertyMap(msg asicdClntDefs.VlanNotifyMsg) {
	if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_CREATE { // Create Vlan
		for _, portNum := range msg.UntagPorts {
			ent := server.portPropertyMap[portNum]
			ent.VlanId = msg.VlanId
			ent.VlanName = msg.VlanName
			server.portPropertyMap[portNum] = ent
		}
	} else { // Delete Vlan
		for _, portNum := range msg.UntagPorts {
			ent := server.portPropertyMap[portNum]
			ent.VlanId = 0
			ent.VlanName = ""
			server.portPropertyMap[portNum] = ent
		}
	}
}

func (server *BFDServer) BuildPortPropertyMap() error {
	currMarker := int(asicdCommonDefs.MIN_SYS_PORTS)
	server.logger.Info("Calling asicd for port property")
	count := 10
	for {
		server.logger.Info("Calling bulkget port ", currMarker, count)
		bulkInfo, _ := server.AsicdPlugin.GetBulkPortState(currMarker, count)
		if bulkInfo == nil {
			server.logger.Info("Bulkget port got nothing")
			return nil
		}
		objCount := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		server.logger.Info("Bulkget port got ", objCount, more)
		currMarker = int(bulkInfo.EndIdx)
		for i := 0; i < objCount; i++ {
			ifIndex := bulkInfo.PortStateList[i].IfIndex
			ent := server.portPropertyMap[ifIndex]
			ent.Name = bulkInfo.PortStateList[i].Name
			ent.VlanId = 0
			ent.VlanName = ""
			server.portPropertyMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
	return nil
}

func (server *BFDServer) BuildLagPropertyMap() error {
	server.logger.Info("Get configured lags ... TBD")
	return nil
}

func (server *BFDServer) updateLagPropertyMap(msg asicdClntDefs.LagNotifyMsg) {
	_, exists := server.lagPropertyMap[msg.IfIndex]
	if msg.MsgType == asicdClntDefs.NOTIFY_LAG_CREATE { // Create LAG
		if exists {
			server.logger.Info("CreateLag: already exists", msg.IfIndex)
		} else {
			server.logger.Info("Creating lag ", msg.IfIndex)
			lagEntry := LagProperty{}
			lagEntry.Links = make([]int32, 0)
			for _, linkNum := range msg.IfIndexList {
				lagEntry.Links = append(lagEntry.Links, linkNum)
			}
			server.lagPropertyMap[msg.IfIndex] = lagEntry
		}
	} else if msg.MsgType == asicdClntDefs.NOTIFY_LAG_DELETE { // Delete Lag
		if exists {
			server.logger.Info("Deleting lag ", msg.IfIndex)
			delete(server.lagPropertyMap, msg.IfIndex)
		} else {
			server.logger.Info("DeleteLag: Does not exist ", msg.IfIndex)
		}
	}
}

func (server *BFDServer) getLinuxIntfName(ifIndex int32) (ifName string, err error) {
	ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(ifIndex)
	if ifType == asicdClntDefs.IfTypeVlan { // Vlan
		ifName = server.vlanPropertyMap[ifIndex].Name
	} else if ifType == asicdClntDefs.IfTypePort { // PHY
		ifName = server.portPropertyMap[int32(ifIndex)].Name
	} else {
		ifName = ""
		err = errors.New("Invalid Interface Type")
	}
	return ifName, err
}

func (server *BFDServer) getMacAddrFromIntfName(ifName string) (macAddr net.HardwareAddr, err error) {
	ifi, err := net.InterfaceByName(ifName)
	if err != nil {
		return macAddr, err
	}
	macAddr = ifi.HardwareAddr
	return macAddr, nil
}
