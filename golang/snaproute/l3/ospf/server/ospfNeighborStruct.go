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
	"sync"
	"time"
)

type NbrMsgType uint32

const (
	NBRADD = 0
	NBRUPD = 1
	NBRDEL = 2
)

const (
	RxDBDInterval = 5
)

type LsaOp uint8

const (
	LSAFLOOD        = 0 // flood when FULL state reached
	LSASELFLOOD     = 1 // flood for received LSA
	LSAINTF         = 2 // Send LSA on the interface in reply to LSAREQ
	LSAAGE          = 3 // flood aged LSAs.
	LSASUMMARYFLOOD = 4 //flood summary LSAs in different areas.
	LSAEXTFLOOD     = 5 //flood AS External summary LSA
	LSAROUTERFLOOD  = 6 //flood only router LSA
)

type NeighborConfKey struct {
	IPAddr  config.IpAddress
	IntfIdx config.InterfaceIndexOrZero
	//	OspfNbrRtrId uint32
}

var INVALID_NEIGHBOR_CONF_KEY uint32
var neighborBulkSlice []NeighborConfKey

type OspfNeighborEntry struct {
	OspfNbrRtrId           uint32
	OspfNbrIPAddr          net.IP
	OspfRtrPrio            uint8
	intfConfKey            IntfConfKey
	OspfNbrOptions         int
	OspfNbrState           config.NbrState
	isStateUpdate          bool
	OspfNbrInactivityTimer time.Time
	OspfNbrDeadTimer       time.Duration
	NbrDeadTimer           *time.Timer
	isDRBDR                bool
	ospfNbrSeqNum          uint32
	nbrEvent               config.NbrEvent
	isSeqNumUpdate         bool
	isMaster               bool
	isMasterUpdate         bool
	ospfNbrDBDTickerCh     *time.Ticker

	ospfNbrLsaIndex        uint8       // db_summary list index
	ospfNbrLsaReqIndex     uint8       //req_list index
	ospfNeighborLsaRxTimer *time.Timer // retx interval timer
	req_list_mutex         *sync.Mutex
	db_summary_list_mutex  *sync.Mutex
	retx_list_mutex        *sync.Mutex
}

/* LSA lists */
type ospfNeighborReq struct {
	lsa_headers ospfLSAHeader
	valid       bool // entry is valid or not
}

func newospfNeighborReq() *ospfNeighborReq {
	return &ospfNeighborReq{}
}

type ospfNeighborDBSummary struct {
	lsa_headers ospfLSAHeader
	valid       bool
}

func newospfNeighborDBSummary() *ospfNeighborDBSummary {
	return &ospfNeighborDBSummary{}
}

type ospfNeighborRetx struct {
	lsa_headers ospfLSAHeader
	valid       bool
}

func newospfNeighborRetx() *ospfNeighborRetx {
	return &ospfNeighborRetx{}
}

type ospfNeighborConfMsg struct {
	ospfNbrConfKey NeighborConfKey
	ospfNbrEntry   OspfNeighborEntry
	nbrMsgType     NbrMsgType
}

type ospfNeighborDBDMsg struct {
	ospfNbrConfKey NeighborConfKey
	nbrFull        bool
	ospfNbrDBDData ospfDatabaseDescriptionData
}

type ospfNbrMdata struct {
	isDR    bool
	areaId  uint32
	intf    IntfConfKey
	nbrList []NeighborConfKey
}

func newospfNbrMdata() *ospfNbrMdata {
	return &ospfNbrMdata{}
}

/*
	Global structures for Neighbor
*/
var OspfNeighborLastDbd map[NeighborConfKey]ospfDatabaseDescriptionData
var ospfNeighborIPToMAC map[NeighborConfKey]net.HardwareAddr

/* neighbor lists each indexed by neighbor router id. */
var ospfNeighborRequest_list map[NeighborConfKey][]*ospfNeighborReq
var ospfNeighborDBSummary_list map[NeighborConfKey][]*ospfNeighborDBSummary
var ospfNeighborRetx_list map[NeighborConfKey][]*ospfNeighborRetx

/* List of Neighbors per interface instance */
var ospfIntfToNbrMap map[IntfConfKey]ospfNbrMdata

func (server *OSPFServer) InitNeighborStateMachine() {

	server.neighborBulkSlice = []NeighborConfKey{}
	INVALID_NEIGHBOR_CONF_KEY = 0
	OspfNeighborLastDbd = make(map[NeighborConfKey]ospfDatabaseDescriptionData)
	ospfNeighborIPToMAC = make(map[NeighborConfKey]net.HardwareAddr)
	ospfIntfToNbrMap = make(map[IntfConfKey]ospfNbrMdata)
	ospfNeighborRequest_list = make(map[NeighborConfKey][]*ospfNeighborReq)
	ospfNeighborDBSummary_list = make(map[NeighborConfKey][]*ospfNeighborDBSummary)
	ospfNeighborRetx_list = make(map[NeighborConfKey][]*ospfNeighborRetx)

	server.neighborSliceRefCh = time.NewTicker(server.RefreshDuration)
	go server.refreshNeighborSlice()
	server.neighborSliceStartCh <- true
	server.logger.Info("NBRINIT: Neighbor FSM init done..")
}

func calculateMaxLsaHeaders() (max_headers uint8) {
	rem := INTF_MTU_MIN - (OSPF_DBD_MIN_SIZE + OSPF_HEADER_SIZE)
	max_headers = uint8(rem / OSPF_LSA_HEADER_SIZE)
	return max_headers
}

func calculateMaxLsaReq() (max_req uint8) {
	rem := INTF_MTU_MIN - OSPF_HEADER_SIZE
	max_req = uint8(rem / OSPF_LSA_REQ_SIZE)
	return max_req
}

/*@fn UpdateNeighborConf
Thread to update/add/delete neighbor global struct.
*/
func (server *OSPFServer) UpdateNeighborConf() {
	for {
		select {
		case nbrMsg := <-(server.neighborConfCh):
			var nbrConf OspfNeighborEntry
			intfConf, _ := server.IntfConfMap[nbrMsg.ospfNbrEntry.intfConfKey]
			//server.logger.Info(fmt.Sprintln("Update neighbor conf.  received"))
			if nbrMsg.nbrMsgType == NBRDEL {
				delete(server.NeighborConfigMap, nbrMsg.ospfNbrConfKey)
				server.logger.Info(fmt.Sprintln("DELETE neighbor with nbr id - ",
					nbrMsg.ospfNbrConfKey.IPAddr, nbrMsg.ospfNbrConfKey.IntfIdx))
				continue
			}
			if nbrMsg.nbrMsgType == NBRUPD {
				nbrConf = server.NeighborConfigMap[nbrMsg.ospfNbrConfKey]
			}
			if nbrMsg.ospfNbrEntry.isStateUpdate {
				nbrConf.OspfNbrState = nbrMsg.ospfNbrEntry.OspfNbrState
			}
			nbrConf.OspfNbrDeadTimer = nbrMsg.ospfNbrEntry.OspfNbrDeadTimer
			nbrConf.OspfNbrInactivityTimer = time.Now()
			if nbrMsg.ospfNbrEntry.isSeqNumUpdate {
				nbrConf.ospfNbrSeqNum = nbrMsg.ospfNbrEntry.ospfNbrSeqNum
			}
			if nbrMsg.ospfNbrEntry.isMasterUpdate {
				nbrConf.isMaster = nbrMsg.ospfNbrEntry.isMaster
			}
			nbrConf.OspfNbrRtrId = nbrMsg.ospfNbrEntry.OspfNbrRtrId
			nbrConf.ospfNbrDBDTickerCh = nbrMsg.ospfNbrEntry.ospfNbrDBDTickerCh
			nbrConf.ospfNbrLsaReqIndex = nbrMsg.ospfNbrEntry.ospfNbrLsaReqIndex
			nbrConf.nbrEvent = nbrMsg.ospfNbrEntry.nbrEvent

			if nbrMsg.nbrMsgType == NBRADD {
				nbrConf.OspfNbrIPAddr = nbrMsg.ospfNbrEntry.OspfNbrIPAddr
				nbrConf.OspfRtrPrio = nbrMsg.ospfNbrEntry.OspfRtrPrio
				nbrConf.intfConfKey = nbrMsg.ospfNbrEntry.intfConfKey
				nbrConf.OspfNbrOptions = 0
				if nbrMsg.ospfNbrEntry.isMasterUpdate {
					nbrConf.isMaster = nbrMsg.ospfNbrEntry.isMaster
				}
				server.neighborBulkSlice = append(server.neighborBulkSlice, nbrMsg.ospfNbrConfKey)
				nbrConf.req_list_mutex = &sync.Mutex{}
				nbrConf.db_summary_list_mutex = &sync.Mutex{}
				nbrConf.retx_list_mutex = &sync.Mutex{}
				updateLSALists(nbrMsg.ospfNbrConfKey)
				server.NeighborConfigMap[nbrMsg.ospfNbrConfKey] = nbrConf
				if nbrMsg.ospfNbrEntry.OspfNbrState >= config.NbrTwoWay {
					seq_num := uint32(time.Now().Nanosecond())
					server.ConstructAndSendDbdPacket(nbrMsg.ospfNbrConfKey, true, true, true,
						INTF_OPTIONS, seq_num, false, false, intfConf.IfMtu)
					nbrConf.OspfNbrState = config.NbrExchangeStart
					nbrConf.nbrEvent = config.Nbr2WayReceived
					nbrConf.ospfNbrSeqNum = seq_num
					server.NeighborConfigMap[nbrMsg.ospfNbrConfKey] = nbrConf
				}
				server.neighborDeadTimerEvent(nbrMsg.ospfNbrConfKey)
				server.logger.Info(fmt.Sprintln("CREATE: Create new neighbor with key ", nbrMsg.ospfNbrConfKey.IPAddr, nbrMsg.ospfNbrConfKey.IntfIdx))
			}

			if nbrMsg.nbrMsgType == NBRUPD {
				server.NeighborConfigMap[nbrMsg.ospfNbrConfKey] = nbrConf
				nbrConf.NbrDeadTimer.Stop()
				nbrConf.NbrDeadTimer.Reset(nbrMsg.ospfNbrEntry.OspfNbrDeadTimer)
			}

			//rtr_id := convertUint32ToIPv4(nbrMsg.ospfNbrEntry.OspfNbrRtrId)
		//	server.logger.Info(fmt.Sprintln("NBR UPDATE: Nbr , state ", rtr_id, " : ", nbrConf.OspfNbrState))

		case state := <-(server.neighborConfStopCh):
			server.logger.Info("Exiting update neighbor config thread..")
			if state == true {
				return
			}
		}
	}
}

func updateLSALists(id NeighborConfKey) {
	ospfNeighborRequest_list[id] = []*ospfNeighborReq{}
	ospfNeighborDBSummary_list[id] = []*ospfNeighborDBSummary{}
	ospfNeighborRetx_list[id] = []*ospfNeighborRetx{}
}

func (server *OSPFServer) neighborExist(nbrKey NeighborConfKey) bool {
	_, exists := server.NeighborConfigMap[nbrKey]
	if exists {
		return true
	}
	return false
}

func (server *OSPFServer) initNeighborMdata(intf IntfConfKey) {
	nbrMdata := newospfNbrMdata()
	intfConf, exist := server.IntfConfMap[intf]
	if !exist {
		server.logger.Err(fmt.Sprintln("Init NbrData:Intf doesnt exsit. Can not initialise nbr mdata. key ", intf))
		return
	}
	nbrMdata.nbrList = []NeighborConfKey{}
	nbrMdata.intf = intf
	nbrMdata.areaId = convertIPv4ToUint32(intfConf.IfAreaId)
	ospfIntfToNbrMap[intf] = *nbrMdata
}

func (server *OSPFServer) updateNeighborMdata(intf IntfConfKey, nbr NeighborConfKey) {
	nbrMdata, exists := ospfIntfToNbrMap[intf]
	intfData := server.IntfConfMap[intf]
	if !exists {
		server.initNeighborMdata(intf)
		nbrMdata = ospfIntfToNbrMap[intf]
	}
	nbrMdata.areaId = binary.BigEndian.Uint32(intfData.IfAreaId)
	routerid := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
	if intfData.IfDRtrId == routerid {
		nbrMdata.isDR = true
	} else {
		nbrMdata.isDR = false
	}

	for inst := range nbrMdata.nbrList {
		if nbrMdata.nbrList[inst] == nbr {
			// nbr already exist no need to add.
			return
		}
	}
	nbrMdata.nbrList = append(nbrMdata.nbrList, nbr)
	ospfIntfToNbrMap[intf] = nbrMdata
}

func (server *OSPFServer) sendLsdbToNeighborEvent(intfKey IntfConfKey, nbrKey NeighborConfKey,
	areaId uint32, lsType uint8, linkId uint32, lsaKey LsaKey, op uint8) {
	msg := ospfFloodMsg{
		intfKey: intfKey,
		nbrKey:  nbrKey,
		areaId:  areaId,
		lsType:  lsType,
		linkid:  linkId,
		lsaKey:  lsaKey,
		lsOp:    op,
	}
	server.ospfNbrLsaUpdSendCh <- msg
}

func (server *OSPFServer) resetNeighborLists(nbr NeighborConfKey, intf IntfConfKey) {
	/* List of Neighbors per interface instance */
	updateLSALists(nbr)
	nbrMdata, exists := ospfIntfToNbrMap[intf]
	if !exists {
		server.logger.Info(fmt.Sprintln("DEAD: Nbr dead but intf-to-nbr map doesnt exist. ", nbr))
		return
	}
	newList := []NeighborConfKey{}
	for inst := range nbrMdata.nbrList {
		if nbrMdata.nbrList[inst] != nbr {
			newList = append(newList, nbrMdata.nbrList[inst])
		}
	}
	nbrMdata.nbrList = newList
	ospfIntfToNbrMap[intf] = nbrMdata
	server.logger.Info(fmt.Sprintln("DEAD: nbrList ", nbrMdata.nbrList))
}

func (server *OSPFServer) CheckNeighborFullEvent(nbrKey NeighborConfKey) {
	nbrConf, exists := server.NeighborConfigMap[nbrKey]
	nbrFull := true
	if exists {
		reqlist := ospfNeighborRequest_list[nbrKey]
		if reqlist != nil {
			for _, ent := range reqlist {
				if ent.valid == true {
					nbrFull = false
				}
			}
		}
		if !nbrFull {
			return
		}
		nbrConfMsg := ospfNeighborConfMsg{
			ospfNbrConfKey: nbrKey,
			ospfNbrEntry: OspfNeighborEntry{
				OspfNbrRtrId:           nbrConf.OspfNbrRtrId,
				OspfNbrIPAddr:          nbrConf.OspfNbrIPAddr,
				OspfRtrPrio:            nbrConf.OspfRtrPrio,
				intfConfKey:            nbrConf.intfConfKey,
				OspfNbrOptions:         0,
				OspfNbrState:           config.NbrFull,
				isStateUpdate:          true,
				OspfNbrInactivityTimer: time.Now(),
				OspfNbrDeadTimer:       nbrConf.OspfNbrDeadTimer,
				isSeqNumUpdate:         false,
				isMasterUpdate:         false,
				nbrEvent:               nbrConf.nbrEvent,
			},
			nbrMsgType: NBRUPD,
		}
		server.neighborConfCh <- nbrConfMsg
		server.logger.Info(fmt.Sprintln("NBREVENT: Nbr FULL ", nbrKey.IPAddr))
	}
}

func (server *OSPFServer) UpdateNeighborList(nbrKey NeighborConfKey) {
	nbrConf, exists := server.NeighborConfigMap[nbrKey]
	if exists {
		if nbrConf.OspfNbrState == config.NbrFull {
			return
		}
		server.CheckNeighborFullEvent(nbrKey)
	}
}
