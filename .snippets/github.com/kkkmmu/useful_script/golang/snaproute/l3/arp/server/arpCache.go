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
	"models/events"
	"time"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/eventUtils"
)

type UpdateArpEntryMsg struct {
	PortIfIdx int
	IpAddr    string
	MacAddr   string
	Type      bool // True: RIB False: Rx
	VlanId    int
	L3IfIdx   int
	LagIfIdx  int
}

const (
	DeleteBasedOnPort uint8 = 0
	DeleteBasedOnVlan uint8 = 1
	DeleteBasedOnL3   uint8 = 2
)

type DeleteArpEntryMsg struct {
	Type  uint8
	IfIdx int
}

type DeleteResolvedIPv4 struct {
	IpAddr string
}

type DeleteArpEntry struct {
	IpAddr string
}

type EventData struct {
	IpAddr  string
	MacAddr string
	IfName  string
}

func (server *ARPServer) PublishEvents(IpAddr, MacAddr, IfName string, eventId events.EventId) {
	evtKey := events.ArpEntryKey{
		IpAddr: IpAddr,
	}
	evtData := EventData{
		IpAddr:  IpAddr,
		MacAddr: MacAddr,
		IfName:  IfName,
	}
	txEvent := eventUtils.TxEvent{
		EventId:        eventId,
		Key:            evtKey,
		AdditionalInfo: "",
		AdditionalData: evtData,
	}
	err := eventUtils.PublishEvents(&txEvent)
	if err != nil {
		server.logger.Err("Error in publishing ArpEntryLearned Event")
	}
}

func (server *ARPServer) updateArpCache() {
	for {
		select {
		case msg := <-server.arpEntryUpdateCh:
			server.processArpEntryUpdateMsg(msg)
		case msg := <-server.arpEntryDeleteCh:
			server.processArpEntryDeleteMsg(msg)
		case <-server.arpSliceRefreshStartCh:
			server.processArpSliceRefreshMsg()
		case <-server.arpCounterUpdateCh:
			server.processArpCounterUpdateMsg()
		case cnt := <-server.arpEntryCntUpdateCh:
			server.processArpEntryCntUpdateMsg(cnt)
		case msg := <-server.arpEntryMacMoveCh:
			server.processArpEntryMacMoveMsg(msg)
		case msg := <-server.arpDeleteArpEntryFromRibCh:
			server.processArpEntryDeleteMsgFromRib(msg)
		case msg := <-server.arpActionProcessCh:
			server.processArpActionMsg(msg)
		case msg, ok := <-server.arpDeleteArpEntryIntCh:
			if ok {
				server.processArpEntryDeleteInt(msg)
			}
		case garpInfo, ok := <-server.GarpEntryCh:
			if ok {
				server.processGarp(garpInfo)
			}
		}
	}
}

func (server *ARPServer) processArpEntryDeleteInt(ipAddr string) {
	arpEnt, exists := server.arpCache[ipAddr]
	if !exists {
		server.logger.Err("Cannot perform arp delete action as no arp entry found for ipAddr:", ipAddr)
		return
	}
	if arpEnt.Type {
		server.logger.Debug("Arp Entry for ipAddr:", ipAddr, "was installed by RIB, hence cannot be deleted")
		return
	}

	if arpEnt.MacAddr != "incomplete" {
		server.logger.Debug("4 Calling Asicd Delete Ip:", ipAddr)
		asicdMsg := AsicdMsg{
			MsgType: DeleteAsicdEntry,
			IpAddr:  ipAddr,
		}
		err := server.processAsicdMsg(asicdMsg)
		if err != nil {
			return
		}
	}

	delete(server.arpCache, ipAddr)
	server.deleteLinuxArp(ipAddr)
}

func (server *ARPServer) processArpEntryDeleteMsgFromRib(ipAddr string) {
	arpEnt, exist := server.arpCache[ipAddr]
	if !exist {
		server.logger.Warning("Cannot perform Arp delete action as Arp Entry does exist for ipAddr:", ipAddr)
		return
	}

	if arpEnt.Type == false {
		server.logger.Debug("Arp Entry for IpAddr:", ipAddr, "was no installed by RIB, hence cannot be delete")
		return
	}

	if arpEnt.MacAddr != "incomplete" {
		server.logger.Debug("4 Calling Asicd Delete Ip:", ipAddr)
		asicdMsg := AsicdMsg{
			MsgType: DeleteAsicdEntry,
			IpAddr:  ipAddr,
		}
		err := server.processAsicdMsg(asicdMsg)
		if err != nil {
			return
		}
	}
	delete(server.arpCache, ipAddr)
	server.deleteLinuxArp(ipAddr)
}

func (server *ARPServer) processArpEntryCntUpdateMsg(cnt int) {
	for key, ent := range server.arpCache {
		if ent.Counter > cnt {
			ent.Counter = cnt
			ent.TimeStamp = time.Now()
			server.arpCache[key] = ent
		}
	}
}

func (server *ARPServer) processArpEntryMacMoveMsg(msg asicdClntDefs.IPv4NbrMacMoveNotifyMsg) {
	if entry, ok := server.arpCache[msg.IpAddr]; ok {
		entry.PortNum = int(msg.IfIndex)
		server.arpCache[msg.IpAddr] = entry
		server.PublishEvents(msg.IpAddr, entry.MacAddr, entry.IfName, events.ArpEntryUpdated)
	} else {
		server.logger.Debug("Mac move message received. Neighbor IP does not exist in arp cache - %x", msg.IpAddr)
	}
}

func (server *ARPServer) processArpEntryDeleteMsg(msg DeleteArpEntryMsg) {
	server.logger.Debug("DeleteArpEntryMsg:", msg)
	for key, ent := range server.arpCache {
		flag := false
		switch msg.Type {
		case DeleteBasedOnL3:
			if msg.IfIdx == ent.L3IfIdx {
				flag = true
			}
		case DeleteBasedOnVlan:
			if msg.IfIdx == ent.VlanId {
				flag = true
			}
		case DeleteBasedOnPort:
			if msg.IfIdx == ent.PortNum {
				flag = true
			}
		}
		if flag == true {
			server.logger.Debug("1 Calling Asicd Delete Ip:", key)
			asicdMsg := AsicdMsg{
				MsgType: DeleteAsicdEntry,
				IpAddr:  key,
			}
			err := server.processAsicdMsg(asicdMsg)
			if err != nil {
				return
			}
			delete(server.arpCache, key)
			server.deleteArpEntryInDB(key)
			server.deleteLinuxArp(key)
		}
	}
}

func (server *ARPServer) processArpIncompleteUpdateMsg(IpAddr string, Type bool, l3IfIdx int) {
	arpEnt, exist := server.arpCache[IpAddr]
	if exist {
		if l3IfIdx == arpEnt.L3IfIdx {
			if arpEnt.Type == false {
				arpEnt.Type = Type
			}
			if arpEnt.MacAddr == "incomplete" {
				arpEnt.Counter = server.timeoutCounter
			}
		} else {
			server.logger.Err("Error: l3IfIdx does not match arp cache entry", l3IfIdx, arpEnt.L3IfIdx)
			return
		}
	} else {
		arpEnt.Type = Type
		arpEnt.L3IfIdx = l3IfIdx
		arpEnt.Counter = server.timeoutCounter
		arpEnt.MacAddr = "incomplete"
		arpEnt.PortNum = -1
		arpEnt.VlanId = -1
		arpEnt.IfName = ""
		server.storeArpEntryInDB(IpAddr, l3IfIdx)
	}
	server.arpCache[IpAddr] = arpEnt
	server.updateArpSlice(IpAddr)
}

func (server *ARPServer) processLearnedArpEntryUpdateMsg(msg UpdateArpEntryMsg) {
	var ifIdx int
	portEnt, _ := server.portPropMap[msg.PortIfIdx]
	ifName := portEnt.IfName
	arpEnt, exist := server.arpCache[msg.IpAddr]
	if exist {
		if arpEnt.L3IfIdx == msg.L3IfIdx &&
			arpEnt.MacAddr == msg.MacAddr &&
			arpEnt.VlanId == msg.VlanId &&
			arpEnt.PortNum == msg.PortIfIdx {
			arpEnt.Counter = server.timeoutCounter
			arpEnt.TimeStamp = time.Now()
			server.arpCache[msg.IpAddr] = arpEnt
			return
		}
	}
	if msg.LagIfIdx != -1 {
		ifIdx = msg.LagIfIdx
	} else {
		ifIdx = msg.PortIfIdx
	}
	server.logger.Debug("3 Calling Asicd Create Ip:", msg.IpAddr, "mac:", msg.MacAddr, "vlanId:", msg.VlanId, "IfIndex:", ifIdx)
	asicdMsg := AsicdMsg{
		MsgType: CreateAsicdEntry,
		IpAddr:  msg.IpAddr,
		MacAddr: msg.MacAddr,
		VlanId:  int32(msg.VlanId),
		IfIdx:   int32(ifIdx),
	}
	err := server.processAsicdMsg(asicdMsg)
	if err != nil {
		return
	}
	if !exist {
		server.PublishEvents(msg.IpAddr, msg.MacAddr, ifName, events.ArpEntryLearned)
		server.storeArpEntryInDB(msg.IpAddr, msg.L3IfIdx)
	}
	arpEnt.Counter = server.timeoutCounter
	arpEnt.IfName = ifName
	arpEnt.L3IfIdx = msg.L3IfIdx
	arpEnt.MacAddr = msg.MacAddr
	arpEnt.TimeStamp = time.Now()
	arpEnt.PortNum = msg.PortIfIdx
	arpEnt.VlanId = msg.VlanId
	if arpEnt.Type == false {
		arpEnt.Type = msg.Type
	}
	server.arpCache[msg.IpAddr] = arpEnt
	server.updateArpSlice(msg.IpAddr)
}

func (server *ARPServer) updateArpSlice(IpAddr string) {
	for i := 0; i < len(server.arpSlice); i++ {
		if server.arpSlice[i] == IpAddr {
			return
		}
	}
	server.arpSlice = append(server.arpSlice, IpAddr)
}

func (server *ARPServer) processArpEntryUpdateMsg(msg UpdateArpEntryMsg) {
	if msg.MacAddr == "incomplete" {
		server.processArpIncompleteUpdateMsg(msg.IpAddr, msg.Type, msg.L3IfIdx)
	} else {
		server.processLearnedArpEntryUpdateMsg(msg)
	}
}

func (server *ARPServer) processArpCounterUpdateMsg() {
	oneMinCnt := (60 / server.timerGranularity)
	thirtySecCnt := (30 / server.timerGranularity)
	for ip, arpEnt := range server.arpCache {
		if arpEnt.Counter <= server.minCnt {
			if arpEnt.Type == false {
				server.deleteArpEntryInDB(ip)
				delete(server.arpCache, ip)
				if arpEnt.MacAddr != "incomplete" {
					server.logger.Debug("5 Calling Asicd Delete Ip:", ip)
					asicdMsg := AsicdMsg{
						MsgType: DeleteAsicdEntry,
						IpAddr:  ip,
					}
					err := server.processAsicdMsg(asicdMsg)
					if err != nil {
						continue
					}
					server.deleteLinuxArp(ip)
					server.PublishEvents(ip, arpEnt.MacAddr, arpEnt.IfName, events.ArpEntryDeleted)
				}
				server.printArpEntries()
			} else {
				//server.logger.Debug("Nexthop", ip, " installed by Rib hence not deleting it")
				if arpEnt.MacAddr != "incomplete" {
					server.logger.Debug("5 Calling Asicd Delete Ip:", ip)
					asicdMsg := AsicdMsg{
						MsgType: DeleteAsicdEntry,
						IpAddr:  ip,
					}
					server.deleteLinuxArp(ip)
					err := server.processAsicdMsg(asicdMsg)
					if err != nil {
						continue
					}
				}
				//server.logger.Debug("Reseting the counter to max", ip)
				arpEnt.MacAddr = "incomplete"
				arpEnt.Counter = server.timeoutCounter
				server.arpCache[ip] = arpEnt
			}
		} else {
			arpEnt.Counter--
			server.arpCache[ip] = arpEnt
			if arpEnt.Counter <= (server.minCnt+server.retryCnt+1) ||
				arpEnt.Counter == (server.timeoutCounter/2) ||
				arpEnt.Counter == (server.timeoutCounter/4) ||
				arpEnt.Counter == oneMinCnt ||
				arpEnt.Counter == thirtySecCnt {
				if arpEnt.MacAddr == "incomplete" {
					server.retryForArpEntry(ip, arpEnt.L3IfIdx)
				} else {
					server.refreshArpEntry(ip, arpEnt.L3IfIdx)
				}
			} else if arpEnt.Counter <= server.timeoutCounter &&
				arpEnt.Counter > (server.timeoutCounter-server.retryCnt) &&
				arpEnt.MacAddr == "incomplete" {
				server.retryForArpEntry(ip, arpEnt.L3IfIdx)
			} else if arpEnt.Counter > (server.minCnt+server.retryCnt+1) &&
				arpEnt.MacAddr != "incomplete" {
				continue
			} else {
				if arpEnt.Type == false {
					server.deleteArpEntryInDB(ip)
					delete(server.arpCache, ip)
					server.printArpEntries()
				} else {
					//server.logger.Debug("Nexthop", ip, " installed by Rib hence not deleting it")
				}
			}
		}
	}
}

func (server *ARPServer) refreshArpEntry(ipAddr string, l3IfIdx int) {
	// TimeoutCounter set to retryCnt
	server.logger.Debug("Refreshing Arp entry for IP:", ipAddr, "on port:", l3IfIdx)
	server.sendArpReq(ipAddr, l3IfIdx)
}

func (server *ARPServer) retryForArpEntry(ipAddr string, l3IfIdx int) {
	server.logger.Debug("Retrying Arp entry for IP:", ipAddr, "on port:", l3IfIdx)
	server.sendArpReq(ipAddr, l3IfIdx)
}

func (server *ARPServer) processArpSliceRefreshMsg() {
	server.logger.Debug("Refresh Arp Slice used for Getbulk")
	server.arpSlice = server.arpSlice[:0]
	server.arpSlice = nil
	server.arpSlice = make([]string, 0)
	for ip, _ := range server.arpCache {
		server.arpSlice = append(server.arpSlice, ip)
	}
	server.arpSliceRefreshDoneCh <- true
}

func (server *ARPServer) refreshArpSlice() {
	refreshArpSlicefunc := func() {
		server.arpSliceRefreshStartCh <- true
		msg := <-server.arpSliceRefreshDoneCh
		if msg == true {
			server.logger.Debug("ARP Entry refresh done")
		} else {
			server.logger.Err("ARP Entry refresh not done")
		}

		server.arpSliceRefreshTimer.Reset(server.arpSliceRefreshDuration)
	}

	server.arpSliceRefreshTimer = time.AfterFunc(server.arpSliceRefreshDuration, refreshArpSlicefunc)
}

func (server *ARPServer) arpCacheTimeout() {
	var count int
	for {
		time.Sleep(server.timeout)
		count++
		if server.dumpArpTable == true &&
			(count%60) == 0 {
			server.logger.Debug("===============Message from ARP Timeout Thread==============")
			server.printArpEntries()
			server.logger.Debug("========================================================")
			server.logger.Debug("Arp Slice: ", server.arpSlice)
			server.dumpInfra()
		}
		server.arpCounterUpdateCh <- true
	}
}
