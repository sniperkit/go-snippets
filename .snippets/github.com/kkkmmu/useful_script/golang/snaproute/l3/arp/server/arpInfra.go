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
	"fmt"
	"github.com/google/gopacket/pcap"
	"net"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/commonDefs"
)

const (
	UNTAGGED bool = false
	TAGGED   bool = true
)

type L3VlanStruct struct {
	IpAddr  string
	IpNet   net.IPMask
	L3IfIdx int
}

type L3LagStruct struct {
	L3Vlan   L3VlanStruct
	LagIfIdx int
	VlanId   int
	TagFlag  bool
}

type L3PortStruct struct {
	L3Lag     L3LagStruct
	PortIfIdx int
}

type L3IntfProperty struct {
	Netmask   net.IPMask
	IpAddr    string
	IfName    string
	OperState bool
}

type L3PortProp struct {
	L3IfIdx  int
	IpAddr   string
	Netmask  net.IPMask
	LagIfIdx int
}

type PortProperty struct {
	IfName        string
	MacAddr       string
	UntagVlanId   int
	L3PortPropMap map[int]L3PortProp //VlanId
	CtrlCh        chan bool
	CtrlReplyCh   chan bool
	PcapHdl       *pcap.Handle
	baseFilter    string
	OperState     bool
}

type VlanProperty struct {
	IfName        string
	UntagIfIdxMap map[int]bool
	TagIfIdxMap   map[int]bool
}

type LagProperty struct {
	IfName      string
	UntagVlanId int
	VlanIdMap   map[int]bool
	PortMap     map[int]bool
}

type VirtualIntfProperty struct {
	IfName        string
	IfIndex       int32
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
}

func (server *ARPServer) isL3Intf(ifIdx int) bool {
	_, exist := server.l3IntfPropMap[ifIdx]
	return exist
}

func (server *ARPServer) buildArpInfra() {
	server.constructPortInfra()
	server.constructLagInfra()
	server.constructVlanInfra()
	server.constructL3Infra()
	server.constructVirtualIpInfra()
	server.dumpInfra()
}

func (server *ARPServer) constructPortInfra() {
	server.getBulkPortState()
	server.getBulkPortConfig()
	for portIfIdx, portEnt := range server.portPropMap {
		portEnt.baseFilter = server.constructBaseFilter(portEnt.MacAddr)
		server.portPropMap[portIfIdx] = portEnt
		if portEnt.OperState == true {
			server.EnableRxOnPort(portIfIdx)
		}
	}
}

func (server *ARPServer) GetMacAddr(l3IfIdx int) string {
	ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(l3IfIdx))
	switch ifType {
	case asicdClntDefs.IfTypeVlan:
		vlanEnt, _ := server.vlanPropMap[l3IfIdx]
		for uIfIdx, _ := range vlanEnt.UntagIfIdxMap {
			uIfType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(uIfIdx))
			switch uIfType {
			case asicdClntDefs.IfTypeLag:
				lagEnt, _ := server.lagPropMap[uIfIdx]
				for portIfIdx, _ := range lagEnt.PortMap {
					portEnt, _ := server.portPropMap[portIfIdx]
					return portEnt.MacAddr
				}
			case asicdClntDefs.IfTypePort:
				portEnt, _ := server.portPropMap[uIfIdx]
				return portEnt.MacAddr
			}
		}
		for tIfIdx, _ := range vlanEnt.TagIfIdxMap {
			tIfType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(tIfIdx))
			switch tIfType {
			case asicdClntDefs.IfTypeLag:
				lagEnt, _ := server.lagPropMap[tIfIdx]
				for portIfIdx, _ := range lagEnt.PortMap {
					portEnt, _ := server.portPropMap[portIfIdx]
					return portEnt.MacAddr
				}
			case asicdClntDefs.IfTypePort:
				portEnt, _ := server.portPropMap[tIfIdx]
				return portEnt.MacAddr
			}
		}
	case asicdClntDefs.IfTypeLag:
		lagEnt, _ := server.lagPropMap[l3IfIdx]
		for portIfIdx, _ := range lagEnt.PortMap {
			portEnt, _ := server.portPropMap[portIfIdx]
			return portEnt.MacAddr
		}
	case asicdClntDefs.IfTypePort:
		portEnt, _ := server.portPropMap[l3IfIdx]
		return portEnt.MacAddr
	}
	return ""
}

func (server *ARPServer) getBulkPortState() {
	curMark := server.AsicdPlugin.GetMinSysPort()
	server.logger.Debug("Calling Asicd for getting Port Property")
	count := 100
	for {
		bulkInfo, _ := server.AsicdPlugin.GetBulkPortState(curMark, count)
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIndex := int(bulkInfo.PortStateList[i].IfIndex)
			ent := server.portPropMap[ifIndex]
			ent.IfName = bulkInfo.PortStateList[i].Name
			ent.UntagVlanId = server.SysRsvdVlan
			ent.L3PortPropMap = make(map[int]L3PortProp)
			ent.PcapHdl = nil
			switch bulkInfo.PortStateList[i].OperState {
			case "UP":
				ent.OperState = true
			case "DOWN":
				ent.OperState = false
			default:
				server.logger.Err("Invalid OperState for the port",
					bulkInfo.PortStateList[i].OperState, ent.IfName)
			}
			server.portPropMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
}

func (server *ARPServer) getBulkPortConfig() {
	curMark := server.AsicdPlugin.GetMinSysPort()
	server.logger.Debug("Calling Asicd for getting Port Property")
	count := 100
	for {
		bulkInfo, _ := server.AsicdPlugin.GetBulkPort(curMark, count)
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIndex := int(bulkInfo.PortList[i].IfIndex)
			ent := server.portPropMap[ifIndex]
			ent.MacAddr = bulkInfo.PortList[i].MacAddr
			ent.CtrlCh = make(chan bool)
			ent.CtrlReplyCh = make(chan bool)
			server.portPropMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
}

func (server *ARPServer) constructLagInfra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting Lag Property")
	count := 100
	for {
		bulkLagInfo, _ := server.AsicdPlugin.GetBulkLag(curMark, count)
		if bulkLagInfo == nil {
			break
		}
		objCnt := int(bulkLagInfo.Count)
		more := bool(bulkLagInfo.More)
		curMark = int(bulkLagInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			lagIfIdx := int(bulkLagInfo.LagList[i].LagIfIndex)
			lagEnt := server.lagPropMap[lagIfIdx]
			lagEnt.IfName = bulkLagInfo.LagList[i].LagName
			lagEnt.PortMap = make(map[int]bool)
			lagEnt.VlanIdMap = make(map[int]bool)
			lagEnt.UntagVlanId = server.SysRsvdVlan
			ifIndexList := bulkLagInfo.LagList[i].IfIndexList
			for idx := 0; idx < len(ifIndexList); idx++ {
				ifIdx := int(ifIndexList[idx])
				lagEnt.PortMap[ifIdx] = true
			}
			server.lagPropMap[lagIfIdx] = lagEnt
		}
		if more == false {
			break
		}
	}
}

func (server *ARPServer) constructVlanInfra() {
	curMark := 0
	server.logger.Debug("Calling Asicd for getting Vlan Property")
	count := 100
	for {
		bulkVlanInfo, _ := server.AsicdPlugin.GetBulkVlan(curMark, count)
		if bulkVlanInfo == nil {
			break
		}
		/*
		* Getbulk on vlan state assuming indexes are one to on mapped
		* between state and config object
		 */
		bulkVlanStateInfo, _ := server.AsicdPlugin.GetBulkVlanState(curMark, count)
		if bulkVlanStateInfo == nil {
			break
		}
		objCnt := int(bulkVlanInfo.Count)
		more := bool(bulkVlanInfo.More)
		curMark = int(bulkVlanInfo.EndIdx)
		for idx := 0; idx < objCnt; idx++ {
			vlanIfIdx := int(bulkVlanStateInfo.VlanStateList[idx].IfIndex)
			vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
			vlanEnt := server.vlanPropMap[vlanIfIdx]
			vlanEnt.IfName = bulkVlanStateInfo.VlanStateList[idx].VlanName
			uIfIdxList := bulkVlanInfo.VlanList[idx].UntagIfIndexList
			vlanEnt.UntagIfIdxMap = make(map[int]bool)
			for i := 0; i < len(uIfIdxList); i++ {
				ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(uIfIdxList[i])
				if ifType == asicdClntDefs.IfTypeLag {
					server.updateLagForVlan(int(uIfIdxList[i]), vlanId, UNTAGGED, true)
				}
				vlanEnt.UntagIfIdxMap[int(uIfIdxList[i])] = true
			}
			tIfIdxList := bulkVlanInfo.VlanList[idx].IfIndexList
			vlanEnt.TagIfIdxMap = make(map[int]bool)
			for i := 0; i < len(tIfIdxList); i++ {
				ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(tIfIdxList[i])
				if ifType == asicdClntDefs.IfTypeLag {
					server.updateLagForVlan(int(tIfIdxList[i]), vlanId, TAGGED, true)
				}
				vlanEnt.TagIfIdxMap[int(tIfIdxList[i])] = true
			}
			server.vlanPropMap[vlanIfIdx] = vlanEnt
		}
		if more == false {
			break
		}
	}
}

func (server *ARPServer) constructL3Infra() {
	currMark := 0
	server.logger.Debug("Calling Asicd for getting L3 Interfaces")
	count := 100
	for {
		bulkIPv4Info, _ := server.AsicdPlugin.GetBulkIPv4IntfState(currMark, count)
		if bulkIPv4Info == nil {
			break
		}

		objCnt := int(bulkIPv4Info.Count)
		more := bool(bulkIPv4Info.More)
		currMark = int(bulkIPv4Info.EndIdx)
		for i := 0; i < objCnt; i++ {
			server.processIPv4IntfCreate(bulkIPv4Info.IPv4IntfStateList[i].IpAddr, bulkIPv4Info.IPv4IntfStateList[i].IfIndex)
			l3IfIdx := int(bulkIPv4Info.IPv4IntfStateList[i].IfIndex)
			l3Ent := server.l3IntfPropMap[l3IfIdx]
			switch bulkIPv4Info.IPv4IntfStateList[i].OperState {
			case "UP":
				l3Ent.OperState = true
			case "DOWN":
				l3Ent.OperState = false
			default:
				server.logger.Err("Invalid OperState for the L3 Interface")
			}
			server.l3IntfPropMap[l3IfIdx] = l3Ent
		}
		if more == false {
			break
		}
	}
}

func (server *ARPServer) updateVlanWithL3(l3Vlan L3VlanStruct) string {
	vlanEnt := server.vlanPropMap[l3Vlan.L3IfIdx]
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(l3Vlan.L3IfIdx))
	l3Lag := L3LagStruct{
		L3Vlan: l3Vlan,
		VlanId: vlanId,
	}
	for uIfIdx, _ := range vlanEnt.UntagIfIdxMap {
		l3Lag.TagFlag = UNTAGGED
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(uIfIdx))
		switch ifType {
		case asicdClntDefs.IfTypeLag:
			l3Lag.LagIfIdx = uIfIdx
			server.updateLagWithL3(l3Lag)
		case asicdClntDefs.IfTypePort:
			l3Lag.LagIfIdx = -1
			l3Port := L3PortStruct{
				L3Lag:     l3Lag,
				PortIfIdx: uIfIdx,
			}
			server.updatePortWithL3(l3Port)
		}
	}
	for tIfIdx, _ := range vlanEnt.TagIfIdxMap {
		l3Lag.TagFlag = TAGGED
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(tIfIdx))
		switch ifType {
		case asicdClntDefs.IfTypeLag:
			l3Lag.LagIfIdx = tIfIdx
			server.updateLagWithL3(l3Lag)
		case asicdClntDefs.IfTypePort:
			l3Lag.LagIfIdx = -1
			l3Port := L3PortStruct{
				L3Lag:     l3Lag,
				PortIfIdx: tIfIdx,
			}
			server.updatePortWithL3(l3Port)
		}
	}
	return vlanEnt.IfName
}

func (server *ARPServer) updateLagWithL3(l3Lag L3LagStruct) string {
	lagEnt := server.lagPropMap[l3Lag.LagIfIdx]
	l3Port := L3PortStruct{
		L3Lag: l3Lag,
	}
	//lagEnt.VlanIdMap[l3Lag.VlanId] = true
	for portIfIdx, _ := range lagEnt.PortMap {
		l3Port.PortIfIdx = portIfIdx
		server.updatePortWithL3(l3Port)
	}
	return lagEnt.IfName
}

func (server *ARPServer) updatePortWithL3(l3Port L3PortStruct) string {
	portEnt := server.portPropMap[l3Port.PortIfIdx]
	ifName := portEnt.IfName
	l3PortPropEnt := portEnt.L3PortPropMap[l3Port.L3Lag.VlanId]
	l3PortPropEnt.IpAddr = l3Port.L3Lag.L3Vlan.IpAddr
	l3PortPropEnt.Netmask = l3Port.L3Lag.L3Vlan.IpNet
	l3PortPropEnt.LagIfIdx = l3Port.L3Lag.LagIfIdx
	l3PortPropEnt.L3IfIdx = l3Port.L3Lag.L3Vlan.L3IfIdx
	portEnt.L3PortPropMap[l3Port.L3Lag.VlanId] = l3PortPropEnt
	if l3Port.L3Lag.TagFlag == UNTAGGED {
		portEnt.UntagVlanId = l3Port.L3Lag.VlanId
	}
	server.portPropMap[l3Port.PortIfIdx] = portEnt
	return ifName
}

func (server *ARPServer) processIPv4NbrMacMove(msg asicdClntDefs.IPv4NbrMacMoveNotifyMsg) {
	server.arpEntryMacMoveCh <- msg
}

func (server *ARPServer) updateIPv4Infra(msg asicdClntDefs.IPv4IntfNotifyMsg) {
	if msg.MsgType == asicdClntDefs.NOTIFY_IPV4INTF_CREATE {
		server.processIPv4IntfCreate(msg.IpAddr, msg.IfIndex)
	} else {
		server.processIPv4IntfDelete(msg.IpAddr, msg.IfIndex)
	}
}

func (server *ARPServer) processVirtualIntfEvent(msg asicdClntDefs.IPv4VirtualIntfNotifyMsg) {
	if msg.MsgType == asicdClntDefs.NOTIFY_IPV4VIRTUAL_INTF_CREATE {
		virEntry := server.virtualIntfPropMap[msg.IfIndex]
		virEntry.IfIndex = msg.IfIndex
		virEntry.ParentIfIndex = msg.ParentIfIndex
		virEntry.IpAddr = msg.IpAddr
		virEntry.MacAddr = msg.MacAddr
		virEntry.IfName = msg.IfName
		server.virtualIntfPropMap[msg.IfIndex] = virEntry
	} else {
		delete(server.virtualIntfPropMap, msg.IfIndex)
	}
}

func (server *ARPServer) processIPv4IntfCreate(IpAddr string, IfIndex int32) {
	var ifName string
	ip, ipNet, _ := net.ParseCIDR(IpAddr)
	l3IfIdx := int(IfIndex)
	ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(IfIndex)
	ipAddr := ip.String()
	l3Vlan := L3VlanStruct{
		IpAddr:  ipAddr,
		IpNet:   ipNet.Mask,
		L3IfIdx: l3IfIdx,
	}
	switch ifType {
	case asicdClntDefs.IfTypeVlan:
		ifName = server.updateVlanWithL3(l3Vlan)
	case asicdClntDefs.IfTypeLag:
		l3Lag := L3LagStruct{
			L3Vlan:   l3Vlan,
			LagIfIdx: l3IfIdx,
			VlanId:   server.SysRsvdVlan,
			TagFlag:  UNTAGGED,
		}
		ifName = server.updateLagWithL3(l3Lag)
	case asicdClntDefs.IfTypePort:
		l3Lag := L3LagStruct{
			L3Vlan:   l3Vlan,
			LagIfIdx: -1,
			VlanId:   server.SysRsvdVlan,
			TagFlag:  UNTAGGED,
		}
		l3Port := L3PortStruct{
			L3Lag:     l3Lag,
			PortIfIdx: l3IfIdx,
		}
		ifName = server.updatePortWithL3(l3Port)
	}
	l3Ent := server.l3IntfPropMap[l3IfIdx]
	l3Ent.Netmask = ipNet.Mask
	l3Ent.IpAddr = ipAddr
	l3Ent.IfName = ifName
	server.l3IntfPropMap[l3IfIdx] = l3Ent
}

func (server *ARPServer) processIPv4IntfDelete(IpAddr string, IfIndex int32) {
	l3IfIdx := int(IfIndex)
	ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(IfIndex)
	switch ifType {
	case asicdClntDefs.IfTypeVlan:
		server.updateVlanForL3Delete(l3IfIdx)
	case asicdClntDefs.IfTypeLag:
		server.updateLagForL3Delete(l3IfIdx, server.SysRsvdVlan, l3IfIdx)
	case asicdClntDefs.IfTypePort:
		server.updatePortForL3Delete(l3IfIdx, server.SysRsvdVlan, -1, l3IfIdx)
	}
	delete(server.l3IntfPropMap, l3IfIdx)
}

func (server *ARPServer) updateVlanForL3Delete(vlanIfIdx int) {
	vlanEnt := server.vlanPropMap[vlanIfIdx]
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
	for uIfIdx, _ := range vlanEnt.UntagIfIdxMap {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(uIfIdx))
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForL3Delete(vlanIfIdx, vlanId, uIfIdx)
		} else {
			server.updatePortForL3Delete(vlanIfIdx, vlanId, -1, uIfIdx)
		}
	}
	for tIfIdx, _ := range vlanEnt.TagIfIdxMap {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(int32(tIfIdx))
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForL3Delete(vlanIfIdx, vlanId, tIfIdx)
		} else {
			server.updatePortForL3Delete(vlanIfIdx, vlanId, -1, tIfIdx)
		}
	}
}

func (server *ARPServer) updateLagForL3Delete(l3IfIdx int, vlanId int, lagIfIdx int) {
	lagEnt := server.lagPropMap[lagIfIdx]
	//delete(lagEnt.VlanIdMap, vlanId)
	for portIfIdx, _ := range lagEnt.PortMap {
		server.updatePortForL3Delete(l3IfIdx, vlanId, lagIfIdx, portIfIdx)
	}
}

func (server *ARPServer) updatePortForL3Delete(l3IfIdx, vlanId, lagIfIdx, portIfIdx int) {
	portEnt := server.portPropMap[portIfIdx]
	delete(portEnt.L3PortPropMap, vlanId)
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) processIPv4L3StateChange(msg asicdClntDefs.IPv4L3IntfStateNotifyMsg) {
	ifIdx := int(msg.IfIndex)
	if msg.IfState == 0 {
		server.DisableArpOnL3(ifIdx)
		server.DisableL3(ifIdx)
	} else {
		server.EnableL3(ifIdx)
		go server.SendArpProbe(ifIdx)
	}
}

func (server *ARPServer) sendVipGarp(ifIndex int32) {
	virEntry, exists := server.virtualIntfPropMap[ifIndex]
	if exists {
		ip, _, _ := net.ParseCIDR(virEntry.IpAddr)
		ip = ip.To4()
		server.GarpEntryCh <- &GarpEntry{virEntry.IfName, virEntry.MacAddr, ip.String()}
		//go server.SendGarp(virEntry.IfName, virEntry.MacAddr, ip.String())
	}
}

func (server *ARPServer) processVirtualIntfStateEvent(msg asicdClntDefs.IPv4VirtualIntfStateNotifyMsg) {
	if msg.IfState == 0 {
		server.RestoreBPFFilter(msg)
	} else {
		server.UpdateBPFFilter(msg)
		server.sendVipGarp(msg.IfIndex)
	}
}

func (server *ARPServer) processL2StateChange(msg asicdClntDefs.L2IntfStateNotifyMsg) {
	ifIdx := int(msg.IfIndex)
	ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(msg.IfIndex)
	if msg.IfState == 0 {
		switch ifType {
		case asicdClntDefs.IfTypeVlan:
			//server.DisableArpOnVlan(ifIdx)
		case asicdClntDefs.IfTypeLag:
			//server.DisableArpOnLag(ifIdx)
		case asicdClntDefs.IfTypePort:
			server.DisableRxOnPort(ifIdx)
			server.DisablePort(ifIdx)
		}
	} else {
		switch ifType {
		case asicdClntDefs.IfTypePort:
			server.EnablePort(ifIdx)
			server.EnableRxOnPort(ifIdx)
		}

	}
}

func (server *ARPServer) DisableArpOnL3(l3IfIdx int) {
	server.arpEntryDeleteCh <- DeleteArpEntryMsg{
		Type:  DeleteBasedOnL3,
		IfIdx: l3IfIdx,
	}
}

/*
func (server *ARPServer) DisableArpOnVlan(l3IfIdx int) {
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(l3IfIdx))
	server.arpEntryDeleteCh <- DeleteArpEntryMsg{
		Type:  DeleteBasedOnVlan,
		IfIdx: vlanId,
	}
}

func (server *ARPServer) DisableArpOnLag(lagIfIdx int) {
	lagEnt := server.lagPropMap[lagIfIdx]
	for portIfIdx, _ := range lagEnt.PortMap {
		server.DisableArpOnPort(portIfIdx)
	}
}
*/

func (server *ARPServer) DisableRxOnPort(portIfIdx int) {
	server.logger.Debug("Disabling Arp on Port", portIfIdx)
	portEnt := server.portPropMap[portIfIdx]
	if portEnt.OperState == true {
		server.logger.Debug("Disabling Arp on Port: OperState=true", portIfIdx)
		if portEnt.PcapHdl != nil {
			server.logger.Debug("Disabling Arp on Port: PcapHdl != nil", portIfIdx)
			portEnt.CtrlCh <- true
			<-portEnt.CtrlReplyCh
			portEnt.PcapHdl.Close()
			portEnt.PcapHdl = nil
			server.arpEntryDeleteCh <- DeleteArpEntryMsg{
				Type:  DeleteBasedOnPort,
				IfIdx: portIfIdx,
			}
		}
	}
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) constructBaseFilter(macAddr string) string {
	return fmt.Sprintf("(not ether proto 0x8809 and not (ether src %s", macAddr)
}

func (server *ARPServer) EnableRxOnPort(portIfIdx int) {
	var err error
	portEnt := server.portPropMap[portIfIdx]
	if portEnt.OperState == true {
		if portEnt.PcapHdl == nil {
			portEnt.PcapHdl, err = server.StartArpRxTx(portEnt.IfName, portEnt.MacAddr, portEnt.baseFilter)
			if err != nil {
				server.logger.Err("Error opening pcap handle on", portEnt.IfName, err)
				return
			}
			server.portPropMap[portIfIdx] = portEnt
			go server.processRxPkts(portIfIdx)
		}
	}
}

func (server *ARPServer) DisablePort(portIfIdx int) {
	portEnt := server.portPropMap[portIfIdx]
	portEnt.OperState = false
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) EnablePort(portIfIdx int) {
	portEnt := server.portPropMap[portIfIdx]
	portEnt.OperState = true
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) EnableL3(l3IfIdx int) {
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	l3Ent.OperState = true
	server.l3IntfPropMap[l3IfIdx] = l3Ent
}

func (server *ARPServer) DisableL3(l3IfIdx int) {
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	l3Ent.OperState = false
	server.l3IntfPropMap[l3IfIdx] = l3Ent
}

func (server *ARPServer) updateVlanInfra(msg asicdClntDefs.VlanNotifyMsg) {
	vlanId := int(msg.VlanId)
	vlanIfIdx := int(server.AsicdPlugin.GetIfIndexFromIntfIdAndIntfType(vlanId, asicdClntDefs.IfTypeVlan))
	uIfIdxList := msg.UntagPorts
	tIfIdxList := msg.TagPorts
	vlanName := msg.VlanName
	if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_CREATE {
		server.processVlanCreate(vlanName, vlanIfIdx, uIfIdxList, tIfIdxList, true)
	} else if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_UPDATE {
		server.processVlanUpdate(vlanName, vlanIfIdx, uIfIdxList, tIfIdxList)
	} else if msg.MsgType == asicdClntDefs.NOTIFY_VLAN_DELETE {
		server.processVlanDelete(vlanName, vlanIfIdx, uIfIdxList, tIfIdxList, true)
	}
}

func (server *ARPServer) processVlanCreate(vlanName string, vlanIfIdx int, uIfIdxList, tIfIdxList []int32, flag bool) {
	vlanEnt, _ := server.vlanPropMap[vlanIfIdx]
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
	if flag == true {
		vlanEnt.UntagIfIdxMap = nil
		vlanEnt.UntagIfIdxMap = make(map[int]bool)
		vlanEnt.TagIfIdxMap = nil
		vlanEnt.TagIfIdxMap = make(map[int]bool)
	}
	for idx := 0; idx < len(uIfIdxList); idx++ {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(uIfIdxList[idx])
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForVlan(int(uIfIdxList[idx]), vlanId, UNTAGGED, true)
		}
		vlanEnt.UntagIfIdxMap[int(uIfIdxList[idx])] = true
	}
	for idx := 0; idx < len(tIfIdxList); idx++ {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(tIfIdxList[idx])
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForVlan(int(tIfIdxList[idx]), vlanId, TAGGED, true)
		}
		vlanEnt.TagIfIdxMap[int(tIfIdxList[idx])] = true
	}
	vlanEnt.IfName = vlanName
	server.vlanPropMap[vlanIfIdx] = vlanEnt
}

func (server *ARPServer) updateLagForVlan(lagIfIdx, vlanId int, TagFlag, create bool) {
	lagEnt, _ := server.lagPropMap[lagIfIdx]
	if create == true {
		if TagFlag == UNTAGGED {
			lagEnt.UntagVlanId = vlanId
		}
		lagEnt.VlanIdMap[vlanId] = true
	} else {
		if TagFlag == UNTAGGED {
			lagEnt.UntagVlanId = server.SysRsvdVlan
		}
		delete(lagEnt.VlanIdMap, vlanId)
	}
	server.lagPropMap[lagIfIdx] = lagEnt
}

func (server *ARPServer) processVlanDelete(vlanName string, vlanIfIdx int, uIfIdxList, tIfIdxList []int32, flag bool) {
	vlanEnt, _ := server.vlanPropMap[vlanIfIdx]
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
	for _, uIfIdx := range uIfIdxList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(uIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForVlan(int(uIfIdx), vlanId, UNTAGGED, false)
		}
		delete(vlanEnt.UntagIfIdxMap, int(uIfIdx))
	}
	for _, tIfIdx := range tIfIdxList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(tIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.updateLagForVlan(int(tIfIdx), vlanId, TAGGED, false)
		}
		delete(vlanEnt.TagIfIdxMap, int(tIfIdx))
	}
	if flag == true {
		vlanEnt.UntagIfIdxMap = nil
		vlanEnt.TagIfIdxMap = nil
		delete(server.vlanPropMap, vlanIfIdx)
	}
}

func (server *ARPServer) processVlanUpdate(vlanName string, vlanIfIdx int, uIfIdxList, tIfIdxList []int32) {
	var uIfIdxDelList []int32
	var tIfIdxDelList []int32
	var uIfIdxCreateList []int32
	var tIfIdxCreateList []int32
	newUIfIdxMap := make(map[int]bool)
	for idx := 0; idx < len(uIfIdxList); idx++ {
		newUIfIdxMap[int(uIfIdxList[idx])] = true
	}
	newTIfIdxMap := make(map[int]bool)
	for idx := 0; idx < len(tIfIdxList); idx++ {
		newTIfIdxMap[int(tIfIdxList[idx])] = true
	}

	vlanEnt, _ := server.vlanPropMap[vlanIfIdx]
	for oldUIfIdx, _ := range vlanEnt.UntagIfIdxMap {
		_, exist := newUIfIdxMap[oldUIfIdx]
		if !exist {
			uIfIdxDelList = append(uIfIdxDelList, int32(oldUIfIdx))
		} else {
			delete(newUIfIdxMap, oldUIfIdx)
		}
	}
	for newUIfIdx, _ := range newUIfIdxMap {
		uIfIdxCreateList = append(uIfIdxCreateList, int32(newUIfIdx))
	}
	for oldTIfIdx, _ := range vlanEnt.TagIfIdxMap {
		_, exist := newTIfIdxMap[oldTIfIdx]
		if !exist {
			tIfIdxDelList = append(tIfIdxDelList, int32(oldTIfIdx))
		} else {
			delete(newTIfIdxMap, oldTIfIdx)
		}
	}
	for newTIfIdx, _ := range newTIfIdxMap {
		tIfIdxCreateList = append(tIfIdxCreateList, int32(newTIfIdx))
	}
	server.processVlanDelete(vlanName, vlanIfIdx, uIfIdxDelList, tIfIdxDelList, false)
	server.processVlanCreate(vlanName, vlanIfIdx, uIfIdxCreateList, tIfIdxCreateList, false)
	if server.isL3Intf(vlanIfIdx) {
		server.updateInfraWithL3VlanUpdate(vlanIfIdx, uIfIdxDelList, tIfIdxDelList, uIfIdxCreateList, tIfIdxCreateList)
	}
}

func (server *ARPServer) updateInfraWithL3VlanUpdate(vlanIfIdx int, uIfIdxDelList, tIfIdxDelList, uIfIdxCreateList, tIfIdxCreateList []int32) {
	for _, uIfIdx := range uIfIdxDelList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(uIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.deleteLagFromL3Vlan(vlanIfIdx, int(uIfIdx), UNTAGGED)
		} else {
			server.deletePortFromL3Vlan(vlanIfIdx, -1, int(uIfIdx), UNTAGGED)
		}
	}
	for _, tIfIdx := range tIfIdxDelList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(tIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.deleteLagFromL3Vlan(vlanIfIdx, int(tIfIdx), TAGGED)
		} else {
			server.deletePortFromL3Vlan(vlanIfIdx, -1, int(tIfIdx), TAGGED)
		}
	}
	for _, uIfIdx := range uIfIdxCreateList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(uIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.createLagFromL3Vlan(vlanIfIdx, int(uIfIdx), UNTAGGED)
		} else {
			server.createPortFromL3Vlan(vlanIfIdx, -1, int(uIfIdx), UNTAGGED)
		}
	}
	for _, tIfIdx := range tIfIdxDelList {
		ifType := server.AsicdPlugin.GetIntfTypeFromIfIndex(tIfIdx)
		if ifType == asicdClntDefs.IfTypeLag {
			server.createLagFromL3Vlan(vlanIfIdx, int(tIfIdx), TAGGED)
		} else {
			server.createPortFromL3Vlan(vlanIfIdx, -1, int(tIfIdx), TAGGED)
		}
	}
}

func (server *ARPServer) deleteLagFromL3Vlan(vlanIfIdx, lagIfIdx int, TagFlag bool) {
	lagEnt, exist := server.lagPropMap[lagIfIdx]
	if !exist {
		return
	}
	for portIfIdx, _ := range lagEnt.PortMap {
		server.deletePortFromL3Vlan(vlanIfIdx, lagIfIdx, portIfIdx, TagFlag)
	}
}

func (server *ARPServer) createLagFromL3Vlan(vlanIfIdx, lagIfIdx int, TagFlag bool) {
	lagEnt, exist := server.lagPropMap[lagIfIdx]
	if !exist {
		return
	}
	for portIfIdx, _ := range lagEnt.PortMap {
		server.createPortFromL3Vlan(vlanIfIdx, lagIfIdx, portIfIdx, TagFlag)
	}
}

func (server *ARPServer) deletePortFromL3Vlan(vlanIfIdx, lagIfIdx, portIfIdx int, TagFlag bool) {
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
	portEnt, _ := server.portPropMap[portIfIdx]
	if TagFlag == UNTAGGED {
		portEnt.UntagVlanId = server.SysRsvdVlan
	}
	delete(portEnt.L3PortPropMap, vlanId)
	server.portPropMap[portIfIdx] = portEnt
	server.arpEntryDeleteCh <- DeleteArpEntryMsg{
		Type:  DeleteBasedOnPort,
		IfIdx: portIfIdx,
	}
}

func (server *ARPServer) createPortFromL3Vlan(vlanIfIdx, lagIfIdx, portIfIdx int, TagFlag bool) {
	vlanId := server.AsicdPlugin.GetIntfIdFromIfIndex(int32(vlanIfIdx))
	l3Ent, _ := server.l3IntfPropMap[vlanIfIdx]
	portEnt, _ := server.portPropMap[portIfIdx]
	l3PortPropEnt, _ := portEnt.L3PortPropMap[vlanId]
	l3PortPropEnt.IpAddr = l3Ent.IpAddr
	l3PortPropEnt.L3IfIdx = vlanIfIdx
	l3PortPropEnt.Netmask = l3Ent.Netmask
	l3PortPropEnt.LagIfIdx = lagIfIdx
	portEnt.L3PortPropMap[vlanId] = l3PortPropEnt
	if TagFlag == UNTAGGED {
		portEnt.UntagVlanId = vlanId
	}
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) updateLagInfra(msg asicdClntDefs.LagNotifyMsg) {
	lagIfIdx := int(msg.IfIndex)
	portList := msg.IfIndexList
	if msg.MsgType == asicdClntDefs.NOTIFY_LAG_CREATE {
		server.processLagCreate(msg.LagName, lagIfIdx, portList, true)
	} else if msg.MsgType == asicdClntDefs.NOTIFY_LAG_UPDATE {
		server.processLagUpdate(lagIfIdx, portList)
	} else if msg.MsgType == asicdClntDefs.NOTIFY_LAG_DELETE {
		server.processLagDelete(lagIfIdx, portList, true)
	}
}

func (server *ARPServer) processLagCreate(lagName string, lagIfIdx int, portList []int32, createFlag bool) {
	lagEnt, _ := server.lagPropMap[lagIfIdx]
	if createFlag == true {
		lagEnt.PortMap = nil
		lagEnt.PortMap = make(map[int]bool)
		lagEnt.VlanIdMap = nil
		lagEnt.VlanIdMap = make(map[int]bool)
		lagEnt.IfName = ""
		lagEnt.IfName = lagName
	}
	for _, portIfIdx := range portList {
		lagEnt.PortMap[int(portIfIdx)] = true
	}
	server.lagPropMap[lagIfIdx] = lagEnt
}

func (server *ARPServer) processLagDelete(lagIfIdx int, portList []int32, deleteFlag bool) {
	lagEnt, _ := server.lagPropMap[lagIfIdx]
	for _, portIfIdx := range portList {
		delete(lagEnt.PortMap, int(portIfIdx))
	}
	if deleteFlag == true {
		lagEnt.PortMap = nil
		lagEnt.VlanIdMap = nil
		delete(server.lagPropMap, lagIfIdx)
	}
}

func (server *ARPServer) processLagUpdate(lagIfIdx int, portList []int32) {
	var delPortList []int32
	var createPortList []int32

	newPortMap := make(map[int]bool)
	for _, portIfIdx := range portList {
		newPortMap[int(portIfIdx)] = true
	}

	lagEnt, _ := server.lagPropMap[lagIfIdx]
	for oldPortIfIdx, _ := range lagEnt.PortMap {
		_, exist := newPortMap[oldPortIfIdx]
		if !exist {
			delPortList = append(delPortList, int32(oldPortIfIdx))
		} else {
			delete(newPortMap, oldPortIfIdx)
		}
	}
	for portIfdx, _ := range newPortMap {
		createPortList = append(createPortList, int32(portIfdx))
	}

	//Dummy arg not used during updates
	lagName := ""
	server.processLagCreate(lagName, lagIfIdx, createPortList, false)
	server.processLagDelete(lagIfIdx, delPortList, false)
	server.updateInfraWithL3LagUpdate(lagIfIdx, createPortList, delPortList)
}

func (server *ARPServer) updateInfraWithL3LagUpdate(lagIfIdx int, createPortList, delPortList []int32) {
	lagEnt, _ := server.lagPropMap[lagIfIdx]
	_, exist := server.l3IntfPropMap[lagIfIdx]
	if exist {
		for _, portIfIdx := range createPortList {
			server.createPortFromL3Lag(lagIfIdx, server.SysRsvdVlan, lagIfIdx, int(portIfIdx), UNTAGGED)
		}
		for _, portIfIdx := range delPortList {
			server.deletePortFromL3Lag(lagIfIdx, server.SysRsvdVlan, lagIfIdx, int(portIfIdx), UNTAGGED)
		}
	} else {
		for vlanId, _ := range lagEnt.VlanIdMap {
			vlanIfIdx := int(server.AsicdPlugin.GetIfIndexFromIntfIdAndIntfType(vlanId, asicdClntDefs.IfTypeVlan))
			_, exist := server.l3IntfPropMap[vlanIfIdx]
			if exist {
				if vlanId == lagEnt.UntagVlanId {
					for _, portIfIdx := range createPortList {
						server.createPortFromL3Lag(vlanIfIdx, vlanId, lagIfIdx, int(portIfIdx), UNTAGGED)
					}
					for _, portIfIdx := range delPortList {
						server.deletePortFromL3Lag(vlanIfIdx, vlanId, lagIfIdx, int(portIfIdx), UNTAGGED)
					}
				} else {
					for _, portIfIdx := range createPortList {
						server.createPortFromL3Lag(vlanIfIdx, vlanId, lagIfIdx, int(portIfIdx), TAGGED)
					}
					for _, portIfIdx := range delPortList {
						server.deletePortFromL3Lag(vlanIfIdx, vlanId, lagIfIdx, int(portIfIdx), TAGGED)
					}

				}
			}
		}
	}
}

func (server *ARPServer) createPortFromL3Lag(l3IfIdx, vlanId, lagIfIdx, portIfIdx int, TagFlag bool) {
	portEnt, _ := server.portPropMap[portIfIdx]
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	if TagFlag == UNTAGGED {
		portEnt.UntagVlanId = vlanId
	}
	l3PortPropEnt, _ := portEnt.L3PortPropMap[vlanId]
	l3PortPropEnt.IpAddr = l3Ent.IpAddr
	l3PortPropEnt.L3IfIdx = l3IfIdx
	l3PortPropEnt.LagIfIdx = lagIfIdx
	l3PortPropEnt.Netmask = l3Ent.Netmask
	portEnt.L3PortPropMap[vlanId] = l3PortPropEnt
	server.portPropMap[portIfIdx] = portEnt
}

func (server *ARPServer) deletePortFromL3Lag(l3IfIdx, vlanId, lagIfIdx, portIfIdx int, TagFlag bool) {
	portEnt, _ := server.portPropMap[portIfIdx]
	l3PortPropEnt, exist := portEnt.L3PortPropMap[vlanId]
	if exist {
		l3PortPropEnt.LagIfIdx = -1
		portEnt.L3PortPropMap[vlanId] = l3PortPropEnt
		server.portPropMap[portIfIdx] = portEnt
		server.arpEntryDeleteCh <- DeleteArpEntryMsg{
			Type:  DeleteBasedOnPort,
			IfIdx: portIfIdx,
		}
	}
}

func (server *ARPServer) restorePortFilter(ifIndex int) error {
	port := server.portPropMap[ifIndex]
	port.baseFilter = server.constructBaseFilter(port.MacAddr)
	if port.PcapHdl != nil {
		filter := port.baseFilter + CLOSE_FILTER
		err := port.PcapHdl.SetBPFFilter(filter)
		if err != nil {
			server.logger.Err("Failed to Restore Pcap Filter:", filter, "for port:", port.IfName)
		}
	}
	server.portPropMap[ifIndex] = port
	return nil
}

func (server *ARPServer) restoreVlanMembersFilters(vlanIfIndex int) error {
	vlan := server.vlanPropMap[vlanIfIndex]
	for ifIndex, _ := range vlan.UntagIfIdxMap {
		server.restorePortFilter(ifIndex)
	}
	for ifIndex, _ := range vlan.TagIfIdxMap {
		server.restorePortFilter(ifIndex)
	}
	server.vlanPropMap[vlanIfIndex] = vlan
	return nil
}

func (server *ARPServer) RestoreBPFFilter(msg asicdClntDefs.IPv4VirtualIntfStateNotifyMsg) error {
	virEntry, exists := server.virtualIntfPropMap[msg.IfIndex]
	if !exists {
		server.logger.Err("No entry found for virtual interface during state down and hence filter cannot be restored:", msg.IfIndex, msg.IpAddr)
		return nil
	}
	_, exists = server.vlanPropMap[int(virEntry.ParentIfIndex)]
	if exists {
		return server.restoreVlanMembersFilters(int(virEntry.ParentIfIndex))
	}
	_, exists = server.portPropMap[int(virEntry.ParentIfIndex)]
	if exists {
		return server.restorePortFilter(int(virEntry.ParentIfIndex))
	}
	return nil
}

func (server *ARPServer) updatePortFilter(ifIndex int, macAddr string) error {
	port := server.portPropMap[ifIndex]
	if port.PcapHdl != nil {
		filter := fmt.Sprintf("%s%s%s%s", port.baseFilter, OR_ETHER_SRC, macAddr, CLOSE_FILTER)
		err := port.PcapHdl.SetBPFFilter(filter)
		if err != nil {
			server.logger.Err("Failed to Update Pcap Filter:", filter, "for port:", port.IfName)
		} else {
			server.logger.Debug("Port:", port.IfName, "filter updated to:", filter, "successfully")
		}
	}
	server.portPropMap[ifIndex] = port
	return nil
}

func (server *ARPServer) updateVlanMembersFilters(vlanIfIndex int, macAddr string) error {
	vlan := server.vlanPropMap[vlanIfIndex]
	for ifIndex, _ := range vlan.UntagIfIdxMap {
		server.updatePortFilter(ifIndex, macAddr)
	}
	for ifIndex, _ := range vlan.TagIfIdxMap {
		server.updatePortFilter(ifIndex, macAddr)
	}
	server.vlanPropMap[vlanIfIndex] = vlan
	return nil
}

func (server *ARPServer) updateFilter(ifIndex int, macAddr string) error {
	_, exists := server.vlanPropMap[ifIndex]
	if exists {
		return server.updateVlanMembersFilters(ifIndex, macAddr)
	}
	_, exists = server.portPropMap[ifIndex]
	if exists {
		return server.updatePortFilter(ifIndex, macAddr)
	}
	return nil
}

func (server *ARPServer) UpdateBPFFilter(msg asicdClntDefs.IPv4VirtualIntfStateNotifyMsg) error {
	virEntry, exists := server.virtualIntfPropMap[msg.IfIndex]
	if !exists {
		server.logger.Err("No entry found for virtual interface during state down and hence filter cannot be restored:", msg.IfIndex, msg.IpAddr)
		return nil
	}

	return server.updateFilter(int(virEntry.ParentIfIndex), virEntry.MacAddr)
}

func (server *ARPServer) constructVirtualIpInfra() {
	ipsInfo, err := server.AsicdPlugin.GetAllSubIPv4IntfState()
	if err != nil {
		server.logger.Warning("Failed to get all sub ipv4 interfaces from system, ERROR:", err)
		return
	}
	for _, ipInfo := range ipsInfo {
		if ipInfo.Type == commonDefs.SUB_INTF_SECONDARY_TYPE {
			continue
		}
		virEntry := server.virtualIntfPropMap[ipInfo.IfIndex]
		virEntry.IfIndex = ipInfo.IfIndex
		virEntry.ParentIfIndex = ipInfo.ParentIfIndex
		virEntry.IpAddr = ipInfo.IpAddr
		virEntry.MacAddr = ipInfo.MacAddr
		virEntry.IfName = ipInfo.IfName
		server.virtualIntfPropMap[ipInfo.IfIndex] = virEntry
		if ipInfo.OperState == "UP" {
			server.updateFilter(int(virEntry.ParentIfIndex), virEntry.MacAddr)
		}
	}
}
