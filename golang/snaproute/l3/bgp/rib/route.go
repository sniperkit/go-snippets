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
	"bgpd"
	"l3/bgp/packet"
	"time"
)

type RouteAction uint8

const (
	RouteActionNone RouteAction = iota
	RouteActionAdd
	RouteActionReplace
	RouteActionDelete
)

type Route struct {
	PathInfo         *bgpd.PathInfo
	Dest             *Destination
	path             *Path
	routeListIdx     int
	time             time.Time
	action           RouteAction
	OutPathId        uint32
	PolicyList       []string
	PolicyHitCounter int
}

func NewRoute(dest *Destination, path *Path, action RouteAction, inPathId, outPathId uint32) *Route {
	currTime := time.Now()
	pathInfo := &bgpd.PathInfo{
		NextHop:        path.GetNextHop(dest.protoFamily).String(),
		Metric:         int32(path.MED),
		LocalPref:      int32(path.LocalPref),
		Path:           path.GetAS4ByteList(),
		PathId:         int32(inPathId),
		UpdatedTime:    currTime.String(),
		ValidPath:      path.IsReachable(dest.protoFamily),
		BestPath:       false,
		MultiPath:      false,
		AdditionalPath: false,
		Origin:         packet.GetOriginTypeStr(path.GetOrigin()),
		PathType:       path.GetSourceStr(),
	}
	return &Route{
		PathInfo:         pathInfo,
		Dest:             dest,
		path:             path,
		routeListIdx:     -1,
		action:           action,
		OutPathId:        outPathId,
		PolicyList:       make([]string, 0),
		PolicyHitCounter: 0,
	}
}

func (r *Route) setAction(action RouteAction) {
	r.action = action
}

func (r *Route) setIdx(idx int) {
	r.routeListIdx = idx
}

func (r *Route) SetBestPath() {
	r.PathInfo.BestPath = true
}

func (r *Route) ResetBestPath() {
	r.PathInfo.BestPath = false
}

func (r *Route) SetMultiPath() {
	r.PathInfo.MultiPath = true
}

func (r *Route) ResetMultiPath() {
	r.PathInfo.MultiPath = false
}

func (r *Route) SetAdditionalPath() {
	r.PathInfo.AdditionalPath = true
}

func (r *Route) ResetAdditionalPath() {
	r.PathInfo.AdditionalPath = false
}
