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

// ppm_test.go
// This is a test file to test the port protocol migration state machine
package stp

import (
	"fmt"
	//"net"
	//"strconv"
	//"strings"
	"testing"
	//"time"
	"utils/fsm"
)

func UsedForTestOnlyPpmmInitPortConfigTest() {

	if PortConfigMap == nil {
		PortConfigMap = make(map[int32]portConfig)
	}
	// In order to test a packet we must listen on loopback interface
	// and send on interface we expect to receive on.  In order
	// to do this a couple of things must occur the PortConfig
	// must be updated with "dummy" ifindex pointing to 'lo'
	TEST_RX_PORT_CONFIG_IFINDEX = 0x0ADDBEEF
	PortConfigMap[TEST_RX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo"}
	PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo"}
	/*
		intfs, err := net.Interfaces()
		if err == nil {
			for _, intf := range intfs {
				if strings.Contains(intf.Name, "eth") {
					ifindex, _ := strconv.Atoi(strings.Split(intf.Name, "eth")[1])
					if ifindex == 0 {
						TEST_TX_PORT_CONFIG_IFINDEX = int32(ifindex)
					}
					PortConfigMap[int32(ifindex)] = portConfig{Name: intf.Name}
				}
			}
		}
	*/
	UsedForTestOnlySetupAsicDPlugin()
}

func UsedForTestOnlyCheckPpmStateSensing(p *StpPort, t *testing.T, trace string) {

	if p.RcvdRSTP != false {
		t.Error(fmt.Sprintf("Failed RcvdRSTP was set%s\n", trace))
		t.FailNow()
	}
	if p.RcvdSTP != false {
		t.Error(fmt.Sprintf("Failed RcvdSTP was set%s\n", trace))
		t.FailNow()
	}
}

func UsedForTestOnlyCheckPpmCheckingRSTP(p *StpPort, t *testing.T, trace string) {
	if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateCheckingRSTP {
		t.Error(fmt.Sprintf("Failed to transition to Checking RSTP State current state %d trace %s\n", p.PpmmMachineFsm.Machine.Curr.CurrentState(), trace))
		t.FailNow()
	}

	if p.Mcheck == true {
		t.Error(fmt.Sprintf("Failed mcheck not set to false %s\n", trace))
		t.FailNow()
	}

	// default rstpVersion is RSTP
	if p.SendRSTP != true {
		t.Error(fmt.Sprintf("Failed sendRstp is not set%s\n", trace))
		t.FailNow()
	}

	if p.MdelayWhiletimer.count != MigrateTimeDefault {
		t.Error(fmt.Sprintf("Failed MdelayWhiletimer tick count not set to MigrateTimeDefault %s\n", trace))
		t.FailNow()
	}
}

func UsedForTestOnlyCheckPpmStateSelectingSTP(p *StpPort, t *testing.T, trace string) {
	if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateSelectingSTP {
		t.Error(fmt.Sprintf("Failed to transition to Checking RSTP State current state %d trace %s\n", p.PpmmMachineFsm.Machine.Curr.CurrentState(), trace))
		t.FailNow()
	}

	// default rstpVersion is RSTP
	if p.SendRSTP != false {
		t.Error(fmt.Sprintf("Failed sendRstp is set%s\n", trace))
		t.FailNow()
	}

	if p.MdelayWhiletimer.count != MigrateTimeDefault {
		t.Error(fmt.Sprintf("Failed MdelayWhiletimer tick count not set to MigrateTimeDefault %s\n", trace))
		t.FailNow()
	}
}

func UsedForTestOnlyStartPpmInSelectingSTPState(p *StpPort, t *testing.T, trace string) {
	testChan := make(chan string)
	// Reset State machine BEGIN
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 1
	p.PortEnabled = true
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventBegin,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, fmt.Sprintf("1 %s\n", trace))

	p.MdelayWhiletimer.count = 0
	// lets transition the state machine to SELECTING_STP
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	p.SendRSTP = true
	p.RcvdSTP = true

	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventSendRSTPAndRcvdSTP,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmStateSelectingSTP(p, t, fmt.Sprintf("1 %s\n", trace))

}

func UsedForTestOnlyStartPpmInSensingState(p *StpPort, t *testing.T, viaselectingstp int, trace string) {
	testChan := make(chan string)

	// lets transition the state machine to CHECKING RSTP
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventBegin,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	p.MdelayWhiletimer.count = 0
	// lets transition the state machine to SELECTING_STP
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmStateSensing(p, t, fmt.Sprintf("1 %s\n", trace))

	if viaselectingstp == 1 {
		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventSendRSTPAndRcvdSTP,
			src:          "TestPpmCheckingRSTPStateTransitions",
			responseChan: testChan}
		<-testChan

		UsedForTestOnlyCheckPpmStateSelectingSTP(p, t, fmt.Sprintf("1 %s\n", trace))

		p.MdelayWhiletimer.count = 0
		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
			src:          "TestPpmCheckingRSTPStateTransitions",
			responseChan: testChan}

		<-testChan
		UsedForTestOnlyCheckPpmStateSensing(p, t, fmt.Sprintf("1 %s\n", trace))
	}
}

func TestPpmmBEGIN(t *testing.T) {
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)

	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PpmmMachineFsm.Machine.Curr.PreviousState() != PrxmStateNone {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

func TestPpmCheckingRSTPInvalidStateTransitions(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PpmmMachineFsm.Machine.Curr.PreviousState() != PrxmStateNone {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	PtxmMachineFSMBuild(p)

	invalidStateMap := [4]fsm.Event{
		PpmmEventNotPortEnabled,
		PpmmEventMcheck,
		PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
		PpmmEventSendRSTPAndRcvdSTP}

	// test the invalid states
	for _, e := range invalidStateMap {
		// one way to validate that a transition did not occur is to make sure
		// that mdelayWhiletimer is not set to MigrateTime
		p.MdelayWhiletimer.count = 1

		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: e,
			src:          "TestPpmCheckingRSTPStateTransitions",
			responseChan: testChan}

		<-testChan
		if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateCheckingRSTP {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a new state state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
		if p.MdelayWhiletimer.count != 1 {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a back to same state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
	}

	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPpmCheckingRSTPMdelayWhileNotEqualMigrateTimeAndNotPortEnable(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PpmmMachineFsm.Machine.Curr.PreviousState() != PrxmStateNone {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	PtxmMachineFSMBuild(p)

	// test valid states

	// mdelaywhile != migratetime and not port enabled
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 1
	p.PortEnabled = false
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "2")

	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

/*
func TestPpmCheckingRSTP(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:               TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
	}

	// create a port
	p := NewStpPort(stpconfig)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PpmmMachineFsm.Machine.Curr.PreviousState() != PrxmStateNone {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	PtxmMachineFSMBuild(p)

	// reset state back to checking rstp BEGIN
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 1
	p.PortEnabled = true
	// this should send event to Port Trasnmit Machine
	p.SendRSTP = false
	p.NewInfo = true
	p.TxCount = 0
	p.HelloWhenTimer.count = BridgeHelloTimeDefault
	p.Selected = true
	p.UpdtInfo = false
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)

	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventBegin,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	// TX Machine should have received this event
	// Test will assume that tx is in Idle state which means a couple of port params need to be updated

	event, _ := <-p.PtxmMachineFsm.PtxmEvents
	if event.e != PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Failed PTXM failed to received event %d", PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo))
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "3")

	// Mdelaywhile == 0
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 0
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmStateSensing(p, t, "1")

	p.PtxmMachineFsm = nil

	DelStpPort(p)

}
*/
func TestPpmSelectingSTPInvalidStateTransitions(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// lets transition the state machine to SELECTING_STP
	UsedForTestOnlyStartPpmInSelectingSTPState(p, t, "1")

	invalidStateMap := [4]fsm.Event{
		PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled,
		PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
		PpmmEventSendRSTPAndRcvdSTP}

	// test the invalid states
	for _, e := range invalidStateMap {
		// one way to validate that a transition did not occur is to make sure
		// that mdelayWhiletimer is not set to MigrateTime
		p.MdelayWhiletimer.count = 1

		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: e,
			src:          "TestPpmCheckingRSTPStateTransitions",
			responseChan: testChan}

		<-testChan
		if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateSelectingSTP {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a new state state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
		if p.MdelayWhiletimer.count != 1 {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a back to same state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
	}
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

func TestPpmSelectingSTPMdelayWhileNotEqualMigrateTime(t *testing.T) {
	testChan := make(chan string)
	// test valid states
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// lets transition the state machine to SELECTING_STP
	UsedForTestOnlyStartPpmInSelectingSTPState(p, t, "1")

	// TEST mdelaywhile != migratetime and not port enabled
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 0
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateSensing {
		t.Error(fmt.Sprintf("Failed event %d did not transition state expected %d", PpmmEventMdelayWhileEqualZero, PpmmStateSensing))
		t.FailNow()
	}

	if p.Mcheck != false {
		t.Error("Failed mcheck not set to FALSE")
		t.FailNow()
	}
	if p.SendRSTP != false {
		t.Error("Failed sendRSTP set to true")
		t.FailNow()
	}
	if p.MdelayWhiletimer.count != 0 {
		t.Error("Failed MdelayWhile was not reset")
		t.FailNow()
	}
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPpmSelectingSTPMdelayWhileEqualZero(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// lets transition the state machine to SELECTING_STP
	UsedForTestOnlyStartPpmInSelectingSTPState(p, t, "1")

	// TEST Mdelaywhile == 0
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 0
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMdelayWhileEqualZero,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateSensing {
		t.Error(fmt.Sprintf("Failed event %d did not transition state expected %d\n", PpmmEventMdelayWhileEqualZero, PpmmStateSensing))
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmStateSensing(p, t, "1")
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPpmSelectingSTPMcheck(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// lets transition the state machine to SELECTING_STP
	UsedForTestOnlyStartPpmInSelectingSTPState(p, t, "1")

	// mcheck == true
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 1
	p.Mcheck = true
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMcheck,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateCheckingRSTP {
		t.Error(fmt.Sprintf("Failed event %d did not transition state expected %d\n", PpmmEventMcheck, PpmmStateCheckingRSTP))
		t.FailNow()
	}

	if p.Mcheck != false {
		t.Error("Failed mcheck not set to FALSE")
		t.FailNow()
	}
	if p.SendRSTP != true {
		t.Error("Failed sendRSTP not set to true")
		t.FailNow()
	}

	if p.MdelayWhiletimer.count != MigrateTimeDefault {
		t.Error("Failed MdelayWhile was not reset")
		t.FailNow()
	}

	if p.RcvdRSTP != false {
		t.Error("Failed RcvdRSTP was set")
		t.FailNow()
	}
	if p.RcvdSTP != false {
		t.Error("Failed RcvdSTP was set")
		t.FailNow()
	}
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPpmSelectingSTPNotPortEnabled(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// lets transition the state machine to SELECTING_STP
	UsedForTestOnlyStartPpmInSelectingSTPState(p, t, "1")

	// Port Not Enabled
	// set pre-existing conditions for event to be raised
	p.MdelayWhiletimer.count = 1
	p.PortEnabled = false
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventNotPortEnabled,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	if p.PpmmMachineFsm.Machine.Curr.PreviousState() != PpmmStateSensing {
		t.Error(fmt.Sprintf("Failed event %d did not transition state expected %d\n", PpmmEventNotPortEnabled, PpmmStateSensing))
		t.FailNow()

	}

	// from transition from Sensing state
	if p.RcvdRSTP != false {
		t.Error("Failed RcvdRSTP was set")
		t.FailNow()
	}
	if p.RcvdSTP != false {
		t.Error("Failed RcvdSTP was set")
		t.FailNow()
	}

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPpmSensingInvalidStateTransitions(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	PtxmMachineFSMBuild(p)

	// lets transition the state machine to SENSING
	UsedForTestOnlyStartPpmInSensingState(p, t, 0, "1")

	invalidStateMap := [4]fsm.Event{
		PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled,
		PpmmEventMdelayWhileEqualZero}

	// test the invalid states
	for _, e := range invalidStateMap {
		// one way to validate that a transition did not occur is to make sure
		// that mdelayWhiletimer is not set to MigrateTime
		p.MdelayWhiletimer.count = 0

		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: e,
			src:          "TestPpmCheckingRSTPStateTransitions",
			responseChan: testChan}

		<-testChan
		if p.PpmmMachineFsm.Machine.Curr.CurrentState() != PpmmStateSensing {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a new state state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
		if p.MdelayWhiletimer.count != 0 {
			t.Error(fmt.Sprintf("Failed e[%d] transitioned PPM to a back to same state %d", e, p.PpmmMachineFsm.Machine.Curr.CurrentState()))
			t.FailNow()
		}
	}
	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

// two part test will based on previous state
func TestPpmSensingDesignatedPortSendRSTPAndRcvdSTP(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	PtxmMachineFSMBuild(p)

	// lets transition the state machine to SENSING
	UsedForTestOnlyStartPpmInSensingState(p, t, 0, "1")

	// Send RSTP and Received STP
	// previous state was checking RSTP
	// set pre-existing conditions for event to be raised
	p.RcvdSTP = true                // set by port receive state machine
	p.Selected = true               // assumed port role selection state machine set this
	p.Role = PortRoleDesignatedPort // port role selection state machine sets this
	p.UpdtInfo = false              // assume port role selection state machine already set this
	p.TxCount = 0
	p.NewInfo = true // owned by topo change machine, port role and port info
	p.HelloWhenTimer.count = 1
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventSendRSTPAndRcvdSTP,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmStateSelectingSTP(p, t, "1")

	// should receive a tx event cause send RSTP is being
	// set to false
	event, _ := <-p.PtxmMachineFsm.PtxmEvents

	if event.e != PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo {
		t.Error("Failed to send event to tx machine")
		t.FailNow()

	}

	// lets transition the state machine to SENSING via SELECTING
	// force the previous state to be that of
	UsedForTestOnlyStartPpmInSensingState(p, t, 1, "1")

	// previous state was checking RSTP
	// set pre-existing conditions for event to be raised
	p.RcvdRSTP = true               // set by port receive state machine
	p.Selected = true               // assumed port role selection state machine set this
	p.Role = PortRoleDesignatedPort // port role selection state machine sets this
	p.UpdtInfo = false              // assume port role selection state machine already set this
	p.TxCount = 0
	p.NewInfo = true // owned by topo change machine, port role and port info
	p.HelloWhenTimer.count = 1
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)

	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	// should receive a tx event cause send RSTP is being
	// set to false
	event, _ = <-p.PtxmMachineFsm.PtxmEvents

	if event.e != PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo {
		t.Error("Failed to send event to tx machine")
		t.FailNow()

	}

	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

// two part test will based on previous state
func TestPpmSensingRootPortSendRSTPAndRcvdSTP(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	PtxmMachineFSMBuild(p)

	// lets transition the state machine to SENSING
	UsedForTestOnlyStartPpmInSensingState(p, t, 0, "1")

	// Send RSTP and Received STP
	// previous state was checking RSTP
	// set pre-existing conditions for event to be raised
	p.RcvdSTP = true          // set by port receive state machine
	p.Selected = true         // assumed port role selection state machine set this
	p.Role = PortRoleRootPort // port role selection state machine sets this
	p.UpdtInfo = false        // assume port role selection state machine already set this
	p.TxCount = 0
	p.NewInfo = true // owned by topo change machine, port role and port info
	p.HelloWhenTimer.count = 1
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventSendRSTPAndRcvdSTP,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmStateSelectingSTP(p, t, "1")

	// should receive a tx event cause send RSTP is being
	// set to false
	event, _ := <-p.PtxmMachineFsm.PtxmEvents

	if event.e != PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo {
		t.Error("Failed to send event to tx machine")
		t.FailNow()

	}

	// lets transition the state machine to SENSING via SELECTING
	// force the previous state to be that of
	UsedForTestOnlyStartPpmInSensingState(p, t, 0, "1")

	// previous state was checking RSTP
	// set pre-existing conditions for event to be raised
	p.RcvdRSTP = true               // set by port receive state machine
	p.Selected = true               // assumed port role selection state machine set this
	p.Role = PortRoleDesignatedPort // port role selection state machine sets this
	p.UpdtInfo = false              // assume port role selection state machine already set this
	p.TxCount = 0
	p.NewInfo = true // owned by topo change machine, port role and port info
	p.HelloWhenTimer.count = 1
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)

	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	// should receive a tx event cause send RSTP is being
	// set to false
	event, _ = <-p.PtxmMachineFsm.PtxmEvents

	if event.e != PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo {
		t.Error("Failed to send event to tx machine")
		t.FailNow()

	}

	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

// two part test will based on previous state
func TestPpmSensingNotPortEnabled(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPpmmInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p := NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PpmmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	PtxmMachineFSMBuild(p)

	// lets transition the state machine to SENSING
	UsedForTestOnlyStartPpmInSensingState(p, t, 0, "1")

	// Send RSTP and Received STP
	// previous state was checking RSTP
	// set pre-existing conditions for event to be raised
	p.RcvdSTP = true          // set by port receive state machine
	p.Selected = true         // assumed port role selection state machine set this
	p.Role = PortRoleRootPort // port role selection state machine sets this
	p.UpdtInfo = false        // assume port role selection state machine already set this
	p.TxCount = 0
	p.NewInfo = true // owned by topo change machine, port role and port info
	p.HelloWhenTimer.count = 1
	p.PtxmMachineFsm.Machine.Curr.SetState(PtxmStateIdle)
	// send test event
	p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventNotPortEnabled,
		src:          "TestPpmCheckingRSTPStateTransitions",
		responseChan: testChan}

	<-testChan

	UsedForTestOnlyCheckPpmCheckingRSTP(p, t, "1")

	p.PtxmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}
