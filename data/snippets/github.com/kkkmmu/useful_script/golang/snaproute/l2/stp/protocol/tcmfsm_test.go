// timer_test.go
package stp

import (
	"net"
	"testing"
	"time"
	asicdmock "utils/asicdClient/mock"
	"utils/fsm"
)

type MockTcAsicdClientMgr struct {
	asicdmock.MockAsicdClientMgr
}

func (asicdClientMgr *MockTcAsicdClientMgr) FlushStgFdb(stgid, port int32) error {
	// lets pretend it took a second to flush
	time.Sleep(time.Second * 1)
	return nil
}

func UsedForTestOnlyTcmInitPortConfigTest() {

	if PortConfigMap == nil {
		PortConfigMap = make(map[int32]portConfig)
	}
	// In order to test a packet we must listen on loopback interface
	// and send on interface we expect to receive on.  In order
	// to do this a couple of things must occur the PortConfig
	// must be updated with "dummy" ifindex pointing to 'lo'
	TEST_RX_PORT_CONFIG_IFINDEX = 0x0ADDBEEF
	PortConfigMap[TEST_RX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo",
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo",
		HardwareAddr: net.HardwareAddr{0x00, 0x33, 0x22, 0x22, 0x11, 0x11},
	}
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
	SetAsicDPlugin(&MockTcAsicdClientMgr{})
}

func UsedForTestOnlyTcmTestSetup(t *testing.T) (p *StpPort) {
	UsedForTestOnlyTcmInitPortConfigTest()

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
		BrgIfIndex:        100,
	}

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		Vlan:         100,
	}

	//StpBridgeCreate
	b := NewStpBridge(bridgeconfig)
	PrsMachineFSMBuild(b)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventBegin, nil)
	b.PrsMachineFsm.Machine.ProcessEvent("TEST", PrsEventUnconditionallFallThrough, nil)

	// create a port
	p = NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// start timer tick machine
	p.TcMachineMain()

	// lets not start the main routines for the other state machines
	p.BEGIN(true)

	if p.PollingTimer != nil {
		p.PollingTimer.Stop()
	}

	// NOTE: must be called after BEGIN
	// Lets Instatiate but not run the following Machines
	// 1) Port Information Machine
	// 2) Port Protocol Migration Machine
	PrxmMachineFSMBuild(p)
	PrtMachineFSMBuild(p)
	PimMachineFSMBuild(p)
	PtxmMachineFSMBuild(p)
	BdmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	PstMachineFSMBuild(p)
	PpmmMachineFSMBuild(p)

	return p

}

func UsedForTestOnlyTcmTestTeardown(p *StpPort, t *testing.T) {

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
	if len(p.PrxmMachineFsm.PrxmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PtmMachineFsm.PtmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PstMachineFsm.PstEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PpmmMachineFsm.PpmmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	p.PrtMachineFsm = nil
	p.PimMachineFsm = nil
	p.PtxmMachineFsm = nil
	p.BdmMachineFsm = nil
	p.PrxmMachineFsm = nil
	p.PtmMachineFsm = nil
	p.PstMachineFsm = nil
	p.PpmmMachineFsm = nil

	b := p.b
	p.b.PrsMachineFsm = nil
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)
}

func ValidateInvalidState(p *StpPort, t *testing.T) {

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateInactive {
		t.Error("ERROR state not Inactive")
	}
	if !p.FdbFlush {
		t.Error("Flush Db not set to true")
	}
	if p.TcWhileTimer.count != 0 {
		t.Error("TC While not cleared")
	}
	if p.TcAck {
		t.Error("TC Ack not cleared")
	}

}

func ValidateLearningState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR state not Learning")
	}
	if p.RcvdTc {
		t.Error("Error RcvdTc is set")
	}
	if p.RcvdTcn {
		t.Error("Error RcvdTcn is set")
	}
	if p.RcvdTcAck {
		t.Error("Error RcvdTcAck is set")
	}
	if p.TcProp {
		t.Error("Error TcProp is set")
	}
}

func ValidateDetectedState(p *StpPort, t *testing.T) {

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateDetected {
		t.Error("ERROR state not Learning")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR state not Learning")
	}
	// newtcwhile
	// value of tcwhile ==0 and sendrstp -> tcwhile == hello time + 1 newInfo = true
	// value of tcwhile == 0 and !sendrstp tcwhile == maxage + fwddelay newInfo not changed
	if p.TcWhileTimer.count == BridgeHelloTimeDefault+1 {
		if !p.NewInfo {
			t.Error("ERROR NewInfo not set")
		}
		if !p.SendRSTP {
			t.Error("ERROR tchwile value should not be hello + 1 sendrstp is set to false")
		}
	} else if p.TcWhileTimer.count == int32(p.b.RootTimes.MaxAge+p.b.RootTimes.ForwardingDelay) {
		// do nothing
	} else {
		t.Error("ERROR Tcwhile not set properly")
	}

	// setTcPropTree set tcprop for all ports except caller
	if p.TcProp {
		t.Error("ERROR tcprop is set")
	}
	if !p.NewInfo {
		t.Error("ERROR new info not set")
	}
}

func ValidateActiveState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR state not Active")
	}
}

func ValidateNotifyTcnState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateNotifiedTcn {
		t.Error("ERROR state not NotifiedTcn")
	}
	// newtcwhile
	// value of tcwhile ==0 and sendrstp -> tcwhile == hello time + 1 newInfo = true
	// value of tcwhile == 0 and !sendrstp tcwhile == maxage + fwddelay newInfo not changed
	if p.TcWhileTimer.count == BridgeHelloTimeDefault+1 {
		if !p.NewInfo {
			t.Error("ERROR NewInfo not set")
		}
		if !p.SendRSTP {
			t.Error("ERROR tchwile value should not be hello + 1 sendrstp is set to false")
		}
	} else if p.TcWhileTimer.count == int32(p.b.RootTimes.MaxAge+p.b.RootTimes.ForwardingDelay) {
		// do nothing
	} else {
		t.Error("ERROR Tcwhile not set properly")
	}
}

func ValidateNotifyTcState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateNotifiedTc {
		t.Error("ERROR state not NotifiedTc")
	}
	if p.RcvdTcn {
		t.Error("ERROR RcvdTcn is set")
	}
	if p.RcvdTc {
		t.Error("ERROR RcvdTc is set")
	}
	if p.Role == PortRoleDesignatedPort &&
		!p.TcAck {
		t.Error("ERROR TcAck not set")
	}
	// setTcPropTree set tcprop for all ports except caller
	if p.TcProp {
		t.Error("ERROR tcprop is set")
	}
}

func ValidatePropagatingState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStatePropagating {
		t.Error("ERROR state not NotifiedTc")
	}
	// newtcwhile
	// value of tcwhile ==0 and sendrstp -> tcwhile == hello time + 1 newInfo = true
	// value of tcwhile == 0 and !sendrstp tcwhile == maxage + fwddelay newInfo not changed
	if p.TcWhileTimer.count == BridgeHelloTimeDefault+1 {
		if !p.NewInfo {
			t.Error("ERROR NewInfo not set")
		}
		if !p.SendRSTP {
			t.Error("ERROR tchwile value should not be hello + 1 sendrstp is set to false")
		}
	} else if p.TcWhileTimer.count == int32(p.b.RootTimes.MaxAge+p.b.RootTimes.ForwardingDelay) {
		// do nothing
	} else {
		t.Error("ERROR Tcwhile not set properly")
	}

	if !p.FdbFlush {
		t.Error("ERROR fdbflush not set")
	}
	if p.TcProp {
		t.Error("ERROR tcprop set")
	}
}

func ValidateAcknowledgedState(p *StpPort, t *testing.T) {
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateAcknowledged {
		t.Error("ERROR state not NotifiedTc")
	}

	if p.TcWhileTimer.count != 0 {
		t.Error("ERROR tcwhile not zero")
	}

	if p.RcvdTcAck {
		t.Error("ERROR RcvdTcAck not set")
	}
}

func UsedForTestStartTcLearningState(p *StpPort, t *testing.T) {
	testChan := make(chan string)
	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true
	p.TcProp = true
	p.FdbFlush = false

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

}

func UsedForTestStartTcDetectedState(p *StpPort, t *testing.T) {
	testChan := make(chan string)
	UsedForTestStartTcLearningState(p, t)

	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDesignatedPort
	p.Forward = true
	p.OperEdge = false

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventRoleEqualDesignatedPortAndForwardAndNotOperEdge,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	// note new state will be active as the next state is fallthrough

}

func TestTcmInactiveInvalidStateTest(t *testing.T) {
	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	invalidStateMap := [11]fsm.Event{
		TcEventUnconditionalFallThrough,
		TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
		TcEventRcvdTc,
		TcEventRcvdTcn,
		TcEventRcvdTcAck,
		TcEventTcProp,
		TcEventRoleEqualRootPortAndForwardAndNotOperEdge,
		TcEventRoleEqualDesignatedPortAndForwardAndNotOperEdge,
		TcEventTcPropAndNotOperEdge,
		TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPort,
		TcEventOperEdge,
	}

	// test the invalid states
	for _, e := range invalidStateMap {

		p.TcMachineFsm.TcEvents <- MachineEvent{e: e,
			src:          "TEST",
			responseChan: testChan}

		<-testChan

		if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateNone {
			t.Error("Previous State not as expected", p.TcMachineFsm.Machine.Curr.PreviousState())
		}

		ValidateInvalidState(p, t)
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactivePostStateProcessing(t *testing.T) {
	p := UsedForTestOnlyTcmTestSetup(t)

	ValidateInvalidState(p, t)

	// lets set what could be set
	p.Learn = true
	p.FdbFlush = false

	p.TcMachineFsm.ProcessPostStateInactive()
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR Did not transition to proper state")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningState(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_1(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleRootPort
	p.Forward = true
	p.OperEdge = false

	p.TcMachineFsm.ProcessPostStateProcessing()

	ValidateDetectedState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_2(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.Forward = true
	p.OperEdge = false

	p.TcMachineFsm.ProcessPostStateProcessing()

	ValidateDetectedState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_3(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTc = true

	p.TcMachineFsm.ProcessPostStateLearning()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state", p.TcMachineFsm.Machine.Curr.PreviousState())
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_4(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTcn = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_5(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTcAck = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmInactiveToLearningStatePostStateProcessing_6(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.TcProp = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateLearning {
		t.Error("ERROR: Did not transition to proper state")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveState(t *testing.T) {

	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	UsedForTestStartTcDetectedState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_1(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleAlternatePort

	p.TcMachineFsm.ProcessPostStateProcessing()

	ValidateLearningState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_2(t *testing.T) {

	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	p.RcvdTc = true
	p.RcvdTcn = true
	p.RcvdTcAck = true

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventLearnAndNotFdbFlush,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	ValidateLearningState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.OperEdge = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	ValidateLearningState(p, t)

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_rcvdTcn(t *testing.T) {

	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	UsedForTestStartTcDetectedState(p, t)

	ValidateDetectedState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTcn = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateNotifiedTc {
		t.Error("ERROR previous state not notify tcn", p.TcMachineFsm.Machine.Curr.PreviousState())
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active", p.TcMachineFsm.Machine.Curr.CurrentState())
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_rcvdTc(t *testing.T) {

	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	UsedForTestStartTcDetectedState(p, t)

	ValidateDetectedState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTc = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateNotifiedTc {
		t.Error("ERROR previous state not notify tcn")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_tcPropAndNotOperEdge(t *testing.T) {

	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	UsedForTestStartTcDetectedState(p, t)

	ValidateDetectedState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.TcProp = true
	p.OperEdge = false

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStatePropagating {
		t.Error("ERROR previous state not notify tcn")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToDetectedToActiveStatePostStateProcessing_rcvdTcAck(t *testing.T) {

	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	// sanity to ensure that the variables are cleared
	UsedForTestStartTcDetectedState(p, t)

	ValidateDetectedState(p, t)

	// setup
	p.Role = PortRoleDesignatedPort
	p.RcvdTcAck = true

	p.TcMachineFsm.ProcessPostStateProcessing()

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateAcknowledged {
		t.Error("ERROR previous state not acknowledged")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToActiveToNotifyTcnState(t *testing.T) {
	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	UsedForTestStartTcDetectedState(p, t)

	p.RcvdTcn = true
	p.TcProp = false
	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventRcvdTcn,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateNotifiedTc {
		t.Error("ERROR previous state not notify tcn", p.TcMachineFsm.Machine.Curr.PreviousState())
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active", p.TcMachineFsm.Machine.Curr.CurrentState())
	}

	// newtcwhile
	// value of tcwhile ==0 and sendrstp -> tcwhile == hello time + 1 newInfo = true
	// value of tcwhile == 0 and !sendrstp tcwhile == maxage + fwddelay newInfo not changed
	if p.TcWhileTimer.count == BridgeHelloTimeDefault+1 {
		if !p.NewInfo {
			t.Error("ERROR NewInfo not set")
		}
		if !p.SendRSTP {
			t.Error("ERROR tchwile value should not be hello + 1 sendrstp is set to false")
		}
	} else if p.TcWhileTimer.count == int32(p.b.RootTimes.MaxAge+p.b.RootTimes.ForwardingDelay) {
		// do nothing
	} else {
		t.Error("ERROR Tcwhile not set properly")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToActiveToNotifyTcState(t *testing.T) {
	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.OperEdge = false
	p.Learn = true
	p.FdbFlush = false

	UsedForTestStartTcDetectedState(p, t)

	// we should be in active state now
	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}
	p.RcvdTc = true
	p.TcProp = false
	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDesignatedPort

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventRcvdTc,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateNotifiedTc {
		t.Error("ERROR previous state not notify tcn")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}
	if p.RcvdTc {
		t.Error("ERROR RcvdTc set")
	}
	if p.RcvdTcn {
		t.Error("ERROR RcvdTcn set")
	}
	if !p.TcAck {
		t.Error("ERROR TcAck not set")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToActiveToPropagatingState(t *testing.T) {
	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	UsedForTestStartTcDetectedState(p, t)

	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDesignatedPort
	p.TcProp = true
	p.OperEdge = false
	p.FdbFlush = false

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventTcPropAndNotOperEdge,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStatePropagating {
		t.Error("ERROR previous state not notify tcn")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}
	if !p.FdbFlush {
		t.Error("ERROR FdbFlush not set")
	}
	if p.TcProp {
		t.Error("ERROR RcvdTcn set")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}

func TestTcmLearningToActiveToAcknowledgedState(t *testing.T) {
	testChan := make(chan string)
	p := UsedForTestOnlyTcmTestSetup(t)

	// setup the variables which should cause the state machine to transition
	p.Learn = true
	p.FdbFlush = false

	UsedForTestStartTcDetectedState(p, t)

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}

	p.RcvdTcAck = true
	p.Role = PortRoleDesignatedPort
	p.SelectedRole = PortRoleDesignatedPort

	p.TcMachineFsm.TcEvents <- MachineEvent{e: TcEventRcvdTcAck,
		src:          "TEST",
		responseChan: testChan}

	<-testChan

	if p.TcMachineFsm.Machine.Curr.PreviousState() != TcStateAcknowledged {
		t.Error("ERROR previous state not acknowledged")
	}

	if p.TcMachineFsm.Machine.Curr.CurrentState() != TcStateActive {
		t.Error("ERROR previous state not active")
	}
	if p.RcvdTcAck {
		t.Error("ERROR RcvdTcAck set")
	}
	if p.TcWhileTimer.count != 0 {
		t.Error("ERROR RcvdTcn tcwhile not cleared")
	}

	UsedForTestOnlyTcmTestTeardown(p, t)
}
