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
	"fmt"
	"l3/ospfv2/objects"
)

func (server *OSPFV2Server) flushNetworkLSA(intfKey IntfConfKey) error {
	intfConfEnt, err := server.GetIntfConfForGivenIntfKey(intfKey)
	if err != nil {
		return err
	}
	if intfConfEnt.Type != objects.INTF_TYPE_BROADCAST {
		return errors.New("Network LSA doesnot exist for Non Broadcast Network")
	}
	lsdbKey := LsdbKey{
		AreaId: intfConfEnt.AreaId,
	}
	lsaKey := LsaKey{
		LSType:    NetworkLSA,
		LSId:      intfConfEnt.IpAddr,
		AdvRouter: server.globalData.RouterId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		return errors.New(fmt.Sprintln("Error: No lsDbEnt found for", lsdbKey))
	}
	lsaEnt, exist := lsdbEnt.NetworkLsaMap[lsaKey]
	if !exist {
		return errors.New(fmt.Sprintln("Error: No Network found for", lsaKey))
	}
	lsaEnt.LsaMd.LSAge = MAX_AGE

	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	delete(lsdbEnt.NetworkLsaMap, lsaKey)
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	delete(selfOrigLsaEnt, lsaKey)
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	//Flood new Network LSA (areaId, lsaEnt, lsaKey)
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, lsaEnt)
	return nil
}

func (server *OSPFV2Server) generateNetworkLSA(intfKey IntfConfKey, nbrList []uint32) error {
	var lsaKey LsaKey
	intfConfEnt, err := server.GetIntfConfForGivenIntfKey(intfKey)
	if err != nil {
		return err
	}
	if intfConfEnt.Type != objects.INTF_TYPE_BROADCAST {
		return errors.New("Network LSA won't be generated for Non Broadcast Network")
	}
	lsdbKey := LsdbKey{
		AreaId: intfConfEnt.AreaId,
	}
	lsaKey = LsaKey{
		LSType:    NetworkLSA,
		LSId:      intfConfEnt.IpAddr,
		AdvRouter: server.globalData.RouterId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		return errors.New(fmt.Sprintln("Error: No lsDbEnt found for", lsdbKey))
	}
	lsaEnt, exist := lsdbEnt.NetworkLsaMap[lsaKey]
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	if len(nbrList) == 0 {
		return server.flushNetworkLSA(intfKey)
	}
	lsaEnt.AttachedRtr = nil
	lsaEnt.AttachedRtr = append(lsaEnt.AttachedRtr, server.globalData.RouterId)
	for _, nbr := range nbrList {
		lsaEnt.AttachedRtr = append(lsaEnt.AttachedRtr, nbr)
	}
	lsaEnt.Netmask = intfConfEnt.Netmask
	lsaEnt.LsaMd.LSAge = 0
	lsaEnt.LsaMd.LSChecksum = 0
	lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (4 * len(lsaEnt.AttachedRtr)))
	if !exist {
		lsaEnt.LsaMd.LSSequenceNum = int(InitialSequenceNum)
	} else {
		lsaEnt.LsaMd.LSSequenceNum += 1
	}
	lsaEnt.LsaMd.Options = EOption
	lsaEnc := encodeNetworkLsa(lsaEnt, lsaKey)
	checksumOffset := uint16(14)
	lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
	lsdbEnt.NetworkLsaMap[lsaKey] = lsaEnt
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	selfOrigLsaEnt[lsaKey] = true
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	//Flood new Network LSA (areaId, lsaEnt, lsaKey)
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, lsaEnt)
	if !exist {
		lsdbSlice := LsdbSliceStruct{
			LsdbKey: lsdbKey,
			LsaKey:  lsaKey,
		}
		server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
	}
	return nil
}

func (server *OSPFV2Server) processUpdateSelfNetworkLSA(msg UpdateSelfNetworkLSAMsg) error {
	if msg.Op == FLUSH {
		return server.flushNetworkLSA(msg.IntfKey)
	} else if msg.Op == GENERATE {
		return server.generateNetworkLSA(msg.IntfKey, msg.NbrList)
	}
	return errors.New("Invalid Op value")
}

func (server *OSPFV2Server) processRecvdSelfNetworkLSA(msg RecvdSelfLsaMsg) error {
	lsa, ok := msg.LsaData.(NetworkLsa)
	if !ok {
		server.logger.Err("Unable to assert given Network lsa")
		return nil
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	lsaEnt, exist := lsdbEnt.NetworkLsaMap[msg.LsaKey]
	if !exist {
		server.logger.Err("No such Network LSA exist", msg.LsaKey)
		// Mark as Max Age and flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	selfOrigLsaEnt, exist := server.LsdbData.AreaSelfOrigLsa[msg.LsdbKey]
	if !exist {
		server.logger.Err("No self originated LSA exist", msg.LsdbKey)
		return nil
	}
	_, exist = selfOrigLsaEnt[msg.LsaKey]
	if !exist {
		server.logger.Err("No such self originated Network LSA exist", msg.LsaKey)
		// Mark as Max Age and flood
		lsa.LsaMd.LSAge = MAX_AGE
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsa)
		return nil
	}
	if lsaEnt.LsaMd.LSSequenceNum < lsa.LsaMd.LSSequenceNum {
		checksumOffset := uint16(14)
		lsaEnt.LsaMd.LSSequenceNum = lsa.LsaMd.LSSequenceNum + 1
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnc := encodeNetworkLsa(lsaEnt, msg.LsaKey)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.NetworkLsaMap[msg.LsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
		// Flood new Self Network LSA (lsaEnt, msg.LsaKey)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
		return nil
	} else {
		//Flood existing Self Network LSA (lsaEnt, msg.LsaKey)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
	}

	return nil
}

func (server *OSPFV2Server) processRecvdNetworkLSA(msg RecvdLsaMsg) error {
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	if msg.MsgType == LSA_ADD {
		lsa, ok := msg.LsaData.(NetworkLsa)
		if !ok {
			server.logger.Err("Unable to assert given Network lsa")
			return nil
		}
		_, exist = lsdbEnt.NetworkLsaMap[msg.LsaKey]
		lsdbEnt.NetworkLsaMap[msg.LsaKey] = lsa
		if !exist {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: msg.LsdbKey,
				LsaKey:  msg.LsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	} else if msg.MsgType == LSA_DEL {
		delete(lsdbEnt.NetworkLsaMap, msg.LsaKey)
	}
	server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
	return nil
}

func (server *OSPFV2Server) reGenerateNetworkLSA(lsaKey LsaKey, nLsa NetworkLsa, lsdbKey LsdbKey) bool {
	server.logger.Info("Regenerating Network LSA")
	flag := false
	var intfConfKey IntfConfKey
	var intfEnt IntfConf
	for intfKey, intfConfEnt := range server.IntfConfMap {
		if intfConfEnt.FSMState != objects.INTF_FSM_STATE_DR {
			continue
		}
		if intfConfEnt.Type != objects.INTF_TYPE_BROADCAST {
			continue
		}
		if nLsa.Netmask == intfConfEnt.Netmask &&
			lsaKey.LSId == intfConfEnt.IpAddr &&
			lsaKey.AdvRouter == server.globalData.RouterId {
			flag = true
			intfConfKey = intfKey
			intfEnt = intfConfEnt
			break
		}
	}

	if flag == false {
		return false
	}
	nbrList, err := server.getFullNbrList(intfConfKey)
	if err != nil {
		return false
	}
	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	lsaEnt, _ := lsdbEnt.NetworkLsaMap[lsaKey]
	lsaEnt.AttachedRtr = nil
	lsaEnt.AttachedRtr = append(lsaEnt.AttachedRtr, server.globalData.RouterId)
	for _, nbr := range nbrList {
		lsaEnt.AttachedRtr = append(lsaEnt.AttachedRtr, nbr)
	}
	lsaEnt.Netmask = intfEnt.Netmask
	lsaEnt.LsaMd.LSAge = 0
	lsaEnt.LsaMd.LSChecksum = 0
	lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (4 * len(lsaEnt.AttachedRtr)))
	lsaEnt.LsaMd.LSSequenceNum = lsaEnt.LsaMd.LSSequenceNum + 1
	lsaEnt.LsaMd.Options = EOption
	lsaEnc := encodeNetworkLsa(lsaEnt, lsaKey)
	checksumOffset := uint16(14)
	lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
	lsdbEnt.NetworkLsaMap[lsaKey] = lsaEnt
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	//Flood new Network LSA (areaId, lsaEnt, lsaKey)
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, lsaEnt)
	return true
}
