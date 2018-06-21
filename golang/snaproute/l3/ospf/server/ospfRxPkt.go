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
	"l3/ospf/config"
	"net"
	"time"
)

func (server *OSPFServer) processIPv4Layer(ipLayer gopacket.Layer, ipAddr net.IP, ipHdrMd *IpHdrMetadata) error {
	ipLayerContents := ipLayer.LayerContents()
	ipChkSum := binary.BigEndian.Uint16(ipLayerContents[10:12])
	binary.BigEndian.PutUint16(ipLayerContents[10:], 0)
	allSPFRouter := net.ParseIP(ALLSPFROUTER)
	allDRouter := net.ParseIP(ALLDROUTER)

	csum := computeCheckSum(ipLayerContents)
	if csum != ipChkSum {
		err := errors.New("Incorrect IPv4 checksum, hence dicarding the packet")
		return err
	}

	ipPkt := ipLayer.(*layers.IPv4)
	if ipAddr.Equal(ipPkt.SrcIP) == true {
		err := errors.New(fmt.Sprintln("locally generated pkt", ipPkt.SrcIP, "hence dicarding the packet"))
		return err
	}

	if ipAddr.Equal(ipPkt.DstIP) == false &&
		allSPFRouter.Equal(ipPkt.DstIP) == false &&
		allDRouter.Equal(ipPkt.DstIP) == false {
		err := errors.New(fmt.Sprintln("Incorrect DstIP", ipPkt.DstIP, "hence dicarding the packet"))
		return err
	}

	if ipPkt.Protocol != layers.IPProtocol(OSPF_PROTO_ID) {
		err := errors.New(fmt.Sprintln("Incorrect ProtocolID", ipPkt.Protocol, "hence dicarding the packet"))
		return err
	}

	ipHdrMd.srcIP = ipPkt.SrcIP.To4()
	ipHdrMd.dstIP = ipPkt.DstIP.To4()
	if allSPFRouter.Equal(ipPkt.DstIP) {
		ipHdrMd.dstIPType = AllSPFRouter
	} else if allDRouter.Equal(ipPkt.DstIP) {
		ipHdrMd.dstIPType = AllDRouter
	} else {
		ipHdrMd.dstIPType = Normal
	}
	/*
	   ToDo:
	   RFC 2328 Section 8.2
	   1. Destination IP (AllDRouters)
	*/
	return nil
}

func (server *OSPFServer) processOspfHeader(ospfPkt []byte, key IntfConfKey, md *OspfHdrMetadata) error {
	if len(ospfPkt) < OSPF_HEADER_SIZE {
		err := errors.New("Invalid length of Ospf Header")
		return err
	}

	ent, exist := server.IntfConfMap[key]
	if !exist {
		err := errors.New("Dropped because of interface no more valid")
		return err
	}

	ospfHdr := NewOSPFHeader()

	decodeOspfHdr(ospfPkt, ospfHdr)

	if server.ospfGlobalConf.Version != ospfHdr.ver {
		err := errors.New("Dropped because of Ospf Version not matching")
		return err
	}

	if ent.IfType != config.NumberedP2P || ent.IfType != config.UnnumberedP2P {
		if bytesEqual(ent.IfAreaId, ospfHdr.areaId) == false &&
			isInSubnet(net.IP(ent.IfAreaId), net.IP(ospfHdr.areaId), net.IPMask(ent.IfNetmask)) == false {
			err := errors.New("Dropped because of Src IP is not in subnet or Area ID not matching")
			return err
		}
	}

	// Todo: when areaId is of backbone
	if bytesEqual(ospfHdr.areaId, []byte{0, 0, 0, 0}) == true {
		md.backbone = true
	} else {
		md.backbone = false
	}

	//OSPF Auth Type
	if ent.IfAuthType != ospfHdr.authType {
		err := errors.New("Dropped because of Router Id not matching")
		return err
	}

	//OSPF Header CheckSum
	binary.BigEndian.PutUint16(ospfPkt[12:14], 0)
	copy(ospfPkt[16:OSPF_HEADER_SIZE], []byte{0, 0, 0, 0, 0, 0, 0, 0})
	csum := computeCheckSum(ospfPkt)
	if csum != ospfHdr.chksum {
		err := errors.New("Dropped because of invalid checksum")
		return err
	}

	/*
	   ToDo:
	   RFC 2328 Section 8.2
	   1. Complete AreaID check
	   2. Authentication
	*/
	md.pktType = OspfType(ospfHdr.pktType)
	md.pktlen = ospfHdr.pktlen
	md.routerId = ospfHdr.routerId
	md.areaId = binary.BigEndian.Uint32(ospfHdr.areaId)
	return nil
}

func (server *OSPFServer) ProcessOspfRecvPkt(key IntfConfKey, pkt gopacket.Packet) {
	//server.logger.Info(fmt.Sprintf("Recevied Ospf Packet"))

	ethLayer := pkt.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		server.logger.Err("Not an Ethernet frame")
		return
	}
	eth := ethLayer.(*layers.Ethernet)

	ethHdrMd := NewEthHdrMetadata()
	ethHdrMd.srcMAC = eth.SrcMAC

	ipLayer := pkt.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		server.logger.Err("Not an IP packet")
		return
	}

	ent, exist := server.IntfConfMap[key]
	if !exist {
		server.logger.Err("Dropped because of interface no more valid")
		return
	}

	ipHdrMd := NewIpHdrMetadata()
	//@liwei: Do IP layer validation check.
	err := server.processIPv4Layer(ipLayer, ent.IfIpAddr, ipHdrMd)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Dropped because of IPv4 layer processing", err))
		return
	} else {
		//server.logger.Info("IPv4 Header is processed succesfully")
	}

	ospfHdrMd := NewOspfHdrMetadata()
	ospfPkt := ipLayer.LayerPayload()
	//@liwei: Do OSPF header validation check.
	err = server.processOspfHeader(ospfPkt, key, ospfHdrMd)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Dropped because of Ospf Header processing", err))
		return
	} else {
		//server.logger.Info("Ospfv2 Header is processed successfully")
	}

	//@liwei: Realy handle the packet.
	ospfData := ospfPkt[OSPF_HEADER_SIZE:]
	err = server.processOspfData(ospfData, ethHdrMd, ipHdrMd, ospfHdrMd, key)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Dropped because of Ospf Header processing", err))
		return
	} else {
		//server.logger.Info("Ospfv2 Header is processed successfully")
	}
	return
}

func (server *OSPFServer) processOspfData(data []byte, ethHdrMd *EthHdrMetadata, ipHdrMd *IpHdrMetadata, ospfHdrMd *OspfHdrMetadata, key IntfConfKey) error {
	var err error = nil
	//routerid := binary.BigEndian.Uint32(ospfHdrMd.routerId)
	NeighborIP := net.IPv4(ipHdrMd.srcIP[0], ipHdrMd.srcIP[1], ipHdrMd.srcIP[2], ipHdrMd.srcIP[3])
	//ipaddr := convertByteToOctetString(ipHdrMd.srcIP)
	ospfNbrConfKey := NeighborConfKey{
		IPAddr:  config.IpAddress(NeighborIP.String()),
		IntfIdx: key.IntfIdx,
	}
	exist := server.neighborExist(ospfNbrConfKey)
	if !exist {
		server.logger.Info(fmt.Sprintln("PACKET: neighbor doesnt exist..", NeighborIP, key.IntfIdx))
	}
	//@liwei: Handle different packet by type.
	switch ospfHdrMd.pktType {
	case HelloType:
		err = server.processRxHelloPkt(data, ospfHdrMd, ipHdrMd, ethHdrMd, key)
	case DBDescriptionType:
		if exist {
			err = server.ProcessRxDbdPkt(data, ospfHdrMd, ipHdrMd, key, ethHdrMd.srcMAC)
		}
	case LSRequestType:
		if exist {
			err = server.ProcessRxLSAReqPkt(data, ospfHdrMd, ipHdrMd, key)
		}
	case LSUpdateType:
		if exist {
			err = server.ProcessRxLsaUpdPkt(data, ospfHdrMd, ipHdrMd, key)
		}
	case LSAckType:
		if exist {
			err = server.ProcessRxLSAAckPkt(data, ospfHdrMd, ipHdrMd, key)
		}
	default:
		err = errors.New("Invalid Ospf packet type")
	}
	return err
}

func (server *OSPFServer) StartOspfRecvPkts(key IntfConfKey) {
	ent, _ := server.IntfRxMap[key]
	handle := ent.RecvPcapHdl
	recv := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := recv.Packets()
	//@liwei: Start OSPF packet handle logic.
	for {
		select {
		case packet, ok := <-in:
			if ok {
				//server.logger.Info("Got Some Ospf Packet on the Recv Thread")
				//@liwei: OSPF packet handler
				go server.ProcessOspfRecvPkt(key, packet)
			}
		case state := <-ent.PktRecvCh:
			if state == false {
				server.logger.Info("Stopping the Recv Ospf packet thread")
				ent.PktRecvStatusCh <- false
				return
			}
		}
	}
}

func (server *OSPFServer) StopOspfRecvPkts(key IntfConfKey) {
	ent, _ := server.IntfRxMap[key]
	ent.PktRecvCh <- false
	cnt := 0
	for {
		select {
		case status := <-ent.PktRecvStatusCh:
			if status == false { // False Means Recv Pkt Thread Stopped
				server.logger.Info("Stopped Recv Pkt thread")
				return
			}
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt == 100 {
				server.logger.Err("Unable to stop the Rx thread")
				return
			}
		}
	}
}
