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

// rib.go
package rib

import (
	"bgpd"
	_ "fmt"
	"l3/bgp/baseobjects"
	"l3/bgp/config"
	"l3/bgp/packet"
	"models/objects"
	"net"
	"sync"
	"time"
	"utils/logging"
	"utils/statedbclient"
)

var totalRoutes int

const ResetTime int = 120
const AggregatePathId uint32 = 0

type ReachabilityInfo struct {
	NextHop       string
	NextHopIfType int32
	NextHopIfIdx  int32
	Metric        int32
}

func NewReachabilityInfo(nextHop string, nhIfType, nhIfIdx, metric int32) *ReachabilityInfo {
	return &ReachabilityInfo{
		NextHop:       nextHop,
		NextHopIfType: nhIfType,
		NextHopIfIdx:  nhIfIdx,
		Metric:        metric,
	}
}

type LocRib struct {
	logger           *logging.Writer
	gConf            *config.GlobalConfig
	routeMgr         config.RouteMgrIntf
	stateDBMgr       statedbclient.StateDBClient
	destPathMap      map[uint32]map[string]*Destination
	reachabilityMap  map[string]*ReachabilityInfo
	unreachablePaths map[string]map[*Path]map[*Destination][]uint32
	routesCount      map[uint32]uint32
	routeList        map[uint32][]*Destination
	routeMutex       sync.RWMutex
	routeListDirty   map[uint32]bool
	activeGet        map[uint32]bool
	timer            map[uint32]*time.Timer
}

func NewLocRib(logger *logging.Writer, rMgr config.RouteMgrIntf, sDBMgr statedbclient.StateDBClient,
	gConf *config.GlobalConfig) *LocRib {
	rib := &LocRib{
		logger:           logger,
		gConf:            gConf,
		routeMgr:         rMgr,
		stateDBMgr:       sDBMgr,
		destPathMap:      make(map[uint32]map[string]*Destination),
		reachabilityMap:  make(map[string]*ReachabilityInfo),
		unreachablePaths: make(map[string]map[*Path]map[*Destination][]uint32),
		routesCount:      make(map[uint32]uint32),
		routeList:        make(map[uint32][]*Destination),
		routeListDirty:   make(map[uint32]bool),
		activeGet:        make(map[uint32]bool),
		routeMutex:       sync.RWMutex{},
		timer:            make(map[uint32]*time.Timer),
	}

	return rib
}

func isIpInList(prefixes []packet.NLRI, ip packet.NLRI) bool {
	for _, nlri := range prefixes {
		if nlri.GetPathId() == ip.GetPathId() &&
			nlri.GetPrefix().Equal(ip.GetPrefix()) {
			return true
		}
	}
	return false
}

func (l *LocRib) GetRoutesCount() map[uint32]uint32 {
	return l.routesCount
}

func (l *LocRib) GetReachabilityInfo(ipStr string) *ReachabilityInfo {
	if reachabilityInfo, ok := l.reachabilityMap[ipStr]; ok {
		return reachabilityInfo
	}

	l.logger.Infof("GetReachabilityInfo: Reachability info not cached for Next hop %s", ipStr)
	ribdReachabilityInfo, err := l.routeMgr.GetNextHopInfo(ipStr, -1)
	if err != nil {
		l.logger.Infof("NEXT_HOP[%s] is not reachable", ipStr)
		return nil
	}
	nextHop := ribdReachabilityInfo.NextHopIp
	if nextHop == "" || nextHop[0] == '0' || nextHop == "::" {
		l.logger.Infof("Next hop for %s is %s. Using %s as the next hop", ipStr, nextHop, ipStr)
		nextHop = ipStr
	}

	reachabilityInfo := NewReachabilityInfo(nextHop, ribdReachabilityInfo.NextHopIfType,
		ribdReachabilityInfo.NextHopIfIndex, ribdReachabilityInfo.Metric)
	l.reachabilityMap[ipStr] = reachabilityInfo
	return reachabilityInfo
}

func (l *LocRib) GetDestFromIPAndLen(protoFamily uint32, ip string, cidrLen uint32) *Destination {
	if nlriDestMap, ok := l.destPathMap[protoFamily]; ok {
		if dest, ok := nlriDestMap[ip]; ok {
			return dest
		}
	}

	return nil
}

func (l *LocRib) GetDest(nlri packet.NLRI, protoFamily uint32, createIfNotExist bool) (dest *Destination, ok bool) {
	nlriDestMap, ok := l.destPathMap[protoFamily]
	if ok || createIfNotExist {
		if !ok {
			l.destPathMap[protoFamily] = make(map[string]*Destination)
			nlriDestMap = l.destPathMap[protoFamily]
		}
		dest, ok = nlriDestMap[nlri.GetCIDR()]
		if !ok && createIfNotExist {
			dest = NewDestination(l, nlri, protoFamily, l.gConf)
			l.destPathMap[protoFamily][nlri.GetCIDR()] = dest
			l.addRoutesToRouteList(dest, protoFamily)
			if _, found := l.routesCount[protoFamily]; !found {
				l.routesCount[protoFamily] = 0
			}
			l.routesCount[protoFamily]++
		}
	}

	return dest, ok
}

func (l *LocRib) updateRibOutInfo(action RouteAction, addPathsMod bool, addRoutes, updRoutes, delRoutes []*Route,
	dest *Destination, updated map[uint32]map[*Path][]*Destination, withdrawn, updatedAddPaths []*Destination) (
	map[uint32]map[*Path][]*Destination, []*Destination, []*Destination) {
	if action == RouteActionAdd || action == RouteActionReplace {
		if _, ok := updated[dest.protoFamily]; !ok {
			updated[dest.protoFamily] = make(map[*Path][]*Destination)
		}
		if _, ok := updated[dest.protoFamily][dest.LocRibPath]; !ok {
			updated[dest.protoFamily][dest.LocRibPath] = make([]*Destination, 0)
		}
		updated[dest.protoFamily][dest.LocRibPath] = append(updated[dest.protoFamily][dest.LocRibPath], dest)
	} else if action == RouteActionDelete {
		withdrawn = append(withdrawn, dest)
	} else if addPathsMod {
		updatedAddPaths = append(updatedAddPaths, dest)
	}

	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) GetRouteStateConfigObj(route config.ModelRouteIntf) objects.ConfigObj {
	return route.GetModelObject()
}

func (l *LocRib) ProcessRoutes(peerIP string, add, rem []packet.NLRI, addPath, remPath *Path, addPathCount int,
	protoFamily uint32, updated map[uint32]map[*Path][]*Destination, withdrawn []*Destination,
	updatedAddPaths []*Destination) (map[uint32]map[*Path][]*Destination, []*Destination, []*Destination, bool) {
	addedAllPrefixes := true

	// process withdrawn routes
	for _, nlri := range rem {
		if !isIpInList(add, nlri) {
			l.logger.Info("Processing withdraw destination", nlri.GetCIDR())
			dest, ok := l.GetDest(nlri, protoFamily, false)
			if !ok {
				l.logger.Warning("Can't process withdraw field, Destination does not exist, Dest:",
					nlri.GetCIDR())
				continue
			}
			op := l.stateDBMgr.UpdateObject
			oldPath := dest.RemovePath(peerIP, nlri.GetPathId(), remPath)
			if oldPath != nil && !oldPath.IsReachable(dest.protoFamily) {
				nextHop := oldPath.GetNextHop(dest.protoFamily)
				if nextHop != nil {
					nextHopStr := nextHop.String()
					if _, ok := l.unreachablePaths[nextHopStr]; ok {
						if _, ok := l.unreachablePaths[nextHopStr][oldPath]; ok {
							if pathIds, ok := l.unreachablePaths[nextHopStr][oldPath][dest]; ok {
								for idx, pathId := range pathIds {
									if pathId == nlri.GetPathId() {
										l.unreachablePaths[nextHopStr][oldPath][dest][idx] = pathIds[len(pathIds)-1]
										l.unreachablePaths[nextHopStr][oldPath][dest] =
											l.unreachablePaths[nextHopStr][oldPath][dest][:len(pathIds)-1]
										break
									}
								}
								if len(l.unreachablePaths[nextHopStr][oldPath][dest]) == 0 {
									delete(l.unreachablePaths[nextHopStr][oldPath], dest)
								}
							}
							if len(l.unreachablePaths[nextHopStr][oldPath]) == 0 {
								delete(l.unreachablePaths[nextHopStr], oldPath)
							}
						}
						if len(l.unreachablePaths[nextHopStr]) == 0 {
							delete(l.unreachablePaths, nextHopStr)
						}
					}
				}
			}
			action, addPathsMod, addRoutes, updRoutes, delRoutes := dest.SelectRouteForLocRib(addPathCount)
			updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes,
				delRoutes, dest, updated, withdrawn, updatedAddPaths)

			if oldPath != nil && remPath != nil {
				if neighborConf := remPath.GetNeighborConf(); neighborConf != nil {
					l.logger.Infof("Decrement prefix count for destination %s from Peer %s",
						nlri.GetCIDR(), peerIP)
					neighborConf.DecrPrefixCount()
				}
			}
			if action == RouteActionDelete {
				if dest.IsEmpty() {
					op = l.stateDBMgr.DeleteObject
					l.removeRoutesFromRouteList(dest, protoFamily)
					delete(l.destPathMap[protoFamily], nlri.GetCIDR())
					l.routesCount[protoFamily]--
				}
			}
			op(l.GetRouteStateConfigObj(dest.GetBGPRoute()))
		} else {
			l.logger.Info("Can't withdraw destination", nlri.GetCIDR(),
				"Destination is part of NLRI in the UDPATE")
		}
	}

	nextHopStr := addPath.GetNextHop(protoFamily).String()
	for _, nlri := range add {
		l.logger.Info("Processing nlri", nlri.GetCIDR())
		op := l.stateDBMgr.UpdateObject
		dest, alreadyCreated := l.GetDest(nlri, protoFamily, true)
		if !alreadyCreated {
			op = l.stateDBMgr.AddObject
		}
		if oldPath := dest.getPathForIP(peerIP, nlri.GetPathId()); oldPath == nil && addPath.NeighborConf != nil {
			if !addPath.NeighborConf.CanAcceptNewPrefix() {
				l.logger.Infof("Max prefixes limit reached for peer %s, can't process %s", peerIP,
					nlri.GetCIDR())
				addedAllPrefixes = false
				continue
			}
			l.logger.Infof("Increment prefix count for destination %s from Peer %s", nlri.GetCIDR(), peerIP)
			addPath.NeighborConf.IncrPrefixCount()
		}

		dest.AddOrUpdatePath(peerIP, nlri.GetPathId(), addPath)
		if !addPath.IsReachable(protoFamily) {
			if _, ok := l.unreachablePaths[nextHopStr][addPath][dest]; !ok {
				l.unreachablePaths[nextHopStr][addPath][dest] = make([]uint32, 0)
			}

			l.unreachablePaths[nextHopStr][addPath][dest] = append(l.unreachablePaths[nextHopStr][addPath][dest],
				nlri.GetPathId())
			continue
		}

		action, addPathsMod, addRoutes, updRoutes, delRoutes := dest.SelectRouteForLocRib(addPathCount)
		updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes, delRoutes,
			dest, updated, withdrawn, updatedAddPaths)
		op(l.GetRouteStateConfigObj(dest.GetBGPRoute()))
	}

	return updated, withdrawn, updatedAddPaths, addedAllPrefixes
}

func (l *LocRib) ProcessRoutesForReachableRoutes(nextHop string, reachabilityInfo *ReachabilityInfo, addPathCount int,
	updated map[uint32]map[*Path][]*Destination, withdrawn []*Destination, updatedAddPaths []*Destination) (
	map[uint32]map[*Path][]*Destination, []*Destination, []*Destination) {
	if _, ok := l.unreachablePaths[nextHop]; ok {
		for path, destinations := range l.unreachablePaths[nextHop] {
			path.SetReachabilityForNextHop(nextHop, reachabilityInfo)
			peerIP := path.GetPeerIP()
			if peerIP == "" {
				l.logger.Errf("ProcessRoutesForReachableRoutes: nexthop %s peer ip not found for path %+v", nextHop,
					path)
				continue
			}

			for dest, pathIds := range destinations {
				l.logger.Info("Processing dest", dest.NLRI.GetCIDR())
				for _, pathId := range pathIds {
					dest.AddOrUpdatePath(peerIP, pathId, path)
				}
				action, addPathsMod, addRoutes, updRoutes, delRoutes := dest.SelectRouteForLocRib(addPathCount)
				updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes,
					delRoutes, dest, updated, withdrawn, updatedAddPaths)
				l.stateDBMgr.AddObject(l.GetRouteStateConfigObj(dest.GetBGPRoute()))
			}
		}
	}

	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) TestNHAndProcessRoutes(peerIP string, add, remove []packet.NLRI, addPath, remPath *Path,
	addPathCount int, protoFamily uint32, updated map[uint32]map[*Path][]*Destination, withdrawn,
	updatedAddPaths []*Destination) (map[uint32]map[*Path][]*Destination, []*Destination, []*Destination, bool) {
	var reachabilityInfo *ReachabilityInfo
	nextHopStr := ""
	if len(add) > 0 {
		nextHop := addPath.GetNextHop(protoFamily)
		if nextHop == nil {
			l.logger.Errf("RIB - Next hop not found for protocol family %d", protoFamily)
			return updated, withdrawn, updatedAddPaths, true
		}
		nextHopStr = nextHop.String()
		reachabilityInfo = l.GetReachabilityInfo(nextHopStr)
		addPath.SetReachabilityForFamily(protoFamily, reachabilityInfo)

		//addPath.GetReachabilityInfo()
		if !addPath.IsValid() {
			l.logger.Infof("Received a update with our cluster id %d, Discarding the update.",
				addPath.NeighborConf.RunningConf.RouteReflectorClusterId)
			return updated, withdrawn, updatedAddPaths, true
		}

		if reachabilityInfo == nil {
			l.logger.Infof("TestNHAndProcessRoutes - next hop %s is not reachable", nextHopStr)

			if _, ok := l.unreachablePaths[nextHopStr]; !ok {
				l.unreachablePaths[nextHopStr] = make(map[*Path]map[*Destination][]uint32)
			}

			if _, ok := l.unreachablePaths[nextHopStr][addPath]; !ok {
				l.unreachablePaths[nextHopStr][addPath] = make(map[*Destination][]uint32)
			}
		}
	}

	updated, withdrawn, updatedAddPaths, addedAllPrefixes := l.ProcessRoutes(peerIP, add, remove, addPath, remPath,
		addPathCount, protoFamily, updated, withdrawn, updatedAddPaths)

	if reachabilityInfo != nil {
		l.logger.Infof("TestNHAndProcessRoutes - next hop %s is reachable, so process previously unreachable routes",
			nextHopStr)
		updated, withdrawn, updatedAddPaths = l.ProcessRoutesForReachableRoutes(nextHopStr, reachabilityInfo,
			addPathCount, updated, withdrawn, updatedAddPaths)
	}
	return updated, withdrawn, updatedAddPaths, addedAllPrefixes
}

func (l *LocRib) ProcessUpdate(neighborConf *base.NeighborConf, path *Path, add, rem []packet.NLRI, protoFamily uint32,
	addPathCount int, updated map[uint32]map[*Path][]*Destination, withdrawn []*Destination,
	updatedAddPaths []*Destination) (map[uint32]map[*Path][]*Destination, []*Destination, []*Destination, bool) {
	//body := pktInfo.Msg.Body.(*packet.BGPUpdate)
	//	updated := make(map[uint32]map[*Path][]*Destination)
	//	withdrawn := make([]*Destination, 0)
	//	updatedAddPaths := make([]*Destination, 0)
	addedAllPrefixes := true
	remPath := path.Clone()

	//	mpReach, mpUnreach := packet.RemoveMPAttrs(&body.PathAttributes)
	//	remPath := NewPath(l, neighborConf, body.PathAttributes, mpReach, RouteTypeEGP)
	//	addPath := NewPath(l, neighborConf, body.PathAttributes, mpReach, RouteTypeEGP)

	if len(add) > 0 || len(rem) > 0 {
		//		protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.TestNHAndProcessRoutes(
			neighborConf.Neighbor.NeighborAddress.String(), add, rem, path, remPath, addPathCount, protoFamily,
			updated, withdrawn, updatedAddPaths)
	}

	//	reachNLRIDone := false
	//	if mpUnreach != nil {
	//		var reachNLRI []packet.NLRI
	//		if mpReach != nil && mpReach.AFI == mpUnreach.AFI && mpReach.SAFI == mpUnreach.SAFI {
	//			reachNLRIDone = true
	//			reachNLRI = mpReach.NLRI
	//		}
	//		protoFamily := packet.GetProtocolFamily(mpUnreach.AFI, mpUnreach.SAFI)
	//		updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.TestNHAndProcessRoutes(pktInfo.Src, reachNLRI,
	//			mpUnreach.NLRI, addPath, remPath, addPathCount, protoFamily, updated, withdrawn, updatedAddPaths)
	//	}

	//	if !reachNLRIDone && mpReach != nil {
	//		protoFamily := packet.GetProtocolFamily(mpReach.AFI, mpReach.SAFI)
	//		updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.TestNHAndProcessRoutes(pktInfo.Src, mpReach.NLRI,
	//			nil, addPath, remPath, addPathCount, protoFamily, updated, withdrawn, updatedAddPaths)
	//	}
	return updated, withdrawn, updatedAddPaths, addedAllPrefixes
}

func (l *LocRib) ProcessFilteredRoutes(neighborConf *base.NeighborConf,
	filteredRoutes map[*Path]map[uint32]*FilteredRoutes, addPathCount int) (map[uint32]map[*Path][]*Destination,
	[]*Destination, []*Destination, bool) {
	updated := make(map[uint32]map[*Path][]*Destination)
	withdrawn := make([]*Destination, 0)
	updatedAddPaths := make([]*Destination, 0)
	addedAllPrefixes := true

	for path, pfNLRIs := range filteredRoutes {
		for protoFamily, routes := range pfNLRIs {
			updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.TestNHAndProcessRoutes(
				neighborConf.Neighbor.NeighborAddress.String(), routes.Add, routes.Remove, path, path, addPathCount,
				protoFamily, updated, withdrawn, updatedAddPaths)
			if !addedAllPrefixes {
				break
			}
		}
		if !addedAllPrefixes {
			break
		}
	}

	return updated, withdrawn, updatedAddPaths, addedAllPrefixes
}

func (l *LocRib) ProcessConnectedRoutes(src string, path *Path, add, remove map[uint32][]packet.NLRI,
	addPathCount int) (map[uint32]map[*Path][]*Destination, []*Destination, []*Destination) {
	var removePath *Path
	var addedAllPrefixes bool
	removePath = path.Clone()
	updated := make(map[uint32]map[*Path][]*Destination)
	withdrawn := make([]*Destination, 0)
	updatedAddPaths := make([]*Destination, 0)

	for protoFamily, withdrawnNLRI := range remove {
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.ProcessRoutes(src, add[protoFamily], withdrawnNLRI,
			path, removePath, addPathCount, protoFamily, updated, withdrawn, updatedAddPaths)
		delete(add, protoFamily)
		if !addedAllPrefixes {
			l.logger.Errf("Failed to add connected routes... max prefixes exceeded for connected routes!")
		}
	}

	for protoFamily, updatedNLRI := range add {
		updated, withdrawn, updatedAddPaths, addedAllPrefixes = l.ProcessRoutes(src, updatedNLRI, nil, path,
			removePath, addPathCount, protoFamily, updated, withdrawn, updatedAddPaths)
		if !addedAllPrefixes {
			l.logger.Errf("Failed to add connected routes... max prefixes exceeded for connected routes!")
		}
	}
	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) RemoveUpdatesFromNeighbor(peerIP string, neighborConf *base.NeighborConf, addPathCount int) (
	map[uint32]map[*Path][]*Destination, []*Destination, []*Destination) {
	remPath := NewPath(l, neighborConf, nil, nil, RouteTypeEGP)
	withdrawn := make([]*Destination, 0)
	updated := make(map[uint32]map[*Path][]*Destination)
	updatedAddPaths := make([]*Destination, 0)

	for protoFamily, ipDestMap := range l.destPathMap {
		for destIP, dest := range ipDestMap {
			op := l.stateDBMgr.UpdateObject
			dest.RemoveAllPaths(peerIP, remPath)
			action, addPathsMod, addRoutes, updRoutes, delRoutes := dest.SelectRouteForLocRib(addPathCount)
			l.logger.Info("RemoveUpdatesFromNeighbor - dest", dest.NLRI.GetCIDR(),
				"SelectRouteForLocRib returned action", action, "addRoutes", addRoutes, "updRoutes", updRoutes,
				"delRoutes", delRoutes)
			updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes,
				delRoutes, dest, updated, withdrawn, updatedAddPaths)
			if action == RouteActionDelete && dest.IsEmpty() {
				l.logger.Info("All routes removed for dest", dest.NLRI.GetCIDR())
				l.removeRoutesFromRouteList(dest, protoFamily)
				delete(l.destPathMap[protoFamily], destIP)
				l.routesCount[protoFamily]--
				op = l.stateDBMgr.DeleteObject
			}
			op(l.GetRouteStateConfigObj(dest.GetBGPRoute()))
		}
	}

	if neighborConf != nil {
		neighborConf.SetPrefixCount(0)
	}
	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) RemoveUpdatesFromAllNeighbors(addPathCount int) {
	withdrawn := make([]*Destination, 0)
	updated := make(map[uint32]map[*Path][]*Destination)
	updatedAddPaths := make([]*Destination, 0)

	for protoFamily, ipDestMap := range l.destPathMap {
		for destIP, dest := range ipDestMap {
			op := l.stateDBMgr.UpdateObject
			dest.RemoveAllNeighborPaths()
			action, addPathsMod, addRoutes, updRoutes, delRoutes := dest.SelectRouteForLocRib(addPathCount)
			l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes, delRoutes, dest, updated, withdrawn,
				updatedAddPaths)
			if action == RouteActionDelete && dest.IsEmpty() {
				l.removeRoutesFromRouteList(dest, protoFamily)
				delete(l.destPathMap[protoFamily], destIP)
				l.routesCount[protoFamily]--
				op = l.stateDBMgr.DeleteObject
			}
			op(l.GetRouteStateConfigObj(dest.GetBGPRoute()))
		}
	}
}

func (l *LocRib) GetLocRib() map[uint32]map[*Path][]*Destination {
	updated := make(map[uint32]map[*Path][]*Destination)
	for protoFamily, ipDestMap := range l.destPathMap {
		if _, ok := updated[protoFamily]; !ok {
			updated[protoFamily] = make(map[*Path][]*Destination)
		}
		for _, dest := range ipDestMap {
			if dest.LocRibPath != nil {
				updated[protoFamily][dest.LocRibPath] = append(updated[protoFamily][dest.LocRibPath], dest)
			}
		}
	}

	return updated
}

func (l *LocRib) RemoveRouteFromAggregate(ip *packet.IPPrefix, aggIP *packet.IPPrefix, srcIP string,
	protoFamily uint32, bgpAgg *config.BGPAggregate, ipDest *Destination, addPathCount int) (
	map[uint32]map[*Path][]*Destination, []*Destination, []*Destination) {
	var aggPath *Path
	var dest *Destination
	var aggDest *Destination
	var ok bool
	withdrawn := make([]*Destination, 0)
	updated := make(map[uint32]map[*Path][]*Destination)
	updatedAddPaths := make([]*Destination, 0)

	l.logger.Infof("LocRib:RemoveRouteFromAggregate - ip %v, aggIP %v", ip, aggIP)
	if dest, ok = l.GetDest(ip, protoFamily, false); !ok {
		if ipDest == nil {
			l.logger.Info("RemoveRouteFromAggregate: routes ip", ip, "not found")
			return updated, withdrawn, nil
		}
		dest = ipDest
	}
	l.logger.Info("RemoveRouteFromAggregate: locRibPath", dest.LocRibPath, "locRibRoutePath",
		dest.LocRibPathRoute.path)
	op := l.stateDBMgr.UpdateObject

	if aggDest, ok = l.GetDest(aggIP, protoFamily, false); !ok {
		l.logger.Infof("LocRib:RemoveRouteFromAggregate - dest not found for aggIP %v", aggIP)
		return updated, withdrawn, nil
	}

	if aggPath = aggDest.getPathForIP(srcIP, AggregatePathId); aggPath == nil {
		l.logger.Infof("LocRib:RemoveRouteFromAggregate - path not found for dest, aggIP %v", aggIP)
		return updated, withdrawn, nil
	}

	aggPath.removePathFromAggregate(ip.Prefix.String(), bgpAgg.GenerateASSet)
	if aggPath.isAggregatePathEmpty() {
		aggDest.RemovePath(srcIP, AggregatePathId, aggPath)
	} else {
		aggDest.setUpdateAggPath(srcIP, AggregatePathId)
	}
	aggDest.removeAggregatedDests(ip.Prefix.String())
	action, addPathsMod, addRoutes, updRoutes, delRoutes := aggDest.SelectRouteForLocRib(addPathCount)
	updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes, delRoutes,
		aggDest, updated, withdrawn, updatedAddPaths)
	if action == RouteActionAdd || action == RouteActionReplace {
		dest.aggPath = aggPath
	}
	if action == RouteActionDelete && aggDest.IsEmpty() {
		l.removeRoutesFromRouteList(dest, protoFamily)
		delete(l.destPathMap[protoFamily], aggIP.Prefix.String())
		l.routesCount[protoFamily]--
		op = l.stateDBMgr.DeleteObject
	}
	op(l.GetRouteStateConfigObj(dest.GetBGPRoute()))

	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) AddRouteToAggregate(ip *packet.IPPrefix, aggIP *packet.IPPrefix, srcIP string, protoFamily uint32,
	ifaceIP net.IP, bgpAgg *config.BGPAggregate, addPathCount int) (map[uint32]map[*Path][]*Destination,
	[]*Destination, []*Destination) {
	var aggPath, path *Path
	var dest *Destination
	var aggDest *Destination
	var ok bool
	withdrawn := make([]*Destination, 0)
	updated := make(map[uint32]map[*Path][]*Destination)
	updatedAddPaths := make([]*Destination, 0)

	nextHop := packet.GetZeroNextHopForFamily(protoFamily)
	if nextHop == nil {
		l.logger.Info("AddRouteToAggregate: Did not find next hop for protocol family", protoFamily)
		return updated, withdrawn, nil
	}

	l.logger.Infof("LocRib:AddRouteToAggregate - ip %v, aggIP %v", ip, aggIP)
	if dest, ok = l.GetDest(ip, protoFamily, false); !ok {
		l.logger.Info("AddRouteToAggregate: routes ip", ip, "not found")
		return updated, withdrawn, nil
	}
	path = dest.LocRibPath

	op := l.stateDBMgr.UpdateObject
	if aggDest, ok = l.GetDest(aggIP, protoFamily, true); ok {
		aggPath = aggDest.getPathForIP(srcIP, AggregatePathId)
		l.logger.Infof("LocRib:AddRouteToAggregate - aggIP %v found in dest, agg path %v", aggIP, aggPath)
	}

	if aggPath != nil {
		l.logger.Infof("LocRib:AddRouteToAggregate - aggIP %v, agg path found, update path attrs", aggIP)
		aggPath.addPathToAggregate(ip.Prefix.String(), path, bgpAgg.GenerateASSet)
		aggDest.setUpdateAggPath(srcIP, AggregatePathId)
		aggDest.addAggregatedDests(ip.Prefix.String(), dest)
	} else {
		l.logger.Infof("LocRib:AddRouteToAggregate - aggIP %v, agg path NOT found, create new path", aggIP)
		op = l.stateDBMgr.AddObject
		pathAttrs := packet.ConstructPathAttrForAggRoutes(path.PathAttrs, bgpAgg.GenerateASSet)
		packet.SetNextHopPathAttrs(pathAttrs, net.IPv4zero)
		packet.SetPathAttrAggregator(pathAttrs, l.gConf.AS, l.gConf.RouterId)
		mpReachNLRI := packet.ConstructMPReachNLRIForAggRoutes(packet.GetProtocolFamily(packet.AfiIP6,
			packet.SafiUnicast))
		aggPath = NewPath(path.rib, nil, pathAttrs, mpReachNLRI, RouteTypeAgg)
		aggPath.setAggregatedPath(ip.Prefix.String(), path)
		aggDest, _ := l.GetDest(aggIP, protoFamily, true)
		aggDest.AddOrUpdatePath(srcIP, AggregatePathId, aggPath)
		aggDest.addAggregatedDests(ip.Prefix.String(), dest)
		reachabilityInfo := NewReachabilityInfo(nextHop.String(), 0, 0, 0)
		aggPath.SetReachabilityForFamily(protoFamily, reachabilityInfo)
	}

	action, addPathsMod, addRoutes, updRoutes, delRoutes := aggDest.SelectRouteForLocRib(addPathCount)
	updated, withdrawn, updatedAddPaths = l.updateRibOutInfo(action, addPathsMod, addRoutes, updRoutes, delRoutes,
		aggDest, updated, withdrawn, updatedAddPaths)
	if action == RouteActionAdd || action == RouteActionReplace {
		dest.aggPath = aggPath
	}

	op(l.GetRouteStateConfigObj(aggDest.GetBGPRoute()))
	return updated, withdrawn, updatedAddPaths
}

func (l *LocRib) removeRoutesFromRouteList(dest *Destination, protoFamily uint32) {
	defer l.routeMutex.Unlock()
	l.routeMutex.Lock()
	if _, ok := l.routeList[protoFamily]; !ok {
		l.logger.Err("Protocol family", protoFamily, "not found in RIB route list")
		return
	}

	idx := dest.routeListIdx
	if idx != -1 {
		l.logger.Info("removeRoutesFromRouteList: remove dest at idx", idx)
		if !l.activeGet[protoFamily] {
			l.routeList[protoFamily][idx] = l.routeList[protoFamily][len(l.routeList[protoFamily])-1]
			l.routeList[protoFamily][idx].routeListIdx = idx
			l.routeList[protoFamily][len(l.routeList[protoFamily])-1] = nil
			l.routeList[protoFamily] = l.routeList[protoFamily][:len(l.routeList[protoFamily])-1]
		} else {
			l.routeList[protoFamily][idx] = nil
			l.routeListDirty[protoFamily] = true
		}
	}
}

func (l *LocRib) addRoutesToRouteList(dest *Destination, protoFamily uint32) {
	defer l.routeMutex.Unlock()
	l.routeMutex.Lock()

	if _, ok := l.timer[protoFamily]; !ok {
		l.timer[protoFamily] = time.AfterFunc(time.Duration(100)*time.Second, func() { l.ResetRouteList(protoFamily) })
		l.timer[protoFamily].Stop()
	}

	l.routeList[protoFamily] = append(l.routeList[protoFamily], dest)
	l.logger.Info("addRoutesToRouteList: added dest at idx", len(l.routeList[protoFamily])-1)
	dest.routeListIdx = len(l.routeList[protoFamily]) - 1
}

func (l *LocRib) ResetAllRouteLists() {
	for protoFamily, _ := range l.routeList {
		l.ResetRouteList(protoFamily)
	}
}

func (l *LocRib) ResetRouteList(protoFamily uint32) {
	defer l.routeMutex.Unlock()
	l.routeMutex.Lock()

	l.activeGet[protoFamily] = false

	if !l.routeListDirty[protoFamily] {
		return
	}

	lastIdx := len(l.routeList[protoFamily]) - 1
	var modIdx, idx int
	for idx = 0; idx < len(l.routeList[protoFamily]); idx++ {
		if l.routeList[protoFamily][idx] == nil {
			for modIdx = lastIdx; modIdx > idx && l.routeList[protoFamily][modIdx] == nil; modIdx-- {
			}
			if modIdx <= idx {
				lastIdx = idx
				break
			}
			l.routeList[protoFamily][idx] = l.routeList[protoFamily][modIdx]
			l.routeList[protoFamily][idx].routeListIdx = idx
			l.routeList[protoFamily][modIdx] = nil
			lastIdx = modIdx
		}
	}
	l.routeList[protoFamily] = l.routeList[protoFamily][:idx]
	l.routeListDirty[protoFamily] = false
}

func (l *LocRib) GetBGPRoute(prefix string, protocolFamily uint32) interface{} {
	defer l.routeMutex.RUnlock()
	l.routeMutex.RLock()

	for pf, ipDestMap := range l.destPathMap {
		if protocolFamily == pf {
			if dest, ok := ipDestMap[prefix]; ok {
				return dest.GetBGPRoute().GetThriftObject()
			}
		}
	}

	return nil
}

func (l *LocRib) GetBGPv4Route(prefix string) *bgpd.BGPv4RouteState {
	route := l.GetBGPRoute(prefix, packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast))
	return route.(*bgpd.BGPv4RouteState)
}

func (l *LocRib) GetBGPv6Route(prefix string) *bgpd.BGPv6RouteState {
	route := l.GetBGPRoute(prefix, packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast))
	return route.(*bgpd.BGPv6RouteState)
}

func (l *LocRib) BulkGetBGPRoutes(index int, count int, protoFamily uint32) (int, int, []interface{}) {
	var i int
	n := 0
	result := make([]interface{}, count)

	if _, ok := l.timer[protoFamily]; !ok {
		l.logger.Err("BulkGetBGPRoutes - protocol family", protoFamily, "not found in routes list")
		return i, n, result
	}

	l.timer[protoFamily].Stop()
	if index == 0 && l.activeGet[protoFamily] {
		l.ResetRouteList(protoFamily)
	}
	l.activeGet[protoFamily] = true

	defer l.routeMutex.RUnlock()
	l.routeMutex.RLock()

	for i = index; i < len(l.routeList[protoFamily]) && n < count; i++ {
		if l.routeList[protoFamily][i] != nil && len(l.routeList[protoFamily][i].BGPRouteState.GetPaths()) > 0 {
			result[n] = l.routeList[protoFamily][i].GetBGPRoute().GetThriftObject()
			n++
		}
	}
	result = result[:n]

	if i >= len(l.routeList[protoFamily]) {
		i = 0
	}

	l.timer[protoFamily].Reset(time.Duration(ResetTime) * time.Second)
	return i, n, result
}

func (l *LocRib) BulkGetBGPv4Routes(index int, count int) (int, int, []*bgpd.BGPv4RouteState) {
	i, n, routes := l.BulkGetBGPRoutes(index, count, packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast))
	thriftRoutes := make([]*bgpd.BGPv4RouteState, len(routes))
	for idx, route := range routes {
		thriftRoutes[idx] = route.(*bgpd.BGPv4RouteState)
	}
	return i, n, thriftRoutes
}

func (l *LocRib) BulkGetBGPv6Routes(index int, count int) (int, int, []*bgpd.BGPv6RouteState) {
	i, n, routes := l.BulkGetBGPRoutes(index, count, packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast))
	thriftRoutes := make([]*bgpd.BGPv6RouteState, len(routes))
	for idx, route := range routes {
		thriftRoutes[idx] = route.(*bgpd.BGPv6RouteState)
	}
	return i, n, thriftRoutes
}
