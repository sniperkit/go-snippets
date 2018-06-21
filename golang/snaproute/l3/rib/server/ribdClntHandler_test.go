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

package server

import (
	"fmt"
	"testing"
)

func TestInitRtClntHandlerServer(t *testing.T) {
	fmt.Println("****Init ClntHandler Server****")
	StartTestServer()
	fmt.Println("****************")
}

func TestDeleteRoutesOfType(t *testing.T) {
	fmt.Println("****Test delete routes of type****")
	TestProcessLogicalIntfCreateEvent(t)
	TestIPv4IntfCreateEvent(t)
	TestProcessV4RouteCreateConfig(t)
	DeleteRoutesOfType("EBGP")
	DeleteRoutesOfType("STATIC")
	DeleteRoutesOfType("OSPF")
	fmt.Println("route reachability after delete routes")
	TestGetRouteReachability(t)
	TestProcessv4RouteDeleteConfig(t)
	fmt.Println("**********************************")
}
