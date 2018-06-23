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

func (server *OSPFV2Server) constructStubLinkP2P(intfConf IntfConf) LinkDetail {
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

	linkDetail.LinkId = intfConf.IpAddr & intfConf.Netmask
	linkDetail.LinkData = intfConf.Netmask
	linkDetail.LinkType = STUB_LINK
	linkDetail.NumOfTOS = 0
	linkDetail.LinkMetric = uint16(intfConf.Cost)

	return linkDetail
}

func (server *OSPFV2Server) GetLinkDetails(areaId uint32, areaEnt AreaConf) []LinkDetail {
	var linkDetails []LinkDetail = nil
	for intfKey, _ := range areaEnt.IntfMap {
		intfConf, err := server.GetIntfConfForGivenIntfKey(intfKey)
		if err != nil {
			continue
		}
		if intfConf.FSMState < objects.INTF_FSM_STATE_WAITING {
			server.logger.Debug("Router LSA: Skipping interface", intfKey)
			continue
		}
		var linkDetail LinkDetail
		if intfConf.FSMState == objects.INTF_FSM_STATE_LOOPBACK {
			if intfKey.IpAddr != 0 { //Un numbered Interface
				linkDetail.LinkType = STUB_LINK
				linkDetail.LinkData = 0xffffffff
				linkDetail.LinkId = intfConf.IpAddr
				linkDetail.LinkMetric = uint16(0)
				linkDetail.NumOfTOS = 0
			}
		} else {
			switch intfConf.Type {
			case objects.INTF_TYPE_BROADCAST:
				if len(intfConf.NbrMap) == 0 { //Stub Network
					server.logger.Debug("Stub Network")
					linkDetail.LinkType = STUB_LINK
					linkDetail.LinkData = intfConf.Netmask
					linkDetail.LinkId = intfConf.IpAddr & intfConf.Netmask
				} else { //Transit Link
					server.logger.Debug("Transit Network")
					linkDetail.LinkType = TRANSIT_LINK
					linkDetail.LinkData = intfConf.IpAddr
					linkDetail.LinkId = intfConf.DRIpAddr
				}
				linkDetail.LinkMetric = uint16(intfConf.Cost)
				linkDetail.NumOfTOS = 0
			case objects.INTF_TYPE_POINT2POINT:
				server.logger.Debug("P2P Network")
				stubLink := server.constructStubLinkP2P(intfConf)
				linkDetails = append(linkDetails, stubLink)
				if len(intfConf.NbrMap) == 0 {
					continue
				}
				linkDetail.LinkType = P2P_LINK
				for _, nbr := range intfConf.NbrMap {
					linkDetail.LinkId = nbr.RtrId
				}
				if intfKey.IpAddr == 0 { //Un-numbered P2P
					linkDetail.LinkData = intfKey.IntfIdx
				} else { // Numbered P2P
					linkDetail.LinkData = intfKey.IpAddr
				}
				linkDetail.NumOfTOS = 0
				linkDetail.LinkMetric = uint16(intfConf.Cost)
			}
		}
		linkDetails = append(linkDetails, linkDetail)
	}
	return linkDetails
}

func (server *OSPFV2Server) GenerateRouterLSA(msg GenerateRouterLSAMsg) error {
	var lsaKey LsaKey
	areaEnt, err := server.GetAreaConfForGivenArea(msg.AreaId)
	if err != nil {
		return err
	}
	var linkDetails []LinkDetail = nil
	linkDetails = append(linkDetails, server.GetLinkDetails(msg.AreaId, areaEnt)...)
	numOfLinks := len(linkDetails)
	BitE := false
	if server.globalData.ASBdrRtrStatus == true {
		BitE = true
	}
	BitB := false
	if server.globalData.AreaBdrRtrStatus == true {
		BitB = true
	}
	lsaKey = LsaKey{
		LSType:    RouterLSA,
		LSId:      server.globalData.RouterId,
		AdvRouter: server.globalData.RouterId,
	}
	lsdbKey := LsdbKey{
		AreaId: msg.AreaId,
	}

	lsdbEnt, lsdbExist := server.LsdbData.AreaLsdb[lsdbKey]
	if !lsdbExist {
		return errors.New(fmt.Sprintln("Area doesnot exist. No router LSA will be generated", lsdbKey))
	}
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	if numOfLinks == 0 {
		delete(lsdbEnt.RouterLsaMap, lsaKey)
		delete(selfOrigLsaEnt, lsaKey)
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
		return nil
	}
	lsaEnt, exist := lsdbEnt.RouterLsaMap[lsaKey]
	lsaEnt.LsaMd.LSAge = 0
	lsaEnt.LsaMd.LSChecksum = 0
	lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (12 * numOfLinks))
	lsaEnt.LsaMd.Options = EOption
	if !exist {
		lsaEnt.LsaMd.LSSequenceNum = int(InitialSequenceNum)
	} else {
		lsaEnt.LsaMd.LSSequenceNum = lsaEnt.LsaMd.LSSequenceNum + 1
	}
	lsaEnt.BitB = BitB
	lsaEnt.BitE = BitE
	lsaEnt.BitV = false
	lsaEnt.NumOfLinks = uint16(numOfLinks)
	lsaEnt.LinkDetails = nil
	lsaEnt.LinkDetails = append(lsaEnt.LinkDetails, linkDetails...)
	lsaEnc := encodeRouterLsa(lsaEnt, lsaKey)
	checksumOffset := uint16(14)
	lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
	lsdbEnt.RouterLsaMap[lsaKey] = lsaEnt
	selfOrigLsaEnt[lsaKey] = true
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
	//Flood new Self Router LSA (areaId, lsaEnt, lsaKey)
	server.logger.Info("Calling CreateAndSendMsgFromLsdbToFloodLsa():", lsdbKey.AreaId, lsaKey, lsaEnt)
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

func (server *OSPFV2Server) reGenerateRouterLSA(msg GenerateRouterLSAMsg) {
	var lsaKey LsaKey
	areaEnt, err := server.GetAreaConfForGivenArea(msg.AreaId)
	if err != nil {
		server.logger.Err("Unable to find the Area", msg.AreaId)
		return
	}
	var linkDetails []LinkDetail = nil
	linkDetails = append(linkDetails, server.GetLinkDetails(msg.AreaId, areaEnt)...)
	numOfLinks := len(linkDetails)
	BitE := false
	if server.globalData.ASBdrRtrStatus == true {
		BitE = true
	}
	BitB := false
	if server.globalData.AreaBdrRtrStatus == true {
		BitB = true
	}
	lsaKey = LsaKey{
		LSType:    RouterLSA,
		LSId:      server.globalData.RouterId,
		AdvRouter: server.globalData.RouterId,
	}
	lsdbKey := LsdbKey{
		AreaId: msg.AreaId,
	}

	lsdbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	selfOrigLsaEnt, _ := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
	lsaEnt, _ := lsdbEnt.RouterLsaMap[lsaKey]
	if numOfLinks == 0 {
		delete(lsdbEnt.RouterLsaMap, lsaKey)
		delete(selfOrigLsaEnt, lsaKey)
		server.LsdbData.AreaSelfOrigLsa[lsdbKey] = selfOrigLsaEnt
		server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
		return
	}
	lsaEnt.LsaMd.LSAge = 0
	lsaEnt.LsaMd.LSChecksum = 0
	lsaEnt.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 4 + (12 * numOfLinks))
	lsaEnt.LsaMd.Options = EOption
	lsaEnt.LsaMd.LSSequenceNum = lsaEnt.LsaMd.LSSequenceNum + 1
	lsaEnt.BitB = BitB
	lsaEnt.BitE = BitE
	lsaEnt.BitV = false
	lsaEnt.NumOfLinks = uint16(numOfLinks)
	lsaEnt.LinkDetails = nil
	lsaEnt.LinkDetails = append(lsaEnt.LinkDetails, linkDetails...)
	lsaEnc := encodeRouterLsa(lsaEnt, lsaKey)
	checksumOffset := uint16(14)
	lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
	lsdbEnt.RouterLsaMap[lsaKey] = lsaEnt
	server.LsdbData.AreaLsdb[lsdbKey] = lsdbEnt
	//Flood new Self Router LSA (areaId, lsaEnt, lsaKey)
	server.CreateAndSendMsgFromLsdbToFloodLsa(lsdbKey.AreaId, lsaKey, lsaEnt)
	return
}

func (server *OSPFV2Server) processRecvdSelfRouterLSA(msg RecvdSelfLsaMsg) error {
	lsa, ok := msg.LsaData.(RouterLsa)
	if !ok {
		server.logger.Err("Unable to assert given router lsa")
		return nil
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	lsaEnt, exist := lsdbEnt.RouterLsaMap[msg.LsaKey]
	if !exist {
		server.logger.Err("No such router LSA exist", msg.LsaKey)
		return nil
	}
	selfOrigLsaEnt, exist := server.LsdbData.AreaSelfOrigLsa[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such self originated LSA exist", msg.LsdbKey)
		return nil
	}
	_, exist = selfOrigLsaEnt[msg.LsaKey]
	if !exist {
		server.logger.Err("No such self originated router LSA exist", msg.LsaKey)
		return nil
	}
	if lsaEnt.LsaMd.LSSequenceNum < lsa.LsaMd.LSSequenceNum {
		checksumOffset := uint16(14)
		lsaEnt.LsaMd.LSSequenceNum = lsa.LsaMd.LSSequenceNum + 1
		lsaEnt.LsaMd.LSAge = 0
		lsaEnt.LsaMd.LSChecksum = 0
		lsaEnc := encodeRouterLsa(lsaEnt, msg.LsaKey)
		lsaEnt.LsaMd.LSChecksum = computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		lsdbEnt.RouterLsaMap[msg.LsaKey] = lsaEnt
		server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
		//Flood new Self Router LSA (areaId, lsaEnt, lsaKey)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
		return nil
	} else {
		//Flood existing Self Router LSA (areaId, lsaEnt, lsaKey)
		server.CreateAndSendMsgFromLsdbToFloodLsa(msg.LsdbKey.AreaId, msg.LsaKey, lsaEnt)
	}

	return nil
}

func (server *OSPFV2Server) processRecvdRouterLSA(msg RecvdLsaMsg) error {
	lsdbEnt, exist := server.LsdbData.AreaLsdb[msg.LsdbKey]
	if !exist {
		server.logger.Err("No such Area exist", msg.LsdbKey)
		return nil
	}
	if msg.MsgType == LSA_ADD {
		lsa, ok := msg.LsaData.(RouterLsa)
		if !ok {
			server.logger.Err("Unable to assert given router lsa")
			return nil
		}
		_, exist = lsdbEnt.RouterLsaMap[msg.LsaKey]
		lsdbEnt.RouterLsaMap[msg.LsaKey] = lsa
		if !exist {
			lsdbSlice := LsdbSliceStruct{
				LsdbKey: msg.LsdbKey,
				LsaKey:  msg.LsaKey,
			}
			server.GetBulkData.LsdbSlice = append(server.GetBulkData.LsdbSlice, lsdbSlice)
		}
	} else if msg.MsgType == LSA_DEL {
		delete(lsdbEnt.RouterLsaMap, msg.LsaKey)
	}
	server.LsdbData.AreaLsdb[msg.LsdbKey] = lsdbEnt
	return nil
}
