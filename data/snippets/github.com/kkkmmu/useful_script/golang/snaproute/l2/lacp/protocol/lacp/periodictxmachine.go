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

// The Periodic Transmission Machine is described in the 802.1ax-2014 Section 6.4.13
package lacp

import (
	//"fmt"
	"l2/lacp/protocol/utils"
	"time"
	"utils/fsm"
)

const PtxMachineModuleStr = "Periodic TX Machine"

const (
	LacpPtxmStateNone = iota + 1
	LacpPtxmStateNoPeriodic
	LacpPtxmStateFastPeriodic
	LacpPtxmStateSlowPeriodic
	LacpPtxmStatePeriodicTx
)

var PtxmStateStrMap map[fsm.State]string

func PtxMachineStrStateMapCreate() {
	PtxmStateStrMap = make(map[fsm.State]string)
	PtxmStateStrMap[LacpPtxmStateNone] = "None"
	PtxmStateStrMap[LacpPtxmStateNoPeriodic] = "NoPeriodic"
	PtxmStateStrMap[LacpPtxmStateFastPeriodic] = "FastPeriodic"
	PtxmStateStrMap[LacpPtxmStateSlowPeriodic] = "SlowPeriodic"
	PtxmStateStrMap[LacpPtxmStatePeriodicTx] = "PeriodicTx"
}

const (
	LacpPtxmEventBegin = iota + 1
	LacpPtxmEventLacpDisabled
	LacpPtxmEventNotPortEnabled
	LacpPtxmEventActorPartnerOperActivityPassiveMode
	LacpPtxmEventUnconditionalFallthrough
	LacpPtxmEventPartnerOperStateTimeoutLong
	LacpPtxmEventPeriodicTimerExpired
	LacpPtxmEventPartnerOperStateTimeoutShort
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type LacpPtxMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	// Reference to LaAggPort
	p *LaAggPort

	// current tx interval LONG/SHORT
	PeriodicTxTimerInterval time.Duration

	// timer
	periodicTxTimer *time.Timer

	// machine specific events
	PtxmEvents chan utils.MachineEvent
	// enable logging
	PtxmLogEnableEvent chan bool
}

func (ptxm *LacpPtxMachine) PrevState() fsm.State { return ptxm.PreviousState }

// PrevStateSet will set the previous State
func (ptxm *LacpPtxMachine) PrevStateSet(s fsm.State) { ptxm.PreviousState = s }

func (ptxm *LacpPtxMachine) Stop() {

	ptxm.PeriodicTimerStop()

	close(ptxm.PtxmEvents)
	close(ptxm.PtxmLogEnableEvent)
}

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewLacpPtxMachine(port *LaAggPort) *LacpPtxMachine {
	ptxm := &LacpPtxMachine{
		p:                       port,
		PreviousState:           LacpPtxmStateNone,
		PeriodicTxTimerInterval: LacpSlowPeriodicTime,
		PtxmEvents:              make(chan utils.MachineEvent, 10),
		PtxmLogEnableEvent:      make(chan bool)}

	port.PtxMachineFsm = ptxm

	// start then stop
	ptxm.PeriodicTimerStart()
	ptxm.PeriodicTimerStop()

	return ptxm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (ptxm *LacpPtxMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if ptxm.Machine == nil {
		ptxm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	ptxm.Machine.Rules = r
	ptxm.Machine.Curr = &utils.StateEvent{
		StrStateMap: PtxmStateStrMap,
		//logEna:      ptxm.p.logEna,
		LogEna: false,
		Logger: ptxm.LacpPtxmLog,
		Owner:  PtxMachineModuleStr,
	}

	return ptxm.Machine
}

// LacpPtxMachineNoPeriodic stops the periodic transmission of packets
func (ptxm *LacpPtxMachine) LacpPtxMachineNoPeriodic(m fsm.Machine, data interface{}) fsm.State {
	ptxm.PeriodicTimerStop()
	return LacpPtxmStateNoPeriodic
}

// LacpPtxMachineFastPeriodic sets the periodic transmission time to fast
// and starts the timer
func (ptxm *LacpPtxMachine) LacpPtxMachineFastPeriodic(m fsm.Machine, data interface{}) fsm.State {
	ptxm.PeriodicTimerIntervalSet(LacpFastPeriodicTime)
	ptxm.PeriodicTimerStart()
	return LacpPtxmStateFastPeriodic
}

// LacpPtxMachineSlowPeriodic sets the periodic transmission time to slow
// and starts the timer
func (ptxm *LacpPtxMachine) LacpPtxMachineSlowPeriodic(m fsm.Machine, data interface{}) fsm.State {
	ptxm.PeriodicTimerIntervalSet(LacpSlowPeriodicTime)
	ptxm.PeriodicTimerStart()
	return LacpPtxmStateSlowPeriodic
}

// LacpPtxMachinePeriodicTx informs the tx machine that a packet should be transmitted by setting
// ntt = true
func (ptxm *LacpPtxMachine) LacpPtxMachinePeriodicTx(m fsm.Machine, data interface{}) fsm.State {
	// inform the tx machine that ntt should change to true which should transmit a
	// packet
	if ptxm.p.TxMachineFsm.Machine.Curr.CurrentState() != LacpTxmStateOff {
		ptxm.p.TxMachineFsm.TxmEvents <- utils.MachineEvent{
			E:   LacpTxmEventNtt,
			Src: PtxMachineModuleStr}
	}

	return LacpPtxmStatePeriodicTx
}

func LacpPtxMachineFSMBuild(p *LaAggPort) *LacpPtxMachine {

	rules := fsm.Ruleset{}

	PtxMachineStrStateMapCreate()
	p.wg.Add(1)

	// Instantiate a new LacpPtxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the NO PERIODIC State
	ptxm := NewLacpPtxMachine(p)

	//BEGIN -> NO PERIODIC
	rules.AddRule(LacpPtxmStateNone, LacpPtxmEventBegin, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventBegin, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventBegin, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventBegin, ptxm.LacpPtxMachineNoPeriodic)
	// LACP DISABLED -> NO PERIODIC
	rules.AddRule(LacpPtxmStateNone, LacpPtxmEventLacpDisabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventLacpDisabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventLacpDisabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventLacpDisabled, ptxm.LacpPtxMachineNoPeriodic)
	// PORT DISABLED -> NO PERIODIC
	rules.AddRule(LacpPtxmStateNone, LacpPtxmEventNotPortEnabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventNotPortEnabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventNotPortEnabled, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventNotPortEnabled, ptxm.LacpPtxMachineNoPeriodic)
	// ACTOR/PARTNER OPER State ACTIVITY MODE == PASSIVE -> NO PERIODIC
	rules.AddRule(LacpPtxmStateNone, LacpPtxmEventActorPartnerOperActivityPassiveMode, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventActorPartnerOperActivityPassiveMode, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventActorPartnerOperActivityPassiveMode, ptxm.LacpPtxMachineNoPeriodic)
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventActorPartnerOperActivityPassiveMode, ptxm.LacpPtxMachineNoPeriodic)
	// INTENTIONAL FALL THROUGH -> FAST PERIODIC
	rules.AddRule(LacpPtxmStateNoPeriodic, LacpPtxmEventUnconditionalFallthrough, ptxm.LacpPtxMachineFastPeriodic)
	// PARTNER OPER STAT LACP TIMEOUT == LONG -> SLOW PERIODIC
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventPartnerOperStateTimeoutLong, ptxm.LacpPtxMachineSlowPeriodic)
	// PERIODIC TIMER EXPIRED -> PERIODIC TX
	rules.AddRule(LacpPtxmStateFastPeriodic, LacpPtxmEventPeriodicTimerExpired, ptxm.LacpPtxMachinePeriodicTx)
	// PARTNER OPER STAT LACP TIMEOUT == SHORT ->  PERIODIC TX
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventPartnerOperStateTimeoutShort, ptxm.LacpPtxMachinePeriodicTx)
	// PERIODIC TIMER EXPIRED -> PERIODIC TX
	rules.AddRule(LacpPtxmStateSlowPeriodic, LacpPtxmEventPeriodicTimerExpired, ptxm.LacpPtxMachinePeriodicTx)
	// PARTNER OPER STAT LACP TIMEOUT == SHORT ->  FAST PERIODIC
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventPartnerOperStateTimeoutShort, ptxm.LacpPtxMachineFastPeriodic)
	// PARTNER OPER STAT LACP TIMEOUT == LONG -> SLOW PERIODIC
	rules.AddRule(LacpPtxmStatePeriodicTx, LacpPtxmEventPartnerOperStateTimeoutLong, ptxm.LacpPtxMachineSlowPeriodic)

	// Create a new FSM and apply the rules
	ptxm.Apply(&rules)

	return ptxm
}

// LacpRxMachineMain:  802.1ax-2014 Table 6-18
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpPtxMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.13 Periodic Transmission Machine
	ptxm := LacpPtxMachineFSMBuild(p)

	// set the inital State
	ptxm.Machine.Start(ptxm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpPtxMachine) {
		m.LacpPtxmLog("PTXM: Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case <-m.periodicTxTimer.C:
				//m.LacpPtxmLog("Timer expired current State")
				//m.LacpPtxmLog(PtxmStateStrMap[m.Machine.Curr.CurrentState()])
				m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventPeriodicTimerExpired, nil)

				if m.Machine.Curr.CurrentState() == LacpPtxmStatePeriodicTx {
					if LacpStateIsSet(m.p.PartnerOper.State, LacpStateTimeoutBit) {
						m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventPartnerOperStateTimeoutShort, nil)
					} else {
						m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventPartnerOperStateTimeoutLong, nil)
					}
				}

			case event, ok := <-m.PtxmEvents:

				if ok {
					tmpLogEna := false
					if !m.Machine.Curr.IsLoggerEna() {
						tmpLogEna = true
						m.Machine.Curr.EnableLogging(true)
					}
					m.Machine.ProcessEvent(event.Src, event.E, nil)
					/* special case */
					if m.LacpPtxIsNoPeriodicExitCondition() {
						m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventUnconditionalFallthrough, nil)
					} else if m.Machine.Curr.CurrentState() == LacpPtxmStatePeriodicTx {
						if LacpStateIsSet(m.p.PartnerOper.State, LacpStateTimeoutBit) {
							m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventPartnerOperStateTimeoutShort, nil)
						} else {
							m.Machine.ProcessEvent(PtxMachineModuleStr, LacpPtxmEventPartnerOperStateTimeoutLong, nil)
						}
					}

					if event.ResponseChan != nil {
						utils.SendResponse(PtxMachineModuleStr, event.ResponseChan)
					}

					if tmpLogEna {
						m.Machine.Curr.EnableLogging(false)
						tmpLogEna = false
					}
				} else {
					m.LacpPtxmLog("Machine End")
					return
				}

			case ena := <-m.PtxmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(ptxm)
}

// LacpPtxIsNoPeriodicExitCondition is meant to check if the UTC
// condition has been met when the State is NO PERIODIC
func (m *LacpPtxMachine) LacpPtxIsNoPeriodicExitCondition() bool {
	p := m.p
	/*
		m.LacpPtxmLog(fmt.Sprintf("LacpPtxIsNoPeriodicExitCondition: State %d ena %d lacpEna %d mode set %d State 0x%x",
			m.Machine.Curr.CurrentState(),
			p.portEnabled,
			p.lacpEnabled,
			LacpModeGet(p.ActorOper.State, p.lacpEnabled),
			p.ActorOper.State))
	*/
	return m.Machine.Curr.CurrentState() == LacpPtxmStateNoPeriodic &&
		p.lacpEnabled &&
		p.IsPortEnabled() &&
		(LacpModeGet(p.ActorOper.State, p.lacpEnabled) == LacpModeActive ||
			LacpModeGet(p.PartnerOper.State, p.lacpEnabled) == LacpModeActive)
}
