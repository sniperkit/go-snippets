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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/ospf/config"
	"math"
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

type ospfNeighborLSAreqMsg struct {
	lsa_slice []ospfLSAReq
	nbrKey    NeighborConfKey
}

type ospfNeighborLSDBMsg struct {
	areaId uint32
	data   []byte
}

type ospfNeighborLSAAckMsg struct {
	lsa_headers []ospfLSAHeader
	nbrKey      NeighborConfKey
}

/* ACK message uses the LSA header byte
  received from LSA UPD packet. Therefore
new message type to tx message is added
*/
type ospfNeighborAckTxMsg struct {
	lsa_headers_byte []byte
	nbrKey           NeighborConfKey
}

func newospfNeighborAckTxMsg() *ospfNeighborAckTxMsg {
	return &ospfNeighborAckTxMsg{}
}

func NewospfNeighborLSDBMsg() *ospfNeighborLSDBMsg {
	return &ospfNeighborLSDBMsg{}
}

func newospfNeighborLSAAckMsg() *ospfNeighborLSAAckMsg {
	return &ospfNeighborLSAAckMsg{}
}

type ospfNeighborLSAUpdMsg struct {
	nbrKey NeighborConfKey
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
	adv_router_id uint32, ls_sequence_num uint32, ls_checksum uint16, ls_len uint16) ospfLSAHeader {

	var lsa_header ospfLSAHeader
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

func (server *OSPFServer) EncodeLSAReqPkt(intfKey IntfConfKey, ent IntfConf,
	nbrConf OspfNeighborEntry, lsa_req_pkt []ospfLSAReq, dstMAC net.HardwareAddr) (data []byte) {
	ospfHdr := OSPFHeader{
		ver:      OSPF_VERSION_2,
		pktType:  uint8(LSRequestType),
		pktlen:   0,
		routerId: server.ospfGlobalConf.RouterId,
		areaId:   ent.IfAreaId,
		chksum:   0,
		authType: ent.IfAuthType,
	}

	lsaDataEnc := encodeLSAReq(lsa_req_pkt)
	ospfPktlen := OSPF_HEADER_SIZE
	ospfPktlen = ospfPktlen + len(lsaDataEnc)

	ospfHdr.pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	server.logger.Info(fmt.Sprintln("ospfEncHdr:", ospfEncHdr))
	server.logger.Info(fmt.Sprintln("lsa Pkt:", lsaDataEnc))

	ospf := append(ospfEncHdr, lsaDataEnc...)
	server.logger.Info(fmt.Sprintln("OSPF LSA REQ:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	copy(ospf[16:24], ent.IfAuthKey)

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.pktlen
	var dstIp net.IP
	if ent.IfType == config.NumberedP2P {
		dstIp = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(config.McastMAC)
	} else {
		dstIp = nbrConf.OspfNbrIPAddr
	}
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    ent.IfIpAddr,
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
	server.logger.Info(fmt.Sprintln("lsaPkt: ", lsaPkt))

	return lsaPkt

}

func (server *OSPFServer) BuildAndSendLSAReq(nbrId NeighborConfKey, nbrConf OspfNeighborEntry) (curr_index uint8) {
	/* calculate max no of requests that can be added
	for req packet */

	var add_items uint8

	var req ospfLSAReq
	var i uint8

	msg := ospfNeighborLSAreqMsg{}
	msg.lsa_slice = []ospfLSAReq{}
	msg.nbrKey = nbrId

	reqlist := ospfNeighborRequest_list[nbrId]
	req_list_items := uint8(len(reqlist)) - nbrConf.ospfNbrLsaReqIndex
	max_req := calculateMaxLsaReq()
	if max_req > req_list_items {
		add_items = req_list_items
		nbrConf.ospfNbrLsaReqIndex = uint8(len(reqlist))

	} else {
		add_items = uint8(max_req)
		nbrConf.ospfNbrLsaReqIndex += max_req
	}
	server.logger.Info(fmt.Sprintln("LSAREQ: nbrIndex ",
		nbrConf.ospfNbrLsaReqIndex, " add_items ", add_items, " req_list len ", len(reqlist)))
	index := nbrConf.ospfNbrLsaReqIndex
	for i = 0; i < uint8(len(reqlist)); i++ {
		req.ls_type = uint32(reqlist[i].lsa_headers.ls_type)
		req.link_state_id = reqlist[i].lsa_headers.link_state_id
		req.adv_router_id = reqlist[i].lsa_headers.adv_router_id
		nbrConf.req_list_mutex.Lock()
		msg.lsa_slice = append(msg.lsa_slice, req)
		reqlist[i].valid = false
		nbrConf.req_list_mutex.Unlock()
		/* update LSA Retx list */
		reTxNbr := newospfNeighborRetx()
		reTxNbr.lsa_headers = reqlist[i].lsa_headers
		reTxNbr.valid = true
		nbrConf.retx_list_mutex.Lock()
		reTxList := ospfNeighborRetx_list[nbrId]
		reTxList = append(reTxList, reTxNbr)
		nbrConf.retx_list_mutex.Unlock()

		lsid := convertUint32ToIPv4(req.link_state_id)
		adv_rtr := convertUint32ToIPv4(req.adv_router_id)
		server.logger.Info(fmt.Sprintln("LSA request: Send req to nbr ", nbrId.IPAddr,
			" lsid ", lsid, " rtrid ", adv_rtr, " lstype ", req.ls_type))
	}
	server.logger.Info(fmt.Sprintln("LSA request: total requests out, req_list_len, current req_list_index ", add_items, len(msg.lsa_slice), nbrConf.ospfNbrLsaReqIndex))
	server.logger.Info(fmt.Sprintln("LSA request: lsa_req", msg.lsa_slice))

	if len(msg.lsa_slice) != 0 {
		server.ospfNbrLsaReqSendCh <- msg
		index += add_items
	}
	return index
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

func (server *OSPFServer) BuildLsaUpdPkt(intfKey IntfConfKey, ent IntfConf,
	dstMAC net.HardwareAddr, dstIp net.IP, lsa_pkt_size int, lsaUpdEnc []byte) (data []byte) {
	ospfHdr := OSPFHeader{
		ver:      OSPF_VERSION_2,
		pktType:  uint8(LSUpdateType),
		pktlen:   0,
		routerId: server.ospfGlobalConf.RouterId,
		areaId:   ent.IfAreaId,
		chksum:   0,
		authType: ent.IfAuthType,
		//authKey:        ent.IfAuthKey,
	}

	ospfPktlen := OSPF_HEADER_SIZE

	ospfPktlen = ospfPktlen + len(lsaUpdEnc)
	ospfHdr.pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	//server.logger.Info(fmt.Sprintln("ospfEncHdr:", ospfEncHdr))

	//server.logger.Info(fmt.Sprintln("LSA upd Pkt:", lsaUpdEnc))

	ospf := append(ospfEncHdr, lsaUpdEnc...)
	//server.logger.Info(fmt.Sprintln("OSPF LSA UPD:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	copy(ospf[16:24], ent.IfAuthKey)

	if ent.IfType == config.NumberedP2P {
		dstIp = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(config.McastMAC)
	}

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.pktlen
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    ent.IfIpAddr,
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
	//server.logger.Info(fmt.Sprintln("buffer: ", buffer))
	lsaUpd := buffer.Bytes()
	//server.logger.Info(fmt.Sprintln("flood Pkt: ", lsaUpd))

	return lsaUpd

}

func (server *OSPFServer) ProcessRxLsaUpdPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, key IntfConfKey) error {

	routerId := convertIPInByteToString(ospfHdrMd.routerId)
	ipaddr := net.IPv4(ipHdrMd.srcIP[0], ipHdrMd.srcIP[1], ipHdrMd.srcIP[2], ipHdrMd.srcIP[3])
	ospfNbrConfKey := NeighborConfKey{
		IPAddr:  config.IpAddress(ipaddr.String()),
		IntfIdx: key.IntfIdx,
	}

	msg := ospfNeighborLSAUpdMsg{
		nbrKey: ospfNbrConfKey,
		areaId: ospfHdrMd.areaId,
		data:   data,
	}

	server.neighborLSAUpdEventCh <- msg
	server.logger.Info(fmt.Sprintln("LSA update: Received LSA update with router_id , lentgh ", routerId, ospfHdrMd.pktlen))
	//	server.logger.Info(fmt.Sprintln("LSA update: pkt byte[]: ", data))
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

func (server *OSPFServer) DecodeLSAUpd(msg ospfNeighborLSAUpdMsg) {
	nbr, exists := server.NeighborConfigMap[msg.nbrKey]
	op := LsdbNoAction
	discard := true
	if !exists {
		return
	}

	lsop := uint8(LSASELFLOOD)
	intf := server.IntfConfMap[nbr.intfConfKey]
	lsa_max_age := false
	discard = server.lsaUpdDiscardCheck(nbr, msg.data)
	if discard {
		server.logger.Err(fmt.Sprintln("LSAUPD: Discard. nbr ", msg.nbrKey))
		return
	}

	no_lsa := binary.BigEndian.Uint32(msg.data[0:4])
	server.logger.Info(fmt.Sprintln("LSAUPD: Nbr, No of LSAs ", msg.nbrKey, no_lsa, "  len  ", len(msg.data)))
	//server.logger.Info(fmt.Sprintln("LSUPD:LSA pkt ", msg.data))
	lsa_header := NewLsaHeader()
	/* decode each LSA and send to lsdb
	 */
	index := 4
	end_index := 0
	lsa_header_byte := make([]byte, OSPF_LSA_HEADER_SIZE)
	for i := 0; i < int(no_lsa); i++ {

		decodeLsaHeader(msg.data[index:index+OSPF_LSA_HEADER_SIZE], lsa_header)
		copy(lsa_header_byte, msg.data[index:index+OSPF_LSA_HEADER_SIZE])
		server.logger.Info(fmt.Sprintln("LSAUPD: lsaheader decoded adv_rter ", lsa_header.Adv_router,
			" linkid ", lsa_header.LinkId, " lsage ", lsa_header.LSAge,
			" checksum ", lsa_header.LSChecksum, " seq num ", lsa_header.LSSequenceNum,
			" LSTYPE ", lsa_header.LSType,
			" len ", lsa_header.length))
		end_index = int(lsa_header.length) + index /* length includes data + header */
		if lsa_header.LSAge == LSA_MAX_AGE {
			lsa_max_age = true
		}
		/* send message to lsdb */
		lsdb_msg := NewLsdbUpdateMsg()
		lsdb_msg.AreaId = msg.areaId
		lsdb_msg.Data = make([]byte, end_index-i)
		copy(lsdb_msg.Data, msg.data[index:end_index])
		valid := validateChecksum(lsdb_msg.Data)
		if !valid {
			server.logger.Info(fmt.Sprintln("LSAUPD: Invalid checksum. Nbr",
				server.NeighborConfigMap[msg.nbrKey]))
			//continue
		}
		lsa_key := NewLsaKey()

		switch lsa_header.LSType {
		case RouterLSA:
			rlsa := NewRouterLsa()
			decodeRouterLsa(lsdb_msg.Data, rlsa, lsa_key)

			drlsa, ret := server.getRouterLsaFromLsdb(msg.areaId, *lsa_key)
			discard, op = server.sanityCheckRouterLsa(*rlsa, drlsa, nbr, intf, ret, lsa_max_age)

		case NetworkLSA:
			nlsa := NewNetworkLsa()
			decodeNetworkLsa(lsdb_msg.Data, nlsa, lsa_key)
			dnlsa, ret := server.getNetworkLsaFromLsdb(msg.areaId, *lsa_key)
			discard, op = server.sanityCheckNetworkLsa(*lsa_key, *nlsa, dnlsa, nbr, intf, ret, lsa_max_age)

		case Summary3LSA, Summary4LSA:
			server.logger.Info(fmt.Sprintln("Received summary Lsa Packet :", lsdb_msg.Data))
			slsa := NewSummaryLsa()
			decodeSummaryLsa(lsdb_msg.Data, slsa, lsa_key)
			server.logger.Info(fmt.Sprintln("Decoded summary Lsa Packet :", slsa))
			dslsa, ret := server.getSummaryLsaFromLsdb(msg.areaId, *lsa_key)
			discard, op = server.sanityCheckSummaryLsa(*slsa, dslsa, nbr, intf, ret, lsa_max_age)

		case ASExternalLSA:
			alsa := NewASExternalLsa()
			decodeASExternalLsa(lsdb_msg.Data, alsa, lsa_key)
			dalsa, ret := server.getASExternalLsaFromLsdb(msg.areaId, *lsa_key)
			discard, op = server.sanityCheckASExternalLsa(*alsa, dalsa, nbr, intf, intf.IfAreaId, ret, lsa_max_age)

		}
		lsid := convertUint32ToIPv4(lsa_header.LinkId)
		router_id := convertUint32ToIPv4(lsa_header.Adv_router)

		self_gen := false
		self_gen = server.selfGenLsaCheck(*lsa_key)
		if self_gen {
			server.logger.Info(fmt.Sprintln("LSAUPD: discard . Received self generated. ", lsa_key))

		}

		if !discard && !self_gen && op == FloodLsa {
			server.logger.Info(fmt.Sprintln("LSAUPD: add to lsdb lsid ", lsid, " router_id ", router_id, " lstype ", lsa_header.LSType))
			lsdb_msg.MsgType = LsdbAdd
			server.LsdbUpdateCh <- *lsdb_msg

		}

		flood_pkt := ospfFloodMsg{
			nbrKey: msg.nbrKey,
			areaId: msg.areaId,
			lsType: lsa_header.LSType,
			linkid: lsa_header.LinkId,
			lsOp:   lsop,
		}
		flood_pkt.pkt = make([]byte, end_index-index)
		copy(flood_pkt.pkt, lsdb_msg.Data)
		if lsop != LSASUMMARYFLOOD && !self_gen { // for ABR summary lsa is flooded after LSDB/SPF changes are done.
			server.ospfNbrLsaUpdSendCh <- flood_pkt
		}

		/* send ACK */
		lsaAckMsg := newospfNeighborAckTxMsg()
		lsaAckMsg.lsa_headers_byte = append(lsaAckMsg.lsa_headers_byte, lsa_header_byte...)
		lsaAckMsg.nbrKey = msg.nbrKey
		server.logger.Info(fmt.Sprintln("ACK TX: nbr ", msg.nbrKey, " ack ", lsaAckMsg.lsa_headers_byte))
		server.ospfNbrLsaAckSendCh <- *lsaAckMsg

		index = end_index
		server.UpdateNeighborList(msg.nbrKey)

	}
}

func (server *OSPFServer) selfGenLsaCheck(key LsaKey) bool {
	rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
	if key.AdvRouter == rtr_id {
		return true
	}
	return false
}
func (server *OSPFServer) lsaUpdDiscardCheck(nbrConf OspfNeighborEntry, data []byte) bool {
	if nbrConf.OspfNbrState < config.NbrExchange {
		server.logger.Info(fmt.Sprintln("LSAUPD: Discard .. Nbrstate (expected less than exchange)", nbrConf.OspfNbrState))
		return true
	}

	return false
}
func (server *OSPFServer) lsAgeCheck(intf IntfConfKey, lsa_max_age bool, exist int) bool {

	send_ack := true
	/*
				if the LSA's LS age is equal to MaxAge, and there is
			    currently no instance of the LSA in the router's link state
			    database, and none of router's neighbors are in states Exchange
			    or Loading, then take the following actions: a) Acknowledge the
			    receipt of the LSA by sending a Link State Acknowledgment packet
			    back to the sending neighbor (see Section 13.5), and b) Discard
			    the LSA and examine the next LSA (if any) listed in the Link
		        State Update packet.
	*/
	data := ospfIntfToNbrMap[intf]
	for _, nbrKey := range data.nbrList {
		nbr := server.NeighborConfigMap[nbrKey]
		if nbr.OspfNbrState == config.NbrExchange || nbr.OspfNbrState == config.NbrLoading {
			continue
		} else {
			send_ack = false
		}
	}
	if send_ack && exist == LsdbEntryNotFound && lsa_max_age {
		return true
	}
	return false
}

func (server *OSPFServer) sanityCheckRouterLsa(rlsa RouterLsa, drlsa RouterLsa, nbr OspfNeighborEntry, intf IntfConf, exist int, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.intfConfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Router LSA Discard. link details", rlsa.LinkDetails, " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(rlsa.LsaMd, drlsa.LsaMd)
		// TODO check if lsa is installed before MinLSArrival
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			server.logger.Info(fmt.Sprintln("LSAUPD: Router LSA Discard.Already present in lsdb. link details", rlsa.LinkDetails, " nbr ", nbr))
			discard = true
			op = LsdbNoAction
		}
	}

	return discard, op
}

func (server *OSPFServer) sanityCheckNetworkLsa(lsaKey LsaKey, nlsa NetworkLsa, dnlsa NetworkLsa, nbr OspfNeighborEntry, intf IntfConf, exist int, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.intfConfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Network LSA Discard. ", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(nlsa.LsaMd, dnlsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			discard = true
			op = LsdbNoAction
		}
	}
	//if i am DR and receive nw LSA from neighbor discard it.
	rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
	if intf.IfDRtrId == rtr_id {
		nbrIp := convertAreaOrRouterIdUint32(nbr.OspfNbrIPAddr.String())
		if lsaKey.LSId == nbrIp {
			server.logger.Info(fmt.Sprintln("DISCARD: I am dr. received nw LSA from nbr . LSA id ", nbr.OspfNbrIPAddr))
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}

func (server *OSPFServer) sanityCheckSummaryLsa(slsa SummaryLsa, dslsa SummaryLsa, nbr OspfNeighborEntry, intf IntfConf, exist int, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.intfConfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Summary LSA Discard. ", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(slsa.LsaMd, dslsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			server.logger.Info(fmt.Sprintln("LSAUPD: Discard Summary LSA slsa from nbr"))
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}

func (server *OSPFServer) sanityCheckASExternalLsa(alsa ASExternalLsa, dalsa ASExternalLsa, nbr OspfNeighborEntry, intf IntfConf, areaid []byte, exist int, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	// TODO Reject this lsa if area is configured as stub area.
	send_ack := server.lsAgeCheck(nbr.intfConfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: As external LSA Discard.", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(alsa.LsaMd, dalsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}

func validateChecksum(data []byte) bool {

	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		//server.logger.Err("LSAUPD: Invalid Router LSA Checksum")
		return false
	}
	return true
}

func (server *OSPFServer) validateLsaIsNew(rlsamd LsaMetadata, dlsamd LsaMetadata) bool {
	if rlsamd.LSSequenceNum > dlsamd.LSSequenceNum {
		server.logger.Info(fmt.Sprintln("LSA: received lsseq num > db seq num. "))
		return true
	}
	if rlsamd.LSChecksum > dlsamd.LSChecksum {
		server.logger.Info(fmt.Sprintln("LSA: received lsa checksum > db chceksum "))
		return true
	}
	if rlsamd.LSAge == LSA_MAX_AGE {
		server.logger.Info(fmt.Sprintln("LSA: LSA is maxage "))
		return true
	}
	age_diff := math.Abs(float64(rlsamd.LSAge - dlsamd.LSAge))
	if age_diff > float64(LSA_MAX_AGE_DIFF) &&
		rlsamd.LSAge < rlsamd.LSAge {
		return true
	}
	/* Debug further - currently it doesnt return true for latest LSA */
	return true
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

//func (server *OSPFServer) encodeLSAAck
func (server *OSPFServer) BuildLSAAckPkt(intfKey IntfConfKey, ent IntfConf,
	nbrConf OspfNeighborEntry, dstMAC net.HardwareAddr, dstIp net.IP, lsa_pkt_size int, lsaAckEnc []byte) (data []byte) {
	ospfHdr := OSPFHeader{
		ver:      OSPF_VERSION_2,
		pktType:  uint8(LSAckType),
		pktlen:   0,
		routerId: server.ospfGlobalConf.RouterId,
		areaId:   ent.IfAreaId,
		chksum:   0,
		authType: ent.IfAuthType,
	}

	ospfPktlen := OSPF_HEADER_SIZE
	ospfPktlen = ospfPktlen + lsa_pkt_size
	ospfHdr.pktlen = uint16(ospfPktlen)
	server.logger.Info(fmt.Sprintln("LSAACK : packet legth header(24) + ack ", ospfPktlen))
	ospfEncHdr := encodeOspfHdr(ospfHdr)
	//server.logger.Info(fmt.Sprintln("ospfEncHdr:", ospfEncHdr))

	//server.logger.Info(fmt.Sprintln("LSA upd Pkt:", lsaAckEnc))

	ospf := append(ospfEncHdr, lsaAckEnc...)
	//server.logger.Info(fmt.Sprintln("OSPF LSA ACK:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	copy(ospf[16:24], ent.IfAuthKey)

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.pktlen
	if ent.IfType == config.NumberedP2P {
		dstIp = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(config.McastMAC)
	}
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    ent.IfIpAddr,
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
	//server.logger.Info(fmt.Sprintln("buffer: ", buffer))
	lsaAck := buffer.Bytes()
	//server.logger.Info(fmt.Sprintln("Send  Ack: ", lsaAck))

	return lsaAck
}

func (server *OSPFServer) ProcessRxLSAAckPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, key IntfConfKey) error {

	link_ack := newospfNeighborLSAAckMsg()
	headers_len := ospfHdrMd.pktlen - OSPF_HEADER_SIZE
	if headers_len >= 20 && headers_len < ospfHdrMd.pktlen {
		server.logger.Info(fmt.Sprintln("LSAACK: LSA headers length ", headers_len))
		num_headers := int(headers_len / 20)
		server.logger.Info(fmt.Sprintln("LSAACK: Received ", num_headers, " LSA headers."))
		header_byte := make([]byte, num_headers*OSPF_LSA_HEADER_SIZE)
		var start_index uint32
		var lsa_header ospfLSAHeader
		for i := 0; i < num_headers; i++ {
			start_index = uint32(i * OSPF_LSA_HEADER_SIZE)
			copy(header_byte, data[start_index:start_index+20])
			lsa_header = decodeLSAHeader(header_byte)
			server.logger.Info(fmt.Sprintln("LSAACK: Header decoded ",
				"ls_age:ls_type:link_state_id:adv_rtr:ls_seq:ls_checksum ",
				lsa_header.ls_age, lsa_header.ls_type, lsa_header.link_state_id,
				lsa_header.adv_router_id, lsa_header.ls_sequence_num,
				lsa_header.ls_checksum))
			link_ack.lsa_headers = append(link_ack.lsa_headers, lsa_header)
		}
	}
	ipaddr := convertByteToOctetString(ipHdrMd.srcIP)
	ospfNbrConfKey := NeighborConfKey{
		IPAddr:  config.IpAddress(ipaddr),
		IntfIdx: key.IntfIdx,
	}
	link_ack.nbrKey = ospfNbrConfKey
	server.neighborLSAACKEventCh <- *link_ack
	return nil
}

func (server *OSPFServer) DecodeLSAAck(msg ospfNeighborLSAAckMsg) {
	server.logger.Info(fmt.Sprintln("LSAACK: Received LSA ACK pkt ", msg))
	nbr, exists := server.NeighborConfigMap[msg.nbrKey]
	if !exists {
		server.logger.Info(fmt.Sprintln("LSAACK: Nbr doesnt exist", msg.nbrKey))
		return
	}
	discard := server.lsaAckPacketDiscardCheck(nbr)
	if discard {
		return
	}
	/* process each LSA and update request list */
	for index := range msg.lsa_headers {
		req_list := ospfNeighborRequest_list[msg.nbrKey]
		reTx_list := ospfNeighborRetx_list[msg.nbrKey]
		for in := range req_list {
			if req_list[in].lsa_headers.link_state_id == msg.lsa_headers[index].link_state_id {
				/* invalidate from request list */
				req := newospfNeighborReq()
				req.lsa_headers = msg.lsa_headers[index]

				nbr.req_list_mutex.Lock()
				req_list[in].valid = false
				nbr.req_list_mutex.Unlock()
			}
			/* update the reTxList */
			for in = range reTx_list {
				if reTx_list[in].lsa_headers.link_state_id == msg.lsa_headers[index].link_state_id {
					nbr.retx_list_mutex.Lock()
					reTx_list[in].valid = false
					nbr.retx_list_mutex.Unlock()
				}
			}

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
func (server *OSPFServer) ProcessRxLSAReqPkt(data []byte, ospfHdrMd *OspfHdrMetadata, ipHdrMd *IpHdrMetadata, key IntfConfKey) error {
	//server.logger.Info(fmt.Sprintln("LSAREQ: Received lsa req with length ", ospfHdrMd.pktlen))
	lsa_req := decodeLSAReqPkt(data, ospfHdrMd.pktlen)
	ipaddr := net.IPv4(ipHdrMd.srcIP[0], ipHdrMd.srcIP[1], ipHdrMd.srcIP[2], ipHdrMd.srcIP[3])
	ospfNbrConfKey := NeighborConfKey{
		IPAddr:  config.IpAddress(ipaddr.String()),
		IntfIdx: key.IntfIdx,
	}

	lsa_req_msg := ospfNeighborLSAreqMsg{
		nbrKey:    ospfNbrConfKey,
		lsa_slice: lsa_req,
	}
	// send the req list to Nbr
	server.logger.Info(fmt.Sprintln("LSAREQ: Decoded LSA packet - ", lsa_req_msg))
	server.neighborLSAReqEventCh <- lsa_req_msg
	return nil
}

/*@fn processLSAReqEvent
Process message for lsa req. Unicast LSA to the neighbor if needed.
*/

func (server *OSPFServer) DecodeLSAReq(msg ospfNeighborLSAreqMsg) {
	server.logger.Info(fmt.Sprintln("LSAREQ: Receieved lsa_req packet for nbr ", msg.nbrKey, " data ", msg.lsa_slice))
	nbrConf, exists := server.NeighborConfigMap[msg.nbrKey]
	if exists {
		intf := server.IntfConfMap[nbrConf.intfConfKey]
		for index := range msg.lsa_slice {
			req := msg.lsa_slice[index]
			lsid := convertUint32ToIPv4(req.link_state_id)
			adv_router := convertUint32ToIPv4(req.adv_router_id)
			isDiscard := server.lsaReqPacketDiscardCheck(nbrConf, req)
			server.logger.Info(fmt.Sprintln("LSAREQ: adv_router ", adv_router, " lsid ", lsid, " discard ", isDiscard))
			if !isDiscard {
				areaid := convertIPv4ToUint32(intf.IfAreaId)
				server.generateLsaUpdUnicast(req, msg.nbrKey, areaid)
				server.logger.Info(fmt.Sprintln("LSAREQ: send LSAUPD . adv_router  ", adv_router, " lsid ", lsid, " discard ", isDiscard))
			} else {
				server.logger.Info(fmt.Sprintln("LSAREQ: DONT flood . adv_router  ", adv_router, " lsid ", lsid, " discard ", isDiscard))
			}
		} // enf of for slice
	} // end of exists
}

func (server *OSPFServer) generateLsaUpdUnicast(req ospfLSAReq, nbrKey NeighborConfKey, areaid uint32) {
	lsa_key := NewLsaKey()
	nbrConf := server.NeighborConfigMap[nbrKey]
	var lsa_pkt []byte
	flood := false

	lsa_key.AdvRouter = req.adv_router_id
	lsa_key.LSId = req.link_state_id
	lsa_key.LSType = uint8(req.ls_type)
	server.logger.Info(fmt.Sprintln("LSAREQ: Generate LSA unicast for LSA type ",
		req.ls_type, " linkid ", req.link_state_id, " areaid ", areaid))
	switch lsa_key.LSType {
	case RouterLSA:
		drlsa, ret := server.getRouterLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeRouterLsa(drlsa, *lsa_key)
			flood = true
		} else {
			server.logger.Info(fmt.Sprintln("LSAREQ: router lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType))
		}
	case NetworkLSA:
		dnlsa, ret := server.getNetworkLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeNetworkLsa(dnlsa, *lsa_key)
			flood = true
		} else {
			server.logger.Info(fmt.Sprintln("LSAREQ: Network lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType))
		}
	case Summary3LSA, Summary4LSA:
		dslsa, ret := server.getSummaryLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeSummaryLsa(dslsa, *lsa_key)
			flood = true
		} else {
			server.logger.Info(fmt.Sprintln("LSAREQ: Summary lsa not found. lsaid ", req.link_state_id, " lstype ", lsa_key.LSType))
		}
	case ASExternalLSA:
		dalsa, ret := server.getASExternalLsaFromLsdb(areaid, *lsa_key)
		if ret == LsdbEntryFound {
			lsa_pkt = encodeASExternalLsa(dalsa, *lsa_key)
			flood = true
		} else {
			server.logger.Info(fmt.Sprintln("LSAREQ: AS external lsa not fount. lsaid ",
				req.link_state_id, " lstype ", lsa_key.LSType, " adv_router ", lsa_key.AdvRouter, " areaid ", areaid))
		}
	}
	lsid := convertUint32ToIPv4(req.link_state_id)
	router_id := convertUint32ToIPv4(req.adv_router_id)

	server.logger.Info(fmt.Sprintln("LSAUPD: lsid ", lsid, " router_id ", router_id, " flood ", flood))

	if flood {
		checksumOffset := uint16(14)
		checkSum := computeFletcherChecksum(lsa_pkt[2:], checksumOffset)
		binary.BigEndian.PutUint16(lsa_pkt[16:18], checkSum)
		flood_pkt := ospfFloodMsg{
			nbrKey:  nbrKey,
			intfKey: nbrConf.intfConfKey,
			areaId:  areaid,
			lsType:  uint8(req.ls_type),
			linkid:  req.link_state_id,
			lsOp:    LSAINTF,
		}
		flood_pkt.pkt = make([]byte, len(lsa_pkt))
		copy(flood_pkt.pkt, lsa_pkt)
		server.ospfNbrLsaUpdSendCh <- flood_pkt
	}
}

func (server *OSPFServer) lsaReqPacketDiscardCheck(nbrConf OspfNeighborEntry, req ospfLSAReq) bool {
	if nbrConf.OspfNbrState < config.NbrExchange {
		server.logger.Info(fmt.Sprintln("LSAREQ: Discard .. Nbrstate (expected less than exchange)", nbrConf.OspfNbrState))
		return true
	}
	/* TODO
	check the router DB if packet needs to be updated.
	if not found in LSDB generate LSAReqEvent */

	return false
}

func (server *OSPFServer) lsaAckPacketDiscardCheck(nbrConf OspfNeighborEntry) bool {
	if nbrConf.OspfNbrState < config.NbrExchange {
		server.logger.Info(fmt.Sprintln("LSAACK: Discard .. Nbrstate (expected less than exchange)", nbrConf.OspfNbrState))
		return true
	}
	/* TODO
	check the router DB if packet needs to be updated.
	if not found in LSDB generate LSAReqEvent */

	return false
}

/*
@fn lsaAddCheck
	This API checks if the LSA header received in DBD
	is to be added in the req list.
*/

func (server *OSPFServer) lsaAddCheck(lsaheader ospfLSAHeader,
	nbr OspfNeighborEntry) (result bool) {

	lsa_max_age := false
	intf := server.IntfConfMap[nbr.intfConfKey]
	areaId := convertIPv4ToUint32(intf.IfAreaId)
	if lsaheader.ls_age == LSA_MAX_AGE {
		lsa_max_age = true
	}
	lsa_key := NewLsaKey()
	lsa_key.AdvRouter = lsaheader.adv_router_id
	lsa_key.LSId = lsaheader.link_state_id
	lsa_key.LSType = lsaheader.ls_type
	adv_router := convertUint32ToIPv4(lsa_key.AdvRouter)
	discard := true
	var op uint8
	discard = server.selfGenLsaCheck(*lsa_key)
	if discard {
		server.logger.Info(fmt.Sprintln("DBD: Db received self originated LSA . discard. lsa key ", *lsa_key))
		return false
	}

	switch lsaheader.ls_type {
	case RouterLSA:
		rlsa := NewRouterLsa()
		drlsa, ret := server.getRouterLsaFromLsdb(areaId, *lsa_key)
		discard, op = server.sanityCheckRouterLsa(*rlsa, drlsa, nbr, intf, ret, lsa_max_age)

	case NetworkLSA:
		nlsa := NewNetworkLsa()
		dnlsa, ret := server.getNetworkLsaFromLsdb(areaId, *lsa_key)
		discard, op = server.sanityCheckNetworkLsa(*lsa_key, *nlsa, dnlsa, nbr, intf, ret, lsa_max_age)

	case Summary3LSA, Summary4LSA:
		slsa := NewSummaryLsa()
		dslsa, ret := server.getSummaryLsaFromLsdb(areaId, *lsa_key)
		discard, op = server.sanityCheckSummaryLsa(*slsa, dslsa, nbr, intf, ret, lsa_max_age)

	case ASExternalLSA:
		alsa := NewASExternalLsa()
		dalsa, ret := server.getASExternalLsaFromLsdb(areaId, *lsa_key)
		discard, op = server.sanityCheckASExternalLsa(*alsa, dalsa, nbr, intf, intf.IfAreaId, ret, lsa_max_age)

	}
	if discard {
		server.logger.Info(fmt.Sprintln("DBD: LSA is not added in the request list. Adv router ", adv_router,
			" ls_type ", lsaheader.ls_type, " op ", op))
		return false
	}
	server.logger.Info(fmt.Sprintln("DBD: LSA append to req_list adv_router ", adv_router,
		" Lsid ", lsa_key.LSId, " lstype ", lsa_key.LSType))
	return true
}

/*@fn lsaReTxTimerCheck
 */
func (server *OSPFServer) lsaReTxTimerCheck(nbrKey NeighborConfKey) {
	var lsa_re_tx_check_func func()
	lsa_re_tx_check_func = func() {
		server.logger.Info(fmt.Sprintln("LSARETIMER: Check for rx. Nbr ", nbrKey))
		// check for retx list
		re_list := ospfNeighborRetx_list[nbrKey]
		if len(re_list) > 0 {
			// retransmit packet
			server.logger.Info(fmt.Sprintln("LSATIMER: Send the retx packets. "))
		}
	}
	_, exists := server.NeighborConfigMap[nbrKey]
	if exists {
		nbrConf := server.NeighborConfigMap[nbrKey]
		nbrConf.ospfNeighborLsaRxTimer = time.AfterFunc(RxDBDInterval, lsa_re_tx_check_func)
		//op := NBRUPD
		//server.sendNeighborConf(nbrKey, nbrConf, NbrMsgType(op))
	}
}

func (server *OSPFServer) processTxLsaAck(lsa_data ospfNeighborAckTxMsg) {
	ack_len := len(lsa_data.lsa_headers_byte)
	total_ack := ack_len / OSPF_LSA_ACK_SIZE
	if total_ack < 0 {
		server.logger.Info(fmt.Sprintln("TX ACK: malformed message. total_ack ", total_ack, " pkt_size ", ack_len))
		return
	}
	nbrConf, exists := server.NeighborConfigMap[lsa_data.nbrKey]
	if !exists {
		server.logger.Warning(fmt.Sprintln("TX ACK: neighbor doesnt exist  to send ack ", lsa_data.nbrKey))
		return
	}
	intf, _ := server.IntfConfMap[nbrConf.intfConfKey]

	dstMac, _ := ospfNeighborIPToMAC[lsa_data.nbrKey]
	dstIp := nbrConf.OspfNbrIPAddr
	pkt := server.BuildLSAAckPkt(nbrConf.intfConfKey, intf, nbrConf, dstMac, dstIp,
		ack_len, lsa_data.lsa_headers_byte)
	// send ack over the pcap.
	server.logger.Info(fmt.Sprintln("ACK SEND: ", pkt))
	server.SendOspfPkt(nbrConf.intfConfKey, pkt)

}
