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

// conversationId.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"net"
	"testing"
	"time"
	asicdmock "utils/asicdClient/mock"
	"utils/commonDefs"
	"utils/fsm"
	"utils/logging"
)

type MyTestMock2 struct {
	asicdmock.MockAsicdClientMgr
}

func SliceEqual(a []uint8, b []uint8) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func (m *MyTestMock2) GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error) {

	return nil, nil
}

func OnlyForConversationIdTestSetup() {
	logger, _ := logging.NewLogger("lacpd", "TEST", false)
	utils.SetLaLogger(logger)
	utils.DeleteAllAsicDPlugins()
	utils.SetAsicDPlugin(&MyTestMock2{})
	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		ConversationIdMap[i].Valid = false
		ConversationIdMap[i].PortList = nil
		ConversationIdMap[i].Cvlan = 0
		ConversationIdMap[i].Refcnt = 0
		ConversationIdMap[i].Idtype = [4]uint8{}
	}
	// fill in conversations
	//GetAllCVIDConversations()
}

func OnlyForConversationIdTestTeardown() {

	utils.SetLaLogger(nil)
	utils.DeleteAllAsicDPlugins()
	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		ConversationIdMap[i].Valid = false
		ConversationIdMap[i].PortList = nil
		ConversationIdMap[i].Cvlan = 0
		ConversationIdMap[i].Refcnt = 0
		ConversationIdMap[i].Idtype = [4]uint8{}
	}
}

func OnlyForConversationIdTestSetupCreateAggGroup(aggId uint32) *lacp.LaAggregator {
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
		Prio:    128,
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

func ConversationIdTestSetup() {
	OnlyForConversationIdTestSetup()
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
func ConversationIdTestTeardwon() {

	OnlyForConversationIdTestTeardown()
	delete(utils.PortConfigMap, ipplink1)
	delete(utils.PortConfigMap, aggport1)
	delete(utils.PortConfigMap, aggport2)
	delete(utils.PortConfigMap, ipplink2)
	delete(utils.PortConfigMap, aggport3)
	delete(utils.PortConfigMap, aggport4)
}

func TestConversationIdVlanMembershipCreateNoPorts(t *testing.T) {

	ConversationIdTestSetup()
	a := OnlyForConversationIdTestSetupCreateAggGroup(200)

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

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	for _, disport := range a.DistributedPortNumList {
		var aggp *lacp.LaAggPort
		foundPort := false
		for lacp.LaGetPortNext(&aggp) && !foundPort {
			if aggp.IntfNum == disport {
				dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList, int32(aggp.PortNum))
			}
		}
	}
	// set gateway info and digest
	dr.SetTimeSharingPortAndGatwewayDigest()

	//dr.DRFHomeConversationPortListDigest = drcp.PortalConfigInfo.PortDigest
	//dr.DRFHomeConversationGatewayListDigest = drcp.PortalConfigInfo.GatewayDigest

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)
	DrcpRxMachineFSMBuild(ipp)

	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	// SETUP create a vlan who does not have any port members
	conversationCfg := &DRConversationConfig{
		DrniName: dr.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    100,
	}

	CreateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if dr.DrniConvAdminGateway[100] != nil {
		t.Error("ERRRO DrniConvAdminGateway values have been set", dr.DrniConvAdminGateway[100])
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	ConversationIdTestTeardwon()
}

func TestConversationIdVlanMembershipCreateNoPortsThenAddDelPort(t *testing.T) {

	ConversationIdTestSetup()
	a := OnlyForConversationIdTestSetupCreateAggGroup(200)

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

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	for _, disport := range a.DistributedPortNumList {
		var aggp *lacp.LaAggPort
		foundPort := false
		for lacp.LaGetPortNext(&aggp) && !foundPort {
			if aggp.IntfNum == disport {
				dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList, int32(aggp.PortNum))
			}
		}
	}
	// set gateway info and digest
	dr.SetTimeSharingPortAndGatwewayDigest()

	//dr.DRFHomeConversationPortListDigest = drcp.PortalConfigInfo.PortDigest
	//dr.DRFHomeConversationGatewayListDigest = drcp.PortalConfigInfo.GatewayDigest

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)
	DrcpRxMachineFSMBuild(ipp)

	dr.GMachineFsm.Machine.Curr.SetState(GmStatePsGatewayUpdate)

	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	// SETUP create a vlan who does not have any port members
	conversationCfg := &DRConversationConfig{
		DrniName: dr.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    100,
	}

	CreateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if dr.DrniConvAdminGateway[100] != nil {
		t.Error("ERRRO DrniConvAdminGateway values have been set")
	}

	// Add Ports
	conversationCfg.PortList = []int32{aggport1, aggport2}
	UpdateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if !SliceEqual(dr.DrniConvAdminGateway[100], []uint8{2, 1}) {
		t.Error("ERRRO DrniConvAdminGateway values have not been set as expected", dr.DrniConvAdminGateway[100])
	}

	eventReceived := false
	go func(evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !(*evrx) {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "CONVERSATION ID: FORCE TEST FAIL",
			}
		}
	}(&eventReceived)

	evt := <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERRRO Failed to received portal change event")
	}
	eventReceived = true

	// Del agg port
	conversationCfg.PortList = []int32{aggport2}
	UpdateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if dr.DrniConvAdminGateway[100] != nil {
		t.Error("ERRRO DrniConvAdminGateway values have not been cleared as expected")
	}

	// Del random port
	conversationCfg.PortList = []int32{aggport2}
	UpdateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated/cleared as expected")
	}

	// admin gateway should be empty
	if dr.DrniConvAdminGateway[100] != nil {
		t.Error("ERRRO DrniConvAdminGateway values have not been set as expected")
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestConversationIdVlanMembershipCreateWithPortsThenDelPorts(t *testing.T) {

	ConversationIdTestSetup()
	a := OnlyForConversationIdTestSetupCreateAggGroup(200)

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

	err := DistributedRelayConfigParamCheck(cfg)
	if err != nil {
		t.Error("Parameter check failed for what was expected to be a valid config", err)
	}
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	// add the port to the local distributed list so that the digests can be
	// calculated
	dr.DRAggregatorDistributedList = make([]int32, 0)
	for _, disport := range a.DistributedPortNumList {
		var aggp *lacp.LaAggPort
		foundPort := false
		for lacp.LaGetPortNext(&aggp) && !foundPort {
			if aggp.IntfNum == disport {
				dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList, int32(aggp.PortNum))
			}
		}
	}
	// set gateway info and digest
	dr.SetTimeSharingPortAndGatwewayDigest()

	//dr.DRFHomeConversationPortListDigest = drcp.PortalConfigInfo.PortDigest
	//dr.DRFHomeConversationGatewayListDigest = drcp.PortalConfigInfo.GatewayDigest

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)
	DrcpRxMachineFSMBuild(ipp)

	dr.PsMachineFsm.Machine.Curr.SetState(PsmStatePortalSystemInitialize)

	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	//dr.DRFHomeConversationPortListDigest = drcp.PortalConfigInfo.PortDigest
	//dr.DRFHomeConversationGatewayListDigest = drcp.PortalConfigInfo.GatewayDigest

	// lets get the IPP
	//ipp := dr.Ipplinks[0]

	// SETUP create a vlan who does not have any port members
	conversationCfg := &DRConversationConfig{
		DrniName: dr.DrniName,
		Idtype:   GATEWAY_ALGORITHM_CVID,
		Cvlan:    100,
		PortList: []int32{aggport1, aggport2},
	}

	CreateConversationId(conversationCfg)

	if !ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if !SliceEqual(dr.DrniConvAdminGateway[100], []uint8{2, 1}) {
		t.Error("ERRRO DrniConvAdminGateway values have not been set as expected", dr.DrniConvAdminGateway[100])
	}

	eventReceived := false
	go func(evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !(*evrx) {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "CONVERSATION ID: FORCE TEST FAIL",
			}
		}
	}(&eventReceived)

	evt := <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERRRO Failed to received portal change event")
	}
	eventReceived = true

	// Del vlan
	conversationCfg.PortList = []int32{aggport1, aggport2}
	DeleteConversationId(conversationCfg, false)

	if ConversationIdMap[100].Valid {
		t.Error("ERRRO Conversation Map was not updated as expected")
	}

	// admin gateway should be empty
	if dr.DrniConvAdminGateway[100] != nil {
		t.Error("ERRRO DrniConvAdminGateway values have not been cleared as expected")
	}

	eventReceived = false
	go func(evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !(*evrx) {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "CONVERSATION ID: FORCE TEST FAIL",
			}
		}
	}(&eventReceived)

	evt = <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERRRO Failed to received portal change event")
	}
	eventReceived = true
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}
