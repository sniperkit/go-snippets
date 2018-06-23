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
	"fmt"
	"time"
)

func (server *DHCPServer) StartStaleEntryHandler(port int32, ipAddr uint32, macAddr string) {
	server.logger.Info(fmt.Sprintln("Start Stale entry Handler....", ipAddr, macAddr))
	portEnt, _ := server.portPropertyMap[port]
	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/
	l3IfIdx := portEnt.L3IfIndex
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	dhcpIntfKey := l3Ent.DhcpIfKey

	removeStaleFunc := func() {
		server.logger.Info(fmt.Sprintln("Removing the stale entry ", ipAddr, macAddr))
		dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
		uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]
		uIPEnt.StaleTimer.Stop()
		delete(dhcpIntfEnt.usedIpPool, ipAddr)
		delete(dhcpIntfEnt.usedIpToMac, macAddr)
		server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
	}

	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]
	uIPEnt.StaleTimer = time.AfterFunc(time.Duration(server.DhcpGlobalConf.DefaultLeaseTime/2)*time.Second, removeStaleFunc)
	//server.logger.Info(fmt.Sprintln("3 uIPEnt: ", uIPEnt))
	dhcpIntfEnt.usedIpPool[ipAddr] = uIPEnt
	server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
}

func (server *DHCPServer) StartLeaseEntryHandler(port int32, ipAddr uint32, macAddr string) {
	server.logger.Info(fmt.Sprintln("Start lease entry Handler....", ipAddr, macAddr))
	portEnt, _ := server.portPropertyMap[port]
	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/

	l3IfIdx := portEnt.L3IfIndex
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	dhcpIntfKey := l3Ent.DhcpIfKey
	removeLeaseExpireFunc := func() {
		server.logger.Info(fmt.Sprintln("Removing the lease expiry entry ", ipAddr, macAddr))
		dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
		uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]
		//server.logger.Info(fmt.Sprintln("2 uIPEnt: ", uIPEnt))
		if uIPEnt.StaleTimer != nil {
			uIPEnt.StaleTimer.Stop()
		}
		uIPEnt.RefreshTimer.Stop()
		delete(dhcpIntfEnt.usedIpPool, ipAddr)
		delete(dhcpIntfEnt.usedIpToMac, macAddr)
		server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
	}
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]
	uIPEnt.RefreshTimer = time.AfterFunc(time.Duration(uIPEnt.LeaseTime)*time.Second, removeLeaseExpireFunc)
	uIPEnt.State = OFFERED
	//server.logger.Info(fmt.Sprintln("3 uIPEnt: ", uIPEnt))
	dhcpIntfEnt.usedIpPool[ipAddr] = uIPEnt
	server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt

}
