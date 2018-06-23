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

// defs
package lacp

import (
	"time"
)

// 6.4.4 Constants
// number of seconds between periodic trasmissions using Short Timeouts
const LacpFastPeriodicTime time.Duration = (time.Second * 1)

// number of seconds etween periodic transmissions using Long timeouts
const LacpSlowPeriodicTime time.Duration = (time.Second * 30)

// number of seconds before invalidating received LACPDU info when using
// Short Timeouts (3 x LacpFastPeriodicTime)
// Lacp State Timeout == 1
const LacpShortTimeoutTime time.Duration = (time.Second * 3)

// number of seconds before invalidating received LACPDU info when using
// Long Timeouts (3 x LacpSlowPeriodicTime)
// Lacp State Timeout == 0
const LacpLongTimeoutTime time.Duration = (time.Second * 90)

// number of seconds that the Actor and Partner Churn State machines
// wait for the Actor or Partner Sync State to stabilize
const LacpChurnDetectionTime time.Duration = (time.Second * 60)

// number of seconds to delay aggregation to allow multiple links to
// aggregate simultaneously
const LacpAggregateWaitTime time.Duration = (time.Second * 2)

// the version number of the Actor LACP implementation
const LacpActorSystemLacpVersion int = 0x01

const LacpPortDuplexFull int = 1
const LacpPortDuplexHalf int = 2

const LacpIsEnabled bool = true
const LacpIsDisabled bool = false

const (
	LacpStateActivityBit = 1 << iota
	LacpStateTimeoutBit
	LacpStateAggregationBit
	LacpStateSyncBit
	LacpStateCollectingBit
	LacpStateDistributingBit
	LacpStateDefaultedBit
	LacpStateExpiredBit
)

// default actor
const LacpStateIndividual uint8 = (LacpStateDefaultedBit | LacpStateActivityBit)

// default partner
const LacpStateAggregatibleUp uint8 = (LacpStateActivityBit |
	LacpStateAggregationBit |
	LacpStateSyncBit |
	LacpStateCollectingBit |
	LacpStateDistributingBit |
	LacpStateDefaultedBit)

// default partner State after lacp pdu's received
const LacpStateAggregatibleDown uint8 = (LacpStateActivityBit |
	LacpStateAggregationBit |
	LacpStateDefaultedBit)

const (
	// also known as manual mode
	LacpModeOn = iota + 1
	// lacp State Activity == TRUE
	// considered lacp enabled
	LacpModeActive
	// lacp State Activity == FALSE
	// considered lacp enabled
	LacpModePassive
)

func LacpStateSet(currState *uint8, StateBits uint8) {
	*currState |= StateBits
}

func LacpStateClear(currState *uint8, StateBits uint8) {
	*currState &= ^(StateBits)
}

func LacpStateIsSet(currState uint8, StateBits uint8) bool {
	return (currState & StateBits) == StateBits
}

func LacpModeGet(currState uint8, lacpEnabled bool) int {
	mode := LacpModeOn
	if lacpEnabled {
		mode = LacpModePassive
		if LacpStateIsSet(currState, LacpStateActivityBit) {
			mode = LacpModeActive
		}
	}
	return mode
}

func DefsStrMapsCreate() {
	MuxStateStrMap = make(map[uint8]string)
	MuxStateStrMap[LacpStateActivityBit] = "Activity"
	MuxStateStrMap[LacpStateTimeoutBit] = "Timeout"
	MuxStateStrMap[LacpStateAggregationBit] = "Aggregation"
	MuxStateStrMap[LacpStateSyncBit] = "Sync"
	MuxStateStrMap[LacpStateCollectingBit] = "Collecting"
	MuxStateStrMap[LacpStateDistributingBit] = "Distributing"
	MuxStateStrMap[LacpStateDefaultedBit] = "Defaulted"
	MuxStateStrMap[LacpStateExpiredBit] = "Expired"

	ModeStrMap = make(map[uint8]string)
	ModeStrMap[LacpModeActive] = "Active"
	ModeStrMap[LacpModePassive] = "Passive"
	ModeStrMap[LacpModeOn] = "On"
}

func LacpStateToStr(state uint8) string {

	var statestr = ""
	if LacpStateIsSet(state, LacpStateActivityBit) {
		statestr += MuxStateStrMap[LacpStateActivityBit] + ","
	}
	if LacpStateIsSet(state, LacpStateTimeoutBit) {
		statestr += MuxStateStrMap[LacpStateTimeoutBit] + ","
	}
	if LacpStateIsSet(state, LacpStateAggregationBit) {
		statestr += MuxStateStrMap[LacpStateAggregationBit] + ","
	}
	if LacpStateIsSet(state, LacpStateSyncBit) {
		statestr += MuxStateStrMap[LacpStateSyncBit] + ","
	}
	if LacpStateIsSet(state, LacpStateCollectingBit) {
		statestr += MuxStateStrMap[LacpStateCollectingBit] + ","
	}
	if LacpStateIsSet(state, LacpStateDistributingBit) {
		statestr += MuxStateStrMap[LacpStateDistributingBit] + ","
	}
	if LacpStateIsSet(state, LacpStateDefaultedBit) {
		statestr += MuxStateStrMap[LacpStateDefaultedBit] + ","
	}
	if LacpStateIsSet(state, LacpStateExpiredBit) {
		statestr += MuxStateStrMap[LacpStateExpiredBit] + ","
	}
	return statestr
}
