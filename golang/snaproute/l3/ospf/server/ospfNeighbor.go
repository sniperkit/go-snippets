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
	"time"
)

/* @fn exchangePacketDiscardCheck
    Function to check SeqNumberMismatch
	for dbd exchange state packets.
*/
func (server *OSPFServer) exchangePacketDiscardCheck(nbrConf OspfNeighborEntry, nbrDbPkt ospfDatabaseDescriptionData) (isDiscard bool) {
	msg := DbEventMsg{
		eventType: config.ADJACENCY,
	}

	if nbrDbPkt.msbit != nbrConf.isMaster {
		msg.eventInfo = "SeqNumberMismatch. Nbr should be master " + nbrConf.OspfNbrIPAddr.String()
		server.logger.Info(fmt.Sprintln("NBREVENT: SeqNumberMismatch. Nbr should be master  dbdmsbit ", nbrDbPkt.msbit,
			" isMaster ", nbrConf.isMaster))
		//server.DbEventOp <- msg
		return true
	}

	if nbrDbPkt.ibit == true {
		server.logger.Info("NBREVENT:SeqNumberMismatch . Nbr ibit is true ")
		msg.eventInfo = "SeqNumberMismatch . Nbr ibit is true " + nbrConf.OspfNbrIPAddr.String()
		//server.DbEventOp <- msg
		return true
	}
	/*
		if nbrDbPkt.options != INTF_OPTIONS {
			server.logger.Info(fmt.Sprintln("NBREVENT:SeqNumberMismatch. Nbr options dont match. Nbr options ", INTF_OPTIONS,
				" dbd oackts options", nbrDbPkt.options))
			return true
		}*/

	if nbrConf.isMaster {
		if nbrDbPkt.dd_sequence_number != nbrConf.ospfNbrSeqNum {
			if nbrDbPkt.dd_sequence_number+1 == nbrConf.ospfNbrSeqNum {
				server.logger.Debug(fmt.Sprintln("Duplicate: This is db duplicate packet. Ignore."))
				return false
			}
			server.logger.Info(fmt.Sprintln("NBREVENT:SeqNumberMismatch : Nbr is master but dbd packet seq no doesnt match. dbd seq ",
				nbrDbPkt.dd_sequence_number, "nbr seq ", nbrConf.ospfNbrSeqNum))
			msg.eventInfo = "SeqNumberMismatch . " + nbrConf.OspfNbrIPAddr.String()
			//	server.DbEventOp <- msg
			return true
		}
	} else {
		if nbrDbPkt.dd_sequence_number != nbrConf.ospfNbrSeqNum {
			server.logger.Info(fmt.Sprintln("NBREVENT:SeqNumberMismatch : Nbr is slave but dbd packet seq no doesnt match.dbd seq ",
				nbrDbPkt.dd_sequence_number, "nbr seq ", nbrConf.ospfNbrSeqNum))
			msg.eventInfo = "SeqNumberMismatch . " + nbrConf.OspfNbrIPAddr.String()
			//	server.DbEventOp <- msg
			return true
		}
	}

	return false
}

func (server *OSPFServer) verifyDuplicatePacket(nbrConf OspfNeighborEntry, nbrDbPkt ospfDatabaseDescriptionData) (isDup bool) {
	if nbrConf.isMaster {
		if nbrDbPkt.dd_sequence_number+1 == nbrConf.ospfNbrSeqNum {
			isDup = true
			server.logger.Info(fmt.Sprintln("NBREVENT: Duplicate packet Dont do anything. dbdseq ",
				nbrDbPkt.dd_sequence_number, " nbrseq ", nbrConf.ospfNbrSeqNum))
			return
		}
	}
	isDup = false
	return
}
func (server *OSPFServer) adjacancyEstablishementCheck(isNbrDRBDR bool, isRtrDRBDR bool) (result bool) {
	if isNbrDRBDR || isRtrDRBDR {
		return true
	}
	/* TODO - check if n/w is p2p , p2mp, virtual link */
	return false
}

func (server *OSPFServer) processNeighborExstart(nbrKey NeighborConfKey,
	nbrConf OspfNeighborEntry, nbrDbPkt ospfDatabaseDescriptionData, ifMtu int32) {
	var dbd_mdata ospfDatabaseDescriptionData
	last_exchange := true
	var isAdjacent bool
	var negotiationDone bool
	isAdjacent = server.adjacancyEstablishementCheck(nbrConf.isDRBDR, true)
	if isAdjacent || nbrConf.OspfNbrState == config.NbrExchangeStart {
		// change nbr state
		nbrConf.OspfNbrState = config.NbrExchangeStart
		// decide master slave relation
		if nbrConf.OspfNbrRtrId > binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId) {
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
		server.logger.Debug(fmt.Sprintln("NBRDBD: nbr ip ", nbrKey.IPAddr,
			" my router id ", binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId),
			" nbr_seq ", nbrConf.ospfNbrSeqNum, "dbd_seq no ", nbrDbPkt.dd_sequence_number))
		if nbrDbPkt.ibit && nbrDbPkt.mbit && nbrDbPkt.msbit &&
			nbrConf.OspfNbrRtrId > binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId) {
			server.logger.Debug(fmt.Sprintln("DBD: (ExStart/slave) SLAVE = self,  MASTER = ", nbrKey.IPAddr))
			nbrConf.isMaster = true
			server.logger.Debug("NBREVENT: Negotiation done..")
			negotiationDone = true
			nbrConf.OspfNbrState = config.NbrExchange
			nbrConf.nbrEvent = config.NbrNegotiationDone
		}
		if nbrDbPkt.msbit && nbrConf.OspfNbrRtrId > binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId) {
			server.logger.Debug(fmt.Sprintln("DBD: (ExStart/slave) SLAVE = self,  MASTER = ", nbrKey.IPAddr))
			nbrConf.isMaster = true
			server.logger.Debug("NBREVENT: Negotiation done..")
			negotiationDone = true
			nbrConf.OspfNbrState = config.NbrExchange
			nbrConf.nbrEvent = config.NbrNegotiationDone
		}

		/*   The initialize(I) and master(MS) bits are off, the
		     packet's DD sequence number equals the neighbor data
		     structure's DD sequence number (indicating
		     acknowledgment) and the neighbor's Router ID is smaller
		     than the router's own.  In this case the router is
		     Master.
		*/
		if nbrDbPkt.ibit == false && nbrDbPkt.msbit == false &&
			nbrDbPkt.dd_sequence_number == nbrConf.ospfNbrSeqNum &&
			nbrConf.OspfNbrRtrId < binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId) {
			nbrConf.isMaster = false
			server.logger.Debug(fmt.Sprintln("DBD:(ExStart) SLAVE = ", nbrKey.IPAddr, "MASTER = SELF"))
			server.logger.Debug("NBREVENT: Negotiation done..")
			negotiationDone = true
			nbrConf.OspfNbrState = config.NbrExchange
			nbrConf.nbrEvent = config.NbrNegotiationDone
		}

	} else {
		nbrConf.OspfNbrState = config.NbrTwoWay
	}

	var lsa_attach uint8
	if negotiationDone {
		//server.logger.Debug(fmt.Sprintln("DBD: (Exstart) lsa_headers = ", len(nbrDbPkt.lsa_headers)))
		server.generateDbSummaryList(nbrKey)
		if nbrConf.isMaster != true { // i am the master
			dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, false, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number+1, true, false, ifMtu)
		} else {
			// send acknowledgement DBD with I and MS bit false , mbit = 1
			dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, false, true, false,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, true, false, ifMtu)
			dbd_mdata.dd_sequence_number++
		}

		if last_exchange {
			nbrConf.nbrEvent = config.NbrExchangeDone
		}
		server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)

	} else { // negotiation not done
		nbrConf.OspfNbrState = config.NbrExchangeStart
		if nbrConf.isMaster &&
			nbrConf.OspfNbrRtrId > binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId) {
			dbd_mdata.dd_sequence_number = nbrDbPkt.dd_sequence_number
			dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, true, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false, ifMtu)
			dbd_mdata.dd_sequence_number++
		} else {
			//start with new seq number
			dbd_mdata.dd_sequence_number = uint32(time.Now().Nanosecond()) //nbrConf.ospfNbrSeqNum
			dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, true, true, true,
				nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false, ifMtu)
		}
	}

	nbrConfMsg := ospfNeighborConfMsg{
		ospfNbrConfKey: nbrKey,
		ospfNbrEntry: OspfNeighborEntry{
			OspfNbrRtrId:           nbrConf.OspfNbrRtrId,
			OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
			OspfRtrPrio:            nbrConf.OspfRtrPrio,
			intfConfKey:            nbrConf.intfConfKey,
			OspfNbrOptions:         0,
			OspfNbrState:           nbrConf.OspfNbrState,
			isStateUpdate:          true,
			OspfNbrInactivityTimer: time.Now(),
			OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
			ospfNbrSeqNum:          dbd_mdata.dd_sequence_number,
			isSeqNumUpdate:         true,
			isMaster:               nbrConf.isMaster,
			isMasterUpdate:         true,
			nbrEvent:               nbrConf.nbrEvent,
			ospfNbrLsaIndex:        nbrConf.ospfNbrLsaIndex + lsa_attach,
		},
		nbrMsgType: NBRUPD,
	}
	server.neighborConfCh <- nbrConfMsg
	OspfNeighborLastDbd[nbrKey] = dbd_mdata
}

func (server *OSPFServer) processDBDEvent(nbrKey NeighborConfKey, nbrDbPkt ospfDatabaseDescriptionData) {
	_, exists := server.NeighborConfigMap[nbrKey]
	var dbd_mdata ospfDatabaseDescriptionData
	last_exchange := true
	if exists {
		nbrConf := server.NeighborConfigMap[nbrKey]
		intfConf, _ := server.IntfConfMap[nbrConf.intfConfKey]
		switch nbrConf.OspfNbrState {
		case config.NbrAttempt:
			/* reject packet */
			return
		case config.NbrInit, config.NbrExchangeStart:
			server.processNeighborExstart(nbrKey, nbrConf, nbrDbPkt, intfConf.IfMtu)

		case config.NbrExchange:
			var nbrState config.NbrState
			isDiscard := server.exchangePacketDiscardCheck(nbrConf, nbrDbPkt)
			if isDiscard {
				server.logger.Debug(fmt.Sprintln("NBRDBD: (Exchange)Discard packet. nbr", nbrKey.IPAddr,
					" nbr state ", nbrConf.OspfNbrState))

				nbrState = config.NbrExchangeStart
				nbrConf.OspfNbrState = config.NbrExchangeStart
				server.processNeighborExstart(nbrKey, nbrConf, nbrDbPkt, intfConf.IfMtu)

				//invalidate all lists.
				newDbdMsg(nbrKey, OspfNeighborLastDbd[nbrKey])
				return
			} else { // process exchange state
				/* 2) Add lsa_headers to db packet from db_summary list */

				if nbrConf.isMaster != true { // i am master
					/* Send the DBD only if packet has mbit =1 or event != NbrExchangeDone
						send DBD with seq num + 1 , ibit = 0 ,  ms = 1
					 * if this is the last DBD for LSA description set mbit = 0
					*/
					server.logger.Debug(fmt.Sprintln("DBD:(master/Exchange) nbr_event ", nbrConf.nbrEvent, " mbit ", nbrDbPkt.mbit))
					if nbrDbPkt.dd_sequence_number == nbrConf.ospfNbrSeqNum &&
						(nbrConf.nbrEvent != config.NbrExchangeDone ||
							nbrDbPkt.mbit) {
						server.logger.Debug(fmt.Sprintln("DBD: (master/Exchange) Send next packet in the exchange  to nbr ", nbrKey.IPAddr))
						dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, false, false, true,
							nbrDbPkt.options, nbrDbPkt.dd_sequence_number+1, true, false, intfConf.IfMtu)
						OspfNeighborLastDbd[nbrKey] = dbd_mdata
					}

					// Genrate request list
					server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)
					server.logger.Debug(fmt.Sprintln("DBD:(Exchange) Total elements in req_list ", len(ospfNeighborRequest_list[nbrKey])))

				} else { // i am slave
					/* send acknowledgement DBD with I and MS bit false and mbit same as
					rx packet
					 if mbit is 0 && last_exchange == true generate NbrExchangeDone*/
					if nbrDbPkt.dd_sequence_number == nbrConf.ospfNbrSeqNum {
						server.logger.Debug(fmt.Sprintln("DBD: (slave/Exchange) Send next packet in the exchange  to nbr ", nbrKey.IPAddr))
						server.generateRequestList(nbrKey, nbrConf, nbrDbPkt)
						dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, false, nbrDbPkt.mbit, false,
							nbrDbPkt.options, nbrDbPkt.dd_sequence_number, true, false, intfConf.IfMtu)
						OspfNeighborLastDbd[nbrKey] = dbd_mdata
						dbd_mdata.dd_sequence_number++
					} else {
						server.logger.Debug(fmt.Sprintln("DBD: (slave/exchange) Duplicated dbd.  . dbd_seq , nbr_seq_num ",
							nbrDbPkt.dd_sequence_number, nbrConf.ospfNbrSeqNum))
						if !nbrDbPkt.msbit && !nbrDbPkt.ibit {
							// the last exchange packet so we need not send duplicate response
							last_exchange = true
						}
						// send old ACK
						data := newDbdMsg(nbrKey, OspfNeighborLastDbd[nbrKey])
						server.ospfNbrDBDSendCh <- data

						dbd_mdata = OspfNeighborLastDbd[nbrKey]

					}
					if !nbrDbPkt.mbit && last_exchange {
						nbrConf.nbrEvent = config.NbrExchangeDone
					}
				}
				if !nbrDbPkt.mbit || last_exchange {
					server.logger.Debug(fmt.Sprintln("DBD: Exchange done with nbr ", nbrKey.IPAddr))
					nbrState = config.NbrLoading
					server.lsaReTxTimerCheck(nbrKey)
					if !nbrConf.isMaster {
						server.updateNeighborMdata(nbrConf.intfConfKey, nbrKey)
						server.CreateNetworkLSACh <- ospfIntfToNbrMap[nbrConf.intfConfKey]

					}
				}
				if !nbrDbPkt.mbit && last_exchange {
					nbrState = config.NbrLoading
					nbrConf.ospfNbrLsaReqIndex = server.BuildAndSendLSAReq(nbrKey, nbrConf)
					server.logger.Debug(fmt.Sprintln("DBD: Loading , nbr ", nbrKey.IPAddr))
					server.updateNeighborMdata(nbrConf.intfConfKey, nbrKey)
				}
			}

			nbrConfMsg := ospfNeighborConfMsg{
				ospfNbrConfKey: nbrKey,
				ospfNbrEntry: OspfNeighborEntry{
					OspfNbrRtrId:           nbrConf.OspfNbrRtrId,
					OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
					OspfRtrPrio:            nbrConf.OspfRtrPrio,
					intfConfKey:            nbrConf.intfConfKey,
					OspfNbrOptions:         0,
					OspfNbrState:           nbrState,
					isStateUpdate:          true,
					OspfNbrInactivityTimer: time.Now(),
					OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
					ospfNbrSeqNum:          dbd_mdata.dd_sequence_number,
					isSeqNumUpdate:         true,
					isMasterUpdate:         false,
					nbrEvent:               nbrConf.nbrEvent,
					ospfNbrLsaReqIndex:     nbrConf.ospfNbrLsaReqIndex,
				},
				nbrMsgType: NBRUPD,
			}
			server.neighborConfCh <- nbrConfMsg

		case config.NbrLoading, config.NbrFull:

			var seq_num uint32
			server.logger.Debug(fmt.Sprintln("DBD: Loading . Nbr ", nbrKey.IPAddr))
			isDiscard := server.exchangePacketDiscardCheck(nbrConf, nbrDbPkt)
			isDuplicate := server.verifyDuplicatePacket(nbrConf, nbrDbPkt)

			if isDiscard {
				server.logger.Debug(fmt.Sprintln("NBRDBD:Loading  Discard packet. nbr", nbrKey.IPAddr,
					" nbr state ", nbrConf.OspfNbrState))
				//update neighbor to exchange start state and send dbd

				nbrConf.OspfNbrState = config.NbrExchangeStart
				nbrConf.nbrEvent = config.Nbr2WayReceived
				nbrConf.isMaster = false
				dbd_mdata, last_exchange = server.ConstructAndSendDbdPacket(nbrKey, true, true, true,
					nbrDbPkt.options, nbrConf.ospfNbrSeqNum+1, false, false, intfConf.IfMtu)
				seq_num = dbd_mdata.dd_sequence_number
			} else if !isDuplicate {
				/*
					    slave - Send the old dbd packet.
						master - discard
				*/
				if nbrConf.isMaster {
					dbd_mdata, _ := server.ConstructAndSendDbdPacket(nbrKey, false, nbrDbPkt.mbit, false,
						nbrDbPkt.options, nbrDbPkt.dd_sequence_number, false, false, intfConf.IfMtu)
					seq_num = dbd_mdata.dd_sequence_number + 1
				}
				nbrConf.ospfNbrLsaReqIndex = server.BuildAndSendLSAReq(nbrKey, nbrConf)
				seq_num = OspfNeighborLastDbd[nbrKey].dd_sequence_number
				nbrConf.OspfNbrState = config.NbrFull
			} else {

				nbrConf.ospfNbrLsaReqIndex = server.BuildAndSendLSAReq(nbrKey, nbrConf)
				seq_num = OspfNeighborLastDbd[nbrKey].dd_sequence_number
				nbrConf.OspfNbrState = config.NbrFull
				server.updateNeighborMdata(nbrConf.intfConfKey, nbrKey)
				server.CreateNetworkLSACh <- ospfIntfToNbrMap[nbrConf.intfConfKey]
			}

			nbrConfMsg := ospfNeighborConfMsg{
				ospfNbrConfKey: nbrKey,
				ospfNbrEntry: OspfNeighborEntry{
					OspfNbrRtrId:           nbrConf.OspfNbrRtrId,
					OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
					OspfRtrPrio:            nbrConf.OspfRtrPrio,
					intfConfKey:            nbrConf.intfConfKey,
					OspfNbrOptions:         0,
					OspfNbrState:           nbrConf.OspfNbrState,
					isStateUpdate:          true,
					OspfNbrInactivityTimer: time.Now(),
					OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
					ospfNbrSeqNum:          seq_num,
					isSeqNumUpdate:         true,
					isMasterUpdate:         false,
					ospfNbrLsaReqIndex:     nbrConf.ospfNbrLsaReqIndex,
				},
				nbrMsgType: NBRUPD,
			}
			server.neighborConfCh <- nbrConfMsg
			server.updateNeighborMdata(nbrConf.intfConfKey, nbrKey)
			server.logger.Debug(fmt.Sprintln("NBREVENT: Flood the LSA. nbr full state ", nbrKey.IPAddr))
		//	server.CreateNetworkLSACh <- ospfIntfToNbrMap[nbrConf.intfConfKey]
		case config.NbrTwoWay:
			/* ignore packet */
			server.logger.Debug(fmt.Sprintln("NBRDBD: Ignore packet as NBR state is two way"))
			return
		case config.NbrDown:
			/* ignore packet. */
			server.logger.Debug(fmt.Sprintln("NBRDBD: Ignore packet . NBR is down"))
			return
		} // end of switch
	} else { //nbr doesnt exist
		server.logger.Debug(fmt.Sprintln("Ignore DB packet. Nbr doesnt exist ", nbrKey))
		return
	}
}

func (server *OSPFServer) ProcessNbrStateMachine() {
	for {

		select {
		case nbrData := <-(server.neighborHelloEventCh):
			//server.logger.Debug(fmt.Sprintln("NBREVENT: Received hellopkt event for nbrId ", nbrData.RouterId, " two_way", nbrData.TwoWayStatus))
			var nbrConf OspfNeighborEntry
			var send_dbd bool
			var seq_update bool
			var dbd_mdata ospfDatabaseDescriptionData
			nbrKey := NeighborConfKey{
				IPAddr:  config.IpAddress(nbrData.NeighborIP.String()),
				IntfIdx: nbrData.IntfConfKey.IntfIdx,
			}
			//server.logger.Debug(fmt.Sprintln("NBREVET: Nbr key ", nbrData.NeighborIP.String(), nbrData.IntfConfKey.IntfIdx))
			//Check if neighbor exists
			_, exists := server.NeighborConfigMap[nbrKey]
			intfConf, _ := server.IntfConfMap[nbrData.IntfConfKey]
			send_dbd = false
			seq_update = false
			isStateUpdate := false
			if exists {
				nbrConf = server.NeighborConfigMap[nbrKey]
				if nbrData.TwoWayStatus { // update the state
					startAdjacency := server.adjacancyEstablishementCheck(nbrConf.isDRBDR, true)
					if startAdjacency && nbrConf.OspfNbrState == config.NbrTwoWay {
						nbrConf.OspfNbrState = config.NbrExchangeStart
						isStateUpdate = true
						if nbrConf.ospfNbrSeqNum == 0 {
							nbrConf.ospfNbrSeqNum = uint32(time.Now().Unix())
							dbd_mdata.dd_sequence_number = nbrConf.ospfNbrSeqNum
							dbd_mdata.msbit = true // i am master
							dbd_mdata.ibit = true
							dbd_mdata.mbit = true
							nbrConf.isMaster = false
							seq_update = true
						} else {
							dbd_mdata.dd_sequence_number = nbrConf.ospfNbrSeqNum
							dbd_mdata.msbit = true
							nbrConf.isMaster = false
						}
						dbd_mdata.interface_mtu = INTF_MTU_MIN
						server.logger.Debug(fmt.Sprintln("NBRHELLO: Send, seq no ", dbd_mdata.dd_sequence_number,
							"msbit ", dbd_mdata.msbit))
						send_dbd = true
					} else { // no adjacency
						if nbrConf.OspfNbrState < config.NbrTwoWay {
							nbrConf.OspfNbrState = config.NbrTwoWay
							isStateUpdate = true
						}
					}
				} else {
					nbrConf.OspfNbrState = config.NbrInit
					isStateUpdate = true
				}

				nbrConfMsg := ospfNeighborConfMsg{
					ospfNbrConfKey: nbrKey,
					ospfNbrEntry: OspfNeighborEntry{
						OspfNbrRtrId:           nbrData.RouterId,
						OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
						OspfRtrPrio:            nbrConf.OspfRtrPrio,
						intfConfKey:            nbrConf.intfConfKey,
						OspfNbrOptions:         0,
						OspfNbrState:           nbrConf.OspfNbrState,
						isStateUpdate:          isStateUpdate,
						OspfNbrInactivityTimer: time.Now(),
						OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
						ospfNbrDBDTickerCh:     nbrConf.ospfNbrDBDTickerCh,
						ospfNbrSeqNum:          nbrConf.ospfNbrSeqNum,
						isSeqNumUpdate:         seq_update,
						isMasterUpdate:         false,
						nbrEvent:               nbrConf.nbrEvent,
					},
					nbrMsgType: NBRUPD,
				}
				server.neighborConfCh <- nbrConfMsg

				if send_dbd {
					server.ConstructAndSendDbdPacket(nbrConfMsg.ospfNbrConfKey, true, true, true,
						INTF_OPTIONS, nbrConf.ospfNbrSeqNum, false, false, intfConf.IfMtu)
				}

			} else { //neighbor doesnt exist
				var ticker *time.Ticker
				var nbrState config.NbrState
				var dbd_mdata ospfDatabaseDescriptionData
				var send_dbd bool
				msg := DbEventMsg{
					eventType: config.ADJACENCY,
				}
				msg.eventInfo = "Create new neighbor with id " + nbrData.NeighborIP.String()
				server.DbEventOp <- msg
				server.logger.Info(fmt.Sprintln("NBREVENT: Create new neighbor with id ", nbrData.RouterId))

				if nbrData.TwoWayStatus { // update the state
					startAdjacency := server.adjacancyEstablishementCheck(false, true)
					if startAdjacency {
						nbrState = config.NbrExchangeStart
						dbd_mdata.dd_sequence_number = uint32(time.Now().Nanosecond())
						seq_update = true
						// send dbd packets
						ticker = time.NewTicker(time.Second * 10)
						send_dbd = true
						server.logger.Info(fmt.Sprintln("NBRHELLO: Send, seq no ", dbd_mdata.dd_sequence_number,
							"msbit ", dbd_mdata.msbit))
					} else { // no adjacency
						nbrState = config.NbrTwoWay
						send_dbd = false
						isStateUpdate = true
					}
				} else {
					nbrState = config.NbrInit
					send_dbd = false
					isStateUpdate = true
				}

				nbrConfMsg := ospfNeighborConfMsg{
					ospfNbrConfKey: nbrKey,
					ospfNbrEntry: OspfNeighborEntry{
						OspfNbrRtrId:           nbrData.RouterId,
						OspfNbrIPAddr:          nbrData.NeighborIP,
						OspfRtrPrio:            nbrData.RtrPrio,
						intfConfKey:            nbrData.IntfConfKey,
						OspfNbrOptions:         0,
						OspfNbrState:           nbrState,
						isStateUpdate:          isStateUpdate,
						OspfNbrInactivityTimer: time.Now(),
						OspfNbrDeadTimer:       nbrData.nbrDeadTimer,
						ospfNbrSeqNum:          dbd_mdata.dd_sequence_number,
						isSeqNumUpdate:         seq_update,
						ospfNbrDBDTickerCh:     ticker,
						isMasterUpdate:         false,
					},
					nbrMsgType: NBRADD,
				}
				/* add the stub entry so that till the update thread updates the data
				valid entry will be present in the map */
				server.NeighborConfigMap[nbrKey] = nbrConfMsg.ospfNbrEntry
				server.initNeighborMdata(nbrData.IntfConfKey)
				server.neighborConfCh <- nbrConfMsg
				if send_dbd {
					dbd_mdata.ibit = true
					dbd_mdata.mbit = true
					dbd_mdata.msbit = true

					dbd_mdata.interface_mtu = uint16(intfConf.IfMtu)
					dbd_mdata.options = INTF_OPTIONS
				}
				server.logger.Info(fmt.Sprintln("NBREVENT: ADD Nbr ", nbrData.RouterId, "state ", nbrState))
			}

		case nbrDbPkt := <-(server.neighborDBDEventCh):
			server.logger.Info(fmt.Sprintln("NBREVENT: DBD received  ", nbrDbPkt))
			server.processDBDEvent(nbrDbPkt.ospfNbrConfKey, nbrDbPkt.ospfNbrDBDData)

		case intfKey := <-(server.neighborIntfEventCh):
			server.logger.Info(fmt.Sprintln("NBREVENT: Interface down ", intfKey))
			server.processIntfDownEvent(intfKey)

		case state := <-server.neighborFSMCtrlCh:
			if state == false {
				return
			}
		}
	} // end of for
}

/* @fn ProcessRxNbrPkt
Nbr packet Rx thread. It handles LSA REQ/UPD/ACK in.
*/
func (server *OSPFServer) ProcessRxNbrPkt() {
	for {
		select {
		case nbrLSAReqPkt := <-(server.neighborLSAReqEventCh):
			nbr, exists := server.NeighborConfigMap[nbrLSAReqPkt.nbrKey]
			if exists && nbr.OspfNbrState >= config.NbrExchange {
				server.DecodeLSAReq(nbrLSAReqPkt)
			}

		case nbrLSAUpdPkt := <-(server.neighborLSAUpdEventCh):
			nbr, exists := server.NeighborConfigMap[nbrLSAUpdPkt.nbrKey]

			if exists && nbr.OspfNbrState >= config.NbrExchange {
				server.DecodeLSAUpd(nbrLSAUpdPkt)
			}

		case nbrLSAAckPkt := <-(server.neighborLSAACKEventCh):
			nbr, exists := server.NeighborConfigMap[nbrLSAAckPkt.nbrKey]

			if exists && nbr.OspfNbrState >= config.NbrExchange {
				server.logger.Info(fmt.Sprintln("ACK : received - ", nbrLSAAckPkt))
				//server.DecodeLSAAck(nbrLSAAckPkt)
			}

		case stop := <-(server.ospfRxNbrPktStopCh):
			if stop {
				return
			}

		}

	}
}

/* @fn ProcessTxNbrPkt
Nbr packet out thread. It handles LSA REQ/UPD/ACK out , DBD packets out
signalling the LSA generation
*/
func (server *OSPFServer) ProcessTxNbrPkt() {
	for {
		select {
		case dbd_mdata := <-server.ospfNbrDBDSendCh:
			nbrConf, exists := server.NeighborConfigMap[dbd_mdata.ospfNbrConfKey]
			if exists {
				intConf, exist := server.IntfConfMap[nbrConf.intfConfKey]
				if exist {
					dstMac, _ := ospfNeighborIPToMAC[dbd_mdata.ospfNbrConfKey]
					data := server.BuildDBDPkt(nbrConf.intfConfKey, intConf, nbrConf,
						dbd_mdata.ospfNbrDBDData, dstMac)
					server.SendOspfPkt(nbrConf.intfConfKey, data)
				}
				/* This ensures Flood packet is sent only after last DBD.  */
				if nbrConf.isMaster && !dbd_mdata.ospfNbrDBDData.ibit && !dbd_mdata.ospfNbrDBDData.mbit {
					server.CreateNetworkLSACh <- ospfIntfToNbrMap[nbrConf.intfConfKey]
				}
			}

		case lsa_data := <-server.ospfNbrLsaReqSendCh:
			nbrConf, exists := server.NeighborConfigMap[lsa_data.nbrKey]
			if exists {
				intConf, exist := server.IntfConfMap[nbrConf.intfConfKey]
				if exist {
					dstMac, _ := ospfNeighborIPToMAC[lsa_data.nbrKey]
					data := server.EncodeLSAReqPkt(nbrConf.intfConfKey, intConf, nbrConf, lsa_data.lsa_slice, dstMac)
					server.SendOspfPkt(nbrConf.intfConfKey, data)
				}

			}

		case msg := <-server.ospfNbrLsaUpdSendCh:
			server.processFloodMsg(msg)

		case msg := <-server.ospfNbrLsaAckSendCh:
			server.processTxLsaAck(msg)

		case stop := <-server.ospfTxNbrPktStopCh:
			if stop == true {
				return
			}
		}
	}

}

func (server *OSPFServer) generateDbSummaryList(nbrConfKey NeighborConfKey) {
	nbrConf, exists := server.NeighborConfigMap[nbrConfKey]

	if !exists {
		server.logger.Err(fmt.Sprintln("negotiation: db_list Nbr  doesnt exist. nbr ", nbrConfKey))
		return
	}
	intf, _ := server.IntfConfMap[nbrConf.intfConfKey]

	areaId := convertIPv4ToUint32(intf.IfAreaId)
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	area_lsa, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("negotiation: db_list self originated lsas dont exist. Nbr , lsdb_key ", nbrConfKey, lsdbKey))
		return
	}
	router_lsdb := area_lsa.RouterLsaMap
	network_lsa := area_lsa.NetworkLsaMap
	ospfNeighborDBSummary_list[nbrConfKey] = nil
	db_list := []*ospfNeighborDBSummary{}
	for lsaKey, _ := range router_lsdb {
		// check if lsa instance is marked true
		db_summary := newospfNeighborDBSummary()
		drlsa, ret := server.getRouterLsaFromLsdb(areaId, lsaKey)
		if ret == LsdbEntryNotFound {
			continue
		}
		db_summary.lsa_headers = getLsaHeaderFromLsa(drlsa.LsaMd.LSAge, drlsa.LsaMd.Options,
			RouterLSA, lsaKey.LSId, lsaKey.AdvRouter,
			uint32(drlsa.LsaMd.LSSequenceNum), drlsa.LsaMd.LSChecksum,
			drlsa.LsaMd.LSLen)
		db_summary.valid = true
		/* add entry to the db summary list  */
		db_list = append(db_list, db_summary)
		//lsid := convertUint32ToIPv4(lsaKey.LSId)
		//server.logger.Info(fmt.Sprintln("negotiation: db_list append router lsid  ", lsid))
	} // end of for

	for networkKey, _ := range network_lsa {
		// check if lsa instance is marked true
		db_summary := newospfNeighborDBSummary()
		//if nbrMdata.isDR {
		dnlsa, ret := server.getNetworkLsaFromLsdb(areaId, networkKey)
		if ret == LsdbEntryNotFound {
			continue
		}
		db_summary.lsa_headers = getLsaHeaderFromLsa(dnlsa.LsaMd.LSAge, dnlsa.LsaMd.Options,
			NetworkLSA, networkKey.LSId, networkKey.AdvRouter,
			uint32(dnlsa.LsaMd.LSSequenceNum), dnlsa.LsaMd.LSChecksum,
			dnlsa.LsaMd.LSLen)
		db_summary.valid = true
		/* add entry to the db summary list  */
		db_list = append(db_list, db_summary)

	} // end of for

	/*   attach summary list */

	summary3_list := server.generateDbsummary3LsaList(areaId)
	if summary3_list != nil {
		db_list = append(db_list, summary3_list...)
	}

	summary4_list := server.generateDbsummary4LsaList(areaId)
	if summary4_list != nil {
		db_list = append(db_list, summary4_list...)
	}

	asExternal_list := server.generateDbasExternalList(areaId)
	if asExternal_list != nil {
		db_list = append(db_list, asExternal_list...)
	}

	for lsa := range db_list {
		rtr_id := convertUint32ToIPv4(db_list[lsa].lsa_headers.adv_router_id)
		server.logger.Info(fmt.Sprintln(lsa, ": ", rtr_id, " lsatype ", db_list[lsa].lsa_headers.ls_type))
	}
	nbrConf.db_summary_list_mutex.Lock()
	ospfNeighborDBSummary_list[nbrConfKey] = db_list
	nbrConf.db_summary_list_mutex.Unlock()
}

/*@fn generateDbasExternalList
This function generates As external list if the router is ASBR
*/
func (server *OSPFServer) generateDbasExternalList(self_areaId uint32) []*ospfNeighborDBSummary {
	if !server.ospfGlobalConf.AreaBdrRtrStatus {
		return nil // dont add self gen LSA if I am not ASBR
	}
	db_list := []*ospfNeighborDBSummary{}
	lsdbKey := LsdbKey{
		AreaId: self_areaId,
	}

	area_lsa, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("negotiation: As external LSA doesnt exist"))
		return nil
	}
	as_lsdb := area_lsa.ASExternalLsaMap

	for lsaKey, _ := range as_lsdb {
		// check if lsa instance is marked true
		db_as := newospfNeighborDBSummary()
		drlsa, ret := server.getASExternalLsaFromLsdb(self_areaId, lsaKey)
		if ret == LsdbEntryNotFound {
			continue
		}
		db_as.lsa_headers = getLsaHeaderFromLsa(drlsa.LsaMd.LSAge, drlsa.LsaMd.Options,
			ASExternalLSA, lsaKey.LSId, lsaKey.AdvRouter,
			uint32(drlsa.LsaMd.LSSequenceNum), drlsa.LsaMd.LSChecksum,
			drlsa.LsaMd.LSLen)
		db_as.valid = true
		/* add entry to the db summary list  */
		db_list = append(db_list, db_as)
		lsid := convertUint32ToIPv4(lsaKey.LSId)
		server.logger.Info(fmt.Sprintln("negotiation: db_list AS ext append router lsid  ", lsid))
	}
	return db_list
}

/* @fn generateDbsummaryLsaList
This function will attach summary LSAs if the router is ABR
*/
func (server *OSPFServer) generateDbsummary3LsaList(self_areaId uint32) []*ospfNeighborDBSummary {
	db_list := []*ospfNeighborDBSummary{}

	lsdbKey := LsdbKey{
		AreaId: self_areaId,
	}

	area_lsa, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("negotiation: Summary LSA 3 doesnt exist"))
		return nil
	}
	summary_lsdb := area_lsa.Summary3LsaMap
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]

	for lsaKey, _ := range summary_lsdb {
		_, exist := selfOrigLsaEnt[lsaKey]
		if exist && !server.ospfGlobalConf.isABR {
			continue // dont add self gen LSA if I am not ABR
		}
		// check if lsa instance is marked true
		db_summary := newospfNeighborDBSummary()
		drlsa, ret := server.getSummaryLsaFromLsdb(self_areaId, lsaKey)
		if ret == LsdbEntryNotFound {
			continue
		}
		db_summary.lsa_headers = getLsaHeaderFromLsa(drlsa.LsaMd.LSAge, drlsa.LsaMd.Options,
			Summary3LSA, lsaKey.LSId, lsaKey.AdvRouter,
			uint32(drlsa.LsaMd.LSSequenceNum), drlsa.LsaMd.LSChecksum,
			drlsa.LsaMd.LSLen)
		db_summary.valid = true
		/* add entry to the db summary list  */
		db_list = append(db_list, db_summary)
		lsid := convertUint32ToIPv4(lsaKey.LSId)
		server.logger.Info(fmt.Sprintln("negotiation: db_list summary 3 append router lsid  ", lsid))
	}
	return db_list
}

/* @fn generateDbsummaryLsaList
This function will attach summary LSAs if the router is ABR
*/
func (server *OSPFServer) generateDbsummary4LsaList(self_areaId uint32) []*ospfNeighborDBSummary {
	db_list := []*ospfNeighborDBSummary{}

	lsdbKey := LsdbKey{
		AreaId: self_areaId,
	}

	area_lsa, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("negotiation: Summary LSA 4 doesnt exist"))
		return nil
	}
	summary_lsdb := area_lsa.Summary4LsaMap
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]

	for lsaKey, _ := range summary_lsdb {
		_, exist := selfOrigLsaEnt[lsaKey]
		if exist && !server.ospfGlobalConf.isABR {
			continue // dont add self gen LSA if I am not ABR
		}
		// check if lsa instance is marked true
		db_summary := newospfNeighborDBSummary()
		drlsa, ret := server.getSummaryLsaFromLsdb(self_areaId, lsaKey)
		if ret == LsdbEntryNotFound {
			continue
		}
		db_summary.lsa_headers = getLsaHeaderFromLsa(drlsa.LsaMd.LSAge, drlsa.LsaMd.Options,
			Summary4LSA, lsaKey.LSId, lsaKey.AdvRouter,
			uint32(drlsa.LsaMd.LSSequenceNum), drlsa.LsaMd.LSChecksum,
			drlsa.LsaMd.LSLen)
		db_summary.valid = true
		/* add entry to the db summary list  */
		db_list = append(db_list, db_summary)
		lsid := convertUint32ToIPv4(lsaKey.LSId)
		server.logger.Info(fmt.Sprintln("negotiation: db_list summary 4 append router lsid  ", lsid))
	}
	return db_list
}

func (server *OSPFServer) generateRequestList(nbrKey NeighborConfKey, nbrConf OspfNeighborEntry, nbrDbPkt ospfDatabaseDescriptionData) {
	/* 1) get lsa headers update in req_list */
	headers_len := len(nbrDbPkt.lsa_headers)
	server.logger.Info(fmt.Sprintln("REQ_LIST: Received lsa headers for nbr ", nbrKey,
		" no of header ", headers_len))
	req_list := ospfNeighborRequest_list[nbrKey]
	for i := 0; i < headers_len; i++ {
		var lsaheader ospfLSAHeader
		lsaheader = nbrDbPkt.lsa_headers[i]
		result := server.lsaAddCheck(lsaheader, nbrConf) // check lsdb
		if result {
			req := newospfNeighborReq()
			req.lsa_headers = lsaheader
			req.valid = true
			nbrConf.req_list_mutex.Lock()
			req_list = append(req_list, req)
			nbrConf.req_list_mutex.Unlock()
		}
	}
	ospfNeighborRequest_list[nbrKey] = req_list
	server.logger.Info(fmt.Sprintln("REQ_LIST: updated req_list for nbr ",
		nbrKey, " req_list ", req_list))
}

func (server *OSPFServer) neighborDeadTimerEvent(nbrConfKey NeighborConfKey) {
	var nbr_entry_dead_func func()

	nbr_entry_dead_func = func() {
		server.logger.Info(fmt.Sprintln("NBRSCAN: DEAD ", nbrConfKey.IPAddr))

		_, exists := server.NeighborConfigMap[nbrConfKey]
		if exists {
			nbrConf := server.NeighborConfigMap[nbrConfKey]
			nbrConfMsg := ospfNeighborConfMsg{
				ospfNbrConfKey: nbrConfKey,
				ospfNbrEntry: OspfNeighborEntry{
					OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
					OspfRtrPrio:            nbrConf.OspfRtrPrio,
					intfConfKey:            nbrConf.intfConfKey,
					OspfNbrOptions:         0,
					OspfNbrState:           config.NbrDown,
					isStateUpdate:          true,
					OspfNbrInactivityTimer: time.Now(),
					OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
				},
				nbrMsgType: NBRDEL,
			}
			// update neighbor map
			server.processNeighborDeadEvent(nbrConfKey, nbrConf.intfConfKey)
			server.neighborConfCh <- nbrConfMsg
		}
	} // end of afterFunc callback

	_, exists := server.NeighborConfigMap[nbrConfKey]
	if exists {
		nbrConf := server.NeighborConfigMap[nbrConfKey]
		nbrConf.NbrDeadTimer = time.AfterFunc(nbrConf.OspfNbrDeadTimer, nbr_entry_dead_func)
		server.NeighborConfigMap[nbrConfKey] = nbrConf
	}

}

/*@fn refreshNeighborSlice
Refresh get bulk slice for all keys.
*/
func (server *OSPFServer) refreshNeighborSlice() {
	for {
		select {
		case start := <-server.neighborSliceStartCh:
			if start {
				server.neighborSliceRefCh = time.NewTicker(server.RefreshDuration)
			} else {
				if server.neighborSliceRefCh != nil {
					server.neighborSliceRefCh.Stop()
				}
			}
		case <-server.neighborSliceRefCh.C:
			server.neighborBulkSlice = []NeighborConfKey{}
			idx := 0
			for nbrKey, _ := range server.NeighborConfigMap {
				server.neighborBulkSlice = append(server.neighborBulkSlice, nbrKey)
				idx++
			}

		}
	}

}

/* @fn processNeighborDeadEvent
	1) clear retransmission list.
	2) send message to IFFSM
	3) Updade If to neighbor map.
	   Send message to LSDB  - To update router LSA
	   and if i am DR update network LSA . Flood these LSAs.
	4) delete neighbor from neighbor global map.
IfFSM takes care of electing DR BDR and sending message to LSDB to
update LSAs . therefore following APIs takes action 1.
Note - From RFC -
        If an adjacent router goes down, retransmissions may occur until
        the adjacency is destroyed by OSPF's Hello Protocol.  When the
        adjacency is destroyed, the Link state retransmission list is
        cleared.

*/
func (server *OSPFServer) processNeighborDeadEvent(nbrKey NeighborConfKey, intfKey IntfConfKey) {
	/* Age LSAs */
	server.logger.Info(fmt.Sprintln("DEAD: start processing nbr dead ", nbrKey))
	server.resetNeighborLists(nbrKey, intfKey)
	nbrStateChangeData := NbrStateChangeMsg{
		nbrKey: nbrKey,
	}

	intfConf := server.IntfConfMap[intfKey]
	intfConf.NbrStateChangeCh <- nbrStateChangeData
	server.logger.Info(fmt.Sprintln("DEAD: end processing nbr dead ", nbrKey))
}

/*@fn processIntfDownEvent
 This API takes care of interface down/delete event
1) Clear all lists.
2) send delete message for neighbor struct.
*/
func (server *OSPFServer) processIntfDownEvent(intfKey IntfConfKey) {
	server.logger.Info(fmt.Sprintln("DEAD: Intf down ", intfKey))
	_, exist := server.IntfConfMap[intfKey]
	if !exist {
		server.logger.Info(fmt.Sprintln("DEAD: Intf map doesnt exist ", intfKey))
		return
	}
	nbrMdata, valid := ospfIntfToNbrMap[intfKey]
	if !valid {
		server.logger.Info(fmt.Sprintln("DEAD: Intf deleted but intf-to-nbr map doesnt exist. ", intfKey))
		return
	}
	for _, nbr := range nbrMdata.nbrList { // delete all lists
		updateLSALists(nbr)
		nbrConfMsg := ospfNeighborConfMsg{
			ospfNbrConfKey: nbr,
			nbrMsgType:     NBRDEL,
		}
		server.neighborConfCh <- nbrConfMsg
	}
	//delete interface to nbr mapping.
	delete(ospfIntfToNbrMap, intfKey)

}
