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
	//"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	//"strconv"
	//"strings"
	"testing"
	"time"
)

var TEST_RX_PORT_CONFIG_IFINDEX int32
var TEST_RX_PORT2_CONFIG_IFINDEX int32
var TEST_TX_PORT_CONFIG_IFINDEX int32

const TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL = time.Millisecond * 1
const NUM_DELAY_TRIES = 60

func UsedForTestOnlyRxInitPortConfigTest() {

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
	UsedForTestOnlySetupAsicDPlugin()
}

func UsedForTestOnlyPrxTestSetup(stpconfig *StpPortConfig, t *testing.T) (p *StpPort) {
	UsedForTestOnlyRxInitPortConfigTest()

	bridgeconfig := &StpBridgeConfig{
		Address:      "00:55:55:55:55:55",
		Priority:     0x20,
		MaxAge:       BridgeMaxAgeDefault,
		HelloTime:    BridgeHelloTimeDefault,
		ForwardDelay: BridgeForwardDelayDefault,
		ForceVersion: 2,
		TxHoldCount:  TransmitHoldCountDefault,
	}

	if stpconfig.BrgIfIndex != 0 {
		bridgeconfig.Vlan = uint16(stpconfig.BrgIfIndex)
	} else {
		stpconfig.BrgIfIndex = DEFAULT_STP_BRIDGE_VLAN
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
	p.PrxmMachineMain()
	p.PimMachineMain()

	// going just send event and not start main as we just did above
	p.BEGIN(true)

	if p.PrxmMachineFsm.Machine.Curr.PreviousState() != PrxmStateNone {
		t.Error("Failed to Initial Rx machine state not set correctly", p.PrxmMachineFsm.Machine.Curr.PreviousState())
		t.FailNow()
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateDiscard {
		t.Error("Failed to transition from None to Discard State")
		t.FailNow()
	}

	// lets advance the PIM machine to the current state
	// since pim depends on rx messages
	responseChan := make(chan string)
	p.Selected = true
	p.UpdtInfo = true
	p.PimMachineFsm.PimEvents <- MachineEvent{
		e:            PimEventSelectedAndUpdtInfo,
		src:          "TEST",
		responseChan: responseChan,
	}
	<-responseChan

	// NOTE: must be called after BEGIN
	// Lets Instatiate but not run the following Machines
	// 1) Port Information Machine
	// 2) Port Protocol Migration Machine
	PrtMachineFSMBuild(p)
	PtxmMachineFSMBuild(p)
	BdmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	PtmMachineFSMBuild(p)
	TcMachineFSMBuild(p)
	PstMachineFSMBuild(p)
	PpmmMachineFSMBuild(p)

	return p

}

func UsedForTestOnlyPrxTestTeardown(p *StpPort, t *testing.T) {

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
	p.PrtMachineFsm = nil
	p.PtxmMachineFsm = nil
	p.BdmMachineFsm = nil
	p.PtmMachineFsm = nil
	p.TcMachineFsm = nil
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

func UsedForTestOnlySendValidStpTopoFrame(txifindex int32, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.BPDUTopologyLength + 3), // found from PCAP from packetlife.net
	}

	llc := layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}

	topo := layers.BPDUTopology{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.STPProtocolVersion,
		BPDUType:          layers.StpBpduType(layers.BPDUTypeTopoChange),
	}

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &llc, &topo)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}

	handle.Close()
	handle = nil
}

func UsedForTestOnlySendValidStpFrame(txifindex int32, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.STPProtocolLength + 3), // found from PCAP from packetlife.net
	}

	llc := layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}

	stp := layers.STP{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.STPProtocolVersion,
		BPDUType:          layers.StpBpduType(layers.BPDUTypeSTP),
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            0,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
	}

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &llc, &stp)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}

	handle.Close()
	handle = nil
}

func UsedForTestOnlySendValidRStpFrame(txifindex int32, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.RSTPProtocolLength + 3), // found from PCAP from packetlife.net
	}

	llc := layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}

	stp := layers.RSTP{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.RSTPProtocolVersion,
		BPDUType:          layers.BPDUTypeRSTP,
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            0,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
		Version1Length:    0,
	}

	var flags uint8
	StpSetBpduFlags(0, 0, 0, 0, ConvertRoleToPktRole(PortRoleDesignatedPort), 1, 0, &flags)

	stp.Flags = layers.StpFlags(flags)

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &llc, &stp)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}
	handle.Close()
	handle = nil
}

func UsedForTestOnlySendValidPVSTFrame(txifindex int32, pvstvlan uint16, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduPVSTDMAC,
		// length
		EthernetType: layers.EthernetTypeDot1Q,
	}

	vlan := layers.Dot1Q{
		Priority:       PVST_VLAN_PRIORITY,
		DropEligible:   false,
		VLANIdentifier: pvstvlan,
		Type:           layers.EthernetType(layers.PVSTProtocolLength + 3 + 5), // length
	}

	llc := layers.LLC{
		DSAP:    0xAA,
		IG:      false,
		SSAP:    0xAA,
		CR:      false,
		Control: 0x03,
	}

	snap := layers.SNAP{
		OrganizationalCode: []byte{0x00, 0x00, 0x0C},
		Type:               0x010b,
	}

	stp := layers.PVST{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.RSTPProtocolVersion,
		BPDUType:          layers.BPDUTypePVST,
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x64, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            0,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
		Version1Length:    0,
		OriginatingVlan: layers.STPOriginatingVlanTlv{
			Type:     0,
			Length:   2,
			OrigVlan: pvstvlan,
		},
	}

	var flags uint8
	StpSetBpduFlags(0, 0, 0, 0, ConvertRoleToPktRole(PortRoleDesignatedPort), 1, 0, &flags)

	stp.Flags = layers.StpFlags(flags)

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &vlan, &llc, &snap, &stp)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}
	handle.Close()
	handle = nil
}

func UsedForTestOnlySendInvalidStpFrame(txifindex int32, stp *layers.STP, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.STPProtocolLength + 3), // found from PCAP from packetlife.net
	}

	llc := layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}
	var flags uint8
	StpSetBpduFlags(0, 0, 0, 0, ConvertRoleToPktRole(PortRoleDesignatedPort), 1, 0, &flags)

	stp.Flags = layers.StpFlags(flags)

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &llc, stp)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}
	handle.Close()
	handle = nil
}

func UsedForTestOnlySendInvalidRStpFrame(txifindex int32, rstp *layers.RSTP, t *testing.T) {
	ifname, _ := PortConfigMap[TEST_TX_PORT_CONFIG_IFINDEX]
	handle, err := pcap.OpenLive(ifname.Name, 65536, true, 50*time.Millisecond)
	if err != nil {
		t.Error("Error opening pcap TX interface", TEST_TX_PORT_CONFIG_IFINDEX, ifname.Name, err)
		return
	}
	//txIface, _ := net.InterfaceByName(ifname.Name)

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.RSTPProtocolLength + 3), // found from PCAP from packetlife.net
	}

	llc := layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}

	var flags uint8
	StpSetBpduFlags(0, 0, 0, 0, ConvertRoleToPktRole(PortRoleDesignatedPort), 1, 0, &flags)

	rstp.Flags = layers.StpFlags(flags)

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &llc, rstp)
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		t.Error("Error writing packet to interface")
	}
	handle.Close()
	handle = nil
}

func TestRxValidStpPacket(t *testing.T) {

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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)
	// force timeout, if this does not happen then event will not be sent to PPM
	p.MdelayWhiletimer.count = 0

	// send a packet
	UsedForTestOnlySendValidStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed receive RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed  OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != true {
		t.Error("Failed RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP == true {
		t.Error("Failed RcvdRSTP is set")
		t.FailNow()
	}

	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*

		if p.RcvdMsg != true {
			t.Error("Failed RcvdMsg not set")
			t.FailNow()
		}
	*/

	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateReceive {
		t.Error("Failed to transition state to Receive")
		t.FailNow()
	}

	// we should have received an event from rx machine
	rx, _ := <-p.PpmmMachineFsm.PpmmEvents
	if rx.e != PpmmEventSendRSTPAndRcvdSTP {
		t.Error("Failed to transition state to Receive")
		t.FailNow()
	}

	// TODO add Pim event to test

	// remove reference to fsm allocated above
	UsedForTestOnlyPrxTestTeardown(p, t)

}

func TestRxValidRStpPacket(t *testing.T) {
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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// setup pre-condition, lets fake out and pretent we were in send STP mode
	p.SendRSTP = false

	// send a packet
	UsedForTestOnlySendValidRStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != false {
		t.Error("Failed RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP != true {
		t.Error("Failed RcvdRSTP not set")
		t.FailNow()
	}

	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg != true {
			t.Error("Failed RcvdMsg not set")
			t.FailNow()
		}
	*/

	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateReceive {
		t.Error("Failed to transition state to Receive")
		t.FailNow()
	}

	// we should have received an event from rx machine
	rx, _ := <-p.PpmmMachineFsm.PpmmEvents
	if rx.e != PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP {
		t.Error("Failed PPMM received invalid event")
		t.FailNow()
	}

	// remove reference to fsm allocated above
	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxValidPVSTPacket(t *testing.T) {
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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// setup pre-condition, lets fake out and pretent we were in send STP mode
	p.SendRSTP = false

	// send a packet
	UsedForTestOnlySendValidPVSTFrame(TEST_TX_PORT_CONFIG_IFINDEX, 100, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.PvstRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.PvstRx == 0 {
		t.Errorf("Failed Rx PVST count did not increment")
	}

	if p.RcvdBPDU == true {
		t.Errorf("Failed RcvdBPDU is set")
	}

	if p.OperEdge == true {
		t.Errorf("Failed OperEdge is set")
	}

	if p.RcvdSTP != false {
		t.Errorf("Failed RcvdSTP is set")
	}
	if p.RcvdRSTP != true {
		t.Errorf("Failed RcvdRSTP not set")
	}

	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*

		if p.RcvdMsg != true {
			t.Errorf("Failed RcvdMsg not set")
		}
	*/
	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Errorf("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateReceive {
		t.Errorf("Failed to transition state to Receive")
	}

	// we should have received an event from rx machine
	rx, _ := <-p.PpmmMachineFsm.PpmmEvents
	if rx.e != PpmmEventRstpVersionAndNotSendRSTPAndRcvdRSTP {
		t.Errorf("Failed PPMM received invalid event")
	}

	// remove reference to fsm allocated above
	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxInvalidRStpPacketBPDUTypeInvalid(t *testing.T) {
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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)
	// send a packet
	rstp := layers.RSTP{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.RSTPProtocolVersion,
		BPDUType:          layers.BPDUTypeSTP,
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            0,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
		Version1Length:    0,
	}

	UsedForTestOnlySendInvalidRStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, &rstp, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed to receive RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed to receive OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != false {
		t.Error("Failed to receive RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP != false {
		t.Error("Failed received RcvdRSTP is set")
		t.FailNow()
	}
	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg != false {
			t.Error("Failed received RcvdMsg not set")
			t.FailNow()
		}
	*/
	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxInvalidRStpPacketProtocolVersionInvalid(t *testing.T) {
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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// send a packet
	rstp := layers.RSTP{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.STPProtocolVersion,
		BPDUType:          layers.BPDUTypeRSTP,
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            0,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
		Version1Length:    0,
	}

	UsedForTestOnlySendInvalidRStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, &rstp, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed to receive RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed to receive OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != false {
		t.Error("Failed to receive RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP != false {
		t.Error("Failed received RcvdRSTP is set")
		t.FailNow()
	}
	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg != false {
			t.Error("Failed received RcvdMsg not set")
			t.FailNow()
		}
	*/
	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxInvalidStpPacketMsgAgeGreaterMaxAge(t *testing.T) {
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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// send a packet
	stp := layers.STP{
		ProtocolId:        layers.RSTPProtocolIdentifier,
		ProtocolVersionId: layers.RSTPProtocolVersion,
		BPDUType:          layers.BPDUTypeRSTP,
		Flags:             0,
		RootId:            [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		RootPathCost:      1,
		BridgeId:          [8]byte{0x80, 0x01, 0x00, 0x19, 0x06, 0xEA, 0xB8, 0x80},
		PortId:            0x1111,
		MsgAge:            21,
		MaxAge:            20,
		HelloTime:         2,
		FwdDelay:          15,
	}

	UsedForTestOnlySendInvalidStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, &stp, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed to receive RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed to receive OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != false {
		t.Error("Failed to receive RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP != false {
		t.Error("Failed received RcvdRSTP is set")
		t.FailNow()
	}
	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg != false {
			t.Error("Failed received RcvdMsg not set")
			t.FailNow()
		}
	*/

	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxSendValidRstpPacketOnDisabledPort(t *testing.T) {

	// configure a port
	stpconfig := &StpPortConfig{
		IfIndex:           TEST_RX_PORT_CONFIG_IFINDEX,
		Priority:          0x80,
		Enable:            false,
		PathCost:          1,
		ProtocolMigration: 0,
		AdminPointToPoint: StpPointToPointForceFalse,
		AdminEdgePort:     false,
		AdminPathCost:     0,
	}

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// send a packet
	UsedForTestOnlySendValidRStpFrame(TEST_TX_PORT_CONFIG_IFINDEX, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU == true {
		t.Error("Failed to receive RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge == true {
		t.Error("Failed to receive OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP == true {
		t.Error("Failed to receive RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP == true {
		t.Error("Failed received RcvdRSTP is set")
		t.FailNow()
	}
	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg == true {
			t.Error("Failed received RcvdMsg not set")
			t.FailNow()
		}
	*/

	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateDiscard {
		t.Error("Failed to state transitioned out of Discard State")
		t.FailNow()
	}

	UsedForTestOnlyPrxTestTeardown(p, t)
}

func TestRxValidTopoChange(t *testing.T) {
	UsedForTestOnlyRxInitPortConfigTest()

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

	p := UsedForTestOnlyPrxTestSetup(stpconfig, t)

	// force state of tc state machine
	p.TcMachineFsm.Machine.Curr.SetState(TcStateLearning)
	p.SendRSTP = true
	// force timeout of mdelay
	p.MdelayWhiletimer.count = 0

	// send a packet
	UsedForTestOnlySendValidStpTopoFrame(TEST_TX_PORT_CONFIG_IFINDEX, t)

	testWait := make(chan bool)
	// may need to delay a bit in order to allow for packet to be receive
	// by pcap
	go func() {
		for i := 0; i < NUM_DELAY_TRIES &&
			(p.BpduRx == 0); i++ {
			time.Sleep(TIME_TO_DELAY_TO_WAIT_FOR_PACKET_ARRIVAL)
		}
		testWait <- true
	}()

	<-testWait

	if p.RcvdBPDU != false {
		t.Error("Failed  RcvdBPDU is set")
		t.FailNow()
	}

	if p.OperEdge != false {
		t.Error("Failed OperEdge is set")
		t.FailNow()
	}

	if p.RcvdSTP != true {
		t.Error("Failed RcvdSTP is set")
		t.FailNow()
	}
	if p.RcvdRSTP != false {
		t.Error("Failed RcvdRSTP is set")
		t.FailNow()
	}

	// Utilizing PIM in this test and it is clearning the RcvdMsg
	/*
		if p.RcvdMsg != true {
			t.Error("Failed RcvdMsg not set")
			t.FailNow()
		}
	*/

	if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
		t.Error("Failed EdgeDelayWhiletimer tick count not set to MigrateTimeDefault")
		t.FailNow()
	}

	if p.PrxmMachineFsm.Machine.Curr.CurrentState() != PrxmStateReceive {
		t.Error("Failed to transition state to Receive")
		t.FailNow()
	}

	// we should have received an event from rx machine
	rx, _ := <-p.PpmmMachineFsm.PpmmEvents
	if rx.e != PpmmEventSendRSTPAndRcvdSTP {
		t.Error("Failed PPMM received invalid event")
		t.FailNow()
	}

	tc, _ := <-p.TcMachineFsm.TcEvents
	if tc.e != TcEventRcvdTcn {
		t.Error("Failed to get proper tc event")
	}

	UsedForTestOnlyPrxTestTeardown(p, t)
}
