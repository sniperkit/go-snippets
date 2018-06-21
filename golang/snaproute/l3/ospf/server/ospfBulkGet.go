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
	"fmt"
	"l3/ospf/config"
	"net"
)

func (server *OSPFServer) GetBulkOspfAreaEntryState(idx int, cnt int) (int, int, []config.AreaState) {
	var nextIdx int
	var count int

	//server.AreaStateMutex.RLock()
	ret := server.AreaStateTimer.Stop()
	if ret == false {
		server.logger.Err("Ospf is busy refreshing the cache")
		return nextIdx, count, nil
	}
	length := len(server.AreaStateSlice)
	if idx+cnt > length {
		count = length - idx
		nextIdx = 0
	}
	result := make([]config.AreaState, count)

	for i := 0; i < count; i++ {
		key := server.AreaStateSlice[idx+i]
		result[i].AreaId = key.AreaId
		ent, exist := server.AreaStateMap[key]
		if exist {
			result[i].SpfRuns = ent.SpfRuns
			result[i].AreaBdrRtrCount = ent.AreaBdrRtrCount
			result[i].AsBdrRtrCount = ent.AsBdrRtrCount
			result[i].AreaLsaCount = ent.AreaLsaCount
			result[i].AreaLsaCksumSum = ent.AreaLsaCksumSum
			result[i].AreaNssaTranslatorState = ent.AreaNssaTranslatorState
			result[i].AreaNssaTranslatorEvents = ent.AreaNssaTranslatorEvents
		} else {
			result[i].SpfRuns = -1
			result[i].AreaBdrRtrCount = -1
			result[i].AsBdrRtrCount = -1
			result[i].AreaLsaCount = -1
			result[i].AreaLsaCksumSum = -1
			result[i].AreaNssaTranslatorState = -1
			result[i].AreaNssaTranslatorEvents = -1
		}

	}

	//server.AreaStateMutex.RUnlock()
	server.AreaStateTimer.Reset(server.RefreshDuration)
	server.logger.Info(fmt.Sprintln("length:", length, "count:", count, "nextIdx:", nextIdx, "result:", result))
	return nextIdx, count, result
}

func (server *OSPFServer) GetBulkOspfLsdbEntryState(idx int, cnt int) (int, int, []config.LsdbState) {
	var nextIdx int
	var count int

	ret := server.LsdbStateTimer.Stop()
	if ret == false {
		server.logger.Err("Ospf is busy refreshing the Lsdb cache")
		return nextIdx, count, nil
	}
	length := len(server.LsdbSlice)

	result := make([]config.LsdbState, cnt)
	var i int
	var j int
	//server.logger.Info(fmt.Sprintln("idx:", idx, "cnt:", cnt, "length of Ls DB Slice:", length))
	for i, j = 0, idx; i < cnt && j < length; j++ {
		var lsaEnc []byte
		var lsaMd LsaMetadata
		lsdbSliceEnt := server.LsdbSlice[j]
		lsdbKey := LsdbKey{
			AreaId: lsdbSliceEnt.AreaId,
		}
		lsDbEnt, exist := server.AreaLsdb[lsdbKey]
		if !exist {
			continue
		}

		lsaKey := LsaKey{
			LSType:    lsdbSliceEnt.LSType,
			LSId:      lsdbSliceEnt.LSId,
			AdvRouter: lsdbSliceEnt.AdvRtr,
		}

		if lsdbSliceEnt.LSType == RouterLSA {
			lsa, exist := lsDbEnt.RouterLsaMap[lsaKey]
			if !exist {
				continue
			}
			lsaEnc = encodeRouterLsa(lsa, lsaKey)
			lsaMd = lsa.LsaMd
		} else if lsdbSliceEnt.LSType == NetworkLSA {
			lsa, exist := lsDbEnt.NetworkLsaMap[lsaKey]
			if !exist {
				continue
			}
			lsaEnc = encodeNetworkLsa(lsa, lsaKey)
			lsaMd = lsa.LsaMd
		} else if lsdbSliceEnt.LSType == Summary3LSA {
			lsa, exist := lsDbEnt.Summary3LsaMap[lsaKey]
			if !exist {
				continue
			}
			lsaEnc = encodeSummaryLsa(lsa, lsaKey)
			lsaMd = lsa.LsaMd
		} else if lsdbSliceEnt.LSType == Summary4LSA {
			lsa, exist := lsDbEnt.Summary4LsaMap[lsaKey]
			if !exist {
				continue
			}
			lsaEnc = encodeSummaryLsa(lsa, lsaKey)
			lsaMd = lsa.LsaMd
		} else if lsdbSliceEnt.LSType == ASExternalLSA {
			lsa, exist := lsDbEnt.ASExternalLsaMap[lsaKey]
			if !exist {
				continue
			}
			lsaEnc = encodeASExternalLsa(lsa, lsaKey)
			lsaMd = lsa.LsaMd
		}

		server.logger.Info(fmt.Sprintln(lsaEnc))
		//server.logger.Info(fmt.Sprintln("lsaEnc:", lsaEnc))
		adv := convertByteToOctetString(lsaEnc[OSPF_LSA_HEADER_SIZE:])
		//server.logger.Info(fmt.Sprintln("adv:", adv))
		result[i].LsdbAreaId = config.AreaId(convertUint32ToIPv4(lsdbKey.AreaId))
		result[i].LsdbType = config.LsaType(lsaKey.LSType)
		result[i].LsdbLsid = config.IpAddress(convertUint32ToIPv4(lsaKey.LSId))
		result[i].LsdbRouterId = config.RouterId(convertUint32ToIPv4(lsaKey.AdvRouter))
		result[i].LsdbSequence = lsaMd.LSSequenceNum
		result[i].LsdbAge = int(lsaMd.LSAge)
		result[i].LsdbCheckSum = int(lsaMd.LSChecksum)
		result[i].LsdbAdvertisement = adv
		//server.logger.Info(fmt.Sprintln("Result of GetBulk:", result))
		i++
	}

	if j == length {
		nextIdx = 0
	}
	count = i

	server.LsdbStateTimer.Reset(server.RefreshDuration)
	//server.logger.Info(fmt.Sprintln("length:", length, "count:", count, "nextIdx:", nextIdx, "result:", result))
	return nextIdx, count, result
}

func (server *OSPFServer) GetBulkOspfIfEntryState(idx int, cnt int) (int, int, []config.InterfaceState) {
	var nextIdx int
	var count int

	ret := server.IntfStateTimer.Stop()
	if ret == false {
		server.logger.Err("Ospf is busy refreshing the cache")
		return nextIdx, count, nil
	}
	length := len(server.IntfKeySlice)
	if idx+cnt > length {
		count = length - idx
		nextIdx = 0
	}
	result := make([]config.InterfaceState, count)

	for i := 0; i < count; i++ {
		key := server.IntfKeySlice[idx+i]
		result[i].IfIpAddress = key.IPAddr
		result[i].AddressLessIf = key.IntfIdx
		if server.IntfKeyToSliceIdxMap[key] == true {
			//if exist {
			ent, _ := server.IntfConfMap[key]
			result[i].IfState = ent.IfFSMState
			ip := net.IPv4(ent.IfDRIp[0], ent.IfDRIp[1], ent.IfDRIp[2], ent.IfDRIp[3])
			result[i].IfDesignatedRouter = config.IpAddress(ip.String())
			ip = net.IPv4(ent.IfBDRIp[0], ent.IfBDRIp[1], ent.IfBDRIp[2], ent.IfBDRIp[3])
			result[i].IfBackupDesignatedRouter = config.IpAddress(ip.String())
			result[i].IfEvents = ent.IfEvents
			result[i].IfLsaCount = ent.IfLsaCount
			result[i].IfLsaCksumSum = ent.IfLsaCksumSum
			result[i].IfDesignatedRouterId = config.RouterId(convertUint32ToIPv4(ent.IfDRtrId))
			result[i].IfBackupDesignatedRouterId = config.RouterId(convertUint32ToIPv4(ent.IfBDRtrId))
		} else {
			result[i].IfState = 0
			result[i].IfDesignatedRouter = "0.0.0.0"
			result[i].IfBackupDesignatedRouter = "0.0.0.0"
			result[i].IfEvents = 0
			result[i].IfLsaCount = 0
			result[i].IfLsaCksumSum = 0
			result[i].IfDesignatedRouterId = "0.0.0.0"
			result[i].IfBackupDesignatedRouterId = "0.0.0.0"
		}
	}

	server.IntfStateTimer.Reset(server.RefreshDuration)
	server.logger.Info(fmt.Sprintln("length:", length, "count:", count, "nextIdx:", nextIdx, "result:", result))
	return nextIdx, count, result
}

func (server *OSPFServer) GetOspfGlobalState() *config.GlobalState {
	result := new(config.GlobalState)
	ent := server.ospfGlobalConf

	ip := net.IPv4(ent.RouterId[0], ent.RouterId[1], ent.RouterId[2], ent.RouterId[3])
	result.RouterId = config.RouterId(ip.String())
	result.VersionNumber = int32(ent.Version)
	result.AreaBdrRtrStatus = ent.AreaBdrRtrStatus
	result.ExternLsaCount = ent.ExternLsaCount
	result.ExternLsaChecksum = ent.ExternLsaChecksum
	result.OriginateNewLsas = ent.OriginateNewLsas
	result.RxNewLsas = ent.RxNewLsas
	result.OpaqueLsaSupport = ent.OpaqueLsaSupport
	result.RestartStatus = ent.RestartStatus
	result.RestartAge = ent.RestartAge
	result.RestartExitReason = ent.RestartExitReason
	result.AsLsaCount = ent.AsLsaCount
	result.AsLsaCksumSum = ent.AsLsaCksumSum
	result.StubRouterSupport = ent.StubRouterSupport
	result.DiscontinuityTime = ent.DiscontinuityTime
	server.logger.Info(fmt.Sprintln("Global State:", result))
	return result
}

func (server *OSPFServer) GetBulkOspfNbrEntryState(idx int, cnt int) (int, int, []config.NeighborState) {
	server.logger.Info(fmt.Sprintln("Getbulk: nbr states called."))
	var nextIdx int
	var count int
	if len(server.neighborBulkSlice) < 1 {
		return 0, 0, nil
	}
	server.neighborSliceStartCh <- false
	/*	if ret == false {
		server.logger.Err("Ospf is busy refreshing the cache")
		return nextIdx, count, nil
	} */
	NbrStateLen := len(config.NbrStateList)
	length := len(server.neighborBulkSlice)
	if idx+cnt > length {
		count = length - idx
		nextIdx = 0
	}
	result := make([]config.NeighborState, count)

	for i := 0; i < count; i++ {
		key := server.neighborBulkSlice[idx+i]
		server.logger.Info(fmt.Sprintln("Key ", key))
		/* get map entries.
		 */

		if ent, ok := server.NeighborConfigMap[key]; ok {
			result[i].NbrIpAddress = config.IpAddress(ent.OspfNbrIPAddr.String())
			result[i].NbrAddressLessIndex = int(ent.intfConfKey.IntfIdx)
			result[i].NbrRtrId = convertUint32ToIPv4(ent.OspfNbrRtrId)
			result[i].NbrOptions = ent.OspfNbrOptions
			result[i].NbrPriority = uint8(ent.OspfRtrPrio)
			result[i].NbrState = config.NbrStateList[int(ent.OspfNbrState)%NbrStateLen]
			result[i].NbrEvents = int(ent.nbrEvent)
			result[i].NbrLsRetransQLen = 0
			result[i].NbmaNbrPermanence = 0
			result[i].NbrHelloSuppressed = false
			result[i].NbrRestartHelperStatus = 0
			result[i].NbrRestartHelperAge = 0
			result[i].NbrRestartHelperExitReason = 0
		}

	}

	server.neighborSliceStartCh <- true
	server.logger.Info(fmt.Sprintln("length:", length, "count:", count, "nextIdx:", nextIdx, "result:", result))
	return nextIdx, count, result
}
