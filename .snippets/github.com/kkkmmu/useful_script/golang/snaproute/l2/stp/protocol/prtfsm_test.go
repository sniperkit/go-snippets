// prtfsm_test.go
package stp

import (
	"fmt"
	"testing"
	"utils/fsm"
)

// TODO This file is not doing much checking against the events that should
// be generated to other state machines according to Figure 17-12

func UsedForTestOnlyPrtInitPortConfigTest() {

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

func UsedForTestPrtBridgeAndPortSetup() *StpPort {
	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	StpBrgConfigSave(bridgeconfig)
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
	// add port to bridge list
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	StpPortConfigSave(stpconfig, false)

	// lets only start the Port Information State Machine
	p.PrtMachineMain()
	// going just send event and not start main as we just did above
	p.BEGIN(true)

	// only instanciated object not starting go routine
	PimMachineFSMBuild(p)
	PrxmMachineFSMBuild(p)
	PtxmMachineFSMBuild(p)
	BdmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	TcMachineFSMBuild(p)
	PstMachineFSMBuild(p)
	PpmmMachineFSMBuild(p)

	return p
}

func UsedForTestPrtFsmStartDisablePort(p *StpPort) {

	responseChan := make(chan string)

	// transition from designated to disabled
	// SelectedRole set by prs state machine
	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDisabledPort
	// selected set by prs state machine
	p.UpdtInfo = false
	p.Selected = true
	// normally set by pst state machine
	p.Learning = true
	p.Forwarding = true

	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
}

func UsedForTestPrtFsmStartDisabledPort(p *StpPort) {

	responseChan := make(chan string)

	// transition from designated to disabled
	// SelectedRole set by prs state machine
	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDisabledPort
	// selected set by prs state machine
	p.UpdtInfo = false
	p.Selected = true
	// normally set by pst state machine
	p.Learning = false
	p.Forwarding = false

	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
}

func TestPrtFsmDisablePortInvalidEvents(t *testing.T) {

	testChan := make(chan string)
	UsedForTestOnlyPrtInitPortConfigTest()
	p := UsedForTestPrtBridgeAndPortSetup()

	UsedForTestPrtFsmStartDisablePort(p)

	invalidStateMap := [44]fsm.Event{

		PrtEventUnconditionallFallThrough,
		PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndSelectedAndNotUpdtInfo,
		PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
		PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
		PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
		PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo,
		PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
		PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
		PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo,
		PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
		PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo,
	}

	// test the invalid states
	for _, e := range invalidStateMap {

		p.PrtMachineFsm.PrtEvents <- MachineEvent{e: e,
			src:          "TestPrtFsmDisablePortInvalidEvents",
			responseChan: testChan}

		<-testChan
		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisablePort {
			t.Errorf(fmt.Sprintf("Failed e[%d] transitioned PRT to a new state state %d", e, p.PrtMachineFsm.Machine.Curr.CurrentState()))
		}
		if p.Role != PortRoleDisabledPort {
			t.Error("ERROR: Role did not change", e)
		}
		if p.Learn || p.Forward {
			t.Error("ERROR: Learn or Forward not cleared", p.Learn, p.Forward, e)
		}
	}

	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p)
	DelStpBridge(b, true)

}

func TestPrtFsmDisabledPortInvalidStates(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPrtInitPortConfigTest()
	p := UsedForTestPrtBridgeAndPortSetup()

	UsedForTestPrtFsmStartDisabledPort(p)

	invalidStateMap := [45]fsm.Event{

		PrtEventUnconditionallFallThrough,
		PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndSelectedAndNotUpdtInfo,
		PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
		PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
		PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
		PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo,
		PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
		PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
		PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
		PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo,
		PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
		PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
		PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
		PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo,
		PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo,
	}

	// test the invalid states
	for _, e := range invalidStateMap {

		p.PrtMachineFsm.PrtEvents <- MachineEvent{e: e,
			src:          "TestPrtFsmDisablePortInvalidEvents",
			responseChan: testChan}

		<-testChan
		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
			t.Errorf(fmt.Sprintf("Failed e[%d] transitioned PRT to a new state state %d", e, p.PrtMachineFsm.Machine.Curr.CurrentState()))
		}
		if p.Role != PortRoleDisabledPort {
			t.Error("ERROR: Role did not change")
		}
		if p.Learn || p.Forward {
			t.Error("ERROR: Learn or Forward not cleared", p.Learn, p.Forward)
		}
	}

	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p)
	DelStpBridge(b, true)
}

func TestPrtFsmDisablePortStates(t *testing.T) {

	responseChan := make(chan string)
	UsedForTestOnlyPrtInitPortConfigTest()
	p := UsedForTestPrtBridgeAndPortSetup()

	// two states that pst can be in when in order for events
	// to be forwarded from prt
	for _, s := range [2]fsm.State{PstStateLearning, PstStateForwarding} {

		// lets assume that pst machine is in forwarding state
		// lets force it
		p.PstMachineFsm.Machine.Curr.SetState(s)
		p.Learn = true
		p.Learning = true
		p.Forward = true
		p.Forwarding = true

		UsedForTestPrtFsmStartDisablePort(p)

		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisablePort {
			t.Error("ERROR: Not in Disable Port state")
		}
		if p.Role != PortRoleDisabledPort {
			t.Error("ERROR: Role did not change")
		}
		if p.Learn || p.Forward {
			t.Error("ERROR: Learn or Forward not cleared", p.Learn, p.Forward)
		}

		// check that the proper event was received by this machine
		event := <-p.PstMachineFsm.PstEvents
		if s == PstStateLearning {
			if event.e != PstEventNotLearn {
				t.Error("ERROR: Did not get no forward event")
			}

		} else {
			if event.e != PstEventNotForward {
				t.Error("ERROR: Did not get no forward event")
			}
		}
	}
	// lets transition to next state disabled
	p.Learning = false
	p.Forwarding = false

	// lets ensure that we transition
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}
	<-responseChan

	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
		t.Error("ERROR: Error did not transition to correct state")
	}

	if p.FdWhileTimer.count != int32(p.b.BridgeTimes.MaxAge) {
		t.Error("ERROR: FdWhile counter not set properly")
	}
	if !p.Synced {
		t.Error("ERROR: Synced should be set")
	}

	if p.RrWhileTimer.count != 0 {
		t.Error("ERROR: rrwhile timer should be cleared")
	}

	if p.Sync {
		t.Error("ERROR: Sync was not cleared")
	}

	if p.ReRoot {
		t.Error("ERROR: reroot to cleared")
	}

	// lets validate that all of these states do not cause any changes in the state
	// transitions
	validDisabledEventsStateMap := [6]fsm.Event{
		PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
		PrtEventSyncAndSelectedAndNotUpdtInfo,
		PrtEventReRootAndSelectedAndNotUpdtInfo,
		PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
		PrtEventBegin,
		PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
	}

	for _, e := range validDisabledEventsStateMap {

		p.PrtMachineFsm.PrtEvents <- MachineEvent{e: e,
			src:          "TestPrtFsmDisablePortInvalidEvents",
			responseChan: responseChan}

		<-responseChan
		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
			t.Errorf(fmt.Sprintf("Failed e[%d] transitioned PRT to a new state state %d", e, p.PrtMachineFsm.Machine.Curr.CurrentState()))
		}

		if p.FdWhileTimer.count != int32(p.b.BridgeTimes.MaxAge) {
			t.Error("ERROR: FdWhile counter not set properly", e)
		}
		if !p.Synced {
			t.Error("ERROR: Synced should be set", e)
		}

		if p.RrWhileTimer.count != 0 {
			t.Error("ERROR: rrwhile timer should be cleared", e)
		}

		if p.Sync {
			t.Error("ERROR: Sync was not cleared", e)
		}

		if p.ReRoot {
			t.Error("ERROR: reroot to cleared", e)
		}
	}

	// reset state so that we can delete the resources
	p.PstMachineFsm.Machine.Curr.SetState(PstStateNone)

	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p)
	DelStpBridge(b, true)

}

type statevalidatefn func(p *StpPort) (rstr string, rv bool)

func ValidateRootProposed(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.Proposed {
		rstr = "ERROR: Proposed is set, should have been cleared"
		rv = false
	}

	// setSyncTree call verify
	if rv {
		var port *StpPort
		for _, IfIndex := range p.b.StpPorts {
			if rv {
				if StpFindPortByIfIndex(IfIndex, p.BrgIfIndex, &port) {
					if !port.Sync {
						rstr = fmt.Sprintln("ERROR: Sync is not set for port", port.IfIndex)
						rv = false
					}
				} else {
					rstr = "ERROR: could not find port"
					rv = false
				}
			}
		}
	}

	return rstr, rv
}

func ValidateRootPort(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.Role != PortRoleRootPort {
		rstr = "ERROR Role not set to Root Port"
		rv = false
	}
	if rv && p.RrWhileTimer.count != int32(p.b.RootTimes.ForwardingDelay) {
		rstr = "ERROR rrwhile != FwdDelay"
		rv = false
	}
	return rstr, rv
}

func ValidateRootAgreed(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.Proposed {
		rstr = "ERROR proposed is set, should have been cleared"
		rv = false
	}
	if rv && p.Sync {
		rstr = "ERROR Sync is set, should have been cleared"
		rv = false
	}
	if rv && !p.Agree {
		rstr = "ERROR Agree is not set, should have been set"
		rv = false
	}
	if rv && !p.NewInfo {
		rstr = "ERROR NewInfo is not set, should have been set"
		rv = false
	}
	return rstr, rv
}

func ValidateReRoot(p *StpPort) (rstr string, rv bool) {
	// setReRootTree call verify
	rv = true
	var port *StpPort
	for _, IfIndex := range p.b.StpPorts {
		if rv {
			if StpFindPortByIfIndex(IfIndex, p.BrgIfIndex, &port) {
				if !port.ReRoot {
					rstr = fmt.Sprintln("ERROR: ReRoot is not set for port", port.IfIndex)
					rv = false
				}
			}
		}
	}
	return rstr, rv
}

func ValidateRootForward(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.FdWhileTimer.count != 0 {
		rstr = "ERROR: FdwWhile != 0, should be set to zero"
		rv = false
	}
	if rv && !p.Forward {
		rstr = "ERROR: Forward not set, should be set"
		rv = false
	}
	return rstr, rv
}

func ValidateRootLearn(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.FdWhileTimer.count != int32(p.b.RootTimes.ForwardingDelay) {
		rstr = "ERROR: FdwWhile != ForwardDelay, should be set to ForwardDelay"
		rv = false
	}
	if rv && !p.Learn {
		rstr = "ERROR: Learn not set, should be set"
		rv = false
	}
	return rstr, rv
}

func ValidateReRooted(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.ReRoot {
		rstr = "ERROR: ReRoot set, should be cleared"
		rv = false
	}
	return rstr, rv
}

func TestRootPortValidateStateTransitions(t *testing.T) {
	var tests = []struct {
		e  fsm.Event
		s  fsm.State
		fn statevalidatefn
	}{
		{PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, PrtStateRootPort, ValidateRootPort},
		{PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, PrtStateRootProposed, ValidateRootProposed},
		{PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, PrtStateRootAgreed, ValidateRootAgreed},
		{PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, PrtStateRootAgreed, ValidateRootAgreed},
		{PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo, PrtStateReRoot, ValidateReRoot},
		{PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo, PrtStateRootPort, ValidateRootPort},
		{PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo, PrtStateReRooted, ValidateReRooted},
		{PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, PrtStateRootLearn, ValidateRootLearn},
		{PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, PrtStateRootLearn, ValidateRootLearn},
		{PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, PrtStateRootForward, ValidateRootForward},
		{PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, PrtStateRootForward, ValidateRootForward},
	}

	UsedForTestOnlyPrtInitPortConfigTest()
	p := UsedForTestPrtBridgeAndPortSetup()

	// setup
	p.Synced = true
	p.Selected = true
	p.UpdtInfo = true
	p.SelectedRole = PortRoleRootPort
	p.Role = PortRoleDesignatedPort
	responseChan := make(chan string)
	// start the prt machine at root port
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TestRootPortValidateStateTransitions",
		responseChan: responseChan,
	}
	<-responseChan

	// validate all state transitions
	for _, test := range tests {
		p.PrtMachineFsm.PrtEvents <- MachineEvent{
			e:            test.e,
			src:          "TestRootPortValidateStateTransitions",
			responseChan: responseChan,
		}
		<-responseChan
		if p.PrtMachineFsm.Machine.Curr.PreviousState() != test.s {
			t.Error("ERROR: Previous state not as expected", p.PrtMachineFsm.Machine.Curr.PreviousState())
		}
		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
			t.Error("ERROR: Previous state not as expected", p.PrtMachineFsm.Machine.Curr.CurrentState())
		}
		rstr, rv := test.fn(p)
		if !rv {
			t.Error(rstr)
		}

		rstr, rv = ValidateRootPort(p)
		if !rv {
			t.Error(rstr)
		}
	}

	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

func ValidateBlockPort(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.Role != p.SelectedRole {
		rstr = "ERROR: Role != Selected Role, should have been set"
		rv = false
	}
	if rv && p.Learn {
		rstr = "ERROR: Learn not cleared"
		rv = false
	}
	if rv && p.Forward {
		rstr = "ERROR: Forward not cleared"
		rv = false
	}
	return rstr, rv
}

func ValidateBackupPort(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.RbWhileTimer.count != 2*int32(p.b.RootTimes.HelloTime) {
		rstr = "ERROR: rbwhile timer not set correctly"
		rv = false
	}
	return rstr, rv
}

func ValidateAlternatePort(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.FdWhileTimer.count != int32(p.b.RootTimes.ForwardingDelay) {
		rstr = "ERROR: fdwhile timer not set correctly"
		rv = false
	}
	if rv && !p.Synced {
		rstr = "ERROR: Error synced was not set"
		rv = false
	}
	if rv && p.RrWhileTimer.count != 0 {
		rstr = "ERROR: rrwhile timer was not cleared"
		rv = false
	}
	if rv && p.Sync {
		rstr = "ERROR: Sync not cleared"
		rv = false
	}
	if rv && p.ReRoot {
		rstr = "ERROR: ReRoot not cleared"
		rv = false
	}
	return rstr, rv
}

func ValidateAlternateAgreed(p *StpPort) (rstr string, rv bool) {
	rv = true
	if p.Proposed {
		rstr = "ERROR: proposed not cleared"
		rv = false
	}
	if rv && !p.Agree {
		rstr = "ERROR: Agree not set"
		rv = false
	}
	if rv && !p.NewInfo {
		rstr = "ERROR: NewInfo not set"
		rv = false
	}
	return rstr, rv
}

func ValidateAlternateProposed(p *StpPort) (rstr string, rv bool) {
	rv = true
	// SetSyncTree
	/* This can't be verified as the Alternate Port state clears
	   the sync bit, perhaps we can check if we add two ports in the test
	var port *StpPort
	for _, IfIndex := range p.b.StpPorts {
		if rv {
			if StpFindPortByIfIndex(IfIndex, p.BrgIfIndex, &port) {
				if !port.Sync {
					rstr = fmt.Sprintln("ERROR: Sync is not set for port", port.IfIndex)
					rv = false
				}
			}
		}
	}
	*/
	if rv && p.Proposed {
		rstr = "ERROR: Proposed not cleared"
		rv = false
	}
	return rstr, rv
}

func TestAlternateBackupValidateStateTransitions(t *testing.T) {
	var tests = []struct {
		e  fsm.Event
		s  fsm.State
		fn statevalidatefn
	}{
		{PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, PrtStateBlockPort, ValidateBlockPort},
		{PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, PrtStateAlternateProposed, ValidateAlternateProposed},
		{PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, PrtStateAlternateAgreed, ValidateAlternateAgreed},
		{PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, PrtStateAlternateAgreed, ValidateAlternateAgreed},
		{PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo, PrtStateBackupPort, ValidateBackupPort},
		{PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo, PrtStateAlternatePort, ValidateAlternatePort},
		{PrtEventSyncAndSelectedAndNotUpdtInfo, PrtStateAlternatePort, ValidateAlternatePort},
		{PrtEventReRootAndSelectedAndNotUpdtInfo, PrtStateAlternatePort, ValidateAlternatePort},
		{PrtEventNotSyncedAndSelectedAndNotUpdtInfo, PrtStateAlternatePort, ValidateAlternatePort},
	}

	UsedForTestOnlyPrtInitPortConfigTest()
	p := UsedForTestPrtBridgeAndPortSetup()

	// setup
	p.Synced = true
	p.Learn = true
	p.Forward = true
	p.Proposed = true
	p.Agree = false

	// set to that we have a post state fall through when in blocked port state
	p.Learning = false
	p.Forwarding = false
	// required for all states
	p.Selected = true
	p.UpdtInfo = false

	p.SelectedRole = PortRoleBackupPort
	p.Role = PortRoleDesignatedPort
	responseChan := make(chan string)
	// start the prt machine at root port
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualBackupPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TestRootPortValidateStateTransitions",
		responseChan: responseChan,
	}
	<-responseChan

	// lets try the other port
	p.SelectedRole = PortRoleAlternatePort
	// validate all state transitions
	for _, test := range tests {

		// setup
		if test.e == PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo {
			p.Learn = true
			p.Forward = true
		} else {
			p.Learn = false
			p.Forward = false
		}
		p.Synced = false

		// set to that we have a post state fall through when in blocked port state
		p.Learning = false
		p.Forwarding = false

		if test.e == PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo {
			p.Proposed = true
			p.Agree = false
			p.Sync = false
		} else if test.e == PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo {
			p.Proposed = true
			p.Agree = true
			p.Sync = true
		} else {
			p.Proposed = false
			p.Agree = true
		}

		// these are always set to different values every time we enter alternate
		p.RrWhileTimer.count = 10
		p.ReRoot = true
		p.FdWhileTimer.count = 0

		p.PrtMachineFsm.PrtEvents <- MachineEvent{
			e:            test.e,
			src:          "TestRootPortValidateStateTransitions",
			responseChan: responseChan,
		}
		<-responseChan
		if test.e != PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo {
			if p.PrtMachineFsm.Machine.Curr.PreviousState() != test.s {
				t.Error("ERROR: Previous state not as expected", p.PrtMachineFsm.Machine.Curr.PreviousState(), test.e)
			}
		} else {
			// because we are only tsting with one port this state will
			// auto matically cause a transition to Alternate Agreed state
			// beause the Alternate Port state sets this port to synced thus
			// the post event of allSynced && !agree will be true
			if p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateAlternateAgreed {
				t.Error("ERROR: Previous state is not as expected", p.PrtMachineFsm.Machine.Curr.PreviousState(), test.e)
			}
		}
		if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateAlternatePort {
			t.Error("ERROR: Previous state not as expected", p.PrtMachineFsm.Machine.Curr.CurrentState(), test.e)
		}
		rstr, rv := test.fn(p)
		if !rv {
			t.Error(rstr)
		}
	}

	b := p.b
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p)
	DelStpBridge(b, true)

}
