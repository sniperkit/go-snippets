// prsfsm_test.go
package stp

import (
	"testing"
	"time"
)

var WAIT_FOR_EVENT_TIME time.Duration = time.Millisecond * 75

func UsedForTestOnlyPrsInitPortConfigTest() {

	if PortConfigMap == nil {
		PortConfigMap = make(map[int32]portConfig)
	}
	// In order to test a packet we must listen on loopback interface
	// and send on interface we expect to receive on.  In order
	// to do this a couple of things must occur the PortConfig
	// must be updated with "dummy" ifindex pointing to 'lo'
	TEST_RX_PORT_CONFIG_IFINDEX = 0x0ADDBEEF
	TEST_RX_PORT2_CONFIG_IFINDEX = 0x0ADDBEF0
	PortConfigMap[TEST_RX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo"}
	PortConfigMap[TEST_RX_PORT2_CONFIG_IFINDEX] = portConfig{Name: "lo"}
	PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo"}
	StpBridgeMac = [6]uint8{0x00, 0x55, 0x55, 0x55, 0x55, 0x55}
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

// TestPrsUpdtRolesTreePortIsRootBridge1 update role should update root as the behind the port msg was received
// because bridge id is 'superior'
func TestPrsUpdtRolesTreePortIsRootBridge1(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	p.PortPriority.RootBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	p.PortPriority.RootPathCost = 20000
	p.PortPriority.DesignatedBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	p.PortPriority.DesignatedPortId = 10
	p.PortPriority.BridgePortId = 10

	p.RcvdRSTP = true // set by port receive state machine
	p.InfoIs = PortInfoStateReceived
	p.Selected = true // assumed port role selection state machine set this
	p.PortEnabled = true

	// configure a second port
	stpconfig = &StpPortConfig{
		IfIndex:           TEST_RX_PORT2_CONFIG_IFINDEX,
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
	p2 := NewStpPort(stpconfig)
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p2.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p2.PortPriority.RootBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	p2.PortPriority.RootPathCost = 200
	p2.PortPriority.DesignatedBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x00, 0x11, 0x22, 0x33, 0x44}
	p2.PortPriority.DesignatedPortId = 10
	p2.PortPriority.BridgePortId = 10

	p2.RcvdRSTP = true // set by port receive state machine
	p2.InfoIs = PortInfoStateReceived
	p2.Selected = true // assumed port role selection state machine set this
	p2.PortEnabled = true

	b.PrsMachineFsm.PrsEvents <- MachineEvent{
		e:            PrsEventReselect,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if b.BridgePriority.RootBridgeId != p.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}
	if b.BridgePriority.RootBridgeId != p2.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}

	if p.SelectedRole != PortRoleDesignatedPort {
		t.Error("ERROR: port should be a designated port", p.SelectedRole)
	}

	if p2.SelectedRole != PortRoleRootPort {
		t.Error("ERROR: port should be a designated port", p2.SelectedRole)
	}

	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	for idx, ifindex := range b.StpPorts {
		if ifindex == p2.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}
	DelStpPort(p2)
	DelStpBridge(b, true)
}

// TestPrsUpdtRolesTreePortIsRootBridge2 we are the root thus all ports are designated
func TestPrsUpdtRolesTreePortIsRootBridge2(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p.PortPriority.RootBridgeId = [8]uint8{0x0, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55}
	p.PortPriority.RootPathCost = 2000
	p.PortPriority.DesignatedBridgeId = [8]uint8{0x0, 0x00, 0x00, 0x22, 0x22, 0x22, 0x22, 0x22}
	p.PortPriority.DesignatedPortId = 10
	p.PortPriority.BridgePortId = 10

	p.RcvdRSTP = true // set by port receive state machine
	p.InfoIs = PortInfoStateReceived
	p.Selected = true // assumed port role selection state machine set this

	// configure a port
	stpconfig = &StpPortConfig{
		IfIndex:           TEST_RX_PORT2_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          10,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p2 := NewStpPort(stpconfig)
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p2.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p2.PortPriority.RootBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55}
	p2.PortPriority.RootPathCost = 200
	p2.PortPriority.DesignatedBridgeId = [8]uint8{0x0, 0x00, 0x00, 0x22, 0x22, 0x22, 0x22, 0x22}
	p2.PortPriority.DesignatedPortId = 11
	p2.PortPriority.BridgePortId = 11

	p2.RcvdRSTP = true // set by port receive state machine
	p2.InfoIs = PortInfoStateReceived
	p2.Selected = true // assumed port role selection state machine set this
	p2.PortEnabled = true

	b.PrsMachineFsm.PrsEvents <- MachineEvent{
		e:            PrsEventReselect,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	// trigger selection again as the first time the path cost is not set against the proper port yet
	b.PrsMachineFsm.PrsEvents <- MachineEvent{
		e:            PrsEventReselect,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if b.BridgePriority.RootBridgeId != p.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}
	if b.BridgePriority.RootBridgeId != p2.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}

	if p.SelectedRole != PortRoleDesignatedPort {
		t.Error("ERROR: port should be a designated port", p.SelectedRole)
	}

	if p2.SelectedRole != PortRoleDesignatedPort {
		t.Error("ERROR: port should be a designated port", p2.SelectedRole)
	}
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	for idx, ifindex := range b.StpPorts {
		if ifindex == p2.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p2)
	DelStpBridge(b, true)
}

// TestPrsUpdtRolesTreePortIsRootBridge3 we are the root one port is disabled
func TestPrsUpdtRolesTreePortIsRootBridge3(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p.PortPriority.RootBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55}
	p.PortPriority.RootPathCost = 2000
	p.PortPriority.DesignatedBridgeId = [8]uint8{0x00, 0x00, 0x00, 0x22, 0x22, 0x22, 0x22, 0x22}
	p.PortPriority.DesignatedPortId = 10
	p.PortPriority.BridgePortId = 10

	p.RcvdRSTP = true // set by port receive state machine
	p.InfoIs = PortInfoStateReceived
	p.Selected = true // assumed port role selection state machine set this

	// configure a port
	stpconfig = &StpPortConfig{
		IfIndex:           TEST_RX_PORT2_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          10,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p2 := NewStpPort(stpconfig)
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p2.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	//p2.PortPriority.RootBridgeId = [8]uint8{0xf, 0xff, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55}
	//p2.PortPriority.RootPathCost = 200
	//p2.PortPriority.DesignatedBridgeId = [8]uint8{0xf, 0xff, 0x00, 0x22, 0x22, 0x22, 0x22, 0x22}
	//p2.PortPriority.DesignatedPortId = 11
	//p2.PortPriority.BridgePortId = 11

	p2.RcvdRSTP = true // set by port receive state machine
	p2.InfoIs = PortInfoStateReceived
	p2.Selected = true // assumed port role selection state machine set this
	p2.PortEnabled = false

	b.PrsMachineFsm.PrsEvents <- MachineEvent{
		e:            PrsEventReselect,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if b.BridgePriority.RootBridgeId != p.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}
	if b.BridgePriority.RootBridgeId != p2.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}

	if p.SelectedRole != PortRoleDesignatedPort {
		t.Error("ERROR: port should be a designated port", p.SelectedRole)
	}

	if p2.SelectedRole != PortRoleDisabledPort {
		t.Error("ERROR: port should be a designated port", p2.SelectedRole)
	}
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	for idx, ifindex := range b.StpPorts {
		if ifindex == p2.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p2)
	DelStpBridge(b, true)
}

// TestPrsUpdtRolesTreePortIsRootBridge4 bridge priority vector not higher than port priority
func xxxxTestPrsUpdtRolesTreePortIsRootBridge4(t *testing.T) {
	testChan := make(chan string)
	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p.PortPriority.RootBridgeId = [8]uint8{0x5, 0xff, 0x00, 0x33, 0x33, 0x33, 0x33, 0x33}
	p.PortPriority.RootPathCost = 20
	p.PortPriority.DesignatedBridgeId = [8]uint8{0x5, 0xff, 0x00, 0x33, 0x33, 0x33, 0x33, 0x33}
	p.PortPriority.DesignatedPortId = 10
	p.PortPriority.BridgePortId = 10

	p.RcvdRSTP = true // set by port receive state machine
	p.InfoIs = PortInfoStateReceived
	p.Selected = true // assumed port role selection state machine set this

	// configure a port
	stpconfig = &StpPortConfig{
		IfIndex:           TEST_RX_PORT2_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            true,
		PathCost:          10,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	// create a port
	p2 := NewStpPort(stpconfig)
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p2.IfIndex)

	// UpdtRolesTree
	// (a)
	// assume message received and recorded in PortPriority
	p2.PortPriority.RootBridgeId = [8]uint8{0xf, 0xff, 0x00, 0x33, 0x33, 0x33, 0x33, 0x33}
	p2.PortPriority.RootPathCost = 0
	p2.PortPriority.DesignatedBridgeId = [8]uint8{0xf, 0xff, 0x00, 0x33, 0x33, 0x33, 0x33, 0x33}
	p2.PortPriority.DesignatedPortId = 9
	p2.PortPriority.BridgePortId = 9

	p2.RcvdRSTP = true // set by port receive state machine
	p2.InfoIs = PortInfoStateReceived
	p2.Selected = true // assumed port role selection state machine set this
	p2.PortEnabled = true

	b.PrsMachineFsm.PrsEvents <- MachineEvent{
		e:            PrsEventReselect,
		src:          "TEST",
		responseChan: testChan,
	}
	<-testChan

	if b.BridgePriority.RootBridgeId != p.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}
	if b.BridgePriority.RootBridgeId == p2.PortPriority.RootBridgeId {
		t.Error("ERROR: received bridge is better than current why was it not stored")
	}

	if p.SelectedRole != PortRoleRootPort {
		t.Error("ERROR: port should be a designated port", p.SelectedRole)
	}

	if p2.SelectedRole != PortRoleDesignatedPort {
		t.Error("ERROR: port should be a designated port", p2.SelectedRole)
	}

	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	for idx, ifindex := range b.StpPorts {
		if ifindex == p2.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p2)
	DelStpBridge(b, true)
}

// disable port -> diabled port via notify selected event
func TestPrsSetSelectedTreeEventNotify_DisabledPortStates_1(t *testing.T) {
	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	//p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// disabled port -> diabled port via notify selected event
func TestPrsSetSelectedTreeEventNotify_DisabledPortStates_2(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.FdWhileTimer.count = 5 // random choice
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	go func() {
		time.Sleep(WAIT_FOR_EVENT_TIME)
		if p.PrtMachineFsm != nil {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   0, // invalid event
				src: "TEST",
			}
		}
	}()

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// disabled port -> diabled port via notify selected event
func TestPrsSetSelectedTreeEventNotify_DisabledPortStates_3(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Sync = true
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// disabled port -> diabled port via notify selected event
func TestPrsSetSelectedTreeEventNotify_DisabledPortStates_4(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.ReRoot = true
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// disabled port -> diabled port via notify selected event
func TestPrsSetSelectedTreeEventNotify_DisabledPortStates_5(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Synced = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateDisabledPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_1(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.RrWhileTimer.count = 2 // random value
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil && p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in root port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_2(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.ReRoot = true
	p.Forward = true
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateReRooted {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in root port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_3(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.FdWhileTimer.count = 10
	p.RbWhileTimer.count = 0
	p.RstpVersion = true
	p.Learn = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootLearn {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in root port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_4(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.FdWhileTimer.count = 0
	p.RstpVersion = true
	p.Learn = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootLearn {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_5(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.FdWhileTimer.count = 4
	p.RbWhileTimer.count = 0
	p.RstpVersion = true
	p.ReRoot = true
	p.Learn = true
	p.Forward = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootForward {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_6(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.FdWhileTimer.count = 0
	p.RrWhileTimer.count = int32(p.b.RootTimes.ForwardingDelay)
	p.RbWhileTimer.count = 0
	p.ReRoot = true
	p.RstpVersion = true
	p.Learn = true
	p.Forward = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootForward {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}
	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_7(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Proposed = true
	p.Agree = false
	p.ReRoot = true
	p.RstpVersion = true
	p.Learn = true
	p.Forward = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootProposed {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_8(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Proposed = true
	p.Agree = true
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootAgreed {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_9(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Synced = true
	p.Agree = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateRootAgreed {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// root port -> root port via notify selected event
func TestPrsSetSelectedTreeEventNotify_RootPortStates_10(t *testing.T) {

	UsedForTestOnlyPrsInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
		DebugLevel:   2,
		Vlan:         DEFAULT_STP_BRIDGE_VLAN,
	}

	//StpBridgeCreate
	StpBridgeCreate(bridgeconfig)
	key := BridgeKey{
		Vlan: bridgeconfig.Vlan,
	}

	var b *Bridge
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: did not find bridge that was just created")
	}

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
	// Don't want to trigger the BEGIN call so going to just add the port to bridge manually
	//	StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
	b.StpPorts = append(b.StpPorts, p.IfIndex)
	p.PrtMachineMain()
	p.BEGIN(true)

	// simulate message call with proper port attributes set
	p.Role = PortRoleDisabledPort
	p.SelectedRole = PortRoleRootPort
	p.Learning = false
	p.Forwarding = false
	p.Selected = true
	p.UpdtInfo = false
	p.Synced = true

	// transition to disabled port
	responseChan := make(chan string)
	p.PrtMachineFsm.PrtEvents <- MachineEvent{
		e:            PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}

	<-responseChan
	// check that we transitioned
	if p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort {
		t.Error("ERROR failed to transition state")
	}

	// start test
	p.Forward = false
	p.ReRoot = false
	p.Selected = true
	p.UpdtInfo = false

	// call method
	p.NotifySelectedChanged(PrsMachineModuleStr, false, true)

	testwait := make(chan bool)

	go func(tw chan bool) {

		for i := 0; i < 10; i++ {
			if p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() != PrtStateRootPort &&
				p.PrtMachineFsm.Machine.Curr.PreviousState() != PrtStateReRoot {
				time.Sleep(WAIT_FOR_EVENT_TIME)
			} else {
				tw <- true
			}
		}
		tw <- false

	}(testwait)

	result := <-testwait
	if !result {
		t.Error("ERROR: PRT state transition did not occur port should in disabled port state")
	}

	// teardown
	for idx, ifindex := range b.StpPorts {
		if ifindex == p.IfIndex {
			b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
		}
	}

	DelStpPort(p)
	DelStpBridge(b, true)

}

// TODO added notify selected for designated port
