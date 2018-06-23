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
// drIpp.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// DRNI - Distributed Resilient Network Interconnect

var DRCPIppDB map[IppDbKey]*DRCPIpp
var DRCPIppDBList []*DRCPIpp

type IppDbKey struct {
	Name   string
	DrName string
}

// 802.1ax-2014 7.4.2.1.1
type DistributedRelayIPP struct {
	Name                         string
	Id                           uint32
	PortConversationPasses       [MAX_CONVERSATION_IDS]bool
	GatewayConversationDirection [MAX_CONVERSATION_IDS]bool
	AdminState                   bool
	OperState                    bool
	TimeOfLstOperChange          time.Time
}

// 802.1ax-2014 7.4.3.1.1
type DistributedRelayIPPCounters struct {
	StatId    uint32
	DRCPDUsRX uint32
	IllegalRX uint32
	DRCPDUsTX uint32
}

// 802.1ax-2014 7.4.4.1.1
type DistributedRelayIPPDebug struct {
	InfoId             uint32
	DRCPRXState        string
	LastRXTime         time.Time
	DifferPortalReason string
}

type GatewayVectorEntry struct {
	Sequence uint32
	// MAX_CONVERSATION_IDS
	Vector []bool
}

type StateVectorInfo struct {
	mutex *sync.Mutex

	OpState bool
	// indexed by the received Home_Gateway_Sequence in
	// increasing sequence number order
	GatewayVector []GatewayVectorEntry
	PortIdList    []uint32
}

// 802.1ax-2014 9.4.9 Per IPP Intra-Portal Variables
type DRCPIntraPortal struct {
	CCTimeShared                 bool
	CCEncTagShared               bool
	DifferConfPortal             bool
	DifferConfPortalSystemNumber bool
	DifferGatewayDigest          bool
	DifferPortDigest             bool
	DifferPortal                 bool
	// range 1..3
	DRFHomeConfNeighborPortalSystemNumber uint8
	DRFHomeNetworkIPLIPLEncapDigest       Md5Digest
	DRFHomeNetworkIPLIPLNetEncapDigest    Md5Digest
	DRFHomeNetworkIPLSharingMethod        EncapMethod
	// defines for state can be found in "github.com/google/gopacket/layers"
	DRFNeighborAdminAggregatorKey            uint16
	DRFNeighborAggregatorId                  [6]uint8
	DRFNeighborAggregatorPriority            uint16
	DRFNeighborConversationGatewayListDigest Md5Digest
	DRFNeighborConversationPortListDigest    Md5Digest
	DRFNeighborGatewayAlgorithm              [4]uint8
	DRFNeighborGatewayConversationMask       [MAX_CONVERSATION_IDS]bool
	DRFNeighborGatewaySequence               uint32
	DRFNeighborNetworkIPLIPLEncapDigest      Md5Digest
	DRFNeighborNetworkIPLNetEncapDigest      Md5Digest
	DRFNeighborNetworkIPLSharingMethod       EncapMethod
	DRFNeighborOperAggregatorKey             uint16
	DRFNeighborOperPartnerAggregatorKey      uint16
	// defines for state can be found in "github.com/google/gopacket/layers"
	DRFNeighborOperDRCPState layers.DRCPState
	// range 1..3
	DRFNeighborConfPortalSystemNumber uint8
	DRFNeighborPortAlgorithm          [4]uint8
	// range 1..3
	DRFNeighborPortalSystemNumber            uint8
	DRFNeighborState                         StateVectorInfo
	DRFOtherNeighborAdminAggregatorKey       uint16
	DRFOtherNeighborGatewayConversationMask  [MAX_CONVERSATION_IDS]bool
	DRFOtherNeighborGatewaySequence          uint16
	DRFOtherNeighborOperPartnerAggregatorKey uint16
	DRFOtherNeighborState                    StateVectorInfo
	DRFRcvHomeGatewayConversationMask        [MAX_CONVERSATION_IDS]bool
	DRFRcvHomeGatewaySequence                uint32
	DRFRcvNeighborGatewayConversationMask    [MAX_CONVERSATION_IDS]bool
	DRFRcvNeighborGatewaySequence            uint16
	DRFRcvOtherGatewayConversationMask       [MAX_CONVERSATION_IDS]bool
	DRFRcvOtherGatewaySequence               uint16
	DrniNeighborCommonMethods                bool
	DrniNeighborGatewayConversation          [1024]uint8
	DrniNeighborPortConversation             [1024]uint8
	DrniNeighborONN                          bool
	DrniNeighborPortalAddr                   [6]uint8
	DrniNeighborPortalPriority               uint16
	DrniNeighborState                        [4]StateVectorInfo
	// This should always be false as we will not support 3 portal system initially
	DrniNeighborThreeSystemPortal        bool
	EnabledTimeShared                    bool
	EnabledEncTagShared                  bool
	IppOtherGatewayConversation          [MAX_CONVERSATION_IDS]uint8
	IppOtherPortConversationPortalSystem [MAX_CONVERSATION_IDS]uint8
	IppPortEnabled                       bool
	IppPortalSystemState                 [4]StateVectorInfo
	MissingRcvGatewayConVector           bool
	MissingRcvPortConVector              bool
	NTTDRCPDU                            bool
	ONN                                  bool

	// 9.4.10
	Begin                       bool
	DRCPEnabled                 bool
	GatewayConversationTransmit bool
	IppAllUpdate                bool
	IppGatewayUpdate            bool
	IppPortUpdate               bool
	PortConversationTransmit    bool

	// 9.3.4.3
	IppGatewayConversationPasses [MAX_CONVERSATION_IDS]bool
	IppPortconversationPasses    [MAX_CONVERSATION_IDS]bool
}

type DRCPIpp struct {
	DistributedRelayIPP
	DRCPIntraPortal
	DistributedRelayIPPCounters
	DistributedRelayIPPDebug

	// reference to the distributed relay object
	dr *DistributedRelay

	// sync creation and deletion
	wg sync.WaitGroup

	// handle used to tx packets to linux if
	handle *pcap.Handle

	// channel used to wait on response from distributed event send
	ippEvtResponseChan chan string

	// FSMs
	RxMachineFsm          *RxMachine
	PtxMachineFsm         *PtxMachine
	TxMachineFsm          *TxMachine
	NetIplShareMachineFsm *NetIplShareMachine
	IAMachineFsm          *IAMachine
	IGMachineFsm          *IGMachine
}

func NewDRCPIpp(id uint32, dr *DistributedRelay) *DRCPIpp {

	neighborPortalSystemNum := uint8(2)
	if dr.DrniPortalSystemNumber == 2 {
		neighborPortalSystemNum = 1
	}

	ipp := &DRCPIpp{
		DistributedRelayIPP: DistributedRelayIPP{
			Name:       utils.PortConfigMap[int32(id&0xffff)].Name,
			Id:         id & 0xffff,
			AdminState: true,
		},
		DRCPIntraPortal: DRCPIntraPortal{
			DRCPEnabled: true,
			// neighbor system id contained in the port id
			DRFHomeConfNeighborPortalSystemNumber: neighborPortalSystemNum,
			DRFHomeNetworkIPLSharingMethod:        dr.DrniEncapMethod,
			DRFNeighborState:                      StateVectorInfo{mutex: &sync.Mutex{}},
			DRFOtherNeighborState:                 StateVectorInfo{mutex: &sync.Mutex{}},
		},
		dr:                 dr,
		ippEvtResponseChan: make(chan string),
	}

	for i, _ := range ipp.DRCPIntraPortal.DrniNeighborState {
		ipp.DRCPIntraPortal.DrniNeighborState[i].mutex = &sync.Mutex{}
	}
	for i, _ := range ipp.IppPortalSystemState {
		ipp.IppPortalSystemState[i].mutex = &sync.Mutex{}
	}

	key := IppDbKey{
		Name:   ipp.Name,
		DrName: ipp.dr.DrniName,
	}

	// add port to port db
	DRCPIppDB[key] = ipp
	DRCPIppDBList = append(DRCPIppDBList, ipp)

	// check the link status
	for _, client := range utils.GetAsicDPluginList() {
		ipp.OperState = client.GetPortLinkStatus(int32(ipp.Id))
		ipp.IppPortEnabled = ipp.OperState
		ipp.LaIppLog(fmt.Sprintln("Initial IPP Link State", ipp.Name, ipp.IppPortEnabled))
	}

	// create the packet capture rule if it does not already exist
	ipp.SetupDRCPMacCapture(ipp.dr.DrniPortalPortProtocolIDA.String())

	ipp.LaIppLog(fmt.Sprintf("Created IPP port %+v\n", ipp))

	if ipp.OperState {
		ipp.CreateRxTx()
	}
	return ipp
}

func (p *DRCPIpp) CreateRxTx() {
	if p.handle == nil {
		handle, err := pcap.OpenLive(p.Name, 65536, true, 50*time.Millisecond)
		if err != nil {
			// failure here may be ok as this may be SIM
			if !strings.Contains(p.Name, "SIM") {
				p.LaIppLog(fmt.Sprintf("Error creating pcap OpenLive handle for port", p.Id, p.Name, err))
			}
			return
		}
		fmt.Println("Creating Listener for intf ", p.Name)
		p.handle = handle
		src := gopacket.NewPacketSource(p.handle, layers.LayerTypeEthernet)
		in := src.Packets()
		// start rx routine
		DrRxMain(uint16(p.Id), p.dr.DrniPortalAddr.String(), in)
		p.LaIppLog(fmt.Sprintf("Rx Main Started for ipp link port %s", p.Name))

		key := IppDbKey{
			Name:   p.Name,
			DrName: p.dr.DrniName,
		}

		// register the tx func
		DRGlobalSystem.DRSystemGlobalRegisterTxCallback(key, TxViaLinuxIf)
	}
}

func (p *DRCPIpp) DeleteRxTx() {

	key := IppDbKey{
		Name:   p.Name,
		DrName: p.dr.DrniName,
	}
	// De-register the tx function
	DRGlobalSystem.DRSystemGlobalDeRegisterTxCallback(key)

	// close rx/tx processing
	if p.handle != nil {
		p.handle.Close()
		p.LaIppLog(fmt.Sprintf("RX/TX handle closed for port", p.Id))

	}

}

//
func (p *DRCPIpp) DeleteDRCPIpp() {
	dr := p.dr
	// remove the packet capture rule if this is the last reference to it
	p.TeardownDRCPMacCapture(p.dr.DrniPortalPortProtocolIDA.String())

	// stop all state machines
	// if agg is attached but maybe we are just deleting the ipp link
	if dr.a != nil {
		p.Stop()
	}

	// cleanup the global tables hosting the port
	key := IppDbKey{
		Name:   p.Name,
		DrName: p.dr.DrniName,
	}
	// cleanup the tables
	if _, ok := DRCPIppDB[key]; ok {
		delete(DRCPIppDB, key)
		for i, delipp := range DRCPIppDBList {
			if delipp == p {
				DRCPIppDBList = append(DRCPIppDBList[:i], DRCPIppDBList[i+1:]...)
			}
		}
	}
}

// SetupDRCPMacCapture will create an pkt capture rule in the hw
func (p *DRCPIpp) SetupDRCPMacCapture(mac string) {
	intfref := utils.GetNameFromIfIndex(int32(p.Id))
	key := MacCaptureKey{
		mac:     mac,
		intfref: intfref,
	}

	if _, ok := MacCaptureCount[key]; ok {
		MacCaptureCount[key]++
	} else {
		MacCaptureCount[key] = 1
		for _, client := range utils.GetAsicDPluginList() {
			p.LaIppLog(fmt.Sprintf("Enabling Pkt Capture in HW IPP port %s with mac %s", p.Name, mac))
			client.EnablePacketReception(mac, 0, int32(p.Id))
		}
	}

}

// TeardownDRCPMacCapture will delete a pkt capture rule in the hw
func (p *DRCPIpp) TeardownDRCPMacCapture(mac string) {
	intfref := utils.GetNameFromIfIndex(int32(p.Id))
	key := MacCaptureKey{
		mac:     mac,
		intfref: intfref,
	}

	if _, ok := MacCaptureCount[key]; ok {
		MacCaptureCount[key]--
		if MacCaptureCount[key] == 0 {
			for _, client := range utils.GetAsicDPluginList() {
				p.LaIppLog(fmt.Sprintf("Disabling Pkt Capture in HW IPP port %s with mac %s", p.Name, mac))
				client.DisablePacketReception(mac, 0, int32(p.Id))
			}
		}
	}

}

// Stop the port services and state machines
func (p *DRCPIpp) Stop() {

	p.DeleteRxTx()
	// Stop the State Machines

	// Ptxm
	if p.PtxMachineFsm != nil {
		p.PtxMachineFsm.Stop()
		p.PtxMachineFsm = nil
	}
	// Rxm
	if p.RxMachineFsm != nil {
		p.RxMachineFsm.Stop()
		p.RxMachineFsm = nil
	}
	// Txm
	if p.TxMachineFsm != nil {
		p.TxMachineFsm.Stop()
		p.TxMachineFsm = nil
	}
	// NetIplShare
	if p.NetIplShareMachineFsm != nil {
		p.NetIplShareMachineFsm.Stop()
		p.NetIplShareMachineFsm = nil
	}
	// IAm
	if p.IAMachineFsm != nil {
		p.IAMachineFsm.Stop()
		p.IAMachineFsm = nil
	}
	// IGm
	if p.IGMachineFsm != nil {
		p.IGMachineFsm.Stop()
		p.IGMachineFsm = nil
	}
	// lets wait for all the State machines to have stopped
	p.wg.Wait()

	close(p.ippEvtResponseChan)
}

// BEGIN this will send the start event to the start the state machines
func (p *DRCPIpp) BEGIN(restart bool) {

	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	// there is a case in which we have only called
	// restart and called main functions outside
	// of this scope (TEST for example)
	//prevBegin := p.begin

	// System in being initalized
	//p.begin = true

	if !restart {
		// start all the State machines
		// Order here matters as Rx machine
		// will send event to Mux machine
		// thus machine must be up and
		// running first
		// Periodic Tx Machine
		p.DrcpPtxMachineMain()
		// Net/IPL Sharing Machine
		p.NetIplShareMachineMain()
		// IPP Aggregator machine
		p.DrcpIAMachineMain()
		// IPP Gateway Machine
		p.DrcpIGMachineMain()
		// Tx Machine
		p.TxMachineMain()
		// Rx Machine
		p.DrcpRxMachineMain()
	}

	// wait group used when stopping all the
	// State mahines associated with this port.
	// want to ensure that all routines are stopped
	// before proceeding with cleanup thus why not
	// create the wg as part of a BEGIN process
	// 1) Rx Machine
	// 2) Tx Machine
	// 3) Periodic Tx Machine
	// 4) Net/IPL Sharing Machine
	// 5) IPP Aggregator Machine
	// 6) IPP Gateway Machine

	// Rxm
	if p.RxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   RxmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// Txm
	if p.TxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   TxmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// Ptxm
	if p.PtxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   PtxmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// NetIplShare
	if p.NetIplShareMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.NetIplShareMachineFsm.NetIplSharemEvents)
		evt = append(evt, utils.MachineEvent{
			E:   NetIplSharemEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// IAm
	if p.IAMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.IAMachineFsm.IAmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   IAmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// IGm
	if p.IGMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.IGMachineFsm.IGmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   IGmEventBegin,
			Src: DRCPConfigModuleStr})
	}

	// call the begin event for each
	// distribute the port disable event to various machines
	p.DistributeMachineEvents(mEvtChan, evt, true)

}

// DrIppLinkUp distribute link up event
func (p *DRCPIpp) DrIppLinkUp() {

	p.CreateRxTx()

	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	p.IppPortEnabled = true

	if p.DRCPEnabled {
		mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   RxmEventNotIPPPortEnabled,
			Src: DRCPConfigModuleStr})

		mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   IGmEventBegin,
			Src: DRCPConfigModuleStr})

	}
	p.DistributeMachineEvents(mEvtChan, evt, false)

}

// DrIppLinkDown distributelink down event
func (p *DRCPIpp) DrIppLinkDown() {

	p.DeleteRxTx()

	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	p.IppPortEnabled = false

	mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
	evt = append(evt, utils.MachineEvent{
		E:   RxmEventNotIPPPortEnabled,
		Src: DRCPConfigModuleStr})

	mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
	evt = append(evt, utils.MachineEvent{
		E:   IGmEventBegin,
		Src: DRCPConfigModuleStr})

	p.DistributeMachineEvents(mEvtChan, evt, false)

}

// DistributeMachineEvents will distribute the events in parrallel
// to each machine
func (p *DRCPIpp) DistributeMachineEvents(mec []chan utils.MachineEvent, e []utils.MachineEvent, waitForResponse bool) {

	length := len(mec)
	if len(mec) != len(e) {
		p.LaIppLog("LADR: Distributing of events failed")
		return
	}

	// send all begin events to each machine in parrallel
	for j := 0; j < length; j++ {
		go func(port *DRCPIpp, w bool, idx int, machineEventChannel []chan utils.MachineEvent, event []utils.MachineEvent) {
			if w {
				event[idx].ResponseChan = p.ippEvtResponseChan
			}
			event[idx].Src = DRCPConfigModuleStr
			machineEventChannel[idx] <- event[idx]
		}(p, waitForResponse, j, mec, e)
	}

	if waitForResponse {
		i := 0
		// lets wait for all the machines to respond
		for {
			select {
			case mStr := <-p.ippEvtResponseChan:
				i++
				p.LaIppLog(strings.Join([]string{"LADRIPP:", mStr, "response received"}, " "))
				//fmt.Println("LAPORT: Waiting for response Delayed", length, "curr", i, time.Now())
				if i >= length {
					// 10/24/15 fixed hack by sending response after Machine.ProcessEvent
					// HACK, found that port is pre-empting the State machine callback return
					// lets delay for a short period to allow for event to be received
					// and other routines to process their events
					/*
						if p.logEna {
							time.Sleep(time.Millisecond * 3)
						} else {
							time.Sleep(time.Millisecond * 1)
						}
					*/
					return
				}
			}
		}
	}
}

// NotifyNTTDRCPUDChange
func (p *DRCPIpp) NotifyNTTDRCPUDChange(src string, oldval, newval bool) {
	if oldval != newval &&
		newval &&
		p.TxMachineFsm != nil {
		p.TxMachineFsm.TxmEvents <- utils.MachineEvent{
			E:   TxmEventNtt,
			Src: src,
		}
	}
}

// ReportToManagement send events for various reason to infor management of something
// is wrong.
func (p *DRCPIpp) reportToManagement() {

	if p.DifferPortalReason != "" {
		p.LaIppLog(fmt.Sprintf("Report Failure to Management: %s", p.DifferPortalReason))
		// TODO send event
	}
}

// DRFindPortByKey find ipp port by key
func DRFindPortByKey(key IppDbKey, p **DRCPIpp) bool {
	if ipp, ok := DRCPIppDB[key]; ok {
		*p = ipp
		return true
	}
	return false
}

// updateGatewayVector will update the vector, indexed by the received
// Gateway_Sequence in increasing sequence number order
func (nsi *StateVectorInfo) updateGatewayVector(sequence uint32, vector []bool) {

	//fmt.Printf("updateGatewayVector: GatewayVector sequence %d vector[100]=%t\n", sequence, vector[100])

	if len(nsi.GatewayVector) > 0 {

		nsi.OpState = true
		if nsi.GatewayVector[0].Sequence != sequence {
			obj := GatewayVectorEntry{
				Sequence: sequence,
				Vector:   make([]bool, 4096)}
			// save off the vector information
			for j, val := range vector {
				obj.Vector[j] = val
			}
			// insert sequence/vecotor at front of list
			nsi.GatewayVector = append([]GatewayVectorEntry{obj}, nsi.GatewayVector...)

			// lets only store the last 5 sequences
			if len(nsi.GatewayVector) > 5 {
				nsi.GatewayVector = append(nsi.GatewayVector[:0], nsi.GatewayVector[0:5]...)
			}
			//fmt.Printf("updateGatewayVector: prepend vector[100] %t\n", nsi.GatewayVector[0].Vector[100])
		}
	} else {
		nsi.GatewayVector = make([]GatewayVectorEntry, 1)
		tmp := GatewayVectorEntry{
			Sequence: sequence,
			Vector:   make([]bool, 4096),
		}
		for j, val := range vector {
			tmp.Vector[j] = val
		}
		nsi.OpState = true
		nsi.GatewayVector[0] = tmp
		//fmt.Printf("updateGatewayVector: new vector[100] %t\n", nsi.GatewayVector[0].Vector[100])
	}
}

// getNeighborVectorGatwaySequenceIndex get the index for the entry whos
// sequence number is equal.
func (nsi *StateVectorInfo) getNeighborVectorGatwaySequenceIndex(sequence uint32, vector []bool) int32 {

	if len(nsi.GatewayVector) > 0 {
		for i, seqVector := range nsi.GatewayVector {
			if seqVector.Sequence == sequence {
				return int32(i)
			}
		}
	}
	return -1
}

func (nsi *StateVectorInfo) getGatewayVectorByIndex(index int32) *GatewayVectorEntry {
	length := len(nsi.GatewayVector)
	if length > 0 && index < int32(length) {
		return &nsi.GatewayVector[index]
	}
	return nil
}
