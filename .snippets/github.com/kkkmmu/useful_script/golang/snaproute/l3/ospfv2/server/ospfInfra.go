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
	//"errors"
	"net"
	"utils/commonDefs"
)

type PortProperty struct {
	Name       string
	Mtu        int32
	Speed      uint32 //Unit Mbps
	IpIfIdxMap map[int32]bool
}

type LagProperty struct {
	Name    string
	PortMap map[int32]bool
}

type VlanProperty struct {
	Name          string
	UntagIfIdxMap map[int32]bool
	TagIfIdxMap   map[int32]bool
}

type LogicalIntfProperty struct {
	IfName string
}

type IpProperty struct {
	IfId    uint32
	IfType  uint32
	IfName  string
	IpAddr  uint32
	NetMask uint32
	MacAddr net.HardwareAddr
	Mtu     int32
	//Cost    uint32
	State bool
}

type InfraStruct struct {
	portPropertyMap        map[int32]PortProperty
	lagPropertyMap         map[int32]LagProperty
	vlanPropertyMap        map[int32]VlanProperty
	logicalIntfPropertyMap map[int32]LogicalIntfProperty
	ipPropertyMap          map[int32]IpProperty
	ipToIfIdxMap           map[uint32]int32
}

func (server *OSPFV2Server) initInfra() {
	server.infraData.portPropertyMap = make(map[int32]PortProperty)
	server.infraData.vlanPropertyMap = make(map[int32]VlanProperty)
	server.infraData.lagPropertyMap = make(map[int32]LagProperty)
	server.infraData.logicalIntfPropertyMap = make(map[int32]LogicalIntfProperty)
	server.infraData.ipPropertyMap = make(map[int32]IpProperty)
	server.infraData.ipToIfIdxMap = make(map[uint32]int32)
}

func (server *OSPFV2Server) buildInfra() {
	server.constructPortInfra()
	server.constructLagInfra()
	server.constructVlanInfra()
	//server.constructLogicalInfra()
	server.constructL3Infra()
	server.processInfra()
}

func (server *OSPFV2Server) constructPortInfra() {
	server.getBulkPortState()
	server.getBulkPortConfig()
}

func getMACAddr(ifName string) (macAddr net.HardwareAddr, err error) {
	ifi, err := net.InterfaceByName(ifName)
	if err != nil {
		return macAddr, err
	}
	macAddr = ifi.HardwareAddr
	return macAddr, nil
}

func (server *OSPFV2Server) getMTU(ifType uint32, ifIdx int32, flag bool) int32 {
	var minMtu int32 = 10000             //in bytes
	if ifType == commonDefs.IfTypePort { // PHY
		ent, _ := server.infraData.portPropertyMap[ifIdx]
		if flag == true {
			minMtu = ent.Mtu
			ent.IpIfIdxMap[ifIdx] = true
		} else {
			delete(ent.IpIfIdxMap, ifIdx)
		}
		server.infraData.portPropertyMap[ifIdx] = ent
	} else if ifType == commonDefs.IfTypeLag { //Lag
		ent, _ := server.infraData.lagPropertyMap[ifIdx]
		for portIfIdx, _ := range ent.PortMap {
			entry, _ := server.infraData.portPropertyMap[portIfIdx]
			if flag == true {
				if minMtu > entry.Mtu {
					minMtu = entry.Mtu
				}
				entry.IpIfIdxMap[ifIdx] = true
			} else {
				delete(entry.IpIfIdxMap, ifIdx)
			}
			server.infraData.portPropertyMap[portIfIdx] = entry
		}
	} else if ifType == commonDefs.IfTypeVlan { // Vlan
		ent, _ := server.infraData.vlanPropertyMap[ifIdx]
		for idx, _ := range ent.UntagIfIdxMap {
			iType := asicdCommonDefs.GetIntfTypeFromIfIndex(idx)
			if iType == commonDefs.IfTypeLag {
				lagEnt, _ := server.infraData.lagPropertyMap[idx]
				for portIfIdx, _ := range lagEnt.PortMap {
					entry, _ := server.infraData.portPropertyMap[portIfIdx]
					if flag == true {
						if minMtu > entry.Mtu {
							minMtu = entry.Mtu
						}
						entry.IpIfIdxMap[ifIdx] = true
					} else {
						delete(entry.IpIfIdxMap, ifIdx)
					}
					server.infraData.portPropertyMap[portIfIdx] = entry
				}
			} else {
				entry, _ := server.infraData.portPropertyMap[idx]
				if flag == true {
					if minMtu > entry.Mtu {
						minMtu = entry.Mtu
					}
					entry.IpIfIdxMap[ifIdx] = true
				} else {
					delete(entry.IpIfIdxMap, ifIdx)
				}
				server.infraData.portPropertyMap[idx] = entry
			}
		}
		for idx, _ := range ent.TagIfIdxMap {
			iType := asicdCommonDefs.GetIntfTypeFromIfIndex(idx)
			if iType == commonDefs.IfTypeLag {
				lagEnt, _ := server.infraData.lagPropertyMap[idx]
				for portIfIdx, _ := range lagEnt.PortMap {
					entry, _ := server.infraData.portPropertyMap[portIfIdx]
					if flag == true {
						if minMtu > entry.Mtu {
							minMtu = entry.Mtu
						}
						entry.IpIfIdxMap[ifIdx] = true
					} else {
						delete(entry.IpIfIdxMap, ifIdx)
					}
					server.infraData.portPropertyMap[portIfIdx] = entry
				}
			} else {
				entry, _ := server.infraData.portPropertyMap[idx]
				if flag == true {
					if minMtu > entry.Mtu {
						minMtu = entry.Mtu
					}
					entry.IpIfIdxMap[ifIdx] = true
				} else {
					delete(entry.IpIfIdxMap, ifIdx)
				}
				server.infraData.portPropertyMap[idx] = entry
			}
		}
	} else if ifType == commonDefs.IfTypeLoopback {
		minMtu = LOGICAL_INTF_MTU
	}
	return minMtu
}

func (server *OSPFV2Server) processInfra() {
	var err error
	for ifIdx, ipEnt := range server.infraData.ipPropertyMap {
		ipEnt.MacAddr, err = getMACAddr(ipEnt.IfName)
		if err != nil {
			server.logger.Err("Error getting MAC Address")
			continue
		}
		ipEnt.Mtu = server.getMTU(ipEnt.IfType, ifIdx, true)
		server.infraData.ipPropertyMap[ifIdx] = ipEnt
	}
}

func (server *OSPFV2Server) getBulkPortState() {
	currMarker := asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS)
	if server.asicdComm.asicdClient.IsConnected {
		server.logger.Info("Calling asicd for getting port state")
		count := 100
		for {
			bulkInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkPortState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo == nil {
				break
			}
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = asicdServices.Int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifIndex := bulkInfo.PortStateList[i].IfIndex
				ent := server.infraData.portPropertyMap[ifIndex]
				ent.Name = bulkInfo.PortStateList[i].Name
				server.infraData.portPropertyMap[ifIndex] = ent
			}
			if more == false {
				break
			}
		}
	}
}

func (server *OSPFV2Server) getBulkPortConfig() {
	currMarker := asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS)
	if server.asicdComm.asicdClient.IsConnected {
		server.logger.Info("Calling asicd for getting the Port Config")
		count := 100
		for {
			bulkInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkPort(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo == nil {
				break
			}
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = asicdServices.Int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifIndex := bulkInfo.PortList[i].IfIndex
				ent := server.infraData.portPropertyMap[ifIndex]
				ent.Mtu = bulkInfo.PortList[i].Mtu
				ent.Speed = uint32(bulkInfo.PortList[i].Speed)
				ent.IpIfIdxMap = make(map[int32]bool)
				server.infraData.portPropertyMap[ifIndex] = ent
			}
			if more == false {
				break
			}
		}
	}
}

func (server *OSPFV2Server) constructLagInfra() {
	server.logger.Info("Calling asicd for getting Lag")
	currMarker := 0
	count := 100
	if server.asicdComm.asicdClient.IsConnected {
		for {
			bulkInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkLag(asicdInt.Int(currMarker), asicdInt.Int(count))
			if bulkInfo == nil {
				break
			}
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifIndex := bulkInfo.LagList[i].LagIfIndex
				ent := server.infraData.lagPropertyMap[ifIndex]
				ent.Name = bulkInfo.LagList[i].LagName
				ent.PortMap = make(map[int32]bool)
				ifIndexList := bulkInfo.LagList[i].IfIndexList
				for i := 0; i < len(ifIndexList); i++ {
					ifIdx := ifIndexList[i]
					ent.PortMap[ifIdx] = true
				}
				server.infraData.lagPropertyMap[ifIndex] = ent
			}
			if more == false {
				break
			}
		}
	}
}

func (server *OSPFV2Server) constructVlanInfra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting Vlan Property")
	count := 100
	for {
		if server.asicdComm.asicdClient.ClientHdl == nil {
			server.logger.Err("Infra: Null client handle for asicd ")
			return
		}
		bulkVlanInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
		// Get bulk on vlan state can re-use curMark and count used
		// by get bulk vlan, as there is a 1:1 mapping in terms of cfg/state objs
		bulkVlanStateInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkVlanStateInfo == nil &&
			bulkVlanInfo == nil {
			break
		}
		objCnt := int(bulkVlanInfo.Count)
		more := bool(bulkVlanInfo.More)
		curMark = int(bulkVlanInfo.EndIdx)
		for idx := 0; idx < objCnt; idx++ {
			vlanIfIdx := bulkVlanStateInfo.VlanStateList[idx].IfIndex
			vlanEnt := server.infraData.vlanPropertyMap[vlanIfIdx]
			vlanEnt.Name = bulkVlanStateInfo.VlanStateList[idx].VlanName
			uIfIdxList := bulkVlanInfo.VlanList[idx].UntagIfIndexList
			vlanEnt.UntagIfIdxMap = make(map[int32]bool)
			for i := 0; i < len(uIfIdxList); i++ {
				vlanEnt.UntagIfIdxMap[uIfIdxList[i]] = true
			}
			tIfIdxList := bulkVlanInfo.VlanList[idx].IfIndexList
			vlanEnt.TagIfIdxMap = make(map[int32]bool)
			for i := 0; i < len(tIfIdxList); i++ {
				vlanEnt.TagIfIdxMap[tIfIdxList[i]] = true
			}
			server.infraData.vlanPropertyMap[vlanIfIdx] = vlanEnt
		}
		if more == false {
			break
		}
	}
}

func (server *OSPFV2Server) constructL3Infra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting L3 Interfaces")
	count := 100
	for {
		if server.asicdComm.asicdClient.ClientHdl == nil {
			server.logger.Err("Infra: Null asicd client handle")
			return
		}
		bulkInfo, _ := server.asicdComm.asicdClient.ClientHdl.GetBulkIPv4IntfState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}

		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIdx := bulkInfo.IPv4IntfStateList[i].IfIndex
			ipEnt := server.infraData.ipPropertyMap[ifIdx]
			ifType := uint32(asicdCommonDefs.GetIntfTypeFromIfIndex(ifIdx))
			ifId := uint32(asicdCommonDefs.GetIntfIdFromIfIndex(ifIdx))
			ip, mask, err := ParseCIDRToUint32(bulkInfo.IPv4IntfStateList[i].IpAddr)
			if err != nil {
				server.logger.Err("Error Parsing IP Address", err)
				continue
			}
			ipEnt.IpAddr = ip
			ipEnt.NetMask = mask
			ipEnt.IfId = ifId
			ipEnt.IfType = ifType
			ipEnt.IfName = bulkInfo.IPv4IntfStateList[i].IntfRef
			if bulkInfo.IPv4IntfStateList[i].OperState == "UP" {
				ipEnt.State = true
			} else {
				ipEnt.State = false
			}
			server.infraData.ipToIfIdxMap[ip] = ifIdx
			server.infraData.ipPropertyMap[ifIdx] = ipEnt
		}
		if more == false {
			break
		}
	}
}

func (server *OSPFV2Server) UpdateMtu(msg asicdCommonDefs.PortAttrChangeNotifyMsg) {
	ent, _ := server.infraData.portPropertyMap[msg.IfIndex]
	ent.Mtu = msg.Mtu
	server.infraData.portPropertyMap[msg.IfIndex] = ent
	for ifIdx, _ := range ent.IpIfIdxMap {
		ipEnt, exist := server.infraData.ipPropertyMap[ifIdx]
		if !exist {
			server.logger.Err("Something bad this should not happen", server.infraData.ipPropertyMap, server.infraData.portPropertyMap)
			continue
		}
		oldMtu := ipEnt.Mtu
		newMtu := server.getMTU(ipEnt.IfType, ifIdx, true)
		if newMtu != oldMtu {
			ipEnt.Mtu = newMtu
			server.infraData.ipPropertyMap[ifIdx] = ipEnt
			intfConfKey := IntfConfKey{
				IpAddr:  ipEnt.IpAddr,
				IntfIdx: 0,
			}
			intfConfEnt, exist := server.IntfConfMap[intfConfKey]
			if exist {
				if intfConfEnt.AdminState == true {
					server.logger.Info("StopIntfFSM ():")
					server.StopIntfFSM(intfConfKey)
				}
				intfConfEnt.Mtu = uint32(newMtu)
				server.IntfConfMap[intfConfKey] = intfConfEnt
				if intfConfEnt.AdminState == true {
					server.logger.Info("StartIntfFSM ():")
					server.StartIntfFSM(intfConfKey)
				}
			}
		}
	}
}

func (server *OSPFV2Server) UpdateLogicalIntfInfra(msg asicdCommonDefs.LogicalIntfNotifyMsg, msgType uint8) {
	//TODO: Loopback
	ifIdx := msg.IfIndex
	if msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE {
		lEnt, _ := server.infraData.logicalIntfPropertyMap[ifIdx]
		lEnt.IfName = msg.LogicalIntfName
		server.infraData.logicalIntfPropertyMap[ifIdx] = lEnt
	} else if msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_DELETE {
		delete(server.infraData.logicalIntfPropertyMap, ifIdx)
	} else if msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_UPDATE {
		lEnt, _ := server.infraData.logicalIntfPropertyMap[ifIdx]
		lEnt.IfName = msg.LogicalIntfName
		server.infraData.logicalIntfPropertyMap[ifIdx] = lEnt
	}
}

func (server *OSPFV2Server) UpdateIPv4Infra(msg asicdCommonDefs.IPv4IntfNotifyMsg, msgType uint8) {
	ifIdx := msg.IfIndex
	ifType := uint32(asicdCommonDefs.GetIntfTypeFromIfIndex(ifIdx))
	ip, mask, err := ParseCIDRToUint32(msg.IpAddr)
	if err != nil {
		server.logger.Err("Error Parsing IP Address", err)
		return
	}

	if msgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE {
		//msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE {
		server.logger.Info("Receive IPV4INTF_CREATE", msg)
		ifId := uint32(asicdCommonDefs.GetIntfIdFromIfIndex(ifIdx))
		ifName := msg.IntfRef
		mtu := server.getMTU(ifType, ifIdx, true)
		macAddr, err := getMACAddr(ifName)
		if err != nil {
			server.logger.Err("Unable to get the Mac Address", err)
			return
		}
		ipEnt, _ := server.infraData.ipPropertyMap[ifIdx]
		ipEnt.IfId = ifId
		ipEnt.IfName = ifName
		ipEnt.IfType = ifType
		ipEnt.IpAddr = ip
		ipEnt.NetMask = mask
		ipEnt.Mtu = mtu
		ipEnt.State = false
		ipEnt.MacAddr = macAddr
		server.infraData.ipToIfIdxMap[ip] = ifIdx
		server.infraData.ipPropertyMap[ifIdx] = ipEnt
		server.logger.Info("Ip to IfIdx Map:", server.infraData.ipToIfIdxMap)
		server.logger.Info("Ip property Map:", server.infraData.ipPropertyMap)
	} else {
		server.logger.Info("Receive IPV4INTF_DELETE", msg)
		server.getMTU(ifType, ifIdx, false)
		//Delete Interface Conf Map if exist
		delete(server.infraData.ipToIfIdxMap, ip)
		delete(server.infraData.ipPropertyMap, ifIdx)
	}

}

func (server *OSPFV2Server) ProcessIPv4StateChange(msg asicdCommonDefs.IPv4L3IntfStateNotifyMsg) {
	ifIdx := msg.IfIndex
	server.logger.Info("ProcessIPv4StateChange:", msg)
	ipEnt, exist := server.infraData.ipPropertyMap[ifIdx]
	if !exist {
		server.logger.Err("IPv4 entry doesnot exist in Ip Property Map")
		return
	}
	if msg.IfState == asicdCommonDefs.INTF_STATE_UP {
		ipEnt.State = true
		intfConfKey := IntfConfKey{
			IpAddr:  ipEnt.IpAddr,
			IntfIdx: 0,
		}
		intfConfEnt, exist := server.IntfConfMap[intfConfKey]
		if exist {
			intfConfEnt.OperState = true
			server.IntfConfMap[intfConfKey] = intfConfEnt
			if intfConfEnt.AdminState == true {
				server.logger.Info("StartIntfFSM ():")
				server.StartIntfFSM(intfConfKey)
			}
		}
	} else {
		//Stop State Machine if it exist and running
		//Stop Tx and Rx Packet
		ipEnt.State = false
		intfConfKey := IntfConfKey{
			IpAddr:  ipEnt.IpAddr,
			IntfIdx: 0,
		}
		intfConfEnt, exist := server.IntfConfMap[intfConfKey]
		if exist {
			if intfConfEnt.AdminState == true {
				server.logger.Info("StopIntfFSM ():")
				server.StopIntfFSM(intfConfKey)
			}
			intfConfEnt.OperState = false
			server.IntfConfMap[intfConfKey] = intfConfEnt
		}
	}
	server.infraData.ipPropertyMap[ifIdx] = ipEnt
}

func (server *OSPFV2Server) ProcessVlanNotify(msg asicdCommonDefs.VlanNotifyMsg, msgType uint8) {
	vlanIfIdx := msg.VlanIfIndex
	VlanName := msg.VlanName
	uIfIdxList := msg.UntagPorts
	tIfIdxList := msg.TagPorts
	if msgType == asicdCommonDefs.NOTIFY_VLAN_CREATE {
		vlanEnt, _ := server.infraData.vlanPropertyMap[vlanIfIdx]
		vlanEnt.UntagIfIdxMap = make(map[int32]bool)
		for idx := 0; idx < len(uIfIdxList); idx++ {
			vlanEnt.UntagIfIdxMap[uIfIdxList[idx]] = true
		}
		vlanEnt.TagIfIdxMap = make(map[int32]bool)
		for idx := 0; idx < len(tIfIdxList); idx++ {
			vlanEnt.TagIfIdxMap[tIfIdxList[idx]] = true
		}
		vlanEnt.Name = VlanName
		server.infraData.vlanPropertyMap[vlanIfIdx] = vlanEnt
	} else if msgType == asicdCommonDefs.NOTIFY_VLAN_DELETE {
		delete(server.infraData.vlanPropertyMap, vlanIfIdx)
	} else if msgType == asicdCommonDefs.NOTIFY_VLAN_UPDATE {
		//TODO
	}
}

func (server *OSPFV2Server) ProcessLagNotify(msg asicdCommonDefs.LagNotifyMsg, msgType uint8) {
	//TODO
}

/*
func (server *OSPFServer) computeMinMTU(IfType uint8, IfId uint16) int32 {
	var minMtu int32 = 10000             //in bytes
	if IfType == commonDefs.IfTypePort { // PHY
		ent, _ := server.portPropertyMap[int32(IfId)]
		minMtu = ent.Mtu
	} else if IfType == commonDefs.IfTypeVlan { // Vlan
		ent, _ := server.vlanPropertyMap[IfId]
		for _, portNum := range ent.UntagPorts {
			entry, _ := server.portPropertyMap[portNum]
			if minMtu > entry.Mtu {
				minMtu = entry.Mtu
			}
		}
	}
	return minMtu
}

func (server *OSPFServer) UpdateMtu(ifIndex int32, mtu int32) {
	ent, _ := server.portPropertyMap[ifIndex]
	ent.Mtu = mtu
	server.portPropertyMap[ifIndex] = ent
}

func (server *OSPFServer) updateIpPropertyMap(msg IPv4IntfNotifyMsg, msgType uint8) {
	ipAddr, _, _ := net.ParseCIDR(msg.IpAddr)
	ip := convertAreaOrRouterIdUint32(ipAddr.String())
	if msgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE { // Create IP
		ent := server.ipPropertyMap[ip]
		ent.IfId = msg.IfId
		ent.IfType = msg.IfType
		ent.IpAddr = msg.IpAddr
		ent.IpState = msg.IfState
		server.ipPropertyMap[ip] = ent
	} else { // Delete IP
		delete(server.ipPropertyMap, ip)
	}
}

func (server *OSPFServer) updateVlanPropertyMap(vlanNotifyMsg asicdCommonDefs.VlanNotifyMsg, msgType uint8) {
	if msgType == asicdCommonDefs.NOTIFY_VLAN_CREATE { // Create Vlan
		ent := server.vlanPropertyMap[vlanNotifyMsg.VlanId]
		ent.Name = vlanNotifyMsg.VlanName
		ent.UntagPorts = vlanNotifyMsg.UntagPorts
		server.vlanPropertyMap[vlanNotifyMsg.VlanId] = ent
	} else { // Delete Vlan
		delete(server.vlanPropertyMap, vlanNotifyMsg.VlanId)
	}
}

func (server *OSPFServer) updateLogicalIntfPropertyMap(logicalNotifyMsg asicdCommonDefs.LogicalIntfNotifyMsg,
	msgType uint8) {
	ifid := int32(asicdCommonDefs.GetIntfIdFromIfIndex(logicalNotifyMsg.IfIndex))
	if msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE {
		ent := server.logicalIntfPropertyMap[ifid]
		ent.Name = logicalNotifyMsg.LogicalIntfName
		server.logicalIntfPropertyMap[ifid] = ent
	} else { // delete interface
		delete(server.logicalIntfPropertyMap, ifid)
	}
}

func (server *OSPFServer) BuildOspfInfra() {
	server.constructPortInfra()
	server.constructVlanInfra()
	server.constructL3Infra()
}

func (server *OSPFServer) constructPortInfra() {
	server.getBulkPortState()
	server.getBulkPortConfig()
}

func (server *OSPFServer) constructVlanInfra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting Vlan Property")
	count := 100
	for {
		if server.asicdClient.ClientHdl == nil {
			server.logger.Err("Infra: Null client handle for asicd ")
			return
		}
		bulkVlanInfo, _ := server.asicdClient.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
		// Get bulk on vlan state can re-use curMark and count used
		// by get bulk vlan, as there is a 1:1 mapping in terms of cfg/state objs
		bulkVlanStateInfo, _ := server.asicdClient.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkVlanStateInfo == nil &&
			bulkVlanInfo == nil {
			break
		}
		objCnt := int(bulkVlanInfo.Count)
		more := bool(bulkVlanInfo.More)
		curMark = int(bulkVlanInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			vlanId := uint16(bulkVlanInfo.VlanList[i].VlanId)
			ent := server.vlanPropertyMap[vlanId]
			ent.UntagPorts = bulkVlanInfo.VlanList[i].UntagIfIndexList
			ent.Name = bulkVlanStateInfo.VlanStateList[i].VlanName
			server.vlanPropertyMap[vlanId] = ent
		}
		if more == false {
			break
		}
	}

}

func (server *OSPFServer) constructL3Infra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting L3 Interfaces")
	count := 100
	for {
		if server.asicdClient.ClientHdl == nil {
			server.logger.Err("Infra: Null asicd client handle")
			return
		}
		bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkIPv4IntfState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}

		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIdx := bulkInfo.IPv4IntfStateList[i].IfIndex
			ifType := uint8(asicdCommonDefs.GetIntfTypeFromIfIndex(ifIdx))
			ifId := uint16(asicdCommonDefs.GetIntfIdFromIfIndex(ifIdx))
			var ipv4IntfMsg IPv4IntfNotifyMsg
			ipv4IntfMsg.IpAddr = bulkInfo.IPv4IntfStateList[i].IpAddr
			ipv4IntfMsg.IfType = ifType
			ipv4IntfMsg.IfId = ifId
			if bulkInfo.IPv4IntfStateList[i].OperState == "UP" {
				ipv4IntfMsg.IfState = config.Intf_Up
			} else {
				ipv4IntfMsg.IfState = config.Intf_Down
			}
			server.updateIpPropertyMap(ipv4IntfMsg, asicdCommonDefs.NOTIFY_IPV4INTF_CREATE)
			mtu := server.computeMinMTU(ipv4IntfMsg.IfType, ipv4IntfMsg.IfId)
			server.createIPIntfConfMap(ipv4IntfMsg, mtu, ifIdx, broadcast)
		}
		if more == false {
			break
		}
	}
}

func (server *OSPFServer) getBulkPortState() {
	currMarker := asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS)
	if server.asicdClient.IsConnected {
		server.logger.Info("Calling asicd for getting port state")
		count := 100
		for {
			bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkPortState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo == nil {
				break
			}
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = asicdServices.Int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifIndex := bulkInfo.PortStateList[i].IfIndex
				ent := server.portPropertyMap[ifIndex]
				ent.Name = bulkInfo.PortStateList[i].Name
				server.portPropertyMap[ifIndex] = ent
			}
			if more == false {
				break
			}
		}
	}
}

func (server *OSPFServer) getBulkPortConfig() {
	currMarker := asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS)
	if server.asicdClient.IsConnected {
		server.logger.Info("Calling asicd for getting the Port Config")
		count := 100
		for {
			bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkPort(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo == nil {
				break
			}
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = asicdServices.Int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifIndex := bulkInfo.PortList[i].IfIndex
				ent := server.portPropertyMap[ifIndex]
				ent.Mtu = bulkInfo.PortList[i].Mtu
				ent.Speed = uint32(bulkInfo.PortList[i].Speed)
				server.portPropertyMap[ifIndex] = ent
			}
			if more == false {
				break
			}
		}
	}
}

func (server *OSPFServer) getLinuxIntfName(ifId int32, ifType uint8) (ifName string, err error) {
	server.logger.Err(fmt.Sprintln("IF : if id ", ifId, " ifType ", ifType))

	if ifType == commonDefs.IfTypeVlan { // Vlan
		ifName = server.vlanPropertyMap[uint16(ifId)].Name
	} else if ifType == commonDefs.IfTypePort { // PHY
		ifName = server.portPropertyMap[ifId].Name
	} else if ifType == commonDefs.IfTypeLoopback {
		ifName = server.logicalIntfPropertyMap[ifId].Name
	} else {
		ifName = ""
		err = errors.New("Invalid Interface Type")
	}
	return ifName, err
}

func (server *OSPFServer) getIntfCost(ifId uint16, ifType uint8) (ifCost uint32, err error) {
	if ifType == commonDefs.IfTypeVlan { // Vlan
		ifCost = DEFAULT_VLAN_COST
	} else if ifType == commonDefs.IfTypePort { // PHY
		speed := server.portPropertyMap[int32(ifId)].Speed
		if speed != 0 {
			ifCost = server.ospfGlobalConf.ReferenceBandwidth / speed
		} else {
			server.logger.Err(fmt.Sprintln("Port Speed for port = ", server.portPropertyMap[int32(ifId)].Name, " is zero, so something wrong"))
			ifCost = 0xff00
		}
	} else if ifType == commonDefs.IfTypeLoopback {
		ifCost = DEFAULT_VLAN_COST

	} else {
		ifCost = 0xff00
		err = errors.New("Invalid Interface Type")
	}
	return ifCost, err
}

func getMacAddrIntfName(ifName string) (macAddr net.HardwareAddr, err error) {

	ifi, err := net.InterfaceByName(ifName)
	if err != nil {
		return macAddr, err
	}
	macAddr = ifi.HardwareAddr
	return macAddr, nil
}

func (server *OSPFServer) getMacAddrLogicalIntf(ifName string) (macAddr net.HardwareAddr, err error) {
	if server.asicdClient.ClientHdl == nil {
		server.logger.Err("Infra: Null asicd client handle")
		return macAddr, errors.New("Null asicd handle")
	}
	portState, err := server.asicdClient.ClientHdl.GetLogicalIntfState(ifName)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Infra : Failed to get logical port config ", ifName))
		return macAddr, errors.New("Failed to get logical port config")
	}
	macAddr, err = net.ParseMAC(portState.SrcMac)
	if err != nil {
		server.logger.Err("Infra: Can not convert string to mac addr ", portState.SrcMac)
		return macAddr, errors.New("Infra : Failed to parse mac addr")
	}
	return macAddr, nil
}
func (server *OSPFServer) UpdateLogicalIntfInfra(msg asicdCommonDefs.LogicalIntfNotifyMsg,
	msgType uint8) {
	server.updateLogicalIntfPropertyMap(msg, msgType)
}

func (server *OSPFServer) UpdateVlanInfra(msg asicdCommonDefs.VlanNotifyMsg, msgType uint8) {
	server.updateVlanPropertyMap(msg, msgType)
}

func (server *OSPFServer) UpdateIPv4Infra(msg asicdCommonDefs.IPv4IntfNotifyMsg, msgType uint8) {
	var ipv4IntfMsg IPv4IntfNotifyMsg
	ipv4IntfMsg.IpAddr = msg.IpAddr
	ipv4IntfMsg.IfType = uint8(asicdCommonDefs.GetIntfTypeFromIfIndex(msg.IfIndex))
	ipv4IntfMsg.IfId = uint16(asicdCommonDefs.GetIntfIdFromIfIndex(msg.IfIndex))
	if msgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE ||
		msgType == asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE {
		server.logger.Info(fmt.Sprintln("Receive IPV4INTF_CREATE", msg))
		mtu := server.computeMinMTU(ipv4IntfMsg.IfType, ipv4IntfMsg.IfId)
		ipv4IntfMsg.IfState = config.Intf_Down
		server.createIPIntfConfMap(ipv4IntfMsg, mtu, msg.IfIndex, broadcast)
		server.updateIpPropertyMap(ipv4IntfMsg, msgType)
	} else {
		server.logger.Info(fmt.Sprintln("Receive IPV4INTF_DELETE", ipv4IntfMsg))
		server.deleteIPIntfConfMap(ipv4IntfMsg, msg.IfIndex)
		server.updateIpPropertyMap(ipv4IntfMsg, msgType)
	}

}
*/
