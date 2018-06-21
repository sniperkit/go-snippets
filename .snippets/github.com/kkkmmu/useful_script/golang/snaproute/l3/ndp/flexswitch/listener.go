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
	"l3/ndp/api"
	"l3/ndp/config"
	"l3/ndp/debug"
	"ndpd"
	"strconv"
	_ "utils/netUtils"
)

func (h *ConfigHandler) CreateNDPGlobal(config *ndpd.NDPGlobal) (bool, error) {
	return api.CreateGlobalConfig(config.Vrf, uint32(config.RetransmitInterval), uint32(config.ReachableTime),
		uint8(config.RouterAdvertisementInterval))
}

func (h *ConfigHandler) UpdateNDPGlobal(orgCfg *ndpd.NDPGlobal, newCfg *ndpd.NDPGlobal, attrset []bool, op []*ndpd.PatchOpInfo) (bool, error) {
	return true, nil
}

func (h *ConfigHandler) DeleteNDPGlobal(config *ndpd.NDPGlobal) (bool, error) {
	return false, errors.New("Delete of Global Object is not supported")
}

func convertNDPEntryStateToThriftEntry(state config.NeighborConfig) *ndpd.NDPEntryState {
	entry := ndpd.NewNDPEntryState()
	entry.IpAddr = state.IpAddr
	entry.MacAddr = state.MacAddr
	if state.VlanId == -1 {
		entry.Vlan = "Internal Vlan"
	} else {
		entry.Vlan = strconv.Itoa(int(state.VlanId))
	}
	entry.Intf = state.Intf
	entry.IfIndex = state.IfIndex
	return entry
}

func (h *ConfigHandler) GetBulkNDPEntryState(fromIdx ndpd.Int, count ndpd.Int) (*ndpd.NDPEntryStateGetInfo, error) {
	debug.Logger.Debug("Get Bulk NDP Entry state called")
	nextIdx, currCount, ndpEntries := api.GetAllNeigborEntries(int(fromIdx), int(count))
	if len(ndpEntries) == 0 || ndpEntries == nil {
		return nil, errors.New("No Neighbor Found")
	}
	ndpResp := make([]*ndpd.NDPEntryState, len(ndpEntries))
	for idx, ndpEntry := range ndpEntries {
		ndpResp[idx] = convertNDPEntryStateToThriftEntry(ndpEntry)
	}
	ndpEntryBulk := ndpd.NewNDPEntryStateGetInfo()
	ndpEntryBulk.StartIdx = fromIdx
	ndpEntryBulk.EndIdx = ndpd.Int(nextIdx)
	ndpEntryBulk.Count = ndpd.Int(currCount)
	ndpEntryBulk.More = (nextIdx != 0)
	ndpEntryBulk.NDPEntryStateList = ndpResp
	return ndpEntryBulk, nil
}

func (h *ConfigHandler) GetNDPEntryState(ipAddr string) (*ndpd.NDPEntryState, error) {
	ndpEntry := api.GetNeighborEntry(ipAddr)
	if ndpEntry == nil {
		return nil, errors.New("No Neighbor Found for Ip Address:" + ipAddr)
	}
	return convertNDPEntryStateToThriftEntry(*ndpEntry), nil
}

func (h *ConfigHandler) convertNDPGlobalStateEntryToThriftEntry(gblState *config.GlobalState) *ndpd.NDPGlobalState {
	entry := ndpd.NewNDPGlobalState()
	entry.Vrf = gblState.Vrf
	entry.TotalTxPackets = gblState.TotalTxPackets
	entry.TotalRxPackets = gblState.TotalRxPackets
	entry.Neighbors = gblState.Neighbors
	entry.RetransmitInterval = gblState.RetransmitInterval
	entry.ReachableTime = gblState.ReachableTime
	entry.RouterAdvertisementInterval = gblState.RouterAdvertisementInterval
	return entry
}

func (h *ConfigHandler) GetBulkNDPGlobalState(fromIndex ndpd.Int, count ndpd.Int) (*ndpd.NDPGlobalStateGetInfo, error) {
	ndpGlobalStateBulk := ndpd.NewNDPGlobalStateGetInfo()
	ndpGlobalStateBulk.EndIdx = ndpd.Int(0)
	ndpGlobalStateBulk.Count = ndpd.Int(1)
	ndpGlobalStateBulk.More = false
	ndpGlobalStateBulk.NDPGlobalStateList = make([]*ndpd.NDPGlobalState, 1)
	gblEntry, _ := api.GetNDPGlobalState("default")
	ndpGlobalStateBulk.NDPGlobalStateList[0] = h.convertNDPGlobalStateEntryToThriftEntry(gblEntry)

	return ndpGlobalStateBulk, nil
}

func (h *ConfigHandler) GetNDPGlobalState(vrf string) (*ndpd.NDPGlobalState, error) {
	gblEntry, _ := api.GetNDPGlobalState(vrf)
	return h.convertNDPGlobalStateEntryToThriftEntry(gblEntry), nil
}

func (h *ConfigHandler) convertIPV6AdjStateEntryToThriftEntry(state config.InterfaceEntries) *ndpd.IPV6AdjState {
	entry := ndpd.NewIPV6AdjState()
	entry.IntfRef = state.IntfRef
	entry.IfIndex = state.IfIndex
	entry.LinkScopeIp = state.LinkScopeIp
	entry.GlobalScopeIp = state.GlobalScopeIp
	entry.ReceivedPackets = state.ReceivedPackets
	entry.SendPackets = state.SendPackets
	for _, nbrEntry := range state.Neighbor {
		nbr := ndpd.NewNeighborEntry()
		nbr.IpAddr = nbrEntry.IpAddr
		nbr.MacAddr = nbrEntry.MacAddr
		nbr.ExpiryTimeLeft = nbrEntry.ExpiryTimeLeft
		nbr.SendPackets = nbrEntry.SendPackets
		nbr.ReceivedPackets = nbrEntry.ReceivedPackets
		nbr.State = nbrEntry.State
		entry.Neighbors = append(entry.Neighbors, nbr)
	}
	return entry
}

func (h *ConfigHandler) GetBulkIPV6AdjState(fromIdx ndpd.Int, count ndpd.Int) (*ndpd.IPV6AdjStateGetInfo, error) {
	nextIdx, currCount, ndpIntfEntries := api.GetAllNdpIntfState(int(fromIdx), int(count))
	if len(ndpIntfEntries) == 0 || ndpIntfEntries == nil {
		return nil, errors.New("No Neighbor Entries Learned")
	}
	ndpResp := make([]*ndpd.IPV6AdjState, len(ndpIntfEntries))
	for idx, ndpEntry := range ndpIntfEntries {
		ndpResp[idx] = h.convertIPV6AdjStateEntryToThriftEntry(ndpEntry)
	}
	ndpIntfStateBulk := ndpd.NewIPV6AdjStateGetInfo()
	ndpIntfStateBulk.StartIdx = fromIdx
	ndpIntfStateBulk.EndIdx = ndpd.Int(nextIdx)
	ndpIntfStateBulk.Count = ndpd.Int(currCount)
	ndpIntfStateBulk.More = (nextIdx != 0)
	ndpIntfStateBulk.IPV6AdjStateList = ndpResp
	return ndpIntfStateBulk, nil
}

func (h *ConfigHandler) GetIPV6AdjState(intfRef string) (*ndpd.IPV6AdjState, error) {
	ndpEntry := api.GetNdpIntfState(intfRef)
	if ndpEntry == nil {
		return nil, errors.New("No Neighbor Found for Port " + intfRef)
	}
	return h.convertIPV6AdjStateEntryToThriftEntry(*ndpEntry), nil
}

func (h *ConfigHandler) DeleteNdpEntry(ipAddr string) error {
	api.SendDeleteByNeighborIp(ipAddr)
	return nil
}
