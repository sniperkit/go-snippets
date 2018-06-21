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
	"l3/ndp/config"
	"testing"
)

var nbrTest *NeighborInfo
var testPktDataCh chan config.PacketData

func initTestNbrInfo() {
	initServerBasic()
	nbrTest = &NeighborInfo{}
	testPktDataCh = make(chan config.PacketData)
	nbrTest.InitCache(testReachableTimerValue, testReTransmitTimerValue,
		testMyAbsLinkScopeIP+"_"+testSrcMac, testPktDataCh, testIfIndex)
}

func deinitTestNbrInfo() {
	nbrTest.DeInit()
}

func TestReTransmitTimer(t *testing.T) {
	initTestNbrInfo()
	if nbrTest.RetransTimer != nil {
		t.Error("Re-Transmit timer should not be started until reachable timer fires")
		return
	}
	nbrTest.Timer()
	nbrTest.StopReTransmitTimer()
	if nbrTest.RetransTimer != nil {
		t.Error("Failed to stop re-transmit timer")
		return
	}

	deinitTestNbrInfo()
	if nbrTest.RetransTimer != nil {
		t.Error("Failed to de-init Re-transmit timer via neigborInfo DeInit api")
		return
	}
	nbrTest = nil
}

func TestReachableTimer(t *testing.T) {
	initTestNbrInfo()
	if nbrTest.ReachableTimer == nil {
		t.Error("Failed to start reachable timer")
		return
	}
	nbrTest.RchTimer()
	nbrTest.StopReachableTimer()
	if nbrTest.ReachableTimer != nil {
		t.Error("Failed to stop Reachable timer")
		return
	}
}

func TestReComputerTimer(t *testing.T) {
	initTestNbrInfo()
	if nbrTest.RecomputeBaseTimer == nil {
		t.Error("Failed to start re-compute base timer")
		return
	}
	nbrTest.ReComputeBaseReachableTimer()
	nbrTest.StopReComputeBaseTimer()
	if nbrTest.RetransTimer != nil {
		t.Error("Failed to stop recompute base timer")
		return
	}
}

func TestDelayFirstTimer(t *testing.T) {
	initTestNbrInfo()
	nbrTest.DelayProbe()
	if nbrTest.DelayFirstProbeTimer == nil {
		t.Error("Failed to start delay probe timer")
		return
	}
	nbrTest.DelayProbe()
	if nbrTest.DelayFirstProbeTimer == nil {
		t.Error("Failed to reset delay probe timer")
		return
	}

	nbrTest.StopDelayProbeTimer()
	if nbrTest.DelayFirstProbeTimer != nil {
		t.Error("Failed to stop delay probe timer")
		return
	}
	// delete the nbrTest and then restart the proble
	nbrTest.DeInit()
	if nbrTest.DelayFirstProbeTimer != nil {
		t.Error("Failed to stop delay probe timer")
		return
	}
	nbrTest = nil

	initTestNbrInfo()
	nbrTest.DelayProbe()
	if nbrTest.DelayFirstProbeTimer == nil {
		t.Error("Failed to start delay probe timer")
		return
	}

	nbrTest.UpdateProbe()

	if nbrTest.ProbesSent != 0 || nbrTest.DelayFirstProbeTimer != nil || nbrTest.RetransTimer != nil {
		t.Error("failed to update probes")
		return
	}
	// delete the nbrTest and make sure delay probe is stopped in that
	nbrTest.DeInit()
	if nbrTest.DelayFirstProbeTimer != nil {
		t.Error("Failed to stop delay probe timer")
		return
	}
	nbrTest = nil
}

func TestInvaliTimer(t *testing.T) {
	initTestNbrInfo()
	if nbrTest.InvalidationTimer != nil {
		t.Error("Invalid timer should not be started during init neigbor info")
		return
	}
	nbrTest.InValidTimer(1800)

	if nbrTest.InvalidationTimer == nil {
		t.Error("Invalidation Timer should have been started but its not")
		return
	}
	nbrTest.StopInvalidTimer()
	if nbrTest.InvalidationTimer != nil {
		t.Error("Failed to stop invalidation timer")
		return
	}
}

func _TestFastProbeTimer(t *testing.T) {
	//initServerBasic()
	nbrTest = &NeighborInfo{}
	TestIPv6IntfCreate(t)
	testPktDataCh = make(chan config.PacketData)
	nbrTest.InitCache(1, testReTransmitTimerValue,
		testMyAbsLinkScopeIP+"_"+testSrcMac, testPktDataCh, testIfIndex)
	probes := 0
	for {
		select {
		case pktDataInfo, ok := <-testPktDataCh:
			if !ok {
				break
			}
			probes++
			//t.Log("Received probe", probes)
			testNdpServer.counter.Rcvd++
			if pktDataInfo.FastProbe == false {
				t.Error("For Fast Probe timer isFastProbe should be set to true")
				return
			}
			//t.Log("Processing Timer Expiry:", pktDataInfo)
			testNdpServer.ProcessTimerExpiry(pktDataInfo)
		}
		break
	}
}
