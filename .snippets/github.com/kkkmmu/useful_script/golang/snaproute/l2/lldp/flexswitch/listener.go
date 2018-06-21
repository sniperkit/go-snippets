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

package flexswitch

import (
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"l2/lldp/api"
	"l2/lldp/config"
	"l2/lldp/utils"
	"lldpd"
	"strconv"
)

type ConfigHandler struct {
}

type NBPlugin struct {
	handler  *ConfigHandler
	fileName string
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func NewNBPlugin(handler *ConfigHandler, fileName string) *NBPlugin {
	l := &NBPlugin{handler, fileName}
	return l
}

func (p *NBPlugin) Start() error {
	fileName := p.fileName + CLIENTS_FILE_NAME

	clientJson, err := getClient(fileName, "lldpd")
	if err != nil || clientJson == nil {
		return err
	}
	debug.Logger.Info("Got Client Info for", clientJson.Name, " port", clientJson.Port)
	// create processor, transport and protocol for server
	processor := lldpd.NewLLDPDServicesProcessor(p.handler)
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTServerSocket("localhost:" +
		strconv.Itoa(clientJson.Port))
	if err != nil {
		debug.Logger.Info("StartServer: NewTServerSocket failed with error:", err)
		return err
	}
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	err = server.Serve()
	if err != nil {
		debug.Logger.Err("Failed to start the listener, err:", err)
		return err
	}
	return nil
}

func (h *ConfigHandler) CreateLLDPIntf(config *lldpd.LLDPIntf) (r bool, err error) {
	return api.SendIntfConfig(config.IntfRef, config.TxRxMode, config.Enable)
}

func (h *ConfigHandler) DeleteLLDPIntf(config *lldpd.LLDPIntf) (r bool, err error) {
	return false, errors.New("LLDP Enable/Disable is only supported via Update")
}

func (h *ConfigHandler) UpdateLLDPIntf(origconfig *lldpd.LLDPIntf,
	newconfig *lldpd.LLDPIntf, attrset []bool, op []*lldpd.PatchOpInfo) (r bool, err error) {
	// On update we do not care for old config... just push the new config to api layer
	// and let the api layer handle the information
	return api.UpdateIntfConfig(newconfig.IntfRef, newconfig.TxRxMode, newconfig.Enable)
}

func (h *ConfigHandler) GetLLDPIntf(intfRef string) (*lldpd.LLDPIntf, error) {
	return nil, nil
}

func (h *ConfigHandler) GetBulkLLDPIntf(fromIndex lldpd.Int, count lldpd.Int) (*lldpd.LLDPIntfGetInfo, error) {
	nextIdx, currCount, lldpIntfEntries := api.GetIntfs(int(fromIndex), int(count))
	if lldpIntfEntries == nil {
		return nil, errors.New("No neighbor found")
	}

	lldpEntryResp := make([]*lldpd.LLDPIntf, len(lldpIntfEntries))

	for idx, item := range lldpIntfEntries {
		lldpEntryResp[idx] = h.convertLLDPIntfEntryToThriftEntry(item)
	}

	lldpEntryBulk := lldpd.NewLLDPIntfGetInfo()
	lldpEntryBulk.StartIdx = fromIndex
	lldpEntryBulk.EndIdx = lldpd.Int(nextIdx)
	lldpEntryBulk.Count = lldpd.Int(currCount)
	lldpEntryBulk.More = (nextIdx != 0)
	lldpEntryBulk.LLDPIntfList = lldpEntryResp

	return lldpEntryBulk, nil
}

func (h *ConfigHandler) CreateLLDPGlobal(config *lldpd.LLDPGlobal) (r bool, err error) {
	debug.Logger.Debug("LLDP listener received create lldp global config", *config)
	return api.SendGlobalConfig(config.Vrf, config.TxRxMode, config.Enable, config.SnoopAndDrop, config.TranmitInterval)
}

func (h *ConfigHandler) DeleteLLDPGlobal(config *lldpd.LLDPGlobal) (r bool, err error) {
	return false, errors.New("LLDP Enable/Disable Globally is only supported via Update")
}

func (h *ConfigHandler) UpdateLLDPGlobal(origconfig *lldpd.LLDPGlobal,
	newconfig *lldpd.LLDPGlobal, attrset []bool, op []*lldpd.PatchOpInfo) (r bool, err error) {
	// On update we do not care for old config... just push the new config to api layer
	// and let the api layer handle the information
	debug.Logger.Debug("LLDP listener received update lldp global orig", *origconfig, "new config", *newconfig)
	return api.UpdateGlobalConfig(newconfig.Vrf, newconfig.TxRxMode, newconfig.Enable, newconfig.SnoopAndDrop, newconfig.TranmitInterval)
}

func (h *ConfigHandler) convertLLDPIntfEntryToThriftEntry(state config.Intf) *lldpd.LLDPIntf {
	entry := lldpd.NewLLDPIntf()
	entry.Enable = state.Enable
	entry.IntfRef = state.IntfRef
	entry.TxRxMode = state.TxRxMode
	return entry
}

func (h *ConfigHandler) convertLLDPIntfStateEntryToThriftEntry(state config.IntfState) *lldpd.LLDPIntfState {
	entry := lldpd.NewLLDPIntfState()
	entry.LocalPort = state.LocalPort
	entry.PeerMac = state.PeerMac
	entry.PeerPort = state.PeerPort
	entry.PeerHostName = state.PeerHostName
	entry.PortDescription = state.PortDescription
	entry.VlanId = state.VlanId
	entry.HoldTime = state.HoldTime
	entry.Enable = state.Enable
	entry.IfIndex = state.IfIndex
	entry.IntfRef = state.IntfRef
	entry.SendFrames = state.SendFrames
	entry.ReceivedFrames = state.ReceivedFrames
	entry.SystemDescription = state.SystemDescription
	entry.SystemCapabilities = state.SystemCapabilities
	entry.EnabledCapabilities = state.EnabledCapabilities
	return entry
}

func (h *ConfigHandler) GetBulkLLDPIntfState(fromIndex lldpd.Int, count lldpd.Int) (*lldpd.LLDPIntfStateGetInfo, error) {

	nextIdx, currCount, lldpIntfStateEntries := api.GetIntfStates(int(fromIndex), int(count))
	if lldpIntfStateEntries == nil {
		return nil, errors.New("No neighbor found")
	}

	lldpEntryResp := make([]*lldpd.LLDPIntfState, len(lldpIntfStateEntries))

	for idx, item := range lldpIntfStateEntries {
		lldpEntryResp[idx] = h.convertLLDPIntfStateEntryToThriftEntry(item)
	}

	lldpEntryBulk := lldpd.NewLLDPIntfStateGetInfo()
	lldpEntryBulk.StartIdx = fromIndex
	lldpEntryBulk.EndIdx = lldpd.Int(nextIdx)
	lldpEntryBulk.Count = lldpd.Int(currCount)
	lldpEntryBulk.More = (nextIdx != 0)
	lldpEntryBulk.LLDPIntfStateList = lldpEntryResp

	return lldpEntryBulk, nil
}

func (h *ConfigHandler) GetLLDPIntfState(intfRef string) (*lldpd.LLDPIntfState, error) {
	lldpIntf := api.GetIntfState(intfRef)
	if lldpIntf == nil {
		return nil, errors.New(fmt.Sprintln("No Information found for", intfRef))
	}

	return h.convertLLDPIntfStateEntryToThriftEntry(*lldpIntf), nil
}

func (h *ConfigHandler) convertLLDPGlobalStateEntryToThriftEntry(gblState *config.GlobalState) *lldpd.LLDPGlobalState {
	entry := lldpd.NewLLDPGlobalState()
	entry.Vrf = gblState.Vrf
	entry.Enable = gblState.Enable
	entry.TranmitInterval = gblState.TranmitInterval
	entry.TotalTxFrames = gblState.TotalTxFrames
	entry.TotalRxFrames = gblState.TotalRxFrames
	entry.Neighbors = gblState.Neighbors
	return entry
}

func (h *ConfigHandler) GetBulkLLDPGlobalState(fromIndex lldpd.Int, count lldpd.Int) (*lldpd.LLDPGlobalStateGetInfo, error) {
	lldpGlobalStateBulk := lldpd.NewLLDPGlobalStateGetInfo()
	lldpGlobalStateBulk.EndIdx = lldpd.Int(0)
	lldpGlobalStateBulk.Count = lldpd.Int(1)
	lldpGlobalStateBulk.More = false
	lldpGlobalStateBulk.LLDPGlobalStateList = make([]*lldpd.LLDPGlobalState, 1)
	gblEntry, _ := api.GetLLDPGlobalState("default")
	lldpGlobalStateBulk.LLDPGlobalStateList[0] = h.convertLLDPGlobalStateEntryToThriftEntry(gblEntry)

	return lldpGlobalStateBulk, nil
}

func (h *ConfigHandler) GetLLDPGlobalState(vrf string) (*lldpd.LLDPGlobalState, error) {
	gblEntry, _ := api.GetLLDPGlobalState(vrf)
	return h.convertLLDPGlobalStateEntryToThriftEntry(gblEntry), nil
}
