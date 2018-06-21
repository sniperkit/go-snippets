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
// 802.1ax-2014 Section 9.4.14 DRCPDU Receive machine
// rxmachine.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"sort"
	"strconv"
	"strings"
	"time"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

const RxMachineModuleStr = "DRCP Rx Machine"

// drxm States
const (
	RxmStateNone = iota + 1
	RxmStateInitialize
	RxmStateExpired
	RxmStatePortalCheck
	RxmStateCompatibilityCheck
	RxmStateDefaulted
	RxmStateDiscard // REPORT_TO_MANAGEMENT state
	RxmStateCurrent
)

var RxmStateStrMap map[fsm.State]string

func RxMachineStrStateMapCreate() {
	RxmStateStrMap = make(map[fsm.State]string)
	RxmStateStrMap[RxmStateNone] = "None"
	RxmStateStrMap[RxmStateInitialize] = "Initialize"
	RxmStateStrMap[RxmStateExpired] = "Expired"
	RxmStateStrMap[RxmStatePortalCheck] = "Portal Check"
	RxmStateStrMap[RxmStateCompatibilityCheck] = "Compatibility Check"
	RxmStateStrMap[RxmStateDefaulted] = "Defaulted"
	RxmStateStrMap[RxmStateDiscard] = "Discard"
	RxmStateStrMap[RxmStateCurrent] = "Current"
}

// rxm events
const (
	RxmEventBegin = iota + 1
	RxmEventNotIPPPortEnabled
	RxmEventNotDRCPEnabled
	RxmEventIPPPortEnabledAndDRCPEnabled
	RxmEventDRCPDURx
	RxmEventDRCPCurrentWhileTimerExpired
	RxmEventNotDifferPortal
	RxmEventDifferPortal
	RxmEventNotDifferConfPortal
	RxmEventDifferConfPortal
)

type RxDrcpPdu struct {
	pdu          *layers.DRCP
	src          string
	responseChan chan string
}

// DrcpRxMachine holds FSM and current State
// and event channels for State transitions
type RxMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	p *DRCPIpp

	MissConfiguredState bool

	// debug log
	//log chan string

	// timer interval
	currentWhileTimerTimeout time.Duration

	// timers
	currentWhileTimer *time.Timer

	// machine specific events
	RxmEvents     chan utils.MachineEvent
	RxmPktRxEvent chan RxDrcpPdu
}

func (rxm *RxMachine) PrevState() fsm.State { return rxm.PreviousState }

// PrevStateSet will set the previous State
func (rxm *RxMachine) PrevStateSet(s fsm.State) { rxm.PreviousState = s }

// Stop should clean up all resources
func (rxm *RxMachine) Stop() {
	rxm.CurrentWhileTimerStop()

	close(rxm.RxmEvents)
	close(rxm.RxmPktRxEvent)
}

// NewDrcpRxMachine will create a new instance of the RxMachine
func NewDrcpRxMachine(port *DRCPIpp) *RxMachine {
	rxm := &RxMachine{
		p:             port,
		PreviousState: RxmStateNone,
		RxmEvents:     make(chan utils.MachineEvent, 10),
		RxmPktRxEvent: make(chan RxDrcpPdu, 100)}

	port.RxMachineFsm = rxm

	// create then stop
	rxm.CurrentWhileTimerStart()
	rxm.CurrentWhileTimerStop()

	return rxm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (rxm *RxMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if rxm.Machine == nil {
		rxm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	rxm.Machine.Rules = r
	rxm.Machine.Curr = &utils.StateEvent{
		StrStateMap: RxmStateStrMap,
		LogEna:      false,
		Logger:      rxm.DrcpRxmLog,
		Owner:       RxMachineModuleStr,
	}

	return rxm.Machine
}

// DrcpRxMachineInitialize function to be called after
// State transition to INITIALIZE
func (rxm *RxMachine) DrcpRxMachineInitialize(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p
	dr := p.dr

	// Record default params
	rxm.recordDefaultDRCPDU()

	// should not be set but lets be complete according to definition of ChangePortal
	isset := p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity)
	if isset {
		defer rxm.NotifyChangePortalChanged(dr.ChangePortal, true)
		dr.ChangePortal = true
	}
	p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStateIPPActivity)
	// next State
	return RxmStateInitialize
}

// DrcpRxMachineExpired function to be called after
// State transition to EXPIRED
func (rxm *RxMachine) DrcpRxMachineExpired(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p
	dr := p.dr

	// should not be set but lets be complete according to definition of ChangePortal
	isset := p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity)
	if isset {
		defer rxm.NotifyChangePortalChanged(dr.ChangePortal, true)
		dr.ChangePortal = true
	}
	p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStateIPPActivity)
	defer rxm.NotifyDRCPStateTimeoutChange(p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout), true)
	// short timeout
	p.DRFNeighborOperDRCPState.SetState(layers.DRCPStateDRCPTimeout)
	// start the timer
	rxm.CurrentWhileTimerTimeoutSet(DrniShortTimeoutTime)
	rxm.CurrentWhileTimerStart()
	return RxmStateExpired
}

// DrcpRxMachinePortalCheck function to be called after
// State transition to PORTAL CHECK
func (rxm *RxMachine) DrcpRxMachinePortalCheck(m fsm.Machine, data interface{}) fsm.State {

	drcpPduInfo := data.(*layers.DRCP)

	rxm.recordPortalValues(drcpPduInfo)
	return RxmStatePortalCheck
}

// DrcpRxMachineCompatibilityCheck function to be called after
// State transition to COMPATIBILITY CHECK
func (rxm *RxMachine) DrcpRxMachineCompatibilityCheck(m fsm.Machine, data interface{}) fsm.State {
	drcpPduInfo := data.(*layers.DRCP)

	rxm.recordPortalConfValues(drcpPduInfo)
	return RxmStateCompatibilityCheck
}

// DrcpRxMachineDiscard function to be called after
// State transition to DISCARD
func (rxm *RxMachine) DrcpRxMachineDiscard(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p
	p.reportToManagement()

	// If the Portal System continues to receive DRCPDUs that do not
	// match the administratively configured expectations for a period longer than twice the Short Timeout the
	// state machine will transit to the DEFAULTED state
	if !rxm.MissConfiguredState {
		rxm.MissConfiguredState = true
		rxm.CurrentWhileTimerTimeoutSet(DrniShortTimeoutTime * 2)
		rxm.CurrentWhileTimerStart()
	}

	return RxmStateDiscard
}

// DrcpRxMachineDefaulted function to be called after
// State transition to DEFAULTED
func (rxm *RxMachine) DrcpRxMachineDefaulted(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p
	dr := p.dr

	dr.DRFHomeOperDRCPState.SetState(layers.DRCPStateExpired)
	rxm.recordDefaultDRCPDU()
	p.reportToManagement()

	return RxmStateDefaulted
}

// DrcpRxMachineCurrent function to be called after
// State transition to CURRENT
func (rxm *RxMachine) DrcpRxMachineCurrent(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p
	dr := p.dr

	rxm.MissConfiguredState = false

	drcpPduInfo := data.(*layers.DRCP)

	rxm.recordNeighborState(drcpPduInfo)
	rxm.updateNTT()

	// 1 short , 0 long
	if dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		rxm.CurrentWhileTimerTimeoutSet(DrniShortTimeoutTime)
	} else {
		rxm.CurrentWhileTimerTimeoutSet(DrniLongTimeoutTime)
	}

	rxm.CurrentWhileTimerStart()
	dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStateExpired)

	return RxmStateCurrent
}

func DrcpRxMachineFSMBuild(p *DRCPIpp) *RxMachine {

	RxMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new RxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	rxm := NewDrcpRxMachine(p)

	//BEGIN -> INITITIALIZE
	rules.AddRule(RxmStateNone, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateInitialize, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateExpired, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStatePortalCheck, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCompatibilityCheck, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDefaulted, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDiscard, RxmEventBegin, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCurrent, RxmEventBegin, rxm.DrcpRxMachineInitialize)

	// NOT IPP PORT ENABLED  > INITIALIZE
	rules.AddRule(RxmStateNone, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateInitialize, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateExpired, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStatePortalCheck, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCompatibilityCheck, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDefaulted, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDiscard, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCurrent, RxmEventNotIPPPortEnabled, rxm.DrcpRxMachineInitialize)

	// NOT DRCP ENABLED  > INITIALIZE
	rules.AddRule(RxmStateNone, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateInitialize, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateExpired, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStatePortalCheck, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCompatibilityCheck, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDefaulted, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateDiscard, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)
	rules.AddRule(RxmStateCurrent, RxmEventNotDRCPEnabled, rxm.DrcpRxMachineInitialize)

	// IPP PORT ENABLED AND DRCP ENABLED -> EXPIRED
	rules.AddRule(RxmStateInitialize, RxmEventIPPPortEnabledAndDRCPEnabled, rxm.DrcpRxMachineExpired)

	// DRCPDU RX -> PORTAL CHECK
	rules.AddRule(RxmStateExpired, RxmEventDRCPDURx, rxm.DrcpRxMachinePortalCheck)
	rules.AddRule(RxmStateDiscard, RxmEventDRCPDURx, rxm.DrcpRxMachinePortalCheck)
	rules.AddRule(RxmStateCurrent, RxmEventDRCPDURx, rxm.DrcpRxMachinePortalCheck)
	rules.AddRule(RxmStateDefaulted, RxmEventDRCPDURx, rxm.DrcpRxMachinePortalCheck)

	// NOT DIFFER PORTAL -> COMPATIBILITY CHECK
	rules.AddRule(RxmStatePortalCheck, RxmEventNotDifferPortal, rxm.DrcpRxMachineCompatibilityCheck)

	// DIFFER PORTAL -> DISCARD
	rules.AddRule(RxmStatePortalCheck, RxmEventDifferPortal, rxm.DrcpRxMachineDiscard)

	// NOT DIFFER CONF PORTAL -> CURRENT
	rules.AddRule(RxmStateCompatibilityCheck, RxmEventNotDifferConfPortal, rxm.DrcpRxMachineCurrent)

	// DIFFER CONF PORTAL -> DISCARD
	rules.AddRule(RxmStateCompatibilityCheck, RxmEventDifferConfPortal, rxm.DrcpRxMachineDiscard)

	// DRCP CURRENT WHILE TIMER EXPIRED -> DEFAULTED
	rules.AddRule(RxmStateExpired, RxmEventDRCPCurrentWhileTimerExpired, rxm.DrcpRxMachineDefaulted)
	rules.AddRule(RxmStateDiscard, RxmEventDRCPCurrentWhileTimerExpired, rxm.DrcpRxMachineDefaulted)
	rules.AddRule(RxmStateCurrent, RxmEventDRCPCurrentWhileTimerExpired, rxm.DrcpRxMachineExpired)

	// Create a new FSM and apply the rules
	rxm.Apply(&rules)

	return rxm
}

// DrcpRxMachineMain:  802.1ax-2014 Figure 9-23
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *DRCPIpp) DrcpRxMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.12 Receive Machine
	rxm := DrcpRxMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	rxm.Machine.Start(rxm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *RxMachine) {
		m.DrcpRxmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case <-m.currentWhileTimer.C:
				// special case if we have pending packets in the queue
				// by the time this expires we want to ensure the packet
				// gets processed first as this will clear/restart the timer
				if len(m.RxmPktRxEvent) == 0 {
					m.DrcpRxmLog("Current While Timer Expired")
					m.Machine.ProcessEvent(RxMachineModuleStr, RxmEventDRCPCurrentWhileTimerExpired, nil)
				}

			case event, ok := <-m.RxmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					if rv == nil {
						/* continue State transition */
						m.processPostStates(nil)
					}

					if rv != nil {
						m.DrcpRxmLog(strings.Join([]string{error.Error(rv), event.Src, RxmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(RxMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.DrcpRxmLog("Machine End")
					return
				}
			case rx, ok := <-m.RxmPktRxEvent:
				if ok {
					rv := m.Machine.ProcessEvent(RxMachineModuleStr, RxmEventDRCPDURx, rx.pdu)
					if rv == nil {
						/* continue State transition */
						m.processPostStates(rx.pdu)
					}

					// respond to caller if necessary so that we don't have a deadlock
					if rx.responseChan != nil {
						utils.SendResponse(RxMachineModuleStr, rx.responseChan)
					}
				} else {
					m.DrcpRxmLog("Machine End")
					return
				}
			}
		}
	}(rxm)
}

// NotifyDRCPStateTimeoutChange notify the Periodic Transmit Machine of a neighbor state
// timeout change
func (rxm *RxMachine) NotifyDRCPStateTimeoutChange(oldval, newval bool) {
	p := rxm.p
	if oldval != newval {
		//layers.DRCPShortTimeout
		if newval {
			if p.PtxMachineFsm != nil {
				p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
					E:   PtxmEventDRFNeighborOPerDRCPStateTimeoutEqualShortTimeout,
					Src: RxMachineModuleStr,
				}
			}
		} else {
			if p.PtxMachineFsm != nil {
				//layers.DRCPLongTimeout
				p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
					E:   PtxmEventDRFNeighborOPerDRCPStateTimeoutEqualLongTimeout,
					Src: RxMachineModuleStr,
				}
			}
		}
	}
}

// processPostStates will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStates(drcpPduInfo *layers.DRCP) {
	rxm.processPostStateInitialize()
	rxm.processPostStateExpired()
	rxm.processPostStatePortalCheck(drcpPduInfo)
	rxm.processPostStateCompatibilityCheck(drcpPduInfo)
	rxm.processPostStateDiscard()
	rxm.processPostStateCurrent()
	rxm.processPostStateDefaulted()
}

// processPostStateInitialize will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateInitialize() {
	p := rxm.p
	if rxm.Machine.Curr.CurrentState() == RxmStateInitialize &&
		p.IppPortEnabled && p.DRCPEnabled {
		rv := rxm.Machine.ProcessEvent(RxMachineModuleStr, RxmEventIPPPortEnabledAndDRCPEnabled, nil)
		if rv != nil {
			rxm.DrcpRxmLog(strings.Join([]string{error.Error(rv), RxMachineModuleStr, RxmStateStrMap[rxm.Machine.Curr.CurrentState()], strconv.Itoa(int(RxmEventIPPPortEnabledAndDRCPEnabled))}, ":"))
		} else {
			rxm.processPostStates(nil)
		}
	}
}

// processPostStateExpired will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateExpired() {
	// nothin to do events are triggered by rx packet or current while timer expired
}

// processPostStatePortalCheck will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStatePortalCheck(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	if rxm.Machine.Curr.CurrentState() == RxmStatePortalCheck {
		if p.DifferPortal {
			rv := rxm.Machine.ProcessEvent(RxMachineModuleStr, RxmEventDifferPortal, drcpPduInfo)
			if rv != nil {
				rxm.DrcpRxmLog(strings.Join([]string{error.Error(rv), RxMachineModuleStr, RxmStateStrMap[rxm.Machine.Curr.CurrentState()], strconv.Itoa(int(RxmEventDifferPortal))}, ":"))
			} else {
				rxm.processPostStates(drcpPduInfo)
			}
		} else {
			rv := rxm.Machine.ProcessEvent(RxMachineModuleStr, RxmEventNotDifferPortal, drcpPduInfo)
			if rv != nil {
				rxm.DrcpRxmLog(strings.Join([]string{error.Error(rv), RxMachineModuleStr, RxmStateStrMap[rxm.Machine.Curr.CurrentState()], strconv.Itoa(int(RxmEventNotDifferPortal))}, ":"))
			} else {
				rxm.processPostStates(drcpPduInfo)
			}
		}
	}
}

// processPostStateCompatibilityCheck will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateCompatibilityCheck(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	if rxm.Machine.Curr.CurrentState() == RxmStateCompatibilityCheck {
		if p.DifferConfPortal ||
			p.DifferConfPortalSystemNumber ||
			p.DifferGatewayDigest ||
			p.DifferPortDigest {
			rv := rxm.Machine.ProcessEvent(RxMachineModuleStr, RxmEventDifferConfPortal, drcpPduInfo)
			if rv != nil {
				rxm.DrcpRxmLog(strings.Join([]string{error.Error(rv), RxMachineModuleStr, RxmStateStrMap[rxm.Machine.Curr.CurrentState()], strconv.Itoa(int(RxmEventDifferConfPortal))}, ":"))
			} else {
				rxm.processPostStates(drcpPduInfo)
			}
		} else {
			rv := rxm.Machine.ProcessEvent(RxMachineModuleStr, RxmEventNotDifferConfPortal, drcpPduInfo)
			if rv != nil {
				rxm.DrcpRxmLog(strings.Join([]string{error.Error(rv), RxMachineModuleStr, RxmStateStrMap[rxm.Machine.Curr.CurrentState()], strconv.Itoa(int(RxmEventNotDifferConfPortal))}, ":"))
			} else {
				rxm.processPostStates(drcpPduInfo)
			}
		}
	}
}

// processPostStateDiscard will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateDiscard() {
	// nothin to do events are triggered by rx packet or current while timer expired
}

// processPostStateCurrent will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateCurrent() {
	// nothin to do events are triggered by rx packet or current while timer expired
}

// processPostStateCurrent will check local params to see if any conditions
// are met in order to transition to next state
func (rxm *RxMachine) processPostStateDefaulted() {
	// nothin to do events are triggered by rx packet or current while timer expired
}

func (rxm *RxMachine) NotifyChangePortalChanged(oldval, newval bool) {
	p := rxm.p
	dr := p.dr
	if newval {
		if dr != nil &&
			dr.PsMachineFsm != nil {
			rxm.DrcpRxmLog("Sending Change Portal Event to PSM")
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   PsmEventChangePortal,
				Src: RxMachineModuleStr,
			}
		}
	}
}

func (rxm *RxMachine) NotifyGatewayConversationUpdate(oldval, newval bool) {
	p := rxm.p
	dr := p.dr
	if newval &&
		dr != nil &&
		dr.GMachineFsm != nil &&
		dr.GMachineFsm.Machine.Curr.CurrentState() == GmStateDRNIGatewayInitialize ||
		dr.GMachineFsm.Machine.Curr.CurrentState() == GmStateDRNIGatewayUpdate ||
		dr.GMachineFsm.Machine.Curr.CurrentState() == GmStatePsGatewayUpdate {
		dr.GMachineFsm.GmEvents <- utils.MachineEvent{
			E:   GmEventGatewayConversationUpdate,
			Src: RxMachineModuleStr,
		}
	}
}

func (rxm *RxMachine) NotifyPortConversationUpdate(oldval, newval bool) {
	p := rxm.p
	dr := p.dr
	if newval &&
		dr != nil &&
		dr.AMachineFsm != nil &&
		dr.AMachineFsm.Machine.Curr.CurrentState() == AmStateDRNIPortInitialize ||
		dr.AMachineFsm.Machine.Curr.CurrentState() == AmStateDRNIPortUpdate ||
		dr.AMachineFsm.Machine.Curr.CurrentState() == AmStatePsPortUpdate {
		dr.AMachineFsm.AmEvents <- utils.MachineEvent{
			E:   AmEventPortConversationUpdate,
			Src: RxMachineModuleStr,
		}
	}
}

// updateNTT This function sets NTTDRCPDU to TRUE, if any of:
func (rxm *RxMachine) updateNTT() {

	// TODO this does not work at the moment
	p := rxm.p
	dr := p.dr
	if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
		!p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!p.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {

		rxm.DrcpRxmLog(fmt.Sprintf("Home Gateway Sync %t Home Port Sync %t Neighbor Gateway Sync %t Neighbor Port Sync %t",
			dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync),
			dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync),
			p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync),
			p.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync)))

		defer p.NotifyNTTDRCPUDChange(RxMachineModuleStr, p.NTTDRCPDU, true)
		p.NTTDRCPDU = true
	}

}

// recordDefaultDRCPDU: 802.1ax Section 9.4.1.1
//
// This function sets the current Neighbor Portal System's operational
// parameter values to the default parameter values provided by the
// administrator as follows
func (rxm *RxMachine) recordDefaultDRCPDU() {
	p := rxm.p
	dr := p.dr
	a := dr.a

	p.DRFNeighborPortAlgorithm = dr.DrniNeighborAdminPortAlgorithm
	p.DRFNeighborGatewayAlgorithm = dr.DrniNeighborAdminGatewayAlgorithm
	p.DRFNeighborConversationPortListDigest = dr.DRFNeighborAdminConversationPortListDigest
	p.DRFNeighborConversationGatewayListDigest = dr.DRFNeighborAdminConversationGatewayListDigest
	p.DRFNeighborOperDRCPState = dr.DRFNeighborAdminDRCPState
	if a != nil {
		p.DRFNeighborAggregatorPriority = uint16(a.PartnerSystemPriority)
		p.DRFNeighborAggregatorId = a.PartnerSystemId
		p.DrniNeighborPortalPriority = uint16(a.PartnerSystemPriority)
		p.DrniNeighborPortalAddr = a.PartnerSystemId
	}
	p.DRFNeighborPortalSystemNumber = p.DRFHomeConfNeighborPortalSystemNumber
	p.DRFNeighborConfPortalSystemNumber = dr.DRFPortalSystemNumber
	p.DRFNeighborState.mutex.Lock()
	p.DRFNeighborState.OpState = false
	p.DRFNeighborState.GatewayVector = nil
	p.DRFNeighborState.PortIdList = nil
	p.DRFNeighborState.mutex.Unlock()
	p.DRFOtherNeighborState.mutex.Lock()
	p.DRFOtherNeighborState.OpState = false
	p.DRFOtherNeighborState.GatewayVector = nil
	p.DRFOtherNeighborState.PortIdList = nil
	p.DRFOtherNeighborState.mutex.Unlock()
	p.DRFNeighborAdminAggregatorKey = 0
	p.DRFOtherNeighborAdminAggregatorKey = 0
	p.DRFNeighborOperPartnerAggregatorKey = 0
	p.DRFOtherNeighborOperPartnerAggregatorKey = 0
	if dr.DrniThreeSystemPortal {
		// TODO may need to chage this logic when 3P system supported
	} else {

		for i := 0; i < 1024; i++ {
			// Boolean vector set to 1
			p.DrniNeighborGatewayConversation[i] = 1
			if dr.ChangePortal {
				p.DrniNeighborPortConversation[i] = 1
			}
		}
	}
	p.CCTimeShared = false
	p.CCEncTagShared = false

	defer rxm.NotifyChangePortalChanged(dr.ChangePortal, true)
	dr.ChangePortal = true

}

// recordPortalValues: 802.1ax-2014 Section 9.4.11 Functions
//
// This function records the Neighbor’s Portal parameter values carried
// in a received DRCPDU (9.4.3.2) from an IPP, as the current operational
// parameter values for the immediate Neighbor Portal System on this IPP,
// as follows
func (rxm *RxMachine) recordPortalValues(drcpPduInfo *layers.DRCP) {

	rxm.recordPortalValuesSavePortalInfo(drcpPduInfo)
	rxm.recordPortalValuesCompareNeighborToPortal()
}

// recordPortalValuesSavePortalInfo records the Neighbor’s Portal parameter
// values carried in a received DRCPDU (9.4.3.2) from an IPP, as the current
// operational parameter values for the immediate Neighbor Portal System on
// this IPP,
func (rxm *RxMachine) recordPortalValuesSavePortalInfo(drcpPduInfo *layers.DRCP) {
	p := rxm.p

	p.DRFNeighborAggregatorPriority = drcpPduInfo.PortalInfo.AggPriority
	p.DRFNeighborAggregatorId = drcpPduInfo.PortalInfo.AggId
	p.DrniNeighborPortalPriority = drcpPduInfo.PortalInfo.PortalPriority
	p.DrniNeighborPortalAddr = drcpPduInfo.PortalInfo.PortalAddr

}

// recordPortalValuesCompareNeighborToPortal compares the newly updated values
// of the Neighbor Portal System to this Portal System’s expectations and if
func (rxm *RxMachine) recordPortalValuesCompareNeighborToPortal() {
	p := rxm.p
	dr := p.dr

	aggPrioEqual := p.DRFNeighborAggregatorPriority == dr.DrniAggregatorPriority
	aggIdEqual := p.DRFNeighborAggregatorId == dr.DrniAggregatorId
	portalPrioEqual := p.DrniNeighborPortalPriority == dr.DrniPortalPriority
	portalAddrEqual := (p.DrniNeighborPortalAddr[0] == dr.DrniPortalAddr[0] &&
		p.DrniNeighborPortalAddr[1] == dr.DrniPortalAddr[1] &&
		p.DrniNeighborPortalAddr[2] == dr.DrniPortalAddr[2] &&
		p.DrniNeighborPortalAddr[3] == dr.DrniPortalAddr[3] &&
		p.DrniNeighborPortalAddr[4] == dr.DrniPortalAddr[4] &&
		p.DrniNeighborPortalAddr[5] == dr.DrniPortalAddr[5])

	if aggPrioEqual &&
		aggIdEqual &&
		portalPrioEqual &&
		portalAddrEqual {
		p.DifferPortal = false
	} else {
		p.DifferPortal = true
		p.DifferPortalReason = ""
		if !aggPrioEqual {
			//fmt.Println(p.DRFNeighborAggregatorPriority, dr.DrniAggregatorPriority)
			rxm.DrcpRxmLog(fmt.Sprintf("Aggregator Priority Diff local %d neighbor %d", dr.DrniAggregatorPriority, p.DRFNeighborAggregatorPriority))
			p.DifferPortalReason += "Neighbor Aggregator Priority, "
		}
		if !aggIdEqual {
			//fmt.Println(p.DRFNeighborAggregatorId, dr.DrniAggregatorId)
			rxm.DrcpRxmLog(fmt.Sprintf("Aggregator Id Diff local %d neighbor %d", dr.DrniAggregatorId, p.DRFNeighborAggregatorId))
			p.DifferPortalReason += "Neighbor Aggregator Id, "
		}
		if !portalPrioEqual {
			rxm.DrcpRxmLog(fmt.Sprintf("Portal Priority Diff local %d neighbor %d", dr.DrniPortalPriority, p.DrniNeighborPortalPriority))
			p.DifferPortalReason += "Neighbor Portal Priority, "
		}
		if !portalAddrEqual {
			rxm.DrcpRxmLog(fmt.Sprintf("Portal Address Diff local %+v neighbor %+v", dr.DrniPortalAddr, p.DrniNeighborPortalAddr))
			p.DifferPortalReason += "Neighbor Portal Addr, "
		}
	}
}

// recordPortalConfValues: 802.1ax-2014 Section 9.4.11 Functions
//
// This function records the Neighbor Portal System’s values carried in
// the Portal Configuration Information TLV of a received DRCPDU (9.4.3.2)
// from an IPP, as the current operational parameter values for the immediate
// Neighbor Portal System on this IPP as follows
func (rxm *RxMachine) recordPortalConfValues(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	dr := p.dr

	// Save config from pkt
	rxm.recordPortalConfValuesSavePortalConfInfo(drcpPduInfo)

	// which conversation vector is preset
	TwoPGatewayConverationVectorPresent := drcpPduInfo.TwoPortalGatewayConversationVector.TlvTypeLength.GetTlv() == layers.DRCPTLV2PGatewayConversationVector
	ThreePGatewayConversationVectorPresent := drcpPduInfo.ThreePortalGatewayConversationVector1.TlvTypeLength.GetTlv() == layers.DRCPTLV3PGatewayConversationVector1 &&
		drcpPduInfo.ThreePortalGatewayConversationVector2.TlvTypeLength.GetTlv() == layers.DRCPTLV3PGatewayConversationVector2
	TwoPPortConversationVectorPresent := drcpPduInfo.TwoPortalPortConversationVector.TlvTypeLength.GetTlv() == layers.DRCPTLV2PPortConversationVector
	ThreePPortConversationVectorPresent := drcpPduInfo.ThreePortalPortConversationVector1.TlvTypeLength.GetTlv() == layers.DRCPTLV3PPortConversationVector1 &&
		drcpPduInfo.ThreePortalPortConversationVector2.TlvTypeLength.GetTlv() == layers.DRCPTLV3PPortConversationVector2

	// It then compares the newly updated values of the Neighbor Portal System to this Portal
	// System’s expectations and if the comparison of
	portalSystemNumEqual := p.DRFNeighborPortalSystemNumber == p.DRFHomeConfNeighborPortalSystemNumber
	confPortalSystemNumEqual := p.DRFNeighborConfPortalSystemNumber == dr.DRFPortalSystemNumber
	threeSystemPortalEqual := p.DrniNeighborThreeSystemPortal == dr.DrniThreeSystemPortal
	commonMethodsEqual := p.DrniNeighborCommonMethods == dr.DrniCommonMethods
	operAggKeyEqual := p.DRFNeighborOperAggregatorKey&0x3fff == dr.DRFHomeOperAggregatorKey&0x3fff
	operAggKeyFullEqual := p.DRFNeighborOperAggregatorKey == dr.DRFHomeOperAggregatorKey
	portAlgorithmEqual := p.DRFNeighborPortAlgorithm == dr.DRFHomePortAlgorithm
	conversationPortListDigestEqual := p.DRFNeighborConversationPortListDigest == dr.DRFHomeConversationPortListDigest
	gatewayAlgorithmEqual := p.DRFNeighborGatewayAlgorithm == dr.DRFHomeGatewayAlgorithm
	//fmt.Println("RX: GatewayListDigest:", p.DRFNeighborConversationGatewayListDigest, dr.DRFHomeConversationGatewayListDigest)
	conversationGatewayListDigestEqual := p.DRFNeighborConversationGatewayListDigest == dr.DRFHomeConversationGatewayListDigest

	// lets set this as it will be cleared later if the fields differ
	p.DifferConfPortalSystemNumber = false
	dr.ChangePortal = false
	p.MissingRcvGatewayConVector = true
	// The event post state procesing should catch this event change if it does not get
	// changed below
	p.DifferConfPortal = false

	p.DifferPortalReason = ""
	if !portalSystemNumEqual {
		p.DifferConfPortalSystemNumber = true
		rxm.DrcpRxmLog(fmt.Sprintf("Local System Number %d Neighbor System Number %d", p.DRFNeighborPortalSystemNumber, p.DRFHomeConfNeighborPortalSystemNumber))
		p.DifferPortalReason += "Portal System Number, "
	}
	if !confPortalSystemNumEqual {
		p.DifferConfPortalSystemNumber = true
		rxm.DrcpRxmLog(fmt.Sprintf("Conf Local System Number %d Neighbor System Number %d", p.DRFNeighborConfPortalSystemNumber, dr.DRFPortalSystemNumber))
		p.DifferPortalReason += "Conf Portal System Number, "
	}
	if operAggKeyEqual {
		p.DifferConfPortal = false
		if !operAggKeyFullEqual {
			rxm.DrcpRxmLog(fmt.Sprintf("Agg Key Full Diff Local 0x%x Neighbor 0x%x", p.DRFNeighborOperAggregatorKey, dr.DRFHomeOperAggregatorKey))
		}
	} else {
		// The event post state procesing should catch this event change
		rxm.DrcpRxmLog(fmt.Sprintf("Agg Key Diff Local 0x%x Neighbor 0x%x", p.DRFNeighborOperAggregatorKey&0x3fff, dr.DRFHomeOperAggregatorKey&0x3fff))
		p.DifferConfPortal = true
		p.DifferPortalReason += "Oper Aggregator Key, "

	}

	if !threeSystemPortalEqual {
		rxm.DrcpRxmLog(fmt.Sprintf("Three SystemPortal Diff Local %t Neighbor %t", p.DrniNeighborThreeSystemPortal, dr.DrniThreeSystemPortal))
		p.DifferPortalReason += "Three System Portal, "
	}

	if !gatewayAlgorithmEqual {
		rxm.DrcpRxmLog(fmt.Sprintf("Gateway Algorithm Diff Local %+v Neighbor %+v", p.DRFNeighborGatewayAlgorithm, dr.DRFHomeGatewayAlgorithm))
		p.DifferPortalReason += "Gateway Algorithm, "
	}

	if !p.DifferConfPortal &&
		(!threeSystemPortalEqual || !gatewayAlgorithmEqual) {
		for i := 0; i < 1024; i++ {
			// boolean vector
			p.DrniNeighborGatewayConversation[i] = 0xff
		}
		p.DifferGatewayDigest = true

	} else if !p.DifferConfPortal &&
		(threeSystemPortalEqual && gatewayAlgorithmEqual) {

		if conversationGatewayListDigestEqual {
			p.DifferGatewayDigest = false
			p.GatewayConversationTransmit = false
			p.MissingRcvGatewayConVector = false
		} else {
			p.DifferGatewayDigest = true
			p.GatewayConversationTransmit = true
			if TwoPGatewayConverationVectorPresent &&
				!p.dr.DrniThreeSystemPortal {
				if drcpPduInfo.TwoPortalGatewayConversationVector.Vector != nil {
					for i := 0; i < 512; i++ {
						p.DrniNeighborGatewayConversation[i] = drcpPduInfo.TwoPortalGatewayConversationVector.Vector[i]

					}
				}
				p.MissingRcvGatewayConVector = false
			} else if ThreePGatewayConversationVectorPresent &&
				p.dr.DrniThreeSystemPortal {
				// TODO need to change the logic to concatinate the two 3P vector
			} else if !TwoPGatewayConverationVectorPresent &&
				!ThreePGatewayConversationVectorPresent &&
				commonMethodsEqual &&
				p.dr.DrniCommonMethods &&
				(TwoPPortConversationVectorPresent ||
					ThreePPortConversationVectorPresent) {
				if TwoPPortConversationVectorPresent &&
					!p.dr.DrniThreeSystemPortal {
					for i := 0; i < 512; i++ {
						p.DrniNeighborGatewayConversation[i] = drcpPduInfo.TwoPortalPortConversationVector.Vector[i]
					}
					p.MissingRcvGatewayConVector = false
				} else if ThreePPortConversationVectorPresent &&
					p.dr.DrniThreeSystemPortal {
					// TODO need to change logic to concatinate the two 3P vector
				}
			} else {
				p.MissingRcvGatewayConVector = true
			}
		}
	}

	if !p.DifferConfPortal &&
		(!threeSystemPortalEqual || !portAlgorithmEqual) {
		if p.dr.DrniThreeSystemPortal {
			// TODO when 3P system supported
		} else {
			for i := 0; i < 1024; i++ {
				// boolean vector
				p.DrniNeighborGatewayConversation[i] = 0xff
			}
		}
	} else if !p.DifferConfPortal &&
		threeSystemPortalEqual &&
		portAlgorithmEqual {
		if conversationPortListDigestEqual {
			p.DifferPortDigest = false
			p.PortConversationTransmit = false
			p.MissingRcvPortConVector = false
		} else {
			p.DifferPortDigest = true
			p.PortConversationTransmit = true
			if TwoPPortConversationVectorPresent &&
				!p.dr.DrniThreeSystemPortal {
				for i := 0; i < 512; i++ {
					p.DrniNeighborPortConversation[i] = drcpPduInfo.TwoPortalPortConversationVector.Vector[i]
				}
				p.MissingRcvPortConVector = false
			} else if ThreePPortConversationVectorPresent &&
				p.dr.DrniThreeSystemPortal {
				// TODO when support 3P system
			} else if !TwoPPortConversationVectorPresent &&
				!ThreePPortConversationVectorPresent &&
				p.DrniNeighborCommonMethods == p.dr.DrniCommonMethods &&
				p.dr.DrniCommonMethods &&
				TwoPGatewayConverationVectorPresent &&
				ThreePGatewayConversationVectorPresent {
				if !dr.DrniThreeSystemPortal {
					for i := 0; i < 512; i++ {
						// boolean vector
						p.DrniNeighborPortConversation[i] = drcpPduInfo.TwoPortalGatewayConversationVector.Vector[i]
					}
				} else {
					// TODO when 3P system supported
				}
				p.MissingRcvPortConVector = false
			}
		}
	} else {
		p.MissingRcvPortConVector = true
	}

	if p.DifferConfPortalSystemNumber {
		if !p.MissingRcvGatewayConVector {
			for i := 0; i < 1023; i += 2 {
				p.DrniNeighborGatewayConversation[i] = p.DRFHomeConfNeighborPortalSystemNumber << 6
				p.DrniNeighborGatewayConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 4
				p.DrniNeighborGatewayConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 2
				p.DrniNeighborGatewayConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 0
				p.DrniNeighborGatewayConversation[i+1] = p.DRFHomeConfNeighborPortalSystemNumber << 6
				p.DrniNeighborGatewayConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 4
				p.DrniNeighborGatewayConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 2
				p.DrniNeighborGatewayConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 0
			}
		} else if p.MissingRcvPortConVector {
			for i := 0; i < 1023; i += 2 {
				p.DrniNeighborPortConversation[i] = p.DRFHomeConfNeighborPortalSystemNumber << 6
				p.DrniNeighborPortConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 4
				p.DrniNeighborPortConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 2
				p.DrniNeighborPortConversation[i] |= p.DRFHomeConfNeighborPortalSystemNumber << 0
				p.DrniNeighborPortConversation[i+1] = p.DRFHomeConfNeighborPortalSystemNumber << 6
				p.DrniNeighborPortConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 4
				p.DrniNeighborPortConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 2
				p.DrniNeighborPortConversation[i+1] |= p.DRFHomeConfNeighborPortalSystemNumber << 0
			}
		}
	}
	if !commonMethodsEqual {
		p.DifferPortalReason += "Common Methods, "
	}
	if !portAlgorithmEqual {
		p.DifferPortalReason += "Port Algorithm, "
	}

	// Sharing by time would mean according to Annex G:
	// There is no agreement on symmetric Port Conversation IDs across the DRNI
	// So, the portList Digest should just be the conversationid's as values
	// within each conversation id
	if !conversationPortListDigestEqual {
		p.DifferPortalReason += "Conversation Port List Digest, "
	}

	// There is an agreement for ownership of Gateway
	// CVID:
	// Odd VID owned by portal 1
	// Even VID owned by portal 2
	if !conversationGatewayListDigestEqual {
		p.DifferPortalReason += "Conversation Gateway List Digest, "
	}

	if p.DifferPortalReason != "" {
		defer rxm.NotifyChangePortalChanged(dr.ChangePortal, true)
		dr.ChangePortal = true
	}
}

// recordPortalConfValuesSavePortalConfInfo This function records the Neighbor
// Portal System’s values carried in the Portal Configuration Information TLV
// of a received DRCPDU (9.4.3.2) from an IPP, as the current operational
// parameter values for the immediate Neighbor Portal System on this IPP as follows
func (rxm *RxMachine) recordPortalConfValuesSavePortalConfInfo(drcpPduInfo *layers.DRCP) {
	p := rxm.p

	p.DRFNeighborPortalSystemNumber =
		uint8(drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyStatePortalSystemNum))
	p.DRFNeighborConfPortalSystemNumber =
		uint8(drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyStateNeighborConfPortalSystemNumber))
	p.DrniNeighborThreeSystemPortal =
		drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyState3SystemPortal) == 1
	p.DrniNeighborCommonMethods =
		drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyStateCommonMethods) == 1
	p.DrniNeighborONN =
		drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyStateOtherNonNeighbor) == 1
	p.DRFNeighborOperAggregatorKey =
		drcpPduInfo.PortalConfigInfo.OperAggKey
	p.DRFNeighborPortAlgorithm =
		drcpPduInfo.PortalConfigInfo.PortAlgorithm
	p.DRFNeighborConversationPortListDigest =
		drcpPduInfo.PortalConfigInfo.PortDigest
	p.DRFNeighborGatewayAlgorithm =
		drcpPduInfo.PortalConfigInfo.GatewayAlgorithm
	p.DRFNeighborConversationGatewayListDigest =
		drcpPduInfo.PortalConfigInfo.GatewayDigest

	// Portal System Machine expects DRFNeighborAdminAggregatorKey to be set, so lets
	// grab it from the HomePortsInfo, Possible bug in standard
	p.DRFNeighborAdminAggregatorKey = drcpPduInfo.HomePortsInfo.AdminAggKey
	p.DRFNeighborOperPartnerAggregatorKey = drcpPduInfo.HomePortsInfo.OperPartnerAggKey

}

// recordNeighborState: 802.1ax-2014 Section 9.4.11 Functions
//
// This function sets DRF_Neighbor_Oper_DRCP_State.IPP_Activity to TRUE and records the
//parameter values for the Drni_Portal_System_State[] and DRF_Home_Oper_DRCP_State
//carried in a received DRCPDU [item s) in 9.4.3.2] on the IPP, as the current parameter
// values for Drni_Neighbor_State[] and DRF_Neighbor_Oper_DRCP_State associated with
// this IPP respectively. In particular, the operational Boolean Gateway Vectors for
// each Portal System in the Drni_Neighbor_State[] are extracted from the received
// DRCPDU as follows
func (rxm *RxMachine) recordNeighborState(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	dr := p.dr
	a := dr.a

	// lets be complete according to definition of ChangePortal
	isset := p.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity)
	if !isset {
		defer rxm.NotifyChangePortalChanged(dr.ChangePortal, true)
		dr.ChangePortal = true
	}

	// Ipp Activity
	p.DRFNeighborOperDRCPState.SetState(layers.DRCPStateIPPActivity)

	// extract the neighbor system number
	//neighborSystemNum := uint8(drcpPduInfo.PortalConfigInfo.TopologyState.GetState(layers.DRCPTopologyStatePortalSystemNum))
	//p.DRFNeighborPortalSystemNumber = neighborSystemNum

	// save gateway vector information
	rxm.saveRcvNeighborGatewayVector(p.DRFNeighborPortalSystemNumber, drcpPduInfo)
	rxm.saveRcvHomeGatewayVector(dr.DrniPortalSystemNumber, drcpPduInfo)
	rxm.saveRcvOtherGatewayVector(drcpPduInfo)

	if drcpPduInfo.State.State.GetState(layers.DRCPStatePortSync) {
		p.DRFNeighborOperDRCPState.SetState(layers.DRCPStatePortSync)
	} else {
		p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStatePortSync)
	}

	if drcpPduInfo.HomePortsInfo.TlvTypeLength.GetTlv() == layers.DRCPTLVTypeHomePortsInfo {
		// Active_Home_Ports in the Home Ports Information TLV, carried in a
		// received DRCPDU on the IPP, are used as the current values for the DRF_Neighbor_State on
		// this IPP and are associated with the Portal System identified by DRF_Neighbor_Portal_System_Number;
		p.DrniNeighborState[p.DRFNeighborPortalSystemNumber].mutex.Lock()
		p.DrniNeighborState[p.DRFNeighborPortalSystemNumber].PortIdList = drcpPduInfo.HomePortsInfo.ActiveHomePorts
		p.DrniNeighborState[p.DRFNeighborPortalSystemNumber].mutex.Unlock()

		p.DRFNeighborState.mutex.Lock()
		// TODO
		// Is this the correct place to unlock the block?
		// This would be the fastest point of entry
		if len(p.DRFNeighborState.PortIdList) > 0 &&
			len(drcpPduInfo.HomePortsInfo.ActiveHomePorts) == 0 &&
			a != nil {
			for _, client := range utils.GetAsicDPluginList() {
				for _, aggport := range a.PortNumList {
					/* TEMP - add actual port names */
					err := client.IppIngressEgressPass("fpPort1", "fpPort2")
					if err != nil {
						dr.LaDrLog(fmt.Sprintf("ERROR (AttachAgg) setting Block from %s tolag port %s", utils.GetNameFromIfIndex(int32(p.Id)), int32(aggport)))
					}
				}
			}
		}
		p.DRFNeighborState.PortIdList = drcpPduInfo.HomePortsInfo.ActiveHomePorts
		p.DRFNeighborState.mutex.Unlock()
	}

	if drcpPduInfo.NeighborPortsInfo.TlvTypeLength.GetTlv() == layers.DRCPTLVTypeNeighborPortsInfo {
		// Active_Home_Ports in the Home Ports Information TLV, carried in a
		// received DRCPDU on the IPP, are used as the current values for the DRF_Neighbor_State on
		// this IPP and are associated with the Portal System identified by DRF_Neighbor_Portal_System_Number;
		p.DrniNeighborState[dr.DrniPortalSystemNumber].mutex.Lock()
		p.DrniNeighborState[dr.DrniPortalSystemNumber].PortIdList = drcpPduInfo.NeighborPortsInfo.ActiveNeighborPorts
		p.DrniNeighborState[dr.DrniPortalSystemNumber].mutex.Unlock()

	}

	rxm.compareOtherPortsInfo(drcpPduInfo)

	// Network / IPL sharing by time (9.3.2.1) is supported
	rxm.compareNetworkIPLMethod(drcpPduInfo)
	rxm.compareNetworkIPLSharingEncapsulation(drcpPduInfo)
	rxm.compareGatewayOperGatewayVector()
	rxm.comparePortIds()
}

// saveRcvNeighborGatewayVector will save the received Home Gateway Vector Info if the following
// logic is occurs:
//
// For the DRF_Rcv_Neighbor_Gateway_Conversation_Mask, if the Home_Gateway bit in the
// DRF_Home_Oper_DRCP_State carried in the received DRCPDU is 0;
// DRF_Rcv_Neighbor_Gateway_Conversation_Mask is set to NULL;
// Otherwise if the Home_Gateway_Vector field [item al) in 9.4.3.2] is present in the received
// Home Gateway Vector TLV [item ai) in 9.4.3.2];
// DRF_Rcv_Neighbor_Gateway_Conversation_Mask = Home_Gateway_Vector;
// The tuple (Home_Gateway_Sequence, Home_Gateway_Vector) in the received Home
// Gateway Vector TLV is stored as an entry in the Gateway Vector database for the
// Neighbor Portal System on this IPP, indexed by the received Home_Gateway_Sequence in
// increasing sequence number order, and;
// The OtherGatewayVectorTransmit on the other IPP, if it exists and is operational, is set to
// TRUE;
// Otherwise if the Home_Gateway_Vector field is not present in the received Home Gateway
// Vector TLV;
// The OtherGatewayVectorTransmit on the other IPP, if it exists and is operational, is set to
// FALSE, and;
// The Home_Gateway_Sequence [item ak) in 9.4.3.2] is used as an index for a query in the
// Gateway Vector database for the Neighbor Portal System on this IPP and;
// If the tuple (Home_Gateway_Sequence, Neighbor_Gateway_Vector) is stored as the first
// entry in the database, then;
// DRF_Rcv_Neighbor_Gateway_Conversation_Mask = Neighbor_Gateway_Vector;
// Otherwise
// DRF_Rcv_Neighbor_Gateway_Conversation_Mask = 1, where 1 is a Boolean
// Vector with all its 4096 elements set to 1.
func (rxm *RxMachine) saveRcvNeighborGatewayVector(portalSystemNum uint8, drcpPduInfo *layers.DRCP) {
	p := rxm.p
	dr := p.dr

	if !drcpPduInfo.State.State.GetState(layers.DRCPStateHomeGatewayBit) {
		p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStateHomeGatewayBit)
		// clear all entries == NULL
		//fmt.Printf("saveRcvNeighborGatewayVector: DRCPStateHomeGatewayBit not set in pkt\n")
		p.DRFRcvNeighborGatewayConversationMask = [MAX_CONVERSATION_IDS]bool{}
	} else {
		p.DRFNeighborOperDRCPState.SetState(layers.DRCPStateHomeGatewayBit)
		dr.DRFHomeOperDRCPState.SetState(layers.DRCPStateNeighborGatewayBit)

		if drcpPduInfo.State.State.GetState(layers.DRCPStateGatewaySync) {
			p.DRFNeighborOperDRCPState.SetState(layers.DRCPStateGatewaySync)
		} else {
			p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStateGatewaySync)
		}
		//fmt.Printf("saveRcvNeighborGatewayVector: pkg %+v", drcpPduInfo)
		// TLV being present is indicator that it existed in the pkt
		//fmt.Printf("saveRcvNeighborGatewayVector: GetTlv %d", drcpPduInfo.HomeGatewayVector.TlvTypeLength.GetTlv())
		// HomeGateway is not required to be sent, only if something has changed
		if drcpPduInfo.HomeGatewayVector.TlvTypeLength.GetTlv() == layers.DRCPTLVTypeHomeGatewayVector {
			vectorlen := len(drcpPduInfo.HomeGatewayVector.Vector)
			rxm.DrcpRxmLog(fmt.Sprintf("saveRcvNeighborGatewayVector: Home Gateway Vector update seq %d", drcpPduInfo.HomeGatewayVector.Sequence))
			if vectorlen == 512 {
				vector := make([]bool, MAX_CONVERSATION_IDS)
				for i, j := 0, 0; i < 512; i, j = i+1, j+8 {
					vector[j] = drcpPduInfo.HomeGatewayVector.Vector[i]>>7&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j] = drcpPduInfo.HomeGatewayVector.Vector[i]>>7&0x1 == 1
					vector[j+1] = drcpPduInfo.HomeGatewayVector.Vector[i]>>6&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+1] = drcpPduInfo.HomeGatewayVector.Vector[i]>>6&0x1 == 1
					vector[j+2] = drcpPduInfo.HomeGatewayVector.Vector[i]>>5&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+2] = drcpPduInfo.HomeGatewayVector.Vector[i]>>5&0x1 == 1
					vector[j+3] = drcpPduInfo.HomeGatewayVector.Vector[i]>>4&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+3] = drcpPduInfo.HomeGatewayVector.Vector[i]>>4&0x1 == 1
					vector[j+4] = drcpPduInfo.HomeGatewayVector.Vector[i]>>3&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+4] = drcpPduInfo.HomeGatewayVector.Vector[i]>>3&0x1 == 1
					vector[j+5] = drcpPduInfo.HomeGatewayVector.Vector[i]>>2&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+5] = drcpPduInfo.HomeGatewayVector.Vector[i]>>2&0x1 == 1
					vector[j+6] = drcpPduInfo.HomeGatewayVector.Vector[i]>>1&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+6] = drcpPduInfo.HomeGatewayVector.Vector[i]>>1&0x1 == 1
					vector[j+7] = drcpPduInfo.HomeGatewayVector.Vector[i]>>0&0x1 == 1
					p.DRFRcvNeighborGatewayConversationMask[j+7] = drcpPduInfo.HomeGatewayVector.Vector[i]>>0&0x1 == 1
				}

				p.DrniNeighborState[portalSystemNum].mutex.Lock()
				//fmt.Printf("saveRcvNeighborGatewayVector: DrniNeighborState[%d] homeportal[%d] from pkt homegatewayvector Sequence %d\n", portalSystemNum, dr.DrniPortalSystemNumber, drcpPduInfo.HomeGatewayVector.Sequence)
				p.DrniNeighborState[portalSystemNum].OpState = true
				rxm.DrcpRxmLog(fmt.Sprintf("saveRcvNeighborGatewayVector: portal[%d] Updateing Gateway Vector %+v", portalSystemNum, drcpPduInfo.HomeGatewayVector))
				p.DrniNeighborState[portalSystemNum].updateGatewayVector(drcpPduInfo.HomeGatewayVector.Sequence, vector)
				p.DrniNeighborState[portalSystemNum].mutex.Unlock()

				// set the current Gateway Sequence from the packet
				p.DRFNeighborGatewaySequence = drcpPduInfo.HomeGatewayVector.Sequence

				// record the immediate member
				p.DRFNeighborState.mutex.Lock()
				p.DRFNeighborState.OpState = true
				p.DRFNeighborState.updateGatewayVector(drcpPduInfo.HomeGatewayVector.Sequence, vector)
				p.DRFNeighborState.mutex.Unlock()

				// The OtherGatewayVectorTransmit on the other IPP, if it exists and is operational, is set to
				// TRUE, which mean we need to get the other ipp
				for _, ipp := range dr.Ipplinks {
					if ipp != p && ipp.DRCPEnabled {
						dr.OtherGatewayVectorTransmit = true
					}
				}
			} else if vectorlen == 0 {
				// The OtherGatewayVectorTransmit on the other IPP, if it exists and is operational, is set to
				// TRUE
				p.DrniNeighborState[portalSystemNum].mutex.Lock()
				for _, ipp := range dr.Ipplinks {
					if ipp != p && ipp.DRCPEnabled {
						dr.OtherGatewayVectorTransmit = false
					}
					// If the tuple (Home_Gateway_Sequence, Neighbor_Gateway_Vector) is stored as the first
					// entry in the database, then
					vector := make([]bool, MAX_CONVERSATION_IDS)
					index := p.DrniNeighborState[portalSystemNum].getNeighborVectorGatwaySequenceIndex(drcpPduInfo.HomeGatewayVector.Sequence, vector)
					if index == 0 {
						for i := 0; i < MAX_CONVERSATION_IDS; i++ {
							p.DRFRcvNeighborGatewayConversationMask[i] = p.DrniNeighborState[portalSystemNum].GatewayVector[index].Vector[i]
						}
					} else {
						for i := 0; i < MAX_CONVERSATION_IDS; i++ {
							p.DRFRcvNeighborGatewayConversationMask[i] = true
						}
					}
				}
				p.DrniNeighborState[portalSystemNum].mutex.Unlock()
			}
		}
	}
}

// saveRcvHomeGatewayVector will record the neighbor gateway vectore as follows
//
// For the DRF_Rcv_Home_Gateway_Conversation_Mask, if the Neighbor_Gateway bit in the
// DRF_Home_Oper_DRCP_State carried in the received DRCPDU is 0;
// DRF_Rcv_Home_Gateway_Conversation_Mask is set to NULL;
// Otherwise;
// The Neighbor_Gateway_Sequence [item ao) in 9.4.3.2] carried in the received Neighbor
// Gateway Vector TLV [item am) in 9.4.3.2] is used as an index for a query in the Gateway
// Vector database of this Portal System and;
// If the tuple (Neighbor_Gateway_Sequence, Home_Gateway_Vector) is stored in the
// database, then;
//     DRF_Rcv_Home_Gateway_Conversation_Mask = Home_Gateway_Vector;
//     In addition, if that is the first entry in the database, then;
//     The HomeGatewayVectorTransmit on this IPP is set to FALSE;
//     Otherwise;
//     The HomeGatewayVectorTransmit on this IPP is set to TRUE, and if the
//     Neighbor_Gateway_Sequence value is larger than the currently used
//     Home_Gateway_Sequence a new entry is created in Gateway Vector database
//     of this Portal System with the tuple values (Neighbor_Gateway_Sequence + 1,
//     Home_Gateway_Vector);
// Otherwise
// DRF_Rcv_Home_Gateway_Conversation_Mask = 1, where 1 is a Boolean Vector
// with all its 4096 elements set to 1, and;
// The HomeGatewayVectorTransmit is set to TRUE.
func (rxm *RxMachine) saveRcvHomeGatewayVector(portalSystemNum uint8, drcpPduInfo *layers.DRCP) {
	p := rxm.p
	dr := p.dr

	//fmt.Printf("saveRcvHomeGatewayVector: Pkt, DRCPStateNeighborGatewayBit %t\n", drcpPduInfo.State.State.GetState(layers.DRCPStateNeighborGatewayBit))

	if !drcpPduInfo.State.State.GetState(layers.DRCPStateNeighborGatewayBit) {
		// clear all entries == NULL
		// Since the neighbor has not received our gateway info
		//fmt.Printf("saveRcvHomeGatewayVector: clearing all homeGatewayConversationmask\n")

		p.DRFRcvHomeGatewayConversationMask = [MAX_CONVERSATION_IDS]bool{}
	} else {
		vector := make([]bool, MAX_CONVERSATION_IDS)

		dr.DrniPortalSystemState[portalSystemNum].mutex.Lock()
		p.DrniNeighborState[portalSystemNum].mutex.Lock()
		index := dr.DrniPortalSystemState[portalSystemNum].getNeighborVectorGatwaySequenceIndex(drcpPduInfo.NeighborGatewayVector.Sequence, vector)
		if index != -1 {
			//rxm.DrcpRxmLog(fmt.Sprintf("saveRcvHomeGatewayVector found in sequence %d in DrniPortalSystemState[%d] at index %d", drcpPduInfo.NeighborGatewayVector.Sequence, portalSystemNum, index))
			for i := 0; i < MAX_CONVERSATION_IDS; i++ {
				p.DRFRcvHomeGatewayConversationMask[i] = dr.DrniPortalSystemState[portalSystemNum].GatewayVector[index].Vector[i]
			}

			if index == 0 {
				//fmt.Printf("saveRcvHomeGatewayVector sequence is the first entry in vector, no update necessary setting sync\n")

				// gateway vector has not changed, thus no reason to transmit the gateway to the neighbor portals
				dr.HomeGatewayVectorTransmit = false
				homeindex := p.DrniNeighborState[dr.DrniPortalSystemNumber].getNeighborVectorGatwaySequenceIndex(drcpPduInfo.NeighborGatewayVector.Sequence, vector)
				// lets add the vector the the drni, since it agrees with the global
				//fmt.Printf("saveRcvHomeGatewayVector DrniNeighborState looking for sequence %d found at index %d\n", drcpPduInfo.NeighborGatewayVector.Sequence, homeindex)

				if homeindex != 0 {
					//fmt.Printf("saveRcvHomeGatewayVector updating DrniNeighborState[%d]\n", portalSystemNum)

					p.DrniNeighborState[portalSystemNum].OpState = true
					rxm.DrcpRxmLog(fmt.Sprintf("saveRcvHomeGatewayVector: portal[%d] Updateing Gateway Vector [200] = %v", portalSystemNum, dr.DrniPortalSystemState[portalSystemNum].GatewayVector[index].Vector[200]))

					p.DrniNeighborState[portalSystemNum].updateGatewayVector(drcpPduInfo.NeighborGatewayVector.Sequence,
						dr.DrniPortalSystemState[portalSystemNum].GatewayVector[0].Vector)
					p.DrniNeighborState[portalSystemNum].PortIdList = drcpPduInfo.HomePortsInfo.ActiveHomePorts
				}
			} else {
				dr.HomeGatewayVectorTransmit = true
				// check if received neighbor gateway vector has changed
				// TODO verify this
				rxm.DrcpRxmLog(fmt.Sprintf("saveRcvHomeGatewayVector neighbor seq %d, currseq %d\n", drcpPduInfo.NeighborGatewayVector.Sequence, dr.DrniPortalSystemState[portalSystemNum].GatewayVector[0].Sequence))

				if drcpPduInfo.NeighborGatewayVector.Sequence < dr.DrniPortalSystemState[portalSystemNum].GatewayVector[0].Sequence {
					//fmt.Sprintln("saveRcvHomeGatewayVector: NeighborState updating vector sequence", drcpPduInfo.NeighborGatewayVector.Sequence+1)
					rxm.DrcpRxmLog(fmt.Sprintln("saveRcvHomeGatewayVector: NeighborState updating vector sequence", drcpPduInfo.NeighborGatewayVector.Sequence+1))
					p.DrniNeighborState[portalSystemNum].OpState = true
					p.DrniNeighborState[portalSystemNum].updateGatewayVector(drcpPduInfo.NeighborGatewayVector.Sequence+1,
						dr.DrniPortalSystemState[portalSystemNum].GatewayVector[0].Vector)
				}
			}
		} else {
			currSeqList := make([]uint32, 0)
			for _, seqvector := range dr.DrniPortalSystemState[portalSystemNum].GatewayVector {
				currSeqList = append(currSeqList, seqvector.Sequence)
			}
			rxm.DrcpRxmLog(fmt.Sprintf("saveRcvHomeGatewayVector did not find sequence %d in currList %v", drcpPduInfo.NeighborGatewayVector.Sequence, currSeqList))

			for i := 0; i < 4096; i++ {
				p.DRFRcvHomeGatewayConversationMask[i] = true
			}
			dr.HomeGatewayVectorTransmit = true
		}
		p.DrniNeighborState[portalSystemNum].mutex.Unlock()
		dr.DrniPortalSystemState[portalSystemNum].mutex.Unlock()
	}
}

// saveRcvOtherGatewayVector will record the Other Gateway Vector info as follows
//
// For the DRF_Rcv_Other_Gateway_Conversation_Mask, if the Other_Gateway bit in the
// DRF_Home_Oper_DRCP_State carried in the received DRCPDU is 0;
// DRF_Rcv_Other_Gateway_Conversation_Mask is set to NULL;
// Otherwise if the Other_Gateway_Vector field [item as) in 9.4.3.2] is present in the received
// Other Gateway Vector TLV [item ap) in 9.4.3.2];
//    DRF_Rcv_Other_Gateway_Conversation_Mask = Other_Gateway_Vector; and
//    If on this IPP, Drni_Neighbor_ONN == FALSE;
//    The tuple (Other_Gateway_Sequence, Other_Gateway_Vector) in the received Other
//    Gateway Vector TLV is stored as an entry in the Gateway Vector database for the
//    Other neighbor Portal System on this IPP indexed by the received
//    Other_Gateway_Sequence in increasing sequence number order;
// Otherwise if the Other_Gateway_Vector field is not present in the received Other Gateway
// Vector TLV;
//    The Other_Gateway_Sequence [item ar) in 9.4.3.2] is used as an index for a query in the
//    Gateway Vector database for the Other neighbor Portal System on this IPP and;
//    If the tuple (Other_Gateway_Sequence, Other_Gateway_Vector) is stored in the database,
//    then;
//        DRF_Rcv_Other_Gateway_Conversation_Mask = Other_Gateway_Vector;
//        In addition, if that is the first entry in the database, then;
//        The OtherGatewayVectorTransmit on this IPP is set to FALSE;
//        Otherwise;
//        The OtherGatewayVectorTransmit on this IPP is set to TRUE;
//    Otherwise
//        DRF_Rcv_Other_Gateway_Conversation_Mask = 1, where 1 is a Boolean Vector
//        with all its 4096 elements set to 1, and;
//        The OtherGatewayVectorTransmit on this IPP is set to TRUE.
func (rxm *RxMachine) saveRcvOtherGatewayVector(drcpPduInfo *layers.DRCP) {
	/* TODO when 3P portal system supported

	p := rxm.p
	if !drcpPduInfo.State.State.GetState(layers.DRCPStateOtherGatewayBit) {
		p.DRFNeighborOperDRCPState.ClearState(layers.DRCPStateOtherGatewayBit)
		// clear all entries == NULL
		p.DRFRcvOtherGatewayConversationMask = [MAX_CONVERSATION_IDS]bool{}
	} else {
		p.DRFNeighborOperDRCPState.SetState(layers.DRCPStateOtherGatewayBit)
		if drcpPduInfo.OtherGatewayVector.TlvTypeLength.GetTlv() == layers.DRCPTLVTypeOtherGatewayVector {
			if len(drcpPduInfo.OtherGatewayVector.Vector) > 0 {
				for i, j := 0, 0; i < 512; i, j = i+1, j+8 {
					p.DRFRcvOtherGatewayConversationMask[j] = drcpPduInfo.OtherGatewayVector.Vector[j]>>7&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+1] = drcpPduInfo.OtherGatewayVector.Vector[j]>>6&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+2] = drcpPduInfo.OtherGatewayVector.Vector[j]>>5&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+3] = drcpPduInfo.OtherGatewayVector.Vector[j]>>4&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+4] = drcpPduInfo.OtherGatewayVector.Vector[j]>>3&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+5] = drcpPduInfo.OtherGatewayVector.Vector[j]>>2&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+6] = drcpPduInfo.OtherGatewayVector.Vector[j]>>1&0x1 == 1
					p.DRFRcvOtherGatewayConversationMask[j+7] = drcpPduInfo.OtherGatewayVector.Vector[j]>>0&0x1 == 1
				}
				if !p.DRNINeighborONN {
					p.DrniNeighborState[neighborSystemNum].updateGatewayVector(drcpPduInfo.OtherGatewayVector.Sequence,
						drcpPduInfo.OtherGatewayVecto.Vector)
				}
			} else {
				index := p.DrniNeighborState[neighborSystemNum].getNeighborVectorGatwaySequenceIndex(drcpPduInfo.OtherGatewayVector.Sequence,
					drcpPduInfo.NeighborGatewayVector.Vector)
				if index != -1 {
					p.DRFRcvOtherGatewayConversationMask = p.DRFRcvNeighborGatewayConversationMask[index]
					if index == 1 {
						p.OtherGatewayVectorTransmit = false
					} else {
						p.OtherGatewayVectorTransmit = true
					}
				} else {
					for i := 0; i < 4096; i++ {
						p.DRFRcvOtherGatewayConversationMask[i] = true
					}
					p.OtherGatewayVectorTransmit = true
				}
			}
		}
	}
	*/

}

// compareNetworkIPLSharingEncapsulation will compare the local portal encap method
// and sharing method with what the neighbor has configured
// NOTE: this implementation only supports SHARING_BY_TIME, thus this logic
//       should have not affect
func (rxm *RxMachine) compareNetworkIPLSharingEncapsulation(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	dr := p.dr
	if (dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TAG ||
		dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_ITAG ||
		dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_BTAG ||
		dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_PSEUDOWIRE) &&
		drcpPduInfo.NetworkIPLEncapsulation.TlvTypeLength.GetTlv() == layers.DRCPTLVNetworkIPLSharingEncapsulation &&
		drcpPduInfo.NetworkIPLMethod.TlvTypeLength.GetTlv() == layers.DRCPTLVNetworkIPLSharingMethod {
		// record the neighbor ipl share method
		p.DRFNeighborNetworkIPLSharingMethod[0] = drcpPduInfo.NetworkIPLMethod.Method[0]
		p.DRFNeighborNetworkIPLSharingMethod[1] = drcpPduInfo.NetworkIPLMethod.Method[1]
		p.DRFNeighborNetworkIPLSharingMethod[2] = drcpPduInfo.NetworkIPLMethod.Method[2]
		p.DRFNeighborNetworkIPLSharingMethod[3] = drcpPduInfo.NetworkIPLMethod.Method[3]

		// record the neighbor ipl encap digest
		p.DRFNeighborNetworkIPLIPLEncapDigest = drcpPduInfo.NetworkIPLEncapsulation.IplEncapDigest
		p.DRFNeighborNetworkIPLNetEncapDigest = drcpPduInfo.NetworkIPLEncapsulation.NetEncapDigest

		if p.DRFHomeNetworkIPLSharingMethod == p.DRFNeighborNetworkIPLSharingMethod &&
			p.DRFHomeNetworkIPLIPLEncapDigest == p.DRFNeighborNetworkIPLIPLEncapDigest &&
			p.DRFHomeNetworkIPLIPLNetEncapDigest == p.DRFNeighborNetworkIPLNetEncapDigest {
			rxm.DrcpRxmLog("Neighbor and Home IPL Sharing And Encap Do not agree")
			p.CCEncTagShared = true
		} else {
			rxm.DrcpRxmLog(fmt.Sprintf("Neighbor and Home IPL Sharing And Encap Do not agree local method[%+v] ipldigest[%+v] netdigest[%+v] neighbor method[%+v] ipldigest[%+v] netdigest[%+v]",
				p.DRFHomeNetworkIPLSharingMethod,
				p.DRFHomeNetworkIPLIPLEncapDigest,
				p.DRFHomeNetworkIPLIPLNetEncapDigest,
				p.DRFNeighborNetworkIPLSharingMethod,
				p.DRFNeighborNetworkIPLIPLEncapDigest,
				p.DRFNeighborNetworkIPLNetEncapDigest))
			p.CCEncTagShared = false
		}
	}
}

// compareOtherPortsInfo function will be used to compare the Other Ports info logic
// but since we are only supporting a 2P System then this function logic is not needed
// will be a placeholder and know that it is incomplete
func (rxm *RxMachine) compareOtherPortsInfo(drcpPduInfo *layers.DRCP) {
	/*
		p := rxm.p
		if drcpPduInfo.OtherPortsInfo.TlvTypeLength.GetTlv() == layers.DRCPTLVTypeOtherPortsInfo {
			// the Other_Neighbor_Ports in the Other Ports Information TLV,
			// carried in a received DRCPDU on the IPP, are used as the current values for the
			// DRF_Other_Neighbor_State on this IPP and are associated with the Portal System identified
			// by the value assigned to the two most significant bits of the
			// DRF_Other_Neighbor_Admin_Aggregator_Key carried within the Other Ports Information
			// TLV in the received DRCPDU
			p.DRFOtherNeighborState.PortIdList = drcpPduInfo.OtherPortsInfo.NeighborPorts
		} else if p.dr.DrniThreeSystemPortal {
			p.DRFOtherNeighborState.OpState = false
			p.DRFOtherNeighborState.GatewayVector = nil
			p.DRFOtherNeighborState.PortIdList = nil
			// no Portal System state information is available on this IPP for the distant
			// Neighbor Portal System on the IPP
			// TODO what does this mean that the info is not present for now going to clear/set
			// the keys
			p.DRFNeighborAdminAggregatorKey = p.dr.DRFHomeAdminAggregatorKey
			p.DRFNeighborOperPartnerAggregatorKey = p.dr.DRFHomeOperPartnerAggregatorKey
			// TODO Document states to set to this but if we reached this if the other should be
			// set to NULL
			//p.DRFOtherNeighborAdminAggregatorKey = p.DRFOtherNeighborAdminAggregatorKey
			//p.DRFOtherNeighborOperPartnerAggregatorKey = p.DRFOtherNeighborOperPartnerAggregatorKey
			p.DRFOtherNeighborAdminAggregatorKey = 0
			p.DRFOtherNeighborOperPartnerAggregatorKey = 0
		}
	*/
}

// compareNetworkIPLMethod will compare the network sharing method between what is configured
// between local portal and neighbor portal.
// NOTE: only SHARING_BY_TIME supported in this implementation
func (rxm *RxMachine) compareNetworkIPLMethod(drcpPduInfo *layers.DRCP) {
	p := rxm.p
	if drcpPduInfo.NetworkIPLMethod.TlvTypeLength.GetTlv() == layers.DRCPTLVNetworkIPLSharingMethod {
		p.DRFNeighborNetworkIPLSharingMethod = drcpPduInfo.NetworkIPLMethod.Method
		if p.DRFNeighborNetworkIPLSharingMethod == p.DRFHomeNetworkIPLSharingMethod &&
			p.DRFHomeNetworkIPLSharingMethod == ENCAP_METHOD_SHARING_BY_TIME {
			p.CCTimeShared = true
		} else {
			rxm.DrcpRxmLog(fmt.Sprintf("Neighbor and Home IPL Sharing by Time differ local method[%+v] neighbor method[%+v]",
				p.DRFNeighborNetworkIPLSharingMethod,
				p.DRFHomeNetworkIPLSharingMethod))
			p.CCTimeShared = false
		}
	}
}

// compareGatewayOperGatewayVector will compare the operational Gateway Vector info
// amongst all portals in the system.
func (rxm *RxMachine) compareGatewayOperGatewayVector() {

	p := rxm.p
	dr := p.dr

	operOrVectorDiffer := false
	for i := 1; i <= MAX_PORTAL_SYSTEM_IDS && !operOrVectorDiffer; i++ {
		dr.DrniPortalSystemState[i].mutex.Lock()
		p.DrniNeighborState[i].mutex.Lock()
		if dr.DrniPortalSystemState[i].OpState != p.DrniNeighborState[i].OpState {
			rxm.DrcpRxmLog(fmt.Sprintf("Neighbor Gateway OpState Different system portal num[%d] Prev[%t] New[%t]", i, dr.DrniPortalSystemState[i].OpState, p.DrniNeighborState[i].OpState))
			operOrVectorDiffer = true
		} else if dr.DrniPortalSystemState[i].OpState {
			// lets only compare the most recent gateway vector
			if (dr.DrniPortalSystemState[i].GatewayVector != nil && p.DrniNeighborState[i].GatewayVector != nil) &&
				((dr.DrniPortalSystemState[i].GatewayVector[0].Sequence != p.DrniNeighborState[i].GatewayVector[0].Sequence) ||
					(len(dr.DrniPortalSystemState[i].GatewayVector[0].Vector) != len(p.DrniNeighborState[i].GatewayVector[0].Vector))) {
				rxm.DrcpRxmLog(fmt.Sprintf("Local[%d] vs Neighbor Sequence/Vector Gateway Vector %d Different seq(%d) len(%d) != seq(%d) len(%d)",
					dr.DrniPortalSystemNumber, i, dr.DrniPortalSystemState[i].GatewayVector[0].Sequence, len(dr.DrniPortalSystemState[i].GatewayVector[0].Vector),
					p.DrniNeighborState[i].GatewayVector[0].Sequence, len(p.DrniNeighborState[i].GatewayVector[0].Vector)))
				//fmt.Printf("compareGatewayOperGatewayVector: Local[%d] vs Neighbor Sequence/Vector Gateway Vector %d Different seq(%d) len(%d) != seq(%d) len(%d)\n",
				//dr.DrniPortalSystemNumber, i, dr.DrniPortalSystemState[i].GatewayVector[0].Sequence, len(dr.DrniPortalSystemState[i].GatewayVector[0].Vector),
				//p.DrniNeighborState[i].GatewayVector[0].Sequence, len(p.DrniNeighborState[i].GatewayVector[0].Vector))
				operOrVectorDiffer = true
			} else {

				length := len(dr.DrniPortalSystemState[i].GatewayVector)
				//for j := 0; j < length; j++ {
				gatewayvector := dr.DrniPortalSystemState[i].getGatewayVectorByIndex(0)
				neighborvector := p.DrniNeighborState[i].getGatewayVectorByIndex(0)
				if neighborvector == nil {
					rxm.DrcpRxmLog(fmt.Sprintf("Neighbor Gateway Vector length(%d) length neighbor(%d) for portal %d does not exist", length, len(p.DrniNeighborState[i].GatewayVector), 0))
					operOrVectorDiffer = true
				}
				for k := 0; k < MAX_CONVERSATION_IDS && !operOrVectorDiffer; k++ {
					if gatewayvector.Vector[k] != neighborvector.Vector[k] {
						rxm.DrcpRxmLog(fmt.Sprintf("Neighbor Gateway Vector localsystem[%d] portal[%d] idx[%d] Different Prev[%d]=%t New[%d]=%t", dr.DrniPortalSystemNumber, i, 0, k, gatewayvector.Vector[k], k, neighborvector.Vector[k]))
						//rxm.DrcpRxmLog(fmt.Sprintf("Neighbor Gateway Vector idx[%d] Different Prev %v New %v", j, k, gatewayvector.Vector, k, neighborvector.Vector))
						operOrVectorDiffer = true
					}
				}
				//}
			}
		}
		p.DrniNeighborState[i].mutex.Unlock()
		dr.DrniPortalSystemState[i].mutex.Unlock()
	}

	if operOrVectorDiffer {
		defer rxm.NotifyGatewayConversationUpdate(dr.GatewayConversationUpdate, true)
		dr.GatewayConversationUpdate = true
		rxm.DrcpRxmLog("Clearing Gateway Sync, OperState or Vector Differ")
		dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStateGatewaySync)
		if p.MissingRcvGatewayConVector {
			for i := 0; i < 1024; i++ {
				if p.dr.DrniThreeSystemPortal {
					// TOOD need to change logic
					p.DrniNeighborGatewayConversation[i] = p.DRFHomeConfNeighborPortalSystemNumber
				} else {
					p.DrniNeighborGatewayConversation[i] = 0xff
				}
			}
		}
	} else {
		if p.DifferGatewayDigest {
			defer rxm.NotifyGatewayConversationUpdate(dr.GatewayConversationUpdate, true)
			dr.GatewayConversationUpdate = true
			dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStateGatewaySync)
			rxm.DrcpRxmLog("Clearing Gateway Sync, Gateway Digest different")
		} else if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) {
			defer rxm.NotifyGatewayConversationUpdate(dr.GatewayConversationUpdate, true)
			dr.GatewayConversationUpdate = true
			rxm.DrcpRxmLog("Setting Gateway Sync, Digest and Vector are the same")
			dr.DRFHomeOperDRCPState.SetState(layers.DRCPStateGatewaySync)
		} /*else {
			dr.DRFHomeOperDRCPState.SetState(layers.DRCPStateGatewaySync)
		}*/
	}
}

type sortPortList []uint32

// Helper Functions For sort
func (s sortPortList) Len() int           { return len(s) }
func (s sortPortList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortPortList) Less(i, j int) bool { return s[i] < s[j] }

// comparePortIds will compare the previously recorded that the portal system knowledge
// of the active ports in the systems agree.
func (rxm *RxMachine) comparePortIds() {
	p := rxm.p
	dr := p.dr
	portListDiffer := false

	// only need to check 1, 2 as we only support 2P systems
	for i := 1; i < MAX_PORTAL_SYSTEM_IDS && !portListDiffer; i++ {
		dr.DrniPortalSystemState[i].mutex.Lock()
		p.DrniNeighborState[i].mutex.Lock()
		//rxm.DrcpRxmLog(fmt.Sprintf("comparePortIds: localPortal[%d] Portal PortList %v neighbor view PortList %v\n",
		//	dr.DrniPortalSystemNumber, dr.DrniPortalSystemState[i].PortIdList, p.DrniNeighborState[i].PortIdList))
		if len(dr.DrniPortalSystemState[i].PortIdList) == len(p.DrniNeighborState[i].PortIdList) {
			// make sure the lists are sorted
			sort.Sort(sortPortList(dr.DrniPortalSystemState[i].PortIdList))
			sort.Sort(sortPortList(p.DrniNeighborState[i].PortIdList))
			for j, val := range dr.DrniPortalSystemState[i].PortIdList {
				if val != p.DrniNeighborState[i].PortIdList[j] {
					rxm.DrcpRxmLog(fmt.Sprintf("PortList Different Prev %v New %v", dr.DrniPortalSystemState[i].PortIdList, p.DrniNeighborState[i].PortIdList[j]))
					portListDiffer = true
				}
			}
		} else {
			rxm.DrcpRxmLog(fmt.Sprintf("PortList Different Prev %v New %v", dr.DrniPortalSystemState[i].PortIdList, p.DrniNeighborState[i].PortIdList))
			portListDiffer = true
		}
		p.DrniNeighborState[i].mutex.Unlock()
		dr.DrniPortalSystemState[i].mutex.Unlock()
	}

	if portListDiffer {
		defer rxm.NotifyPortConversationUpdate(p.dr.PortConversationUpdate, true)
		p.dr.PortConversationUpdate = true
		dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStatePortSync)
		if p.MissingRcvPortConVector {
			for i := 0; i < 1024; i++ {
				if p.dr.DrniThreeSystemPortal {
					// TODO when 3P system supported
				} else {
					p.DrniNeighborPortConversation[i] = 0xff
				}
			}
		}
	} else {
		if p.DifferPortDigest {
			defer rxm.NotifyPortConversationUpdate(p.dr.PortConversationUpdate, true)
			p.dr.PortConversationUpdate = true
			dr.DRFHomeOperDRCPState.ClearState(layers.DRCPStatePortSync)
			if dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
				rxm.DrcpRxmLog("Clearing Port Sync")
			}
		} else if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
			defer rxm.NotifyPortConversationUpdate(p.dr.PortConversationUpdate, true)
			p.dr.PortConversationUpdate = true
			rxm.DrcpRxmLog("Setting Port Sync")
			dr.DRFHomeOperDRCPState.SetState(layers.DRCPStatePortSync)
		} else {
			dr.DRFHomeOperDRCPState.SetState(layers.DRCPStatePortSync)
		}
	}
}
