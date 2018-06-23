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
	"fmt"
	"ospfd"
)

func (h *OSPFHandler) GetOspfGlobalState(routerId string) (*ospfd.OspfGlobalState, error) {
	h.logger.Info(fmt.Sprintln("Get global attrs"))
	ospfGlobalResponse := ospfd.NewOspfGlobalState()
	ospfGlobalState := h.server.GetOspfGlobalState()
	ospfGlobalResponse = h.convertGlobalStateToThrift(*ospfGlobalState)
	return ospfGlobalResponse, nil
}

func (h *OSPFHandler) GetOspfAreaEntryState(areaId string) (*ospfd.OspfAreaEntryState, error) {
	h.logger.Info(fmt.Sprintln("Get Area attrs"))
	ospfAreaResponse := ospfd.NewOspfAreaEntryState()
	return ospfAreaResponse, nil
}

/*
func (h *OSPFHandler) GetOspfStubAreaEntryState(stubAreaId string, stubTOS int32) (*ospfd.OspfStubAreaState, error) {
    h.logger.Info(fmt.Sprintln("Get Area Stub attrs"))
    ospfStubAreaResponse := ospfd.NewOspfStubAreaState()
    return ospfStubAreaResponse, nil
}
*/

func (h *OSPFHandler) GetOspfLsdbEntryState(lsdbType int32, lsdbLsid string, lsdbAreaId string, lsdbRouterId string) (*ospfd.OspfLsdbEntryState, error) {
	h.logger.Info(fmt.Sprintln("Get Link State Database attrs"))
	ospfLsdbResponse := ospfd.NewOspfLsdbEntryState()
	return ospfLsdbResponse, nil
}

/*
func (h *OSPFHandler) GetOspfAreaRangeEntryState(rangeAreaId string, areaRangeNet string) (*ospfd.OspfAreaRangeState, error) {
    h.logger.Info(fmt.Sprintln("Get Address range attrs"))
    ospfAreaRangeResponse := ospfd.NewOspfAreaRangeState()
    return ospfAreaRangeResponse, nil
}
*/

func (h *OSPFHandler) GetOspfIfEntryState(ifIpAddress string, addressLessIf int32) (*ospfd.OspfIfEntryState, error) {
	h.logger.Info(fmt.Sprintln("Get Interface attrs"))
	ospfIfResponse := ospfd.NewOspfIfEntryState()
	return ospfIfResponse, nil
}

/*
func (h *OSPFHandler) GetOspfIfMetricState(ifMetricIpAddress string, ifMetricAddressLessIf int32, ifMetricTOS int32) (*ospfd.OspfIfMetricState, error) {
    h.logger.Info(fmt.Sprintln("Get Interface Metric attrs"))
    ospfIfMetricResponse := ospfd.NewOspfIfMetricState()
    return ospfIfMetricResponse, nil
}

func (h *OSPFHandler) GetOspfVirtIfState(virtIfAreaId string, virtIfNeighbor string) (*ospfd.OspfVirtIfState, error) {
    h.logger.Info(fmt.Sprintln("Get Virtual Interface attrs"))
    ospfVirtIfResponse := ospfd.NewOspfVirtIfState()
    return ospfVirtIfResponse, nil
}
*/

func (h *OSPFHandler) GetOspfNbrEntryState(nbrIpAddr string, nbrAddressLessIndex int32) (*ospfd.OspfNbrEntryState, error) {
	h.logger.Info(fmt.Sprintln("Get Neighbor attrs"))
	ospfNbrResponse := ospfd.NewOspfNbrEntryState()
	return ospfNbrResponse, nil
}

/* Not supported
func (h *OSPFHandler) GetOspfVirtNbrEntryState(virtNbrRtrId string, virtNbrArea string) (*ospfd.OspfVirtNbrEntryState, error) {
	h.logger.Info(fmt.Sprintln("Get Virtual Neighbor attrs"))
	ospfVirtNbrResponse := ospfd.NewOspfVirtNbrEntryState()
	return ospfVirtNbrResponse, nil
}
*/

func (h *OSPFHandler) GetOspfIPv4RouteState(destId string, addrMask string, destType string) (*ospfd.OspfIPv4RouteState, error) {
	return nil, nil
}

func (h *OSPFHandler) GetOspfEventState(Index int32) (*ospfd.OspfEventState, error) {
	return nil, nil
}
