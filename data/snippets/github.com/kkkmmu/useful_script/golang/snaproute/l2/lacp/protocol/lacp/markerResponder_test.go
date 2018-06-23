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

// markerResponder_test.go
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestTwoAggsBackToBackSinglePortInjectMarkerPdu(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	const LaAggPortActor = 10
	const LaAggPortPeer = 20
	LaAggPortActorIf := "SIMeth0"
	LaAggPortPeerIf := "SIM2eth0"
	// must be called to initialize the global
	LaSystemActor := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x64}}
	LaSystemPeer := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	bridge := SimulationBridge{
		Port1:       LaAggPortActor,
		Port2:       LaAggPortPeer,
		RxLacpPort1: make(chan gopacket.Packet, 10),
		RxLacpPort2: make(chan gopacket.Packet, 10),
	}

	ActorSystem := LacpSysGlobalInfoInit(LaSystemActor)
	PeerSystem := LacpSysGlobalInfoInit(LaSystemPeer)
	ActorSystem.LaSysGlobalRegisterTxCallback(LaAggPortActorIf, bridge.TxViaGoChannel)
	PeerSystem.LaSysGlobalRegisterTxCallback(LaAggPortPeerIf, bridge.TxViaGoChannel)

	p1conf := &LaAggPortConfig{
		Id:     LaAggPortActor,
		Prio:   0x80,
		Key:    100,
		AggId:  100,
		Enable: true,
		Mode:   LacpModeActive,
		//Timeout: LacpFastPeriodicTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortActorIf,
		TraceEna: true,
	}

	utils.PortConfigMap[int32(p1conf.Id)] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	p2conf := &LaAggPortConfig{
		Id:     LaAggPortPeer,
		Prio:   0x80,
		Key:    200,
		AggId:  200,
		Enable: true,
		Mode:   LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortPeer, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortPeerIf,
		TraceEna: true,
	}

	utils.PortConfigMap[int32(p2conf.Id)] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	// lets create a port and start the machines
	CreateLaAggPort(p1conf)
	CreateLaAggPort(p2conf)

	// port 1
	LaRxMain(bridge.Port1, bridge.RxLacpPort1)
	// port 2
	LaRxMain(bridge.Port2, bridge.RxLacpPort2)

	a1conf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x01, 0x01, 0x01},
		Id:   100,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:64",
			SystemPriority: 128},
	}

	a2conf := &LaAggConfig{
		Name: "agg2",
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x02, 0x02, 0x02},
		Id:   200,
		Key:  200,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:C8",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(a1conf)
	CreateLaAgg(a2conf)

	// Add port to agg
	//AddLaAggPortToAgg(a1conf.Id, p1conf.Id)
	//AddLaAggPortToAgg(a2conf.Id, p2conf.Id)

	//time.Sleep(time.Second * 30)
	testWait := make(chan bool)

	var p1 *LaAggPort
	var p2 *LaAggPort
	if LaFindPortById(p1conf.Id, &p1) &&
		LaFindPortById(p2conf.Id, &p2) {

		go func() {
			for i := 0; i < 10 &&
				(p1.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDistributing ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDistributing); i++ {
				time.Sleep(time.Second * 1)
			}
			testWait <- true
		}()

		<-testWait

		State1 := GetLaAggPortActorOperState(p1conf.Id)
		State2 := GetLaAggPortActorOperState(p2conf.Id)

		const portUpState = LacpStateActivityBit | LacpStateAggregationBit |
			LacpStateSyncBit | LacpStateCollectingBit | LacpStateDistributingBit

		if !LacpStateIsSet(State1, portUpState) {
			t.Error(fmt.Sprintf("Actor Port State 0x%x did not come up properly with peer expected 0x%x", State1, portUpState))
		}
		if !LacpStateIsSet(State2, portUpState) {
			t.Error(fmt.Sprintf("Peer Port State 0x%x did not come up properly with actor expected 0x%x", State2, portUpState))
		}

		// lets inject a marker pdu now
		lamp := &layers.LAMP{
			Version: layers.LAMPVersion1,
			Marker: layers.LAMPMarkerTlv{TlvType: layers.LAMPTLVMarkerInfo,
				Length:                 layers.LAMPMarkerTlvLength,
				RequesterPort:          LaAggPortActor,
				RequesterSystem:        LaSystemActor.Actor_System,
				RequesterTransactionId: 10,
			},
			Terminator: layers.LAMPTerminatorTlv{},
		}
		// don't support generating marker but we certainly can inject one
		// for this test
		bridge.TxViaGoChannel(LaAggPortActor, lamp)
		go func() {
			for i := 0; i < 10 &&
				(p1.LacpCounter.AggPortStatsMarkerResponsePDUsRx == 0 ||
					p2.LacpCounter.AggPortStatsMarkerPDUsRx == 0 ||
					p2.LacpCounter.AggPortStatsMarkerResponsePDUsTx == 0); i++ {
				time.Sleep(time.Second * 1)
			}
			testWait <- true
		}()
		<-testWait

		if p2.LacpCounter.AggPortStatsMarkerPDUsRx == 0 {
			t.Error("ERROR: Failed to receive Marker PDU")
		}
		if p2.LacpCounter.AggPortStatsMarkerResponsePDUsTx == 0 {
			t.Error("ERROR: Failed to send Marker Response PDU")
		}
		if p1.LacpCounter.AggPortStatsMarkerResponsePDUsRx == 0 {
			t.Error("ERROR: Failed to transmit a Marker PDU")
		}

	} else {
		t.Error("Unable to find port just created")
	}

	// cleanup the provisioning
	close(bridge.RxLacpPort1)
	close(bridge.RxLacpPort2)
	bridge.RxLacpPort1 = nil
	bridge.RxLacpPort2 = nil
	DeleteLaAgg(a1conf.Id)
	DeleteLaAgg(a2conf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	OnlyForTestTeardown()
	LacpSysGlobalInfoDestroy(LaSystemActor)
	LacpSysGlobalInfoDestroy(LaSystemPeer)

}
