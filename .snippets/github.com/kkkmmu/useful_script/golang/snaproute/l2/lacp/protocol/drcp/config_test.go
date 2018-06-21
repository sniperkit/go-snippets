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
// config_test.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"net"
	//"sort"
	"testing"
	"time"
	asicdmock "utils/asicdClient/mock"
	"utils/commonDefs"
	"utils/logging"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// first two bytes are priority but in case of ipp this is the neighbor system number
const ipplink1 int32 = 3
const aggport1 int32 = 1
const aggport2 int32 = 2

// first two bytes are priority but in case of ipp this is the neighbor system number
const ipplink2 int32 = 4
const aggport3 int32 = 5
const aggport4 int32 = 6

const LaAggPort1NeighborActor = 10
const LaAggPort2NeighborActor = 30
const LaAggPort1Peer = 20
const LaAggPort2Peer = 21
const DRNeighborIpp1 = 11
const DRNeighborIpp2 = 31
const LaAggPortNeighborActor1If = "SIMeth0"
const LaAggPortNeighborActor2If = "SIMeth1"
const LaAggPortPeerIf1 = "SIM2eth0"
const LaAggPortPeerIf2 = "SIM2eth1"
const DRNeighborIppIf1 = "SIMIPPeth2"
const DRNeighborIppIf2 = "SIMIPPeth3"

type MyTestMock struct {
	asicdmock.MockAsicdClientMgr
}

func (m *MyTestMock) GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error) {

	getinfo := &commonDefs.VlanGetInfo{
		StartIdx: 1,
		EndIdx:   1,
		Count:    1,
		More:     false,
		VlanList: make([]commonDefs.Vlan, 1),
	}

	getinfo.VlanList[0] = commonDefs.Vlan{
		VlanId:      100,
		IfIndexList: []int32{aggport1, aggport2, LaAggPort1NeighborActor, LaAggPort2NeighborActor, LaAggPort1Peer, LaAggPort2Peer},
	}
	return getinfo, nil
}

func (m *MyTestMock) GetPortLinkStatus(port int32) bool {
	return true
}

var testBlockMap map[string][]string = make(map[string][]string, 0)

func (m *MyTestMock) IppIngressEgressDrop(inport, aggport string) error {

	if _, ok := testBlockMap[inport]; ok {
		testBlockMap[inport] = append(testBlockMap[inport], aggport)
	}
	return nil
}

func (m *MyTestMock) IppIngressEgressPass(inport, aggport string) error {

	if _, ok := testBlockMap[inport]; ok {
		for i, p := range testBlockMap[inport] {
			if p == aggport {
				testBlockMap[inport] = append(testBlockMap[inport][:i], testBlockMap[inport][i+1:]...)
				if len(testBlockMap[inport]) == 0 {
					delete(testBlockMap, inport)
				}
			}
		}
	}
	return nil
}

func OnlyForTestSetup() {
	logger, _ := logging.NewLogger("lacpd", "TEST", false)
	utils.SetLaLogger(logger)
	utils.DeleteAllAsicDPlugins()
	utils.SetAsicDPlugin(&MyTestMock{})
	// fill in conversations
	GetAllCVIDConversations()
}

func OnlyForTestTeardown(t *testing.T) {

	utils.SetLaLogger(nil)
	utils.DeleteAllAsicDPlugins()
	ConversationIdMap[100].Valid = false
	ConversationIdMap[100].PortList = nil
	ConversationIdMap[100].Cvlan = 0
	ConversationIdMap[100].Refcnt = 0
	ConversationIdMap[100].Idtype = [4]uint8{}

	ConversationIdMap[200].Valid = false
	ConversationIdMap[200].PortList = nil
	ConversationIdMap[200].Cvlan = 0
	ConversationIdMap[200].Refcnt = 0
	ConversationIdMap[200].Idtype = [4]uint8{}

	lacp.ConfigAggMap = nil
	lacp.ConfigAggList = nil

	// validate that the
	if len(DistributedRelayDB) != 0 {
		t.Error("Error DR objects not deleted")
	}
	if len(DistributedRelayDBList) != 0 {
		t.Error("Error DR objects not deleted")
	}
	if len(DRCPIppDB) != 0 {
		t.Error("Error IPP objects not deleted")
	}
	if len(DRCPIppDBList) != 0 {
		t.Error("Error IPP objects not deleted")
	}
}

func OnlyForTestSetupCreateAggGroup(aggId uint32) *lacp.LaAggregator {
	a1conf := &lacp.LaAggConfig{
		Name: fmt.Sprintf("agg%d", aggId),
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x01, 0x01, 0x01},
		Id:   int(aggId),
		Key:  uint16(aggId),
		Lacp: lacp.LacpConfigInfo{Interval: lacp.LacpFastPeriodicTime,
			Mode:           lacp.LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:64",
			SystemPriority: 128},
	}
	lacp.CreateLaAgg(a1conf)

	p1conf := &lacp.LaAggPortConfig{
		Id:      uint16(aggport1),
		Prio:    0x80,
		Key:     uint16(aggId),
		AggId:   int(aggId),
		Enable:  true,
		Mode:    lacp.LacpModeActive,
		Timeout: lacp.LacpShortTimeoutTime,
		Properties: lacp.PortProperties{
			Mac:    net.HardwareAddr{0x00, byte(aggport1), 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: lacp.LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   utils.PortConfigMap[aggport1].Name,
		TraceEna: false,
	}

	lacp.CreateLaAggPort(p1conf)
	lacp.AddLaAggPortToAgg(a1conf.Key, p1conf.Id)

	var a *lacp.LaAggregator
	if lacp.LaFindAggById(a1conf.Id, &a) {
		a.DistributedPortNumList = append(a.DistributedPortNumList, utils.PortConfigMap[aggport1].Name)

		return a
	}
	return nil
}

func ConfigTestSetup() {
	OnlyForTestSetup()
	utils.PortConfigMap[ipplink1] = utils.PortConfig{Name: "SIMeth1.3",
		HardwareAddr: net.HardwareAddr{0x00, 0x33, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[aggport1] = utils.PortConfig{Name: "SIMeth1.1",
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[aggport2] = utils.PortConfig{Name: "SIMeth1.2",
		HardwareAddr: net.HardwareAddr{0x00, 0x22, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[ipplink2] = utils.PortConfig{Name: "SIMeth0.3",
		HardwareAddr: net.HardwareAddr{0x00, 0x44, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[aggport3] = utils.PortConfig{Name: "SIMeth0.1",
		HardwareAddr: net.HardwareAddr{0x00, 0x55, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[aggport4] = utils.PortConfig{Name: "SIMeth0.2",
		HardwareAddr: net.HardwareAddr{0x00, 0x66, 0x11, 0x22, 0x22, 0x33},
	}
}
func ConfigTestTeardwon(t *testing.T) {

	OnlyForTestTeardown(t)
	delete(utils.PortConfigMap, ipplink1)
	delete(utils.PortConfigMap, aggport1)
	delete(utils.PortConfigMap, aggport2)
	delete(utils.PortConfigMap, ipplink2)
	delete(utils.PortConfigMap, aggport3)
	delete(utils.PortConfigMap, aggport4)
}

func TestConfigDistributedRelayValidCreateAggWithPortsThenCreateDR(t *testing.T) {

	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    uint32(a.AggId),
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}
	// map vlan 100 to this system
	// in real system this should be filled in by vlan membership
	cfg.DrniConvAdminGateway[100][0] = cfg.DrniPortalSystemNumber

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}

	CreateDistributedRelay(cfg)

	// configuration is incomplete because lag has not been created as part of this
	if len(DistributedRelayDB) == 0 ||
		len(DistributedRelayDBList) == 0 {
		t.Error("ERROR Distributed Relay Object was not added to global DB's")
	}
	dr, ok := DistributedRelayDB[cfg.DrniName]
	if !ok {
		t.Error("ERROR Distributed Relay Object was not found in global DB's")
	}

	// check the inital state of each of the state machines
	if dr.a == nil {
		t.Error("ERROR BEGIN was called before an Agg has been attached")
	}

	waitChan := make(chan bool)
	go func(wc *chan bool) {
		for i := 0; i < 10 && (dr.PsMachineFsm.Machine.Curr.CurrentState() != PsmStatePortalSystemUpdate); i++ {
			time.Sleep(time.Millisecond * 10)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	// Rx machine sets the change portal which should inform the psm
	if dr.PsMachineFsm == nil ||
		dr.PsMachineFsm.Machine.Curr.CurrentState() != GmStateDRNIGatewayUpdate {
		t.Error("ERROR BEGIN Initial Portal System Machine state is not correct", PsmStateStrMap[dr.PsMachineFsm.Machine.Curr.CurrentState()])
	}

	go func(wc *chan bool) {
		for i := 0; i < 10 && dr.GMachineFsm.Machine.Curr.CurrentState() != GmStateDRNIGatewayUpdate; i++ {
			time.Sleep(time.Millisecond * 10)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	// Ps Machine updates the Gm based on the rx gateway update
	if dr.GMachineFsm == nil ||
		dr.GMachineFsm.Machine.Curr.CurrentState() != GmStateDRNIGatewayUpdate {
		t.Error("ERROR BEGIN Initial Gateway Machine state is not correct", GmStateStrMap[dr.GMachineFsm.Machine.Curr.CurrentState()])
	}

	go func(wc *chan bool) {
		for i := 0; i < 4 && dr.AMachineFsm.Machine.Curr.CurrentState() != AmStateDRNIPortUpdate; i++ {
			time.Sleep(time.Second * 1)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	if dr.AMachineFsm == nil ||
		dr.AMachineFsm.Machine.Curr.CurrentState() != AmStateDRNIPortUpdate {
		t.Error("ERROR BEGIN Initial Aggregator System Machine state is not correct", AmStateStrMap[dr.AMachineFsm.Machine.Curr.CurrentState()])
	}

	if len(dr.Ipplinks) == 0 {
		t.Error("ERROR Why did the IPL IPP link not get created")
	}

	for _, ipp := range dr.Ipplinks {
		// IPL should be disabled thus state should be in initialized state
		if ipp.RxMachineFsm == nil ||
			ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateExpired {
			t.Error("ERROR BEGIN Initial Receive Machine state is not correct", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
		}
		// port is enabled and drcp is enabled thus we should be in fast periodic state
		if ipp.PtxMachineFsm == nil {
			t.Error("ERROR BEGIN Initial Periodic Tx Machine state is not correct")
		}
		if ipp.TxMachineFsm == nil ||
			ipp.TxMachineFsm.Machine.Curr.CurrentState() != TxmStateOff {
			t.Error("ERROR BEGIN Initial Tx Machine state is not correct", TxmStateStrMap[ipp.TxMachineFsm.Machine.Curr.CurrentState()])
		}
		if ipp.NetIplShareMachineFsm == nil ||
			ipp.NetIplShareMachineFsm.Machine.Curr.CurrentState() != NetIplSharemStateNoManipulatedFramesSent {
			t.Error("ERROR BEGIN Initial Net/IPL Sharing Machine state is not correct", ipp.NetIplShareMachineFsm.Machine.Curr.CurrentState())
		}
		if ipp.IAMachineFsm == nil ||
			ipp.IAMachineFsm.Machine.Curr.CurrentState() != IAmStateIPPPortInitialize {
			t.Error("ERROR BEGIN Initial IPP Aggregator state is not correct", IAmStateStrMap[ipp.IAMachineFsm.Machine.Curr.CurrentState()])
		}

		go func(wc *chan bool) {
			for i := 0; i < 10 && ipp.IGMachineFsm.Machine.Curr.CurrentState() != IGmStateIPPGatewayUpdate; i++ {
				time.Sleep(time.Millisecond * 10)
			}
			*wc <- true
		}(&waitChan)

		<-waitChan
		/*
			TODO when gateway sync is fixed uncomment this
			if ipp.IGMachineFsm == nil ||
				ipp.IGMachineFsm.Machine.Curr.CurrentState() != IGmStateIPPGatewayUpdate {
				t.Error("ERROR BEGIN Initial IPP Gateway Machine state is not correct", GmStateStrMap[ipp.IGMachineFsm.Machine.Curr.CurrentState()])
			}

		*/
	}
	DeleteDistributedRelay(cfg.DrniName)

	if len(DistributedRelayDB) != 0 ||
		len(DistributedRelayDBList) != 0 {
		t.Error("ERROR Distributed Relay DB was not cleaned up")
	}
	if len(DRCPIppDB) != 0 ||
		len(DRCPIppDBList) != 0 {
		t.Error("ERROR IPP DB was not cleaned up")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigDistributedRelayCreateDRThenCreateAgg(t *testing.T) {

	ConfigTestSetup()

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    200,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}
	// map vlan 100 to this system
	// in real system this should be filled in by vlan membership
	cfg.DrniConvAdminGateway[100][0] = cfg.DrniPortalSystemNumber

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}

	CreateDistributedRelay(cfg)

	// configuration is incomplete because lag has not been created as part of this
	if len(DistributedRelayDB) == 0 ||
		len(DistributedRelayDBList) == 0 {
		t.Error("ERROR Distributed Relay Object was not added to global DB's")
	}
	dr, ok := DistributedRelayDB[cfg.DrniName]
	if !ok {
		t.Error("ERROR Distributed Relay Object was not found in global DB's")
	}

	// check the inital state of each of the state machines
	if dr.a != nil {
		t.Error("ERROR Agg does not exist no associated should exist")
	}
	if dr.PsMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Portal System Machine state was created, provisioning incomplete")
	}
	if dr.GMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Gateway Machine state was created, provisioning incomplete")
	}
	if dr.AMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Aggregator System Machine state was created, provisioning incomplete")
	}

	for _, ipp := range dr.Ipplinks {
		if ipp.DRCPEnabled {
			t.Error("ERROR Why is the IPL IPP link DRCP enabled")
		}
	}

	// TEST aggregator created after DR
	a := OnlyForTestSetupCreateAggGroup(200)

	// check the inital state of each of the state machines
	if dr.a == nil {
		t.Error("ERROR BEGIN was called before an Agg has been attached")
	}
	waitChan := make(chan bool)
	go func(wc *chan bool) {
		for i := 0; i < 20 && dr.PsMachineFsm.Machine.Curr.CurrentState() != PsmStatePortalSystemUpdate; i++ {
			time.Sleep(time.Millisecond * 10)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	// Rx machine sets the change portal which should inform the psm
	if dr.PsMachineFsm == nil ||
		dr.PsMachineFsm.Machine.Curr.CurrentState() != PsmStatePortalSystemUpdate {
		t.Error("ERROR BEGIN Initial Portal System Machine state is not correct", PsmStateStrMap[dr.PsMachineFsm.Machine.Curr.CurrentState()])
	}
	waitChan = make(chan bool)
	go func(wc *chan bool) {
		for i := 0; i < 20 && dr.GMachineFsm.Machine.Curr.CurrentState() != GmStateDRNIGatewayUpdate; i++ {
			time.Sleep(time.Millisecond * 10)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	// Ps Machine updates the Gm based on the rx gateway update
	if dr.GMachineFsm == nil ||
		dr.GMachineFsm.Machine.Curr.CurrentState() != GmStateDRNIGatewayUpdate {
		t.Error("ERROR BEGIN Initial Gateway Machine state is not correct", GmStateStrMap[dr.GMachineFsm.Machine.Curr.CurrentState()])
	}

	waitChan = make(chan bool)
	go func(wc *chan bool) {
		// takes time for agg port to be in distributing state so lets give it time
		for i := 0; i < 5 && dr.AMachineFsm.Machine.Curr.CurrentState() != AmStateDRNIPortUpdate; i++ {
			time.Sleep(time.Second * 1)
		}
		*wc <- true
	}(&waitChan)

	<-waitChan

	if dr.AMachineFsm == nil ||
		dr.AMachineFsm.Machine.Curr.CurrentState() != AmStateDRNIPortUpdate {
		t.Error("ERROR BEGIN Initial Aggregator System Machine state is not correct", AmStateStrMap[dr.AMachineFsm.Machine.Curr.CurrentState()])
	}

	if len(dr.Ipplinks) == 0 {
		t.Error("ERROR Why did the IPL IPP link not get created")
	}

	for _, ipp := range dr.Ipplinks {
		// IPL should be disabled thus state should be in initialized state
		if ipp.RxMachineFsm == nil ||
			ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDefaulted {
			t.Error("ERROR BEGIN Initial Receive Machine state is not correct", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
		}
		// port is enabled and drcp is enabled thus we should be in fast periodic state
		if ipp.PtxMachineFsm == nil {
			t.Error("ERROR BEGIN Initial Periodic Tx Machine state is not correct")
		}
		if ipp.TxMachineFsm == nil ||
			ipp.TxMachineFsm.Machine.Curr.CurrentState() != TxmStateOff {
			t.Error("ERROR BEGIN Initial Tx Machine state is not correct", TxmStateStrMap[ipp.TxMachineFsm.Machine.Curr.CurrentState()])
		}
		if ipp.NetIplShareMachineFsm == nil ||
			ipp.NetIplShareMachineFsm.Machine.Curr.CurrentState() != NetIplSharemStateNoManipulatedFramesSent {
			t.Error("ERROR BEGIN Initial Net/IPL Sharing Machine state is not correct", ipp.NetIplShareMachineFsm.Machine.Curr.CurrentState())
		}
		if ipp.IAMachineFsm == nil ||
			ipp.IAMachineFsm.Machine.Curr.CurrentState() != IAmStateIPPPortInitialize {
			t.Error("ERROR BEGIN Initial IPP Aggregator state is not correct", IAmStateStrMap[ipp.IAMachineFsm.Machine.Curr.CurrentState()])
		}

		go func(wc *chan bool) {
			for i := 0; i < 10 && ipp.IGMachineFsm.Machine.Curr.CurrentState() != IGmStateIPPGatewayUpdate; i++ {
				time.Sleep(time.Millisecond * 10)
			}
			*wc <- true
		}(&waitChan)

		<-waitChan

		if ipp.IGMachineFsm == nil ||
			ipp.IGMachineFsm.Machine.Curr.CurrentState() != IGmStateIPPGatewayUpdate {
			t.Error("ERROR BEGIN Initial IPP Gateway Machine state is not correct", GmStateStrMap[ipp.IGMachineFsm.Machine.Curr.CurrentState()])
		}
	}

	DeleteDistributedRelay(cfg.DrniName)

	if len(DistributedRelayDB) != 0 ||
		len(DistributedRelayDBList) != 0 {
		t.Error("ERROR Distributed Relay DB was not cleaned up")
	}
	if len(DRCPIppDB) != 0 ||
		len(DRCPIppDBList) != 0 {
		t.Error("ERROR IPP DB was not cleaned up")
	}
	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigDistributedRelayInValidCreateDRNoAgg(t *testing.T) {

	ConfigTestSetup()
	//a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    200,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}
	// map vlan 100 to this system
	// in real system this should be filled in by vlan membership
	cfg.DrniConvAdminGateway[100][0] = cfg.DrniPortalSystemNumber

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}

	CreateDistributedRelay(cfg)

	// configuration is incomplete because lag has not been created as part of this
	if len(DistributedRelayDB) == 0 ||
		len(DistributedRelayDBList) == 0 {
		t.Error("ERROR Distributed Relay Object was not added to global DB's")
	}
	dr, ok := DistributedRelayDB[cfg.DrniName]
	if !ok {
		t.Error("ERROR Distributed Relay Object was not found in global DB's")
	}

	// check the inital state of each of the state machines
	if dr.a != nil {
		t.Error("ERROR Agg exists no associated should exist")
	}
	if dr.PsMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Portal System Machine state was created, provisioning incomplete")
	}
	if dr.GMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Gateway Machine state was created, provisioning incomplete")
	}
	if dr.AMachineFsm != nil {
		t.Error("ERROR BEGIN Initial Aggregator System Machine state was created, provisioning incomplete")
	}

	for _, ipp := range dr.Ipplinks {
		if ipp.DRCPEnabled {
			t.Error("ERROR Why is the IPL IPP link DRCP enabled")
		}
	}

	DeleteDistributedRelay(cfg.DrniName)

	if len(DistributedRelayDB) != 0 ||
		len(DistributedRelayDBList) != 0 {
		t.Error("ERROR Distributed Relay DB was not cleaned up")
	}
	if len(DRCPIppDB) != 0 ||
		len(DRCPIppDBList) != 0 {
		t.Error("ERROR IPP DB was not cleaned up")
	}

	ConfigTestTeardwon(t)
}

func TestConfigInvalidPortalAddressString(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)
	/*




	 */
	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE", // invalid!!!
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail for bad Portal Address")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidThreePortalSystemSet(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             true, // invalid not supported
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail setting 3P system")
	}
	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidPortalSytemNumber(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            0, // invalid in 2P system
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)},
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail portal system number 0")
	}
	// invalid in 2P system
	cfg.DrniPortalSystemNumber = 3
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail portal system number 3")
	}
	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidIntraPortalLink(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail IPP link not supplied")
	}
	// invalid ipp link
	cfg.DrniIntraPortalLinkList[0] = 300
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail invalid port")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidGatewayAlgorithm(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:88:C2:01", // invalid string
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Gateway Algorithm")
	}

	cfg.DrniGatewayAlgorithm = ""
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Gateway Algorithm empty string")
	}

	cfg.DrniGatewayAlgorithm = "00:80:C2"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Gateway Algorithm wrong format to short missing actual type byte")
	}

	cfg.DrniGatewayAlgorithm = "00-80:C2-02"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Gateway Algorithm separator")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidNeighborGatewayAlgorithm(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "01:80:C2:01", // invalid string
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Gateway Algorithm")
	}

	cfg.DrniNeighborAdminGatewayAlgorithm = ""
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Gateway Algorithm empty string")
	}

	cfg.DrniNeighborAdminGatewayAlgorithm = "00:80:C2"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Gateway Algorithm wrong format to short missing actual type byte")
	}

	cfg.DrniNeighborAdminGatewayAlgorithm = "00-80:C2-02"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Gateway Algorithm separator")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidNeighborPortAlgorithm(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "10:80:C2:01", // invalid string
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Port Algorithm")
	}

	cfg.DrniNeighborAdminPortAlgorithm = ""
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Port Algorithm empty string")
	}

	cfg.DrniNeighborAdminPortAlgorithm = "00:80:C2"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Port Algorithm wrong format to short missing actual type byte")
	}

	cfg.DrniNeighborAdminPortAlgorithm = "00-80:C2-02"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Neighbor Port Algorithm separator")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidEncapMethod(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "22:80:C2:01", // invalid string
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Encap Method")
	}

	cfg.DrniEncapMethod = ""
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Encap Method empty string")
	}

	cfg.DrniEncapMethod = "00:80:C2"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Encap Method wrong format to short missing actual type byte")
	}

	cfg.DrniEncapMethod = "00-80:C2-02"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Encap Method separator")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func TestConfigInvalidPortalPortProtocolDA(t *testing.T) {
	ConfigTestSetup()
	a := OnlyForTestSetupCreateAggGroup(100)

	cfg := &DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(ipplink1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "00:80:C2:01",
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01", // invalid string
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C0:00:00:03", // invalid
	}

	err := DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail invalid Portal Port Potocol DA")
	}

	cfg.DrniIntraPortalPortProtocolDA = "01-80-C2-00-00-11"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Portal Port Potocol DA different format")
	}

	cfg.DrniIntraPortalPortProtocolDA = "80-C2-00-00-11"
	err = DistributedRelayConfigParamCheck(cfg)
	if err == nil {
		t.Error("Parameter check did not fail Invalid Portal Port Potocol DA not enough bytes")
	}

	lacp.DeleteLaAgg(a.AggId)
	ConfigTestTeardwon(t)
}

func FullBackToBackConfigTestSetup() {
	OnlyForTestSetup()
	utils.PortConfigMap[LaAggPort1NeighborActor] = utils.PortConfig{Name: LaAggPortNeighborActor1If,
		HardwareAddr: net.HardwareAddr{0x00, LaAggPort1NeighborActor, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPort2NeighborActor] = utils.PortConfig{Name: LaAggPortNeighborActor2If,
		HardwareAddr: net.HardwareAddr{0x00, LaAggPort2NeighborActor, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPort1Peer] = utils.PortConfig{Name: LaAggPortPeerIf1,
		HardwareAddr: net.HardwareAddr{0x00, LaAggPort1Peer, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[LaAggPort2Peer] = utils.PortConfig{Name: LaAggPortPeerIf2,
		HardwareAddr: net.HardwareAddr{0x00, LaAggPort2Peer, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[DRNeighborIpp1] = utils.PortConfig{Name: DRNeighborIppIf1,
		HardwareAddr: net.HardwareAddr{0x00, DRNeighborIpp1, 0x11, 0x22, 0x22, 0x33},
	}
	utils.PortConfigMap[DRNeighborIpp2] = utils.PortConfig{Name: DRNeighborIppIf2,
		HardwareAddr: net.HardwareAddr{0x00, DRNeighborIpp2, 0x11, 0x22, 0x22, 0x33},
	}
}

func FullBackToBackConfigTestTeardown(t *testing.T) {
	OnlyForTestTeardown(t)
	delete(utils.PortConfigMap, LaAggPort1NeighborActor)
	delete(utils.PortConfigMap, LaAggPort2NeighborActor)
	delete(utils.PortConfigMap, LaAggPort1Peer)
	delete(utils.PortConfigMap, LaAggPort2Peer)
	delete(utils.PortConfigMap, DRNeighborIpp1)
	delete(utils.PortConfigMap, DRNeighborIpp2)
}

type ThreeNodeConfig struct {
	neighborbridge SimulationNeighborBridge
	bridge1        lacp.SimulationBridge
	bridge2        lacp.SimulationBridge
	cfg            DistributedRelayConfig
	cfg2           DistributedRelayConfig
	a1conf         *lacp.LaAggConfig
	a2conf         *lacp.LaAggConfig
	a3conf         *lacp.LaAggConfig
	p1conf         *lacp.LaAggPortConfig
	p2conf         *lacp.LaAggPortConfig
	p3conf         *lacp.LaAggPortConfig
	p4conf         *lacp.LaAggPortConfig
}

func Teardown3NodeMlag(mlagcfg *ThreeNodeConfig, t *testing.T) {
	// cleanup the provisioning
	close(mlagcfg.bridge1.RxLacpPort1)
	close(mlagcfg.bridge1.RxLacpPort2)
	close(mlagcfg.bridge2.RxLacpPort1)
	close(mlagcfg.bridge2.RxLacpPort2)
	mlagcfg.bridge1.RxLacpPort1 = nil
	mlagcfg.bridge1.RxLacpPort2 = nil
	mlagcfg.bridge2.RxLacpPort1 = nil
	mlagcfg.bridge2.RxLacpPort2 = nil
	lacp.DeleteLaAgg(mlagcfg.a1conf.Id)
	lacp.DeleteLaAgg(mlagcfg.a2conf.Id)
	lacp.DeleteLaAgg(mlagcfg.a3conf.Id)
	for _, sgi := range lacp.LacpSysGlobalInfoGet() {
		if len(sgi.AggList) > 0 || len(sgi.AggMap) > 0 {
			t.Error("System Agg List or Map is not empty", sgi.AggList, sgi.AggMap)
		}
		if len(sgi.PortList) > 0 || len(sgi.PortMap) > 0 {
			t.Error("System Port List or Map is not empty", sgi.PortList, sgi.PortMap)
		}
	}

	DeleteDistributedRelay(mlagcfg.cfg.DrniName)
	DeleteDistributedRelay(mlagcfg.cfg2.DrniName)

	// must be called to initialize the global
	LaSystem1NeighborActor := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x01, 0x64}}
	LaSystem2NeighborActor := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x02, 0x64}}

	LaSystemPeer := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	lacp.LacpSysGlobalInfoDestroy(LaSystem1NeighborActor)
	lacp.LacpSysGlobalInfoDestroy(LaSystem2NeighborActor)
	lacp.LacpSysGlobalInfoDestroy(LaSystemPeer)

	//delete(testBlockMap, LaAggPort1NeighborActor)
	//delete(testBlockMap, LaAggPort2NeighborActor)
}
func Setup3NodeMlag() *ThreeNodeConfig {
	threenodecfg := &ThreeNodeConfig{}
	// actor1 to actor2
	threenodecfg.neighborbridge = SimulationNeighborBridge{
		Port1:      DRNeighborIpp1,
		Port2:      DRNeighborIpp2,
		RxIppPort1: make(chan gopacket.Packet, 10),
		RxIppPort2: make(chan gopacket.Packet, 10),
	}
	ipp1Key := IppDbKey{
		Name:   DRNeighborIppIf1,
		DrName: "DR-1",
	}
	ipp2Key := IppDbKey{
		Name:   DRNeighborIppIf2,
		DrName: "DR-2",
	}

	DRGlobalSystem.DRSystemGlobalRegisterTxCallback(ipp1Key, threenodecfg.neighborbridge.TxViaGoChannel)
	DRGlobalSystem.DRSystemGlobalRegisterTxCallback(ipp2Key, threenodecfg.neighborbridge.TxViaGoChannel)

	DrRxMain(uint16(DRNeighborIpp1), "00:00:DE:AD:BE:EF", threenodecfg.neighborbridge.RxIppPort1)
	DrRxMain(uint16(DRNeighborIpp2), "00:00:DE:AD:BE:EF", threenodecfg.neighborbridge.RxIppPort2)

	// Lets create the Distributed Relay
	threenodecfg.cfg = DistributedRelayConfig{
		DrniName:                          "DR-1",
		DrniPortalAddress:                 "00:00:DE:AD:BE:EF",
		DrniPortalPriority:                128,
		DrniThreePortalSystem:             false,
		DrniPortalSystemNumber:            1,
		DrniIntraPortalLinkList:           [3]uint32{uint32(DRNeighborIpp1)}, // no link supplied is invalid
		DrniAggregator:                    100,
		DrniGatewayAlgorithm:              "00:80:C2:01",
		DrniNeighborAdminGatewayAlgorithm: "00:80:C2:01",
		DrniNeighborAdminPortAlgorithm:    "10:80:C2:01", // invalid string
		DrniNeighborAdminDRCPState:        "00000000",
		DrniEncapMethod:                   "00:80:C2:01",
		DrniPortConversationControl:       false,
		DrniIntraPortalPortProtocolDA:     "01:80:C2:00:00:03", // only supported value that we are going to support
	}

	// create first drni
	CreateDistributedRelay(&threenodecfg.cfg)

	threenodecfg.cfg2 = threenodecfg.cfg

	// create second drni
	threenodecfg.cfg2.DrniName = "DR-2"
	threenodecfg.cfg2.DrniPortalSystemNumber = 2
	threenodecfg.cfg2.DrniAggregator = 500
	threenodecfg.cfg2.DrniIntraPortalLinkList = [3]uint32{uint32(DRNeighborIpp2)}
	CreateDistributedRelay(&threenodecfg.cfg2)

	// must be called to initialize the global
	LaSystem1NeighborActor := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x01, 0x64}}
	LaSystem2NeighborActor := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x02, 0x64}}

	LaSystemPeer := lacp.LacpSystem{Actor_System_priority: 128,
		Actor_System: [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xC8}}

	// actor1 to peer
	threenodecfg.bridge1 = lacp.SimulationBridge{
		Port1:       LaAggPort1NeighborActor,
		Port2:       LaAggPort1Peer,
		RxLacpPort1: make(chan gopacket.Packet, 10),
		RxLacpPort2: make(chan gopacket.Packet, 10),
	}
	// actor2 to peer
	threenodecfg.bridge2 = lacp.SimulationBridge{
		Port1:       LaAggPort2NeighborActor,
		Port2:       LaAggPort2Peer,
		RxLacpPort1: make(chan gopacket.Packet, 10),
		RxLacpPort2: make(chan gopacket.Packet, 10),
	}

	Actor1System := lacp.LacpSysGlobalInfoInit(LaSystem1NeighborActor)
	Actor2System := lacp.LacpSysGlobalInfoInit(LaSystem2NeighborActor)
	PeerSystem := lacp.LacpSysGlobalInfoInit(LaSystemPeer)
	// la and ipp ports
	Actor1System.LaSysGlobalRegisterTxCallback(LaAggPortNeighborActor1If, threenodecfg.bridge1.TxViaGoChannel)
	// la and ipp ports
	Actor2System.LaSysGlobalRegisterTxCallback(LaAggPortNeighborActor2If, threenodecfg.bridge2.TxViaGoChannel)
	// la ports
	PeerSystem.LaSysGlobalRegisterTxCallback(LaAggPortPeerIf1, threenodecfg.bridge1.TxViaGoChannel)
	PeerSystem.LaSysGlobalRegisterTxCallback(LaAggPortPeerIf2, threenodecfg.bridge2.TxViaGoChannel)

	// port 1
	lacp.LaRxMain(threenodecfg.bridge1.Port1, threenodecfg.bridge1.RxLacpPort1)
	// port 2
	lacp.LaRxMain(threenodecfg.bridge1.Port2, threenodecfg.bridge1.RxLacpPort2)

	// port 1
	lacp.LaRxMain(threenodecfg.bridge2.Port1, threenodecfg.bridge2.RxLacpPort1)
	// port 2
	lacp.LaRxMain(threenodecfg.bridge2.Port2, threenodecfg.bridge2.RxLacpPort2)

	// lag system 1 actor && lag system 2 actor
	threenodecfg.a1conf = &lacp.LaAggConfig{
		Mac:  [6]uint8{0x00, 0x00, 0x01, 0x01, 0x01, 0x01},
		Id:   100,
		Name: "agg1",
		Key:  100,
		Lacp: lacp.LacpConfigInfo{Interval: lacp.LacpSlowPeriodicTime,
			Mode:           lacp.LacpModeActive,
			SystemIdMac:    "00:00:00:00:01:64",
			SystemPriority: 128},
	}

	threenodecfg.a2conf = &lacp.LaAggConfig{
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x01, 0x01, 0x01},
		Id:   500,
		Name: "agg2",
		Key:  500,
		Lacp: lacp.LacpConfigInfo{Interval: lacp.LacpSlowPeriodicTime,
			Mode:           lacp.LacpModeActive,
			SystemIdMac:    "00:00:00:00:02:64",
			SystemPriority: 128},
	}

	// lag sytem 3 peer
	threenodecfg.a3conf = &lacp.LaAggConfig{
		Mac:  [6]uint8{0x00, 0x00, 0x02, 0x02, 0x02, 0x02},
		Id:   200,
		Name: "agg3",
		Key:  200,
		Lacp: lacp.LacpConfigInfo{Interval: lacp.LacpSlowPeriodicTime,
			Mode:           lacp.LacpModeActive,
			SystemIdMac:    "00:00:00:00:00:C8",
			SystemPriority: 128},
	}

	// Create Aggregation
	lacp.CreateLaAgg(threenodecfg.a1conf)
	lacp.CreateLaAgg(threenodecfg.a2conf)
	lacp.CreateLaAgg(threenodecfg.a3conf)

	// actor        peer
	// p1config <-> p2config
	// p3config <-> p4config
	threenodecfg.p1conf = &lacp.LaAggPortConfig{
		Id:     LaAggPort1NeighborActor,
		Prio:   0x80,
		Key:    100,
		AggId:  100,
		Enable: true,
		Mode:   lacp.LacpModeActive,
		//Timeout: LacpFastPeriodicTime,
		Properties: lacp.PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPort1NeighborActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: lacp.LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortNeighborActor1If,
		TraceEna: true,
	}
	threenodecfg.p2conf = &lacp.LaAggPortConfig{
		Id:     LaAggPort1Peer,
		Prio:   0x80,
		Key:    200,
		AggId:  200,
		Enable: true,
		Mode:   lacp.LacpModeActive,
		Properties: lacp.PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPort1Peer, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: lacp.LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortPeerIf1,
		TraceEna: true,
	}

	threenodecfg.p3conf = &lacp.LaAggPortConfig{
		Id:     LaAggPort2NeighborActor,
		Prio:   0x80,
		Key:    500,
		AggId:  500,
		Enable: true,
		Mode:   lacp.LacpModeActive,
		//Timeout: LacpFastPeriodicTime,
		Properties: lacp.PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPort2NeighborActor, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: lacp.LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortNeighborActor2If,
		TraceEna: true,
	}

	threenodecfg.p4conf = &lacp.LaAggPortConfig{
		Id:     LaAggPort2Peer,
		Prio:   0x80,
		Key:    200,
		AggId:  200,
		Enable: true,
		Mode:   lacp.LacpModeActive,
		Properties: lacp.PortProperties{
			Mac:    net.HardwareAddr{0x00, LaAggPort2Peer, 0xDE, 0xAD, 0xBE, 0xEF},
			Speed:  1000000000,
			Duplex: lacp.LacpPortDuplexFull,
			Mtu:    1500,
		},
		IntfId:   LaAggPortPeerIf2,
		TraceEna: true,
	}

	// actor / neighbor
	lacp.CreateLaAggPort(threenodecfg.p1conf)
	// peer
	lacp.CreateLaAggPort(threenodecfg.p2conf)

	// actor / neighbor
	lacp.CreateLaAggPort(threenodecfg.p3conf)
	// peer
	lacp.CreateLaAggPort(threenodecfg.p4conf)

	return threenodecfg
}

func Verify3NodeMlag(mlagcfg *ThreeNodeConfig, step string, convlist []uint16, t *testing.T) {
	testWait := make(chan bool)

	var p1 *lacp.LaAggPort
	var p2 *lacp.LaAggPort
	// TODO this should fail as the ports should not sync up with the peer because the agg key does not agree between
	// the ports
	if lacp.LaFindPortById(mlagcfg.p1conf.Id, &p1) &&
		lacp.LaFindPortById(mlagcfg.p2conf.Id, &p2) {
		//fmt.Println("Checking for port to come up in distributed state (0)")
		go func(wc chan bool) {
			for i := 0; i < 10 &&
				(p1.MuxMachineFsm.Machine.Curr.CurrentState() != lacp.LacpMuxmStateDistributing ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != lacp.LacpMuxmStateDistributing); i++ {
				//fmt.Println("waiting for distrubuted (0)")
				time.Sleep(time.Second * 1)
			}
			wc <- true
		}(testWait)

		<-testWait
		close(testWait)

		State1 := lacp.GetLaAggPortActorOperState(mlagcfg.p1conf.Id)
		State2 := lacp.GetLaAggPortActorOperState(mlagcfg.p2conf.Id)

		const portUpState = lacp.LacpStateActivityBit | lacp.LacpStateAggregationBit |
			lacp.LacpStateSyncBit | lacp.LacpStateCollectingBit | lacp.LacpStateDistributingBit

		if !lacp.LacpStateIsSet(State1, portUpState) {
			t.Error(fmt.Sprintf("step: %s Actor Port State %s did not come up properly with peer expected %s", step, lacp.LacpStateToStr(State1), lacp.LacpStateToStr(portUpState)))
		}
		if !lacp.LacpStateIsSet(State2, portUpState) {
			t.Error(fmt.Sprintf("step: %s Peer Port State %s did not come up properly with actor expected %s", step, lacp.LacpStateToStr(State2), lacp.LacpStateToStr(portUpState)))
		}

		// TODO check the States of the other State machines
	} else {
		t.Error(fmt.Sprintf("step: %s Unable to find port just created", step))
	}

	testWait = make(chan bool)
	// TODO this should fail as the ports should not sync up with the peer because the agg key does not agree between
	// the ports
	if lacp.LaFindPortById(mlagcfg.p3conf.Id, &p1) &&
		lacp.LaFindPortById(mlagcfg.p4conf.Id, &p2) {
		//fmt.Println("Checking for port to come up in distributed state (1)")
		go func(wc chan bool) {
			for i := 0; i < 10 &&
				(p1.MuxMachineFsm.Machine.Curr.CurrentState() != lacp.LacpMuxmStateDistributing ||
					p2.MuxMachineFsm.Machine.Curr.CurrentState() != lacp.LacpMuxmStateDistributing); i++ {
				//fmt.Println("waiting for distrubuted (1)")
				time.Sleep(time.Second * 1)
			}
			wc <- true
		}(testWait)

		<-testWait
		close(testWait)

		State1 := lacp.GetLaAggPortActorOperState(mlagcfg.p3conf.Id)
		State2 := lacp.GetLaAggPortActorOperState(mlagcfg.p4conf.Id)

		const portUpState = lacp.LacpStateActivityBit | lacp.LacpStateAggregationBit |
			lacp.LacpStateSyncBit | lacp.LacpStateCollectingBit | lacp.LacpStateDistributingBit

		if !lacp.LacpStateIsSet(State1, portUpState) {
			t.Error(fmt.Sprintf("step: %s Actor Port State %s did not come up properly with peer expected %s", step, lacp.LacpStateToStr(State1), lacp.LacpStateToStr(portUpState)))
		}
		if !lacp.LacpStateIsSet(State2, portUpState) {
			t.Error(fmt.Sprintf("step: %s Peer Port State %s did not come up properly with actor expected %s", step, lacp.LacpStateToStr(State2), lacp.LacpStateToStr(portUpState)))
		}

		// TODO check the States of the other State machines
	} else {
		t.Error(fmt.Sprintf("step: %s Unable to find port just created", step))
	}
	var dr *DistributedRelay
	if !DrFindByAggregator(int32(mlagcfg.cfg.DrniAggregator), &dr) {
		t.Error(fmt.Sprintf("step: %s Error could not find te DR by local aggregator", step))
	}
	var dr2 *DistributedRelay
	if !DrFindByAggregator(int32(mlagcfg.cfg2.DrniAggregator), &dr2) {
		t.Error(fmt.Sprintf("step: %s Error could not find te DR by local aggregator", step))
	}
	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10 &&
			(!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr state to converge", dr.DRFHomeOperDRCPState.String())
	close(testWait)
	testWait = make(chan bool)

	go func(wc chan bool) {
		for i := 0; i < 10 &&
			(!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr2.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr2.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr2 state to converge", dr2.DRFHomeOperDRCPState.String())
	close(testWait)

	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10; i++ {
			for _, ipp := range dr.Ipplinks {
				convnotfound := false
				if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
					time.Sleep(time.Second * 1)
				} else {
					for _, cid := range convlist {
						if !ipp.IppGatewayConversationPasses[cid] {
							convnotfound = true
						}
					}
					if convnotfound {
						time.Sleep(time.Second * 1)
					}
				}
			}
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr2 state to converge", dr2.DRFHomeOperDRCPState.String())
	close(testWait)

	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10; i++ {
			for _, ipp := range dr2.Ipplinks {
				convnotfound := false
				if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
					!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
					time.Sleep(time.Second * 1)
				} else {
					for _, cid := range convlist {
						if !ipp.IppGatewayConversationPasses[cid] {
							convnotfound = true
						}
					}
					if convnotfound {
						time.Sleep(time.Second * 1)
					}
				}
			}
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr2 state to converge", dr2.DRFHomeOperDRCPState.String())
	close(testWait)

	if len(dr.DRAggregatorDistributedList) != 1 {
		t.Error(fmt.Sprintf("step: %s Error Distributed Ports does not equal %v", step, dr.DRAggregatorDistributedList))
	} else {
		if dr.DRAggregatorDistributedList[0] != LaAggPort1NeighborActor {
			t.Error(fmt.Sprintf("step: %s Error Distributed Ports Incorrect port found %v", step, dr.DRAggregatorDistributedList[0]))
		}
	}

	if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error(fmt.Sprintf("step: %s Error IPP HOME did not sync up as expected current state %v", step, dr.DRFHomeOperDRCPState.String()))
	}

	for _, ipp := range dr.Ipplinks {
		if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
			t.Error(fmt.Sprintf("step: %s Error IPP NEIGHBOR did not sync up as expected current state %v", ipp.DRFNeighborOperDRCPState.String()))
		}
		for _, cid := range convlist {
			if !ipp.IppGatewayConversationPasses[cid] {
				t.Error("Error IPP Neighbor did not set conversation passes for 100 ", ipp.Id)
			}
		}
		//sort.Sort(sortPortList(testBlockMap[ipp.Name]))
		tmpPortList := make([]string, 0)

		for _, p := range dr.a.DistributedPortNumList {
			tmpPortList = append(tmpPortList, p)
		}
		//sort.Sort(sortPortList(tmpPortList))

		//if len(testBlockMap[int32(ipp.Id)]) != len(tmpPortList) {
		//	t.Error("Error Block Map not set correctly, expected ", dr.a.PortNumList, "found", tmpPortList, "len1", len(testBlockMap[int32(ipp.Id)]), "len2", len(tmpPortList))
		//}

		for _, p1 := range testBlockMap[ipp.Name] {
			for _, p2 := range tmpPortList {
				if p1 != p2 {
					t.Error("step:", step, "Error (2) Block Map not set correctly, expected ", dr.a.PortNumList, "found", testBlockMap[ipp.Name])
				}
			}
		}
	}

	if len(dr2.DRAggregatorDistributedList) != 1 {
		t.Error("step:", step, "Error Distributed Ports does not equal", dr2.DRAggregatorDistributedList)
	} else {

		if dr2.DRAggregatorDistributedList[0] != LaAggPort2NeighborActor {
			t.Error("step: ", step, "Error Distributed Ports Incorrect port found", dr2.DRAggregatorDistributedList[0])
		}
	}

	if !dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error("step: ", step, "Error IPP HOME did not sync up as expected current state ", dr2.DRFHomeOperDRCPState.String())
	}

	for _, ipp := range dr2.Ipplinks {
		if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
			t.Error("step:", step, "Error IPP NEIGHBOR did not sync up as expected current state ", ipp.DRFNeighborOperDRCPState.String())
		}

		for _, cid := range convlist {
			if !ipp.IppGatewayConversationPasses[cid] {
				t.Error("step:", step, "Error IPP Neighbor did not set conversation passes for 100 ")
			}
		}
		//sort.Sort(sortPortList(testBlockMap[ipp.Name]))
		tmpPortList := make([]string, 0)

		for _, p := range dr2.a.DistributedPortNumList {
			tmpPortList = append(tmpPortList, p)
		}
		//sort.Sort(sortPortList(tmpPortList))

		//if len(testBlockMap[int32(ipp.Id)]) != len(tmpPortList) {
		//	t.Error("Error Block Map not set correctly, expected ", dr2.a.PortNumList, "found", tmpPortList)
		//}

		for _, p1 := range testBlockMap[ipp.Name] {
			for _, p2 := range tmpPortList {
				if p1 != p2 {
					t.Error("step:", step, "Error Block Map not set correctly, expected ", dr2.a.PortNumList, "found", testBlockMap[ipp.Name])
				}
			}
		}
	}
}

// 3 node system where two neighbors are connected to 1 peer device
func TestConfigCreateBackToBackMLagAndPeer1(t *testing.T) {

	FullBackToBackConfigTestSetup()

	mlagcfg := Setup3NodeMlag()
	//time.Sleep(time.Second * 20)

	// basic verify
	Verify3NodeMlag(mlagcfg, "basic", []uint16{100}, t)
	Teardown3NodeMlag(mlagcfg, t)

	FullBackToBackConfigTestTeardown(t)
}

// Add a new conversation to both MLAG's
func TestConfigCreateBackToBackMLagAndPeerValidAddDelVlan(t *testing.T) {

	FullBackToBackConfigTestSetup()

	mlagcfg := Setup3NodeMlag()
	//time.Sleep(time.Second * 20)

	// basic verify
	Verify3NodeMlag(mlagcfg, "basic", []uint16{100}, t)

	var dr *DistributedRelay
	if !DrFindByAggregator(int32(mlagcfg.cfg.DrniAggregator), &dr) {
		t.Error("Error could not find te DR by local aggregator")
	}
	var dr2 *DistributedRelay
	if !DrFindByAggregator(int32(mlagcfg.cfg2.DrniAggregator), &dr2) {
		t.Error("Error could not find te DR by local aggregator")
	}

	// Add a new conversation vlan 200, with port list created with conversation
	cfg := &DRConversationConfig{
		DrniName: mlagcfg.cfg.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    200,
	}

	for _, aggport := range dr.a.PortNumList {
		cfg.PortList = append(cfg.PortList, int32(aggport))
	}

	CreateConversationId(cfg)

	cfg = &DRConversationConfig{
		DrniName: mlagcfg.cfg2.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    200,
	}
	for _, aggport := range dr2.a.PortNumList {
		cfg.PortList = append(cfg.PortList, int32(aggport))
	}

	CreateConversationId(cfg)

	Verify3NodeMlag(mlagcfg, "after vlan add", []uint16{100, 200}, t)

	testWait := make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10 &&
			(!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr state to converge", dr.DRFHomeOperDRCPState.String())
	close(testWait)
	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10 &&
			(!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr2.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr2.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr2.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr2 state to converge", dr2.DRFHomeOperDRCPState.String())
	close(testWait)

	allippList := make([]*DRCPIpp, 0)
	for _, p := range dr.Ipplinks {
		allippList = append(allippList, p)
	}
	for _, p := range dr2.Ipplinks {
		allippList = append(allippList, p)
	}

	for _, ipp := range allippList {
		testWait = make(chan bool)

		go func(wc chan bool) {

			for i := 0; i < 10 &&
				(ipp.IppGatewayConversationPasses[100] ||
					!ipp.IppGatewayConversationPasses[200]); i++ {
				//fmt.Println("waiting for dr2 state to converge", dr2.DRFHomeOperDRCPState.String(), i)
				time.Sleep(time.Second * 1)
			}
			wc <- true
		}(testWait)

		<-testWait

		if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
			t.Error("Error IPP NEIGHBOR did not sync up as expected current state ", ipp.DRFNeighborOperDRCPState.String())
		}

		if !ipp.IppGatewayConversationPasses[100] {
			t.Error("Error IPP Neighbor did not set conversation passes for 100 ", ipp.Id)
		}
		if !ipp.IppGatewayConversationPasses[200] {
			t.Error("Error IPP Neighbor did not set conversation passes for 200 ", ipp.Id)
		}
	}

	if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error("Error IPP HOME did not sync up as expected current state ", dr.DRFHomeOperDRCPState.String())
	}

	if !dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error("Error IPP HOME did not sync up as expected current state ", dr2.DRFHomeOperDRCPState.String())
	}

	// Del a conversation vlan 100, with port list created with conversation
	// Add a new conversation vlan 200, with port list created with conversation
	cfg = &DRConversationConfig{
		DrniName: mlagcfg.cfg.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    100,
	}

	for _, aggport := range dr.a.PortNumList {
		cfg.PortList = append(cfg.PortList, int32(aggport))
	}

	DeleteConversationId(cfg, true)

	cfg = &DRConversationConfig{
		DrniName: mlagcfg.cfg2.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    100,
	}
	for _, aggport := range dr2.a.PortNumList {
		cfg.PortList = append(cfg.PortList, int32(aggport))
	}

	DeleteConversationId(cfg, true)

	Verify3NodeMlag(mlagcfg, "after vlan del", []uint16{200}, t)

	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10 &&
			(!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr state to converge", dr.DRFHomeOperDRCPState.String())
	close(testWait)
	testWait = make(chan bool)

	go func(wc chan bool) {

		for i := 0; i < 10 &&
			(!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
				!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) ||
				len(dr2.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList) != 1 ||
				len(dr2.DrniPortalSystemState[dr.Ipplinks[0].DRFNeighborPortalSystemNumber].PortIdList) != 1); i++ {
			//fmt.Println("waiting for dr2 state to converge", dr2.DRFHomeOperDRCPState.String(), i)
			time.Sleep(time.Second * 1)
		}
		wc <- true
	}(testWait)

	<-testWait
	//fmt.Println("after wait for dr2 state to converge", dr2.DRFHomeOperDRCPState.String())
	close(testWait)

	allippList = make([]*DRCPIpp, 0)
	for _, p := range dr.Ipplinks {
		allippList = append(allippList, p)
	}
	for _, p := range dr2.Ipplinks {
		allippList = append(allippList, p)
	}

	for _, ipp := range allippList {
		testWait = make(chan bool)

		go func(wc chan bool) {

			for i := 0; i < 10 &&
				(ipp.IppGatewayConversationPasses[100] ||
					!ipp.IppGatewayConversationPasses[200]); i++ {
				//fmt.Println("waiting for dr2 state to converge", dr2.DRFHomeOperDRCPState.String(), i)
				time.Sleep(time.Second * 1)
			}
			wc <- true
		}(testWait)

		<-testWait

		if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
			!ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStatePortSync) {
			t.Error("Error IPP NEIGHBOR did not sync up as expected current state ", ipp.DRFNeighborOperDRCPState.String())
		}

		if ipp.IppGatewayConversationPasses[100] {
			t.Error("Error IPP Neighbor did not clear conversation passes for 100 ", ipp.Id)
		}
		if !ipp.IppGatewayConversationPasses[200] {
			t.Error("Error IPP Neighbor did not set conversation passes for 200 ", ipp.Id)
		}
	}

	if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error("Error IPP HOME did not sync up as expected current state ", dr.DRFHomeOperDRCPState.String())
	}

	if !dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateIPPActivity) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateHomeGatewayBit) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStateGatewaySync) ||
		!dr2.DRFHomeOperDRCPState.GetState(layers.DRCPStatePortSync) {
		t.Error("Error IPP HOME did not sync up as expected current state ", dr2.DRFHomeOperDRCPState.String())
	}

	Teardown3NodeMlag(mlagcfg, t)

	FullBackToBackConfigTestTeardown(t)
}
