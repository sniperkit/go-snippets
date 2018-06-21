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
package api

import (
	"l3/vrrp/common"
	"l3/vrrp/server"
	"sync"
)

var vrrpApi *VRRPApiLayer = nil
var once sync.Once

type VRRPApiLayer struct {
	server *server.VrrpServer
}

func InitComplete() bool {
	if vrrpApi == nil {
		return false
	}
	if vrrpApi.server == nil {
		return false
	}
	return true
}

/*  Singleton instance should be accessible only within api
 */
func getApiInstance() *VRRPApiLayer {
	once.Do(func() {
		vrrpApi = &VRRPApiLayer{}
	})
	return vrrpApi
}

func Init(svr *server.VrrpServer) {
	vrrpApi = getApiInstance()
	vrrpApi.server = svr
}

func VrrpIntfConfig(cfg *common.IntfCfg) (bool, error) {
	rv, err := vrrpApi.server.ValidConfiguration(cfg)
	if rv == false {
		return rv, err
	}
	vrrpApi.server.CfgCh <- cfg
	return true, nil
}

func CreateVrrpGbl(cfg *common.GlobalConfig) {
	vrrpApi.server.GblCfgCh <- cfg
}

func UpdateVrrpGbl(cfg *common.GlobalConfig) {
	vrrpApi.server.GblCfgCh <- cfg
}

// this includes both v4 & v6 and create & delete
func SendIpIntfNotification(ipNotify *common.BaseIpInfo) {
	vrrpApi.server.L3IntfNotifyCh <- ipNotify
}

func GetAllV4IntfStates(from, count int) (n int, c int, result []common.State) {
	n, c, result = vrrpApi.server.GetV4Intfs(from, count)
	return n, c, result
}

func GetAllV6IntfStates(from, count int) (n int, c int, result []common.State) {
	n, c, result = vrrpApi.server.GetV6Intfs(from, count)
	return n, c, result
}

func GetVrrpIntfEntry(intfRef string, vrid int32, ipType int) *common.State {
	return vrrpApi.server.GetEntry(server.KeyInfo{intfRef, vrid, ipType})
}

func GetVrrpGlobalStateInfo(vrf string) (*common.GlobalState, error) {
	return vrrpApi.server.GetGlobalState(vrf), nil
}
