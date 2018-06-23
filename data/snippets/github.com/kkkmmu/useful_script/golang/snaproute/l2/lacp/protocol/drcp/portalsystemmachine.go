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
// 802.1ax-2014 Section 9.4.15 DRCPDU Periodic Transmission machine
// rxmachine.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"strconv"
	"strings"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

const PsMachineModuleStr = "Portal System Machine"

// psm States
const (
	PsmStateNone = iota + 1
	PsmStatePortalSystemInitialize
	PsmStatePortalSystemUpdate
)

var PsmStateStrMap map[fsm.State]string

func PsMachineStrStateMapCreate() {
	PsmStateStrMap = make(map[fsm.State]string)
	PsmStateStrMap[PsmStateNone] = "None"
	PsmStateStrMap[PsmStatePortalSystemInitialize] = "Portal System Initialize"
	PsmStateStrMap[PsmStatePortalSystemUpdate] = "Portal System Update"
}

// psm events
const (
	PsmEventBegin = iota + 1
	PsmEventChangePortal
	PsmEventChangeDRFPorts
)

// PsMachine holds FSM and current State
// and event channels for State transitions
type PsMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	dr *DistributedRelay

	// machine specific events
	PsmEvents chan utils.MachineEvent
}

// PrevState will get the previous State
func (psm *PsMachine) PrevState() fsm.State { return psm.PreviousState }

// PrevStateSet will set the previous State
func (psm *PsMachine) PrevStateSet(s fsm.State) { psm.PreviousState = s }

// Stop should clean up all resources
func (psm *PsMachine) Stop() {
	close(psm.PsmEvents)
}

// NewDrcpPsMachine will create a new instance of the PsMachine
func NewDrcpPsMachine(dr *DistributedRelay) *PsMachine {
	psm := &PsMachine{
		dr:            dr,
		PreviousState: PsmStateNone,
		PsmEvents:     make(chan utils.MachineEvent, 10),
	}

	dr.PsMachineFsm = psm

	return psm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (psm *PsMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if psm.Machine == nil {
		psm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	psm.Machine.Rules = r
	psm.Machine.Curr = &utils.StateEvent{
		StrStateMap: PsmStateStrMap,
		LogEna:      true,
		Logger:      psm.DrcpPsmLog,
		Owner:       PsMachineModuleStr,
	}

	return psm.Machine
}

// DrcpPsMachinePortalSystemInitialize function to be called after
// State transition to PORTAL_SYSTEM_INITIALIZE
func (psm *PsMachine) DrcpPsMachinePortalSystemInitialize(m fsm.Machine, data interface{}) fsm.State {
	psm.setDefaultPortalSystemParameters()
	psm.updateKey()
	return PsmStatePortalSystemInitialize
}

// DrcpPsMachineFastPeriodic function to be called after
// State transition to FAST_PERIODIC
func (psm *PsMachine) DrcpPsMachinePortalSystemUpdate(m fsm.Machine, data interface{}) fsm.State {
	dr := psm.dr

	psm.updateDRFHomeState(dr.ChangePortal, dr.ChangeDRFPorts)
	psm.updateKey()
	dr.ChangePortal = false
	dr.ChangeDRFPorts = false

	// next State
	return PsmStatePortalSystemUpdate
}

func DrcpPsMachineFSMBuild(dr *DistributedRelay) *PsMachine {

	PsMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new PsMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	psm := NewDrcpPsMachine(dr)

	//BEGIN -> PORTAL SYSTEM INITITIALIZE
	rules.AddRule(PsmStateNone, PsmEventBegin, psm.DrcpPsMachinePortalSystemInitialize)
	rules.AddRule(PsmStatePortalSystemInitialize, PsmEventBegin, psm.DrcpPsMachinePortalSystemInitialize)
	rules.AddRule(PsmStatePortalSystemUpdate, PsmEventBegin, psm.DrcpPsMachinePortalSystemInitialize)

	// CHANGE PORTAL  > PORTAL SYSTEM UPDATE
	rules.AddRule(PsmStatePortalSystemInitialize, PsmEventChangePortal, psm.DrcpPsMachinePortalSystemUpdate)
	rules.AddRule(PsmStatePortalSystemUpdate, PsmEventChangePortal, psm.DrcpPsMachinePortalSystemUpdate)

	// CHANGE DRF PORTS  > PORTAL SYSTEM UPDATE
	rules.AddRule(PsmStatePortalSystemInitialize, PsmEventChangeDRFPorts, psm.DrcpPsMachinePortalSystemUpdate)
	rules.AddRule(PsmStatePortalSystemUpdate, PsmEventChangeDRFPorts, psm.DrcpPsMachinePortalSystemUpdate)

	// Create a new FSM and apply the rules
	psm.Apply(&rules)

	return psm
}

// DrcpPsMachineMain:  802.1ax-2014 Figure 9-25
// Creation of Portal System Machine State transitions and callbacks
// and create go routine to pend on events
func (dr *DistributedRelay) DrcpPsMachineMain() {

	// Build the State machine for  DRCP Portal System Machine according to
	// 802.1ax-2014 Section 9.4.16 Portal System Machine
	psm := DrcpPsMachineFSMBuild(dr)
	//dr.wg.Add(1)
	dr.waitgroupadd("PSM")

	// set the inital State
	psm.Machine.Start(psm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the PsMachine should handle.
	go func(m *PsMachine) {
		m.DrcpPsmLog("Machine Start")
		//defer m.dr.wg.Done()
		defer m.dr.waitgroupstop("PSM")
		for {
			select {
			case event, ok := <-m.PsmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					if rv == nil {
						dr := m.dr
						// port state processing
						if dr.ChangePortal {
							rv = m.Machine.ProcessEvent(PsMachineModuleStr, PsmEventChangePortal, nil)
						} else if dr.ChangeDRFPorts {
							rv = m.Machine.ProcessEvent(PsMachineModuleStr, PsmEventChangeDRFPorts, nil)
						}
					}

					if rv != nil {
						m.DrcpPsmLog(strings.Join([]string{error.Error(rv), event.Src, PsmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(PsMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.DrcpPsmLog("Machine End")
					return
				}
			}
		}
	}(psm)
}

// setDefaultPortalSystemParameters This function sets this Portal System’s variables
// to administrative set values as follows
func (psm PsMachine) setDefaultPortalSystemParameters() {
	dr := psm.dr
	a := dr.a

	//fmt.Println("AggPriorityDefault:", a.AggPriority)
	//fmt.Println("AggIdDefault:", a.AggMacAddr)
	//fmt.Println("AdminKey:", a.ActorAdminKey)
	//fmt.Println("PortAlgorithm:", a.PortAlgorithm)

	dr.DrniAggregatorPriority = a.AggPriority
	dr.DrniAggregatorId = a.AggMacAddr
	dr.DrniPortalPriority = dr.DrniPortalPriority
	dr.DrniPortalAddr = dr.DrniPortalAddr
	dr.DRFPortalSystemNumber = dr.DrniPortalSystemNumber
	dr.DRFHomeAdminAggregatorKey = a.ActorAdminKey
	dr.DRFHomePortAlgorithm = a.PortAlgorithm
	dr.DRFHomeGatewayAlgorithm = dr.DrniGatewayAlgorithm
	dr.DRFHomeOperDRCPState = dr.DRFNeighborAdminDRCPState

	// set during config do not want to clear this to a default because
	// there is only one default
	// dr.Ipplinks = nil
	for i := 1; i <= MAX_PORTAL_SYSTEM_IDS; i++ {
		dr.DrniPortalSystemState[i].mutex.Lock()
		psm.DrcpPsmLog(fmt.Sprintf("DrniPortalSystemState[%d] setting OpState to false setting default params", i))
		dr.DrniPortalSystemState[i].OpState = false
		dr.DrniPortalSystemState[i].GatewayVector = nil
		dr.DrniPortalSystemState[i].PortIdList = nil
		dr.DrniPortalSystemState[i].mutex.Unlock()
		// Don't want to clear this as it is set by config
		//dr.DrniIntraPortalLinkList[i] = 0
	}

	// Because we only support Sharing by Time we don't
	// need to worry about how the other system is
	// provisioned thus we only need to worry about filling
	// in the priority list for our system.
	//
	// This should be updated with the valid 'conversation ids'
	//
	//dr.DRFHomeConversationPortListDigest
	dr.SetTimeSharingPortAndGatwewayDigest()
}

//updateKey This function updates the operational Aggregator Key,
//DRF_Home_Oper_Aggregator_Key, as follows
func (psm *PsMachine) updateKey() {
	dr := psm.dr
	a := dr.a
	//fmt.Println("updateKey: calling method")
	if a != nil &&
		a.PartnerDWC {
		dr.DRFHomeOperAggregatorKey = ((dr.DRFHomeAdminAggregatorKey & 0x3fff) | 0x6000)
	} else if dr.DrniPSI &&
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) {
		dr.DRFHomeOperAggregatorKey = (dr.DRFHomeAdminAggregatorKey & 0x3fff)
	} else {
		// lets get the lowest admin key and set this to the
		// operational key
		var operKey uint16
		for _, ipp := range dr.Ipplinks {

			if dr.DRFHomeAdminAggregatorKey != 0 &&
				(dr.DRFHomeAdminAggregatorKey <= ipp.DRFNeighborAdminAggregatorKey || ipp.DRFNeighborAdminAggregatorKey == 0) &&
				(dr.DRFHomeAdminAggregatorKey <= ipp.DRFOtherNeighborAdminAggregatorKey || ipp.DRFOtherNeighborAdminAggregatorKey == 0) {
				if operKey == 0 || dr.DRFHomeAdminAggregatorKey < operKey {
					operKey = dr.DRFHomeAdminAggregatorKey
				}
			} else if ipp.DRFNeighborAdminAggregatorKey != 0 &&
				(ipp.DRFNeighborAdminAggregatorKey <= dr.DRFHomeAdminAggregatorKey || dr.DRFHomeAdminAggregatorKey == 0) &&
				(ipp.DRFNeighborAdminAggregatorKey <= ipp.DRFOtherNeighborAdminAggregatorKey || ipp.DRFOtherNeighborAdminAggregatorKey == 0) {
				if operKey == 0 || ipp.DRFNeighborAdminAggregatorKey < operKey {
					operKey = ipp.DRFNeighborAdminAggregatorKey
				}
			} else if ipp.DRFOtherNeighborAdminAggregatorKey != 0 &&
				(ipp.DRFOtherNeighborAdminAggregatorKey <= dr.DRFHomeAdminAggregatorKey || dr.DRFHomeAdminAggregatorKey == 0) &&
				(ipp.DRFOtherNeighborAdminAggregatorKey <= ipp.DRFNeighborAdminAggregatorKey || ipp.DRFNeighborAdminAggregatorKey == 0) {
				if operKey == 0 || ipp.DRFOtherNeighborAdminAggregatorKey <= operKey {
					operKey = ipp.DRFOtherNeighborAdminAggregatorKey
				}
			}

			dr.DRFHomeOperAggregatorKey = operKey

			// oper key has been successfully been negotiated because the
			// neighbor
			if a.ActorOperKey != operKey {
				psm.DrcpPsmLog(fmt.Sprintf("updateKey: Admin Aggregator Key is updated from %d to %d source(home[%d]neighbor[%d]other[%d] updateing ports %+v)",
					dr.DRFHomeAdminAggregatorKey, operKey, dr.DRFHomeAdminAggregatorKey, ipp.DRFNeighborAdminAggregatorKey, ipp.DRFOtherNeighborAdminAggregatorKey, a.PortNumList))

				// required so that the checkForSelection succeeds in finding the
				// new port info as being set below
				a.ActorOperKey = operKey

				for _, aggport := range a.PortNumList {
					if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
						for _, client := range utils.GetAsicDPluginList() {
							for _, ippid := range dr.DrniIntraPortalLinkList {
								inport := ippid & 0xffff
								if inport > 0 {
									dr.LaDrLog(fmt.Sprintf("updateKey: Blocking IPP %d to AggPort %d", inport, aggport))
									/* TEMP - add actual port names */
									err := client.IppIngressEgressDrop("fpPort1", "fpPort2")
									if err != nil {
										dr.LaDrLog(fmt.Sprintf("ERROR (updateKey) setting Block from %s tolag port %s", utils.GetNameFromIfIndex(int32(inport)), int32(aggport)))
									}
								}
							}
						}
					}
					lacp.SetLaAggPortSystemInfoFromDistributedRelay(
						uint16(aggport),
						fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
							dr.DrniPortalAddr[0],
							dr.DrniPortalAddr[1],
							dr.DrniPortalAddr[2],
							dr.DrniPortalAddr[3],
							dr.DrniPortalAddr[4],
							dr.DrniPortalAddr[5]),
						dr.DrniPortalPriority,
						dr.DRFHomeOperAggregatorKey,
						dr.DrniName,
						false)
				}
			}
		}
	}
}

//updateDRFHomeState This function updates the DRF_Home_State based on the operational
//state of the local ports as follows
func (psm *PsMachine) updateDRFHomeState(changePortal, changeDRFPorts bool) {
	// TODO need to understand the logic better
	dr := psm.dr
	a := dr.a
	psm.DrcpPsmLog(fmt.Sprintln("updateDRFHomeState: changePortal changeDRFPorts", dr.DrniName, a.AggName, changePortal, changeDRFPorts))
	if changeDRFPorts {
		//fmt.Println("updateDRFHomeState called with changedDRFPorts")
		dr.SetTimeSharingPortAndGatwewayDigest()
	}

	// Update the Home Gateway State
	psm.setDRFHomeState(changeDRFPorts)

	allippactivitynotset := true
	for _, ipp := range dr.Ipplinks {
		if ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
			allippactivitynotset = false
		}
	}
	if allippactivitynotset {
		dr.DrniPSI = true
	} else {
		dr.DrniPSI = false
	}

	// update the digests
	if changePortal {
		// If change portal event was received means that the gateway
		// has been updated so lets update the home state
		// update the Home State Gateway Vector sequence
		vector := make([]bool, 4096)
		if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
			for cid, portalsystemnumbers := range dr.DrniConvAdminGateway {
				if len(portalsystemnumbers) > 1 {
					psm.DrcpPsmLog(fmt.Sprintf("DrniConvAdminGateway[%d] == true", cid))
					vector[cid] = true
				}
			}
		}
		dr.DRFHomeState.mutex.Lock()
		if len(dr.DRFHomeState.GatewayVector) > 0 {
			dr.DRFHomeState.OpState = true
			dr.DRFHomeState.updateGatewayVector(dr.DRFHomeState.GatewayVector[0].Sequence+1, vector)
		} else {
			dr.DRFHomeState.OpState = true
			dr.DRFHomeState.updateGatewayVector(1, vector)
		}

		dr.DRFHomeState.mutex.Unlock()
		dr.HomeGatewayVectorTransmit = true

		dr.GatewayConversationUpdate = true
		if dr.GMachineFsm != nil {
			psm.DrcpPsmLog("Sending Update Gateway Conversation to GM machine")
			dr.GMachineFsm.GmEvents <- utils.MachineEvent{
				E:   GmEventGatewayConversationUpdate,
				Src: PsMachineModuleStr,
			}
		}
		// if this event occured it means that a conversation has been updated
		// state has changeds.  Lets trigger and NTT update so that the partner
		// has the correct view of the system
		dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStateGatewaySync)
		psm.DrcpPsmLog("Clearing Gateway Sync")
	}
	if changeDRFPorts {
		dr.DRFHomeState.mutex.Lock()
		// Assumption is made that if DRFPorts have changed that the aggregators
		// distributing ports have changed
		dr.DRFHomeState.PortIdList = make([]uint32, len(dr.DRAggregatorDistributedList))
		for i, portId := range dr.DRAggregatorDistributedList {
			dr.DRFHomeState.PortIdList[i] = uint32(portId)
		}
		dr.DRFHomeState.mutex.Unlock()
		dr.PortConversationUpdate = true
		if dr.AMachineFsm != nil {
			psm.DrcpPsmLog("Sending Update Port Conversation to AM machine")
			dr.AMachineFsm.AmEvents <- utils.MachineEvent{
				E:   AmEventPortConversationUpdate,
				Src: PsMachineModuleStr,
			}
		}
		// if this event occured it means that the aggregator ports now in distributing
		// state has changeds.  Lets trigger and NTT update so that the partner
		// has the correct view of the system
		dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStatePortSync)
		psm.DrcpPsmLog("Clearing Port Sync")
	}

	// Portal System is Isolated, which means IPP link is down
	// take down the AGGs
	if dr.DrniPSI &&
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) {
		// LaAggregator processing
		for _, intfref := range a.PortNumList {
			var p *lacp.LaAggPort
			if lacp.LaFindPortById(intfref, &p) {
				// Lets set the port to unselected
				lacp.SetLaAggPortCheckSelectionDistributedRelayIsSynced(intfref, false)
			}
		}
	}
}

// setDRFHomeState TRUE indicates operable (i.e., the local DR
// Function is able to relay traffic through its Gateway Port and at least one of its other Ports—
// IPP(s) or Aggregator) and that connectivity through the local Gateway is enabled by the
// operation of the network control protocol
func (psm *PsMachine) setDRFHomeState(changeDRFPorts bool) (gatewayChanged bool) {
	dr := psm.dr
	// 1) ipp link is down
	// 2) No agg ports are in distributed state
	// 3) No conversations exist
	distributedPortsValid := (dr.DRAggregatorDistributedList != nil &&
		len(dr.DRAggregatorDistributedList) > 0)
	isGatewaySet := dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit)
	for _, ipp := range dr.Ipplinks {
		if (ipp.IppPortEnabled || distributedPortsValid) &&
			!psm.isGatewayConvNull() {
			dr.DRFHomeOperDRCPState.SetState(layers.DRCPStateHomeGatewayBit)
			if !isGatewaySet {
				gatewayChanged = true
				psm.DrcpPsmLog("Home Gateway Set")
			}
		} else {
			dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStateHomeGatewayBit)
			if isGatewaySet {
				gatewayChanged = true
				psm.DrcpPsmLog("Home Gateway Cleared")
			}
		}
	}

	if !distributedPortsValid &&
		changeDRFPorts &&
		dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
		// TODO flush the mac table so no mac is forwarded
		// to neighbor card and not this local lag
		psm.DrcpPsmLog(fmt.Sprintln("Flush DB based on AGG", dr.DrniAggregator))
	}

	return gatewayChanged
}

// IsGatewayConvNull will check that there are no conversations provisioned
// on this system
func (psm *PsMachine) isGatewayConvNull() bool {
	dr := psm.dr
	conversationdoesnotexist := true
	for _, portalsystemnumbers := range dr.DrniConvAdminGateway {
		if portalsystemnumbers != nil {
			conversationdoesnotexist = false
			break
		}
	}
	return conversationdoesnotexist
}
