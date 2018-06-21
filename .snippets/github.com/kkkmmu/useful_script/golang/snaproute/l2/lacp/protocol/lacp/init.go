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

// init
package lacp

var LaSystemIdDefault LacpSystem
var MuxStateStrMap map[uint8]string
var ModeStrMap map[uint8]string

func init() {

	DefsStrMapsCreate()

	// Default System Id is all zero's
	// this will be used by all static lags, as well as initial
	// aggregation configs.
	LaSystemIdDefault = LacpSystem{
		Actor_System_priority: 0,
		Actor_System:          [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	LacpSysGlobalInfoInit(LaSystemIdDefault)

	LacpCbDb = LacpCbDbEntry{
		PortCreateDbList:  make(map[string]LacpPortEvtCb),
		PortDeleteDbList:  make(map[string]LacpPortEvtCb),
		PortUpDbList:      make(map[string]LacpPortEvtCb),
		PortDownDbList:    make(map[string]LacpPortEvtCb),
		AggCreateDbList:   make(map[string]LacpAggEvtCb),
		AggDeleteDbList:   make(map[string]LacpAggEvtCb),
		AggOperUpDbList:   make(map[string]LacpAggEvtCb),
		AggOperDownDbList: make(map[string]LacpAggEvtCb),
	}

	initHwPortCreateDelCb()

	ConfigAggMap = make(map[string]*LaAggConfig)
	ConfigAggList = make([]*LaAggConfig, 0)
}
