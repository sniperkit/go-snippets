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
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

func (server *DHCPServer) processDhcpGlobalConf(conf DhcpGlobalConfig) {
	server.logger.Info(fmt.Sprintln("Received DHCP Global Configuration:", conf))
	server.DhcpGlobalConf.DefaultLeaseTime = conf.DefaultLeaseTime
	server.DhcpGlobalConf.MaxLeaseTime = conf.MaxLeaseTime
	if server.DhcpGlobalConf.Enable != conf.Enable {
		server.DhcpGlobalConf.Enable = conf.Enable
		server.handleDhcpGlobalConf()
	}
	return
}

func (server *DHCPServer) processDhcpIntfConf(conf DhcpIntfConfig) (error, int32) {
	server.logger.Info(fmt.Sprintln("Received DHCP Interface Configuration:", conf, server.l3IntfPropMap))
	intfRef := conf.IntfRef
	l3IntfIdx, _ := strconv.Atoi(intfRef) // Need to re visited
	l3IfIdx := int32(l3IntfIdx)
	l3Ent, exist := server.l3IntfPropMap[l3IfIdx]
	if !exist {
		err := errors.New("No L3 Interface on the given IntfRef, first configure L3 Interface")
		return err, l3IfIdx
	}

	if l3Ent.DhcpConfig == true {
		err := errors.New("Dhcp Server already Configured on this L3 interface")
		return err, l3IfIdx
	}

	dhcpIntfKey := DhcpIntfKey{
		subnet:     conf.Subnet,
		subnetMask: conf.SubnetMask,
	}
	/*
		if _, ok := server.l3PropertyMap[dhcpIntfKey]; !ok {
			err := errors.New("No L3 Interface in the given subnet, first configure L3 Interface")
			return err, dhcpIntfKey
		}
	*/

	l3Ent.DhcpConfig = true
	l3Ent.DhcpIfKey = dhcpIntfKey
	dhcpIntfEnt := server.DhcpIntfConfMap[dhcpIntfKey]
	dhcpIntfEnt.l3IfIdx = l3IfIdx
	dhcpIntfEnt.enable = conf.Enable
	dhcpIntfEnt.lowerIPBound = conf.LowerIPBound
	dhcpIntfEnt.higherIPBound = conf.HigherIPBound
	dhcpIntfEnt.rtrAddr = conf.RtrAddr
	dhcpIntfEnt.dnsAddr = conf.DnsAddr
	dhcpIntfEnt.domainName = conf.DomainName
	dhcpIntfEnt.usedIpPool = make(map[uint32]DhcpOfferedData)
	dhcpIntfEnt.usedIpToMac = make(map[string]uint32)
	server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
	server.l3IntfPropMap[l3IfIdx] = l3Ent
	return nil, l3IfIdx
}

func (server *DHCPServer) handleDhcpGlobalConf() {
	if server.DhcpGlobalConf.Enable == false {
		server.StopAllDhcpServer()
	} else {
		server.StartAllDhcpServer()
	}
}

func (server *DHCPServer) handleDhcpIntfConf(l3IfIdx int32) {
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	dhcpIntfKey := l3Ent.DhcpIfKey
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	if server.DhcpGlobalConf.Enable == true &&
		dhcpIntfEnt.enable == true {
		//server.GetAllStaticIpFromArp(server.l3PropertyMap[dhcpIntfKey])
		server.StartDhcpServer(l3IfIdx)
	}
}

func (server *DHCPServer) constructDhcpMsg(l3IfIdx int32) {
	l3Ent, _ := server.l3IntfPropMap[l3IfIdx]
	/*
		dhcpKey := DhcpIntfKey{
			subnet:     l3Ent.IpAddr & l3Ent.Mask,
			subnetMask: l3Ent.Mask,
		}
	*/
	dhcpKey := l3Ent.DhcpIfKey
	dhcpDataEnt, _ := server.DhcpIntfConfMap[dhcpKey]
	dhcpDataEnt.dhcpMsg = make([]byte, 36)
	binary.BigEndian.PutUint32(dhcpDataEnt.dhcpMsg[0:4], 0x63825363)
	dhcpDataEnt.dhcpMsg[4] = DhcpMsgTypeOptCode
	dhcpDataEnt.dhcpMsg[5] = uint8(1)
	dhcpDataEnt.dhcpMsg[6] = uint8(0)
	dhcpDataEnt.dhcpMsg[7] = ServerIdOptCode
	dhcpDataEnt.dhcpMsg[8] = uint8(4)
	binary.BigEndian.PutUint32(dhcpDataEnt.dhcpMsg[9:13], l3Ent.IpAddr)
	dhcpDataEnt.dhcpMsg[13] = IPAddrLeaseTimeOptCode
	dhcpDataEnt.dhcpMsg[14] = uint8(4)
	binary.BigEndian.PutUint32(dhcpDataEnt.dhcpMsg[15:19], server.DhcpGlobalConf.DefaultLeaseTime)
	dhcpDataEnt.dhcpMsg[19] = SubnetMaskOptCode
	dhcpDataEnt.dhcpMsg[20] = uint8(4)
	binary.BigEndian.PutUint32(dhcpDataEnt.dhcpMsg[21:25], dhcpKey.subnetMask)
	dhcpDataEnt.dhcpMsg[25] = RouterOptCode
	dhcpDataEnt.dhcpMsg[26] = uint8(4)
	binary.BigEndian.PutUint32(dhcpDataEnt.dhcpMsg[27:31], dhcpDataEnt.rtrAddr)
	dhcpDataEnt.dhcpMsg[31] = EndOptCode
	server.DhcpIntfConfMap[dhcpKey] = dhcpDataEnt
	//server.logger.Info(fmt.Sprintln("====Hello=======", dhcpDataEnt))
}
