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
	"fmt"
	"github.com/google/gopacket/pcap"
	"net"
	"utils/commonDefs"
)

type L3Property struct {
	//State      uint8
	IpAddr     uint32
	Mask       uint32
	DhcpIfKey  DhcpIntfKey
	DhcpConfig bool
}

type PortProperty struct {
	IfName    string
	MacAddr   string
	IpAddr    uint32
	Mask      uint32
	L3IfIndex int32
	CtrlCh    chan bool
	PcapHdl   *pcap.Handle
}

type VlanProperty struct {
	UntagPortMap map[int32]bool
}

type LagProperty struct {
	PortMap map[int32]bool
}

func (server *DHCPServer) processIPv4IntfCreate(msg asicdCommonDefs.IPv4IntfNotifyMsg) {
	ip, ipNet, _ := net.ParseCIDR(msg.IpAddr)
	server.logger.Debug(fmt.Sprintln("Received IPv4 Create Notification for IP:", ip))

	ipAddr, _ := convertIPStrToUint32(ip.String())
	mask := convertIPv4ToUint32(ipNet.Mask)

	ifType := asicdCommonDefs.GetIntfTypeFromIfIndex(msg.IfIndex)
	ifIdx := msg.IfIndex
	ent, _ := server.l3IntfPropMap[ifIdx]
	ent.IpAddr = ipAddr
	ent.Mask = mask
	//ent.State = uint8(asicdCommonDefs.INTF_STATE_UP)
	ent.DhcpConfig = false
	server.l3IntfPropMap[ifIdx] = ent
	/*
		dhcpKey := DhcpIntfKey{
			subnet:     ipAddr & mask,
			subnetMask: mask,
		}
		server.l3PropertyMap[dhcpKey] = ifIdx
	*/
	if ifType == commonDefs.IfTypeVlan {
		vlanEnt, _ := server.vlanPropertyMap[ifIdx]
		server.logger.Debug(fmt.Sprintln("Received IPv4 Create Notification for Untag Port List:", vlanEnt.UntagPortMap))
		for port, _ := range vlanEnt.UntagPortMap {
			portEnt := server.portPropertyMap[port]
			portEnt.IpAddr = ipAddr
			portEnt.Mask = mask
			portEnt.L3IfIndex = ifIdx
			server.portPropertyMap[port] = portEnt
		}
	} else if ifType == commonDefs.IfTypeLag {
		lagEnt, _ := server.lagPropertyMap[ifIdx]
		for port, _ := range lagEnt.PortMap {
			portEnt := server.portPropertyMap[port]
			portEnt.IpAddr = ipAddr
			portEnt.Mask = mask
			portEnt.L3IfIndex = ifIdx
			server.portPropertyMap[port] = portEnt

		}
	} else {
		portEnt := server.portPropertyMap[ifIdx]
		portEnt.IpAddr = ipAddr
		portEnt.Mask = mask
		portEnt.L3IfIndex = ifIdx
		server.portPropertyMap[ifIdx] = portEnt
	}
}

func (server *DHCPServer) processIPv4IntfDelete(msg asicdCommonDefs.IPv4IntfNotifyMsg) {
	ip, ipNet, _ := net.ParseCIDR(msg.IpAddr)
	server.logger.Debug(fmt.Sprintln("Received IPv4 Delete Notification for IP:", ip, ipNet))
	/*

		ipAddr, _ := convertIPStrToUint32(ip.String())
		mask := convertIPv4ToUint32(ipNet.Mask)
		dhcpKey := DhcpIntfKey{
			subnet:     ipAddr & mask,
			subnetMask: mask,
		}
	*/
	ifType := asicdCommonDefs.GetIntfTypeFromIfIndex(msg.IfIndex)
	ifIdx := msg.IfIndex
	if ifType == commonDefs.IfTypeVlan {
		vlanEnt, _ := server.vlanPropertyMap[ifIdx]
		for port, _ := range vlanEnt.UntagPortMap {
			portEnt := server.portPropertyMap[port]
			/*
				if portEnt.PcapHdl != nil {
					portEnt.CtrlCh <- true
					<-portEnt.CtrlCh
					// Delete Dhcp configuration for this L3Interface
				}
			*/
			portEnt.IpAddr = 0
			portEnt.Mask = 0
			portEnt.L3IfIndex = -1
			server.portPropertyMap[port] = portEnt
		}
	} else if ifType == commonDefs.IfTypeLag {
		lagEnt, _ := server.lagPropertyMap[ifIdx]
		for port, _ := range lagEnt.PortMap {
			portEnt := server.portPropertyMap[port]
			/*
				if portEnt.PcapHdl != nil {
					portEnt.CtrlCh <- true
					<-portEnt.CtrlCh
					// Delete Dhcp configuration for this L3Interface
				}
			*/
			portEnt.IpAddr = 0
			portEnt.Mask = 0
			portEnt.L3IfIndex = -1
			server.portPropertyMap[port] = portEnt
		}
	} else {
		portEnt := server.portPropertyMap[ifIdx]
		/*
			if portEnt.PcapHdl != nil {
				portEnt.CtrlCh <- true
				<-portEnt.CtrlCh
				// Delete Dhcp configuration for this L3Interface
			}
		*/
		portEnt.IpAddr = 0
		portEnt.Mask = 0
		portEnt.L3IfIndex = -1
		server.portPropertyMap[ifIdx] = portEnt
	}

	//delete(server.l3PropertyMap, dhcpKey)
	delete(server.l3IntfPropMap, ifIdx)
}

func (server *DHCPServer) updateIpv4Infra(msg asicdCommonDefs.IPv4IntfNotifyMsg, msgType uint8) {
	if msgType == asicdCommonDefs.NOTIFY_IPV4INTF_CREATE {
		server.processIPv4IntfCreate(msg)
	} else {
		server.processIPv4IntfDelete(msg)
	}
}

/*
func (server *DHCPServer) processL3StateChange(msg asicdCommonDefs.L3IntfStateNotifyMsg) {
	ifIdx := msg.IfIndex
	ifType := asicdCommonDefs.GetIntfTypeFromIfIndex(msg.IfIndex)
	l3Ent, _ := server.l3IntfPropMap[ifIdx]
	if msg.IfState == 0 {
		if ifType == commonDefs.IfTypeVlan {
			vlanEnt := server.vlanPropertyMap[ifIdx]
			for port, _ := range vlanEnt.UntagPortMap {
				server.logger.Debug(fmt.Sprintln("Stopping Dhcp server port:", port))
				portEnt, _ := server.portPropertyMap[port]
				if portEnt.PcapHdl != nil {
					portEnt.CtrlCh <- true
					<-portEnt.CtrlCh
					portEnt.PcapHdl = nil
					server.portPropertyMap[port] = portEnt
				}
			}
		} else if ifType == commonDefs.IfTypeLag {
			lagEnt := server.lagPropertyMap[ifIdx]
			for port, _ := range lagEnt.PortMap {
				server.logger.Info(fmt.Sprintln("Stopping Dhcp server port:", port))
				portEnt, _ := server.portPropertyMap[port]
				if portEnt.PcapHdl != nil {
					portEnt.CtrlCh <- true
					<-portEnt.CtrlCh
					portEnt.PcapHdl = nil
					server.portPropertyMap[port] = portEnt
				}
			}
		} else if ifType == commonDefs.IfTypePort {
			port := ifIdx
			server.logger.Info(fmt.Sprintln("Stopping Dhcp server port:", port))
			portEnt, _ := server.portPropertyMap[port]
			if portEnt.PcapHdl != nil {
				portEnt.CtrlCh <- true
				<-portEnt.CtrlCh
				portEnt.PcapHdl = nil
				server.portPropertyMap[port] = portEnt
			}
		}
		l3Ent.State = 0
	} else {
		l3Ent.State = 1
	}
	server.l3IntfPropMap[ifIdx] = l3Ent
}
*/

func (server *DHCPServer) processDhcpInfra() {
	// TODO
	// Used during restart
	// Go thru all the L3 interfaces
	// Enbale and start rx and tx thread for which it is configured
}

func (server *DHCPServer) buildDhcpInfra() {
	server.constructPortInfra()
	server.constructVlanInfra()
	//server.constructLagInfra()
	server.constructL3Infra()
}

func (server *DHCPServer) constructL3Infra() {
	curMark := 0
	server.logger.Debug("Calling Asicd for getting L3 Interfaces")
	count := 100
	for {
		bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkIPv4IntfState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ip, ipNet, _ := net.ParseCIDR(bulkInfo.IPv4IntfStateList[i].IpAddr)
			ipAddr, _ := convertIPStrToUint32(ip.String())
			mask := convertIPv4ToUint32(ipNet.Mask)
			ifIdx := bulkInfo.IPv4IntfStateList[i].IfIndex
			ifType := asicdCommonDefs.GetIntfTypeFromIfIndex(ifIdx)
			ent, _ := server.l3IntfPropMap[ifIdx]
			if ifType == commonDefs.IfTypeVlan {
				vlanEnt, _ := server.vlanPropertyMap[ifIdx]
				for port, _ := range vlanEnt.UntagPortMap {
					portEnt, _ := server.portPropertyMap[port]
					portEnt.IpAddr = ipAddr
					portEnt.Mask = mask
					portEnt.L3IfIndex = ifIdx
					server.portPropertyMap[port] = portEnt
				}
			} else if ifType == commonDefs.IfTypeLag {
				lagEnt, _ := server.lagPropertyMap[ifIdx]
				for port, _ := range lagEnt.PortMap {
					portEnt, _ := server.portPropertyMap[port]
					portEnt.IpAddr = ipAddr
					portEnt.Mask = mask
					portEnt.L3IfIndex = ifIdx
					server.portPropertyMap[port] = portEnt
				}
			} else {
				port := ifIdx
				portEnt, _ := server.portPropertyMap[port]
				portEnt.IpAddr = ipAddr
				portEnt.Mask = mask
				portEnt.L3IfIndex = ifIdx
				server.portPropertyMap[port] = portEnt

			}
			/*
				if bulkInfo.IPv4IntfStateList[i].OperState == "UP" {
					ent.State = uint8(1)
				} else {
					ent.State = uint8(0)
				}
			*/
			ent.DhcpConfig = false
			ent.IpAddr = ipAddr
			ent.Mask = mask
			server.l3IntfPropMap[ifIdx] = ent
		}
		if more == false {
			break
		}
	}
}

func (server *DHCPServer) constructPortInfra() {
	server.getBulkPortState()
	server.getBulkPortConfig()
}

func (server *DHCPServer) getBulkPortConfig() {
	curMark := int(asicdCommonDefs.MIN_SYS_PORTS)
	server.logger.Debug("Calling Asicd for getting Port Property")
	count := 100
	for {
		bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkPort(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIndex := bulkInfo.PortList[i].IfIndex
			ent := server.portPropertyMap[ifIndex]
			ent.MacAddr = bulkInfo.PortList[i].MacAddr
			server.portPropertyMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
}

func (server *DHCPServer) getBulkPortState() {
	curMark := int(asicdCommonDefs.MIN_SYS_PORTS)
	server.logger.Debug("Calling Asicd for getting Port Property")
	count := 100
	for {
		bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkPortState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIndex := bulkInfo.PortStateList[i].IfIndex
			ent := server.portPropertyMap[ifIndex]
			ent.IfName = bulkInfo.PortStateList[i].Name
			ent.L3IfIndex = -1
			ent.CtrlCh = make(chan bool)
			ent.PcapHdl = nil
			ent.IpAddr = 0
			ent.Mask = 0
			server.portPropertyMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
}

/*
func (server *DHCPServer) constructLagInfra() {
	curMark := 0
	server.logger.Info("Calling Asicd for getting Lag Property")
	count := 100
	for {
		bulkInfo, _ := server.asicdClient.ClientHdl.GetBulkLag(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkInfo == nil {
			break
		}
		objCnt := int(bulkInfo.Count)
		more := bool(bulkInfo.More)
		curMark = asicdServices.Int(bulkInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIdx := bulkInfo.LagList[i].IfIndex
			ent := server.lagPropertyMap[ifIdx]
			ifIndexList := ParseUsrPortStrToPortList(bulkInfo.LagList[i].IfIndexList)
			for i := 0; i < len(ifIndexList); i++ {
				ent.PortMap[port] = true
			}
			server.lagPropertyMap[ifIdx] = ent
		}
		if more == false {
			break
		}
	}
}
*/

func (server *DHCPServer) constructVlanInfra() {
	curMark := 0
	server.logger.Debug("Calling Asicd for getting Vlan Property")
	count := 100
	for {
		bulkVlanInfo, _ := server.asicdClient.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
		if bulkVlanInfo == nil {
			break
		}
		/* Get bulk on vlan state can re-use curMark and count used by get bulk vlan, as there is a 1:1 mapping in terms of cfg/state objs */
		bulkVlanStateInfo, _ := server.asicdClient.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
		if bulkVlanStateInfo == nil {
			break
		}
		objCnt := int(bulkVlanInfo.Count)
		more := bool(bulkVlanInfo.More)
		curMark = int(bulkVlanInfo.EndIdx)
		for i := 0; i < objCnt; i++ {
			ifIndex := bulkVlanStateInfo.VlanStateList[i].IfIndex
			ent := server.vlanPropertyMap[ifIndex]
			untaggedIfIndexList := bulkVlanInfo.VlanList[i].UntagIfIndexList
			ent.UntagPortMap = make(map[int32]bool)
			for i := 0; i < len(untaggedIfIndexList); i++ {
				ent.UntagPortMap[untaggedIfIndexList[i]] = true
			}
			server.vlanPropertyMap[ifIndex] = ent
		}
		if more == false {
			break
		}
	}
}

func (server *DHCPServer) updateVlanInfra(msg asicdCommonDefs.VlanNotifyMsg, msgType uint8) {
	vlanId := int(msg.VlanId)
	ifIdx := asicdCommonDefs.GetIfIndexFromIntfIdAndIntfType(vlanId, commonDefs.IfTypeVlan)
	portList := msg.UntagPorts
	vlanEnt, _ := server.vlanPropertyMap[ifIdx]
	if msgType == asicdCommonDefs.NOTIFY_VLAN_CREATE ||
		msgType == asicdCommonDefs.NOTIFY_VLAN_UPDATE { // VLAN CREATE
		server.logger.Debug(fmt.Sprintln("Received Vlan Create or Update Notification Vlan:", vlanId, "PortList:", portList))
		vlanEnt.UntagPortMap = nil
		vlanEnt.UntagPortMap = make(map[int32]bool)
		for i := 0; i < len(portList); i++ {
			port := portList[i]
			vlanEnt.UntagPortMap[port] = true
		}
		server.vlanPropertyMap[ifIdx] = vlanEnt
	} else { // VLAN DELETE
		server.logger.Debug(fmt.Sprintln("Received Vlan Delete Notification Vlan:", vlanId, "PortList:", portList))
		vlanEnt.UntagPortMap = nil
		delete(server.vlanPropertyMap, ifIdx)
	}
}

/*
func (server *DHCPServer) updateLagInfra(msg asicdCommonDefs.LagNotifyMsg, msgType uint8) {
        ifIdx := int(msg.IfIndex)
        portList := msg.IfIndexList
        //server.logger.Info(fmt.Sprintln("Lag Property Map:", server.lagPropMap))
        lagEnt, _ := server.lagPropMap[ifIdx]
        if msgType == asicdCommonDefs.NOTIFY_LAG_CREATE {
                server.logger.Info(fmt.Sprintln("Received Lag Create Notification IfIdx:", ifIdx, "PortList:", portList))
                lagEnt.PortMap = nil
                lagEnt.PortMap = make(map[int]bool)
                for i := 0; i < len(portList); i++ {
                        port := int(portList[i])
                        portEnt, _ := server.portPropMap[port]
                        portEnt.LagIfIdx = ifIdx
                        server.portPropMap[port] = portEnt
                        lagEnt.PortMap[port] = true
                }
                server.lagPropMap[ifIdx] = lagEnt
        } else if msgType == asicdCommonDefs.NOTIFY_LAG_UPDATE {
                newPortMap := make(map[int]bool)
                for i := 0; i < len(portList); i++ {
                        port := int(portList[i])
                        newPortMap[port] = true
                }
                for oldPort, _ := range lagEnt.PortMap {
                        _, exist := newPortMap[oldPort]
                        if !exist { // There in Old but Not in New so flush arp cache
				//TODO: Update port Property map
                                portEnt, _ := server.portPropMap[oldPort]
                                portEnt.LagIfIdx = -1
                                server.portPropMap[oldPort] = portEnt
                        } else { //Intersecting Ports (already there in PortMap)
                                delete(newPortMap, oldPort)
                        }
                }
                for newPort, _ := range newPortMap { // All new ports need to be added
                        portEnt, _ := server.portPropMap[newPort]
                        portEnt.LagIfIdx = ifIdx
                        server.portPropMap[newPort] = portEnt
                        lagEnt.PortMap[newPort] = true
                }
                server.lagPropMap[ifIdx] = lagEnt
        } else {
               server.logger.Info(fmt.Sprintln("Received Lag Delete Notification IfIdx:", ifIdx, "PortList:", portList))
                for i := 0; i < len(portList); i++ {
                        port := int(portList[i])
                        portEnt, _ := server.portPropMap[port]
                        portEnt.LagIfIdx = -1
                        server.portPropMap[port] = portEnt
                }
		//TODO: Update port Property map
                lagEnt.PortMap = nil
                delete(server.lagPropMap, ifIdx)
        }
        //server.logger.Info(fmt.Sprintln("Lag Property Map:", server.lagPropMap))
}
*/
