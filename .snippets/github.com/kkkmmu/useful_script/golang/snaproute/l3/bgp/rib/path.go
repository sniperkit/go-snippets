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

// path.go
package rib

import (
	"encoding/binary"
	_ "fmt"
	"l3/bgp/baseobjects"
	"l3/bgp/packet"
	"net"
	_ "ribd"
	"strconv"
	"strings"
	"utils/logging"
)

const (
	RouteTypeAgg uint8 = 1 << iota
	RouteTypeConnected
	RouteTypeDefault
	RouteTypeStatic
	RouteTypeIGP
	RouteTypeEGP
	RouteTypeMax
)

const (
	RouteSrcLocal uint8 = 1 << iota
	RouteSrcExternal
	RouteSrcUnknown
)

var RouteTypeToSource = map[uint8]uint8{
	RouteTypeAgg:       RouteSrcLocal,
	RouteTypeConnected: RouteSrcLocal,
	RouteTypeDefault:   RouteSrcLocal,
	RouteTypeStatic:    RouteSrcLocal,
	RouteTypeIGP:       RouteSrcLocal,
	RouteTypeEGP:       RouteSrcExternal,
}

func getRouteSource(routeType uint8) uint8 {
	if routeSource, ok := RouteTypeToSource[routeType]; ok {
		return routeSource
	}

	return RouteSrcUnknown
}

type NHReachabilityInfo struct {
	nextHop          net.IP
	reachabilityInfo *ReachabilityInfo
}

type Path struct {
	rib                *LocRib
	logger             *logging.Writer
	NeighborConf       *base.NeighborConf
	PathAttrs          []packet.BGPPathAttr
	Pref               uint32
	nhReachabilityInfo map[uint32]*NHReachabilityInfo
	routeType          uint8
	MED                uint32
	LocalPref          uint32
	AggregatedPaths    map[string]*Path
}

func NewPath(locRib *LocRib, peer *base.NeighborConf, pa []packet.BGPPathAttr,
	mpReach *packet.BGPPathAttrMPReachNLRI, routeType uint8) *Path {
	path := &Path{
		rib:                locRib,
		logger:             locRib.logger,
		NeighborConf:       peer,
		PathAttrs:          pa,
		nhReachabilityInfo: make(map[uint32]*NHReachabilityInfo),
		routeType:          routeType,
		AggregatedPaths:    make(map[string]*Path),
	}

	path.logger.Info("Path:NewPath - path attr =", pa, "path.path attrs =", path.PathAttrs)
	path.Pref = path.calculatePref()
	path.constructNHReachabilityInfo(mpReach)
	return path
}

func (p *Path) Clone() *Path {
	path := &Path{
		rib:                p.rib,
		logger:             p.rib.logger,
		NeighborConf:       p.NeighborConf,
		PathAttrs:          p.PathAttrs,
		Pref:               p.Pref,
		nhReachabilityInfo: p.nhReachabilityInfo,
		routeType:          p.routeType,
		MED:                p.MED,
		LocalPref:          p.LocalPref,
	}

	return path
}

func (p *Path) calculatePref() uint32 {
	var pref uint32

	pref = BGP_INTERNAL_PREF

	for _, attr := range p.PathAttrs {
		if attr.GetCode() == packet.BGPPathAttrTypeLocalPref {
			p.LocalPref = attr.(*packet.BGPPathAttrLocalPref).Value
			pref = p.LocalPref
		} else if attr.GetCode() == packet.BGPPathAttrTypeMultiExitDisc {
			p.MED = attr.(*packet.BGPPathAttrMultiExitDisc).Value
		}
	}

	if p.IsExternal() {
		pref = BGP_EXTERNAL_PREF
	}

	return pref
}

func (p *Path) constructNHReachabilityInfo(mpReach *packet.BGPPathAttrMPReachNLRI) {
	for _, attr := range p.PathAttrs {
		if attr.GetCode() == packet.BGPPathAttrTypeNextHop {
			nextHop := attr.(*packet.BGPPathAttrNextHop)
			protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
			if _, ok := p.nhReachabilityInfo[protoFamily]; !ok {
				p.nhReachabilityInfo[protoFamily] = &NHReachabilityInfo{}
			}
			p.nhReachabilityInfo[protoFamily].nextHop = nextHop.Value
		}
	}

	if mpReach != nil {
		protoFamily := packet.GetProtocolFamily(mpReach.AFI, mpReach.SAFI)
		if _, ok := p.nhReachabilityInfo[protoFamily]; !ok {
			p.nhReachabilityInfo[protoFamily] = &NHReachabilityInfo{}
		}
		p.nhReachabilityInfo[protoFamily].nextHop = mpReach.NextHop.GetNextHop()
	}
}

func (p *Path) IsValid() bool {
	for _, attr := range p.PathAttrs {
		if attr.GetCode() == packet.BGPPathAttrTypeOriginatorId {
			if p.NeighborConf.Global.RouterId.Equal(attr.(*packet.BGPPathAttrOriginatorId).Value) {
				return false
			}
		}

		if attr.GetCode() == packet.BGPPathAttrTypeClusterList {
			clusters := attr.(*packet.BGPPathAttrClusterList).Value
			for _, clusterId := range clusters {
				if clusterId == p.NeighborConf.RunningConf.RouteReflectorClusterId {
					return false
				}
			}
		}
	}

	return true
}

func (p *Path) GetNeighborConf() *base.NeighborConf {
	return p.NeighborConf
}

func (p *Path) GetPeerIP() string {
	if p.NeighborConf != nil {
		return p.NeighborConf.Neighbor.NeighborAddress.String()
	}
	return ""
}

func (p *Path) GetPreference() uint32 {
	return p.Pref
}

func (p *Path) GetAS4ByteList() []string {
	asList := make([]string, 0)
	for _, attr := range p.PathAttrs {
		if attr.GetCode() == packet.BGPPathAttrTypeASPath {
			asPaths := attr.(*packet.BGPPathAttrASPath).Value
			asSize := attr.(*packet.BGPPathAttrASPath).ASSize
			for _, asSegment := range asPaths {
				if asSize == 4 {
					seg := asSegment.(*packet.BGPAS4PathSegment)
					if seg.Type == packet.BGPASPathSegmentSet {
						asSetList := make([]string, 0, len(seg.AS))
						for _, as := range seg.AS {
							asSetList = append(asSetList, strconv.Itoa(int(as)))
						}
						asSetStr := strings.Join(asSetList, ", ")
						asSetStr = "{ " + asSetStr + " }"
						asList = append(asList, asSetStr)
					} else if seg.Type == packet.BGPASPathSegmentSequence {
						for _, as := range seg.AS {
							asList = append(asList, strconv.Itoa(int(as)))
						}
					}
				} else {
					seg := asSegment.(*packet.BGPAS2PathSegment)
					if seg.Type == packet.BGPASPathSegmentSet {
						asSetList := make([]string, 0, len(seg.AS))
						for _, as := range seg.AS {
							asSetList = append(asSetList, strconv.Itoa(int(as)))
						}
						asSetStr := strings.Join(asSetList, ", ")
						asSetStr = "{ " + asSetStr + " }"
						asList = append(asList, asSetStr)
					} else if seg.Type == packet.BGPASPathSegmentSequence {
						for _, as := range seg.AS {
							asList = append(asList, strconv.Itoa(int(as)))
						}
					}
				}
			}
			break
		}
	}

	return asList
}

func (p *Path) HasASLoop() bool {
	if p.NeighborConf == nil {
		return false
	}
	return packet.HasASLoop(p.PathAttrs, p.NeighborConf.RunningConf.LocalAS)
}

func (p *Path) IsLocal() bool {
	return getRouteSource(p.routeType) == RouteSrcLocal
}

func (p *Path) IsAggregate() bool {
	return p.routeType == RouteTypeAgg
}

func (p *Path) IsExternal() bool {
	return p.NeighborConf != nil && p.NeighborConf.IsExternal()
}

func (p *Path) IsInternal() bool {
	return p.NeighborConf != nil && p.NeighborConf.IsInternal()
}

func (p *Path) GetSourceStr() string {
	return ""
}

func (p *Path) GetNumASes() uint32 {
	p.logger.Info("Path:GetNumASes - path attrs =", p.PathAttrs)
	return packet.GetNumASes(p.PathAttrs)
}

func (p *Path) GetOrigin() uint8 {
	return packet.GetOrigin(p.PathAttrs)
}

func (p *Path) GetNextHop(protoFamily uint32) net.IP {
	if nhReachInfo, ok := p.nhReachabilityInfo[protoFamily]; ok {
		return nhReachInfo.nextHop
	}
	return nil
}

func (p *Path) GetBGPId() uint32 {
	for _, attr := range p.PathAttrs {
		if attr.GetCode() == packet.BGPPathAttrTypeOriginatorId {
			return binary.BigEndian.Uint32(attr.(*packet.BGPPathAttrOriginatorId).Value.To4())
		}
	}

	return binary.BigEndian.Uint32(p.NeighborConf.BGPId.To4())
}

func (p *Path) GetNumClusters() uint16 {
	return packet.GetNumClusters(p.PathAttrs)
}

func (p *Path) SetReachabilityForFamily(protoFamily uint32, reachabilityInfo *ReachabilityInfo) {
	if nhReachInfo, ok := p.nhReachabilityInfo[protoFamily]; ok {
		nhReachInfo.reachabilityInfo = reachabilityInfo
	}
}

func (p *Path) SetReachabilityForNextHop(nextHop string, reachabilityInfo *ReachabilityInfo) {
	for _, nhReachInfo := range p.nhReachabilityInfo {
		if nhReachInfo.nextHop.String() == nextHop {
			nhReachInfo.reachabilityInfo = reachabilityInfo
		}
	}
}

func (p *Path) GetReachability(protoFamily uint32) *ReachabilityInfo {
	if nhReachInfo, ok := p.nhReachabilityInfo[protoFamily]; ok {
		return nhReachInfo.reachabilityInfo
	}
	return nil
}

func (p *Path) IsReachable(protoFamily uint32) bool {
	if p.routeType == RouteTypeStatic || p.routeType == RouteTypeConnected || p.routeType == RouteTypeDefault ||
		p.routeType == RouteTypeIGP {
		return true
	}

	if nhReachInfo, ok := p.nhReachabilityInfo[protoFamily]; ok {
		if nhReachInfo.reachabilityInfo != nil {
			return true
		}
	}
	return false
}

func (p *Path) setAggregatedPath(destIP string, path *Path) {
	if _, ok := p.AggregatedPaths[destIP]; ok {
		p.logger.Errf("Path from %s is already added to the aggregated paths %v", destIP, p.AggregatedPaths)
	}
	p.AggregatedPaths[destIP] = path
}

func (p *Path) checkMEDForAggregation(path *Path) (uint32, uint32, bool) {
	aggMED, aggOK := packet.GetMED(p.PathAttrs)
	med, ok := packet.GetMED(path.PathAttrs)
	if aggOK == ok && aggMED == med {
		return aggMED, med, true
	}

	return aggMED, med, false
}

func (p *Path) addPathToAggregate(destIP string, path *Path, generateASSet bool) bool {
	aggMED, med, isMEDEqual := p.checkMEDForAggregation(path)

	if _, ok := p.AggregatedPaths[destIP]; ok {
		if !isMEDEqual {
			p.logger.Info("addPathToAggregate: MED", med, "in the new path", path,
				"is not the same as the MED", aggMED, "in the agg path, remove the old path...")
			delete(p.AggregatedPaths, destIP)
			p.removePathFromAggregate(destIP, generateASSet)
		} else {
			p.logger.Infof("addPathToAggregatePath from %s is already aggregated, replace it...", destIP)
			p.AggregatedPaths[destIP] = path
			p.aggregateAllPaths(generateASSet)
		}

		return true
	}

	if !isMEDEqual {
		p.logger.Info("addPathToAggregate: Can't aggregate new path MEDs not equal, new path MED =", med,
			"Agg path MED =", aggMED)
		return false
	}

	origin := packet.BGPPathAttrOriginType(packet.GetOrigin(path.PathAttrs))
	atomicAggregate := packet.GetAtomicAggregatePathAttr(path.PathAttrs)

	idx := 0
	foundAtomicAgg := false
	p.logger.Debugf("addPathToAggregatePath - len=%d, p.PathAttrs=%+v", len(p.PathAttrs), p.PathAttrs)
	p.logger.Debugf("addPathToAggregatePath - len=%d, path.PathAttrs=%+v", len(path.PathAttrs), path.PathAttrs)
	for idx = 0; idx < len(p.PathAttrs); idx++ {
		if p.PathAttrs[idx].GetCode() == packet.BGPPathAttrTypeAtomicAggregate {
			foundAtomicAgg = true
			continue
		}

		if p.PathAttrs[idx].GetCode() == packet.BGPPathAttrTypeOrigin && origin != packet.BGPPathAttrOriginMax {
			if origin > p.PathAttrs[idx].(*packet.BGPPathAttrOrigin).Value {
				p.PathAttrs[idx].(*packet.BGPPathAttrOrigin).Value = origin
			}
		}
	}

	if !foundAtomicAgg && atomicAggregate != nil {
		atomicAggregate := packet.NewBGPPathAttrAtomicAggregate()
		p.PathAttrs = packet.AddPathAttrToPathAttrsByCode(p.PathAttrs, packet.BGPPathAttrTypeAtomicAggregate, atomicAggregate)
	}

	p.AggregatedPaths[destIP] = path
	return true
}

func (p *Path) removePathFromAggregate(destIP string, generateASSet bool) {
	delete(p.AggregatedPaths, destIP)
	p.aggregateAllPaths(generateASSet)
}

func (p *Path) IsAggregatePath() bool {
	return (len(p.AggregatedPaths) > 0)
}

func (p *Path) isAggregatePathEmpty() bool {
	return (len(p.AggregatedPaths) == 0)
}

func (p *Path) aggregateAllPaths(generateASSet bool) {
	var origin, atomicAggregate packet.BGPPathAttr
	asPathList := make([]*packet.BGPPathAttrASPath, 0, len(p.AggregatedPaths))
	var aggASPath *packet.BGPPathAttrASPath
	for _, individualPath := range p.AggregatedPaths {
		for _, pathAttr := range individualPath.PathAttrs {
			if pathAttr.GetCode() == packet.BGPPathAttrTypeOrigin {
				if origin == nil || pathAttr.(*packet.BGPPathAttrOrigin).Value > origin.(*packet.BGPPathAttrOrigin).Value {
					origin = pathAttr
				}
			}

			if pathAttr.GetCode() == packet.BGPPathAttrTypeAtomicAggregate {
				if atomicAggregate == nil {
					atomicAggregate = pathAttr
				}
			}

			if pathAttr.GetCode() == packet.BGPPathAttrTypeASPath {
				asPathList = append(asPathList, pathAttr.(*packet.BGPPathAttrASPath))
			}
		}
	}

	if generateASSet {
		aggASPath = packet.AggregateASPaths(asPathList)
	}

	for idx, pathAttr := range p.PathAttrs {
		if pathAttr.GetCode() == packet.BGPPathAttrTypeOrigin && origin != nil {
			p.PathAttrs[idx] = origin
		}

		if pathAttr.GetCode() == packet.BGPPathAttrTypeASPath && aggASPath != nil {
			p.PathAttrs[idx] = aggASPath
		}

		if pathAttr.GetCode() == packet.BGPPathAttrTypeAtomicAggregate && atomicAggregate != nil {
			p.PathAttrs[idx] = atomicAggregate
		}
	}
}
