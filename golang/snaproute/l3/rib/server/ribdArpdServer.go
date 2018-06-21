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

// ribdArpdServer.go
package server

import (
	//"arpdInt"
	//	"fmt"
	defs "l3/rib/ribdCommonDefs"
)

func (ribdServiceHandler *RIBDServer) arpdResolveRoute(routeInfoRecord RouteInfoRecord) {
	/*
		if arpdclnt.IsConnected == false {
			return
		}
	*/
	if routeInfoRecord.ipType == defs.IPv6 {
		return
	}
	ribdServiceHandler.ArpdClntPlugin.ResolveArpIPv4(routeInfoRecord.resolvedNextHopIpIntf.NextHopIp, int32(routeInfoRecord.nextHopIfIndex))
}
func (ribdServiceHandler *RIBDServer) arpdRemoveRoute(routeInfoRecord RouteInfoRecord) {
	/*
		if arpdclnt.IsConnected == false {
			return
		}
	*/
	if routeInfoRecord.ipType == defs.IPv6 {
		return
	}
	ribdServiceHandler.ArpdClntPlugin.DeleteResolveArpIPv4(routeInfoRecord.resolvedNextHopIpIntf.NextHopIp)
}
func (ribdServiceHandler *RIBDServer) StartArpdServer() {
	logger.Info("Starting the arpdserver loop")
	for {
		select {
		case route := <-ribdServiceHandler.ArpdRouteCh:
			logger.Debug(" received message on ArpdRouteCh, op:", route.Op)
			if route.Op == defs.Add {
				ribdServiceHandler.arpdResolveRoute(route.OrigConfigObject.(RouteInfoRecord))
			} else if route.Op == defs.Del {
				ribdServiceHandler.arpdRemoveRoute(route.OrigConfigObject.(RouteInfoRecord))
			}
		}
	}
}
