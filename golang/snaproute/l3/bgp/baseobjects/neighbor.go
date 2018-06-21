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

// neighbor.go
package base

import (
	"fmt"
	"l3/bgp/config"
	"l3/bgp/packet"
	"models/events"
	"net"
	"time"
	"utils/eventUtils"
	"utils/logging"
)

const IgnoreBfdFaultsDefaultTime uint32 = 300 // seconds

type NeighborConf struct {
	logger               *logging.Writer
	Global               *config.GlobalConfig
	Group                *config.PeerGroupConfig
	Neighbor             *config.Neighbor
	RunningConf          config.NeighborConfig
	BGPId                net.IP
	ASSize               uint8
	AfiSafiMap           map[uint32]bool
	MaxPrefixesThreshold uint32
	ignoreBfdFaultsTimer *time.Timer
}

func NewNeighborConf(logger *logging.Writer, globalConf *config.GlobalConfig, peerGroup *config.PeerGroupConfig,
	peerConf config.NeighborConfig) *NeighborConf {
	conf := NeighborConf{
		logger:               logger,
		Global:               globalConf,
		Group:                peerGroup,
		AfiSafiMap:           make(map[uint32]bool),
		BGPId:                net.IP{},
		MaxPrefixesThreshold: 0,
		RunningConf:          config.NeighborConfig{},
		Neighbor: &config.Neighbor{
			NeighborAddress: peerConf.NeighborAddress,
			Config:          peerConf,
		},
	}

	conf.SetRunningConf(peerGroup, &conf.RunningConf)
	conf.SetNeighborState(&conf.RunningConf)
	conf.setOtherStates()
	conf.AfiSafiMap, _ = packet.GetProtocolFromConfig(&conf.Neighbor.AfiSafis, conf.Neighbor.NeighborAddress)
	return &conf
}

func (n *NeighborConf) SetNeighborAddress(ip net.IP) {
	n.Neighbor.NeighborAddress = ip
	n.Neighbor.Config.NeighborAddress = ip
	n.Neighbor.State.NeighborAddress = ip
	n.RunningConf.NeighborAddress = ip
}

func (n *NeighborConf) ResetNeighborAddress() {
	n.RunningConf.NeighborAddress = nil
	n.Neighbor.NeighborAddress = nil
	n.Neighbor.Config.NeighborAddress = nil
	n.Neighbor.State.NeighborAddress = nil
}

func (n *NeighborConf) SetNeighborState(peerConf *config.NeighborConfig) {
	n.Neighbor.State = config.NeighborState{
		Disabled:                peerConf.Disabled,
		PeerAS:                  peerConf.PeerAS,
		LocalAS:                 peerConf.LocalAS,
		UpdateSource:            peerConf.UpdateSource,
		NextHopSelf:             peerConf.NextHopSelf,
		AuthPassword:            peerConf.AuthPassword,
		Description:             peerConf.Description,
		NeighborAddress:         peerConf.NeighborAddress,
		IfIndex:                 peerConf.IfIndex,
		RouteReflectorClusterId: peerConf.RouteReflectorClusterId,
		RouteReflectorClient:    peerConf.RouteReflectorClient,
		MultiHopEnable:          peerConf.MultiHopEnable,
		MultiHopTTL:             peerConf.MultiHopTTL,
		ConnectRetryTime:        peerConf.ConnectRetryTime,
		HoldTime:                peerConf.HoldTime,
		KeepaliveTime:           peerConf.KeepaliveTime,
		PeerGroup:               peerConf.PeerGroup,
		AddPathsRx:              false,
		AddPathsMaxTx:           0,
		MaxPrefixes:             peerConf.MaxPrefixes,
		MaxPrefixesThresholdPct: peerConf.MaxPrefixesThresholdPct,
		MaxPrefixesDisconnect:   peerConf.MaxPrefixesDisconnect,
		MaxPrefixesRestartTimer: peerConf.MaxPrefixesRestartTimer,
		TotalPrefixes:           0,
		AdjRIBInFilter:          peerConf.AdjRIBInFilter,
		AdjRIBOutFilter:         peerConf.AdjRIBOutFilter,
	}
	n.MaxPrefixesThreshold = uint32(float64(peerConf.MaxPrefixes*uint32(peerConf.MaxPrefixesThresholdPct)) / 100)
}

func (n *NeighborConf) setOtherStates() {
	if n.RunningConf.LocalAS == n.RunningConf.PeerAS {
		n.Neighbor.State.PeerType = config.PeerTypeInternal
	} else {
		n.Neighbor.State.PeerType = config.PeerTypeExternal
	}
	if n.RunningConf.BfdEnable {
		n.Neighbor.State.BfdNeighborState = "up"
	} else {
		n.Neighbor.State.BfdNeighborState = "down"
	}
}

func (n *NeighborConf) copyNonKeyNeighConfAttrs(nConf config.NeighborConfig) {
	n.Neighbor.Config.BaseConfig = nConf.BaseConfig
	n.Neighbor.Config.PeerGroup = nConf.PeerGroup
	n.Neighbor.Config.Disabled = nConf.Disabled
}

func (n *NeighborConf) UpdateNeighborConf(nConf config.NeighborConfig, bgp *config.Bgp) {
	n.logger.Infof("UpdateNeighborConf - nConf=%+v", nConf)
	n.copyNonKeyNeighConfAttrs(nConf)
	n.logger.Infof("UpdateNeighborConf - updated conf=%+v", n.Neighbor.Config)
	n.RunningConf = config.NeighborConfig{}
	if (n.Group == nil && nConf.PeerGroup != "") || (n.Group != nil && nConf.PeerGroup != n.Group.Name) {
		protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(nConf.PeerAddressType)
		if _, ok := bgp.PeerGroups[protoFamily]; ok {
			if peerGroup, ok := bgp.PeerGroups[protoFamily][nConf.PeerGroup]; ok {
				n.GetNeighConfFromPeerGroup(&peerGroup.Config, &n.RunningConf)
			} else {
				n.logger.Err("Peer group", nConf.PeerGroup, "not found in BGP config")
			}
		}
	}
	n.GetConfFromNeighbor(&n.Neighbor.Config, &n.RunningConf)
	n.logger.Infof("UpdateNeighborConf - running conf=%+v", n.Neighbor.Config)
	n.SetNeighborState(&n.RunningConf)
	n.logger.Infof("UpdateNeighborConf - neigh state=%+v", n.Neighbor.State)
}

func (n *NeighborConf) UpdateGlobal(gConf *config.GlobalConfig) {
	n.logger.Infof("UpdateGlobal - gConf=%+v", gConf)
	n.Global = gConf
	n.RunningConf = config.NeighborConfig{}
	n.SetRunningConf(n.Group, &n.RunningConf)
	n.SetNeighborState(&n.RunningConf)
	n.setOtherStates()
}

func (n *NeighborConf) UpdatePeerGroup(peerGroup *config.PeerGroupConfig) {
	n.logger.Infof("UpdatePeerGroup - peergroup=%+v", peerGroup)
	n.Group = peerGroup
	n.RunningConf = config.NeighborConfig{}
	n.SetRunningConf(peerGroup, &n.RunningConf)
	n.SetNeighborState(&n.RunningConf)
}

func (n *NeighborConf) SetRunningConf(peerGroup *config.PeerGroupConfig, peerConf *config.NeighborConfig) {
	n.GetNeighConfFromGlobal(peerConf)
	n.GetNeighConfFromPeerGroup(peerGroup, peerConf)
	n.GetConfFromNeighbor(&n.Neighbor.Config, peerConf)
}

func (n *NeighborConf) GetNeighConfFromGlobal(peerConf *config.NeighborConfig) {
	peerConf.LocalAS = n.Global.AS
}

func (n *NeighborConf) GetNeighConfFromPeerGroup(groupConf *config.PeerGroupConfig, peerConf *config.NeighborConfig) {
	globalAS := peerConf.LocalAS
	if groupConf != nil {
		peerConf.BaseConfig = groupConf.BaseConfig
	}
	if peerConf.LocalAS == 0 {
		peerConf.LocalAS = globalAS
	}
}

func (n *NeighborConf) GetConfFromNeighbor(inConf *config.NeighborConfig, outConf *config.NeighborConfig) {
	if inConf.PeerAS != 0 {
		outConf.PeerAS = inConf.PeerAS
	}

	if inConf.LocalAS != 0 {
		outConf.LocalAS = inConf.LocalAS
	}

	if inConf.UpdateSource != "" {
		outConf.UpdateSource = inConf.UpdateSource
	}

	if inConf.NextHopSelf != false {
		outConf.NextHopSelf = inConf.NextHopSelf
	}

	if inConf.AuthPassword != "" {
		outConf.AuthPassword = inConf.AuthPassword
	}

	if inConf.Description != "" {
		outConf.Description = inConf.Description
	}

	if inConf.RouteReflectorClusterId != 0 {
		outConf.RouteReflectorClusterId = inConf.RouteReflectorClusterId
	}

	if inConf.RouteReflectorClient != false {
		outConf.RouteReflectorClient = inConf.RouteReflectorClient
	}

	if inConf.MultiHopEnable != false {
		outConf.MultiHopEnable = inConf.MultiHopEnable
	}

	if inConf.MultiHopTTL != 0 {
		outConf.MultiHopTTL = inConf.MultiHopTTL
	}

	if inConf.ConnectRetryTime != 0 {
		outConf.ConnectRetryTime = inConf.ConnectRetryTime
	}

	if inConf.HoldTime != 0 {
		outConf.HoldTime = inConf.HoldTime
	}

	if inConf.KeepaliveTime != 0 {
		outConf.KeepaliveTime = inConf.KeepaliveTime
	}

	if inConf.AddPathsRx != false {
		outConf.AddPathsRx = inConf.AddPathsRx
	}

	if inConf.AddPathsMaxTx != 0 {
		outConf.AddPathsMaxTx = inConf.AddPathsMaxTx
	}

	if inConf.BfdEnable != false {
		outConf.BfdEnable = inConf.BfdEnable
	}

	if inConf.BfdSessionParam != "" {
		outConf.BfdSessionParam = inConf.BfdSessionParam
	}

	if inConf.MaxPrefixes != 0 {
		outConf.MaxPrefixes = inConf.MaxPrefixes
	}

	if inConf.MaxPrefixesThresholdPct != 0 {
		outConf.MaxPrefixesThresholdPct = inConf.MaxPrefixesThresholdPct
	}

	if inConf.MaxPrefixesDisconnect != false {
		outConf.MaxPrefixesDisconnect = inConf.MaxPrefixesDisconnect
	}

	if inConf.MaxPrefixesRestartTimer != 0 {
		outConf.MaxPrefixesRestartTimer = inConf.MaxPrefixesRestartTimer
	}

	if inConf.AdjRIBInFilter != "" {
		outConf.AdjRIBInFilter = inConf.AdjRIBInFilter
	}

	if inConf.AdjRIBOutFilter != "" {
		outConf.AdjRIBOutFilter = inConf.AdjRIBOutFilter
	}

	n.setDefaults(outConf)
	outConf.PeerAddressType = inConf.PeerAddressType
	outConf.NeighborAddress = inConf.NeighborAddress
	outConf.IfIndex = inConf.IfIndex
	outConf.IfName = inConf.IfName
	outConf.PeerGroup = inConf.PeerGroup
	outConf.Disabled = inConf.Disabled
}

func (n *NeighborConf) setDefaults(nConf *config.NeighborConfig) {
	if nConf.ConnectRetryTime == 0 {
		nConf.ConnectRetryTime = config.BGPConnectRetryTime
	}

	if nConf.HoldTime == 0 { // default hold time is 180 seconds
		nConf.HoldTime = config.BGPHoldTimeDefault
	} else if nConf.HoldTime < 3 {
		nConf.HoldTime = 3
	}

	if nConf.KeepaliveTime == 0 { // default keep alive time is 60 seconds
		nConf.KeepaliveTime = nConf.HoldTime / 3
	}
}

func (n *NeighborConf) IsInternal() bool {
	return n.RunningConf.PeerAS == n.RunningConf.LocalAS
}

func (n *NeighborConf) IsExternal() bool {
	return n.RunningConf.LocalAS != n.RunningConf.PeerAS
}

func (n *NeighborConf) IsRouteReflectorClient() bool {
	return n.RunningConf.RouteReflectorClient
}

func (n *NeighborConf) IncrPrefixCount() {
	n.Neighbor.State.TotalPrefixes++
}

func (n *NeighborConf) DecrPrefixCount() {
	n.Neighbor.State.TotalPrefixes--
}

func (n *NeighborConf) SetPrefixCount(count uint32) {
	n.Neighbor.State.TotalPrefixes = 0
}

func (n *NeighborConf) CanAcceptNewPrefix() bool {
	if n.RunningConf.MaxPrefixes > 0 {
		if n.Neighbor.State.TotalPrefixes >= n.RunningConf.MaxPrefixes {
			n.logger.Warningf("Neighbor %s Number of prefixes received %d exceeds the max prefix limit %d",
				n.RunningConf.NeighborAddress, n.Neighbor.State.TotalPrefixes, n.MaxPrefixesThreshold)
			return false
		}

		if n.Neighbor.State.TotalPrefixes >= n.MaxPrefixesThreshold {
			n.logger.Warningf("Neighbor %s Number of prefixes received %d reached the threshold limit %d",
				n.RunningConf.NeighborAddress, n.Neighbor.State.TotalPrefixes, n.MaxPrefixesThreshold)
		}
	}

	return true
}

func (n *NeighborConf) PublishEvents(stateId uint32) {
	oldState := config.GetBGPStateToStr(config.BGPFSMState(n.Neighbor.State.SessionState))
	newState := config.GetBGPStateToStr(config.BGPFSMState(stateId))
	evtKey := events.BGPNeighborKey{
		NeighborAddress: n.Neighbor.NeighborAddress.String(),
		IfIndex:         n.Neighbor.Config.IfIndex,
	}
	additionalInfo := fmt.Sprintf("State change from %s to %s", oldState, newState)
	txEvt := eventUtils.TxEvent{
		EventId:        events.BGPNeighborStateChange,
		Key:            evtKey,
		AdditionalInfo: additionalInfo,
		AdditionalData: nil,
	}
	err := eventUtils.PublishEvents(&txEvt)
	if err != nil {
		n.logger.Err("Error publish new events for BGPNeighborStateChange")
	}
}

func (n *NeighborConf) FSMStateChange(state uint32) {
	n.logger.Infof("Neighbor %s: FSMStateChange %d", n.Neighbor.NeighborAddress, state)
	n.PublishEvents(state)
	n.Neighbor.State.SessionState = uint32(state)
	n.Neighbor.State.SessionStateUpdatedTime = time.Now()
}

func (n *NeighborConf) SetPeerAttrs(bgpId net.IP, asSize uint8, holdTime uint32, keepaliveTime uint32,
	addPathFamily map[packet.AFI]map[packet.SAFI]uint8) {
	n.BGPId = bgpId
	n.ASSize = asSize
	n.Neighbor.State.HoldTime = holdTime
	n.Neighbor.State.KeepaliveTime = keepaliveTime
	for afi, safiMap := range addPathFamily {
		if afi == packet.AfiIP || afi == packet.AfiIP6 {
			for _, val := range safiMap {
				if (val & packet.BGPCapAddPathRx) != 0 {
					n.logger.Infof("SetPeerAttrs - Neighbor %s set add paths maxtx to %d",
						n.Neighbor.NeighborAddress, n.RunningConf.AddPathsMaxTx)
					n.Neighbor.State.AddPathsMaxTx = n.RunningConf.AddPathsMaxTx
				}
				if (val & packet.BGPCapAddPathTx) != 0 {
					n.logger.Infof("SetPeerAttrs - Neighbor %s set add paths rx to %s",
						n.Neighbor.NeighborAddress, n.RunningConf.AddPathsRx)
					n.Neighbor.State.AddPathsRx = true
				}
			}
		}
	}
}

func (n *NeighborConf) BfdFaultSet() {
	n.Neighbor.State.BfdNeighborState = "down"
	if n.ignoreBfdFaultsTimer != nil {
		n.ignoreBfdFaultsTimer.Stop()
	}
	n.ignoreBfdFaultsTimer = time.AfterFunc(time.Duration(IgnoreBfdFaultsDefaultTime)*time.Second,
		n.IgnoreBfdFaultsTimerExpired)
}

func (n *NeighborConf) BfdFaultCleared() {
	if n.IgnoreBfdFaultsTimerExpired != nil {
		n.ignoreBfdFaultsTimer.Stop()
	}
	n.Neighbor.State.BfdNeighborState = "up"
}

func (n *NeighborConf) IgnoreBfdFaultsTimerExpired() {
	n.Neighbor.State.UseBfdState = false
}

func (n *NeighborConf) PeerConnEstablished() {
	n.Neighbor.State.UseBfdState = true
}

func (n *NeighborConf) PeerConnBroken() {
	n.Neighbor.State.ConnectRetryTime = n.RunningConf.ConnectRetryTime
	n.Neighbor.State.HoldTime = n.RunningConf.HoldTime
	n.Neighbor.State.KeepaliveTime = n.RunningConf.KeepaliveTime
	n.Neighbor.State.AddPathsRx = false
	n.Neighbor.State.AddPathsMaxTx = 0
	n.Neighbor.State.TotalPrefixes = 0
}
