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

// CHURN DETECTION MACHINE 802.1ax-2014 Section 6.4.17
package lacp

import (
	//"fmt"
	"l2/lacp/protocol/utils"
	"net"
	"testing"
	"time"
	"utils/fsm"
)

const LaAggChurnAgg1 = 100
const LaAggChurnAgg2 = 200
const LaAggChurnPortActor = 10
const LaAggChurnPortPeer = 21
const LaAggChurnPortActorIf = "SIMeth0"
const LaAggChurnPortPeerIf = "SIMeth1"

func ChurnDetectionStateMachineTeardown() {

	DeleteLaAggPort(LaAggChurnPortActor)
	OnlyForTestTeardown()
}

func ChurnDetectionStateMachineSetup() {

	OnlyForTestSetup()
	// must be called to initialize the global
	//LaSystemActor := LacpSystem{Actor_System_priority: 128,
	//	Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x64}}
	//LaSystemPeer := LacpSystem{Actor_System_priority: 128,
	//	Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	p1conf := &LaAggPortConfig{
		Id:      LaAggChurnPortActor,
		Prio:    0x80,
		Key:     100,
		AggId:   LaAggChurnAgg1,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggChurnPortActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggChurnPortActorIf,
		TraceEna: true,
	}

	utils.PortConfigMap[int32(p1conf.Id)] = utils.PortConfig{Name: LaAggChurnPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	utils.PortConfigMap[LaAggChurnPortPeer] = utils.PortConfig{Name: LaAggChurnPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x44, 0x44, 0x22, 0x22, 0x33},
	}

	// lets create a port and start the machines
	p1 := NewLaAggPort(p1conf)
	// disable the logger as it is dependent on
	// log server
	p1.logEna = false
	p1.LacpActorCdMachineMain()
	p1.LacpPartnerCdMachineMain()

	// since we are only concerned about the CDM machine lets set
	// reload to true so that we don't initialize all the other machines.
	p1.BEGIN(true)

}

func TestCdmNoActorChurnInvalidEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	invalidStateMap := [2]fsm.Event{
		LacpCdmEventActorOperPortStateSyncOn,
		LacpCdmEventActorChurnTimerExpired,
	}

	var p1 *LaAggPort
	responseChannel := make(chan string)
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStateActorChurnMonitor], "found ", CdmStateStrMap[p1.CdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			// set sync bit
			LacpStateSet(&p1.ActorOper.State, LacpStateSyncBit)
			p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            LacpCdmEventActorOperPortStateSyncOn,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
		}

		for _, evt := range invalidStateMap {
			// force state to no actor churn
			p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            evt,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
			if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateNoActorChurn {
				t.Error("Error Cdm State Machine transitioned on a invalid event")
			}
		}
	}
}

func TestCdmNoPartnerChurnInvalidEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	invalidStateMap := [2]fsm.Event{
		LacpCdmEventPartnerOperPortStateSyncOn,
		LacpCdmEventPartnerChurnTimerExpired,
	}

	var p1 *LaAggPort
	responseChannel := make(chan string)
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStatePartnerChurnMonitor], "found ", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			// set sync bit
			LacpStateSet(&p1.PartnerOper.State, LacpStateSyncBit)
			p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            LacpCdmEventPartnerOperPortStateSyncOn,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
		}
		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateNoPartnerChurn {
			t.Error("Error Cdm State Machine transition did not happen currstate is", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		}
		for _, evt := range invalidStateMap {
			// force state to no actor churn
			p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            evt,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
			if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateNoPartnerChurn {
				t.Error("Error Cdm State Machine transitioned on a invalid event")
			}
		}
	}
}

func TestCdmActorChurnInvalidEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	invalidStateMap := [2]fsm.Event{
		LacpCdmEventActorChurnTimerExpired,
		LacpCdmEventActorOperPortStateSyncOff,
	}

	var p1 *LaAggPort
	responseChannel := make(chan string)
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStateActorChurnMonitor], "found ", CdmStateStrMap[p1.CdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            LacpCdmEventActorChurnTimerExpired,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
		}
		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurn {
			t.Error("Error Cdm State Machine did not transition to actor churn")
		}

		for _, evt := range invalidStateMap {
			// force state to no actor churn
			p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            evt,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
			if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurn {
				t.Error("Error Cdm State Machine transitioned on a invalid event")
			}
		}
	}
}

func TestCdmPartnerChurnInvalidEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	invalidStateMap := [2]fsm.Event{
		LacpCdmEventPartnerChurnTimerExpired,
		LacpCdmEventPartnerOperPortStateSyncOff,
	}

	var p1 *LaAggPort
	responseChannel := make(chan string)
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStatePartnerChurnMonitor], "found ", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            LacpCdmEventActorChurnTimerExpired,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
		}

		for _, evt := range invalidStateMap {
			// force state to no actor churn
			p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:            evt,
				Src:          "TEST",
				ResponseChan: responseChannel,
			}
			<-responseChannel
			if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurn {
				t.Error("Error Cdm State Machine transitioned on a invalid event")
			}
		}
	}
}

func TestCdmActorChurnDetectionExpireEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStateActorChurnMonitor], "found ", CdmStateStrMap[p1.CdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			p1.CdMachineFsm.ChurnDetectionTimerIntervalSet(time.Millisecond * 10)
			p1.CdMachineFsm.ChurnDetectionTimerStart()
			time.Sleep(time.Millisecond * 11)
		}

		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurn {
			t.Error("Error Cdm State Machine did not transition on timer expired event")
		}
		if p1.actorChurn != true {
			t.Error("Error actorChurn not set")
		}
	}
}

func TestCdmPartnerChurnDetectionExpireEvents(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStatePartnerChurnMonitor], "found ", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			p1.PCdMachineFsm.ChurnDetectionTimerIntervalSet(time.Millisecond * 10)
			p1.PCdMachineFsm.ChurnDetectionTimerStart()
			time.Sleep(time.Millisecond * 11)
		}

		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurn {
			t.Error("Error Cdm State Machine did not transition on timer expired event")
		}

		if p1.partnerChurn != true {
			t.Error("Error partnerChurn not set")
		}
	}
}

func TestCdmActorChurnDebugCountDoesNotIncrementWhenStateReachedMoreThanFiveTimesInASecond(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {
		responseChannel := make(chan string)
		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStateActorChurnMonitor], "found ", CdmStateStrMap[p1.CdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			currentChurnActorTimestamp := p1.CdMachineFsm.churnCountTimestamp
			for i := 0; i < 7; i++ {
				// Sync set Actor Churn
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorChurnTimerExpired,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				// Actor Churn -> No actor churn
				LacpStateSet(&p1.ActorOper.State, LacpStateSyncBit)
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorOperPortStateSyncOn,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				// No Actory Churn -> Actor Churn Monitor
				LacpStateClear(&p1.ActorOper.State, LacpStateSyncBit)
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorOperPortStateSyncOff,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel
			}

			if p1.CdMachineFsm.churnCountTimestamp.Nanosecond()-currentChurnActorTimestamp.Nanosecond() < 1000000000 &&
				p1.AggPortDebug.AggPortDebugActorChurnCount > 5 {
				t.Error("Error Churn Count incremented more than 5 times in less than a second")
			}
		}
	}
}

func TestCdmPartnerChurnDebugCountDoesNotIncrementWhenStateReachedMoreThanFiveTimesInASecond(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		responseChannel := make(chan string)
		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStatePartnerChurnMonitor], "found ", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			currentChurnPartnerTimestamp := p1.PCdMachineFsm.churnCountTimestamp
			for i := 0; i < 7; i++ {
				// Sync set Actor Churn
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerChurnTimerExpired,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurn {
					t.Error("ERROR Did not transition to Partner Churn State")
				}

				// Actor Churn -> No actor churn
				LacpStateSet(&p1.PartnerOper.State, LacpStateSyncBit)
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerOperPortStateSyncOn,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateNoPartnerChurn {
					t.Error("ERROR Did not transition to No Partner Churn State")
				}

				// No Actory Churn -> Actor Churn Monitor
				LacpStateClear(&p1.PartnerOper.State, LacpStateSyncBit)
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerOperPortStateSyncOff,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel
				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
					t.Error("ERROR Did not transition to Partner Churn Monitor State")
				}

			}

			if p1.PCdMachineFsm.churnCountTimestamp.Nanosecond()-currentChurnPartnerTimestamp.Nanosecond() < 1000000000 &&
				p1.AggPortDebug.AggPortDebugPartnerChurnCount > 5 {
				t.Error("Error Churn Count incremented more than 5 times in less than a second")
			}
		}
	}
}

func TestCdmActorChurnDebugCountIncrementsWhenStateReachMoreThanOneSecond(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {
		responseChannel := make(chan string)
		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStateActorChurnMonitor], "found ", CdmStateStrMap[p1.CdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			for i := 0; i < 3; i++ {
				// Sync set Actor Churn
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorChurnTimerExpired,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				// Actor Churn -> No actor churn
				LacpStateSet(&p1.ActorOper.State, LacpStateSyncBit)
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorOperPortStateSyncOn,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				// No Actory Churn -> Actor Churn Monitor
				LacpStateClear(&p1.ActorOper.State, LacpStateSyncBit)
				p1.CdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventActorOperPortStateSyncOff,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel
				time.Sleep(time.Second * 2)

			}

			if p1.AggPortDebug.AggPortDebugActorChurnCount != 3 {
				t.Error("Error Churn Count did not increment every time Actor Churn State Reached")
			}
		}
	}
}

func TestCdmPartnerChurnDebugCountIncrementsWhenStateReachMoreThanOneSecond(t *testing.T) {
	defer MemoryCheck(t)
	ChurnDetectionStateMachineSetup()
	defer ChurnDetectionStateMachineTeardown()

	var p1 *LaAggPort
	if LaFindPortById(LaAggChurnPortActor, &p1) {

		responseChannel := make(chan string)
		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("Error State machine is not initalized to the proper state of ",
				CdmStateStrMap[LacpCdmStatePartnerChurnMonitor], "found ", CdmStateStrMap[p1.PCdMachineFsm.Machine.Curr.CurrentState()])
		} else {
			for i := 0; i < 3; i++ {
				// Sync set Actor Churn
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerChurnTimerExpired,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurn {
					t.Error("ERROR Did not transition to Partner Churn State")
				}

				// Actor Churn -> No actor churn
				LacpStateSet(&p1.PartnerOper.State, LacpStateSyncBit)
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerOperPortStateSyncOn,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel

				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateNoPartnerChurn {
					t.Error("ERROR Did not transition to No Partner Churn State")
				}

				// No Actory Churn -> Actor Churn Monitor
				LacpStateClear(&p1.PartnerOper.State, LacpStateSyncBit)
				p1.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:            LacpCdmEventPartnerOperPortStateSyncOff,
					Src:          "TEST",
					ResponseChan: responseChannel,
				}
				<-responseChannel
				if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
					t.Error("ERROR Did not transition to Partner Churn Monitor State")
				}
				time.Sleep(time.Second * 2)
			}

			if p1.AggPortDebug.AggPortDebugPartnerChurnCount != 3 {
				t.Error("Error Churn Count did not increment every time Partner Churn State reached")
			}
		}
	}
}
