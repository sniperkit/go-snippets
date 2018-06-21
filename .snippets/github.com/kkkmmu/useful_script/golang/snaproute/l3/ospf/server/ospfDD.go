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
const INTF_OPTIONS = 66

type ospfDatabaseDescriptionData struct {
	options            uint8
	interface_mtu      uint16
	dd_sequence_number uint32
	ibit               bool
	mbit               bool
	msbit              bool
	lsa_headers        []ospfLSAHeader
}

type ospfLSAHeader struct {
	ls_age          uint16
	options         uint8
	ls_type         uint8
	link_state_id   uint32
	adv_router_id   uint32
	ls_sequence_num uint32
	ls_checksum     uint16
	ls_len          uint16
}

func NewOspfDatabaseDescriptionData() *ospfDatabaseDescriptionData {
	return &ospfDatabaseDescriptionData{}
}

func newospfLSAHeader() *ospfLSAHeader {
	return &ospfLSAHeader{}
}

func newDbdMsg(key NeighborConfKey, dbd_data ospfDatabaseDescriptionData) ospfNeighborDBDMsg {
	dbdNbrMsg := ospfNeighborDBDMsg{
		ospfNbrConfKey: key,
		ospfNbrDBDData: dbd_data,
	}
	return dbdNbrMsg
}

func DecodeDatabaseDescriptionData(data []byte, dbd_data *ospfDatabaseDescriptionData, pktlen uint16) {
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

		headers_len := pktlen - (OSPF_DBD_MIN_SIZE + OSPF_HEADER_SIZE)
		fmt.Println("DBD: Received headers_len ", headers_len, " pktLen", pktlen, " data len ", len(data))
		if headers_len >= 20 && headers_len < pktlen {
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

func encodeLSAHeader(dd_data ospfDatabaseDescriptionData) []byte {
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
func encodeDatabaseDescriptionData(dd_data ospfDatabaseDescriptionData) []byte {
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

func (server *OSPFServer) BuildDBDPkt(intfKey IntfConfKey, ent IntfConf,
	nbrConf OspfNeighborEntry, dbdData ospfDatabaseDescriptionData, dstMAC net.HardwareAddr) (data []byte) {
	ospfHdr := OSPFHeader{
		ver:      OSPF_VERSION_2,
		pktType:  uint8(DBDescriptionType),
		pktlen:   0,
		routerId: server.ospfGlobalConf.RouterId,
		areaId:   ent.IfAreaId,
		chksum:   0,
		authType: ent.IfAuthType,
	}

	ospfPktlen := OSPF_HEADER_SIZE
	lsa_header_size := OSPF_LSA_HEADER_SIZE * len(dbdData.lsa_headers)
	ospfPktlen = ospfPktlen + OSPF_DBD_MIN_SIZE + lsa_header_size

	ospfHdr.pktlen = uint16(ospfPktlen)

	ospfEncHdr := encodeOspfHdr(ospfHdr)
	//server.logger.Info(fmt.Sprintln("ospfEncHdr:", ospfEncHdr))
	dbdDataEnc := encodeDatabaseDescriptionData(dbdData)
	//server.logger.Info(fmt.Sprintln("DBD Pkt:", dbdDataEnc))

	ospf := append(ospfEncHdr, dbdDataEnc...)
	//server.logger.Info(fmt.Sprintln("OSPF DBD:", ospf))
	csum := computeCheckSum(ospf)
	binary.BigEndian.PutUint16(ospf[12:14], csum)
	copy(ospf[16:24], ent.IfAuthKey)

	var DstIP net.IP
	var DstMAC net.HardwareAddr

	ipPktlen := IP_HEADER_MIN_LEN + ospfHdr.pktlen
	SrcIP := ent.IfIpAddr

	if ent.IfType == config.NumberedP2P {
		DstIP = net.ParseIP(config.AllSPFRouters)
		DstMAC, _ = net.ParseMAC(config.McastMAC)
	} else {
		DstIP = nbrConf.OspfNbrIPAddr
		DstMAC = dstMAC
	}

	ipLayer := layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(IP_HEADER_MIN_LEN),
		TOS:      uint8(0xc0),
		Length:   uint16(ipPktlen),
		TTL:      uint8(1),
		Protocol: layers.IPProtocol(OSPF_PROTO_ID),
		SrcIP:    SrcIP,
		DstIP:    DstIP,
	}

	ethLayer := layers.Ethernet{
		SrcMAC:       ent.IfMacAddr,
		DstMAC:       DstMAC,
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

	return dbdPkt

}

func (server *OSPFServer) ProcessRxDbdPkt(data []byte, ospfHdrMd *OspfHdrMetadata,
	ipHdrMd *IpHdrMetadata, key IntfConfKey, srcMAC net.HardwareAddr) error {
	ospfdbd_data := NewOspfDatabaseDescriptionData()
	ospfdbd_data.lsa_headers = []ospfLSAHeader{}
	//routerId := convertIPv4ToUint32(ospfHdrMd.routerId)
	pktlen := ospfHdrMd.pktlen

	if pktlen < OSPF_DBD_MIN_SIZE+OSPF_HEADER_SIZE {
		server.logger.Warning(fmt.Sprintln("DBD WARNING: Packet < min DBD length. pktlen ", pktlen,
			" min_dbd_len ", OSPF_DBD_MIN_SIZE+OSPF_HEADER_SIZE))
	}

	DecodeDatabaseDescriptionData(data, ospfdbd_data, pktlen)
	//ipaddr := convertIPInByteToString(ipHdrMd.srcIP)
	ipaddr := net.IPv4(ipHdrMd.srcIP[0], ipHdrMd.srcIP[1], ipHdrMd.srcIP[2], ipHdrMd.srcIP[3])

	dbdNbrMsg := ospfNeighborDBDMsg{
		ospfNbrConfKey: NeighborConfKey{
			IPAddr:  config.IpAddress(ipaddr.String()),
			IntfIdx: key.IntfIdx,
		},
		ospfNbrDBDData: *ospfdbd_data,
	}
	server.logger.Debug(fmt.Sprintln("DBD: nbr key ", ipaddr, key.IntfIdx))
	if ospfNeighborIPToMAC == nil {
		server.logger.Info(fmt.Sprintln("DBD: ospfNeighborIPToMAC is NULL. Check if nbr thread is running."))
		return nil
	}
	ospfNeighborIPToMAC[dbdNbrMsg.ospfNbrConfKey] = srcMAC
	//fmt.Println(" lsa_header length = ", len(ospfdbd_data.lsa_headers))
	dbdNbrMsg.ospfNbrDBDData.lsa_headers = []ospfLSAHeader{}

	copy(dbdNbrMsg.ospfNbrDBDData.lsa_headers, ospfdbd_data.lsa_headers)
	for i := 0; i < len(ospfdbd_data.lsa_headers); i++ {
		dbdNbrMsg.ospfNbrDBDData.lsa_headers = append(dbdNbrMsg.ospfNbrDBDData.lsa_headers,
			ospfdbd_data.lsa_headers[i])
	}

	server.neighborDBDEventCh <- dbdNbrMsg
	//fmt.Println("msg lsa_header length = ", len(dbdNbrMsg.ospfNbrDBDData.lsa_headers))
	return nil
}

func (server *OSPFServer) ConstructAndSendDbdPacket(nbrKey NeighborConfKey,
	ibit bool, mbit bool, msbit bool, options uint8,
	seq uint32, append_lsa bool, is_duplicate bool, ifMtu int32) (dbd_mdata ospfDatabaseDescriptionData, last_exchange bool) {
	last_exchange = true
	nbrCon, exists := server.NeighborConfigMap[nbrKey]
	if !exists {
		server.logger.Err(fmt.Sprintln("DBD: Failed to send initial dbd packet as nbr doesnt exist. nbr",
			nbrKey.IPAddr))
		return dbd_mdata, last_exchange
	}

	dbd_mdata.ibit = ibit
	dbd_mdata.mbit = mbit
	dbd_mdata.msbit = msbit

	server.logger.Debug(fmt.Sprintln("DBD: MTU ", ifMtu))
	dbd_mdata.interface_mtu = uint16(ifMtu)
	dbd_mdata.options = options
	dbd_mdata.dd_sequence_number = seq

	lsa_count_done := 0
	lsa_count_att := 0
	if append_lsa && exists {

		dbd_mdata.lsa_headers = []ospfLSAHeader{}
		var index uint8

		nbrCon.db_summary_list_mutex.Lock()
		db_list, exist := ospfNeighborDBSummary_list[nbrKey]
		server.logger.Debug(fmt.Sprintln("DBD: db_list ", db_list))
		if exist {
			for index = 0; index < uint8(len(db_list)); index++ {
				if db_list[index].valid {
					dbd_mdata.lsa_headers = append(dbd_mdata.lsa_headers, db_list[index].lsa_headers)
					lsa_count_att++
				} else {
					lsa_count_done++
				}
				db_list[index].valid = false
			}
		}
		nbrCon.db_summary_list_mutex.Unlock()
		if (lsa_count_att + lsa_count_done) == len(db_list) {
			dbd_mdata.mbit = false
			last_exchange = true
		}
	}

	server.logger.Debug(fmt.Sprintln("DBDSEND: nbr state ", nbrCon.OspfNbrState,
		" imms ", dbd_mdata.ibit, dbd_mdata.mbit, dbd_mdata.msbit,
		" seq num ", seq, "options ", dbd_mdata.options, " headers_list ", dbd_mdata.lsa_headers))

	data := newDbdMsg(nbrKey, dbd_mdata)
	server.ospfNbrDBDSendCh <- data
	return dbd_mdata, last_exchange
}

/*
 @fn calculateDBLsaAttach
	This API detects how many LSA headers can be added in
	the DB packet
*/
func (server *OSPFServer) calculateDBLsaAttach(nbrKey NeighborConfKey, nbrConf OspfNeighborEntry) (last_exchange bool, lsa_attach uint8) {
	last_exchange = true
	lsa_attach = 0

	max_lsa_headers := calculateMaxLsaHeaders()
	db_list := ospfNeighborDBSummary_list[nbrKey]
	slice_len := len(db_list)
	server.logger.Info(fmt.Sprintln("DBD: slice_len ", slice_len, "max_lsa_header ", max_lsa_headers,
		"nbrConf.lsa_index ", nbrConf.ospfNbrLsaIndex))
	if slice_len == int(nbrConf.ospfNbrLsaIndex) {
		return
	}
	if max_lsa_headers > (uint8(slice_len) - uint8(nbrConf.ospfNbrLsaIndex)) {
		lsa_attach = uint8(slice_len) - uint8(nbrConf.ospfNbrLsaIndex)
	} else {
		lsa_attach = max_lsa_headers
	}
	if (nbrConf.ospfNbrLsaIndex + lsa_attach) >= uint8(slice_len) {
		// the last slice in the list being sent
		server.logger.Info(fmt.Sprintln("DBD:  Send the last dd packet with nbr/state ", nbrKey.IPAddr, nbrConf.OspfNbrState))
		last_exchange = true
	}
	return last_exchange, 0
}
