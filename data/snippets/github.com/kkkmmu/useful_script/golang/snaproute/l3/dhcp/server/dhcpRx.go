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
	"math/rand"
	"net"
	"time"
)

const (
	DHCPServerPort        uint16 = 67
	DHCPClientPort        uint16 = 68
	UDPProto              uint8  = 17
	BootPReply            uint8  = 2
	BootPRequest          uint8  = 1
	UDP_HDR_LEN           uint16 = 8
	CLIENT_HW_ADDR_LEN    int    = 16
	MAX_LEN_SERVER_NAME   int    = 64
	MAX_LEN_BOOT_FILENAME int    = 128
	DHCPDISCOVER          uint8  = 1
	DHCPOFFER             uint8  = 2
	DHCPREQUEST           uint8  = 3
	DHCPACK               uint8  = 5
	DHCPRELEASE           uint8  = 7
	BOOTP_MSG_SIZE        uint16 = 236
)

const (
	OFFERED uint8 = 1
)

func (server *DHCPServer) StartRxDhcpPkt(port int32) {
	portEnt, _ := server.portPropertyMap[port]
	filter := fmt.Sprintf(`not ether src %s and dst port %d and udp`, portEnt.MacAddr, int(DHCPServerPort))
	server.logger.Debug(fmt.Sprintln("Port:", port, "Filter:", filter))
	pcapHdl, err := pcap.OpenLive(portEnt.IfName, server.snapshotLen, server.promiscuous, server.pcapTimeout)
	if pcapHdl == nil {
		server.logger.Err(fmt.Sprintln("Unable to open pcap handler on:", portEnt.IfName, "err:", err))
		return
	} else {
		err := pcapHdl.SetBPFFilter(filter)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to set filter on port:", port))
		}
	}

	portEnt.PcapHdl = pcapHdl
	server.portPropertyMap[port] = portEnt
	go server.processRxPkts(port)
	server.logger.Debug(fmt.Sprintln("Starting Rx on port:", port))
	return
}

func (server *DHCPServer) deinitProcessRxPkt(port int32) {
	portEnt, _ := server.portPropertyMap[port]
	portEnt.PcapHdl.Close()
	portEnt.PcapHdl = nil
	server.portPropertyMap[port] = portEnt
	portEnt.CtrlCh <- true
}

func (server *DHCPServer) processRxPkts(port int32) {
	portEnt, _ := server.portPropertyMap[port]
	recv := gopacket.NewPacketSource(portEnt.PcapHdl, layers.LayerTypeEthernet)
	in := recv.Packets()
	for {
		select {
		case packet, ok := <-in:
			if ok {
				go server.ProcessRecvPkts(port, packet)
			}
		case <-portEnt.CtrlCh:
			server.deinitProcessRxPkt(port)
			return
		}
	}
	return
}

func (server *DHCPServer) ProcessRecvPkts(port int32, pkt gopacket.Packet) {
	//server.logger.Info(fmt.Sprintln("Port:", port, "packet:", pkt))
	ethLayer := pkt.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		server.logger.Err("Not an ethernet frame")
		return
	}

	eth := ethLayer.(*layers.Ethernet)
	ethHdrMd := NewEthHdrMetadata()
	ethHdrMd.srcMAC = eth.SrcMAC
	ethHdrMd.dstMAC = eth.DstMAC

	ipLayer := pkt.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		server.logger.Err("Not an IP Packet")
		return
	}

	portEnt, exist := server.portPropertyMap[port]
	if !exist {
		server.logger.Err(fmt.Sprintln("Port:", port, "doesnot exist anymore"))
		return
	}
	if portEnt.L3IfIndex == -1 {
		server.logger.Err("Port is not member of any L3 interface")
		return
	}

	ipHdrMd := NewIPHdrMetadata()
	err := server.processIPv4Layer(ipLayer, ipHdrMd)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Dropped because of IPv4 Processing", err))
		return
	}

	if ipHdrMd.Protocol != UDPProto {
		server.logger.Err(fmt.Sprintln("Not an UDP Protocol"))
		return
	}

	udpLayer := pkt.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		server.logger.Err("Not an UDP Packet")
		return
	}

	udpHdrMd := NewUDPHdrMetadata()
	err = server.processUDPLayer(udpLayer, udpHdrMd, ipLayer.LayerPayload(), ipLayer)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Dropped because of UDP Processing", err))
		return
	}

	if udpHdrMd.DstPort != DHCPServerPort ||
		udpHdrMd.SrcPort != DHCPClientPort {
		server.logger.Err(fmt.Sprintln("Packet recvd from invalid src port", udpHdrMd.SrcPort, "dst port", udpHdrMd.DstPort))
		return
	}

	bootPData := udpLayer.LayerPayload()
	//server.logger.Info(fmt.Sprintln("DHCP Data:", bootPData))
	/*
		dhcpMsgData := NewDhcpMsgStruct()
	*/
	pktMd := NewPktMetadata()
	pktMd.ethHdrMd = ethHdrMd
	pktMd.ipHdrMd = ipHdrMd
	pktMd.udpHdrMd = udpHdrMd
	if len(bootPData) != int(udpHdrMd.Length-UDP_HDR_LEN) {
		server.logger.Err(fmt.Sprintln("Invalid data Length"))
		return
	}
	err = server.processBootPMsg(port, bootPData, pktMd)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Error in processing DHCP Message", err))
		return
	}
	return
}

func (server *DHCPServer) processBootPMsg(port int32, data []byte, pktMd *PktMetadata) error {
	//server.logger.Debug(fmt.Sprintln("Pkt Metadata:", pktMd))
	//server.logger.Debug(fmt.Sprintln("data:", data))

	msgType := uint8(data[0])
	switch msgType {
	case BootPReply:
		err := server.processBootPReply(port, data, pktMd)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Error Handling Boot Reply", err))
			return err
		}
	case BootPRequest:
		err := server.processBootPRequest(port, data, pktMd)
		if err != nil {
			server.logger.Err(fmt.Sprintln("Error handling Boot Request", err))
		}
	default:
		server.logger.Err(fmt.Sprintln("Invalid Msg Type"))
		err := errors.New("Invalid Msg Type")
		return err
	}

	return nil
}

func (server *DHCPServer) processBootPReply(port int32, data []byte, pktMd *PktMetadata) error {
	server.logger.Debug(fmt.Sprintln("Port:", port, "processBootReply()"))
	/*
		dhcpMsgData := NewDhcpMsgStruct()
		err := decodeDhcpMsg(data, dhcpMsgStruct)
		if err != nil {
			server.logger.Err(fmt.Sprintln("There was decoding error for DHCP Boot Reply", err))
		}
	*/
	// TODO: If we are seeing this msg which means there is another DHCP Server
	return nil
}

func (server *DHCPServer) processBootPRequest(port int32, data []byte, pktMd *PktMetadata) error {
	server.logger.Debug(fmt.Sprintln("Port:", port, "processBootRequest()"))
	bootPMsgData := NewBootPMsgStruct()
	err := decodeBootPMsg(data, bootPMsgData)
	if err != nil {
		server.logger.Err(fmt.Sprintln("There was decoding error for DHCP Boot Request", err))
		return err
	}
	//server.logger.Info(fmt.Sprintln("OptionFields:", bootPMsgData.DhcpOptionMap))
	if bootPMsgData.MagicCookie != 0x63825363 {
		server.logger.Err(fmt.Sprintln("Invaild magic cookie:", bootPMsgData.MagicCookie))
		err := errors.New(fmt.Sprintln("Invaild magic cookie:", bootPMsgData.MagicCookie))
		return err
	}
	if bootPMsgData.DhcpOptionMap != nil {
		ent, exist := bootPMsgData.DhcpOptionMap[DhcpMsgTypeOptCode]
		if !exist {
			server.logger.Err("Dhcp Msg Type does not exist")
			err := errors.New("Dhcp Msg Type does not exist")
			return err
		}

		if ent.Length != 1 {
			server.logger.Err("Invalid length for DHCP Message Type")
		}
		msgType := uint8(ent.Data[0])
		switch msgType {
		case DHCPDISCOVER:
			server.sendDhcpOffer(port, pktMd, bootPMsgData, data)
		case DHCPREQUEST:
			server.processDhcpRequest(port, pktMd, bootPMsgData, data)
		//case DHCPDECLINE:
		case DHCPRELEASE:
			server.processDhcpRelease(port, pktMd, bootPMsgData)
		//case DHCPINFORM:
		default:
			server.logger.Err("Dhcp Server will not handle other DHCP message apart from DHCPDISCOVER, DHCPREQUEST")
		}
	}
	return nil
}

func (server *DHCPServer) processDhcpRelease(port int32, pktMd *PktMetadata, bootPMsgData *BootPMsgStruct) {
	server.logger.Info("Handle Dhcp Release msg")
	clientMac := (net.HardwareAddr(bootPMsgData.ClientHWAddr)).String()
	portEnt, _ := server.portPropertyMap[port]
	l3IfIdx := portEnt.L3IfIndex
	l3Ent := server.l3IntfPropMap[l3IfIdx]
	dhcpIntfKey := l3Ent.DhcpIfKey
	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	ipAddr, exist := dhcpIntfEnt.usedIpToMac[clientMac]
	if !exist {
		server.logger.Info("This request is not for our DHCP Offer")
		return
	}

	uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]

	if uIPEnt.RefreshTimer != nil {
		uIPEnt.RefreshTimer.Stop()
	}
	delete(dhcpIntfEnt.usedIpPool, ipAddr)
	delete(dhcpIntfEnt.usedIpToMac, clientMac)
	server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt
}

func (server *DHCPServer) processDhcpRequest(port int32, pktMd *PktMetadata, bootPMsgData *BootPMsgStruct, data []byte) {
	server.logger.Info("Handle Dhcp Request msg")
	clientMac := (net.HardwareAddr(bootPMsgData.ClientHWAddr)).String()
	//server.logger.Info(fmt.Sprintln("1 Handle Dhcp Request msg", clientMac, "data:", data, "bootPMsgData", bootPMsgData.ClientHWAddr))
	portEnt, _ := server.portPropertyMap[port]
	l3IfIdx := portEnt.L3IfIndex
	l3Ent := server.l3IntfPropMap[l3IfIdx]
	dhcpIntfKey := l3Ent.DhcpIfKey
	/*
		dhcpIntfKey := DhcpIntfKey{
			subnet:     portEnt.IpAddr & portEnt.Mask,
			subnetMask: portEnt.Mask,
		}
	*/
	dhcpIntfEnt, _ := server.DhcpIntfConfMap[dhcpIntfKey]
	ipAddr, exist := dhcpIntfEnt.usedIpToMac[clientMac]
	if !exist {
		server.logger.Info("This request is not for our DHCP Offer")
		return
	}

	uIPEnt, _ := dhcpIntfEnt.usedIpPool[ipAddr]
	// Check server ID
	// Check Transaction ID
	// Check Requested IP Address
	// Check Lease time
	// Check Subnet Address
	// Check RtrAddr
	// Check for Client ID TODO
	serverId, exist := bootPMsgData.DhcpOptionMap[ServerIdOptCode]
	if !exist {
		server.logger.Info(fmt.Sprintln("Server Id doesnot exist"))
	} else {
		sId := convertIPv4ToUint32(serverId.Data)
		if sId != portEnt.IpAddr {
			server.logger.Info(fmt.Sprintln("Server Id ", sId, "is not equal to ", portEnt.IpAddr))
			return
		}
	}

	if bootPMsgData.TransactionId != uIPEnt.TransactionId {
		// DHCP NACK
		server.logger.Info(fmt.Sprintln("TransactionId are not equal"))
		return
	}

	reqIPAddr, exist := bootPMsgData.DhcpOptionMap[ReqIPAddrOptCode]
	if !exist {
		if bootPMsgData.ClientIPAddr != ipAddr {
			// Send DHCP NACK
			server.logger.Info(fmt.Sprintln("no reqIpAddr and client ip Addr"))
			return
		}
	} else {
		if reqIPAddr.Length != 4 {
			server.logger.Info("Hello1")
			//Send DHCP NACK
		}
		reqIP := convertIPv4ToUint32(reqIPAddr.Data)
		if reqIP != ipAddr {
			server.logger.Info("Hello2")
			// Send DHCP NACK
		}
	}
	leaseTime, exist := bootPMsgData.DhcpOptionMap[IPAddrLeaseTimeOptCode]
	if exist {
		if leaseTime.Length != 4 {
			server.logger.Info("Hello3")
			//Send DHCP NACK
		}
		lTime := convertIPv4ToUint32(leaseTime.Data)
		if lTime != uIPEnt.LeaseTime {
			server.logger.Info("Hello4")
			// Send DHCP NACK
		}
	}
	subnetAddr, exist := bootPMsgData.DhcpOptionMap[SubnetMaskOptCode]
	if exist {
		if subnetAddr.Length != 4 {
			server.logger.Info("Hello5")
			//Send DHCP NACK
		}
		sAddr := convertIPv4ToUint32(subnetAddr.Data)
		if sAddr != dhcpIntfKey.subnetMask {
			server.logger.Info("Hello6")
			// Send DHCP NACK
		}
	}
	rtrAddr, exist := bootPMsgData.DhcpOptionMap[RouterOptCode]
	if exist {
		if rtrAddr.Length != 4 {
			server.logger.Info("Hello7")
			//Send DHCP NACK
		}
		rtrAddr := convertIPv4ToUint32(rtrAddr.Data)
		if rtrAddr != dhcpIntfEnt.rtrAddr {
			// Send DHCP NACK
		}
	}

	server.sendDhcpAck(port, bootPMsgData, data, ipAddr)
	dhcpIntfEnt, _ = server.DhcpIntfConfMap[dhcpIntfKey]
	ipAddr, exist = dhcpIntfEnt.usedIpToMac[clientMac]
	if !exist {
		server.logger.Info("This request is not for our DHCP Offer")
		return
	}

	uIPEnt, _ = dhcpIntfEnt.usedIpPool[ipAddr]
	//server.logger.Info(fmt.Sprintln("1 uIPEnt: ", uIPEnt))
	if uIPEnt.State == OFFERED {
		server.logger.Info(fmt.Sprintln("Reseting refresh timer...."))
		uIPEnt.RefreshTimer.Reset(time.Duration(uIPEnt.LeaseTime) * time.Second)
	}
	if uIPEnt.StaleTimer != nil {
		server.logger.Info(fmt.Sprintln("Stopping stale timer...."))
		uIPEnt.StaleTimer.Stop()
		uIPEnt.StaleTimer = nil
	}
	dhcpIntfEnt.usedIpPool[ipAddr] = uIPEnt
	server.DhcpIntfConfMap[dhcpIntfKey] = dhcpIntfEnt

}

/*
func (server *DHCPServer) sendDhcpAck(port int32, bootPMsgData *BootPMsgStruct, data []byte, ipAddr uint32) {
	server.logger.Info("Sending Dhcp Ack  msg")
	clientMac := (net.HardwareAddr(bootPMsgData.ClientHWAddr)).String()
	portEnt, _ := server.portPropertyMap[port]
	dhcpIntfKey := DhcpIntfKey{
		subnet:     portEnt.IpAddr & portEnt.Mask,
		subnetMask: portEnt.Mask,
	}
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
	dhcpIntfKey := DhcpIntfKey{
		subnet:     portEnt.IpAddr & portEnt.Mask,
		subnetMask: portEnt.Mask,
	}
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

func (server *DHCPServer) buildDhcpAckPkt(portEnt PortProperty, dhcpAck []byte, bootPMsgData *BootPMsgStruct) []byte {
	udpLen := len(dhcpAck) + 8
	udpHdr := layers.UDP{
		SrcPort: layers.UDPPort(67),
		DstPort: layers.UDPPort(68),
		Length:  uint16(udpLen),
	}

	ipPktLen := 20 + udpLen
	srcIP := convertUint32ToNetIPv4(portEnt.IpAddr)
	var dstIp net.IP
	var dstMac net.HardwareAddr
	if bootPMsgData.ClientIPAddr != 0 {
		dstIp = convertUint32ToNetIPv4(bootPMsgData.ClientIPAddr)
		cHWAddr := bootPMsgData.ClientHWAddr
		//server.logger.Info(fmt.Sprintln("************clientHWaddr:*******", bootPMsgData.ClientHWAddr, cHWAddr))
		dstMac = net.HardwareAddr{cHWAddr[0], cHWAddr[1], cHWAddr[2], cHWAddr[3], cHWAddr[4], cHWAddr[5]}
	} else {
		dstIp = convertUint32ToNetIPv4(0xffffffff)
		dstMac = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	}
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(20),
		TOS:      uint8(0),
		TTL:      uint8(255),
		Length:   uint16(ipPktLen),
		Protocol: layers.IPProtocol(17),
		SrcIP:    srcIP,
		DstIP:    dstIp,
	}

	srcMac := getHWAddr(portEnt.MacAddr)
	if srcMac == nil {
		server.logger.Err(fmt.Sprintln("Corrupt MAC Addr:", portEnt.MacAddr))
		return nil
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       srcMac,
		DstMAC:       dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	udpHdr.SetNetworkLayerForChecksum(&ipLayer)
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, &udpHdr, gopacket.Payload(dhcpAck))
	dhcpAckPkt := buffer.Bytes()
	return dhcpAckPkt
}

func (server *DHCPServer) buildDhcpOfferPkt(portEnt PortProperty, dhcpOffer []byte) []byte {
	udpLen := len(dhcpOffer) + 8
	udpHdr := layers.UDP{
		SrcPort: layers.UDPPort(67),
		DstPort: layers.UDPPort(68),
		Length:  uint16(udpLen),
	}

	ipPktLen := 20 + udpLen
	srcIP := convertUint32ToNetIPv4(portEnt.IpAddr)
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(20),
		TOS:      uint8(0),
		TTL:      uint8(255),
		Length:   uint16(ipPktLen),
		Protocol: layers.IPProtocol(17),
		SrcIP:    srcIP,
		DstIP:    net.IP{255, 255, 255, 255},
	}

	srcMac := getHWAddr(portEnt.MacAddr)
	if srcMac == nil {
		server.logger.Err(fmt.Sprintln("Corrupt MAC Addr:", portEnt.MacAddr))
		return nil
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       srcMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	udpHdr.SetNetworkLayerForChecksum(&ipLayer)
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, &udpHdr, gopacket.Payload(dhcpOffer))
	dhcpOfferPkt := buffer.Bytes()
	return dhcpOfferPkt
}

*/

func (server *DHCPServer) findUnusedIP(dhcpIntfData DhcpIntfData) (uint32, bool) {
	diff := int(dhcpIntfData.higherIPBound - dhcpIntfData.lowerIPBound + 1)
	if len(dhcpIntfData.usedIpPool) == diff {
		return 0, false
	}
	for {
		offset := rand.Intn(diff)
		_, exist := dhcpIntfData.usedIpPool[dhcpIntfData.lowerIPBound+uint32(offset)]
		if !exist {
			return dhcpIntfData.lowerIPBound + uint32(offset), true
		}
	}
}
