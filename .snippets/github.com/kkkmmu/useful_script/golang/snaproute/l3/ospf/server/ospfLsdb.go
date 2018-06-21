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
	"utils/commonDefs"
)

type LsdbUpdateMsg struct {
	MsgType uint8
	AreaId  uint32
	Data    []byte
}

func NewLsdbUpdateMsg() *LsdbUpdateMsg {
	return &LsdbUpdateMsg{}
}

type LSAChangeMsg struct {
	areaId uint32
}

type NetworkLSAChangeMsg struct {
	areaId    uint32
	intfKey   IntfConfKey
	intfState config.Status
}

type DrChangeMsg struct {
	areaId   uint32
	intfKey  IntfConfKey
	oldstate config.IfState
	newstate config.IfState
}

const (
	LsdbAdd      uint8 = 0
	LsdbDel      uint8 = 1
	LsdbUpdate   uint8 = 2
	LsdbNoAction uint8 = 3
)

const (
	P2PLink     uint8 = 1
	TransitLink uint8 = 2
	StubLink    uint8 = 3
	VirtualLink uint8 = 4
)

var lsdbTickerCh *time.Timer
var lsdbRefreshTickerCh *time.Timer

func (server *OSPFServer) initLSDatabase(areaId uint32) {
	server.logger.Debug(fmt.Sprintln("LSDB: Initialise LSDB for area id ", areaId))
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		lsDbEnt.RouterLsaMap = make(map[LsaKey]RouterLsa)
		lsDbEnt.NetworkLsaMap = make(map[LsaKey]NetworkLsa)
		lsDbEnt.Summary3LsaMap = make(map[LsaKey]SummaryLsa)
		lsDbEnt.Summary4LsaMap = make(map[LsaKey]SummaryLsa)
		lsDbEnt.ASExternalLsaMap = make(map[LsaKey]ASExternalLsa)
		server.AreaLsdb[lsdbKey] = lsDbEnt
	}
	selfOrigLsaEnt, exist := server.AreaSelfOrigLsa[lsdbKey]
	if !exist {
		selfOrigLsaEnt = make(map[LsaKey]bool)
		server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	}

}

func (server *OSPFServer) StartLSDatabase() {
	server.logger.Debug("Initializing LSA Database")
	for key, _ := range server.AreaConfMap {
		areaId := convertAreaOrRouterIdUint32(string(key.AreaId))
		server.initLSDatabase(areaId)
	}

	server.lsdbStateRefresh()
	maxAgeLsaMap = make(map[LsaKey][]byte)
	// start LSDB aging ticker
	lsdbTickerCh = time.NewTimer(time.Second * 1)
	lsdbRefreshTickerCh = time.NewTimer(time.Second * time.Duration(config.LSRefreshTime))
	go server.processLSDatabaseUpdates()
	return
}

func (server *OSPFServer) StopLSDatabase() {
	lsdbTickerCh.Stop()
	lsdbRefreshTickerCh.Stop()
}

func (server *OSPFServer) compareSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) bool {
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsDbEnt.Summary3LsaMap[lsaKey]
	} else if lsaKey.LSType == Summary4LSA {
		sLsa, _ = lsDbEnt.Summary4LsaMap[lsaKey]
	}
	if sLsa.Netmask != lsaEnt.Netmask {
		return false
	}
	if sLsa.Metric != lsaEnt.Metric {
		return false
	}
	// TODO : More garnular comparision
	return true
}

func (server *OSPFServer) updateSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) {
	server.logger.Debug(fmt.Sprintln("Need to update Summary Lsa in LSDB:", lsdbKey, lsaKey, lsaEnt))
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsDbEnt.Summary3LsaMap[lsaKey]
	} else if lsaKey.LSType == Summary4LSA {
		sLsa, _ = lsDbEnt.Summary4LsaMap[lsaKey]
	}

	sLsa.Netmask = lsaEnt.Netmask
	sLsa.Metric = lsaEnt.Metric

	sLsa.LsaMd.LSAge = 0
	sLsa.LsaMd.LSSequenceNum = sLsa.LsaMd.LSSequenceNum + 1
	sLsa.LsaMd.LSLen = lsaEnt.LsaMd.LSLen
	sLsa.LsaMd.Options = lsaEnt.LsaMd.Options
	sLsa.LsaMd.LSChecksum = 0
	LsaEnc := encodeSummaryLsa(sLsa, lsaKey)
	checksumOffset := uint16(14)
	sLsa.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	if lsaKey.LSType == Summary3LSA {
		lsDbEnt.Summary3LsaMap[lsaKey] = sLsa
	} else if lsaKey.LSType == Summary4LSA {
		lsDbEnt.Summary4LsaMap[lsaKey] = sLsa
	}
	server.AreaLsdb[lsdbKey] = lsDbEnt
}

func (server *OSPFServer) insertSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) {
	server.logger.Debug(fmt.Sprintln("Need to Insert Summary Lsa in LSDB:", lsdbKey, lsaKey, lsaEnt))
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsDbEnt.Summary3LsaMap[lsaKey]
	} else if lsaKey.LSType == Summary4LSA {
		sLsa, _ = lsDbEnt.Summary4LsaMap[lsaKey]
	}

	sLsa.Netmask = lsaEnt.Netmask
	sLsa.Metric = lsaEnt.Metric

	sLsa.LsaMd.LSAge = 0
	sLsa.LsaMd.LSSequenceNum = InitialSequenceNumber
	sLsa.LsaMd.LSLen = lsaEnt.LsaMd.LSLen
	sLsa.LsaMd.Options = lsaEnt.LsaMd.Options
	sLsa.LsaMd.LSChecksum = 0
	LsaEnc := encodeSummaryLsa(sLsa, lsaKey)
	checksumOffset := uint16(14)
	sLsa.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	if lsaKey.LSType == Summary3LSA {
		lsDbEnt.Summary3LsaMap[lsaKey] = sLsa
	} else if lsaKey.LSType == Summary4LSA {
		lsDbEnt.Summary4LsaMap[lsaKey] = sLsa
	}
	server.AreaLsdb[lsdbKey] = lsDbEnt
	selfOrigLsaEnt[lsaKey] = true
	server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	var val LsdbSliceEnt
	val.AreaId = lsdbKey.AreaId
	val.LSType = lsaKey.LSType
	val.LSId = lsaKey.LSId
	val.AdvRtr = lsaKey.AdvRouter
	server.LsdbSlice = append(server.LsdbSlice, val)
	msg := DbLsdbMsg{
		entry: val,
		op:    true,
	}
	server.DbLsdbOp <- msg
}

func (server *OSPFServer) flushSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey) {
	server.logger.Debug(fmt.Sprintln("Need to flush Summary Lsa:", lsdbKey, lsaKey))
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	delete(selfOrigLsaEnt, lsaKey)
	server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	if lsaKey.LSType == Summary3LSA {
		delete(lsDbEnt.Summary3LsaMap, lsaKey)
	} else if lsaKey.LSType == Summary4LSA {
		delete(lsDbEnt.Summary4LsaMap, lsaKey)
	}
	server.AreaLsdb[lsdbKey] = lsDbEnt

}

func (server *OSPFServer) installSummaryLsa() {
	ifkey := IntfConfKey{}
	nbr := NeighborConfKey{}
	server.logger.Debug("Installing summary Lsa...")
	for lsdbKey, sLsa := range server.SummaryLsDb {
		selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
		oldSelfOrigSummaryLsa := make(map[LsaKey]bool)
		for sKey, _ := range selfOrigLsaEnt {
			if sKey.LSType == Summary3LSA ||
				sKey.LSType == Summary4LSA {
				oldSelfOrigSummaryLsa[sKey] = true
			}
		}

		for sKey, sEnt := range sLsa {
			if selfOrigLsaEnt[sKey] == true {
				oldSelfOrigSummaryLsa[sKey] = false
				ret := server.compareSummaryLsa(lsdbKey, sKey, sEnt)
				if ret == false {
					server.updateSummaryLsa(lsdbKey, sKey, sEnt)
					//Flood Updated Summary LSA
					server.logger.Debug(fmt.Sprintln("LSDB: Send message to flood ", sKey, lsdbKey.AreaId))
					server.sendLsdbToNeighborEvent(ifkey, nbr, lsdbKey.AreaId, 0, 0, sKey, LSASUMMARYFLOOD)
				} else {
					continue
				}
			} else {
				server.insertSummaryLsa(lsdbKey, sKey, sEnt)
				//Flood New Summary LSA
				server.sendLsdbToNeighborEvent(ifkey, nbr, lsdbKey.AreaId, 0, 0, sKey, LSASUMMARYFLOOD)
				server.logger.Debug(fmt.Sprintln("LSDB: Send message to flood ", sKey, lsdbKey.AreaId))

			}
		}
		sLsa = nil
		server.SummaryLsDb[lsdbKey] = sLsa
		for sKey, ent := range oldSelfOrigSummaryLsa {
			if ent == true {
				server.flushSummaryLsa(lsdbKey, sKey)
			}
		}
		oldSelfOrigSummaryLsa = nil
	}
	server.SummaryLsDb = nil
}

func (server *OSPFServer) flushNetworkLSA(areaId uint32, key IntfConfKey) {
	ent := server.IntfConfMap[key]
	AreaId := convertIPv4ToUint32(ent.IfAreaId)
	if areaId != AreaId {
		return
	}
	if ent.IfFSMState <= config.Waiting {
		return
	}

	LSType := NetworkLSA
	LSId := convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
	AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]

	lsa, exist := server.getNetworkLsaFromLsdb(areaId, lsaKey)
	if exist == LsdbEntryNotFound {
		return
	}

	server.logger.Debug(fmt.Sprintln("FLUSH: Network lsa lsid ",
		ent.IfIpAddr, " adv_router ", server.ospfGlobalConf.RouterId))
	lsa.LsaMd.LSAge = config.MaxAge
	lsa_pkt := encodeNetworkLsa(lsa, lsaKey)
	// Add entry to the flush map which will be flooded to all neighbors
	maxAgeLsaMap[lsaKey] = lsa_pkt
	// Need to Flush these entries
	delete(lsDbEnt.NetworkLsaMap, lsaKey)
	delete(selfOrigLsaEnt, lsaKey)
	server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	server.AreaLsdb[lsdbKey] = lsDbEnt

}

func (server *OSPFServer) generateNetworkLSA(areaId uint32, key IntfConfKey, isDR bool) {

	//routerId := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	ent := server.IntfConfMap[key]
	AreaId := convertIPv4ToUint32(ent.IfAreaId)
	nbrmdata := ospfIntfToNbrMap[key]

	if areaId != AreaId {
		return
	}
	if ent.IfFSMState <= config.Waiting {
		return
	}

	netmask := convertIPv4ToUint32(ent.IfNetmask)
	attachedRtr := make([]uint32, 0)

	for index := range nbrmdata.nbrList {
		nbrKey := nbrmdata.nbrList[index]
		nbrConf := server.NeighborConfigMap[nbrKey]
		flag := false
		for i := 0; i < len(attachedRtr); i++ {

			if nbrConf.OspfNbrRtrId == attachedRtr[i] {
				flag = true
			}
		}
		if flag == false {
			attachedRtr = append(attachedRtr, nbrConf.OspfNbrRtrId)
		}
	}
	selfRtrId := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	attachedRtr = append(attachedRtr, selfRtrId)

	numOfAttachedRtr := len(attachedRtr)
	if numOfAttachedRtr <= 1 {
		return
	}

	server.logger.Debug(fmt.Sprintln("NetworkLSA: attached router"))
	for i := range attachedRtr {
		server.logger.Debug(fmt.Sprintln("NetworkLSA: ", i, " ", attachedRtr[i]))
	}
	LSType := NetworkLSA
	LSId := convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
	Options := uint8(2) // Need to be revisited
	LSAge := 0
	AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}

	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	entry, exist := lsDbEnt.NetworkLsaMap[lsaKey]
	entry.LsaMd.LSAge = 0
	entry.LsaMd.Options = Options
	if !exist {
		entry.LsaMd.LSSequenceNum = InitialSequenceNumber
	} else {
		entry.LsaMd.LSSequenceNum = entry.LsaMd.LSSequenceNum + 1
	}
	entry.LsaMd.LSChecksum = 0
	// Length of Network LSA Metadata (netmask)  = 4 bytes
	entry.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (4 * numOfAttachedRtr))
	entry.Netmask = netmask
	entry.AttachedRtr = make([]uint32, numOfAttachedRtr)
	for i := 0; i < numOfAttachedRtr; i++ {
		entry.AttachedRtr[i] = attachedRtr[i]
	}
	server.logger.Debug(fmt.Sprintln("Attached Routers:", entry.AttachedRtr))
	//server.logger.Debug(fmt.Sprintln("Self Originated Router LSA Key:", server.AreaSelfOrigLsa[lsdbKey]))
	LsaEnc := encodeNetworkLsa(entry, lsaKey)
	checksumOffset := uint16(14)
	entry.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	entry.LsaMd.LSAge = uint16(LSAge)
	lsDbEnt.NetworkLsaMap[lsaKey] = entry
	server.AreaLsdb[lsdbKey] = lsDbEnt
	selfOrigLsaEnt[lsaKey] = true
	server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt

	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsaKey.LSType
		val.LSId = lsaKey.LSId
		val.AdvRtr = lsaKey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}
	msg := DbEventMsg{
		eventType: config.LSA,
		eventInfo: "Generate network LSA" + ent.IfIpAddr.String(),
	}
	server.DbEventOp <- msg

	return
}

func (server *OSPFServer) constructStubLinkP2P(ent IntfConf, likType config.IfType) LinkDetail {
	var linkDetail LinkDetail
	/*
	   There are two forms that this stub link can take:

	   Option 1
	   Assuming that the neighboring router's IP
	   address is known, set the Link ID of the Type 3
	   link to the neighbor's IP address, the Link Data
	   to the mask 0xffffffff (indicating a host
	   route), and the cost to the interface's
	   configured output cost.[15]

	   Option 2
	   If a subnet has been assigned to the point-to-
	   point link, set the Link ID of the Type 3 link
	   to the subnet's IP address, the Link Data to the
	   subnet's mask, and the cost to the interface's
	   configured output cost.[16]

	*/

	ipAddr := convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
	netmask := convertIPv4ToUint32(ent.IfNetmask)
	linkDetail.LinkId = ipAddr & netmask
	linkDetail.LinkData = netmask
	linkDetail.LinkType = StubLink
	linkDetail.NumOfTOS = 0
	linkDetail.LinkMetric = uint16(ent.IfCost)

	return linkDetail
}

func (server *OSPFServer) generateRouterLSA(areaId uint32) {
	var linkDetails []LinkDetail = nil
	for key, ent := range server.IntfConfMap {
		ipKey := convertAreaOrRouterIdUint32(string(key.IPAddr))
		ip, _ := server.ipPropertyMap[ipKey]

		AreaId := convertIPv4ToUint32(ent.IfAreaId)
		if areaId != AreaId {
			server.logger.Warning(fmt.Sprintln("LSDB: Area id not matching. i/p ", areaId, "if areaid ", AreaId, ent.IfIpAddr))
			continue
		}
		msg := DbEventMsg{
			eventType: config.LSA,
			eventInfo: "Generate router LSA " + ent.IfIpAddr.String(),
		}
		server.DbEventOp <- msg

		if ent.IfFSMState <= config.Waiting &&
			ip.IfType != commonDefs.IfTypeLoopback {
			server.logger.Debug(fmt.Sprintln("LSDB: If is in waiting. Skip.", ent.IfIpAddr))
			continue
		}
		var linkDetail LinkDetail
		switch ent.IfType {
		case config.Broadcast:
			if len(ent.NeighborMap) == 0 { // Stub Network
				server.logger.Debug("Stub Network")
				ipAddr := convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
				netmask := convertIPv4ToUint32(ent.IfNetmask)
				linkDetail.LinkId = ipAddr & netmask
				/* For links to stub networks, this field specifies the stub
				   network’s IP address mask. */
				linkDetail.LinkData = netmask
				linkDetail.LinkType = StubLink
				/* Todo: Need to handle IfMetricConf */
				linkDetail.NumOfTOS = 0
				linkDetail.LinkMetric = uint16(ent.IfCost)
			} else { // Transit Network
				server.logger.Debug("Transit Network")
				linkDetail.LinkId = convertIPv4ToUint32(ent.IfDRIp)
				/* For links to transit networks, numbered point-to-point links
				   and virtual links, this field specifies the IP interface
				   address of the associated router interface*/
				linkDetail.LinkData = convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
				linkDetail.LinkType = TransitLink
				/* Todo: Need to handle IfMetricConf */
				linkDetail.NumOfTOS = 0
				linkDetail.LinkMetric = uint16(ent.IfCost)
				server.logger.Debug(fmt.Sprintln("LinkDetail: linkid ", ent.IfDRIp,
					" linkdata ", ent.IfIpAddr))
			}
		case config.NumberedP2P:
			/*
						 In addition, as long as the state of the interface
				                    is "Point-to-Point" (and regardless of the
				                    neighboring router state), a Type 3 link (stub
				                    network) should be added.

			*/
			stub_link := server.constructStubLinkP2P(ent, config.NumberedP2P)
			linkDetails = append(linkDetails, stub_link)
			/* The Link ID should be
			   set to the Router ID of the neighboring router. For
			   numbered point-to-point networks, the Link Data
			   should specify the IP interface address. For
			   unnumbered point-to-point networks, the Link Data
			   field should specify the interface's MIB-II [Ref8]
			   ifIndex value. The cost should be set to the output
			   cost of the point-to-point interface.
			*/
			if len(ent.NeighborMap) == 0 {
				server.logger.Debug(fmt.Sprintln("LSDB: No neighbor detected for P2P link ", ent.IfIpAddr))
				continue
			}
			if nbrData, exist := ospfIntfToNbrMap[key]; exist {
				if len(nbrData.nbrList) != 0 {
					nbr := server.NeighborConfigMap[nbrData.nbrList[0]]
					server.logger.Debug(fmt.Sprintln("LSDB: Numbered P2P Router LSA with link id ", nbr.OspfNbrRtrId))
					linkDetail.LinkId = nbr.OspfNbrRtrId
				}

			}
			linkDetail.LinkData = convertAreaOrRouterIdUint32(ent.IfIpAddr.String())
			linkDetail.LinkType = P2PLink
			linkDetail.NumOfTOS = 0
			linkDetail.LinkMetric = uint16(ent.IfCost)

		case config.UnnumberedP2P:
			stub_link := server.constructStubLinkP2P(ent, config.UnnumberedP2P)
			linkDetails = append(linkDetails, stub_link)
			if len(ent.NeighborMap) == 0 {
				server.logger.Debug(fmt.Sprintln("LSDB: No neighbor detected for P2P link ", ent.IfIpAddr))
				continue
			}
			if nbrData, exist := ospfIntfToNbrMap[key]; exist {
				if len(nbrData.nbrList) != 0 {
					nbr := server.NeighborConfigMap[nbrData.nbrList[0]]
					linkDetail.LinkId = nbr.OspfNbrRtrId
					server.logger.Debug(fmt.Sprintln("LSDB: Unnumbered P2P Router LSA with link id ", nbr.OspfNbrRtrId))
				}

			}
			linkDetail.LinkData = uint32(key.IntfIdx)
			linkDetail.LinkType = P2PLink
			linkDetail.NumOfTOS = 0
			linkDetail.LinkMetric = uint16(ent.IfCost)
		}
		linkDetails = append(linkDetails, linkDetail)
	}

	numOfLinks := len(linkDetails)

	LSType := RouterLSA
	LSId := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	Options := uint8(2) // Need to be revisited
	LSAge := 0
	AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	BitE := false //not an AS boundary router (Todo)
	BitB := false
	if server.ospfGlobalConf.AreaBdrRtrStatus == true {
		BitB = true
	}
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}

	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, lsdbExist := server.AreaLsdb[lsdbKey]
	if !lsdbExist {
		server.logger.Err(fmt.Sprintln("LSDB: Area LSDB doesnt exist. No router LSA will be generated .. ", lsdbKey))
		return
	}
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]

	if numOfLinks == 0 {
		delete(lsDbEnt.RouterLsaMap, lsaKey)
		delete(selfOrigLsaEnt, lsaKey)
		server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.AreaLsdb[lsdbKey] = lsDbEnt
		return
	}
	ent, exist := lsDbEnt.RouterLsaMap[lsaKey]
	ent.LsaMd.LSAge = 0
	ent.LsaMd.Options = Options
	if !exist {
		ent.LsaMd.LSSequenceNum = InitialSequenceNumber
	} else {
		ent.LsaMd.LSSequenceNum = ent.LsaMd.LSSequenceNum + 1
	}
	ent.LsaMd.LSChecksum = 0
	// Length of Per Link Details = 12 bytes
	// Length of Router LSA Metadata (BitE, BitB, NumofLinks)  = 4 bytes
	ent.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (12 * numOfLinks))
	ent.BitE = BitE
	ent.BitB = BitB
	ent.NumOfLinks = uint16(numOfLinks)
	ent.LinkDetails = make([]LinkDetail, numOfLinks)
	copy(ent.LinkDetails, linkDetails[0:])
	server.logger.Debug(fmt.Sprintln("LinkDetails:", ent.LinkDetails))
	LsaEnc := encodeRouterLsa(ent, lsaKey)
	checksumOffset := uint16(14)
	ent.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	ent.LsaMd.LSAge = uint16(LSAge)
	lsDbEnt.RouterLsaMap[lsaKey] = ent
	server.AreaLsdb[lsdbKey] = lsDbEnt

	selfOrigLsaEnt[lsaKey] = true
	server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	server.logger.Debug(fmt.Sprintln("Self Originated Router LSA Key:", server.AreaSelfOrigLsa[lsdbKey]))
	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsaKey.LSType
		val.LSId = lsaKey.LSId
		val.AdvRtr = lsaKey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}
	return
}

func (server *OSPFServer) generateASExternalLsa(route RouteMdata) LsaKey {
	server.logger.Debug(fmt.Sprintln("LSDB: Generating AS External LSA routemdata ", route))

	LSType := ASExternalLSA
	LSId := route.ipaddr & route.mask
	AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)

	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}

	BitE := true
	for lsdbKey, _ := range server.AreaLsdb {
		lsDbEnt, _ := server.AreaLsdb[lsdbKey]
		ent, exist := lsDbEnt.ASExternalLsaMap[lsaKey]
		LSAge := 0
		ent.LsaMd.LSChecksum = 0
		ent.LsaMd.Options = 0x20
		ent.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 16)
		if !exist {
			ent.LsaMd.LSSequenceNum = InitialSequenceNumber
		} else {
			if route.isDel {
				ent.LsaMd.LSAge = LSA_MAX_AGE
			} else {
				ent.LsaMd.LSSequenceNum = ent.LsaMd.LSSequenceNum + 1
			}
		}
		ent.BitE = BitE
		ent.FwdAddr = convertAreaOrRouterIdUint32("0.0.0.0")
		ent.Metric = route.metric
		ent.Netmask = route.mask
		ent.ExtRouteTag = 0

		LsaEnc := encodeASExternalLsa(ent, lsaKey)
		checksumOffset := uint16(14)
		ent.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
		ent.LsaMd.LSAge = uint16(LSAge)
		lsDbEnt.ASExternalLsaMap[lsaKey] = ent
		server.AreaLsdb[lsdbKey] = lsDbEnt

		selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
		if !route.isDel {
			selfOrigLsaEnt[lsaKey] = true
			server.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		} else {
			selfOrigLsaEnt[lsaKey] = false
		}
		server.logger.Debug(fmt.Sprintln("ASBR: Added LSA to area ", lsdbKey, " lsaKey ", lsaKey))
		if !exist {
			var val LsdbSliceEnt
			val.AreaId = lsdbKey.AreaId
			val.LSType = lsaKey.LSType
			val.LSId = lsaKey.LSId
			val.AdvRtr = lsaKey.AdvRouter
			server.LsdbSlice = append(server.LsdbSlice, val)
			msg := DbLsdbMsg{
				entry: val,
				op:    true,
			}
			server.DbLsdbOp <- msg
		}
	}

	return lsaKey
}

func (server *OSPFServer) updateAsExternalLSA(lsdbKey LsdbKey, lsaKey LsaKey) error {
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if exist {
		ent, valid := lsDbEnt.ASExternalLsaMap[lsaKey]
		if !valid {
			server.logger.Warning(fmt.Sprintln("LSDB: AS external LSA doesnt exist lsdb ", lsdbKey, lsaKey))
			return nil
		}
		LsaEnc := encodeASExternalLsa(ent, lsaKey)
		checksumOffset := uint16(14)
		ent.LsaMd.LSChecksum = computeFletcherChecksum(LsaEnc[2:], checksumOffset)
		LSAge := 0
		ent.LsaMd.LSAge = uint16(LSAge)
		lsDbEnt.ASExternalLsaMap[lsaKey] = ent
		server.AreaLsdb[lsdbKey] = lsDbEnt

		//update db entry
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsaKey.LSType
		val.LSId = lsaKey.LSId
		val.AdvRtr = lsaKey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}
	return nil
}

func (server *OSPFServer) processDeleteRouterLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	routerLsa := NewRouterLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeRouterLsa(data, routerLsa, lsakey)
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	delete(lsDbEnt.RouterLsaMap, *lsakey)
	server.AreaLsdb[lsdbKey] = lsDbEnt
	return true
}

func (server *OSPFServer) processRecvdRouterLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	routerLsa := NewRouterLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeRouterLsa(data, routerLsa, lsakey)
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	_, exist := selfOrigLsaEnt[*lsakey]
	if exist {
		server.logger.Debug("Recvd a self generated Router LSA increment sequence number and flood. ")
		//increment sequence number and flood .
		//	return false
	}
	//Check Checksum
	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		server.logger.Err("Invalid Router LSA Checksum")
		return false
	}
	//Todo: If there is already existing entry Verify the seq num
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]

	//Add entry in LSADatabase
	lsDbEnt.RouterLsaMap[*lsakey] = *routerLsa
	server.AreaLsdb[lsdbKey] = lsDbEnt
	server.printRouterLsa()
	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsakey.LSType
		val.LSId = lsakey.LSId
		val.AdvRtr = lsakey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}
	server.logger.Debug(fmt.Sprintln("Router LSA: added to LSDB lsid ",
		lsakey.LSId, " adv_router ", lsakey.AdvRouter, " lstype ", lsakey.LSType))
	return true
}

func (server *OSPFServer) printRouterLsa() {
	server.logger.Debug("AREA  LSDB")
	for key, val := range server.AreaLsdb {
		server.logger.Debug(fmt.Sprintln("key ", key, " LSA ", val))
	}
}

func (server *OSPFServer) processDeleteNetworkLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	networkLsa := NewNetworkLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeNetworkLsa(data, networkLsa, lsakey)
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	delete(lsDbEnt.NetworkLsaMap, *lsakey)
	server.AreaLsdb[lsdbKey] = lsDbEnt

	return true
}

func (server *OSPFServer) processRecvdNetworkLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	networkLsa := NewNetworkLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeNetworkLsa(data, networkLsa, lsakey)
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	_, exist := selfOrigLsaEnt[*lsakey]
	if exist {
		server.logger.Debug("Recvd a self generated Network LSA")
		return false
	}

	//Check Checksum
	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		server.logger.Err("Invalid Network LSA Checksum")
		return false
	}
	//Todo: If there is already existing entry Verify the seq num
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	ent, exist := lsDbEnt.NetworkLsaMap[*lsakey]
	if exist {
		if ent.LsaMd.LSSequenceNum >= networkLsa.LsaMd.LSSequenceNum {
			server.logger.Err("Old instance of Network LSA Recvd")
			return false
		}
	}
	//Handle LsaAge
	//Add entry in LSADatabase
	lsDbEnt.NetworkLsaMap[*lsakey] = *networkLsa
	server.AreaLsdb[lsdbKey] = lsDbEnt
	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsakey.LSType
		val.LSId = lsakey.LSId
		val.AdvRtr = lsakey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
	}
	server.logger.Debug(fmt.Sprintln("Network LSA: added to LSDB lsid ",
		lsakey.LSId, " adv_router ", lsakey.AdvRouter, " lstype ", lsakey.LSType))
	return true
}

func (server *OSPFServer) processDeleteSummaryLsa(data []byte, areaId uint32, lsaType uint8) bool {
	lsakey := NewLsaKey()
	var val LsdbSliceEnt
	summaryLsa := NewSummaryLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	val.AreaId = lsdbKey.AreaId

	decodeSummaryLsa(data, summaryLsa, lsakey)

	val.LSType = lsakey.LSType
	val.LSId = lsakey.LSId
	val.AdvRtr = lsakey.AdvRouter

	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	if lsaType == Summary3LSA {
		delete(lsDbEnt.Summary3LsaMap, *lsakey)
	} else if lsaType == Summary4LSA {
		delete(lsDbEnt.Summary4LsaMap, *lsakey)
	}
	server.AreaLsdb[lsdbKey] = lsDbEnt
	server.printRouterLsa()
	server.DelLsdbEntry(val)
	return true
}

func (server *OSPFServer) processRecvdSummaryLsa(data []byte, areaId uint32, lsaType uint8) bool {
	lsakey := NewLsaKey()
	summaryLsa := NewSummaryLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeSummaryLsa(data, summaryLsa, lsakey)

	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	_, exist := selfOrigLsaEnt[*lsakey]
	if exist {
		server.logger.Debug("Recvd a self generated Summary LSA")
		return false
	}

	//Check Checksum
	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		server.logger.Err("Invalid Summary LSA Checksum")
		return false
	}
	//Todo: If there is already existing entry Verify the seq num
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	if lsaType == Summary3LSA {
		ent, exist := lsDbEnt.Summary3LsaMap[*lsakey]
		if exist {
			if ent.LsaMd.LSSequenceNum >= summaryLsa.LsaMd.LSSequenceNum {
				server.logger.Err("Old instance of Summary 3 LSA Recvd")
				return false
			}
		}
		//Handle LsaAge
		//Add entry in LSADatabase
		lsDbEnt.Summary3LsaMap[*lsakey] = *summaryLsa
	} else if lsaType == Summary4LSA {
		ent, exist := lsDbEnt.Summary4LsaMap[*lsakey]
		if exist {
			if ent.LsaMd.LSSequenceNum >= summaryLsa.LsaMd.LSSequenceNum {
				server.logger.Err("Old instance of Summary 4 LSA Recvd")
				return false
			}
		}
		//Handle LsaAge
		//Add entry in LSADatabase
		lsDbEnt.Summary4LsaMap[*lsakey] = *summaryLsa
	} else {
		return false
	}
	server.AreaLsdb[lsdbKey] = lsDbEnt
	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsakey.LSType
		val.LSId = lsakey.LSId
		val.AdvRtr = lsakey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}
	return true
}

func (server *OSPFServer) processDeleteASExternalLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	var val LsdbSliceEnt
	asExtLsa := NewASExternalLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeASExternalLsa(data, asExtLsa, lsakey)
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	delete(lsDbEnt.ASExternalLsaMap, *lsakey)
	server.AreaLsdb[lsdbKey] = lsDbEnt

	val.AreaId = lsdbKey.AreaId
	val.LSType = lsakey.LSType
	val.LSId = lsakey.LSId
	val.AdvRtr = lsakey.AdvRouter
	err := server.DelLsdbEntry(val)
	if err != nil {
		server.logger.Debug(fmt.Sprintln("DB: Failed to delete entry from db ", lsakey))
	}
	return true
}

func (server *OSPFServer) processRecvdASExternalLsa(data []byte, areaId uint32) bool {
	lsakey := NewLsaKey()
	asExtLsa := NewASExternalLsa()
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	decodeASExternalLsa(data, asExtLsa, lsakey)
	selfOrigLsaEnt, _ := server.AreaSelfOrigLsa[lsdbKey]
	_, exist := selfOrigLsaEnt[*lsakey]
	if exist {
		server.logger.Debug("Recvd a self generated AS External LSA")
		return false
	}

	//Check Checksum
	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		server.logger.Err("Invalid AS External LSA Checksum")
		return false
	}
	//Todo: If there is already existing entry Verify the seq num
	lsDbEnt, _ := server.AreaLsdb[lsdbKey]
	ent, exist := lsDbEnt.ASExternalLsaMap[*lsakey]
	if exist {
		if ent.LsaMd.LSSequenceNum >= asExtLsa.LsaMd.LSSequenceNum {
			server.logger.Err("Old instance of AS External LSA Recvd")
			return false
		}
	}
	//Handle LsaAge
	//Add entry in LSADatabase
	lsDbEnt.ASExternalLsaMap[*lsakey] = *asExtLsa
	server.AreaLsdb[lsdbKey] = lsDbEnt
	if !exist {
		var val LsdbSliceEnt
		val.AreaId = lsdbKey.AreaId
		val.LSType = lsakey.LSType
		val.LSId = lsakey.LSId
		val.AdvRtr = lsakey.AdvRouter
		server.LsdbSlice = append(server.LsdbSlice, val)
		msg := DbLsdbMsg{
			entry: val,
			op:    true,
		}
		server.DbLsdbOp <- msg
	}

	return true
}

func (server *OSPFServer) processRecvdLsa(data []byte, areaId uint32) bool {
	LSType := uint8(data[3])
	if LSType == RouterLSA {
		server.logger.Debug("LSDB: Received router lsa")
		return server.processRecvdRouterLsa(data, areaId)
	} else if LSType == NetworkLSA {
		server.logger.Debug("LSDB: Received network lsa")
		return server.processRecvdNetworkLsa(data, areaId)
	} else if LSType == Summary3LSA {
		server.logger.Debug("LSDB: Received summary3 lsa")
		return server.processRecvdSummaryLsa(data, areaId, LSType)
	} else if LSType == Summary4LSA {
		server.logger.Debug("LSDB: Received summary4 lsa")
		return server.processRecvdSummaryLsa(data, areaId, LSType)
	} else if LSType == ASExternalLSA {
		return server.processRecvdASExternalLsa(data, areaId)
	} else {
		server.logger.Debug("LSDB: Invalid LSA packet from nbr")
		return false
	}
}

func (server *OSPFServer) processDeleteLsa(data []byte, areaId uint32) bool {
	LSType := uint8(data[3])

	if LSType == RouterLSA {
		return server.processDeleteRouterLsa(data, areaId)
	} else if LSType == NetworkLSA {
		return server.processDeleteNetworkLsa(data, areaId)
	} else if LSType == Summary3LSA {
		return server.processDeleteSummaryLsa(data, areaId, LSType)
	} else if LSType == Summary4LSA {
		return server.processDeleteSummaryLsa(data, areaId, LSType)
	} else if LSType == ASExternalLSA {
		return server.processDeleteASExternalLsa(data, areaId)
	} else {
		return false
	}
}

func (server *OSPFServer) processLSDatabaseUpdates() {
	for {
		select {
		case msg := <-server.LsdbUpdateCh:
			if msg.MsgType == LsdbAdd {
				server.logger.Debug("Received New LSA")
				ret := server.processRecvdLsa(msg.Data, msg.AreaId)
				server.logger.Debug(fmt.Sprintln("Return Code:", ret))
				//server.LsaUpdateRetCodeCh <- ret
				server.StartCalcSPFCh <- true
				spfStatus := <-server.DoneCalcSPFCh
				server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
				if server.ospfGlobalConf.AreaBdrRtrStatus == true {
					server.installSummaryLsa()
				}
			} else if msg.MsgType == LsdbDel {
				server.logger.Debug("Deleting LS in the Lsdb")
				ret := server.processDeleteLsa(msg.Data, msg.AreaId)
				//server.LsaUpdateRetCodeCh <- ret
				server.logger.Debug(fmt.Sprintln("Return Code:", ret))
				server.StartCalcSPFCh <- true
				spfStatus := <-server.DoneCalcSPFCh
				server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
				if server.ospfGlobalConf.AreaBdrRtrStatus == true {
					server.installSummaryLsa()
				}
			} else if msg.MsgType == LsdbUpdate {
				server.logger.Debug("Deleting LS in the Lsdb")
				ret := server.processRecvdLsa(msg.Data, msg.AreaId)
				//server.LsaUpdateRetCodeCh <- ret
				server.logger.Debug(fmt.Sprintln("Return Code:", ret))
				server.StartCalcSPFCh <- true
				spfStatus := <-server.DoneCalcSPFCh
				server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
				if server.ospfGlobalConf.AreaBdrRtrStatus == true {
					server.installSummaryLsa()
				}
			}
		case msg := <-server.IntfStateChangeCh:
			server.logger.Debug(fmt.Sprintf("Interface State change msg", msg))
			server.generateRouterLSA(msg.areaId)
			//server.logger.Debug(fmt.Sprintln("LS Database", server.AreaLsdb))
			server.StartCalcSPFCh <- true
			spfStatus := <-server.DoneCalcSPFCh
			server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
			server.processInterfaceChangeMsg(msg)
			if server.ospfGlobalConf.AreaBdrRtrStatus == true {
				server.installSummaryLsa()
			}
		case msg := <-server.NetworkDRChangeCh:
			server.logger.Debug(fmt.Sprintf("Network DR change msg", msg))
			// Create a new router LSA
			//server.logger.Debug(fmt.Sprintln("LS Database", server.AreaLsdb))
			server.processDrBdrChangeMsg(msg)
			server.StartCalcSPFCh <- true
			spfStatus := <-server.DoneCalcSPFCh
			server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
			if server.ospfGlobalConf.AreaBdrRtrStatus == true {
				server.installSummaryLsa()
			}
		case msg := <-server.CreateNetworkLSACh:
			server.logger.Debug(fmt.Sprintf("Create Network LSA msg", msg))
			server.processNeighborFullEvent(msg)
			//server.generateNetworkLSA(msg.areaId, msg.intf, msg.isDR)
			// Flush the old Network LSA
			// Check if link is broadcast or not
			// If link is broadcast
			// Create Network LSA
			//server.logger.Debug(fmt.Sprintln("LS Database", server.AreaLsdb))
			server.StartCalcSPFCh <- true
			spfStatus := <-server.DoneCalcSPFCh
			server.logger.Debug(fmt.Sprintln("SPF Calculation Return Status", spfStatus))
			if server.ospfGlobalConf.AreaBdrRtrStatus == true {
				server.installSummaryLsa()
			}

		case msg := <-server.ExternalRouteNotif: //Generate external LSA
			server.processExtRouteUpd(msg)

		case msg := <-server.maxAgeLsaCh: //Flood MaxAge LSA
			server.processMaxAgeLsaMsg(msg)

		case <-lsdbTickerCh.C: //Increment LSA AGE
			lsdbTickerCh.Stop()
			server.processLSDatabaseTicker()
			lsdbTickerCh.Reset(time.Duration(1) * time.Second)

		case <-lsdbRefreshTickerCh.C: //Regenerate LSA
			lsdbRefreshTickerCh.Stop()
			server.lsdbSelfLsaRefresh()
			lsdbRefreshTickerCh.Reset(time.Duration(config.LSRefreshTime) * time.Second)
		}
	}
}

/*@fn processExtRouteUpd
Generate / delete As external LSA.
Send flood message if new route is added.
*/
func (server *OSPFServer) processExtRouteUpd(msg RouteMdata) {
	ifkey := IntfConfKey{}
	nbr := NeighborConfKey{}
	lsaKey := server.generateASExternalLsa(msg)
	if !msg.isDel {
		server.sendLsdbToNeighborEvent(ifkey, nbr, 0, 0, 0, lsaKey, LSAEXTFLOOD)
	}
}

/*
@fn processIntfStateEventNbr
- Cleaning for neighbor data structures.
-Flood nw/router lsas based on
interface up/down add/delete events.
*/
func (server *OSPFServer) processInterfaceChangeMsg(msg NetworkLSAChangeMsg) {
	if msg.intfState == config.Intf_Down {
		server.neighborIntfEventCh <- msg.intfKey
	}
	nbr := NeighborConfKey{}
	LSType := RouterLSA
	LSId := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}

	// send message for flooding.
	server.sendLsdbToNeighborEvent(msg.intfKey, nbr, msg.areaId, 0, 0, lsaKey, LSAROUTERFLOOD)
}

/* @fn processNeighborFullEvent
Generate network LSA if the router is DR.
Send message for LSAFLOOD which will flood
router (and network) LSA
*/
func (server *OSPFServer) processNeighborFullEvent(msg ospfNbrMdata) {
	lsaKey := LsaKey{}
	nbr := NeighborConfKey{}

	rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
	intConf := server.IntfConfMap[msg.intf]
	server.logger.Debug(fmt.Sprintln("LSDB: Nbr full. Generate router and network LSA  area id  ",
		msg.areaId, " intf ", intConf.IfIpAddr))
	if intConf.IfDRtrId == rtr_id && intConf.IfType == config.Broadcast {
		server.logger.Debug(fmt.Sprintln("Generate network LSA ", msg.intf))
		server.generateNetworkLSA(msg.areaId, msg.intf, true)
	}
	server.generateRouterLSA(msg.areaId)
	server.sendLsdbToNeighborEvent(msg.intf, nbr, msg.areaId, 0, 0, lsaKey, LSAFLOOD)
}

/* @fn processDrBdrChangeMsg
when DR changes
generate network and router LSA if I am DR.
generate router LSA if I am not DR. Also flush
network LSA is I am become DR to no DR.
*/
func (server *OSPFServer) processDrBdrChangeMsg(msg DrChangeMsg) {
	/* check if any nbr attached to the intf if not dont generate network LSA
	 */
	lsaKey := LsaKey{}
	nbr := NeighborConfKey{}
	intf, _ := server.IntfConfMap[msg.intfKey]
	server.logger.Debug(fmt.Sprintln("LSDB: received DR BDR change message ",
		intf.IfIpAddr, "dr ip ", intf.IfDRIp, " bdr ip ", intf.IfBDRIp))
	nbrExists := false
	if len(intf.NeighborMap) > 0 {
		nbrExists = true
	}
	if nbrExists {
		rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
		if intf.IfDRtrId == rtr_id {
			server.logger.Debug(fmt.Sprintln("Generate network LSA ", intf.IfIpAddr))
			server.generateNetworkLSA(msg.areaId, msg.intfKey, true)
		}
	}
	server.generateRouterLSA(msg.areaId)
	server.sendLsdbToNeighborEvent(msg.intfKey, nbr, msg.areaId, 0, 0, lsaKey, LSAFLOOD)

}
