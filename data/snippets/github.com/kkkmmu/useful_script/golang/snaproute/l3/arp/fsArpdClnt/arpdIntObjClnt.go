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
	"arpdInt"
	"errors"
)

func (arpdClientMgr *FSArpdClntMgr) ResolveArpIPv4(destNetIp string, ifIdx int32) error {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		err := arpdClientMgr.ClientHdl.ResolveArpIPv4(destNetIp, arpdInt.Int(ifIdx))
		arpdMutex.Unlock()
		return err
	}
	return errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) DeleteResolveArpIPv4(NbrIP string) error {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		err := arpdClientMgr.ClientHdl.DeleteResolveArpIPv4(NbrIP)
		arpdMutex.Unlock()
		return err
	}
	return errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) DeleteArpEntry(ipAddr string) error {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		err := arpdClientMgr.ClientHdl.DeleteArpEntry(ipAddr)
		arpdMutex.Unlock()
		return err
	}
	return errors.New("Arpd Client Handle is nil")
}

func (arpdClientMgr *FSArpdClntMgr) SendGarp(ifName string, macAddr string, ipAddr string) error {
	if arpdClientMgr.ClientHdl != nil {
		arpdMutex.Lock()
		err := arpdClientMgr.ClientHdl.SendGarp(ifName, macAddr, ipAddr)
		arpdMutex.Unlock()
		return err
	}
	return errors.New("Arpd Client Handle is nil")
}
