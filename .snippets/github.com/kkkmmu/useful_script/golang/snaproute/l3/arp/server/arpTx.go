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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"math/rand"
	"net"
	"time"
)

func getIP(ipAddr string) (ip net.IP) {
	ip = net.ParseIP(ipAddr)
	if ip == nil {
		return ip
	}
	ip = ip.To4()
	return ip
}

func getHWAddr(macAddr string) (mac net.HardwareAddr) {
	mac, err := net.ParseMAC(macAddr)
	if mac == nil || err != nil {
		return nil
	}

	return mac
}

/*
 *@fn sendArpReq
 *  Send the ARP request for ip targetIP
 */
func (server *ARPServer) sendArpReq(targetIp string, l3IfIdx int) {
	server.logger.Debug("sendArpReq(): sending arp requeust for targetIp ", targetIp, "to port:", l3IfIdx)

	l3Ent, exist := server.l3IntfPropMap[l3IfIdx]
	if !exist {
		server.logger.Err("SendArpReq(): L3 interface does not exist", l3IfIdx)
		return
	}
	macAddr := server.GetMacAddr(l3IfIdx)
	pcapHdl, err := pcap.OpenLive(l3Ent.IfName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err("Unable to open pcap handle on:", l3Ent.IfName, "error:", err)
		return
	}
	defer pcapHdl.Close()
	srcIpAddr := getIP(l3Ent.IpAddr)
	if srcIpAddr == nil {
		server.logger.Err(fmt.Sprintf("Corrupted source ip :  ", l3Ent.IpAddr))
		return
	}

	destIpAddr := getIP(targetIp)
	if destIpAddr == nil {
		server.logger.Err(fmt.Sprintf("Corrupted destination ip :  ", targetIp))
		return
	}

	myMacAddr := getHWAddr(macAddr)
	if myMacAddr == nil {
		server.logger.Err(fmt.Sprintf("corrupted my mac : ", macAddr))
		return
	}
	arp_layer := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   myMacAddr,
		SourceProtAddress: srcIpAddr,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	eth_layer := layers.Ethernet{
		SrcMAC:       myMacAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	arp_layer.DstProtAddress = destIpAddr
	gopacket.SerializeLayers(buffer, options, &eth_layer, &arp_layer)

	//logger.Println("Buffer : ", buffer)
	// send arp request and retry after timeout if arp cache is not updated
	if err := pcapHdl.WritePacketData(buffer.Bytes()); err != nil {
		server.logger.Err(fmt.Sprintln("Error writing data to packet buffer for l3Intf:", l3IfIdx))
		return
	}
	return
}

/*
 *@fn sendArpProbe
 *  Send the ARP Probe for ip localIP
 */
func (server *ARPServer) sendArpProbe(l3IfIdx int) {
	l3Ent, exist := server.l3IntfPropMap[l3IfIdx]
	if !exist {
		server.logger.Err("SendArpProbe(): L3 interface does not exist", l3IfIdx)
		return
	}
	macAddr := server.GetMacAddr(l3IfIdx)
	localIp := l3Ent.IpAddr
	server.logger.Debug("sendArpReq(): sending arp requeust for localIp ", localIp, "to port:", l3Ent.IfName)

	if l3Ent.OperState == false {
		return
	}

	pcapHdl, err := pcap.OpenLive(l3Ent.IfName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err("Unable to open pcap handle on:", l3Ent.IfName, "error:", err)
		return
	}
	defer pcapHdl.Close()

	srcIpAddr := getIP("0.0.0.0")
	if srcIpAddr == nil {
		return
	}

	destIpAddr := getIP(localIp)
	if destIpAddr == nil {
		server.logger.Err("Corrupted destination ip :  ", localIp)
		return
	}

	myMacAddr := getHWAddr(macAddr)
	if myMacAddr == nil {
		server.logger.Err(fmt.Sprintf("corrupted my mac : ", macAddr))
		return
	}
	arp_layer := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   myMacAddr,
		SourceProtAddress: srcIpAddr,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	eth_layer := layers.Ethernet{
		SrcMAC:       myMacAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	arp_layer.DstProtAddress = destIpAddr
	gopacket.SerializeLayers(buffer, options, &eth_layer, &arp_layer)

	// send arp request and retry after timeout if arp cache is not updated
	if err := pcapHdl.WritePacketData(buffer.Bytes()); err != nil {
		server.logger.Err("Error writing data to packet buffer for port:", l3Ent.IfName)
		return
	}
	return
}

func (server *ARPServer) SendArpProbe(l3IfIdx int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)
	wait := r1.Intn(server.probeWait)
	time.Sleep(time.Duration(wait) * time.Second)
	for i := 0; i < server.probeNum; i++ {
		server.sendArpProbe(l3IfIdx)
		diff := r2.Intn(server.probeMax - server.probeMin)
		diff = diff + server.probeMin
		time.Sleep(time.Duration(diff) * time.Second)
	}
	return
}

func (server *ARPServer) SendGarp(ifName, macAddr, ipAddr string) {
	server.logger.Info("Sending garp for:", ifName, macAddr, ipAddr)
	// src mac
	myMacAddr := getHWAddr(macAddr)
	if myMacAddr == nil {
		server.logger.Err("corrupted my mac : ", macAddr)
		return
	}
	// src ip
	srcIpAddr := getIP(ipAddr)
	if srcIpAddr == nil {
		server.logger.Err("Corrupted source ip :  ", ipAddr)
		return
	}
	// broadcast mac
	bMac := getHWAddr("ff:ff:ff:ff:ff:ff")
	if myMacAddr == nil {
		server.logger.Err("corrupted my mac : ff:ff:ff:ff:ff:ff")
		return
	}

	pcapHdl, err := pcap.OpenLive(ifName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err("Unable to open pcap handle on:", ifName, "error:", err)
		return
	}
	defer pcapHdl.Close()

	arp_layer := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   myMacAddr,
		SourceProtAddress: srcIpAddr,
		DstProtAddress:    srcIpAddr,
		DstHwAddress:      bMac,
	}

	eth_layer := layers.Ethernet{
		SrcMAC:       myMacAddr,
		DstMAC:       bMac,
		EthernetType: layers.EthernetTypeARP,
	}
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buffer, options, &eth_layer, &arp_layer)

	// send arp request and retry after timeout if arp cache is not updated
	if err := pcapHdl.WritePacketData(buffer.Bytes()); err != nil {
		server.logger.Err("Error writing data to packet buffer for port:", err)
		return
	}
	server.logger.Info("GARP send out successfully for:", ifName, macAddr, ipAddr)
}
