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
	"errors"
	"l3/ospfv2/objects"
	"time"
)

func (server *OSPFV2Server) InitLsdbData() {
	server.LsdbData.AreaLsdb = make(map[LsdbKey]LSDatabase)
	server.LsdbData.AreaSelfOrigLsa = make(map[LsdbKey]SelfOrigLsa)
	server.LsdbData.LsdbAgingTicker = nil
	server.LsdbData.ExtRouteInfoMap = make(map[RouteInfo]bool)
}

func (server *OSPFV2Server) DeinitLsdb() {
	server.LsdbData.LsdbAgingTicker = nil
	for lsdbKey, _ := range server.LsdbData.AreaLsdb {
		delete(server.LsdbData.AreaLsdb, lsdbKey)
	}
	for lsdbKey, _ := range server.LsdbData.AreaSelfOrigLsa {
		delete(server.LsdbData.AreaSelfOrigLsa, lsdbKey)
	}
	server.LsdbData.AreaLsdb = nil
	server.LsdbData.AreaSelfOrigLsa = nil
	server.LsdbData.ExtRouteInfoMap = nil
}

func (server *OSPFV2Server) GetExtRouteInfo() {
	routeInfoList := server.getBulkRoutesFromRibd()
	for _, route := range routeInfoList {
		server.LsdbData.ExtRouteInfoMap[*route] = true
		server.generateASExternalLSA(*route)
	}
}

func (server *OSPFV2Server) InitAreaLsdb(areaId uint32) {
	server.logger.Debug("LSDB: Initialise LSDB for area id ", areaId)
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		lsDbEnt.RouterLsaMap = make(map[LsaKey]RouterLsa)
		lsDbEnt.NetworkLsaMap = make(map[LsaKey]NetworkLsa)
		lsDbEnt.Summary3LsaMap = make(map[LsaKey]SummaryLsa)
		lsDbEnt.Summary4LsaMap = make(map[LsaKey]SummaryLsa)
		lsDbEnt.ASExternalLsaMap = make(map[LsaKey]ASExternalLsa)
		server.LsdbData.AreaLsdb[lsdbKey] = lsDbEnt
	}
	selfOrigLsaEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	if !exist {
		selfOrigLsaEnt = make(map[LsaKey]bool)
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	}

}
func (server *OSPFV2Server) DeinitAreaLsdb(areaId uint32) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if exist {
		lsDbEnt.RouterLsaMap = nil
		lsDbEnt.NetworkLsaMap = nil
		lsDbEnt.Summary3LsaMap = nil
		lsDbEnt.Summary4LsaMap = nil
		lsDbEnt.ASExternalLsaMap = nil
		delete(server.LsdbData.AreaLsdb, lsdbKey)
	}
	_, exist = server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	if exist {
		delete(server.LsdbData.AreaSelfOrigLsa, lsdbKey)
	}
}

func (server *OSPFV2Server) FlushAreaLsdb(areaId uint32) {
	server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlCh <- areaId
	<-server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlReplyCh
}

func (server *OSPFV2Server) StartLsdbRoutine() {
	server.LsdbData.LsdbCtrlChData.LsdbGblCtrlCh = make(chan bool)
	server.LsdbData.LsdbCtrlChData.LsdbGblCtrlReplyCh = make(chan bool)
	server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlCh = make(chan uint32)
	server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlReplyCh = make(chan uint32)
	initDoneCh := make(chan bool)
	go server.ProcessLsdb(initDoneCh)
	<-initDoneCh
}

func (server *OSPFV2Server) StopLsdbRoutine() {
	server.LsdbData.LsdbCtrlChData.LsdbGblCtrlCh <- true
	cnt := 0
	for {
		select {
		case _ = <-server.LsdbData.LsdbCtrlChData.LsdbGblCtrlReplyCh:
			server.logger.Info("Successfully Stopped ProcessLsdb routine")
			server.LsdbData.LsdbCtrlChData.LsdbGblCtrlCh = nil
			server.LsdbData.LsdbCtrlChData.LsdbGblCtrlReplyCh = nil
			server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlCh = nil
			server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlReplyCh = nil
			return
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt%1000 == 0 {
				server.logger.Err("Trying to stop the ProcessLsdb routine")
				//return
			}
		}
	}
}

func (server *OSPFV2Server) processRecvdLSA(msg RecvdLsaMsg) error {
	switch msg.LsaKey.LSType {
	case RouterLSA:
		server.processRecvdRouterLSA(msg)
	case NetworkLSA:
		server.processRecvdNetworkLSA(msg)
	case Summary3LSA:
		server.processRecvdSummaryLSA(msg)
	case Summary4LSA:
		server.processRecvdSummaryLSA(msg)
	case ASExternalLSA:
		server.processRecvdASExternalLSA(msg)
	default:
		server.logger.Err("Invalid LsaType:", msg)
	}
	return nil
}

func (server *OSPFV2Server) processRecvdSelfLSA(msg RecvdSelfLsaMsg) error {
	switch msg.LsaKey.LSType {
	case RouterLSA:
		server.processRecvdSelfRouterLSA(msg)
	case NetworkLSA:
		server.processRecvdSelfNetworkLSA(msg)
	case Summary3LSA:
		server.processRecvdSelfSummaryLSA(msg)
	case Summary4LSA:
		server.processRecvdSelfSummaryLSA(msg)
	case ASExternalLSA:
		server.processRecvdSelfASExternalLSA(msg)
	default:
		server.logger.Err("Invalid LsaType:", msg)
	}
	return nil
}

func (server *OSPFV2Server) ProcessRouteInfoData(msg RouteInfoDataUpdateMsg) {
	if msg.MsgType == ROUTE_INFO_ADD {
		for _, routeInfo := range msg.RouteInfoList {
			server.LsdbData.ExtRouteInfoMap[routeInfo] = true
			server.generateASExternalLSA(routeInfo)
		}
	} else if msg.MsgType == ROUTE_INFO_DEL {
		for _, routeInfo := range msg.RouteInfoList {
			delete(server.LsdbData.ExtRouteInfoMap, routeInfo)
			server.flushASExternalLSA(routeInfo)
		}
	} else {
		server.logger.Err("Invalid MsgType for RouteInfoDataUpdateMsg")
	}
}

func (server *OSPFV2Server) processNbrDead(msg NbrDeadMsg) bool {
	lsdbKey := LsdbKey{
		AreaId: msg.AreaId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("No Lsdb Exist for given Area:", lsdbKey)
		return false
	}
	flag := false
	for lsaKey, lsaEnt := range lsdbEnt.RouterLsaMap {
		if lsaKey.AdvRouter == msg.NbrRtrId {
			flag = true
			delete(lsdbEnt.RouterLsaMap, lsaKey)
			lsaEnt.LsaMd.LSAge = MAX_AGE
			server.CreateAndSendMsgFromLsdbToFloodLsa(msg.AreaId, lsaKey, lsaEnt)
			break
		}
	}
	for lsaKey, lsaEnt := range lsdbEnt.NetworkLsaMap {
		if lsaKey.AdvRouter == msg.NbrRtrId {
			flag = true
			delete(lsdbEnt.NetworkLsaMap, lsaKey)
			lsaEnt.LsaMd.LSAge = MAX_AGE
			server.CreateAndSendMsgFromLsdbToFloodLsa(msg.AreaId, lsaKey, lsaEnt)
		}
	}
	for lsaKey, lsaEnt := range lsdbEnt.Summary3LsaMap {
		if lsaKey.AdvRouter == msg.NbrRtrId {
			flag = true
			delete(lsdbEnt.Summary3LsaMap, lsaKey)
			lsaEnt.LsaMd.LSAge = MAX_AGE
			server.CreateAndSendMsgFromLsdbToFloodLsa(msg.AreaId, lsaKey, lsaEnt)
		}
	}
	for lsaKey, lsaEnt := range lsdbEnt.Summary4LsaMap {
		if lsaKey.AdvRouter == msg.NbrRtrId {
			flag = true
			delete(lsdbEnt.Summary4LsaMap, lsaKey)
			lsaEnt.LsaMd.LSAge = MAX_AGE
			server.CreateAndSendMsgFromLsdbToFloodLsa(msg.AreaId, lsaKey, lsaEnt)
		}
	}
	for lsaKey, lsaEnt := range lsdbEnt.ASExternalLsaMap {
		if lsaKey.AdvRouter == msg.NbrRtrId {
			flag = true
			delete(lsdbEnt.ASExternalLsaMap, lsaKey)
			lsaEnt.LsaMd.LSAge = MAX_AGE
			server.CreateAndSendMsgFromLsdbToFloodLsa(msg.AreaId, lsaKey, lsaEnt)
		}
	}
	if flag == true {
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
		return true
	}
	return false
}

func (server *OSPFV2Server) ProcessLsdb(initDoneCh chan bool) {
	server.InitLsdbData()
	for areaId, areaEnt := range server.AreaConfMap {
		if areaEnt.AdminState == true {
			server.InitAreaLsdb(areaId)
		}
	}
	server.GetExtRouteInfo()
	server.LsdbData.LsdbAgingTicker = time.NewTicker(LsaAgingTimeGranularity)
	initDoneCh <- true
	for {
		select {
		case _ = <-server.LsdbData.LsdbCtrlChData.LsdbGblCtrlCh:
			server.logger.Info("Stopping ProcessLsdb routine")
			server.LsdbData.LsdbAgingTicker.Stop()
			server.DeinitLsdb()
			server.LsdbData.LsdbCtrlChData.LsdbGblCtrlReplyCh <- true
			return
		case areaId := <-server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlCh:
			server.DeinitAreaLsdb(areaId)
			server.RefreshLsdbSlice()
			server.CalcSPFAndRoutingTbl()
			server.LsdbData.LsdbCtrlChData.LsdbAreaCtrlReplyCh <- areaId
		case areaId := <-server.MessagingChData.ServerToLsdbChData.InitAreaLsdbCh:
			server.logger.Info("InitAreaLsdb...")
			server.InitAreaLsdb(areaId)
			server.logger.Info("InitAreaLsdb...")
			server.SendMsgFromLsdbToServerForInitAreaLsdbDone()
			server.GenerateAllASExternalLSA(areaId)
		case msg := <-server.MessagingChData.IntfFSMToLsdbChData.GenerateRouterLSACh:
			server.logger.Info("Generate self originated Router LSA", msg)
			err := server.GenerateRouterLSA(msg)
			if err != nil {
				continue
			}
			server.logger.Info("Successfully Generated Router LSA")
			server.CalcSPFAndRoutingTbl()
			server.logger.Info("Successfully Calculated SPF")
		case msg := <-server.MessagingChData.NbrFSMToLsdbChData.UpdateSelfNetworkLSACh:
			server.logger.Info("Update self originated Network LSA", msg)
			err := server.processUpdateSelfNetworkLSA(msg)
			if err != nil {
				continue
			}
			server.CalcSPFAndRoutingTbl()
		case msg := <-server.MessagingChData.NbrFSMToLsdbChData.RecvdLsaMsgCh:
			server.logger.Info("Update LSA", msg)
			server.processRecvdLSA(msg)
			server.CalcSPFAndRoutingTbl()
		case msg := <-server.MessagingChData.NbrFSMToLsdbChData.RecvdSelfLsaMsgCh:
			server.logger.Info("Recvd Self LSA", msg)
			server.processRecvdSelfLSA(msg)
			server.CalcSPFAndRoutingTbl()
		case msg := <-server.MessagingChData.NbrFSMToLsdbChData.NbrDeadMsgCh:
			server.logger.Info("Recvd Nbr Dead in Lsdb:", msg)
			ret := server.processNbrDead(msg)
			if ret == true {
				server.CalcSPFAndRoutingTbl()
			}
		case msg := <-server.MessagingChData.ServerToLsdbChData.RouteInfoDataUpdateCh:
			//TODO: Handle AS External
			server.ProcessRouteInfoData(msg)
		case <-server.LsdbData.LsdbAgingTicker.C:
			server.processLsdbAgingTicker()
		case <-server.MessagingChData.ServerToLsdbChData.RefreshLsdbSliceCh:
			server.RefreshLsdbSlice()
			server.SendMsgFromLsdbToServerForRefreshDone()
		}
	}
}

func (server *OSPFV2Server) CalcSPFAndRoutingTbl() {
	server.SummaryLsDb = nil
	server.SendMsgToStartSpf()
	spfState := <-server.MessagingChData.SPFToLsdbChData.DoneSPF
	server.logger.Debug("SPF Calculation Return Status", spfState)
	if server.globalData.AreaBdrRtrStatus == true {
		server.logger.Info("Examine transit areas, Summary LSA...")
		server.HandleTransitAreaSummaryLsa()
		server.logger.Info("Generate Summary LSA...")
		server.GenerateSummaryLsa()
		server.logger.Info("========", server.SummaryLsDb, "==========")
		//Summary LSA
		server.installSummaryLsa()
	}
}

func (server *OSPFV2Server) RefreshLsdbSlice() {
	if len(server.GetBulkData.LsdbSlice) == 0 {
		return
	}
	server.GetBulkData.LsdbSlice = server.GetBulkData.LsdbSlice[:len(server.GetBulkData.LsdbSlice)-1]
	server.GetBulkData.LsdbSlice = nil
	for lsdbKey, lsDbEnt := range server.LsdbData.AreaLsdb {
		for lsaKey, _ := range lsDbEnt.RouterLsaMap {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
		for lsaKey, _ := range lsDbEnt.NetworkLsaMap {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
		for lsaKey, _ := range lsDbEnt.Summary3LsaMap {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
		for lsaKey, _ := range lsDbEnt.Summary4LsaMap {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
		for lsaKey, _ := range lsDbEnt.ASExternalLsaMap {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	}
}

func (server *OSPFV2Server) getLsdbState(lsaType uint8, lsId, areaId, advRtrId uint32) (*objects.Ospfv2LsdbState, error) {
	var retObj objects.Ospfv2LsdbState
	server.logger.Info("Lsdb Get for ", lsaType, lsId, areaId, advRtrId)
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		return nil, errors.New("No such area exist")
	}
	lsaKey := LsaKey{
		LSType:    lsaType,
		LSId:      lsId,
		AdvRouter: advRtrId,
	}
	var lsaMd LsaMetadata
	var lsaEnc []byte
	switch lsaType {
	case RouterLSA:
		lsaEnt, exist := lsdbEnt.RouterLsaMap[lsaKey]
		if !exist {
			return nil, errors.New("No such LSA exist")
		}
		lsaMd = lsaEnt.LsaMd
		lsaEnc = encodeRouterLsa(lsaEnt, lsaKey)
	case NetworkLSA:
		lsaEnt, exist := lsdbEnt.NetworkLsaMap[lsaKey]
		if !exist {
			return nil, errors.New("No such LSA exist")
		}
		lsaMd = lsaEnt.LsaMd
		lsaEnc = encodeNetworkLsa(lsaEnt, lsaKey)
	case Summary3LSA:
		lsaEnt, exist := lsdbEnt.Summary3LsaMap[lsaKey]
		if !exist {
			return nil, errors.New("No such LSA exist")
		}
		lsaMd = lsaEnt.LsaMd
		lsaEnc = encodeSummaryLsa(lsaEnt, lsaKey)
	case Summary4LSA:
		lsaEnt, exist := lsdbEnt.Summary4LsaMap[lsaKey]
		if !exist {
			return nil, errors.New("No such LSA exist")
		}
		lsaMd = lsaEnt.LsaMd
		lsaEnc = encodeSummaryLsa(lsaEnt, lsaKey)
	case ASExternalLSA:
		lsaEnt, exist := lsdbEnt.ASExternalLsaMap[lsaKey]
		if !exist {
			return nil, errors.New("No such LSA exist")
		}
		lsaMd = lsaEnt.LsaMd
		lsaEnc = encodeASExternalLsa(lsaEnt, lsaKey)
	default:
		return nil, errors.New("Invalid LSType")
	}
	retObj.LSType = lsaType
	retObj.LSId = lsId
	retObj.AdvRouterId = advRtrId
	retObj.AreaId = areaId
	retObj.SequenceNum = uint32(lsaMd.LSSequenceNum)
	retObj.Age = lsaMd.LSAge
	retObj.Checksum = lsaMd.LSChecksum
	retObj.Options = lsaMd.Options
	retObj.Length = lsaMd.LSLen
	retObj.Advertisement = convertByteToOctetString(lsaEnc[OSPF_LSA_HEADER_SIZE:])
	return &retObj, nil
}

func (server *OSPFV2Server) getBulkLsdbState(fromIdx, cnt int) (*objects.Ospfv2LsdbStateGetInfo, error) {
	var retObj objects.Ospfv2LsdbStateGetInfo
	var lsdbSliceMap map[LsdbSliceStruct]bool
	count := 0
	sliceLen := len(server.GetBulkData.LsdbSlice)
	if fromIdx >= sliceLen {
		return nil, errors.New("Invalid Range")
	}
	lsdbSliceMap = make(map[LsdbSliceStruct]bool)
	for idx := 0; idx < fromIdx; idx++ {
		lsdbSliceMap[server.GetBulkData.LsdbSlice[idx]] = true
	}
	idx := fromIdx
	for count < cnt {
		if idx == sliceLen {
			break
		}
		lsdbSlice := server.GetBulkData.LsdbSlice[idx]
		_, exist := lsdbSliceMap[lsdbSlice]
		if exist {
			idx++
			continue
		}
		lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbSlice.LsdbKey]
		if !exist {
			idx++
			continue
		}
		var lsaMd LsaMetadata
		var lsaEnc []byte
		switch lsdbSlice.LsaKey.LSType {
		case RouterLSA:
			lsaEnt, exist := lsdbEnt.RouterLsaMap[lsdbSlice.LsaKey]
			if !exist {
				idx++
				continue
			}
			lsaMd = lsaEnt.LsaMd
			lsaEnc = encodeRouterLsa(lsaEnt, lsdbSlice.LsaKey)
		case NetworkLSA:
			lsaEnt, exist := lsdbEnt.NetworkLsaMap[lsdbSlice.LsaKey]
			if !exist {
				idx++
				continue
			}
			lsaMd = lsaEnt.LsaMd
			lsaEnc = encodeNetworkLsa(lsaEnt, lsdbSlice.LsaKey)
		case Summary3LSA:
			lsaEnt, exist := lsdbEnt.Summary3LsaMap[lsdbSlice.LsaKey]
			if !exist {
				idx++
				continue
			}
			lsaMd = lsaEnt.LsaMd
			lsaEnc = encodeSummaryLsa(lsaEnt, lsdbSlice.LsaKey)
		case Summary4LSA:
			lsaEnt, exist := lsdbEnt.Summary4LsaMap[lsdbSlice.LsaKey]
			if !exist {
				idx++
				continue
			}
			lsaMd = lsaEnt.LsaMd
			lsaEnc = encodeSummaryLsa(lsaEnt, lsdbSlice.LsaKey)
		case ASExternalLSA:
			lsaEnt, exist := lsdbEnt.ASExternalLsaMap[lsdbSlice.LsaKey]
			if !exist {
				idx++
				continue
			}
			lsaMd = lsaEnt.LsaMd
			lsaEnc = encodeASExternalLsa(lsaEnt, lsdbSlice.LsaKey)
		default:
			idx++
			continue
		}
		var obj objects.Ospfv2LsdbState
		obj.LSType = lsdbSlice.LsaKey.LSType
		obj.LSId = lsdbSlice.LsaKey.LSId
		obj.AdvRouterId = lsdbSlice.LsaKey.AdvRouter
		obj.AreaId = lsdbSlice.LsdbKey.AreaId
		obj.SequenceNum = uint32(lsaMd.LSSequenceNum)
		obj.Age = lsaMd.LSAge
		obj.Checksum = lsaMd.LSChecksum
		obj.Options = lsaMd.Options
		obj.Length = lsaMd.LSLen
		obj.Advertisement = convertByteToOctetString(lsaEnc[OSPF_LSA_HEADER_SIZE:])
		retObj.List = append(retObj.List, &obj)
		lsdbSliceMap[lsdbSlice] = true
		count++
		idx++
	}
	retObj.EndIdx = idx
	if retObj.EndIdx == sliceLen {
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.More = true
		retObj.Count = sliceLen - retObj.EndIdx + 1
	}
	return &retObj, nil
}
