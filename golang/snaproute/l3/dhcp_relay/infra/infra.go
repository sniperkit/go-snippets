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
package infra

import (
	"dhcprelayd"
	//	"errors"
	//	"fmt"
	"net"
	"sync"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/logging"
)

type IPv4IntfProperty struct {
	IpAddr  string
	Netmask net.IPMask
	IfIndex int
	IfRef   string
	State   bool
}

type IPv6IntfProperty struct {
	IpAddr  string
	Netmask net.IPMask
	IfIndex int
	IfRef   string
	State   bool
}

type InfraMgr struct {
	Logger        logging.LoggerIntf
	AsicdHdl      asicdClient.AsicdClientIntf
	InfraMgrMutex sync.Mutex

	IPv4IntfProps      map[int]*IPv4IntfProperty
	IPv4IfRefToIfIndex map[string]int
	IPv6IntfProps      map[int]*IPv6IntfProperty
	IPv6IfRefToIfIndex map[string]int
	// (NOTE): LL meaning link local
	IPv6LLIntfProps      map[int]*IPv6IntfProperty
	IPv6LLIfRefToIfIndex map[string]int

	DRAv4Global      *dhcprelayd.DHCPRelayGlobal
	DRAv4Intfs       map[string]*dhcprelayd.DHCPRelayIntf
	ActiveDRAv4Intfs map[int]*dhcprelayd.DHCPRelayIntf
	DRAv6Global      *dhcprelayd.DHCPv6RelayGlobal
	DRAv6Intfs       map[string]*dhcprelayd.DHCPv6RelayIntf
	ActiveDRAv6Intfs map[int]*dhcprelayd.DHCPv6RelayIntf
}

func NewInfraMgr(logger logging.LoggerIntf,
	asicdHdl asicdClient.AsicdClientIntf) *InfraMgr {

	iMgr := &InfraMgr{
		Logger:   logger,
		AsicdHdl: asicdHdl,
	}
	iMgr.IPv4IntfProps = make(map[int]*IPv4IntfProperty)
	iMgr.IPv4IfRefToIfIndex = make(map[string]int)
	iMgr.IPv6IntfProps = make(map[int]*IPv6IntfProperty)
	iMgr.IPv6IfRefToIfIndex = make(map[string]int)
	iMgr.IPv6LLIntfProps = make(map[int]*IPv6IntfProperty)
	iMgr.IPv6LLIfRefToIfIndex = make(map[string]int)
	iMgr.DRAv4Intfs = make(map[string]*dhcprelayd.DHCPRelayIntf)
	iMgr.ActiveDRAv4Intfs = make(map[int]*dhcprelayd.DHCPRelayIntf)
	iMgr.DRAv6Intfs = make(map[string]*dhcprelayd.DHCPv6RelayIntf)
	iMgr.ActiveDRAv6Intfs = make(map[int]*dhcprelayd.DHCPv6RelayIntf)
	return iMgr
}

func (iMgr *InfraMgr) BuildInfra() {
	iMgr.constructIPv4Infra()
	iMgr.constructIPv6Infra()
}

func getBinaryState(OperState string) bool {
	if OperState == "UP" {
		return true
	} else {
		return false
	}
}

func (iMgr *InfraMgr) constructIPv4Infra() {
	iMgr.Logger.Info("Calling Asicd for getting IPv4 Interfaces")
	intfs, err := iMgr.AsicdHdl.GetAllIPv4IntfState()
	if err != nil {
		return
	}
	for _, intf := range intfs {
		ip, ipNet, _ := net.ParseCIDR(intf.IpAddr)
		ifIdx := int(intf.IfIndex)
		ifRef := intf.IntfRef
		intfProp := &IPv4IntfProperty{
			IpAddr:  ip.String(),
			Netmask: ipNet.Mask,
			IfIndex: ifIdx,
			IfRef: ifRef,
			State:   getBinaryState(intf.OperState),
		}
		iMgr.IPv4IntfProps[ifIdx] = intfProp
		iMgr.IPv4IfRefToIfIndex[ifRef] = ifIdx
	}
}

func (iMgr *InfraMgr) constructIPv6Infra() {
	iMgr.Logger.Info("Calling Asicd for getting IPv6 Interfaces")
	intfs, err := iMgr.AsicdHdl.GetAllIPv6IntfState()
	if err != nil {
		return
	}
	for _, intf := range intfs {
		ip, ipNet, _ := net.ParseCIDR(intf.IpAddr)
		ifIdx := int(intf.IfIndex)
		ifRef := intf.IntfRef
		intfProp := &IPv6IntfProperty{
			IpAddr:  ip.String(),
			Netmask: ipNet.Mask,
			IfIndex: ifIdx,
			IfRef:   ifRef,
			State:   getBinaryState(intf.OperState),
		}
		if ip.IsLinkLocalUnicast() {
			iMgr.IPv6LLIntfProps[ifIdx] = intfProp
			iMgr.IPv6LLIfRefToIfIndex[ifRef] = ifIdx
		} else {
			iMgr.IPv6IntfProps[ifIdx] = intfProp
			iMgr.IPv6IfRefToIfIndex[ifRef] = ifIdx
		}
	}
}

func (iMgr *InfraMgr) GetIPv4Intf(ifIdx int) (*IPv4IntfProperty, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ipv4Intf, ok := iMgr.IPv4IntfProps[ifIdx]
	return ipv4Intf, ok
}

func (iMgr *InfraMgr) GetIPv4IntfIndex(intfRef string) (int, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifIdx, ok := iMgr.IPv4IfRefToIfIndex[intfRef]
	return ifIdx, ok
}

func (iMgr *InfraMgr) GetIPv6Intf(ifIdx int) (*IPv6IntfProperty, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ipv6Intf, ok := iMgr.IPv6IntfProps[ifIdx]
	return ipv6Intf, ok
}

func (iMgr *InfraMgr) GetIPv6IntfIndex(intfRef string) (int, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifIdx, ok := iMgr.IPv6IfRefToIfIndex[intfRef]
	return ifIdx, ok
}

func (iMgr *InfraMgr) GetIPv6LLIntf(ifIdx int) (*IPv6IntfProperty, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ipv6Intf, ok := iMgr.IPv6LLIntfProps[ifIdx]
	return ipv6Intf, ok
}

func (iMgr *InfraMgr) GetIPv6LLIntfIndex(intfRef string) (int, bool) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifIdx, ok := iMgr.IPv6LLIfRefToIfIndex[intfRef]
	return ifIdx, ok
}

// Get from Asicd
// func (iMgr *InfraMgr) asicdGetIPv6State(ifIdx int) (*IPv6IntfProperty, error) {
// 	intfs, err := iMgr.AsicdHdl.GetAllIPv6IntfState()
// 	if err != nil {
// 		return nil, errors.New("Cannot get IPv6 States from asicd")
// 	}

// 	for _, intf := range intfs {
// 		if int(intf.IfIndex) == ifIdx {
// 			ip, ipNet, _ := net.ParseCIDR(intf.IpAddr)
// 			if ip.IsLinkLocalUnicast() {
// 				continue
// 			}
// 			intfProp := &IPv6IntfProperty{
// 				IpAddr:  ip.String(),
// 				Netmask: ipNet.Mask,
// 				IfIndex: ifIdx,
// 				IfRef:   intf.IntfRef,
// 				State:   getBinaryState(intf.OperState),
// 			}
// 			return intfProp, nil
// 		}
// 	}
// 	return nil, errors.New(fmt.Sprintln("Cannot get IPv6 State for intf", ifIdx))
// }
func (iMgr *InfraMgr) ProcessIPv4IntfCreate(
	msg commonDefs.IPv4IntfNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ip, ipNet, _ := net.ParseCIDR(msg.IpAddr)
	ifIdx := int(msg.IfIndex)
	intfProp := &IPv4IntfProperty{
		IpAddr:  ip.String(),
		Netmask: ipNet.Mask,
		IfIndex: ifIdx,
		IfRef: msg.IntfRef,
		State:   false, //Correct assumption?
	}
	iMgr.Logger.Info("In ProcessIPv4IntfCreate")
	iMgr.IPv4IntfProps[ifIdx] = intfProp
	iMgr.IPv4IfRefToIfIndex[intfProp.IfRef] = ifIdx
	if !intfProp.State {
		return
	}
	if iMgr.DRAv4Global == nil || !iMgr.DRAv4Global.Enable {
		return
	}
	draIntf, ok := iMgr.DRAv4Intfs[intfProp.IfRef]
	if !ok || !draIntf.Enable {
		return
	}
	iMgr.ActiveDRAv4Intfs[ifIdx] = draIntf
}

func (iMgr *InfraMgr) ProcessIPv6IntfCreate(
	msg commonDefs.IPv6IntfNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ip, ipNet, _ := net.ParseCIDR(msg.IpAddr)
	ifIdx := int(msg.IfIndex)
	intfProp := &IPv6IntfProperty{
		IpAddr:  ip.String(),
		Netmask: ipNet.Mask,
		IfIndex: ifIdx,
		IfRef:   msg.IntfRef,
		State:   false, //Correct assumption?
	}
	if ip.IsLinkLocalUnicast() {
		iMgr.IPv6LLIntfProps[ifIdx] = intfProp
		iMgr.IPv6LLIfRefToIfIndex[intfProp.IfRef] = ifIdx
		return
	}
	iMgr.IPv6IntfProps[ifIdx] = intfProp
	iMgr.IPv6IfRefToIfIndex[intfProp.IfRef] = ifIdx
	if !intfProp.State {
		return
	}
	if iMgr.DRAv6Global == nil || !iMgr.DRAv6Global.Enable {
		return
	}
	draIntf, ok := iMgr.DRAv6Intfs[intfProp.IfRef]
	if !ok || !draIntf.Enable {
		return
	}
	iMgr.ActiveDRAv6Intfs[ifIdx] = draIntf
}

func (iMgr *InfraMgr) ProcessIPv4IntfDelete(
	msg commonDefs.IPv4IntfNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifIdx := int(msg.IfIndex)
	ipv4Intf, ok := iMgr.IPv4IntfProps[ifIdx]
	if !ok {
		iMgr.Logger.Info("No IPv4 Intf found for index", ifIdx)
		return
	}
	delete(iMgr.IPv4IfRefToIfIndex, ipv4Intf.IfRef)
	delete(iMgr.IPv4IntfProps, ifIdx)
	delete(iMgr.ActiveDRAv4Intfs, ifIdx)
}

func (iMgr *InfraMgr) ProcessIPv6IntfDelete(
	msg commonDefs.IPv6IntfNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ip, _, _ := net.ParseCIDR(msg.IpAddr)
	ifIdx := int(msg.IfIndex)
	if ip.IsLinkLocalUnicast() {
		ipv6Intf, ok := iMgr.IPv6LLIntfProps[ifIdx]
		if !ok {
			iMgr.Logger.Info("No IPv6 Intf found for index", ifIdx)
			return
		}
		delete(iMgr.IPv6LLIfRefToIfIndex, ipv6Intf.IfRef)
		delete(iMgr.IPv6LLIntfProps, ifIdx)
		return
	}
	ipv6Intf, ok := iMgr.IPv6IntfProps[ifIdx]
	if !ok {
		iMgr.Logger.Info("No IPv6 Intf found for index", ifIdx)
		return
	}
	delete(iMgr.IPv6IfRefToIfIndex, ipv6Intf.IfRef)
	delete(iMgr.IPv6IntfProps, ifIdx)
	delete(iMgr.ActiveDRAv6Intfs, ifIdx)
}

func (iMgr *InfraMgr) ProcessIPv4StateChange(
	msg commonDefs.IPv4L3IntfStateNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifIdx := int(msg.IfIndex)
	ipv4Intf, ok := iMgr.IPv4IntfProps[ifIdx]
	if !ok {
		iMgr.Logger.Info("No IPv4 Intf found for index", ifIdx)
		return
	}
	if msg.IfState == 0 {
		ipv4Intf.State = false
	} else {
		ipv4Intf.State = true
	}
	if iMgr.DRAv4Global == nil || !iMgr.DRAv4Global.Enable {
		return
	}
	draIntf, ok := iMgr.DRAv4Intfs[ipv4Intf.IfRef]
	if !ok || !draIntf.Enable {
		return
	}
	if msg.IfState == 0 { // State going DOWN
		delete(iMgr.ActiveDRAv4Intfs, ifIdx)
	} else { // State going UP
		iMgr.ActiveDRAv4Intfs[ifIdx] = draIntf
	}
}

func (iMgr *InfraMgr) ProcessIPv6StateChange(
	msg commonDefs.IPv6L3IntfStateNotifyMsg) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ip, _, _ := net.ParseCIDR(msg.IpAddr)
	ifIdx := int(msg.IfIndex)
	if ip.IsLinkLocalUnicast() {
		ipv6Intf, ok := iMgr.IPv6LLIntfProps[ifIdx]
		if !ok {
			iMgr.Logger.Info("No IPv6 Intf found for index", ifIdx)
			return
		}
		if msg.IfState == 0 {
			ipv6Intf.State = false
		} else {
			ipv6Intf.State = true
		}
	} else {
		ipv6Intf, ok := iMgr.IPv6IntfProps[ifIdx]
		if !ok {
			iMgr.Logger.Info("No IPv6 Intf found for index", ifIdx)
			return
		}
		if msg.IfState == 0 {
			ipv6Intf.State = false
		} else {
			ipv6Intf.State = true
		}
		if iMgr.DRAv6Global == nil || !iMgr.DRAv6Global.Enable {
			return
		}
		draIntf, ok := iMgr.DRAv6Intfs[ipv6Intf.IfRef]
		if !ok || !draIntf.Enable {
			return
		}
		if msg.IfState == 0 { // State going DOWN
			delete(iMgr.ActiveDRAv6Intfs, ifIdx)
		} else { // State going UP
			iMgr.ActiveDRAv6Intfs[ifIdx] = draIntf
		}
	}
}

func (iMgr *InfraMgr) GetActiveDRAv4IntfCount() int {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	return len(iMgr.ActiveDRAv4Intfs)
}

func (iMgr *InfraMgr) GetActiveDRAv6IntfCount() int {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	return len(iMgr.ActiveDRAv6Intfs)
}

func (iMgr *InfraMgr) GetActiveDRAv4Intf(
	ifIdx int) (*dhcprelayd.DHCPRelayIntf, bool) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	draIntf, ok := iMgr.ActiveDRAv4Intfs[ifIdx]
	return draIntf, ok
}

func (iMgr *InfraMgr) GetActiveDRAv6Intf(
	ifIdx int) (*dhcprelayd.DHCPv6RelayIntf, bool) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	draIntf, ok := iMgr.ActiveDRAv6Intfs[ifIdx]
	return draIntf, ok
}

func (iMgr *InfraMgr) GetAllActiveDRAv4Intfs() []int {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	activeIfIdxs := make([]int, 0, len(iMgr.ActiveDRAv4Intfs))
	for idx := range iMgr.ActiveDRAv4Intfs {
		activeIfIdxs = append(activeIfIdxs, idx)
	}
	return activeIfIdxs
}

func (iMgr *InfraMgr) GetAllActiveDRAv6Intfs() []int {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	activeIfIdxs := make([]int, 0, len(iMgr.ActiveDRAv6Intfs))
	for idx := range iMgr.ActiveDRAv6Intfs {
		activeIfIdxs = append(activeIfIdxs, idx)
	}
	return activeIfIdxs
}

func (iMgr *InfraMgr) UpdateDRAv4Global(cfg *dhcprelayd.DHCPRelayGlobal) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	iMgr.DRAv4Global = cfg
	if !iMgr.DRAv4Global.Enable {
		for ifIdx, _ := range iMgr.ActiveDRAv4Intfs {
			delete(iMgr.ActiveDRAv4Intfs, ifIdx)
		}
		return
	}
	for ifRef, draIntf := range iMgr.DRAv4Intfs {
		ifIdx, ok := iMgr.IPv4IfRefToIfIndex[ifRef]
		if !ok {
			continue
		}
		ipv4Intf, _ := iMgr.IPv4IntfProps[ifIdx]
		if draIntf.Enable && ipv4Intf.State {
			iMgr.ActiveDRAv4Intfs[ifIdx] = draIntf
		}
	}
}

func (iMgr *InfraMgr) UpdateDRAv6Global(cfg *dhcprelayd.DHCPv6RelayGlobal) {

	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	iMgr.DRAv6Global = cfg
	if !iMgr.DRAv6Global.Enable {
		for ifIdx, _ := range iMgr.ActiveDRAv6Intfs {
			delete(iMgr.ActiveDRAv6Intfs, ifIdx)
		}
		return
	}
	for ifRef, draIntf := range iMgr.DRAv6Intfs {
		ifIdx, ok := iMgr.IPv6IfRefToIfIndex[ifRef]
		if !ok {
			continue
		}
		ipv6Intf, _ := iMgr.IPv6IntfProps[ifIdx]
		if draIntf.Enable && ipv6Intf.State {
			iMgr.ActiveDRAv6Intfs[ifIdx] = draIntf
		}
	}
}

func (iMgr *InfraMgr) DeleteDRAv4Global() {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	for ifIdx := range iMgr.ActiveDRAv4Intfs {
		delete(iMgr.ActiveDRAv4Intfs, ifIdx)
	}
	iMgr.DRAv4Global = nil
}

func (iMgr *InfraMgr) DeleteDRAv6Global() {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	for ifIdx := range iMgr.ActiveDRAv6Intfs {
		delete(iMgr.ActiveDRAv6Intfs, ifIdx)
	}
	iMgr.DRAv6Global = nil
}

func (iMgr *InfraMgr) UpdateDRAv4Intf(cfg *dhcprelayd.DHCPRelayIntf) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifRef := cfg.IntfRef
	iMgr.DRAv4Intfs[ifRef] = cfg

	if iMgr.DRAv4Global == nil || !iMgr.DRAv4Global.Enable {
		return
	}
	ifIdx, ok := iMgr.IPv4IfRefToIfIndex[ifRef]
	if !ok {
		iMgr.Logger.Info("No IPv4 IfIndex found for IntfRef", ifRef)
		return
	}
	ipv4Intf, _ := iMgr.IPv4IntfProps[ifIdx]
	if cfg.Enable && ipv4Intf.State {
		iMgr.ActiveDRAv4Intfs[ifIdx] = cfg
	} else {
		delete(iMgr.ActiveDRAv4Intfs, ifIdx)
	}
}

func (iMgr *InfraMgr) UpdateDRAv6Intf(cfg *dhcprelayd.DHCPv6RelayIntf) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	ifRef := cfg.IntfRef
	iMgr.DRAv6Intfs[ifRef] = cfg

	if iMgr.DRAv6Global == nil || !iMgr.DRAv6Global.Enable {
		return
	}
	ifIdx, ok := iMgr.IPv6IfRefToIfIndex[ifRef]
	if !ok {
		iMgr.Logger.Info("No IPv6 IfIndex found for IntfRef", ifRef)
		return
	}
	ipv6Intf, _ := iMgr.IPv6IntfProps[ifIdx]
	if cfg.Enable && ipv6Intf.State {
		iMgr.ActiveDRAv6Intfs[ifIdx] = cfg
	} else {
		delete(iMgr.ActiveDRAv6Intfs, ifIdx)
	}
}

func (iMgr *InfraMgr) DeleteDRAv4Intf(ifRef string) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	delete(iMgr.DRAv4Intfs, ifRef)
	if ifIdx, ok := iMgr.IPv4IfRefToIfIndex[ifRef]; ok {
		delete(iMgr.ActiveDRAv4Intfs, ifIdx)
	}
}

func (iMgr *InfraMgr) DeleteDRAv6Intf(ifRef string) {
	defer iMgr.InfraMgrMutex.Unlock()
	iMgr.InfraMgrMutex.Lock()

	delete(iMgr.DRAv6Intfs, ifRef)
	if ifIdx, ok := iMgr.IPv6IfRefToIfIndex[ifRef]; ok {
		delete(iMgr.ActiveDRAv6Intfs, ifIdx)
	}
}
