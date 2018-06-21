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
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"strconv"
	"strings"
	"time"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

const TxMachineModuleStr = "Tx Machine"

const (
	LacpTxmStateNone = iota + 1
	LacpTxmStateOn
	LacpTxmStateOff
	LacpTxmStateDelayed
	LacpTxmStateGuardTimerExpire
)

var TxmStateStrMap map[fsm.State]string

func TxMachineStrStateMapCreate() {

	TxmStateStrMap = make(map[fsm.State]string)
	TxmStateStrMap[LacpTxmStateNone] = "LacpTxmStateNone"
	TxmStateStrMap[LacpTxmStateOn] = "LacpTxmStateOn"
	TxmStateStrMap[LacpTxmStateOff] = "LacpTxmStateOff"
	TxmStateStrMap[LacpTxmStateDelayed] = "LacpTxmStateDelayed"
	TxmStateStrMap[LacpTxmStateGuardTimerExpire] = "LacpTxmStateGuardTimerExpire"
}

const (
	LacpTxmEventBegin = iota + 1
	LacpTxmEventNtt
	LacpTxmEventGuardTimer
	LacpTxmEventDelayTx
	LacpTxmEventLacpDisabled
	LacpTxmEventLacpEnabled
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type LacpTxMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	// Port this Machine is associated with
	p *LaAggPort

	// number of frames that should be transmitted
	// after restriction logic has cleared
	txPending int

	// number of frames transmitted within guard timer interval
	txPkts int

	// ntt, this may be set by external applications
	// the State machine will only clear
	ntt bool

	// timer needed for 802.1ax-20014 section 6.4.16
	txGuardTimer *time.Timer

	// machine specific events
	TxmEvents         chan utils.MachineEvent
	TxmLogEnableEvent chan bool
}

// PrevState will get the previous State from the State transitions
func (txm *LacpTxMachine) PrevState() fsm.State { return txm.PreviousState }

// PrevStateSet will set the previous State
func (txm *LacpTxMachine) PrevStateSet(s fsm.State) { txm.PreviousState = s }

// Stop will stop all timers and close all channels
func (txm *LacpTxMachine) Stop() {
	txm.TxGuardTimerStop()

	close(txm.TxmEvents)
	close(txm.TxmLogEnableEvent)
}

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewLacpTxMachine(port *LaAggPort) *LacpTxMachine {
	txm := &LacpTxMachine{
		p:                 port,
		txPending:         0,
		txPkts:            0,
		ntt:               false,
		PreviousState:     LacpTxmStateNone,
		TxmEvents:         make(chan utils.MachineEvent, 1000),
		TxmLogEnableEvent: make(chan bool)}

	port.TxMachineFsm = txm

	// start then stop
	txm.TxGuardTimerStart()
	txm.TxGuardTimerStop()

	return txm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (txm *LacpTxMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if txm.Machine == nil {
		txm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	txm.Machine.Rules = r
	txm.Machine.Curr = &utils.StateEvent{
		StrStateMap: TxmStateStrMap,
		LogEna:      txm.p.logEna,
		Logger:      txm.LacpTxmLog,
		Owner:       TxMachineModuleStr,
	}

	return txm.Machine
}

// LacpTxMachineOn will either send a packet out or delay transmission of a
// packet
func (txm *LacpTxMachine) LacpTxMachineOn(m fsm.Machine, data interface{}) fsm.State {

	var nextState fsm.State
	p := txm.p

	txm.PrevStateSet(txm.Machine.Curr.CurrentState())

	nextState = LacpTxmStateOn

	// NTT must be set to tx
	if txm.ntt {
		// if more than 3 packets are being transmitted within time interval
		// delay transmission
		if txm.txPkts < 3 {
			if txm.txPkts == 0 {
				txm.TxGuardTimerStart()
			}
			txm.txPkts++

			lacp := &layers.LACP{
				Version: layers.LACPVersion1,
				Actor: layers.LACPInfoTlv{TlvType: layers.LACPTLVActorInfo,
					Length: layers.LACPActorTlvLength,
					Info: layers.LACPPortInfo{
						System: layers.LACPSystem{SystemId: p.ActorOper.System.Actor_System,
							SystemPriority: p.ActorOper.System.Actor_System_priority,
						},
						Key:     p.ActorOper.Key,
						PortPri: p.ActorOper.Port_pri,
						Port:    p.ActorOper.port,
						State:   p.ActorOper.State,
					},
				},
				Partner: layers.LACPInfoTlv{TlvType: layers.LACPTLVPartnerInfo,
					Length: layers.LACPActorTlvLength,
					Info: layers.LACPPortInfo{
						System: layers.LACPSystem{SystemId: p.PartnerOper.System.Actor_System,
							SystemPriority: p.PartnerOper.System.Actor_System_priority,
						},
						Key:     p.PartnerOper.Key,
						PortPri: p.PartnerOper.Port_pri,
						Port:    p.PartnerOper.port,
						State:   p.PartnerOper.State,
					},
				},
				Collector: layers.LACPCollectorInfoTlv{
					TlvType:  layers.LACPTLVCollectorInfo,
					Length:   layers.LACPCollectorTlvLength,
					MaxDelay: 0,
				},
			}

			// transmit the packet
			for _, ftx := range LaSysGlobalTxCallbackListGet(p) {
				//txm.LacpTxmLog(fmt.Sprintf("Sending Tx packet port %d pkts %d", p.PortNum, txm.txPkts))
				ftx(p.PortNum, lacp)
				p.LacpCounter.AggPortStatsLACPDUsTx += 1
			}
			// Version 2 consideration if enable_long_pdu_xmit and
			// LongLACPPDUTransmit are True:
			// LACPDU will be a Long LACPDU formatted by 802.1ax-2014 Section
			// 6.4.2 and including Port Conversation Mask TLV 6.4.2.4.3
			txm.ntt = false

			// lets force another transmit
			if txm.txPending > 0 && txm.txPkts < 3 {
				txm.txPending--
				txm.TxmEvents <- utils.MachineEvent{

					E:   LacpTxmEventNtt,
					Src: TxMachineModuleStr}
			}
		} else {
			if txm.txPending < 5 {
				txm.txPending++
			}
			txm.LacpTxmLog(fmt.Sprintf("ON: Delay packets %d", txm.txPending))
			nextState = LacpTxmStateDelayed
		}
	}
	return nextState
}

// LacpTxMachineDelayed is a State in which a packet is forced to transmit
// regardless of the ntt State
func (txm *LacpTxMachine) LacpTxMachineDelayed(m fsm.Machine, data interface{}) fsm.State {
	var State fsm.State

	txm.PrevStateSet(txm.Machine.Curr.CurrentState())

	State = LacpTxmStateOn

	// if more than 3 packets are being transmitted within time interval
	// Version 2 consideration if enable_long_pdu_xmit and
	// LongLACPPDUTransmit are True:
	// LACPDU will be a Long LACPDU formatted by 802.1ax-2014 Section
	// 6.4.2 and including Port Conversation Mask TLV 6.4.2.4.3
	txm.LacpTxmLog(fmt.Sprintf("Delayed: txPending %d txPkts %d delaying tx", txm.txPending, txm.txPkts))
	if txm.txPending > 0 && txm.txPkts > 3 {
		State = LacpTxmStateDelayed
		//txm.TxmEvents <- utils.MachineEvent{e: LacpTxmEventDelayTx,
		//	src: TxMachineModuleStr}
	} else {
		// transmit packet
		txm.txPending--
		txm.TxmEvents <- utils.MachineEvent{
			E:   LacpTxmEventNtt,
			Src: TxMachineModuleStr}
	}

	return State
}

// LacpTxMachineOff will ensure that no packets are transmitted, typically means that
// lacp has been disabled
func (txm *LacpTxMachine) LacpTxMachineOff(m fsm.Machine, data interface{}) fsm.State {
	txm.txPending = 0
	txm.txPkts = 0
	txm.ntt = false
	txm.TxGuardTimerStop()
	return LacpTxmStateOff
}

// LacpTxMachineGuard will clear the current transmited packet count and
// generate a new event to tx a new packet
func (txm *LacpTxMachine) LacpTxMachineGuard(m fsm.Machine, data interface{}) fsm.State {
	txm.txPkts = 0
	var State fsm.State

	State = LacpTxmStateOn
	if txm.txPending > 0 {
		State = LacpTxmStateGuardTimerExpire
	}

	// no State transition just need to clear the txPkts
	return State
}

// LacpTxMachineFSMBuild will build the State machine with callbacks
func LacpTxMachineFSMBuild(p *LaAggPort) *LacpTxMachine {

	rules := fsm.Ruleset{}

	TxMachineStrStateMapCreate()

	// Instantiate a new LacpRxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the initalize State
	txm := NewLacpTxMachine(p)

	//BEGIN -> TX OFF
	rules.AddRule(LacpTxmStateNone, LacpTxmEventBegin, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateOn, LacpTxmEventBegin, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateOff, LacpTxmEventBegin, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateDelayed, LacpTxmEventBegin, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateGuardTimerExpire, LacpTxmEventBegin, txm.LacpTxMachineOff)

	// NTT -> TX ON
	rules.AddRule(LacpTxmStateOn, LacpTxmEventNtt, txm.LacpTxMachineOn)
	rules.AddRule(LacpTxmStateGuardTimerExpire, LacpTxmEventNtt, txm.LacpTxMachineOn)
	rules.AddRule(LacpTxmStateDelayed, LacpTxmEventNtt, txm.LacpTxMachineOn)
	// DELAY -> TX DELAY
	rules.AddRule(LacpTxmStateOn, LacpTxmEventDelayTx, txm.LacpTxMachineDelayed)
	rules.AddRule(LacpTxmStateDelayed, LacpTxmEventDelayTx, txm.LacpTxMachineDelayed)
	// LACP ON -> TX ON
	rules.AddRule(LacpTxmStateOff, LacpTxmEventLacpEnabled, txm.LacpTxMachineOn)
	// LACP DISABLED -> TX OFF
	rules.AddRule(LacpTxmStateNone, LacpTxmEventLacpDisabled, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateOn, LacpTxmEventLacpDisabled, txm.LacpTxMachineOff)
	rules.AddRule(LacpTxmStateDelayed, LacpTxmEventLacpDisabled, txm.LacpTxMachineOff)
	// GUARD TIMER -> TX ON
	rules.AddRule(LacpTxmStateOn, LacpTxmEventGuardTimer, txm.LacpTxMachineGuard)
	rules.AddRule(LacpTxmStateDelayed, LacpTxmEventGuardTimer, txm.LacpTxMachineGuard)

	// Create a new FSM and apply the rules
	txm.Apply(&rules)

	return txm
}

// LacpRxMachineMain:  802.1ax-2014 Table 6-18
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *LaAggPort) LacpTxMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.13 Periodic Transmission Machine
	txm := LacpTxMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	txm.Machine.Start(txm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the RxMachine should handle.
	go func(m *LacpTxMachine) {
		m.LacpTxmLog("Machine Start")
		defer m.p.wg.Done()
		for {
			select {

			case event, ok := <-m.TxmEvents:
				if ok {
					//m.LacpTxmLog(fmt.Sprintf("Event rx %d %s %s", event.E, event.Src, TxmStateStrMap[m.Machine.Curr.CurrentState()]))
					// special case, another machine has a need to
					// transmit a packet
					if event.E == LacpTxmEventNtt {
						m.ntt = true
					}

					rv := m.Machine.ProcessEvent(event.Src, event.E, nil)

					if rv != nil {
						m.LacpTxmLog(strings.Join([]string{error.Error(rv), event.Src, TxmStateStrMap[m.Machine.Curr.CurrentState()], strconv.Itoa(int(event.E))}, ":"))
					} else {
						if m.Machine.Curr.CurrentState() == LacpTxmStateGuardTimerExpire &&
							m.txPending > 0 && m.txPkts == 0 {

							for m.txPending > 0 && m.txPkts < 3 {
								m.txPending--
								m.ntt = true
								m.LacpTxmLog(fmt.Sprintf("Forcing NTT processing from expire pending pkts %d\n", m.txPending))
								m.Machine.ProcessEvent(TxMachineModuleStr, LacpTxmEventNtt, nil)
							}
						}
					}

					if event.ResponseChan != nil {
						utils.SendResponse(TxMachineModuleStr, event.ResponseChan)
					}
				} else {
					m.LacpTxmLog("Machine End")
					return
				}
			case ena := <-m.TxmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(txm)
}

// LacpTxGuardGeneration will generate an event to the Tx Machine
// in order to clear the txPkts count
func (txm *LacpTxMachine) LacpTxGuardGeneration() {
	//txm.LacpTxmLog("LacpTxGuardGeneration")
	txm.TxmEvents <- utils.MachineEvent{
		E:   LacpTxmEventGuardTimer,
		Src: TxMachineModuleStr}
}
