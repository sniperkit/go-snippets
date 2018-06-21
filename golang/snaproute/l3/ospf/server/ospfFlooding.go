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
	"l3/ospf/config"
	"net"
)

const (
	FloodLsa uint8 = LsdbNoAction + 1
)

/* Flood message structure to be sent
for tx LSAUPD channel
*/
type ospfFloodMsg struct {
	nbrKey  NeighborConfKey
	intfKey IntfConfKey
	areaId  uint32
	lsType  uint8
	linkid  uint32
	lsaKey  LsaKey
	lsOp    uint8  // indicates whether to flood on all interfaces or selective ones.
	pkt     []byte //LSA flood packet received from another neighbor
}

var maxAgeLsaMap map[LsaKey][]byte

/*@fn SendSelfOrigLSA
Api is called
When adjacency is established
DR/BDR change

*/
func (server *OSPFServer) SendSelfOrigLSA(areaId uint32, intfKey IntfConfKey) []byte {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	intConf, _ := server.IntfConfMap[intfKey]
	ospfLsaPkt := newospfNeighborLSAUpdPkt()
	var lsaEncPkt []byte
	LsaEnc := []byte{}
	server.logger.Info(fmt.Sprintln("Flood: Start flooding as Nbr is in full state intf ", intConf.IfIpAddr,
		" areaid ", areaId))

	selfOrigLsaEnt, exist := server.AreaSelfOrigLsa[lsdbKey]
	if !exist {
		return nil
	}
	pktLen := 0
	total_len := 0
	ospfLsaPkt.no_lsas = 0

	for key, valid := range selfOrigLsaEnt {
		if valid {
			switch key.LSType {
			case RouterLSA:
				entry, _ := server.getRouterLsaFromLsdb(areaId, key)
				LsaEnc = encodeRouterLsa(entry, key)
				checksumOffset := uint16(14)
				checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
				binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
				pktLen = len(LsaEnc)
				binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
				lsaid := convertUint32ToIPv4(key.LSId)
				server.logger.Info(fmt.Sprintln("Flood: router  LSA = ", lsaid))
				ospfLsaPkt.lsa = append(ospfLsaPkt.lsa, LsaEnc...)
				ospfLsaPkt.no_lsas++
				total_len += pktLen

			case NetworkLSA:
				rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
				server.logger.Info(fmt.Sprintln("Flood: intfRouterid ", intConf.IfDRtrId, " globalrtrId ", rtr_id))
				if intConf.IfDRtrId == rtr_id {
					server.logger.Info(fmt.Sprintln("Flood: I am DR. Send Nw LSA."))
					entry, _ := server.getNetworkLsaFromLsdb(areaId, key)

					server.logger.Info(fmt.Sprintln("Flood: Network lsa for key ", key, " lsa ", entry))
					LsaEnc = encodeNetworkLsa(entry, key)
					checksumOffset := uint16(14)
					checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
					binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
					pktLen = len(LsaEnc)
					binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
					//server.logger.Info(fmt.Sprintln("Flood: Encoded LSA = ", LsaEnc))
					ospfLsaPkt.lsa = append(ospfLsaPkt.lsa, LsaEnc...)
					ospfLsaPkt.no_lsas++
					total_len += pktLen

				}
			case Summary3LSA, Summary4LSA:
				entry, _ := server.getSummaryLsaFromLsdb(areaId, key)
				LsaEnc = encodeSummaryLsa(entry, key)
				checksumOffset := uint16(14)
				checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
				binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
				pktLen = len(LsaEnc)
				binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
				lsaid := convertUint32ToIPv4(key.LSId)
				server.logger.Info(fmt.Sprintln("Flood: summary  LSA = ", lsaid))
				ospfLsaPkt.lsa = append(ospfLsaPkt.lsa, LsaEnc...)
				ospfLsaPkt.no_lsas++
				total_len += pktLen

			} // end of case
		}
	}

	lsa_pkt_len := total_len + OSPF_NO_OF_LSA_FIELD
	//server.logger.Info(fmt.Sprintln("Flood: Total length ", lsa_pkt_len, "total lsas ", ospfLsaPkt.no_lsas))
	if lsa_pkt_len == OSPF_NO_OF_LSA_FIELD {
		server.logger.Info(fmt.Sprintln("Flood: No LSA to send"))
		return nil
	}
	lsas_enc := make([]byte, 4)

	binary.BigEndian.PutUint32(lsas_enc, ospfLsaPkt.no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, ospfLsaPkt.lsa...)

	//server.logger.Info(fmt.Sprintln("Flood: LSA pkt with #lsas = ", lsaEncPkt))

	return lsaEncPkt
}

/* @fn processFloodMsg
When new LSA is received on the interfaces
flood is based on different checks
LSAFLOOD - Flood router LSA and n/w LSA.
LSASELFLOOD - Flood on selective interface.
LSAINTF - LSA sent over the interface for LSAREQ
*/

func (server *OSPFServer) processFloodMsg(lsa_data ospfFloodMsg) {

	intConf := server.IntfConfMap[lsa_data.intfKey]
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}

	switch lsa_data.lsOp {
	case LSAFLOOD: // flood router LSAs and n/w LSA if DR
		server.logger.Info(fmt.Sprintln("FLOOD: Flood request received from interface key ",
			lsa_data.intfKey, " nbr ", lsa_data.nbrKey))
		for key, intf := range server.IntfConfMap {
			ifAreaId := convertIPv4ToUint32(intf.IfAreaId)
			flood_lsa := server.interfaceFloodCheck(key)
			if !flood_lsa || ifAreaId != lsa_data.areaId {
				continue // dont flood if no nbr is full for this interface
			}
			lsa_upd_pkt := server.SendSelfOrigLSA(lsa_data.areaId, lsa_data.intfKey)
			lsa_pkt_len := len(lsa_upd_pkt)
			if lsa_pkt_len == 0 {
				return
			}
			pkt := server.BuildLsaUpdPkt(key, intf,
				dstMac, dstIp, lsa_pkt_len, lsa_upd_pkt)
			server.SendOspfPkt(key, pkt)
			server.logger.Info(fmt.Sprintln("FLOOD: Nbr FULL intf ", intf.IfIpAddr))
		}

	case LSASELFLOOD: //Flood received LSA on selective interfaces.
		nbrConf := server.NeighborConfigMap[lsa_data.nbrKey]
		rxIntf := server.IntfConfMap[nbrConf.intfConfKey]
		lsid := convertUint32ToIPv4(lsa_data.linkid)
		server.logger.Info(fmt.Sprintln("LSASELFLOOD: Received lsid ", lsid, " lstype ", lsa_data.lsType))
		var lsaEncPkt []byte
		for key, intf := range server.IntfConfMap {
			areaid := convertIPv4ToUint32(intf.IfAreaId)
			if intf.IfIpAddr.Equal(rxIntf.IfIpAddr) || lsa_data.areaId != areaid {
				server.logger.Info(fmt.Sprintln("LSASELFLOOD:Dont flood on rx intf ", rxIntf.IfIpAddr))
				continue // dont flood the LSA on the interface it is received.
			}
			send := server.nbrFloodCheck(lsa_data.nbrKey, key, intf, lsa_data.lsType)
			if send {
				if lsa_data.pkt != nil {
					server.logger.Info(fmt.Sprintln("LSASELFLOOD: Unicast LSA interface ", intf.IfIpAddr, " lsid ", lsid, " lstype ", lsa_data.lsType))
					lsas_enc := make([]byte, 4)
					var no_lsa uint32
					no_lsa = 1
					binary.BigEndian.PutUint32(lsas_enc, no_lsa)
					lsaEncPkt = append(lsaEncPkt, lsas_enc...)
					lsaEncPkt = append(lsaEncPkt, lsa_data.pkt...)
					lsa_pkt_len := len(lsaEncPkt)
					pkt := server.BuildLsaUpdPkt(key, intf,
						dstMac, dstIp, lsa_pkt_len, lsaEncPkt)
					server.SendOspfPkt(key, pkt)
				}
			}
		}
	case LSAINTF: //send the LSA on specific interface for reply to the LSAREQ
		nbrConf, exists := server.NeighborConfigMap[lsa_data.nbrKey]
		if !exists {
			server.logger.Info(fmt.Sprintln("Flood: LSAINTF Neighbor doesnt exist . Dont flood.", lsa_data.nbrKey))
			return
		}
		lsid := convertUint32ToIPv4(lsa_data.linkid)
		var lsaEncPkt []byte
		if lsa_data.pkt != nil {
			lsas_enc := make([]byte, 4)
			var no_lsa uint32
			no_lsa = 1
			binary.BigEndian.PutUint32(lsas_enc, no_lsa)
			lsaEncPkt = append(lsaEncPkt, lsas_enc...)
			lsaEncPkt = append(lsaEncPkt, lsa_data.pkt...)
			lsa_pkt_len := len(lsaEncPkt)
			pkt := server.BuildLsaUpdPkt(nbrConf.intfConfKey, intConf,
				dstMac, dstIp, lsa_pkt_len, lsaEncPkt)
			server.logger.Info(fmt.Sprintln("LSAINTF: Send  LSA to interface ", intConf.IfIpAddr,
				" lsid ", lsid, " lstype ", lsa_data.lsType))
			server.SendOspfPkt(nbrConf.intfConfKey, pkt)

		}

	case LSAROUTERFLOOD: // Flood router LSA.
		server.logger.Info(fmt.Sprintln("FLOOD: Flood for interface event ",
			lsa_data.intfKey))
		for key, intf := range server.IntfConfMap {
			ifAreaId := convertIPv4ToUint32(intf.IfAreaId)
			flood_lsa := server.interfaceFloodCheck(key)
			if !flood_lsa || ifAreaId != lsa_data.areaId {
				continue
			}
			lsa_upd_pkt := server.SendRouterLsa(lsa_data.areaId, lsa_data.intfKey, lsa_data.lsaKey)
			lsa_pkt_len := len(lsa_upd_pkt)
			if lsa_pkt_len == 0 {
				return
			}
			pkt := server.BuildLsaUpdPkt(key, intf,
				dstMac, dstIp, lsa_pkt_len, lsa_upd_pkt)
			server.SendOspfPkt(key, pkt)
		}

	case LSASUMMARYFLOOD:
		server.logger.Info(fmt.Sprintln("Flood: Summary LSA flood msg received."))
		server.processSummaryLSAFlood(lsa_data.areaId, lsa_data.lsaKey)

	case LSAEXTFLOOD: //flood AS External LSA
		server.logger.Info(fmt.Sprintln("LSAEXTFLOOD: Flood external routes for lsa key ", lsa_data.lsaKey))
		server.processAsExternalLSAFlood(lsa_data.lsaKey)

	case LSAAGE: // Flood aged LSAs
		server.constructAndSendLsaAgeFlood()

	}
}

/*@fn sendRouterLsa
At the event of interface down need to flood
updated router LSA.
*/
func (server *OSPFServer) SendRouterLsa(areaId uint32, intfKey IntfConfKey,
	lsaKey LsaKey) []byte {

	ospfLsaPkt := newospfNeighborLSAUpdPkt()
	var lsaEncPkt []byte
	LsaEnc := []byte{}

	pktLen := 0
	total_len := 0
	ospfLsaPkt.no_lsas = 0
	entry, exist := server.getRouterLsaFromLsdb(areaId, lsaKey)
	if exist == LsdbEntryNotFound {
		return nil
	}
	LsaEnc = encodeRouterLsa(entry, lsaKey)
	checksumOffset := uint16(14)
	checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
	pktLen = len(LsaEnc)
	binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
	lsaid := convertUint32ToIPv4(lsaKey.LSId)
	server.logger.Info(fmt.Sprintln("Flood: router  LSA = ", lsaid))
	ospfLsaPkt.lsa = append(ospfLsaPkt.lsa, LsaEnc...)
	ospfLsaPkt.no_lsas++
	total_len += pktLen
	lsa_pkt_len := total_len + OSPF_NO_OF_LSA_FIELD
	//server.logger.Info(fmt.Sprintln("Flood: Total length ", lsa_pkt_len, "total lsas ", ospfLsaPkt.no_lsas))
	if lsa_pkt_len == OSPF_NO_OF_LSA_FIELD {
		server.logger.Info(fmt.Sprintln("Flood: No LSA to send"))
		return nil
	}
	lsas_enc := make([]byte, 4)

	binary.BigEndian.PutUint32(lsas_enc, ospfLsaPkt.no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, ospfLsaPkt.lsa...)

	return lsaEncPkt
}

/*@fn constructAndSendLsaAgeFlood
Flood LSAs which reached max age.
*/
func (server *OSPFServer) constructAndSendLsaAgeFlood() {
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	lsas_enc := make([]byte, 4)
	var lsaEncPkt []byte
	var lsasWithHeader []byte
	var no_lsa uint32
	no_lsa = 0
	total_len := 0
	for lsaKey, lsaPkt := range maxAgeLsaMap {
		if lsaPkt != nil {
			no_lsa++
			checksumOffset := uint16(14)
			checkSum := computeFletcherChecksum(lsaPkt[2:], checksumOffset)
			binary.BigEndian.PutUint16(lsaPkt[16:18], checkSum)
			pktLen := len(lsaPkt)
			binary.BigEndian.PutUint16(lsaPkt[18:20], uint16(pktLen))
			lsaEncPkt = append(lsaEncPkt, lsaPkt...)
			total_len += pktLen
			server.logger.Info(fmt.Sprintln("FLUSH: Added to flush list lsakey",
				lsaKey.AdvRouter, lsaKey.LSId, lsaKey.LSId))
		}
		msg := maxAgeLsaMsg{
			lsaKey:   lsaKey,
			msg_type: delMaxAgeLsa,
		}
		server.maxAgeLsaCh <- msg
	}
	lsa_pkt_len := total_len + OSPF_NO_OF_LSA_FIELD
	if lsa_pkt_len == OSPF_NO_OF_LSA_FIELD {
		return
	}
	binary.BigEndian.PutUint32(lsas_enc, no_lsa)
	lsasWithHeader = append(lsasWithHeader, lsas_enc...)
	lsasWithHeader = append(lsasWithHeader, lsaEncPkt...)

	/* flood on all eligible interfaces */
	for key, intConf := range server.IntfConfMap {
		server.logger.Info(fmt.Sprintln("FLUSH: Send flush message ", intConf.IfIpAddr))
		pkt := server.BuildLsaUpdPkt(key, intConf,
			dstMac, dstIp, lsa_pkt_len, lsasWithHeader)
		server.SendOspfPkt(key, pkt)
	}

}

/* @fn interfaceFloodCheck
Check if we need to flood the LSA on the interface
*/
func (server *OSPFServer) nbrFloodCheck(nbrKey NeighborConfKey, key IntfConfKey, intf IntfConf, lsType uint8) bool {
	/* Check neighbor state */
	flood_check := true
	nbrConf := server.NeighborConfigMap[nbrKey]
	//rtrid := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	if nbrConf.intfConfKey == key && nbrConf.isDRBDR && lsType != Summary3LSA && lsType != Summary4LSA {
		server.logger.Info(fmt.Sprintln("IF FLOOD: Nbr is DR/BDR.   flood on this interface . nbr - ", nbrKey.IPAddr, nbrConf.OspfNbrIPAddr))
		return false
	}
	flood_check = server.interfaceFloodCheck(key)
	return flood_check
}

func (server *OSPFServer) interfaceFloodCheck(key IntfConfKey) bool {
	flood_check := false
	nbrData, exist := ospfIntfToNbrMap[key]
	if !exist {
		server.logger.Info(fmt.Sprintln("FLOOD: Intf to nbr map doesnt exist.Dont flood."))
		return false
	}
	if nbrData.nbrList != nil {
		for index := range nbrData.nbrList {
			nbrId := nbrData.nbrList[index]
			nbrConf := server.NeighborConfigMap[nbrId]
			if nbrConf.OspfNbrState < config.NbrExchange {
				server.logger.Info(fmt.Sprintln("FLOOD: Nbr < exchange . ", nbrConf.OspfNbrIPAddr))
				flood_check = false
				continue
			}
			flood_check = true
			/* TODO - add check if nbrstate is loading - check its retransmission list
			   add LSA to the adjacency list of neighbor with FULL state.*/
		}
	} else {
		server.logger.Info(fmt.Sprintln("FLOOD: nbr list is null for interface ", key.IPAddr))
	}
	return flood_check
}

/*
@fn processSummaryLSAFlood
This API takes care of flooding new summary LSAs that is added in the LSDB
*/
func (server *OSPFServer) processSummaryLSAFlood(areaId uint32, lsaKey LsaKey) {
	var lsaEncPkt []byte
	LsaEnc := []byte{}

	server.logger.Info(fmt.Sprintln("Summary: Start flooding algorithm. Area ",
		areaId, " lsa ", lsaKey))
	LsaEnc = server.encodeSummaryLsa(areaId, lsaKey)
	no_lsas := uint32(1)
	lsas_enc := make([]byte, 4)
	binary.BigEndian.PutUint32(lsas_enc, no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, LsaEnc...)
	lsid := convertUint32ToIPv4(lsaKey.LSId)
	adv_router := convertUint32ToIPv4(lsaKey.AdvRouter)
	server.logger.Info(fmt.Sprintln("SUMMARY: Send for flooding ",
		areaId, " adv_router ", adv_router, " lsid ",
		lsid))
	server.floodSummaryLsa(lsaEncPkt, areaId)
	server.logger.Info(fmt.Sprintln("SUMMARY: End flooding process. lsa", lsaKey))
}

func (server *OSPFServer) encodeSummaryLsa(areaid uint32, lsakey LsaKey) []byte {
	entry, ret := server.getSummaryLsaFromLsdb(areaid, lsakey)
	if ret == LsdbEntryNotFound {
		server.logger.Info(fmt.Sprintln("Summary LSA: Lsa not found . Area",
			areaid, " LSA key ", lsakey))
		return nil
	}
	LsaEnc := encodeSummaryLsa(entry, lsakey)
	pktLen := len(LsaEnc)
	checksumOffset := uint16(14)
	checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
	binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
	return LsaEnc

}

func (server *OSPFServer) floodSummaryLsa(pkt []byte, areaid uint32) {
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	for key, _ := range server.IntfConfMap {
		intf, ok := server.IntfConfMap[key]
		if !ok {
			continue
		}
		ifArea := convertIPv4ToUint32(intf.IfAreaId)
		//isStub := server.isStubArea(areaid)
		if ifArea == areaid {
			// flood to your own area
			nbrMdata, ok := ospfIntfToNbrMap[key]
			if ok && len(nbrMdata.nbrList) > 0 {
				send_pkt := server.BuildLsaUpdPkt(key, intf, dstMac, dstIp, len(pkt), pkt)
				server.logger.Info(fmt.Sprintln("SUMMARY: Send  LSA to interface ", intf.IfIpAddr, " area ", intf.IfAreaId))
				server.SendOspfPkt(key, send_pkt)
			}

		}
	}
}

/*
@fn processAsExternalLSAFlood
	This API takes care of flooding external routes through
	AS external LSA
*/
func (server *OSPFServer) processAsExternalLSAFlood(lsakey LsaKey) {
	areaId := convertAreaOrRouterIdUint32("0.0.0.0")
	for ent, _ := range server.AreaConfMap {
		areaId = convertAreaOrRouterIdUint32(string(ent.AreaId))
	}
	var lsaEncPkt []byte
	LsaEnc := []byte{}

	entry, ret := server.getASExternalLsaFromLsdb(areaId, lsakey)
	if ret == LsdbEntryNotFound {
		server.logger.Info(fmt.Sprintln("ASBR: Lsa not found . Area",
			areaId, " LSA key ", lsakey))
		return
	}
	LsaEnc = encodeASExternalLsa(entry, lsakey)
	pktLen := len(LsaEnc)
	checksumOffset := uint16(14)
	checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
	binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))

	no_lsas := uint32(1)
	lsas_enc := make([]byte, 4)
	binary.BigEndian.PutUint32(lsas_enc, no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, LsaEnc...)
	lsid := convertUint32ToIPv4(lsakey.LSId)
	adv_router := convertUint32ToIPv4(lsakey.AdvRouter)
	server.logger.Info(fmt.Sprintln("ASBR: flood lsid ", lsid, " adv_router ", adv_router))
	server.floodASExternalLsa(lsaEncPkt)
}

func (server *OSPFServer) floodASExternalLsa(pkt []byte) {
	server.logger.Info(fmt.Sprintln("ASBR: Received for flood ", pkt))
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	for key, _ := range server.IntfConfMap {
		intf, ok := server.IntfConfMap[key]
		if !ok {
			continue
		}
		areaId := config.AreaId(convertIPInByteToString(intf.IfAreaId))
		isStub := server.isStubArea(areaId)
		if isStub {
			server.logger.Info(fmt.Sprintln("ASBR: Dont flood AS external as area is stub ", areaId))
			continue
		}
		nbrMdata, ok := ospfIntfToNbrMap[key]
		if ok && len(nbrMdata.nbrList) > 0 {
			send_pkt := server.BuildLsaUpdPkt(key, intf, dstMac, dstIp, len(pkt), pkt)
			server.logger.Info(fmt.Sprintln("ASBR: Send  LSA to interface ", intf.IfIpAddr))
			server.SendOspfPkt(key, send_pkt)
		}
	}
}
