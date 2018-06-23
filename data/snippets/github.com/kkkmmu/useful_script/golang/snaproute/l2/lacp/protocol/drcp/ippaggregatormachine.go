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
// 802.1ax-2014 Section 9.4.18 DRNI IPP machines
// rxmachine.go
package drcp

import (
	"l2/lacp/protocol/utils"
	"strconv"
	"strings"
	"utils/fsm"
)

const IAMachineModuleStr = "DRNI Aggregator Machine"

// drxm States
const (
	IAmStateNone = iota + 1
	IAmStateIPPPortInitialize
	IAmStateIPPPortUpdate
)

var IAmStateStrMap map[fsm.State]string

func IAMachineStrStateMapCreate() {
	IAmStateStrMap = make(map[fsm.State]string)
	IAmStateStrMap[IAmStateNone] = "None"
	IAmStateStrMap[IAmStateIPPPortInitialize] = "IPP Port Initialize"
	IAmStateStrMap[IAmStateIPPPortUpdate] = "IPP Port Update"
}

// am events
const (
	IAmEventBegin = iota + 1
	IAmEventIPPPortUpdate
)

// IAMachine holds FSM and current State
// and event channels for State transitions
type IAMachine struct {
	ConversationIdType int

	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	p *DRCPIpp

	// machine specific events
	IAmEvents chan utils.MachineEvent
}

func (iam *IAMachine) PrevState() fsm.State { return iam.PreviousState }

// PrevStateSet will set the previous State
func (iam *IAMachine) PrevStateSet(s fsm.State) { iam.PreviousState = s }

// Stop should clean up all resources
func (iam *IAMachine) Stop() {
	close(iam.IAmEvents)
}

// NewDrcpAMachine will create a new instance of the IAMachine
func NewDrcpIAMachine(p *DRCPIpp) *IAMachine {
	iam := &IAMachine{
		p:             p,
		PreviousState: AmStateNone,
		IAmEvents:     make(chan utils.MachineEvent, 10),
	}

	p.IAMachineFsm = iam

	return iam
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (iam *IAMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if iam.Machine == nil {
		iam.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	iam.Machine.Rules = r
	iam.Machine.Curr = &utils.StateEvent{
		StrStateMap: IAmStateStrMap,
		LogEna:      true,
		Logger:      iam.DrcpIAmLog,
		Owner:       IAMachineModuleStr,
	}

	return iam.Machine
}

// DrcpIAMachineIPPPortInitialize function to be called after
// State transition to IPP_PORT_INITIALIZE
func (iam *IAMachine) DrcpIAMachineIPPPortInitialize(m fsm.Machine, data interface{}) fsm.State {
	p := iam.p
	iam.initializeIPPPortConversation()
	p.IppPortUpdate = false

	return IAmStateIPPPortInitialize
}

// DrcpIAMachineIPPPortUpdate function to be called after
// State transition to IPP_PORT_UPDATE
func (iam *IAMachine) DrcpIAMachineIPPPortUpdate(m fsm.Machine, data interface{}) fsm.State {
	p := iam.p

	iam.setIPPPortConversation()
	iam.updateIPPPortConversationPasses()
	p.IppPortUpdate = false

	// next State
	return IAmStateIPPPortUpdate
}

func DrcpIAMachineFSMBuild(p *DRCPIpp) *IAMachine {

	IAMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new IAMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	iam := NewDrcpIAMachine(p)

	//BEGIN -> IPP PORT INITIALIZE
	rules.AddRule(IAmStateNone, IAmEventBegin, iam.DrcpIAMachineIPPPortInitialize)
	rules.AddRule(IAmStateIPPPortInitialize, IAmEventBegin, iam.DrcpIAMachineIPPPortInitialize)
	rules.AddRule(IAmStateIPPPortUpdate, IAmEventBegin, iam.DrcpIAMachineIPPPortInitialize)

	// IPP PORT UPDATE  > IPP PORT UPDATE
	rules.AddRule(IAmStateIPPPortInitialize, IAmEventIPPPortUpdate, iam.DrcpIAMachineIPPPortUpdate)
	rules.AddRule(IAmStateIPPPortUpdate, IAmEventIPPPortUpdate, iam.DrcpIAMachineIPPPortUpdate)

	// Create a new FSM and apply the rules
	iam.Apply(&rules)

	return iam
}

// DrcpIAMachineMain:  802.1ax-2014 Figure 9-27
// Creation of DRNI IPP Machines state transitions and callbacks
// and create go routine to pend on events
func (p *DRCPIpp) DrcpIAMachineMain() {

	// Build the State machine for  DRNI IPP Machine according to
	// 802.1ax-2014 Section 9.4.18 DRNI IPP machines
	iam := DrcpIAMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	iam.Machine.Start(iam.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the AMachine should handle.
	go func(m *IAMachine) {
		m.DrcpIAmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case event, ok := <-m.IAmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					p := m.p
					// post state processing
					if rv == nil &&
						p.IppPortUpdate {
						rv = m.Machine.ProcessEvent(IAMachineModuleStr, IAmEventIPPPortUpdate, nil)
					}

					if rv != nil {
						m.DrcpIAmLog(strings.Join([]string{error.Error(rv), event.Src, IAmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(IAMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.DrcpIAmLog("Machine End")
					return
				}
			}
		}
	}(iam)
}

// initializeIPPPortConversation This function sets the Ipp_Port_Conversation_Passes to a sequence of zeros, indexed by Port
// Conversation ID.
func (iam *IAMachine) initializeIPPPortConversation() {
	p := iam.p

	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		p.PortConversationPasses[i] = false
	}
}

// setIPPPortConversation This function sets Ipp_Other_Port_Conversation_Portal_System as follows
func (iam *IAMachine) setIPPPortConversation() {
	p := iam.p

	if p.DifferPortDigest &&
		p.dr.DrniThreeSystemPortal {
		// TODO when a 3P system is supported
		//for i, j := 0, 512; i < MAX_CONVERSATION_IDS; i, j = i+8, i+1 {
		//}
	} else if p.DifferPortDigest &&
		!p.dr.DrniThreeSystemPortal {
		var neighborConversationSystemNumbers [1024]uint8
		for i, j := 0, 0; i < MAX_CONVERSATION_IDS; i, j = i+8, j+2 {
			if (p.DrniNeighborPortConversation[i] >> 7 & 1) == 0 {
				neighborConversationSystemNumbers[j] = p.dr.DRFPortalSystemNumber << 6
			} else {
				neighborConversationSystemNumbers[j] = p.DRFHomeConfNeighborPortalSystemNumber << 6
			}
			if (p.DrniNeighborPortConversation[i] >> 6 & 1) == 0 {
				neighborConversationSystemNumbers[j] = p.dr.DRFPortalSystemNumber << 4
			} else {
				neighborConversationSystemNumbers[j] = p.DRFHomeConfNeighborPortalSystemNumber << 4
			}
			if (p.DrniNeighborPortConversation[i] >> 5 & 1) == 0 {
				neighborConversationSystemNumbers[j] = p.dr.DRFPortalSystemNumber << 2
			} else {
				neighborConversationSystemNumbers[j] = p.DRFHomeConfNeighborPortalSystemNumber << 2
			}
			if (p.DrniNeighborPortConversation[i] >> 4 & 1) == 0 {
				neighborConversationSystemNumbers[j] = p.dr.DRFPortalSystemNumber << 0
			} else {
				neighborConversationSystemNumbers[j] = p.DRFHomeConfNeighborPortalSystemNumber << 0
			}
			if (p.DrniNeighborPortConversation[i] >> 3 & 1) == 0 {
				neighborConversationSystemNumbers[j+1] = p.dr.DRFPortalSystemNumber << 6
			} else {
				neighborConversationSystemNumbers[j+1] = p.DRFHomeConfNeighborPortalSystemNumber << 6
			}
			if (p.DrniNeighborPortConversation[i] >> 2 & 1) == 0 {
				neighborConversationSystemNumbers[j+1] = p.dr.DRFPortalSystemNumber << 4
			} else {
				neighborConversationSystemNumbers[j+1] = p.DRFHomeConfNeighborPortalSystemNumber << 4
			}
			if (p.DrniNeighborPortConversation[i] >> 1 & 1) == 0 {
				neighborConversationSystemNumbers[j+1] = p.dr.DRFPortalSystemNumber << 2
			} else {
				neighborConversationSystemNumbers[j+1] = p.DRFHomeConfNeighborPortalSystemNumber << 2
			}
			if (p.DrniNeighborPortConversation[i] >> 0 & 1) == 0 {
				neighborConversationSystemNumbers[j+1] = p.dr.DRFPortalSystemNumber << 0
			} else {
				neighborConversationSystemNumbers[j+1] = p.DRFHomeConfNeighborPortalSystemNumber << 0
			}
		}
		// now lets save the values
		p.DrniNeighborPortConversation = neighborConversationSystemNumbers
	} else if !p.DifferPortDigest {
		// TODO Revisit logic here DrniNeighborState
		//
		// This function sets Ipp_Other_Port_Conversation_Portal_System to the values computed from
		// Conversation_PortList[] and the Drni_Neighbor_State[] as follows:
		// For every indexed Port Conversation ID, a Portal System Number is identified by choosing the
		// highest priority Portal System Number in the list of Portal Systems Numbers provided by
		// Conversation_PortList[] when only the operational Aggregation Ports, as provided by the
		//associated Lists of the Drni_Neighbor_State[] variable, are included.
	}
}

func (iam IAMachine) updateIPPPortConversationPasses() {

}
