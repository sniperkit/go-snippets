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
	"errors"
	"fmt"
	"l3/ospfv2/objects"
)

func newDbdMsg(key NbrConfKey, dbd_data NbrDbdData) NbrDbdMsg {
	dbdNbrMsg := NbrDbdMsg{
		nbrConfKey: key,
		nbrDbdData: dbd_data,
	}
	return dbdNbrMsg
}

func (server *OSPFV2Server) UpdateNbrConf(nbrKey NbrConfKey, conf NbrConf, flags int) {
	nbrE, valid := server.NbrConfMap[nbrKey]
	if !valid {
		server.logger.Err("Nbr : Nbr conf does not exist . Not updated ", nbrKey)
		return
	}
	if flags&NBR_FLAG_STATE == NBR_FLAG_STATE {
		nbrE.State = conf.State
	}
	if flags&NBR_FLAG_DEAD_TIMER == NBR_FLAG_DEAD_TIMER {
		server.logger.Debug("Nbr : Nbr inactivity reset ", nbrKey)
	}
	if flags&NBR_FLAG_SEQ_NUMBER == NBR_FLAG_SEQ_NUMBER {
		nbrE.DDSequenceNum = conf.DDSequenceNum
	}

	if flags&NBR_FLAG_IS_MASTER == NBR_FLAG_IS_MASTER {
		nbrE.isMaster = conf.isMaster
	}
	if flags&NBR_FLAG_PRIORITY == NBR_FLAG_PRIORITY {
		nbrE.NbrPriority = conf.NbrPriority
	}
	if flags&NBR_FLAG_OPTION == NBR_FLAG_STATE {
		nbrE.NbrOption = conf.NbrOption
	}
	if flags&NBR_FLAG_REQ_LIST == NBR_FLAG_REQ_LIST {
		if len(conf.NbrReqList) > 0 {
			nbrE.NbrReqList = conf.NbrReqList
			server.logger.Debug("nbr: updated req list ", len(nbrE.NbrReqList))
		}
	}
	if flags&NBR_FLAG_REQ_LIST_INDEX == NBR_FLAG_REQ_LIST_INDEX {
		nbrE.NbrReqListIndex = conf.NbrReqListIndex
	}
	server.NbrConfMap[nbrKey] = nbrE
}

func (server *OSPFV2Server) UpdateIntfToNbrMap(nbrKey NbrConfKey) {
	var newList []NbrConfKey
	nbrConf := server.NbrConfMap[nbrKey]
	_, exists := server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey]
	if !exists {
		newList = []NbrConfKey{}
	} else {
		newList = server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey]
	}
	for _, nbr := range newList {
		if nbr.NbrAddressLessIfIdx == nbrKey.NbrAddressLessIfIdx &&
			nbr.NbrIdentity == nbrKey.NbrIdentity {
			return
		}
	}
	newList = append(newList, nbrKey)
	server.NbrConfData.IntfToNbrMap[nbrConf.IntfKey] = newList
	server.logger.Debug("Nbr : Intf to nbr list updated ", newList)
}

func (server *OSPFV2Server) NbrDbPacketDiscardCheck(nbrDbPkt NbrDbdData, nbrConf NbrConf) bool {
	if nbrDbPkt.msbit != nbrConf.isMaster {
		server.logger.Info("NBREVENT: SeqNumberMismatch. Nbr should be master  dbdmsbit ", nbrDbPkt.msbit,
			" isMaster ", nbrConf.isMaster)
		return true
	}

	if nbrDbPkt.ibit == true {
		server.logger.Info("NBREVENT:SeqNumberMismatch . Nbr ibit is true ", nbrConf.NbrIP)
		return true
	}

	if nbrConf.isMaster {
		if nbrDbPkt.dd_sequence_number != nbrConf.DDSequenceNum+1 {
			server.logger.Info(fmt.Sprintln("NBREVENT:SeqNumberMismatch : Nbr is master but dbd packet seq no doesnt match. dbd seq ",
				nbrDbPkt.dd_sequence_number, "nbr seq ", nbrConf.DDSequenceNum))

			return true
		}
	} else {
		if nbrDbPkt.dd_sequence_number != nbrConf.DDSequenceNum+1 {
			server.logger.Info(fmt.Sprintln("NBREVENT:SeqNumberMismatch : Nbr is slave but dbd packet seq no doesnt match.dbd seq ",
				nbrDbPkt.dd_sequence_number, "nbr seq ", nbrConf.DDSequenceNum))
			return true
		}
	}

	return false
}

func (server *OSPFV2Server) verifyDuplicatePacket(nbrConf NbrConf, nbrDbPkt NbrDbdData) (isDup bool) {
	if nbrConf.isMaster {
		if nbrDbPkt.dd_sequence_number+1 == nbrConf.DDSequenceNum {
			isDup = true
			server.logger.Info(fmt.Sprintln("NBREVENT: Duplicate packet Dont do anything. dbdseq ",
				nbrDbPkt.dd_sequence_number, " nbrseq ", nbrConf.DDSequenceNum))
			return
		}
	}
	isDup = false
	return
}

func calculateMaxLsaHeaders() (max_headers uint8) {
	rem := INTF_MTU_MIN - (OSPF_DBD_MIN_SIZE + OSPF_HEADER_SIZE)
	max_headers = uint8(rem / OSPF_LSA_HEADER_SIZE)
	return max_headers
}

func calculateMaxLsaReq() (max_req int) {
	rem := INTF_MTU_MIN - OSPF_HEADER_SIZE
	max_req = rem / OSPF_LSA_REQ_SIZE
	return max_req
}

func (server *OSPFV2Server) addNbrToSlice(nbrKey NbrConfKey) {
	add := true
	for _, nbr := range server.GetBulkData.NbrConfSlice {
		if nbr.NbrIdentity == nbrKey.NbrIdentity &&
			nbr.NbrAddressLessIfIdx == nbrKey.NbrAddressLessIfIdx {
			add = false
		}
	}
	if add {
		server.GetBulkData.NbrConfSlice = append(server.GetBulkData.NbrConfSlice, nbrKey)
	}
}

func (server *OSPFV2Server) delNbrFromSlice(nbrKey NbrConfKey) {
	for index, nbr := range server.GetBulkData.NbrConfSlice {
		if nbr.NbrIdentity == nbrKey.NbrIdentity &&
			nbr.NbrAddressLessIfIdx == nbrKey.NbrAddressLessIfIdx {
			server.GetBulkData.NbrConfSlice = append(server.GetBulkData.NbrConfSlice[:index],
				server.GetBulkData.NbrConfSlice[index+1:]...)
		}
	}
}

/**** Get bulk APis ***/
func (server *OSPFV2Server) RefreshNbrConfSlice() {
	if len(server.GetBulkData.NbrConfSlice) == 0 {
		return
	}
	server.GetBulkData.NbrConfSlice = server.GetBulkData.NbrConfSlice[:len(server.GetBulkData.NbrConfSlice)-1]
	server.GetBulkData.NbrConfSlice = nil
	for nbrKey, _ := range server.NbrConfMap {
		server.GetBulkData.NbrConfSlice = append(server.GetBulkData.NbrConfSlice, nbrKey)
	}
}

func (server *OSPFV2Server) getNbrState(ipAddr, addressLessIfIdx uint32) (*objects.Ospfv2NbrState, error) {
	var retObj objects.Ospfv2NbrState

	nbrKey := NbrConfKey{
		NbrIdentity:         ipAddr,
		NbrAddressLessIfIdx: addressLessIfIdx,
	}
	nbr, valid := server.NbrConfMap[nbrKey]
	if !valid {
		return nil, errors.New("Nbr does not exist ")
	}
	retObj.AddressLessIfIdx = addressLessIfIdx
	retObj.IpAddr = nbr.NbrIP
	retObj.RtrId = nbr.NbrRtrId
	retObj.State = uint8(nbr.State)
	retObj.Options = int32(nbr.NbrOption)

	return &retObj, nil
}

func (server *OSPFV2Server) getBulkNbrState(fromIdx, cnt int) (*objects.Ospfv2NbrStateGetInfo, error) {
	var retObj objects.Ospfv2NbrStateGetInfo
	count := 0
	idx := fromIdx
	sliceLen := len(server.GetBulkData.NbrConfSlice)
	server.logger.Debug("Nbr : Total elements in nbr slice ", sliceLen)
	if fromIdx >= sliceLen {
		return nil, errors.New("Invalid Range")
	}
	for count < cnt {
		if idx == sliceLen {
			break
		}
		nbrKey := server.GetBulkData.NbrConfSlice[idx]
		nbrEnt, exist := server.NbrConfMap[nbrKey]
		if !exist {
			idx++
			continue
		}
		var obj objects.Ospfv2NbrState
		obj.AddressLessIfIdx = nbrKey.NbrAddressLessIfIdx
		obj.IpAddr = nbrEnt.NbrIP
		obj.Options = int32(nbrEnt.NbrOption)
		obj.RtrId = uint32(nbrEnt.NbrRtrId)
		obj.State = uint8(nbrEnt.State)
		retObj.List = append(retObj.List, &obj)
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
