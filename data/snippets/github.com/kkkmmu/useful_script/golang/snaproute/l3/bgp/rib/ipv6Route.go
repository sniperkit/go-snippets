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
	bgputils "l3/bgp/utils"
	"models/objects"
	"strconv"
)

type IPv6Route struct {
	*bgpd.BGPv6RouteState
}

func NewIPv6Route(network string, cidrLen int16) *IPv6Route {
	return &IPv6Route{
		&bgpd.BGPv6RouteState{
			Network: network,
			CIDRLen: cidrLen,
		},
	}
}

func (i *IPv6Route) SetNetwork(network string) {
	i.Network = network
}

func (i *IPv6Route) GetNetwork() string {
	return i.Network
}

func (i *IPv6Route) SetCIDRLen(cidrLen int16) {
	i.CIDRLen = cidrLen
}

func (i *IPv6Route) GetCIDRLen() int16 {
	return i.CIDRLen
}

func (i *IPv6Route) GetPaths() []*bgpd.PathInfo {
	return i.Paths
}

func (i *IPv6Route) AppendPath(path *bgpd.PathInfo) {
	i.Paths = append(i.Paths, path)
}

func (i *IPv6Route) SetPath(path *bgpd.PathInfo, idx int) {
	i.Paths[idx] = path
}

func (i *IPv6Route) GetPath(idx int) *bgpd.PathInfo {
	return i.Paths[idx]
}

func (i *IPv6Route) GetLastPath() *bgpd.PathInfo {
	return i.Paths[len(i.Paths)-1]
}

func (i *IPv6Route) RemovePathAndSetLast(idx int) {
	if idx < len(i.Paths) {
		i.Paths[idx] = i.Paths[len(i.Paths)-1]
		i.Paths[len(i.Paths)-1] = nil
		i.Paths = i.Paths[:len(i.Paths)-1]
	}
}

func (i *IPv6Route) GetModelObject() objects.ConfigObj {
	var dbObj objects.BGPv6RouteState
	objects.ConvertThriftTobgpdBGPv6RouteStateObj(i.BGPv6RouteState, &dbObj)
	for idx1 := 0; idx1 < len(dbObj.Paths); idx1++ {
		for idx2 := 0; idx2 < len(dbObj.Paths[idx1].Path); idx2++ {
			asdoPlain, _ := strconv.Atoi(dbObj.Paths[idx1].Path[idx2])
			asdotPath, _ := bgputils.GetAsDot(asdoPlain)
			dbObj.Paths[idx1].Path[idx2] = asdotPath
		}
	}
	return &dbObj
}

func (i *IPv6Route) GetThriftObject() interface{} {
	return i.BGPv6RouteState
}
