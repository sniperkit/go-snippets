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
	//"fmt"
	"l3/bgp/config"
)

/*  Constructor for route manager
 */
func NewOvsRouteMgr() *OvsRouteMgr {
	mgr := &OvsRouteMgr{
		plugin: "ovsdb",
	}

	return mgr
}

func (mgr *OvsRouteMgr) Start() {

}

func (mgr *OvsRouteMgr) CreateRoute(cfg *config.RouteConfig) {
	//fmt.Println("Create Route called in", mgr.plugin, "with configs", cfg)
}

func (mgr *OvsRouteMgr) DeleteRoute(cfg *config.RouteConfig) {

}

func (mgr *OvsRouteMgr) UpdateRoute(cfg *config.RouteConfig, op string) {

}

func (mgr *OvsRouteMgr) GetNextHopInfo(ipAddr string, ifIndex int32) (*config.NextHopInfo, error) {
	return nil, nil
}
func (mgr *OvsRouteMgr) ApplyPolicy(applyList []*config.ApplyPolicyInfo, undoList []*config.ApplyPolicyInfo) {

	return
}
func (mgr *OvsRouteMgr) GetRoutes() ([]*config.RouteInfo, []*config.RouteInfo) {
	return nil, nil
}
