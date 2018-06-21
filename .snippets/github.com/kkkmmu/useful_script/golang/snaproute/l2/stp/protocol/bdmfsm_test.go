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

// bdm_test.go
package stp

import (
	"fmt"
	"testing"
)

func UsedForTestOnlyBdmTestSetup(stpconfig *StpPortConfig, t *testing.T) (p *StpPort) {

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

	stpconfig.BrgIfIndex = DEFAULT_STP_BRIDGE_VLAN

	// create a port
	p = NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the BDM State Machine
	p.BdmMachineMain()
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// only instanciated object not starting go routine
	PrxmMachineFSMBuild(p)
	PtxmMachineFSMBuild(p)
	PimMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	TcMachineFSMBuild(p)
	PstMachineFSMBuild(p)
	PpmmMachineFSMBuild(p)
	PrtMachineFSMBuild(p)

	if p.BdmMachineFsm.Machine.Curr.PreviousState() != BdmStateNone {
		t.Error("Failed to Initial Bridge Detection machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	if p.AdminEdge {
		if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
			t.Error("Failed to Initial Port Information machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.CurrentState())
			t.FailNow()
		}
	} else {
		if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
			t.Error("Failed to Initial Port Information machine state not set correctly", p.PpmmMachineFsm.Machine.Curr.CurrentState())
			t.FailNow()
		}
	}
	return p
}

func UsedForTestOnlyBdmTestTeardown(p *StpPort, t *testing.T) {

	if len(p.b.PrsMachineFsm.PrsEvents) > 0 {
		t.Error("Failed to check event sent")
	}

	if len(p.PpmmMachineFsm.PpmmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PimMachineFsm.PimEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PtxmMachineFsm.PtxmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.BdmMachineFsm.BdmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PtmMachineFsm.PtmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.TcMachineFsm.TcEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PstMachineFsm.PstEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PpmmMachineFsm.PpmmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PrxmMachineFsm.PrxmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	p.PrtMachineFsm = nil
	p.PimMachineFsm = nil
	p.PtxmMachineFsm = nil
	p.PrxmMachineFsm = nil
	p.PtmMachineFsm = nil
	p.TcMachineFsm = nil
	p.PstMachineFsm = nil
	p.PpmmMachineFsm = nil
	p.b.PrsMachineFsm = nil
	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestBdmMachineBEGINAdminEdge(t *testing.T) {

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_1(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Forwarding = false
	p.Agreed = false
	p.Proposing = false
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo {
		t.Error("Did not receive expected event")
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_2(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = false
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_3(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = false
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_4(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = true
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_5(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = true
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_6(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = false
	p.Disputed = true
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotPortEnabledAndNotAdminEdge_7(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = false
	p.PortEnabled = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = false
	p.Disputed = true
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndNotAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// should have received an event from bdm
	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event received %d", e.e))
	}
	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}

	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_1(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Forwarding = false
	p.Agreed = false
	p.Proposing = false
	p.Selected = true
	p.UpdtInfo = false

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}

	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_2(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = false
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_3(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = false
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_4(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = true
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_5(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = true
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_6(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = false
	p.Disputed = true
	p.Learn = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is not set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineNotOperEdge_7(t *testing.T) {
	testChan := make(chan string)
	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     true,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != true {
		t.Error("Failed operedge not set")
	}

	// this test should transition this state machine and transition
	// Prt state machine
	// The following settings should force an event
	p.AdminEdge = true
	p.PortEnabled = true
	p.OperEdge = false

	// Prt state variables
	p.Sync = true
	p.Synced = true
	p.ReRoot = false
	p.Disputed = true
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false
	p.Agreed = true

	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotOperEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("State not as expected")
	}
	UsedForTestOnlyBdmTestTeardown(p, t)

}

func TestBdmMachineBEGINAndNotAdminEdge(t *testing.T) {
	//testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndAdminEdge_1(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = false
	p.AdminEdge = true

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = false
	p.Sync = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndAdminEdge_2(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = false
	p.AdminEdge = true

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = true
	p.Sync = false
	p.Learn = false
	p.RrWhileTimer.count = 0
	p.ReRoot = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndAdminEdge_3(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = false
	p.AdminEdge = true

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = true
	p.Sync = false
	p.ReRoot = false
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndAdminEdge_4(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = false
	p.AdminEdge = true

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Learn = true
	p.Synced = true
	p.RrWhileTimer.count = 0
	p.Sync = false
	p.ReRoot = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineNotPortEnabledAndAdminEdge_5(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = false
	p.AdminEdge = true

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Learn = true
	p.Synced = true
	p.Sync = false
	p.ReRoot = false
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventNotPortEnabledAndAdminEdge,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRstpAndProposing_1(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = true
	p.AdminEdge = true
	p.AutoEdgePort = true
	p.SendRSTP = true
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = 0

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = false
	p.Sync = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRstpAndProposing_2(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = true
	p.AdminEdge = true
	p.AutoEdgePort = true
	p.SendRSTP = true
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = 0

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = true
	p.Sync = false
	p.Learn = false
	p.RrWhileTimer.count = 0
	p.ReRoot = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRstpAndProposing_3(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = true
	p.AdminEdge = true
	p.AutoEdgePort = true
	p.SendRSTP = true
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = 0

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Synced = true
	p.Sync = false
	p.ReRoot = false
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRstpAndProposing_4(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = true
	p.AdminEdge = true
	p.AutoEdgePort = true
	p.SendRSTP = true
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = 0

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Learn = true
	p.Synced = true
	p.RrWhileTimer.count = 0
	p.Sync = false
	p.ReRoot = true
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}

func TestBdmMachineEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRstpAndProposing_5(t *testing.T) {
	testChan := make(chan string)
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
	}

	p := UsedForTestOnlyBdmTestSetup(stpconfig, t)

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
	}

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateNotEdge {
		t.Error("Failed state is not as expected")
	}

	// setup test params
	p.PortEnabled = true
	p.AdminEdge = true
	p.AutoEdgePort = true
	p.SendRSTP = true
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = 0

	// specific event will be sent to port role transition state machine if the following is set
	p.Learning = true
	p.Learn = true
	p.Synced = true
	p.Sync = false
	p.ReRoot = false
	p.Selected = true
	p.UpdtInfo = false
	p.Role = PortRoleDesignatedPort
	// Do i need to set up the disabled port state variables?
	p.PrtMachineFsm.Machine.Curr.SetState(PrtStateDesignatedPort)

	p.BdmMachineFsm.BdmEvents <- MachineEvent{
		e:            BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if p.BdmMachineFsm.Machine.Curr.CurrentState() != BdmStateEdge {
		t.Error("State not as expected")
	}

	e := <-p.PrtMachineFsm.PrtEvents
	if e.e != PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo {
		t.Error(fmt.Sprintf("Did not receive expected event %d", e.e))
	}
	UsedForTestOnlyBdmTestTeardown(p, t)
}
