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

// 802.1D-2004 17.29 Port State Transition State Machine
//The Port State Transition state machine shall implement the function specified by the state diagram in Figure
//17-24, the definitions in 17.13, 17.16, 17.20, and 17.21, and variable declarations in 17.17, 17.18, and 17.19.
//This state machine models changes in the Port State (7.4). The Port Role Transitions state machine requests
//changes by setting the learn and forward variables; the Port State Transitions machine updates the learning
//and forwarding variables as the actual Port State changes. The disableLearning(), disableForwarding(),
//enableLearning(), and enableForwarding() procedures model the system-dependent actions and delays that
//take place; these procedures do not complete until the desired behavior has been achieved.
package stp

import (
	"fmt"
	//"time"
	"asicd/pluginManager/pluginCommon"
	"utils/fsm"
)

const PstMachineModuleStr = "PSTM"

const (
	PstStateNone = iota + 1
	PstStateDiscarding
	PstStateLearning
	PstStateForwarding
)

var PstStateStrMap map[fsm.State]string

func PstMachineStrStateMapInit() {
	PstStateStrMap = make(map[fsm.State]string)
	PstStateStrMap[PstStateNone] = "None"
	PstStateStrMap[PstStateDiscarding] = "Discarding"
	PstStateStrMap[PstStateLearning] = "Learning"
	PstStateStrMap[PstStateForwarding] = "Forwarding"
}

const (
	PstEventBegin = iota + 1
	PstEventLearn
	PstEventNotLearn
	PstEventForward
	PstEventNotForward
)

// PstMachine holds FSM and current State
// and event channels for State transitions
type PstMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PstEvents chan MachineEvent
	// enable logging
	PstLogEnableEvent chan bool
}

func (m *PstMachine) GetCurrStateStr() string {
	return PstStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PstMachine) GetPrevStateStr() string {
	return PstStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewStpPrtMachine will create a new instance of the LacpRxMachine
func NewStpPstMachine(p *StpPort) *PstMachine {
	pstm := &PstMachine{
		p:                 p,
		PstEvents:         make(chan MachineEvent, 50),
		PstLogEnableEvent: make(chan bool)}

	p.PstMachineFsm = pstm

	return pstm
}

func (pstm *PstMachine) PstmLogger(s string) {
	StpMachineLogger("DEBUG", PstMachineModuleStr, pstm.p.IfIndex, pstm.p.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (pstm *PstMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if pstm.Machine == nil {
		pstm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	pstm.Machine.Rules = r
	pstm.Machine.Curr = &StpStateEvent{
		strStateMap: PstStateStrMap,
		logEna:      true,
		logger:      pstm.PstmLogger,
		owner:       PstMachineModuleStr,
		ps:          PstStateNone,
		s:           PstStateNone,
	}

	return pstm.Machine
}

// Stop should clean up all resources
func (pstm *PstMachine) Stop() {

	close(pstm.PstEvents)
	close(pstm.PstLogEnableEvent)
}

// PstMachineDiscarding
func (pstm *PstMachine) PstMachineDiscarding(m fsm.Machine, data interface{}) fsm.State {
	p := pstm.p
	pstm.disableLearning()
	defer pstm.NotifyLearningChanged(p.Learning, false)
	p.Learning = false
	pstm.disableForwarding()
	defer pstm.NotifyForwardingChanged(p.Forwarding, false)
	p.Forwarding = false
	return PstStateDiscarding
}

// PstMachineLearning
func (pstm *PstMachine) PstMachineLearning(m fsm.Machine, data interface{}) fsm.State {
	p := pstm.p
	pstm.enableLearning()
	defer pstm.NotifyLearningChanged(p.Learning, true)
	p.Learning = true
	return PstStateLearning
}

// PstMachineForwarding
func (pstm *PstMachine) PstMachineForwarding(m fsm.Machine, data interface{}) fsm.State {
	p := pstm.p
	pstm.enableForwarding()
	defer pstm.NotifyForwardingChanged(p.Forwarding, true)
	p.Forwarding = true
	return PstStateForwarding
}

func PstMachineFSMBuild(p *StpPort) *PstMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrtMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	pstm := NewStpPstMachine(p)

	// BEGIN -> DISCARDING
	rules.AddRule(PstStateNone, PstEventBegin, pstm.PstMachineDiscarding)
	rules.AddRule(PstStateDiscarding, PstEventBegin, pstm.PstMachineDiscarding)
	rules.AddRule(PstStateLearning, PstEventBegin, pstm.PstMachineDiscarding)
	rules.AddRule(PstStateForwarding, PstEventBegin, pstm.PstMachineDiscarding)

	// LEARN -> LEARNING
	rules.AddRule(PstStateDiscarding, PstEventLearn, pstm.PstMachineLearning)

	// NOT LEARN -> LEARNING
	rules.AddRule(PstStateLearning, PstEventNotLearn, pstm.PstMachineDiscarding)

	// FORWARD -> FORWARDING
	rules.AddRule(PstStateLearning, PstEventForward, pstm.PstMachineForwarding)

	// NOT FORWARD -> FORWARDING
	rules.AddRule(PstStateForwarding, PstEventNotForward, pstm.PstMachineDiscarding)

	// Create a new FSM and apply the rules
	pstm.Apply(&rules)

	return pstm
}

// PstmMachineMain:
func (p *StpPort) PstMachineMain() {

	// Build the State machine for STP Port Role Transitions State Machine according to
	// 802.1d Section 17.29
	pstm := PstMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	pstm.Machine.Start(pstm.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	go func(m *PstMachine) {
		StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case event, ok := <-m.PstEvents:

				if ok {
					if m.Machine.Curr.CurrentState() == PstStateNone && event.e != PstEventBegin {
						m.PstEvents <- event
						break
					}

					//fmt.Println("Event Rx", event.src, event.e)
					rv := m.Machine.ProcessEvent(event.src, event.e, nil)
					if rv != nil {
						StpMachineLogger("ERROR", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, event.src, PstStateStrMap[m.Machine.Curr.CurrentState()], event.e))
					} else {
						m.ProcessPostStateProcessing()
					}

					if event.responseChan != nil {
						SendResponse(PstMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}

			case ena := <-m.PstLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(pstm)
}

func (pstm *PstMachine) ProcessPostStateDiscarding() {
	p := pstm.p
	if pstm.Machine.Curr.CurrentState() == PstStateDiscarding {
		if p.Learn {
			rv := pstm.Machine.ProcessEvent(PstMachineModuleStr, PstEventLearn, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, PstMachineModuleStr, PstStateStrMap[pstm.Machine.Curr.CurrentState()], PstEventLearn))
			} else {
				pstm.ProcessPostStateProcessing()
			}

		}
	}
}

func (pstm *PstMachine) ProcessPostStateLearning() {
	p := pstm.p
	if pstm.Machine.Curr.CurrentState() == PstStateLearning {
		if !p.Learn {
			rv := pstm.Machine.ProcessEvent(PstMachineModuleStr, PstEventNotLearn, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, PstMachineModuleStr, PstStateStrMap[pstm.Machine.Curr.CurrentState()], PstEventNotLearn))
			} else {
				pstm.ProcessPostStateProcessing()
			}
		} else if p.Forward {
			rv := pstm.Machine.ProcessEvent(PstMachineModuleStr, PstEventForward, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, PstMachineModuleStr, PstStateStrMap[pstm.Machine.Curr.CurrentState()], PstEventForward))
			} else {
				pstm.ProcessPostStateProcessing()
			}
		}
	}
}

func (pstm *PstMachine) ProcessPostStateForwarding() {
	p := pstm.p
	if pstm.Machine.Curr.CurrentState() == PstStateForwarding {
		if !p.Forward {
			rv := pstm.Machine.ProcessEvent(PstMachineModuleStr, PstEventNotForward, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, PstMachineModuleStr, PstStateStrMap[pstm.Machine.Curr.CurrentState()], PstEventNotForward))
			} else {
				pstm.ProcessPostStateProcessing()
			}
		}
	}
}

func (pstm *PstMachine) ProcessPostStateProcessing() {
	pstm.ProcessPostStateDiscarding()
	pstm.ProcessPostStateLearning()
	pstm.ProcessPostStateForwarding()
}

func (pstm *PstMachine) NotifyLearningChanged(oldlearning bool, newlearning bool) {
	p := pstm.p
	if oldlearning != newlearning {

		// Prt
		if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisablePort {
			if !p.Learning &&
				!p.Forwarding &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
			if !p.Learning &&
				!p.Forwarding &&
				!p.Synced &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBlockPort {
			if !p.Learning &&
				!p.Forwarding &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		}

		// Tc
		if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning {
			if p.Role != PortRoleRootPort &&
				p.Role != PortRoleDesignatedPort &&
				!p.Learning &&
				!(p.RcvdTc || p.RcvdTcn || p.RcvdTcAck || p.TcProp) {
				p.TcMachineFsm.TcEvents <- MachineEvent{
					e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
					src: PstMachineModuleStr,
				}
			}
		}
	}
}

func (pstm *PstMachine) NotifyForwardingChanged(oldforwarding bool, newforwarding bool) {
	p := pstm.p
	if oldforwarding != newforwarding {

		// Prt
		if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisablePort {
			if !p.Learning &&
				!p.Forwarding &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
			if !p.Learning &&
				!p.Forwarding &&
				!p.Synced &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBlockPort {
			if !p.Learning &&
				!p.Forwarding &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
					src: PstMachineModuleStr,
				}
			}
		}

		// Tc
		if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning {
			if p.Role != PortRoleRootPort &&
				p.Role != PortRoleDesignatedPort &&
				!p.Learning &&
				!(p.RcvdTc || p.RcvdTcn || p.RcvdTcAck || p.TcProp) {
				p.TcMachineFsm.TcEvents <- MachineEvent{
					e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
					src: PstMachineModuleStr,
				}
			}
		}
	}
}

func (pstm *PstMachine) disableLearning() {
	p := pstm.p
	StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Calling Asic to do disable learning")
	for _, client := range GetAsicDPluginList() {
		client.SetStgPortState(p.b.StgId, p.IfIndex, pluginCommon.STP_PORT_STATE_BLOCKING)
	}
}

func (pstm *PstMachine) disableForwarding() {
	p := pstm.p
	StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Calling Asic to do disable forwarding")
	for _, client := range GetAsicDPluginList() {
		client.SetStgPortState(p.b.StgId, p.IfIndex, pluginCommon.STP_PORT_STATE_BLOCKING)
	}
}

func (pstm *PstMachine) enableLearning() {
	p := pstm.p
	StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Calling Asic to do enable learning")
	for _, client := range GetAsicDPluginList() {
		client.SetStgPortState(p.b.StgId, p.IfIndex, pluginCommon.STP_PORT_STATE_LEARNING)
	}
}

func (pstm *PstMachine) enableForwarding() {
	p := pstm.p
	StpMachineLogger("DEBUG", PstMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Calling Asic to do enable forwarding")
	for _, client := range GetAsicDPluginList() {
		client.SetStgPortState(p.b.StgId, p.IfIndex, pluginCommon.STP_PORT_STATE_FORWARDING)
	}
	p.ForwardingTransitions += 1
}
