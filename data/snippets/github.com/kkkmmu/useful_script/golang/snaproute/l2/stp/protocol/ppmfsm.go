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

// 802.1D-2004 17.24 Port Protocol Migration state machine
// The Port Protocol Migration state machine shall implement the function specified by the state diagram in
// Figure 17-15, the definitions in 17.16, 17.20, and 17.21, and the variable declarations in 17.17, 17.18, and
// 17.19. It updates sendRSTP (17.19.38) to tell the Port Transmit state machine (17.26) which BPDU types
// (9.3) to transmit, to support interoperability (17.4) with the Spanning Tree Algorithm and Protocol specified
// in previous revisions of this standard.
package stp

import (
	"fmt"
	//"time"
	"utils/fsm"
)

const PpmmMachineModuleStr = "PPMM"

const (
	PpmmStateNone = iota + 1
	PpmmStateCheckingRSTP
	PpmmStateSelectingSTP
	PpmmStateSensing
)

var PpmmStateStrMap map[fsm.State]string

func PpmmMachineStrStateMapInit() {
	PpmmStateStrMap = make(map[fsm.State]string)
	PpmmStateStrMap[PpmmStateNone] = "None"
	PpmmStateStrMap[PpmmStateCheckingRSTP] = "Checking RSTP"
	PpmmStateStrMap[PpmmStateSelectingSTP] = "Selecting STP"
	PpmmStateStrMap[PpmmStateSensing] = "Sensing"
}

const (
	PpmmEventBegin = iota + 1
	PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled
	PpmmEventNotPortEnabled
	PpmmEventMcheck
	PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP
	PpmmEventMdelayWhileEqualZero
	PpmmEventSendRSTPAndRcvdSTP
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type PpmmMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PpmmEvents chan MachineEvent
	// enable logging
	PpmmLogEnableEvent chan bool
}

func (m *PpmmMachine) GetCurrStateStr() string {
	return PpmmStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PpmmMachine) GetPrevStateStr() string {
	return PpmmStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewStpPpmmMachine(p *StpPort) *PpmmMachine {
	ppmm := &PpmmMachine{
		p:                  p,
		PpmmEvents:         make(chan MachineEvent, 50),
		PpmmLogEnableEvent: make(chan bool)}

	p.PpmmMachineFsm = ppmm

	return ppmm
}

func (ppm *PpmmMachine) PpmLogger(s string) {
	StpMachineLogger("DEBUG", PpmmMachineModuleStr, ppm.p.IfIndex, ppm.p.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (ppmm *PpmmMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if ppmm.Machine == nil {
		ppmm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	ppmm.Machine.Rules = r
	ppmm.Machine.Curr = &StpStateEvent{
		strStateMap: PpmmStateStrMap,
		logEna:      true,
		logger:      ppmm.PpmLogger,
		owner:       PpmmMachineModuleStr,
		ps:          PpmmStateNone,
		s:           PpmmStateNone,
	}

	return ppmm.Machine
}

// Stop should clean up all resources
func (ppmm *PpmmMachine) Stop() {

	close(ppmm.PpmmEvents)
	close(ppmm.PpmmLogEnableEvent)

}

func (ppmm *PpmmMachine) InformPtxMachineSendRSTPChanged() {
	p := ppmm.p
	if p.PtxmMachineFsm != nil &&
		p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle {
		/*fmt.Println(p.SendRSTP,
		p.NewInfo,
		p.TxCount < TransmitHoldCountDefault,
		p.HelloWhenTimer.count != 0,
		p.Selected,
		p.UpdtInfo)*/
		if p.SendRSTP == true &&
			p.NewInfo == true &&
			p.TxCount < p.b.TxHoldCount &&
			p.HelloWhenTimer.count != 0 &&
			p.Selected == true &&
			p.UpdtInfo == false {
			p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
				e:   PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
				src: PpmmMachineModuleStr}
		} else if p.SendRSTP == false &&
			p.NewInfo == true &&
			p.Role == PortRoleRootPort &&
			p.TxCount < p.b.TxHoldCount &&
			p.HelloWhenTimer.count != 0 &&
			p.Selected == true &&
			p.UpdtInfo == false {
			p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
				e:   PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
				src: PpmmMachineModuleStr}
		} else if p.SendRSTP == false &&
			p.NewInfo == true &&
			p.Role == PortRoleDesignatedPort &&
			p.TxCount < p.b.TxHoldCount &&
			p.HelloWhenTimer.count != 0 &&
			p.Selected == true &&
			p.UpdtInfo == false {
			p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
				e:   PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
				src: PpmmMachineModuleStr}
		}
	}
}

// PpmMachineCheckingRSTP
func (ppmm *PpmmMachine) PpmmMachineCheckingRSTP(m fsm.Machine, data interface{}) fsm.State {
	p := ppmm.p
	p.Mcheck = false

	sendRSTPchanged := p.SendRSTP != p.RstpVersion
	p.MdelayWhiletimer.count = MigrateTimeDefault

	if sendRSTPchanged {
		p.SendRSTP = p.RstpVersion
		// 17.24
		// Inform Port Transmit State Machine what STP version to send and which BPDU types
		// to support interoperability
		ppmm.InformPtxMachineSendRSTPChanged()
	}

	return PpmmStateCheckingRSTP
}

// PpmmMachineSelectingSTP
func (ppmm *PpmmMachine) PpmmMachineSelectingSTP(m fsm.Machine, data interface{}) fsm.State {
	p := ppmm.p

	sendRSTPchanged := p.SendRSTP != false
	p.MdelayWhiletimer.count = MigrateTimeDefault
	if sendRSTPchanged {
		p.SendRSTP = false
		// 17.24
		// Inform Port Transmit State Machine what STP version to send and which BPDU types
		// to support interoperability
		// TODO change naming
		ppmm.InformPtxMachineSendRSTPChanged()
	}

	return PpmmStateSelectingSTP
}

// PpmmMachineSensing
func (ppmm *PpmmMachine) PpmmMachineSensing(m fsm.Machine, data interface{}) fsm.State {
	p := ppmm.p

	p.RcvdRSTP = false
	p.RcvdSTP = false

	return PpmmStateSensing
}

func PpmmMachineFSMBuild(p *StpPort) *PpmmMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrxmMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	ppmm := NewStpPpmmMachine(p)

	//BEGIN -> CHECKING RSTP
	rules.AddRule(PpmmStateNone, PpmmEventBegin, ppmm.PpmmMachineCheckingRSTP)
	rules.AddRule(PpmmStateCheckingRSTP, PpmmEventBegin, ppmm.PpmmMachineCheckingRSTP)
	rules.AddRule(PpmmStateSelectingSTP, PpmmEventBegin, ppmm.PpmmMachineCheckingRSTP)
	rules.AddRule(PpmmStateSensing, PpmmEventBegin, ppmm.PpmmMachineCheckingRSTP)

	// mdelayWhile != MigrateTime and not portEnable	 -> CHECKING RSTP
	rules.AddRule(PpmmStateCheckingRSTP, PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled, ppmm.PpmmMachineCheckingRSTP)

	// NOT PORT ENABLED -> SENSING/CHECKING RSTP
	rules.AddRule(PpmmStateSelectingSTP, PpmmEventNotPortEnabled, ppmm.PpmmMachineSensing)
	rules.AddRule(PpmmStateSensing, PpmmEventNotPortEnabled, ppmm.PpmmMachineCheckingRSTP)

	// MCHECK ->  SENSING/CHECKING RSTP
	rules.AddRule(PpmmStateSelectingSTP, PpmmEventMcheck, ppmm.PpmmMachineSensing)
	rules.AddRule(PpmmStateSensing, PpmmEventMcheck, ppmm.PpmmMachineCheckingRSTP)

	// RSTP VERSION and NOT SENDING RSTP AND RCVD RSTP -> CHECKING RSTP
	rules.AddRule(PpmmStateSensing, PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP, ppmm.PpmmMachineCheckingRSTP)

	//MDELAYWHILE EQUALS ZERO -> SENSING
	rules.AddRule(PpmmStateCheckingRSTP, PpmmEventMdelayWhileEqualZero, ppmm.PpmmMachineSensing)
	rules.AddRule(PpmmStateSelectingSTP, PpmmEventMdelayWhileEqualZero, ppmm.PpmmMachineSensing)

	// SEND RSTP and RCVD STP
	rules.AddRule(PpmmStateSensing, PpmmEventSendRSTPAndRcvdSTP, ppmm.PpmmMachineSelectingSTP)

	// Create a new FSM and apply the rules
	ppmm.Apply(&rules)

	return ppmm
}

// PrxmMachineMain:
func (p *StpPort) PpmmMachineMain() {

	// Build the State machine for STP Receive Machine according to
	// 802.1d Section 17.23
	ppmm := PpmmMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	ppmm.Machine.Start(ppmm.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PpmmMachine) {
		StpMachineLogger("DEBUG", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {

			case event, ok := <-m.PpmmEvents:

				if ok {
					if m.Machine.Curr.CurrentState() == PpmmStateNone && event.e != PpmmEventBegin {
						m.PpmmEvents <- event
						break
					}

					//fmt.Println("Event Rx", event.src, event.e, PpmmStateStrMap[m.Machine.Curr.CurrentState()])
					rv := m.Machine.ProcessEvent(event.src, event.e, nil)
					if rv != nil {
						StpMachineLogger("DEBUG", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, event.e, PpmmStateStrMap[m.Machine.Curr.CurrentState()]))
					} else {

						// post processing
						m.ProcessPostStateProcessing()
					}

					if event.responseChan != nil {
						SendResponse(PpmmMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}
			case ena := <-m.PpmmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(ppmm)
}

func (ppmm *PpmmMachine) ProcessPostStateCheckingRSTP() {
	// nothing to be done as entry to this state will not allow for state to transition
}
func (ppmm *PpmmMachine) ProcessPostStateSensing() {
	p := ppmm.p
	if ppmm.Machine.Curr.CurrentState() == PpmmStateSensing {
		if !p.PortEnabled {
			rv := ppmm.Machine.ProcessEvent(PpmmMachineModuleStr, PpmmEventNotPortEnabled, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PpmmEventNotPortEnabled, PpmmStateStrMap[ppmm.Machine.Curr.CurrentState()]))
			} else {
				ppmm.ProcessPostStateProcessing()
			}
		} else if p.Mcheck {
			rv := ppmm.Machine.ProcessEvent(PpmmMachineModuleStr, PpmmEventMcheck, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PpmmEventMcheck, PpmmStateStrMap[ppmm.Machine.Curr.CurrentState()]))
			} else {
				ppmm.ProcessPostStateProcessing()
			}
		} else if p.RstpVersion &&
			!p.SendRSTP &&
			p.RcvdRSTP {
			rv := ppmm.Machine.ProcessEvent(PpmmMachineModuleStr, PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP, PpmmStateStrMap[ppmm.Machine.Curr.CurrentState()]))
			} else {
				ppmm.ProcessPostStateProcessing()
			}
		}
	}
}

func (ppmm *PpmmMachine) ProcessPostStateSelectingSTP() {
	p := ppmm.p
	if ppmm.Machine.Curr.CurrentState() == PpmmStateSelectingSTP {
		if !p.PortEnabled {
			rv := ppmm.Machine.ProcessEvent(PpmmMachineModuleStr, PpmmEventNotPortEnabled, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PpmmEventNotPortEnabled, PpmmStateStrMap[ppmm.Machine.Curr.CurrentState()]))
			} else {
				ppmm.ProcessPostStateProcessing()
			}
		} else if p.Mcheck {
			rv := ppmm.Machine.ProcessEvent(PpmmMachineModuleStr, PpmmEventMcheck, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PpmmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PpmmEventMcheck, PpmmStateStrMap[ppmm.Machine.Curr.CurrentState()]))
			} else {
				ppmm.ProcessPostStateProcessing()
			}
		}
	}
}

func (ppmm *PpmmMachine) ProcessPostStateProcessing() {

	ppmm.ProcessPostStateSensing()
	ppmm.ProcessPostStateCheckingRSTP()
	ppmm.ProcessPostStateSelectingSTP()
}
