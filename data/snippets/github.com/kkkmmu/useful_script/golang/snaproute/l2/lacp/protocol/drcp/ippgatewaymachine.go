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
// ippgatewaymachine.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"strconv"
	"strings"
	"utils/fsm"
)

const IGMachineModuleStr = "IPP Gateway Machine"

// igm States
const (
	IGmStateNone = iota + 1
	IGmStateIPPGatewayInitialize
	IGmStateIPPGatewayUpdate
)

var IGmStateStrMap map[fsm.State]string

func IGMachineStrStateMapCreate() {
	IGmStateStrMap = make(map[fsm.State]string)
	IGmStateStrMap[IGmStateNone] = "None"
	IGmStateStrMap[IGmStateIPPGatewayInitialize] = "IPP Gateway Initialize"
	IGmStateStrMap[IGmStateIPPGatewayUpdate] = "IPP Gateway Update"
}

// igm events
const (
	IGmEventBegin = iota + 1
	IGmEventGatewayUpdate
)

// IGMachine holds FSM and current State
// and event channels for State transitions
type IGMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	p *DRCPIpp

	// machine specific events
	IGmEvents chan utils.MachineEvent
}

func (igm *IGMachine) PrevState() fsm.State { return igm.PreviousState }

// PrevStateSet will set the previous State
func (igm *IGMachine) PrevStateSet(s fsm.State) { igm.PreviousState = s }

// Stop should clean up all resources
func (igm *IGMachine) Stop() {
	close(igm.IGmEvents)
}

// NewDrcpIGMachine will create a new instance of the IGMachine
func NewDrcpIGMachine(p *DRCPIpp) *IGMachine {
	igm := &IGMachine{
		p:             p,
		PreviousState: IGmStateNone,
		IGmEvents:     make(chan utils.MachineEvent, 10),
	}

	p.IGMachineFsm = igm

	return igm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (igm *IGMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if igm.Machine == nil {
		igm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	igm.Machine.Rules = r
	igm.Machine.Curr = &utils.StateEvent{
		StrStateMap: IGmStateStrMap,
		LogEna:      true,
		Logger:      igm.DrcpIGmLog,
		Owner:       IGMachineModuleStr,
	}

	return igm.Machine
}

// DrcpGMachineIPPGatewayInitialize function to be called after
// State transition to IPP_GATEWAY_INITIALIZE
func (igm *IGMachine) DrcpIGMachineIPPGatewayInitialize(m fsm.Machine, data interface{}) fsm.State {
	p := igm.p
	igm.initializeIPPGatewayConversation()
	p.IppGatewayUpdate = false

	defer igm.NotifyIppAllGatewayUpdate()
	return IGmStateIPPGatewayInitialize
}

// DrcpIGMachineIPPGatewayUpdate function to be called after
// State transition to IPP_GATEWAY_UPDATE
func (igm *IGMachine) DrcpIGMachineIPPGatewayUpdate(m fsm.Machine, data interface{}) fsm.State {
	p := igm.p

	igm.setIPPGatewayConversation()
	igm.updateIPPGatewayConversationDirection()
	p.IppGatewayUpdate = false

	defer igm.NotifyIppAllGatewayUpdate()
	// next State
	return IGmStateIPPGatewayUpdate
}

func DrcpIGMachineFSMBuild(p *DRCPIpp) *IGMachine {

	IGMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new IGMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	igm := NewDrcpIGMachine(p)

	//BEGIN -> IPP GATEWAY INITIALIZE
	rules.AddRule(IGmStateNone, IGmEventBegin, igm.DrcpIGMachineIPPGatewayInitialize)
	rules.AddRule(IGmStateIPPGatewayInitialize, IGmEventBegin, igm.DrcpIGMachineIPPGatewayInitialize)
	rules.AddRule(IGmStateIPPGatewayUpdate, IGmEventBegin, igm.DrcpIGMachineIPPGatewayInitialize)

	// IPP GATEWAY UPDATE  > IPP GATEWAY UPDATE
	rules.AddRule(IGmStateIPPGatewayInitialize, IGmEventGatewayUpdate, igm.DrcpIGMachineIPPGatewayUpdate)
	rules.AddRule(IGmStateIPPGatewayUpdate, IGmEventGatewayUpdate, igm.DrcpIGMachineIPPGatewayUpdate)

	// Create a new FSM and apply the rules
	igm.Apply(&rules)

	return igm
}

// DrcpIGMachineMain:  802.1ax-2014 Figure 9-27
// Creation of DRNI IPP Gateway State Machine state transitions and callbacks
// and create go routine to pend on events
func (p *DRCPIpp) DrcpIGMachineMain() {

	// Build the State machine for  DRNI Gateway Machine according to
	// 802.1ax-2014 Section 9.4.18 DRNI IPP machines
	igm := DrcpIGMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	igm.Machine.Start(igm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the IGMachine should handle.
	go func(m *IGMachine) {
		m.DrcpIGmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case event, ok := <-m.IGmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					p := m.p
					// post state processing
					if rv == nil &&
						p.IppGatewayUpdate {
						rv = m.Machine.ProcessEvent(IGMachineModuleStr, IGmEventGatewayUpdate, nil)
					}

					if rv != nil {
						m.DrcpIGmLog(strings.Join([]string{error.Error(rv), event.Src, IGmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(IGMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.DrcpIGmLog("Machine End")
					return
				}
			}
		}
	}(igm)
}

// initializeIPPGatewayConversation This function sets the Ipp_Gateway_Conversation_Direction to a sequence of zeros, indexed
// by Gateway Conversation ID.
func (igm *IGMachine) initializeIPPGatewayConversation() {
	p := igm.p
	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		p.GatewayConversationDirection[i] = false
	}
}

// setIPPGatewayConversation This function sets Ipp_Other_Gateway_Conversation as follows
func (igm *IGMachine) setIPPGatewayConversation() {
	p := igm.p
	dr := p.dr

	if p.DifferGatewayDigest &&
		dr.DrniThreeSystemPortal {
		// TODO handle 3 portal system logic
	} else if p.DifferGatewayDigest &&
		!dr.DrniThreeSystemPortal {
		for i, j := 0, 0; i < 512; i, j = i+1, j+8 {
			if (p.DrniNeighborGatewayConversation[i] >> 7 & 1) == 0 {
				p.IppOtherGatewayConversation[j] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 6 & 1) == 0 {
				p.IppOtherGatewayConversation[j+1] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+1] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 5 & 1) == 0 {
				p.IppOtherGatewayConversation[j+2] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+2] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 4 & 1) == 0 {
				p.IppOtherGatewayConversation[j+3] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+3] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 3 & 1) == 0 {
				p.IppOtherGatewayConversation[j+4] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+4] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 2 & 1) == 0 {
				p.IppOtherGatewayConversation[j+5] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+5] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 1 & 1) == 0 {
				p.IppOtherGatewayConversation[j+6] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+6] = p.DRFHomeConfNeighborPortalSystemNumber
			}
			if (p.DrniNeighborGatewayConversation[i] >> 0 & 1) == 0 {
				p.IppOtherGatewayConversation[j+7] = dr.DRFPortalSystemNumber
			} else {
				p.IppOtherGatewayConversation[j+7] = p.DRFHomeConfNeighborPortalSystemNumber
			}
		}
	} else if !p.DifferGatewayDigest {
		// TODO revisit logic
		// This function sets Ipp_Other_Gateway_Conversation to the values computed from
		// aDrniConvAdminGateway[] and the Drni_Neighbor_State[] as follows:
		// For every indexed Gateway Conversation ID, a Portal System Number is identified by
		// choosing the highest priority Portal System Number in the list of Portal System Numbers
		// provided by aDrniConvAdminGateway[] when only the Portal Systems having that
		// Gateway Conversation ID enabled in the Gateway Vectors of the Drni_Neighbor_State[]
		// variable, are included
		for cid, portalsystemnumbers := range dr.DrniConvAdminGateway {
			if portalsystemnumbers != nil {
				for _, portalsystemnumber := range portalsystemnumbers {
					p.DrniNeighborState[portalsystemnumber].mutex.Lock()
					if portalsystemnumber != 0 &&
						p.DrniNeighborState[portalsystemnumber].OpState &&
						p.DrniNeighborState[portalsystemnumber].GatewayVector != nil &&
						p.DrniNeighborState[portalsystemnumber].GatewayVector[0].Vector[cid] {
						p.IppOtherGatewayConversation[cid] = portalsystemnumber
					}
					p.DrniNeighborState[portalsystemnumber].mutex.Unlock()
				}
			}
		}
	}
}

// updateIPPGatewayConversationDirection This function computes a value for Ipp_Gateway_Conversation_Direction as follows
func (igm *IGMachine) updateIPPGatewayConversationDirection() {
	p := igm.p
	dr := p.dr

	// NOTE: logic below is not checking the following case because there is no other gateway
	// in our supported implementation
	//b) Drni_Gateway_Conversation and Ipp_Other_Gateway_Conversation are in agreement
	//   as to which Portal System should get this Gateway Conversation ID.
	//   In addition, if Drni_Gateway_Conversation and Ipp_Other_Gateway_Conversation are in
	//   disagreement for any Gateway Conversation ID:
	//   It sets DRF_HomIe_Oper_DRCP_State.Gateway_Sync to FALSE, and;
	//   NTTDRCPDU to TRUE.
	if !dr.DrniThreeSystemPortal {
		for conid := 0; conid < MAX_CONVERSATION_IDS; conid++ {

			if (dr.DrniGatewayConversation[conid] != nil && !p.IppGatewayConversationPasses[conid]) ||
				(dr.DrniGatewayConversation[conid] == nil && p.IppGatewayConversationPasses[conid]) {
				// lets only check the first portal as this indicates
				// which system according to the gateway conversation
				// owns this conversation, however it should be noted
				// that in the case of sharing by time both systems
				// will own a conversation
				if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
					igm.DrcpIGmLog(fmt.Sprintf("Made it here Conversation Id %d ipp port %d\n", conid, p.Id))

					for _, statevector := range p.IppPortalSystemState {
						statevector.mutex.Lock()
						if statevector.OpState {

							if statevector.GatewayVector != nil {
								seqvector := statevector.GatewayVector[0]
								if seqvector.Vector != nil {
									if seqvector.Vector[conid] &&
										!p.IppGatewayConversationPasses[conid] {
										p.IppGatewayConversationPasses[conid] = true
										if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
											// TODO  Set the Vlan Membership for each conversation that is valid
											igm.DrcpIGmLog(fmt.Sprintf("Setting Vlan Membership for Conversation Id %d ipp port %d\n", conid, p.Id))
											for _, client := range utils.GetAsicDPluginList() {
												err := client.IppVlanConversationSet(uint16(conid), int32(p.Id))
												if err != nil {
													igm.DrcpIGmLog(fmt.Sprintf("ERROR setting Vlan membership %v", err))
												}
											}
										}
									} else if !seqvector.Vector[conid] {
										if p.IppGatewayConversationPasses[conid] &&
											dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
											// TODO  Set the Vlan Membership for each conversation that is valid
											igm.DrcpIGmLog(fmt.Sprintf("Clearing Vlan Membership for Conversation Id %d ipp port %d\n", conid, p.Id))
											for _, client := range utils.GetAsicDPluginList() {
												err := client.IppVlanConversationClear(uint16(conid), int32(p.Id))
												if err != nil {
													igm.DrcpIGmLog(fmt.Sprintf("ERROR clearing Vlan membership %v", err))
												}
											}
										}
										p.IppGatewayConversationPasses[conid] = false
									}
								}
							}
						}
						statevector.mutex.Unlock()

					}
				}
			}
		}
		// TODO
		// In addition, if Drni_Gateway_Conversation and Ipp_Other_Gateway_Conversation are in
		// disagreement for any Gateway Conversation ID:
		// It sets DRF_Home_Oper_DRCP_State.Gateway_Sync to FALSE, and;
		// NTTDRCPDU to TRUE.
		// Otherwise:
		// DRF_Home_Oper_DRCP_State.Gateway_Sync and NTTDRCPDU are left unchanged.
	}
}

// NotifyIppAllGatewayUpdate this should be called each time IppGatewayUpdate is changed
// to false so that the gateway machine can be informed
func (igm *IGMachine) NotifyIppAllGatewayUpdate() {
	p := igm.p
	dr := p.dr

	allgatewayupdate := false
	for _, ipp := range dr.Ipplinks {
		if ipp.IppGatewayUpdate {
			allgatewayupdate = true
		}
	}
	if !allgatewayupdate && dr.IppAllGatewayUpdate {
		dr.IppAllGatewayUpdate = false
		if dr.GMachineFsm != nil {
			dr.GMachineFsm.GmEvents <- utils.MachineEvent{
				E:   GmEventNotIppAllGatewayUpdate,
				Src: IGMachineModuleStr,
			}
		}
	}
}
