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

package rpc

import (
	"errors"
	"fmt"
	"l3/ospf/config"
	"ospfd"
	//    "l3/ospf/server"
	//    "utils/logging"
	//    "net"
)

func (h *OSPFHandler) convertAreaEntryStateToThrift(ent config.AreaState) *ospfd.OspfAreaEntryState {
	areaEntry := ospfd.NewOspfAreaEntryState()
	areaEntry.AreaId = string(ent.AreaId)
	areaEntry.SpfRuns = ent.SpfRuns
	areaEntry.AreaBdrRtrCount = ent.AreaBdrRtrCount
	areaEntry.AsBdrRtrCount = ent.AsBdrRtrCount
	areaEntry.AreaLsaCount = ent.AreaLsaCount

	return areaEntry
}

func (h *OSPFHandler) convertLsdbEntryStateToThrift(ent config.LsdbState) *ospfd.OspfLsdbEntryState {
	//h.logger.Info(fmt.Sprintln("Converting Lsdb entry to Thrift", ent))
	lsdbEntry := ospfd.NewOspfLsdbEntryState()
	lsdbEntry.LsdbType = int32(ent.LsdbType)
	lsdbEntry.LsdbLsid = string(ent.LsdbLsid)
	lsdbEntry.LsdbAreaId = string(ent.LsdbAreaId)
	lsdbEntry.LsdbRouterId = string(ent.LsdbRouterId)
	lsdbEntry.LsdbSequence = int32(ent.LsdbSequence)
	lsdbEntry.LsdbAge = int32(ent.LsdbAge)
	lsdbEntry.LsdbChecksum = int32(ent.LsdbCheckSum)
	lsdbEntry.LsdbAdvertisement = string(ent.LsdbAdvertisement)

	return lsdbEntry
}

func (h *OSPFHandler) convertIfEntryStateToThrift(ent config.InterfaceState) *ospfd.OspfIfEntryState {
	ifEntry := ospfd.NewOspfIfEntryState()
	ifEntry.IfIpAddress = string(ent.IfIpAddress)
	ifEntry.AddressLessIf = int32(ent.AddressLessIf)
	ifEntry.IfState = int32(ent.IfState)
	ifEntry.IfDesignatedRouter = string(ent.IfDesignatedRouter)
	ifEntry.IfBackupDesignatedRouter = string(ent.IfBackupDesignatedRouter)
	ifEntry.IfEvents = int32(ent.IfEvents)
	ifEntry.IfLsaCount = int32(ent.IfLsaCount)
	ifEntry.IfDesignatedRouterId = string(ent.IfDesignatedRouterId)
	ifEntry.IfBackupDesignatedRouterId = string(ent.IfBackupDesignatedRouter)

	return ifEntry
}

func (h *OSPFHandler) convertGlobalStateToThrift(ent config.GlobalState) *ospfd.OspfGlobalState {
	gState := ospfd.NewOspfGlobalState()
	gState.RouterId = string(ent.RouterId)
	gState.VersionNumber = ent.VersionNumber
	gState.AreaBdrRtrStatus = ent.AreaBdrRtrStatus
	gState.ExternLsaCount = ent.ExternLsaCount
	gState.OpaqueLsaSupport = ent.OpaqueLsaSupport

	return gState
}

func (h *OSPFHandler) convertNbrEntryStateToThrift(nbr config.NeighborState) *ospfd.OspfNbrEntryState {
	nbrEntry := ospfd.NewOspfNbrEntryState()
	nbrEntry.NbrIpAddr = string(nbr.NbrIpAddress)
	nbrEntry.NbrAddressLessIndex = int32(nbr.NbrAddressLessIndex)
	nbrEntry.NbrRtrId = string(nbr.NbrRtrId)
	nbrEntry.NbrOptions = int32(nbr.NbrOptions)
	nbrEntry.NbrState = string(nbr.NbrState)
	nbrEntry.NbrEvents = int32(nbr.NbrEvents)
	nbrEntry.NbrHelloSuppressed = bool(nbr.NbrHelloSuppressed)

	return nbrEntry

}

func (h *OSPFHandler) GetBulkOspfAreaEntryState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfAreaEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Area attrs"))

	nextIdx, currCount, ospfAreaEntryStates := h.server.GetBulkOspfAreaEntryState(int(fromIdx), int(count))
	if ospfAreaEntryStates == nil {
		err := errors.New("Ospf is busy refreshing the cache")
		return nil, err
	}
	ospfAreaEntryStateResponse := make([]*ospfd.OspfAreaEntryState, len(ospfAreaEntryStates))
	for idx, item := range ospfAreaEntryStates {
		ospfAreaEntryStateResponse[idx] = h.convertAreaEntryStateToThrift(item)
	}
	ospfAreaEntryStateGetInfo := ospfd.NewOspfAreaEntryStateGetInfo()
	ospfAreaEntryStateGetInfo.Count = ospfd.Int(currCount)
	ospfAreaEntryStateGetInfo.StartIdx = ospfd.Int(fromIdx)
	ospfAreaEntryStateGetInfo.EndIdx = ospfd.Int(nextIdx)
	ospfAreaEntryStateGetInfo.More = (nextIdx != 0)
	ospfAreaEntryStateGetInfo.OspfAreaEntryStateList = ospfAreaEntryStateResponse
	return ospfAreaEntryStateGetInfo, nil
}

func (h *OSPFHandler) GetBulkOspfLsdbEntryState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfLsdbEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Link State Database attrs"))
	nextIdx, currCount, ospfLsdbEntryStates := h.server.GetBulkOspfLsdbEntryState(int(fromIdx), int(count))
	if ospfLsdbEntryStates == nil {
		err := errors.New("Ospf is busy refreshing the cache")
		return nil, err
	}
	ospfLsdbEntryStateResponse := make([]*ospfd.OspfLsdbEntryState, len(ospfLsdbEntryStates))
	for idx, item := range ospfLsdbEntryStates {
		//h.logger.Info(fmt.Sprintln("converting Lsdb Entry into thrift format", item))
		ospfLsdbEntryStateResponse[idx] = h.convertLsdbEntryStateToThrift(item)
		//h.logger.Info(fmt.Sprintln("After converting Lsdb Entry into thrift format", idx))
	}
	ospfLsdbEntryStateGetInfo := ospfd.NewOspfLsdbEntryStateGetInfo()
	ospfLsdbEntryStateGetInfo.Count = ospfd.Int(currCount)
	ospfLsdbEntryStateGetInfo.StartIdx = ospfd.Int(fromIdx)
	ospfLsdbEntryStateGetInfo.EndIdx = ospfd.Int(nextIdx)
	ospfLsdbEntryStateGetInfo.More = (nextIdx != 0)
	ospfLsdbEntryStateGetInfo.OspfLsdbEntryStateList = ospfLsdbEntryStateResponse
	return ospfLsdbEntryStateGetInfo, nil
}

func (h *OSPFHandler) GetBulkOspfIfEntryState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfIfEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Interface attrs"))

	nextIdx, currCount, ospfIfEntryStates := h.server.GetBulkOspfIfEntryState(int(fromIdx), int(count))
	if ospfIfEntryStates == nil {
		err := errors.New("Ospf is busy refreshing the cache")
		return nil, err
	}
	ospfIfEntryStateResponse := make([]*ospfd.OspfIfEntryState, len(ospfIfEntryStates))
	for idx, item := range ospfIfEntryStates {
		ospfIfEntryStateResponse[idx] = h.convertIfEntryStateToThrift(item)
	}
	ospfIfEntryStateGetInfo := ospfd.NewOspfIfEntryStateGetInfo()
	ospfIfEntryStateGetInfo.Count = ospfd.Int(currCount)
	ospfIfEntryStateGetInfo.StartIdx = ospfd.Int(fromIdx)
	ospfIfEntryStateGetInfo.EndIdx = ospfd.Int(nextIdx)
	ospfIfEntryStateGetInfo.More = (nextIdx != 0)
	ospfIfEntryStateGetInfo.OspfIfEntryStateList = ospfIfEntryStateResponse
	return ospfIfEntryStateGetInfo, nil
}

func (h *OSPFHandler) GetBulkOspfNbrEntryState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfNbrEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Neighbor attrs"))
	nextIdx, currCount, ospfNbrEntryStates := h.server.GetBulkOspfNbrEntryState(int(fromIdx), int(count))
	if ospfNbrEntryStates == nil {
		err := errors.New("Ospf NBR is busy refreshing the cache")
		return nil, err
	}
	ospfNbrEntryResponse := make([]*ospfd.OspfNbrEntryState, len(ospfNbrEntryStates))
	for idx, item := range ospfNbrEntryStates {
		ospfNbrEntryResponse[idx] = h.convertNbrEntryStateToThrift(item)
	}
	OspfNbrEntryStateGetInfo := ospfd.NewOspfNbrEntryStateGetInfo()
	OspfNbrEntryStateGetInfo.Count = ospfd.Int(currCount)
	OspfNbrEntryStateGetInfo.StartIdx = ospfd.Int(fromIdx)
	OspfNbrEntryStateGetInfo.EndIdx = ospfd.Int(nextIdx)
	OspfNbrEntryStateGetInfo.More = (nextIdx != 0)
	OspfNbrEntryStateGetInfo.OspfNbrEntryStateList = ospfNbrEntryResponse
	return OspfNbrEntryStateGetInfo, nil
}

/* Not supported
func (h *OSPFHandler) GetBulkOspfVirtNbrEntryState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfVirtNbrEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Virtual Neighbor attrs"))
	ospfVirtNbrResponse := ospfd.NewOspfVirtNbrEntryStateGetInfo()
	return ospfVirtNbrResponse, nil
}
*/

func (h *OSPFHandler) GetBulkOspfGlobalState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfGlobalStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get OSPF global state"))

	if fromIdx != 0 {
		err := errors.New("Invalid range")
		return nil, err
	}
	ospfGlobalState := h.server.GetOspfGlobalState()
	ospfGlobalStateResponse := make([]*ospfd.OspfGlobalState, 1)
	ospfGlobalStateResponse[0] = h.convertGlobalStateToThrift(*ospfGlobalState)
	ospfGlobalStateGetInfo := ospfd.NewOspfGlobalStateGetInfo()
	ospfGlobalStateGetInfo.Count = ospfd.Int(1)
	ospfGlobalStateGetInfo.StartIdx = ospfd.Int(0)
	ospfGlobalStateGetInfo.EndIdx = ospfd.Int(0)
	ospfGlobalStateGetInfo.More = false
	ospfGlobalStateGetInfo.OspfGlobalStateList = ospfGlobalStateResponse
	return ospfGlobalStateGetInfo, nil
}

func (h *OSPFHandler) GetBulkOspfIPv4RouteState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfIPv4RouteStateGetInfo, error) {
	/* This is template API . All OSPF routes are currently installed in the redis db */
	return nil, nil
}

func (h *OSPFHandler) GetBulkOspfEventState(fromIdx ospfd.Int, count ospfd.Int) (*ospfd.OspfEventStateGetInfo, error) {
	/* This is template API. Events are stored in redis-db */
	return nil, nil
}
