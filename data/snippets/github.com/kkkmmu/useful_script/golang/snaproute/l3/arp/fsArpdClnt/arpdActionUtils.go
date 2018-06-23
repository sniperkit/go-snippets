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
	"models/actions"
)

func convertToThriftFromArpDeleteByIfName(obj *actions.ArpDeleteByIfName) *arpd.ArpDeleteByIfName {
	return &arpd.ArpDeleteByIfName{
		IfName: string(obj.IfName),
	}
}

func convertFromThriftToArpDeleteByIfName(obj *arpd.ArpDeleteByIfName) *actions.ArpDeleteByIfName {
	return &actions.ArpDeleteByIfName{
		IfName: string(obj.IfName),
	}
}

func convertToThriftFromArpRefreshByIfName(obj *actions.ArpRefreshByIfName) *arpd.ArpRefreshByIfName {
	return &arpd.ArpRefreshByIfName{
		IfName: string(obj.IfName),
	}
}

func convertFromThriftToArpRefreshByIfName(obj *arpd.ArpRefreshByIfName) *actions.ArpRefreshByIfName {
	return &actions.ArpRefreshByIfName{
		IfName: string(obj.IfName),
	}
}

func convertToThriftFromArpRefreshByIPv4Addr(obj *actions.ArpRefreshByIPv4Addr) *arpd.ArpRefreshByIPv4Addr {
	return &arpd.ArpRefreshByIPv4Addr{
		IpAddr: string(obj.IpAddr),
	}
}

func convertFromThriftToArpRefreshByIPv4Addr(obj *arpd.ArpRefreshByIPv4Addr) *actions.ArpRefreshByIPv4Addr {
	return &actions.ArpRefreshByIPv4Addr{
		IpAddr: string(obj.IpAddr),
	}
}

func convertToThriftFromArpDeleteByIPv4Addr(obj *actions.ArpDeleteByIPv4Addr) *arpd.ArpDeleteByIPv4Addr {
	return &arpd.ArpDeleteByIPv4Addr{
		IpAddr: string(obj.IpAddr),
	}
}

func convertFromThriftToArpDeleteByIPv4Addr(obj *arpd.ArpDeleteByIPv4Addr) *actions.ArpDeleteByIPv4Addr {
	return &actions.ArpDeleteByIPv4Addr{
		IpAddr: string(obj.IpAddr),
	}
}
