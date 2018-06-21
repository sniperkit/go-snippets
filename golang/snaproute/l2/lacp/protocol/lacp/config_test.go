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
	"fmt"
	"l2/lacp/protocol/utils"
	"net"
	"testing"
	"time"
	"utils/fsm"

	"github.com/google/gopacket"
)

func ConfigSetup() {

	OnlyForTestSetup()
	utils.PortConfigMap[3] = utils.PortConfig{Name: "SIMeth1.1",
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
}

func ConfigTeardown() {
	OnlyForTestTeardown()
}

func InvalidStateCheck(p *LaAggPort, invalidStates []fsm.Event, prevState fsm.State, currState fsm.State) (string, bool) {

	var s string
	rc := true

	portchan := p.PortChannelGet()

	// force what State transition should have been
	p.RxMachineFsm.Machine.Curr.SetState(prevState)
	p.RxMachineFsm.Machine.Curr.SetState(currState)

	for _, e := range invalidStates {
		// send PORT MOVED event to Rx Machine
		p.RxMachineFsm.RxmEvents <- utils.MachineEvent{
			E:            e,
			ResponseChan: portchan,
			Src:          "TEST"}

		// wait for response
		if msg := <-portchan; msg != RxMachineModuleStr {
			s = fmt.Sprintf("Expected response from", RxMachineModuleStr)
			rc = false
			return s, rc
		}

		// PORT MOVED
		if p.RxMachineFsm.Machine.Curr.PreviousState() != prevState &&
			p.RxMachineFsm.Machine.Curr.CurrentState() != currState {
			s = fmt.Sprintf("ERROR RX Machine State incorrect expected (prev/curr)",
				prevState,
				currState,
				"actual",
				p.RxMachineFsm.Machine.Curr.PreviousState(),
				p.RxMachineFsm.Machine.Curr.CurrentState())
			rc = false
			return s, rc
		}
	}

	return "", rc
}

func TestLaAggPortCreateAndBeginEvent(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort

	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:     3,
		Prio:   0x80,
		Key:    100,
		AggId:  2000,
		Enable: false,
		Mode:   LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {

		//	fmt.Println("Rx:", p.RxMachineFsm.Machine.Curr.CurrentState(),
		//		"Ptx:", p.PtxMachineFsm.Machine.Curr.CurrentState(),
		//		"Cd:", p.CdMachineFsm.Machine.Curr.CurrentState(),
		//		"Mux:", p.MuxMachineFsm.Machine.Curr.CurrentState(),
		//		"Tx:", p.TxMachineFsm.Machine.Curr.CurrentState())

		// lets test the States, after initialization port moves to Disabled State
		// Rx Machine
		if p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
			t.Error("ERROR RX Machine State incorrect expected",
				LacpRxmStatePortDisabled, "actual",
				p.RxMachineFsm.Machine.Curr.CurrentState())
		}
		// Periodic Tx Machine
		if p.PtxMachineFsm.Machine.Curr.CurrentState() != LacpPtxmStateNoPeriodic {
			t.Error("ERROR PTX Machine State incorrect expected",
				LacpPtxmStateNoPeriodic, "actual",
				p.PtxMachineFsm.Machine.Curr.CurrentState())
		}
		// Churn Detection Machine
		if p.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("ERROR CD Machine State incorrect expected",
				LacpCdmStateActorChurnMonitor, "actual",
				p.CdMachineFsm.Machine.Curr.CurrentState())
		}
		// Mux Machine
		if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached {
			t.Error("ERROR MUX Machine State incorrect expected",
				LacpMuxmStateDetached, "actual",
				p.MuxMachineFsm.Machine.Curr.CurrentState())
		}
		// Tx Machine
		if p.TxMachineFsm.Machine.Curr.CurrentState() != LacpTxmStateOff {
			t.Error("ERROR TX Machine State incorrect expected",
				LacpTxmStateOff, "actual",
				p.TxMachineFsm.Machine.Curr.CurrentState())
		}
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

	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

func TestLaAggPortCreateDifferentModes(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort

	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	modeList := []int{LacpModeOn, LacpModeActive, LacpModePassive}

	for _, mode := range modeList {
		pconf := &LaAggPortConfig{
			Id:     3,
			Prio:   0x80,
			Key:    100,
			AggId:  2000,
			Enable: false,
			Mode:   mode,
			Properties: PortProperties{
				Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
				Speed:  1000000000,
				Duplex: LacpPortDuplexFull,
				Mtu:    1500,
			},
			IntfId:   "SIMeth1.1",
			TraceEna: false,
		}

		// lets create a port and start the machines
		CreateLaAggPort(pconf)

		// if the port is found verify the initial State after begin event
		// which was called as part of create
		if LaFindPortById(pconf.Id, &p) {

			//	fmt.Println("Rx:", p.RxMachineFsm.Machine.Curr.CurrentState(),
			//		"Ptx:", p.PtxMachineFsm.Machine.Curr.CurrentState(),
			//		"Cd:", p.CdMachineFsm.Machine.Curr.CurrentState(),
			//		"Mux:", p.MuxMachineFsm.Machine.Curr.CurrentState(),
			//		"Tx:", p.TxMachineFsm.Machine.Curr.CurrentState())

			// lets test the States, after initialization port moves to Disabled State
			// Rx Machine
			if p.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStatePortDisabled {
				t.Error("ERROR RX Machine State incorrect expected",
					LacpRxmStatePortDisabled, "actual",
					p.RxMachineFsm.Machine.Curr.CurrentState())
			}
			// Periodic Tx Machine
			if p.PtxMachineFsm.Machine.Curr.CurrentState() != LacpPtxmStateNoPeriodic {
				t.Error("ERROR PTX Machine State incorrect expected",
					LacpPtxmStateNoPeriodic, "actual",
					p.PtxMachineFsm.Machine.Curr.CurrentState())
			}
			// Churn Detection Machine
			if p.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
				t.Error("ERROR CD Machine State incorrect expected",
					LacpCdmStateActorChurnMonitor, "actual",
					p.CdMachineFsm.Machine.Curr.CurrentState())
			}
			// Mux Machine
			if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached {
				t.Error("ERROR MUX Machine State incorrect expected",
					LacpMuxmStateDetached, "actual",
					p.MuxMachineFsm.Machine.Curr.CurrentState())
			}
			// Tx Machine
			if p.TxMachineFsm.Machine.Curr.CurrentState() != LacpTxmStateOff {
				t.Error("ERROR TX Machine State incorrect expected",
					LacpTxmStateOff, "actual",
					p.TxMachineFsm.Machine.Curr.CurrentState())
			}
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
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

func TestLaAggPortCreateWithInvalidKeySetWithAgg(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort

	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}

	LacpSysGlobalInfoInit(sysId)

	aconf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
		Id:   2000,
		Key:  50,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:01:02:03:04:05",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(aconf)

	pconf := &LaAggPortConfig{
		Id:     3,
		Prio:   0x80,
		Key:    100, // INVALID
		AggId:  2000,
		Enable: true,
		Mode:   LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x02, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {
		if p.aggSelected == LacpAggSelected {
			t.Error("Port is in SELECTED mode")
		}
	}

	// Delete the port and agg
	DeleteLaAggPort(pconf.Id)
	DeleteLaAgg(aconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.SysKey, sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.SysKey, sgi.PortList, sgi.PortMap)
		}
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

func TestLaAggPortCreateWithoutKeySetNoAgg(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort
	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:     3,
		Prio:   0x80,
		Key:    100,
		AggId:  2000,
		Enable: true,
		Mode:   LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {
		if p.aggSelected == LacpAggSelected {
			t.Error("Port is in SELECTED mode")
		}
	}

	// Delete port
	DeleteLaAggPort(pconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

func TestLaAggPortCreateThenCorrectAggCreate(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort
	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:     3,
		Prio:   0x80,
		Key:    100,
		AggId:  2000,
		Enable: true,
		Mode:   LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {
		if p.aggSelected == LacpAggSelected {
			t.Error("Port is in SELECTED mode should be UNSELECTED")
		}
	}

	aconf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
		Id:   2000,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:01:02:03:04:05",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(aconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if p.aggSelected != LacpAggSelected {
		t.Error("Port is in SELECTED mode (2)")
	}

	if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateAttached {
		t.Error("Mux State expected", MuxmStateStrMap[LacpMuxmStateAttached], "actual", MuxmStateStrMap[p.MuxMachineFsm.Machine.Curr.CurrentState()])
	}

	// TODO Check States of other State machines

	// Delete agg
	DeleteLaAggPort(pconf.Id)
	DeleteLaAgg(aconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

// TestLaAggPortCreateThenCorrectAggCreateThenDetach:
// - create port
// - create lag
// - attach port
// - enable port
func TestLaAggPortCreateThenCorrectAggCreateThenDetach(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort
	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:    3,
		Prio:  0x80,
		Key:   100,
		AggId: 2000,
		Mode:  LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {
		if p.aggSelected == LacpAggSelected {
			t.Error("Port is in SELECTED mode 1")
		}
	}

	aconf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
		Id:   2000,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:01:02:03:04:05",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(aconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create should be disabled since this
	// is the initial state of port config
	if p.aggSelected == LacpAggSelected {
		t.Error("Port is in SELECTED mode 2 mux state", p.MuxMachineFsm.Machine.Curr.CurrentState())
	}

	EnableLaAggPort(pconf.Id)

	if p.aggSelected != LacpAggSelected {
		t.Error("Port is NOT in SELECTED mode 3")
	}

	if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateAttached {
		t.Error("Mux State expected", LacpMuxmStateAttached, "actual", p.MuxMachineFsm.Machine.Curr.CurrentState())
	}
	// Delete port
	DeleteLaAggPortFromAgg(pconf.Key, pconf.Id)
	DeleteLaAggPort(pconf.Id)
	DeleteLaAgg(aconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

// Enable port post creation
func TestLaAggPortEnable(t *testing.T) {
	defer MemoryCheck(t)
	var p *LaAggPort
	ConfigSetup()
	// must be called to initialize the global
	sysId := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}}
	LacpSysGlobalInfoInit(sysId)

	pconf := &LaAggPortConfig{
		Id:    3,
		Prio:  0x80,
		Key:   100,
		AggId: 2000,
		Mode:  LacpModeActive,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   "SIMeth1.1",
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(pconf)

	// if the port is found verify the initial State after begin event
	// which was called as part of create
	if LaFindPortById(pconf.Id, &p) {
		if p.aggSelected == LacpAggSelected {
			t.Error("Port is in SELECTED mode")
		}
	}

	aconf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
		Id:   2000,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpSlowPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:01:02:03:04:05",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(aconf)

	if p.aggSelected == LacpAggSelected {
		t.Error("Port is in SELECTED mode")
	}

	EnableLaAggPort(pconf.Id)

	if p.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateAttached {
		t.Error("Mux State expected", LacpMuxmStateAttached, "actual", p.MuxMachineFsm.Machine.Curr.CurrentState())
	}

	if p.aggSelected != LacpAggSelected {
		t.Error("Port is in SELECTED mode")
	}

	// Delete port
	DeleteLaAggPort(pconf.Id)

	DeleteLaAgg(aconf.Id)
	for _, sgi := range LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}
	ConfigTeardown()
	LacpSysGlobalInfoDestroy(sysId)
}

func TestTwoAggsBackToBackSinglePort(t *testing.T) {
	defer MemoryCheck(t)
	const LaAggPortActor = 10
	const LaAggPortPeer = 20
	LaAggPortActorIf := "SIMeth0"
	LaAggPortPeerIf := "SIM2eth0"
	OnlyForTestSetup()
	utils.PortConfigMap[LaAggPortActor] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPortPeer] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x44, 0x44, 0x22, 0x22, 0x33},
	}

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
		TraceEna: false,
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
		TraceEna: false,
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
		close(testWait)

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

		if !p1.AggAttached.OperState {
			t.Error(fmt.Sprintf("OperState not UP as expected"))
		}
		if !p2.AggAttached.OperState {
			t.Error(fmt.Sprintf("OperState not UP as expected"))
		}

		// TODO check the States of the other State machines
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

// TestTwoAggsBackToBackSinglePortTimeout will allow for
// two ports to sync up then force a timeout by disabling
// one end of the connection by setting the mode to "ON"
func TestTwoAggsBackToBackSinglePortTimeout(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	const LaAggPortActor = 11
	const LaAggPortPeer = 21
	const LaAggPortActorIf = "SIMeth0"
	const LaAggPortPeerIf = "SIMeth1"

	utils.PortConfigMap[LaAggPortActor] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPortPeer] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x44, 0x44, 0x22, 0x22, 0x33},
	}

	// must be called to initialize the global
	LaSystemActor := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x64}}
	LaSystemPeer := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	bridge := SimulationBridge{
		Port1:       LaAggPortActor,
		Port2:       LaAggPortPeer,
		RxLacpPort1: make(chan gopacket.Packet),
		RxLacpPort2: make(chan gopacket.Packet),
	}

	ActorSystem := LacpSysGlobalInfoInit(LaSystemActor)
	PeerSystem := LacpSysGlobalInfoInit(LaSystemPeer)
	ActorSystem.LaSysGlobalRegisterTxCallback(LaAggPortActorIf, bridge.TxViaGoChannel)
	PeerSystem.LaSysGlobalRegisterTxCallback(LaAggPortPeerIf, bridge.TxViaGoChannel)

	// port 1
	go LaRxMain(bridge.Port1, bridge.RxLacpPort1)
	// port 2
	go LaRxMain(bridge.Port2, bridge.RxLacpPort2)

	p1conf := &LaAggPortConfig{
		Id:      LaAggPortActor,
		Prio:    0x80,
		Key:     100,
		AggId:   100,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortActorIf,
		TraceEna: false,
	}

	p2conf := &LaAggPortConfig{
		Id:      LaAggPortPeer,
		Prio:    0x80,
		Key:     200,
		AggId:   200,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortPeer, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortPeerIf,
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(p1conf)
	CreateLaAggPort(p2conf)

	a1conf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x01, 0x01, 0x01},
		Id:   100,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpFastPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:64",
			SystemPriority: 128},
	}

	a2conf := &LaAggConfig{
		Name: "agg2",
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x02, 0x02, 0x02},
		Id:   200,
		Key:  200,
		Lacp: LacpConfigInfo{Interval: LacpFastPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:C8",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(a1conf)
	CreateLaAgg(a2conf)

	// Add port to agg
	AddLaAggPortToAgg(a1conf.Key, p1conf.Id)
	AddLaAggPortToAgg(a2conf.Key, p2conf.Id)

	//time.Sleep(time.Second * 30)
	testWait := make(chan bool)

	var p1 *LaAggPort
	var p2 *LaAggPort
	if LaFindPortById(p1conf.Id, &p1) &&
		LaFindPortById(p2conf.Id, &p2) {

		go func() {
			for i := 0; i < 5 &&
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

		// TODO check the States of the other State machines

		// Lets disable lacp for p1
		SetLaAggPortLacpMode(p1conf.Id, LacpModeOn)

		go func() {
			var i int
			// current while timer must expire twice to get  to expected states fo
			for i = 0; i < 5 &&
				(p1.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted ||
					p1.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached ||
					p2.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached); i++ {

				time.Sleep(time.Second * 1)
			}

			if i == 10 {
				testWait <- false
			} else {
				testWait <- true
			}

		}()

		testResult := <-testWait
		if !testResult {
			t.Error(fmt.Sprintln("Actor and Peer States are not correct Expected P1 RXM/MUX",
				LacpRxmStateLacpDisabled, LacpMuxmStateDistributing, "Actual", p1.RxMachineFsm.Machine.Curr.CurrentState(),
				p1.MuxMachineFsm.Machine.Curr.CurrentState(), "Expected P2 RXM/MUX", LacpRxmStateDefaulted,
				LacpMuxmStateDetached, "Actual", p2.RxMachineFsm.Machine.Curr.CurrentState(),
				p2.MuxMachineFsm.Machine.Curr.CurrentState()))
		}

		// TODO check the States of the other State machines
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

// TestLaAggCallSaveLaAggConfig No logic just coverage
func TestLaAggCallSaveLaAggConfig(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	aconf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x02, 0x02, 0x02},
		Id:   200,
		Key:  200,
		Lacp: LacpConfigInfo{Interval: LacpFastPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:C8",
			SystemPriority: 128},
		LagMembers: []uint16{1, 2, 3, 4},
	}

	// Call method when agg does not exist
	SaveLaAggConfig(aconf)

	// Create Port Config
	CreateLaAgg(aconf)

	// Lets actaully save off the config
	SaveLaAggConfig(aconf)

	// delete the lag
	DeleteLaAgg(aconf.Id)
	OnlyForTestTeardown()

}

// TestTwoAggsBackToBackSingleDisableEnableLaAgg will allow for
// two ports to sync up then disable one end of the lag
func TestTwoAggsBackToBackSingleDisableEnableLaAgg(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	const LaAggPortActor = 11
	const LaAggPortPeer = 21
	const LaAggPortActorIf = "SIMeth0"
	const LaAggPortPeerIf = "SIMeth1"

	utils.PortConfigMap[LaAggPortActor] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPortPeer] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	// must be called to initialize the global
	LaSystemActor := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x64}}
	LaSystemPeer := LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	bridge := SimulationBridge{
		Port1:       LaAggPortActor,
		Port2:       LaAggPortPeer,
		RxLacpPort1: make(chan gopacket.Packet),
		RxLacpPort2: make(chan gopacket.Packet),
	}

	ActorSystem := LacpSysGlobalInfoInit(LaSystemActor)
	PeerSystem := LacpSysGlobalInfoInit(LaSystemPeer)
	ActorSystem.LaSysGlobalRegisterTxCallback(LaAggPortActorIf, bridge.TxViaGoChannel)
	PeerSystem.LaSysGlobalRegisterTxCallback(LaAggPortPeerIf, bridge.TxViaGoChannel)

	// port 1
	go LaRxMain(bridge.Port1, bridge.RxLacpPort1)
	// port 2
	go LaRxMain(bridge.Port2, bridge.RxLacpPort2)

	p1conf := &LaAggPortConfig{
		Id:      LaAggPortActor,
		Prio:    0x80,
		Key:     100,
		AggId:   100,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortActorIf,
		TraceEna: true,
	}

	p2conf := &LaAggPortConfig{
		Id:      LaAggPortPeer,
		Prio:    0x80,
		Key:     200,
		AggId:   200,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPortPeer, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortPeerIf,
		TraceEna: false,
	}

	// lets create a port and start the machines
	CreateLaAggPort(p1conf)
	CreateLaAggPort(p2conf)

	a1conf := &LaAggConfig{
		Name: "agg1",
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x01, 0x01, 0x01},
		Id:   100,
		Key:  100,
		Lacp: LacpConfigInfo{Interval: LacpFastPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:64",
			SystemPriority: 128},
	}

	a2conf := &LaAggConfig{
		Name: "agg2",
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x02, 0x02, 0x02},
		Id:   200,
		Key:  200,
		Lacp: LacpConfigInfo{Interval: LacpFastPeriodicTime,
			Mode:           LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:C8",
			SystemPriority: 128},
	}

	// Create Aggregation
	CreateLaAgg(a1conf)
	CreateLaAgg(a2conf)

	// Add port to agg
	AddLaAggPortToAgg(a1conf.Key, p1conf.Id)
	AddLaAggPortToAgg(a2conf.Key, p2conf.Id)

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

		// TODO check the States of the other State machines

		// Lets disable the first lag
		DisableLaAgg(a1conf.Id)

		go func() {
			var i int
			// current while timer must expire twice to get  to expected states fo
			for i = 0; i < 20 &&
				(p1.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted ||
					p1.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached ||
					p2.RxMachineFsm.Machine.Curr.CurrentState() != LacpRxmStateDefaulted ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDetached); i++ {

				time.Sleep(time.Second * 1)
			}

			if i == 10 {
				testWait <- false
			} else {
				testWait <- true
			}

		}()

		testResult := <-testWait
		if !testResult {
			t.Error(fmt.Sprintln("Actor and Peer States are not correct Expected P1 RXM/MUX",
				LacpRxmStateDefaulted, LacpMuxmStateDistributing, "Actual", p1.RxMachineFsm.Machine.Curr.CurrentState(),
				p1.MuxMachineFsm.Machine.Curr.CurrentState(), "Expected P2 RXM/MUX", LacpRxmStateDefaulted,
				LacpMuxmStateDetached, "Actual", p2.RxMachineFsm.Machine.Curr.CurrentState(),
				p2.MuxMachineFsm.Machine.Curr.CurrentState()))
		}

		if p1.CdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStateActorChurnMonitor {
			t.Error("ERROR: Error Actor Churn Detection Machine did not transition properly")
		}
		if p1.PCdMachineFsm.Machine.Curr.CurrentState() != LacpCdmStatePartnerChurnMonitor {
			t.Error("ERROR: Error Partner Churn Detection Machine did not transition properly")
		}
		if p1.TxMachineFsm.Machine.Curr.CurrentState() != LacpTxmStateOff {
			t.Error("ERROR: Error Transmit machine was not disabled")
		}

		// Lets re-enable the first lag
		EnableLaAgg(a1conf.Id)

		go func() {
			for i := 0; i < 35 &&
				(p1.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDistributing ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != LacpMuxmStateDistributing); i++ {
				time.Sleep(time.Second * 1)
			}
			testWait <- true
		}()

		<-testWait

		State1 = GetLaAggPortActorOperState(p1conf.Id)
		State2 = GetLaAggPortActorOperState(p2conf.Id)

		if !LacpStateIsSet(State1, portUpState) {
			t.Error(fmt.Sprintf("Actor Port State 0x%x did not come up properly with peer expected 0x%x", State1, portUpState))
		}
		if !LacpStateIsSet(State2, portUpState) {
			t.Error(fmt.Sprintf("Peer Port State 0x%x did not come up properly with actor expected 0x%x", State2, portUpState))
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

func TestTwoAggsBackToBackSinglePortValidLacpModeCombo(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	const LaAggPortActor = 10
	const LaAggPortPeer = 20
	LaAggPortActorIf := "SIMeth0"
	LaAggPortPeerIf := "SIM2eth0"
	utils.PortConfigMap[LaAggPortActor] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPortPeer] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

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

	var comboList = []int{LacpModeOn, LacpModeActive, LacpModePassive}

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

		// active state
		const portUpState = LacpStateActivityBit | LacpStateAggregationBit |
			LacpStateSyncBit | LacpStateCollectingBit | LacpStateDistributingBit

		if !LacpStateIsSet(State1, portUpState) {
			t.Error(fmt.Sprintf("Actor Port State %s did not come up properly with peer expected %s", LacpStateToStr(State1), LacpStateToStr(portUpState)))
		}
		if !LacpStateIsSet(State2, portUpState) {
			t.Error(fmt.Sprintf("Peer Port State %s did not come up properly with actor expected %s", LacpStateToStr(State2), LacpStateToStr(portUpState)))
		}

		for _, p1mode := range comboList {
			for _, p2mode := range comboList {
				if (p1mode != LacpModeOn && p2mode != LacpModeOn) ||
					(p1mode == LacpModeOn && p2mode == LacpModeOn) {

					SetLaAggPortLacpMode(p1conf.Id, p1mode)
					SetLaAggPortLacpMode(p2conf.Id, p2mode)

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

					var port1UpState = uint8(LacpStateActivityBit | LacpStateAggregationBit |
						LacpStateSyncBit | LacpStateCollectingBit | LacpStateDistributingBit)

					var port2UpState = uint8(LacpStateActivityBit | LacpStateAggregationBit |
						LacpStateSyncBit | LacpStateCollectingBit | LacpStateDistributingBit)

					if p1mode == LacpModeOn {
						port1UpState |= LacpStateDefaultedBit
					} else {
						port1UpState &= ^(uint8(LacpStateDefaultedBit))
					}
					if p1mode == LacpModePassive {
						port1UpState &= ^(uint8(LacpStateActivityBit))
					}

					if p2mode == LacpModeOn {
						port2UpState |= LacpStateDefaultedBit
					} else {
						port2UpState &= ^(uint8(LacpStateDefaultedBit))
					}
					if p2mode == LacpModePassive {
						port2UpState &= ^(uint8(LacpStateActivityBit))
					}

					if !LacpStateIsSet(State1, port1UpState) {
						t.Error(fmt.Sprintf("Actor Port State %s did not come up properly with peer expected %s p1mode %d p2mode %d", LacpStateToStr(State1), LacpStateToStr(port1UpState), p1mode, p2mode))
					}
					if !LacpStateIsSet(State2, port2UpState) {
						t.Error(fmt.Sprintf("Peer Port State %s did not come up properly with actor expected %s p1mode %d p2mode %d", LacpStateToStr(State2), LacpStateToStr(port2UpState), p1mode, p2mode))
					}
				}
			}
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
	}
	OnlyForTestTeardown()
	LacpSysGlobalInfoDestroy(LaSystemActor)
	LacpSysGlobalInfoDestroy(LaSystemPeer)

}

func TestSetLaAggPortSystemInfo(t *testing.T) {
	defer MemoryCheck(t)
	OnlyForTestSetup()
	p1conf := &LaAggPortConfig{
		Id:      LaAggChurnPortActor,
		Prio:    0x80,
		Key:     100,
		AggId:   LaAggChurnAgg1,
		Enable:  true,
		Mode:    LacpModeActive,
		Timeout: LacpShortTimeoutTime,
		Properties: PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggChurnPortActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggChurnPortActorIf,
		TraceEna: true,
	}

	CreateLaAggPort(p1conf)

	// Send bad mac string
	SetLaAggPortSystemInfo(p1conf.Id, "11:11:11", 100)

	// Send invalid port
	SetLaAggPortSystemInfo(5, "11:11:11", 100)

	// Send valid port and valid info
	SetLaAggPortSystemInfo(p1conf.Id, "11:11:11:11:11:11", 100)

	DeleteLaAggPort(p1conf.Id)
	OnlyForTestTeardown()
}

func TestTwoAggsBackToBackSinglePortDisablePort(t *testing.T) {
	defer MemoryCheck(t)
	const LaAggPortActor = 10
	const LaAggPortPeer = 20
	LaAggPortActorIf := "SIMeth0"
	LaAggPortPeerIf := "SIM2eth0"
	OnlyForTestSetup()
	utils.PortConfigMap[LaAggPortActor] = utils.PortConfig{Name: LaAggPortActorIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPortPeer] = utils.PortConfig{Name: LaAggPortPeerIf,
		HardwareAddr: net.HardwareAddr{0x00, 0x44, 0x44, 0x22, 0x22, 0x33},
	}

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
		TraceEna: false,
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
		TraceEna: false,
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
		close(testWait)

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

		if !p1.AggAttached.OperState {
			t.Error(fmt.Sprintf("OperState not UP as expected"))
		}
		if !p2.AggAttached.OperState {
			t.Error(fmt.Sprintf("OperState not UP as expected"))
		}

		// TODO check the States of the other State machines
	} else {
		t.Error("Unable to find port just created")
	}

	DisableLaAggPort(p1.PortNum)
	testWait = make(chan bool)

	go func() {
		for i := 0; i < 10 &&
			(p1.MuxMachineFsm.Machine.Curr.CurrentState() == LacpMuxmStateDistributing); i++ {
			time.Sleep(time.Second * 1)
		}
		testWait <- true
	}()

	<-testWait

	if p1.AggAttached.OperState {
		t.Error("Error OperState is set when it should not be as the port is down")
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
