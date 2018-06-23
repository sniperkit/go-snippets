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

package ovsMgr

import (
	"l3/bgp/config"
)

/*  Constructor for interface manager
 */
func NewOvsIntfMgr() *OvsIntfMgr {
	mgr := &OvsIntfMgr{
		plugin: "ovsdb",
	}

	return mgr
}

func (mgr *OvsIntfMgr) Start() {

}

func (mgr *OvsIntfMgr) GetIPv4Intfs() []*config.IntfStateInfo {
	return make([]*config.IntfStateInfo, 0)
}

func (mgr *OvsIntfMgr) GetIPv6Intfs() []*config.IntfStateInfo {
	return make([]*config.IntfStateInfo, 0)
}
func (mgr *OvsIntfMgr) GetIPv6Neighbors() []*config.IntfStateInfo {
	return make([]*config.IntfStateInfo, 0)
}
func (mgr *OvsIntfMgr) GetIPv4Information(ifIndex int32) (string, error) {
	return "", nil
}
func (mgr *OvsIntfMgr) GetIPv6Information(ifIndex int32) (string, error) {
	return "", nil
}

func (mgr *OvsIntfMgr) GetIfIndex(ifIndex, ifType int) int32 {
	return 1
}

func (mgr *OvsIntfMgr) PortStateChange() {

}
func (m *OvsIntfMgr) GetLogicalIntfInfo() []config.IntfMapInfo {
	return nil
}
func (m *OvsIntfMgr) GetVlanInfo() []config.IntfMapInfo {
	return nil
}
func (m *OvsIntfMgr) GetPortInfo() []config.IntfMapInfo {
	return nil
}
