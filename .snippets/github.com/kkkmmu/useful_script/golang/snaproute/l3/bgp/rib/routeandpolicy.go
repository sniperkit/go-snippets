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

// route.go
package rib

import (
	"l3/bgp/packet"
	"net"
)

type AdjRIBDir int

const (
	AdjRIBDirIn AdjRIBDir = iota
	AdjRIBDirOut
)

type AdjRIBRouteBase struct {
	Neighbor       net.IP
	ProtocolFamily uint32
	NLRI           packet.NLRI
}

func NewAdjRIBRouteBase(neighbor net.IP, protoFamily uint32, nlri packet.NLRI) *AdjRIBRouteBase {
	return &AdjRIBRouteBase{
		Neighbor:       neighbor,
		ProtocolFamily: protoFamily,
		NLRI:           nlri,
	}
}

type AdjRIBPathIdRoute struct {
	*AdjRIBRouteBase
	PathId           uint32
	Path             *Path
	PolicyList       []string
	PolicyStmt       string
	PolicyHitCounter int
	Accept           bool
}

func NewAdjRIBPathIdRoute(adjRIBRouteBase *AdjRIBRouteBase, pathId uint32, path *Path) *AdjRIBPathIdRoute {
	return &AdjRIBPathIdRoute{
		AdjRIBRouteBase:  adjRIBRouteBase,
		PathId:           pathId,
		Path:             path,
		PolicyList:       make([]string, 0),
		PolicyStmt:       "",
		PolicyHitCounter: 0,
		Accept:           false,
	}
}

type AdjRIBRoute struct {
	*AdjRIBRouteBase
	PathIdRouteMap map[uint32]*AdjRIBPathIdRoute
}

func NewAdjRIBRoute(neighbor net.IP, protoFamily uint32, nlri packet.NLRI) *AdjRIBRoute {
	adjRIBRouteBase := NewAdjRIBRouteBase(neighbor, protoFamily, nlri)
	return &AdjRIBRoute{
		AdjRIBRouteBase: adjRIBRouteBase,
		PathIdRouteMap:  make(map[uint32]*AdjRIBPathIdRoute),
	}
}

func (a *AdjRIBRoute) AddPath(pathId uint32, path *Path) *AdjRIBPathIdRoute {
	adjRIBPathIdRoute := NewAdjRIBPathIdRoute(a.AdjRIBRouteBase, pathId, path)
	a.PathIdRouteMap[pathId] = adjRIBPathIdRoute
	return adjRIBPathIdRoute
}

func (a *AdjRIBRoute) RemovePath(pathId uint32) {
	delete(a.PathIdRouteMap, pathId)
}

func (a *AdjRIBRoute) GetPathIdRoute(pathId uint32) *AdjRIBPathIdRoute {
	if route, ok := a.PathIdRouteMap[pathId]; ok {
		return route
	}

	return nil
}

func (a *AdjRIBRoute) GetPath(pathId uint32) *Path {
	if route, ok := a.PathIdRouteMap[pathId]; ok {
		return route.Path
	}

	return nil
}

func (a *AdjRIBRoute) DoesPathsExist() bool {
	return len(a.PathIdRouteMap) != 0
}

func (a *AdjRIBRoute) GetPathMap() map[uint32]*Path {
	pathMap := make(map[uint32]*Path)
	for pathId, adjRIBPathIdRoute := range a.PathIdRouteMap {
		pathMap[pathId] = adjRIBPathIdRoute.Path
	}
	return pathMap
}

func (a *AdjRIBRoute) RemoveAllPaths() {
	a.PathIdRouteMap = nil
	a.PathIdRouteMap = make(map[uint32]*AdjRIBPathIdRoute)
}

type FilteredRoutes struct {
	Add    []packet.NLRI
	Remove []packet.NLRI
}

func NewFilteredRoutes() *FilteredRoutes {
	f := FilteredRoutes{
		Add:    make([]packet.NLRI, 0),
		Remove: make([]packet.NLRI, 0),
	}

	return &f
}
