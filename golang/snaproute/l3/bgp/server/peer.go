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

// peer.go
package server

import (
	_ "fmt"
	"l3/bgp/baseobjects"
	"l3/bgp/config"
	"l3/bgp/fsm"
	"l3/bgp/packet"
	bgppolicy "l3/bgp/policy"
	bgprib "l3/bgp/rib"
	"net"
	"runtime"
	"strings"
	"sync/atomic"
	"utils/logging"
	"utils/patriciaDB"
	utilspolicy "utils/policy"
	"utils/policy/policyCommonDefs"
)

const (
	Reject int = iota
	Accept
)

type AdjRIBPolicyParams struct {
	CreateType      int
	DeleteType      int
	Route           *bgprib.AdjRIBPathIdRoute
	Peer            *Peer
	Accept          int
	PolicyEngine    *bgppolicy.AdjRibPPolicyEngine
	PolicyStmt      utilspolicy.PolicyStmt
	updated         *(map[uint32]map[*bgprib.Path][]*bgprib.Destination)
	withdrawn       *([]*bgprib.Destination)
	updatedAddPaths *([]*bgprib.Destination)
}

type Peer struct {
	server       *BGPServer
	logger       *logging.Writer
	locRib       *bgprib.LocRib
	NeighborConf *base.NeighborConf
	fsmManager   *fsm.FSMManager
	active       bool
	ifIdx        int32
	ribIn        map[uint32]map[string]*bgprib.AdjRIBRoute
	ribOut       map[uint32]map[string]*bgprib.AdjRIBRoute
}

func NewPeer(server *BGPServer, locRib *bgprib.LocRib, globalConf *config.GlobalConfig,
	peerGroup *config.PeerGroupConfig, peerConf config.NeighborConfig) *Peer {
	server.logger.Info("NewPeer - ip:", peerConf.NeighborAddress, "ifIndex:", peerConf.IfIndex)

	peer := Peer{
		server: server,
		logger: server.logger,
		locRib: locRib,
		active: false,
		ifIdx:  -1,
		ribIn:  make(map[uint32]map[string]*bgprib.AdjRIBRoute),
		ribOut: make(map[uint32]map[string]*bgprib.AdjRIBRoute),
	}

	peer.NeighborConf = base.NewNeighborConf(peer.logger, globalConf, peerGroup, peerConf)

	if !peer.IsConfigured() {
		peer.logger.Infof("NewPeer - Neighbor is not ready to be started, ip:",
			peer.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", peer.NeighborConf.Neighbor.Config.IfIndex)
		return &peer
	}

	if peer.IsDisabled() {
		peer.logger.Info("NewPeer - Neighbor is disabled, ip:", peer.NeighborConf.Neighbor.NeighborAddress,
			"ifIndex:", peer.NeighborConf.Neighbor.Config.IfIndex)
		return &peer
	}

	peer.fsmManager = fsm.NewFSMManager(peer.logger, peer.NeighborConf, server.BGPPktSrcCh,
		server.PeerFSMConnCh, server.ReachabilityCh)
	return &peer
}

func (p *Peer) UpdateGlobal(gConf *config.GlobalConfig) {
	p.NeighborConf.UpdateGlobal(gConf)
}

func (p *Peer) UpdatePeerGroup(peerGroup *config.PeerGroupConfig) {
	p.NeighborConf.UpdatePeerGroup(peerGroup)
}

func (p *Peer) UpdateNeighborConf(nConf config.NeighborConfig, bgp *config.Bgp) {
	p.NeighborConf.UpdateNeighborConf(nConf, bgp)
}

func (p *Peer) initAdjRIBTables() {
	for protoFamily, ok := range p.NeighborConf.AfiSafiMap {
		if ok {
			p.ribIn[protoFamily] = make(map[string]*bgprib.AdjRIBRoute)
			p.ribOut[protoFamily] = make(map[string]*bgprib.AdjRIBRoute)
		}
	}
}

func (p *Peer) IsBfdStateUp() bool {
	up := true
	if p.NeighborConf.Neighbor.State.UseBfdState {
		if p.NeighborConf.RunningConf.BfdEnable &&
			p.NeighborConf.Neighbor.State.BfdNeighborState == "down" {
			p.logger.Infof("Neighbor's bfd state is down for %s", p.NeighborConf.Neighbor.NeighborAddress)
			up = false
		}
	}
	return up
}

func (p *Peer) GetActionType(adjRIBDir bgprib.AdjRIBDir) (int, bool) {
	switch adjRIBDir {
	case bgprib.AdjRIBDirIn:
		return policyCommonDefs.PolicyActionTypeRIBIn, true

	case bgprib.AdjRIBDirOut:
		return policyCommonDefs.PolicyActionTypeRIBOut, true
	}
	return -1, false
}

func (p *Peer) RemoveAdjRIBFilter(pe *bgppolicy.AdjRibPPolicyEngine, policyName string, adjRIBDir bgprib.AdjRIBDir) {
	p.logger.Debug("RemoveAdjRIBFilter")
	if !p.IsConfigured() {
		p.logger.Infof("RemoveAdjRIBFilter - Neighbor is not ready to be started, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if p.IsDisabled() {
		p.logger.Info("RemoveAdjRIBFilter - Neighbor is disabled, ip:", p.NeighborConf.Neighbor.NeighborAddress,
			"ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	policyEngine := pe.GetPolicyEngine()
	policyDB := policyEngine.PolicyDB
	//neighborIP := p.NeighborConf.RunningConf.NeighborAddress.String()

	nodeGet := policyDB.Get(patriciaDB.Prefix(policyName))
	if nodeGet == nil {
		p.logger.Err("RemoveAdjRIBFilter - Policy ", policyName, " not created yet")
		return
	}
	node := nodeGet.(utilspolicy.Policy)
	actionType, ok := p.GetActionType(adjRIBDir)
	if !ok {
		p.logger.Err("RemoveAdjRIBFilter - Action type not found for Adj RIB direction", adjRIBDir)
		return
	}

	neighborIP := p.NeighborConf.RunningConf.NeighborAddress.String()
	conditionNameList := make([]string, 1)
	conditionNameList[0] = neighborIP

	policyAction := utilspolicy.PolicyAction{
		Name:       neighborIP,
		ActionType: actionType,
	}

	p.logger.Debug("RemoveAdjRIBFilter - Calling applypolicy with conditionNameList: ", conditionNameList)
	pe.UpdateUndoApplyPolicy(utilspolicy.ApplyPolicyInfo{node, policyAction, conditionNameList}, true)
}

func (p *Peer) AddAdjRIBFilter(pe *bgppolicy.AdjRibPPolicyEngine, policyName string, adjRIBDir bgprib.AdjRIBDir) {
	p.logger.Debug("AddAdjRIBFilter")
	if !p.IsConfigured() {
		p.logger.Infof("AddAdjRIBFilter - Neighbor is not ready to be started, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if p.IsDisabled() {
		p.logger.Info("AddAdjRIBFilter - Neighbor is disabled, ip:", p.NeighborConf.Neighbor.NeighborAddress,
			"ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	policyEngine := pe.GetPolicyEngine()
	policyDB := policyEngine.PolicyDB
	nodeGet := policyDB.Get(patriciaDB.Prefix(policyName))
	if nodeGet == nil {
		p.logger.Err("AddAdjRIBFilter - Policy ", policyName, " not defined")
		return
	}
	node := nodeGet.(utilspolicy.Policy)

	actionType, ok := p.GetActionType(adjRIBDir)
	if !ok {
		p.logger.Err("AddAdjRIBFilter - Action type not found for Adj RIB direction", adjRIBDir)
		return
	}

	neighborIP := p.NeighborConf.RunningConf.NeighborAddress.String()
	cond := utilspolicy.PolicyConditionConfig{
		Name:                       neighborIP,
		ConditionType:              "MatchNeighbor",
		MatchNeighborConditionInfo: neighborIP,
	}

	_, err := pe.CreatePolicyCondition(cond)
	if err != nil {
		p.logger.Errf("AddAdjRIBFilter - Failed to create policy condition to match neighbor %s with error %s",
			neighborIP, err)
		return
	}

	conditionNameList := make([]string, 1)
	conditionNameList[0] = neighborIP

	policyAction := utilspolicy.PolicyAction{
		Name:       neighborIP,
		ActionType: actionType,
	}

	p.logger.Debug("AddAdjRIBFilter - Calling applypolicy with conditionNameList: ", conditionNameList)
	pe.UpdateApplyPolicy(utilspolicy.ApplyPolicyInfo{node, policyAction, conditionNameList}, true)
}

func (p *Peer) IsConfigured() bool {
	return p.NeighborConf.RunningConf.NeighborAddress != nil
}

func (p *Peer) IsDisabled() bool {
	return p.NeighborConf.RunningConf.Disabled
}

func (p *Peer) IsActive() bool {
	return p.active
}

func (p *Peer) SetNeighborAddress(ip net.IP) {
	p.NeighborConf.SetNeighborAddress(ip)
}

func (p *Peer) ResetNeighborAddress() {
	p.NeighborConf.ResetNeighborAddress()
}

func (p *Peer) Init() {
	var fsmMgr *fsm.FSMManager

	if p.IsActive() {
		p.logger.Info("Init - Neighbor is already active, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if !p.IsConfigured() {
		p.logger.Info("Init - Neighbor is not ready to be started, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if p.IsDisabled() {
		p.logger.Info("Init - Neighbor is disabled, ip:", p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:",
			p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	p.logger.Debug("Init - adjribinfilter:", p.NeighborConf.RunningConf.AdjRIBInFilter, "adjriboutfilter:",
		p.NeighborConf.RunningConf.AdjRIBOutFilter)
	if p.NeighborConf.RunningConf.AdjRIBInFilter != "" {
		p.AddAdjRIBFilter(p.server.ribInPE, p.NeighborConf.RunningConf.AdjRIBInFilter, bgprib.AdjRIBDirIn)
	}

	if p.NeighborConf.RunningConf.AdjRIBOutFilter != "" {
		p.AddAdjRIBFilter(p.server.ribOutPE, p.NeighborConf.RunningConf.AdjRIBOutFilter, bgprib.AdjRIBDirOut)
	}

	if p.fsmManager == nil {
		p.logger.Infof("Init - Instantiating new FSM Manager for neighbor %s", p.NeighborConf.Neighbor.NeighborAddress)
		fsmMgr = fsm.NewFSMManager(p.logger, p.NeighborConf, p.server.BGPPktSrcCh,
			p.server.PeerFSMConnCh, p.server.ReachabilityCh)
	} else {
		fsmMgr = p.fsmManager
	}

	p.clearRibOut()
	go fsmMgr.Init()
	runtime.Gosched()

	p.fsmManager = fsmMgr
	p.active = true
	p.ProcessBfd(true)
}

func (p *Peer) Cleanup() {
	if !p.IsActive() {
		p.logger.Info("Cleanup - Neighbor is not active, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if !p.IsConfigured() {
		p.logger.Infof("Cleanup - Neighbor is not started yet, ip:",
			p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:", p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	if p.IsDisabled() {
		p.logger.Info("Cleanup - Neighbor is disabled, ip:", p.NeighborConf.Neighbor.NeighborAddress, "ifIndex:",
			p.NeighborConf.Neighbor.Config.IfIndex)
		return
	}

	p.active = false

	if p.NeighborConf.RunningConf.AdjRIBInFilter != "" {
		p.RemoveAdjRIBFilter(p.server.ribInPE, p.NeighborConf.RunningConf.AdjRIBInFilter, bgprib.AdjRIBDirIn)
	}

	if p.NeighborConf.RunningConf.AdjRIBOutFilter != "" {
		p.RemoveAdjRIBFilter(p.server.ribOutPE, p.NeighborConf.RunningConf.AdjRIBOutFilter, bgprib.AdjRIBDirOut)
	}

	p.ProcessBfd(false)

	if p.fsmManager == nil {
		p.logger.Errf("Can't cleanup FSM, FSM Manager is not instantiated for neighbor %s",
			p.NeighborConf.Neighbor.NeighborAddress)
		return
	}

	fsmMgr := p.fsmManager
	p.fsmManager = nil
	fsmMgr.CloseCh <- true
}

func (p *Peer) StopFSM(msg string) {
	if p.fsmManager == nil {
		p.logger.Errf("Can't stop FSM, FSM Manager is not instantiated for neighbor %s",
			p.NeighborConf.Neighbor.NeighborAddress)
		return
	}

	p.fsmManager.StopFSMCh <- msg
}

func (p *Peer) MaxPrefixesExceeded() {
	if p.NeighborConf.RunningConf.MaxPrefixesDisconnect {
		p.Command(int(fsm.BGPEventAutoStop), fsm.BGPCmdReasonMaxPrefixExceeded)
	}
}

func (p *Peer) setIfIdx(ifIdx int32) {
	p.ifIdx = ifIdx
}

func (p *Peer) getIfIdx() int32 {
	return p.ifIdx
}

func (p *Peer) AcceptConn(conn *net.TCPConn) {
	if p.fsmManager == nil {
		p.logger.Errf("FSM Manager is not instantiated yet for neighbor %s",
			p.NeighborConf.Neighbor.NeighborAddress)
		conn.Close()
		return
	}
	p.fsmManager.AcceptCh <- conn
}

func (p *Peer) Command(command int, reason int) {
	if p.fsmManager == nil {
		p.logger.Errf("FSM Manager is not instantiated yet for neighbor %s",
			p.NeighborConf.Neighbor.NeighborAddress)
		return
	}
	p.fsmManager.CommandCh <- fsm.PeerFSMCommand{command, reason}
}

func (p *Peer) BfdFaultSet() {
	p.NeighborConf.BfdFaultSet()
	p.fsmManager.BfdStatusCh <- false
}

func (p *Peer) BfdFaultCleared() {
	p.NeighborConf.BfdFaultCleared()
	p.fsmManager.BfdStatusCh <- true
}

func (p *Peer) getAddPathsMaxTx() int {
	return int(p.NeighborConf.Neighbor.State.AddPathsMaxTx)
}

func (p *Peer) clearRibOut() {
	p.ribIn = nil
	p.ribOut = nil
	p.ribIn = make(map[uint32]map[string]*bgprib.AdjRIBRoute)
	p.ribOut = make(map[uint32]map[string]*bgprib.AdjRIBRoute)
	p.initAdjRIBTables()
}

func (p *Peer) ProcessBfd(add bool) {
	ipAddr := p.NeighborConf.Neighbor.NeighborAddress.String()
	iface := p.NeighborConf.RunningConf.IfName
	sessionParam := p.NeighborConf.RunningConf.BfdSessionParam
	if add && p.NeighborConf.RunningConf.BfdEnable {
		p.logger.Info("Bfd enabled on", p.NeighborConf.Neighbor.NeighborAddress)
		ret, err := p.server.bfdMgr.CreateBfdSession(ipAddr, iface, sessionParam)
		if !ret {
			p.logger.Info("BfdSessionConfig FAILED, ret:", ret, "err:", err)
		} else {
			p.logger.Info("Bfd session configured: ", ipAddr, " param: ", sessionParam)
			p.NeighborConf.Neighbor.State.BfdNeighborState = "up"
		}
	} else {
		if p.NeighborConf.Neighbor.State.BfdNeighborState != "" {
			p.logger.Info("Bfd disabled on", p.NeighborConf.Neighbor.NeighborAddress)
			ret, err := p.server.bfdMgr.DeleteBfdSession(ipAddr, iface)
			if !ret {
				p.logger.Info("BfdSessionConfig FAILED, ret:", ret, "err:", err)
			} else {
				p.logger.Info("Bfd session removed for", p.NeighborConf.Neighbor.NeighborAddress)
				p.NeighborConf.Neighbor.State.BfdNeighborState = ""
			}
		}
	}

}

func (p *Peer) PeerConnEstablished(conn *net.Conn) {
	host, _, err := net.SplitHostPort((*conn).LocalAddr().String())
	if err != nil {
		p.logger.Errf("Neighbor %s: Can't find local address from the peer connection: %s",
			p.NeighborConf.Neighbor.NeighborAddress, (*conn).LocalAddr())
		return
	}
	p.logger.Infof("Neighbor %s: Local address %s for conn", p.NeighborConf.Neighbor.NeighborAddress, host)
	hostSplit := strings.Split(host, "%")
	host = hostSplit[0]
	p.NeighborConf.Neighbor.Transport.Config.LocalAddress = net.ParseIP(host)
	p.NeighborConf.PeerConnEstablished()
	p.clearRibOut()
	//p.Server.PeerConnEstCh <- p.Neighbor.NeighborAddress.String()
}

func (p *Peer) PeerConnBroken(fsmCleanup bool) {
	if p.NeighborConf.Neighbor.Transport.Config.LocalAddress != nil {
		p.NeighborConf.Neighbor.Transport.Config.LocalAddress = nil
		//p.Server.PeerConnBrokenCh <- p.Neighbor.NeighborAddress.String()
	}
	p.NeighborConf.PeerConnBroken()
	p.clearRibOut()
}

func (p *Peer) GetAdjRIB(adjRIBDir bgprib.AdjRIBDir) map[uint32]map[string]*bgprib.AdjRIBRoute {
	switch adjRIBDir {
	case bgprib.AdjRIBDirIn:
		return p.ribIn

	case bgprib.AdjRIBDirOut:
		return p.ribOut
	}
	return nil
}

func (p *Peer) processWithdraws(protoFamily uint32, nlris *[]packet.NLRI) {
	var route *bgprib.AdjRIBRoute
	var pathIdRoute *bgprib.AdjRIBPathIdRoute
	var ok bool
	total := len(*nlris)
	last := total - 1
	idx := 0
	for i := 0; i < total; i++ {
		nlri := (*nlris)[idx]
		if nlri == nil {
			if idx >= last {
				break
			}
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}

		ip := nlri.GetCIDR()
		p.logger.Infof("Neighbor %s: Withdraw Prefix %s protocol family=%d RIB-In=%+v",
			p.NeighborConf.Neighbor.NeighborAddress, ip, protoFamily, p.ribIn[protoFamily])
		if route, ok = p.ribIn[protoFamily][ip]; !ok {
			p.logger.Errf("Neighbor %s: Withdraw Prefix %s not found in RIB-In, protocol family=%d",
				p.NeighborConf.Neighbor.NeighborAddress, ip, protoFamily)
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}

		if pathIdRoute = route.GetPathIdRoute(nlri.GetPathId()); pathIdRoute == nil {
			p.logger.Errf("Neighbor %s: Withdraw Prefix %s Path id %d not found in RIB-In, protocol family=%d",
				p.NeighborConf.Neighbor.NeighborAddress, ip, nlri.GetPathId(), protoFamily)
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}

		p.logger.Infof("Neighbor %s: Remove path id %d for nlri %s protocol family %d from RIB-In",
			p.NeighborConf.RunningConf.NeighborAddress, nlri.GetPathId(), ip, protoFamily)
		route.RemovePath(nlri.GetPathId())
		if !route.DoesPathsExist() {
			p.logger.Infof("Neighbor %s: remove nlri %s protocol family %s from RIB-In",
				p.NeighborConf.RunningConf.NeighborAddress, ip, protoFamily)
			peEntity := bgppolicy.GetPolicyEngineFilterEntity(pathIdRoute.Path)
			p.checkRIBInFilter(nlri, pathIdRoute, peEntity, false)
			delete(p.ribIn[protoFamily], ip)
		}

		if !pathIdRoute.Accept {
			p.logger.Errf("Neighbor %s: Withdraw Prefix %s Path id %d was already filtered out",
				p.NeighborConf.Neighbor.NeighborAddress, ip, nlri.GetPathId())
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}

		idx++
	}
	(*nlris) = (*nlris)[:idx]
}

func (p *Peer) checkAdjRIBFilter(nlri packet.NLRI, route *bgprib.AdjRIBPathIdRoute,
	peEntity *utilspolicy.PolicyEngineFilterEntityParams, pe *bgppolicy.AdjRibPPolicyEngine,
	policyDir int, create bool) (bool, *AdjRIBPolicyParams) {
	if route != nil {
		peEntity.DestNetIp = route.NLRI.GetCIDR()
		peEntity.Neighbor = p.NeighborConf.RunningConf.NeighborAddress.String()

		callbackInfo := &AdjRIBPolicyParams{
			Peer:  p,
			Route: route,
		}

		if create {
			peEntity.CreatePath = true
			callbackInfo.CreateType = utilspolicy.Valid
			callbackInfo.DeleteType = utilspolicy.Invalid
		} else {
			peEntity.DeletePath = true
			callbackInfo.CreateType = utilspolicy.Invalid
			callbackInfo.DeleteType = utilspolicy.Valid
		}

		pe.PolicyEngine.PolicyEngineFilter(*peEntity, policyDir, callbackInfo)
		p.logger.Infof("checkAdjRIBFilter - NLRI %s policylist %v hit %v after applying create policy, callbackInfo=%+v",
			nlri.GetCIDR(), route.PolicyList, route.PolicyHitCounter, callbackInfo)
		return callbackInfo.Accept == Accept, callbackInfo
	}
	return false, nil
}

func (p *Peer) checkRIBInFilter(nlri packet.NLRI, route *bgprib.AdjRIBPathIdRoute,
	peEntity *utilspolicy.PolicyEngineFilterEntityParams, create bool) (bool,
	*AdjRIBPolicyParams) {
	if p.NeighborConf.Neighbor.Config.AdjRIBInFilter == "" {
		p.logger.Debugf("Peer %s - RIB In filter is not set", p.NeighborConf.Neighbor.NeighborAddress)
		return true, nil
	}

	return p.checkAdjRIBFilter(nlri, route, peEntity, p.server.ribInPE, policyCommonDefs.PolicyPath_Import, create)
}

func (p *Peer) checkRIBOutFilter(nlri packet.NLRI, route *bgprib.AdjRIBPathIdRoute,
	peEntity *utilspolicy.PolicyEngineFilterEntityParams, create bool) (bool,
	*AdjRIBPolicyParams) {
	if p.NeighborConf.Neighbor.Config.AdjRIBOutFilter == "" {
		p.logger.Debugf("Peer %s - RIB Out filter is not set", p.NeighborConf.Neighbor.NeighborAddress)
		return true, nil
	}

	return p.checkAdjRIBFilter(nlri, route, peEntity, p.server.ribOutPE, policyCommonDefs.PolicyPath_Export, create)
}

func (p *Peer) processUpdates(protoFamily uint32, nlris *[]packet.NLRI, path *bgprib.Path) {
	var ok bool
	var route *bgprib.AdjRIBRoute
	var pathIdRoute *bgprib.AdjRIBPathIdRoute
	total := len(*nlris)
	last := total - 1
	idx := 0

	peEntity := bgppolicy.GetPolicyEngineFilterEntity(path)

	for i := 0; i < total; i++ {
		nlri := (*nlris)[idx]
		if nlri == nil {
			if idx >= last {
				break
			}
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}

		ip := nlri.GetCIDR()
		if route, ok = p.ribIn[protoFamily][ip]; !ok {
			route = bgprib.NewAdjRIBRoute(p.NeighborConf.Neighbor.NeighborAddress, protoFamily, nlri)
			p.ribIn[protoFamily][ip] = route
			p.logger.Infof("Neighbor %s: add nlri %s protocol family %d",
				p.NeighborConf.RunningConf.NeighborAddress, ip, protoFamily)
		}
		pathIdRoute = route.AddPath(nlri.GetPathId(), path)
		p.logger.Infof("Neighbor %s: add path id %d for nlri %s protocol family %d to RIB-In %+v",
			p.NeighborConf.RunningConf.NeighborAddress, nlri.GetPathId(), ip, protoFamily, p.ribIn[protoFamily])

		accept, _ := p.checkRIBInFilter(nlri, pathIdRoute, peEntity, true)
		pathIdRoute.Accept = accept
		if !accept {
			p.logger.Infof("Neighbor %s: filter nlri %s", p.NeighborConf.RunningConf.NeighborAddress, ip)
			(*nlris)[idx] = (*nlris)[last]
			(*nlris)[last] = nil
			last--
			continue
		}
		idx++
	}
	(*nlris) = (*nlris)[:idx]
}

func (p *Peer) AddRouteNLRIs(route *bgprib.AdjRIBPathIdRoute, pathNLRIs map[*bgprib.Path]map[uint32]*bgprib.FilteredRoutes,
	add bool) map[*bgprib.Path]map[uint32]*bgprib.FilteredRoutes {
	var nlris *bgprib.FilteredRoutes
	var ok bool

	if route.Path != nil {
		routePath := route.Path
		if _, ok := pathNLRIs[routePath]; !ok {
			pathNLRIs[routePath] = make(map[uint32]*bgprib.FilteredRoutes)
		}

		if nlris, ok = pathNLRIs[routePath][route.ProtocolFamily]; !ok {
			nlris = bgprib.NewFilteredRoutes()
			pathNLRIs[routePath][route.ProtocolFamily] = nlris
		}

		if add {
			nlris.Add = append(nlris.Add, packet.ConstructNLRIFromPathIdAndNLRI(route.NLRI, route.PathId))
		} else {
			nlris.Remove = append(nlris.Remove, packet.ConstructNLRIFromPathIdAndNLRI(route.NLRI, route.PathId))
		}
	}

	return pathNLRIs
}

func (p *Peer) AdjRIBInPolicyUpdated(adjRibDir bgprib.AdjRIBDir, data interface{},
	updateFunc utilspolicy.PolicyApplyfunc) (map[uint32]map[*bgprib.Path][]*bgprib.Destination, []*bgprib.Destination,
	[]*bgprib.Destination) {
	filteredRoutes := make(map[*bgprib.Path]map[uint32]*bgprib.FilteredRoutes)
	adjRIB := p.GetAdjRIB(adjRibDir)
	for _, prefixRouteMap := range adjRIB {
		for _, adjRIBRoute := range prefixRouteMap {
			for _, adjRoute := range adjRIBRoute.PathIdRouteMap {
				if adjRoute == nil {
					continue
				}

				p.logger.Debugf("Peer %s - NLRI %s policylist %v hit %v before applying create policy",
					adjRoute.NLRI.GetPrefix(), adjRoute.PolicyList, adjRoute.PolicyHitCounter)
				peEntity := utilspolicy.PolicyEngineFilterEntityParams{
					DestNetIp:  adjRoute.NLRI.GetCIDR(),
					Neighbor:   p.NeighborConf.RunningConf.NeighborAddress.String(),
					PolicyList: adjRoute.PolicyList,
				}
				callbackInfo := &AdjRIBPolicyParams{
					CreateType: utilspolicy.Invalid,
					DeleteType: utilspolicy.Invalid,
					Peer:       p,
					Route:      adjRoute,
				}

				updateFunc(peEntity, data, callbackInfo)

				if !adjRoute.Accept && callbackInfo.Accept == Accept {
					filteredRoutes = p.AddRouteNLRIs(adjRoute, filteredRoutes, true)
				} else if adjRoute.Accept && callbackInfo.Accept == Reject {
					filteredRoutes = p.AddRouteNLRIs(adjRoute, filteredRoutes, false)
				}
			}
		}
	}

	updated, withdrawn, updatedAddPaths, addedAllPrefixes := p.locRib.ProcessFilteredRoutes(p.NeighborConf,
		filteredRoutes, p.server.AddPathCount)
	if !addedAllPrefixes {
		p.MaxPrefixesExceeded()
	}

	return updated, withdrawn, updatedAddPaths
}

func (p *Peer) ReceiveUpdate(pktInfo *packet.BGPPktSrc) (map[uint32]map[*bgprib.Path][]*bgprib.Destination,
	[]*bgprib.Destination, []*bgprib.Destination) {
	var mpReachProtoFamily, mpUnreachProtoFamily uint32 = 0, 0
	mpReachNLRI := make([]packet.NLRI, 0)
	updated := make(map[uint32]map[*bgprib.Path][]*bgprib.Destination)
	withdrawn := make([]*bgprib.Destination, 0)
	updatedAddPaths := make([]*bgprib.Destination, 0)
	addedAllPrefixes := true

	atomic.AddUint32(&p.NeighborConf.Neighbor.State.Queues.Input, ^uint32(0))
	p.NeighborConf.Neighbor.State.Messages.Received.Update++

	asLoop := false
	updateMsg := pktInfo.Msg.Body.(*packet.BGPUpdate)
	if packet.HasASLoop(updateMsg.PathAttributes, p.NeighborConf.RunningConf.LocalAS) {
		p.logger.Infof("Neighbor %s: Recived Update message has AS loop", p.NeighborConf.Neighbor.NeighborAddress)
		asLoop = true
	}

	protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
	//mpReach, mpUnreach := packet.GetMPAttrs(updateMsg.PathAttributes)
	//path := bgprib.NewPath(p.locRib, p.NeighborConf, updateMsg.PathAttributes, mpReach, bgprib.RouteTypeEGP)

	mpReach, mpUnreach := packet.RemoveMPAttrs(&updateMsg.PathAttributes)
	//remPath := bgprib.NewPath(p.locRib, p.neighborConf, updateMsg.PathAttributes, mpReach, RouteTypeEGP)
	path := bgprib.NewPath(p.locRib, p.NeighborConf, updateMsg.PathAttributes, mpReach, bgprib.RouteTypeEGP)

	p.processWithdraws(protoFamily, &updateMsg.WithdrawnRoutes)
	if asLoop {
		updateMsg.NLRI = make([]packet.NLRI, 0)
	} else {
		p.processUpdates(protoFamily, &updateMsg.NLRI, path)
	}

	if len(updateMsg.WithdrawnRoutes) > 0 || len(updateMsg.NLRI) > 0 {
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = p.locRib.ProcessUpdate(p.NeighborConf, path,
			updateMsg.NLRI, updateMsg.WithdrawnRoutes, protoFamily, p.server.AddPathCount, updated, withdrawn,
			updatedAddPaths)
		if !addedAllPrefixes {
			p.MaxPrefixesExceeded()
		}
	}

	if mpUnreach != nil {
		mpUnreachProtoFamily := packet.GetProtocolFamily(mpUnreach.AFI, mpUnreach.SAFI)
		p.processWithdraws(mpUnreachProtoFamily, &(mpUnreach.NLRI))
	}

	mpProtoFamilySame := false
	if mpReach != nil {
		if asLoop {
			mpReach.NLRI = make([]packet.NLRI, 0)
		} else {
			mpReachProtoFamily = packet.GetProtocolFamily(mpReach.AFI, mpReach.SAFI)
			p.processUpdates(mpReachProtoFamily, &(mpReach.NLRI), path)
			if mpReachProtoFamily == mpUnreachProtoFamily {
				mpProtoFamilySame = true
				mpReachNLRI = mpReach.NLRI
			}
		}
	}

	if mpUnreach != nil && len(mpUnreach.NLRI) > 0 {
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = p.locRib.ProcessUpdate(p.NeighborConf, path,
			mpReachNLRI, mpUnreach.NLRI, mpUnreachProtoFamily, p.server.AddPathCount, updated, withdrawn,
			updatedAddPaths)
		if !addedAllPrefixes {
			p.MaxPrefixesExceeded()
		}
	}

	if mpReach != nil && !mpProtoFamilySame && len(mpReach.NLRI) > 0 {
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = p.locRib.ProcessUpdate(p.NeighborConf, path,
			mpReach.NLRI, make([]packet.NLRI, 0), mpReachProtoFamily, p.server.AddPathCount, updated, withdrawn,
			updatedAddPaths)
		if !addedAllPrefixes {
			p.MaxPrefixesExceeded()
		}
	}

	return updated, withdrawn, updatedAddPaths
}

func (p *Peer) updatePathAttrs(bgpMsg *packet.BGPMessage, path *bgprib.Path, medUpdated bool) bool {
	if p.NeighborConf.Neighbor.Transport.Config.LocalAddress == nil {
		p.logger.Errf("Neighbor %s: Can't send Update message, FSM is not in Established state",
			p.NeighborConf.Neighbor.NeighborAddress)
		return false
	}

	if bgpMsg == nil || bgpMsg.Body.(*packet.BGPUpdate).PathAttributes == nil {
		p.logger.Errf("Neighbor %s: Path attrs not found in BGP Update message",
			p.NeighborConf.Neighbor.NeighborAddress)
		return false
	}

	updateMsg := bgpMsg.Body.(*packet.BGPUpdate)
	if len(updateMsg.NLRI) == 0 && !packet.HasMPReachNLRI(updateMsg.PathAttributes) {
		return true
	}

	if p.NeighborConf.ASSize == 2 {
		packet.Convert4ByteTo2ByteASPath(bgpMsg)
	}

	removeRRPathAttrs := true
	if p.NeighborConf.IsInternal() {
		if path.NeighborConf != nil && (path.NeighborConf.IsRouteReflectorClient() ||
			p.NeighborConf.IsRouteReflectorClient()) {
			removeRRPathAttrs = false
			packet.AddOriginatorId(bgpMsg, path.NeighborConf.BGPId)
			packet.AddClusterId(bgpMsg, path.NeighborConf.RunningConf.RouteReflectorClusterId)
		} else {
			packet.SetLocalPref(bgpMsg, path.GetPreference(), false)
		}
		if path.NeighborConf == nil || p.NeighborConf.RunningConf.NextHopSelf {
			packet.SetNextHop(bgpMsg, p.NeighborConf.Neighbor.Transport.Config.LocalAddress)
		}
	} else {
		// Do change these path attrs for local routes
		if path.NeighborConf != nil && !medUpdated {
			packet.RemoveMultiExitDisc(bgpMsg)
		}
		packet.PrependAS(bgpMsg, p.NeighborConf.RunningConf.LocalAS, p.NeighborConf.ASSize)
		if updateMsg.NLRI != nil && len(updateMsg.NLRI) > 0 {
			packet.SetNextHop(bgpMsg, p.NeighborConf.Neighbor.Transport.Config.LocalAddress)
		} else if len(updateMsg.PathAttributes) > 0 {
			packet.RemoveNextHop(&(updateMsg.PathAttributes))
		}
		packet.RemoveLocalPref(bgpMsg)
	}

	if removeRRPathAttrs {
		packet.RemoveOriginatorId(bgpMsg)
		packet.RemoveClusterList(bgpMsg)
	}

	return true
}

func (p *Peer) sendUpdateMsg(msg *packet.BGPMessage, path *bgprib.Path, medUpdated bool) {
	if p.fsmManager == nil {
		p.logger.Errf("Can't send update, FSM Manager is not instantiated for neighbor %s",
			p.NeighborConf.Neighbor.NeighborAddress)
		return
	}

	if path != nil && path.NeighborConf != nil {
		if path.NeighborConf.IsInternal() {

			if p.NeighborConf.IsInternal() && !path.NeighborConf.IsRouteReflectorClient() &&
				!p.NeighborConf.IsRouteReflectorClient() {
				return
			}
		}

		// Don't send the update to the peer that sent the update.
		if p.NeighborConf.RunningConf.NeighborAddress.String() ==
			path.NeighborConf.RunningConf.NeighborAddress.String() {
			return
		}
	}

	if p.updatePathAttrs(msg, path, medUpdated) {
		atomic.AddUint32(&p.NeighborConf.Neighbor.State.Queues.Output, 1)
		p.fsmManager.SendUpdateMsg(msg)
	}
}

func (p *Peer) isAdvertisable(path *bgprib.Path) bool {
	if path != nil && path.NeighborConf != nil {
		if path.NeighborConf.IsInternal() {
			if p.NeighborConf.IsInternal() && !path.NeighborConf.IsRouteReflectorClient() &&
				!p.NeighborConf.IsRouteReflectorClient() {
				return false
			}
		}

		// Don't send the update to the peer that sent the update.
		if p.NeighborConf.RunningConf.NeighborAddress.String() ==
			path.NeighborConf.RunningConf.NeighborAddress.String() {
			return false
		}

		if packet.HasASLoop(path.PathAttrs, p.NeighborConf.RunningConf.PeerAS) {
			return false
		}

	}

	return true
}

func (p *Peer) addPathFamilyToUpdated(pathAdded bool, path *bgprib.Path, stmt string, protoFamily uint32,
	updated map[*bgprib.Path]map[string]map[uint32][]packet.NLRI) (bool,
	map[*bgprib.Path]map[string]map[uint32][]packet.NLRI) {
	if !pathAdded {
		updated[path] = make(map[string]map[uint32][]packet.NLRI)
		pathAdded = true
	}

	if _, ok := updated[path][stmt]; !ok {
		updated[path][stmt] = make(map[uint32][]packet.NLRI)
	}

	if _, ok := updated[path][stmt][protoFamily]; !ok {
		updated[path][stmt][protoFamily] = make([]packet.NLRI, 0)
	}

	return pathAdded, updated
}

func (p *Peer) calculateAddPathsAdvertisements(dest *bgprib.Destination, path *bgprib.Path,
	newUpdated map[*bgprib.Path]map[string]map[uint32][]packet.NLRI, withdrawList map[uint32][]packet.NLRI,
	addPathsTx int, policyStmts map[string]utilspolicy.PolicyStmt) (
	map[*bgprib.Path]map[string]map[uint32][]packet.NLRI, map[uint32][]packet.NLRI,
	map[string]utilspolicy.PolicyStmt) {
	pathIdMap := make(map[uint32]*bgprib.Path)
	ip := dest.NLRI.GetCIDR()
	protoFamily := dest.GetProtocolFamily()

	if _, ok := p.ribOut[protoFamily][ip]; !ok {
		p.logger.Info("Neighbor", p.NeighborConf.Neighbor.NeighborAddress,
			"calculateAddPathsAdvertisements - processing updates, dest", ip, "not found in rib out")
		p.ribOut[protoFamily][ip] = bgprib.NewAdjRIBRoute(p.NeighborConf.Neighbor.NeighborAddress, protoFamily,
			dest.NLRI)
	}

	ribOutRoute := p.ribOut[protoFamily][ip]

	pathAdded := false
	if _, ok := newUpdated[path]; ok {
		pathAdded = true
	}

	if p.isAdvertisable(path) {
		route := dest.LocRibPathRoute
		if path != nil { // Loc-RIB path changed
			pathIdRoute := ribOutRoute.AddPath(route.OutPathId, path)
			peEntity := bgppolicy.GetPolicyEngineFilterEntity(path)
			if accept, policyParams := p.checkRIBOutFilter(dest.NLRI, pathIdRoute, peEntity, true); accept {
				stmt := ""
				if policyParams != nil {
					stmt = policyParams.PolicyStmt.Name
					if _, ok := policyStmts[stmt]; !ok {
						policyStmts[stmt] = policyParams.PolicyStmt
					}
				}

				pathAdded, newUpdated = p.addPathFamilyToUpdated(pathAdded, path, stmt, protoFamily, newUpdated)
				nlri := packet.NewExtNLRI(route.OutPathId, dest.NLRI.GetIPPrefix())
				newUpdated[path][stmt][protoFamily] = append(newUpdated[path][stmt][protoFamily], nlri)
			}
		} else {
			path = dest.LocRibPath
		}
		pathIdMap[route.OutPathId] = path
	}

	for i := 0; i < len(dest.AddPaths) && len(pathIdMap) < (addPathsTx-1); i++ {
		route := dest.GetPathRoute(dest.AddPaths[i])
		if route != nil && p.isAdvertisable(dest.AddPaths[i]) {
			pathIdMap[route.OutPathId] = dest.AddPaths[i]
		}
	}

	for ribOutPathId, ribOutPath := range ribOutRoute.GetPathMap() {
		if path, ok := pathIdMap[ribOutPathId]; !ok {
			if p.checkRIBOutWithdraw(ribOutRoute.GetPathIdRoute(ribOutPathId)) {
				nlri := packet.NewExtNLRI(ribOutPathId, dest.NLRI.GetIPPrefix())
				withdrawList[protoFamily] = append(withdrawList[protoFamily], nlri)
			}
			ribOutRoute.RemovePath(ribOutPathId)
		} else if ribOutPath == path {
			delete(pathIdMap, ribOutPathId)
		} else if ribOutPath != path {
			pathIdRoute := ribOutRoute.AddPath(ribOutPathId, path)
			peEntity := bgppolicy.GetPolicyEngineFilterEntity(path)
			if accept, policyParams := p.checkRIBOutFilter(dest.NLRI, pathIdRoute, peEntity, true); accept {
				stmt := ""
				if policyParams != nil {
					stmt = policyParams.PolicyStmt.Name
					if _, ok := policyStmts[stmt]; !ok {
						policyStmts[stmt] = policyParams.PolicyStmt
					}
				}

				pathAdded, newUpdated = p.addPathFamilyToUpdated(pathAdded, path, stmt, protoFamily, newUpdated)
				nlri := packet.NewExtNLRI(ribOutPathId, dest.NLRI.GetIPPrefix())
				newUpdated[path][stmt][protoFamily] = append(newUpdated[path][stmt][protoFamily], nlri)
			}
			delete(pathIdMap, ribOutPathId)
		}
	}

	for pathId, path := range pathIdMap {
		pathIdRoute := ribOutRoute.AddPath(pathId, path)
		peEntity := bgppolicy.GetPolicyEngineFilterEntity(path)
		if accept, policyParams := p.checkRIBOutFilter(dest.NLRI, pathIdRoute, peEntity, true); accept {
			stmt := ""
			if policyParams != nil {
				stmt = policyParams.PolicyStmt.Name
				if _, ok := policyStmts[stmt]; !ok {
					policyStmts[stmt] = policyParams.PolicyStmt
				}
			}

			pathAdded, newUpdated = p.addPathFamilyToUpdated(pathAdded, path, stmt, protoFamily, newUpdated)
			nlri := packet.NewExtNLRI(pathId, dest.NLRI.GetIPPrefix())
			newUpdated[path][stmt][protoFamily] = append(newUpdated[path][stmt][protoFamily], nlri)
		}
		delete(pathIdMap, pathId)
	}

	return newUpdated, withdrawList, policyStmts
}

func (p *Peer) checkRIBOutWithdraw(route *bgprib.AdjRIBPathIdRoute) bool {
	if p.NeighborConf.Neighbor.Config.AdjRIBOutFilter == "" {
		p.logger.Debugf("Peer %s - withdraw %s RIB Out filter is not set", p.NeighborConf.Neighbor.NeighborAddress,
			route.NLRI)
		return true
	}

	if route.Accept {
		p.logger.Debugf("Peer %s - withdraw %s, Network was advertised", p.NeighborConf.Neighbor.NeighborAddress,
			route.NLRI)
		return true
	}
	return false
}

func (p *Peer) SendUpdate(updated map[uint32]map[*bgprib.Path][]*bgprib.Destination, withdrawn,
	updatedAddPaths []*bgprib.Destination) {
	p.logger.Infof("Neighbor %s: Send update message valid routes:%v, withdraw routes:%v",
		p.NeighborConf.Neighbor.NeighborAddress, updated, withdrawn)
	if p.NeighborConf.Neighbor.Transport.Config.LocalAddress == nil {
		p.logger.Errf("Neighbor %s: Can't send Update message, FSM is not in Established state",
			p.NeighborConf.Neighbor.NeighborAddress)
		return
	}

	addPathsTx := p.getAddPathsMaxTx()
	withdrawList := make(map[uint32][]packet.NLRI)
	newUpdated := make(map[*bgprib.Path]map[string]map[uint32][]packet.NLRI)
	policyStmts := make(map[string]utilspolicy.PolicyStmt)

	if len(withdrawn) > 0 {
		for _, dest := range withdrawn {
			if dest != nil {
				protoFamily := dest.GetProtocolFamily()
				if _, ok := withdrawList[protoFamily]; !ok {
					withdrawList[protoFamily] = make([]packet.NLRI, 0)
				}
				ip := dest.NLRI.GetCIDR()
				if p.ribOut[protoFamily] != nil && p.NeighborConf.AfiSafiMap[protoFamily] {
					route, ok := p.ribOut[protoFamily][ip]
					if !ok {
						p.logger.Errf("Neighbor %s: processing withdraws, dest %s not found in rib out",
							p.NeighborConf.Neighbor.NeighborAddress, ip)
						continue
					}

					delete(p.ribOut[protoFamily], ip)

					for pathId, pathIdRoute := range route.PathIdRouteMap {
						if !p.checkRIBOutWithdraw(pathIdRoute) {
							p.logger.Errf("Neighbor %s: processing withdraws, dest %s path id %d not advertised",
								p.NeighborConf.Neighbor.NeighborAddress, ip, pathId)
							continue
						}
						if addPathsTx > 0 {
							nlri := packet.NewExtNLRI(pathId, dest.NLRI.GetIPPrefix())
							withdrawList[protoFamily] = append(withdrawList[protoFamily], nlri)
						} else {
							withdrawList[protoFamily] = append(withdrawList[protoFamily], dest.NLRI)
							break
						}
					}
					route.RemoveAllPaths()
				}
			}
		}
	}

	for protoFamily, pathDestMap := range updated {
		if !p.NeighborConf.AfiSafiMap[protoFamily] {
			continue
		}
		if _, ok := p.ribOut[protoFamily]; !ok {
			p.ribOut[protoFamily] = make(map[string]*bgprib.AdjRIBRoute)
		}
		if _, ok := withdrawList[protoFamily]; !ok {
			withdrawList[protoFamily] = make([]packet.NLRI, 0)
		}
		for path, destinations := range pathDestMap {
			peEntity := bgppolicy.GetPolicyEngineFilterEntity(path)
			for _, dest := range destinations {
				if dest == nil {
					continue
				}
				ip := dest.NLRI.GetCIDR()
				if addPathsTx > 0 {
					newUpdated, withdrawList, policyStmts = p.calculateAddPathsAdvertisements(dest, path, newUpdated,
						withdrawList, addPathsTx, policyStmts)
				} else {
					route := dest.LocRibPathRoute
					pathId := route.OutPathId
					if !p.isAdvertisable(path) {
						if ribOutRoute := p.ribOut[protoFamily][ip]; ribOutRoute != nil {
							for _, pathIdRoute := range ribOutRoute.PathIdRouteMap {
								if p.checkRIBOutWithdraw(pathIdRoute) {
									withdrawList[protoFamily] = append(withdrawList[protoFamily], dest.NLRI)
									p.ribOut[protoFamily][ip].RemoveAllPaths()
									delete(p.ribOut[protoFamily], ip)
									break
								}
							}
						}
					} else {
						if _, ok := p.ribOut[protoFamily][ip]; !ok {
							p.ribOut[protoFamily][ip] = bgprib.NewAdjRIBRoute(p.NeighborConf.Neighbor.NeighborAddress,
								protoFamily, dest.NLRI)
						}

						ribOutRoute := p.ribOut[protoFamily][ip]
						for ribPathId, _ := range ribOutRoute.GetPathMap() {
							if pathId != ribPathId {
								ribOutRoute.RemovePath(ribPathId)
							}
						}

						if ribOutPath := ribOutRoute.GetPath(pathId); ribOutPath == nil || ribOutPath != path {
							ribOutRoute.AddPath(pathId, path)
							if accept, policyParams := p.checkRIBOutFilter(dest.NLRI,
								ribOutRoute.GetPathIdRoute(pathId), peEntity, true); accept {
								stmt := ""
								if policyParams != nil {
									stmt = policyParams.PolicyStmt.Name
									if _, ok := policyStmts[stmt]; !ok {
										policyStmts[stmt] = policyParams.PolicyStmt
									}
								}

								p.logger.Infof("Neighbor %s: path %v stmt %v proto family %v",
									p.NeighborConf.Neighbor.NeighborAddress, path, stmt, protoFamily)
								if _, ok := newUpdated[path]; !ok {
									newUpdated[path] = make(map[string]map[uint32][]packet.NLRI)
								}
								if _, ok := newUpdated[path][stmt]; !ok {
									newUpdated[path][stmt] = make(map[uint32][]packet.NLRI)
								}
								if _, ok := newUpdated[path][stmt][protoFamily]; !ok {
									newUpdated[path][stmt][protoFamily] = make([]packet.NLRI, 0)
								}
								newUpdated[path][stmt][protoFamily] = append(newUpdated[path][stmt][protoFamily],
									dest.NLRI.GetIPPrefix())
							}
						}
					}
				}
			}
		}
	}

	if addPathsTx > 0 {
		for _, dest := range updatedAddPaths {
			newUpdated, withdrawList, policyStmts = p.calculateAddPathsAdvertisements(dest, nil, newUpdated,
				withdrawList, addPathsTx, policyStmts)
		}
	}

	if withdrawList != nil {
		p.logger.Infof("Neighbor %s: Send update message withdraw routes:%+v",
			p.NeighborConf.Neighbor.NeighborAddress, withdrawList)
		var updateMsg *packet.BGPMessage
		var ipv4List []packet.NLRI
		protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
		if nlriList, ok := withdrawList[protoFamily]; ok && len(nlriList) > 0 {
			ipv4List = nlriList
			delete(withdrawList, protoFamily)
		}
		for protoFamily, nlriList := range withdrawList {
			if len(nlriList) > 0 {
				mpUnreachNLRI := packet.ConstructMPUnreachNLRIFromProtoFamily(protoFamily, nlriList)
				pathAtts := make([]packet.BGPPathAttr, 0)
				pathAtts = append(pathAtts, mpUnreachNLRI)
				updateMsg = packet.NewBGPUpdateMessage(ipv4List, pathAtts, nil)
				p.sendUpdateMsg(updateMsg.Clone(), nil, false)
				ipv4List = nil
			}
		}
		if ipv4List != nil {
			updateMsg = packet.NewBGPUpdateMessage(ipv4List, nil, nil)
			p.sendUpdateMsg(updateMsg.Clone(), nil, false)
		}
	}

	localAddress := p.NeighborConf.Neighbor.Transport.Config.LocalAddress
	p.logger.Infof("Neighbor %s: new updated routes:%+v, policy stmts:%+v", p.NeighborConf.Neighbor.NeighborAddress,
		newUpdated, policyStmts)
	for path, stmtPFMap := range newUpdated {
		for stmt, pfNLRIMap := range stmtPFMap {
			var updateMsg *packet.BGPMessage
			var ipv4List []packet.NLRI
			protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
			if nlriList, ok := pfNLRIMap[protoFamily]; ok {
				if len(nlriList) > 0 {
					ipv4List = nlriList
					delete(pfNLRIMap, protoFamily)
				}
			}

			for protoFamily, nlriList := range pfNLRIMap {
				if len(nlriList) > 0 {
					mpReachNLRI := packet.ConstructIPv6MPReachNLRI(protoFamily, localAddress, nil, nlriList)
					pa := packet.ClonePathAttrs(path.PathAttrs)
					pa = packet.AddMPReachNLRIToPathAttrs(pa, mpReachNLRI)
					medUpdated := false
					if policyStmt, ok := policyStmts[stmt]; ok {
						p.logger.Infof("Neighbor %s: apply policy stmt %+v to path attrs",
							p.NeighborConf.Neighbor.NeighborAddress, policyStmt)
						pa, medUpdated = bgppolicy.ApplyActionsToPacket(pa, policyStmt)
					}
					updateMsg = packet.NewBGPUpdateMessage(nil, pa, ipv4List)
					p.logger.Infof("Neighbor %s: Send update message valid routes:%+v, path attrs:%+v",
						p.NeighborConf.Neighbor.NeighborAddress, nlriList, path.PathAttrs)
					p.sendUpdateMsg(updateMsg.Clone(), path, medUpdated)
					ipv4List = nil
				}
			}

			if ipv4List != nil {
				p.logger.Infof("Neighbor %s: Send update message valid routes:%+v, path attrs:%+v",
					p.NeighborConf.Neighbor.NeighborAddress, ipv4List, path.PathAttrs)
				pa := packet.ClonePathAttrs(path.PathAttrs)
				medUpdated := false
				if policyStmt, ok := policyStmts[stmt]; ok {
					p.logger.Infof("Neighbor %s: apply policy stmt %+v to path attrs",
						p.NeighborConf.Neighbor.NeighborAddress, policyStmt)
					pa, medUpdated = bgppolicy.ApplyActionsToPacket(pa, policyStmt)
				}
				updateMsg := packet.NewBGPUpdateMessage(make([]packet.NLRI, 0), pa, ipv4List)
				p.sendUpdateMsg(updateMsg.Clone(), path, medUpdated)
			}
		}
	}
}

func (p *Peer) AdjRIBOutPolicyUpdated(data interface{}, updateFunc utilspolicy.PolicyApplyfunc) {
	filteredRoutes := make(map[*bgprib.Path]map[uint32]*bgprib.FilteredRoutes)
	adjRIB := p.GetAdjRIB(bgprib.AdjRIBDirOut)
	for _, prefixRouteMap := range adjRIB {
		for _, adjRIBRoute := range prefixRouteMap {
			for _, adjRoute := range adjRIBRoute.PathIdRouteMap {
				if adjRIBRoute == nil {
					continue
				}

				p.logger.Debugf("Peer %s - NLRI %s policylist %v hit %v before applying create policy",
					adjRoute.NLRI.GetPrefix(), adjRoute.PolicyList, adjRoute.PolicyHitCounter)
				peEntity := utilspolicy.PolicyEngineFilterEntityParams{
					DestNetIp:  adjRoute.NLRI.GetCIDR(),
					Neighbor:   p.NeighborConf.RunningConf.NeighborAddress.String(),
					PolicyList: adjRoute.PolicyList,
				}
				callbackInfo := &AdjRIBPolicyParams{
					CreateType: utilspolicy.Invalid,
					DeleteType: utilspolicy.Invalid,
					Peer:       p,
					Route:      adjRoute,
				}

				updateFunc(peEntity, data, callbackInfo)

				if !adjRoute.Accept && callbackInfo.Accept == Accept {
					filteredRoutes = p.AddRouteNLRIs(adjRoute, filteredRoutes, true)
				} else if adjRoute.Accept && callbackInfo.Accept == Reject {
					filteredRoutes = p.AddRouteNLRIs(adjRoute, filteredRoutes, false)
				}
			}
		}
	}

	localAddress := p.NeighborConf.Neighbor.Transport.Config.LocalAddress
	p.logger.Infof("Neighbor %s: filtered routes:%+v", p.NeighborConf.Neighbor.NeighborAddress, filteredRoutes)
	for path, pfNLRIMap := range filteredRoutes {
		var updateMsg *packet.BGPMessage
		var updateList, withdrawList []packet.NLRI
		protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
		if routesMap, ok := pfNLRIMap[protoFamily]; ok {
			if len(routesMap.Add) > 0 {
				updateList = routesMap.Add
			}
			if len(routesMap.Remove) > 0 {
				withdrawList = routesMap.Remove
			}
			delete(pfNLRIMap, protoFamily)
		}

		for protoFamily, routesMap := range pfNLRIMap {
			if len(routesMap.Add) > 0 || len(routesMap.Remove) > 0 {
				var pa []packet.BGPPathAttr
				if len(routesMap.Add) > 0 {
					pa = packet.ClonePathAttrs(path.PathAttrs)
					mpReachNLRI := packet.ConstructIPv6MPReachNLRI(protoFamily, localAddress, nil, routesMap.Add)
					pa = packet.AddMPReachNLRIToPathAttrs(pa, mpReachNLRI)
				}

				if len(routesMap.Remove) > 0 {
					if pa == nil {
						pa = make([]packet.BGPPathAttr, 0)
					}

					mpUnreachNLRI := packet.ConstructMPUnreachNLRIFromProtoFamily(protoFamily, routesMap.Remove)
					pa = packet.AddMPUnreachNLRIToPathAttrs(pa, mpUnreachNLRI)
				}

				updateMsg = packet.NewBGPUpdateMessage(withdrawList, pa, updateList)
				p.logger.Infof("Neighbor %s: Send update message valid routes:%+v, path attrs:%+v",
					p.NeighborConf.Neighbor.NeighborAddress, updateList, path.PathAttrs)
				p.sendUpdateMsg(updateMsg.Clone(), path, false)
				updateList = nil
				withdrawList = nil
			}
		}

		if updateList != nil || withdrawList != nil {
			p.logger.Infof("Neighbor %s: Send update message valid routes:%+v, path attrs:%+v",
				p.NeighborConf.Neighbor.NeighborAddress, updateList, path.PathAttrs)
			updateMsg := packet.NewBGPUpdateMessage(withdrawList, path.PathAttrs, updateList)
			p.sendUpdateMsg(updateMsg.Clone(), path, false)
		}
	}

}
