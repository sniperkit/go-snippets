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

// TX MACHINE, this is not really a State machine but going to create a sort of
// State machine to processes events
// TX Machine is described in 802.1ax-2014 6.4.16
package drcp

import (
	"l2/lacp/protocol/utils"
	"strconv"
	"strings"
	"utils/fsm"
)

const NetIplShareMachineModuleStr = "NET/IPL Sharing Machine"

const (
	NetIplSharemStateNone = iota + 1
	NetIplSharemStateNoManipulatedFramesSent
	NetIplSharemStateTimeShareMethod
	NetIplSharemStateManipulatedFramesSent
)

var NetIplSharemStateStrMap map[fsm.State]string

func NetIplShareMachineStrStateMapCreate() {

	NetIplSharemStateStrMap = make(map[fsm.State]string)
	NetIplSharemStateStrMap[NetIplSharemStateNone] = "None"
	NetIplSharemStateStrMap[NetIplSharemStateNoManipulatedFramesSent] = "No Manipulated Frames Sent"
	NetIplSharemStateStrMap[NetIplSharemStateTimeShareMethod] = "Time Share Method"
	NetIplSharemStateStrMap[NetIplSharemStateManipulatedFramesSent] = "Manipulated Frames Sent"

}

const (
	NetIplSharemEventBegin = iota + 1
	NetIplSharemEventCCTimeShare
	NetIplSharemEventCCEncTagShared
	NetIplSharemEventNotCCTimeShare
	NetIplSharemEventNotCCEncTagShared
)

// NetIplShareMachine holds FSM and current State
// and event channels for State transitions
type NetIplShareMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	// Port this Machine is associated with
	p *DRCPIpp

	// machine specific events
	NetIplSharemEvents chan utils.MachineEvent
}

// PrevState will get the previous State from the State transitions
func (nism *NetIplShareMachine) PrevState() fsm.State { return nism.PreviousState }

// PrevStateSet will set the previous State
func (nism *NetIplShareMachine) PrevStateSet(s fsm.State) { nism.PreviousState = s }

// Stop will stop all timers and close all channels
func (nism *NetIplShareMachine) Stop() {

	close(nism.NetIplSharemEvents)
}

// NewDrcpTxMachine will create a new instance of the TxMachine
func NewDrcpNetIplShareMachine(port *DRCPIpp) *NetIplShareMachine {
	nism := &NetIplShareMachine{
		p:                  port,
		PreviousState:      NetIplSharemStateNone,
		NetIplSharemEvents: make(chan utils.MachineEvent, 1000)}

	port.NetIplShareMachineFsm = nism

	return nism
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (nism *NetIplShareMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if nism.Machine == nil {
		nism.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	nism.Machine.Rules = r
	nism.Machine.Curr = &utils.StateEvent{
		StrStateMap: NetIplSharemStateStrMap,
		LogEna:      false,
		Logger:      nism.DrcpNetIplSharemLog,
		Owner:       NetIplShareMachineModuleStr,
	}

	return nism.Machine
}

// DrcpNetIplShareMachineNoManipulatedFramesSent While in this state, the IPL can only be supported by a
// physical or Aggregation Link
func (nism *NetIplShareMachine) DrcpNetIplShareMachineNoManipulatedFramesSent(m fsm.Machine, data interface{}) fsm.State {

	p := nism.p
	p.EnabledTimeShared = false
	p.EnabledEncTagShared = false
	return NetIplSharemStateNoManipulatedFramesSent
}

// DrcpNetIplShareMachineTimeSharedMethod While in this state, the Network / IPL sharing by time methods
// specified in 9.3.2.1 are enabled.
func (nism *NetIplShareMachine) DrcpNetIplShareMachineTimeSharedMethod(m fsm.Machine, data interface{}) fsm.State {
	p := nism.p
	p.EnabledTimeShared = true
	return NetIplSharemStateTimeShareMethod
}

// DrcpNetIplShareMachineManipulatedFramesSent While in this state, the tag manipulation methods of Network /
// IPL sharing by tag or Network / IPL sharing by encapsulation, as dictated by the Network / IPL
// sharing method selected the aDrniEncapsulationMethod (7.4.1.1.17), are enabled
func (nism *NetIplShareMachine) DrcpNetIplShareMachineManipulatedFramesSent(m fsm.Machine, data interface{}) fsm.State {
	p := nism.p
	p.EnabledEncTagShared = true
	return NetIplSharemStateManipulatedFramesSent
}

// DrcpNetIplShareMachineFSMBuild will build the State machine with callbacks
func DrcpNetIplShareMachineFSMBuild(p *DRCPIpp) *NetIplShareMachine {

	rules := fsm.Ruleset{}

	NetIplShareMachineStrStateMapCreate()

	// Instantiate a new NetIplShareMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	nism := NewDrcpNetIplShareMachine(p)

	//BEGIN -> NO MANIPULATED FRAMES SENT
	rules.AddRule(NetIplSharemStateNone, NetIplSharemEventBegin, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)
	rules.AddRule(NetIplSharemStateNoManipulatedFramesSent, NetIplSharemEventBegin, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)
	rules.AddRule(NetIplSharemStateTimeShareMethod, NetIplSharemEventBegin, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)
	rules.AddRule(NetIplSharemStateManipulatedFramesSent, NetIplSharemEventBegin, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)

	// CC TIME SHARED -> TIME SHARED METHOD
	rules.AddRule(NetIplSharemStateNoManipulatedFramesSent, NetIplSharemEventCCTimeShare, nism.DrcpNetIplShareMachineTimeSharedMethod)

	// CC ENC TAG SHARED -> MANIPULATED FRAMES SENT
	rules.AddRule(NetIplSharemStateNoManipulatedFramesSent, NetIplSharemEventCCEncTagShared, nism.DrcpNetIplShareMachineManipulatedFramesSent)

	// NO CC TIME SHARED -> NO MANIPULATED FRAMES SENT
	rules.AddRule(NetIplSharemStateTimeShareMethod, NetIplSharemEventNotCCTimeShare, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)

	// NO CC ENC TAG SHARED -> NO MANIPULATED FRAMES SENT
	rules.AddRule(NetIplSharemStateManipulatedFramesSent, NetIplSharemEventNotCCEncTagShared, nism.DrcpNetIplShareMachineNoManipulatedFramesSent)

	// Create a new FSM and apply the rules
	nism.Apply(&rules)

	return nism
}

// NetIplShareMachineMain:  802.1ax-2014 Section 9.4.20 Network/IPL sharing machine
// Creation of Network/IPL sharing State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *DRCPIpp) NetIplShareMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 9.4.19 DRCPDU Transmit machine
	nism := DrcpNetIplShareMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	nism.Machine.Start(nism.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the NetIplShareMachine should handle.
	go func(m *NetIplShareMachine) {
		m.DrcpNetIplSharemLog("Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case event, ok := <-m.NetIplSharemEvents:
				var rv error
				if ok {
					//m.LacpTxmLog(fmt.Sprintf("Event rx %d %s %s", event.E, event.Src, TxmStateStrMap[m.Machine.Curr.CurrentState()]))
					// special case, another machine has a need to
					// transmit a packet
					if event.E == TxmEventNtt {
						p.NTTDRCPDU = true
					}

					rv = m.Machine.ProcessEvent(event.Src, event.E, nil)
					if rv == nil {
						m.processPostStates()
					}

					if event.ResponseChan != nil {
						utils.SendResponse(NetIplShareMachineModuleStr, event.ResponseChan)
					}

					if rv != nil {
						m.DrcpNetIplSharemLog(strings.Join([]string{error.Error(rv), event.Src, NetIplSharemStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}
				} else {
					m.DrcpNetIplSharemLog("Machine End")
					return
				}
			}
		}
	}(nism)
}

// processPostStates will check state and local variables for additional state transitions
func (nism *NetIplShareMachine) processPostStates() {
	p := nism.p
	if nism.Machine.Curr.CurrentState() == NetIplSharemStateNoManipulatedFramesSent {
		if p.CCTimeShared {
			rv := nism.Machine.ProcessEvent(NetIplShareMachineModuleStr, NetIplSharemEventCCTimeShare, nil)
			if rv == nil {
				nism.processPostStates()
			}
		} else if p.CCEncTagShared {
			rv := nism.Machine.ProcessEvent(NetIplShareMachineModuleStr, NetIplSharemEventCCEncTagShared, nil)
			if rv == nil {
				nism.processPostStates()
			}
		}
	} else if nism.Machine.Curr.CurrentState() == NetIplSharemStateTimeShareMethod {
		if !p.CCTimeShared {
			rv := nism.Machine.ProcessEvent(NetIplShareMachineModuleStr, NetIplSharemEventCCTimeShare, nil)
			if rv == nil {
				nism.processPostStates()
			}
		}
	} else if nism.Machine.Curr.CurrentState() == NetIplSharemStateManipulatedFramesSent {
		if !p.CCEncTagShared {
			rv := nism.Machine.ProcessEvent(NetIplShareMachineModuleStr, NetIplSharemEventNotCCEncTagShared, nil)
			if rv == nil {
				nism.processPostStates()
			}
		}
	}
}
