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
// 802.1ax-2014 Section 9.4.17 DRNI Gateway and Aggregator machines
// rxmachine.go
package drcp

import (
	"utils/fsm"
	//"github.com/google/gopacket/layers"
	"l2/lacp/protocol/utils"
	//"sort"
	//"time"
	"strconv"
	"strings"
)

const AMachineModuleStr = "DRNI Aggregator Machine"

// drxm States
const (
	AmStateNone = iota + 1
	AmStateDRNIPortInitialize
	AmStateDRNIPortUpdate
	AmStatePsPortUpdate
)

var AmStateStrMap map[fsm.State]string

func AMachineStrStateMapCreate() {
	AmStateStrMap = make(map[fsm.State]string)
	AmStateStrMap[AmStateNone] = "None"
	AmStateStrMap[AmStateDRNIPortInitialize] = "DRNI Port Initialize"
	AmStateStrMap[AmStateDRNIPortUpdate] = "DRNI Port Update"
	AmStateStrMap[AmStatePsPortUpdate] = "PS Port Update"
}

// am events
const (
	AmEventBegin = iota + 1
	AmEventPortConversationUpdate
	AmEventNotIppAllPortUpdate
)

// AMachine holds FSM and current State
// and event channels for State transitions
type AMachine struct {
	ConversationIdType int

	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	dr *DistributedRelay

	// machine specific events
	AmEvents chan utils.MachineEvent
}

func (am *AMachine) PrevState() fsm.State { return am.PreviousState }

// PrevStateSet will set the previous State
func (am *AMachine) PrevStateSet(s fsm.State) { am.PreviousState = s }

// Stop should clean up all resources
func (am *AMachine) Stop() {
	close(am.AmEvents)
}

// NewDrcpAMachine will create a new instance of the AMachine
func NewDrcpAMachine(dr *DistributedRelay) *AMachine {
	am := &AMachine{
		dr:            dr,
		PreviousState: AmStateNone,
		AmEvents:      make(chan utils.MachineEvent, 50),
	}

	dr.AMachineFsm = am

	return am
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (am *AMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if am.Machine == nil {
		am.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	am.Machine.Rules = r
	am.Machine.Curr = &utils.StateEvent{
		StrStateMap: AmStateStrMap,
		LogEna:      true,
		Logger:      am.DrcpAmLog,
		Owner:       AMachineModuleStr,
	}

	return am.Machine
}

// DrcpAMachineDRNIPortInitialize function to be called after
// State transition to DRNI_PORT_INITIALIZE
func (am *AMachine) DrcpAMachineDRNIPortInitialize(m fsm.Machine, data interface{}) fsm.State {
	dr := am.dr
	am.initializeDRNIPortConversation()
	dr.PortConversationUpdate = false

	return AmStateDRNIPortInitialize
}

// DrcpAMachineDRNIPortUpdate function to be called after
// State transition to DRNI_PORT_UPDATE
func (am *AMachine) DrcpAMachineDRNIPortUpdate(m fsm.Machine, data interface{}) fsm.State {
	dr := am.dr
	dr.PortConversationUpdate = false
	am.updatePortalState()
	am.setIPPPortUpdate()
	am.setPortConversation()

	// next State
	return AmStateDRNIPortUpdate
}

// DrcpAMachinePSPortUpdate function to be called after
// State transition to PS_PORT_UPDATE
func (am *AMachine) DrcpAMachinePSPortUpdate(m fsm.Machine, data interface{}) fsm.State {
	am.updatePortalSystemPortConversation()

	// next State
	return AmStateDRNIPortUpdate
}

func DrcpAMachineFSMBuild(dr *DistributedRelay) *AMachine {

	AMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new AMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	am := NewDrcpAMachine(dr)

	//BEGIN -> DRNI PORT INITIALIZE
	rules.AddRule(AmStateNone, AmEventBegin, am.DrcpAMachineDRNIPortInitialize)
	rules.AddRule(AmStateDRNIPortUpdate, AmEventBegin, am.DrcpAMachineDRNIPortInitialize)
	rules.AddRule(AmStatePsPortUpdate, AmEventBegin, am.DrcpAMachineDRNIPortInitialize)

	// PORT CONVERSATION UPDATE  > DRNI PORT UPDATE
	rules.AddRule(AmStateDRNIPortInitialize, AmEventPortConversationUpdate, am.DrcpAMachineDRNIPortUpdate)
	rules.AddRule(AmStateDRNIPortUpdate, AmEventPortConversationUpdate, am.DrcpAMachineDRNIPortUpdate)
	rules.AddRule(AmStatePsPortUpdate, AmEventPortConversationUpdate, am.DrcpAMachineDRNIPortUpdate)

	// NOT IPP ALL GATEWAY UPDATE  > PS PORT UPDATE
	rules.AddRule(AmStateDRNIPortUpdate, AmEventNotIppAllPortUpdate, am.DrcpAMachinePSPortUpdate)

	// Create a new FSM and apply the rules
	am.Apply(&rules)

	return am
}

// DrcpAMachineMain:  802.1ax-2014 Figure 9-26
// Creation of DRNI Aggregator State Machine state transitions and callbacks
// and create go routine to pend on events
func (dr *DistributedRelay) DrcpAMachineMain() {

	// Build the State machine for  DRNI Aggregator Machine according to
	// 802.1ax-2014 Section 9.4.17 DRNI Gateway and Aggregator machines
	am := DrcpAMachineFSMBuild(dr)
	//dr.wg.Add(1)
	dr.waitgroupadd("AM")

	// set the inital State
	am.Machine.Start(am.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the AMachine should handle.
	go func(m *AMachine) {
		m.DrcpAmLog("Machine Start")
		//defer m.dr.wg.Done()
		defer m.dr.waitgroupstop("AM")
		for {
			select {
			case event, ok := <-m.AmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					dr := m.dr
					// post state processing
					if rv == nil &&
						dr.PortConversationUpdate {
						rv = m.Machine.ProcessEvent(AMachineModuleStr, AmEventPortConversationUpdate, nil)
					}
					if rv == nil &&
						m.Machine.Curr.CurrentState() == AmStateDRNIPortUpdate &&
						!m.IppAllPortUpdateCheck() {
						rv = m.Machine.ProcessEvent(AMachineModuleStr, AmEventNotIppAllPortUpdate, nil)
					}

					if rv != nil {
						m.DrcpAmLog(strings.Join([]string{error.Error(rv), event.Src, AmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(AMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.DrcpAmLog("Machine End")
					return
				}
			}
		}
	}(am)
}

// IppAllPortUpdateCheck Check is made in order to try and perform the following logic
// needed for IppAllPortUpdate; This variable is the logical OR of the IppPortUpdate
// variables for all IPPs in this Portal System.
func (am *AMachine) IppAllPortUpdateCheck() bool {
	dr := am.dr
	for _, ippid := range dr.DrniIntraPortalLinkList {
		for _, ipp := range DRCPIppDBList {
			if ipp.Id == ippid&0xffff {
				if ipp.IppPortUpdate {
					dr.IppAllPortUpdate = true
				}
			}
		}
	}
	return dr.IppAllPortUpdate
}

// initializeDRNIPortConversation This function sets the Drni_Portal_System_Port_Conversation to a sequence of zeros, indexed
// by Port Conversation ID.
func (am *AMachine) initializeDRNIPortConversation() {
	dr := am.dr

	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		dr.DrniPortalSystemPortConversation[i] = false
	}
}

// updatePortalState This function updates the Drni_Portal_System_State[] as follows
func (am *AMachine) updatePortalState() {
	dr := am.dr

	dr.updatePortalState(AMachineModuleStr)
}

// setIPPPortUpdate This function sets the IppPortUpdate on every IPP on this Portal System to TRUE.
func (am *AMachine) setIPPPortUpdate() {
	dr := am.dr
	for _, ippid := range dr.DrniIntraPortalLinkList {
		for _, ipp := range DRCPIppDBList {
			if ipp.Id == ippid {
				ipp.IppPortUpdate = true
			}
		}
	}
}

// setPortConversation This function sets Drni_Port_Conversation to the values computed from
// Conversation_PortList[] and the current Drni_Portal_System_State[] as follows
func (am *AMachine) setPortConversation() {
	//dr := am.dr
	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		// TODO
		// For every indexed Gateway Conversation ID, a Portal System Number is identified by
		// choosing the highest priority Portal System Number in the list of Portal System Numbers
		// provided by aDrniConvAdminGateway[] when only the Portal Systems having that Gateway
		// Conversation ID enabled in the Gateway Vectors of the Drni_Portal_System_State[] variable,
		// are included.
		//dr.a.ConversationPortList
	}
}

// updatePortalSystemPortConversation This function sets Drni_Portal_System_Port_Conversation as follows
func (am *AMachine) updatePortalSystemPortConversation() {
	dr := am.dr

	for _, ipp := range dr.Ipplinks {
		if ipp.DifferPortDigest &&
			!dr.DrniThreeSystemPortal {
			for i := 0; i < 512; i++ {
				dr.DrniPortalSystemPortConversation[i] = ipp.DrniNeighborPortConversation[i]>>7&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+1] = ipp.DrniNeighborPortConversation[i]>>6&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+2] = ipp.DrniNeighborPortConversation[i]>>5&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+3] = ipp.DrniNeighborPortConversation[i]>>4&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+4] = ipp.DrniNeighborPortConversation[i]>>3&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+5] = ipp.DrniNeighborPortConversation[i]>>2&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+6] = ipp.DrniNeighborPortConversation[i]>>1&0x1 == 1
				dr.DrniPortalSystemPortConversation[i+7] = ipp.DrniNeighborPortConversation[i]>>0&0x1 == 1
			}
		} else {

			// TODO not sure this a applicable as OTHER should always be false
			// This function sets the Drni_Portal_System_Port_Conversation to the result of the logical
			// AND operation between, the Boolean vector constructed from the
			// Drni_Port_Conversation, by setting to FALSE all the indexed Port Conversation ID
			// entries that are associated with other Portal Systems in the Portal, and the Boolean vector
			// constructed from the Ipp_Other_Port_Conversation_Portal_System, by setting to FALSE
			// all the indexed Port Conversation ID entries that are associated with other Portal Systems
			// in the Portal.
		}
	}
}
