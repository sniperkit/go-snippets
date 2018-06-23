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

package manager

import (
	"dhcprelayd"
	"errors"
	"fmt"
	"l3/dhcp_relay/infra"
	"l3/dhcp_relay/protocol/dhcp4"
	"l3/dhcp_relay/protocol/dhcp6"
	"net"
	"utils/commonDefs"
	"utils/dbutils"
	"utils/logging"
)

type DRAMgr struct {
	DbHdl  dbutils.DBIntf
	Logger logging.LoggerIntf
	IMgr   *infra.InfraMgr

	PProc4 IPv4ProcessorIntf
	PProc6 IPv6ProcessorIntf
}

func NewDRAMgr(logger logging.LoggerIntf,
	dbHdl dbutils.DBIntf, infraMgr *infra.InfraMgr) *DRAMgr {

	draMgr := &DRAMgr{}
	draMgr.Logger = logger
	draMgr.DbHdl = dbHdl
	draMgr.IMgr = infraMgr
	draMgr.PProc4 = dhcp4.NewProcessor(
		&dhcp4.ProcessorInitParams{
			Logger:   draMgr.Logger,
			InfraMgr: draMgr.IMgr,
		},
	)
	draMgr.PProc6 = dhcp6.NewProcessor(
		&dhcp6.ProcessorInitParams{
			Logger:   draMgr.Logger,
			InfraMgr: draMgr.IMgr,
		},
	)
	return draMgr
}

func (draMgr *DRAMgr) InitDRAMgr() bool {
	draMgr.IMgr.BuildInfra()
	return draMgr.initDRAv4() && draMgr.initDRAv6()
}

func (draMgr *DRAMgr) initDRAv4() bool {
	draGbl, err := draMgr.readDRAv4GlobalConfig()
	if err != nil {
		draMgr.Logger.Err("DB Read failed:", err)
		return false
	}
	draIntfs, err := draMgr.readDRAv4IntfConfig()
	if err != nil {
		draMgr.Logger.Err("DB Read failed:", err)
		return false
	}
	if draGbl != nil {
		draMgr.IMgr.UpdateDRAv4Global(draGbl)
	}
	for _, draIntf := range draIntfs {
		draMgr.IMgr.UpdateDRAv4Intf(draIntf)
		draMgr.PProc4.ProcessCreateDRAIntf(draIntf.IntfRef)
	}
	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		draMgr.PProc4.StopRxTx()
		return true
	}
	draMgr.PProc4.StartRxTx()
	for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv4Intfs() {
		draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
	}
	return true
}

func (draMgr *DRAMgr) initDRAv6() bool {
	draGbl, err := draMgr.readDRAv6GlobalConfig()
	if err != nil {
		draMgr.Logger.Err("DB Read failed:", err)
		return false
	}
	draIntfs, err := draMgr.readDRAv6IntfConfig()
	if err != nil {
		draMgr.Logger.Err("DB Read failed:", err)
		return false
	}
	if draGbl != nil {
		draMgr.IMgr.UpdateDRAv6Global(draGbl)
	}
	for _, draIntf := range draIntfs {
		draMgr.IMgr.UpdateDRAv6Intf(draIntf)
		draMgr.PProc6.ProcessCreateDRAIntf(draIntf.IntfRef)
	}
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		draMgr.PProc6.StopRxTx()
		return true
	}
	draMgr.PProc6.StartRxTx()
	for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv6Intfs() {
		draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
	}
	return true
}

func (draMgr *DRAMgr) CreateDRAv4Global(
	cfg *dhcprelayd.DHCPRelayGlobal) (bool, error) {

	if cfg.Vrf != "default" {
		errMsg := fmt.Sprintln("DRA: only \"default\" Vrf supported")
		draMgr.Logger.Err(errMsg)
		return false, errors.New(errMsg)
	}
	draMgr.IMgr.UpdateDRAv4Global(cfg)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		return true, nil
	}
	draMgr.PProc4.StartRxTx()
	for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv4Intfs() {
		if _, ok := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx); ok {
			draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
		}
	}
	return true, nil
}

func (draMgr *DRAMgr) UpdateDRAv4Global(oldCfg *dhcprelayd.DHCPRelayGlobal,
	newCfg *dhcprelayd.DHCPRelayGlobal, attrSet []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	cfg := oldCfg
	if attrSet[2] {
		cfg.HopCountLimit = newCfg.HopCountLimit
	}
	if attrSet[1] {
		cfg.Enable = newCfg.Enable
	}
	draPreState := (draMgr.IMgr.GetActiveDRAv4IntfCount() > 0)
	draMgr.IMgr.UpdateDRAv4Global(cfg)
	draPostState := (draMgr.IMgr.GetActiveDRAv4IntfCount() > 0)
	if !draPostState {
		draMgr.PProc4.StopRxTx()
		return true, nil
	}
	if !draPreState && draPostState {
		draMgr.PProc4.StartRxTx()
		for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv4Intfs() {
			if _, ok := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx); ok {
				draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
			}
		}
	}
	return true, nil
}

func (draMgr *DRAMgr) DeleteDRAv4Global(Vrf string) (bool, error) {
	draMgr.IMgr.DeleteDRAv4Global()
	draMgr.PProc4.StopRxTx()
	return true, nil
}

func (draMgr *DRAMgr) CreateDRAv4Interface(
	cfg *dhcprelayd.DHCPRelayIntf) (bool, error) {

	for _, val := range cfg.ServerIp {
		ip := net.ParseIP(val)
		if ip == nil || ip.To4() == nil {
			errMsg := fmt.Sprintln(
				"DRA: Parsing of ServerIp address failed:", val,
			)
			draMgr.Logger.Err(errMsg)
			return false, errors.New(errMsg)
		}
	}

	ifIdx, ok := draMgr.IMgr.GetIPv4IntfIndex(cfg.IntfRef)
	if !ok {
		errMsg := fmt.Sprintln(
			"DRA: No IPv4 Intf exists with reference", cfg.IntfRef)
		draMgr.Logger.Err(errMsg)
		return false, errors.New(errMsg)
	}
	draMgr.IMgr.UpdateDRAv4Intf(cfg)
	draMgr.PProc4.ProcessCreateDRAIntf(cfg.IntfRef)
	if _, ok := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx); ok {
		draMgr.PProc4.StartRxTx()
		draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) UpdateDRAv4Interface(oldCfg *dhcprelayd.DHCPRelayIntf,
	newCfg *dhcprelayd.DHCPRelayIntf, attrSet []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	cfg := oldCfg
	if attrSet[2] {
		for _, val := range newCfg.ServerIp {
			ip := net.ParseIP(val)
			if ip == nil || ip.To4() == nil {
				errMsg := fmt.Sprintln(
					"DRA: Parsing of ServerIp address failed:", val)
				draMgr.Logger.Err(errMsg)
				return false, errors.New(errMsg)
			}
		}
		cfg.ServerIp = newCfg.ServerIp
	}
	if attrSet[1] {
		cfg.Enable = newCfg.Enable
	}
	ifIdx, ok := draMgr.IMgr.GetIPv4IntfIndex(newCfg.IntfRef)
	if !ok { // No valid interface
		draMgr.IMgr.UpdateDRAv4Intf(cfg)
		return true, nil
	}

	_, draPreState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)
	draMgr.IMgr.UpdateDRAv4Intf(cfg)
	_, draPostState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		draMgr.PProc4.StopRxTx()
		return true, nil
	}
	draMgr.PProc4.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc4.ProcessInactiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) DeleteDRAv4Interface(intfRef string) (bool, error) {
	ifIdx, ok := draMgr.IMgr.GetIPv4IntfIndex(intfRef)
	if !ok {
		draMgr.Logger.Info(
			"DRA: Cannot find interface with reference,", intfRef)
	}
	_, draPreState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)
	draMgr.IMgr.DeleteDRAv4Intf(intfRef)
	draMgr.PProc4.ProcessDeleteDRAIntf(intfRef)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		draMgr.PProc4.StopRxTx()
		return true, nil
	}
	if draPreState {
		draMgr.PProc4.ProcessInactiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) CreateDRAv6Global(
	cfg *dhcprelayd.DHCPv6RelayGlobal) (bool, error) {

	if cfg.Vrf != "default" {
		errMsg := fmt.Sprintln("DRA: only \"default\" Vrf supported")
		draMgr.Logger.Err(errMsg)
		return false, errors.New(errMsg)
	}
	draMgr.IMgr.UpdateDRAv6Global(cfg)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		return true, nil
	}
	draMgr.PProc6.StartRxTx()
	for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv6Intfs() {
		if _, ok := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx); ok {
			draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
		}
	}
	return true, nil
}

func (draMgr *DRAMgr) UpdateDRAv6Global(oldCfg *dhcprelayd.DHCPv6RelayGlobal,
	newCfg *dhcprelayd.DHCPv6RelayGlobal, attrSet []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	cfg := oldCfg
	if attrSet[2] {
		cfg.HopCountLimit = newCfg.HopCountLimit
	}
	if attrSet[1] {
		cfg.Enable = newCfg.Enable
	}
	draPreState := (draMgr.IMgr.GetActiveDRAv6IntfCount() > 0)
	draMgr.IMgr.UpdateDRAv6Global(cfg)
	draPostState := (draMgr.IMgr.GetActiveDRAv6IntfCount() > 0)
	if !draPostState {
		draMgr.PProc6.StopRxTx()
		return true, nil
	}
	if !draPreState && draPostState {
		draMgr.PProc6.StartRxTx()
		for _, ifIdx := range draMgr.IMgr.GetAllActiveDRAv6Intfs() {
			if _, ok := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx); ok {
				draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
			}
		}
	}
	return true, nil
}

func (draMgr *DRAMgr) DeleteDRAv6Global(Vrf string) (bool, error) {
	draMgr.IMgr.DeleteDRAv6Global()
	draMgr.PProc6.StopRxTx()
	return true, nil
}

func checkStringSliceUnique(stringList []string) bool {
	uniqueMap := make(map[string]bool)
	for _, s := range stringList {
		if _, ok := uniqueMap[s]; ok {
			return false
		}
	}
	return true
}

func (draMgr *DRAMgr) CreateDRAv6Interface(
	cfg *dhcprelayd.DHCPv6RelayIntf) (bool, error) {

	for _, ifRef := range cfg.UpstreamIntfs {
		_, ok := draMgr.IMgr.GetIPv6LLIntfIndex(ifRef)
		if !ok {
			errMsg := fmt.Sprintln(
				"DRA: No IPv6 link-local Intf exists with reference", ifRef)
			draMgr.Logger.Err(errMsg)
			return false, errors.New(errMsg)
		}
	}
	if !checkStringSliceUnique(cfg.UpstreamIntfs) {
		errMsg := fmt.Sprintln(
			"DRA: Non unique upstream interfaces")
		draMgr.Logger.Err(errMsg)
		return false, errors.New(errMsg)
	}
	for _, val := range cfg.ServerIp {
		ip := net.ParseIP(val)
		if ip == nil || ip.To4() != nil || ip.IsLinkLocalUnicast() {
			errMsg := fmt.Sprintln(
				"DRA: Parsing of ServerIp address failed:", val,
			)
			draMgr.Logger.Err(errMsg)
			return false, errors.New(errMsg)
		}
	}
	ifIdx, ok := draMgr.IMgr.GetIPv6IntfIndex(cfg.IntfRef)
	if !ok {
		errMsg := fmt.Sprintln(
			"DRA: No IPv6 Intf exists with reference", cfg.IntfRef)
		draMgr.Logger.Err(errMsg)
		return false, errors.New(errMsg)
	}
	draMgr.IMgr.UpdateDRAv6Intf(cfg)
	draMgr.PProc6.ProcessCreateDRAIntf(cfg.IntfRef)
	if _, ok := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx); ok {
		draMgr.PProc6.StartRxTx()
		draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) UpdateDRAv6Interface(oldCfg *dhcprelayd.DHCPv6RelayIntf,
	newCfg *dhcprelayd.DHCPv6RelayIntf, attrSet []bool,
	op []*dhcprelayd.PatchOpInfo) (bool, error) {

	cfg := oldCfg
	if attrSet[3] {
		for _, ifRef := range cfg.UpstreamIntfs {
			_, ok := draMgr.IMgr.GetIPv6LLIntfIndex(ifRef)
			if !ok {
				errMsg := fmt.Sprintln(
					"DRA: No IPv6 Intf link-local exists with reference", ifRef)
				draMgr.Logger.Err(errMsg)
				return false, errors.New(errMsg)
			}
		}
		cfg.UpstreamIntfs = newCfg.UpstreamIntfs
	}
	if attrSet[2] {
		for _, val := range newCfg.ServerIp {
			ip := net.ParseIP(val)
			if ip == nil || ip.To4() != nil || ip.IsLinkLocalUnicast() {
				errMsg := fmt.Sprintln(
					"DRA: Parsing of ServerIp address failed:", val)
				draMgr.Logger.Err(errMsg)
				return false, errors.New(errMsg)
			}
		}
		cfg.ServerIp = newCfg.ServerIp
	}
	if attrSet[1] {
		cfg.Enable = newCfg.Enable
	}
	ifIdx, ok := draMgr.IMgr.GetIPv6IntfIndex(newCfg.IntfRef)
	if !ok { // No valid interface
		draMgr.IMgr.UpdateDRAv6Intf(cfg)
		return true, nil
	}

	_, draPreState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)
	draMgr.IMgr.UpdateDRAv6Intf(cfg)
	_, draPostState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		draMgr.PProc6.StopRxTx()
		return true, nil
	}
	draMgr.PProc6.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc6.ProcessInactiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) DeleteDRAv6Interface(intfRef string) (bool, error) {
	ifIdx, ok := draMgr.IMgr.GetIPv6IntfIndex(intfRef)
	if !ok {
		draMgr.Logger.Info(
			"DRA: Cannot find interface with reference,", intfRef)
	}
	_, draPreState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)
	draMgr.IMgr.DeleteDRAv6Intf(intfRef)
	draMgr.PProc6.ProcessDeleteDRAIntf(intfRef)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		draMgr.PProc6.StopRxTx()
		return true, nil
	}
	if draPreState {
		draMgr.PProc6.ProcessInactiveDRAIntf(ifIdx)
	}
	return true, nil
}

func (draMgr *DRAMgr) ProcessAsicdNotification(
	msg commonDefs.AsicdNotifyMsg) {

	switch msg.(type) {
	case commonDefs.IPv4L3IntfStateNotifyMsg:
		stateMsg := msg.(commonDefs.IPv4L3IntfStateNotifyMsg)
		draMgr.processIPv4StateChange(stateMsg)
	case commonDefs.IPv6L3IntfStateNotifyMsg:
		stateMsg := msg.(commonDefs.IPv6L3IntfStateNotifyMsg)
		draMgr.processIPv6StateChange(stateMsg)
	case commonDefs.IPv4IntfNotifyMsg:
		ipv4Msg := msg.(commonDefs.IPv4IntfNotifyMsg)
		draMgr.processIPv4Notification(ipv4Msg)
	case commonDefs.IPv6IntfNotifyMsg:
		ipv6Msg := msg.(commonDefs.IPv6IntfNotifyMsg)
		draMgr.processIPv6Notification(ipv6Msg)
	}
}

func (draMgr *DRAMgr) processIPv4Notification(
	msg commonDefs.IPv4IntfNotifyMsg) {

	ifIdx := int(msg.IfIndex)
	_, draPreState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)
	if msg.MsgType == commonDefs.NOTIFY_IPV4INTF_CREATE {
		draMgr.IMgr.ProcessIPv4IntfCreate(msg)
	} else {
		draMgr.IMgr.ProcessIPv4IntfDelete(msg)
	}
	_, draPostState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)

	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		draMgr.PProc4.StopRxTx()
		return
	}
	draMgr.PProc4.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc4.ProcessInactiveDRAIntf(ifIdx)
	}
}

func (draMgr *DRAMgr) processIPv4StateChange(
	msg commonDefs.IPv4L3IntfStateNotifyMsg) {

	ifIdx := int(msg.IfIndex)
	_, draPreState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)
	draMgr.IMgr.ProcessIPv4StateChange(msg)
	_, draPostState := draMgr.IMgr.GetActiveDRAv4Intf(ifIdx)

	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		draMgr.PProc4.StopRxTx()
		return
	}
	draMgr.PProc4.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc4.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc4.ProcessInactiveDRAIntf(ifIdx)
	}
}

func (draMgr *DRAMgr) processIPv6Notification(
	msg commonDefs.IPv6IntfNotifyMsg) {

	ifIdx := int(msg.IfIndex)
	_, draPreState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)
	if msg.MsgType == commonDefs.NOTIFY_IPV6INTF_CREATE {
		draMgr.IMgr.ProcessIPv6IntfCreate(msg)
	} else {
		draMgr.IMgr.ProcessIPv6IntfDelete(msg)
	}
	_, draPostState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)

	// (TODO): Optimization check to ignore linklocal intf messages?
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		draMgr.PProc6.StopRxTx()
		return
	}
	draMgr.PProc6.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc6.ProcessInactiveDRAIntf(ifIdx)
	}
}

func (draMgr *DRAMgr) processIPv6StateChange(
	msg commonDefs.IPv6L3IntfStateNotifyMsg) {

	ifIdx := int(msg.IfIndex)
	_, draPreState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)
	draMgr.IMgr.ProcessIPv6StateChange(msg)
	_, draPostState := draMgr.IMgr.GetActiveDRAv6Intf(ifIdx)

	// (TODO): Optimization check to ignore linklocal intf messages?
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		draMgr.PProc6.StopRxTx()
		return
	}
	draMgr.PProc6.StartRxTx()
	if !draPreState && draPostState {
		draMgr.PProc6.ProcessActiveDRAIntf(ifIdx)
	} else if draPreState && !draPostState {
		draMgr.PProc6.ProcessInactiveDRAIntf(ifIdx)
	}
}

// State
func (draMgr *DRAMgr) GetDRAv4ClientState(macAddr string) (
	*dhcprelayd.DHCPRelayClientState, error) {

	if val, ok := draMgr.PProc4.GetClientState(macAddr); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv4ClientState(
	fromIdx, count int) *dhcprelayd.DHCPRelayClientStateGetInfo {

	result := &dhcprelayd.DHCPRelayClientStateGetInfo{}
	nextIdx, actualCount, more, clientStateSlice :=
		draMgr.PProc4.GetClientStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPRelayClientStateList = clientStateSlice
	return result
}

func (draMgr *DRAMgr) GetDRAv6ClientState(macAddr string) (
	*dhcprelayd.DHCPv6RelayClientState, error) {

	if val, ok := draMgr.PProc6.GetClientState(macAddr); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv6ClientState(
	fromIdx, count int) *dhcprelayd.DHCPv6RelayClientStateGetInfo {

	result := &dhcprelayd.DHCPv6RelayClientStateGetInfo{}
	nextIdx, actualCount, more, clientStateSlice :=
		draMgr.PProc6.GetClientStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPv6RelayClientStateList = clientStateSlice
	return result
}

func (draMgr *DRAMgr) GetDRAv4IntfState(intfRef string) (
	*dhcprelayd.DHCPRelayIntfState, error) {

	_, ok := draMgr.IMgr.GetIPv4IntfIndex(intfRef)
	if !ok {
		errMsg := fmt.Sprintln(
			"Could not find ipv4 interface", intfRef)
		draMgr.Logger.Err(errMsg)
		return nil, errors.New(errMsg)
	}

	if val, ok := draMgr.PProc4.GetIntfState(intfRef); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv4IntfState(
	fromIdx, count int) *dhcprelayd.DHCPRelayIntfStateGetInfo {

	result := &dhcprelayd.DHCPRelayIntfStateGetInfo{}
	nextIdx, actualCount, more, intfStateSlice :=
		draMgr.PProc4.GetIntfStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPRelayIntfStateList = intfStateSlice
	return result
}

func (draMgr *DRAMgr) GetDRAv6IntfState(intfRef string) (
	*dhcprelayd.DHCPv6RelayIntfState, error) {

	_, ok := draMgr.IMgr.GetIPv6IntfIndex(intfRef)
	if !ok {
		errMsg := fmt.Sprintln(
			"Could not find non link-local interface", intfRef)
		draMgr.Logger.Err(errMsg)
		return nil, errors.New(errMsg)
	}

	if val, ok := draMgr.PProc6.GetIntfState(intfRef); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv6IntfState(
	fromIdx, count int) *dhcprelayd.DHCPv6RelayIntfStateGetInfo {

	result := &dhcprelayd.DHCPv6RelayIntfStateGetInfo{}
	nextIdx, actualCount, more, intfStateSlice :=
		draMgr.PProc6.GetIntfStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPv6RelayIntfStateList = intfStateSlice
	return result
}

func (draMgr *DRAMgr) GetDRAv4IntfServerState(
	intfRef string, serverIp string) (
	*dhcprelayd.DHCPRelayIntfServerState, error) {

	if val, ok := draMgr.PProc4.GetIntfServerState(intfRef, serverIp); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv4IntfServerState(
	fromIdx, count int) *dhcprelayd.DHCPRelayIntfServerStateGetInfo {

	result := &dhcprelayd.DHCPRelayIntfServerStateGetInfo{}
	nextIdx, actualCount, more, intfServerStateSlice :=
		draMgr.PProc4.GetIntfServerStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPRelayIntfServerStateList = intfServerStateSlice
	return result
}

func (draMgr *DRAMgr) GetDRAv6IntfServerState(
	intfRef string, serverIp string) (
	*dhcprelayd.DHCPv6RelayIntfServerState, error) {

	if val, ok := draMgr.PProc6.GetIntfServerState(intfRef, serverIp); ok {
		return val, nil
	}
	return nil, errors.New("Could not find entry")
}

func (draMgr *DRAMgr) GetBulkDRAv6IntfServerState(
	fromIdx, count int) *dhcprelayd.DHCPv6RelayIntfServerStateGetInfo {

	result := &dhcprelayd.DHCPv6RelayIntfServerStateGetInfo{}
	nextIdx, actualCount, more, intfServerStateSlice :=
		draMgr.PProc6.GetIntfServerStateSlice(fromIdx, count)

	result.StartIdx = dhcprelayd.Int(fromIdx)
	result.EndIdx = dhcprelayd.Int(nextIdx)
	result.Count = dhcprelayd.Int(actualCount)
	result.More = more
	result.DHCPv6RelayIntfServerStateList = intfServerStateSlice
	return result
}
