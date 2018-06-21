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

func (h *ARPHandler) SendDeleteResolveArpIPv4(NextHopIp string) {
	rConf := server.DeleteResolvedIPv4{
		IpAddr: NextHopIp,
	}

	h.server.DeleteResolvedIPv4Ch <- rConf
	return
}

func (h *ARPHandler) DeleteArpGlobal(conf *arpd.ArpGlobal) (bool, error) {
	h.logger.Info("Delete Arp config attrs:", conf)
	return true, nil
}

func (h *ARPHandler) DeleteResolveArpIPv4(NextHopIp string) error {
	h.logger.Info("Received DeleteResolveArpIPv4 call with NextHopIp:", NextHopIp)
	if h.server.IsLinuxOnly {
		return nil
	}
	h.SendDeleteResolveArpIPv4(NextHopIp)
	return nil
}

func (h *ARPHandler) DeleteArpEntry(ipAddr string) error {
	h.logger.Info("Received Delete Arp Entry for:", ipAddr)
	if h.server.IsLinuxOnly {
		return errors.New("Not supported for linux only plugin")
	}
	h.server.DeleteArpEntryCh <- &server.DeleteArpEntry{ipAddr}
	return nil
}
