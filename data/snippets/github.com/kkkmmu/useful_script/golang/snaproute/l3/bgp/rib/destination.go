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

// destination.go
package rib

import (
	"bgpd"
	"fmt"
	"l3/bgp/config"
	"l3/bgp/packet"
	"math"
	"net"
	"sort"
	"strconv"
	"utils/logging"
)

const BGP_INTERNAL_PREF = 100
const BGP_EXTERNAL_PREF = 100

type PathAndRoute struct {
	Path
}

type Destination struct {
	rib               *LocRib
	logger            *logging.Writer
	gConf             *config.GlobalConfig
	NLRI              packet.NLRI
	protoFamily       uint32
	peerPathMap       map[string]map[uint32]*Path
	LocRibPath        *Path
	LocRibPathRoute   *Route
	aggPath           *Path
	aggregatedDestMap map[string]*Destination
	ecmpPaths         map[*Path]*Route
	pathRouteMap      map[*Path]*Route
	AddPaths          []*Path
	maxPathId         uint32
	pathIds           []uint32
	recalculate       bool
	BGPRouteState     config.ModelRouteIntf
	PathInfoRouteMap  map[*bgpd.PathInfo]*Route
	routeListIdx      int
}

func NewDestination(rib *LocRib, nlri packet.NLRI, protoFamily uint32, gConf *config.GlobalConfig) *Destination {
	dest := &Destination{
		rib:               rib,
		logger:            rib.logger,
		gConf:             gConf,
		NLRI:              nlri,
		protoFamily:       protoFamily,
		peerPathMap:       make(map[string]map[uint32]*Path),
		ecmpPaths:         make(map[*Path]*Route),
		aggregatedDestMap: make(map[string]*Destination),
		pathRouteMap:      make(map[*Path]*Route),
		AddPaths:          make([]*Path, 0),
		maxPathId:         1,
		pathIds:           make([]uint32, 0),
		routeListIdx:      -1,
		PathInfoRouteMap:  make(map[*bgpd.PathInfo]*Route),
	}

	dest.setBGPRouteState(protoFamily, nlri.GetPrefix().String(), int16(nlri.GetLength()))
	return dest
}

func (d *Destination) setBGPRouteState(protoFamily uint32, network string, cidrLen int16) {
	afi, _ := packet.GetAfiSafi(protoFamily)
	if afi == packet.AfiIP6 {
		d.BGPRouteState = NewIPv6Route(network, cidrLen)
	} else {
		d.BGPRouteState = NewIPv4Route(network, cidrLen)
	}
}

func (d *Destination) GetLocRibPathRoute() *Route {
	d.logger.Infof("GetLocRibPathRoute for %s", d.NLRI.GetCIDR())
	return d.LocRibPathRoute
}

func (d *Destination) GetBGPRoute() config.ModelRouteIntf {
	return d.BGPRouteState
}

func (d *Destination) GetPathRoute(path *Path) *Route {
	if route, ok := d.pathRouteMap[path]; ok {
		return route
	}

	return nil
}

func (d *Destination) GetProtocolFamily() uint32 {
	return d.protoFamily
}

func (d *Destination) String() string {
	return d.NLRI.String()
}

func (d *Destination) IsEmpty() bool {
	return len(d.peerPathMap) == 0
}

func (d *Destination) getNextPathId() uint32 {
	var pathId uint32
	if len(d.pathIds) > 0 {
		pathId = d.pathIds[len(d.pathIds)-1]
		d.pathIds = d.pathIds[:len(d.pathIds)-1]
		return pathId
	}

	pathId = d.maxPathId
	d.maxPathId++
	return pathId
}

func (d *Destination) releasePathId(pathId uint32) {
	if pathId+1 == d.maxPathId {
		d.maxPathId--
		return
	}

	d.pathIds = append(d.pathIds, pathId)
}

func (d *Destination) updateAddPaths(addPaths []*Path) (modified bool) {
	if len(d.AddPaths) != len(addPaths) {
		modified = true
	} else {
		for i := 0; i < len(d.AddPaths); i++ {
			if d.AddPaths[i] != addPaths[i] {
				modified = true
			}
		}
	}
	if modified {
		for i := 0; i < len(d.AddPaths); i++ {
			if route, ok := d.pathRouteMap[d.AddPaths[i]]; ok {
				route.ResetAdditionalPath()
			}
			d.AddPaths[i] = nil
		}
		d.AddPaths = addPaths
		for _, path := range d.AddPaths {
			if route, ok := d.pathRouteMap[path]; ok {
				route.SetAdditionalPath()
			}
		}
	}
	return modified
}

func (d *Destination) getPathForIP(peerIP string, pathId uint32) (path *Path) {
	if pathMap, ok := d.peerPathMap[peerIP]; ok {
		path = pathMap[pathId]
	}
	return path
}

func (d *Destination) getPathIdForPath(path *Path) (uint32, bool) {
	for _, pathMap := range d.peerPathMap {
		for pathId, peerPath := range pathMap {
			if path == peerPath {
				return pathId, true
			}
		}
	}

	d.logger.Errf("Destination:getPathIdForPath - path id not found for path %v", path)
	return 0, false
}

func (d *Destination) setUpdateAggPath(peerIP string, pathId uint32) {
	pathMap, ok := d.peerPathMap[peerIP]
	if !ok {
		d.logger.Errf("setUpdateAggPath - peer ip %s not found in peer path map", peerIP)
	} else {
		path, ok := pathMap[pathId]
		if !ok {
			d.logger.Errf("setUpdateAggPath - pathId %d not found in peer %s path map", pathId, peerIP)
		} else if d.LocRibPath == nil || path == d.LocRibPath ||
			getRouteSource(d.LocRibPath.routeType) >= getRouteSource(path.routeType) {
			d.recalculate = true
		}
	}

	if d.LocRibPath == nil {
		d.recalculate = true
	}
}

func (d *Destination) setAggPath(path *Path) {
	d.aggPath = path
}

func (d *Destination) addAggregatedDests(peerIP string, dest *Destination) {
	d.aggregatedDestMap[peerIP] = dest
}

func (d *Destination) removeAggregatedDests(peerIP string) {
	delete(d.aggregatedDestMap, peerIP)
}

func (d *Destination) AddOrUpdatePath(peerIp string, pathId uint32, path *Path) bool {
	var pathMap map[uint32]*Path
	added := false
	ok := false
	idx := -1

	d.logger.Infof("AddOrUpdatePath: Destination %s peerIP %s pathId %d, path %v", d.NLRI.GetCIDR(), peerIp, pathId, path)
	if pathMap, ok = d.peerPathMap[peerIp]; !ok {
		d.peerPathMap[peerIp] = make(map[uint32]*Path)
	}

	if oldPath, ok := pathMap[pathId]; ok {
		d.logger.Infof("Destination %s Update path from %s, id %d", d.NLRI.GetPrefix(), peerIp, pathId)
		if oldPath == path {
			d.logger.Infof("Destination %s new path from %s, id %d is the same as the old path", d.NLRI.GetPrefix(),
				peerIp, pathId)
			return added
		}
		if route, ok := d.pathRouteMap[oldPath]; ok {
			idx = route.routeListIdx
			delete(d.PathInfoRouteMap, route.PathInfo)
		}
		if d.LocRibPath == oldPath {
			d.LocRibPath = nil
		}
	} else {
		d.logger.Infof("Destination %s New path from %s, id %d", d.NLRI.GetPrefix(), peerIp, pathId)
		added = true
	}

	if d.LocRibPath == nil ||
		getRouteSource(d.LocRibPath.routeType) >= getRouteSource(path.routeType) {
		d.recalculate = true
	}

	outPathId := d.getNextPathId()
	route := NewRoute(d, path, RouteActionNone, pathId, outPathId)
	d.pathRouteMap[path] = route
	if idx != -1 {
		d.BGPRouteState.SetPath(route.PathInfo, idx)
	} else {
		idx = len(d.BGPRouteState.GetPaths())
		d.BGPRouteState.AppendPath(route.PathInfo)
	}
	d.PathInfoRouteMap[route.PathInfo] = route
	route.setIdx(idx)
	d.peerPathMap[peerIp][pathId] = path
	return added
}

func (d *Destination) RemovePath(peerIP string, pathId uint32, path *Path) *Path {
	var pathMap map[uint32]*Path
	var oldPath *Path
	ok := false
	if pathMap, ok = d.peerPathMap[peerIP]; !ok {
		d.logger.Errf("Destination %s Path not found from peer %s", d.NLRI.GetCIDR(), peerIP)
		return oldPath
	}

	if oldPath, ok = pathMap[pathId]; ok {
		for ecmpPath, _ := range d.ecmpPaths {
			if ecmpPath == oldPath {
				d.recalculate = true
			}
		}

		if d.LocRibPath == oldPath {
			d.recalculate = true
			d.LocRibPath = nil
		}

		route := d.pathRouteMap[oldPath]
		d.releasePathId(route.OutPathId)
		delete(d.pathRouteMap, oldPath)
		if route.routeListIdx != -1 {
			newPath := d.BGPRouteState.GetLastPath()
			if newRoute, ok := d.PathInfoRouteMap[newPath]; ok {
				delete(d.PathInfoRouteMap, route.PathInfo)
				d.BGPRouteState.RemovePathAndSetLast(route.routeListIdx)
				newRoute.setIdx(route.routeListIdx)
			} else {
				d.logger.Errf("Could not find path %v in PathInfoRouteMap %v", d.BGPRouteState.GetLastPath(),
					d.PathInfoRouteMap)
			}
		}
		delete(d.peerPathMap[peerIP], pathId)
		if len(d.peerPathMap[peerIP]) == 0 {
			delete(d.peerPathMap, peerIP)
		}
	} else {
		d.logger.Err("Destination", d.NLRI.GetCIDR(), "Path with path id", pathId, "not found from peer",
			peerIP)
	}
	return oldPath
}

func (d *Destination) RemoveAllPaths(peerIP string, path *Path) {
	var pathMap map[uint32]*Path
	ok := false
	if pathMap, ok = d.peerPathMap[peerIP]; !ok {
		d.logger.Err("Can't remove paths for", d.NLRI.GetCIDR(), "peer not found", peerIP)
		return
	}

	d.logger.Info("Remove all paths for", d.NLRI.GetCIDR(), "from peer", peerIP)
	for pathId, _ := range pathMap {
		d.logger.Info("Remove path id", pathId, "from peer", peerIP)
		d.RemovePath(peerIP, pathId, path)
	}
}

func (d *Destination) RemoveAllNeighborPaths() {
	for peerIP, pathMap := range d.peerPathMap {
		for pathId, path := range pathMap {
			if path.NeighborConf != nil {
				delete(d.peerPathMap[peerIP], pathId)
				if len(d.peerPathMap[peerIP]) == 0 {
					delete(d.peerPathMap, peerIP)
				}
			}
		}
	}

	if d.LocRibPath != nil {
		if d.LocRibPath.NeighborConf != nil {
			d.recalculate = true
			d.LocRibPath = nil
		}
	}
}

func (d *Destination) constructNetmaskFromLen(ones, bits int) net.IP {
	ip := make(net.IP, bits/8)
	bytes := ones / 8
	i := 0
	for ; i < bytes; i++ {
		ip[i] = 255
	}
	rem := ones % 8
	if rem != 0 {
		ip[i] = (255 << uint(8-rem))
	}
	return ip
}

func (d *Destination) removeAndPrepend(pathsList *[][]*Path, item *Path) {
	idx := 0
	found := false
	var paths []*Path

	for idx, paths = range *pathsList {
		var path *Path
		pathIdx := 0
		for pathIdx, path = range paths {
			if path == item {
				found = true
				break
			}
		}
		if found {
			copy(paths[1:pathIdx+1], paths[:pathIdx])
			paths[0] = path
			break
		}
	}

	if !found {
		paths = make([]*Path, 1)
		paths[0] = item
		*pathsList = append(*pathsList, paths)
	}
	copy((*pathsList)[1:idx+1], (*pathsList)[:idx])
	(*pathsList)[0] = paths
}

func (d *Destination) ConstructRouteConfig(path *Path, reachInfo *ReachabilityInfo, ipLength int) *config.RouteConfig {
	protocol := "IBGP"
	if path.IsExternal() {
		protocol = "EBGP"
	}
	nullRoute := false
	if path.IsAggregate() {
		nullRoute = true
	}
	isIPv6 := false
	if ipLength == 16 {
		isIPv6 = true
	}

	cfg := config.RouteConfig{
		Cost:              int32(reachInfo.Metric),
		IntfType:          int32(reachInfo.NextHopIfType),
		Protocol:          protocol,
		NextHopIp:         reachInfo.NextHop,
		NetworkMask:       d.constructNetmaskFromLen(int(d.NLRI.GetLength()), ipLength*8).String(),
		DestinationNw:     d.NLRI.GetPrefix().String(),
		OutgoingInterface: strconv.Itoa(int(reachInfo.NextHopIfIdx)),
		IsIPv6:            isIPv6,
		NullRoute:         nullRoute,
	}

	return &cfg
}

func (d *Destination) SelectRouteForLocRib(addPathCount int) (RouteAction, bool, []*Route, []*Route, []*Route) {
	updatedPaths := make([]*Path, 0)
	removedPaths := make([]*Path, 0)
	addedRoutes := make([]*Route, 0)
	updatedRoutes := make([]*Route, 0)
	deletedRoutes := make([]*Route, 0)
	createRibRoutes := make([]*Path, 0)
	routeSrc := RouteSrcUnknown
	locRibAction := RouteActionNone
	addPathsUpdated := false
	ipLength := packet.GetAddressLengthForFamily(d.protoFamily)

	d.logger.Infof("Destination - selecting best path for prefix %s", d.NLRI.GetPrefix())
	if !d.recalculate {
		return locRibAction, addPathsUpdated, addedRoutes, updatedRoutes, deletedRoutes
	}
	d.recalculate = false

	if d.LocRibPath != nil {
		var peerIP string
		if d.LocRibPath.NeighborConf != nil {
			peerIP = d.LocRibPath.NeighborConf.Neighbor.NeighborAddress.String()
		} else {
			peerIP = d.gConf.RouterId.String()
		}
		routeSrc = getRouteSource(d.LocRibPath.routeType)
		updatedPaths = append(updatedPaths, d.LocRibPath)
		d.logger.Infof("Destination %s Add loc rib path %v from %s to path selection, source=%d",
			d.NLRI.GetPrefix(), d.LocRibPath, peerIP, routeSrc)
	}

	for peerIP, pathMap := range d.peerPathMap {
		for _, path := range pathMap {
			if d.LocRibPath == nil || d.LocRibPath != path {
				if !path.IsReachable(d.protoFamily) {
					d.logger.Infof("Destination %s peer %s, NEXT_HOP %s is not reachable", d.NLRI.GetPrefix(), peerIP,
						path.GetNextHop(d.protoFamily))
					continue
				}

				if path.HasASLoop() {
					d.logger.Infof("Destination %s peer %s, path has AS %d loop", d.NLRI.GetPrefix(),
						peerIP, path.NeighborConf.RunningConf.LocalAS)
					continue
				}

				currPathSource := getRouteSource(path.routeType)
				if currPathSource > routeSrc {
					removedPaths = append(removedPaths, path)
					continue
				} else if currPathSource < routeSrc {
					if len(updatedPaths) > 0 {
						removedPaths = append(removedPaths, updatedPaths...)
						// For garbage collection
						for i := 0; i < len(updatedPaths); i++ {
							updatedPaths[i] = nil
						}
						updatedPaths[0] = path
						updatedPaths = updatedPaths[:1]
					} else {
						updatedPaths = append(updatedPaths, path)
					}
					d.logger.Info("Destination", d.NLRI.GetPrefix(), "route from", peerIP,
						"is from a better source type, old type=", routeSrc, "new type=", currPathSource, "path=", path)
					routeSrc = currPathSource
					continue
				} else {
					updatedPaths = append(updatedPaths, path)
				}
			}
		}
	}

	d.logger.Infof("Destination %s, ECMP routes %v updated paths %v", d.NLRI.GetPrefix(), d.ecmpPaths, updatedPaths)
	firstRoute := true
	if len(updatedPaths) > 0 {
		var ecmpPaths [][]*Path
		var addPaths []*Path
		if len(updatedPaths) > 1 || (addPathCount > 0) {
			d.logger.Infof("Found multiple paths with same pref, run path selection algorithm")
			if d.gConf.UseMultiplePaths {
				updatedPaths, ecmpPaths, addPaths =
					d.calculateBestPath(updatedPaths, removedPaths, d.gConf.EBGPMaxPaths > 1, d.gConf.IBGPMaxPaths > 1,
						addPathCount)
			} else {
				updatedPaths, ecmpPaths, addPaths = d.calculateBestPath(updatedPaths, removedPaths, false, false,
					addPathCount)
			}
		}

		if len(updatedPaths) > 1 {
			d.logger.Errf("Have more than one route after the tie breaking rules... using the first one, routes[%s]",
				updatedPaths)
		}

		d.logger.Info("before mod, ecmpPaths =", ecmpPaths)
		addPathsUpdated = d.updateAddPaths(addPaths)
		d.removeAndPrepend(&ecmpPaths, updatedPaths[0])
		d.logger.Info("after mod, ecmpPaths =", ecmpPaths)

		for idx, paths := range ecmpPaths {
			found := false
			for pathIdx, path := range paths {
				// If the first path (best path) in the first sub list is not already installed, break out
				if idx == 0 && pathIdx > 0 {
					break
				}
				if route, ok := d.ecmpPaths[path]; ok {
					// Update path
					d.logger.Infof("Destination %s path %v at [%d][%d] found in ecmp paths %v", d.NLRI.GetPrefix(),
						path, idx, pathIdx, d.ecmpPaths)
					found = true
					firstRoute = false
					if (idx == 0) && path.IsAggregate() {
						locRibAction = RouteActionReplace
					}
					updatedRoutes = append(updatedRoutes, route)
					route.setAction(RouteActionReplace)
					if idx != 0 {
						route.ResetBestPath()
					}
					break
				}
			}

			if !found {
				// Add route
				newRoute := d.pathRouteMap[paths[0]]
				if newRoute == nil {
					d.logger.Infof("Destination %s path %v NOT found in path route map %v", d.NLRI.GetPrefix(),
						paths[0], d.pathRouteMap)
					continue
				}
				newRoute.setAction(RouteActionAdd)
				newRoute.SetMultiPath()

				if !paths[0].IsLocal() || paths[0].IsAggregate() {
					d.logger.Infof("Add route for ip=%s, mask=%s, next hop=%s", d.NLRI.GetPrefix(),
						d.constructNetmaskFromLen(int(d.NLRI.GetLength()), ipLength*8),
						paths[0].GetReachability(d.protoFamily).NextHop)
					createRibRoutes = append(createRibRoutes, paths[0])
				}
				if idx == 0 {
					locRibAction = RouteActionAdd
					newRoute.SetBestPath()
				}
				d.ecmpPaths[paths[0]] = newRoute
				addedRoutes = append(addedRoutes, newRoute)
			}
		}

		d.LocRibPath = ecmpPaths[0][0]
		d.LocRibPathRoute = d.ecmpPaths[d.LocRibPath]
		d.logger.Infof("Destination %s loc rib path %v route %v, d.ecmpPaths %v ecmpPaths %v",
			d.NLRI.GetPrefix(), d.LocRibPath, d.LocRibPathRoute, d.ecmpPaths, ecmpPaths)
	} else {
		// Remove route
		for _, route := range d.ecmpPaths {
			route.setAction(RouteActionDelete)
		}
		locRibAction = RouteActionDelete
		d.LocRibPath = nil
	}

	for path, route := range d.ecmpPaths {
		if route.action == RouteActionNone || route.action == RouteActionDelete {
			if !path.IsLocal() || path.IsAggregate() {
				reachInfo := path.GetReachability(d.protoFamily)
				d.logger.Info("Remove route from ECMP paths, route =", route, "ip =",
					d.NLRI.GetCIDR(), "next hop =", reachInfo.NextHop)
				cfg := d.ConstructRouteConfig(path, reachInfo, ipLength)
				d.rib.routeMgr.UpdateRoute(cfg, "remove")
				d.logger.Info("DeleteV4Route from ECMP paths, route =", route, "ip =",
					d.NLRI.GetCIDR(), "next hop =", reachInfo.NextHop, "DONE")
			}
			route.ResetBestPath()
			route.ResetMultiPath()
			deletedRoutes = append(deletedRoutes, route)
			delete(d.ecmpPaths, path)
		} else {
			route.setAction(RouteActionNone)
		}
	}

	for _, path := range createRibRoutes {
		reachInfo := path.GetReachability(d.protoFamily)
		d.logger.Infof("Add route for ip=%s, mask=%s, next hop=%s", d.NLRI.GetCIDR(),
			d.constructNetmaskFromLen(int(d.NLRI.GetLength()), ipLength*8).String(), reachInfo.NextHop)
		cfg := d.ConstructRouteConfig(path, reachInfo, ipLength)
		if firstRoute {
			d.rib.routeMgr.CreateRoute(cfg)
			firstRoute = false
		} else {
			d.rib.routeMgr.UpdateRoute(cfg, "add")
		}
	}
	return locRibAction, addPathsUpdated, addedRoutes, updatedRoutes, deletedRoutes
}

func (d *Destination) getRoutesWithHighestPref(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	maxPref := uint32(0)
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths)
	idx := 0

	for i := 0; i < n; i++ {
		currPref := updatedPaths[i].GetPreference()
		from := ""
		if updatedPaths[i].NeighborConf != nil {
			from = updatedPaths[i].NeighborConf.Neighbor.NeighborAddress.String()
		} else {
			from = d.gConf.RouterId.String()
		}
		d.logger.Infof("Destination %s path pref %d from %s", d.NLRI.GetPrefix(), currPref, from)
		if currPref < maxPref {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if currPref > maxPref {
			d.logger.Infof("Destination %s route from %s has more preference, old pref=%d, new pref=%d",
				d.NLRI.GetPrefix(), from, maxPref, currPref)
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			maxPref = currPref
			updatedPaths[0] = updatedPaths[i]
			idx = 1
		} else if currPref == maxPref {
			d.logger.Infof("Destination %s route from %s has same preference, pref=%d", d.NLRI.GetPrefix(), from,
				maxPref)
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByPref{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func (d *Destination) getRoutesWithSmallestAS(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	minASNums := uint32(4096)
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths)
	idx := 0

	for i := 0; i < n; i++ {
		d.logger.Info("get num ASes from path", updatedPaths[i])
		asNums := updatedPaths[i].GetNumASes()
		from := ""
		if updatedPaths[i].NeighborConf != nil {
			from = updatedPaths[i].NeighborConf.Neighbor.NeighborAddress.String()
		}
		d.logger.Info("Dest =", d.NLRI.GetPrefix(), "number of ASes =", asNums, "from", from)
		if asNums > minASNums {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if asNums < minASNums {
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			minASNums = asNums
			updatedPaths[0] = updatedPaths[i]
			idx = 1
		} else if asNums == minASNums {
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: BySmallestAS{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func (d *Destination) getRoutesWithLowestOrigin(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	minOrigin := uint8(packet.BGPPathAttrOriginMax)
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths)
	idx := 0

	for i := 0; i < n; i++ {
		origin := updatedPaths[i].GetOrigin()
		if origin > minOrigin {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if origin < minOrigin {
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			minOrigin = origin
			updatedPaths[0] = updatedPaths[i]
			idx++
		} else if origin == minOrigin {
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByLowestOrigin{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func deleteIBGPRoutes(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path, []PathSortIface) {
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths) - 1
	i := 0

	for i <= n {
		if updatedPaths[i].NeighborConf.IsInternal() {
			removedPaths = append(removedPaths, updatedPaths[i])
			updatedPaths[i] = updatedPaths[n]
			updatedPaths[n] = nil
			n--
			continue
		}
		i++
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByIBGPOrEBGPRoutes{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	return updatedPaths[:i], prunedPaths
}

func (d *Destination) removeIBGPRoutesIfEBGPExist(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	for _, path := range updatedPaths {
		if path.NeighborConf != nil && path.NeighborConf.IsExternal() {
			return deleteIBGPRoutes(updatedPaths, prunedPaths)
		}
	}

	return updatedPaths, prunedPaths
}

func (d *Destination) isEBGPRoute(path *Path) bool {
	if path.NeighborConf != nil && path.NeighborConf.IsExternal() {
		return true
	}

	return false
}

func (d *Destination) isIBGPRoute(path *Path) bool {
	if path.NeighborConf != nil && path.NeighborConf.IsInternal() {
		return true
	}

	return false
}

func (d *Destination) getRoutesWithLowestBGPId(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths)
	lowestBGPId := uint32(math.MaxUint32)
	idx := 0

	for i := 0; i < n; i++ {
		bgpId := updatedPaths[i].GetBGPId()
		if bgpId > lowestBGPId {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if bgpId < lowestBGPId {
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			lowestBGPId = bgpId
			updatedPaths[0] = updatedPaths[i]
			idx = 1
		} else if bgpId == lowestBGPId {
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByLowestBGPId{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func (d *Destination) getRoutesWithShorterClusterLen(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	removedPaths := make([]*Path, 0)
	minClusterLen := uint16(math.MaxUint16)
	n := len(updatedPaths)
	idx := 0

	for i := 0; i < n; i++ {
		clusterLen := updatedPaths[i].GetNumClusters()
		if clusterLen > minClusterLen {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if clusterLen < minClusterLen {
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			minClusterLen = clusterLen
			updatedPaths[0] = updatedPaths[i]
			idx = 1
		} else if clusterLen == minClusterLen {
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByShorterClusterLen{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func CompareNeighborAddress(a net.IP, b net.IP) (int, error) {
	if len(a) != len(b) {
		return 0, config.AddressError{fmt.Sprintf("Address lenght not equal, Neighbor Address: %s, compare address: %s",
			a.String(), b.String())}
	}

	for i, val := range a {
		if val < b[i] {
			return -1, nil
		} else if val > b[i] {
			return 1, nil
		}
	}

	return 0, nil
}

func (d *Destination) getRoutesWithLowestPeerAddress(updatedPaths []*Path, prunedPaths []PathSortIface) ([]*Path,
	[]PathSortIface) {
	removedPaths := make([]*Path, 0)
	n := len(updatedPaths)
	idx := 0

	for i, path := range updatedPaths {
		val, err := CompareNeighborAddress(path.NeighborConf.Neighbor.NeighborAddress,
			updatedPaths[0].NeighborConf.Neighbor.NeighborAddress)
		if err != nil {
			d.logger.Errf("CompareNeighborAddress failed with %s", err)
		}

		if val > 0 {
			removedPaths = append(removedPaths, updatedPaths[i])
		} else if val < 0 {
			removedPaths = append(removedPaths, updatedPaths[:idx]...)
			updatedPaths[0] = updatedPaths[i]
			idx = 1
		} else if val == 0 {
			updatedPaths[idx] = updatedPaths[i]
			idx++
		}
	}

	if len(removedPaths) > 0 {
		pathSortIface := PathSortIface{
			paths: removedPaths,
			iface: ByLowestPeerAddress{removedPaths},
		}
		prunedPaths = append(prunedPaths, pathSortIface)
	}

	if idx > 0 {
		for i := idx; i < n; i++ {
			updatedPaths[i] = nil
		}
		updatedPaths = updatedPaths[:idx]
	}

	return updatedPaths, prunedPaths
}

func (d *Destination) getECMPPaths(updatedPaths []*Path) [][]*Path {
	ecmpPathMap := make(map[string][]*Path)

	for _, path := range updatedPaths {
		reachInfo := path.GetReachability(d.protoFamily)
		d.logger.Info("getECMPPaths: path =", path, "next hop =", reachInfo.NextHop)
		if _, ok := ecmpPathMap[reachInfo.NextHop]; !ok {
			ecmpPathMap[reachInfo.NextHop] = make([]*Path, 1)
			ecmpPathMap[reachInfo.NextHop][0] = path
		} else {
			ecmpPathMap[reachInfo.NextHop] = append(ecmpPathMap[reachInfo.NextHop], path)
		}
	}

	d.logger.Info("getECMPPaths: update paths =", updatedPaths, "ecmpPathsMap =", ecmpPathMap)
	ecmpPaths := make([][]*Path, 0)
	for _, paths := range ecmpPathMap {
		ecmpPaths = append(ecmpPaths, paths)
	}
	return ecmpPaths
}

func (d *Destination) addAddPaths(addPaths, currPaths []*Path, pathMap map[string]*Path) ([]*Path, map[string]*Path) {
	currPathMap := make(map[string]*Path)
	for _, path := range currPaths {
		reachInfo := path.GetReachability(d.protoFamily)
		if _, ok := pathMap[reachInfo.NextHop]; !ok {
			currPathMap[reachInfo.NextHop] = path
			pathMap[reachInfo.NextHop] = path
		}
	}

	d.logger.Info("getAddPaths: add paths =", addPaths, "pathMap =", pathMap)
	for _, path := range currPathMap {
		addPaths = append(addPaths, path)
	}
	return addPaths, pathMap
}

func (d *Destination) calculateBestPath(updatedPaths, removedPaths []*Path, ebgpMultiPath, ibgpMultiPath bool,
	addPathCount int) ([]*Path, [][]*Path, []*Path) {
	var ecmpPaths [][]*Path
	prunedPaths := make([]PathSortIface, 0)
	pathSortIface := PathSortIface{
		paths: removedPaths,
		iface: ByRouteSrc{removedPaths},
	}
	prunedPaths = append(prunedPaths, pathSortIface)

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithHighestPref, update paths =", updatedPaths)
		updatedPaths, prunedPaths = d.getRoutesWithHighestPref(updatedPaths, prunedPaths)
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithSmallestAS, update paths =", updatedPaths)
		updatedPaths, prunedPaths = d.getRoutesWithSmallestAS(updatedPaths, prunedPaths)
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithLowestOrigin, update paths =", updatedPaths)
		updatedPaths, prunedPaths = d.getRoutesWithLowestOrigin(updatedPaths, prunedPaths)
	}

	if (len(updatedPaths) > 1) && ebgpMultiPath && ibgpMultiPath {
		ecmpPaths = d.getECMPPaths(updatedPaths)
		d.logger.Info("calculateBestPath: IBGP & EBGP multi paths =", ecmpPaths)
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling removeIBGPRoutesIfEBGPExist, update paths =", updatedPaths)
		updatedPaths, prunedPaths = d.removeIBGPRoutesIfEBGPExist(updatedPaths, prunedPaths)
	}

	if len(updatedPaths) > 1 && ibgpMultiPath != ebgpMultiPath {
		if ebgpMultiPath && d.isEBGPRoute(updatedPaths[0]) {
			ecmpPaths = d.getECMPPaths(updatedPaths)
			d.logger.Infof("calculateBestPath: EBGP multi paths =", ecmpPaths)
		} else if ibgpMultiPath && d.isIBGPRoute(updatedPaths[0]) {
			ecmpPaths = d.getECMPPaths(updatedPaths)
			d.logger.Infof("calculateBestPath: IBGP multi paths =", ecmpPaths)
		}
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithLowestBGPId, update paths =", updatedPaths)
		updatedPaths, prunedPaths = d.getRoutesWithLowestBGPId(updatedPaths, prunedPaths)
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithShorterClusterLen")
		updatedPaths, prunedPaths = d.getRoutesWithShorterClusterLen(updatedPaths, prunedPaths)
	}

	if len(updatedPaths) > 1 {
		d.logger.Info("calling getRoutesWithLowestPeerAddress")
		updatedPaths, prunedPaths = d.getRoutesWithLowestPeerAddress(updatedPaths, prunedPaths)
	}

	pathMap := make(map[string]*Path)
	addPaths := make([]*Path, 0)
	if len(addPaths) < addPathCount && len(updatedPaths) > 1 {
		addPaths, pathMap = d.addAddPaths(addPaths, updatedPaths[1:], pathMap)
	}

	if len(addPaths) < addPathCount {
		for i := len(prunedPaths) - 1; i >= 0; i-- {
			sort.Sort(prunedPaths[i].iface)
			currPaths := prunedPaths[i].paths
			addPaths, pathMap = d.addAddPaths(addPaths, currPaths, pathMap)
			if len(addPaths) >= addPathCount {
				for idx := addPathCount; idx < len(addPaths); idx++ {
					addPaths[idx] = nil
				}
				addPaths = addPaths[:addPathCount]
				break
			}
		}
	}
	return updatedPaths, ecmpPaths, addPaths
}
