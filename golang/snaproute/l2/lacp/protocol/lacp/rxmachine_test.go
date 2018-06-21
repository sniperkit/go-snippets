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

// lacp tests
// go test
// go test -coverageprofile lacpcov.out
// go tool cover -html=lacpcov.out
package lacp

import (
	"l2/lacp/protocol/utils"
	"net"
	"testing"
	"time"
	"utils/fsm"

	"github.com/google/gopacket/layers"
)

func TestLaAggPortRxMachineStateTransitions(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	var msg string
	var portchan chan string
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:     1,
		Prio:   0x80,
		IntfId: "SIMeth1.1",
		Key:    100,
	}

	utils.PortConfigMap[int32(pconf.Id)] = utils.PortConfig{Name: pconf.IntfId,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	// not calling Create because we don't want to launch all State machines
	p := NewLaAggPort(pconf)

	// lets start the Rx Machine only
	p.LacpRxMachineMain()

	// Rx Machine
	if p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateNone {
		t.Error("ERROR RX Machine State incorrect expected",
			LacpRxmStateNone, "actual",
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	p.BEGIN(true)
	portchan = p.PortChannelGet()

	// port is initally disabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateInitialize &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateInitialize,
			LacpRxmStatePortDisabled,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	// check State info
	if p.aggSelected != LacpAggUnSelected {
		t.Error("expected UNSELECTED", LacpAggUnSelected, "actual", p.aggSelected)
	}
	if LacpStateIsSet(p.ActorOper.State, LacpStateExpiredBit) {
		t.Error("expected State Expired to be cleared")
	}
	if p.portMoved != false {
		t.Error("expected port moved to be false")
	}
	// TODO check actor oper State

	p.portMoved = true
	// send PORT MOVED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventPortMoved,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// PORT MOVED
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateInitialize &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateInitialize,
			LacpRxmStatePortDisabled,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	p.aggSelected = LacpAggSelected
	p.portMoved = false
	p.PortEnabled = true
	p.LinkOperStatus = true
	p.lacpEnabled = false
	// send PORT ENABLED && LACP DISABLED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventPortEnabledAndLacpDisabled,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port is initally disabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStatePortDisabled &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateLacpDisabled {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStatePortDisabled,
			LacpRxmStateLacpDisabled,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	p.lacpEnabled = true
	p.LinkOperStatus = true
	// send LACP ENABLED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventLacpEnabled,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was lacp disabled, but then transitioned to port disabled
	// then expired
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStatePortDisabled &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateExpired {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStatePortDisabled,
			LacpRxmStateExpired,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	if LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
		t.Error("Expected partner Sync Bit to not be set")
	}
	if !LacpStateIsSet(p.PartnerOper.State, LacpStateTimeoutBit) {
		t.Error("Expected partner Timeout bit to be set since we are in short timeout")
	}
	if p.RxMachineFsm.currentWhileTimerTimeout != LacpShortTimeoutTime {
		t.Error("Expected timer to be set to short timeout")
	}
	if !LacpStateIsSet(p.ActorOper.State, LacpStateExpiredBit) {
		t.Error("Expected actor expired bit to be set")
	}

	p.PortEnabled = false
	p.LinkOperStatus = false
	// send NOT ENABLED AND NOT MOVED event to Rx Machine from Expired State
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventNotPortEnabledAndNotPortMoved,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateExpired &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStatePortDisabled,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	if LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
		t.Error("Expected partner Sync Bit to not be set")
	}

	p.PortEnabled = true
	p.LinkOperStatus = true
	p.lacpEnabled = false
	// send NOT ENABLED AND NOT MOVED event to Rx Machine from Expired State
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventPortEnabledAndLacpDisabled,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	p.PortEnabled = false
	p.LinkOperStatus = false
	// send NOT ENABLED AND NOT MOVED event to Rx Machine from LACP DISABLED
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventNotPortEnabledAndNotPortMoved,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateLacpDisabled &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateLacpDisabled,
			LacpRxmStatePortDisabled,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	p.PortEnabled = true
	p.LinkOperStatus = true
	p.lacpEnabled = true
	// send PORT ENABLE LACP ENABLED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventPortEnabledAndLacpEnabled,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// send CURRENT WHILE TIMER event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventCurrentWhileTimerExpired,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateExpired &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStateDefaulted,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	// TODO check default selected, record default, expired == false

	// LETS GET THE State BACK TO EXPIRED

	p.PortEnabled = false
	p.LinkOperStatus = false
	// send NOT PORT ENABLE NOT PORT MOVED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventNotPortEnabledAndNotPortMoved,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	p.PortEnabled = true
	p.LinkOperStatus = true
	// send PORT ENABLE LACP ENABLED event to Rx Machine
	p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            LacpRxmEventPortEnabledAndLacpEnabled,
		ResponseChan: portchan,
		Src:          "TEST"}

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// lets adjust the ActorOper timeout State
	// TODO Assume a method was called to adjust this
	LacpStateSet(&p.ActorAdmin.State, LacpStateTimeoutBit)
	LacpStateSet(&p.ActorOper.State, LacpStateTimeoutBit)

	//slow := &layers.SlowProtocol{
	//	SubType: layers.SlowProtocolTypeLACP,
	//}
	// send valid pdu
	lacppdu := &layers.LACP{
		Version: layers.LACPVersion2,
		Actor: layers.LACPInfoTlv{TlvType: layers.LACPTLVActorInfo,
			Length: layers.LACPActorTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: [6]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
					SystemPriority: 1},
				Key:     100,
				PortPri: 0x80,
				Port:    10,
				State:   LacpStateActivityBit | LacpStateAggregationBit | LacpStateTimeoutBit},
		},
		Partner: layers.LACPInfoTlv{TlvType: layers.LACPTLVPartnerInfo,
			Length: layers.LACPPartnerTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: p.ActorOper.System.Actor_System,
					SystemPriority: p.ActorOper.System.Actor_System_priority},
				Key:     p.Key,
				PortPri: p.portPriority,
				Port:    p.PortNum,
				State:   p.ActorOper.State},
		},
	}

	rx := LacpRxLacpPdu{
		pdu:          lacppdu,
		responseChan: portchan,
		src:          "TEST"}
	p.RxMachineFsm.RxmPktRxEvent <- rx

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateExpired &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateCurrent {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStateCurrent,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	// allow for current while timer to expire
	time.Sleep(time.Second * 4)

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateCurrent &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateExpired {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateCurrent,
			LacpRxmStateExpired,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	// allow for current while timer to expire
	time.Sleep(time.Second * 4)

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateExpired &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStateDefaulted,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	lacppdu = &layers.LACP{
		Version: layers.LACPVersion2,
		Actor: layers.LACPInfoTlv{TlvType: layers.LACPTLVActorInfo,
			Length: layers.LACPActorTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: [6]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
					SystemPriority: 1},
				Key:     100,
				PortPri: 0x80,
				Port:    10,
				State:   LacpStateActivityBit | LacpStateAggregationBit | LacpStateTimeoutBit},
		},
		Partner: layers.LACPInfoTlv{TlvType: layers.LACPTLVPartnerInfo,
			Length: layers.LACPPartnerTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: p.ActorOper.System.Actor_System,
					SystemPriority: p.ActorOper.System.Actor_System_priority},
				Key:     p.Key,
				PortPri: p.portPriority,
				Port:    p.PortNum,
				State:   p.ActorOper.State},
		},
	}

	rx = LacpRxLacpPdu{
		pdu:          lacppdu,
		responseChan: portchan,
		src:          "TEST"}
	p.RxMachineFsm.RxmPktRxEvent <- rx

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateDefaulted &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateCurrent {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStateCurrent,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	lacppdu = &layers.LACP{
		Version: layers.LACPVersion2,
		Actor: layers.LACPInfoTlv{TlvType: layers.LACPTLVActorInfo,
			Length: layers.LACPActorTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: [6]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
					SystemPriority: 1},
				Key:     100,
				PortPri: 0x80,
				Port:    10,
				State:   LacpStateActivityBit | LacpStateAggregationBit | LacpStateTimeoutBit},
		},
		Partner: layers.LACPInfoTlv{TlvType: layers.LACPTLVPartnerInfo,
			Length: layers.LACPPartnerTlvLength,
			Info: layers.LACPPortInfo{
				System: layers.LACPSystem{SystemId: p.ActorOper.System.Actor_System,
					SystemPriority: p.ActorOper.System.Actor_System_priority},
				Key:     p.Key,
				PortPri: p.portPriority,
				Port:    p.PortNum,
				State:   p.ActorOper.State},
		},
	}

	rx = LacpRxLacpPdu{
		pdu:          lacppdu,
		responseChan: portchan,
		src:          "TEST"}
	p.RxMachineFsm.RxmPktRxEvent <- rx

	// wait for response
	msg = <-portchan
	if msg != RxMachineModuleStr {
		t.Error("Expected response from", RxMachineModuleStr)
	}

	// port was enabled and lacp is disabled
	if p.RxMachineFsm.Machine.Curr.PreviousState() != LacpRxmStateCurrent &&
		p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateCurrent {
		t.Error("ERROR RX Machine State incorrect expected (prev/curr)",
			LacpRxmStateExpired,
			LacpRxmStateCurrent,
			"actual",
			p.RxMachineFsm.Machine.Curr.PreviousState(),
			p.RxMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteLaAggPort(pconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	OnlyForTestTeardown()
	LacpSysGlobalInfoDestroy(sysId)

}

func TestLaAggPortRxMachineInvalidStateTransitions(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	// must be called to initialize the global
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:     1,
		Prio:   0x80,
		IntfId: "SIMeth1.1",
		Key:    100,
	}

	utils.PortConfigMap[int32(pconf.Id)] = utils.PortConfig{Name: pconf.IntfId,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	// not calling Create because we don't want to launch all State machines
	p := NewLaAggPort(pconf)

	// lets start the Rx Machine only
	p.LacpRxMachineMain()

	p.BEGIN(true)

	// turn timer off so that we do not accidentally transition States
	p.RxMachineFsm.CurrentWhileTimerStop()
	/*
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
	*/

	// BEGIN -> INITIALIZE automatically falls through to PORT_DISABLED so no
	// need to tests

	// PORT_DISABLED
	portDisableInvalidStates := [4]fsm.Event{LacpRxmEventUnconditionalFallthrough,
		LacpRxmEventCurrentWhileTimerExpired,
		LacpRxmEventLacpEnabled,
		LacpRxmEventLacpPktRx}

	str, ok := InvalidStateCheck(p, portDisableInvalidStates[:], LacpRxmStateInitialize, LacpRxmStatePortDisabled)
	if !ok {
		t.Error(str)
	}

	// EXPIRED - note disabling current while timer so State does not change
	expiredInvalidStates := [5]fsm.Event{LacpRxmEventUnconditionalFallthrough,
		LacpRxmEventPortMoved,
		LacpRxmEventPortEnabledAndLacpEnabled,
		LacpRxmEventPortEnabledAndLacpDisabled,
		LacpRxmEventLacpEnabled}

	str, ok = InvalidStateCheck(p, expiredInvalidStates[:], LacpRxmStatePortDisabled, LacpRxmStateExpired)
	if !ok {
		t.Error(str)
	}

	// LACP_DISABLED
	lacpDisabledInvalidStates := [6]fsm.Event{LacpRxmEventUnconditionalFallthrough,
		LacpRxmEventPortMoved,
		LacpRxmEventPortEnabledAndLacpEnabled,
		LacpRxmEventPortEnabledAndLacpDisabled,
		LacpRxmEventCurrentWhileTimerExpired,
		LacpRxmEventLacpPktRx}

	str, ok = InvalidStateCheck(p, lacpDisabledInvalidStates[:], LacpRxmStatePortDisabled, LacpRxmStateLacpDisabled)
	if !ok {
		t.Error(str)
	}

	// DEFAULTED
	defaultedInvalidStates := [6]fsm.Event{LacpRxmEventUnconditionalFallthrough,
		LacpRxmEventPortMoved,
		LacpRxmEventPortEnabledAndLacpEnabled,
		LacpRxmEventPortEnabledAndLacpDisabled,
		LacpRxmEventCurrentWhileTimerExpired,
		LacpRxmEventLacpEnabled}

	str, ok = InvalidStateCheck(p, defaultedInvalidStates[:], LacpRxmStateExpired, LacpRxmStateDefaulted)
	if !ok {
		t.Error(str)
	}

	// DEFAULTED
	currentInvalidStates := [5]fsm.Event{LacpRxmEventUnconditionalFallthrough,
		LacpRxmEventPortMoved,
		LacpRxmEventPortEnabledAndLacpEnabled,
		LacpRxmEventPortEnabledAndLacpDisabled,
		LacpRxmEventLacpEnabled}

	str, ok = InvalidStateCheck(p, currentInvalidStates[:], LacpRxmStateExpired, LacpRxmStateCurrent)
	if !ok {
		t.Error(str)
	}

	p.LaAggPortDelete()
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	OnlyForTestTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}
