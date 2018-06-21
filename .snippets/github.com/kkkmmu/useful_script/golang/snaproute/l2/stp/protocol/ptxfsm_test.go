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

// rx_test.go
// This is a test file to test the rx/portrcvfsm
package stp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	//"strconv"
	//"strings"
	"fmt"
	"net"
	"testing"
	"time"
)

var PTX_TEST_RX_PORT_CONFIG_IFINDEX int32
var PTX_TEST_TX_PORT_CONFIG_IFINDEX int32

const PTX_TEST_TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL = time.Millisecond * 1
const PTX_TEST_NUM_DELAY_TRIES = 60

func UsedForTestOnlyTxInitPortConfigTest() {

	if PortConfigMap == nil {
		PortConfigMap = make(map[int32]portConfig)
	}
	// In order to test a packet we must listen on loopback interface
	// and send on interface we expect to receive on.  In order
	// to do this a couple of things must occur the PortConfig
	// must be updated with "dummy" ifindex pointing to 'lo'
	PTX_TEST_RX_PORT_CONFIG_IFINDEX = 0x0ADDBEEF
	PortConfigMap[PTX_TEST_RX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo",
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	PortConfigMap[PTX_TEST_TX_PORT_CONFIG_IFINDEX] = portConfig{Name: "lo",
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
	UsedForTestOnlySetupAsicDPlugin()
}

func UsedForTestOnlyPtxTestSetup(stpconfig *StpPortConfig, t *testing.T) (p *StpPort) {
	UsedForTestOnlyTxInitPortConfigTest()

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

	// create a port
	p = NewStpPort(stpconfig)
	b.StpPorts = append(b.StpPorts, p.IfIndex)

	// lets only start the Port Receive State Machine
	p.PtxmMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PtxmMachineFsm.Machine.Curr.PreviousState() != PtxmStateTransmitInit {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PtxmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	if p.PtxmMachineFsm.Machine.Curr.CurrentState() != PtxmStateIdle {
		t.Error("Failed to transition from None to Discard State")
		t.FailNow()
	}

	// NOTE: must be called after BEGIN
	// Lets Instatiate but not run the following Machines
	// 1) Port Information Machine
	// 2) Port Protocol Migration Machine
	PrtMachineFSMBuild(p)
	PimMachineFSMBuild(p)
	PrxmMachineFSMBuild(p)
	BdmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	TcMachineFSMBuild(p)
	PstMachineFSMBuild(p)
	PpmmMachineFSMBuild(p)

	return p

}

func UsedForTestOnlyPtxTestTeardown(p *StpPort, t *testing.T) {

	if len(p.PpmmMachineFsm.PpmmEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PimMachineFsm.PimEvents) > 0 {
		t.Error("Failed to check event sent")
	}
	if len(p.PrxmMachineFsm.PrxmEvents) > 0 {
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

	p.PrtMachineFsm = nil
	p.PimMachineFsm = nil
	p.PrxmMachineFsm = nil
	p.BdmMachineFsm = nil
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

func TestTxHelloWhenEqualZeroTransmitRSTP(t *testing.T) {

	UsedForTestOnlyTxInitPortConfigTest()

	testWait := make(chan bool)
	testChan := make(chan string)
	ifname, _ := PortConfigMap[PTX_TEST_RX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", PTX_TEST_RX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

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
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	p := UsedForTestOnlyPtxTestSetup(stpconfig, t)

	// setup event change, hellowhen expired should
	// transmit an rstp packet
	p.HelloWhenTimer.count = 0
	p.TxCount = 1
	p.Selected = true
	p.UpdtInfo = false
	p.NewInfo = false
	p.SendRSTP = true
	p.Role = PortRoleDesignatedPort

	p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
		e:            PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: testChan,
	}

	<-testChan

	if p.PtxmMachineFsm.Machine.Curr.CurrentState() != PtxmStateIdle {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.CurrentState()))
		t.FailNow()
	}

	if p.PtxmMachineFsm.Machine.Curr.PreviousState() != PtxmStateTransmitRSTP {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.PreviousState()))
		t.FailNow()
	}
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func(rx chan gopacket.Packet, t *testing.T) {
		for {
			select {
			case packet, _ := <-rx:
				if packet != nil {
					bpduLayer := packet.Layer(layers.LayerTypeBPDU)
					if bpduLayer == nil {
						continue
					}
					rstp := bpduLayer.(*layers.RSTP)
					if rstp == nil {
						t.Error("Did not get rstp packet as expected")
					}
					testWait <- true
					return
				}
			}
		}
	}(in, t)

	<-testWait

	// remove reference to fsm allocated above
	UsedForTestOnlyPtxTestTeardown(p, t)
	handle.Close()

	close(testWait)
	close(testChan)
}

func TestTxHelloWhenEqualZeroTransmitSTP(t *testing.T) {

	UsedForTestOnlyTxInitPortConfigTest()

	testWait := make(chan bool)
	testChan := make(chan string)
	ifname, _ := PortConfigMap[PTX_TEST_RX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", PTX_TEST_RX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

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
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	p := UsedForTestOnlyPtxTestSetup(stpconfig, t)

	// setup event change, hellowhen expired should
	// transmit an rstp packet
	p.HelloWhenTimer.count = 0
	p.TxCount = 1
	p.Selected = true
	p.UpdtInfo = false
	p.NewInfo = false
	p.SendRSTP = false
	p.Role = PortRoleDesignatedPort

	p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
		e:            PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: testChan,
	}

	<-testChan

	if p.PtxmMachineFsm.Machine.Curr.CurrentState() != PtxmStateIdle {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.CurrentState()))
		t.FailNow()
	}

	if p.PtxmMachineFsm.Machine.Curr.PreviousState() != PtxmStateTransmitConfig {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.PreviousState()))
		t.FailNow()
	}
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func(rx chan gopacket.Packet, t *testing.T) {
		for {
			select {
			case packet, _ := <-rx:
				if packet != nil {
					bpduLayer := packet.Layer(layers.LayerTypeBPDU)
					if bpduLayer == nil {
						continue
					}
					stp := bpduLayer.(*layers.STP)
					if stp == nil {
						t.Error("Did not get rstp packet as expected")
					}
					testWait <- true
					return
				}
			}
		}
	}(in, t)

	<-testWait

	// remove reference to fsm allocated above
	UsedForTestOnlyPtxTestTeardown(p, t)
	handle.Close()

	close(testWait)
	close(testChan)
}

func TestTxHelloWhenEqualZeroTransmitTCN(t *testing.T) {

	UsedForTestOnlyTxInitPortConfigTest()

	testWait := make(chan bool)
	testChan := make(chan string)
	ifname, _ := PortConfigMap[PTX_TEST_RX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", PTX_TEST_RX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

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
		BrgIfIndex:        DEFAULT_STP_BRIDGE_VLAN,
	}

	p := UsedForTestOnlyPtxTestSetup(stpconfig, t)

	// setup event change, hellowhen expired should
	// transmit an rstp packet
	p.HelloWhenTimer.count = 0
	p.TxCount = 1
	p.Selected = true
	p.UpdtInfo = false
	p.NewInfo = false
	p.SendRSTP = false
	p.Role = PortRoleRootPort

	p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
		e:            PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
		src:          "TEST",
		responseChan: testChan,
	}

	<-testChan

	if p.PtxmMachineFsm.Machine.Curr.CurrentState() != PtxmStateIdle {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.CurrentState()))
		t.FailNow()
	}

	if p.PtxmMachineFsm.Machine.Curr.PreviousState() != PtxmStateTransmitTCN {
		t.Error(fmt.Sprintf("Previous state not valid %d", p.PtxmMachineFsm.Machine.Curr.PreviousState()))
		t.FailNow()
	}
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func(rx chan gopacket.Packet, t *testing.T) {
		for {
			select {
			case packet, _ := <-rx:
				if packet != nil {
					bpduLayer := packet.Layer(layers.LayerTypeBPDU)
					if bpduLayer == nil {
						continue
					}
					bpdu := bpduLayer.(*layers.BPDUTopology)
					if bpdu == nil {
						t.Error("Did not get rstp packet as expected")
					}
					testWait <- true
					return
				}
			}
		}
	}(in, t)

	<-testWait

	// remove reference to fsm allocated above
	UsedForTestOnlyPtxTestTeardown(p, t)
	handle.Close()

	close(testWait)
	close(testChan)
}
