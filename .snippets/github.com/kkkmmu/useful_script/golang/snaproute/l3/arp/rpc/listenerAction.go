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
	"l3/arp/server"
)

func (h *ARPHandler) ExecuteActionArpDeleteByIfName(config *arpd.ArpDeleteByIfName) (bool, error) {
	h.logger.Info("Received ArpDeleteByIfName for", config)
	if h.server.IsLinuxOnly {
		return false, errors.New("Not supported for linux only plugin")
	}
	msg := server.ArpActionMsg{
		Type: server.DeleteByIfName,
		Obj:  config.IfName,
	}
	h.server.ArpActionCh <- msg
	return true, nil
}

func (h *ARPHandler) ExecuteActionArpDeleteByIPv4Addr(config *arpd.ArpDeleteByIPv4Addr) (bool, error) {
	h.logger.Info("Received ArpDeleteByIPv4Addr for", config)
	if h.server.IsLinuxOnly {
		return false, errors.New("Not supported for linux only plugin")
	}
	msg := server.ArpActionMsg{
		Type: server.DeleteByIPAddr,
		Obj:  config.IpAddr,
	}
	h.server.ArpActionCh <- msg
	return true, nil
}

func (h *ARPHandler) ExecuteActionArpRefreshByIfName(config *arpd.ArpRefreshByIfName) (bool, error) {
	h.logger.Info("Received ArpRefreshByIfName for", config)
	if h.server.IsLinuxOnly {
		return false, errors.New("Not supported for linux only plugin")
	}
	msg := server.ArpActionMsg{
		Type: server.RefreshByIfName,
		Obj:  config.IfName,
	}
	h.server.ArpActionCh <- msg
	return true, nil
}

func (h *ARPHandler) ExecuteActionArpRefreshByIPv4Addr(config *arpd.ArpRefreshByIPv4Addr) (bool, error) {
	h.logger.Info("Received ArpRefreshByIPv4Addr for", config)
	if h.server.IsLinuxOnly {
		return false, errors.New("Not supported for linux only plugin")
	}
	msg := server.ArpActionMsg{
		Type: server.RefreshByIPAddr,
		Obj:  config.IpAddr,
	}
	h.server.ArpActionCh <- msg
	return true, nil
}
