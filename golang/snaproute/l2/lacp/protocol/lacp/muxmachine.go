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

// MUX MACHINE 802.1ax-2014 Section 6.4.15
// This implementation will assume that bot State machines in Section 6.4.15 are
// implemented with an extra flag indicating the capabilities of the port
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"sort"
	"strconv"
	"strings"
	"time"
	"utils/fsm"
)

const MuxMachineModuleStr = "Mux Machine"

const (
	LacpMuxmStateNone = iota
	LacpMuxmStateDetached
	LacpMuxmStateWaiting
	LacpMuxmStateAttached
	LacpMuxmStateCollecting
	LacpMuxmStateDistributing
	// Coupled control - Collecting and Distributing can't be controlled independently
	LacpMuxmStateCNone
	LacpMuxmStateCDetached
	LacpMuxmStateCWaiting
	LacpMuxmStateCAttached
	LacpMuxStateCCollectingDistributing
)

var MuxmStateStrMap map[fsm.State]string
var MuxmEventStrMap map[int]string

func MuxxMachineStrStateMapCreate() {

	MuxmStateStrMap = make(map[fsm.State]string)
	MuxmStateStrMap[LacpMuxmStateNone] = "None"
	MuxmStateStrMap[LacpMuxmStateDetached] = "Detached"
	MuxmStateStrMap[LacpMuxmStateWaiting] = "Waiting"
	MuxmStateStrMap[LacpMuxmStateAttached] = "Attached"
	MuxmStateStrMap[LacpMuxmStateCollecting] = "Collecting"
	MuxmStateStrMap[LacpMuxmStateDistributing] = "Distributing"
	MuxmStateStrMap[LacpMuxmStateCNone] = "LacpMuxmStateCNone"
	MuxmStateStrMap[LacpMuxmStateCDetached] = "CDetached"
	MuxmStateStrMap[LacpMuxmStateCWaiting] = "CWaiting"
	MuxmStateStrMap[LacpMuxmStateCAttached] = "CAttached"
	MuxmStateStrMap[LacpMuxStateCCollectingDistributing] = "CCollectingDistributing"

	MuxmEventStrMap = make(map[int]string)
	MuxmEventStrMap[LacpMuxmEventBegin] = "Event Begin"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualSelected] = "Event AggSelected equals Selected"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualStandby] = "Event AggSlected equals Standby"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualUnselected] = "Event AggSelected equals Unselected"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualSelectedAndReady] = "Event AggSelected equals Selected and Agg is Ready"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualSelectedAndPartnerSync] = "Event Selected equals Selected and Partner Oper Sync state set"
	MuxmEventStrMap[LacpMuxmEventNotPartnerSync] = "Event Partner Oper Sync state is NOT set"
	MuxmEventStrMap[LacpMuxmEventNotPartnerCollecting] = "Event Partner Oper Collecting state is not set"
	MuxmEventStrMap[LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting] = "Event Selected equals Selected and Partner Oper Sync and Collecting state is set"

}

const (
	LacpMuxmEventBegin = iota + 1
	LacpMuxmEventSelectedEqualSelected
	LacpMuxmEventSelectedEqualStandby
	LacpMuxmEventSelectedEqualUnselected
	LacpMuxmEventSelectedEqualSelectedAndReady
	LacpMuxmEventSelectedEqualSelectedAndPartnerSync
	LacpMuxmEventNotPartnerSync
	LacpMuxmEventNotPartnerCollecting
	LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type LacpMuxMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	p *LaAggPort

	// debug log
	log    chan string
	logEna bool

	collDistCoupled bool

	// timer interval
	waitWhileTimerTimeout time.Duration
	waitWhileTimerRunning bool

	// timers
	waitWhileTimer *time.Timer

	// machine specific events
	MuxmEvents         chan utils.MachineEvent
	MuxmLogEnableEvent chan bool

	actorSyncTransitionTimestamp time.Time
}

func (muxm *LacpMuxMachine) Stop() {

	close(muxm.MuxmEvents)
	close(muxm.MuxmLogEnableEvent)
}

func (muxm *LacpMuxMachine) PrevState() fsm.State { return muxm.PreviousState }

// PrevStateSet will set the previous State
func (muxm *LacpMuxMachine) PrevStateSet(s fsm.State) { muxm.PreviousState = s }

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewLacpMuxMachine(port *LaAggPort) *LacpMuxMachine {
	muxm := &LacpMuxMachine{
		p:                     port,
		collDistCoupled:       false,
		waitWhileTimerTimeout: LacpAggregateWaitTime,
		PreviousState:         LacpMuxmStateNone,
		MuxmEvents:            make(chan utils.MachineEvent, 10),
		MuxmLogEnableEvent:    make(chan bool)}

	port.MuxMachineFsm = muxm

	// start then stop
	muxm.WaitWhileTimerStart()
	muxm.WaitWhileTimerStop()

	return muxm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (muxm *LacpMuxMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if muxm.Machine == nil {
		muxm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	muxm.Machine.Rules = r
	muxm.Machine.Curr = &utils.StateEvent{
		StrStateMap: MuxmStateStrMap,
		LogEna:      muxm.p.logEna,
		Logger:      muxm.LacpMuxmLog,
		Owner:       MuxMachineModuleStr,
	}

	return muxm.Machine
}

func (muxm *LacpMuxMachine) SendTxMachineNtt() {

	if muxm.p.TxMachineFsm.Machine.Curr.CurrentState() != LacpTxmStateOff {
		muxm.p.TxMachineFsm.TxmEvents <- utils.MachineEvent{
			E:   LacpTxmEventNtt,
			Src: MuxMachineModuleStr}
	}
}

// LacpMuxmDetached
func (muxm *LacpMuxMachine) LacpMuxmDetached(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// DETACH MUX FROM AGGREGATOR
	muxm.DetachMuxFromAggregator()

	// Actor Oper State Sync = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateSyncBit)
	// inform cdm
	if p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateNoActorChurn {
		p.CdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventActorOperPortStateSyncOff,
			Src: MuxMachineModuleStr}
	}

	// Disable Distributing
	muxm.DisableDistributing()

	// Actor Oper State Distributing = FALSE
	// Actor Oper State Collecting = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateDistributingBit|LacpStateCollectingBit)

	// Disable Collecting
	muxm.DisableCollecting()

	// NTT = TRUE
	// TODO: is this necessary? May only want to let TxMachine
	//       set ntt to true based on NTT event
	//p.TxMachineFsm.ntt = true

	// indicate that NTT = TRUE
	if muxm.Machine.Curr.CurrentState() != LacpMuxmStateNone {
		muxm.SendTxMachineNtt()
	}

	return LacpMuxmStateDetached
}

// LacpMuxmWaiting
func (muxm *LacpMuxMachine) LacpMuxmWaiting(m fsm.Machine, data interface{}) fsm.State {
	var a *LaAggregator
	//var State fsm.State
	p := muxm.p

	skipWaitWhileTimer := false

	// only need to kick off the timer if ready is not true
	// ready will be true if all other ports are attached
	// or this is the the first
	// or lacp is not enabled
	if LaFindAggById(p.AggId, &a) {
		if a.ready || LacpModeGet(p.ActorAdmin.State, p.lacpEnabled) == LacpModeOn {
			skipWaitWhileTimer = true
			a.ready = false
		}
	}

	//State = LacpMuxmStateWaiting
	if !skipWaitWhileTimer {
		muxm.WaitWhileTimerStart()
	} else {
		muxm.LacpMuxmLog("Force Stopping Wait While Timer")
		muxm.WaitWhileTimerStop()
		// force the the next State to attach
		//muxm.LacpMuxmWaitingEvaluateSelected(true)
		//muxm.Machine.Curr.CurrentState()
	}

	return LacpMuxmStateWaiting
}

// LacpMuxmAttached
func (muxm *LacpMuxMachine) LacpMuxmAttached(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// NTT = TRUE
	defer muxm.SendTxMachineNtt()
	// send event to user port and partner info don't
	defer utils.ProcessLacpPortPartnerInfoSync(int32(p.PortNum))

	// Attach Mux to Aggregator
	muxm.AttachMuxToAggregator()

	// Actor Oper State Sync = TRUE
	//muxm.LacpMuxmLog("Setting Actor Sync Bit")
	LacpStateSet(&p.ActorOper.State, LacpStateSyncBit)

	// inform cdm
	if p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateActorChurnMonitor ||
		p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateActorChurn {
		p.CdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventActorOperPortStateSyncOn,
			Src: MuxMachineModuleStr}
	}

	// debug
	if p.AggPortDebug.AggPortDebugActorSyncTransitionCount == 0 {
		muxm.actorSyncTransitionTimestamp = time.Now()
		p.AggPortDebug.AggPortDebugActorSyncTransitionCount++
	} else if time.Now().Second()-muxm.actorSyncTransitionTimestamp.Second() > 5 {
		p.AggPortDebug.AggPortDebugActorSyncTransitionCount++
		muxm.actorSyncTransitionTimestamp = time.Now()
	}

	// Actor Oper State Collecting = FALSE
	//muxm.LacpMuxmLog("Clearing Actor Collecting Bit")
	LacpStateClear(&p.ActorOper.State, LacpStateCollectingBit)

	// Disable Collecting
	muxm.DisableCollecting()

	return LacpMuxmStateAttached
}

// LacpMuxmCollecting
func (muxm *LacpMuxMachine) LacpMuxmCollecting(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// Enabled Collecting
	muxm.EnableCollecting()

	// Actor Oper State Sync == TRUE
	//muxm.LacpMuxmLog("Setting Actor Collecting Bit")
	LacpStateSet(&p.ActorOper.State, LacpStateCollectingBit)

	// Disable Distributing
	muxm.DisableDistributing()

	// Actor Oper State Distributing = FALSE
	//muxm.LacpMuxmLog("Clearing Actor Distributing Bit")
	LacpStateClear(&p.ActorOper.State, LacpStateDistributingBit)

	if p.AggAttached != nil &&
		len(p.AggAttached.DistributedPortNumList) == 0 {
		p.AggAttached.OperState = false
	}

	// indicate that NTT = TRUE
	defer muxm.SendTxMachineNtt()

	return LacpMuxmStateCollecting
}

// LacpMuxmDistributing
func (muxm *LacpMuxMachine) LacpMuxmDistributing(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// Actor Oper State Sync == TRUE
	//muxm.LacpMuxmLog("Setting Actor Distributing Bit")
	LacpStateSet(&p.ActorOper.State, LacpStateDistributingBit)

	// Enabled Distributing
	muxm.EnableDistributing()
	if p.AggAttached != nil {
		p.AggAttached.OperState = true
	}

	// indicate that NTT = TRUE
	defer muxm.SendTxMachineNtt()

	return LacpMuxmStateDistributing
}

// LacpMuxmCDetached
func (muxm *LacpMuxMachine) LacpMuxmCDetached(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// indicate that NTT = TRUE
	defer muxm.SendTxMachineNtt()
	// send event to user port and partner 	info don't
	defer utils.ProcessLacpPortPartnerInfoMismatch(int32(p.PortNum))

	// DETACH MUX FROM AGGREGATOR
	muxm.DetachMuxFromAggregator()

	// Actor Oper State Sync = FALSE
	// Actor Oper State Collecting = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateSyncBit|LacpStateCollectingBit)

	if p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateNoActorChurn {
		// inform cdm
		p.CdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventActorOperPortStateSyncOff,
			Src: MuxMachineModuleStr}
	}

	// Disable Collecting && Distributing
	muxm.DisableCollectingDistributing()

	// Actor Oper State Distributing = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateDistributingBit)

	return LacpMuxmStateDetached
}

// LacpMuxmCWaiting
func (muxm *LacpMuxMachine) LacpMuxmCWaiting(m fsm.Machine, data interface{}) fsm.State {
	//p := muxm.p

	muxm.WaitWhileTimerStart()

	return LacpMuxmStateWaiting
}

// LacpMuxmAttached
func (muxm *LacpMuxMachine) LacpMuxmCAttached(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// NTT = TRUE
	defer muxm.SendTxMachineNtt()
	// send event to user port and partner info don't
	defer utils.ProcessLacpPortPartnerInfoSync(int32(p.PortNum))

	// Attach Mux to Aggregator
	muxm.AttachMuxToAggregator()

	// Actor Oper State Sync = TRUE
	LacpStateSet(&p.ActorOper.State, LacpStateSyncBit)

	// inform cdm
	if p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateActorChurnMonitor ||
		p.CdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateActorChurn {
		p.CdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventActorOperPortStateSyncOn,
			Src: MuxMachineModuleStr}
	}

	// Actor Oper State Collecting = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateCollectingBit)

	// Disable Collecting && Distributing
	muxm.DisableCollectingDistributing()

	// Actor Oper State Distributing = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateDistributingBit)

	return LacpMuxmStateWaiting
}

// LacpMuxmCollecting
func (muxm *LacpMuxMachine) LacpMuxmCCollectingDistributing(m fsm.Machine, data interface{}) fsm.State {
	p := muxm.p

	// Actor Oper State Distributing = TRUE
	LacpStateSet(&p.ActorOper.State, LacpStateDistributingBit)

	// Enable Collecting && Distributing
	muxm.EnableCollectingDistributing()

	// Actor Oper State Distributing == FALSE
	LacpStateSet(&p.ActorOper.State, LacpStateDistributingBit)

	// indicate that NTT = TRUE
	defer muxm.SendTxMachineNtt()

	return LacpMuxmStateWaiting
}

// LacpMuxMachineFSMBuild:  802.1ax-2014 Figure 6-21 && 6-22
func (p *LaAggPort) LacpMuxMachineFSMBuild() *LacpMuxMachine {

	rules := fsm.Ruleset{}

	MuxxMachineStrStateMapCreate()

	// Instantiate a new LacpRxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	muxm := NewLacpMuxMachine(p)

	// MUX
	//BEGIN -> DETACHED
	rules.AddRule(LacpMuxmStateNone, LacpMuxmEventBegin, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateDetached, LacpMuxmEventBegin, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateWaiting, LacpMuxmEventBegin, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateAttached, LacpMuxmEventBegin, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventBegin, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateDistributing, LacpMuxmEventBegin, muxm.LacpMuxmDetached)

	// SELECTED or STANDBY -> WAITING
	rules.AddRule(LacpMuxmStateDetached, LacpMuxmEventSelectedEqualSelected, muxm.LacpMuxmWaiting)
	rules.AddRule(LacpMuxmStateDetached, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmWaiting)
	// UNSELECTED -> DETACHED
	rules.AddRule(LacpMuxmStateWaiting, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmDetached)
	// SELECTED && READY -> ATTACHED
	rules.AddRule(LacpMuxmStateWaiting, LacpMuxmEventSelectedEqualSelectedAndReady, muxm.LacpMuxmAttached)
	// UNSELECTED or STANDBY -> DETACHED
	rules.AddRule(LacpMuxmStateAttached, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmDetached)
	rules.AddRule(LacpMuxmStateAttached, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmDetached)
	// SELECTED && PARTNER SYNC -> COLLECTING
	rules.AddRule(LacpMuxmStateAttached, LacpMuxmEventSelectedEqualSelectedAndPartnerSync, muxm.LacpMuxmCollecting)
	// UNSELECTED or STANDBY or NOT PARTNER SYNC -> ATTACHED
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmAttached)
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmAttached)
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventNotPartnerSync, muxm.LacpMuxmAttached)
	// SELECTED && PARTNER SYNC && PARTNER COLLECTING -> DISTRIBUTING
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting, muxm.LacpMuxmDistributing)
	// UNSELECTED or STANDBY or NOT PARTNER SYNC or NOT PARTNER COLLECTING -> COLLECTING
	rules.AddRule(LacpMuxmStateDistributing, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmCollecting)
	rules.AddRule(LacpMuxmStateDistributing, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmCollecting)
	rules.AddRule(LacpMuxmStateDistributing, LacpMuxmEventNotPartnerSync, muxm.LacpMuxmCollecting)
	rules.AddRule(LacpMuxmStateDistributing, LacpMuxmEventNotPartnerCollecting, muxm.LacpMuxmCollecting)

	// MUX Coupled
	//BEGIN -> DETACHED
	rules.AddRule(LacpMuxmStateNone, LacpMuxmEventBegin, muxm.LacpMuxmCDetached)
	rules.AddRule(LacpMuxmStateCDetached, LacpMuxmEventBegin, muxm.LacpMuxmCDetached)
	rules.AddRule(LacpMuxmStateCWaiting, LacpMuxmEventBegin, muxm.LacpMuxmCDetached)
	rules.AddRule(LacpMuxmStateCAttached, LacpMuxmEventBegin, muxm.LacpMuxmCDetached)
	rules.AddRule(LacpMuxStateCCollectingDistributing, LacpMuxmEventBegin, muxm.LacpMuxmCDetached)

	// SELECTED or STANDBY -> WAITING
	rules.AddRule(LacpMuxmStateCDetached, LacpMuxmEventSelectedEqualSelected, muxm.LacpMuxmCWaiting)
	rules.AddRule(LacpMuxmStateCDetached, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmCWaiting)
	// UNSELECTED -> DETACHED
	rules.AddRule(LacpMuxmStateCWaiting, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmCDetached)
	// SELECTED && READY -> ATTACHED
	rules.AddRule(LacpMuxmStateCWaiting, LacpMuxmEventSelectedEqualSelectedAndReady, muxm.LacpMuxmAttached)
	// UNSELECTED or STANDBY -> DETACHED
	rules.AddRule(LacpMuxmStateCAttached, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmCDetached)
	rules.AddRule(LacpMuxmStateCAttached, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmCDetached)
	// SELECTED && PARTNER SYNC -> COLLECTING-DISTRIBUTING
	rules.AddRule(LacpMuxmStateCAttached, LacpMuxmEventSelectedEqualSelectedAndPartnerSync, muxm.LacpMuxmCCollectingDistributing)
	// UNSELECTED or STANDBY or NOT PARTNER SYNC -> ATTACHED
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventSelectedEqualUnselected, muxm.LacpMuxmCAttached)
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventSelectedEqualStandby, muxm.LacpMuxmCAttached)
	rules.AddRule(LacpMuxmStateCollecting, LacpMuxmEventNotPartnerSync, muxm.LacpMuxmCAttached)

	// Create a new FSM and apply the rules
	muxm.Apply(&rules)

	return muxm
}

// LacpMuxMachineMain:  802.1ax-2014 Figure 6-21 && 6-22
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpMuxMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.13 Periodic Transmission Machine
	muxm := p.LacpMuxMachineFSMBuild()
	p.wg.Add(1)

	// TODO: Hw only supports mux coupling, this should be a param file for lacp
	//if LacpSysGlobalInfoGet(LacpSystem{Actor_System: p.AggAttached.Config.SystemIdMac,
	//	Actor_System_priority: p.AggAttached.Config.SystemPriority}).muxCoupling {
	//	muxm.PrevStateSet(LacpMuxmStateCNone)
	//}
	// set the inital State
	muxm.Machine.Start(muxm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpMuxMachine) {
		m.LacpMuxmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			// save the current machine state
			p.AggPortDebug.AggPortDebugMuxState = int(m.Machine.Curr.CurrentState())
			select {

			case <-m.waitWhileTimer.C:
				m.LacpMuxmLog("MUXM: Wait While Timer Expired")
				// lets evaluate selection
				if m.Machine.Curr.CurrentState() == LacpMuxmStateWaiting ||
					m.Machine.Curr.CurrentState() == LacpMuxmStateCWaiting {
					m.LacpMuxmWaitingEvaluateSelected(false)
				}

			case event, ok := <-m.MuxmEvents:

				if ok {
					p := m.p
					//m.LacpMuxmLog(fmt.Sprintf("Event received %d src %s", event.E, event.Src))
					eventStr := strings.Join([]string{"from", event.Src, MuxmEventStrMap[int(event.E)]}, " ")

					// process the event
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)

					if rv != nil {
						m.LacpMuxmLog(strings.Join([]string{error.Error(rv), event.Src, MuxmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					} else {

						// continuation events
						if m.Machine.Curr.CurrentState() == LacpMuxmStateDetached ||
							m.Machine.Curr.CurrentState() == LacpMuxmStateCDetached {
							// if port is attached then we know that provisioning found
							// a valid agg thus port should be attached.
							if p.AggAttached != nil &&
								p.IsPortEnabled() &&
								p.lacpEnabled {
								// change the selection to be Selected
								p.aggSelected = LacpAggSelected
								//muxm.LacpMuxmLog("Setting Actor Aggregation Bit")
								LacpStateSet(&p.ActorOper.State, LacpStateAggregationBit)

								eventStr = strings.Join([]string{eventStr,
									"and\nfrom", MuxMachineModuleStr, MuxmEventStrMap[LacpMuxmEventSelectedEqualSelected]}, " ")

								m.Machine.ProcessEvent(MuxMachineModuleStr, LacpMuxmEventSelectedEqualSelected, nil)
								event.E = LacpMuxmEventSelectedEqualSelected
							}
						}
						if event.E == LacpMuxmEventSelectedEqualSelected &&
							(m.Machine.Curr.CurrentState() == LacpMuxmStateWaiting ||
								m.Machine.Curr.CurrentState() == LacpMuxmStateCWaiting) &&
							!m.waitWhileTimerRunning {
							// special case we may have a delayed event which will do a fast transition to next State
							// Attached, trigger is the fact that the timer is not running
							m.LacpMuxmWaitingEvaluateSelected(true)
						}
						if (m.Machine.Curr.CurrentState() == LacpMuxmStateAttached ||
							m.Machine.Curr.CurrentState() == LacpMuxmStateCAttached) &&
							p.aggSelected == LacpAggSelected &&
							LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {

							eventStr = strings.Join([]string{eventStr,
								"and\nfrom", MuxMachineModuleStr, MuxmEventStrMap[LacpMuxmEventSelectedEqualSelectedAndPartnerSync]}, " ")

							m.Machine.ProcessEvent(MuxMachineModuleStr, LacpMuxmEventSelectedEqualSelectedAndPartnerSync, nil)
						}
						if m.Machine.Curr.CurrentState() == LacpMuxmStateCollecting &&
							p.aggSelected == LacpAggSelected &&
							LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) &&
							LacpStateIsSet(p.PartnerOper.State, LacpStateCollectingBit) {

							eventStr = strings.Join([]string{eventStr,
								"and\nfrom", MuxMachineModuleStr, MuxmEventStrMap[LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting]}, " ")
							m.Machine.ProcessEvent(MuxMachineModuleStr, LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting, nil)
						}
						if event.E == LacpMuxmEventSelectedEqualUnselected &&
							(m.Machine.Curr.CurrentState() != LacpMuxmStateDetached &&
								m.Machine.Curr.CurrentState() != LacpMuxmStateCDetached) {
							// Unselected State will cause a downward transition to detached State
							State := m.Machine.Curr.CurrentState()
							endState := fsm.State(LacpMuxmStateDetached)
							if m.Machine.Curr.CurrentState() > LacpMuxmStateDistributing {
								endState = LacpMuxmStateCDetached
							}
							eventStr = strings.Join([]string{eventStr,
								"and\nfrom", MuxMachineModuleStr, MuxmEventStrMap[LacpMuxmEventSelectedEqualUnselected]}, " ")

							for ; State > endState; State-- {

								m.Machine.ProcessEvent(MuxMachineModuleStr, LacpMuxmEventSelectedEqualUnselected, nil)
							}
						}
					}

					if len(eventStr) > 255 {
						fmt.Println("WARNING string to long for MuxReason:", eventStr)
						fmt.Println(eventStr)
					}
					p.AggPortDebug.AggPortDebugMuxReason = eventStr

					if event.ResponseChan != nil {
						//m.LacpMuxmLog("Sending response")
						utils.SendResponse(MuxMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.LacpMuxmLog("Machine End")
					return
				}

			case ena := <-m.MuxmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(muxm)
}

// LacpMuxmEvaluateSelected 802.1ax-2014 Section 6.4.15
// d) If Selected is SELECTED, the wait_while_timer forces a delay to allow
// for the possibility that other Aggregation Ports may be reconfiguring
// at the same time. Once the wait_while_timer expires, and once the wait_
// while_timers of all other Aggregation Ports that are ready to attach to
// the same Aggregator have expired, the process of attaching the Aggregation
// Port to the Aggregator can proceed, and the State machine enters the
// ATTACHED State. During the waiting time, changes in selection parameters
// can occur that will result in a re-evaluation of Selected. If Selected
// becomes UNSELECTED, then the State machine reenters the DETACHED State.
// If Selected becomes STANDBY, the operation is as described in item e).
//
// NOTE—This waiting period reduces the disturbance that will be visible
// to higher layers; for example, on start-up events. However, the selection
// need not wait for the entire waiting period in cases where it is known that
// no other Aggregation Ports will attach; for example, where all other
// Aggregation Ports with the same operational Key are already attached to the
// Aggregator.
//
// e) If Selected is STANDBY, the Aggregation Port is held in the WAITING
// State until such a time as the selection parameters change, resulting in a
// re-evaluation of the Selected variable. If Selected becomes UNSELECTED,
// the State machine reenters the DETACHED State. If SELECTED becomes SELECTED,
// then the operation is as described in item d). The latter case allows an
// Aggregation Port to be brought into operation from STANDBY with minimum
// delay once Selected becomes SELECTED.
func (muxm *LacpMuxMachine) LacpMuxmWaitingEvaluateSelected(sendResponse bool) {
	var a *LaAggregator
	p := muxm.p
	muxm.LacpMuxmLog(strings.Join([]string{"Selected", strconv.Itoa(LacpAggSelected), "actual", strconv.Itoa(p.aggSelected)}, "="))
	// current port should be in selected State
	if p.aggSelected == LacpAggSelected ||
		p.aggSelected == LacpAggStandby {
		p.readyN = true
		if LaFindAggById(p.AggId, &a) {
			a.LacpMuxCheckSelectionLogic(p, sendResponse)
		} else {
			muxm.LacpMuxmLog(fmt.Sprintf("Unable to find Aggrigator %d", p.AggId))
		}
	}
}

// AttachMuxToAggregator is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregation Port’s Control Parser/Multiplexer
// to be attached to the Aggregator Parser/Multiplexer of the selected
// Aggregator, in preparation for collecting and distributing frames.
func (muxm *LacpMuxMachine) AttachMuxToAggregator() {
	// TODO send message to asic deamon  create
	p := muxm.p
	if LaFindAggById(p.AggId, &p.AggAttached) {
		LacpStateSet(&p.ActorOper.State, LacpStateAggregationBit)
		muxm.LacpMuxmLog("Attach Mux To Aggregator Enter")
	}
}

// DetachMuxFromAggregator is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregation Port’s Control Parser/Multiplexer
// to be detached from the Aggregator Parser/Multiplexer of the Aggregator
// to which the Aggregation Port is currently attached.
func (muxm *LacpMuxMachine) DetachMuxFromAggregator() {
	// TODO send message to asic deamon delete
	muxm.LacpMuxmLog("Detach Mux From Aggregator Enter")
	//p := muxm.p
	//p.AggAttached = nil
	// should already be in unselected State
	//p.aggSelected = LacpAggUnSelected

	// Remove port from HW lag group
}

// EnableCollecting is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Parser of the Aggregator to which
// the Aggregation Port is attached to start collecting frames from the
// Aggregation Port.
func (muxm *LacpMuxMachine) EnableCollecting() {
	// TODO send message to asic deamon
	muxm.LacpMuxmLog("Sending Collection Enable to ASICD")
}

// DisableCollecting is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Parser of the Aggregator to which
// the Aggregation Port is attached to stop collecting frames from the
// Aggregation Port.
func (muxm *LacpMuxMachine) DisableCollecting() {
	// TODO send message to asic deamon
	p := muxm.p
	if LacpStateIsSet(p.ActorOper.State, LacpStateCollectingBit) {
		muxm.LacpMuxmLog("Sending Collection Disable to ASICD")
	}
}

// EnableDistributing is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Multiplexer of the Aggregator
// to which the Aggregation Port is attached to start distributing frames
// to the Aggregation Port.
func (muxm *LacpMuxMachine) EnableDistributing() {
	p := muxm.p
	a := muxm.p.AggAttached

	if a != nil {

		// asicd expects the port list to be a bitmap in string format

		a.DistributedPortNumList = append(a.DistributedPortNumList, p.IntfNum)
		sort.Strings(a.DistributedPortNumList)

		muxm.LacpMuxmLog(fmt.Sprintf("Agg %d hwAggId %d EnableDistributing PortsListLen %d PortList %v", p.AggId, a.HwAggId, len(a.DistributedPortNumList), a.DistributedPortNumList))
		for _, client := range utils.GetAsicDPluginList() {
			err := client.UpdateLag(a.HwAggId, asicDHashModeGet(a.LagHash), asicDPortBmpFormatGet(a.DistributedPortNumList))
			if err != nil {
				a.LacpAggLog(fmt.Sprintln("EnableDistributing: Error updating LAG in HW", err))
			}
		}

		// notify DR that port has been created
		for name, upcb := range LacpCbDb.PortUpDbList {
			a.LacpAggLog(fmt.Sprintf("Checking %s if it cares about port up for port %s", name, p.IntfNum))
			upcb(int32(p.PortNum))
		}

		if len(a.DistributedPortNumList) == 1 {
			for name, upcb := range LacpCbDb.AggOperUpDbList {
				a.LacpAggLog(fmt.Sprintf("Notify %s Agg OperState UP %s", name, a.AggName))
				upcb(int32(a.AggId))
			}
		}
	}
}

// DisableDistributing is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Multiplexer of the Aggregator
// to which the Aggregation Port is attached to stop distributing frames
// to the Aggregation Port.
func (muxm *LacpMuxMachine) DisableDistributing() {
	var portFound bool
	p := muxm.p
	a := muxm.p.AggAttached

	if a != nil {

		portFound = false
		for j := 0; j < len(a.DistributedPortNumList) && !portFound; j++ {
			if p.IntfNum == a.DistributedPortNumList[j] {
				portFound = true
				a.DistributedPortNumList = append(a.DistributedPortNumList[:j], a.DistributedPortNumList[j+1:]...)
			}
		}
		// only send info to hw if port is in distributed list
		if portFound {
			sort.Strings(a.DistributedPortNumList)

			muxm.LacpMuxmLog(fmt.Sprintf("Agg %d HwId %d DisableDistributing PortsListLen %d PortList %v", p.AggId, a.HwAggId, len(a.DistributedPortNumList), a.DistributedPortNumList))

			for _, client := range utils.GetAsicDPluginList() {
				err := client.UpdateLag(a.HwAggId, asicDHashModeGet(a.LagHash), asicDPortBmpFormatGet(a.DistributedPortNumList))
				if err != nil {
					muxm.LacpMuxmLog(fmt.Sprintln("ERROR Updating Lag in HW", err))
					return
				}
			}

			if len(a.DistributedPortNumList) == 0 {
				for name, upcb := range LacpCbDb.AggOperUpDbList {
					a.LacpAggLog(fmt.Sprintf("Notify %s Agg OperState DOWN %s\n", name, a.AggName))
					upcb(int32(a.AggId))
				}
			}

			// notify DR that port has been created
			for _, downcb := range LacpCbDb.PortDownDbList {
				downcb(int32(p.PortNum))
			}
		}
	}
}

// EnableCollectingDistributing is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Parser of the Aggregator to which
// the Aggregation Port is attached to start collecting frames from the
// Aggregation Port, and the Aggregator Multiplexer to start distributing
// frames to the Aggregation Port.
func (muxm *LacpMuxMachine) EnableCollectingDistributing() {
	// TODO send message to asic deamon
	muxm.LacpMuxmLog("Sending Collection-Distributing Enable to ASICD")
}

// DisableCollectingDistributing is a required function defined in 802.1ax-2014
// Section 6.4.9
// This function causes the Aggregator Parser of the Aggregator to which the
// Aggregation Port is attached to stop collecting frames from the Aggregation
// Port, and the Aggregator Multiplexer to stop distributing frames to the
// Aggregation Port.
func (muxm *LacpMuxMachine) DisableCollectingDistributing() {
	// TODO send message to asic deamon
	muxm.LacpMuxmLog("Sending Collection-Distributing Disable to ASICD")
}
