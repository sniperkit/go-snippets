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
	"l3/ndp/api"
	_ "l3/ndp/config"
	_ "l3/ndp/debug"
	"ndpd"
)

func (h *ConfigHandler) ExecuteActionNdpDeleteByIfName(config *ndpd.NdpDeleteByIfName) (bool, error) {
	api.SendDeleteByIfName(config.IfName)
	return true, nil
}

func (h *ConfigHandler) ExecuteActionNdpDeleteByIPv6Addr(config *ndpd.NdpDeleteByIPv6Addr) (bool, error) {
	api.SendDeleteByNeighborIp(config.IpAddr)
	return true, nil
}

func (h *ConfigHandler) ExecuteActionNdpRefreshByIfName(config *ndpd.NdpRefreshByIfName) (bool, error) {
	api.SendRefreshByIfName(config.IfName)
	return true, nil
}

func (h *ConfigHandler) ExecuteActionNdpRefreshByIPv6Addr(config *ndpd.NdpRefreshByIPv6Addr) (bool, error) {
	api.SendRefreshByNeighborIp(config.IpAddr)
	return true, nil
}
