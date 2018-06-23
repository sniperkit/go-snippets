//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

/* ospfLsdb_test
   This test covers
   1) LSA encode and decode APIs.
   2) LSDB operations.
   3) LSDB timer related APIs.
   4) Area border router APIs.
   5) SPF calculation APIs which are driven by LSDB.
   6) Flooding operations that are driven by LSDB events.
*/
package server

import (
	"fmt"
	"l3/ospf/config"
	"net"
	"testing"
)

func initLsdbTestParams() {
	fmt.Println("\n Get Server object")
	ospf = getServerObject()
	initAttr()
	ospf.IntfConfMap[key] = intf
	ospf.processGlobalConfig(gConf)
	ospf.InitNeighborStateMachine()
	ospfNeighborIPToMAC = make(map[NeighborConfKey]net.HardwareAddr)
	ospfNeighborIPToMAC[nbrKey] = dstMAC
	go startDummyChannels(ospf)
}

func _TestOspfLsdb(t *testing.T) {
	fmt.Println("\n**************** LSDB ************\n")
	initLsdbTestParams()
	for index := 1; index < 21; index++ {
		err := lsdbTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test case for interface FSM")
		}
	}
}

func lsdbTestLogic(tNum int) int {
	areaId := uint32(2)
	switch tNum {
	case 1:
		fmt.Println(tNum, ": Running initLSDatabase ")
		ospf.initLSDatabase(areaId)
	case 2:
		fmt.Println(tNum, ": Running insertSummaryLsa ")
		ospf.initLSDatabase(lsdbKey.AreaId)
		ospf.insertSummaryLsa(lsdbKey, summaryKey, summaryLsa)
	case 3:
		fmt.Println(tNum, ": Running processRecvdRouterLsa ")
		ospf.initLSDatabase(lsdbKey.AreaId)
		ospf.processRecvdRouterLsa(lsa_router, lsdbKey.AreaId)
		ospf.processRecvdLsa(lsa_router, lsdbKey.AreaId)
		ospf.processRecvdLsa(lsa_fake, lsdbKey.AreaId)
		ospf.processDeleteLsa(lsa_router, lsdbKey.AreaId)
		ospf.processDeleteLsa(lsa_fake, lsdbKey.AreaId)
	case 4:
		fmt.Println(tNum, ": Running processRecvdNetworkLsa")
		ospf.initLSDatabase(lsdbKey.AreaId)

		ospf.processRecvdNetworkLsa(lsa_network, lsdbKey.AreaId)
		ospf.processRecvdLsa(lsa_network, lsdbKey.AreaId)
		ospf.processDeleteLsa(lsa_network, lsdbKey.AreaId)
	case 5:
		fmt.Println(tNum, ": Running processRecvdSummaryLsa ")
		ospf.initLSDatabase(lsdbKey.AreaId)
		ospf.processRecvdSummaryLsa(lsa_summary, lsdbKey.AreaId, Summary3LSA)
		ospf.processRecvdLsa(lsa_summary, lsdbKey.AreaId)
		ospf.processDeleteLsa(lsa_summary, lsdbKey.AreaId)
	case 6:
		fmt.Println(tNum, ": Running processRecvdASExternalLsa ")
		ospf.initLSDatabase(lsdbKey.AreaId)
		ospf.processRecvdASExternalLsa(lsa_asExt, lsdbKey.AreaId)
		ospf.processRecvdLsa(lsa_asExt, lsdbKey.AreaId)
		ospf.processDeleteLsa(lsa_asExt, lsdbKey.AreaId)

	case 7:
		fmt.Println(tNum, ": Running processLSDatabaseUpdates")
		checkLSDatabaseUpdates()

	case 8:
		fmt.Println(tNum, ": Running LSAPKT tests ")
		checkLsaPktApis()
	case 9:
		fmt.Println(tNum, ": Running LSA decode tests ")
		checkFloodAPIs()
	}

	return SUCCESS

}

func checkLSDatabaseUpdates() {
	ospf.StartLSDatabase()
	ospf.initLSDatabase(lsdbKey.AreaId)
	lsdb_msg := NewLsdbUpdateMsg()
	lsdb_msg.AreaId = lsdbKey.AreaId
	lsdb_msg.Data = make([]byte, len(lsa_router))
	copy(lsdb_msg.Data, lsa_router)
	lsdb_msg.MsgType = LsdbAdd

	ospf.LsdbUpdateCh <- *lsdb_msg

	lsdb_msg.MsgType = LsdbDel
	ospf.LsdbUpdateCh <- *lsdb_msg
	lsdb_msg.MsgType = LsdbUpdate
	ospf.LsdbUpdateCh <- *lsdb_msg
	ospf.ospfGlobalConf.AreaBdrRtrStatus = true
	msg := NetworkLSAChangeMsg{
		areaId:  lsdbKey.AreaId,
		intfKey: key,
	}
	ospf.IntfStateChangeCh <- msg
	ospf.processInterfaceChangeMsg(msg)

	nbrMdata := newospfNbrMdata()
	nbrMdata.areaId = lsdbKey.AreaId
	nbrMdata.intf = key
	nbrMdata.isDR = true
	nbrMdata.nbrList = nil
	/* CHECK why blocked */
	//ospf.CreateNetworkLSACh <- *nbrMdata
	ospf.processNeighborFullEvent(*nbrMdata)

	msg1 := DrChangeMsg{
		areaId:   lsdbKey.AreaId,
		intfKey:  key,
		oldstate: config.OtherDesignatedRouter,
		newstate: config.BackupDesignatedRouter,
	}
	//ospf.NetworkDRChangeCh <- msg1
	ospf.processDrBdrChangeMsg(msg1)

	routemdata := RouteMdata{
		ipaddr: 2,
		mask:   100,
		metric: 10,
		isDel:  false,
	}
	//	ospf.ExternalRouteNotif <- routemdata
	ospf.processExtRouteUpd(routemdata)

	msg2 := maxAgeLsaMsg{
		lsaKey:   summaryKey,
		msg_type: delMaxAgeLsa,
	}
	//	ospf.maxAgeLsaCh <- msg2
	ospf.processMaxAgeLsaMsg(msg2)

}

func checkLsaPktApis() {

	lsaHeader := getLsaHeaderFromLsa(routerLsa.LsaMd.LSAge, routerLsa.LsaMd.Options, routerKey.LSType,
		routerKey.LSId, routerKey.AdvRouter, uint32(routerLsa.LsaMd.LSSequenceNum),
		routerLsa.LsaMd.LSChecksum, routerLsa.LsaMd.LSLen)
	fmt.Println("Decoded LSA header ", lsaHeader)
	lsaDecode := decodeLSAReq(lsareq)
	fmt.Println("Decoded LSA req ", lsaDecode)
	decodeLSAReqPkt(lsa_router, uint16(len(lsa_router)))

	/* LSA req */
	encodeLSAReq(lsa_reqs)
	lsaPkt := ospf.EncodeLSAReqPkt(key, intf, nbrConf, lsa_reqs, dstMAC)
	fmt.Println("Encoded LSA packet ", lsaPkt)
	nbr_req = &ospfNeighborReq{}
	nbr_req.lsa_headers = lsaHeader
	nbr_req.valid = true
	nbr_req_list = []*ospfNeighborReq{}
	nbr_req_list = append(nbr_req_list, nbr_req)
	ospfNeighborRequest_list[nbrKey] = nbr_req_list
	index := ospf.BuildAndSendLSAReq(nbrKey, nbrConf)
	fmt.Println("Nbr lsa req list index ", index)

	lsaPkt = ospf.BuildLsaUpdPkt(key, intf, dstMAC, dstIP, len(lsa_router), lsa_router)
	fmt.Println("Encoded LSA pkt :", lsaPkt)

	err := ospf.ProcessRxLsaUpdPkt(lsa_router, &ospfHdrMd, &ipHdrMd, key)
	if err != nil {
		fmt.Println("Failed to process received Rx LSA packet.", err)
	}

	/* LSA upd */
	lsaupd_msg := ospfNeighborLSAUpdMsg{
		nbrKey: nbrKey,
		data:   lsa_update,
		areaId: lsdbKey.AreaId,
	}

	nbrConf.OspfNbrState = config.NbrFull
	ospf.NeighborConfigMap[nbrKey] = nbrConf

	ospf.DecodeLSAUpd(lsaupd_msg)

	ospf.selfGenLsaCheck(routerKey)
	ospf.lsaUpdDiscardCheck(nbrConf, lsa_router)

	/* LSA ack */
	lsaAck := ospf.BuildLSAAckPkt(key, intf, nbrConf, dstMAC, dstIP, len(lsaack), lsaack)
	fmt.Println("Encoded lsa ack packet ", lsaAck)

	ospf.ProcessRxLSAAckPkt(lsaack, &ospfHdrMd, &ipHdrMd, key)
	lsa_headers := []ospfLSAHeader{}
	lsa_headers = append(lsa_headers, lsaHeader)
	lsaHeader = getLsaHeaderFromLsa(routerLsa.LsaMd.LSAge, routerLsa.LsaMd.Options, NetworkLSA,
		routerKey.LSId, routerKey.AdvRouter, uint32(routerLsa.LsaMd.LSSequenceNum),
		routerLsa.LsaMd.LSChecksum, routerLsa.LsaMd.LSLen)
	lsa_headers = append(lsa_headers, lsaHeader)
	ack_msg.lsa_headers = lsa_headers
	ack_msg.nbrKey = nbrKey

	ospf.DecodeLSAAck(*ack_msg)

	/* LSA req */
	ospfHdrMd.pktlen = uint16(len(lsareq) + OSPF_HEADER_SIZE)
	err = ospf.ProcessRxLSAReqPkt(lsareq, &ospfHdrMd, &ipHdrMd, key)
	if err != nil {
		fmt.Println("Failed to process rx LSA req pkt ", err)
	}
	nbrConf.OspfNbrState = config.NbrFull
	ospf.DecodeLSAReq(nbrLsaReqMsg)
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)
	req.ls_type = uint32(NetworkLSA)
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)
	req.ls_type = uint32(Summary4LSA)
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)
	req.ls_type = uint32(Summary3LSA)
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)
	req.ls_type = uint32(ASExternalLSA)
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)
	req.ls_type = uint32(10) //fake type
	ospf.generateLsaUpdUnicast(req, nbrKey, lsdbKey.AreaId)

	discard := ospf.lsaReqPacketDiscardCheck(nbrConf, req)
	if discard {
		fmt.Println("Discard this packet ")
	}
	discard = ospf.lsaAckPacketDiscardCheck(nbrConf)

	/* LSA sanity checks */
	discard = ospf.lsaAddCheck(lsaHeader, nbrConf)

	ospf.lsaReTxTimerCheck(nbrKey)

	ospf.processTxLsaAck(ackTxMsg)

}

/* UT for LSA decode routines and
   Flooding */
func checkFloodAPIs() {
	ospf.initLSDatabase(areaId)
	lsDbEnt, _ := ospf.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := ospf.AreaSelfOrigLsa[lsdbKey]

	routerLsa := &RouterLsa{}
	lsaKey := &LsaKey{}
	decodeRouterLsa(lsa_router, routerLsa, lsaKey)
	routerLsa_byte := encodeRouterLsa(*routerLsa, *lsaKey)
	fmt.Println("Encoded router LSA ", routerLsa_byte)
	rlsa, ret := ospf.getRouterLsaFromLsdb(lsdbKey.AreaId, *lsaKey)
	fmt.Println("Lsa from db(lsa, ret) ", rlsa, ret)
	lsDbEnt.RouterLsaMap[*lsaKey] = *routerLsa
	selfOrigLsaEnt[*lsaKey] = true
	ospf.regenerateLsa(lsdbKey, *lsaKey)

	networkLsa := &NetworkLsa{}
	decodeNetworkLsa(lsa_network, networkLsa, lsaKey)
	networkLsa_byte := encodeNetworkLsa(*networkLsa, *lsaKey)
	fmt.Println("Encoded network LSA ", networkLsa_byte)
	nlsa, ret := ospf.getRouterLsaFromLsdb(lsdbKey.AreaId, *lsaKey)
	fmt.Println("Lsa from db(lsa, ret) ", nlsa, ret)
	lsDbEnt.NetworkLsaMap[*lsaKey] = *networkLsa
	selfOrigLsaEnt[*lsaKey] = true
	ospf.regenerateLsa(lsdbKey, *lsaKey)

	summaryLsa := &SummaryLsa{}
	decodeSummaryLsa(lsa_summary, summaryLsa, lsaKey)
	symmaryLsa_byte := encodeSummaryLsa(*summaryLsa, *lsaKey)
	fmt.Println("Encoded summary LSA ", symmaryLsa_byte)
	slsa, ret := ospf.getSummaryLsaFromLsdb(lsdbKey.AreaId, *lsaKey)
	fmt.Println("Lsa from db(lsa, ret) ", slsa, ret)
	lsDbEnt.Summary3LsaMap[*lsaKey] = *summaryLsa
	selfOrigLsaEnt[*lsaKey] = true
	ospf.regenerateLsa(lsdbKey, *lsaKey)

	asexternalLsa := &ASExternalLsa{}
	decodeASExternalLsa(lsa_asExt, asexternalLsa, lsaKey)
	asexternalLsa_byte := encodeASExternalLsa(*asexternalLsa, *lsaKey)
	fmt.Println("Encoded as external LSA ", asexternalLsa_byte)
	alsa, ret := ospf.getASExternalLsaFromLsdb(lsdbKey.AreaId, *lsaKey)
	fmt.Println("Lsa from db(lsa, ret) ", alsa, ret)
	lsDbEnt.ASExternalLsaMap[*lsaKey] = *asexternalLsa
	selfOrigLsaEnt[*lsaKey] = true
	ospf.regenerateLsa(lsdbKey, *lsaKey)
	ospf.processAsExternalLSAFlood(*lsaKey)
	ospf.floodASExternalLsa(lsa_asExt)

	ospf.AreaLsdb[lsdbKey] = lsDbEnt
	ospf.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt

	ospf.processMaxAgeLSA(lsdbKey, lsDbEnt)
	ospf.AddLsdbEntry(val)
	ospf.DelLsdbEntry(val)
	ospf.AddOspfEventState("Del ls entry", "LSDB")
	ospf.lsdbStateRefresh()
	ospf.lsdbSelfLsaRefresh()
	ospf.processLSDatabaseTicker()
	maxAgeLsaMap = make(map[LsaKey][]byte)
	maxAgeMsg := maxAgeLsaMsg{
		lsaKey:   *lsaKey,
		msg_type: delMaxAgeLsa,
	}
	ospf.processMaxAgeLsaMsg(maxAgeMsg)

	/* AS External LSAs */
	ospf.HandleASExternalLsa(lsdbKey.AreaId)
	ospf.CalcASBorderRoutes(lsdbKey.AreaId)
	ospf.GenerateType4SummaryLSA(rKey, rEntry, lsdbKey)
	/* Flooding */
	ospf.generateDbsummary4LsaList(lsdbKey.AreaId)
	ospf.generateDbsummary3LsaList(lsdbKey.AreaId)
	ospf.generateDbasExternalList(lsdbKey.AreaId)
	ospf.generateDbSummaryList(nbrKey)
	ospf.SendSelfOrigLSA(lsdbKey.AreaId, key)
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSASELFLOOD
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSAINTF
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSAROUTERFLOOD
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSASUMMARYFLOOD
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSAEXTFLOOD
	ospf.processFloodMsg(floodMsg)
	floodMsg.lsOp = LSAAGE
	ospf.processFloodMsg(floodMsg)

	ospf.SendRouterLsa(lsdbKey.AreaId, key, routerKey)
	ospf.constructAndSendLsaAgeFlood()
	ospf.nbrFloodCheck(nbrKey, key, intf, RouterLSA)
	ospf.interfaceFloodCheck(key)
	ospf.processSummaryLSAFlood(lsdbKey.AreaId, routerKey)
	ospf.floodSummaryLsa(lsa_summary, lsdbKey.AreaId)

	ospf.CalcInterAreaRoutes(lsdbKey.AreaId)

	/* SPF */
	checkSPFAPIs(selfOrigLsaEnt)
	//	ospf.GenerateSummaryLsa()
}

func checkSPFAPIs(selfOrMap map[LsaKey]bool) {
	lsaKey, err := findSelfOrigRouterLsaKey(selfOrMap)
	fmt.Println("Found router Key ", lsaKey, " error ", err)

	network := link1.LinkId & netmask
	cons := ospf.checkRouterLsaConsistency(lsdbKey.AreaId, lsid, network, netmask)
	if cons {
		fmt.Println("Router LSA is not consitent with given lsid, network , netmask ",
			lsid, ", ", network, " ,", netmask)
	}

	var vkey VertexKey
	//start SPF calculation
	//ospf.spfCalculation()
	//fmt.Println(" Done initialising SPF tables.", done)
	ospf.initialiseSPFStructs()
	vkey, err = ospf.CreateAreaGraph(lsdbKey.AreaId)
	fmt.Println("Area graph created with vertex kay ", vkey)
	ospf.AreaStubs[vKeyR] = sVertex1
	ospf.AreaStubs[vKeyN] = sVertex2
	ospf.AreaStubs[vKeyT] = sVertex3

	fmt.Println("Handle stub networks ")
	ospf.HandleStubs(vKeyR, lsdbKey.AreaId)
	err = ospf.UpdateAreaGraphNetworkLsa(networkLsa, networkKey, lsdbKey.AreaId)
	fmt.Println("Area graph updated with error ", err)

	lsaKey, err = ospf.findNetworkLsa(lsdbKey.AreaId, networkKey.LSId)
	err = ospf.findRouterLsa(lsdbKey.AreaId, routerKey.LSId)

	ospf.UpdateAreaGraphRouterLsa(routerLsa, routerKey, lsdbKey.AreaId)
	ospf.UpdateRoutingTblForRouter(areaidkey, vKeyR, treeVertex, vKeyN)

	ip, nhip, success := ospf.findP2PNextHopIP(vKeyR, vKeyN, areaidkey)
	fmt.Println("Calculated P2P nexthop ip ", ip, nhip, success)
}

func initialiseSPFData() {
	ospf.AreaGraph = make(map[VertexKey]Vertex)
	ospf.SPFTree = make(map[VertexKey]TreeVertex)

	ospf.AreaGraph[vKeyR] = vertexR
	ospf.AreaGraph[vKeyN] = vertexN
	ospf.AreaGraph[vKeyT] = vertexT
}
