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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/ospf/config"
	"l3/ospfv2/objects"
	"net"
	"time"
)

/*
LSA request
 0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |   Version #   |       3       |         Packet length         |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          Router ID                            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                           Area ID                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |           Checksum            |             AuType            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          LS type                              |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Link State ID                           |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |
*/
type ospfLSAReq struct {
	ls_type       uint32
	link_state_id uint32
	adv_router_id uint32
}

type ospfNeighborLSDBMsg struct {
	areaId uint32
	data   []byte
}

func newNbrAckTxMsg() *NbrAckTxMsg {
	return &NbrAckTxMsg{}
}

func NewospfNeighborLSDBMsg() *ospfNeighborLSDBMsg {
	return &ospfNeighborLSDBMsg{}
}

func newNbrLsaAckMsg() *NbrLsaAckMsg {
	return &NbrLsaAckMsg{}
}

type NbrLsaUpdMsg struct {
	nbrKey NbrConfKey
	data   []byte
	areaId uint32
}

type ospfNeighborLSAUpdPkt struct {
	no_lsas uint32
	lsa     []byte
}

func newospfNeighborLSAUpdPkt() *ospfNeighborLSAUpdPkt {
	return &ospfNeighborLSAUpdPkt{}
}

func getLsaHeaderFromLsa(ls_age uint16, options uint8, ls_type uint8, link_state_id uint32,
	adv_router_id uint32, ls_sequence_num uint32, ls_checksum uint16, ls_len uint16) *ospfLSAHeader {

	lsa_header := &ospfLSAHeader{}
	lsa_header.ls_age = ls_age
	lsa_header.options = options
	lsa_header.ls_type = ls_type
	lsa_header.link_state_id = link_state_id
	lsa_header.adv_router_id = adv_router_id
	lsa_header.ls_sequence_num = ls_sequence_num
	lsa_header.ls_checksum = ls_checksum
	lsa_header.ls_len = ls_len
	return lsa_header
}

func decodeLSAReq(data []byte) (lsa_req ospfLSAReq) {
	lsa_req.ls_type = binary.BigEndian.Uint32(data[0:4])
	lsa_req.link_state_id = binary.BigEndian.Uint32(data[4:8])
	lsa_req.adv_router_id = binary.BigEndian.Uint32(data[8:12])
	return lsa_req
}

func decodeLSAReqPkt(data []byte, pktlen uint16) []ospfLSAReq {
	no_of_lsa := int((pktlen - OSPF_HEADER_SIZE) / OSPF_LSA_REQ_SIZE)
	lsa_req_pkt := []ospfLSAReq{}
	start := 0
	end := OSPF_LSA_REQ_SIZE
	for i := 0; i < no_of_lsa; i++ {
		lsa_req := decodeLSAReq(data[start:end])
		lsa_req_pkt = append(lsa_req_pkt, lsa_req)
		start = end
		end += OSPF_LSA_REQ_SIZE
	}
	return lsa_req_pkt
}

func encodeLSAReq(lsa_data []ospfLSAReq) []byte {
	lsa_pkt := []byte{}
	for i := 0; i < len(lsa_data); i++ {
		pkt := make([]byte, OSPF_LSA_REQ_SIZE)
		binary.BigEndian.PutUint32(pkt[0:4], lsa_data[i].ls_type)
		binary.BigEndian.PutUint32(pkt[4:8], lsa_data[i].link_state_id)
		binary.BigEndian.PutUint32(pkt[8:12], lsa_data[i].adv_router_id)
		lsa_pkt = append(pkt, lsa_pkt...)
	}
	return lsa_pkt
}

func (server *OSPFV2Server) EncodeLSAReqPkt(intfKey IntfConfKey, ent IntfConf,
	nbrConf NbrConf, lsa_req_pkt []ospfLSAReq, dstMAC net.HardwareAddr) (data []byte) {
	ospfHdr := OSPFHeader{
		Ver:      OSPF_VERSION_2,
		PktType:  uint8(LSRequestType),
		Pktlen:   0,
		RouterId: server.globalData.RouterId,
		AreaId:   ent.AreaId,
		Chksum:   0,
		AuthType: ent.AuthType,
	}

	lsaDataEnc := encodeLSAReq(lsa_req_pkt)
	ospfPktlen := OSPF_HEADER_SIZE
	ospfPktlen = ospfPktlen + len(lsaDataEnc)

	ospfHdr.Pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	server.logger.Debug("ospfEncHdr:", ospfEncHdr)
	server.logger.Debug("lsa Pkt:", lsaDataEnc)

	ospf := append(ospfEncHdr, lsaDataEnc...)
	server.logger.Debug("OSPF LSA REQ:", ospf)
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	binary.BigEndian.PutUint64(ospf[16:24], ent.AuthKey)

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.Pktlen
	var dstIp net.IP
	if ent.FSMState == objects.INTF_FSM_STATE_P2P {
		dstIp = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(config.McastMAC)
	} else {
		dstIp = net.ParseIP(convertUint32ToDotNotation(nbrConf.NbrIP))
	}
	srcIp := net.ParseIP(convertUint32ToDotNotation(ent.IpAddr))
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    srcIp,
		DstIP:    dstIp,
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       ent.IfMacAddr,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, gopacket.Payload(ospf))
	lsaPkt := buffer.Bytes()
	server.logger.Debug("lsaPkt: ", lsaPkt)

	return lsaPkt

}

func (server *OSPFV2Server) BuildAndSendLSAReq(nbrId NbrConfKey, nbrConf NbrConf) (curr_index int) {
	/* calculate max no of requests that can be added
	for req packet */

	var add_items int

	var req ospfLSAReq
	var i int

	msg := NbrLsaReqMsg{}
	msg.lsa_slice = []ospfLSAReq{}
	msg.nbrKey = nbrId

	reqlist := nbrConf.NbrReqList
	if len(reqlist) < 1 && nbrConf.State != NbrFull {
		server.logger.Warning("Nbr : Req list is nill. Nbr is full. ", nbrId)
		server.ProcessNbrFull(nbrId)
		return nbrConf.NbrReqListIndex
	}
	if len(reqlist) == 0 {
		return nbrConf.NbrReqListIndex
	}
	req_list_items := len(reqlist) - nbrConf.NbrReqListIndex
	max_req := calculateMaxLsaReq()
	server.logger.Debug("Nbr : max req ", max_req)
	if max_req > req_list_items {
		add_items = req_list_items
		nbrConf.NbrReqListIndex = len(reqlist)

	} else {
		add_items = max_req
		nbrConf.NbrReqListIndex += max_req
	}
	if add_items < 1 {
		server.logger.Debug("Nbr : All requests sent out")
		server.ProcessNbrFull(nbrId)
		return nbrConf.NbrReqListIndex
	}
	server.logger.Debug("LSAREQ: nbrIndex ",
		nbrConf.NbrReqListIndex, " add_items ", add_items, " req_list len ", len(reqlist))
	for i = 0; i < add_items; i++ {
		req.ls_type = uint32(reqlist[i].ls_type)
		req.link_state_id = reqlist[i].link_state_id
		req.adv_router_id = reqlist[i].adv_router_id
		msg.lsa_slice = append(msg.lsa_slice, req)

		lsid := net.ParseIP(convertUint32ToDotNotation(req.link_state_id))
		adv_rtr := net.ParseIP(convertUint32ToDotNotation(req.adv_router_id))
		server.logger.Debug("LSA request: Send req to nbr ", nbrId.NbrIdentity,
			" lsid ", lsid, " rtrid ", adv_rtr, " lstype ", req.ls_type)
	}
	if nbrConf.NbrReqListIndex == len(reqlist) {
		nbrConf.NbrReqList = []*ospfLSAHeader{}
	} else {
		newList := []*ospfLSAHeader{}
		newList = append(newList, reqlist[:add_items]...)
		nbrConf.NbrReqList = newList
	}
	server.logger.Debug("LSA request: total requests out, req_list_len, current req_list_index ", add_items, len(msg.lsa_slice), nbrConf.NbrReqListIndex)
	server.logger.Debug("LSA request: lsa_req", msg.lsa_slice)

	if len(msg.lsa_slice) != 0 {
		intConf, exist := server.IntfConfMap[nbrConf.IntfKey]
		if !exist {
			server.logger.Err("Nbr: Intfconf does not exist.No lsa req sent ", nbrConf.IntfKey)
			return 0
		}
		data := server.EncodeLSAReqPkt(nbrConf.IntfKey, intConf, nbrConf, msg.lsa_slice, nbrConf.NbrMac)
		server.SendOspfPkt(nbrConf.IntfKey, data)
	}
	return nbrConf.NbrReqListIndex
}

/*
LSA update packet
   0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |   Version #   |       4       | d        Packet length         |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          Router ID                            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                           Area ID                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |           Checksum            |             AuType            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                            # LSAs                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                                                               |
       +-                                                            +-+
       |                             LSAs                              |
       +-                                                            +-+
       |                              ...                              |
*/

func (server *OSPFV2Server) BuildLsaUpdPkt(intfKey IntfConfKey, ent IntfConf,
	dstMAC net.HardwareAddr, dstIp net.IP, lsa_pkt_size int, lsaUpdEnc []byte) (data []byte) {
	ospfHdr := OSPFHeader{
		Ver:      OSPF_VERSION_2,
		PktType:  uint8(LSUpdateType),
		Pktlen:   0,
		RouterId: server.globalData.RouterId,
		AreaId:   ent.AreaId,
		Chksum:   0,
		AuthType: ent.AuthType,
		//authKey:        ent.AuthKey,
	}

	ospfPktlen := OSPF_HEADER_SIZE

	ospfPktlen = ospfPktlen + len(lsaUpdEnc)
	ospfHdr.Pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)

	//server.logger.Debug(fmt.Sprintln("LSA upd Pkt:", lsaUpdEnc))

	ospf := append(ospfEncHdr, lsaUpdEnc...)
	//server.logger.Debug(fmt.Sprintln("OSPF LSA UPD:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	binary.BigEndian.PutUint64(ospf[16:24], ent.AuthKey)
	srcIp := net.ParseIP(convertUint32ToDotNotation(ent.IpAddr))
	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.Pktlen
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    srcIp,
		DstIP:    dstIp,
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       ent.IfMacAddr,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, gopacket.Payload(ospf))
	//server.logger.Debug(fmt.Sprintln("buffer: ", buffer))
	lsaUpd := buffer.Bytes()
	//server.logger.Debug(fmt.Sprintln("flood Pkt: ", lsaUpd))

	return lsaUpd

}

/*
Packet from Rx thread is processed here
in rx thread context
*/
func (server *OSPFV2Server) ProcessRxLsaUpdPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, nbrKey NbrConfKey) error {

	msg := NbrLsaUpdMsg{
		nbrKey: nbrKey,
		areaId: ospfHdrMd.AreaId,
		data:   data,
	}

	server.NbrConfData.neighborLSAUpdEventCh <- msg
	server.logger.Debug("LSA update: Received LSA update with router_id , lentgh ", ospfHdrMd.RouterId, ospfHdrMd.Pktlen)
	//	server.logger.Debug(fmt.Sprintln("LSA update: pkt byte[]: ", data))
	return nil
}

/*
@fn processLSAUpdEvent
 Get total lsas.
  For each LSA :
		1) decode LSA
		2) get LSA instance from LSDB
		3) perform sanity check on LSA.
		4) update/delete/reject based on sanity.
		5) Send ACK if needed.

*/

func (server *OSPFV2Server) ProcessLsaUpd(msg NbrLsaUpdMsg) {
	nbr, exists := server.NbrConfMap[msg.nbrKey]
	discard := true
	if !exists {
		return
	}

	intf := server.IntfConfMap[nbr.IntfKey]
	lsa_max_age := false
	discard = server.lsaUpdDiscardCheck(nbr, msg.data)
	if discard {
		server.logger.Err("LSAUPD: Discard. nbr ", msg.nbrKey)
		return
	}

	no_lsa := binary.BigEndian.Uint32(msg.data[0:4])
	server.logger.Debug("LSAUPD: Nbr, No of LSAs ", msg.nbrKey, no_lsa, "  len  ", len(msg.data))
	//server.logger.Debug(fmt.Sprintln("LSUPD:LSA pkt ", msg.data))
	lsa_header := NewLsaHeader()
	/* decode each LSA and send to lsdb
	 */
	index := 4
	end_index := 0
	lsa_header_byte := make([]byte, OSPF_LSA_HEADER_SIZE)
	for i := 0; i < int(no_lsa); i++ {
		decodeLsaHeader(msg.data[index:index+OSPF_LSA_HEADER_SIZE], lsa_header)
		copy(lsa_header_byte, msg.data[index:index+OSPF_LSA_HEADER_SIZE])
		server.logger.Debug("LSAUPD: lsaheader decoded adv_rter ", lsa_header.Adv_router,
			" linkid ", lsa_header.LinkId, " lsage ", lsa_header.LSAge,
			" checksum ", lsa_header.LSChecksum, " seq num ", lsa_header.LSSequenceNum,
			" LSTYPE ", lsa_header.LSType,
			" len ", lsa_header.length)
		end_index = int(lsa_header.length) + index /* length includes data + header */
		currLsa := make([]byte, end_index-index)
		copy(currLsa, msg.data[index:end_index])
		if lsa_header.LSAge == LSA_MAX_AGE {
			lsa_max_age = true
		}
		/* send message to lsdb */
		lsdb_msg := RecvdLsaMsg{}
		lsdbKey := LsdbKey{
			AreaId: msg.areaId,
		}
		lsdb_msg.LsdbKey = lsdbKey
		valid := validateChecksum(currLsa)
		if !valid {
			server.logger.Debug("LSAUPD: Invalid checksum. Nbr",
				server.NbrConfMap[msg.nbrKey])
			//continue
		}
		lsa_key := NewLsaKey()
		selfGenLsaMsg := RecvdSelfLsaMsg{}
		switch lsa_header.LSType {
		case RouterLSA:
			rlsa := NewRouterLsa()
			decodeRouterLsa(currLsa, rlsa, lsa_key)

			drlsa, ret := server.getRouterLsaFromLsdb(msg.areaId, *lsa_key)
			discard, _ = server.sanityCheckRouterLsa(*rlsa, drlsa, nbr, intf, ret, lsa_max_age)
			lsdb_msg.LsaData = *rlsa
			selfGenLsaMsg.LsaData = *rlsa

		case NetworkLSA:
			nlsa := NewNetworkLsa()
			decodeNetworkLsa(currLsa, nlsa, lsa_key)
			dnlsa, ret := server.getNetworkLsaFromLsdb(msg.areaId, *lsa_key)
			discard, _ = server.sanityCheckNetworkLsa(*lsa_key, *nlsa, dnlsa, nbr, intf, ret, lsa_max_age)
			lsdb_msg.LsaData = *nlsa
			selfGenLsaMsg.LsaData = *nlsa

		case Summary3LSA, Summary4LSA:
			server.logger.Debug("Received summary Lsa Packet :", lsdb_msg.LsaData)
			slsa := NewSummaryLsa()
			decodeSummaryLsa(currLsa, slsa, lsa_key)
			server.logger.Debug("Decoded summary Lsa Packet :", slsa)
			dslsa, ret := server.getSummaryLsaFromLsdb(msg.areaId, *lsa_key)
			discard, _ = server.sanityCheckSummaryLsa(*slsa, dslsa, nbr, intf, ret, lsa_max_age)
			lsdb_msg.LsaData = *slsa
			selfGenLsaMsg.LsaData = *slsa

		case ASExternalLSA:
			alsa := NewASExternalLsa()
			decodeASExternalLsa(currLsa, alsa, lsa_key)
			dalsa, ret := server.getASExternalLsaFromLsdb(msg.areaId, *lsa_key)
			discard, _ = server.sanityCheckASExternalLsa(*alsa, dalsa, nbr, intf, ret, lsa_max_age)
			lsdb_msg.LsaData = *alsa
			selfGenLsaMsg.LsaData = *alsa
		}

		lsid := lsa_header.LinkId
		router_id := lsa_header.Adv_router

		self_gen := false
		self_gen = server.selfGenLsaCheck(*lsa_key)
		if self_gen {
			//send message to lsdb for self gen LSA
			selfGenLsaMsg = RecvdSelfLsaMsg{
				LsaKey:  *lsa_key,
				LsdbKey: lsdbKey,
				LsaData: lsdb_msg.LsaData,
			}
			server.MessagingChData.NbrFSMToLsdbChData.RecvdSelfLsaMsgCh <- selfGenLsaMsg
			server.logger.Debug("LSAUPD: discard . Received self generated. ", lsa_key)

		}

		if !discard && !self_gen {
			server.logger.Debug("LSAUPD: add to lsdb lsid ", lsid, " router_id ", router_id, " lstype ", lsa_header.LSType)
			lsdb_msg.MsgType = LSA_ADD
			lsdb_msg.LsaKey = *lsa_key
			server.MessagingChData.NbrFSMToLsdbChData.RecvdLsaMsgCh <- lsdb_msg

		}
		flood_pkt := NbrToFloodMsg{
			NbrKey:  msg.nbrKey,
			MsgType: LSA_FLOOD_ALL,
			LsaType: lsa_header.LSType,
		}
		flood_pkt.LsaPkt = make([]byte, end_index-index)
		copy(flood_pkt.LsaPkt, currLsa)
		if !self_gen { // for ABR summary lsa is flooded after LSDB/SPF changes are done.
			server.MessagingChData.NbrFSMToFloodChData.LsaFloodCh <- flood_pkt
		}

		/* send ACK */
		lsaAckMsg := newNbrAckTxMsg()
		lsaAckMsg.lsa_headers_byte = append(lsaAckMsg.lsa_headers_byte, lsa_header_byte...)
		lsaAckMsg.nbrKey = msg.nbrKey
		server.logger.Debug("ACK TX: nbr ", msg.nbrKey, " ack ", lsaAckMsg.lsa_headers_byte)
		server.processTxLsaAck(*lsaAckMsg)

		index = end_index
		//server.UpdateNbrList(msg.nbrKey)

	}
	server.BuildAndSendLSAReq(msg.nbrKey, nbr)

}

/* link state ACK packet
0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |   Version #   |       5       |         Packet length         |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                          Router ID                            |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                           Area ID                             |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |           Checksum            |             AuType            |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                       Authentication                          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                       Authentication                          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                                                               |
      +-                                                             -+
      |                             A                                 |
      +-                 Link State Advertisement                    -+
      |                           Header                              |
      +-                                                             -+
      |                                                               |
      +-                                                             -+
      |                                                               |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                              ...                              |
*/

//func (server *OSPFV2Server) encodeLSAAck
func (server *OSPFV2Server) BuildLSAAckPkt(intfKey IntfConfKey, ent IntfConf,
	nbrConf NbrConf, dstMAC net.HardwareAddr, dstIp net.IP, lsa_pkt_size int, lsaAckEnc []byte) (data []byte) {
	ospfHdr := OSPFHeader{
		Ver:      OSPF_VERSION_2,
		PktType:  uint8(LSAckType),
		Pktlen:   0,
		RouterId: server.globalData.RouterId,
		AreaId:   ent.AreaId,
		Chksum:   0,
		AuthType: ent.AuthType,
	}

	ospfPktlen := OSPF_HEADER_SIZE
	ospfPktlen = ospfPktlen + lsa_pkt_size
	ospfHdr.Pktlen = uint16(ospfPktlen)
	server.logger.Debug("LSAACK : packet legth header(24) + ack ", ospfPktlen)
	ospfEncHdr := encodeOspfHdr(ospfHdr)
	//server.logger.Debug(fmt.Sprintln("ospfEncHdr:", ospfEncHdr))

	//server.logger.Debug(fmt.Sprintln("LSA upd Pkt:", lsaAckEnc))

	ospf := append(ospfEncHdr, lsaAckEnc...)
	//server.logger.Debug(fmt.Sprintln("OSPF LSA ACK:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	binary.BigEndian.PutUint64(ospf[16:24], ent.AuthKey)

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.Pktlen
	if ent.FSMState == objects.INTF_FSM_STATE_P2P {
		dstIp = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(config.McastMAC)
	}
	srcIp := net.ParseIP(convertUint32ToDotNotation(ent.IpAddr))
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    srcIp,
		DstIP:    dstIp,
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       ent.IfMacAddr,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buffer, options, &ethLayer, &ipLayer, gopacket.Payload(ospf))
	//server.logger.Debug(fmt.Sprintln("buffer: ", buffer))
	lsaAck := buffer.Bytes()
	//server.logger.Debug(fmt.Sprintln("Send  Ack: ", lsaAck))

	return lsaAck
}

func (server *OSPFV2Server) ProcessRxLSAAckPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, nbrKey NbrConfKey) error {

	link_ack := newNbrLsaAckMsg()
	headers_len := ospfHdrMd.Pktlen - OSPF_HEADER_SIZE
	if headers_len >= 20 && headers_len < ospfHdrMd.Pktlen {
		server.logger.Debug("LSAACK: LSA headers length ", headers_len)
		num_headers := int(headers_len / 20)
		server.logger.Debug("LSAACK: Received ", num_headers, " LSA headers.")
		header_byte := make([]byte, num_headers*OSPF_LSA_HEADER_SIZE)
		var start_index uint32
		var lsa_header ospfLSAHeader
		for i := 0; i < num_headers; i++ {
			start_index = uint32(i * OSPF_LSA_HEADER_SIZE)
			copy(header_byte, data[start_index:start_index+20])
			lsa_header = decodeLSAHeader(header_byte)
			server.logger.Debug("LSAACK: Header decoded ",
				"ls_age:ls_type:link_state_id:adv_rtr:ls_seq:ls_checksum ",
				lsa_header.ls_age, lsa_header.ls_type, lsa_header.link_state_id,
				lsa_header.adv_router_id, lsa_header.ls_sequence_num,
				lsa_header.ls_checksum)
			link_ack.lsa_headers = append(link_ack.lsa_headers, lsa_header)
		}
	}
	link_ack.nbrKey = nbrKey
	server.NbrConfData.nbrLsaAckEventCh <- *link_ack
	return nil
}

func (server *OSPFV2Server) DecodeLSAAck(msg NbrLsaAckMsg) {
	server.logger.Debug("LSAACK: Received LSA ACK pkt ", msg)
	nbr, exists := server.NbrConfMap[msg.nbrKey]
	if !exists {
		server.logger.Debug("LSAACK: Nbr doesnt exist", msg.nbrKey)
		return
	}
	discard := server.lsaAckPacketDiscardCheck(nbr)
	if discard {
		return
	}
	/* process each LSA and update request list */
	for index := range msg.lsa_headers {
		req_list := nbr.NbrReqList
		for in := range req_list {
			if req_list[in].link_state_id == msg.lsa_headers[index].link_state_id {
				/* invalidate from request list */
				/* TODO delete req list element */
			}
			/* update the reTxList */

		}
	}
}

/*
Link state request packet
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |   Version #   |       3       |         Packet length         |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          Router ID                            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                           Area ID                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |           Checksum            |             AuType            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Authentication                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          LS type                              |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                       Link State ID                           |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |
*/

/*@fn ProcessRxLSAReqPkt
Send Lsa req packet meta data to Rx packet thread
*/
func (server *OSPFV2Server) ProcessRxLSAReqPkt(data []byte, ospfHdrMd *OspfHdrMetadata, ipHdrMd *IpHdrMetadata, nbrKey NbrConfKey) error {
	//server.logger.Debug(fmt.Sprintln("LSAREQ: Received lsa req with length ", ospfHdrMd.pktlen))
	lsa_req := decodeLSAReqPkt(data, ospfHdrMd.Pktlen)

	lsa_req_msg := NbrLsaReqMsg{
		nbrKey:    nbrKey,
		lsa_slice: lsa_req,
	}
	// send the req list to Nbr
	server.logger.Debug("LSAREQ: Decoded LSA packet - ", lsa_req_msg)
	server.NbrConfData.neighborLSAReqEventCh <- lsa_req_msg
	return nil
}

/*@fn processLSAReqEvent
Process message for lsa req. Unicast LSA to the neighbor if needed.
*/

func (server *OSPFV2Server) ProcessLsaReq(msg NbrLsaReqMsg) {
	server.logger.Debug("LSAREQ: Receieved lsa_req packet for nbr ", msg.nbrKey, " data ", msg.lsa_slice)
	nbrConf, exists := server.NbrConfMap[msg.nbrKey]
	if exists {
		intf := server.IntfConfMap[nbrConf.IntfKey]
		for index := range msg.lsa_slice {
			req := msg.lsa_slice[index]
			lsid := convertUint32ToDotNotation(req.link_state_id)
			adv_router := convertUint32ToDotNotation(req.adv_router_id)
			isDiscard := server.lsaReqPacketDiscardCheck(nbrConf, req)
			server.logger.Debug("LSAREQ: adv_router ", adv_router, " lsid ", lsid, " discard ", isDiscard)
			if !isDiscard {
				areaid := intf.AreaId
				server.generateLsaUpdUnicast(req, msg.nbrKey, areaid)
				server.logger.Debug("LSAREQ: send LSAUPD . adv_router  ", adv_router, " lsid ", lsid, " discard ", isDiscard)
			} else {
				server.logger.Debug("LSAREQ: DONT flood . adv_router  ", adv_router, " lsid ", lsid, " discard ", isDiscard)
			}
		} // enf of for slice
	} // end of exists
}

/*
In response to the LSA request below API generates message
for flooding which sends unicast LSA update packet */
func (server *OSPFV2Server) generateLsaUpdUnicast(req ospfLSAReq, nbrKey NbrConfKey, areaid uint32) {
	lsa_key := NewLsaKey()
	var lsa_pkt []byte
	flood := false

	lsa_key.AdvRouter = req.adv_router_id
	lsa_key.LSId = req.link_state_id
	lsa_key.LSType = uint8(req.ls_type)
	server.logger.Debug("LSAREQ: Generate LSA unicast for LSA type ",
		req.ls_type, " linkid ", req.link_state_id, " areaid ", areaid)
	switch lsa_key.LSType {
	case RouterLSA:
		drlsa, ret := server.getRouterLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeRouterLsa(drlsa, *lsa_key)
			flood = true
		} else {
			server.logger.Debug("LSAREQ: router lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType)
		}
	case NetworkLSA:
		dnlsa, ret := server.getNetworkLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeNetworkLsa(dnlsa, *lsa_key)
			flood = true
		} else {
			server.logger.Debug("LSAREQ: Network lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType)
		}
	case Summary3LSA, Summary4LSA:
		dslsa, ret := server.getSummaryLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeSummaryLsa(dslsa, *lsa_key)
			flood = true
		} else {
			server.logger.Debug("LSAREQ: Summary lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType)
		}
	case ASExternalLSA:
		dalsa, ret := server.getASExternalLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeASExternalLsa(dalsa, *lsa_key)
			flood = true
		} else {
			server.logger.Debug("LSAREQ: AS external lsa not fount. lsaid ",
				req.link_state_id, " lstype ", lsa_key.LSType, " adv_router ", lsa_key.AdvRouter, " areaid ", areaid)
		}
	}
	lsid := req.link_state_id
	router_id := req.adv_router_id

	server.logger.Debug("LSAUPD: lsid ", lsid, " router_id ", router_id, " flood ", flood)

	if flood {
		checksumOffset := uint16(14)
		checkSum := computeFletcherChecksum(lsa_pkt[2:], checksumOffset)
		binary.BigEndian.PutUint16(lsa_pkt[16:18], checkSum)
		flood_msg := NbrToFloodMsg{}
		flood_msg.MsgType = LSA_FLOOD_INTF
		flood_msg.NbrKey = nbrKey
		flood_msg.LsaType = lsa_key.LSType
		flood_msg.LsaPkt = make([]byte, len(lsa_pkt))
		copy(flood_msg.LsaPkt, lsa_pkt)

		server.MessagingChData.NbrFSMToFloodChData.LsaFloodCh <- flood_msg
	}
}

func (server *OSPFV2Server) lsaReqPacketDiscardCheck(nbrConf NbrConf, req ospfLSAReq) bool {
	if nbrConf.State < NbrExchange {
		server.logger.Debug("LSAREQ: Discard .. Nbrstate (expected less than exchange)", nbrConf.State)
		return true
	}
	/* TODO
	check the router DB if packet needs to be updated.
	if not found in LSDB generate LSAReqEvent */

	return false
}

func (server *OSPFV2Server) lsaAckPacketDiscardCheck(nbrConf NbrConf) bool {
	if nbrConf.State < NbrExchange {
		server.logger.Debug("LSAACK: Discard .. Nbrstate (expected less than exchange)", nbrConf.State)
		return true
	}
	/* TODO
	check the router DB if packet needs to be updated.
	if not found in LSDB generate LSAReqEvent */

	return false
}

/*@fn lsaReTxTimerCheck
 */
func (server *OSPFV2Server) lsaReTxTimerCheck(nbrKey NbrConfKey) {
	var lsa_re_tx_check_func func()
	lsa_re_tx_check_func = func() {
		nbrConf, exist := server.NbrConfMap[nbrKey]
		if !exist {
			server.logger.Err("NbrConf does not exist . No retx timer set ", nbrKey)
			return
		}
		server.logger.Debug("LSARETIMER: Check for rx. Nbr ", nbrKey)
		// check for retx list
		re_list := nbrConf.NbrRetxList

		if len(re_list) > 0 {
			// retransmit packet
			server.logger.Debug("LSATIMER: Send the retx packets. ")
		}
	}
	_, exists := server.NbrConfMap[nbrKey]
	if exists {
		nbrConf := server.NbrConfMap[nbrKey]
		nbrConf.NbrLsaRxTimer = time.AfterFunc(NBR_DBD_RTX_INTERVAL, lsa_re_tx_check_func)
		//op := NBRUPD
		//server.sendNeighborConf(nbrKey, nbrConf, NbrMsgType(op))
	}
}

func (server *OSPFV2Server) processTxLsaAck(lsa_data NbrAckTxMsg) {
	ack_len := len(lsa_data.lsa_headers_byte)
	total_ack := ack_len / OSPF_LSA_ACK_SIZE
	if total_ack < 0 {
		server.logger.Debug("TX ACK: malformed message. total_ack ", total_ack, " pkt_size ", ack_len)
		return
	}
	nbrConf, exists := server.NbrConfMap[lsa_data.nbrKey]
	if !exists {
		server.logger.Warning("TX ACK: neighbor doesnt exist  to send ack ", lsa_data.nbrKey)
		return
	}
	intf, valid := server.IntfConfMap[nbrConf.IntfKey]
	if !valid {
		server.logger.Err("Nbr : Intf does not exist. No ack send", nbrConf.IntfKey)
		return
	}

	dstMac := nbrConf.NbrMac
	dstIp := net.ParseIP(convertUint32ToDotNotation(nbrConf.NbrIP))
	pkt := server.BuildLSAAckPkt(nbrConf.IntfKey, intf, nbrConf, dstMac, dstIp,
		ack_len, lsa_data.lsa_headers_byte)
	// send ack over the pcap.
	server.logger.Debug("ACK SEND: ", pkt)
	server.SendOspfPkt(nbrConf.IntfKey, pkt)

}
