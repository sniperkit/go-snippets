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

// 17.22 Port Timers state machine
package stp

import (
	"time"
	"utils/fsm"
)

const PtmMachineModuleStr = "PTIM"

const (
	PtmStateNone = iota + 1
	PtmStateOneSecond
	PtmStateTick
)

var PtmStateStrMap map[fsm.State]string

func PtmMachineStrStateMapInit() {
	PtmStateStrMap = make(map[fsm.State]string)
	PtmStateStrMap[PtmStateNone] = "None"
	PtmStateStrMap[PtmStateOneSecond] = "OneSecond"
	PtmStateStrMap[PtmStateTick] = "Tick"
}

const (
	PtmEventBegin = iota + 1
	PtmEventTickEqualsTrue
	PtmEventUnconditionalFallthrough
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type PtmMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	// State transition log
	log chan string

	// timer type
	TickTimer *time.Timer
	Tick      bool

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PtmEvents chan MachineEvent
	// enable logging
	PtmLogEnableEvent chan bool
}

func (m *PtmMachine) GetCurrStateStr() string {
	return PtmStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PtmMachine) GetPrevStateStr() string {
	return PtmStateStrMap[m.Machine.Curr.PreviousState()]
}

func (ptm *PtmMachine) PrevState() fsm.State { return ptm.PreviousState }

// PrevStateSet will set the previous State
func (ptm *PtmMachine) PrevStateSet(s fsm.State) { ptm.PreviousState = s }

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewStpPtmMachine(p *StpPort) *PtmMachine {
	ptm := &PtmMachine{
		p:                 p,
		PreviousState:     PtmStateNone,
		PtmEvents:         make(chan MachineEvent, 50),
		PtmLogEnableEvent: make(chan bool)}

	// start then stop
	ptm.TickTimerStart()
	ptm.TickTimerStop()

	p.PtmMachineFsm = ptm

	return ptm
}

func (ptm *PtmMachine) PtmLogger(s string) {
	//StpMachineLogger("DEBUG", PtmMachineModuleStr, ptm.p.IfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (ptm *PtmMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if ptm.Machine == nil {
		ptm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	ptm.Machine.Rules = r
	ptm.Machine.Curr = &StpStateEvent{
		strStateMap: PtmStateStrMap,
		logEna:      true, // WARNING do not enable as this will cause a log ever second
		logger:      ptm.PtmLogger,
		owner:       PtmMachineModuleStr,
		ps:          PtmStateNone,
		s:           PtmStateNone,
	}

	return ptm.Machine
}

// Stop should clean up all resources
func (ptm *PtmMachine) Stop() {
	ptm.TickTimerDestroy()
	close(ptm.PtmEvents)

}

// LacpPtxMachineNoPeriodic stops the periodic transmission of packets
func (ptm *PtmMachine) PtmMachineOneSecond(m fsm.Machine, data interface{}) fsm.State {
	ptm.Tick = false
	return PtmStateOneSecond
}

// LacpPtxMachineFastPeriodic sets the periodic transmission time to fast
// and starts the timer
func (ptm *PtmMachine) PtmMachineTick(m fsm.Machine, data interface{}) fsm.State {
	p := ptm.p
	p.DecrementTimerCounters()

	return PtmStateTick
}

func PtmMachineFSMBuild(p *StpPort) *PtmMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new LacpPtxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the NO PERIODIC State
	ptm := NewStpPtmMachine(p)

	//BEGIN -> ONE SECOND
	rules.AddRule(PtmStateNone, PtmEventBegin, ptm.PtmMachineOneSecond)
	rules.AddRule(PtmStateOneSecond, PtmEventBegin, ptm.PtmMachineOneSecond)
	rules.AddRule(PtmStateTick, PtmEventBegin, ptm.PtmMachineOneSecond)

	// TICK EQUALS TRUE	 -> TICK
	rules.AddRule(PtmStateOneSecond, PtmEventTickEqualsTrue, ptm.PtmMachineTick)

	// PORT DISABLED -> NO PERIODIC
	rules.AddRule(PtmStateTick, PtmEventUnconditionalFallthrough, ptm.PtmMachineOneSecond)

	// Create a new FSM and apply the rules
	ptm.Apply(&rules)

	return ptm
}

// LacpRxMachineMain:  802.1ax-2014 Table 6-18
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *StpPort) PtmMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.13 Periodic Transmission Machine
	ptm := PtmMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	ptm.Machine.Start(ptm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PtmMachine) {
		StpMachineLogger("DEBUG", PtmMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case <-m.TickTimer.C:
				m.Tick = true
				m.Machine.ProcessEvent(PtmMachineModuleStr, PtmEventTickEqualsTrue, nil)

				// post state processing
				if m.Machine.Curr.CurrentState() == PtmStateTick {
					m.Machine.ProcessEvent(PtmMachineModuleStr, PtmEventUnconditionalFallthrough, nil)

				}
				// restart the timer
				m.TickTimerStart()

			case event, ok := <-m.PtmEvents:

				if ok {
					m.Machine.ProcessEvent(event.src, event.e, nil)

					if event.responseChan != nil {
						SendResponse(PtmMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PtmMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}
			case ena := <-m.PtmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(ptm)
}
