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
	"arpd"
	"errors"
	"fmt"
	"l3/arp/server"
	"strconv"
)

func (h *ARPHandler) convertArpEntryToThrift(arpState server.ArpState) *arpd.ArpEntryState {
	arpEnt := arpd.NewArpEntryState()
	arpEnt.IpAddr = arpState.IpAddr
	arpEnt.MacAddr = arpState.MacAddr
	if arpState.VlanId == h.server.SysRsvdVlan {
		arpEnt.Vlan = "Internal Vlan"
	} else {
		arpEnt.Vlan = strconv.Itoa((arpState.VlanId))
	}
	arpEnt.Intf = arpState.Intf
	arpEnt.ExpiryTimeLeft = arpState.ExpiryTimeLeft
	return arpEnt
}

func (h *ARPHandler) GetBulkArpEntryState(fromIdx arpd.Int, count arpd.Int) (*arpd.ArpEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("GetBulk call for ArpEntries..."))
	if h.server.IsLinuxOnly {
		return nil, errors.New("Not supported in Linux only plugin")
	}
	nextIdx, currCount, arpEntries := h.server.GetBulkArpEntry(int(fromIdx), int(count))
	if arpEntries == nil {
		err := errors.New("Arp is busy refreshing the Arp Cache")
		return nil, err
	}
	arpEntryResponse := make([]*arpd.ArpEntryState, len(arpEntries))
	for idx, item := range arpEntries {
		arpEntryResponse[idx] = h.convertArpEntryToThrift(item)
	}
	arpEntryBulk := arpd.NewArpEntryStateGetInfo()
	arpEntryBulk.Count = arpd.Int(currCount)
	arpEntryBulk.StartIdx = arpd.Int(fromIdx)
	arpEntryBulk.EndIdx = arpd.Int(nextIdx)
	arpEntryBulk.More = (nextIdx != 0)
	arpEntryBulk.ArpEntryStateList = arpEntryResponse
	return arpEntryBulk, nil
}

func (h *ARPHandler) convertArpLinuxEntryToThrift(arpLinuxState server.ArpLinuxState) *arpd.ArpLinuxEntryState {
	arpEnt := arpd.NewArpLinuxEntryState()
	arpEnt.IpAddr = arpLinuxState.IpAddr
	arpEnt.HWType = arpLinuxState.HWType
	arpEnt.MacAddr = arpLinuxState.MacAddr
	arpEnt.IfName = arpLinuxState.IfName
	return arpEnt
}

func (h *ARPHandler) GetBulkArpLinuxEntryState(fromIdx arpd.Int, count arpd.Int) (*arpd.ArpLinuxEntryStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("GetBulk call for Linux Arp Entry"))
	nextIdx, currCount, arpLinuxEntry := h.server.GetBulkLinuxArpEntry(int(fromIdx), int(count))
	if arpLinuxEntry == nil {
		err := errors.New("Arp server unable to fetch liunx Arp Entry")
		return nil, err
	}
	arpLinuxEntryResponse := make([]*arpd.ArpLinuxEntryState, len(arpLinuxEntry))
	for idx, item := range arpLinuxEntry {
		arpLinuxEntryResponse[idx] = h.convertArpLinuxEntryToThrift(item)
	}
	arpLinuxEntryBulk := arpd.NewArpLinuxEntryStateGetInfo()
	arpLinuxEntryBulk.Count = arpd.Int(currCount)
	arpLinuxEntryBulk.StartIdx = arpd.Int(fromIdx)
	arpLinuxEntryBulk.EndIdx = arpd.Int(nextIdx)
	arpLinuxEntryBulk.More = (nextIdx != 0)
	arpLinuxEntryBulk.ArpLinuxEntryStateList = arpLinuxEntryResponse
	return arpLinuxEntryBulk, nil
}
