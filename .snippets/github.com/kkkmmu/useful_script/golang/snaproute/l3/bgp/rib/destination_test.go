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

// path_test.go
package rib

import (
	"l3/bgp/baseobjects"
	"l3/bgp/config"
	"l3/bgp/packet"
	"net"
	"testing"
	"utils/logging"
)

type RouteMgr struct {
	t *testing.T
}

func (r *RouteMgr) Start() {
	r.t.Log("RouteMgr:Start")
}

func (r *RouteMgr) GetNextHopInfo(ipAddr string, ifIndex int32) (*config.NextHopInfo, error) {
	nh := config.NextHopInfo{}
	nh.Metric = 0
	nh.NextHopIp = "30.1.1.1"
	nh.IsReachable = true
	nh.NextHopIfType = 1
	nh.NextHopIfIndex = 1
	return &nh, nil
}

func (r *RouteMgr) CreateRoute(route *config.RouteConfig) {
	r.t.Log("RouteMgr:CreateRoute:", route)
}
func (r *RouteMgr) DeleteRoute(route *config.RouteConfig) {
	r.t.Log("RouteMgr:DeleteRoute:", route)
}
func (r *RouteMgr) UpdateRoute(cfg *config.RouteConfig, op string) {
	r.t.Log("RouteMgr:UpdateRoute:", cfg, "operation:", op)
}

func (r *RouteMgr) ApplyPolicy(policy, conditions []*config.ApplyPolicyInfo) {
	r.t.Log("RouteMgr:ApplyPolicy")
}
func (r *RouteMgr) GetRoutes() (ri1 []*config.RouteInfo, ri2 []*config.RouteInfo) {
	r.t.Log("RouteMgr:GetRoutes")
	return ri1, ri2
}

func constructRibAndDest(t *testing.T, logger *logging.Writer, gConf *config.GlobalConfig) (*LocRib, *Destination) {
	routeMgr := &RouteMgr{t}
	locRib := NewLocRib(logger, routeMgr, nil, gConf)
	nlri := packet.NewExtNLRI(1001, packet.NewIPPrefix(net.ParseIP("20.1.10.0"), 24))
	protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
	dest := NewDestination(locRib, nlri, protoFamily, gConf)
	return locRib, dest
}

func getGlobalConf(localAS uint32) *config.GlobalConfig {
	gConf := &config.GlobalConfig{}
	gConf.AS = localAS
	gConf.RouterId = net.ParseIP("10.1.10.100")
	return gConf
}

func getNeighborConf(neighbor string, localAS, peerAS uint32) *config.NeighborConfig {
	pConf := &config.NeighborConfig{}
	pConf.NeighborAddress = net.ParseIP(neighbor)
	pConf.LocalAS = localAS
	pConf.PeerAS = peerAS
	return pConf
}

func TestDestination(t *testing.T) {
	logger := getLogger(t)
	gConf, _ := getConfObjects("192.168.0.100", uint32(1234), uint32(4321))
	_, dest := constructRibAndDest(t, logger, gConf)
	if dest != nil {
		t.Log("Destination successfully created")
	}
}

func TestAddOrUpdatePath(t *testing.T) {
	logger := getLogger(t)
	peerIP := "192.168.0.100"
	gConf, pConf := getConfObjects(peerIP, uint32(1234), uint32(4321))
	locRib, dest := constructRibAndDest(t, logger, gConf)

	// Add path with id 2 from neighbor1
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	pathAttrs := constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+1)
	path := NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 2, path)

	// Add path with id 2 from neighbor2
	peerIP2 := "172.16.0.1"
	pConf2 := getNeighborConf(peerIP2, 0, 5432)
	nConf2 := base.NewNeighborConf(logger, gConf, nil, *pConf2)
	pathAttrs2 := constructPathAttrs(pConf2.NeighborAddress, pConf2.PeerAS, pConf2.PeerAS+2)
	path2 := NewPath(locRib, nConf2, pathAttrs2, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP2, 2, path2)

	// Add path with id 1 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+3)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 1, path)

	// Replace path with id 2 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+4)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 2, path)
}

func TestRemovePath(t *testing.T) {
	logger := getLogger(t)
	peerIP := "192.168.0.100"
	gConf, pConf := getConfObjects(peerIP, uint32(1234), uint32(4321))
	locRib, dest := constructRibAndDest(t, logger, gConf)

	// Add path with id 2 from neighbor1
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	pathAttrs := constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+1)
	path := NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 2, path)

	// Add path with id 2 from neighbor2
	peerIP2 := "172.16.0.1"
	pConf2 := getNeighborConf(peerIP2, 0, 5432)
	nConf2 := base.NewNeighborConf(logger, gConf, nil, *pConf2)
	pathAttrs2 := constructPathAttrs(pConf2.NeighborAddress, pConf2.PeerAS, pConf2.PeerAS+2)
	path2 := NewPath(locRib, nConf2, pathAttrs2, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP2, 2, path2)

	// Remove path with id 2 from neighbor1
	dest.RemovePath(peerIP, 2, path)
	// Add path with id 2 from neighbor1 back
	dest.AddOrUpdatePath(peerIP, 2, path)

	// Add path with id 1 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+3)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 1, path)

	// Replace path with id 2 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+4)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	dest.AddOrUpdatePath(peerIP, 2, path)

	// Remove path with id 2 from neighbor1
	dest.RemovePath(peerIP, 2, path)
	// Remove path with id 1 from neighbor1
	dest.RemovePath(peerIP, 1, path)
	// Remove path with id 2 from neighbor2
	dest.RemovePath(peerIP2, 2, path2)
}

func TestRemoveAllPaths(t *testing.T) {
	logger := getLogger(t)
	peerIP := "192.168.0.100"
	gConf, pConf := getConfObjects(peerIP, uint32(1234), uint32(4321))
	locRib, dest := constructRibAndDest(t, logger, gConf)

	// Add path with id 2 from neighbor1
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	pathAttrs := constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+1)
	path := NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	reachInfo := NewReachabilityInfo("192.168.0.101", 0, 0, 0)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 2, path)

	// Add path with id 2 from neighbor2
	peerIP2 := "172.16.0.1"
	pConf2 := getNeighborConf(peerIP2, 0, 5432)
	nConf2 := base.NewNeighborConf(logger, gConf, nil, *pConf2)
	pathAttrs2 := constructPathAttrs(pConf2.NeighborAddress, pConf2.PeerAS, pConf2.PeerAS+2)
	path2 := NewPath(locRib, nConf2, pathAttrs2, nil, RouteTypeEGP)
	reachInfo2 := NewReachabilityInfo("172.16.0.2", 0, 0, 0)
	path.SetReachabilityForNextHop(pConf2.NeighborAddress.String(), reachInfo2)
	dest.AddOrUpdatePath(peerIP2, 2, path2)

	// Add path with id 1 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+3)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 1, path)

	// Replace path with id 2 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+4)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 2, path)

	dest.RemoveAllPaths(peerIP, path)
	dest.RemoveAllPaths(peerIP2, path2)
}

func TestSelectRouteForLocRib(t *testing.T) {
	logger := getLogger(t)
	peerIP := "192.168.0.100"
	gConf, pConf := getConfObjects(peerIP, uint32(1234), uint32(4321))
	locRib, dest := constructRibAndDest(t, logger, gConf)

	// Add path with id 2 from neighbor1
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	nConf.SetPeerAttrs(net.ParseIP(peerIP), 4, 3, 1, nil)
	pathAttrs := constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+1)
	path := NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	reachInfo := NewReachabilityInfo("192.168.0.101", 0, 0, 0)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 2, path)
	action, addPathsMod, _, _, _ := dest.SelectRouteForLocRib(2)
	t.Log("SelectRouteForLocRib returned action:", action, "addPaths updated:", addPathsMod)

	// Add path with id 2 from neighbor2
	peerIP2 := "172.16.0.1"
	pConf2 := getNeighborConf(peerIP2, 0, 5432)
	nConf2 := base.NewNeighborConf(logger, gConf, nil, *pConf2)
	nConf.SetPeerAttrs(net.ParseIP(peerIP2), 4, 3, 1, nil)
	pathAttrs2 := constructPathAttrs(pConf2.NeighborAddress, pConf2.PeerAS, pConf2.PeerAS+2)
	path2 := NewPath(locRib, nConf2, pathAttrs2, nil, RouteTypeEGP)
	reachInfo2 := NewReachabilityInfo("172.16.0.2", 0, 0, 0)
	path.SetReachabilityForNextHop(pConf2.NeighborAddress.String(), reachInfo2)
	dest.AddOrUpdatePath(peerIP2, 2, path2)
	action, addPathsMod, _, _, _ = dest.SelectRouteForLocRib(2)
	t.Log("SelectRouteForLocRib returned action:", action, "addPaths updated:", addPathsMod)

	// Add path with id 1 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+3)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 1, path)
	action, addPathsMod, _, _, _ = dest.SelectRouteForLocRib(2)
	t.Log("SelectRouteForLocRib returned action:", action, "addPaths updated:", addPathsMod)

	// Replace path with id 2 from neighbor1
	pathAttrs = constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+4)
	path = NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	path.SetReachabilityForNextHop(pConf.NeighborAddress.String(), reachInfo)
	dest.AddOrUpdatePath(peerIP, 2, path)
	action, addPathsMod, _, _, _ = dest.SelectRouteForLocRib(2)
	t.Log("SelectRouteForLocRib returned action:", action, "addPaths updated:", addPathsMod)
}
