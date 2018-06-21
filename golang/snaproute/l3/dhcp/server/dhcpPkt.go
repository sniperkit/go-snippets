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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type EthHdrMetadata struct {
	srcMAC net.HardwareAddr
	dstMAC net.HardwareAddr
}

func NewEthHdrMetadata() *EthHdrMetadata {
	return &EthHdrMetadata{}
}

/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  |Type of Service|          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Identification        |Flags|      Fragment Offset    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Time to Live |    Protocol   |         Header Checksum       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Source Address                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Destination Address                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

type IPHdrMetadata struct {
	HdrLen   uint8
	Length   uint16
	Protocol uint8
	SrcIP    uint32
	DstIP    uint32
}

func NewIPHdrMetadata() *IPHdrMetadata {
	return &IPHdrMetadata{}
}

func (server *DHCPServer) processIPv4Layer(ipLayer gopacket.Layer, ipHdrMd *IPHdrMetadata) error {
	ipLayerContents := ipLayer.LayerContents()
	ipChkSum := binary.BigEndian.Uint16(ipLayerContents[10:12])
	binary.BigEndian.PutUint16(ipLayerContents[10:], 0)

	csum := computeChkSum(ipLayerContents)

	if csum != ipChkSum {
		err := errors.New("Incorrect IPv4 Checksum, hence discarding the packet")
		return err
	}

	ipPkt := ipLayer.(*layers.IPv4)
	ipHdrMd.Length = ipPkt.Length
	ipHdrMd.HdrLen = ipPkt.IHL
	ipHdrMd.Protocol = uint8(ipPkt.Protocol)
	ipHdrMd.SrcIP, _ = convertIPStrToUint32(ipPkt.SrcIP.String())
	ipHdrMd.DstIP, _ = convertIPStrToUint32(ipPkt.DstIP.String())

	return nil
}

type UDPHdrMetadata struct {
	SrcPort uint16
	DstPort uint16
	Length  uint16
}

func NewUDPHdrMetadata() *UDPHdrMetadata {
	return &UDPHdrMetadata{}
}

func computePseudoHdrChkSum(ipLayer gopacket.Layer) (csum uint32, err error) {
	ip := ipLayer.(*layers.IPv4)

	csum += (uint32(ip.SrcIP[0]) + uint32(ip.SrcIP[2])) << 8
	csum += uint32(ip.SrcIP[1]) + uint32(ip.SrcIP[3])
	csum += (uint32(ip.DstIP[0]) + uint32(ip.DstIP[2])) << 8
	csum += uint32(ip.DstIP[1]) + uint32(ip.DstIP[3])
	csum += uint32(ip.Protocol)
	return csum, err
}

func (server *DHCPServer) processUDPLayer(udpLayer gopacket.Layer, udpHdrMd *UDPHdrMetadata, data []byte, ipLayer gopacket.Layer) error {
	udpLayerContents := udpLayer.LayerContents()
	udpChkSum := binary.BigEndian.Uint16(udpLayerContents[6:8])
	binary.BigEndian.PutUint16(data[6:], 0)

	length := uint32(len(data))
	pseudoHdrChkSum, _ := computePseudoHdrChkSum(ipLayer)
	pseudoHdrChkSum += length & 0xffff
	pseudoHdrChkSum += length >> 16

	csum := computeTcpIPChkSum(data, pseudoHdrChkSum)
	if csum != udpChkSum {
		err := errors.New("Incorrect UDP Checksum, hence discarding the packet")
		return err
	}
	udpPkt := udpLayer.(*layers.UDP)
	udpHdrMd.SrcPort = uint16(udpPkt.SrcPort)
	udpHdrMd.DstPort = uint16(udpPkt.DstPort)
	udpHdrMd.Length = udpPkt.Length
	return nil
}

/*
   0                   1                   2                   3
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |     op (1)    |   htype (1)   |   hlen (1)    |   hops (1)    |
   +---------------+---------------+---------------+---------------+
   |                            xid (4)                            |
   +-------------------------------+-------------------------------+
   |           secs (2)            |           flags (2)           |
   +-------------------------------+-------------------------------+
   |                          ciaddr  (4)                          |
   +---------------------------------------------------------------+
   |                          yiaddr  (4)                          |
   +---------------------------------------------------------------+
   |                          siaddr  (4)                          |
   +---------------------------------------------------------------+
   |                          giaddr  (4)                          |
   +---------------------------------------------------------------+
   |                                                               |
   |                          chaddr  (16)                         |
   |                                                               |
   |                                                               |
   +---------------------------------------------------------------+
   |                                                               |
   |                          sname   (64)                         |
   +---------------------------------------------------------------+
   |                                                               |
   |                          file    (128)                        |
   +---------------------------------------------------------------+
   |                                                               |
   |                          options (variable)                   |
   +---------------------------------------------------------------+
*/

const (
	SubnetMaskOptCode         uint8 = 1 //ParamReqItemCode = 1
	RouterOptCode             uint8 = 3 //ParamReqItemCode = 3
	DNSOptCode                uint8 = 6 //ParamReqItemCode = 6
	HostNameOptCode           uint8 = 12
	DomainNameOptCode         uint8 = 15 //ParamReqItemCode = 15
	StaticRouteOptCode        uint8 = 33 //ParamReqItemCode = 33
	VendorSpecificInfoOptCode uint8 = 43 //ParamReqItemCode = 43
	NetBIOSServerOptCode      uint8 = 44 //ParamReqItemCode = 44
	ReqIPAddrOptCode          uint8 = 50
	IPAddrLeaseTimeOptCode    uint8 = 51
	DhcpMsgTypeOptCode        uint8 = 53
	ServerIdOptCode           uint8 = 54
	ParamReqListOptCode       uint8 = 55
	MaxDhcpMsgSizeOptCode     uint8 = 57
	RenewalTimeOptCode        uint8 = 58
	RebindingTimeOptCode      uint8 = 59
	ClientIdOptCode           uint8 = 61
	TFTPServerAddrOptCode     uint8 = 150 //ParamReqItemCode = 150
	EndOptCode                uint8 = 255
)

type OptionMetadata struct {
	OPCode uint8
	Length uint8
}

type SubnetMaskStruct struct {
	OptMd      OptionMetadata
	SubnetMask []byte
}

type DNSStruct struct {
	OptMd OptionMetadata
	DNS   []byte
}

type DomainNameStruct struct {
	OptMd      OptionMetadata
	DomainName []byte
}

type NetBIOSServerStruct struct {
	OptMd         OptionMetadata
	NetBIOSServer []byte
}

type RouterStruct struct {
	OptMd  OptionMetadata
	Router []byte
}

type StaticRouteStruct struct {
	OptMd       OptionMetadata
	StaticRoute []byte
}

type TFTPServerAddrStruct struct {
	OptMd          OptionMetadata
	TFTPServerAddr []byte
}

type VendorSpecificInfoStruct struct {
	OptMd              OptionMetadata
	VendorSpecificInfo []byte
}

type HostNameStruct struct {
	OptMd    OptionMetadata
	HostName []byte
}

type DhcpMsgTypeStruct struct {
	OptMd   OptionMetadata
	MsgType []byte
}

type ParamReqListStruct struct {
	OptMd        OptionMetadata
	ParamReqList []byte
}

type MaxDhcpMsgSizeStruct struct {
	OptMd          OptionMetadata
	MaxDhcpMsgSize []byte
}

type ClientIdStruct struct {
	OptMd    OptionMetadata
	ClientId []byte
}

type ServerIdStruct struct {
	OptMd    OptionMetadata
	ServerId []byte
}

type IPAddrLeaseTimeStruct struct {
	OptMd           OptionMetadata
	IPAddrLeaseTime []byte
}

type RenewalTimeStruct struct {
	OptMd       OptionMetadata
	RenewalTime []byte
}

type RebindingTimeStruct struct {
	OptMd         OptionMetadata
	RebindingTime []byte
}

type ReqIPAddrStruct struct {
	OptMd     OptionMetadata
	ReqIPAddr []byte
}

type OptionFields struct {
	SubnetMaskOpt         bool
	SubnetMaskVar         SubnetMaskStruct
	RouterOpt             bool
	RouterVar             RouterStruct
	DNSOpt                bool
	DNSVar                DNSStruct
	HostNameOpt           bool
	HostNameVar           HostNameStruct
	DomainNameOpt         bool
	DomainNameVar         DomainNameStruct
	StaticRouteOpt        bool
	StaticRouteVar        StaticRouteStruct
	VendorSpecificInfoOpt bool
	VendorSpecificInfoVar VendorSpecificInfoStruct
	NetBIOSServerOpt      bool
	NetBIOSServerVar      NetBIOSServerStruct
	ReqIPAddrOpt          bool
	ReqIPAddrVar          ReqIPAddrStruct
	IPAddrLeaseTimeOpt    bool
	IPAddrLeaseTimeVar    IPAddrLeaseTimeStruct
	DhcpMsgTypeOpt        bool
	DhcpMsgTypeVar        DhcpMsgTypeStruct
	ServerIdOpt           bool
	ServerIdVar           ServerIdStruct
	ParamReqListOpt       bool
	ParamReqListVar       ParamReqListStruct
	MaxDhcpMsgSizeOpt     bool
	MaxDhcpMsgSizeVar     MaxDhcpMsgSizeStruct
	RenewalTimeOpt        bool
	RenewalTimeVar        RenewalTimeStruct
	RebindingTimeOpt      bool
	RebindingTimeVar      RebindingTimeStruct
	ClientIdOpt           bool
	ClientIdVar           ClientIdStruct
	TFTPServerAddrOpt     bool
	TFTPServerAddrVar     TFTPServerAddrStruct
}

type DhcpOptionData struct {
	Length uint8
	Data   []byte
}

type BootPMsgStruct struct {
	OPCode           uint8
	HWAddrType       uint8
	HWAddrLen        uint8
	Hops             uint8
	TransactionId    uint32
	TimeElapsed      uint16
	Flags            uint16
	ClientIPAddr     uint32
	YourIPAddr       uint32
	ServerIPAddr     uint32
	RelayAgentIPAddr uint32
	ClientHWAddr     []byte
	ServerName       []byte
	BootFilename     []byte
	MagicCookie      uint32
	//DhcpOptFields    OptionFields
	DhcpOptionMap map[uint8]DhcpOptionData
}

func NewBootPMsgStruct() *BootPMsgStruct {
	return &BootPMsgStruct{}
}

type PktMetadata struct {
	ethHdrMd *EthHdrMetadata
	ipHdrMd  *IPHdrMetadata
	udpHdrMd *UDPHdrMetadata
}

func NewPktMetadata() *PktMetadata {
	return &PktMetadata{}
}

func decodeBootPMsg(data []byte, bootPMsgData *BootPMsgStruct) error {
	length := len(data)
	bootPMsgData.OPCode = uint8(data[0])
	bootPMsgData.HWAddrType = uint8(data[1])
	bootPMsgData.HWAddrLen = uint8(data[2])
	bootPMsgData.Hops = uint8(data[3])
	bootPMsgData.TransactionId = binary.BigEndian.Uint32(data[4:8])
	bootPMsgData.TimeElapsed = binary.BigEndian.Uint16(data[8:10])
	bootPMsgData.Flags = binary.BigEndian.Uint16(data[10:12])
	bootPMsgData.ClientIPAddr = binary.BigEndian.Uint32(data[12:16])
	bootPMsgData.YourIPAddr = binary.BigEndian.Uint32(data[16:20])
	bootPMsgData.ServerIPAddr = binary.BigEndian.Uint32(data[20:24])
	bootPMsgData.RelayAgentIPAddr = binary.BigEndian.Uint32(data[24:28])
	start := 28
	end := start + CLIENT_HW_ADDR_LEN
	temp := start + int(bootPMsgData.HWAddrLen)
	bootPMsgData.ClientHWAddr = make([]byte, bootPMsgData.HWAddrLen)
	copy(bootPMsgData.ClientHWAddr, data[start:temp])
	start = end
	end = start + MAX_LEN_SERVER_NAME
	bootPMsgData.ServerName = make([]byte, MAX_LEN_SERVER_NAME)
	copy(bootPMsgData.ServerName, data[start:end])
	start = end
	end = start + MAX_LEN_BOOT_FILENAME
	bootPMsgData.BootFilename = make([]byte, MAX_LEN_BOOT_FILENAME)
	copy(bootPMsgData.BootFilename, data[start:end])
	start = end

	if length > start+4 {
		end = start + 4
		bootPMsgData.MagicCookie = binary.BigEndian.Uint32(data[start:end])
	} else {
		return nil
	}

	//fmt.Println("Magic Cookie:", bootPMsgData.MagicCookie)
	//fmt.Println("Start:", start, "end:", end)
	start = end
	if start < length {
		bootPMsgData.DhcpOptionMap = make(map[uint8]DhcpOptionData)
		decodeDhcpOptions(bootPMsgData, data[start:])
		//fmt.Println("Data :", bootPMsgData.DhcpOptionMap)
	}
	return nil
}

func decodeDhcpOptions(bootPMsgData *BootPMsgStruct, data []byte) {
	fmt.Println("data :", data)

	start := 0
	end := 0
	length := len(data)

	for start < length {
		opCode := uint8(data[start])
		if opCode == EndOptCode {
			break
		}
		optLen := uint8(data[start+1])
		ent := bootPMsgData.DhcpOptionMap[opCode]
		ent.Length = optLen
		ent.Data = make([]byte, optLen)
		start = start + 2
		end = start + int(optLen)
		copy(ent.Data, data[start:end])
		bootPMsgData.DhcpOptionMap[opCode] = ent
		start = end
	}
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

/*
func decodeDhcpOptionFields(bootPMsgData *BootPMsgStruct, data []byte) {
	start := 0
	end := 0
	length := len(data)
	fmt.Println("Data:", data)
	for start < length {
		opCode := uint8(data[start])
		optLen := uint8(data[start+1])

		fmt.Println("OpCode:", opCode, "optLen:", optLen)
		start = start + 2
		switch opCode {
		case SubnetMaskOptCode:
			bootPMsgData.DhcpOptFields.SubnetMaskOpt = true
			bootPMsgData.DhcpOptFields.SubnetMaskVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.SubnetMaskVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.SubnetMaskVar.SubnetMask = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.SubnetMaskVar.SubnetMask, data[start:end])
			start = end
		case RouterOptCode:
			bootPMsgData.DhcpOptFields.RouterOpt = true
			bootPMsgData.DhcpOptFields.RouterVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.RouterVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.RouterVar.Router = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.RouterVar.Router, data[start:end])
			start = end
		case DNSOptCode:
			bootPMsgData.DhcpOptFields.DNSOpt = true
			bootPMsgData.DhcpOptFields.DNSVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.DNSVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.DNSVar.DNS = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.DNSVar.DNS, data[start:end])
			start = end
		case HostNameOptCode:
			bootPMsgData.DhcpOptFields.HostNameOpt = true
			bootPMsgData.DhcpOptFields.HostNameVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.HostNameVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.HostNameVar.HostName = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.HostNameVar.HostName, data[start:end])
			start = end
		case DomainNameOptCode:
			bootPMsgData.DhcpOptFields.DomainNameOpt = true
			bootPMsgData.DhcpOptFields.DomainNameVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.DomainNameVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.DomainNameVar.DomainName = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.DomainNameVar.DomainName, data[start:end])
			start = end
		case StaticRouteOptCode:
			bootPMsgData.DhcpOptFields.StaticRouteOpt = true
			bootPMsgData.DhcpOptFields.StaticRouteVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.StaticRouteVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.StaticRouteVar.StaticRoute = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.StaticRouteVar.StaticRoute, data[start:end])
			start = end
		case VendorSpecificInfoOptCode:
			bootPMsgData.DhcpOptFields.VendorSpecificInfoOpt = true
			bootPMsgData.DhcpOptFields.VendorSpecificInfoVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.VendorSpecificInfoVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.VendorSpecificInfoVar.VendorSpecificInfo = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.VendorSpecificInfoVar.VendorSpecificInfo, data[start:end])
			start = end
		case NetBIOSServerOptCode:
			bootPMsgData.DhcpOptFields.NetBIOSServerOpt = true
			bootPMsgData.DhcpOptFields.NetBIOSServerVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.NetBIOSServerVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.NetBIOSServerVar.NetBIOSServer = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.NetBIOSServerVar.NetBIOSServer, data[start:end])
			start = end
		case ReqIPAddrOptCode:
			bootPMsgData.DhcpOptFields.ReqIPAddrOpt = true
			bootPMsgData.DhcpOptFields.ReqIPAddrVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.ReqIPAddrVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.ReqIPAddrVar.ReqIPAddr = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.ReqIPAddrVar.ReqIPAddr, data[start:end])
			start = end
		case IPAddrLeaseTimeOptCode:
			bootPMsgData.DhcpOptFields.IPAddrLeaseTimeOpt = true
			bootPMsgData.DhcpOptFields.IPAddrLeaseTimeVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.IPAddrLeaseTimeVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.IPAddrLeaseTimeVar.IPAddrLeaseTime = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.IPAddrLeaseTimeVar.IPAddrLeaseTime, data[start:end])
			start = end
		case DhcpMsgTypeOptCode:
			bootPMsgData.DhcpOptFields.DhcpMsgTypeOpt = true
			bootPMsgData.DhcpOptFields.DhcpMsgTypeVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.DhcpMsgTypeVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.DhcpMsgTypeVar.MsgType = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.DhcpMsgTypeVar.MsgType, data[start:end])
			start = end
		case ServerIdOptCode:
			bootPMsgData.DhcpOptFields.ServerIdOpt = true
			bootPMsgData.DhcpOptFields.ServerIdVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.ServerIdVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.ServerIdVar.ServerId = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.ServerIdVar.ServerId, data[start:end])
			start = end
		case ParamReqListOptCode:
			bootPMsgData.DhcpOptFields.ParamReqListOpt = true
			bootPMsgData.DhcpOptFields.ParamReqListVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.ParamReqListVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.ParamReqListVar.ParamReqList = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.ParamReqListVar.ParamReqList, data[start:end])
			start = end
		case MaxDhcpMsgSizeOptCode:
			bootPMsgData.DhcpOptFields.MaxDhcpMsgSizeOpt = true
			bootPMsgData.DhcpOptFields.MaxDhcpMsgSizeVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.MaxDhcpMsgSizeVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.MaxDhcpMsgSizeVar.MaxDhcpMsgSize = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.MaxDhcpMsgSizeVar.MaxDhcpMsgSize, data[start:end])
			start = end
		case RenewalTimeOptCode:
			bootPMsgData.DhcpOptFields.RenewalTimeOpt = true
			bootPMsgData.DhcpOptFields.RenewalTimeVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.RenewalTimeVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.RenewalTimeVar.RenewalTime = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.RenewalTimeVar.RenewalTime, data[start:end])
			start = end
		case RebindingTimeOptCode:
			bootPMsgData.DhcpOptFields.RebindingTimeOpt = true
			bootPMsgData.DhcpOptFields.RebindingTimeVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.RebindingTimeVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.RebindingTimeVar.RebindingTime = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.RebindingTimeVar.RebindingTime, data[start:end])
			start = end
		case ClientIdOptCode:
			bootPMsgData.DhcpOptFields.ClientIdOpt = true
			bootPMsgData.DhcpOptFields.ClientIdVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.ClientIdVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.ClientIdVar.ClientId = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.ClientIdVar.ClientId, data[start:end])
			start = end
		case TFTPServerAddrOptCode:
			bootPMsgData.DhcpOptFields.TFTPServerAddrOpt = true
			bootPMsgData.DhcpOptFields.TFTPServerAddrVar.OptMd.Length = optLen
			bootPMsgData.DhcpOptFields.TFTPServerAddrVar.OptMd.OPCode = opCode
			bootPMsgData.DhcpOptFields.TFTPServerAddrVar.TFTPServerAddr = make([]byte, optLen)
			end = start + int(optLen)
			copy(bootPMsgData.DhcpOptFields.TFTPServerAddrVar.TFTPServerAddr, data[start:end])
			start = end
		default:
			fmt.Println("Invalid Op Code...", opCode)
		}
	}
	fmt.Println("Hello...", length, start)
}
*/
