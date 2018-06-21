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

// 17.23 Port Receive state machine
package stp

import (
	"fmt"
	//"time"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

const PrxmMachineModuleStr = "PRXM"

const (
	PrxmStateNone = iota + 1
	PrxmStateDiscard
	PrxmStateReceive
)

var PrxmStateStrMap map[fsm.State]string

func PrxmMachineStrStateMapInit() {
	PrxmStateStrMap = make(map[fsm.State]string)
	PrxmStateStrMap[PrxmStateNone] = "None"
	PrxmStateStrMap[PrxmStateDiscard] = "Discard"
	PrxmStateStrMap[PrxmStateReceive] = "Receive"
}

const (
	PrxmEventBegin = iota + 1
	PrxmEventRcvdBpduAndNotPortEnabled
	PrxmEventEdgeDelayWhileNotEqualMigrateTimeAndNotPortEnabled
	PrxmEventRcvdBpduAndPortEnabled
	PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg
)

type RxBpduPdu struct {
	pdu          interface{}
	ptype        BPDURxType
	src          string
	responseChan chan string
}

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type PrxmMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PrxmEvents chan MachineEvent
	// rx pkt
	PrxmRxBpduPkt chan RxBpduPdu

	// enable logging
	PrxmLogEnableEvent chan bool
}

func (m *PrxmMachine) GetCurrStateStr() string {
	return PrxmStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PrxmMachine) GetPrevStateStr() string {
	return PrxmStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewStpPrxmMachine(p *StpPort) *PrxmMachine {
	prxm := &PrxmMachine{
		p:                  p,
		PrxmEvents:         make(chan MachineEvent, 50),
		PrxmRxBpduPkt:      make(chan RxBpduPdu, 50),
		PrxmLogEnableEvent: make(chan bool)}

	p.PrxmMachineFsm = prxm

	return prxm
}

func (prxm *PrxmMachine) PrxmLogger(s string) {
	StpMachineLogger("DEBUG", PrtMachineModuleStr, prxm.p.IfIndex, prxm.p.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (prxm *PrxmMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if prxm.Machine == nil {
		prxm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	prxm.Machine.Rules = r
	prxm.Machine.Curr = &StpStateEvent{
		strStateMap: PrxmStateStrMap,
		logEna:      true,
		logger:      prxm.PrxmLogger,
		owner:       PrxmMachineModuleStr,
		ps:          PrxmStateNone,
		s:           PrxmStateNone,
	}

	return prxm.Machine
}

// Stop should clean up all resources
func (prxm *PrxmMachine) Stop() {

	close(prxm.PrxmEvents)
	close(prxm.PrxmRxBpduPkt)
	close(prxm.PrxmLogEnableEvent)
}

// PrmMachineDiscard
func (prxm *PrxmMachine) PrxmMachineDiscard(m fsm.Machine, data interface{}) fsm.State {
	p := prxm.p
	p.RcvdBPDU = false
	p.RcvdRSTP = false
	p.RcvdSTP = false
	defer p.NotifyRcvdMsgChanged(PrxmMachineModuleStr, p.RcvdMsg, false, data)
	p.RcvdMsg = false
	// set to RSTP performance paramters Migrate Time
	p.EdgeDelayWhileTimer.count = MigrateTimeDefault
	return PrxmStateDiscard
}

// LacpPtxMachineFastPeriodic sets the periodic transmission time to fast
// and starts the timer
func (prxm *PrxmMachine) PrxmMachineReceive(m fsm.Machine, data interface{}) fsm.State {

	p := prxm.p

	// save off the bpdu info
	p.SaveMsgRcvInfo(data)

	//17.23
	// Decoding has been done as part of the Rx logic was a means of filtering
	rcvdMsg := prxm.UpdtBPDUVersion(data)
	// Figure 17-12
	defer p.NotifyRcvdMsgChanged(PrxmMachineModuleStr, p.RcvdMsg, rcvdMsg, data)
	p.RcvdMsg = rcvdMsg

	/* do not transition to NOT OperEdge if AdminEdge is set */
	if !p.AdminEdge && p.AutoEdgePort {
		defer p.NotifyOperEdgeChanged(PrxmMachineModuleStr, p.OperEdge, false)
		p.OperEdge = false
	}

	// Not setting this as it will conflict with bridge assurance / BPDU Guard
	//p.OperEdge = false
	p.RcvdBPDU = false
	p.EdgeDelayWhileTimer.count = MigrateTimeDefault

	return PrxmStateReceive
}

func PrxmMachineFSMBuild(p *StpPort) *PrxmMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrxmMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	prxm := NewStpPrxmMachine(p)

	//BEGIN -> DISCARD
	rules.AddRule(PrxmStateNone, PrxmEventBegin, prxm.PrxmMachineDiscard)
	rules.AddRule(PrxmStateDiscard, PrxmEventBegin, prxm.PrxmMachineDiscard)
	rules.AddRule(PrxmStateReceive, PrxmEventBegin, prxm.PrxmMachineDiscard)

	// RX BPDU && PORT NOT ENABLED	 -> DISCARD
	rules.AddRule(PrxmStateDiscard, PrxmEventRcvdBpduAndNotPortEnabled, prxm.PrxmMachineDiscard)
	rules.AddRule(PrxmStateReceive, PrxmEventRcvdBpduAndNotPortEnabled, prxm.PrxmMachineDiscard)

	// EDGEDELAYWHILE != MIGRATETIME && PORT NOT ENABLED -> DISCARD
	rules.AddRule(PrxmStateDiscard, PrxmEventEdgeDelayWhileNotEqualMigrateTimeAndNotPortEnabled, prxm.PrxmMachineDiscard)
	rules.AddRule(PrxmStateReceive, PrxmEventEdgeDelayWhileNotEqualMigrateTimeAndNotPortEnabled, prxm.PrxmMachineDiscard)

	// RX BPDU && PORT ENABLED -> RECEIVE
	rules.AddRule(PrxmStateDiscard, PrxmEventRcvdBpduAndPortEnabled, prxm.PrxmMachineReceive)

	// RX BPDU && PORT ENABLED && NOT RCVDMSG
	rules.AddRule(PrxmStateReceive, PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg, prxm.PrxmMachineReceive)

	// Create a new FSM and apply the rules
	prxm.Apply(&rules)

	return prxm
}

// PrxmMachineMain:
func (p *StpPort) PrxmMachineMain() {

	// Build the State machine for STP Receive Machine according to
	// 802.1d Section 17.23
	prxm := PrxmMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	prxm.Machine.Start(prxm.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PrxmMachine) {
		StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {

			case event, ok := <-m.PrxmEvents:

				if ok {
					if m.Machine.Curr.CurrentState() == PrxmStateNone && event.e != PrxmEventBegin {
						m.PrxmEvents <- event
						break
					}

					//fmt.Println("Event Rx", event.src, event.e)
					rv := m.Machine.ProcessEvent(event.src, event.e, nil)
					if rv != nil {
						StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s state[%s]event[%d]\n", rv, PrxmStateStrMap[m.Machine.Curr.CurrentState()], event.e))
					} else {
						// for faster state transitions
						m.ProcessPostStateProcessing(event.data)
					}

					if event.responseChan != nil {
						SendResponse(PrxmMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}
			case rx := <-m.PrxmRxBpduPkt:

				if m.Machine.Curr.CurrentState() == PrxmStateNone {
					continue
				}

				if p.BpduGuard &&
					p.AdminEdge {
					if p.BPDUGuardTimer.count != 0 {
						p.BPDUGuardTimer.count = p.BpduGuardInterval
						for _, client := range GetAsicDPluginList() {
							client.BPDUGuardDetected(p.IfIndex, true)
						}
					} else {
						p.BPDUGuardTimer.count = p.BpduGuardInterval
					}
				} else {

					//fmt.Println("Event PKT Rx", p.IfIndex, p.BrgIfIndex, rx.src, PrxmStateStrMap[m.Machine.Curr.CurrentState()], rx.ptype, p.RcvdMsg, p.PortEnabled)
					if m.Machine.Curr.CurrentState() == PrxmStateDiscard {
						if p.PortEnabled {
							rv := m.Machine.ProcessEvent("RX MODULE", PrxmEventRcvdBpduAndPortEnabled, rx)
							if rv != nil {
								StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s state[%s]event[%d]\n", rv, PrxmStateStrMap[m.Machine.Curr.CurrentState()], PrxmEventRcvdBpduAndPortEnabled))
							} else {
								// for faster state transitions
								m.ProcessPostStateProcessing(rx)
							}
						} else {
							rv := m.Machine.ProcessEvent("RX MODULE", PrxmEventRcvdBpduAndNotPortEnabled, rx)
							if rv != nil {
								StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s state[%s]event[%d]\n", rv, PrxmStateStrMap[m.Machine.Curr.CurrentState()], PrxmEventRcvdBpduAndPortEnabled))
							} else {
								// for faster state transitions
								m.ProcessPostStateProcessing(rx)
							}
						}
					} else {
						if p.PortEnabled &&
							!p.RcvdMsg {
							rv := m.Machine.ProcessEvent("RX MODULE", PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg, rx)
							if rv != nil {
								StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s state[%s]event[%d]\n", rv, PrxmStateStrMap[m.Machine.Curr.CurrentState()], PrxmEventRcvdBpduAndPortEnabled))
							} else {
								// for faster state transitions
								m.ProcessPostStateProcessing(rx)
							}
						} else if !p.PortEnabled {
							rv := m.Machine.ProcessEvent("RX MODULE", PrxmEventRcvdBpduAndNotPortEnabled, rx)
							if rv != nil {
								StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s state[%s]event[%d]\n", rv, PrxmStateStrMap[m.Machine.Curr.CurrentState()], PrxmEventRcvdBpduAndPortEnabled))
							} else {
								// for faster state transitions
								m.ProcessPostStateProcessing(rx)
							}
						}
					}
				}
				p.SetRxPortCounters(rx.ptype)
			case ena := <-m.PrxmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(prxm)
}

func (prxm *PrxmMachine) ProcessPostStateDiscard(data interface{}) {
	p := prxm.p
	if prxm.Machine.Curr.CurrentState() == PrxmStateDiscard &&
		p.RcvdBPDU &&
		p.PortEnabled {
		rv := prxm.Machine.ProcessEvent(PrxmMachineModuleStr, PrxmEventRcvdBpduAndPortEnabled, data)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s\n", rv))
		} else {
			prxm.ProcessPostStateProcessing(data)
		}
	}
}

func (prxm *PrxmMachine) ProcessPostStateReceive(data interface{}) {
	p := prxm.p
	if prxm.Machine.Curr.CurrentState() == PrxmStateReceive &&
		p.RcvdBPDU &&
		p.PortEnabled &&
		!p.RcvdMsg {
		rv := prxm.Machine.ProcessEvent(PrxmMachineModuleStr, PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg, data)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s\n", rv))
		} else {
			prxm.ProcessPostStateProcessing(data)
		}
	}
}

func (prxm *PrxmMachine) ProcessPostStateProcessing(data interface{}) {
	// post processing
	prxm.ProcessPostStateDiscard(data)
	prxm.ProcessPostStateReceive(data)
}

// UpdtBPDUVersion:  17.21.22
// This function will also inform Port Migration of
// BPDU type rcvdRSTP or rcvdSTP Figure 17-12
func (prxm *PrxmMachine) UpdtBPDUVersion(data interface{}) bool {
	validPdu := false
	p := prxm.p
	bpdumsg := data.(RxBpduPdu)
	bpduLayer := bpdumsg.pdu
	flags := uint8(0)
	StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("UpdtBPDUVersion: pbduType %#v", bpduLayer))
	switch bpduLayer.(type) {
	case *layers.RSTP:
		// 17.21.22
		// some checks a bit redundant as the layers class has already validated
		// the BPDUType, but for completness going to add the check anyways
		rstp := bpduLayer.(*layers.RSTP)
		flags = uint8(rstp.Flags)
		if rstp.ProtocolVersionId == layers.RSTPProtocolVersion &&
			rstp.BPDUType == layers.BPDUTypeRSTP {
			// Inform the Port Protocol Migration STate machine
			// that we have received a RSTP packet when we were previously
			// sending non-RSTP
			if !p.RcvdRSTP &&
				!p.SendRSTP &&
				p.RstpVersion {
				if p.PpmmMachineFsm != nil {
					p.PpmmMachineFsm.PpmmEvents <- MachineEvent{
						e:    PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
						data: bpduLayer,
						src:  PrxmMachineModuleStr}
				}
			}
			// lets reset the timer as we have received an rstp frame
			p.MdelayWhiletimer.count = MigrateTimeDefault

			p.RcvdRSTP = true
			validPdu = true
		}

		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received RSTP packet flags rcvdRSTP[%t] sendRSTP[%t]", rstp.Flags, p.RcvdRSTP, p.SendRSTP))

		defer p.NotifyRcvdTcRcvdTcnRcvdTcAck(p.RcvdTc, p.RcvdTcn, p.RcvdTcAck, StpGetBpduTopoChange(flags), false, false)
		p.RcvdTc = StpGetBpduTopoChange(flags)
		p.RcvdTcn = false
		p.RcvdTcAck = StpGetBpduTopoChangeAck(flags)

		if p.RcvdTc {
			p.SetRxPortCounters(BPDURxTypeTopo)
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC packet"))
		}
		if p.RcvdTcAck {
			p.SetRxPortCounters(BPDURxTypeTopoAck)
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC Ack packet "))
		}

	case *layers.PVST:
		// 17.21.22
		// some checks a bit redundant as the layers class has already validated
		// the BPDUType, but for completness going to add the check anyways
		pvst := bpduLayer.(*layers.PVST)
		flags = uint8(pvst.Flags)
		//fmt.Println("PVST: Protocol version, bpdu type", pvst.ProtocolVersionId, pvst.BPDUType)
		if pvst.ProtocolVersionId == layers.RSTPProtocolVersion &&
			pvst.BPDUType == layers.BPDUTypeRSTP {
			// Inform the Port Protocol Migration STate machine
			// that we have received a RSTP packet when we were previously
			// sending non-RSTP
			if !p.RcvdRSTP &&
				!p.SendRSTP &&
				p.RstpVersion {
				if p.PpmmMachineFsm != nil {
					p.PpmmMachineFsm.PpmmEvents <- MachineEvent{
						e:    PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
						data: bpduLayer,
						src:  PrxmMachineModuleStr}
				}
			}
			// lets reset the timer as we have received an rstp frame
			p.MdelayWhiletimer.count = MigrateTimeDefault

			p.RcvdRSTP = true
			validPdu = true
		}

		//StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received PVST packet flags", pvst.Flags))

		defer p.NotifyRcvdTcRcvdTcnRcvdTcAck(p.RcvdTc, p.RcvdTcn, p.RcvdTcAck, StpGetBpduTopoChange(flags), false, StpGetBpduTopoChangeAck(flags))
		p.RcvdTc = StpGetBpduTopoChange(flags)
		p.RcvdTcn = false
		p.RcvdTcAck = StpGetBpduTopoChangeAck(flags)

		if p.RcvdTc {
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC packet"))
			p.SetRxPortCounters(BPDURxTypeTopo)
		}
		if p.RcvdTcAck {
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC Ack packet"))
			p.SetRxPortCounters(BPDURxTypeTopoAck)
		}

	case *layers.STP:
		stp := bpduLayer.(*layers.STP)
		flags = uint8(stp.Flags)
		if stp.ProtocolVersionId == layers.STPProtocolVersion &&
			stp.BPDUType == layers.BPDUTypeSTP {

			// Inform the Port Protocol Migration State Machine
			// that we have received an STP packet when we were previously
			// sending RSTP
			// do not transition this to STP true until
			// mdelay while exires, this gives the far end enough
			// time to transition
			if p.MdelayWhiletimer.count == 0 {
				if p.SendRSTP {
					if p.PpmmMachineFsm != nil {
						p.PpmmMachineFsm.PpmmEvents <- MachineEvent{
							e:    PpmmEventSendRSTPAndRcvdSTP,
							data: bpduLayer,
							src:  PrxmMachineModuleStr}
					}
				}
			}
			//}

			p.RcvdSTP = true
			validPdu = true
		}

		StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received STP packet %#v", stp))
		defer p.NotifyRcvdTcRcvdTcnRcvdTcAck(p.RcvdTc, p.RcvdTcn, p.RcvdTcAck, StpGetBpduTopoChange(flags), false, StpGetBpduTopoChangeAck(flags))
		p.RcvdTc = StpGetBpduTopoChange(flags)
		p.RcvdTcn = false
		p.RcvdTcAck = StpGetBpduTopoChangeAck(flags)

		if p.RcvdTc {
			p.SetRxPortCounters(BPDURxTypeTopo)
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC packet"))
		}
		if p.RcvdTcAck {
			p.SetRxPortCounters(BPDURxTypeTopoAck)
			StpMachineLogger("DEBUG", PrxmMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Received TC Ack packet"))
		}

	case *layers.BPDUTopology:
		topo := bpduLayer.(*layers.BPDUTopology)
		if (topo.ProtocolVersionId == layers.STPProtocolVersion &&
			topo.BPDUType == layers.BPDUTypeTopoChange) ||
			(topo.ProtocolVersionId == layers.TCNProtocolVersion &&
				topo.BPDUType == layers.BPDUTypeTopoChange) {
			// Inform the Port Protocol Migration State Machine
			// that we have received an STP packet when we were previously
			// sending RSTP
			if p.MdelayWhiletimer.count == 0 {
				if p.SendRSTP {
					if p.PpmmMachineFsm != nil {
						p.PpmmMachineFsm.PpmmEvents <- MachineEvent{
							e:    PpmmEventSendRSTPAndRcvdSTP,
							data: bpduLayer,
							src:  PrxmMachineModuleStr}
					}
				}
				p.RcvdSTP = true
			}
			validPdu = true
			StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Received TCN packet")
			defer p.NotifyRcvdTcRcvdTcnRcvdTcAck(p.RcvdTc, p.RcvdTcn, p.RcvdTcAck, false, true, false)
			p.RcvdTc = false
			p.RcvdTcn = true
			p.RcvdTcAck = false
			if p.RcvdTc {
				p.SetRxPortCounters(BPDURxTypeTopo)
			}

		}
	}

	return validPdu
}
