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

// 802.1D-2004 17.27 Port Information State Machine
//The Port Information state machine shall implement the function specified by the state diagram in Figure 17-
//18, the definitions in 17.13, 17.16, 17.20, and 17.21, and the variable declarations in 17.17, 17.18, and 17.19.
//This state machine is responsible for updating and recording the source (infoIs, 17.19.10) of the Spanning
//Tree information (portPriority 17.19.21, portTimes 17.19.22) used to test the information conveyed
//(msgPriority, 17.19.14; msgTimes, 17.19.15) by received Configuration Messages. If new, superior,
//information arrives on the port, or the existing information is aged out, it sets the reselect variable to request
//the Port Role Selection state machine to update the spanning tree priority vectors held by the Bridge and the
//Bridge’s Port Roles.
package stp

import (
	"fmt"
	//"time"
	"github.com/google/gopacket/layers"
	//"reflect"
	"utils/fsm"
)

const PimMachineModuleStr = "PIM"

const (
	PimStateNone = iota + 1
	PimStateDisabled
	PimStateAged
	PimStateUpdate
	PimStateSuperiorDesignated
	PimStateRepeatedDesignated
	PimStateInferiorDesignated
	PimStateNotDesignated
	PimStateOther
	PimStateCurrent
	PimStateReceive
)

var PimStateStrMap map[fsm.State]string

func PimMachineStrStateMapInit() {
	PimStateStrMap = make(map[fsm.State]string)
	PimStateStrMap[PimStateNone] = "None"
	PimStateStrMap[PimStateDisabled] = "Disabled"
	PimStateStrMap[PimStateAged] = "Aged"
	PimStateStrMap[PimStateUpdate] = "Updated"
	PimStateStrMap[PimStateSuperiorDesignated] = "Superior Designated"
	PimStateStrMap[PimStateRepeatedDesignated] = "Repeated Designated"
	PimStateStrMap[PimStateInferiorDesignated] = "Inferior Designated"
	PimStateStrMap[PimStateNotDesignated] = "Not Designated"
	PimStateStrMap[PimStateOther] = "Other"
	PimStateStrMap[PimStateCurrent] = "Current"
	PimStateStrMap[PimStateReceive] = "Receive"

}

const (
	PimEventBegin = iota + 1
	PimEventNotPortEnabledInfoIsNotEqualDisabled
	PimEventRcvdMsg
	PimEventRcvdMsgAndNotUpdtInfo
	PimEventPortEnabled
	PimEventSelectedAndUpdtInfo
	PimEventUnconditionalFallThrough
	PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg
	PimEventRcvdInfoEqualSuperiorDesignatedInfo
	PimEventRcvdInfoEqualRepeatedDesignatedInfo
	PimEventRcvdInfoEqualInferiorDesignatedInfo
	PimEventRcvdInfoEqualInferiorRootAlternateInfo
	PimEventRcvdInfoEqualOtherInfo
)

// PimMachine holds FSM and current State
// and event channels for State transitions
type PimMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PimEvents chan MachineEvent
	// enable logging
	PimLogEnableEvent chan bool
}

func (m *PimMachine) GetCurrStateStr() string {
	return PimStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PimMachine) GetPrevStateStr() string {
	return PimStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewStpPimMachine will create a new instance of the LacpRxMachine
func NewStpPimMachine(p *StpPort) *PimMachine {
	pim := &PimMachine{
		p:                 p,
		PimEvents:         make(chan MachineEvent, 50),
		PimLogEnableEvent: make(chan bool)}

	p.PimMachineFsm = pim

	return pim
}

func (pim *PimMachine) PimLogger(s string) {
	StpMachineLogger("DEBUG", PimMachineModuleStr, pim.p.IfIndex, pim.p.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (pim *PimMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if pim.Machine == nil {
		pim.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	pim.Machine.Rules = r
	pim.Machine.Curr = &StpStateEvent{
		strStateMap: PimStateStrMap,
		logEna:      true, // this will produce excessive logging as rx packets cause machine to change states constantly
		logger:      pim.PimLogger,
		owner:       PimMachineModuleStr,
		ps:          PimStateNone,
		s:           PimStateNone,
	}

	return pim.Machine
}

// Stop should clean up all resources
func (pim *PimMachine) Stop() {

	close(pim.PimEvents)
	close(pim.PimLogEnableEvent)
}

// PimMachineDisabled
func (pim *PimMachine) PimMachineDisabled(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p
	//defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	//defer p.NotifyProposingChanged(PimMachineModuleStr, p.Proposing, false)
	p.Proposing = false
	defer pim.NotifyProposedChanged(p.Proposed, false)
	p.Proposed = false
	defer pim.NotifyAgreeChanged(p.Agree, false)
	p.Agree = false
	defer pim.NotifyAgreedChanged(p.Agreed, false)
	p.Agreed = false
	p.RcvdInfoWhiletimer.count = 0
	p.InfoIs = PortInfoStateDisabled
	defer p.NotifySelectedChanged(PimMachineModuleStr, p.Selected, false)
	p.Selected = false
	defer pim.NotifyReselectChanged(p.Reselect, true)
	p.Reselect = true
	return PimStateDisabled
}

// PimMachineAged
func (pim *PimMachine) PimMachineAged(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p
	p.InfoIs = PortInfoStateAged
	defer p.NotifySelectedChanged(PimMachineModuleStr, p.Selected, false)
	p.Selected = false
	defer pim.NotifyReselectChanged(p.Reselect, true)
	p.Reselect = true
	return PimStateAged
}

// PimMachineUpdate
func (pim *PimMachine) PimMachineUpdate(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p
	//defer p.NotifyProposingChanged(PimMachineModuleStr, p.Proposing, false)
	p.Proposing = false
	defer pim.NotifyProposedChanged(p.Proposed, false)
	p.Proposed = false
	tmp := p.Agreed && pim.betterorsameinfo(p.InfoIs)
	defer pim.NotifyAgreedChanged(p.Agreed, tmp)
	p.Agreed = tmp
	tmp = p.Synced && p.Agreed
	defer p.NotifySyncedChanged(PimMachineModuleStr, p.Synced, tmp)
	p.Synced = tmp
	// root has not been assigned yet so lets make this the root
	if p.b.RootPortId == 0 {
		p.PortPriority.RootBridgeId = p.b.BridgeIdentifier
		p.PortPriority.RootPathCost = 0
	} else {
		p.PortPriority.RootBridgeId = p.b.BridgePriority.RootBridgeId
		p.PortPriority.RootPathCost = p.b.BridgePriority.RootPathCost
	}
	p.PortPriority.DesignatedBridgeId = p.b.BridgeIdentifier
	p.PortPriority.DesignatedPortId = uint16(p.Priority<<8 | p.PortId)
	p.PortTimes = p.b.BridgeTimes
	//defer p.NotifyUpdtInfoChanged(PimMachineModuleStr, p.UpdtInfo, false)
	p.UpdtInfo = false
	p.InfoIs = PortInfoStateMine
	defer pim.NotifyNewInfoChange(p.NewInfo, true)
	p.NewInfo = true
	return PimStateUpdate
}

// PimMachineCurrent
func (pim *PimMachine) PimMachineCurrent(m fsm.Machine, data interface{}) fsm.State {
	return PimStateCurrent
}

// PimMachineReceive
func (pim *PimMachine) PimMachineReceive(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p
	p.RcvdInfo = pim.rcvInfo(data)

	// rcvd a valid BPDU
	if p.BridgeAssurance {
		p.BAWhileTimer.count = int32(p.b.RootTimes.HelloTime * 3)
		p.BridgeAssuranceInconsistant = false
	}

	return PimStateReceive
}

// PimMachineSuperiorDesignated
func (pim *PimMachine) PimMachineSuperiorDesignated(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p
	defer pim.NotifyAgreedChanged(p.Agreed, false)
	p.Agreed = false
	//defer p.NotifyProposingChanged(PimMachineModuleStr, p.Proposing, false)
	p.Proposing = false
	flags := pim.getRcvdMsgFlags(data)
	pim.recordProposal(flags)
	pim.setTcFlags(flags, data)
	betterorsame := pim.recordPriority(pim.getRcvdMsgPriority(data))
	pim.recordTimes(pim.getRcvdMsgTimes(data))
	tmp := p.Agree && betterorsame
	//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, fmt.Sprintf("SuperiorDesignated: p.Agree[%t] betterorsame[%t] &&[%t]",
	//	p.Agree, betterorsame, tmp))
	defer pim.NotifyAgreeChanged(p.Agree, tmp)
	p.Agree = tmp
	pim.updtRcvdInfoWhile()
	p.InfoIs = PortInfoStateReceived
	defer p.NotifySelectedChanged(PimMachineModuleStr, p.Selected, false)
	p.Selected = false
	defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	defer pim.NotifyReselectChanged(p.Reselect, true)
	p.Reselect = true
	return PimStateSuperiorDesignated
}

// PimMachineRepeatedDesignated
func (pim *PimMachine) PimMachineRepeatedDesignated(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p

	flags := pim.getRcvdMsgFlags(data)
	pim.recordProposal(flags)
	pim.setTcFlags(flags, data)
	pim.updtRcvdInfoWhile()
	defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	return PimStateRepeatedDesignated
}

// PimMachineInferiorDesignated
func (pim *PimMachine) PimMachineInferiorDesignated(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p

	flags := pim.getRcvdMsgFlags(data)
	pim.recordDispute(flags)
	defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	return PimStateInferiorDesignated
}

// PimMachineNotDesignated
func (pim *PimMachine) PimMachineNotDesignated(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p

	flags := pim.getRcvdMsgFlags(data)
	pim.recordAgreement(flags)
	pim.setTcFlags(flags, data)
	defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	return PimStateNotDesignated
}

// PimMachineOther
func (pim *PimMachine) PimMachineOther(m fsm.Machine, data interface{}) fsm.State {
	p := pim.p

	defer p.NotifyRcvdMsgChanged(PimMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	return PimStateOther
}

func PimMachineFSMBuild(p *StpPort) *PimMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrxmMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	pim := NewStpPimMachine(p)
	/*
		PimStateDisabled
		PimStateAged
		PimStateUpdate
		PimStateSuperiorDesignated
		PimStateRepeatedDesignated
		PimStateInferiorDesignated
		PimStateNotDesignated
		PimStateOther
		PimStateCurrent
		PimStateReceive
	*/
	// BEGIN -> DISABLED
	rules.AddRule(PimStateNone, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateDisabled, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateAged, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateUpdate, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateSuperiorDesignated, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateRepeatedDesignated, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateInferiorDesignated, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateNotDesignated, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateOther, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateCurrent, PimEventBegin, pim.PimMachineDisabled)
	rules.AddRule(PimStateReceive, PimEventBegin, pim.PimMachineDisabled)

	// RCVDMSG -> DISABLED
	rules.AddRule(PimStateDisabled, PimEventRcvdMsg, pim.PimMachineDisabled)

	// NOT PORT ENABLED and INFOLS != DISABLED -> DISABLED
	rules.AddRule(PimStateDisabled, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateAged, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateUpdate, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateSuperiorDesignated, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateRepeatedDesignated, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateInferiorDesignated, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateNotDesignated, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateOther, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateCurrent, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)
	rules.AddRule(PimStateReceive, PimEventNotPortEnabledInfoIsNotEqualDisabled, pim.PimMachineDisabled)

	// PORT ENABLED -> AGED
	rules.AddRule(PimStateDisabled, PimEventPortEnabled, pim.PimMachineAged)

	// SELECTED and  UPDTINFO -> UPDATED
	rules.AddRule(PimStateAged, PimEventSelectedAndUpdtInfo, pim.PimMachineUpdate)
	rules.AddRule(PimStateCurrent, PimEventSelectedAndUpdtInfo, pim.PimMachineUpdate)

	// UNCONDITIONAL FALL THROUGH -> CURRENT
	rules.AddRule(PimStateUpdate, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)
	rules.AddRule(PimStateSuperiorDesignated, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)
	rules.AddRule(PimStateRepeatedDesignated, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)
	rules.AddRule(PimStateInferiorDesignated, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)
	rules.AddRule(PimStateNotDesignated, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)
	rules.AddRule(PimStateOther, PimEventUnconditionalFallThrough, pim.PimMachineCurrent)

	// INFOIS == RECEIVED nad RCVDINFOWHILE == 0 and NOT UPDTINFO and NOT RCVDMSG ->  AGED
	rules.AddRule(PimStateCurrent, PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg, pim.PimMachineAged)

	// RCVDMSG and NOT UPDTINFO -> RECEIVE
	rules.AddRule(PimStateCurrent, PimEventRcvdMsgAndNotUpdtInfo, pim.PimMachineReceive)

	// RCVDINFO == SUPERIORDESIGNATEDINFO -> SUPERIOR DESIGNATED
	rules.AddRule(PimStateReceive, PimEventRcvdInfoEqualSuperiorDesignatedInfo, pim.PimMachineSuperiorDesignated)

	// RCVDINFO == REPEATEDDESIGNATEDINFO -> REPEATED DESIGNATED
	rules.AddRule(PimStateReceive, PimEventRcvdInfoEqualRepeatedDesignatedInfo, pim.PimMachineRepeatedDesignated)

	// RCVDINFO == INFERIORDESIGNATEDINFO -> INFERIOR DESIGNATED
	rules.AddRule(PimStateReceive, PimEventRcvdInfoEqualInferiorDesignatedInfo, pim.PimMachineInferiorDesignated)

	// RCVDINFO == NOTDESIGNATEDINFO -> NOT DESIGNATED
	rules.AddRule(PimStateReceive, PimEventRcvdInfoEqualInferiorRootAlternateInfo, pim.PimMachineNotDesignated)

	// RCVDINFO == INFERIORROOTALTERNATEINFO -> OTHER
	rules.AddRule(PimStateReceive, PimEventRcvdInfoEqualOtherInfo, pim.PimMachineOther)

	// Create a new FSM and apply the rules
	pim.Apply(&rules)

	return pim
}

// PimMachineMain:
func (p *StpPort) PimMachineMain() {

	// Build the State machine for STP Receive Machine according to
	// 802.1d Section 17.27
	pim := PimMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	pim.Machine.Start(pim.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PimMachine) {
		StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case event, ok := <-m.PimEvents:

				if ok {
					if m.Machine.Curr.CurrentState() == PimStateNone && event.e != PimEventBegin {
						m.PimEvents <- event
						break
					}

					//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Event Rx src[%s] event[%d] data[%#v]", event.src, event.e, event.data))
					rv := m.Machine.ProcessEvent(event.src, event.e, event.data)
					if rv != nil {
						StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, event.e, PimStateStrMap[m.Machine.Curr.CurrentState()]))
					} else {
						// POST events
						m.ProcessPostStateProcessing(event.data)
					}
					if event.responseChan != nil {
						SendResponse(PimMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}

			case ena := <-m.PimLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(pim)
}

func (pim *PimMachine) ProcessingPostStateDisabled(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateDisabled {
		if p.PortEnabled {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventPortEnabled, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventPortEnabled, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdMsg {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdMsg, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdMsg, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		}
	}
}

func (pim *PimMachine) ProcessingPostStateAged(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateAged {
		if p.Selected &&
			p.UpdtInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventSelectedAndUpdtInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventSelectedAndUpdtInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		}
	}
}

func (pim *PimMachine) ProcessingPostStateUpdate(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateUpdate {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}

func (pim *PimMachine) ProcessingPostStateSuperiorDesignated(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateSuperiorDesignated {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}
func (pim *PimMachine) ProcessingPostStateRepeatedDesignated(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateRepeatedDesignated {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}
func (pim *PimMachine) ProcessingPostStateInferiorDesignated(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateInferiorDesignated {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}
func (pim *PimMachine) ProcessingPostStateNotDesignated(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateNotDesignated {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}
func (pim *PimMachine) ProcessingPostStateOther(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateOther {
		rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventUnconditionalFallThrough, data)
		if rv != nil {
			StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventUnconditionalFallThrough, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
		} else {
			pim.ProcessPostStateProcessing(data)
		}
	}
}
func (pim *PimMachine) ProcessingPostStateCurrent(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateCurrent {
		if p.Selected &&
			p.UpdtInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventSelectedAndUpdtInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventSelectedAndUpdtInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.InfoIs == PortInfoStateReceived &&
			p.RcvdInfoWhiletimer.count == 0 &&
			!p.UpdtInfo &&
			!p.RcvdMsg {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdMsg &&
			!p.UpdtInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdMsgAndNotUpdtInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdMsgAndNotUpdtInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		}
	}
}

func (pim *PimMachine) ProcessingPostStateReceive(data interface{}) {
	p := pim.p
	if pim.Machine.Curr.CurrentState() == PimStateReceive {
		if p.RcvdInfo == SuperiorDesignatedInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdInfoEqualSuperiorDesignatedInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdInfoEqualSuperiorDesignatedInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdInfo == RepeatedDesignatedInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdInfoEqualRepeatedDesignatedInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdInfoEqualRepeatedDesignatedInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdInfo == InferiorDesignatedInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdInfoEqualInferiorDesignatedInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdInfoEqualRepeatedDesignatedInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdInfo == InferiorRootAlternateInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdInfoEqualInferiorRootAlternateInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdInfoEqualInferiorRootAlternateInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		} else if p.RcvdInfo == OtherInfo {
			rv := pim.Machine.ProcessEvent(PimMachineModuleStr, PimEventRcvdInfoEqualOtherInfo, data)
			if rv != nil {
				StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, PimEventRcvdInfoEqualOtherInfo, PimStateStrMap[pim.Machine.Curr.CurrentState()]))
			} else {
				pim.ProcessPostStateProcessing(data)
			}
		}
	}
}

func (pim *PimMachine) ProcessPostStateProcessing(data interface{}) {

	pim.ProcessingPostStateDisabled(data)
	pim.ProcessingPostStateAged(data)
	pim.ProcessingPostStateUpdate(data)
	pim.ProcessingPostStateSuperiorDesignated(data)
	pim.ProcessingPostStateRepeatedDesignated(data)
	pim.ProcessingPostStateInferiorDesignated(data)
	pim.ProcessingPostStateNotDesignated(data)
	pim.ProcessingPostStateOther(data)
	pim.ProcessingPostStateCurrent(data)
	pim.ProcessingPostStateReceive(data)

}

func (pim *PimMachine) NotifyAgreedChanged(oldagreed bool, newagreed bool) {
	p := pim.p
	if oldagreed != newagreed {
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)
		if p.PrtMachineFsm != nil {

			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
				if !p.Forward &&
					!p.Agreed &&
					!p.Proposing &&
					!p.OperEdge &&
					p.Selected &&
					!p.UpdtInfo {

					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.Synced &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					p.RrWhileTimer.count == 0 &&
					!p.Sync &&
					!p.Learn &&
					p.Selected &&
					!p.UpdtInfo {

					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.ReRoot &&
					!p.Sync &&
					!p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					p.RrWhileTimer.count == 0 &&
					!p.Sync &&
					p.Learn &&
					!p.Forward &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.ReRoot &&
					!p.Sync &&
					p.Learn &&
					!p.Forward &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			}
		}

		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}
	}
}

// notify the prt machine
func (pim *PimMachine) NotifyAgreeChanged(oldagree bool, newagree bool) {
	p := pim.p
	if oldagree != newagree {
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)

		if p.PrtMachineFsm != nil {

			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
				if p.Proposed &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.b.AllSynced() &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Proposed &&
					p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
				if !p.Forward &&
					!p.Agreed &&
					!p.Proposing &&
					!p.OperEdge &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.Synced &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					p.RrWhileTimer.count == 0 &&
					!p.Sync &&
					!p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.ReRoot &&
					!p.Sync &&
					!p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					p.RrWhileTimer.count == 0 &&
					!p.Sync &&
					p.Learn &&
					!p.Forward &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Agreed &&
					!p.ReRoot &&
					!p.Sync &&
					!p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
				if p.Proposed &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.b.AllSynced() &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Proposed &&
					p.Agree &&
					p.Proposed &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			}
		}
		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}

	}
}

func (pim *PimMachine) NotifyProposedChanged(oldproposed bool, newproposed bool) {
	p := pim.p
	if oldproposed != newproposed {
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)

		if p.PrtMachineFsm != nil {
			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
				if p.Proposed &&
					p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Proposed &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
				if p.Proposed &&
					!p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Proposed &&
					p.Agree &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			}
		}
		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}
	}
}

func (pim *PimMachine) NotifyReselectChanged(oldreselect bool, newreselect bool) {
	mEvtChan := make([]chan MachineEvent, 0)
	evt := make([]MachineEvent, 0)

	p := pim.p
	if newreselect {
		if p.b.PrsMachineFsm != nil {
			if p.b.PrsMachineFsm.Machine.Curr.CurrentState() == PrsStateRoleSelection {
				mEvtChan = append(mEvtChan, p.b.PrsMachineFsm.PrsEvents)
				evt = append(evt, MachineEvent{e: PrsEventReselect,
					src: PimMachineModuleStr})

				//p.b.PrsMachineFsm.PrsEvents <- MachineEvent{
				//	e:   PrsEventReselect,
				//	src: PimMachineModuleStr,
				//}
			}
		}
		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}
	}
}

func (pim *PimMachine) NotifyNewInfoChange(oldnewinfo bool, newnewinfo bool) {
	p := pim.p
	if oldnewinfo != newnewinfo {
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)

		if p.PtxmMachineFsm != nil {
			if p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle {
				if !p.SendRSTP &&
					p.NewInfo &&
					p.Role == PortRoleDesignatedPort &&
					p.TxCount < p.b.TxHoldCount &&
					p.HelloWhenTimer.count != 0 &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PtxmMachineFsm.PtxmEvents)
					evt = append(evt, MachineEvent{e: PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
					//	e:   PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if !p.SendRSTP &&
					p.NewInfo &&
					p.Role == PortRoleRootPort &&
					p.HelloWhenTimer.count != 0 &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PtxmMachineFsm.PtxmEvents)
					evt = append(evt, MachineEvent{e: PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
					//	e:   PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				}
			}
		}
		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}
	}
}

func (pim *PimMachine) NotifyDisputedChanged(olddisputed bool, newdisputed bool) {
	p := pim.p
	if olddisputed != newdisputed {
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)
		if p.PrtMachineFsm != nil {
			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
				if p.Disputed &&
					!p.OperEdge &&
					p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					//	src: PimMachineModuleStr,
					//}
				} else if p.Disputed &&
					!p.OperEdge &&
					p.Forward &&
					p.Selected &&
					!p.UpdtInfo {
					mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
					evt = append(evt, MachineEvent{e: PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
						src: PimMachineModuleStr})

					//p.PrtMachineFsm.PrtEvents <- MachineEvent{
					//	e:   PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					//	src: PrsMachineModuleStr,
					//}
				}
			}
		}
		if len(evt) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, false)
		}
	}
}

func (pim *PimMachine) isTcnBPDU(bpduLayer interface{}) bool {
	switch bpduLayer.(type) {
	case *layers.BPDUTopology:
		return true
	}
	return false
}

func (pim *PimMachine) rcvInfo(data interface{}) PortDesignatedRcvInfo {
	p := pim.p
	msgRole := StpGetBpduRole(pim.getRcvdMsgFlags(data))
	msgpriority := pim.getRcvdMsgPriority(data)
	msgtimes := pim.getRcvdMsgTimes(data)
	/*
		StpMachineLogger("DEBUG",
			PimMachineModuleStr,
			p.IfIndex,
			p.BrgIfIndex,
			fmt.Sprintf("role[%d] msgVector[%#v] portVector[%#v] msgTimes[%#v] designatedTimes[%#v]",
				msgRole,
				msgpriority,
				p.PortPriority,
				msgtimes,
				p.PortTimes))
	*/
	if pim.isTcnBPDU(data) {
		return OtherInfo
	}
	// 17.21.8 NOTE
	switch data.(type) {
	case *layers.STP:
		msgRole = PortRoleDesignatedPort
	}

	if CompareBridgeAddr(GetBridgeAddrFromBridgeId(msgpriority.RootBridgeId), GetBridgeAddrFromBridgeId(p.b.BridgeIdentifier)) == 0 {
		if GetBridgePriorityFromBridgeId(msgpriority.RootBridgeId) != GetBridgePriorityFromBridgeId(p.b.BridgeIdentifier) {
			//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: Root bridge addr == bridge addr and root priority != bridge priority")
			return OtherInfo
		}
	}

	//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("rcvInfo: vector superior %t equal %t port times diff %t",
	//	IsMsgPriorityVectorSuperiorThanPortPriorityVector(msgpriority, &p.PortPriority), *msgpriority == p.PortPriority, *msgtimes != p.PortTimes))
	if msgRole == PortRoleDesignatedPort &&
		(*msgpriority == p.PortPriority &&
			*msgtimes != p.PortTimes) {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: msg times != port times")
		return SuperiorDesignatedInfo
	} else if msgRole == PortRoleDesignatedPort &&
		*msgpriority == p.PortPriority &&
		*msgtimes == p.PortTimes {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: msg and port vector and times equal")
		return RepeatedDesignatedInfo
	} else if msgRole == PortRoleDesignatedPort &&
		IsMsgPriorityVectorWorseThanPortPriorityVector(msgpriority, &p.PortPriority) {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: msg priority vector inferior")
		return InferiorDesignatedInfo
	} else if (msgRole == PortRoleRootPort ||
		msgRole == PortRoleAlternatePort ||
		msgRole == PortRoleBackupPort) &&
		IsMsgPriorityVectorTheSameOrWorseThanPortPriorityVector(msgpriority, &p.PortPriority) {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: msg role is root, alternate or backup and msg priority is same or inferior")
		return InferiorRootAlternateInfo
	} else if msgRole == PortRoleDesignatedPort &&
		IsMsgPriorityVectorSuperiorThanPortPriorityVector(msgpriority, &p.PortPriority) {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: msg priority vector superior")
		return SuperiorDesignatedInfo
	} else {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "rcvInfo: default")
		return OtherInfo
	}
}

func (pim *PimMachine) getRcvdMsgFlags(bpduLayer interface{}) uint8 {
	//p := pim.p
	//bpdumsg := data.(RxBpduPdu)
	//packet := bpdumsg.pdu.(gopacket.Packet)
	//bpduLayer := packet.Layer(layers.LayerTypeBPDU)

	var flags uint8
	switch bpduLayer.(type) {
	case *layers.STP:
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, "Found STP frame getting flags")
		stp := bpduLayer.(*layers.STP)
		flags = uint8(stp.Flags)
	case *layers.RSTP:
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, "Found RSTP frame getting flags")
		rstp := bpduLayer.(*layers.RSTP)
		flags = uint8(rstp.Flags)
	case *layers.PVST:
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, "Found PVST frame getting flags")
		pvst := bpduLayer.(*layers.PVST)
		flags = uint8(pvst.Flags)
		//default:
		//	StpMachineLogger("ERROR", PimMachineModuleStr, p.IfIndex, fmt.Sprintf("Error getRcvdMsgFlags rcvd TCN %T\n", bpduLayer))
	}
	return flags
}

func (pim *PimMachine) getRcvdMsgPriority(bpduLayer interface{}) (msgpriority *PriorityVector) {
	msgpriority = &PriorityVector{}

	switch bpduLayer.(type) {
	case *layers.STP:
		stp := bpduLayer.(*layers.STP)
		msgpriority.RootBridgeId = stp.RootId
		msgpriority.RootPathCost = stp.RootPathCost
		msgpriority.DesignatedBridgeId = stp.BridgeId
		msgpriority.DesignatedPortId = stp.PortId
		msgpriority.BridgePortId = stp.PortId
	case *layers.RSTP:
		rstp := bpduLayer.(*layers.RSTP)
		msgpriority.RootBridgeId = rstp.RootId
		msgpriority.RootPathCost = rstp.RootPathCost
		msgpriority.DesignatedBridgeId = rstp.BridgeId
		msgpriority.DesignatedPortId = rstp.PortId
		msgpriority.BridgePortId = rstp.PortId
	case *layers.PVST:
		pvst := bpduLayer.(*layers.PVST)
		msgpriority.RootBridgeId = pvst.RootId
		msgpriority.RootPathCost = pvst.RootPathCost
		msgpriority.DesignatedBridgeId = pvst.BridgeId
		msgpriority.DesignatedPortId = pvst.PortId
		msgpriority.BridgePortId = pvst.PortId
	}
	return msgpriority
}

func (pim *PimMachine) getRcvdMsgTimes(bpduLayer interface{}) (msgtimes *Times) {
	msgtimes = &Times{}

	switch bpduLayer.(type) {
	case *layers.STP:
		stp := bpduLayer.(*layers.STP)
		msgtimes.MessageAge = stp.MsgAge >> 8
		msgtimes.MaxAge = stp.MaxAge >> 8
		msgtimes.HelloTime = stp.HelloTime >> 8
		msgtimes.ForwardingDelay = stp.FwdDelay >> 8
	case *layers.RSTP:
		rstp := bpduLayer.(*layers.RSTP)
		msgtimes.MessageAge = rstp.MsgAge >> 8
		msgtimes.MaxAge = rstp.MaxAge >> 8
		msgtimes.HelloTime = rstp.HelloTime >> 8
		msgtimes.ForwardingDelay = rstp.FwdDelay >> 8
	case *layers.PVST:
		pvst := bpduLayer.(*layers.PVST)
		msgtimes.MessageAge = pvst.MsgAge >> 8
		msgtimes.MaxAge = pvst.MaxAge >> 8
		msgtimes.HelloTime = pvst.HelloTime >> 8
		msgtimes.ForwardingDelay = pvst.FwdDelay >> 8
	}
	return msgtimes
}

// recordProposal(): 17.21.11
func (pim *PimMachine) recordProposal(rcvdMsgFlags uint8) {
	p := pim.p
	if StpGetBpduRole(rcvdMsgFlags) == PortRoleDesignatedPort &&
		StpGetBpduProposal(rcvdMsgFlags) {
		//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, "recording proposal set")
		defer pim.NotifyProposedChanged(p.Proposed, true)
		p.Proposed = true
	}
}

// setTcFlags 17,21,17
//Sets rcvdTc and/or rcvdTcAck if the Topology Change and/or Topology Change Acknowledgment flags,
//respectively, are set in a ConfigBPDU or RST BPDU. Sets rcvdTcn TRUE if the BPDU is a TCN BPDU.
func (pim *PimMachine) setTcFlags(rcvdMsgFlags uint8, bpduLayer interface{}) {
	p := pim.p
	//bpdumsg := data.(RxBpduPdu)
	//packet := bpdumsg.pdu.(gopacket.Packet)
	//bpduLayer := packet.Layer(layers.LayerTypeBPDU)

	//p.NotifyRcvdTcRcvdTcnRcvdTcAck(p.RcvdTc, p.RcvdTcn, p.RcvdTcAck, StpGetBpduTopoChange(rcvdMsgFlags), pim.isTcnBPDU(bpduLayer), StpGetBpduTopoChangeAck(rcvdMsgFlags))
	p.RcvdTcAck = StpGetBpduTopoChangeAck(rcvdMsgFlags)
	p.RcvdTc = StpGetBpduTopoChange(rcvdMsgFlags)
	p.RcvdTcn = pim.isTcnBPDU(bpduLayer)
}

// betterorsameinfo 17.21.1
func (pim *PimMachine) betterorsameinfo(newInfoIs PortInfoState) bool {
	p := pim.p

	StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("betterorsameinfo: newInfoIs[%d] p.InfoIs[%d] msgPriority[%#v] portPriority[%#v]",
		newInfoIs, p.InfoIs, p.MsgPriority, p.PortPriority))

	// recordPriority should be called when this is called from superior designated
	// this way we don't need to pass the message around
	if newInfoIs == PortInfoStateReceived &&
		p.InfoIs == PortInfoStateReceived &&
		IsMsgPriorityVectorSuperiorOrSameThanPortPriorityVector(&p.MsgPriority, &p.PortPriority) {
		StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "betterorsameinfo: UPDATE InfoIs=Receive and msg vector superior or same as port")
		return true
	} else if newInfoIs == PortInfoStateMine &&
		p.InfoIs == PortInfoStateMine &&
		IsMsgPriorityVectorSuperiorOrSameThanPortPriorityVector(&p.b.BridgePriority, &p.PortPriority) {
		StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, "betterorsameinfo: InfoIs=Mine and designated vector superior or same as port")
		return true
	}
	return false
}

// recordPriority  17.21.12
func (pim *PimMachine) recordPriority(rcvdMsgPriority *PriorityVector) bool {
	p := pim.p
	p.MsgPriority = *rcvdMsgPriority
	// 17.6
	//This message priority vector is superior to the port priority vector and will replace it if, and only if, the
	//message priority vector is better than the port priority vector, or the message has been transmitted from the
	//same Designated Bridge and Designated Port as the port priority vector, i.e., if the following is true
	betterorsame := pim.betterorsameinfo(p.InfoIs)
	//if betterorsame {
	//p.PortPriority.RootBridgeId = rcvdMsgPriority.RootBridgeId
	//p.PortPriority.RootPathCost = rcvdMsgPriority.RootPathCost
	//}
	p.PortPriority = *rcvdMsgPriority
	//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, fmt.Sprintf("recordPriority: copying rcvmsg to port %#v", *rcvdMsgPriority))
	return betterorsame
}

// recordTimes 17.21.13
func (pim *PimMachine) recordTimes(rcvdMsgTimes *Times) {
	p := pim.p
	p.PortTimes.ForwardingDelay = rcvdMsgTimes.ForwardingDelay
	if rcvdMsgTimes.HelloTime > BridgeHelloTimeMin {
		p.PortTimes.HelloTime = rcvdMsgTimes.HelloTime
	} else {
		p.PortTimes.HelloTime = BridgeHelloTimeMin
	}
	p.PortTimes.MaxAge = rcvdMsgTimes.MaxAge
	p.PortTimes.MessageAge = rcvdMsgTimes.MessageAge
}

// updtRcvdInfoWhile 17.21.23
func (pim *PimMachine) updtRcvdInfoWhile() {
	p := pim.p
	//StpMachineLogger("DEBUG", PimMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("PortTimes msgAge[%d] maxAge[%d]", p.PortTimes.MessageAge, p.PortTimes.MaxAge))
	if p.PortTimes.MessageAge+1 <= p.PortTimes.MaxAge {
		p.RcvdInfoWhiletimer.count = 3 * int32(p.PortTimes.HelloTime)
	} else {
		p.RcvdInfoWhiletimer.count = 0
		// TODO what happens when this is set
	}
}

func (pim *PimMachine) recordDispute(rcvdMsgFlags uint8) {
	p := pim.p
	msgRole := StpGetBpduRole(rcvdMsgFlags)

	if msgRole == PortRoleDesignatedPort &&
		(StpGetBpduLearning(rcvdMsgFlags) ||
			StpGetBpduForwarding(rcvdMsgFlags)) {
		pim.NotifyDisputedChanged(p.Disputed, true)
		p.Disputed = true
		defer pim.NotifyAgreedChanged(p.Agreed, false)
		p.Agreed = false
	}
	if StpGetBpduLearning(rcvdMsgFlags) {
		defer pim.NotifyAgreedChanged(p.Agreed, true)
		p.Agreed = true
		defer p.NotifyProposingChanged(PimMachineModuleStr, p.Proposing, false)
		p.Proposing = false
	}

}

func (pim *PimMachine) recordAgreement(rcvdMsgFlags uint8) {
	p := pim.p
	if p.RstpVersion &&
		p.OperPointToPointMAC &&
		StpGetBpduAgreement(rcvdMsgFlags) {
		defer pim.NotifyAgreedChanged(p.Agreed, true)
		p.Agreed = true
		//defer p.NotifyProposingChanged(PimMachineModuleStr, p.Proposing, false)
		p.Proposing = false
	} else {
		defer pim.NotifyAgreedChanged(p.Agreed, false)
		p.Agreed = false
	}

}
