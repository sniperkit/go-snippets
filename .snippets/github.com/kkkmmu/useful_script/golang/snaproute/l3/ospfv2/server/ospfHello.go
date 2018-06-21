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
	//"fmt"
	//    "bytes"
	"encoding/binary"
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/ospfv2/objects"
	"net"
	"time"
)

func (server *OSPFV2Server) SendHelloPkt(key IntfConfKey) {
	ospfHelloPkt := server.BuildHelloPkt(key)
	if ospfHelloPkt == nil {
		server.logger.Err("Unable to send the ospf Hello pkt")
		return
	}
	err := server.SendOspfPkt(key, ospfHelloPkt)
	if err != nil {
		server.logger.Err("Unable to send the ospf Hello pkt")
	}
	return
}

func (server *OSPFV2Server) BuildHelloPkt(key IntfConfKey) []byte {
	ent, exist := server.IntfConfMap[key]
	if !exist {
		server.logger.Err("Interface doesnot exist", key)
		return nil
	}
	ospfHdr := OSPFHeader{
		Ver:      OSPF_VERSION_2,
		PktType:  HelloType,
		Pktlen:   0,
		RouterId: server.globalData.RouterId,
		AreaId:   ent.AreaId,
		Chksum:   0,
		AuthType: ent.AuthType,
		//authKey:        ent.IfAuthKey,
	}

	//Rfc 2328 4.5
	isStub, err := server.isStubArea(ent.AreaId)
	if err != nil {
		server.logger.Err("Not sending Hello Packet, ", err)
		return nil
	}
	var option uint8
	if isStub {
		option = 0
	} else {
		option = EOption
	}
	helloData := OSPFHelloData{
		HelloInterval:   ent.HelloInterval,
		Options:         option,
		RtrPrio:         ent.RtrPriority,
		RtrDeadInterval: ent.RtrDeadInterval,
		DRtrIpAddr:      ent.DRIpAddr,
		BDRtrIpAddr:     ent.BDRIpAddr,
	}

	if key.IpAddr == 0 && key.IntfIdx != 0 &&
		ent.Type == objects.INTF_TYPE_POINT2POINT {
		helloData.Netmask = 0
	} else {
		helloData.Netmask = ent.Netmask
	}

	var nbrList []uint32
	for _, nbrEnt := range ent.NbrMap {
		nbr := nbrEnt.RtrId
		nbrList = append(nbrList, nbr)
	}

	ospfPktlen := OSPF_HEADER_SIZE + OSPF_HELLO_MIN_SIZE + len(nbrList)*4

	ospfHdr.Pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	helloDataEnc := encodeOspfHelloData(helloData, nbrList)
	ospf := append(ospfEncHdr, helloDataEnc...)
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	binary.BigEndian.PutUint64(ospf[16:24], ent.AuthKey)

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.Pktlen
	srcIp := net.ParseIP(convertUint32ToDotNotation(ent.IpAddr))
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    srcIp,
		DstIP:    net.IP{224, 0, 0, 5}, //ALLSPFROUTER
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       ent.IfMacAddr,
		DstMAC:       net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05},
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, gopacket.Payload(ospf))
	ospfPkt := buffer.Bytes()
	return ospfPkt
}

func (server *OSPFV2Server) processRxHelloPkt(data []byte,
	ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata,
	ethHdrMd *EthHdrMetadata,
	key IntfConfKey) error {

	ent, _ := server.IntfConfMap[key]
	ospfHelloData := NewOSPFHelloData()
	if len(data) < OSPF_HELLO_MIN_SIZE {
		err := errors.New("Invalid Hello Pkt data length")
		return err
	}
	decodeOspfHelloData(data, ospfHelloData)

	if ent.Type != objects.INTF_TYPE_POINT2POINT {
		if ent.Netmask != ospfHelloData.Netmask {
			server.logger.Info("HELLO: Netmask mismatch. Int mask", ent.Netmask, " Hello mask ", ospfHelloData.Netmask, " ip ", ipHdrMd.SrcIP)
			err := errors.New("Netmask mismatch")
			return err
		}
	}

	if ent.HelloInterval != ospfHelloData.HelloInterval {
		err := errors.New("Hello Interval mismatch")
		return err
	}

	if ent.RtrDeadInterval != ospfHelloData.RtrDeadInterval {
		err := errors.New("Router Dead Interval mismatch")
		return err
	}

	areaEnt, exist := server.AreaConfMap[ent.AreaId]
	if !exist {
		return errors.New("Area does not exist")
	}

	if areaEnt.ImportASExtern == true {
		if (ospfHelloData.Options & EOption) == 0 {
			return errors.New("External Routing Capability mismatch")
		}
	} else {
		if (ospfHelloData.Options & EOption) != 0 {
			return errors.New("External Routing Capability mismatch")
		}

	}

	TwoWayStatus := false
	for _, nbr := range ospfHelloData.NbrList {
		if nbr == server.globalData.RouterId {
			TwoWayStatus = true
			break
		}
	}

	/*
		srcIp := ipHdrMd.srcIP
		nbrKey := NbrConfKey{
			IpAddr:  srcIp,
			IntfIdx: key.IntfIdx,
		}
		ospfNbrIPToMAC[nbrKey] = ethHdrMd.srcMAC
	*/

	server.processOspfHelloNbr(ethHdrMd, ipHdrMd, ospfHdrMd, ospfHelloData, TwoWayStatus, key)

	return nil
}

func (server *OSPFV2Server) processOspfHelloNbr(ethHdrMd *EthHdrMetadata, ipHdrMd *IpHdrMetadata, ospfHdrMd *OspfHdrMetadata, ospfHelloData *OSPFHelloData, TwoWayStatus bool, key IntfConfKey) {

	routerId := ospfHdrMd.RouterId

	//Todo: Find whether one way or two way
	ent, _ := server.IntfConfMap[key]
	var nbrIdentity uint32
	if ent.Type == objects.INTF_TYPE_POINT2POINT {
		/*Currently unumbered p2p is not supported
		 hence nbrIdentity is SrcIp.
		for unnumbered p2p it is router id.
		*/
		nbrIdentity = ipHdrMd.SrcIP
		//nbrIdentity = ospfHdrMd.RouterId
	} else {
		nbrIdentity = ipHdrMd.SrcIP
	}

	nbrKey := NbrConfKey{
		NbrIdentity:         nbrIdentity,
		NbrAddressLessIfIdx: key.IntfIdx,
	}
	nbrEntry, exist := ent.NbrMap[nbrKey]
	if !exist {
		var nbrCreateMsg NbrCreateMsg
		nbrCreateMsg.RouterId = ospfHdrMd.RouterId
		nbrCreateMsg.NbrIP = ipHdrMd.SrcIP
		nbrCreateMsg.RtrPrio = ospfHelloData.RtrPrio
		nbrCreateMsg.TwoWayStatus = TwoWayStatus
		nbrCreateMsg.DRtrIpAddr = ospfHelloData.DRtrIpAddr
		nbrCreateMsg.BDRtrIpAddr = ospfHelloData.BDRtrIpAddr
		nbrCreateMsg.NbrKey = nbrKey
		ent.NbrCreateCh <- nbrCreateMsg
		server.logger.Info("Nbr Entry Created", nbrEntry)
	} else {
		if nbrEntry.TwoWayStatus != TwoWayStatus ||
			nbrEntry.DRtrIpAddr != ospfHelloData.DRtrIpAddr ||
			nbrEntry.BDRtrIpAddr != ospfHelloData.BDRtrIpAddr ||
			nbrEntry.RtrPrio != ospfHelloData.RtrPrio {
			var nbrChangeMsg NbrChangeMsg
			nbrChangeMsg.RouterId = ospfHdrMd.RouterId
			nbrChangeMsg.NbrIP = ipHdrMd.SrcIP
			nbrChangeMsg.TwoWayStatus = TwoWayStatus
			nbrChangeMsg.RtrPrio = ospfHelloData.RtrPrio
			nbrChangeMsg.DRtrIpAddr = ospfHelloData.DRtrIpAddr
			nbrChangeMsg.BDRtrIpAddr = ospfHelloData.BDRtrIpAddr
			nbrChangeMsg.NbrKey = nbrKey
			ent.NbrChangeCh <- nbrChangeMsg
		}
	}

	nbrDeadInterval := time.Duration(ent.RtrDeadInterval) * time.Second
	nbrHelloEventMsg := NbrHelloEventMsg{
		IntfConfKey:  key,
		RouterId:     routerId,
		RtrPrio:      ospfHelloData.RtrPrio,
		NbrIP:        ipHdrMd.SrcIP,
		NbrDeadTime:  nbrDeadInterval,
		TwoWayStatus: TwoWayStatus,
		NbrDRIpAddr:  ospfHelloData.DRtrIpAddr,
		NbrBDRIpAddr: ospfHelloData.BDRtrIpAddr,
		NbrMAC:       ethHdrMd.SrcMAC,
		NbrKey:       nbrKey,
	}
	server.SendHelloEventMsg(nbrHelloEventMsg)

	var backupSeenMsg BackupSeenMsg
	if TwoWayStatus == true && ent.FSMState == objects.INTF_FSM_STATE_WAITING {
		if ipHdrMd.SrcIP == ospfHelloData.DRtrIpAddr {
			if ospfHelloData.BDRtrIpAddr != 0 {
				ret := ent.WaitTimer.Stop()
				if ret == true {
					backupSeenMsg.RouterId = ospfHdrMd.RouterId
					backupSeenMsg.DRtrIpAddr = ipHdrMd.SrcIP
					backupSeenMsg.BDRtrIpAddr = ospfHelloData.BDRtrIpAddr
					server.logger.Info("Neigbor choose itself as Designated Router")
					server.logger.Info("Backup Designated Router also exist")
					ent.BackupSeenCh <- backupSeenMsg
				}
			}
		} else if ipHdrMd.SrcIP == ospfHelloData.BDRtrIpAddr {
			ret := ent.WaitTimer.Stop()
			if ret == true {
				server.logger.Info("Neigbor choose itself as Backup Designated Router")
				backupSeenMsg.RouterId = ospfHdrMd.RouterId
				backupSeenMsg.DRtrIpAddr = ospfHelloData.DRtrIpAddr
				backupSeenMsg.BDRtrIpAddr = ipHdrMd.SrcIP
				ent.BackupSeenCh <- backupSeenMsg
			}
		}
	}
}
