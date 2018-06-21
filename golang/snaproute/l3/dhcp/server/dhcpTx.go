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
	"fmt"
	"github.com/google/gopacket/pcap"
	"net"
)

func (server *DHCPServer) sendDhcpAck(port int32, bootPMsgData *BootPMsgStruct, data []byte, ipAddr uint32) {
	server.logger.Info("Sending Dhcp Ack  msg")
	clientMac := (net.HardwareAddr(bootPMsgData.ClientHWAddr)).String()
	portEnt, _ := server.portPropertyMap[port]
	l3Ent, _ := server.l3IntfPropMap[portEnt.L3IfIndex]
	dhcpIntfKey := l3Ent.DhcpIfKey

	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	dhcpAck := make([]byte, BOOTP_MSG_SIZE+36)
	copy(dhcpAck, data[0:BOOTP_MSG_SIZE])
	dhcpIntfEnt.dhcpMsg[6] = DHCPACK
	copy(dhcpAck[BOOTP_MSG_SIZE:], dhcpIntfEnt.dhcpMsg[0:])
	binary.BigEndian.PutUint32(dhcpAck[16:20], ipAddr)
	dhcpAckPkt := server.buildDhcpAckPkt(portEnt, dhcpAck, bootPMsgData)

	//server.logger.Info(fmt.Sprintln("====DHCP Ack=====", dhcpAckPkt))
	if dhcpAckPkt == nil {
		return
	}
	//dhcpOfferPkt := server.constructDhcpOffer(port, pktMd, bootPMsgData, data)
	pcapHdl, err := pcap.OpenLive(portEnt.IfName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err(fmt.Sprintln("Unable to open pcap handle on:", portEnt.IfName, "error:", err))
		return
	}
	defer pcapHdl.Close()
	if err := pcapHdl.WritePacketData(dhcpAckPkt); err != nil {
		server.logger.Err(fmt.Sprintln("Error writing data to packet buffer for port:", port))
		return
	}
	if bootPMsgData.ClientIPAddr == 0 {
		server.logger.Info("Starting Lease Entry Handler")
		go server.StartLeaseEntryHandler(port, ipAddr, clientMac)
	}
}

func (server *DHCPServer) sendDhcpOffer(port int32, pktMd *PktMetadata, bootPMsgData *BootPMsgStruct, data []byte) {
	server.logger.Info("Handle Dhcp Discover msg")
	clientMac := (net.HardwareAddr(bootPMsgData.ClientHWAddr)).String()
	portEnt, _ := server.portPropertyMap[port]
	l3Ent, _ := server.l3IntfPropMap[portEnt.L3IfIndex]
	dhcpIntfKey := l3Ent.DhcpIfKey
	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	dhcpOffer := make([]byte, BOOTP_MSG_SIZE+36)
	copy(dhcpOffer, data[0:BOOTP_MSG_SIZE])
	//Set DHCP Offer Message
	//Set yiaddr
	dhcpIntfEnt.dhcpMsg[6] = DHCPOFFER
	copy(dhcpOffer[BOOTP_MSG_SIZE:], dhcpIntfEnt.dhcpMsg[0:])
	ipAddr, exist := dhcpIntfEnt.usedIpToMac[clientMac]
	if !exist {
		ip, ret := server.findUnusedIP(dhcpIntfEnt)
		if ret == false {
			server.logger.Err("No available IP Addr")
			return
		}
		dhcpIntfEnt.usedIpToMac[clientMac] = ip
		ipAddr = ip
		uIPEnt, _ := dhcpIntfEnt.usedIpPool[ip]
		uIPEnt.LeaseTime = server.DhcpGlobalConf.DefaultLeaseTime
		uIPEnt.MacAddr = clientMac
		uIPEnt.TransactionId = bootPMsgData.TransactionId
		dhcpIntfEnt.usedIpPool[ip] = uIPEnt
		server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
		server.logger.Info("Starting Stale Entry Handler")
		go server.StartStaleEntryHandler(port, ipAddr, clientMac)
	} else {
		server.logger.Info(fmt.Sprintln("Already offered and IP Address to this client", clientMac))
	}

	binary.BigEndian.PutUint32(dhcpOffer[16:20], ipAddr)
	dhcpOfferPkt := server.buildDhcpOfferPkt(portEnt, dhcpOffer)

	if dhcpOfferPkt == nil {
		return
	}
	//dhcpOfferPkt := server.constructDhcpOffer(port, pktMd, bootPMsgData, data)
	pcapHdl, err := pcap.OpenLive(portEnt.IfName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err(fmt.Sprintln("Unable to open pcap handle on:", portEnt.IfName, "error:", err))
		return
	}
	defer pcapHdl.Close()
	if err := pcapHdl.WritePacketData(dhcpOfferPkt); err != nil {
		server.logger.Err(fmt.Sprintln("Error writing data to packet buffer for port:", port))
		return
	}
	return

}
