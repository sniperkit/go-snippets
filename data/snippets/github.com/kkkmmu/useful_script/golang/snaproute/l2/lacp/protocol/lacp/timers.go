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
package lacp

import (
	"time"
)

// WaitWhileTimerStart
// Start the timer
func (muxm *LacpMuxMachine) WaitWhileTimerStart() {
	if muxm.waitWhileTimer == nil {
		muxm.waitWhileTimer = time.NewTimer(muxm.waitWhileTimerTimeout)
	} else {
		muxm.waitWhileTimer.Reset(muxm.waitWhileTimerTimeout)
	}
	muxm.waitWhileTimerRunning = true
}

// WaitWhileTimerStop
// Stop the timer, which should only happen
// on creation as well as when the lacp mode is "on"
func (muxm *LacpMuxMachine) WaitWhileTimerStop() {
	if muxm.waitWhileTimer != nil {
		muxm.waitWhileTimer.Stop()
		muxm.waitWhileTimerRunning = false
	}
}

func (muxm *LacpMuxMachine) WaitWhileTimerTimeoutSet(timeout time.Duration) {
	muxm.waitWhileTimerTimeout = timeout
}

func (rxm *LacpRxMachine) CurrentWhileTimerStart() {
	if rxm.currentWhileTimer == nil {
		rxm.currentWhileTimer = time.NewTimer(rxm.currentWhileTimerTimeout)
	} else {
		rxm.currentWhileTimer.Reset(rxm.currentWhileTimerTimeout)
	}
}

func (rxm *LacpRxMachine) CurrentWhileTimerStop() {
	if rxm.currentWhileTimer != nil {
		rxm.currentWhileTimer.Stop()
	}
}

func (rxm *LacpRxMachine) CurrentWhileTimerTimeoutSet(timeout time.Duration) {
	rxm.currentWhileTimerTimeout = timeout
}

func (ptxm *LacpPtxMachine) PeriodicTimerStart() {
	if ptxm.periodicTxTimer == nil {
		ptxm.periodicTxTimer = time.NewTimer(ptxm.PeriodicTxTimerInterval)
	} else {
		ptxm.periodicTxTimer.Reset(ptxm.PeriodicTxTimerInterval)
	}
}

func (ptxm *LacpPtxMachine) PeriodicTimerStop() {
	if ptxm.periodicTxTimer != nil {
		ptxm.periodicTxTimer.Stop()
	}
}

func (ptxm *LacpPtxMachine) PeriodicTimerIntervalSet(interval time.Duration) {
	ptxm.PeriodicTxTimerInterval = interval
}

func (cdm *LacpCdMachine) ChurnDetectionTimerStart() {
	if cdm.churnTimer == nil {
		cdm.churnTimer = time.NewTimer(cdm.churnTimerInterval)
	} else {
		cdm.churnTimer.Reset(cdm.churnTimerInterval)
	}
}

func (cdm *LacpCdMachine) ChurnDetectionTimerStop() {
	if cdm.churnTimer != nil {
		cdm.churnTimer.Stop()
	}
}

func (cdm *LacpCdMachine) ChurnDetectionTimerIntervalSet(interval time.Duration) {
	cdm.churnTimerInterval = interval
}

// TxGuardTimerStart used by Tx Machine as described in
// 802.1ax-2014 Section 6.4.17 in order to not transmit
// more than 3 packets in this interval
func (txm *LacpTxMachine) TxGuardTimerStart() {
	//if txm.p.begin == false {
	//	txm.LacpTxmLog("Starting Guard Timer")
	//}
	if txm.txGuardTimer == nil {
		txm.txGuardTimer = time.AfterFunc(LacpFastPeriodicTime, txm.LacpTxGuardGeneration)
	} else {
		txm.txGuardTimer.Reset(LacpFastPeriodicTime)
	}
}

// TxDelayTimerStop to stop the Delay timer
// in case a port is deleted or initialized
func (txm *LacpTxMachine) TxGuardTimerStop() {
	if txm.txGuardTimer != nil {
		txm.txGuardTimer.Stop()
	}
}
