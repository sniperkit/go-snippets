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
	"strings"
)

func (h *OSPFHandler) SendOspfGlobal(ospfGlobalConf *ospfd.OspfGlobal) error {
	gConf := config.GlobalConf{
		Vrf:                ospfGlobalConf.Vrf,
		RouterId:           config.RouterId(ospfGlobalConf.RouterId),
		AdminStat:          config.Status(ospfGlobalConf.AdminStat),
		ASBdrRtrStatus:     ospfGlobalConf.ASBdrRtrStatus,
		TOSSupport:         ospfGlobalConf.TOSSupport,
		RestartSupport:     config.RestartSupport(ospfGlobalConf.RestartSupport),
		RestartInterval:    ospfGlobalConf.RestartInterval,
		ReferenceBandwidth: uint32(ospfGlobalConf.ReferenceBandwidth),
	}
	h.server.GlobalConfigCh <- gConf
	//	retMsg := <-h.server.GlobalConfigRetCh
	//	return retMsg
	return nil
}

func (h *OSPFHandler) SendOspfIfConf(ospfIfConf *ospfd.OspfIfEntry, msg config.MsgType) error {
	ifConf := config.InterfaceConf{
		IfIpAddress:       config.IpAddress(ospfIfConf.IfIpAddress),
		AddressLessIf:     config.InterfaceIndexOrZero(ospfIfConf.AddressLessIf),
		IfAreaId:          config.AreaId(ospfIfConf.IfAreaId),
		IfAdminStat:       config.Status(ospfIfConf.IfAdminStat),
		IfRtrPriority:     config.DesignatedRouterPriority(ospfIfConf.IfRtrPriority),
		IfTransitDelay:    config.UpToMaxAge(ospfIfConf.IfTransitDelay),
		IfRetransInterval: config.UpToMaxAge(ospfIfConf.IfRetransInterval),
		IfHelloInterval:   config.HelloRange(ospfIfConf.IfHelloInterval),
		IfRtrDeadInterval: config.PositiveInteger(ospfIfConf.IfRtrDeadInterval),
		IfPollInterval:    config.PositiveInteger(ospfIfConf.IfPollInterval),
	}

	for index, ifName := range config.IfTypeList {
		if strings.EqualFold(ospfIfConf.IfType, ifName) {
			ifConf.IfType = config.IfType(index)
			break
		}
	}

	ifMsg := config.InterfaceRpcMsg{
		IntfConf: ifConf,
		Op:       msg,
	}
	h.logger.Info(fmt.Sprintln("Before sending interface config "))
	h.server.IntfConfigCh <- ifMsg

	//retMsg := <-h.server.IntfConfigRetCh
	//return retMsg
	h.logger.Info(fmt.Sprintln("After receiving the create interface reply ..."))
	return nil
}

func (h *OSPFHandler) SendOspfAreaConf(ospfAreaConf *ospfd.OspfAreaEntry) error {
	areaConf := config.AreaConf{
		AreaId:          config.AreaId(ospfAreaConf.AreaId),
		AuthType:        config.AuthType(ospfAreaConf.AuthType),
		ImportAsExtern:  config.ImportAsExtern(ospfAreaConf.ImportAsExtern),
		AreaSummary:     config.AreaSummary(ospfAreaConf.AreaSummary),
		StubDefaultCost: ospfAreaConf.StubDefaultCost,
	}

	h.server.AreaConfigCh <- areaConf
	//	retMsg := <-h.server.AreaConfigRetCh
	//	return retMsg
	return nil

}

func (h *OSPFHandler) SendOspfIfMetricConf(ospfIfMetricConf *ospfd.OspfIfMetricEntry) error {
	ifMetricConf := config.IfMetricConf{
		IfMetricIpAddress:     config.IpAddress(ospfIfMetricConf.IfMetricIpAddress),
		IfMetricAddressLessIf: config.InterfaceIndexOrZero(ospfIfMetricConf.IfMetricAddressLessIf),
		IfMetricTOS:           config.TosType(ospfIfMetricConf.IfMetricTOS),
		IfMetricValue:         config.Metric(ospfIfMetricConf.IfMetricValue),
	}

	h.server.IfMetricConfCh <- ifMetricConf
	return nil
}

func (h *OSPFHandler) CreateOspfGlobal(ospfGlobalConf *ospfd.OspfGlobal) (bool, error) {
	if ospfGlobalConf == nil {
		err := errors.New("Invalid Global Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Create global config attrs:", ospfGlobalConf))
	err := h.SendOspfGlobal(ospfGlobalConf)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *OSPFHandler) CreateOspfAreaEntry(ospfAreaConf *ospfd.OspfAreaEntry) (bool, error) {
	if ospfAreaConf == nil {
		err := errors.New("Invalid Area Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Create Area config attrs:", ospfAreaConf))
	err := h.SendOspfAreaConf(ospfAreaConf)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *OSPFHandler) CreateOspfIfEntry(ospfIfConf *ospfd.OspfIfEntry) (bool, error) {
	if ospfIfConf == nil {
		err := errors.New("Invalid Interface Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Create interface config attrs:", ospfIfConf))
	err := h.SendOspfIfConf(ospfIfConf, config.CREATE)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *OSPFHandler) CreateOspfIfMetricEntry(ospfIfMetricConf *ospfd.OspfIfMetricEntry) (bool, error) {
	if ospfIfMetricConf == nil {
		err := errors.New("Invalid Interface Metric Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Create interface metric config attrs:", ospfIfMetricConf))
	err := h.SendOspfIfMetricConf(ospfIfMetricConf)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* This feature is not supported now
func (h *OSPFHandler) CreateOspfVirtIfEntry(ospfVirtIfConf *ospfd.OspfVirtIfEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Create virtual interface config attrs:", ospfVirtIfConf))
	return true, nil
}
*/
