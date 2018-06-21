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

func (server *OSPFV2Server) processRecvdSelfSummary4LSA(msg RecvdSelfLsaMsg) error {
	lsa, ok := msg.LsaData.(SummaryLsa)
	if !ok {
		server.logger.Err("Unable to assert given network lsa")
		return nil
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	lsaEnt, exist := lsdbEnt.Summary4LsaMap[msg.LsaKey]
	if !exist {
		server.logger.Err("No such Summary 4 LSA exist", msg.LsaKey)
		// Mark the recvd lsa as MaxAge and flood
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
		server.logger.Err("No such self originated summary LSA exist")
		// Mark the recvd lsa as MaxAge and flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	if lsaEnt.LsaMd.LSSequenceNum < lsa.LsaMd.LSSequenceNum {
		lsaEnt.LsaMd.LSSequenceNum = lsa.LsaMd.LSSequenceNum + 1
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnc := encodeSummaryLsa(lsaEnt, msg.LsaKey)
		checksumOffset := uint16(14)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.Summary4LsaMap[msg.LsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
		// Flood new Self Summary 4 LSA (areaId, lsaKey, lsaEnt)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
		return nil
	} else {
		//Flood existing Self Summary 4 LSA (areaId, lsaKey, lsaEnt)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
	}
	return nil
}

func (server *OSPFV2Server) processRecvdSelfSummary3LSA(msg RecvdSelfLsaMsg) error {
	lsa, ok := msg.LsaData.(SummaryLsa)
	if !ok {
		server.logger.Err("Unable to assert given network lsa")
		return nil
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	lsaEnt, exist := lsdbEnt.Summary3LsaMap[msg.LsaKey]
	if !exist {
		server.logger.Err("No such Summary 3 LSA exist", msg.LsaKey)
		// Mark the recvd lsa as MaxAge and flood
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
		server.logger.Err("No such self originated summary LSA exist")
		// Mark the recvd lsa as MaxAge and flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	if lsaEnt.LsaMd.LSSequenceNum < lsa.LsaMd.LSSequenceNum {
		lsaEnt.LsaMd.LSSequenceNum = lsa.LsaMd.LSSequenceNum + 1
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnc := encodeSummaryLsa(lsaEnt, msg.LsaKey)
		checksumOffset := uint16(14)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.Summary3LsaMap[msg.LsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
		// Flood new Self Summary 3 LSA (areaId, lsaKey, lsaEnt)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
		return nil
	} else {
		// Flood existing Self Summary 3 LSA (areaId, lsaKey, lsaEnt)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
	}
	return nil
}

func (server *OSPFV2Server) processRecvdSelfSummaryLSA(msg RecvdSelfLsaMsg) error {
	if msg.LsaKey.LSType == Summary3LSA {
		server.processRecvdSelfSummary3LSA(msg)
	} else {
		server.processRecvdSelfSummary4LSA(msg)
	}
	return nil
}

func (server *OSPFV2Server) processRecvdSummaryLSA(msg RecvdLsaMsg) error {
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	if msg.MsgType == LSA_ADD {
		lsa, ok := msg.LsaData.(SummaryLsa)
		if !ok {
			server.logger.Err("Unable to assert given router lsa")
			return nil
		}
		if msg.LsaKey.LSType == Summary3LSA {
			_, exist = lsdbEnt.Summary3LsaMap[msg.LsaKey]
			lsdbEnt.Summary3LsaMap[msg.LsaKey] = lsa
		} else {
			_, exist = lsdbEnt.Summary3LsaMap[msg.LsaKey]
			lsdbEnt.Summary4LsaMap[msg.LsaKey] = lsa
		}
		if !exist {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: msg.LsdbKey,
				LsaKey:  msg.LsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	} else if msg.MsgType == LSA_DEL {
		if msg.LsaKey.LSType == Summary3LSA {
			delete(lsdbEnt.Summary3LsaMap, msg.LsaKey)
		} else {
			delete(lsdbEnt.Summary4LsaMap, msg.LsaKey)
		}
	}
	server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
	return nil
}

func (server *OSPFV2Server) compareSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) bool {
	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsdbEnt.Summary3LsaMap[lsaKey]
	} else {
		sLsa, _ = lsdbEnt.Summary4LsaMap[lsaKey]
	}
	if sLsa.Netmask != lsaEnt.Netmask ||
		sLsa.Metric != lsaEnt.Metric {
		return false
	}
	return true
}

func (server *OSPFV2Server) updateSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) {
	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsdbEnt.Summary3LsaMap[lsaKey]
	} else {
		sLsa, _ = lsdbEnt.Summary4LsaMap[lsaKey]
	}
	sLsa.Metric = lsaEnt.Metric
	sLsa.Netmask = lsaEnt.Netmask
	sLsa.LsaMd.LSAge = 0
	sLsa.LsaMd.LSChecksum = 0
	sLsa.LsaMd.LSLen = lsaEnt.LsaMd.LSLen
	sLsa.LsaMd.LSSequenceNum = sLsa.LsaMd.LSSequenceNum + 1
	sLsa.LsaMd.Options = EOption
	sLsaEnc := encodeSummaryLsa(sLsa, lsaKey)
	checksumOffset := uint16(14)
	sLsa.LsaMd.LSChecksum = computeFletcherChecksum(sLsaEnc[2:], checksumOffset)
	if lsaKey.LSType == Summary3LSA {
		lsdbEnt.Summary3LsaMap[lsaKey] = sLsa
	} else {
		lsdbEnt.Summary4LsaMap[lsaKey] = sLsa
	}
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	// Flood Updated Summary Lsa
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, sLsa)
}

func (server *OSPFV2Server) insertNewSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt SummaryLsa) {
	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsdbEnt.Summary3LsaMap[lsaKey]
	} else {
		sLsa, _ = lsdbEnt.Summary4LsaMap[lsaKey]
	}
	sLsa.Metric = lsaEnt.Metric
	sLsa.Netmask = lsaEnt.Netmask
	sLsa.LsaMd.LSAge = 0
	sLsa.LsaMd.LSChecksum = 0
	sLsa.LsaMd.LSLen = lsaEnt.LsaMd.LSLen
	sLsa.LsaMd.LSSequenceNum = int(InitialSequenceNum)
	sLsa.LsaMd.Options = EOption
	sLsaEnc := encodeSummaryLsa(sLsa, lsaKey)
	checksumOffset := uint16(14)
	sLsa.LsaMd.LSChecksum = computeFletcherChecksum(sLsaEnc[2:], checksumOffset)
	if lsaKey.LSType == Summary3LSA {
		lsdbEnt.Summary3LsaMap[lsaKey] = sLsa
	} else {
		lsdbEnt.Summary4LsaMap[lsaKey] = sLsa
	}
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	selfOrigLsaEnt[lsaKey] = true
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	//Flood New Summary Lsa
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, sLsa)
	lsdbSlice := LsdbSliceStruct{
		LsdbKey: lsdbKey,
		LsaKey:  lsaKey,
	}
	server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
}

func (server *OSPFV2Server) flushSummaryLsa(lsdbKey LsdbKey, lsaKey LsaKey) {
	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	var sLsa SummaryLsa
	if lsaKey.LSType == Summary3LSA {
		sLsa, _ = lsdbEnt.Summary3LsaMap[lsaKey]
	} else {
		sLsa, _ = lsdbEnt.Summary4LsaMap[lsaKey]
	}
	sLsa.LsaMd.LSAge = MAX_AGE
	// Flood LSA to flush
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, sLsa)
	delete(selfOrigLsaEnt, lsaKey)
	if lsaKey.LSType == Summary3LSA {
		delete(lsdbEnt.Summary3LsaMap, lsaKey)
	} else {
		delete(lsdbEnt.Summary4LsaMap, lsaKey)
	}
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
}

func (server *OSPFV2Server) installSummaryLsa() {
	server.logger.Info("Installing Summary LSA")
	for lsdbKey, sLsa := range server.SummaryLsDb {
		selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
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
				}
			} else {
				server.insertNewSummaryLsa(lsdbKey, sKey, sEnt)
			}
		}
		for sKey, val := range oldSelfOrigSummaryLsa {
			if val == true {
				server.flushSummaryLsa(lsdbKey, sKey)
			}
		}
		oldSelfOrigSummaryLsa = nil
	}
}
