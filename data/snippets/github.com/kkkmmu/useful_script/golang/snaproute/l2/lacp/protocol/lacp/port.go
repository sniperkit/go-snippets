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

// port
package lacp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type PortProperties struct {
	Mac    net.HardwareAddr
	Speed  int
	Duplex int
	Mtu    int
}

type LacpPortInfo struct {
	System   LacpSystem
	Key      uint16
	Port_pri uint16
	port     uint16
	State    uint8
}

type LacpCounters struct {
	LacpInPkts        uint64
	LacpOutPkts       uint64
	LacpRxErrors      uint64
	LacpTxErrors      uint64
	LacpUnknownErrors uint64
	LacpErrors        uint64
}

type AggPortStatus struct {
	// GET
	AggPortId int
	// GET
	AggPortActorSystemId [6]uint8
	// GET
	AggPortActorOperKey uint16
	// GET
	AggPortPartnerOperSystemPriority uint16
	// GET
	AggPortPartnerOperSystemId [6]uint8
	// GET
	AggPortPartnerOperKey uint16
	// GET
	AggPortSelectedAggID int
	// GET
	AggPortAttachedAggID int
	// GET
	AggPortActorPort int
	// GET
	AggPortPartnerOperPort int
	// GET
	AggPortPartnerOperPortPriority uint8
	// GET
	AggPortActorOperState uint8
	// GET
	AggPortPartnerOperState uint8
	// GET
	AggPortAggregateOrIndividual bool
	// GET
	AggPortOperConversationPasses bool
	// GET
	AggPortOperConversationCollected bool
	// GET
	AggPortStats AggPortStatsObject
	// GET
	AggPortDebug AggPortDebugInformationObject
}

type AggregatorPortObject struct {
	// GET-SET
	Config AggPortConfig
	// GET
	Status AggPortStatus
	/*
		// GET
		AggPortId int
		// GET-SET
		AggPortActorSystemPriority uint16
		// GET
		AggPortActorSystemId [6]uint8
		// GET-SET
		AggPortActorAdminKey uint16
		// GET
		AggPortActorOperKey uint16
		// GET-SET
		AggPortPartnerAdminSystemPriority uint16
		// GET
		AggPortPartnerOperSystemPriority uint16
		// GET-SET
		AggPortPartnerAdminSystemId [6]uint8
		// GET
		AggPortPartnerOperSystemId [6]uint8
		// GET-SET
		AggPortPartnerAdminKey uint16
		// GET
		AggPortPartnerOperKey uint16
		// GET
		AggPortSelectedAggID int
		// GET
		AggPortAttachedAggID int
		// GET
		AggPortActorPort int
		// GET-SET
		AggPortActorPortPriority uint8
		// GET-SET
		AggPortPartnerAdminPort int
		// GET
		AggPortPartnerOperPort int
		// GET-SET
		AggPortPartnerAdminPortPriority uint8
		// GET
		AggPortPartnerOperPortPriority uint8
		// GET-SET
		AggPortActorAdminState uint8
		// GET
		AggPortActorOperState uint8
		// GET-SET
		AggPortPartnerAdminState uint8
		// GET
		AggPortPartnerOperState uint8
		// GET
		AggPortAggregateOrIndividual bool
		// GET
		AggPortOperConversationPasses bool
		// GET
		AggPortOperConversationCollected bool
		// GET-SET
		AggPortLinkNumberID int
		// GET-SET
		AggPortPartnerAdminLInkNumberID int
		// GET-SET
		AggPortWTRTime int
		// GET-SET
		AggPortProtocolDA [6]uint8
		// GET
		AggPortStats AggPortStatsObject
		// GET
		AggPortDebug AggPortDebugInformationObject
	*/

	Internal AggInternalData
}

type AggInternalData struct {
	// Linux Interface Name
	AggNameStr string
}

// GET
type AggPortStatsObject struct {
	AggPortStatsID                   uint64
	AggPortStatsLACPDUsRx            uint64
	AggPortStatsMarkerPDUsRx         uint64
	AggPortStatsMarkerResponsePDUsRx uint64
	AggPortStatsUnknownRx            uint64
	AggPortStatsIllegalRx            uint64
	AggPortStatsLACPDUsTx            uint64
	AggPortStatsMarkerPDUsTx         uint64
	AggPortStatsMarkerResponsePDUsTx uint64
	AggPortStateMissMatchInfoRx      uint64
}

//GET
type AggPortDebugInformationObject struct {
	// same as AggregationPort
	AggPortDebugInformationID int
	// enum
	AggPortDebugRxState    int
	AggPortDebugLastRxTime int
	// enum
	AggPortDebugMuxState                   int
	AggPortDebugMuxReason                  string
	AggPortDebugActorChurnState            int
	AggPortDebugPartnerChurnState          int
	AggPortDebugActorChurnPrevCnt          int
	AggPortDebugActorChurnCount            int
	AggPortDebugPartnerChurnPrevCount      int
	AggPortDebugPartnerChurnCount          int
	AggPortDebugActorSyncTransitionCount   int
	AggPortDebugPartnerSyncTransitionCount int
	// TODO
	AggPortDebugActorChangeCount     int
	AggPortDebugPartnerChangeCount   int
	AggPortDebugActorCDSChurnState   int
	AggPortDebugPartnerCDSChurnState int
	AggPortDebugActorCDSChurnCount   int
	AggPortDebugPartnerCDSChurnCount int
}

// 802.1ax Section 6.4.7
// Port attributes associated with aggregator
type LaAggPort struct {
	// 802.1ax-2014 Section 6.3.4:
	// Link Aggregation Control uses a Port Identifier (Port ID), comprising
	// the concatenation of a Port Priority (7.3.2.1.15) and a Port Number
	// (7.3.2.1.14), to identify the Aggregation Port....
	// The most significant and second most significant octets are the first
	// and second most significant octets of the Port Priority, respectively.
	// The third and fourth most significant octets are the first and second
	// most significant octets of the Port Number, respectively.
	portId int

	// string id of port
	IntfNum string

	// Key
	Key uint16

	// used to form portId
	PortNum      uint16
	portPriority uint16

	AggId int

	// Once selected reference to agg group will be made
	AggAttached *LaAggregator
	aggSelected int
	// unable to aggregate with other links in an agg
	operIndividual int
	lacpEnabled    bool
	// TRUE - Aggregation port is operable (MAC_Operational == True)
	// FALSE - otherwise
	PortEnabled  bool
	portMoved    bool
	begin        bool
	actorChurn   bool
	partnerChurn bool
	readyN       bool

	macProperties PortProperties

	// determine whether a port is up or down
	LinkOperStatus bool

	// administrative values for State described in 6.4.2.3
	ActorAdmin   LacpPortInfo
	ActorOper    LacpPortInfo
	partnerAdmin LacpPortInfo
	PartnerOper  LacpPortInfo

	// State machines
	RxMachineFsm       *LacpRxMachine
	PtxMachineFsm      *LacpPtxMachine
	TxMachineFsm       *LacpTxMachine
	CdMachineFsm       *LacpActorCdMachine
	PCdMachineFsm      *LacpPartnerCdMachine
	MuxMachineFsm      *LacpMuxMachine
	MarkerResponderFsm *LampMarkerResponderMachine

	// Counters
	LacpCounter AggPortStatsObject

	// GET
	AggPortDebug AggPortDebugInformationObject

	// Distributed Relay reference name
	DrniName   string
	DrniSynced bool

	// on configuration changes need to inform all State
	// machines and wait for a response
	portChan chan string
	logEna   bool
	wg       sync.WaitGroup

	// handle used to tx packets to linux if
	handle *pcap.Handle

	// Version 2
	partnerLacpPduVersionNumber int
	enableLongPduXmit           bool
	// packet is 1 byte, but spec says save as int.
	// going to save as byte
	partnerVersion uint8

	sysId net.HardwareAddr
}

// find a port from the global map table by PortNum
func LaFindPortById(pId uint16, port **LaAggPort) bool {
	for _, sgi := range LacpSysGlobalInfoGet() {
		for _, p := range sgi.LacpSysGlobalAggPortListGet() {
			if p.PortNum == pId {
				*port = p
				return true
			}
		}
	}
	return false
}

func LaConvertPortAndPriToPortId(pId uint16, prio uint16) int {
	return int(pId | prio<<16)
}

func LaGetPortNext(port **LaAggPort) bool {
	returnNext := false
	for _, sgi := range LacpSysGlobalInfoGet() {
		for _, p := range sgi.LacpSysGlobalAggPortListGet() {
			if *port == nil {
				// first port
				*port = p
				return true
			} else if (*port).PortNum == p.PortNum {
				// found port, lets return the next port
				returnNext = true
			} else if returnNext {
				// next port
				*port = p
				return true
			}
		}
	}
	*port = nil
	return false
}

// find a port from the global map table by PortNum
func LaFindPortByPortId(portId int, port **LaAggPort) bool {
	for _, sgi := range LacpSysGlobalInfoGet() {
		for _, p := range sgi.LacpSysGlobalAggPortListGet() {
			if p.portId == portId {
				*port = p
				return true
			}
		}
	}
	return false
}

// LaFindPortByKey will find a port form the global map table by Key
// index value should input 0 for the first value
func LaFindPortByKey(Key uint16, index *int, port **LaAggPort) bool {
	var i int
	for _, sgi := range LacpSysGlobalInfoGet() {
		i = *index
		aggPortList := sgi.LacpSysGlobalAggPortListGet()
		l := len(aggPortList)
		for _, p := range aggPortList {
			if i < l {
				if p.Key == Key {
					*port = p
					*index++
					return true
				}
			}
			i++
		}
	}
	*index = -1
	return false
}

// NewLaAggPort
// Allocate a new lag port, creating appropriate timers
func NewLaAggPort(config *LaAggPortConfig) *LaAggPort {

	// Lets see if the agg exists and add this port to this config
	// otherwise lets use the default
	var a *LaAggregator
	var sysId LacpSystem
	if LaFindAggByKey(config.Key, &a) {
		mac, _ := net.ParseMAC(a.Config.SystemIdMac)
		sysId.Actor_System = convertNetHwAddressToSysIdKey(mac)
		sysId.Actor_System_priority = a.Config.SystemPriority
	}
	sgi := LacpSysGlobalInfoByIdGet(sysId)
	portcfg, ok := utils.PortConfigMap[int32(config.Id)]
	if !ok {
		utils.GlobalLogger.Err(fmt.Sprintln("ERROR could not find port in map", config.Id, utils.PortConfigMap))
		return nil
	}
	config.IntfId = portcfg.Name

	p := &LaAggPort{
		portId:       LaConvertPortAndPriToPortId(config.Id, config.Prio),
		PortNum:      uint16(config.Id),
		portPriority: config.Prio,
		IntfNum:      config.IntfId,
		Key:          config.Key,
		AggId:        0, // this should be set on config AddLaAggPortToAgg
		aggSelected:  LacpAggUnSelected,
		begin:        false,
		portMoved:    false,
		lacpEnabled:  false,
		PortEnabled:  config.Enable,
		macProperties: PortProperties{Mac: config.Properties.Mac,
			Speed:  config.Properties.Speed,
			Duplex: config.Properties.Duplex,
			Mtu:    config.Properties.Mtu},
		logEna:       true,
		portChan:     make(chan string),
		AggPortDebug: AggPortDebugInformationObject{AggPortDebugInformationID: int(config.Id)},
		DrniName:     "",
	}

	// register the events
	utils.CreateEventMap(int32(p.PortNum))
	utils.ProcessLacpPortOperStateDown(int32(p.PortNum))
	RegisterLaPortUpCb("event_"+p.IntfNum, utils.ProcessLacpPortOperStateUp)
	RegisterLaPortDownCb("event_"+p.IntfNum, utils.ProcessLacpPortOperStateDown)

	// default actor admin
	//fmt.Println(config.sysId, gLacpSysGlobalInfo[config.sysId])
	p.ActorAdmin.State = sgi.ActorStateDefaultParams.State
	if a != nil &&
		a.DrniName != "" {
		// DR is the owner of the systemid and priority and it set the aggregator
		// mac and priority
		p.LaPortLog(fmt.Sprintf("Aggregator  %s is an MLAG owned by DR %s, thus using DR SystemId %+v Priority %d", a.AggName, a.DrniName, a.AggMacAddr, a.AggPriority))
		p.ActorAdmin.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(sgi.SystemDefaultParams.Actor_System))
		p.ActorAdmin.System.LacpSystemActorSystemPrioritySet(sgi.SystemDefaultParams.Actor_System_priority)
		p.ActorOper.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(a.AggMacAddr))
		p.ActorOper.System.LacpSystemActorSystemPrioritySet(a.AggPriority)

	} else {
		if a != nil {
			p.LaPortLog(fmt.Sprintf("Aggregator %s is a LAG owned by System, thus using SystemId %+v Priority %d", a.AggName, sgi.SystemDefaultParams.Actor_System, sgi.SystemDefaultParams.Actor_System_priority))
		}
		p.ActorAdmin.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(sgi.SystemDefaultParams.Actor_System))
		p.ActorAdmin.System.LacpSystemActorSystemPrioritySet(sgi.SystemDefaultParams.Actor_System_priority)
		p.ActorOper.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(sgi.SystemDefaultParams.Actor_System))
		p.ActorOper.System.LacpSystemActorSystemPrioritySet(sgi.SystemDefaultParams.Actor_System_priority)
	}
	p.ActorAdmin.Key = p.Key
	p.ActorAdmin.port = p.PortNum
	p.ActorAdmin.Port_pri = p.portPriority

	p.ActorOper.Key = p.ActorAdmin.Key
	p.ActorOper.port = p.ActorAdmin.port
	p.ActorOper.Port_pri = p.ActorAdmin.Port_pri
	p.ActorOper.State = p.ActorAdmin.State

	// port should be forced to unselected until the
	// DRNI has negotiated the Key
	if a != nil &&
		a.DrniName != "" {
		p.ActorOper.Key = 0
	}

	// default partner admin
	p.partnerAdmin.State = sgi.PartnerStateDefaultParams.State
	// default partner oper same as admin
	p.PartnerOper = p.partnerAdmin

	if config.Mode != LacpModeOn {
		p.lacpEnabled = true
	}

	// add port to port map
	sgi.PortMap[PortIdKey{Name: p.IntfNum,
		Id: p.PortNum}] = p

	sgi.PortList = append(sgi.PortList, p)

	return p
}

func (p *LaAggPort) CreateRxTx() {
	if p.handle == nil {
		var a *LaAggregator
		var sysId LacpSystem
		if LaFindAggById(p.AggId, &a) {
			mac, _ := net.ParseMAC(a.Config.SystemIdMac)
			sysId.Actor_System = convertNetHwAddressToSysIdKey(mac)
			sysId.Actor_System_priority = a.Config.SystemPriority

			sgi := LacpSysGlobalInfoByIdGet(sysId)

			handle, err := pcap.OpenLive(p.IntfNum, 65536, true, 50*time.Millisecond)
			if err != nil {
				// failure here may be ok as this may be SIM
				if !strings.Contains(p.IntfNum, "SIM") {
					fmt.Println("Error creating pcap OpenLive handle for port", p.PortNum, p.IntfNum, err)
				}
				return
			}
			filter := fmt.Sprintf(`ether dst 01:80:C2:00:00:02`)
			err = handle.SetBPFFilter(filter)
			if err != nil {
				p.LaPortLog(fmt.Sprintln("Unable to set bpf filter to pcap handler", err))
			}
			p.LaPortLog(fmt.Sprintln("Creating Listener for intf", p.IntfNum))
			//p.LaPortLog(fmt.Sprintf("Creating Listener for intf", p.IntfNum))
			p.handle = handle
			src := gopacket.NewPacketSource(p.handle, layers.LayerTypeEthernet)
			in := src.Packets()
			// start rx routine
			LaRxMain(p.PortNum, in)
			p.LaPortLog(fmt.Sprintln("Rx Main Started for port", p.PortNum, sysId))

			// register the tx func
			if sgi != nil {
				sgi.LaSysGlobalRegisterTxCallback(p.IntfNum, TxViaLinuxIf)
			}
		}
	} else {
		p.LaPortLog("Unabled to find AGG assocaited with this port")
	}
}

func (p *LaAggPort) DeleteRxTx() {
	var a *LaAggregator
	var sysId LacpSystem
	if LaFindAggById(p.AggId, &a) {
		mac, _ := net.ParseMAC(a.Config.SystemIdMac)
		sysId.Actor_System = convertNetHwAddressToSysIdKey(mac)
		sysId.Actor_System_priority = a.Config.SystemPriority
	}

	sgi := LacpSysGlobalInfoByIdGet(sysId)
	if sgi != nil {
		sgi.LaSysGlobalDeRegisterTxCallback(p.IntfNum)
	}

	// close rx/tx processing
	if p.handle != nil {
		p.handle.Close()
		p.LaPortLog(fmt.Sprintf("RX/TX handle closed for port", p.PortNum))
		p.handle = nil
	}
}

func (p *LaAggPort) EnableLogging(ena bool) {
	p.logEna = ena
}

func (p *LaAggPort) PortChannelGet() chan string {
	return p.portChan
}

// IsPortAdminEnabled will check if provisioned port enable
// State is enabled or disabled
func (p *LaAggPort) IsPortAdminEnabled() bool {
	return p.PortEnabled
}

func (p *LaAggPort) IsPortOperStatusUp() bool {
	for _, client := range utils.GetAsicDPluginList() {
		p.LinkOperStatus = client.GetPortLinkStatus(utils.GetIfIndexFromName(p.IntfNum))
	}
	return p.LinkOperStatus
}

// IsPortEnabled will check if port is admin enabled
// and link is operationally up
func (p *LaAggPort) IsPortEnabled() bool {
	return p.IsPortAdminEnabled() && p.LinkOperStatus
}

func (p *LaAggPort) LaAggPortDelete() {
	// notify DR that port has been deleted
	if p.AggAttached != nil &&
		p.AggAttached.DrniName != "" {
		if deletecb, ok := LacpCbDb.PortCreateDbList[p.AggAttached.DrniName]; ok {
			deletecb(int32(p.PortNum))
		}
	}
	utils.DeleteEventMap(int32(p.PortNum))
	p.Stop()
	for _, sgi := range LacpSysGlobalInfoGet() {
		for Key, port := range sgi.PortMap {
			if port.PortNum == p.PortNum ||
				port.IntfNum == p.IntfNum {
				// remove the port from the port map
				delete(sgi.PortMap, Key)
				for i, delPort := range sgi.LacpSysGlobalAggPortListGet() {
					if delPort.PortNum == p.PortNum {
						sgi.PortList = append(sgi.PortList[:i], sgi.PortList[i+1:]...)
					}
				}
				return
			}
		}
	}
}

func (p *LaAggPort) Stop() {

	p.DeleteRxTx()

	//p.BEGIN(true)
	// stop the State machines
	// TODO maybe run these in parrallel?
	if p.RxMachineFsm != nil {
		p.RxMachineFsm.Stop()
	}

	if p.PtxMachineFsm != nil {
		p.PtxMachineFsm.Stop()
	}

	if p.TxMachineFsm != nil {
		p.TxMachineFsm.Stop()
	}

	if p.CdMachineFsm != nil {
		p.CdMachineFsm.Stop()
	}

	if p.PCdMachineFsm != nil {
		p.PCdMachineFsm.Stop()
	}

	if p.MuxMachineFsm != nil {
		p.MuxMachineFsm.Stop()
	}

	if p.MarkerResponderFsm != nil {
		p.MarkerResponderFsm.Stop()
	}

	// lets wait for all the State machines to have stopped
	p.wg.Wait()
	close(p.portChan)
}

//  BEGIN will initiate all the State machines
// and will send an event back to this caller
// to begin processing.
func (p *LaAggPort) BEGIN(restart bool) {
	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	// System in being initalized
	p.begin = true

	if !restart {
		// start all the State machines
		// Order here matters as Rx machine
		// will send event to Mux machine
		// thus machine must be up and
		// running first
		// Mux Machine
		p.LacpMuxMachineMain()
		// Periodic Tx Machine
		p.LacpPtxMachineMain()
		// Churn Detection Machine
		p.LacpActorCdMachineMain()
		// Partner Churn Detection Machine
		p.LacpPartnerCdMachineMain()
		// Rx Machine
		p.LacpRxMachineMain()
		// Tx Machine
		p.LacpTxMachineMain()
		// Marker Responder
		p.LampMarkerResponderMain()
	}

	// wait group used when stopping all the
	// State mahines associated with this port.
	// want to ensure that all routines are stopped
	// before proceeding with cleanup thus why not
	// create the wg as part of a BEGIN process
	// 1) Rx Machine
	// 2) Tx Machine
	// 3) Mux Machine
	// 4) Periodic Tx Machine
	// 5) Churn Detection Machine * 2
	// 6) Marker Responder
	// Rxm
	if p.RxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpRxmEventBegin,
			Src: PortConfigModuleStr})
	}

	// Ptxm
	if p.PtxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpPtxmEventBegin,
			Src: PortConfigModuleStr})
	}
	// Cdm
	if p.CdMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.CdMachineFsm.CdmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpCdmEventBegin,
			Src: PortConfigModuleStr})
	}
	// Cdm
	if p.PCdMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PCdMachineFsm.CdmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpCdmEventBegin,
			Src: PortConfigModuleStr})
	}
	// Muxm
	if p.MuxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.MuxMachineFsm.MuxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpMuxmEventBegin,
			Src: PortConfigModuleStr})
	}
	// Txm
	if p.TxMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpTxmEventBegin,
			Src: PortConfigModuleStr})
	}
	// Marker Responder
	if p.MarkerResponderFsm != nil {
		mEvtChan = append(mEvtChan, p.MarkerResponderFsm.LampMarkerResponderEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LampMarkerResponderEventBegin,
			Src: PortConfigModuleStr})
	}
	// call the begin event for each
	// distribute the port disable event to various machines
	p.DistributeMachineEvents(mEvtChan, evt, true)
}

// DistributeMachineEvents will distribute the events in parrallel
// to each machine
func (p *LaAggPort) DistributeMachineEvents(mec []chan utils.MachineEvent, e []utils.MachineEvent, waitForResponse bool) {

	length := len(mec)
	if len(mec) != len(e) {
		p.LaPortLog("LAPORT: Distributing of events failed")
		return
	}

	// send all begin events to each machine in parrallel
	for j := 0; j < length; j++ {
		go func(port *LaAggPort, w bool, idx int, machineEventChannel []chan utils.MachineEvent, event []utils.MachineEvent) {
			if w {
				event[idx].ResponseChan = p.portChan
			}
			event[idx].Src = PortConfigModuleStr
			machineEventChannel[idx] <- event[idx]
		}(p, waitForResponse, j, mec, e)
	}

	if waitForResponse {
		i := 0
		// lets wait for all the machines to respond
		for {
			select {
			case mStr := <-p.portChan:
				i++
				p.LaPortLog(strings.Join([]string{"LAPORT:", mStr, "response received"}, " "))
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

// LaAggPortDisable will update the status on the port
// as well as inform the appropriate State machines of the
// State change
func (p *LaAggPort) LaAggPortDisable() {
	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)
	p.LaPortLog("LAPORT: Port Disabled")

	// port is disabled
	p.PortEnabled = false

	// Rxm
	if !p.portMoved {
		mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpRxmEventNotPortEnabledAndNotPortMoved,
			Src: PortConfigModuleStr})
	}
	// Ptxm
	mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
	evt = append(evt, utils.MachineEvent{
		E:   LacpPtxmEventNotPortEnabled,
		Src: PortConfigModuleStr})

	// Cdm
	mEvtChan = append(mEvtChan, p.CdMachineFsm.CdmEvents)
	evt = append(evt, utils.MachineEvent{
		E:   LacpCdmEventNotPortEnabled,
		Src: PortConfigModuleStr})

	// Partner Cdm
	mEvtChan = append(mEvtChan, p.PCdMachineFsm.CdmEvents)
	evt = append(evt, utils.MachineEvent{
		E:   LacpCdmEventNotPortEnabled,
		Src: PortConfigModuleStr})

	if p.lacpEnabled {
		// Txm
		mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpTxmEventLacpDisabled,
			Src: PortConfigModuleStr})
	}

	// distribute the port disable event to various machines
	p.DistributeMachineEvents(mEvtChan, evt, true)
}

// LaAggPortEnabled will update the status on the port
// as well as inform the appropriate State machines of the
// State change
// When this is called, it is assumed that all States are
// in their default State.
func (p *LaAggPort) LaAggPortEnabled() {
	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	p.LaPortLog("LAPORT: Port Enabled")

	// port is enabled
	p.PortEnabled = true

	// Rxm
	if p.RxMachineFsm.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled {
		if p.lacpEnabled {
			if p.RxMachineFsm.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled {
				mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
				evt = append(evt, utils.MachineEvent{
					E:   LacpRxmEventPortEnabledAndLacpEnabled,
					Src: PortConfigModuleStr})
			}
		} else {
			mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpRxmEventPortEnabledAndLacpDisabled,
				Src: PortConfigModuleStr})
		}
	}

	// Ptxm
	if p.PtxMachineFsm.LacpPtxIsNoPeriodicExitCondition() {
		mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpPtxmEventUnconditionalFallthrough,
			Src: PortConfigModuleStr})
	}

	// Cdm
	if LacpStateIsSet(p.ActorOper.State, LacpStateSyncBit) {
		mEvtChan = append(mEvtChan, p.CdMachineFsm.CdmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpCdmEventActorOperPortStateSyncOn,
			Src: PortConfigModuleStr})
	}
	// Partner Cdm
	if LacpStateIsSet(p.PartnerOper.State, LacpStateSyncBit) {
		mEvtChan = append(mEvtChan, p.CdMachineFsm.CdmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpCdmEventPartnerOperPortStateSyncOn,
			Src: PortConfigModuleStr})
	}

	// Txm
	if p.lacpEnabled {
		mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
		evt = append(evt, utils.MachineEvent{E: LacpTxmEventLacpEnabled,
			Src: PortConfigModuleStr})
	}

	// distribute the port disable event to various machines
	p.DistributeMachineEvents(mEvtChan, evt, true)
}

// LaAggPortLacpDisable will update the status on the port
// as well as inform the appropriate State machines of the
// State change
func (p *LaAggPort) LaAggPortLacpDisable() {
	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	p.LaPortLog("LAPORT: Port LACP Disabled")

	// port is disabled
	p.lacpEnabled = false

	// Rxm
	if p.PortEnabled {
		if p.RxMachineFsm.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled {
			mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpRxmEventPortEnabledAndLacpDisabled,
				Src: PortConfigModuleStr})
		} else {
			// lets force a port disable on the rx machine
			// Important to note we are not changing the port enable status as it
			// is not really disabled, just lacp is disabled
			// State should transition from PORT_DISBLED to LACP_DISABLED state
			mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpRxmEventNotPortEnabledAndNotPortMoved,
				Src: PortConfigModuleStr})
		}

		// Ptxm
		mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
		evt = append(evt, utils.MachineEvent{E: LacpPtxmEventLacpDisabled,
			Src: PortConfigModuleStr})

		// Txm, if lacp is disabled then should not transmit packets
		mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   LacpTxmEventLacpDisabled,
			Src: PortConfigModuleStr})

		// distribute the port disable event to various machines
		p.DistributeMachineEvents(mEvtChan, evt, true)
	}

	var a *LaAggregator
	var sysId LacpSystem
	if LaFindAggById(p.AggId, &a) {
		mac, _ := net.ParseMAC(a.Config.SystemIdMac)
		sysId.Actor_System = convertNetHwAddressToSysIdKey(mac)
		sysId.Actor_System_priority = a.Config.SystemPriority
	}

	// port is no longer controlling lacp State
	sgi := LacpSysGlobalInfoByIdGet(sysId)
	p.ActorAdmin.State = sgi.ActorStateDefaultParams.State
	p.ActorOper.State = sgi.ActorStateDefaultParams.State
}

// LaAggPortEnabled will update the status on the port
// as well as inform the appropriate State machines of the
// State change
func (p *LaAggPort) LaAggPortLacpEnabled(mode int) {
	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	p.LaPortLog(fmt.Sprintln("LAPORT: Port LACP Enabled selected", p.aggSelected))

	// port has lacp enabled, this must be set priror to notifying
	// the state machines as they may depend on this being enabled
	p.lacpEnabled = true

	// port can be added to aggregator
	LacpStateSet(&p.ActorAdmin.State, LacpStateAggregationBit)

	// Activity mode
	if mode == LacpModeActive {
		LacpStateSet(&p.ActorAdmin.State, LacpStateActivityBit)
	} else {
		LacpStateClear(&p.ActorAdmin.State, LacpStateActivityBit)
	}

	if p.IsPortEnabled() &&
		p.aggSelected == LacpAggSelected {
		// Rxm
		if p.RxMachineFsm.Machine.Curr.CurrentState() == LacpRxmStatePortDisabled {
			mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpRxmEventPortEnabledAndLacpEnabled,
				Src: PortConfigModuleStr})
		} else if p.RxMachineFsm.Machine.Curr.CurrentState() == LacpRxmStateLacpDisabled {
			mEvtChan = append(mEvtChan, p.RxMachineFsm.RxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpRxmEventLacpEnabled,
				Src: PortConfigModuleStr})
		}

		// txm
		if p.TxMachineFsm.Machine.Curr.CurrentState() == LacpTxmStateOff {
			mEvtChan = append(mEvtChan, p.TxMachineFsm.TxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpTxmEventLacpEnabled,
				Src: PortConfigModuleStr})
		}

		// Ptxm
		if p.PtxMachineFsm.LacpPtxIsNoPeriodicExitCondition() {
			mEvtChan = append(mEvtChan, p.PtxMachineFsm.PtxmEvents)
			evt = append(evt, utils.MachineEvent{
				E:   LacpPtxmEventUnconditionalFallthrough,
				Src: PortConfigModuleStr})
		}

		// distribute the port disable event to various machines
		p.DistributeMachineEvents(mEvtChan, evt, true)
	}

}

func (p *LaAggPort) LaAggPortActorAdminInfoSet(sysIdMac [6]uint8, sysPrio uint16) {

	p.LaPortLog(fmt.Sprintf("Changing Actor Admin SystemId: MAC: %+v Priority %d ", sysIdMac, sysPrio))
	p.ActorAdmin.System.Actor_System = sysIdMac
	p.ActorAdmin.System.Actor_System_priority = sysPrio
	// only change the oper status if this is not owned by DR
	// if it is owned by DR then will ignore as the oper status
	// is based on the portal system info
	if p.DrniName == "" {
		p.ActorOper.System.Actor_System = p.ActorAdmin.System.Actor_System
		p.ActorOper.System.Actor_System_priority = p.ActorAdmin.System.Actor_System_priority
	}

	p.aggSelected = LacpAggUnSelected

	if p.ModeGet() == LacpModeOn ||
		!p.lacpEnabled ||
		!p.IsPortEnabled() {
		return
	}

	// partner info should be wrong so lets force sync to be off
	LacpStateClear(&p.PartnerOper.State, LacpStateSyncBit)

	p.checkConfigForSelection()

}

func (p *LaAggPort) LaAggPortActorOperInfoSet(sysIdMac [6]uint8, sysPrio uint16) {

	p.LaPortLog(fmt.Sprintf("Changing Actor Oper SystemId: MAC: %+v Priority %d ", sysIdMac, sysPrio))
	// only change the oper status if this is not owned by DR
	// if it is owned by DR then will ignore as the oper status
	// is based on the portal system info
	p.ActorOper.System.Actor_System = sysIdMac
	p.ActorOper.System.Actor_System_priority = sysPrio

	p.aggSelected = LacpAggUnSelected

	if p.ModeGet() == LacpModeOn ||
		!p.lacpEnabled ||
		!p.IsPortEnabled() {
		return
	}

	// partner info should be wrong so lets force sync to be off
	LacpStateClear(&p.PartnerOper.State, LacpStateSyncBit)

	p.checkConfigForSelection()

}

func (p *LaAggPort) TimeoutGet() time.Duration {
	return p.PtxMachineFsm.PeriodicTxTimerInterval
}

func (p *LaAggPort) ModeGet() int {
	return LacpModeGet(p.ActorOper.State, p.lacpEnabled)
}

func LacpCopyLacpPortInfoFromPkt(fromPortInfoPtr *layers.LACPPortInfo, toPortInfoPtr *LacpPortInfo) {
	toPortInfoPtr.Key = fromPortInfoPtr.Key
	toPortInfoPtr.port = fromPortInfoPtr.Port
	toPortInfoPtr.Port_pri = fromPortInfoPtr.PortPri
	toPortInfoPtr.State = fromPortInfoPtr.State
	toPortInfoPtr.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(fromPortInfoPtr.System.SystemId))
	toPortInfoPtr.System.LacpSystemActorSystemPrioritySet(fromPortInfoPtr.System.SystemPriority)
}

// LacpCopyLacpPortInfo:
// Copy the LacpPortInfo data from->to
func LacpCopyLacpPortInfo(fromPortInfoPtr *LacpPortInfo, toPortInfoPtr *LacpPortInfo) {
	toPortInfoPtr.Key = fromPortInfoPtr.Key
	toPortInfoPtr.port = fromPortInfoPtr.port
	toPortInfoPtr.Port_pri = fromPortInfoPtr.Port_pri
	toPortInfoPtr.State = fromPortInfoPtr.State
	toPortInfoPtr.System.LacpSystemActorSystemIdSet(convertSysIdKeyToNetHwAddress(fromPortInfoPtr.System.Actor_System))
	toPortInfoPtr.System.LacpSystemActorSystemPrioritySet(fromPortInfoPtr.System.Actor_System_priority)
}

func LacpLacpPktPortInfoIsEqual(aPortInfoPtr *layers.LACPPortInfo, bPortInfoPtr *LacpPortInfo, StateBits uint8) bool {
	//utils.GlobalLogger.Info(fmt.Sprintf("LacpLacpPktPortInfoIsEqual: pkt %+v  port %+v %t %t", aPortInfoPtr, bPortInfoPtr, LacpStateIsSet(aPortInfoPtr.State, StateBits), LacpStateIsSet(bPortInfoPtr.State, StateBits)))
	return aPortInfoPtr.System.SystemId == bPortInfoPtr.System.Actor_System &&
		aPortInfoPtr.System.SystemPriority == bPortInfoPtr.System.Actor_System_priority &&
		aPortInfoPtr.Port == bPortInfoPtr.port &&
		aPortInfoPtr.PortPri == bPortInfoPtr.Port_pri &&
		aPortInfoPtr.Key == bPortInfoPtr.Key &&
		(LacpStateIsSet(aPortInfoPtr.State, StateBits) && LacpStateIsSet(bPortInfoPtr.State, StateBits))
}

// LacpLacpPortInfoIsEqual:
// Compare the LacpPortInfo data except be selective
// about the State bits that is being compared against
func LacpLacpPortInfoIsEqual(aPortInfoPtr *LacpPortInfo, bPortInfoPtr *LacpPortInfo, StateBits uint8) bool {

	return aPortInfoPtr.System.Actor_System == bPortInfoPtr.System.Actor_System &&
		aPortInfoPtr.System.Actor_System_priority == bPortInfoPtr.System.Actor_System_priority &&
		aPortInfoPtr.port == bPortInfoPtr.port &&
		aPortInfoPtr.Port_pri == bPortInfoPtr.Port_pri &&
		aPortInfoPtr.Key == bPortInfoPtr.Key &&
		(LacpStateIsSet(aPortInfoPtr.State, StateBits) && LacpStateIsSet(bPortInfoPtr.State, StateBits))
}
