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

// rxmachine
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

const RxMachineModuleStr = "Rx Machine"

// rxm States
const (
	LacpRxmStateNone = iota + 1
	LacpRxmStateInitialize
	LacpRxmStatePortDisabled
	LacpRxmStateExpired
	LacpRxmStateLacpDisabled
	LacpRxmStateDefaulted
	LacpRxmStateCurrent
)

var RxmStateStrMap map[fsm.State]string

func RxMachineStrStateMapCreate() {
	RxmStateStrMap = make(map[fsm.State]string)
	RxmStateStrMap[LacpRxmStateNone] = "None"
	RxmStateStrMap[LacpRxmStateInitialize] = "Initialize"
	RxmStateStrMap[LacpRxmStatePortDisabled] = "PortDisabled"
	RxmStateStrMap[LacpRxmStateExpired] = "Expired"
	RxmStateStrMap[LacpRxmStateLacpDisabled] = "LacpDisabled"
	RxmStateStrMap[LacpRxmStateDefaulted] = "Defaulted"
	RxmStateStrMap[LacpRxmStateCurrent] = "Current"
}

// rxm events
const (
	LacpRxmEventBegin = iota + 1
	LacpRxmEventUnconditionalFallthrough
	LacpRxmEventNotPortEnabledAndNotPortMoved
	LacpRxmEventPortMoved
	LacpRxmEventPortEnabledAndLacpEnabled
	LacpRxmEventPortEnabledAndLacpDisabled
	LacpRxmEventCurrentWhileTimerExpired
	LacpRxmEventLacpEnabled
	LacpRxmEventLacpPktRx
	LacpRxmEventKillSignal
)

type LacpRxLacpPdu struct {
	pdu          *layers.LACP
	src          string
	responseChan chan string
}

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type LacpRxMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	p *LaAggPort

	// timer interval
	currentWhileTimerTimeout time.Duration

	// timers
	currentWhileTimer *time.Timer

	// machine specific events
	RxmEvents         chan utils.MachineEvent
	RxmPktRxEvent     chan LacpRxLacpPdu
	RxmLogEnableEvent chan bool
}

func (rxm *LacpRxMachine) PrevState() fsm.State { return rxm.PreviousState }

// PrevStateSet will set the previous State
func (rxm *LacpRxMachine) PrevStateSet(s fsm.State) { rxm.PreviousState = s }

// Stop should clean up all resources
func (rxm *LacpRxMachine) Stop() {
	rxm.CurrentWhileTimerStop()

	close(rxm.RxmEvents)
	close(rxm.RxmPktRxEvent)
	close(rxm.RxmLogEnableEvent)

}

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewLacpRxMachine(port *LaAggPort) *LacpRxMachine {
	rxm := &LacpRxMachine{
		p:                 port,
		PreviousState:     LacpRxmStateNone,
		RxmEvents:         make(chan utils.MachineEvent, 10),
		RxmPktRxEvent:     make(chan LacpRxLacpPdu, 1000),
		RxmLogEnableEvent: make(chan bool)}

	port.RxMachineFsm = rxm

	// create then stop
	rxm.CurrentWhileTimerStart()
	rxm.CurrentWhileTimerStop()

	return rxm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (rxm *LacpRxMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if rxm.Machine == nil {
		rxm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	rxm.Machine.Rules = r
	rxm.Machine.Curr = &utils.StateEvent{
		StrStateMap: RxmStateStrMap,
		LogEna:      rxm.p.logEna,
		Logger:      rxm.LacpRxmLog,
		Owner:       RxMachineModuleStr,
	}

	return rxm.Machine
}

// LacpRxMachineInitialize function to be called after
// State transition to INITIALIZE
func (rxm *LacpRxMachine) LacpRxMachineInitialize(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	if timeoutTime, ok := rxm.CurrentWhileTimerValid(); !ok {
		rxm.CurrentWhileTimerTimeoutSet(timeoutTime)
		rxm.CurrentWhileTimerStart()
	}
	// Lets ensure that the port moves to the correct defaulted State
	// after initialization.  Default params will change after lacp
	// packets have arrived
	LacpStateSet(&p.PartnerOper.State, p.partnerAdmin.State)

	// set the agg as being unselected
	p.aggSelected = LacpAggUnSelected
	if p.MuxMachineFsm != nil {
		if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached &&
			p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCDetached {
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventSelectedEqualUnselected,
				Src: RxMachineModuleStr}
		}
	}

	// Record default params
	rxm.recordDefault()

	// Actor Port Oper State Expired = False
	LacpStateClear(&p.ActorOper.State, LacpStateExpiredBit)

	// set the port moved to false
	p.portMoved = false

	// next State
	return LacpRxmStateInitialize
}

// LacpRxMachineExpired function to be called after
// State transition to PORT_DISABLED
func (rxm *LacpRxMachine) LacpRxMachinePortDisabled(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	// Partner Port Oper State Sync = False
	LacpStateClear(&p.PartnerOper.State, LacpStateSyncBit)

	// inform partner cdm
	if p.PCdMachineFsm != nil &&
		p.PCdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateNoPartnerChurn {
		p.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventPartnerOperPortStateSyncOff,
			Src: RxMachineModuleStr}
	}
	if p.MuxMachineFsm != nil {
		if p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDistributing ||
			p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCollecting ||
			p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxStateCCollectingDistributing {
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventNotPartnerSync,
				Src: RxMachineModuleStr}
		}
	}

	return LacpRxmStatePortDisabled
}

// LacpRxMachineExpired function to be called after
// State transition to EXPIRED
func (rxm *LacpRxMachine) LacpRxMachineExpired(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	// Partner Port Oper State Sync = FALSE
	//rxm.LacpRxmLog("Clearing Partner Sync Bit")
	LacpStateClear(&p.PartnerOper.State, LacpStateSyncBit)
	// inform partner cdm
	if p.PCdMachineFsm != nil &&
		p.PCdMachineFsm.Machine.Curr.CurrentState() == LacpCdmStateNoPartnerChurn {
		p.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
			E:   LacpCdmEventPartnerOperPortStateSyncOff,
			Src: RxMachineModuleStr}
	}

	if p.MuxMachineFsm != nil {
		p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
			E:   LacpMuxmEventNotPartnerSync,
			Src: RxMachineModuleStr}
	}
	// Short timeout
	//rxm.LacpRxmLog("Setting Partner Timeout Bit")
	LacpStateSet(&p.PartnerOper.State, LacpStateTimeoutBit)

	// Set the Short timeout
	rxm.CurrentWhileTimerTimeoutSet(LacpShortTimeoutTime)

	// Start the Current While timer
	rxm.CurrentWhileTimerStart()

	// Actor Port Oper State Expired = TRUE
	//rxm.LacpRxmLog("Setting Actor Expired Bit")
	LacpStateSet(&p.ActorOper.State, LacpStateExpiredBit)

	return LacpRxmStateExpired
}

// LacpRxMachineLacpDisabled function to be called after
// State transition to LACP_DISABLED
func (rxm *LacpRxMachine) LacpRxMachineLacpDisabled(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	// stop the current while timer as it does not need to run as LACP is now
	// disabled
	rxm.CurrentWhileTimerStop()

	// Unselect the aggregator
	p.aggSelected = LacpAggUnSelected
	if p.MuxMachineFsm != nil {
		if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached &&
			p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateCDetached {
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventSelectedEqualUnselected,
				Src: RxMachineModuleStr}
		}
	}

	// setup the default params
	rxm.recordDefault()

	// Partner Port Oper State Aggregation = FALSE
	LacpStateClear(&p.PartnerOper.State, LacpStateAggregationBit)

	// Actor Port Oper State Expired = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateExpiredBit)

	return LacpRxmStateLacpDisabled
}

// LacpRxMachineDefaulted function to be called after
// State transition to DEFAULTED
func (rxm *LacpRxMachine) LacpRxMachineDefaulted(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	//lacpPduInfo := data.(LacpPdu)

	// Updated the default selected State
	rxm.updateDefaultSelected()

	// Record the default partner info
	rxm.recordDefault()

	// Actor Port Oper State Expired = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateExpiredBit)

	// Lets set the partner admin State to aggregatable and up
	LacpStateSet(&p.partnerAdmin.State, LacpStateAggregatibleUp)

	return LacpRxmStateDefaulted
}

// LacpRxMachineCurrent function to be called after
// State transition to CURRENT
func (rxm *LacpRxMachine) LacpRxMachineCurrent(m fsm.Machine, data interface{}) fsm.State {
	p := rxm.p

	// Version 1, V2 will require a serialize/deserialize routine since TLV's are involved
	lacpPduInfo := data.(*layers.LACP)

	// update selection logic
	rxm.updateSelected(lacpPduInfo)

	// update the ntt
	ntt := rxm.updateNTT(lacpPduInfo)

	// Version 2 or higher check
	if LacpActorSystemLacpVersion >= 0x2 {
		rxm.recordVersionNumber(lacpPduInfo)
	}

	// record the current packet State
	rxm.recordPDU(lacpPduInfo)

	//rxm.LacpRxmLog(fmt.Sprintf("Partner Oper %#v", p.PartnerOper))

	// Current while should already be set to
	// Actors Oper value of Timeout, lets check
	// anyways
	if timeoutTime, ok := rxm.CurrentWhileTimerValid(); !ok {
		rxm.CurrentWhileTimerTimeoutSet(timeoutTime)
	}
	// lets kick off the Current While Timer
	rxm.CurrentWhileTimerStart()

	// Actor_Oper_Port_Sate.Expired = FALSE
	LacpStateClear(&p.ActorOper.State, LacpStateExpiredBit)

	if ntt && p.TxMachineFsm != nil {
		// update ntt, which should trigger a packet transmit
		p.TxMachineFsm.TxmEvents <- utils.MachineEvent{
			E:   LacpTxmEventNtt,
			Src: RxMachineModuleStr}
	}

	// Other machines may need to be informed of the various
	// State info changes
	rxm.InformMachinesOfStateChanges()

	// In the event that the rx machine times out we want to ensure that the port
	// stays down so lets change the default partner admin State
	LacpStateSet(&p.partnerAdmin.State, LacpStateAggregatibleDown)

	return LacpRxmStateCurrent
}

// InformMachinesOfStateChanges will inform other State machines of
// the various event changes made when rx machine receives a packet
func (rxm *LacpRxMachine) InformMachinesOfStateChanges() {
	p := rxm.p

	if p.MuxMachineFsm != nil &&
		p.PtxMachineFsm != nil &&
		p.TxMachineFsm != nil {

		if p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDetached ||
			p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCDetached {
			p.checkConfigForSelection()
		}

		// lets inform the MUX of a possible State change
		if LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
			if p.aggSelected == LacpAggSelected {
				if p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateAttached ||
					p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCAttached {
					p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
						E:   LacpMuxmEventSelectedEqualSelectedAndPartnerSync,
						Src: RxMachineModuleStr}
				} else if p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCollecting {
					p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
						E:   LacpMuxmEventSelectedEqualSelectedPartnerSyncCollecting,
						Src: RxMachineModuleStr}
				}
			}
		} else if !LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) &&
			(p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDistributing ||
				p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCollecting) {
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventNotPartnerSync,
				Src: RxMachineModuleStr}

		} else if !LacpStateIsSet(p.PartnerOper.State, LacpStateCollectingBit) &&
			p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDistributing {
			p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
				E:   LacpMuxmEventNotPartnerCollecting,
				Src: RxMachineModuleStr}
		}

		// if we were in no periodic state because both ends were in passive
		// mode
		if p.TxMachineFsm.Machine.Curr.CurrentState() == LacpPtxmStateNoPeriodic &&
			p.begin &&
			p.lacpEnabled &&
			LacpStateIsSet(p.PartnerOper.State, LacpStateActivityBit) {
			// peer changed to active mode
			p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
				E:   LacpPtxmEventUnconditionalFallthrough,
				Src: RxMachineModuleStr}
		} else if LacpStateIsSet(p.PartnerOper.State, LacpStateTimeoutBit) &&
			p.PtxMachineFsm.PeriodicTxTimerInterval == LacpSlowPeriodicTime {
			p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
				E:   LacpPtxmEventPartnerOperStateTimeoutShort,
				Src: RxMachineModuleStr}
		} else if !LacpStateIsSet(p.PartnerOper.State, LacpStateTimeoutBit) &&
			p.PtxMachineFsm.PeriodicTxTimerInterval == LacpFastPeriodicTime {
			p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
				E:   LacpPtxmEventPartnerOperStateTimeoutLong,
				Src: RxMachineModuleStr}
		}

		// lets inform the PTX machine of change as this is an indication of
		// no tx packets, case should occur on first bring up when transmission
		// is based on admin provisioning.  Peer should respond to initial messages
		if !LacpStateIsSet(p.ActorOper.State, LacpStateActivityBit) &&
			!LacpStateIsSet(p.PartnerOper.State, LacpStateActivityBit) &&
			p.PtxMachineFsm != nil {
			p.PtxMachineFsm.PtxmEvents <- utils.MachineEvent{
				E:   LacpPtxmEventActorPartnerOperActivityPassiveMode,
				Src: RxMachineModuleStr}
		}
	}

}

func LacpRxMachineFSMBuild(p *LaAggPort) *LacpRxMachine {

	RxMachineStrStateMapCreate()

	rules := fsm.Ruleset{}

	// Instantiate a new LacpRxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	rxm := NewLacpRxMachine(p)

	//BEGIN -> INIT
	rules.AddRule(LacpRxmStateNone, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	rules.AddRule(LacpRxmStatePortDisabled, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	rules.AddRule(LacpRxmStateExpired, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	rules.AddRule(LacpRxmStateLacpDisabled, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	rules.AddRule(LacpRxmStateDefaulted, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	rules.AddRule(LacpRxmStateCurrent, LacpRxmEventBegin, rxm.LacpRxMachineInitialize)
	// INIT -> PORT_DISABLE
	rules.AddRule(LacpRxmStateInitialize, LacpRxmEventUnconditionalFallthrough, rxm.LacpRxMachinePortDisabled)
	// NOT PORT ENABLED  && NOT PORT MOVED
	// All States transition to this State
	rules.AddRule(LacpRxmStateInitialize, LacpRxmEventNotPortEnabledAndNotPortMoved, rxm.LacpRxMachinePortDisabled)
	rules.AddRule(LacpRxmStateExpired, LacpRxmEventNotPortEnabledAndNotPortMoved, rxm.LacpRxMachinePortDisabled)
	rules.AddRule(LacpRxmStateLacpDisabled, LacpRxmEventNotPortEnabledAndNotPortMoved, rxm.LacpRxMachinePortDisabled)
	rules.AddRule(LacpRxmStateDefaulted, LacpRxmEventNotPortEnabledAndNotPortMoved, rxm.LacpRxMachinePortDisabled)
	rules.AddRule(LacpRxmStateCurrent, LacpRxmEventNotPortEnabledAndNotPortMoved, rxm.LacpRxMachinePortDisabled)
	// PORT MOVED -> INIT
	rules.AddRule(LacpRxmStatePortDisabled, LacpRxmEventPortMoved, rxm.LacpRxMachineInitialize)
	// PORT ENABLED && LACP ENABLED
	rules.AddRule(LacpRxmStatePortDisabled, LacpRxmEventPortEnabledAndLacpEnabled, rxm.LacpRxMachineExpired)
	// PORT ENABLED && LACP DISABLED
	rules.AddRule(LacpRxmStatePortDisabled, LacpRxmEventPortEnabledAndLacpDisabled, rxm.LacpRxMachineLacpDisabled)
	// CURRENT WHILE TIMER EXPIRED
	rules.AddRule(LacpRxmStateExpired, LacpRxmEventCurrentWhileTimerExpired, rxm.LacpRxMachineDefaulted)
	rules.AddRule(LacpRxmStateCurrent, LacpRxmEventCurrentWhileTimerExpired, rxm.LacpRxMachineExpired)
	// LACP ENABLED
	rules.AddRule(LacpRxmStateLacpDisabled, LacpRxmEventLacpEnabled, rxm.LacpRxMachinePortDisabled)
	// PKT RX
	rules.AddRule(LacpRxmStateExpired, LacpRxmEventLacpPktRx, rxm.LacpRxMachineCurrent)
	rules.AddRule(LacpRxmStateDefaulted, LacpRxmEventLacpPktRx, rxm.LacpRxMachineCurrent)
	rules.AddRule(LacpRxmStateCurrent, LacpRxmEventLacpPktRx, rxm.LacpRxMachineCurrent)

	// Create a new FSM and apply the rules
	rxm.Apply(&rules)

	return rxm
}

// LacpRxMachineMain:  802.1ax-2014 Table 6-18
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpRxMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.12 Receive Machine
	rxm := LacpRxMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	rxm.Machine.Start(rxm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpRxMachine) {
		m.LacpRxmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			// lets set the current state
			m.p.AggPortDebug.AggPortDebugRxState = int(m.Machine.Curr.CurrentState())
			select {

			case <-m.currentWhileTimer.C:
				// special case if we have pending packets in the queue
				// by the time this expires we want to ensure the packet
				// gets processed first as this will clear/restart the timer
				if len(m.RxmPktRxEvent) == 0 {
					m.LacpRxmLog("Current While Timer Expired")
					m.Machine.ProcessEvent(RxMachineModuleStr, LacpRxmEventCurrentWhileTimerExpired, nil)
				}

			case event, ok := <-m.RxmEvents:
				if ok {
					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)
					if rv == nil {
						p := m.p
						/* continue State transition */
						if m.Machine.Curr.CurrentState() == LacpRxmStateInitialize {
							rv = m.Machine.ProcessEvent(RxMachineModuleStr, LacpRxmEventUnconditionalFallthrough, nil)
						}
						if rv == nil {
							m.LacpRxmLog(fmt.Sprintln("Port Enabled, LacpEnabled, State", p.PortEnabled, p.lacpEnabled, m.Machine.Curr.CurrentState()))
							if m.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled {
								if p.lacpEnabled &&
									p.IsPortEnabled() {
									rv = m.Machine.ProcessEvent(RxMachineModuleStr, LacpRxmEventPortEnabledAndLacpEnabled, nil)
								} else if !p.lacpEnabled &&
									p.IsPortEnabled() {
									rv = m.Machine.ProcessEvent(RxMachineModuleStr, LacpRxmEventPortEnabledAndLacpDisabled, nil)
								} else if p.portMoved {
									rv = m.Machine.ProcessEvent(RxMachineModuleStr, LacpRxmEventPortMoved, nil)
								}
							}
						}
					}

					if rv != nil {
						m.LacpRxmLog(strings.Join([]string{error.Error(rv), event.Src, RxmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					}

					// respond to caller if necessary so that we don't have a deadlock
					if event.ResponseChan != nil {
						utils.SendResponse(RxMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.LacpRxmLog("Machine End")
					return
				}
			case rx, ok := <-m.RxmPktRxEvent:
				if ok {
					//m.LacpRxmLog(fmt.Sprintf("RXM: received packet %d %s", m.p.PortNum, rx.src))
					// lets check if the port has moved
					p.LacpCounter.AggPortStatsLACPDUsRx += 1

					// centisecond
					p.AggPortDebug.AggPortDebugLastRxTime = (time.Now().Nanosecond() - LacpStartTime.Nanosecond()) / 10

					if m.CheckPortMoved(&p.PartnerOper, &(rx.pdu.Actor.Info)) {
						m.LacpRxmLog("port moved")
						m.p.portMoved = true
						m.Machine.ProcessEvent(RxModuleStr, LacpRxmEventPortMoved, nil)
					} else {
						// If you rx a packet must be in one
						// of 3 States
						// Expired/Defaulted/Current. each
						// State will transition to current
						// all other States should be ignored.
						m.Machine.ProcessEvent(RxModuleStr, LacpRxmEventLacpPktRx, rx.pdu)
					}

					// respond to caller if necessary so that we don't have a deadlock
					if rx.responseChan != nil {
						utils.SendResponse(RxMachineModuleStr, rx.responseChan)
					}
				}

			case ena := <-m.RxmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)

			}
		}
	}(rxm)
}

// handleRxFrame:
// TBD: First entry point of the raw ethernet frame
//func handleRxFrame(port int, pdu []bytes) {

// TODO
//	lacp := LacpPdu()
//	err := binary.Read(pdu, binary.BigEndian, &lacp)
//	if err != nil {
//		panic(err)
//	}
//}

// recordPDU: 802.1ax Section 6.4.9
//
// Record actor informatio from the packet
// Clear Defaulted Actor Operational State
// Determine Partner Operational Sync State
func (rxm *LacpRxMachine) recordPDU(lacpPduInfo *layers.LACP) {

	p := rxm.p
	collDistMap := map[fsm.State]bool{
		LacpMuxmStateCollecting:             true,
		LacpMuxmStateDistributing:           true,
		LacpMuxStateCCollectingDistributing: true,
	}

	//rxm.LacpRxmLog(fmt.Sprintf("recordPDU: %#v", lacpPduInfo))
	// Record Actor info from packet - store in parter operational
	// Port Number, Port Priority, System, System Priority
	// Key, State variables
	LacpCopyLacpPortInfoFromPkt(&lacpPduInfo.Actor.Info, &p.PartnerOper)

	// Set Actor Oper port State Defaulted to FALSE
	//rxm.LacpRxmLog("Clearing Defaulted Bit")
	LacpStateClear(&p.ActorOper.State, LacpStateDefaultedBit)

	// Set Partner Oper port State Sync State to
	// TRUE if the (1) or (2) is true:
	//
	// 1) Rx pdu: (Partner Port, Partner Port Priority, Partner
	// System, Partner System Priority, Partner Key,
	// Partner State Aggregation) vs 	cooresponding Operational
	// parameters of the Actor and Actor State Sync is TRUE and (3)
	//
	// 2) Rx pdu: Value of Actor State aggregation is FALSE
	// (indicates individual link) and Actor State sync is TRUE
	// and (3)
	//
	// 3) Rx pdu: Actor State LACP_Activity is TRUE
	// or both Actor Oper Port State LACP_Activity and PDU Partner
	// Partner State LACP_Activity is TRUE

	//rxm.LacpRxmLog(fmt.Sprintf("Pkt Partner State %s, Pkt Actor State %s, ActorOper State %s",
	//	LacpStateToStr(lacpPduInfo.Partner.Info.State),
	//	LacpStateToStr(lacpPduInfo.Actor.Info.State),
	//	LacpStateToStr(p.ActorOper.State)))
	// (1)
	if ((LacpLacpPktPortInfoIsEqual(&lacpPduInfo.Partner.Info, &p.ActorOper, LacpStateAggregationBit) &&
		LacpStateIsSet(lacpPduInfo.Actor.Info.State, LacpStateSyncBit)) ||
		//(2)
		(!LacpStateIsSet(lacpPduInfo.Actor.Info.State, LacpStateAggregationBit) &&
			LacpStateIsSet(lacpPduInfo.Actor.Info.State, LacpStateSyncBit))) &&
		// (3)
		(LacpStateIsSet(lacpPduInfo.Actor.Info.State, LacpStateActivityBit) ||
			(LacpStateIsSet(p.ActorOper.State, LacpStateActivityBit) &&
				LacpStateIsSet(lacpPduInfo.Partner.Info.State, LacpStateActivityBit))) {
		if !LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
			rxm.LacpRxmLog("Setting Partner Sync Bit")
			LacpStateSet(&p.PartnerOper.State, LacpStateSyncBit)
			if p.PCdMachineFsm != nil {
				// inform partner cdm
				p.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
					E:   LacpCdmEventPartnerOperPortStateSyncOn,
					Src: RxMachineModuleStr}
			}
			// NOTE Mux will be informed at a later time
		}
	} else {
		if LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
			rxm.LacpRxmLog("Clearing Sync Bit")

			LacpStateClear(&p.PartnerOper.State, LacpStateSyncBit)
			// inform mux of State change
			if p.MuxMachineFsm != nil {
				_, ok := collDistMap[p.MuxMachineFsm.Machine.Curr.CurrentState()]
				if ok {
					p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
						E:   LacpMuxmEventNotPartnerSync,
						Src: RxMachineModuleStr}
				}
			}
		}
	}

	// Optional to validate length of the following:
	// actor, partner, collector
}

// recordDefault: 802.1ax Section 6.4.9
//
// records the default parameter values for the
// partner carried in the partner admin parameters
// (Partner Admin Port Number, Partner Admin Port Priority,
//  Partner Admin System, Partner Admin System Priority,
// Partner Admin Key, and Partner Admin Port State) as the
// current Partner operational parameter values.  Sets Actor
// Oper Port State Default to TRUE and Partner Oper Port State
// Sync to TRUE
func (rxm *LacpRxMachine) recordDefault() {

	p := rxm.p

	LacpCopyLacpPortInfo(&p.partnerAdmin, &p.PartnerOper)
	//rxm.LacpRxmLog("Setting Actor Defaulted Bit")
	LacpStateSet(&p.ActorOper.State, LacpStateDefaultedBit)
	if !LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
		//rxm.LacpRxmLog("Setting Partner Sync Bit")
		LacpStateSet(&p.PartnerOper.State, LacpStateSyncBit)
		// inform partner cdm
		if p.PCdMachineFsm != nil {
			p.PCdMachineFsm.CdmEvents <- utils.MachineEvent{
				E:   LacpCdmEventPartnerOperPortStateSyncOn,
				Src: RxMachineModuleStr}
		}
	}

	if p.MuxMachineFsm != nil &&
		(p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateAttached ||
			p.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateCAttached) &&
		p.aggSelected == LacpAggSelected {
		p.MuxMachineFsm.MuxmEvents <- utils.MachineEvent{
			E:   LacpMuxmEventSelectedEqualSelectedAndPartnerSync,
			Src: RxMachineModuleStr}
	}
}

// updateNTT: 802.1ax Section 6.4.9
//
// Compare that the newly received PDU partner
// info agrees with the local port oper State.
// If it does not agree then set the NTT flag
// such that the Tx machine generates LACPDU
// Activity and Timeout are configurable so need to check that these have
// not changed
func (rxm *LacpRxMachine) updateNTT(lacpPduInfo *layers.LACP) bool {

	p := rxm.p

	const nttStateCompare uint8 = (LacpStateAggregationBit | LacpStateSyncBit)

	if !LacpLacpPktPortInfoIsEqual(&lacpPduInfo.Partner.Info, &p.ActorOper, nttStateCompare) {
		rxm.LacpRxmLog(fmt.Sprintf("PDU/Oper info different: \npdu: %#v state %s\n oper: %#v state %s",
			lacpPduInfo.Partner.Info,
			LacpStateToStr(lacpPduInfo.Partner.Info.State),
			p.ActorOper,
			LacpStateToStr(p.ActorOper.State)))
		// send event to user port and partner info don't
		utils.ProcessLacpPortPartnerInfoMismatch(int32(p.PortNum))

		p.LacpCounter.AggPortStateMissMatchInfoRx++
		return true
	} else if (LacpStateIsSet(lacpPduInfo.Partner.Info.State, LacpStateTimeoutBit) && !LacpStateIsSet(p.ActorOper.State, LacpStateTimeoutBit)) ||
		(!LacpStateIsSet(lacpPduInfo.Partner.Info.State, LacpStateTimeoutBit) && LacpStateIsSet(p.ActorOper.State, LacpStateTimeoutBit)) ||
		(LacpStateIsSet(lacpPduInfo.Partner.Info.State, LacpStateActivityBit) && !LacpStateIsSet(p.ActorOper.State, LacpStateActivityBit)) ||
		(!LacpStateIsSet(lacpPduInfo.Partner.Info.State, LacpStateActivityBit) && LacpStateIsSet(p.ActorOper.State, LacpStateActivityBit)) {
		rxm.LacpRxmLog(fmt.Sprintf("PDU/Oper state Timeout/Activity different: \npdu: %#v\n oper: %#v", lacpPduInfo.Partner.Info, p.ActorOper))
		p.LacpCounter.AggPortStateMissMatchInfoRx++
		utils.ProcessLacpPortPartnerInfoMismatch(int32(p.PortNum))
		return true
	}
	return false
}

func (rxm *LacpRxMachine) recordVersionNumber(lacpPduInfo *layers.LACP) {

	p := rxm.p

	p.partnerVersion = uint8(lacpPduInfo.Version)
}

// currentWhileTimerValid checks the State against
// the Actor Port Oper State Timeout
func (rxm *LacpRxMachine) CurrentWhileTimerValid() (time.Duration, bool) {

	p := rxm.p
	if rxm.currentWhileTimerTimeout == LacpShortTimeoutTime &&
		!LacpStateIsSet(p.ActorOper.State, LacpStateTimeoutBit) {
		rxm.LacpRxmLog("Current While Timer invalid adjusting to LONG TIMEOUT")
		return LacpLongTimeoutTime, false
	}
	if rxm.currentWhileTimerTimeout == LacpLongTimeoutTime &&
		LacpStateIsSet(p.ActorOper.State, LacpStateTimeoutBit) {
		rxm.LacpRxmLog("Current While Timer invalid adjusting to SHORT TIMEOUT")
		return LacpShortTimeoutTime, false
	}
	return 0, true
}

func (rxm *LacpRxMachine) CheckPortMoved(PartnerOper *LacpPortInfo, pktActor *layers.LACPPortInfo) bool {
	return rxm.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled &&
		PartnerOper.port == pktActor.Port &&
		reflect.DeepEqual(PartnerOper.System.Actor_System, pktActor.System.SystemId) &&
		PartnerOper.System.Actor_System_priority == pktActor.System.SystemPriority
}
