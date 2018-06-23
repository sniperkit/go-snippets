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
	"strconv"
	"strings"
	"time"
	"utils/fsm"
)

const CdMachineModuleStr = "Actor Churn Detection Machine"
const PCdMachineModuleStr = "Partner Churn Detection Machine"

const (
	LacpCdmStateNone = iota + 1
	LacpCdmStateNoActorChurn
	LacpCdmStateActorChurnMonitor
	LacpCdmStateActorChurn
	LacpCdmStateNoPartnerChurn
	LacpCdmStatePartnerChurnMonitor
	LacpCdmStatePartnerChurn
)

var CdmStateStrMap map[fsm.State]string

func CdMachineStrStateMapCreate() {
	if CdmStateStrMap == nil {
		CdmStateStrMap = make(map[fsm.State]string)
		CdmStateStrMap[LacpCdmStateNone] = "None"
		CdmStateStrMap[LacpCdmStateNoActorChurn] = "NoActorChurn"
		CdmStateStrMap[LacpCdmStateActorChurnMonitor] = "ActorChurnMonitor"
		CdmStateStrMap[LacpCdmStateActorChurn] = "ActorChurn"
		CdmStateStrMap[LacpCdmStateNoPartnerChurn] = "NoPartnerChurn"
		CdmStateStrMap[LacpCdmStatePartnerChurnMonitor] = "PartnerChurnMonitor"
		CdmStateStrMap[LacpCdmStatePartnerChurn] = "PartnerChurn"
	}
}

const (
	LacpCdmEventBegin = iota + 1
	LacpCdmEventNotPortEnabled
	LacpCdmEventActorOperPortStateSyncOn
	LacpCdmEventActorOperPortStateSyncOff
	LacpCdmEventActorChurnTimerExpired
	LacpCdmEventPartnerOperPortStateSyncOn
	LacpCdmEventPartnerOperPortStateSyncOff
	LacpCdmEventPartnerChurnTimerExpired
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type LacpCdMachine struct {
	// for debugging
	PreviousState fsm.State

	// actor
	Machine *fsm.Machine

	p *LaAggPort

	// timer intervals
	churnTimerInterval time.Duration

	// Interval timers
	churnTimer *time.Timer

	// machine specific events
	CdmEvents            chan utils.MachineEvent
	CdmLogEnableEvent    chan bool
	churnCountTimestamp  time.Time
	pchurnCountTimestamp time.Time
}

type LacpActorCdMachine struct {
	LacpCdMachine
}

type LacpPartnerCdMachine struct {
	LacpCdMachine
}

func (cdm *LacpCdMachine) Stop() {

	close(cdm.CdmEvents)
	close(cdm.CdmLogEnableEvent)
}

func (cdm *LacpCdMachine) PrevState() fsm.State { return cdm.PreviousState }

// PrevStateSet will set the previous State
func (cdm *LacpCdMachine) PrevStateSet(s fsm.State) { cdm.PreviousState = s }

// NewLacpActorCdMachine will create a new instance of the LacpRxMachine
func NewLacpActorCdMachine(port *LaAggPort) *LacpActorCdMachine {
	cdm := &LacpActorCdMachine{
		LacpCdMachine{p: port,
			PreviousState:      LacpCdmStateNone,
			churnTimerInterval: LacpChurnDetectionTime,
			CdmEvents:          make(chan utils.MachineEvent, 10),
			CdmLogEnableEvent:  make(chan bool)}}

	port.CdMachineFsm = cdm
	cdm.ChurnDetectionTimerStart()
	cdm.ChurnDetectionTimerStop()
	return cdm
}

// NewLacpActorCdMachine will create a new instance of the LacpRxMachine
func NewLacpPartnerCdMachine(port *LaAggPort) *LacpPartnerCdMachine {
	cdm := &LacpPartnerCdMachine{
		LacpCdMachine{p: port,
			PreviousState:      LacpCdmStateNone,
			churnTimerInterval: LacpChurnDetectionTime,
			CdmEvents:          make(chan utils.MachineEvent, 10),
			CdmLogEnableEvent:  make(chan bool)}}

	port.PCdMachineFsm = cdm
	cdm.ChurnDetectionTimerStart()
	cdm.ChurnDetectionTimerStop()
	return cdm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (cdm *LacpCdMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if cdm.Machine == nil {
		cdm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	cdm.Machine.Rules = r
	cdm.Machine.Curr = &utils.StateEvent{
		StrStateMap: CdmStateStrMap,
		LogEna:      true, //cdm.p.logEna,
		Logger:      cdm.LacpCdmLog,
		Owner:       CdMachineModuleStr,
	}

	return cdm.Machine
}

// LacpCdMachineNoActorChurn will set the churn State to false
func (cdm *LacpActorCdMachine) LacpCdMachineNoActorChurn(m fsm.Machine, data interface{}) fsm.State {
	p := cdm.p
	p.actorChurn = false
	cdm.ChurnDetectionTimerStop()
	return LacpCdmStateNoActorChurn
}

// LacpCdMachineActorChurn will set the churn State to true
func (cdm *LacpActorCdMachine) LacpCdMachineActorChurn(m fsm.Machine, data interface{}) fsm.State {
	const ONE_SECOND = 1
	p := cdm.p
	p.actorChurn = true
	if cdm.churnCountTimestamp.Nanosecond() == 0 {
		cdm.churnCountTimestamp = time.Now()
	}

	// 802.1ax 7.3.4.1.8 aAggPortDebugActorChurnCount
	// maximum 5 counts per second
	timeDiff := time.Now().Second() - cdm.churnCountTimestamp.Second()
	if timeDiff < ONE_SECOND &&
		(p.AggPortDebug.AggPortDebugActorChurnCount-p.AggPortDebug.AggPortDebugActorChurnPrevCnt) < 5 {
		p.AggPortDebug.AggPortDebugActorChurnCount++
	} else if timeDiff > ONE_SECOND {
		// reset the timestamp
		cdm.churnCountTimestamp = time.Now()
		p.AggPortDebug.AggPortDebugActorChurnPrevCnt = p.AggPortDebug.AggPortDebugActorChurnCount
		p.AggPortDebug.AggPortDebugActorChurnCount++
	}

	cdm.ChurnDetectionTimerStop()
	return LacpCdmStateActorChurn
}

// LacpCdMachineActorChurnMonitor will set the churn State to true
// and kick off the churn detection timer
func (cdm *LacpActorCdMachine) LacpCdMachineActorChurnMonitor(m fsm.Machine, data interface{}) fsm.State {
	p := cdm.p
	p.actorChurn = false
	cdm.ChurnDetectionTimerStart()
	return LacpCdmStateActorChurnMonitor
}

// LacpCdMachineNoActorChurn will set the churn State to false
func (cdm *LacpPartnerCdMachine) LacpCdMachineNoPartnerChurn(m fsm.Machine, data interface{}) fsm.State {
	p := cdm.p
	p.partnerChurn = false
	cdm.ChurnDetectionTimerStop()
	return LacpCdmStateNoPartnerChurn
}

// LacpCdMachineActorChurn will set the churn State to true
func (cdm *LacpPartnerCdMachine) LacpCdMachinePartnerChurn(m fsm.Machine, data interface{}) fsm.State {
	const ONE_SECOND = 1
	p := cdm.p
	p.partnerChurn = true
	if cdm.churnCountTimestamp.Nanosecond() == 0 {
		cdm.churnCountTimestamp = time.Now()
	}

	// 802.1ax 7.3.4.1.9 aAggPortDebugActorChurnCount
	// maximum 5 counts per second
	timeDiff := time.Now().Second() - cdm.churnCountTimestamp.Second()
	if timeDiff < ONE_SECOND &&
		(p.AggPortDebug.AggPortDebugPartnerChurnCount-p.AggPortDebug.AggPortDebugPartnerChurnPrevCount) < 5 {
		p.AggPortDebug.AggPortDebugPartnerChurnCount++
	} else if timeDiff >= ONE_SECOND {
		// reset the timestamp
		cdm.churnCountTimestamp = time.Now()
		p.AggPortDebug.AggPortDebugPartnerChurnPrevCount = p.AggPortDebug.AggPortDebugPartnerChurnCount
		p.AggPortDebug.AggPortDebugPartnerChurnCount++
	}

	cdm.ChurnDetectionTimerStop()
	return LacpCdmStatePartnerChurn
}

// LacpCdMachineActorChurnMonitor will set the churn State to true
// and kick off the churn detection timer
func (cdm *LacpPartnerCdMachine) LacpCdMachinePartnerChurnMonitor(m fsm.Machine, data interface{}) fsm.State {
	p := cdm.p
	p.partnerChurn = false
	cdm.ChurnDetectionTimerStart()
	return LacpCdmStatePartnerChurnMonitor
}

func LacpActorCdMachineFSMBuild(p *LaAggPort) *LacpActorCdMachine {

	rules := fsm.Ruleset{}

	CdMachineStrStateMapCreate()

	// Instantiate a new LacpRxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	cdm := NewLacpActorCdMachine(p)

	//BEGIN -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStateNone, LacpCdmEventBegin, cdm.LacpCdMachineActorChurnMonitor)
	rules.AddRule(LacpCdmStateActorChurn, LacpCdmEventBegin, cdm.LacpCdMachineActorChurnMonitor)
	rules.AddRule(LacpCdmStateNoActorChurn, LacpCdmEventBegin, cdm.LacpCdMachineActorChurnMonitor)
	rules.AddRule(LacpCdmStateActorChurnMonitor, LacpCdmEventBegin, cdm.LacpCdMachineActorChurnMonitor)

	// PORT ENABLED -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStateActorChurn, LacpCdmEventNotPortEnabled, cdm.LacpCdMachineActorChurnMonitor)
	rules.AddRule(LacpCdmStateNoActorChurn, LacpCdmEventNotPortEnabled, cdm.LacpCdMachineActorChurnMonitor)
	rules.AddRule(LacpCdmStateActorChurnMonitor, LacpCdmEventNotPortEnabled, cdm.LacpCdMachineActorChurnMonitor)

	// SYNC ON -> NO CHURN
	rules.AddRule(LacpCdmStateActorChurnMonitor, LacpCdmEventActorOperPortStateSyncOn, cdm.LacpCdMachineNoActorChurn)
	rules.AddRule(LacpCdmStateActorChurn, LacpCdmEventActorOperPortStateSyncOn, cdm.LacpCdMachineNoActorChurn)

	// SYNC OFF -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStateNoActorChurn, LacpCdmEventActorOperPortStateSyncOff, cdm.LacpCdMachineActorChurnMonitor)

	// TIMEOUT -> CHURN
	rules.AddRule(LacpCdmStateActorChurnMonitor, LacpCdmEventActorChurnTimerExpired, cdm.LacpCdMachineActorChurn)

	// Create a new FSM and apply the rules
	cdm.Apply(&rules)

	return cdm
}
func LacpPartnerCdMachineFSMBuild(p *LaAggPort) *LacpPartnerCdMachine {

	rules := fsm.Ruleset{}

	CdMachineStrStateMapCreate()

	// Instantiate a new LacpRxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	cdm := NewLacpPartnerCdMachine(p)

	//BEGIN -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStateNone, LacpCdmEventBegin, cdm.LacpCdMachinePartnerChurnMonitor)
	rules.AddRule(LacpCdmStatePartnerChurn, LacpCdmEventBegin, cdm.LacpCdMachinePartnerChurnMonitor)
	rules.AddRule(LacpCdmStateNoPartnerChurn, LacpCdmEventBegin, cdm.LacpCdMachinePartnerChurnMonitor)
	rules.AddRule(LacpCdmStatePartnerChurnMonitor, LacpCdmEventBegin, cdm.LacpCdMachinePartnerChurnMonitor)

	// PORT ENABLED -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStatePartnerChurn, LacpCdmEventNotPortEnabled, cdm.LacpCdMachinePartnerChurnMonitor)
	rules.AddRule(LacpCdmStateNoPartnerChurn, LacpCdmEventNotPortEnabled, cdm.LacpCdMachinePartnerChurnMonitor)
	rules.AddRule(LacpCdmStatePartnerChurnMonitor, LacpCdmEventNotPortEnabled, cdm.LacpCdMachinePartnerChurnMonitor)

	// SYNC ON -> NO CHURN
	rules.AddRule(LacpCdmStatePartnerChurnMonitor, LacpCdmEventPartnerOperPortStateSyncOn, cdm.LacpCdMachineNoPartnerChurn)
	rules.AddRule(LacpCdmStatePartnerChurn, LacpCdmEventPartnerOperPortStateSyncOn, cdm.LacpCdMachineNoPartnerChurn)

	// SYNC OFF -> ACTOR CHURN MONITOR
	rules.AddRule(LacpCdmStateNoPartnerChurn, LacpCdmEventPartnerOperPortStateSyncOff, cdm.LacpCdMachinePartnerChurnMonitor)

	// TIMEOUT -> CHURN
	rules.AddRule(LacpCdmStatePartnerChurnMonitor, LacpCdmEventPartnerChurnTimerExpired, cdm.LacpCdMachinePartnerChurn)

	// Create a new FSM and apply the rules
	cdm.Apply(&rules)

	return cdm
}

// LacpActorCdMachineMain:  802.1ax-2014
// Creation of Actor Churn Detection State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpActorCdMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.17 Churn Detection machine
	cdm := LacpActorCdMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	cdm.Machine.Start(cdm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpActorCdMachine) {
		m.LacpCdmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			m.p.AggPortDebug.AggPortDebugActorChurnState = int(m.Machine.Curr.CurrentState())
			select {

			case <-m.churnTimer.C:
				rv := m.Machine.ProcessEvent(CdMachineModuleStr, LacpCdmEventActorChurnTimerExpired, nil)
				if rv != nil {
					m.LacpCdmLog(strings.Join([]string{error.Error(rv), CdMachineModuleStr, CdmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(LacpCdmEventActorChurnTimerExpired))}, ":"))
				}
			case event, ok := <-m.CdmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)

					if rv == nil &&
						m.Machine.Curr.CurrentState() == LacpCdmStateNoActorChurn &&
						!LacpStateIsSet(m.p.ActorOper.State, LacpStateSyncBit) {
						rv = m.Machine.ProcessEvent(CdMachineModuleStr, LacpCdmEventActorOperPortStateSyncOff, nil)
					}
					if rv == nil &&
						(m.Machine.Curr.CurrentState() == LacpCdmStateActorChurnMonitor ||
							m.Machine.Curr.CurrentState() == LacpCdmStateActorChurn) &&
						LacpStateIsSet(m.p.ActorOper.State, LacpStateSyncBit) {
						rv = m.Machine.ProcessEvent(CdMachineModuleStr, LacpCdmEventActorOperPortStateSyncOn, nil)
					}

					if rv != nil {
						m.LacpCdmLog(strings.Join([]string{error.Error(rv), event.Src, CdmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					if event.ResponseChan != nil {
						utils.SendResponse(CdMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.LacpCdmLog("Machine End")
					return
				}
			case ena := <-m.CdmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(cdm)
}

// LacpActorCdMachineMain:  802.1ax-2014
// Creation of Actor Churn Detection State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpPartnerCdMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.17 Churn Detection machine
	cdm := LacpPartnerCdMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	cdm.Machine.Start(cdm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpPartnerCdMachine) {
		m.LacpCdmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			m.p.AggPortDebug.AggPortDebugPartnerChurnState = int(m.Machine.Curr.CurrentState())
			select {
			case <-m.churnTimer.C:
				rv := m.Machine.ProcessEvent(PCdMachineModuleStr, LacpCdmEventPartnerChurnTimerExpired, nil)
				if rv != nil {
					m.LacpCdmLog(strings.Join([]string{error.Error(rv), PCdMachineModuleStr, CdmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(LacpCdmEventPartnerChurnTimerExpired))}, ":"))
				}
			case event, ok := <-m.CdmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)

					if rv == nil &&
						m.Machine.Curr.CurrentState() == LacpCdmStateNoActorChurn &&
						!LacpStateIsSet(m.p.PartnerOper.State, LacpStateSyncBit) {
						rv = m.Machine.ProcessEvent(PCdMachineModuleStr, LacpCdmEventActorOperPortStateSyncOff, nil)
					}
					if rv == nil &&
						(m.Machine.Curr.CurrentState() == LacpCdmStateActorChurnMonitor ||
							m.Machine.Curr.CurrentState() == LacpCdmStateActorChurn) &&
						LacpStateIsSet(m.p.PartnerOper.State, LacpStateSyncBit) {
						rv = m.Machine.ProcessEvent(PCdMachineModuleStr, LacpCdmEventActorOperPortStateSyncOn, nil)
					}

					if rv != nil {
						m.LacpCdmLog(strings.Join([]string{error.Error(rv), event.Src, CdmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					if event.ResponseChan != nil {
						utils.SendResponse(PCdMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.LacpCdmLog("Machine End")
					return
				}
			case ena := <-m.CdmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(cdm)
}
