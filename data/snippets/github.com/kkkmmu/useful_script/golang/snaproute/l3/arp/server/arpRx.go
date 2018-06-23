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
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
)

const (
	OR_ETHER_SRC = " or ether src "
	CLOSE_FILTER = "))"
)

func (server *ARPServer) StartArpRxTx(ifName string, macAddr string, filter string) (*pcap.Handle, error) {
	filter = filter + CLOSE_FILTER
	server.logger.Debug("Port: ", ifName, "Pcap filter:", filter)
	pcapHdl, err := pcap.OpenLive(ifName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		return nil, errors.New(fmt.Sprintln("Unable to open pcap handler on", ifName, "error:", err))
	} else {
		err := pcapHdl.SetBPFFilter(filter)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Unable to set bpf filter to pcap handler on", ifName, "error:", err))
		}
	}

	return pcapHdl, nil
}

func (server *ARPServer) processRxPkts(portIfIdx int) {
	portEnt, _ := server.portPropMap[portIfIdx]
	src := gopacket.NewPacketSource(portEnt.PcapHdl, portEnt.PcapHdl.LinkType())
	in := src.Packets()
	for {
		select {
		case packet, ok := <-in:
			if ok {
				ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
				if ethernetLayer == nil {
					continue
				}
				arpLayer := packet.Layer(layers.LayerTypeARP)
				ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
				if arpLayer == nil && ipv4Layer == nil {
					continue
				}
				dot1QLayer := packet.Layer(layers.LayerTypeDot1Q)
				if dot1QLayer != nil {
					dot1Q := dot1QLayer.(*layers.Dot1Q)
					server.processTagPacket(ethernetLayer, arpLayer, ipv4Layer, portIfIdx, int(dot1Q.VLANIdentifier))
				} else {
					server.processUntagPacket(ethernetLayer, arpLayer, ipv4Layer, portIfIdx)
				}
			}
		case <-portEnt.CtrlCh:
			server.logger.Info("Recevd shutdown for:", portIfIdx, portEnt.IfName)
			portEnt.CtrlReplyCh <- true
			return
		}
	}
	return
}

func (server *ARPServer) processTagPacket(ethernetLayer, arpLayer, ipv4Layer gopacket.Layer, portIfIdx int, vlanId int) {
	if arpLayer != nil {
		server.processArpPkt(arpLayer, portIfIdx, vlanId)
	} else if ipv4Layer != nil {
		server.processIpPkt(ethernetLayer, ipv4Layer, portIfIdx, vlanId)
	}

}

func (server *ARPServer) processUntagPacket(ethernetLayer, arpLayer, ipv4Layer gopacket.Layer, portIfIdx int) {
	if arpLayer != nil {
		server.processArpPkt(arpLayer, portIfIdx, -1)
	} else if ipv4Layer != nil {
		server.processIpPkt(ethernetLayer, ipv4Layer, portIfIdx, -1)
	}
}

func (server *ARPServer) processArpPkt(arpLayer gopacket.Layer, portIfIdx, vlanId int) {
	arp := arpLayer.(*layers.ARP)
	if arp == nil {
		return
	}
	portEnt, _ := server.portPropMap[portIfIdx]
	if portEnt.MacAddr == (net.HardwareAddr(arp.SourceHwAddress)).String() {
		return
	}

	if arp.Operation == layers.ARPReply {
		server.processArpReply(arp, portIfIdx, vlanId)
	} else if arp.Operation == layers.ARPRequest {
		server.processArpRequest(arp, portIfIdx, vlanId)
	}
}

func (server *ARPServer) getMyIPAndNetAddr(portIfIdx, vlanId int) (net.IP, net.IP, net.IPMask, error, int, int, int) {
	portEnt, _ := server.portPropMap[portIfIdx]
	if vlanId == -1 {
		vlanId = portEnt.UntagVlanId
	}
	l3PortProp, exist := portEnt.L3PortPropMap[vlanId]
	if !exist {
		return nil, nil, nil, errors.New("L3 Intf doesnot exist"), -1, -1, -1
	}
	myIP := net.ParseIP(l3PortProp.IpAddr)
	myNet := myIP.Mask(l3PortProp.Netmask)
	return myIP, myNet, l3PortProp.Netmask, nil, l3PortProp.L3IfIdx, l3PortProp.LagIfIdx, vlanId
}

func (server *ARPServer) processArpReply(arp *layers.ARP, portIfIdx, vlanId int) {
	srcMac := (net.HardwareAddr(arp.SourceHwAddress)).String()
	srcIpAddr := (net.IP(arp.SourceProtAddress)).String()
	destIpAddr := (net.IP(arp.DstProtAddress)).String()
	if destIpAddr == "0.0.0.0" {
		return
	}

	_, myNet, myMask, err, l3IfIdx, lagIfIdx, vlanId := server.getMyIPAndNetAddr(portIfIdx, vlanId)
	if err != nil {
		return
	}
	srcIP := net.ParseIP(srcIpAddr)
	srcNet := srcIP.Mask(myMask)
	destIP := net.ParseIP(destIpAddr)
	destNet := destIP.Mask(myMask)
	if myNet.Equal(srcNet) == false ||
		myNet.Equal(destNet) == false {
		return
	}
	server.SendArpEntryUpdateMsg(portIfIdx, srcIpAddr, srcMac, false, vlanId, l3IfIdx, lagIfIdx)
}

func (server *ARPServer) SendArpEntryUpdateMsg(portIfIdx int, IpAddr string, macAddr string, Type bool, VlanId, L3IfIdx, LagIfIdx int) {
	server.logger.Debug("Sending Arp Entry Update Msg: portIfIdx:", portIfIdx, "IpAddr:", IpAddr, "macAddr:", macAddr, "Type:", Type, "VlanId:", VlanId, "L3IfIdx:", L3IfIdx, "LagIfIdx:", LagIfIdx)
	server.arpEntryUpdateCh <- UpdateArpEntryMsg{
		PortIfIdx: portIfIdx,
		IpAddr:    IpAddr,
		MacAddr:   macAddr,
		Type:      Type,
		VlanId:    VlanId,
		L3IfIdx:   L3IfIdx,
		LagIfIdx:  LagIfIdx,
	}

}

func (server *ARPServer) processArpRequest(arp *layers.ARP, portIfIdx, vlanId int) {
	srcMac := (net.HardwareAddr(arp.SourceHwAddress)).String()
	srcIpAddr := (net.IP(arp.SourceProtAddress)).String()
	destIpAddr := (net.IP(arp.DstProtAddress)).String()

	myIP, myNet, myMask, err, l3IfIdx, lagIfIdx, vlanId := server.getMyIPAndNetAddr(portIfIdx, vlanId)
	if err != nil {
		return
	}
	srcIP := net.ParseIP(srcIpAddr)
	srcNet := srcIP.Mask(myMask)
	destIP := net.ParseIP(destIpAddr)
	destNet := destIP.Mask(myMask)
	if srcIpAddr != "0.0.0.0" {
		if myNet.Equal(srcNet) == false ||
			myNet.Equal(destNet) == false {
			return
		}
	} else if myNet.Equal(destNet) == false {
		return
	}
	srcIpEqual := myIP.Equal(srcIP)
	destIPEqual := myIP.Equal(destIP)
	if srcIpEqual == true &&
		destIPEqual == true {
		// Our Own Gratuitous ARP
		return
	} else if destIPEqual == true {
		if srcIpAddr != "0.0.0.0" {
			server.SendArpEntryUpdateMsg(portIfIdx, srcIpAddr, srcMac, false, vlanId, l3IfIdx, lagIfIdx)
		}
	} else if srcIpEqual {
		return
	} else {
		if srcIpAddr == destIpAddr &&
			srcIpAddr == "0.0.0.0" {
			// Gratuitous Arp
			server.SendArpEntryUpdateMsg(portIfIdx, srcIpAddr, srcMac, false, vlanId, l3IfIdx, lagIfIdx)
		} else {
			if srcIpAddr != "0.0.0.0" {
				server.SendArpEntryUpdateMsg(portIfIdx, srcIpAddr, srcMac, false, vlanId, l3IfIdx, lagIfIdx)
			}
			server.updateL3ArpEntry(destIpAddr, false, l3IfIdx)
			server.sendArpReq(destIpAddr, l3IfIdx)
		}
	}
}

func (server *ARPServer) processIpPkt(ethernetLayer, ipv4Layer gopacket.Layer, portIfIdx, vlanId int) {
	ethernet := ethernetLayer.(*layers.Ethernet)
	if ethernet == nil {
		return
	}
	ipv4 := ipv4Layer.(*layers.IPv4)
	if ipv4 == nil {
		return
	}
	srcIP := ipv4.SrcIP
	srcIpAddr := srcIP.String()
	destIP := ipv4.DstIP
	destIpAddr := destIP.String()
	l3Intf := server.getL3IntfOnSameSubnet(srcIpAddr)
	if l3Intf != -1 {
		arpEnt, exist := server.arpCache[srcIpAddr]
		if !exist {
			server.updateL3ArpEntry(srcIpAddr, false, l3Intf)
		} else {
			if arpEnt.PortNum != portIfIdx ||
				arpEnt.MacAddr != (ethernet.SrcMAC).String() ||
				arpEnt.L3IfIdx != l3Intf {
				server.sendArpReq(srcIpAddr, l3Intf)
			}
		}
	}
	l3Intf = server.getL3IntfOnSameSubnet(destIpAddr)
	if l3Intf != -1 {
		arpEnt, exist := server.arpCache[destIpAddr]
		if !exist {
			server.updateL3ArpEntry(destIpAddr, false, l3Intf)
		} else {
			if arpEnt.PortNum != portIfIdx ||
				arpEnt.MacAddr != (ethernet.DstMAC).String() ||
				arpEnt.L3IfIdx != l3Intf {
				server.sendArpReq(destIpAddr, l3Intf)
			}
		}
	}
}

func (server *ARPServer) getL3IntfOnSameSubnet(ip string) int {
	ipAddr := net.ParseIP(ip)
	for l3Idx, l3Ent := range server.l3IntfPropMap {
		if l3Ent.IpAddr == ip {
			return -1
		}

		l3IpAddr := net.ParseIP(l3Ent.IpAddr)
		l3Net := l3IpAddr.Mask(l3Ent.Netmask)
		ipNet := ipAddr.Mask(l3Ent.Netmask)
		if l3Net.Equal(ipNet) {
			return l3Idx
		}
	}
	return -1
}

func (server *ARPServer) updateL3ArpEntry(TargetIP string, Type bool, l3IfIdx int) {
	server.arpEntryUpdateCh <- UpdateArpEntryMsg{
		IpAddr:    TargetIP,
		MacAddr:   "incomplete",
		Type:      Type,
		L3IfIdx:   l3IfIdx,
		PortIfIdx: -1,
		VlanId:    -1,
		LagIfIdx:  -1,
	}
}
