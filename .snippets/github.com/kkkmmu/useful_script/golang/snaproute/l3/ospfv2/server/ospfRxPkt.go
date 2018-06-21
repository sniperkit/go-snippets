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
	"github.com/google/gopacket/pcap"
	"l3/ospfv2/objects"
	"time"
)

type OspfPktRecvStruct struct {
	OspfRecvPktCh       chan gopacket.Packet
	OspfRecvCtrlCh      chan bool
	OspfRecvCtrlReplyCh chan bool
	IntfConfKey         IntfConfKey
}

type OspfHelloPktRecvStruct struct {
	OspfRecvHelloPktCh       chan *OspfPktStruct
	OspfRecvHelloCtrlCh      chan bool
	OspfRecvHelloCtrlReplyCh chan bool
	IntfConfKey              IntfConfKey
}

type OspfLsaAndDbdPktRecvStruct struct {
	OspfRecvLsaAndDbdPktCh       chan *OspfPktStruct
	OspfRecvLsaAndDbdCtrlCh      chan bool
	OspfRecvLsaAndDbdCtrlReplyCh chan bool
	IntfConfKey                  IntfConfKey
}

func (server *OSPFV2Server) processOspfData(recvPktData *OspfPktStruct, helloPktMsgCh, lsaAndDbdMsgCh chan *OspfPktStruct) error {
	var err error = nil
	switch recvPktData.OspfHdrMd.PktType {
	case HelloType:
		helloPktMsgCh <- recvPktData
	case DBDescriptionType, LSRequestType, LSUpdateType, LSAckType:
		lsaAndDbdMsgCh <- recvPktData
	default:
		err = errors.New("Invalid Ospf packet type")
	}
	return err
}

func (server *OSPFV2Server) processIPv4Layer(ipLayer gopacket.Layer, IpAddr uint32, ipHdrMd *IpHdrMetadata) error {
	ipLayerContents := ipLayer.LayerContents()
	ipChkSum := binary.BigEndian.Uint16(ipLayerContents[10:12])
	binary.BigEndian.PutUint16(ipLayerContents[10:], 0)

	csum := computeCheckSum(ipLayerContents)
	if csum != ipChkSum {
		err := errors.New("Incorrect IPv4 checksum, hence dicarding the packet")
		return err
	}

	ipPkt := ipLayer.(*layers.IPv4)
	ipHdrMd.SrcIP, _ = convertDotNotationToUint32(ipPkt.SrcIP.To4().String())
	ipHdrMd.DstIP, _ = convertDotNotationToUint32(ipPkt.DstIP.To4().String())
	if IpAddr == ipHdrMd.SrcIP {
		err := errors.New(fmt.Sprintln("locally generated pkt", ipPkt.SrcIP, "hence dicarding the packet"))
		return err
	}

	if IpAddr != ipHdrMd.DstIP &&
		ALLDROUTER != ipHdrMd.DstIP &&
		ALLSPFROUTER != ipHdrMd.DstIP {
		err := errors.New(fmt.Sprintln("Incorrect DstIP", ipPkt.DstIP, "hence dicarding the packet"))
		return err
	}

	if ipPkt.Protocol != layers.IPProtocol(OSPF_PROTO_ID) {
		err := errors.New(fmt.Sprintln("Incorrect ProtocolID", ipPkt.Protocol, "hence dicarding the packet"))
		return err
	}
	if ALLSPFROUTER == ipHdrMd.DstIP {
		ipHdrMd.DstIPType = AllSPFRouterType
	} else if ALLDROUTER == ipHdrMd.DstIP {
		ipHdrMd.DstIPType = AllDRouterType
	} else {
		ipHdrMd.DstIPType = NormalType
	}
	return nil
}

func (server *OSPFV2Server) processOspfHeader(ospfPkt []byte, key IntfConfKey, md *OspfHdrMetadata, ipHdrMd *IpHdrMetadata) error {
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

	if OSPF_VERSION_2 != ospfHdr.Ver {
		err := errors.New("Dropped because of Ospf Version not matching")
		return err
	}

	if ent.AreaId == ospfHdr.AreaId {
		if ent.Type != objects.INTF_TYPE_POINT2POINT {
			if (ent.IpAddr & ent.Netmask) != (ipHdrMd.SrcIP & ent.Netmask) {
				err := errors.New("Dropped because of Src IP is not in subnet and Area ID is matching")
				return err

			}
		}
	} else {
		// We don't support Virtual Link
		err := errors.New("Dropped because Area ID is not matching and we dont support Virtual links, so this should not happend")
		return err

	}

	if ipHdrMd.DstIPType == AllDRouterType {
		if ent.DRtrId != server.globalData.RouterId &&
			ent.BDRtrId != server.globalData.RouterId {
			err := errors.New("Dropped because we should not recv any pkt with ALLDROUTER as we are not DR or BDR")
			return err
		}
	}

	//OSPF Auth Type
	if ent.AuthType != ospfHdr.AuthType {
		err := errors.New("Dropped because of Router Id not matching")
		return err
	}

	//TODO: We don't support Authentication

	if ospfHdr.PktType != HelloType {
		if ent.Type == objects.INTF_TYPE_BROADCAST {
			nbrKey := NbrConfKey{
				NbrIdentity:         ipHdrMd.SrcIP,
				NbrAddressLessIfIdx: key.IntfIdx,
			}
			_, exist := ent.NbrMap[nbrKey]
			if !exist {
				err := errors.New("Adjacency not established with this nbr")
				return err
			}
		} else if ent.Type == objects.INTF_TYPE_POINT2POINT {
			/* For future - For unnumbered P2P the identity will be
			   router id. */

			nbrKey := NbrConfKey{
				NbrIdentity:         ipHdrMd.SrcIP,
				NbrAddressLessIfIdx: key.IntfIdx,
			}
			_, exist := ent.NbrMap[nbrKey]
			if !exist {
				err := errors.New("Adjacency not established with this nbr")
				return err
			}
		}
	}

	//OSPF Header CheckSum
	binary.BigEndian.PutUint16(ospfPkt[12:14], 0)
	copy(ospfPkt[16:OSPF_HEADER_SIZE], []byte{0, 0, 0, 0, 0, 0, 0, 0})
	csum := computeCheckSum(ospfPkt)
	if csum != ospfHdr.Chksum {
		err := errors.New("Dropped because of invalid checksum")
		return err
	}

	md.PktType = ospfHdr.PktType
	md.Pktlen = ospfHdr.Pktlen
	md.RouterId = ospfHdr.RouterId
	md.AreaId = ospfHdr.AreaId
	if ospfHdr.AreaId == 0 {
		md.Backbone = true
	} else {
		md.Backbone = false
	}

	return nil
}

type OspfPktDataStruct struct {
	data      []byte
	ethHdrMd  *EthHdrMetadata
	ipHdrMd   *IpHdrMetadata
	ospfHdrMd *OspfHdrMetadata
	key       IntfConfKey
}

func (server *OSPFV2Server) processOspfPkt(pkt gopacket.Packet, key IntfConfKey, ospfPktData *OspfPktStruct) error {
	server.logger.Debug("Recevied Ospf Packet")
	ent, exist := server.IntfConfMap[key]
	if !exist {
		return errors.New("Dropped because of interface no more valid")
	}

	ethLayer := pkt.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return errors.New("Not an Ethernet frame")
	}
	eth := ethLayer.(*layers.Ethernet)

	ethHdrMd := NewEthHdrMetadata()
	ethHdrMd.SrcMAC = eth.SrcMAC
	ospfPktData.EthHdrMd = ethHdrMd
	ipLayer := pkt.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return errors.New("Not an IP packet")
	}

	ipHdrMd := NewIpHdrMetadata()
	err := server.processIPv4Layer(ipLayer, ent.IpAddr, ipHdrMd)
	if err != nil {
		return errors.New(fmt.Sprintln("Dropped because of IPv4 layer processing", err))
	}
	ospfPktData.IpHdrMd = ipHdrMd

	ospfHdrMd := NewOspfHdrMetadata()
	ospfPkt := ipLayer.LayerPayload()
	err = server.processOspfHeader(ospfPkt, key, ospfHdrMd, ipHdrMd)
	if err != nil {
		return errors.New(fmt.Sprintln("Dropped because of Ospf Header processing", err))
	}
	ospfPktData.OspfHdrMd = ospfHdrMd

	ospfPktData.Data = ospfPkt[OSPF_HEADER_SIZE:]
	return nil
}

func (server *OSPFV2Server) ProcessOspfRecvHelloPkt(recvPktData OspfHelloPktRecvStruct) {
	for {
		select {
		case msg := <-recvPktData.OspfRecvHelloPktCh:
			err := server.processRxHelloPkt(msg.Data, msg.OspfHdrMd, msg.IpHdrMd, msg.EthHdrMd, recvPktData.IntfConfKey)
			if err != nil {
				server.logger.Err("Error Processing Rx Hello Pkt:", err)
			}
		case _ = <-recvPktData.OspfRecvHelloCtrlCh:
			server.logger.Info("Stopping ProcessOspfRecvHelloPkt routine")
			recvPktData.OspfRecvHelloCtrlReplyCh <- true
		}
	}

}

func (server *OSPFV2Server) ProcessOspfRecvLsaAndDbdPkt(recvPktData OspfLsaAndDbdPktRecvStruct) {
	var nbrIdentity uint32
	ent, valid := server.IntfConfMap[recvPktData.IntfConfKey]
	if !valid {
		server.logger.Err("Intf : RecvLsaAndDbdPkt interface entry does not exist ",
			recvPktData.IntfConfKey)
		return
	}

	for {
		select {
		case msg := <-recvPktData.OspfRecvLsaAndDbdPktCh:
			if ent.Type == objects.INTF_TYPE_POINT2POINT {
				/* For future - add routerId as identity
				for unnumbered p2p */
				nbrIdentity = msg.IpHdrMd.SrcIP
			} else {
				nbrIdentity = msg.IpHdrMd.SrcIP
			}

			nbrKey := NbrConfKey{
				NbrIdentity:         nbrIdentity,
				NbrAddressLessIfIdx: recvPktData.IntfConfKey.IntfIdx,
			}

			switch msg.OspfHdrMd.PktType {
			case DBDescriptionType:
				err := server.ProcessRxDbdPkt(msg.Data, msg.OspfHdrMd, msg.IpHdrMd, nbrKey)
				if err != nil {
					server.logger.Err("Failed to process rx dbd pkt ", nbrKey)
				}
			case LSRequestType:
				err := server.ProcessRxLSAReqPkt(msg.Data, msg.OspfHdrMd, msg.IpHdrMd, nbrKey)
				if err != nil {
					server.logger.Err("Failed to process rx dbd pkt ", nbrKey)
				}
			case LSUpdateType:
				err := server.ProcessRxLsaUpdPkt(msg.Data, msg.OspfHdrMd, msg.IpHdrMd, nbrKey)
				if err != nil {
					server.logger.Err("Failed to process rx dbd pkt ", nbrKey)
				}
			case LSAckType:
				err := server.ProcessRxLSAAckPkt(msg.Data, msg.OspfHdrMd, msg.IpHdrMd, nbrKey)
				if err != nil {
					server.logger.Err("Failed to process rx dbd pkt ", nbrKey)
				}
			default:
				server.logger.Err("Invalid Packet type")
			}
		case _ = <-recvPktData.OspfRecvLsaAndDbdCtrlCh:
			server.logger.Info("Stopping ProcessOspfRecvLsaAndDbdPkt routine")
			recvPktData.OspfRecvLsaAndDbdCtrlReplyCh <- true
		}
	}
}

func (server *OSPFV2Server) ProcessOspfRecvPkt(recvPkt OspfPktRecvStruct) {
	ospfRecvHelloCtrlCh := make(chan bool)
	ospfRecvHelloCtrlReplyCh := make(chan bool)
	ospfRecvHelloPktCh := make(chan *OspfPktStruct, 1000)
	recvHelloPkt := OspfHelloPktRecvStruct{
		OspfRecvHelloPktCh:       ospfRecvHelloPktCh,
		OspfRecvHelloCtrlCh:      ospfRecvHelloCtrlCh,
		OspfRecvHelloCtrlReplyCh: ospfRecvHelloCtrlReplyCh,
		IntfConfKey:              recvPkt.IntfConfKey,
	}
	go server.ProcessOspfRecvHelloPkt(recvHelloPkt)
	ospfRecvLsaAndDbdCtrlCh := make(chan bool)
	ospfRecvLsaAndDbdCtrlReplyCh := make(chan bool)
	ospfRecvLsaAndDbdPktCh := make(chan *OspfPktStruct, 1000)
	recvLsaAndDbdPkt := OspfLsaAndDbdPktRecvStruct{
		OspfRecvLsaAndDbdPktCh:       ospfRecvLsaAndDbdPktCh,
		OspfRecvLsaAndDbdCtrlCh:      ospfRecvLsaAndDbdCtrlCh,
		OspfRecvLsaAndDbdCtrlReplyCh: ospfRecvLsaAndDbdCtrlReplyCh,
		IntfConfKey:                  recvPkt.IntfConfKey,
	}
	go server.ProcessOspfRecvLsaAndDbdPkt(recvLsaAndDbdPkt)
	for {
		select {
		case packet := <-recvPkt.OspfRecvPktCh:
			ospfPktData := NewOspfPktStruct()
			err := server.processOspfPkt(packet, recvPkt.IntfConfKey, ospfPktData)
			if err != nil {
				server.logger.Err("Error processing Ospf Pkt:", err)
				continue
			}
			server.processOspfData(ospfPktData, recvHelloPkt.OspfRecvHelloPktCh, recvLsaAndDbdPkt.OspfRecvLsaAndDbdPktCh)
		case _ = <-recvPkt.OspfRecvCtrlCh:
			server.logger.Info("Stopping ProcessOspfRecvPkt")
			recvHelloPkt.OspfRecvHelloCtrlCh <- true
			_ = <-recvHelloPkt.OspfRecvHelloCtrlReplyCh
			recvLsaAndDbdPkt.OspfRecvLsaAndDbdCtrlCh <- true
			_ = recvLsaAndDbdPkt.OspfRecvLsaAndDbdCtrlReplyCh
			recvPkt.OspfRecvCtrlReplyCh <- true
			return
		}
	}
}

func (server *OSPFV2Server) StartOspfRecvPkts(key IntfConfKey) {
	ospfRecvCtrlCh := make(chan bool)
	ospfRecvCtrlReplyCh := make(chan bool)
	ospfRecvPktCh := make(chan gopacket.Packet, 1000)
	recvPkt := OspfPktRecvStruct{
		OspfRecvPktCh:       ospfRecvPktCh,
		OspfRecvCtrlCh:      ospfRecvCtrlCh,
		OspfRecvCtrlReplyCh: ospfRecvCtrlReplyCh,
		IntfConfKey:         key,
	}
	ent, _ := server.IntfConfMap[key]
	go server.ProcessOspfRecvPkt(recvPkt)
	handle := ent.rxHdl.RecvPcapHdl
	recv := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := recv.Packets()
	for {
		select {
		case packet, ok := <-in:
			if ok {
				ipLayer := packet.Layer(layers.LayerTypeIPv4)
				if ipLayer == nil {
					server.logger.Err("Not an IP packet")
					continue
				}

				ipPkt := ipLayer.(*layers.IPv4)
				if ipPkt.Protocol == layers.IPProtocol(OSPF_PROTO_ID) {
					recvPkt.OspfRecvPktCh <- packet
				}
			}
		case _ = <-ent.rxHdl.PktRecvCtrlCh:
			server.logger.Info("Stopping the Recv Ospf packet thread")
			recvPkt.OspfRecvCtrlCh <- true
			_ = <-recvPkt.OspfRecvCtrlReplyCh
			ent.rxHdl.PktRecvCtrlReplyCh <- true
			return
		}
	}
}

func (server *OSPFV2Server) StopOspfRecvPkts(key IntfConfKey) {
	intfEnt, _ := server.IntfConfMap[key]
	intfEnt.rxHdl.PktRecvCtrlCh <- true
	cnt := 0
	for {
		select {
		case _ = <-intfEnt.rxHdl.PktRecvCtrlReplyCh:
			server.logger.Info("Stopped Recv Pkt thread")
			return
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt%1000 == 0 {
				server.logger.Err("Trying to stop the Rx thread")
			}
		}
	}
}

func (server *OSPFV2Server) InitRxPkt(intfKey IntfConfKey) error {
	intfEnt, _ := server.IntfConfMap[intfKey]
	ifName := intfEnt.IfName
	ipAddr := intfEnt.IpAddr
	recvHdl, err := pcap.OpenLive(ifName, snapshotLen, promiscuous, pcapTimeout)
	if err != nil {
		server.logger.Err("Error opening recv pcap handler", ifName)
		return err
	}
	ip := convertUint32ToDotNotation(ipAddr)
	filter := fmt.Sprintf("proto ospf and not src host %s", ip)
	server.logger.Info("Filter:", filter)
	err = recvHdl.SetBPFFilter(filter)
	if err != nil {
		server.logger.Err("Unable to set filter on", ifName)
		return err
	}
	intfEnt.rxHdl.RecvPcapHdl = recvHdl
	intfEnt.rxHdl.PktRecvCtrlCh = make(chan bool)
	intfEnt.rxHdl.PktRecvCtrlReplyCh = make(chan bool)
	server.IntfConfMap[intfKey] = intfEnt
	return nil
}

func (server *OSPFV2Server) DeinitRxPkt(intfKey IntfConfKey) {
	intfEnt, _ := server.IntfConfMap[intfKey]
	intfEnt.rxHdl.RecvPcapHdl.Close()
	intfEnt.rxHdl.RecvPcapHdl = nil
	intfEnt.rxHdl.PktRecvCtrlCh = nil
	intfEnt.rxHdl.PktRecvCtrlReplyCh = nil
	server.IntfConfMap[intfKey] = intfEnt
}
