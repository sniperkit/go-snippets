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
	"l3/ospf/config"
	"sync"
	"testing"
)

func _initNbrTestParams() {
	fmt.Println("\n Get Server object")
	ospf = getServerObject()
	initAttr()
	go startDummyChannels(ospf)
	ospf.InitNeighborStateMachine()
}

func _TestOspfNbrFSM(t *testing.T) {
	fmt.Println("\n**************** NEIGHBOR FSM ************\n")
	initNbrTestParams()
	for index := 1; index < 21; index++ {
		err := nbrFSMTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test case for interface FSM")
		}
	}
}

func _nbrFSMTestLogic(tNum int) int {
	ospf.initDefaultIntfConf(key, ipIntfProp, ifType)
	ospf.updateGlobalConf(gConf)
	switch tNum {
	case 1:
		fmt.Println(tNum, ": Running Neighbor create")
		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		ospf.neighborConfCh <- nbrConfMsg
		ospf.neighborConfStopCh <- true

	case 2:
		fmt.Println(tNum, ": Running updateLSALists")
		updateLSALists(nbrKey)

	case 3:
		fmt.Println(tNum, ": Running initNeighborMdata")
		ospf.initNeighborMdata(key)

	case 4:
		fmt.Println(tNum, ": Running updateNeighborMdata")
		ospf.updateNeighborMdata(key, nbrKey)

	case 5:
		fmt.Println(tNum, ": Running resetNeighborLists")
		ospf.IntfConfMap[key] = intf
		ospf.resetNeighborLists(nbrKey, key)

	case 6:
		fmt.Println(tNum, ": Running UpdateNeighborList")
		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		ospf.UpdateNeighborList(nbrKey)
		ospf.neighborConfStopCh <- true

	case 7:
		fmt.Println(tNum, ": Running exchangePacketDiscardCheck")
		discard := ospf.exchangePacketDiscardCheck(ospfNbrEntry, nbrDbPkt)
		if discard {
			fmt.Println("NbrTest : Packet discarded. ")
		}

	case 8:
		fmt.Println(tNum, ": Running verifyDuplicatePacket")
		isdDup := ospf.verifyDuplicatePacket(ospfNbrEntry, nbrDbPkt)
		if isdDup {
			fmt.Println("NbrTest: Packet is duplicate.")
		}

	case 9:
		fmt.Println(tNum, ": Running adjacancyEstablishementCheck")
		ok := ospf.adjacancyEstablishementCheck(false, false)
		if !ok {
			fmt.Println("NbrTest: Dont establish adjacency as its neither DR or BDR.")
		}

	case 10:
		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		nbrConfMsg.ospfNbrEntry.OspfNbrState = config.NbrExchangeStart
		ospf.neighborConfCh <- nbrConfMsg
		nbrDbPkt.ibit = true
		nbrDbPkt.msbit = true
		fmt.Println(tNum, ": Running processNeighborExstart")

		ospf.processNeighborExstart(nbrKey, ospfNbrEntry, nbrDbPkt)
		ospf.neighborConfStopCh <- true

	case 11:
		fmt.Println(tNum, ": Running processDBDEvent exstart")
		ospfNbrEntry.OspfNbrState = config.NbrExchangeStart

		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = true
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002

		ospfNbrEntry.db_summary_list_mutex = &sync.Mutex{}
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

	case 12:
		fmt.Println(tNum, ": Running processDBDEvent exchange")
		ospfNbrEntry.OspfNbrState = config.NbrExchange

		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

	case 13:
		fmt.Println(tNum, ": Running processDBDEvent NbrLoading")
		ospfNbrEntry.OspfNbrState = config.NbrLoading
		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

	case 14:
		fmt.Println(tNum, ": Running processDBDEvent NbrFull")
		ospfNbrEntry.OspfNbrState = config.NbrFull
		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

	case 15:
		fmt.Println(tNum, ": Running ProcessNbrStateMachine")
		go ospf.ProcessNbrStateMachine()
		ospf.neighborHelloEventCh <- nbrIntfMsg
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry

		fmt.Println("   Add existing Nbr ")
		ospf.neighborHelloEventCh <- nbrIntfMsg

		fmt.Println(" Check dbd processing ")
		ospf.neighborDBDEventCh <- nbrDbdMsg

		fmt.Println(" INtf down event ")
		ospf.neighborIntfEventCh <- key

		fmt.Println("Stop nbr processing routine ")
		ospf.neighborFSMCtrlCh <- false

	case 16:
		fmt.Println(tNum, ": Running ProcessRxNbrPkt ")
		go ospf.ProcessRxNbrPkt()
		ospfNbrEntry.OspfNbrState = config.NbrFull
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.neighborLSAReqEventCh <- nbrLsaReqMsg
	}
	return SUCCESS
}
