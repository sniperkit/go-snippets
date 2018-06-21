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

import ()

type ActionType uint8

const (
	DeleteByIPAddr  ActionType = 1
	RefreshByIPAddr ActionType = 2
	DeleteByIfName  ActionType = 3
	RefreshByIfName ActionType = 4
)

type ArpActionMsg struct {
	Type ActionType
	Obj  string
}

func (server *ARPServer) processDeleteByIPAddr(ipAddr string) {
	server.logger.Info("Delete Arp entry by IpAddr:", ipAddr)
	arpEnt, exist := server.arpCache[ipAddr]
	if !exist {
		server.logger.Warning("Cannot perform Arp delete action as Arp Entry does exist for ipAddr:", ipAddr)
		return
	}

	if arpEnt.Type == true {
		server.logger.Warning("Cannot perform Arp delete action as Arp Entry for", ipAddr, "belong to nexthop of some route, can only be deleted by RIB")
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

func (server *ARPServer) processDeleteByIfName(l3IfName string) {
	server.logger.Info("Delete Arp entry by IfName:", l3IfName)
	for l3IfIdx, l3Ent := range server.l3IntfPropMap {
		if l3Ent.IfName == l3IfName {
			for ip, arpEnt := range server.arpCache {
				if arpEnt.L3IfIdx == l3IfIdx {
					server.processDeleteByIPAddr(ip)
				}
			}
		}
	}
}

func (server *ARPServer) processRefreshByIPAddr(ipAddr string) {
	server.logger.Info("Refresh Arp entry by IpAddr:", ipAddr)
	arpEnt, exist := server.arpCache[ipAddr]
	if !exist {
		server.logger.Warning("Cannot perform Arp refresh action as Arp Entry does exist for ipAddr:", ipAddr)
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
	arpEnt.MacAddr = "incomplete"
	arpEnt.Counter = server.timeoutCounter
	server.arpCache[ipAddr] = arpEnt
	server.deleteLinuxArp(ipAddr)
}

func (server *ARPServer) processRefreshByIfName(l3IfName string) {
	server.logger.Info("Refresh Arp entry by IfName:", l3IfName)
	for l3IfIdx, l3Ent := range server.l3IntfPropMap {
		if l3Ent.IfName == l3IfName {
			for ip, arpEnt := range server.arpCache {
				if arpEnt.L3IfIdx == l3IfIdx {
					server.processRefreshByIPAddr(ip)
				}
			}
		}
	}
}

func (server *ARPServer) processArpActionMsg(msg ArpActionMsg) {
	switch msg.Type {
	case DeleteByIPAddr:
		server.processDeleteByIPAddr(msg.Obj)
	case DeleteByIfName:
		server.processDeleteByIfName(msg.Obj)
	case RefreshByIPAddr:
		server.processRefreshByIPAddr(msg.Obj)
	case RefreshByIfName:
		server.processRefreshByIfName(msg.Obj)
	}

}
