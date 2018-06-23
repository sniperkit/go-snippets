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
	"l3/ospfv2/objects"
	"net"
)

/*
This file decodes database description packets.as per below format
 0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |   Version #   |       2       |         Packet length         |
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
       |       0       |       0       |    Options    |0|0|0|0|0|I|M|MS
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     DD sequence number                        |
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
*/

/* TODO
remote hardcoding and get it while config.
*/
const INTF_MTU_MIN = 1500

func NewOspfDatabaseDescriptionData() *NbrDbdData {
	return &NbrDbdData{}
}

func newospfLSAHeader() *ospfLSAHeader {
	return &ospfLSAHeader{}
}

func DecodeDatabaseDescriptionData(data []byte, dbd_data *NbrDbdData, Pktlen uint16) {
	dbd_data.interface_mtu = binary.BigEndian.Uint16(data[0:2])
	dbd_data.options = data[2]
	dbd_data.dd_sequence_number = binary.BigEndian.Uint32(data[4:8])
	imms_options := data[3]
	dbd_data.ibit = imms_options&0x4 != 0
	dbd_data.mbit = imms_options&0x02 != 0
	dbd_data.msbit = imms_options&0x01 != 0

	fmt.Println("Decoded packet options:", dbd_data.options,
		"IMMS:", dbd_data.ibit, dbd_data.mbit, dbd_data.msbit,
		"seq num:", dbd_data.dd_sequence_number)

	if dbd_data.ibit == false {
		// negotiation is done. Check if we have LSA headers

		headers_len := Pktlen - (OSPF_DBD_MIN_SIZE + OSPF_HEADER_SIZE)
		fmt.Println("DBD: Received headers_len ", headers_len, " PktLen", Pktlen, " data len ", len(data))
		if headers_len >= 20 && headers_len < Pktlen {
			fmt.Println("DBD: LSA headers length ", headers_len)
			num_headers := int(headers_len / 20)
			fmt.Println("DBD: Received ", num_headers, " LSA headers.")
			header_byte := make([]byte, num_headers*OSPF_LSA_HEADER_SIZE)
			var start_index uint16
			var lsa_header ospfLSAHeader
			for i := 0; i < num_headers; i++ {
				start_index = uint16(OSPF_DBD_MIN_SIZE + (i * OSPF_LSA_HEADER_SIZE))
				copy(header_byte, data[start_index:start_index+20])
				lsa_header = decodeLSAHeader(header_byte)
				fmt.Println("DBD: Header decoded ",
					"ls_age:options:ls_type:link_state_id:adv_rtr:ls_seq:ls_checksum ",
					lsa_header.ls_age, lsa_header.ls_type, lsa_header.link_state_id,
					lsa_header.adv_router_id, lsa_header.ls_sequence_num,
					lsa_header.ls_checksum)
				dbd_data.lsa_headers = append(dbd_data.lsa_headers, lsa_header)
			}
		}
	}
}

/*

LSA headers
 0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |           LS Age              |           LS Type             |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                       Link State ID                           |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                    Advertising Router                         |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                    LS Sequence Number                         |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |        LS Checksum            |             Length            |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

func decodeLSAHeader(data []byte) (lsa_header ospfLSAHeader) {
	lsa_header.ls_age = binary.BigEndian.Uint16(data[0:2])
	lsa_header.ls_type = data[3]
	lsa_header.options = data[2]
	lsa_header.link_state_id = binary.BigEndian.Uint32(data[4:8])
	lsa_header.adv_router_id = binary.BigEndian.Uint32(data[8:12])
	lsa_header.ls_sequence_num = binary.BigEndian.Uint32(data[12:16])
	lsa_header.ls_checksum = binary.BigEndian.Uint16(data[16:18])
	lsa_header.ls_len = binary.BigEndian.Uint16(data[18:20])

	return lsa_header
}

func encodeLSAHeader(dd_data NbrDbdData) []byte {
	headers := len(dd_data.lsa_headers)

	if headers == 0 {
		return nil
	}
	//fmt.Sprintln("no of headers ", headers)
	pkt := make([]byte, headers*OSPF_LSA_HEADER_SIZE)
	for index := 0; index < headers; index++ {
		//	fmt.Sprintln("Attached header ", index)
		lsa_header := dd_data.lsa_headers[index]
		pkt_index := 20 * index
		binary.BigEndian.PutUint16(pkt[pkt_index:pkt_index+2], lsa_header.ls_age)
		pkt[pkt_index+2] = lsa_header.options
		pkt[pkt_index+3] = lsa_header.ls_type
		binary.BigEndian.PutUint32(pkt[pkt_index+4:pkt_index+8], lsa_header.link_state_id)
		binary.BigEndian.PutUint32(pkt[pkt_index+8:pkt_index+12], lsa_header.adv_router_id)
		binary.BigEndian.PutUint32(pkt[pkt_index+12:pkt_index+16], lsa_header.ls_sequence_num)
		binary.BigEndian.PutUint16(pkt[pkt_index+16:pkt_index+18], lsa_header.ls_checksum)
		binary.BigEndian.PutUint16(pkt[pkt_index+18:pkt_index+20], lsa_header.ls_len)
	}
	return pkt
}

/*
0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |   Version #   |       2       |         Packet length         |
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
       |       0       |       0       |    Options    |0|0|0|0|0|I|M|MS
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     DD sequence number                        |
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
func encodeDatabaseDescriptionData(dd_data NbrDbdData) []byte {
	pkt := make([]byte, OSPF_DBD_MIN_SIZE)
	binary.BigEndian.PutUint16(pkt[0:2], dd_data.interface_mtu)
	pkt[2] = 0x2
	imms := 0
	if dd_data.ibit {
		imms = imms | 0x4
	}
	if dd_data.mbit {
		imms = imms | 0x2
	}
	if dd_data.msbit {
		imms = imms | 0x1
	}
	pkt[3] = byte(imms)
	//fmt.Println("data imms  ", pkt[3])
	binary.BigEndian.PutUint32(pkt[4:8], dd_data.dd_sequence_number)
	lsa_pkt := encodeLSAHeader(dd_data)
	if lsa_pkt != nil {
		pkt = append(pkt, lsa_pkt...)
	}

	return pkt
}

func (server *OSPFV2Server) BuildAndSendDdBDPkt(nbrConf NbrConf, dbdData NbrDbdData) {
	var dstMAC net.HardwareAddr
	server.logger.Debug("Nbr : Send the dbd packet out ")
	intfKey := nbrConf.IntfKey
	ent, exist := server.IntfConfMap[intfKey]
	if !exist {
		server.logger.Err("Nbr : failed to send db packet ", intfKey)
		return
	}

	ospfHdr := OSPFHeader{
		Ver:      OSPF_VERSION_2,
		PktType:  uint8(DBDescriptionType),
		Pktlen:   0,
		RouterId: server.globalData.RouterId,
		AreaId:   ent.AreaId,
		Chksum:   0,
		AuthType: ent.AuthType,
	}

	ospfPktlen := OSPF_HEADER_SIZE
	lsa_header_size := OSPF_LSA_HEADER_SIZE * len(dbdData.lsa_headers)
	ospfPktlen = ospfPktlen + OSPF_DBD_MIN_SIZE + lsa_header_size

	ospfHdr.Pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	server.logger.Debug("ospfEncHdr:", ospfEncHdr)
	dbdDataEnc := encodeDatabaseDescriptionData(dbdData)
	//server.logger.Info(fmt.Sprintln("DBD Pkt:", dbdDataEnc))

	ospf := append(ospfEncHdr, dbdDataEnc...)
	server.logger.Debug("OSPF DBD:", ospf)
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	binary.BigEndian.PutUint64(ospf[16:24], ent.AuthKey)

	var DstIP net.IP

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.Pktlen
	if ent.FSMState == objects.INTF_FSM_STATE_P2P {
		DstIP = net.ParseIP(config.AllSPFRouters)
		dstMAC, _ = net.ParseMAC(ALLSPFROUTERMAC)
	} else {
		dstMAC = nbrConf.NbrMac
		DstIP = net.ParseIP(convertUint32ToDotNotation(nbrConf.NbrIP))
	}
	SrcIp := net.ParseIP(convertUint32ToDotNotation(ent.IpAddr))
	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    SrcIp,
		DstIP:    DstIP,
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
	//	server.logger.Info(fmt.Sprintln("buffer: ", buffer))
	dbdPkt := buffer.Bytes()
	//	server.logger.Info(fmt.Sprintln("dbdPkt: ", dbdPkt))
	server.logger.Debug("Nbr: send db packet out ", dbdPkt)
	server.SendOspfPkt(intfKey, dbdPkt)
}

func (server *OSPFV2Server) ProcessRxDbdPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, nbrKey NbrConfKey) error {
	ospfdbd_data := NewOspfDatabaseDescriptionData()
	ospfdbd_data.lsa_headers = []ospfLSAHeader{}
	//routerId := convertIPv4ToUint32(ospfHdrMd.routerId)
	Pktlen := ospfHdrMd.Pktlen

	if Pktlen < OSPF_DBD_MIN_SIZE+OSPF_HEADER_SIZE {
		server.logger.Warning(fmt.Sprintln("DBD WARNING: Packet < min DBD length. Pktlen ", Pktlen,
			" min_dbd_len ", OSPF_DBD_MIN_SIZE+OSPF_HEADER_SIZE))
	}

	DecodeDatabaseDescriptionData(data, ospfdbd_data, Pktlen)
	//ipaddr := convertIPInByteToString(ipHdrMd.srcIP)

	dbdNbrMsg := NbrDbdMsg{
		nbrConfKey: nbrKey,
		nbrDbdData: *ospfdbd_data,
	}
	server.logger.Debug(fmt.Sprintln("DBD: nbr key ", nbrKey))
	//fmt.Println(" lsa_header length = ", len(ospfdbd_data.lsa_headers))
	dbdNbrMsg.nbrDbdData.lsa_headers = []ospfLSAHeader{}

	copy(dbdNbrMsg.nbrDbdData.lsa_headers, ospfdbd_data.lsa_headers)
	for i := 0; i < len(ospfdbd_data.lsa_headers); i++ {
		dbdNbrMsg.nbrDbdData.lsa_headers = append(dbdNbrMsg.nbrDbdData.lsa_headers,
			ospfdbd_data.lsa_headers[i])
	}
	//send packet to nbr fsm
	server.NbrConfData.neighborDBDEventCh <- dbdNbrMsg
	return nil
}

func (server *OSPFV2Server) ConstructDbdMdata(nbrKey NbrConfKey,
	ibit bool, mbit bool, msbit bool, options uint8,
	seq uint32, append_lsa bool, is_duplicate bool) (dbd_mdata NbrDbdData, last_exchange bool) {
	last_exchange = true
	nbrConf, exists := server.NbrConfMap[nbrKey]
	if !exists {
		server.logger.Err(fmt.Sprintln("DBD: Failed to send initial dbd packet as nbr doesnt exist. nbr",
			nbrKey.NbrIdentity))
		return dbd_mdata, last_exchange
	}
	intfConf, _ := server.IntfConfMap[nbrConf.IntfKey]
	dbd_mdata.ibit = ibit
	dbd_mdata.mbit = mbit
	dbd_mdata.msbit = msbit

	dbd_mdata.interface_mtu = uint16(intfConf.Mtu)
	dbd_mdata.options = options
	dbd_mdata.dd_sequence_number = seq

	lsa_count_done := 0
	lsa_count_att := 0
	if append_lsa && exists {

		dbd_mdata.lsa_headers = []ospfLSAHeader{}
		var index uint8

		db_list := nbrConf.NbrDBSummaryList
		server.logger.Debug(fmt.Sprintln("DBD: db_list ", db_list))
		if len(db_list) == 0 {
			for index = 0; index < uint8(len(db_list)); index++ {
				dbd_mdata.lsa_headers = append(dbd_mdata.lsa_headers, *db_list[index])
				lsa_count_att++
				lsa_count_done++
			}
		}
		if (lsa_count_att + lsa_count_done) == len(db_list) {
			dbd_mdata.mbit = false
			last_exchange = true
		}
	}

	server.logger.Debug(fmt.Sprintln("DBDSEND: nbr state ", nbrConf.State,
		" imms ", dbd_mdata.ibit, dbd_mdata.mbit, dbd_mdata.msbit,
		" seq num ", seq, "options ", dbd_mdata.options, " headers_list ", dbd_mdata.lsa_headers))

	//	data := newDbdMsg(nbrKey, dbd_mdata)
	// send the data
	return dbd_mdata, last_exchange
}

/*
 @fn calculateDBLsaAttach
	This API detects how many LSA headers can be added in
	the DB packet
*/
func (server *OSPFV2Server) calculateDBLsaAttach(nbrKey NbrConfKey, nbrConf NbrConf) (last_exchange bool, lsa_attach uint8) {
	last_exchange = true
	lsa_attach = 0

	max_lsa_headers := calculateMaxLsaHeaders()
	db_list := nbrConf.NbrDBSummaryList
	slice_len := len(db_list)
	server.logger.Info(fmt.Sprintln("DBD: slice_len ", slice_len, "max_lsa_header ", max_lsa_headers,
		"nbrConf.lsa_index ", nbrConf.NbrLsaIndex))
	if slice_len == int(nbrConf.NbrLsaIndex) {
		return
	}
	if max_lsa_headers > (uint8(slice_len) - uint8(nbrConf.NbrLsaIndex)) {
		lsa_attach = uint8(slice_len) - uint8(nbrConf.NbrLsaIndex)
	} else {
		lsa_attach = max_lsa_headers
	}
	if (uint8(nbrConf.NbrLsaIndex) + lsa_attach) >= uint8(slice_len) {
		// the last slice in the list being sent
		server.logger.Info(fmt.Sprintln("DBD:  Send the last dd packet with nbr/state ", nbrKey.NbrIdentity, nbrConf.State))
		last_exchange = true
	}
	return last_exchange, 0
}
