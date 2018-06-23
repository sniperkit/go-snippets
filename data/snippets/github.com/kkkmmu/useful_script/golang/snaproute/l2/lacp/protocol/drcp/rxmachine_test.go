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
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"net"
	"strings"
	"testing"
	"time"
	"utils/fsm"
	"utils/logging"

	"github.com/google/gopacket/layers"
)

//const ipplink1 int32 = 3
//const aggport1 int32 = 1
//const aggport2 int32 = 2

//const ipplink2 int32 = 4
//const aggport3 int32 = 5
//const aggport4 int32 = 6

//type MyTestMock struct {
//	asicdmock.MockAsicdClientMgr
//}

func OnlyForRxMachineTestSetup() {
	logger, _ := logging.NewLogger("lacpd", "TEST", false)
	utils.SetLaLogger(logger)
	utils.DeleteAllAsicDPlugins()
	utils.SetAsicDPlugin(&MyTestMock{})
	// fill in conversations
	GetAllCVIDConversations()
}

func OnlyForRxMachineTestTeardown(t *testing.T) {

	//utils.SetLaLogger(nil)
	//utils.DeleteAllAsicDPlugins()
	//ConversationIdMap[100].Valid = false
	//ConversationIdMap[100].PortList = nil
	//ConversationIdMap[100].Cvlan = 0
	//ConversationIdMap[100].Refcnt = 0
	//ConversationIdMap[100].Idtype = [4]uint8{}
	OnlyForTestTeardown(t)

}

func OnlyForRxMachineCreateValidDRCPPacket() *layers.DRCP {

	phash := md5.New()

	for i := 0; i < MAX_CONVERSATION_IDS; i++ {
		buf := new(bytes.Buffer)
		// network byte order
		binary.Write(buf, binary.BigEndian, []uint16{uint16(i)})
		phash.Write(buf.Bytes())

	}

	ghash := md5.New()
	for i := float64(0); i < MAX_CONVERSATION_IDS; i++ {

		// we are only provisioning vlan 100 in the syatem
		if i == 100 {
			buf := new(bytes.Buffer)
			// network byte order
			binary.Write(buf, binary.BigEndian, []uint8{2, 1, uint8(uint16(i) >> 8 & 0xff), uint8(uint16(i) & 0xff)})
			ghash.Write(buf.Bytes())
		} else {
			buf := new(bytes.Buffer)
			// network byte order
			binary.Write(buf, binary.BigEndian, []uint16{uint16(i)})
			ghash.Write(buf.Bytes())

		}
	}
	portdigest := phash.Sum(nil)
	gatewaydigest := ghash.Sum(nil)

	drcp := &layers.DRCP{
		PortalInfo: layers.DRCPPortalInfoTlv{
			TlvTypeLength:  layers.DRCPTLVTypePortalInfo | layers.DRCPTLVPortalInfoLength,
			AggPriority:    128,
			AggId:          [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x64},
			PortalPriority: 128,
			PortalAddr:     [6]uint8{0x00, 0x00, 0xDE, 0xAD, 0xBE, 0xEF},
		},
		PortalConfigInfo: layers.DRCPPortalConfigurationInfoTlv{
			TlvTypeLength:    layers.DRCPTLVTypePortalConfigInfo | layers.DRCPTLVPortalConfigurationInfoLength,
			TopologyState:    layers.DRCPTopologyState(0x6),
			OperAggKey:       200,
			PortAlgorithm:    [4]uint8{0x00, 0x80, 0xC2, 0x01}, // C-VID
			GatewayAlgorithm: [4]uint8{0x00, 0x80, 0xC2, 0x01},
			PortDigest: [16]uint8{
				portdigest[0], portdigest[1], portdigest[2], portdigest[3],
				portdigest[4], portdigest[5], portdigest[6], portdigest[7],
				portdigest[8], portdigest[9], portdigest[10], portdigest[11],
				portdigest[12], portdigest[13], portdigest[14], portdigest[15],
			},
			GatewayDigest: [16]uint8{
				gatewaydigest[0], gatewaydigest[1], gatewaydigest[2], gatewaydigest[3],
				gatewaydigest[4], gatewaydigest[5], gatewaydigest[6], gatewaydigest[7],
				gatewaydigest[8], gatewaydigest[9], gatewaydigest[10], gatewaydigest[11],
				gatewaydigest[12], gatewaydigest[13], gatewaydigest[14], gatewaydigest[15],
			},
		},
		State: layers.DRCPStateTlv{
			TlvTypeLength: layers.DRCPTLVTypeDRCPState | layers.DRCPTLVStateLength,
			State:         layers.DRCPState(1 << layers.DRCPStateHomeGatewayBit),
		},
		HomePortsInfo: layers.DRCPHomePortsInfoTlv{
			TlvTypeLength:     layers.DRCPTLVTypeHomePortsInfo | layers.DRCPTlvTypeLength(8),
			AdminAggKey:       100,
			OperPartnerAggKey: 100,
			ActiveHomePorts:   []uint32{uint32(aggport3), uint32(aggport4)},
		},
		//NeighborPortsInfo:                  DRCPNeighborPortsInfoTlv{},
		HomeGatewayVector: layers.DRCPHomeGatewayVectorTlv{
			TlvTypeLength: layers.DRCPTLVTypeHomeGatewayVector | layers.DRCPTLVHomeGatewayVectorLength_2,
			Sequence:      1,
			Vector:        make([]uint8, 512),
		},
		//NeighborGatewayVector: DRCPGatewaNeighborGatewayVector	//	TlvTypeLength: layers.DRCPTLVTypeGatewayVectorEntry | layers.DRCPTLVGatewayVectorEntryLength,
		//	Sequence:      1,
		//},
		//TwoPNeighborGatewayVectorsationVector: DRCP2PGatNeighborGatewayVectorctorTlv{},
		//TwoPortalPortConversationVector:    DRCP2PPortConversationVectorTlv{},
		NetworkIPLMethod: layers.DRCPNetworkIPLSharingMethodTlv{
			TlvTypeLength: layers.DRCPTLVNetworkIPLSharingMethod | layers.DRCPTLVNetworkIPLSharingMethodLength,
			Method:        [4]uint8{0x00, 0x80, 0xC2, 0x1},
		},
		//NetworkIPLEncapsulation:            DRCPNetworkIPLSharingEncapsulationTlv{},
	}
	return drcp
}

func OnlyForRxMachineTestSetupCreateAggGroup(aggId uint32) *lacp.LaAggregator {
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

func RxMachineTestSetup() {
	OnlyForRxMachineTestSetup()
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
func RxMachineTestTeardown(t *testing.T) {

	OnlyForRxMachineTestTeardown(t)
	delete(utils.PortConfigMap, ipplink1)
	delete(utils.PortConfigMap, aggport1)
	delete(utils.PortConfigMap, aggport2)
	delete(utils.PortConfigMap, ipplink2)
	delete(utils.PortConfigMap, aggport3)
	delete(utils.PortConfigMap, aggport4)
}

func TestRxMachineRxValidDRCPDUNeighborPkt(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// create packet which will be sent for this test
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

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

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateCurrent {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was not set when it should be set", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	if ipp.DifferPortal {
		t.Error("ERROR packet portal info should agree with local since they are provisioned the same")
	}
	if ipp.DifferPortalReason != "" {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}
	/*
			// TODO when gateway sync code is fixed this should be uncommented
			eventReceived := false
			go func(evrx *bool) {
				for i := 0; i < 10 && !*evrx; i++ {
					time.Sleep(time.Second * 1)
				}
				if !eventReceived {
					ipp.TxMachineFsm.TxmEvents <- utils.MachineEvent{
						E:   fsm.Event(0),
						Src: "RX MACHINE: FORCE TEST FAIL",
					}
				}
			}(&eventReceived)

			evt := <-ipp.TxMachineFsm.TxmEvents
			if evt.E != TxmEventNtt {
				t.Error("ERROR Invalid event received", evt.E)
			}

		eventReceived = true
	*/
	//ipp.RxMachineFsm.Stop()
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxValidDRCPDUNeighborPktThenTimeout(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
		DrniNeighborAdminDRCPState:        "00010000",
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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// create packet which will be sent for this test
	drcp := OnlyForRxMachineCreateValidDRCPPacket()
	drcp.State.State.SetState(layers.DRCPStateDRCPTimeout)
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

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateCurrent {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was not set when it should be set", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	if ipp.DifferPortal {
		t.Error("ERROR packet portal info should agree with local since they are provisioned the same")
	}
	if ipp.DifferPortalReason != "" {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	/*
			// TODO when gateway sync code is fixed this should be uncommented

		eventReceived := false
		go func(evrx *bool) {
			for i := 0; i < 10 && !*evrx; i++ {
				time.Sleep(time.Second * 1)
			}
			if !eventReceived {
				ipp.TxMachineFsm.TxmEvents <- utils.MachineEvent{
					E:   fsm.Event(0),
					Src: "RX MACHINE: FORCE TEST FAIL",
				}
			}
		}(&eventReceived)

		evt := <-ipp.TxMachineFsm.TxmEvents
		if evt.E != TxmEventNtt {
			t.Error("ERROR Invalid event received", evt.E)
		}

		eventReceived = true
	*/
	// delay to allow for expire
	waitchan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 4)
		waitchan <- true
	}()

	<-waitchan

	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateExpired {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}

	go func() {
		time.Sleep(time.Second * 4)
		waitchan <- true
	}()

	<-waitchan

	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDefaulted {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}

	if !dr.DRFHomeOperDRCPState.GetState(layers.DRCPStateExpired) {
		t.Error("ERROR Oper STate should have Expired Set because we timed out twice")
	}

	//ipp.RxMachineFsm.Stop()
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferAggregatorPriority(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	drcp.PortalInfo.AggPriority = 256

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DifferPortal {
		t.Error("ERROR packet portal info should not agree with local since they are provisioned differently")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Neighbor Aggregator Priority") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferAggregatorAddr(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	drcp.PortalInfo.AggId = [6]uint8{0x00, 0x00, 0x00, 0x11, 0x00, 0x64}

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DifferPortal {
		t.Error("ERROR packet portal info should not agree with local since they are provisioned differently")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Neighbor Aggregator Id") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferOperAggregatorKey(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// In addition if the 14 least significant bits of DRF_Neighbor_Oper_Aggregator_Key ==
	// the 14 least significant bits of DRF_Home_Oper_Aggregator_Key
	drcp.PortalConfigInfo.OperAggKey = 0x8000 | 200

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Config info agrees, now operAgg needs to be negotiated
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateCurrent {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if !dr.ChangePortal {
		t.Error("ERROR packet portal info should not agree with local since they are provisioned differently")
	}
	//if !strings.Contains(ipp.DifferPortalReason, "Oper Aggregator Key") {
	//	t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	//}

	eventReceived := false
	go func(dr *DistributedRelay, evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !eventReceived {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "RX MACHINE: FORCE TEST FAIL",
			}
		}
	}(dr, &eventReceived)
	// event sent from recordDefaultDRCPDU
	evt := <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERROR Invalid event received", evt.E)
	}

	eventReceived = true

	// TEST now send a different oper key
	// case from above: Otherwise
	drcp.PortalConfigInfo.OperAggKey = 1000

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if !ipp.DifferConfPortal {
		t.Error("ERROR packet portal info should not agree with local since they are provisioned differently")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Oper Aggregator Key") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferThreeSystemPortalDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Further, if the variable Differ_Conf_Portal is set to FALSE and one or more of the comparisons
	// Drni_Neighbor_Three_System_Portal to Drni_Three_System_Portal, or
	// DRF_Neighbor_Gateway_Algorithm to DRF_Home_Gateway_Algorithm differ
	drcp.PortalConfigInfo.TopologyState.SetState(layers.DRCPTopologyState3SystemPortal, 1)

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if !ipp.DifferGatewayDigest {
		t.Error("ERROR Differ Gateway Digest is not set as it should be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Three System Portal") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	eventReceived := false
	go func(dr *DistributedRelay, evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !eventReceived {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "RX MACHINE: FORCE TEST FAIL",
			}
		}
	}(dr, &eventReceived)
	// event sent from recordDefaultDRCPDU
	evt := <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERROR Invalid event received", evt.E)
	}

	eventReceived = true
	/*
		eventReceived = false
		go func(dr *DistributedRelay, evrx *bool) {
			for i := 0; i < 10 && !*evrx; i++ {
				time.Sleep(time.Second * 1)
			}
			if !eventReceived {
				dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
					E:   fsm.Event(0),
					Src: "RX MACHINE: FORCE TEST FAIL",
				}
			}
		}(dr, &eventReceived)
		// event sent based on oper aggregator key changed
		evt = <-dr.PsMachineFsm.PsmEvents
		if evt.E != PsmEventChangePortal {
			t.Error("ERROR Invalid event received", evt.E)
		}
		eventReceived = true
	*/
	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferNeighborPortalSystemNumDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Check Portal System and Neighbor Portal System Number differ
	// Neighbor set as part of IPP Id
	drcp.PortalConfigInfo.TopologyState.SetState(layers.DRCPTopologyStatePortalSystemNum, 1)
	drcp.PortalConfigInfo.TopologyState.SetState(layers.DRCPTopologyStateNeighborConfPortalSystemNumber, 2)

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// lets check some settings on the ipp
	if ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateIPPActivity) {
		t.Error("ERROR Neighbor_Oper_DRCP_State IPP_Activity was set when it should be cleared", ipp.DRFNeighborOperDRCPState)
	}
	if !ipp.DRFNeighborOperDRCPState.GetState(layers.DRCPStateDRCPTimeout) {
		t.Error("ERROR Neighbor_Oper_DRCP_State DRCP Timeout was set to LONG when it should be SHORT", ipp.DRFNeighborOperDRCPState)
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if !ipp.DifferConfPortalSystemNumber {
		t.Error("ERROR DifferConfPortalSystemNumber is not set as it should be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Portal System Number") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	eventReceived := false
	go func(dr *DistributedRelay, evrx *bool) {
		for i := 0; i < 10 && !*evrx; i++ {
			time.Sleep(time.Second * 1)
		}
		if !eventReceived {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   fsm.Event(0),
				Src: "RX MACHINE: FORCE TEST FAIL",
			}
		}
	}(dr, &eventReceived)
	// event sent from recordDefaultDRCPDU
	evt := <-dr.PsMachineFsm.PsmEvents
	if evt.E != PsmEventChangePortal {
		t.Error("ERROR Invalid event received", evt.E)
	}

	eventReceived = true

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferGatewayAlgorithmDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Comparison differ
	drcp.PortalConfigInfo.GatewayAlgorithm = [4]uint8{0x00, 0x80, 0xC2, 0x4}

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if ipp.DifferConfPortalSystemNumber {
		t.Error("ERROR DifferConfPortalSystemNumber is set as it should not be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Gateway Algorithm") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

func TestRxMachineRxPktDRCPDUNeighborPortalInfoDifferPortAlgorithmDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Comparison differ
	drcp.PortalConfigInfo.PortAlgorithm = [4]uint8{0x00, 0x80, 0xC2, 0x2}

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateCurrent {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// only the upper two bits of operaggkey is different this does not mean that the rx fails
	if ipp.DifferConfPortalSystemNumber {
		t.Error("ERROR DifferConfPortalSystemNumber is set as it should not be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Port Algorithm") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

// TODO investigate why this test is hanging when running with full test suite
func xTestRxMachineRxPktDRCPDUNeighborPortalInfoDifferGatewayDigestDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Comparison differ
	drcp.PortalConfigInfo.GatewayDigest = [16]uint8{
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
	}

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// gateway digest is different
	if !ipp.DifferGatewayDigest {
		t.Error("ERROR DifferGatewayDigest is not set as it should be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Conversation Gateway List Digest") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}

// TODO investigate why this test is hanging when running with full test suite
func xTestRxMachineRxPktDRCPDUNeighborPortalInfoDifferPortDigestDiff(t *testing.T) {

	RxMachineTestSetup()
	a := OnlyForRxMachineTestSetupCreateAggGroup(200)

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
	// just create instance not starting any state machines
	dr := NewDistributedRelay(cfg)
	dr.a = a
	a.DrniName = dr.DrniName
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

	// lets get the IPP
	ipp := dr.Ipplinks[0]

	// rx machine sends event to each of these machines according to figure 9-22
	DrcpAMachineFSMBuild(dr)
	DrcpGMachineFSMBuild(dr)
	DrcpPsMachineFSMBuild(dr)
	DrcpTxMachineFSMBuild(ipp)
	DrcpPtxMachineFSMBuild(ipp)

	// start RX MAIN
	ipp.DrcpRxMachineMain()
	// enable because aggregator was attached above
	ipp.DRCPEnabled = true

	responseChan := make(chan string)

	// Psm is expected to be in update state
	// initialize will set the dr default values
	// which will be matched against received packets
	dr.PsMachineFsm.DrcpPsMachinePortalSystemInitialize(*dr.PsMachineFsm.Machine, nil)

	ipp.RxMachineFsm.RxmEvents <- utils.MachineEvent{
		E:            RxmEventBegin,
		Src:          "RX MACHINE TEST",
		ResponseChan: responseChan,
	}

	<-responseChan

	// create packet
	drcp := OnlyForRxMachineCreateValidDRCPPacket()

	// TEST:
	// Comparison differ
	drcp.PortalConfigInfo.PortDigest = [16]uint8{
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
	}

	ipp.RxMachineFsm.RxmPktRxEvent <- RxDrcpPdu{
		pdu:          drcp,
		src:          "RX MACHINE TEST",
		responseChan: responseChan,
	}

	<-responseChan

	// Neighbor Admin values not correct, thus should discard as
	// neighbor info is not known yet
	if ipp.RxMachineFsm.Machine.Curr.CurrentState() != RxmStateDiscard {
		t.Error("ERROR Rx Machine is not in expected state from first received PDU actual:", RxmStateStrMap[ipp.RxMachineFsm.Machine.Curr.CurrentState()])
	}
	// gateway digest is different
	if !ipp.DifferPortDigest {
		t.Error("ERROR DifferPortDigest is not set as it should be")
	}
	if !strings.Contains(ipp.DifferPortalReason, "Conversation Port List Digest") {
		t.Error("ERROR Portal Difference Detected", ipp.DifferPortalReason)
	}

	lacp.DeleteLaAgg(a.AggId)
	dr.DeleteDistributedRelay()
	RxMachineTestTeardown(t)
}
