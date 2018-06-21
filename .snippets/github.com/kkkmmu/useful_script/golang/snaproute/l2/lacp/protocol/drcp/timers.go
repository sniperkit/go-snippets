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

// timers
package drcp

import (
	"time"
)

func (rxm *RxMachine) CurrentWhileTimerStart() {
	if rxm.currentWhileTimer == nil {
		rxm.currentWhileTimer = time.NewTimer(rxm.currentWhileTimerTimeout)
	} else {
		rxm.currentWhileTimer.Reset(rxm.currentWhileTimerTimeout)
	}
}

func (rxm *RxMachine) CurrentWhileTimerStop() {
	if rxm.currentWhileTimer != nil {
		rxm.currentWhileTimer.Stop()
	}
}

func (rxm *RxMachine) CurrentWhileTimerTimeoutSet(timeout time.Duration) {
	rxm.currentWhileTimerTimeout = timeout
}

func (ptxm *PtxMachine) PeriodicTimerStart() {
	if ptxm.periodicTimer == nil {
		ptxm.periodicTimer = time.NewTimer(ptxm.periodicTimerInterval)
	} else {
		ptxm.periodicTimer.Reset(ptxm.periodicTimerInterval)
	}
}

func (ptxm *PtxMachine) PeriodicTimerStop() {
	if ptxm.periodicTimer != nil {
		ptxm.periodicTimer.Stop()
	}
}

func (ptxm *PtxMachine) PeriodicTimerIntervalSet(interval time.Duration) {
	ptxm.periodicTimerInterval = interval
}
