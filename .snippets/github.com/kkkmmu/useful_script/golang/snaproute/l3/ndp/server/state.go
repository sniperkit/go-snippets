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
	"l3/ndp/config"
	"l3/ndp/debug"
)

func (svr *NDPServer) GetNeighborEntries(idx, cnt int) (int, int, []config.NeighborConfig) {
	var nextIdx int
	var count int
	var i, j int

	length := len(svr.neighborKey)
	if length == 0 {
		debug.Logger.Err("No Neighbors created by server")
		return 0, 0, nil
	}
	var result []config.NeighborConfig

	svr.NeigborEntryLock.RLock()
	for i, j = 0, idx; i < cnt && j < length; j++ {
		key := svr.neighborKey[j]
		result = append(result, svr.NeighborInfo[key])
		debug.Logger.Debug("Adding Neigbhor:", svr.NeighborInfo[key])
		i++
	}
	svr.NeigborEntryLock.RUnlock()
	if j == length {
		nextIdx = 0
	}
	count = i
	return nextIdx, count, result
}

func (svr *NDPServer) GetNeighborEntry(ipAddr string) *config.NeighborConfig {
	svr.NeigborEntryLock.RLock()
	defer svr.NeigborEntryLock.RUnlock()
	for _, nbrKey := range svr.neighborKey {
		splitString := splitNeighborKey(nbrKey)
		if splitString[1] == ipAddr {
			nbrEntry, exists := svr.NeighborInfo[nbrKey]
			if exists {
				return &nbrEntry
			}
		}
	}
	return nil
}

func (svr *NDPServer) GetGlobalState(vrf string) *config.GlobalState {
	gblState := config.GlobalState{}
	gblState.Vrf = vrf
	gblState.TotalRxPackets = svr.counter.Rcvd
	gblState.TotalTxPackets = svr.counter.Send
	gblState.Neighbors = int32(len(svr.neighborKey))
	gblState.RetransmitInterval = int32(svr.NdpConfig.RetransTime)
	gblState.ReachableTime = int32(svr.NdpConfig.ReachableTime)
	gblState.RouterAdvertisementInterval = int32(svr.NdpConfig.RaRestransmitTime)
	return &gblState
}

func (svr *NDPServer) PopulateInterfaceInfo(ifIndex int32, entry *config.InterfaceEntries) {
	intf, exists := svr.L3Port[ifIndex]
	if !exists {
		return
	}
	entry.IntfRef = intf.IntfRef
	entry.IfIndex = intf.IfIndex
	entry.LinkScopeIp = intf.LinkLocalIp
	entry.GlobalScopeIp = intf.IpAddr
	entry.SendPackets = intf.counter.Send
	entry.ReceivedPackets = intf.counter.Rcvd
	for _, nbrInfo := range intf.Neighbor {
		nbrEntry := config.NeighborEntry{}
		intf.PopulateNeighborInfo(nbrInfo, &nbrEntry)
		entry.Neighbor = append(entry.Neighbor, nbrEntry)
	}
}

func (svr *NDPServer) GetInterfaceNeighborEntries(idx, cnt int) (int, int, []config.InterfaceEntries) {
	var nextIdx int
	var count int
	var i, j int

	length := len(svr.ndpUpL3IntfStateSlice)
	if length == 0 {
		return 0, 0, nil
	}

	var result []config.InterfaceEntries
	svr.NeigborEntryLock.RLock()
	for i, j = 0, idx; i < cnt && j < length; j++ {
		ifIndex := svr.ndpUpL3IntfStateSlice[j]
		entry := config.InterfaceEntries{}
		svr.PopulateInterfaceInfo(ifIndex, &entry)
		result = append(result, entry)
		i++
	}
	svr.NeigborEntryLock.RUnlock()
	if j == length {
		nextIdx = 0
	}
	count = i
	return nextIdx, count, result
}

func (svr *NDPServer) GetInterfaceNeighborEntry(intfRef string) *config.InterfaceEntries {
	result := config.InterfaceEntries{}
	ifIndex, exists := svr.L3IfIntfRefToIfIndex[intfRef]
	if !exists {
		return &result
	}
	svr.PopulateInterfaceInfo(ifIndex, &result)
	return &result
}
