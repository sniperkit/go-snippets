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

package server

import (
	"time"
)

type LsdbKey struct {
	AreaId uint32
}

const (
	RouterLSA     uint8 = 1
	NetworkLSA    uint8 = 2
	Summary3LSA   uint8 = 3
	Summary4LSA   uint8 = 4
	ASExternalLSA uint8 = 5
)

type LsaKey struct {
	LSType    uint8  /* LS Type */
	LSId      uint32 /* Link State Id */
	AdvRouter uint32 /* Avertising Router */
}

func NewLsaKey() *LsaKey {
	return &LsaKey{}
}

type LSDatabase struct {
	RouterLsaMap     map[LsaKey]RouterLsa
	NetworkLsaMap    map[LsaKey]NetworkLsa
	Summary3LsaMap   map[LsaKey]SummaryLsa
	Summary4LsaMap   map[LsaKey]SummaryLsa
	ASExternalLsaMap map[LsaKey]ASExternalLsa
}

type SelfOrigLsa map[LsaKey]bool

type LsdbCtrlChStruct struct {
	LsdbGblCtrlCh       chan bool
	LsdbGblCtrlReplyCh  chan bool
	LsdbAreaCtrlCh      chan uint32
	LsdbAreaCtrlReplyCh chan uint32
}

type RouteInfo struct {
	NwAddr  uint32
	Netmask uint32
	Metric  uint32
}

type LsdbStruct struct {
	AreaLsdb        map[LsdbKey]LSDatabase
	AreaSelfOrigLsa map[LsdbKey]SelfOrigLsa
	LsdbCtrlChData  LsdbCtrlChStruct
	LsdbAgingTicker *time.Ticker
	ExtRouteInfoMap map[RouteInfo]bool
}
