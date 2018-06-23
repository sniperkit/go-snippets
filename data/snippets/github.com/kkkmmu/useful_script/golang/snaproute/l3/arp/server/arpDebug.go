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

import ()

func (server *ARPServer) printArpEntries() {
	server.logger.Debug("************")
	server.logger.Debug("Time Out Counter:", server.timeoutCounter, "Timer granularity:", server.timerGranularity)
	for ip, arpEnt := range server.arpCache {
		server.logger.Debug("IP:", ip, "MAC:", arpEnt.MacAddr, "VlanId:", arpEnt.VlanId, "IfName:", arpEnt.IfName, "IfIndex", arpEnt.L3IfIdx, "Counter:", arpEnt.Counter, "Timestamp:", arpEnt.TimeStamp, "PortNum:", arpEnt.PortNum, "Type:", arpEnt.Type)
	}
	server.logger.Debug("************")
}

func (server *ARPServer) dumpInfra() {
	server.dumpL3IntfProp()
	server.dumpVlanProp()
	server.dumpLagProp()
	server.dumpPortProp()
}

func (server *ARPServer) dumpL3IntfProp() {
	server.logger.Debug("==================================================")
	server.logger.Debug("L3 Interface Property Map:")
	for l3IfIndex, l3Ent := range server.l3IntfPropMap {
		server.logger.Debug("L3 IfIndex:", l3IfIndex, "IpAddr:", l3Ent.IpAddr, "Netmask:", l3Ent.Netmask, "IfName:", l3Ent.IfName)
	}
	server.logger.Debug("==================================================")
}

func (server *ARPServer) dumpVlanProp() {
	server.logger.Debug("==================================================")
	server.logger.Debug("Vlan Property Map:")
	for vlanIfIdx, vlanEnt := range server.vlanPropMap {
		server.logger.Debug("Vlan IfIdx:", vlanIfIdx, "Vlan Name:", vlanEnt.IfName)
		server.logger.Debug("Untagged IfIndex Map:")
		for uIfIdx, _ := range vlanEnt.UntagIfIdxMap {
			server.logger.Debug(uIfIdx)
		}
		server.logger.Debug("Tagged IfIndex Map:")
		for tIfIdx, _ := range vlanEnt.TagIfIdxMap {
			server.logger.Debug(tIfIdx)
		}
	}
	server.logger.Debug("==================================================")
}

func (server *ARPServer) dumpLagProp() {
	server.logger.Debug("==================================================")
	server.logger.Debug("Lag Property Map:")
	for lagIfIdx, lagEnt := range server.lagPropMap {
		server.logger.Debug("Lag IfIdx:", lagIfIdx, "Lag Name:", lagEnt.IfName, "UntagVlanId:", lagEnt.UntagVlanId)
		server.logger.Debug("Port IfIndex Map:")
		for ifIdx, _ := range lagEnt.PortMap {
			server.logger.Debug(ifIdx)
		}
		server.logger.Debug("Vlan Id Map:")
		for id, _ := range lagEnt.VlanIdMap {
			server.logger.Debug(id)
		}
	}
	server.logger.Debug("==================================================")
}

func (server *ARPServer) dumpPortProp() {
	server.logger.Debug("==================================================")
	server.logger.Debug(server.portPropMap)
	server.logger.Debug("==================================================")
}
