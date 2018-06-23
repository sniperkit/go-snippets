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

func initNbrTestParams() {
	fmt.Println("\n Get Server object")
	ospf = getServerObject()
	initAttr()
	go startDummyChannels(ospf)
	ospf.InitNeighborStateMachine()
}

func TestOspfNbrFSM(t *testing.T) {
	fmt.Println("\n**************** NEIGHBOR FSM ************\n")
	initNbrTestParams()
	for index := 1; index < 21; index++ {
		err := nbrFSMTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test case for interface FSM")
		}
	}
}

func nbrFSMTestLogic(tNum int) int {
	ospf.initDefaultIntfConf(key, ipIntfProp, ifType)
	ospf.updateGlobalConf(gConf)
	switch tNum {
	case 1:
		fmt.Println(tNum, ": Running Neighbor create")
		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		ospf.neighborConfCh <- nbrConfMsg
		ospf.neighborConfStopCh <- true

		fmt.Println(tNum, ": Running updateLSALists")
		updateLSALists(nbrKey)

		fmt.Println(tNum, ": Running initNeighborMdata")
		ospf.initNeighborMdata(key)

		fmt.Println(tNum, ": Running updateNeighborMdata")
		ospf.updateNeighborMdata(key, nbrKey)

		fmt.Println(tNum, ": Running resetNeighborLists")
		ospf.IntfConfMap[key] = intf
		ospf.resetNeighborLists(nbrKey, key)

		fmt.Println(tNum, ": Running UpdateNeighborList")
		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		ospf.UpdateNeighborList(nbrKey)
		ospf.neighborConfStopCh <- true

		fmt.Println(tNum, ": Running exchangePacketDiscardCheck")
		discard := ospf.exchangePacketDiscardCheck(ospfNbrEntry, nbrDbPkt)
		if discard {
			fmt.Println("NbrTest : Packet discarded. ")
		}

		fmt.Println(tNum, ": Running verifyDuplicatePacket")
		isdDup := ospf.verifyDuplicatePacket(ospfNbrEntry, nbrDbPkt)
		if isdDup {
			fmt.Println("NbrTest: Packet is duplicate.")
		}

		fmt.Println(tNum, ": Running adjacancyEstablishementCheck")
		ok := ospf.adjacancyEstablishementCheck(false, false)
		if !ok {
			fmt.Println("NbrTest: Dont establish adjacency as its neither DR or BDR.")
		}

		ospf.IntfConfMap[key] = intf
		go ospf.UpdateNeighborConf()
		nbrConfMsg.ospfNbrEntry.OspfNbrState = config.NbrExchangeStart
		ospf.neighborConfCh <- nbrConfMsg
		nbrDbPkt.ibit = true
		nbrDbPkt.msbit = true
		fmt.Println(tNum, ": Running processNeighborExstart")

		ospf.processNeighborExstart(nbrKey, ospfNbrEntry, nbrDbPkt)
		ospf.neighborConfStopCh <- true

		fmt.Println(tNum, ": Running processDBDEvent exstart")
		ospfNbrEntry.OspfNbrState = config.NbrExchangeStart

		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = true
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002

		ospfNbrEntry.db_summary_list_mutex = &sync.Mutex{}
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

		fmt.Println(tNum, ": Running processDBDEvent exchange")
		ospfNbrEntry.OspfNbrState = config.NbrExchange

		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

		fmt.Println(tNum, ": Running processDBDEvent NbrLoading")
		ospfNbrEntry.OspfNbrState = config.NbrLoading
		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

		fmt.Println(tNum, ": Running processDBDEvent NbrFull")
		ospfNbrEntry.OspfNbrState = config.NbrFull
		ospfNbrEntry.ospfNbrSeqNum = 2002
		nbrDbPkt.ibit = false
		nbrDbPkt.msbit = true
		nbrDbPkt.dd_sequence_number = 2002
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.processDBDEvent(nbrKey, nbrDbPkt)

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

		fmt.Println(tNum, ": Running ProcessRxNbrPkt ")
		go ospf.ProcessRxNbrPkt()
		ospfNbrEntry.OspfNbrState = config.NbrFull
		ospf.NeighborConfigMap[nbrKey] = ospfNbrEntry
		ospf.neighborLSAReqEventCh <- nbrLsaReqMsg
		ospf.neighborBulkSlice = append(ospf.neighborBulkSlice, nbrKey)
		nextid, cnt, res := ospf.GetBulkOspfNbrEntryState(0, 10)
		fmt.Println("Getbulk nextid, cnt, res", nextid, " ", cnt, " ", res)
	}
	return SUCCESS
}
