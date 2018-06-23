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
	"l3/bgp/config"
	"l3/bgp/server"
	"sync"
)

type ApiLayer struct {
	server *server.BGPServer
}

var bgpapi *ApiLayer = nil
var once sync.Once

/*  Singleton instance should be accesible only within api
 */
func getInstance() *ApiLayer {
	once.Do(func() {
		bgpapi = &ApiLayer{}
	})
	return bgpapi
}

/*  Initialize bgp api layer with the channels that will be used for communicating
 *  with the server
 */
func Init(server *server.BGPServer) {
	bgpapi = getInstance()
	bgpapi.server = server
}

/*  Send bfd state information from bfd manager to server
 */
func SendBfdNotification(DestIp string, State bool, Oper config.Operation) {
	bgpapi.server.BfdCh <- config.BfdInfo{
		DestIp: DestIp,
		State:  State,
		Oper:   Oper,
	}
}

/*  Send interface state notification to server
 */
func SendIntfNotification(ifIndex int32, ipAddr string, linklocalIp string, state config.Operation) {
	bgpapi.server.IntfCh <- config.IntfStateInfo{
		Idx:         ifIndex,
		IPAddr:      ipAddr,
		LinkLocalIP: linklocalIp,
		State:       state,
	}
}

/*  Send interface map notification to server
 */
func SendIntfMapNotification(ifIndex int32, ifName string) {
	bgpapi.server.IntfMapCh <- config.IntfMapInfo{
		Idx:    ifIndex,
		IfName: ifName,
	}
}

/*  Send Routes information to server
 */
func SendRouteNotification(add []*config.RouteInfo, remove []*config.RouteInfo) {
	bgpapi.server.RoutesCh <- &config.RouteCh{
		Add:    add,
		Remove: remove,
	}
}
