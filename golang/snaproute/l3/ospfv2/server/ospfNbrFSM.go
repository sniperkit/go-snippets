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

package server

import (
	"fmt"
	"l3/ospfv2/objects"
	"time"
)

func (server *OSPFV2Server) StartNbrFSM() {
	server.InitNbrStruct()
	go server.ProcessNbrFSM()
}

/* Handle neighbor events
 */
func (server *OSPFV2Server) ProcessNbrFSM() {
	for {
		select {
		//Hello packet recieved.
		case nbrData := <-server.MessagingChData.IntfToNbrFSMChData.NbrHelloEventCh:
			server.logger.Debug("Nbr : Received hello event. ", nbrData.NbrIP)
			server.ProcessNbrHello(nbrData)
			//DBD received
		case msg := <-server.MessagingChData.IntfToNbrFSMChData.NetworkDRChangeCh:
			server.logger.Info("Network DR Change Msg", msg)
			server.ProcessNetworkDRChangeMsg(msg)

		case dbdData := <-server.NbrConfData.neighborDBDEventCh:
			server.logger.Debug("Nbr: Received dbd event ", dbdData)
			server.ProcessNbrDbdMsg(dbdData)

		case lsaAckData := <-server.NbrConfData.nbrLsaAckEventCh:
			server.logger.Debug("Nbr: Received ack  ", lsaAckData)

			//LSAReq received
		case lsaReqData := <-server.NbrConfData.neighborLSAReqEventCh:
			nbr, exists := server.NbrConfMap[lsaReqData.nbrKey]
			if exists && nbr.State >= NbrExchange {
				server.ProcessLsaReq(lsaReqData)
			}

		case lsaData := <-server.NbrConfData.neighborLSAUpdEventCh:
			nbr, exists := server.NbrConfMap[lsaData.nbrKey]

			if exists && nbr.State >= NbrExchange {
				server.ProcessLsaUpd(lsaData)
			}

			//Intf state change
		case msg := <-server.MessagingChData.IntfToNbrFSMChData.DeleteNbrCh:
			_, exist := server.NbrConfData.IntfToNbrMap[msg.IntfKey]
			if !exist {
				server.logger.Info("Nbr : Intf down . No nbrs on this interface ", msg.IntfKey)
			} else {
				server.ProcessNbrDeadFromIntf(msg.IntfKey)
			}

			//NbrFsmCtrlCh
		case _ = <-server.NbrConfData.nbrFSMCtrlCh:
			server.logger.Debug("Nbr : FSM stopping.. ")
			server.DeinitNbrStruct()
			server.NbrConfData.nbrFSMCtrlReplyCh <- true
			return
		}
	}
}

func (server *OSPFV2Server) StopNbrFSM() {
	server.NbrConfData.nbrFSMCtrlCh <- true
	cnt := 0
	for {
		select {
		case _ = <-server.NbrConfData.nbrFSMCtrlReplyCh:
			server.logger.Info("Successfully Stopped Nbr FSM")
			server.NbrConfData.IntfToNbrMap = nil
			server.NbrConfData.neighborDBDEventCh = nil
			server.NbrConfData.neighborLSAUpdEventCh = nil
			server.NbrConfData.neighborLSAReqEventCh = nil
			server.NbrConfData.nbrLsaAckEventCh = nil

			return
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt == 100 {
				server.logger.Err("Unable to stop the  Neighbor FSM")
				return
			}
		}
	}

}

/**** handle neighbor states ***/
func (server *OSPFV2Server) ProcessNbrHello(nbrData NbrHelloEventMsg) {
	var oldState NbrState
	var newState NbrState
	oldState = NbrDown
	/*
		nbrKey := NbrConfKey{
			NbrIdentity:         nbrData.NbrIpAddr,
			NbrAddressLessIfIdx: 0,
		} */
	nbrKey := nbrData.NbrKey
	nbrConf, valid := server.NbrConfMap[nbrKey]
	if !valid {
		server.logger.Debug("Nbr: add new neighbor ", nbrData.NbrIP)
		// add new neighbor
		server.CreateNewNbr(nbrData)
		if nbrData.TwoWayStatus {
			newState = NbrTwoWay
		} else {
			newState = NbrInit
		}
	} else {
		oldState = nbrConf.State
		if nbrData.TwoWayStatus {
			newState = NbrTwoWay
		} else {
			newState = NbrInit
		}
		//nbrConf.NbrDeadTimer.Reset(nbrConf.NbrDeadTimeDuration)
	}
	server.logger.Debug("Nbr : oldstate", oldState, " newState ", newState, " state ", nbrConf.State)
	if (oldState == NbrDown || oldState == NbrInit) &&
		newState == NbrTwoWay {
		nbrConf.State = NbrTwoWay
		server.ProcessNbrTwoway(nbrKey)
	} else if newState == NbrInit &&
		(oldState != NbrInit) { //&& oldState != NbrDown) {
		nbrConf.State = NbrInit
		server.ProcessNbrInit(nbrKey)
	} else {
		server.ProcessNbrUpdate(nbrKey, nbrConf)
	}
}

func (server *OSPFV2Server) ProcessNbrDbdMsg(dbdMsg NbrDbdMsg) {
	nbrConf, exists := server.NbrConfMap[dbdMsg.nbrConfKey]
	if !exists {
		server.logger.Err("Nbr : ProcessNbrDbdMsg Nbrkey does not exist ", dbdMsg.nbrConfKey)
		return
	}
	switch nbrConf.State {
	case NbrInit, NbrExchangeStart:
		server.ProcessNbrExstart(dbdMsg.nbrConfKey, nbrConf, dbdMsg.nbrDbdData)
	case NbrExchange:
		server.ProcessNbrExchange(dbdMsg.nbrConfKey, nbrConf, dbdMsg.nbrDbdData)
	case NbrLoading:
		server.ProcessNbrLoading(dbdMsg.nbrConfKey, nbrConf, dbdMsg.nbrDbdData)
	case NbrFull:
		server.logger.Err("Nbr: Received dbd packet when nbr is full . Restart FSM", dbdMsg.nbrConfKey)
		nbrConf.State = NbrExchangeStart
		server.ProcessNbrExstart(dbdMsg.nbrConfKey, nbrConf, dbdMsg.nbrDbdData)
	case NbrDown:
		server.logger.Warning("Nbr: Nbr is down state. Dont process dbd ", dbdMsg.nbrConfKey)
	case NbrTwoWay:
		server.logger.Warning("Nbr: Nbr is two way state.Dont process dbd ", dbdMsg.nbrConfKey)
	}

}

func (server *OSPFV2Server) CreateNewNbr(nbrData NbrHelloEventMsg) {
	var nbrConf NbrConf
	nbrKey := nbrData.NbrKey
	nbrConf.NbrIP = nbrData.NbrIP
	nbrConf.NbrMac = nbrData.NbrMAC
	nbrConf.NbrDR = nbrData.NbrDRIpAddr
	nbrConf.NbrBdr = nbrData.NbrBDRIpAddr
	nbrConf.IntfKey = nbrData.IntfConfKey
	nbrConf.NbrRtrId = nbrData.RouterId
	nbrConf.NbrReqListIndex = 0
	nbrConf.NbrReqList = []*ospfLSAHeader{}
	nbrConf.NbrRetxList = []*ospfLSAHeader{}
	nbrConf.NbrDBSummaryList = []*ospfLSAHeader{}
	nbrConf.NbrDeadTimeDuration = nbrData.NbrDeadTime
	if nbrData.TwoWayStatus {
		nbrConf.State = NbrTwoWay
		server.NbrConfMap[nbrKey] = nbrConf
	} else {
		nbrConf.State = NbrInit
		server.NbrConfMap[nbrKey] = nbrConf
	}
	server.ProcessNbrDead(nbrKey)
	//	server.ProcessNbrFsmStart(nbrKey, nbrConf)
	server.logger.Debug("Nbr : Add to slice ", nbrKey)
	server.addNbrToSlice(nbrKey)
}

func (server *OSPFV2Server) ProcessNbrFsmStart(nbrKey NbrConfKey) {
	var dbd_mdata NbrDbdData
	server.logger.Debug("Nbr: Nbr fsm start ", nbrKey)
	nbrConf, _ := server.NbrConfMap[nbrKey]
	isAdjacent := server.AdjacencyCheck(nbrKey)
	if isAdjacent {
		nbrConf.State = NbrExchangeStart
		dbd_mdata.dd_sequence_number = uint32(time.Now().Nanosecond())
		// send dbd packets
		server.ConstructDbdMdata(nbrKey, true, true, true,
			INTF_OPTIONS, nbrConf.DDSequenceNum, false, false)
		server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
		server.logger.Debug("Nbr: Exstart seq ", dbd_mdata.dd_sequence_number)

	} else { // no adjacency
		server.logger.Debug("Nbr: Twoway  ", nbrKey)
		nbrConf.State = NbrTwoWay
	}

	server.ProcessNbrUpdate(nbrKey, nbrConf)

}

func (server *OSPFV2Server) ProcessNbrInit(nbrKey NbrConfKey) {
	nbrConf, exist := server.NbrConfMap[nbrKey]
	if !exist {
		server.logger.Err("Nbr : Does not exist. Init state ", nbrKey)
		return
	}
	nbrConf.NbrReqListIndex = -1
	nbrConf.NbrReqList = nil
	nbrConf.NbrRetxList = nil
	nbrConf.NbrDBSummaryList = nil
	nbrConf.NbrLsaIndex = -1
	server.ProcessNbrUpdate(nbrKey, nbrConf)
}

func (server *OSPFV2Server) ProcessNbrTwoway(nbrKey NbrConfKey) {
	server.ProcessNbrFsmStart(nbrKey)
}

func (server *OSPFV2Server) ProcessNbrExstart(nbrKey NbrConfKey, nbrConf NbrConf, nbrDbPkt NbrDbdData) {
	var dbd_mdata NbrDbdData
	var isAdjacent bool
	var negotiationDone bool
	isAdjacent = server.AdjacencyCheck(nbrKey)
	if isAdjacent || nbrConf.State == NbrExchangeStart {
		// decide master slave relation
		if nbrConf.NbrRtrId > server.globalData.RouterId {
			nbrConf.isMaster = true
		} else {
			nbrConf.isMaster = false
		}
		/* The initialize(I), more (M) and master(MS) bits are set,
		   the contents of the packet are empty, and the neighbor's
		   Router ID is larger than the router's own.  In this case
		   the router is now Slave.  Set the master/slave bit to
		   slave, and set the neighbor data structure's DD sequence
		   number to that specified by the master.
		*/
		server.logger.Debug("NBRDBD: nbr rtr id ", nbrConf.NbrRtrId,
			" my router id ", server.globalData.RouterId,
			" nbr_seq ", nbrConf.DDSequenceNum, "dbd_seq no ", nbrDbPkt.dd_sequence_number)
		if nbrDbPkt.ibit && nbrDbPkt.mbit && nbrDbPkt.msbit &&
			nbrConf.NbrRtrId > server.globalData.RouterId {
			server.logger.Debug("DBD: (ExStart/slave) SLAVE = self,  MASTER = ", nbrConf.NbrRtrId)
			nbrConf.isMaster = true
			server.logger.Debug("NBREVENT: Negotiation done..")
			negotiationDone = true
			nbrConf.State = NbrExchange
		}
		/*
			if nbrDbPkt.msbit && nbrConf.NbrRtrId > server.globalData.RouterId {
				server.logger.Debug("DBD: (ExStart/slave) SLAVE = self,  MASTER = ", nbrKey.NbrIdentity)
				nbrConf.isMaster = true
				server.logger.Debug("NBREVENT: Negotiation done..")
				negotiationDone = true
				nbrConf.State = NbrExchange
			} */

		/*   The initialize(I) and master(MS) bits are off, the
		     packet's DD sequence number equals the neighbor data
		     structure's DD sequence number (indicating
		     acknowledgment) and the neighbor's Router ID is smaller
		     than the router's own.  In this case the router is
		     Master.
		*/
		if nbrDbPkt.ibit == false && nbrDbPkt.msbit == false &&
			nbrDbPkt.dd_sequence_number == nbrConf.DDSequenceNum &&
			nbrConf.NbrRtrId < server.globalData.RouterId {
			nbrConf.isMaster = false
			server.logger.Debug("DBD:(ExStart) SLAVE = ", nbrKey.NbrIdentity, "MASTER = SELF")
			server.logger.Debug("NBREVENT: Negotiation done..")
			negotiationDone = true
			nbrConf.State = NbrExchange
		}

	} else {
		nbrConf.State = NbrTwoWay
	}

	if negotiationDone {
		//server.logger.Debug(fmt.Sprintln("DBD: (Exstart) lsa_headers = ", len(nbrDbPkt.lsa_headers)))
		server.generateDbSummaryList(nbrKey)

		if nbrConf.isMaster != true { // i am the master
			dbd_mdata, _ = server.ConstructDbdMdata(nbrKey, false, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number+1, true, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
		} else {
			// send acknowledgement DBD with I and MS bit false , mbit = 1
			dbd_mdata, _ = server.ConstructDbdMdata(nbrKey, false, true, false,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, true, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
			dbd_mdata.dd_sequence_number++
		}

		req_list := server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)
		server.logger.Debug("Nbr: received nbr req len ", len(req_list))
		nbrConf.NbrReqList = req_list
	} else { // negotiation not done
		server.logger.Debug("Nbr: Negotiation not done. ", nbrConf.NbrIP)
		nbrConf.State = NbrExchangeStart
		if nbrConf.isMaster &&
			nbrConf.NbrRtrId > server.globalData.RouterId {
			dbd_mdata.dd_sequence_number = nbrDbPkt.dd_sequence_number
			dbd_mdata, _ = server.ConstructDbdMdata(nbrKey, true, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
		} else {
			dbd_mdata, _ = server.ConstructDbdMdata(nbrKey, true, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
		}
	}
	nbrConf.DDSequenceNum = nbrDbPkt.dd_sequence_number
	nbrConf.NbrOption = uint32(nbrDbPkt.options)
	server.ProcessNbrUpdate(nbrKey, nbrConf)

}

func (server *OSPFV2Server) ProcessNbrExchange(nbrKey NbrConfKey, nbrConf NbrConf, nbrDbPkt NbrDbdData) {
	var last_exchange bool
	var dbd_mdata NbrDbdData
	isDiscard := server.NbrDbPacketDiscardCheck(nbrDbPkt, nbrConf)
	if isDiscard {
		server.logger.Debug(fmt.Sprintln("NBRDBD: (Exchange)Discard packet. nbr", nbrConf.NbrIP,
			" nbr state ", nbrConf.State))

		nbrConf.State = NbrExchangeStart
		server.ProcessNbrExstart(nbrKey, nbrConf, nbrDbPkt)

		return
	} else { // process exchange state
		/* 2) Add lsa_headers to db packet from db_summary list */

		if nbrConf.isMaster != true { // i am master
			/* Send the DBD only if packet has mbit =1 or event != NbrExchangeDone
			          send DBD with seq num + 1 , ibit = 0 ,  ms = 1
			   * if this is the last DBD for LSA description set mbit = 0
			*/
			server.logger.Debug(fmt.Sprintln("DBD:(master/Exchange) mbit ", nbrDbPkt.mbit))
			if nbrDbPkt.mbit {
				server.logger.Debug(fmt.Sprintln("DBD: (master/Exchange) Send next packet in the exchange  to nbr ", nbrKey.NbrIdentity))
				dbd_mdata, _ := server.ConstructDbdMdata(nbrKey, false, false, true,
					nbrDbPkt.options, nbrDbPkt.dd_sequence_number+1, true, false)
				server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
				nbrConf.NbrLastDbd = dbd_mdata
			}

			// Genrate request list
			req_list := server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)
			nbrConf.NbrReqList = append(nbrConf.NbrReqList, req_list...)
			server.logger.Debug("DBD:(Exchange) Total elements in req_list ", len(nbrConf.NbrReqList))

		} else { // i am slave
			/* send acknowledgement DBD with I and MS bit false and mbit same as
			   rx packet
			    if mbit is 0 && last_exchange == true generate NbrExchangeDone*/
			server.logger.Debug(fmt.Sprintln("DBD: (slave/Exchange) Send next packet in the exchange  to nbr ", nbrKey.NbrIdentity))
			req_list := server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)
			nbrConf.NbrReqList = append(nbrConf.NbrReqList, req_list...)
			dbd_mdata, last_exchange = server.ConstructDbdMdata(nbrKey, false, nbrDbPkt.mbit, false,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, true, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
			nbrConf.NbrLastDbd = dbd_mdata
			dbd_mdata.dd_sequence_number++
		}
		if !nbrDbPkt.mbit && last_exchange {
			nbrConf.State = NbrLoading
			nbrConf.NbrReqListIndex = server.BuildAndSendLSAReq(nbrKey, nbrConf)
			server.logger.Debug(fmt.Sprintln("DBD: Loading , nbr ", nbrKey.NbrIdentity))
		}

	}
	nbrConf.DDSequenceNum = nbrDbPkt.dd_sequence_number
	server.ProcessNbrUpdate(nbrKey, nbrConf)

}

func (server *OSPFV2Server) ProcessNbrLoading(nbrKey NbrConfKey, nbrConf NbrConf, nbrDbPkt NbrDbdData) {
	var seq_num uint32
	server.logger.Debug(fmt.Sprintln("DBD: Loading . Nbr ", nbrKey.NbrIdentity))
	isDiscard := server.NbrDbPacketDiscardCheck(nbrDbPkt, nbrConf)
	isDuplicate := server.verifyDuplicatePacket(nbrConf, nbrDbPkt)
	nbrConf.State = NbrLoading
	if isDiscard {
		server.logger.Debug(fmt.Sprintln("NBRDBD:Loading  Discard packet. nbr", nbrKey.NbrIdentity,
			" nbr state ", nbrConf.State))
		//update neighbor to exchange start state and send dbd

		nbrConf.State = NbrExchangeStart
		nbrConf.isMaster = false
		/*
		dbd_mdata, _ = server.ConstructDbdMdata(nbrKey, true, true, true,
			nbrDbPkt.options, nbrConf.DDSequenceNum+1, false, false)
		server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
		seq_num = dbd_mdata.dd_sequence_number */
		nbrConf.State = NbrExchangeStart
                server.ProcessNbrExstart(nbrKey, nbrConf, nbrDbPkt)

	} else if !isDuplicate {
		/*
		   slave - Send the old dbd packet.
		       master - discard
		*/
		if nbrConf.isMaster {
			dbd_mdata, _ := server.ConstructDbdMdata(nbrKey, false, nbrDbPkt.mbit, false,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false)
			server.BuildAndSendDdBDPkt(nbrConf, dbd_mdata)
			seq_num = dbd_mdata.dd_sequence_number + 1
		}
		seq_num = nbrConf.NbrLastDbd.dd_sequence_number
	} else {
		seq_num = nbrConf.NbrLastDbd.dd_sequence_number
	}
	nbrConf.DDSequenceNum = seq_num
	server.ProcessNbrUpdate(nbrKey, nbrConf)
}

func (server *OSPFV2Server) ProcessNbrFull(nbrKey NbrConfKey) {
	nbrConf, valid := server.NbrConfMap[nbrKey]
	if !valid {
		server.logger.Err("Nbr: Full event , nbr key does not exist ", nbrKey)
	}
	nbrConf.State = NbrFull
	server.UpdateIntfToNbrMap(nbrKey)
	server.ProcessNbrUpdate(nbrKey, nbrConf)
	server.logger.Debug("Nbr: Nbr full event ", nbrKey)
	intf, valid := server.IntfConfMap[nbrConf.IntfKey]
	if !valid {
		server.logger.Err("Nbr : Intf does not exist. ", nbrKey)
		return
	}
	server.logger.Debug("Nbr : intf rtr ", intf.DRtrId, " global rtr ", server.globalData.RouterId)
	if intf.DRtrId == server.globalData.RouterId {
		server.logger.Debug("Nbr : Send message to lsdb to generate nw lsa ", nbrConf.NbrIP)
		nbrList := []uint32{}
		for _, nbr := range server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey] {
			nbrC, exist := server.NbrConfMap[nbr]
			if !exist {
				server.logger.Info("Nbr : Generate nw lsa . nbr does not exist ", nbr)
				continue
			}
			nbrList = append(nbrList, nbrC.NbrRtrId)
		}
		msg := UpdateSelfNetworkLSAMsg{
			Op:      GENERATE,
			IntfKey: nbrConf.IntfKey,
			NbrList: nbrList,
		}

		server.SendMsgFromNbrToLsdb(msg)
	}
	floodMsg := NbrToFloodMsg{
		NbrKey:  nbrKey,
		MsgType: LSA_FLOOD_NBR_FULL,
	}
	server.MessagingChData.NbrFSMToFloodChData.LsaFloodCh <- floodMsg
}

func (server *OSPFV2Server) ProcessNetworkDRChangeMsg(msg NetworkDRChangeMsg) {
	server.logger.Debug("Nbr: Received Network DR change message ")
	var Op LsaOp
	if msg.OldIntfFSMState == objects.INTF_FSM_STATE_DR &&
		msg.NewIntfFSMState == objects.INTF_FSM_STATE_OTHER_DR {
		Op = FLUSH
	} else if msg.OldIntfFSMState == objects.INTF_FSM_STATE_OTHER_DR &&
		msg.NewIntfFSMState == objects.INTF_FSM_STATE_DR {
		Op = GENERATE
	} else {
		return
	}

	nbrList := []uint32{}
	for _, nbr := range server.NbrConfData.IntfToNbrMap[msg.IntfKey] {
		nbrC, exist := server.NbrConfMap[nbr]
		if !exist {
			server.logger.Info("Nbr : Generate nw lsa . nbr does not exist ", nbr)
			continue
		}
		nbrList = append(nbrList, nbrC.NbrRtrId)
	}

	floodMsg := UpdateSelfNetworkLSAMsg{
		Op:      Op,
		IntfKey: msg.IntfKey,
		NbrList: nbrList,
	}

	server.SendMsgFromNbrToLsdb(floodMsg)

}

func (server *OSPFV2Server) ProcessNbrDeadFromIntf(key IntfConfKey) {
	server.logger.Debug("Nbr : intf down. Process nbr down", key)
	nbrList, valid := server.NbrConfData.IntfToNbrMap[key]
	if !valid {
		server.logger.Info("Nbr : No nbr on this interface. ", key)
		return
	}

	for _, nbr := range nbrList {
		nbrConf, exist := server.NbrConfMap[nbr]
		if !exist {
			continue
		}
		nbrConf.NbrDeadTimer.Stop()
		nbrConf.NbrDeadTimer = nil
		if len(nbrConf.NbrReqList) > 0 {
			nbrConf.NbrReqList = nbrConf.NbrReqList[:len(nbrConf.NbrReqList)-1]
		}
		if len(nbrConf.NbrRetxList) > 0 {
			nbrConf.NbrRetxList = nbrConf.NbrRetxList[:len(nbrConf.NbrRetxList)-1]
		}
		if len(nbrConf.NbrDBSummaryList) > 0 {
			nbrConf.NbrDBSummaryList = nbrConf.NbrDBSummaryList[:len(nbrConf.NbrDBSummaryList)-1]
		}

		nbrConf.NbrReqList = nil
		nbrConf.NbrRetxList = nil
		nbrConf.NbrDBSummaryList = nil
		delete(server.NbrConfMap, nbr)
		server.delNbrFromSlice(nbr)
		server.logger.Info("Nbr: Deleted", nbr)
	}

	intf, valid := server.IntfConfMap[key]
	if !valid {
		server.logger.Info("Nbr : intf does not exist . Dont send msg to lsdb", key)
	} else {
		if intf.DRtrId == server.globalData.RouterId {
			lsdbMsg := UpdateSelfNetworkLSAMsg{
				Op:      FLUSH,
				IntfKey: key,
			}

			server.SendMsgFromNbrToLsdb(lsdbMsg)
		}
	}
	server.logger.Debug("nbr : Intf down processing done. ")
}

func (server *OSPFV2Server) ProcessNbrDead(nbrKey NbrConfKey) {
	server.logger.Debug("Nbr : nbr dead called")
	var nbr_entry_dead_func func()
	nbr_entry_dead_func = func() {
		nbrConf, _ := server.NbrConfMap[nbrKey]

		server.logger.Info(fmt.Sprintln("NBRSCAN: DEAD ", nbrKey))
		server.logger.Info(fmt.Sprintln("DEAD: start processing nbr dead ", nbrKey))
		server.ResetNbrData(nbrKey, nbrConf.IntfKey)

		server.logger.Info(fmt.Sprintln("DEAD: end processing nbr dead ", nbrKey))

		nbrConf, exists := server.NbrConfMap[nbrKey]
		if exists {
			//update interface to neighbor map
			nbrList, valid := server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey]
			if valid {
				for i, nbrKeyT := range nbrList {
					if nbrKeyT.NbrIdentity == nbrKey.NbrIdentity {
						nbrList = append(nbrList[:i], nbrList[i+1:]...)
						break
					}
				}
				server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey] = nbrList
				server.logger.Debug("Nbr : Int to nbr list updated ", nbrList)
			}
			//send signal to intf fsm
			nbrDownMsg := NbrDownMsg{
				NbrKey: nbrKey,
			}
			server.MessagingChData.NbrToIntfFSMChData.NbrDownMsgChMap[nbrConf.IntfKey] <- nbrDownMsg
			//Send Msg to Lsdb for Nbr Dead
			intfConfEnt, exist := server.IntfConfMap[nbrConf.IntfKey]
			if exist {
				nbrDeadMsg := NbrDeadMsg{
					AreaId:   intfConfEnt.AreaId,
					NbrRtrId: nbrConf.NbrRtrId,
				}
				server.SendMsgToLsdbFromNbrFSMForNbrDead(nbrDeadMsg)
			}
			//send message to lsdb if I am DR.
			intf, valid := server.IntfConfMap[nbrConf.IntfKey]
			if !valid {
				server.logger.Info("Nbr : intf does not exist . Dont send msg to lsdb", nbrConf.IntfKey)
			} else {
				if intf.DRtrId == server.globalData.RouterId {
					nbrList := []uint32{}
					for _, nbr := range server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey] {
						nbrC, exist := server.NbrConfMap[nbr]
						if !exist {
							server.logger.Info("Nbr : Generate nw lsa . nbr does not exist ", nbr)
							continue
						}
						nbrList = append(nbrList, nbrC.NbrRtrId)
					}

					lsdbMsg := UpdateSelfNetworkLSAMsg{
						Op:      GENERATE,
						IntfKey: nbrConf.IntfKey,
						NbrList: nbrList,
					}

					server.SendMsgFromNbrToLsdb(lsdbMsg)
				}

			}
			if len(nbrConf.NbrReqList) > 0 {
				nbrConf.NbrReqList = nbrConf.NbrReqList[:len(nbrConf.NbrReqList)-1]
			}
			if len(nbrConf.NbrRetxList) > 0 {
				nbrConf.NbrRetxList = nbrConf.NbrRetxList[:len(nbrConf.NbrRetxList)-1]
			}
			if len(nbrConf.NbrDBSummaryList) > 0 {
				nbrConf.NbrDBSummaryList = nbrConf.NbrDBSummaryList[:len(nbrConf.NbrDBSummaryList)-1]
			}

			nbrConf.NbrReqList = nil
			nbrConf.NbrRetxList = nil
			nbrConf.NbrDBSummaryList = nil
			//delete neighbor from map
			delete(server.NbrConfMap, nbrKey)
			server.logger.Info("Nbr: Deleted ", nbrKey)
		}
	} // end of afterFunc callback
	nbrConf, exists := server.NbrConfMap[nbrKey]
	if exists {
		nbrConf.NbrDeadTimer = time.AfterFunc(nbrConf.NbrDeadTimeDuration, nbr_entry_dead_func)
		server.NbrConfMap[nbrKey] = nbrConf
		server.logger.Debug("Nbr : nbr dead updated ")
	}

}

func (server *OSPFV2Server) ProcessNbrUpdate(nbrKey NbrConfKey, nbrConf NbrConf) {
	server.logger.Debug("Nbr : ", nbrConf)
	if nbrConf.NbrDeadTimer != nil {
		nbrConf.NbrDeadTimer.Reset(nbrConf.NbrDeadTimeDuration)
	}
	server.NbrConfMap[nbrKey] = nbrConf
	server.logger.Debug("Nbr: Nbr conf updated ", nbrKey)
}

/**** Utils APis *****/
func (server *OSPFV2Server) dbPacketDiscardCheck() bool {
	return false
}
func (server *OSPFV2Server) AdjacencyCheck(nbrKey NbrConfKey) bool {
	_, valid := server.NbrConfMap[nbrKey]
	if !valid {
		server.logger.Err("Nbr : Nbr does not exist . No adjacency.", nbrKey)
		return false
	}
	/*
			 o   The underlying network type is point-to-point

		        o   The underlying network type is Point-to-MultiPoint

		        o   The underlying network type is virtual link

		        o   The router itself is the Designated Router

		        o   The router itself is the Backup Designated Router

		        o   The neighboring router is the Designated Router

		        o   The neighboring router is the Backup Designated Router
	*/
	return true
}

func (server *OSPFV2Server) ResetNbrData(nbr NbrConfKey, intf IntfConfKey) {
	/* List of Neighbors per interface instance */
	nbrConf, exist := server.NbrConfMap[nbr]
	if exist {
		nbrConf.NbrReqList = nil
		nbrConf.NbrDBSummaryList = nil
		nbrConf.NbrRetxList = nil
	}

	nbrList, exists := server.NbrConfData.IntfToNbrMap[intf]
	if !exists {
		server.logger.Info(fmt.Sprintln("DEAD: Nbr dead but intf-to-nbr map doesnt exist. ", nbr))
		return
	}
	newList := []NbrConfKey{}
	for inst := range nbrList {
		if nbrList[inst] != nbr {
			newList = append(newList, nbr)
		}
	}
	server.NbrConfData.IntfToNbrMap[intf] = newList
	server.logger.Info(fmt.Sprintln("Nbr: nbrList ", newList))
}
