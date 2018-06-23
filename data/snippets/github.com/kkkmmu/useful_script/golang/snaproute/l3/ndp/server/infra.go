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
	"errors"
	"fmt"
	"l3/ndp/config"
	"l3/ndp/debug"
	"net"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/commonDefs"
)

/*
 * API: will return all system port information
 */
func (svr *NDPServer) GetPorts() {
	debug.Logger.Info("Get Port State List")
	curMark := svr.SwitchPlugin.GetMinSysPort()
	count := 512
	for {
		bulkInfo, _ := svr.SwitchPlugin.GetBulkPortState(curMark, count)
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			obj := bulkInfo.PortStateList[i]
			var empty struct{}
			port := config.PortInfo{
				IntfRef:   obj.IntfRef,
				IfIndex:   obj.IfIndex,
				OperState: obj.OperState,
				Name:      obj.Name,
			}
			pObj, err := svr.SwitchPlugin.GetPort(obj.Name)
			if err != nil {
				debug.Logger.Err("Getting mac address for", obj.Name, "failed, error:", err)
			} else {
				port.MacAddr = pObj.MacAddr
				port.Description = pObj.Description
			}
			l2Port := svr.L2Port[port.IfIndex]
			l2Port.Info = port
			l2Port.RX = nil
			debug.Logger.Info("L2 IfIndex:", port.IfIndex, "Information is:", l2Port.Info)
			svr.L2Port[port.IfIndex] = l2Port
			svr.SwitchMacMapEntries[port.MacAddr] = empty
			svr.SwitchMac = port.MacAddr // @HACK.... need better solution
		}
		if more == false {
			break
		}
	}
	debug.Logger.Info("Done with Port State list")
	return
}

/*
 * API: will return all system vlan information
 */
func (svr *NDPServer) GetVlans() {
	debug.Logger.Info("Get Vlan Information")
	curMark := svr.SwitchPlugin.GetMinSysPort()
	count := 1024
	vlanConfigList := make([]*asicdClntDefs.Vlan, 0)
	// get all vlan configuration
	for {
		bulkVlanInfo, _ := svr.SwitchPlugin.GetBulkVlan(curMark, count)
		if bulkVlanInfo == nil {
			break
		}

		//objCnt := int(bulkVlanInfo.Count)
		more := bool(bulkVlanInfo.More)
		curMark = int(bulkVlanInfo.EndIdx)
		vlanConfigList = append(vlanConfigList, bulkVlanInfo.VlanList...)
		if more == false {
			break
		}
	}
	curMark = 0
	// get all vlan state information
	for {
		bulkVlanStateInfo, _ := svr.SwitchPlugin.GetBulkVlanState(curMark, count)
		if bulkVlanStateInfo == nil {
			break
		}
		//objCnt := int(bulkVlanStateInfo.Count)
		more := bool(bulkVlanStateInfo.More)
		curMark = int(bulkVlanStateInfo.EndIdx)
		// store vlan state information like name, ifIndex, operstate
		for _, vlanState := range bulkVlanStateInfo.VlanStateList {
			debug.Logger.Info("vlan:", *vlanState)
			entry, _ := svr.VlanInfo[vlanState.IfIndex]
			entry.VlanId = vlanState.VlanId
			entry.VlanIfIndex = vlanState.IfIndex
			entry.Name = vlanState.VlanName
			entry.OperState = vlanState.OperState
			for _, vlanconfig := range vlanConfigList {
				if entry.VlanId == vlanconfig.VlanId {
					entry.UntagPortsMap = make(map[int32]bool)
					for _, untagintf := range vlanconfig.UntagIfIndexList {
						entry.UntagPortsMap[untagintf] = true
					}
					entry.TagPortsMap = make(map[int32]bool)
					for _, tagIntf := range vlanconfig.IfIndexList {
						entry.TagPortsMap[tagIntf] = true
					}
				}
			}
			svr.VlanInfo[vlanState.IfIndex] = entry
			svr.VlanIfIdxVlanIdMap[vlanState.VlanName] = vlanState.VlanId
			svr.Dot1QToVlanIfIndex[vlanState.VlanId] = vlanState.IfIndex
		}
		if more == false {
			break
		}
	}
	debug.Logger.Info("Done with Vlan List")
	return
}

/*
 * API: will return all system L3 interfaces information
 */
func (svr *NDPServer) GetIPIntf() {
	debug.Logger.Info("Get IPv6 Interface List")
	curMark := svr.SwitchPlugin.GetMinSysPort()
	count := 1024
	for {
		bulkIPv6Info, _ := svr.SwitchPlugin.GetBulkIPv6IntfState(curMark, count)
		if bulkIPv6Info == nil {
			break
		}
		//objCnt := int(bulkIPv6Info.Count)
		more := bool(bulkIPv6Info.More)
		curMark = int(bulkIPv6Info.EndIdx)
		for _, obj := range bulkIPv6Info.IPv6IntfStateList {
			// ndp will not listen on loopback interfaces
			if obj.L2IntfType == commonDefs.LOOPBACK {
				continue
			}
			ipInfo, exists := svr.L3Port[obj.IfIndex]
			if !exists {
				ipStObj := &config.IPIntfNotification{
					IntfRef:   obj.IntfRef,
					IfIndex:   obj.IfIndex,
					Operation: obj.OperState,
					IpAddr:    obj.IpAddr,
				}
				ipInfo.InitIntf(ipStObj, svr.PktDataCh, svr.NdpConfig)
				ipInfo.SetIfType(svr.GetIfType(obj.IfIndex))
				// cache reverse map from intfref to ifIndex, used mainly during state
				svr.L3IfIntfRefToIfIndex[obj.IntfRef] = obj.IfIndex
			} else {
				ipInfo.UpdateIntf(obj.IpAddr)
			}
			svr.L3Port[ipInfo.IfIndex] = ipInfo
			if !exists {
				svr.ndpL3IntfStateSlice = append(svr.ndpL3IntfStateSlice, ipInfo.IfIndex)
			}
		}
		if more == false {
			break
		}
	}
	debug.Logger.Info("Done with IPv6 State list")
	return
}

func (svr *NDPServer) getVirtualIpIntf() {
	debug.Logger.Info("Get Virtual IPv6 Interface List")
	curMark := svr.SwitchPlugin.GetMinSysPort()
	count := 1024
	for {
		bulkIPv6Info, _ := svr.SwitchPlugin.GetBulkSubIPv6IntfState(curMark, count)
		if bulkIPv6Info == nil {
			break
		}
		//objCnt := int(bulkIPv6Info.Count)
		more := bool(bulkIPv6Info.More)
		curMark = int(bulkIPv6Info.EndIdx)
		for _, ipInfo := range bulkIPv6Info.SubIPv6IntfStateList {
			if ipInfo.Type == commonDefs.SUB_INTF_VIRTUAL_TYPE {
				virEntry := &config.VirtualIpInfo{}
				virEntry.IfIndex = ipInfo.IfIndex
				virEntry.ParentIfIndex = ipInfo.ParentIfIndex
				virEntry.IpAddr = ipInfo.IpAddr
				virEntry.MacAddr = ipInfo.MacAddr
				virEntry.IfName = ipInfo.IfName
				svr.virtualIp[ipInfo.IfIndex] = virEntry
				if ipInfo.OperState == config.STATE_UP {
					svr.updateFilters(virEntry)
				}
			}
		}
		if more == false {
			break
		}
	}
	debug.Logger.Info("Done with virtual ipv6 interface list")
}

func (svr *NDPServer) GetIfType(ifIndex int32) int {
	debug.Logger.Info("get ifType for ifIndex:", ifIndex)
	if _, ok := svr.L2Port[ifIndex]; ok {
		debug.Logger.Info("L3 Port is of IfTypePort")
		return commonDefs.IfTypePort
	}

	if _, ok := svr.VlanInfo[ifIndex]; ok {
		debug.Logger.Info("L3 Port is of IfTypeVlan")
		return commonDefs.IfTypeVlan
	}
	debug.Logger.Info("no valid ifIndex found")
	return -1
}

/*
 *	insertNeighborInfo: Helper API to update list of neighbor keys that are created by ndp
 *	@NOTE: caller is responsible for acquiring the lock to access slice
 */
func (svr *NDPServer) insertNeigborInfo(nbrInfo *config.NeighborConfig, hwIfIndex int32, learnedIntf string) {
	nbrEntry := *nbrInfo
	nbrEntry.IfIndex = hwIfIndex
	nbrEntry.Intf = learnedIntf
	nbrKey := createNeighborKey(nbrInfo.MacAddr, nbrInfo.IpAddr, learnedIntf)
	debug.Logger.Debug("server state nbrKey is:", nbrKey)
	svr.NeighborInfo[nbrKey] = nbrEntry
	svr.neighborKey = append(svr.neighborKey, nbrKey)
}

/*
 *	deleteSvrStateNbrInfo: Helper API to update list of neighbor keys that are deleted by ndp
 *	@NOTE: caller is responsible for acquiring the lock to access slice
 */
func (svr *NDPServer) deleteSvrStateNbrInfo(nbrKey string) {
	// delete the entry from neighbor map
	delete(svr.NeighborInfo, nbrKey)

	for idx, _ := range svr.neighborKey {
		if svr.neighborKey[idx] == nbrKey {
			svr.neighborKey = append(svr.neighborKey[:idx], svr.neighborKey[idx+1:]...)
			break
		}
	}
}

/*
 *	updateStateNbrInfo: Helper API to update list of neighbor keys that are being updated by Mac Move
 *	@NOTE: caller is responsible for acquiring the lock to access slice
 */
func (svr *NDPServer) updateStateNeighborInfo(nbrKey, ipAddr string, ifIndex, vlanId int32) string {
	nbrEntry, exists := svr.NeighborInfo[nbrKey]
	if !exists {
		return ""
	}
	debug.Logger.Info("found nbrEntry is:", nbrEntry)
	l2Port, exists := svr.L2Port[ifIndex]
	if exists {
		debug.Logger.Info("Updating entry:", nbrEntry, "old key is:", nbrKey)
		newNbrEntry := nbrEntry
		//svr.insertNeigborInfo(&newNbrEntry, msg.IfIndex, l2Port.Info.Name)
		newNbrEntry.IfIndex = ifIndex
		newNbrEntry.Intf = l2Port.Info.Name
		newNbrKey := createNeighborKey(newNbrEntry.MacAddr, newNbrEntry.IpAddr, newNbrEntry.Intf)
		debug.Logger.Info("Updated entry:", newNbrEntry, "new key is:", newNbrKey)
		svr.NeighborInfo[newNbrKey] = newNbrEntry
		svr.neighborKey = append(svr.neighborKey, newNbrKey)
		debug.Logger.Info("Deleting old nbr entry at old location:", nbrKey)
		svr.deleteSvrStateNbrInfo(nbrKey)
		return newNbrKey
	}

	return ""
}

/*
 *	 CreateNeighborInfo
 *			a) It will first check whether a neighbor exists in the neighbor cache
 *			b) If it doesn't exists then we create neighbor in the platform
 *		        a) It will update ndp server neighbor info cache with the latest information
 */
func (svr *NDPServer) CreateNeighborInfo(nbrInfo *config.NeighborConfig, hwIfIndex int32, learnedIntf string) {
	debug.Logger.Debug("Calling create ipv6 neighbor for global nbrinfo is", nbrInfo.IpAddr, nbrInfo.MacAddr,
		"vlanId:", nbrInfo.VlanId, "ifIndex:", hwIfIndex, "interface:", learnedIntf)
	if net.ParseIP(nbrInfo.IpAddr).IsLinkLocalUnicast() == false {
		_, err := svr.SwitchPlugin.CreateIPv6Neighbor(nbrInfo.IpAddr, nbrInfo.MacAddr, nbrInfo.VlanId, hwIfIndex)
		if err != nil {
			debug.Logger.Err("create ipv6 global neigbor failed for", nbrInfo, "error is", err)
			// do not enter that neighbor in our neigbor map
			return
		}
	}
	// for bgp send out l3 ifIndex only do not use hwIfIndex
	svr.SendIPv6CreateNotification(nbrInfo.IpAddr, nbrInfo.IfIndex)
	// after sending notification update ifIndex to hwIfIndex
	svr.NeigborEntryLock.Lock()
	svr.insertNeigborInfo(nbrInfo, hwIfIndex, learnedIntf)
	svr.NeigborEntryLock.Unlock()
}

func (svr *NDPServer) deleteNeighbor(nbrKey string, ifIndex int32) {
	debug.Logger.Debug("deleteNeighbor called for nbrKey:", nbrKey)
	// Inform clients that neighbor is gonna be deleted
	splitString := splitNeighborKey(nbrKey)
	nbrIp := splitString[1]
	svr.SendIPv6DeleteNotification(nbrIp, ifIndex)
	// Request asicd to delete the neighbor
	if net.ParseIP(nbrIp).IsLinkLocalUnicast() == false {
		_, err := svr.SwitchPlugin.DeleteIPv6Neighbor(nbrIp, "00:00:00:00:00:00", 0, 0)
		if err != nil {
			debug.Logger.Err("delete ipv6 neigbor failed for", nbrIp, "error is", err)
		}
	}
	svr.NeigborEntryLock.Lock()
	svr.deleteSvrStateNbrInfo(nbrKey)
	svr.NeigborEntryLock.Unlock()
}

/*
 *	 DeleteNeighborInfo
 *			a) It will first check whether a neighbor exists in the neighbor cache
 *			b) If it doesn't exists then we will move on to next neighbor
 *		        c) If exists then we will call DeleteIPV6Neighbor for that entry and remove
 *			   the entry from our runtime information
 *	ifIndex is always l3 ifIndex
 */
func (svr *NDPServer) DeleteNeighborInfo(deleteEntries []string, ifIndex int32) {
	for _, nbrKey := range deleteEntries {
		debug.Logger.Debug("Calling delete ipv6 neighbor for nbr:", nbrKey, "ifIndex:", ifIndex)
		svr.deleteNeighbor(nbrKey, ifIndex)
	}
}

/*  API: will handle IPv6 notifications received from switch/asicd
 *      Msg types
 *	    1) Create:
 *		    Create an entry in the map
 *	    2) Delete:
 *		    delete an entry from the map
 */
func (svr *NDPServer) HandleIPIntfCreateDelete(obj *config.IPIntfNotification) {
	ipInfo, exists := svr.L3Port[obj.IfIndex]
	switch obj.Operation {
	case config.CONFIG_CREATE:
		// Done during Init
		if exists {
			ipInfo.UpdateIntf(obj.IpAddr)
			svr.L3Port[obj.IfIndex] = ipInfo
			return
		}

		ipInfo = Interface{}
		ipInfo.CreateIntf(obj, svr.PktDataCh, svr.NdpConfig)
		ipInfo.SetIfType(svr.GetIfType(obj.IfIndex))
		// cache reverse map from intfref to ifIndex, used mainly during state
		svr.L3IfIntfRefToIfIndex[obj.IntfRef] = obj.IfIndex
		svr.ndpL3IntfStateSlice = append(svr.ndpL3IntfStateSlice, ipInfo.IfIndex)
	case config.CONFIG_DELETE:
		if !exists {
			debug.Logger.Err("Got Delete request for non existing l3 port", obj.IfIndex)
			return
		}
		// stop rx/tx on the deleted interface
		debug.Logger.Info("Delete IP interface received for", ipInfo.IntfRef, "ifIndex:", ipInfo.IfIndex)
		deleteEntries := ipInfo.DeInitIntf()
		if len(deleteEntries) > 0 {
			svr.DeleteNeighborInfo(deleteEntries, obj.IfIndex)
		}
		delete(svr.L3IfIntfRefToIfIndex, obj.IntfRef)
		// @TODO: need to take care for ifTYpe vlan
		//@TODO: need to remove ndp l3 interface from up slice, but that is taken care of by stop rx/tx
	}
	svr.L3Port[ipInfo.IfIndex] = ipInfo
}

/*  API: will handle l2/physical notifications received from switch/asicd
 *	  Update map entry and then call state notification
 *
 */
func (svr *NDPServer) HandlePhyPortStateNotification(msg *config.PortState) {
	debug.Logger.Info("Handling L2 Port State:", msg.IfState, "for ifIndex:", msg.IfIndex)
	svr.updateL2Operstate(msg.IfIndex, msg.IfState)
}

/*  API: will handle Vlan Create/Delete/Update notifications received from switch/asicd
 */
func (svr *NDPServer) HandleVlanNotification(msg *config.VlanNotification) {
	debug.Logger.Info("Handle Vlan Notfication:", msg.Operation, "for vlanId:", msg.VlanId, "vlan:", msg.VlanName,
		"vlanIfIndex:", msg.VlanIfIndex, "tagList:", msg.TagPorts, "unTagList:", msg.UntagPorts)
	vlan, exists := svr.VlanInfo[msg.VlanIfIndex]
	switch msg.Operation {
	case config.CONFIG_CREATE:
		debug.Logger.Info("Received Vlan Create:", *msg)
		svr.VlanIfIdxVlanIdMap[msg.VlanName] = msg.VlanId
		vlan.Name = msg.VlanName
		vlan.VlanId = msg.VlanId
		vlan.VlanIfIndex = msg.VlanIfIndex
		// Store untag port information
		for _, untagIntf := range msg.UntagPorts {
			if vlan.UntagPortsMap == nil {
				vlan.UntagPortsMap = make(map[int32]bool)
			}
			vlan.UntagPortsMap[untagIntf] = true
		}
		// Store untag port information
		for _, tagIntf := range msg.TagPorts {
			if vlan.TagPortsMap == nil {
				vlan.TagPortsMap = make(map[int32]bool)
			}
			vlan.TagPortsMap[tagIntf] = true
		}
		svr.VlanInfo[msg.VlanIfIndex] = vlan
		svr.Dot1QToVlanIfIndex[msg.VlanId] = msg.VlanIfIndex
	case config.CONFIG_DELETE:
		debug.Logger.Info("Received Vlan Delete:", *msg)
		if exists {
			vlan.UntagPortsMap = nil
			vlan.TagPortsMap = nil
			delete(svr.VlanInfo, msg.VlanIfIndex)
		}
		delete(svr.Dot1QToVlanIfIndex, msg.VlanId)
	case config.CONFIG_UPDATE:
		debug.Logger.Info("Received Vlan Update:", *msg)
		for tagIntf, _ := range vlan.TagPortsMap {
			del := true
			for _, msgTagIntf := range msg.TagPorts {
				if msgTagIntf == tagIntf {
					del = false
					break
				}
			}
			if del {
				delete(vlan.TagPortsMap, tagIntf)
			}
		}
		for unTagIntf, _ := range vlan.UntagPortsMap {
			del := true
			for _, msgUnTagIntf := range msg.UntagPorts {
				if msgUnTagIntf == unTagIntf {
					del = false
					break
				}
			}
			if del {
				delete(vlan.UntagPortsMap, unTagIntf)
			}
		}
		// Store untag port information
		for _, untagIntf := range msg.UntagPorts {
			if vlan.UntagPortsMap == nil {
				vlan.UntagPortsMap = make(map[int32]bool)
			}
			vlan.UntagPortsMap[untagIntf] = true
		}
		// Store untag port information
		for _, tagIntf := range msg.TagPorts {
			if vlan.TagPortsMap == nil {
				vlan.TagPortsMap = make(map[int32]bool)
			}
			vlan.TagPortsMap[tagIntf] = true
		}
		svr.VlanInfo[msg.VlanIfIndex] = vlan
		svr.UpdatePhyPortToVlanInfo(msg)
	}
}

/*  API: will handle IPv6 notifications received from switch/asicd
 *      Msg types
 *	    1) Create:
 *		     Start Rx/Tx in this case
 *	    2) Delete:
 *		     Stop Rx/Tx in this case
 */
func (svr *NDPServer) HandleStateNotification(msg *config.IPIntfNotification) {
	debug.Logger.Info("Handling L3 State:", msg.Operation, "for port:", msg.IntfRef, "ifIndex:", msg.IfIndex, "ipAddr:", msg.IpAddr)
	switch msg.Operation {
	case config.STATE_UP:
		debug.Logger.Info("Create pkt handler for port:", msg.IntfRef, "ifIndex:", msg.IfIndex, "IpAddr:", msg.IpAddr)
		svr.StartRxTx(msg.IfIndex)
	case config.STATE_DOWN:
		debug.Logger.Info("Delete pkt handler for port:", msg.IntfRef, "ifIndex:", msg.IfIndex, "IpAddr:", msg.IpAddr)
		svr.StopRxTx(msg.IfIndex, msg.IpAddr)
	}
}

/*
 *    API: helper function to update ifIndex & port information for software. Hardware is already taken care
 *	   off
 *	   NOTE:
 *         Below Scenario will only happen when mac move happens between a physical port.. L3 port remains
 *	   the same and hence do not need to update clients
 */
func (svr *NDPServer) SoftwareUpdateNbrEntry(msg *config.MacMoveNotification) {
	debug.Logger.Info("Received Mac Move Notification for IPV6 entry:", *msg)
	nbrIp := msg.IpAddr
	svr.NeigborEntryLock.Lock()
	defer svr.NeigborEntryLock.Unlock()
	for _, nbrKey := range svr.neighborKey {
		splitString := splitNeighborKey(nbrKey)
		if splitString[1] == nbrIp {
			debug.Logger.Info("Updating Neigbor information for:", nbrIp)
			debug.Logger.Info("Neighbor Key is:", nbrKey)
			// updating svr neighbor state information
			newNbrKey := svr.updateStateNeighborInfo(nbrKey, msg.IpAddr, msg.IfIndex, msg.VlanId)
			if newNbrKey == "" {
				return
			}
			// updating interface neighbor state information
			l3IfIndex, exists := svr.Dot1QToVlanIfIndex[msg.VlanId]
			if exists {
				l3Port, exists := svr.L3Port[l3IfIndex]
				if exists {
					debug.Logger.Info("Updating Neighbor Cache in ipv6 adj")
					oldLsKey, linkScopeIp := l3Port.UpdateNbrEntry(nbrKey, newNbrKey)
					svr.L3Port[l3IfIndex] = l3Port
					if oldLsKey != "" && linkScopeIp != "" {
						debug.Logger.Info("Updating Link Scope Neighbor Information:", linkScopeIp)
						svr.updateStateNeighborInfo(oldLsKey, linkScopeIp, msg.IfIndex, msg.VlanId)
					}
				}
			}
			break
		}
	}
}

/*
 *    API: handle action request coming from the user
 */
func (svr *NDPServer) HandleAction(action *config.ActionData) {
	debug.Logger.Debug("Handle Action:", *action)

	switch action.Type {
	case config.DELETE_BY_IFNAME:
		svr.ActionDeleteByIntf(action.IntfRef)

	case config.DELETE_BY_IPADDR:
		svr.ActionDeleteByNbrIp(action.NbrIp)

	case config.REFRESH_BY_IFNAME:
		svr.ActionRefreshByIntf(action.IntfRef)

	case config.REFRESH_BY_IPADDR:
		svr.ActionRefreshByNbrIp(action.NbrIp)
	}
}

/*
 *    API: It will remove any deleted ip port from the up state slice list
 */
func (svr *NDPServer) DeleteL3IntfFromUpState(ifIndex int32) {
	for idx, entry := range svr.ndpUpL3IntfStateSlice {
		if entry == ifIndex {
			//@TODO: need to optimize this
			svr.ndpUpL3IntfStateSlice = append(svr.ndpUpL3IntfStateSlice[:idx],
				svr.ndpUpL3IntfStateSlice[idx+1:]...)
			break
		}
	}
}

/*
 *    API: It will populate correct vlan information which will be used for ipv6 neighbor create
 */
func (svr *NDPServer) PopulateVlanInfo(nbrInfo *config.NeighborConfig, intfRef string) {
	// check if the ifIndex is present in the reverse map..
	vlanId, exists := svr.VlanIfIdxVlanIdMap[intfRef]
	if exists {
		// if the entry exists then use the vlanId from reverse map
		nbrInfo.VlanId = vlanId
	} else {
		// @TODO: move this to plugin specific
		// in this case use system reserved Vlan id which is -1
		nbrInfo.VlanId = config.INTERNAL_VLAN
	}
}

/*
 *    API: send ipv6 neighbor create notification
 */
func (svr *NDPServer) SendIPv6CreateNotification(ipAddr string, ifIndex int32) {
	msgBuf, err := createNotificationMsg(ipAddr, ifIndex)
	if err != nil {
		return
	}

	notification := commonDefs.NdpNotification{
		MsgType: commonDefs.NOTIFY_IPV6_NEIGHBOR_CREATE,
		Msg:     msgBuf,
	}
	debug.Logger.Info("Sending Create notification for ip address:", ipAddr, "and ifIndex:", ifIndex, "to other protocols")
	svr.pushNotification(notification)
}

/*
 *    API: send ipv6 neighbor delete notification
 */
func (svr *NDPServer) SendIPv6DeleteNotification(ipAddr string, ifIndex int32) {
	msgBuf, err := createNotificationMsg(ipAddr, ifIndex)
	if err != nil {
		return
	}

	notification := commonDefs.NdpNotification{
		MsgType: commonDefs.NOTIFY_IPV6_NEIGHBOR_DELETE,
		Msg:     msgBuf,
	}
	debug.Logger.Info("Sending Delete notification for ip address:", ipAddr, "and ifIndex:", ifIndex, "to other protocols")
	svr.pushNotification(notification)
}

/*
 * helper function to create notification msg
 */
func createNotificationMsg(ipAddr string, ifIndex int32) ([]byte, error) {
	msg := commonDefs.Ipv6NeighborNotification{
		IpAddr:  ipAddr,
		IfIndex: ifIndex,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		debug.Logger.Err("Failed to marshal IPv6 Neighbor Notification message", msg, "error:", err)
		return msgBuf, err
	}

	return msgBuf, nil
}

/*
 * helper function to marshal notification and push it on to the channel
 */
func (svr *NDPServer) pushNotification(notification commonDefs.NdpNotification) {
	notifyBuf, err := json.Marshal(notification)
	if err != nil {
		debug.Logger.Err("Failed to marshal ipv6 notification before pushing it on channel error:", err)
		return
	}
	svr.notifyChan <- notifyBuf
}

/*
 *  Change L2 port state from switch asicd notification
 */
func (svr *NDPServer) updateL2Operstate(ifIndex int32, state string) {
	l2Port, exists := svr.L2Port[ifIndex]
	if !exists {
		debug.Logger.Err("No L2 Port found for ifIndex:", ifIndex, "hence nothing to update on OperState")
		return
	}
	l2Port.Info.OperState = state
	/* HANDLE PORT FLAP SCENARIOS */
	switch state {
	case config.STATE_UP:
		// NO-OP Just change the state
	case config.STATE_DOWN:
		debug.Logger.Info("L2 Port is down and hence deleting pcap handler for port:", l2Port.Info.Name)
		l2Port.deletePcap()
	}
	debug.Logger.Info("During L2 State Notification L2 IfIndex:", ifIndex, "Information is:", l2Port.Info)
	svr.L2Port[ifIndex] = l2Port
}

/*
 *  Creating Pcap handlers for l2 port which are marked as tag/untag for l3 vlan port and are in UP state
 *  only l3 CreatePcap should update l2Port.L3 information
 */
func (svr *NDPServer) CreatePcap(ifIndex int32) error {
	debug.Logger.Info("Creating Physical Port Pcap to L3 Vlan, ifIndex:", ifIndex)
	vlan, exists := svr.VlanInfo[ifIndex]
	if !exists {
		debug.Logger.Err("No matching vlan found for ifIndex:", ifIndex)
		return errors.New(fmt.Sprintln("No matching vlan found for ifIndex:", ifIndex))
	}
	debug.Logger.Debug("Creating Pcap Handlers for tag ports:", vlan.TagPortsMap)
	// open rx pcap handler for tagged ports
	for pIfIndex, _ := range vlan.TagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			debug.Logger.Debug("L2:", l2Port.Info.Name, "operstate is", l2Port.Info.OperState)
			if l2Port.Info.OperState == config.STATE_UP {
				l2Port.createPortPcap(svr.RxPktCh, l2Port.Info.Name)
				svr.L2Port[pIfIndex] = l2Port
			}
		}
	}
	debug.Logger.Debug("Creating pcap handlers for unTag ports:", vlan.UntagPortsMap)
	// open rx pcap handler for untagged ports
	for pIfIndex, _ := range vlan.UntagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			debug.Logger.Debug("L2:", l2Port.Info.Name, "operstate is", l2Port.Info.OperState)
			if l2Port.Info.OperState == config.STATE_UP {
				l2Port.createPortPcap(svr.RxPktCh, l2Port.Info.Name)
				svr.L2Port[pIfIndex] = l2Port
			}
			// reverse map updated
			l3Info := L3Info{
				Name:    vlan.Name,
				IfIndex: ifIndex,
			}
			svr.L2Port[pIfIndex] = l2Port
			svr.PhyPortToL3PortMap[pIfIndex] = l3Info
		}
	}
	debug.Logger.Debug("PhyPortToL3PortMap after create vlan is:", svr.PhyPortToL3PortMap)
	return nil
}

/*
 *  Deleting Pcap handlers for l2 port which are marked as tag/untag for l3 vlan port and are in UP state
 *  only l3 CreatePcap should update l2Port.L3 information
 */
func (svr *NDPServer) DeletePcap(ifIndex int32) {
	debug.Logger.Info("Deleting Physical Port Pcap RX Handlers for L3 Vlan, ifIndex:", ifIndex)
	vlan, exists := svr.VlanInfo[ifIndex]
	if !exists {
		debug.Logger.Err("No matching vlan found for ifIndex:", ifIndex)
		return //errors.New(fmt.Sprintln("No matching vlan found for ifIndex:", ifIndex))
	}
	debug.Logger.Debug("Deleting pcap handlers for Tag ports:", vlan.TagPortsMap)
	// open rx pcap handler for tagged ports
	for pIfIndex, _ := range vlan.TagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.deletePcap()
			svr.L2Port[pIfIndex] = l2Port
		}
	}
	// open rx pcap handler for untagged ports
	debug.Logger.Debug("Deleting pcap handlers for unTag ports:", vlan.UntagPortsMap)
	for pIfIndex, _ := range vlan.UntagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.deletePcap()
			svr.L2Port[pIfIndex] = l2Port
			delete(svr.PhyPortToL3PortMap, pIfIndex)
		}
	}
}

func (svr *NDPServer) UpdatePhyPortToVlanInfo(msg *config.VlanNotification) {
	ifIndex := msg.VlanIfIndex
	debug.Logger.Info("Updating Phy port to vlan information, for ifIndex:", msg.VlanIfIndex)
	vlan, exists := svr.VlanInfo[msg.VlanIfIndex]
	if !exists {
		debug.Logger.Err("no matching vlan found for update msg:", *msg)
		return
	}
	debug.Logger.Info("vlan tag port information is:", msg.TagPorts)
	// iterating over slice
	for _, pIfIndex := range msg.UntagPorts {
		debug.Logger.Info("Untag port ifIndex:", pIfIndex)
		l3Info, exists := svr.PhyPortToL3PortMap[pIfIndex]
		if !exists {
			l3Info = L3Info{
				Name:    vlan.Name,
				IfIndex: ifIndex,
				Updated: true,
			}
			debug.Logger.Info("new untag port received for pIfIndex:", pIfIndex, "L3Info:", l3Info)
		} else {
			debug.Logger.Info("existing untag port setting update to true for pIfIndex:", pIfIndex, "L3Info:", l3Info)
			l3Info.Updated = true
		}
		svr.PhyPortToL3PortMap[pIfIndex] = l3Info
	}
	l3Port, l3exists := svr.L3Port[ifIndex]
	// iterating over map
	for pIfIndex, l3Info := range svr.PhyPortToL3PortMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if l3Info.Updated == false {
			debug.Logger.Debug("before deleteing pcap for:", pIfIndex, "check if it belongs to tag port map")
			if _, ok := vlan.TagPortsMap[pIfIndex]; !ok {
				debug.Logger.Info("pIfIndex:", pIfIndex, "is removed from vlan:", ifIndex, vlan.Name, "as un-tag member",
					"hence deleting its pcap")
				l2Port.deletePcap()
			}
			delete(svr.PhyPortToL3PortMap, pIfIndex)
		} else {
			debug.Logger.Info("pIfIndex:", pIfIndex, "is part of vlan:", ifIndex, vlan.Name, "as un-tag member",
				"check if l3 is up or not and start rx/tx if needed")
			l3Info.Updated = false
			svr.PhyPortToL3PortMap[pIfIndex] = l3Info
			// if the entry is updated then check for l3 entry
			if exists {
				/*
				 * if l3 entry exists then check whether l3 is running or not
				 * if l3 is started then incoming should be create
				 * if not then delete the pcap for incoming notification
				 */
				if l3exists {
					if l2Port.Info.OperState == config.STATE_UP && l3Port.PcapBase.Tx != nil {
						l2Port.createPortPcap(svr.RxPktCh, l2Port.Info.Name)
					} else {
						l2Port.deletePcap()
					}
				}
			}
		}
		svr.L2Port[pIfIndex] = l2Port
	}
	// open rx pcap handler for tagged ports
	for pIfIndex, _ := range vlan.TagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			debug.Logger.Debug("L2:", l2Port.Info.Name, "operstate is", l2Port.Info.OperState)
			if l2Port.Info.OperState == config.STATE_UP && l3Port.PcapBase.Tx != nil {
				l2Port.createPortPcap(svr.RxPktCh, l2Port.Info.Name)
			} else {
				l2Port.deletePcap()
			}
		}
		svr.L2Port[pIfIndex] = l2Port
	}
}

/*
 *  Utility Action function to delete ndp entries by L3 Port interface name
 */
func (svr *NDPServer) ActionDeleteByIntf(intfRef string) {
	ifIndex, exists := svr.L3IfIntfRefToIfIndex[intfRef]
	if !exists {
		debug.Logger.Err("Refresh Action by Interface Name:", intfRef,
			"cannot be performed as no ifIndex found for L3 interface")
		return
	}
	l3Port, exists := svr.L3Port[ifIndex]
	if !exists {
		debug.Logger.Err("Delete Action by Interface Name:", intfRef,
			"cannot be performed as no such L3 interface exists")
		return
	}
	deleteEntries, err := l3Port.FlushNeighbors()
	if len(deleteEntries) > 0 && err == nil {
		debug.Logger.Info("Server Action Delete by Intf:", l3Port.IntfRef, "Neighbors:", deleteEntries)
		svr.DeleteNeighborInfo(deleteEntries, ifIndex)
	}
	svr.L3Port[ifIndex] = l3Port
}

/*
 *  Utility Action function to refreshndp entries by L3 Port interface name
 */
func (svr *NDPServer) ActionRefreshByIntf(intfRef string) {
	ifIndex, exists := svr.L3IfIntfRefToIfIndex[intfRef]
	if !exists {
		debug.Logger.Err("Refresh Action by Interface Name:", intfRef,
			"cannot be performed as no ifIndex found for L3 interface")
		return
	}
	l3Port, exists := svr.L3Port[ifIndex]
	if !exists {
		debug.Logger.Err("Refresh Action by Interface Name:", intfRef,
			"cannot be performed as no such L3 interface exists")
		return
	}

	l3Port.RefreshAllNeighbors(svr.SwitchMac)
	svr.L3Port[ifIndex] = l3Port
}

/*
 *  Utility Action function to delete ndp entries by Neighbor Ip Address
 */
func (svr *NDPServer) ActionDeleteByNbrIp(ipAddr string) {
	debug.Logger.Info("performing delete action by ip address for ipAddr:", ipAddr)
	var nbrKey string
	found := false
	for _, nbrKey = range svr.neighborKey {
		debug.Logger.Debug("neighbor key in server neighbor keys:", nbrKey)
		splitString := splitNeighborKey(nbrKey)
		if splitString[1] == ipAddr {
			found = true
			break
		}
	}
	if !found {
		debug.Logger.Err("Delete Action by Ip Address:", ipAddr, "as no such neighbor is learned")
		return
	}
	debug.Logger.Debug("entry found in neighbor keys:", nbrKey)
	nbrEntry, exists := svr.NeighborInfo[nbrKey]
	if !exists {
		debug.Logger.Err("Delete Action by Ip Address:", ipAddr, "as no such neighbor is learned")
		return
	}
	debug.Logger.Debug("nbr entry found in NeighborInfo", nbrEntry)
	l3IfIndex := nbrEntry.IfIndex
	// if valid vlan then get l3 ifIndex from PhyPortToL3PortMap
	if nbrEntry.VlanId != config.INTERNAL_VLAN {
		l3Info, exists := svr.PhyPortToL3PortMap[nbrEntry.IfIndex]
		if !exists {
			debug.Logger.Err("Delete Action by Ip Address:", ipAddr,
				"cannot be performed as no l3IfIndex mapping found for", nbrEntry.IfIndex,
				"vlan:", nbrEntry.VlanId)
			return
		}
		l3IfIndex = l3Info.IfIndex
	}
	debug.Logger.Debug("l3Ifindex where the neighbor entry was learned:", l3IfIndex)
	l3Port, exists := svr.L3Port[l3IfIndex]
	if !exists {
		debug.Logger.Err("Delete Action by Ip Address:", ipAddr, "as no L3 Port found where this neighbor is learned")
		return
	}
	deleteEntries, err := l3Port.DeleteNeighbor(nbrEntry)
	if err != nil {
		debug.Logger.Info("Server Action Delete by NbrIp:", ipAddr, "L3 Port:", l3Port.IntfRef,
			"Neighbors:", deleteEntries, "failed")
	} else {
		debug.Logger.Info("Server Action Delete by NbrIp:", ipAddr, "L3 Port:", l3Port.IntfRef,
			"Neighbors:", deleteEntries)
		svr.deleteNeighbor(deleteEntries[0], l3Port.IfIndex)
		debug.Logger.Debug("delete action by ipAddr performed successfully")
	}
	svr.L3Port[l3IfIndex] = l3Port
}

/*
 *  Utility Action function to refresh ndp entries by Neighbor Ip Address
 */
func (svr *NDPServer) ActionRefreshByNbrIp(ipAddr string) {
	var nbrKey string
	found := false
	for _, nbrKey = range svr.neighborKey {
		splitString := splitNeighborKey(nbrKey)
		if splitString[1] == ipAddr {
			found = true
		}
	}
	if !found {
		debug.Logger.Err("Delete Action by Ip Address:", ipAddr, "as no such neighbor is learned")
		return
	}
	nbrEntry, exists := svr.NeighborInfo[nbrKey]
	if !exists {
		debug.Logger.Err("Refresh Action by Ip Address:", ipAddr, "as no such neighbor is learned")
		return
	}
	l3IfIndex := nbrEntry.IfIndex
	// if valid vlan then get l3 ifIndex from PhyPortToL3PortMap
	if nbrEntry.VlanId != config.INTERNAL_VLAN {
		l3Info, exists := svr.PhyPortToL3PortMap[nbrEntry.IfIndex]
		if !exists {
			debug.Logger.Err("Refresh Action by Ip Address:", ipAddr,
				"cannot be performed as no l3IfIndex mapping found for", nbrEntry.IfIndex,
				"vlan:", nbrEntry.VlanId)
			return
		}
		l3IfIndex = l3Info.IfIndex
	}

	l3Port, exists := svr.L3Port[l3IfIndex]
	if !exists {
		debug.Logger.Err("Delete Action by Ip Address:", ipAddr, "as no L3 Port found where this neighbor is learned")
		return
	}
	l3Port.SendNS(svr.SwitchMac, nbrEntry.MacAddr, nbrEntry.IpAddr, false /*isFastProbe*/)
	svr.L3Port[l3IfIndex] = l3Port
}

func (svr *NDPServer) updateVlanMemberFilters(vlan config.VlanInfo, macAddr string) {
	for pIfIndex, _ := range vlan.TagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.updateFilter(macAddr)
			svr.L2Port[pIfIndex] = l2Port
		}
	}
	for pIfIndex, _ := range vlan.UntagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.updateFilter(macAddr)
			svr.L2Port[pIfIndex] = l2Port
		}
	}
}

func (svr *NDPServer) resetVlanMemberFilters(vlan config.VlanInfo) {
	for pIfIndex, _ := range vlan.TagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.resetFilter()
			svr.L2Port[pIfIndex] = l2Port
		}
	}
	for pIfIndex, _ := range vlan.UntagPortsMap {
		l2Port, exists := svr.L2Port[pIfIndex]
		if exists {
			l2Port.resetFilter()
			svr.L2Port[pIfIndex] = l2Port
		}
	}
}

func (svr *NDPServer) updateFilters(virEntry *config.VirtualIpInfo) {
	l3Port, exists := svr.L3Port[virEntry.ParentIfIndex]
	if !exists {
		debug.Logger.Err("Should have never happened how can a virtual interface exists without an l3Port for:", *virEntry)
		return
	}
	switch l3Port.IfType {
	case commonDefs.IfTypePort:
		l3Port.updateFilter(virEntry.MacAddr)
		svr.L3Port[virEntry.ParentIfIndex] = l3Port
	case commonDefs.IfTypeVlan:
		vlan, exists := svr.VlanInfo[virEntry.ParentIfIndex]
		if !exists {
			debug.Logger.Err("Should have never happened how can a virtual interface exists without an l3Port for:", *virEntry)
			return
		}
		svr.updateVlanMemberFilters(vlan, virEntry.MacAddr)
	}
}

func (svr *NDPServer) resetFilters(virEntry *config.VirtualIpInfo) {
	l3Port, exists := svr.L3Port[virEntry.ParentIfIndex]
	if !exists {
		debug.Logger.Err("Should have never happened how can a virtual interface exists without an l3Port for:", *virEntry)
		return
	}
	switch l3Port.IfType {
	case commonDefs.IfTypePort:
		l3Port.resetFilter()
		svr.L3Port[virEntry.ParentIfIndex] = l3Port
	case commonDefs.IfTypeVlan:
		vlan, exists := svr.VlanInfo[virEntry.ParentIfIndex]
		if !exists {
			debug.Logger.Err("Should have never happened how can a virtual interface exists without an l3Port for:", *virEntry)
			return
		}
		svr.resetVlanMemberFilters(vlan)
	}

}

func (svr *NDPServer) HandleVirtualIpNotification(vipInfo *config.VirtualIpInfo) {
	switch vipInfo.MsgType {
	case config.VIRTUAL_CREATE:
		_, exists := svr.virtualIp[vipInfo.IfIndex]
		if !exists {
			virEntry := &config.VirtualIpInfo{}
			virEntry.IfIndex = vipInfo.IfIndex
			virEntry.ParentIfIndex = vipInfo.ParentIfIndex
			virEntry.IpAddr = vipInfo.IpAddr
			virEntry.MacAddr = vipInfo.MacAddr
			virEntry.IfName = vipInfo.IfName
			svr.virtualIp[vipInfo.IfIndex] = virEntry
			debug.Logger.Info("Added Virtual Ip Information to Runtime info:", *virEntry)
		}

	case config.VIRTUAL_DELETE:
		virEntry, exists := svr.virtualIp[vipInfo.IfIndex]
		if virEntry != nil {
			virEntry = nil
		}
		if exists {
			delete(svr.virtualIp, vipInfo.IfIndex)
		}
	case config.STATE_UP:
		virEntry, exists := svr.virtualIp[vipInfo.IfIndex]
		if exists {
			svr.updateFilters(virEntry)
		}

	case config.STATE_DOWN:
		virEntry, exists := svr.virtualIp[vipInfo.IfIndex]
		if exists {
			svr.resetFilters(virEntry)
		}
	}
}
