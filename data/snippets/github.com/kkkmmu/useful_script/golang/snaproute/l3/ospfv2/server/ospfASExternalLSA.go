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
//"errors"
//"fmt"
//"l3/ospfv2/objects"
)

func (server *OSPFV2Server) processRecvdSelfASExternalLSA(msg RecvdSelfLsaMsg) error {
	lsa, ok := msg.LsaData.(ASExternalLsa)
	if !ok {
		server.logger.Err("Unable to assert given ASExternal lsa")
		return nil
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	lsaEnt, exist := lsdbEnt.ASExternalLsaMap[msg.LsaKey]
	if !exist {
		server.logger.Err("No such ASExternal LSA exist", msg.LsaKey)
		// Mark the recvd LSA as MAX_AGE and Flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	selfOrigLsaEnt, exist := server.LsdbData.AreaSelfOrigLsa[msg.LsdbKey]
	if !exist {
		server.logger.Err("No self originated LSA exist")
		return nil
	}
	_, exist = selfOrigLsaEnt[msg.LsaKey]
	if !exist {
		server.logger.Err("No such self originated ASExternal LSA exist", msg.LsaKey)
		// Mark the recvd LSA as MAX_AGE and Flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	if lsaEnt.LsaMd.LSSequenceNum < lsa.LsaMd.LSSequenceNum {
		lsaEnt.LsaMd.LSSequenceNum = lsa.LsaMd.LSSequenceNum + 1
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnc := encodeASExternalLsa(lsaEnt, msg.LsaKey)
		checksumOffset := uint16(14)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.ASExternalLsaMap[msg.LsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
		// Flood new Self ASExternal LSA
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
		return nil
	} else {
		// Flood existing Self ASExternal LSA
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
	}

	return nil
}

func (server *OSPFV2Server) processRecvdASExternalLSA(msg RecvdLsaMsg) error {
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	if msg.MsgType == LSA_ADD {
		lsa, ok := msg.LsaData.(ASExternalLsa)
		if !ok {
			server.logger.Err("Unable to assert given ASExternal lsa")
			return nil
		}
		_, exist = lsdbEnt.ASExternalLsaMap[msg.LsaKey]
		lsdbEnt.ASExternalLsaMap[msg.LsaKey] = lsa
		if !exist {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: msg.LsdbKey,
				LsaKey:  msg.LsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	} else if msg.MsgType == LSA_DEL {
		delete(lsdbEnt.ASExternalLsaMap, msg.LsaKey)
	}
	server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
	return nil
}

func (server *OSPFV2Server) generateASExternalLSA(routeInfo RouteInfo) {
	if server.globalData.ASBdrRtrStatus == false {
		return
	}

	var areaIdList []uint32
	for areaId, areaEnt := range server.AreaConfMap {
		if areaEnt.AdminState == false ||
			areaEnt.ImportASExtern == false {
			continue
		}
		areaIdList = append(areaIdList, areaId)
	}
	if len(areaIdList) == 0 {
		return
	}
	LSType := ASExternalLSA
	LSId := routeInfo.NwAddr & routeInfo.Netmask
	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}
	checksumOffset := uint16(14)
	var asLsaEnt ASExternalLsa
	asLsaEnt.LsaMd.LSAge = 0
	asLsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 16)
	asLsaEnt.BitE = true
	asLsaEnt.ExtRouteTag = 0
	asLsaEnt.FwdAddr = 0
	asLsaEnt.Metric = routeInfo.Metric
	asLsaEnt.Netmask = routeInfo.Netmask
	asLsaEnt.LsaMd.Options = EOption
	for _, areaId := range areaIdList {
		lsdbKey := LsdbKey{
			AreaId: areaId,
		}
		lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
		if !exist {
			server.logger.Err("No Lsdb Exist for:", lsdbKey)
			continue
		}
		lsaEnt, exist := lsdbEnt.ASExternalLsaMap[lsaKey]
		lsaEnt = asLsaEnt
		lsaEnt.LsaMd.LSChecksum = 0
		if exist {
			lsaEnt.LsaMd.LSSequenceNum = int(InitialSequenceNum)
		} else {
			lsaEnt.LsaMd.LSSequenceNum = lsaEnt.LsaMd.LSSequenceNum + 1
		}
		lsaEnc := encodeASExternalLsa(lsaEnt, lsaKey)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.ASExternalLsaMap[lsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
		selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
		selfOrigLsaEnt[lsaKey] = true
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.CreateAndSendMsgFromLsdbToFloodLsa(areaId, lsaKey, lsaEnt)
		if !exist {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: lsdbKey,
				LsaKey:  lsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	}
}

func (server *OSPFV2Server) flushASExternalLSA(routeInfo RouteInfo) {
	if server.globalData.ASBdrRtrStatus == false {
		return
	}
	var areaIdList []uint32
	for areaId, areaEnt := range server.AreaConfMap {
		if areaEnt.AdminState == false ||
			areaEnt.ImportASExtern == false {
			continue
		}
		areaIdList = append(areaIdList, areaId)
	}
	if len(areaIdList) == 0 {
		return
	}
	LSType := ASExternalLSA
	LSId := routeInfo.NwAddr & routeInfo.Netmask
	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}
	for _, areaId := range areaIdList {
		lsdbKey := LsdbKey{
			AreaId: areaId,
		}
		lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
		if !exist {
			server.logger.Err("No Lsdb Exist for:", lsdbKey)
			continue
		}
		lsaEnt, exist := lsdbEnt.ASExternalLsaMap[lsaKey]
		if !exist {
			server.logger.Err("No LSA exist:", lsaKey)
			continue
		}
		selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
		lsaEnt.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(areaId, lsaKey, lsaEnt)
		delete(selfOrigLsaEnt, lsaKey)
		delete(lsdbEnt.ASExternalLsaMap, lsaKey)
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	}
}

func (server *OSPFV2Server) GenerateAllASExternalLSA(areaId uint32) {
	if server.globalData.ASBdrRtrStatus == false {
		return
	}
	areaEnt, err := server.GetAreaConfForGivenArea(areaId)
	if err != nil {
		server.logger.Err("No such area exist")
		return
	}
	if areaEnt.AdminState == false ||
		areaEnt.ImportASExtern == false {
		return
	}
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("No Lsdb Exist for Area:", areaId)
		return
	}
	for route, _ := range server.LsdbData.ExtRouteInfoMap {
		LSType := ASExternalLSA
		LSId := route.NwAddr & route.Netmask
		AdvRouter := server.globalData.RouterId
		lsaKey := LsaKey{
			LSType:    LSType,
			LSId:      LSId,
			AdvRouter: AdvRouter,
		}
		lsaEnt, _ := lsdbEnt.ASExternalLsaMap[lsaKey]
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 16)
		lsaEnt.LsaMd.LSSequenceNum = int(InitialSequenceNum)
		lsaEnt.LsaMd.Options = EOption
		lsaEnt.BitE = true
		lsaEnt.ExtRouteTag = 0
		lsaEnt.FwdAddr = 0
		lsaEnt.Metric = route.Metric
		lsaEnt.Netmask = route.Netmask
		checksumOffset := uint16(14)
		lsaEnc := encodeASExternalLsa(lsaEnt, lsaKey)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.ASExternalLsaMap[lsaKey] = lsaEnt
		selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
		selfOrigLsaEnt[lsaKey] = true
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
		lsdbSlice := LsdbSliceStruct{
			LsdbKey: lsdbKey,
			LsaKey:  lsaKey,
		}
		server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
	}
}

func (server *OSPFV2Server) reGenerateASExternalLSAForGivenArea(routeInfo RouteInfo, areaId uint32) {
	if server.globalData.ASBdrRtrStatus == false {
		return
	}

	areaEnt, err := server.GetAreaConfForGivenArea(areaId)
	if err != nil {
		server.logger.Err("Error: Unable to find the areaConf for:", areaId)
		return
	}
	if areaEnt.AdminState == false ||
		areaEnt.ImportASExtern == false {
		return
	}
	LSType := ASExternalLSA
	LSId := routeInfo.NwAddr & routeInfo.Netmask
	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    LSType,
		LSId:      LSId,
		AdvRouter: AdvRouter,
	}
	checksumOffset := uint16(14)
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("No Lsdb Exist for:", lsdbKey)
		return
	}
	lsaEnt, exist := lsdbEnt.ASExternalLsaMap[lsaKey]
	if !exist {
		server.logger.Err("Error: Unable to find lsa hence cannot regenerate", lsaKey)
		return
	}
	lsaEnt.LsaMd.LSAge = 0
	lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 16)
	lsaEnt.BitE = true
	lsaEnt.ExtRouteTag = 0
	lsaEnt.FwdAddr = 0
	lsaEnt.Metric = routeInfo.Metric
	lsaEnt.Netmask = routeInfo.Netmask
	lsaEnt.LsaMd.Options = EOption
	lsaEnt.LsaMd.LSChecksum = 0
	lsaEnt.LsaMd.LSSequenceNum = lsaEnt.LsaMd.LSSequenceNum + 1
	lsaEnc := encodeASExternalLsa(lsaEnt, lsaKey)
	lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
	lsdbEnt.ASExternalLsaMap[lsaKey] = lsaEnt
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	selfOrigLsaEnt[lsaKey] = true
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	server.CreateAndSendMsgFromLsdbToFloodLsa(areaId, lsaKey, lsaEnt)
}
