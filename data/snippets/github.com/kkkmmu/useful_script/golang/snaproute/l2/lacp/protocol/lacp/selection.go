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

// selection
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"sync"

	"github.com/google/gopacket/layers"
)

/*
Aggregation is represented by an Aggregation Port selecting an appropriate Aggregator, and then attaching
to that Aggregator. The following are required for correct operation of the selection and attachment logic:

a) The implementation shall support at least one Aggregator per System.
b) Each Aggregation Port shall be assigned an operational Key (6.3.5). Aggregation Ports that can
   aggregate together are assigned the same operational Key as the other Aggregation Ports with which
   they can aggregate; Aggregation Ports that cannot aggregate with any other Aggregation Port are
   allocated unique operational Keys.
c) Each Aggregator shall be assigned an operational Key.
d) Each Aggregator shall be assigned an identifier that distinguishes it among the set of Aggregators in
   the System.
e) An Aggregation Port shall only select an Aggregator that has the same operational Key assignment
   as its own operational Key.
f) Subject to the exception Stated in item g), Aggregation Ports that are members of the same LAG
   (i.e., two or more Aggregation Ports that have the same Actor System ID, Actor Key, Partner System
   ID, and Partner Key, and that are not required to be Individual) shall select the same Aggregator.
g) Any pair of Aggregation Ports that are members of the same LAG, but are connected together by the
   same link, shall not select the same Aggregator (i.e., if a loopback condition exists between two
   Aggregation Ports, they shall not be aggregated together. For both Aggregation Ports, the Actor
   System ID is the same as the Partner System ID; also, for Aggregation Port A, the Partner’s Port
   Identifier is Aggregation Port B, and for Aggregation Port B, the Partner’s Port Identifier is
   Aggregation Port A).
   NOTE 1—This exception condition prevents the formation of an aggregated link, comprising two ends of the
   same link aggregated together, in which all frames transmitted through an Aggregator are immediately received
   through the same Aggregator. However, it permits the aggregation of multiple links that are in loopback; for
   example, if Aggregation Port A is looped back to Aggregation Port C and Aggregation Port B is looped back to
   Aggregation Port D, then it is permissible for A and B (or A and D) to aggregate together, and for C and D (or
   B and C) to aggregate together.
h) Any Aggregation Port that is required to be Individual (i.e., the operational State for the Actor or the
   Partner indicates that the Aggregation Port is Individual) shall not select the same Aggregator as any
   other Aggregation Port.
i) Any Aggregation Port that is Aggregateable shall not select an Aggregator to which an Individual
   Aggregation Port is already attached.
j) If the preceding conditions result in a given Aggregation Port being unable to select an Aggregator,
   then that Aggregation Port shall not be attached to any Aggregator.
k) If there are further constraints on the attachment of Aggregation Ports that have selected an
   Aggregator, those Aggregation Ports may be selected as standby in accordance with the rules
   specified in 6.7.1. Selection or deselection of that Aggregator can cause the Selection Logic to
   reevaluate the Aggregation Ports to be selected as standby.
l) The Selection Logic operates upon the operational information recorded by the Receive State
   machine, along with knowledge of the Actor’s own operational configuration and State. The
   Selection Logic uses the LAG ID for the Aggregation Port, determined from these operational
   parameters, to locate the correct Aggregator to which to attach the Aggregation Port.
m) The Selection Logic is invoked whenever an Aggregation Port is not attached to and has not selected
   an Aggregator, and executes continuously until it has determined the correct Aggregator for the
   Aggregation Port.
   NOTE 2—The Selection Logic may take a significant time to complete its determination of the correct
   Aggregator, as a suitable Aggregator may not be immediately available, due to configuration restrictions or the
   time taken to reallocate Aggregation Ports to other Aggregators.
n) Once the correct Aggregator has been determined, the variable Selected shall be set to SELECTED
   or to STANDBY (6.4.8, 6.7.1).
   NOTE 3—If Selected is SELECTED, the Mux machine will start the process of attaching the Aggregation Port
   to the selected Aggregator. If Selected is STANDBY, the Mux machine holds the Aggregation Port in the
   WAITING State, ready to be attached to its Aggregator once its Selected State changes to SELECTED.
o) The Selection Logic is responsible for computing the value of the Ready variable from the values of
   the Ready_N variable(s) associated with the set of Aggregation Ports that are waiting to attach to the
   same Aggregator (see 6.4.8).
p) Where the selection of a new Aggregator by an Aggregation Port, as a result of changes to the
   selection parameters, results in other Aggregation Ports in the System being required to reselect their
   Aggregators in turn, this is achieved by setting Selected to UNSELECTED for those other
   Aggregation Ports that are required to reselect their Aggregators.
   NOTE 4—The value of Selected is set to UNSELECTED by the Receive machine for the Aggregation Port
   when a change of LAG ID is detected.
q) An Aggregation Port shall not be enabled for use by the Aggregator Client until it has both selected
   and attached to an Aggregator.
r) An Aggregation Port shall not select an Aggregator, which has been assigned to a Portal (Clause 9),
   unless the Partner_Oper_Key of the associated LAG ID is equal to the lowest numerical value of the
   set comprising the values of the DRF_Home_Oper_Partner_Aggregator_Key, the
   DRF_Neighbor_Oper_Partner_Aggregator_Key,
   and the DRF_Other_Neighbor_Oper_Partner_Aggregator_Key, on each IPP on the Portal System, where
   any variable having its most significant bits set to 00 is excluded (corresponding to accepting only
   variables that have an integer value that is larger than 16383).
s) An Aggregation Port shall not select an Aggregator, which has been assigned to a Portal (Clause 9),
   if its Portal’s System Identifier is set to a value that is numerically lower than the Partner’s System
   Identifier, PSI == TRUE, the most significant two bits of Partner_Oper_Key are equal to the value 2
   or 3 and the two least significant bits of the Aggregation Port’s Partner_Oper_Port_priority are
   equal to the value 2 or 3. This is to prevent network partition due to isolation of the Portal Systems
   in the interconnected Portals (Clause 9)
*/

// LacpMuxCheckSelectionLogic will be called after the
// wait while timer has expired.  If this is the last
// port to have its wait while timer expire then
// will transition the mux State from waiting to
// attached
func (a *LaAggregator) LacpMuxCheckSelectionLogic(p *LaAggPort, sendResponse bool) {

	readyChan := make(chan bool)
	// lets do a this work in parrallel
	for _, pId := range a.PortNumList {

		go func(id uint16) {
			var port *LaAggPort
			if LaFindPortById(id, &port) {
				readyChan <- port.readyN
			}
		}(pId)
	}

	// lets set the agg.ready flag to true
	// until we know that at least one port
	// is not ready
	a.ready = true
	for range a.PortNumList {
		select {
		case readyN := <-readyChan:
			if !readyN {
				p.MuxMachineFsm.LacpMuxmLog("LacpMuxCheckSelectionLogic:Setting ready to false ")
				a.ready = false
			}
		}
	}

	close(readyChan)

	// if agg is ready then lets attach the
	// ports which are not already attached
	if a.ready {
		var wg sync.WaitGroup
		// lets do this work in parrallel
		for _, pId := range a.PortNumList {
			wg.Add(1)
			go func(id uint16) {
				defer wg.Done()
				var port *LaAggPort
				p.MuxMachineFsm.LacpMuxmLog(fmt.Sprintf("LacpMuxCheckSelectionLogic: looking for port %d", id))
				if LaFindPortById(id, &port) &&
					port.readyN {
					// trigger event to mux
					// event should be defered in the processing
					port.MuxMachineFsm.Machine.ProcessEvent(MuxMachineModuleStr, LacpMuxmEventSelectedEqualSelectedAndReady, nil)
				} else {
					p.MuxMachineFsm.LacpMuxmLog("LacpMuxCheckSelectionLogic: port not found or readyN is false")
				}
			}(pId)
		}
		wg.Wait()
	}
}

// updateSelected:  802.1ax Section 6.4.9
// Sets the value of the Selected variable based on the following:
//
// Rx pdu: (Actor: Port, Priority, System, System Priority, Key
// and State Aggregation) vs local recorded: (Partner Oper: Port, Priority,
// System, System Priority, Key, State Aggregation).  If values
// have changed then Selected is set to UNSELECTED, othewise
// SELECTED
func (rxm *LacpRxMachine) updateSelected(lacpPduInfo *layers.LACP) {

	p := rxm.p
	unselectedCondition := false
	//rxm.LacpRxmLog(fmt.Sprintf("PDU actor info %#v", lacpPduInfo.Actor.Info))
	//rxm.LacpRxmLog(fmt.Sprintf("Port partner oper info %#v", p.PartnerOper))

	// lets check a few conditions from Selection logic 802.1ax Section 6.4.14.1
	/*
		if rxm.detectLoopbackCondition(lacpPduInfo) {
			rxm.LacpRxmLog("ERROR Loopback condition detected")
			unselectedCondition = true
		}

		if rxm.detectInvalidPeerSelection(lacpPduInfo) {
			rxm.LacpRxmLog("ERROR Peer Selection Invalid")
			unselectedCondition = true
		}
	*/

	if !LacpLacpPktPortInfoIsEqual(&lacpPduInfo.Actor.Info, &p.PartnerOper, LacpStateAggregationBit) &&
		p.aggSelected == LacpAggSelected {
		rxm.LacpRxmLog("ERROR PDU and Oper States do not agree, moving port to UnSelected")
		unselectedCondition = true
	}

	if unselectedCondition {
		// lets trigger the event only if mux is not in waiting State as
		// the wait while timer expiration will trigger the unselected event

		if p.MuxMachineFsm != nil &&
			(p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateWaiting &&
				p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCWaiting) {
			p.aggSelected = LacpAggUnSelected
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventSelectedEqualUnselected,
				Src: RxMachineModuleStr}
		}
	}
}

// detectLoopbackCondition: 802.1ax Section 6.4.14.1 (g)
// Any pair of Aggregation Ports that are members of the same LAG, but are connected together by the
// same link, shall not select the same Aggregator (i.e., if a loopback condition exists between two
// Aggregation Ports, they shall not be aggregated together. For both Aggregation Ports, the Actor
// System ID is the same as the Partner System ID; also, for Aggregation Port A, the Partner’s Port
// Identifier is Aggregation Port B, and for Aggregation Port B, the Partner’s Port Identifier is
//Aggregation Port A).
func (rxm *LacpRxMachine) detectLoopbackCondition(lacpPduInfo *layers.LACP) bool {
	p := rxm.p
	// will allow ports to talk to each other on same system but on different aggregator
	// keys
	if (p.PartnerOper.System.Actor_System == lacpPduInfo.Actor.Info.System.SystemId &&
		p.PartnerOper.System.Actor_System_priority == lacpPduInfo.Actor.Info.System.SystemPriority) &&
		(lacpPduInfo.Actor.Info.Key == p.ActorOper.Key ||
			lacpPduInfo.Actor.Info.Port == p.ActorOper.port) {
		return true
	}
	return false
}

// detectInvalidPeerAggregators: detect that if multiple aggregation ports are connected
// to mulitiple peers which don't have the same aggregator info then all ports should be taken down
func (rxm *LacpRxMachine) detectInvalidPeerSelection(lacpPduInfo *layers.LACP) bool {

	p := rxm.p
	a := p.AggAttached
	selectionInvalid := false
	if a != nil {
		for _, pId := range a.PortNumList {
			if pId != p.PortNum {
				var aggport *LaAggPort
				if LaFindPortById(pId, &aggport) {
					if (aggport.PartnerOper.System.Actor_System != lacpPduInfo.Actor.Info.System.SystemId ||
						aggport.PartnerOper.System.Actor_System_priority != lacpPduInfo.Actor.Info.System.SystemPriority) ||
						aggport.PartnerOper.Key != lacpPduInfo.Actor.Info.Key {
						selectionInvalid = true
						if p.MuxMachineFsm != nil &&
							(aggport.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateWaiting &&
								aggport.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCWaiting) {
							aggport.aggSelected = LacpAggUnSelected
							aggport.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
								E:   LacpMuxmEventSelectedEqualUnselected,
								Src: RxMachineModuleStr}
						}
					}
				}
			}
		}
	}
	return selectionInvalid
}

// updateDefaultedSelected: 802.1ax Section 6.4.9
//
// Update the value of the Selected variable comparing
// the Partner admin info based with the partner
// operational info
// (port num, port priority, System, System priority,
//  Key, stat.Aggregation)
func (rxm *LacpRxMachine) updateDefaultSelected() {

	p := rxm.p

	rxm.LacpRxmLog(fmt.Sprintf("Port partner admin info %#v", p.partnerAdmin))
	rxm.LacpRxmLog(fmt.Sprintf("Port partner oper info %#v", p.PartnerOper))
	if !LacpLacpPortInfoIsEqual(&p.partnerAdmin, &p.PartnerOper, LacpStateAggregationBit) {
		//p.aggSelected = LacpAggUnSelected
		// lets trigger the event only if mux is not in waiting State as
		// the wait while timer expiration will trigger the unselected event
		if p.MuxMachineFsm != nil &&
			(p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateWaiting &&
				p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCWaiting) {
			p.aggSelected = LacpAggUnSelected
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventSelectedEqualUnselected,
				Src: RxMachineModuleStr}
		}
	}
}

// checkConfigForSelection will send selection bit to State machine
// and return to the user true
func (p *LaAggPort) checkConfigForSelection() bool {
	var a *LaAggregator

	// check to see if aggrigator exists
	// and that the Keys match
	if p.AggId != 0 && LaFindAggById(p.AggId, &a) {
		if p.MuxMachineFsm != nil {
			if (p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDetached ||
				p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCDetached) &&
				p.ActorOper.Key == a.ActorOperKey &&
				//p.PortEnabled &&
				(p.DrniName == "" ||
					p.DrniName != "" && p.DrniSynced) {

				p.LaPortLog("checkConfigForSelection: selected")

				// set port as selected
				p.aggSelected = LacpAggSelected
				LacpStateSet(&p.ActorOper.State, LacpStateAggregationBit)

				mEvtChan := make([]chan utils.MachineEvent, 0)
				evt := make([]utils.MachineEvent, 0)

				mEvtChan = append(mEvtChan, p.MuxMachineFsm.MuxmEvents)
				evt = append(evt, utils.MachineEvent{
					E:   LacpMuxmEventSelectedEqualSelected,
					Src: PortConfigModuleStr})
				// inform mux that port has been selected
				// wait for response
				p.DistributeMachineEvents(mEvtChan, evt, true)
				//msg := <-p.portChan
				return true
			} else if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached &&
				p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCDetached {

				p.LaPortLog("checkConfigForSelection: unselected")
				// set port as selected
				p.aggSelected = LacpAggUnSelected
				// attach the agg to the port
				//p.AggAttached = a

				mEvtChan := make([]chan utils.MachineEvent, 0)
				evt := make([]utils.MachineEvent, 0)

				mEvtChan = append(mEvtChan, p.MuxMachineFsm.MuxmEvents)
				evt = append(evt, utils.MachineEvent{
					E:   LacpMuxmEventSelectedEqualUnselected,
					Src: PortConfigModuleStr})
				// inform mux that port has been selected
				// wait for response
				p.DistributeMachineEvents(mEvtChan, evt, true)
			}
		}
	}
	return false
}
