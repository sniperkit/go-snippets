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
	_ "fmt"
	"l3/bgp/utils"
)

type ASPathIter struct {
	asPath     *BGPPathAttrASPath
	segments   []*BGPAS4PathSegment
	segmentLen int
	segIdx     int
	asValIdx   int
}

func NewASPathIter(asPath *BGPPathAttrASPath) *ASPathIter {
	iter := ASPathIter{
		asPath:     asPath,
		segmentLen: len(asPath.Value),
	}

	iter.segments = make([]*BGPAS4PathSegment, 0, len(asPath.Value))
	for idx := 0; idx < len(asPath.Value); idx++ {
		var as4Seg *BGPAS4PathSegment
		var ok bool
		if as4Seg, ok = asPath.Value[idx].(*BGPAS4PathSegment); !ok {
			utils.Logger.Err("AS path segment", idx, "is not AS4PathSegment")
			return nil
		}
		iter.segments = append(iter.segments, as4Seg)
	}
	return &iter
}

func RemoveNilItemsFromList(iterList []*ASPathIter) []*ASPathIter {
	lastIdx := len(iterList) - 1
	var modIdx, idx int
	for idx = 0; idx < len(iterList); idx++ {
		if iterList[idx] == nil {
			for modIdx = lastIdx; modIdx > idx && iterList[modIdx] == nil; modIdx-- {
			}
			if modIdx <= idx {
				break
			}
			iterList[idx] = iterList[modIdx]
			iterList[modIdx] = nil
			lastIdx = modIdx
		}
	}

	return iterList[:idx]
}

func (a *ASPathIter) Next() (val uint32, segType BGPASPathSegmentType, flag bool) {
	if a.segIdx >= a.segmentLen {
		return val, segType, flag
	}

	val = a.segments[a.segIdx].AS[a.asValIdx]
	segType = a.segments[a.segIdx].Type
	flag = true

	a.asValIdx++
	if a.asValIdx >= len(a.segments[a.segIdx].AS) {
		a.segIdx++
		a.asValIdx = 0
	}

	return val, segType, flag
}
