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

// bgp.go
package packet

import (
	"errors"
	_ "fmt"
	"l3/bgp/utils"
	"math"
	"net"
	"sort"
	"strings"
)

func PrependAS(updateMsg *BGPMessage, AS uint32, asSize uint8) {
	body := updateMsg.Body.(*BGPUpdate)

	for _, pa := range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeASPath {
			asPathSegments := pa.(*BGPPathAttrASPath).Value
			var newASPathSegment BGPASPathSegment
			if len(asPathSegments) == 0 || asPathSegments[0].GetType() == BGPASPathSegmentSet || asPathSegments[0].GetLen() >= 255 {
				if asSize == 4 {
					newASPathSegment = NewBGPAS4PathSegmentSeq()
				} else {
					newASPathSegment = NewBGPAS2PathSegmentSeq()
					if asSize == 2 {
						if AS > math.MaxUint16 {
							AS = uint32(BGPASTrans)
						}
					}
				}
				pa.(*BGPPathAttrASPath).PrependASPathSegment(newASPathSegment)
			}
			asPathSegments = pa.(*BGPPathAttrASPath).Value
			asPathSegments[0].PrependAS(AS)
			pa.(*BGPPathAttrASPath).BGPPathAttrBase.Length += uint16(asSize)
		} else if pa.GetCode() == BGPPathAttrTypeAS4Path {
			asPathSegments := pa.(*BGPPathAttrAS4Path).Value
			var newAS4PathSegment *BGPAS4PathSegment
			if len(asPathSegments) == 0 || asPathSegments[0].GetType() == BGPASPathSegmentSet || asPathSegments[0].GetLen() >= 255 {
				newAS4PathSegment = NewBGPAS4PathSegmentSeq()
				pa.(*BGPPathAttrAS4Path).AddASPathSegment(newAS4PathSegment)
			}
			asPathSegments = pa.(*BGPPathAttrAS4Path).Value
			asPathSegments[0].PrependAS(AS)
			pa.(*BGPPathAttrASPath).BGPPathAttrBase.Length += uint16(asSize)
		}
	}
}

func PrependASList(pathAttrs []BGPPathAttr, asList []uint32, asSize uint8) []BGPPathAttr {
	newASes := 0
	for _, pa := range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeASPath {
			asPathSegments := pa.(*BGPPathAttrASPath).Value
			var newASPathSegment BGPASPathSegment
			if len(asPathSegments) == 0 || asPathSegments[0].GetType() == BGPASPathSegmentSet || asPathSegments[0].GetLen() >= 255 {
				if asSize == 4 {
					newASPathSegment = NewBGPAS4PathSegmentSeq()
				} else {
					newASPathSegment = NewBGPAS2PathSegmentSeq()
					if asSize == 2 {
						for idx, as := range asList {
							if as > math.MaxUint16 {
								asList[idx] = uint32(BGPASTrans)
							}
						}
					}
				}
				pa.(*BGPPathAttrASPath).PrependASPathSegment(newASPathSegment)
			}
			asPathSegments = pa.(*BGPPathAttrASPath).Value
			numASes := int(asPathSegments[0].GetLen())
			newASes = len(asList)
			if newASes > (asPathSegments[0].GetMaxASes() - numASes) {
				newASes = (asPathSegments[0].GetMaxASes() - numASes)
			}
			asPathSegments[0].PrependASList(asList[:newASes])
			pa.(*BGPPathAttrASPath).BGPPathAttrBase.Length += uint16(int(asSize) * newASes)
		} else if pa.GetCode() == BGPPathAttrTypeAS4Path {
			asPathSegments := pa.(*BGPPathAttrAS4Path).Value
			var newAS4PathSegment *BGPAS4PathSegment
			if len(asPathSegments) == 0 || asPathSegments[0].GetType() == BGPASPathSegmentSet || asPathSegments[0].GetLen() >= 255 {
				newAS4PathSegment = NewBGPAS4PathSegmentSeq()
				pa.(*BGPPathAttrAS4Path).AddASPathSegment(newAS4PathSegment)
			}
			asPathSegments = pa.(*BGPPathAttrAS4Path).Value
			numASes := int(asPathSegments[0].GetLen())
			newASes = len(asList)
			if newASes > (asPathSegments[0].GetMaxASes() - numASes) {
				newASes = (asPathSegments[0].GetMaxASes() - numASes)
			}
			asPathSegments[0].PrependASList(asList[:newASes])
			pa.(*BGPPathAttrASPath).BGPPathAttrBase.Length += uint16(int(asSize) * newASes)
		}
	}

	if newASes > 0 {
		pathAttrs = PrependASList(pathAttrs, asList[newASes:], asSize)
	}

	return pathAttrs
}

func AppendASToAS4PathSeg(asPath *BGPPathAttrASPath, pathSeg BGPASPathSegment, asPathType BGPASPathSegmentType,
	asNum uint32) BGPASPathSegment {
	if pathSeg == nil {
		pathSeg = NewBGPAS4PathSegment(asPathType)
	} else if pathSeg.GetType() != asPathType {
		asPath.AppendASPathSegment(pathSeg)
	}

	if !pathSeg.AppendAS(asNum) {
		asPath.AppendASPathSegment(pathSeg)
		pathSeg = NewBGPAS4PathSegment(asPathType)
		pathSeg.AppendAS(asNum)
	}

	return pathSeg
}

func GetASPathAsString(pathAttrs []BGPPathAttr) (asPathStr string, err error) {
	err = errors.New("AS path not found")
	paASPath := getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeASPath)
	if paASPath != nil {
		asPath := paASPath.(*BGPPathAttrASPath)
		for _, asSeg := range asPath.Value {
			asSegStr := asSeg.String()
			utils.Logger.Infof("GetASPathAsString - asSegStr {%s}", asSegStr)
			asSegStr = strings.TrimLeft(asSegStr, "[ ")
			asSegStr = strings.TrimRight(asSegStr, "] ")
			utils.Logger.Infof("GetASPathAsString - asSegStr after trim {%s}", asSegStr)
			asPathStr = asPathStr + asSegStr + " "
			utils.Logger.Infof("GetASPathAsString - asPathStr {%s}", asPathStr)
		}
		asPathStr = strings.TrimSpace(asPathStr)
		err = nil
	} else {
		paASPath = getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeAS4Path)
		if paASPath != nil {
			asPath := paASPath.(*BGPPathAttrAS4Path)
			for _, asSeg := range asPath.Value {
				asSegStr := asSeg.String()
				utils.Logger.Infof("GetASPathAsString - asSegStr {%s}", asSegStr)
				asSegStr = strings.TrimLeft(asSegStr, "[ ")
				asSegStr = strings.TrimRight(asSegStr, "] ")
				utils.Logger.Infof("GetASPathAsString - asSegStr after trim {%s}", asSegStr)
				asPathStr = asPathStr + asSegStr + " "
				utils.Logger.Infof("GetASPathAsString - asPathStr {%s}", asPathStr)
			}
			asPathStr = strings.TrimSpace(asPathStr)
			err = nil
		}
	}
	utils.Logger.Infof("GetASPathAsString - final asPathStr {%s}", asPathStr)
	return asPathStr, err
}

func AddPathAttrToPathAttrsByCode(pathAttrs []BGPPathAttr, code BGPPathAttrType, attr BGPPathAttr) []BGPPathAttr {
	addIdx := -1
	for idx, pa := range pathAttrs {
		if pa.GetCode() > code {
			addIdx = idx
		}
	}

	if addIdx == -1 {
		addIdx = len(pathAttrs)
	}

	pathAttrs = append(pathAttrs, attr)
	copy(pathAttrs[addIdx+1:], pathAttrs[addIdx:])
	pathAttrs[addIdx] = attr
	return pathAttrs
}

func addPathAttrToPathAttrs(pathAttrs []BGPPathAttr, attr BGPPathAttr) []BGPPathAttr {
	if attr != nil {
		return AddPathAttrToPathAttrsByCode(pathAttrs, attr.GetCode(), attr)
	}

	return pathAttrs
}

func AddMPReachNLRIToPathAttrs(pathAttrs []BGPPathAttr, mpReach *BGPPathAttrMPReachNLRI) []BGPPathAttr {
	return AddPathAttrToPathAttrsByCode(pathAttrs, BGPPathAttrTypeMPReachNLRI, mpReach)
}

func AddMPUnreachNLRIToPathAttrs(pathAttrs []BGPPathAttr, mpUnreach *BGPPathAttrMPUnreachNLRI) []BGPPathAttr {
	return AddPathAttrToPathAttrsByCode(pathAttrs, BGPPathAttrTypeMPUnreachNLRI, mpUnreach)
}

func addPathAttr(updateMsg *BGPMessage, code BGPPathAttrType, attr BGPPathAttr) {
	body := updateMsg.Body.(*BGPUpdate)
	body.PathAttributes = AddPathAttrToPathAttrsByCode(body.PathAttributes, code, attr)
	return
}

func removePathAttr(updateMsg *BGPMessage, code BGPPathAttrType) BGPPathAttr {
	body := updateMsg.Body.(*BGPUpdate)
	return removeTypeFromPathAttrs(&body.PathAttributes, code)
}

func removeTypeFromPathAttrs(pathAttrs *[]BGPPathAttr, code BGPPathAttrType) BGPPathAttr {
	for idx, pa := range *pathAttrs {
		if pa.GetCode() == code {
			(*pathAttrs) = append((*pathAttrs)[:idx], (*pathAttrs)[idx+1:]...)
			return pa
		}
	}
	return nil
}

func RemoveNextHop(pathAttrs *[]BGPPathAttr) {
	removeTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeNextHop)
}

func RemoveMultiExitDisc(updateMsg *BGPMessage) BGPPathAttr {
	return removePathAttr(updateMsg, BGPPathAttrTypeMultiExitDisc)
}

func RemoveLocalPref(updateMsg *BGPMessage) BGPPathAttr {
	return removePathAttr(updateMsg, BGPPathAttrTypeLocalPref)
}

func RemoveMPAttrs(pathAttrs *[]BGPPathAttr) (mpReach *BGPPathAttrMPReachNLRI, mpUnreach *BGPPathAttrMPUnreachNLRI) {
	reach := removeTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeMPReachNLRI)
	if reach != nil {
		mpReach = reach.(*BGPPathAttrMPReachNLRI)
	}

	unreach := removeTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeMPUnreachNLRI)
	if unreach != nil {
		mpUnreach = unreach.(*BGPPathAttrMPUnreachNLRI)
	}
	return mpReach, mpUnreach
}

func getPathAttr(updateMsg *BGPMessage, code BGPPathAttrType) BGPPathAttr {
	body := updateMsg.Body.(*BGPUpdate)
	return getTypeFromPathAttrs(body.PathAttributes, code)
}

func getTypeFromPathAttrs(pathAttrs []BGPPathAttr, code BGPPathAttrType) BGPPathAttr {
	for _, pa := range pathAttrs {
		if pa.GetCode() == code {
			return pa
		}
	}
	return nil
}

func GetAtomicAggregatePathAttr(pathAttrs []BGPPathAttr) BGPPathAttr {
	return getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeAtomicAggregate)
}

func GetMPAttrs(pathAttrs []BGPPathAttr) (mpReach *BGPPathAttrMPReachNLRI, mpUnreach *BGPPathAttrMPUnreachNLRI) {
	reach := getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeMPReachNLRI)
	if reach != nil {
		mpReach = reach.(*BGPPathAttrMPReachNLRI)
	}

	unreach := getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeMPUnreachNLRI)
	if unreach != nil {
		mpUnreach = unreach.(*BGPPathAttrMPUnreachNLRI)
	}
	return mpReach, mpUnreach
}

func SetLocalPrefToPathAttrs(pathAttrs []BGPPathAttr, pref uint32, overwrite bool) []BGPPathAttr {
	var idx int
	var pa BGPPathAttr
	for idx, pa = range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeLocalPref {
			if overwrite {
				pathAttrs[idx].(*BGPPathAttrLocalPref).Value = pref
			}
			return pathAttrs
		}
	}

	idx = -1
	for idx, pa = range pathAttrs {
		if pa.GetCode() > BGPPathAttrTypeLocalPref {
			break
		} else if idx == len(pathAttrs)-1 {
			idx += 1
		}
	}

	if idx >= 0 {
		paLocalPref := NewBGPPathAttrLocalPref()
		paLocalPref.Value = pref
		pathAttrs = append(pathAttrs, paLocalPref)
		if idx < len(pathAttrs)-1 {
			copy(pathAttrs[idx+1:], pathAttrs[idx:])
			pathAttrs[idx] = paLocalPref
		}
	}

	return pathAttrs
}

func SetLocalPref(updateMsg *BGPMessage, pref uint32, overwrite bool) {
	body := updateMsg.Body.(*BGPUpdate)
	body.PathAttributes = SetLocalPrefToPathAttrs(body.PathAttributes, pref, overwrite)
}

func SetMEDToPathAttrs(pathAttrs []BGPPathAttr, med uint32, overwrite bool) []BGPPathAttr {
	var idx int
	var pa BGPPathAttr
	for idx, pa = range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeMultiExitDisc {
			if overwrite {
				pathAttrs[idx].(*BGPPathAttrMultiExitDisc).Value = med
			}
			return pathAttrs
		}
	}

	idx = -1
	for idx, pa = range pathAttrs {
		if pa.GetCode() > BGPPathAttrTypeMultiExitDisc {
			break
		} else if idx == len(pathAttrs)-1 {
			idx += 1
		}
	}

	if idx >= 0 {
		paMED := NewBGPPathAttrMultiExitDisc()
		paMED.Value = med
		pathAttrs = append(pathAttrs, paMED)
		if idx < len(pathAttrs)-1 {
			copy(pathAttrs[idx+1:], pathAttrs[idx:])
			pathAttrs[idx] = paMED
		}
	}

	return pathAttrs
}

func SetMED(updateMsg *BGPMessage, med uint32, overwrite bool) {
	body := updateMsg.Body.(*BGPUpdate)
	body.PathAttributes = SetMEDToPathAttrs(body.PathAttributes, med, overwrite)
}

func SetNextHop(updateMsg *BGPMessage, nextHop net.IP) {
	body := updateMsg.Body.(*BGPUpdate)
	SetNextHopPathAttrs(body.PathAttributes, nextHop)
}

func SetNextHopPathAttrs(pathAttrs []BGPPathAttr, nextHopIP net.IP) {
	for idx, pa := range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeNextHop {
			pathAttrs[idx].(*BGPPathAttrNextHop).Value = nextHopIP
		}
	}
}

func SetPathAttrAggregator(pathAttrs []BGPPathAttr, as uint32, ip net.IP) {
	for idx, pa := range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeAggregator {
			aggAS := NewBGPAggregator4ByteAS()
			aggAS.AS = as
			pathAttrs[idx].(*BGPPathAttrAggregator).SetBGPAggregatorAS(aggAS)
			pathAttrs[idx].(*BGPPathAttrAggregator).IP = ip
		}
	}
}

func HasMPAttrs(pathAttrs []BGPPathAttr) bool {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeMPReachNLRI || attr.GetCode() == BGPPathAttrTypeMPUnreachNLRI {
			return true
		}
	}
	return false
}

func HasMPReachNLRI(pathAttrs []BGPPathAttr) bool {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeMPReachNLRI {
			return true
		}
	}
	return false
}

func HasASLoop(pathAttrs []BGPPathAttr, localAS uint32) bool {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeASPath {
			asPaths := attr.(*BGPPathAttrASPath).Value
			asSize := attr.(*BGPPathAttrASPath).ASSize
			for _, asSegment := range asPaths {
				if asSize == 4 {
					seg := asSegment.(*BGPAS4PathSegment)
					for _, as := range seg.AS {
						if as == localAS {
							return true
						}
					}
				} else {
					seg := asSegment.(*BGPAS2PathSegment)
					for _, as := range seg.AS {
						if as == uint16(localAS) {
							return true
						}
					}
				}
			}
			break
		}
	}

	return false
}

func GetNumASes(pathAttrs []BGPPathAttr) uint32 {
	var total uint32 = 0
	utils.Logger.Info("helpers:GetNumASes - path attrs =", pathAttrs)
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeASPath {
			asPaths := attr.(*BGPPathAttrASPath).Value
			for _, asPath := range asPaths {
				total += uint32(asPath.GetNumASes())
			}
			break
		}
	}

	return total
}

func GetOrigin(pathAttrs []BGPPathAttr) uint8 {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeOrigin {
			return uint8(attr.(*BGPPathAttrOrigin).Value)
		}
	}

	return uint8(BGPPathAttrOriginMax)
}

func GetOriginTypeStr(origin uint8) string {
	if val, ok := BGPPathAttrOriginToStrMap[BGPPathAttrOriginType(origin)]; ok {
		return val
	}

	return BGPPathAttrOriginToStrMap[BGPPathAttrOriginMax]
}

func GetMED(pathAttrs []BGPPathAttr) (uint32, bool) {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeMultiExitDisc {
			return attr.(*BGPPathAttrMultiExitDisc).Value, true
		}
	}

	return uint32(0), false
}

func GetLocalPref(pathAttrs []BGPPathAttr) (uint32, error) {
	pa := getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeLocalPref)
	if pa != nil {
		localPrefPathAttr := pa.(*BGPPathAttrLocalPref)
		return localPrefPathAttr.Value, nil
	}

	return 0, errors.New("Local Pref path attr not found")
}

func GetNextHop(pathAttrs []BGPPathAttr) net.IP {
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeNextHop {
			return attr.(*BGPPathAttrNextHop).Value
		}
	}

	return net.IPv4zero
}

func GetNumClusters(pathAttrs []BGPPathAttr) uint16 {
	var total uint16 = 0
	for _, attr := range pathAttrs {
		if attr.GetCode() == BGPPathAttrTypeClusterList {
			length := attr.(*BGPPathAttrClusterList).Length
			total = length / 4
			break
		}
	}

	return total
}

var AggRoutesDefaultBGPPathAttr = map[BGPPathAttrType]BGPPathAttr{
	BGPPathAttrTypeOrigin:     NewBGPPathAttrOrigin(BGPPathAttrOriginIncomplete),
	BGPPathAttrTypeASPath:     NewBGPPathAttrASPath(),
	BGPPathAttrTypeNextHop:    NewBGPPathAttrNextHop(),
	BGPPathAttrTypeAggregator: NewBGPPathAttrAggregator(),
	//BGPPathAttrTypeAtomicAggregate: NewBGPPathAttrAtomicAggregate(),
}

func AggregateASPaths(asPathList []*BGPPathAttrASPath) *BGPPathAttrASPath {
	aggASPath := NewBGPPathAttrASPath()
	if len(asPathList) > 0 {
		asNumMap := make(map[uint32]bool, 10)
		asPathIterList := make([]*ASPathIter, 0, len(asPathList))
		for i := 0; i < len(asPathList); i++ {
			asPathIterList = append(asPathIterList, NewASPathIter(asPathList[i]))
		}

		var asPathSeg BGPASPathSegment
		asSeqDone := false
		var asPathVal, iterASPathVal uint32
		var asPathType, iterASPathType BGPASPathSegmentType
		var flag, iterFlag bool
		var idx int
		for {
			idx = 0
			asPathVal, asPathType, flag = asPathIterList[idx].Next()
			if !flag {
				break
			}

			for idx = 1; idx < len(asPathIterList); idx++ {
				iterASPathVal, iterASPathType, iterFlag = asPathIterList[idx].Next()
				if !iterFlag || iterASPathType != asPathType || iterASPathVal != asPathVal {
					asSeqDone = true
					break
				}
			}

			if asSeqDone {
				break
			}

			asPathSeg = AppendASToAS4PathSeg(aggASPath, asPathSeg, asPathType, asPathVal)
			asNumMap[asPathVal] = true
		}
		if asPathSeg != nil && asPathSeg.GetNumASes() > 0 {
			aggASPath.AppendASPathSegment(asPathSeg)
		}
		if !flag || !iterFlag {
			asPathIterList[idx] = nil
		}
		asPathSeg = NewBGPAS4PathSegmentSet()

		if flag {
			if !asNumMap[asPathVal] {
				asPathSeg = AppendASToAS4PathSeg(aggASPath, asPathSeg, asPathType, asPathVal)
				asNumMap[asPathVal] = true
			}
			if iterFlag {
				if !asNumMap[iterASPathVal] {
					asPathSeg = AppendASToAS4PathSeg(aggASPath, asPathSeg, iterASPathType, iterASPathVal)
					asNumMap[iterASPathVal] = true
				}
			}
			for idx = idx + 1; idx < len(asPathIterList); idx++ {
				asPathVal, asPathType, flag = asPathIterList[idx].Next()
				if flag {
					if !asNumMap[asPathVal] {
						asPathSeg = AppendASToAS4PathSeg(aggASPath, asPathSeg, asPathType, asPathVal)
						asNumMap[asPathVal] = true
					}
				} else {
					asPathIterList[idx] = nil
				}
			}
		}
		asPathIterList = RemoveNilItemsFromList(asPathIterList)
		for idx = 0; idx < len(asPathIterList); idx++ {
			for asPathVal, asPathType, flag = asPathIterList[idx].Next(); flag; {
				if !asNumMap[asPathVal] {
					asPathSeg = AppendASToAS4PathSeg(aggASPath, asPathSeg, asPathType, asPathVal)
					asNumMap[asPathVal] = true
				} else {
					asPathIterList[idx] = nil
					break
				}
			}
		}
	}
	return aggASPath
}

func ConstructPathAttrForAggRoutes(pathAttrs []BGPPathAttr, generateASSet bool) []BGPPathAttr {
	newPathAttrs := make([]BGPPathAttr, 0)
	reqAttrs := []BGPPathAttrType{BGPPathAttrTypeOrigin, BGPPathAttrTypeASPath, BGPPathAttrTypeNextHop,
		BGPPathAttrTypeAtomicAggregate, BGPPathAttrTypeAggregator}

	for _, pa := range pathAttrs {
		if pa.GetCode() == BGPPathAttrTypeNextHop || pa.GetCode() == BGPPathAttrTypeOrigin ||
			pa.GetCode() == BGPPathAttrTypeASPath || pa.GetCode() == BGPPathAttrTypeAtomicAggregate ||
			pa.GetCode() == BGPPathAttrTypeAggregator || pa.GetCode() == BGPPathAttrTypeMultiExitDisc {
			if pa.GetCode() == BGPPathAttrTypeASPath && !generateASSet {
				asPath := NewBGPPathAttrASPath()
				newPathAttrs = append(newPathAttrs, asPath)
			} else {
				newPathAttrs = append(newPathAttrs, pa.Clone())
			}
		}
	}

	sort.Sort(PathAttrs(newPathAttrs))

	total := len(newPathAttrs)
	idx := 0
	for _, pa := range newPathAttrs {
		for i, paType := range reqAttrs[idx:total] {
			if paType < pa.GetCode() && (paType != BGPPathAttrTypeASPath || generateASSet) {
				pathAttr := AggRoutesDefaultBGPPathAttr[paType]
				newPathAttrs = append(newPathAttrs, pathAttr)
			} else {
				if paType == pa.GetCode() {
					idx += (i + 1)
				} else {
					idx += i
				}
				break
			}
		}
	}
	return newPathAttrs
}

func ConstructMPReachNLRIForAggRoutes(protoFamily uint32) *BGPPathAttrMPReachNLRI {
	nh := NewMPNextHopIP()
	nh.SetNextHop(net.IPv6zero)
	afi, safi := GetAfiSafi(protoFamily)
	pa := NewBGPPathAttrMPReachNLRI()
	pa.AFI = afi
	pa.SAFI = safi
	pa.SetNextHop(nh)
	return pa
}

func ConstructPathAttrForLocalRoutes(originType BGPPathAttrOriginType, as, med uint32) []BGPPathAttr {
	pathAttrs := make([]BGPPathAttr, 0)

	origin := NewBGPPathAttrOrigin(originType)
	pathAttrs = append(pathAttrs, origin)

	asPath := NewBGPPathAttrASPath()
	pathAttrs = append(pathAttrs, asPath)

	nextHop := NewBGPPathAttrNextHop()
	nextHop.Value = net.IPv4zero
	pathAttrs = append(pathAttrs, nextHop)

	paMED := NewBGPPathAttrMultiExitDisc()
	paMED.Value = med
	pathAttrs = append(pathAttrs, paMED)

	return pathAttrs
}

func ConstructPathAttrForConnRoutes(as, med uint32) []BGPPathAttr {
	return ConstructPathAttrForLocalRoutes(BGPPathAttrOriginIncomplete, as, med)
}

func ConstructPathAttrForDefaultRoute(as, med uint32) []BGPPathAttr {
	return ConstructPathAttrForLocalRoutes(BGPPathAttrOriginEGP, as, med)
}

func ClonePathAttrs(pathAttrs []BGPPathAttr) []BGPPathAttr {
	newPathAttrs := make([]BGPPathAttr, len(pathAttrs))
	for idx := range pathAttrs {
		newPathAttrs[idx] = pathAttrs[idx].Clone()
	}
	return newPathAttrs
}

func CopyPathAttrs(pathAttrs []BGPPathAttr) []BGPPathAttr {
	newPathAttrs := make([]BGPPathAttr, len(pathAttrs))
	copy(newPathAttrs, pathAttrs)
	return newPathAttrs
}

func ConstructNLRIFromPathIdAndNLRI(nlri NLRI, pathId uint32) NLRI {
	newNLRI := nlri.Clone()
	if extNLRI, ok := newNLRI.(*ExtNLRI); ok {
		extNLRI.PathId = pathId
	}

	return newNLRI
}

func ConstructIPPrefix(ipStr string, maskStr string) *IPPrefix {
	ip := net.ParseIP(ipStr)
	var mask net.IPMask
	if ip.To4() != nil {
		utils.Logger.Infof("ConstructIPPrefix IPv4 - mask ip %+v mask ip mask %+v", net.ParseIP(maskStr),
			net.IPMask(net.ParseIP(maskStr).To4()))
		mask = net.IPMask(net.ParseIP(maskStr).To4())
	} else {
		utils.Logger.Infof("ConstructIPPrefix IPv6 - mask ip %+v mask ip mask %+v", net.ParseIP(maskStr),
			net.IPMask(net.ParseIP(maskStr).To16()))
		mask = net.IPMask(net.ParseIP(maskStr).To16())
	}
	ones, _ := mask.Size()
	return NewIPPrefix(ip.Mask(mask), uint8(ones))
}

func ConstructIPPrefixFromCIDR(cidr string) (*IPPrefix, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		utils.Logger.Info("ConstructIPPrefixFromCIDR: ParseCIDR for IPPrefix", cidr, "failed with err", err)
		return nil, err
	}

	ones, _ := ipNet.Mask.Size()
	return NewIPPrefix(ipNet.IP, uint8(ones)), nil
}

func ConstructMPUnreachNLRI(afi AFI, safi SAFI, nlriList []NLRI) *BGPPathAttrMPUnreachNLRI {
	mpUnreachNLRI := NewBGPPathAttrMPUnreachNLRI()
	mpUnreachNLRI.AFI = afi
	mpUnreachNLRI.SAFI = safi
	mpUnreachNLRI.AddNLRIList(nlriList)
	return mpUnreachNLRI
}

func ConstructMPUnreachNLRIFromProtoFamily(protoFamily uint32, nlriList []NLRI) *BGPPathAttrMPUnreachNLRI {
	afi, safi := GetAfiSafi(protoFamily)
	return ConstructMPUnreachNLRI(afi, safi, nlriList)
}

func ConstructIPv6MPReachNLRI(protoFamily uint32, nextHop, nextHopLinkLocal net.IP,
	nlriList []NLRI) *BGPPathAttrMPReachNLRI {
	afi, safi := GetAfiSafi(protoFamily)
	mpReachNLRI := NewBGPPathAttrMPReachNLRI()
	mpReachNLRI.AFI = afi
	mpReachNLRI.SAFI = safi
	mpNextHop := NewMPNextHopIP6()
	mpNextHop.SetGlobalNextHop(nextHop)
	if nextHopLinkLocal != nil && nextHopLinkLocal.To16() == nil {
		mpNextHop.SetLinkLocalNextHop(nextHopLinkLocal)
	}
	mpReachNLRI.SetNextHop(mpNextHop)
	mpReachNLRI.SetNLRIList(nlriList)
	return mpReachNLRI
}

func CloneMPReachNLRIWithNewNLRI(mpReachNLRI *BGPPathAttrMPReachNLRI, nlri []NLRI) *BGPPathAttrMPReachNLRI {
	newMPReachNLRI := NewBGPPathAttrMPReachNLRI()
	newMPReachNLRI.AFI = mpReachNLRI.AFI
	newMPReachNLRI.SAFI = mpReachNLRI.SAFI
	newMPReachNLRI.SetNextHop(mpReachNLRI.NextHop)
	newMPReachNLRI.SetNLRIList(nlri)
	return newMPReachNLRI
}

func ConstructIPv6MPReachNLRIForConnRoutes(protoFamily uint32) *BGPPathAttrMPReachNLRI {
	return ConstructIPv6MPReachNLRI(protoFamily, net.IPv6zero, nil, nil)
}

func AddOriginatorId(updateMsg *BGPMessage, id net.IP) bool {
	body := updateMsg.Body.(*BGPUpdate)
	var pa BGPPathAttr

	for _, pa = range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeOriginatorId {
			return false
		}
	}

	idx := -1
	for idx, pa = range body.PathAttributes {
		if pa.GetCode() > BGPPathAttrTypeOriginatorId {
			break
		} else if idx == len(body.PathAttributes)-1 {
			idx += 1
		}
	}

	if idx >= 0 {
		paOriginatorId := NewBGPPathAttrOriginatorId(id)
		body.PathAttributes = append(body.PathAttributes[:idx], paOriginatorId)
		copy(body.PathAttributes[idx+1:], body.PathAttributes[idx:])
		body.PathAttributes[idx] = paOriginatorId
	}

	return true
}

func RemoveOriginatorId(updateMsg *BGPMessage) {
	removePathAttr(updateMsg, BGPPathAttrTypeOriginatorId)
}

func AddClusterId(updateMsg *BGPMessage, id uint32) bool {
	body := updateMsg.Body.(*BGPUpdate)
	var pa BGPPathAttr
	var i int
	found := false
	idx := -1

	for i, pa = range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeClusterList {
			idx = i
			found = true
			break
		} else if idx == -1 {
			if pa.GetCode() > BGPPathAttrTypeClusterList {
				idx = i
			} else if i == len(body.PathAttributes)-1 {
				idx = i + 1
			}
		}
	}

	if !found && idx >= 0 {
		clusterList := NewBGPPathAttrClusterList()
		body.PathAttributes = append(body.PathAttributes[:idx], clusterList)
		copy(body.PathAttributes[idx+1:], body.PathAttributes[idx:])
		body.PathAttributes[idx] = clusterList
	}

	if idx >= 0 {
		body.PathAttributes[idx].(*BGPPathAttrClusterList).PrependId(id)
		return true
	}

	return false
}

func RemoveClusterList(updateMsg *BGPMessage) {
	removePathAttr(updateMsg, BGPPathAttrTypeClusterList)
}

func ConvertIPBytesToUint(bytes []byte) uint32 {
	return uint32(bytes[0])<<24 | uint32(bytes[1]<<16) | uint32(bytes[2]<<8) | uint32(bytes[3])
}

func ConstructOptParams(as uint32, afiSAfiMap map[uint32]bool, addPathsRx bool, addPathsMaxTx uint8) []BGPOptParam {
	optParams := make([]BGPOptParam, 0)
	capParams := make([]BGPCapability, 0)

	cap4ByteASPath := NewBGPCap4ByteASPath(as)
	capParams = append(capParams, cap4ByteASPath)
	capAddPaths := NewBGPCapAddPath()
	addPathFlags := uint8(0)
	if addPathsRx {
		addPathFlags |= BGPCapAddPathRx
	}
	if addPathsMaxTx > 0 {
		addPathFlags |= BGPCapAddPathTx
	}

	for protoFamily, _ := range afiSAfiMap {
		afi, safi := GetAfiSafi(protoFamily)
		utils.Logger.Infof("Advertising capability for afi %d safi %d", afi, safi)
		capAfiSafi := NewBGPCapMPExt(afi, safi)
		capParams = append(capParams, capAfiSafi)

		addPathAfiSafi := NewAddPathAFISAFI(afi, safi, addPathFlags)
		capAddPaths.AddAddPathAFISAFI(addPathAfiSafi)
	}

	if addPathFlags != 0 {
		utils.Logger.Infof("Advertising capability for addPaths %+v", capAddPaths.Value)
		capParams = append(capParams, capAddPaths)
	}

	optCapability := NewBGPOptParamCapability(capParams)
	optParams = append(optParams, optCapability)

	return optParams
}

func GetASSize(openMsg *BGPOpen) uint8 {
	for _, optParam := range openMsg.OptParams {
		if optParam.GetCode() == BGPOptParamTypeCapability {
			capabilities := optParam.(*BGPOptParamCapability)
			for _, capability := range capabilities.Value {
				if capability.GetCode() == BGPCapTypeAS4Path {
					return 4
				}
			}
		}
	}

	return 2
}

func GetPeerAS(openMsg *BGPOpen) uint32 {
	var as uint32
	as = openMsg.MyAS
	for _, optParam := range openMsg.OptParams {
		if optParam.GetCode() == BGPOptParamTypeCapability {
			capabilities := optParam.(*BGPOptParamCapability)
			for _, capability := range capabilities.Value {
				if capability.GetCode() == BGPCapTypeAS4Path {
					as4Path := capability.(*BGPCapAS4Path)
					as = as4Path.Value
				}
			}
		}
	}

	return as
}

func GetAddPathFamily(openMsg *BGPOpen) map[AFI]map[SAFI]uint8 {
	addPathFamily := make(map[AFI]map[SAFI]uint8)
	for _, optParam := range openMsg.OptParams {
		if capabilities, ok := optParam.(*BGPOptParamCapability); ok {
			for _, capability := range capabilities.Value {
				if addPathCap, ok := capability.(*BGPCapAddPath); ok {
					utils.Logger.Infof("add path capability = %+v", addPathCap)
					for _, val := range addPathCap.Value {
						if _, ok := addPathFamily[val.AFI]; !ok {
							addPathFamily[val.AFI] = make(map[SAFI]uint8)
						}
						if _, ok := addPathFamily[val.AFI][val.SAFI]; !ok {
							addPathFamily[val.AFI][val.SAFI] = val.Flags
						}
					}
					return addPathFamily
				}
			}
		}
	}
	return addPathFamily
}

func IsAddPathsTxEnabledForIPv4(addPathFamily map[AFI]map[SAFI]uint8) bool {
	enabled := false
	/*
		if _, ok := addPathFamily[AfiIP]; ok {
			for safi, flags := range addPathFamily[AfiIP] {
				if (safi == SafiUnicast || safi == SafiMulticast) && (flags&BGPCapAddPathTx != 0) {
					utils.Logger.Infof("isAddPathsTxEnabledForIPv4 - add path Tx enabled for IPv4")
					enabled = true
				}
			}
		}
	*/
	for afi, _ := range addPathFamily {
		for _, flags := range addPathFamily[afi] {
			if flags&BGPCapAddPathTx != 0 {
				utils.Logger.Infof("isAddPathsTxEnabledForIPv4 - add path Tx enabled for IPv4")
				enabled = true
			}
		}
	}

	return enabled
}

func GetNumASesByASType(updateMsg *BGPMessage, asType BGPPathAttrType) uint32 {
	var total uint32 = 0

	if asType != BGPPathAttrTypeASPath && asType != BGPPathAttrTypeAS4Path {
		return total
	}

	body := updateMsg.Body.(*BGPUpdate)
	for _, pa := range body.PathAttributes {
		if pa.GetCode() == asType {
			if pa.GetCode() == BGPPathAttrTypeASPath {
				asPaths := pa.(*BGPPathAttrASPath).Value
				for _, asPath := range asPaths {
					if asPath.GetType() == BGPASPathSegmentSet {
						total += 1
					} else if asPath.GetType() == BGPASPathSegmentSequence {
						total += uint32(asPath.GetLen())
					}
				}
				break
			} else if pa.GetCode() == BGPPathAttrTypeAS4Path {
				asPaths := pa.(*BGPPathAttrAS4Path).Value
				for _, asPath := range asPaths {
					if asPath.GetType() == BGPASPathSegmentSet {
						total += 1
					} else if asPath.GetType() == BGPASPathSegmentSequence {
						total += uint32(asPath.GetLen())
					}
				}
				break
			}
		}
	}

	return total
}

func ConvertAS2ToAS4(updateMsg *BGPMessage) {
	body := updateMsg.Body.(*BGPUpdate)
	for idx, pa := range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeASPath {
			asPath := pa.(*BGPPathAttrASPath)
			newASPath := NewBGPPathAttrASPath()
			newASPath.ASSize = 4
			for _, seg := range asPath.Value {
				as2Seg := seg.(*BGPAS2PathSegment)
				as4Seg := NewBGPAS4PathSegmentSeq()
				as4Seg.Type = as2Seg.Type
				as4Seg.Length = as2Seg.GetLen()
				as4Seg.BGPASPathSegmentLen += (uint16(as2Seg.GetLen()) * 4)
				as4Seg.AS = make([]uint32, as4Seg.Length)
				for i, as := range as2Seg.AS {
					as4Seg.AS[i] = uint32(as)
				}
				newASPath.AppendASPathSegment(as4Seg)
			}
			body.PathAttributes[idx] = nil
			body.PathAttributes[idx] = newASPath
			break
		}
	}
}

func ConstructASPathFromAS4Path(asPath *BGPPathAttrASPath, as4Path *BGPPathAttrAS4Path, skip uint16) *BGPPathAttrASPath {
	utils.Logger.Info("ConstructASPathFromAS4Path")
	var segIdx int
	var segment BGPASPathSegment
	var asNum uint16 = 0
	newASPath := NewBGPPathAttrASPath()
	newASPath.ASSize = 4
	for segIdx, segment = range asPath.Value {
		if (uint16(segment.GetNumASes()) + asNum) > skip {
			break
		}
		newASPath.AppendASPathSegment(segment)
		asNum += uint16(segment.GetNumASes())
	}

	for idx, segment := range as4Path.Value {
		if idx == 0 {
			seg := asPath.Value[segIdx].(*BGPAS4PathSegment)
			segment.PrependASList(seg.AS[:skip-asNum])
		}
		newASPath.AppendASPathSegment(segment)
	}

	return newASPath
}

func Convert4ByteTo2ByteASPath(updateMsg *BGPMessage) {
	body := updateMsg.Body.(*BGPUpdate)
	for idx, pa := range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeASPath {
			asPath := pa.(*BGPPathAttrASPath)
			addAS4Path := false
			newAS4Path := asPath.CloneAsAS4Path()
			newAS2Path := NewBGPPathAttrASPath()
			for _, seg := range asPath.Value {
				as4Seg := seg.(*BGPAS4PathSegment)
				as2Seg, mappable := as4Seg.CloneAsAS2PathSegment()
				if !mappable {
					addAS4Path = true
				}
				newAS2Path.AppendASPathSegment(as2Seg)
			}
			body.PathAttributes[idx] = nil
			body.PathAttributes[idx] = newAS2Path
			if addAS4Path {
				addPathAttr(updateMsg, BGPPathAttrTypeAS4Path, newAS4Path)
			}
			break
		}
	}
}

func NormalizeASPath(updateMsg *BGPMessage, data interface{}) {
	var asPath *BGPPathAttrASPath
	var as4Path *BGPPathAttrAS4Path
	var asAggregator *BGPPathAttrAggregator
	var as4Aggregator *BGPPathAttrAS4Aggregator

	body := updateMsg.Body.(*BGPUpdate)
	if body.TotalPathAttrLen == 0 {
		return
	}

	for _, pa := range body.PathAttributes {
		if pa.GetCode() == BGPPathAttrTypeASPath {
			utils.Logger.Debug("NormalizeASPath ASPath is set")
			asPath = pa.(*BGPPathAttrASPath)
		} else if pa.GetCode() == BGPPathAttrTypeAS4Path {
			utils.Logger.Debug("NormalizeASPath AS4Path is set")
			as4Path = pa.(*BGPPathAttrAS4Path)
		} else if pa.GetCode() == BGPPathAttrTypeAggregator {
			utils.Logger.Debug("NormalizeASPath Aggregator is set")
			asAggregator = pa.(*BGPPathAttrAggregator)
		} else if pa.GetCode() == BGPPathAttrTypeAS4Aggregator {
			utils.Logger.Debug("NormalizeASPath AS4Aggregator is set")
			as4Aggregator = pa.(*BGPPathAttrAS4Aggregator)
		}
	}

	if asPath == nil {
		utils.Logger.Err("***** BGP update message does not have AS path *****")
		return
	}

	if asPath.ASSize == 2 {
		if asAggregator != nil && as4Aggregator != nil && uint16(asAggregator.AS.GetAS()) != BGPASTrans {
			utils.Logger.Info("NormalizeASPath agg and as4agg is nil and agg.AS != BGPASTrans")
			removePathAttr(updateMsg, BGPPathAttrTypeAS4Aggregator)
			removePathAttr(updateMsg, BGPPathAttrTypeAS4Path)
		} else {
			utils.Logger.Info("NormalizeASPath calling ConvertAS2ToAS4")
			ConvertAS2ToAS4(updateMsg)
			if as4Path != nil {
				numASes := GetNumASesByASType(updateMsg, BGPPathAttrTypeASPath)
				numAS4es := GetNumASesByASType(updateMsg, BGPPathAttrTypeAS4Path)
				utils.Logger.Info("NormalizeASPath as4Path != nil, numASes=", numASes, "numAS4es=", numAS4es)
				if numASes >= numAS4es {
					utils.Logger.Info("NormalizeASPath numASes >= numAS4es")
					newASPath := ConstructASPathFromAS4Path(asPath, as4Path, uint16(numASes-numAS4es))
					removePathAttr(updateMsg, BGPPathAttrTypeAS4Path)
					removePathAttr(updateMsg, BGPPathAttrTypeASPath)
					addPathAttr(updateMsg, BGPPathAttrTypeASPath, newASPath)
				}
			}
		}
	} else if asPath.ASSize == 4 {
		utils.Logger.Debug("NormalizeASPath asPath.ASSize == 4")
		if as4Aggregator != nil {
			removePathAttr(updateMsg, BGPPathAttrTypeAS4Aggregator)
		}
		if as4Path != nil {
			removePathAttr(updateMsg, BGPPathAttrTypeAS4Path)
		}
	}
}

func ConstructMaxSizedUpdatePackets(bgpMsg *BGPMessage) []*BGPMessage {
	var withdrawnRoutes []NLRI
	newUpdateMsgs := make([]*BGPMessage, 0)
	pktLen := uint32(BGPUpdateMsgMinLen)
	startIdx := 0
	lastIdx := 0
	updateMsg := bgpMsg.Body.(*BGPUpdate)
	pathAttrs := make([]BGPPathAttr, 0)

	if updateMsg.WithdrawnRoutes != nil {
		for lastIdx = 0; lastIdx < len(updateMsg.WithdrawnRoutes); lastIdx++ {
			nlriLen := updateMsg.WithdrawnRoutes[lastIdx].Len()
			if nlriLen+pktLen > BGPMsgMaxLen {
				newMsg := NewBGPUpdateMessage(updateMsg.WithdrawnRoutes[startIdx:lastIdx], nil, nil)
				utils.Logger.Debug("ConstructMaxSizedUpdatePackets - Constructed BGP update message for withdrawn routes")
				newUpdateMsgs = append(newUpdateMsgs, newMsg)
				startIdx = lastIdx
				pktLen = uint32(BGPUpdateMsgMinLen)
			}
			pktLen += nlriLen
		}
	}

	if lastIdx > startIdx {
		withdrawnRoutes = updateMsg.WithdrawnRoutes[startIdx:lastIdx]
	}

	paLen := uint32(0)
	for i := 0; i < len(updateMsg.PathAttributes); i++ {
		paLen += updateMsg.PathAttributes[i].TotalLen()
	}

	if pktLen+paLen > BGPMsgMaxLen {
		mpReach, mpUnreach := GetMPAttrs(updateMsg.PathAttributes)
		utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - pktLen %d mpReach %+v, mpUnreach %+v", pktLen, mpReach, mpUnreach)
		if mpUnreach != nil {
			utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - mpUnreach is not nil total len %d", mpUnreach.TotalLen())
			removeTypeFromPathAttrs(&updateMsg.PathAttributes, BGPPathAttrTypeMPUnreachNLRI)
			mpUnreachBaseLen := uint32(mpUnreach.BGPPathAttrLen) + 3 // Path attr base len + (AFI + SAFI) len
			startIdx = 0
			lastIdx = 0
			if len(mpUnreach.NLRI) > 0 && pktLen+mpUnreach.TotalLen() > BGPMsgMaxLen {
				for lastIdx = 0; lastIdx < len(mpUnreach.NLRI); lastIdx++ {
					nlriLen := mpUnreach.NLRI[lastIdx].Len()
					if nlriLen+pktLen+mpUnreachBaseLen > BGPMsgMaxLen {
						pa := make([]BGPPathAttr, 0)
						if lastIdx != 0 {
							mpUnreachPA := ConstructMPUnreachNLRI(mpUnreach.AFI, mpUnreach.SAFI, mpUnreach.NLRI[startIdx:lastIdx])
							pa = append(pa, mpUnreachPA)
						}
						newMsg := NewBGPUpdateMessage(withdrawnRoutes, pa, nil)
						utils.Logger.Debug("ConstructMaxSizedUpdatePackets - Constructed BGP update message for MPUnreach routes")
						newUpdateMsgs = append(newUpdateMsgs, newMsg)
						withdrawnRoutes = nil
						startIdx = lastIdx
						pktLen = uint32(BGPUpdateMsgMinLen)
					}
					pktLen += nlriLen
				}
			}
			if startIdx != lastIdx {
				utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - mpUnreach startIdx %d != lastIdx %d", startIdx, lastIdx)
				mpUnreachPA := ConstructMPUnreachNLRI(mpUnreach.AFI, mpUnreach.SAFI, mpUnreach.NLRI[startIdx:lastIdx])
				pathAttrs = append(pathAttrs, mpUnreachPA)
				pktLen += mpUnreachBaseLen
			}
		}

		if mpReach != nil {
			removeTypeFromPathAttrs(&updateMsg.PathAttributes, BGPPathAttrTypeMPReachNLRI)
			startIdx = 0
			lastIdx = 0
			otherPAs := CopyPathAttrs(updateMsg.PathAttributes)
			otherPAsLen := uint32(0)
			for i := 0; i < len(updateMsg.PathAttributes); i++ {
				otherPAsLen += updateMsg.PathAttributes[i].TotalLen()
			}
			utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - mpReach is not nil total len %d, pkt len %d, other PAs len %d",
				mpReach.TotalLen(), pktLen, otherPAsLen)
			if len(mpReach.NLRI) > 0 && pktLen+otherPAsLen+mpReach.TotalLen() > BGPMsgMaxLen {
				// Path attr base len + (AFI + SAFI + Reserved) len + Next hop len + path attrs len
				mpReachBaseLen := uint32(mpReach.BGPPathAttrLen) + 4 + uint32(mpReach.NextHop.Len()) + otherPAsLen
				for lastIdx = 0; lastIdx < len(mpReach.NLRI); lastIdx++ {
					nlriLen := mpReach.NLRI[lastIdx].Len()
					if nlriLen+pktLen+mpReachBaseLen > BGPMsgMaxLen {
						pa := make([]BGPPathAttr, 0)
						if lastIdx != 0 {
							mpReachPA := CloneMPReachNLRIWithNewNLRI(mpReach, mpReach.NLRI[startIdx:lastIdx])
							pa = AddMPReachNLRIToPathAttrs(otherPAs, mpReachPA)
						}
						if len(pathAttrs) > 0 {
							for idx, _ := range pathAttrs {
								pa = addPathAttrToPathAttrs(pa, pathAttrs[idx])
							}
							pathAttrs = make([]BGPPathAttr, 0)
						}
						newMsg := NewBGPUpdateMessage(withdrawnRoutes, pa, nil)
						utils.Logger.Debug("ConstructMaxSizedUpdatePackets - Constructed BGP update message for MPReach routes")
						newUpdateMsgs = append(newUpdateMsgs, newMsg)
						withdrawnRoutes = nil
						startIdx = lastIdx
						pktLen = uint32(BGPUpdateMsgMinLen)
					}
					pktLen += nlriLen
				}
			}
			if startIdx != lastIdx {
				utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - mpReach startIdx %d != lastIdx %d", startIdx, lastIdx)
				mpReachPA := CloneMPReachNLRIWithNewNLRI(mpReach, mpReach.NLRI[startIdx:lastIdx])
				pathAttrs = AddMPReachNLRIToPathAttrs(pathAttrs, mpReachPA)
				pktLen += uint32(mpReach.BGPPathAttrLen) + 4 + uint32(mpReach.NextHop.Len())
			}
		}
	}

	paLen = uint32(0)
	for i := 0; i < len(updateMsg.PathAttributes); i++ {
		paLen += updateMsg.PathAttributes[i].TotalLen()
	}
	utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - paLen %d pktLen %d", paLen, pktLen)

	if pktLen+paLen > BGPMsgMaxLen {
		utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - pktLen+paLen > BGPMsgMaxLen")
		otherPAs := make([]BGPPathAttr, 0)
		mpReachPA := getTypeFromPathAttrs(pathAttrs, BGPPathAttrTypeMPReachNLRI)
		if mpReachPA != nil {
			utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - mpReachPA is not nil, copy path attrs")
			otherPAs = CopyPathAttrs(updateMsg.PathAttributes)
		}
		if len(pathAttrs) > 0 {
			utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - pathAttrs %+v", pathAttrs)
			for idx, _ := range pathAttrs {
				otherPAs = addPathAttrToPathAttrs(otherPAs, pathAttrs[idx])
			}
			pathAttrs = make([]BGPPathAttr, 0)
		}
		newMsg := NewBGPUpdateMessage(withdrawnRoutes, otherPAs, nil)
		utils.Logger.Debugf("ConstructMaxSizedUpdatePackets - Constructed BGP update message for withdrawn routes... last message")
		withdrawnRoutes = nil
		newUpdateMsgs = append(newUpdateMsgs, newMsg)
		pktLen = BGPUpdateMsgMinLen
	}
	mpAttsFound := HasMPAttrs(updateMsg.PathAttributes)

	startIdx = 0
	lastIdx = 0
	for lastIdx = 0; lastIdx < len(updateMsg.NLRI); lastIdx++ {
		nlriLen := updateMsg.NLRI[lastIdx].Len()
		if nlriLen+pktLen+paLen > BGPMsgMaxLen {
			pa := updateMsg.PathAttributes
			if len(pathAttrs) > 0 {
				for idx, _ := range pathAttrs {
					pa = addPathAttrToPathAttrs(pa, pathAttrs[idx])
				}
				pathAttrs = make([]BGPPathAttr, 0)
			}
			newMsg := NewBGPUpdateMessage(withdrawnRoutes, pa, updateMsg.NLRI[startIdx:lastIdx])
			utils.Logger.Debug("ConstructMaxSizedUpdatePackets - Constructed BGP update message for updated routes")
			newUpdateMsgs = append(newUpdateMsgs, newMsg)
			withdrawnRoutes = nil
			startIdx = lastIdx
			pktLen = uint32(BGPUpdateMsgMinLen)
			mpAttsFound = false
		}
		pktLen += nlriLen
	}

	if (withdrawnRoutes != nil && len(withdrawnRoutes) > 0) || (lastIdx > startIdx) || len(pathAttrs) > 0 || mpAttsFound {
		utils.Logger.Debug("ConstructMaxSizedUpdatePackets - Constructed BGP update message for updated routes... last message")
		pa := updateMsg.PathAttributes
		if len(pathAttrs) > 0 {
			for idx, _ := range pathAttrs {
				pa = addPathAttrToPathAttrs(pa, pathAttrs[idx])
			}
			pathAttrs = make([]BGPPathAttr, 0)
		}
		newMsg := NewBGPUpdateMessage(withdrawnRoutes, pa, updateMsg.NLRI[startIdx:lastIdx])
		newUpdateMsgs = append(newUpdateMsgs, newMsg)
	}

	return newUpdateMsgs
}
