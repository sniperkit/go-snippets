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

package fsArpdClnt

import (
	"arpd"
	"models/objects"
)

func convertToThriftFromArpEntryState(obj *objects.ArpEntryState) *arpd.ArpEntryState {
	return &arpd.ArpEntryState{
		Vlan:           string(obj.Vlan),
		MacAddr:        string(obj.MacAddr),
		Intf:           string(obj.Intf),
		IpAddr:         string(obj.IpAddr),
		ExpiryTimeLeft: string(obj.ExpiryTimeLeft),
	}
}

func convertFromThriftToArpEntryState(obj *arpd.ArpEntryState) *objects.ArpEntryState {
	return &objects.ArpEntryState{
		Vlan:           string(obj.Vlan),
		MacAddr:        string(obj.MacAddr),
		Intf:           string(obj.Intf),
		IpAddr:         string(obj.IpAddr),
		ExpiryTimeLeft: string(obj.ExpiryTimeLeft),
	}
}

func convertToThriftToArpLinuxEntryState(obj *objects.ArpLinuxEntryState) *arpd.ArpLinuxEntryState {
	return &arpd.ArpLinuxEntryState{
		HWType:  string(obj.HWType),
		IfName:  string(obj.IfName),
		MacAddr: string(obj.MacAddr),
		IpAddr:  string(obj.IpAddr),
	}
}

func convertFromThriftToArpLinuxEntryState(obj *arpd.ArpLinuxEntryState) *objects.ArpLinuxEntryState {
	return &objects.ArpLinuxEntryState{
		HWType:  string(obj.HWType),
		IfName:  string(obj.IfName),
		MacAddr: string(obj.MacAddr),
		IpAddr:  string(obj.IpAddr),
	}
}

func convertToThriftFromArpGlobal(obj *objects.ArpGlobal) *arpd.ArpGlobal {
	return &arpd.ArpGlobal{
		Vrf:     string(obj.Vrf),
		Timeout: int32(obj.Timeout),
	}
}

func convertFromThriftToArpGlobal(obj *arpd.ArpGlobal) *objects.ArpGlobal {
	return &objects.ArpGlobal{
		Vrf:     string(obj.Vrf),
		Timeout: int32(obj.Timeout),
	}
}
