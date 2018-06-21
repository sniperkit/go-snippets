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
// dr.go
package drcp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/google/gopacket/layers"
)

var DistributedRelayDB map[string]*DistributedRelay
var DistributedRelayDBList []*DistributedRelay

var MacCaptureCount map[MacCaptureKey]int

type MacCaptureKey struct {
	mac     string
	intfref string
}

// 802.1ax-2014 7.4.1.1
type DistributedRelay struct {
	DistributedRelayFunction
	DrniId          uint32
	DrniDescription string
	DrniName        string

	// Also defined in 9.4.7
	DrniAggregatorId        [6]uint8
	DrniAggregatorPriority  uint16
	DrniPortalAddr          net.HardwareAddr
	DrniPortalPriority      uint16
	DrniThreeSystemPortal   bool
	DrniPortConversation    [MAX_CONVERSATION_IDS][4]uint16
	DrniGatewayConversation [MAX_CONVERSATION_IDS][]uint8
	// End also defined in 9.4.7

	// save the origional values from the aggregator
	PrevAggregatorId       [6]uint8
	PrevAggregatorPriority uint16

	DrniPortalSystemNumber  uint8                 // 1-3
	DrniIntraPortalLinkList [MAX_IPP_LINKS]uint32 // ifindex
	DrniAggregator          int32
	DrniConvAdminGateway    [MAX_CONVERSATION_IDS][]uint8
	// conversation id -> gateway
	DrniNeighborAdminConvGatewayListDigest Md5Digest
	DrniNeighborAdminConvPortListDigest    Md5Digest
	DrniGatewayAlgorithm                   GatewayAlgorithm
	DrniNeighborAdminGatewayAlgorithm      GatewayAlgorithm
	DrniNeighborAdminPortAlgorithm         GatewayAlgorithm
	DrniNeighborAdminDRCPState             uint8
	DrniEncapMethod                        EncapMethod
	DrniIPLEncapMap                        map[uint32]uint32
	DrniNetEncapMap                        map[uint32]uint32
	DrniPSI                                bool
	DrniPortConversationControl            bool
	DrniPortalPortProtocolIDA              net.HardwareAddr

	// 9.4.10
	PortConversationUpdate     bool
	IppAllPortUpdate           bool
	GatewayConversationUpdate  bool
	IppAllGatewayUpdate        bool
	HomeGatewayVectorTransmit  bool
	OtherGatewayVectorTransmit bool

	// channel used to wait on response from distributed event send
	drEvtResponseChan chan string

	a *lacp.LaAggregator

	// Local list to keep track of distributed port list
	// Server will only indicate that a change has occured
	// updateDRFHomeState will determine what actions
	// to perform based on the differences between
	// what DR and Aggregator distributed port list
	DRAggregatorDistributedList []int32

	// sync creation and deletion
	wg sync.WaitGroup

	// state machines
	PsMachineFsm *PsMachine
	GMachineFsm  *GMachine
	AMachineFsm  *AMachine

	Ipplinks []*DRCPIpp
}

// 802.1ax-2014 Section 9.4.8 Per-DR Function variables
type DistributedRelayFunction struct {
	ChangeDRFPorts                                bool
	ChangePortal                                  bool
	DrniCommonMethods                             bool
	DrniConversationGatewayList                   [MAX_CONVERSATION_IDS]uint32
	DrniPortalSystemState                         [4]StateVectorInfo
	DRFHomeAdminAggregatorKey                     uint16
	DRFHomeConversationGatewayListDigest          Md5Digest
	DRFHomeConversationPortListDigest             Md5Digest
	DRFHomeGatewayAlgorithm                       [4]uint8
	DRFHomeGatewayConversationMask                [MAX_CONVERSATION_IDS]bool
	DRFHomeGatewaySequence                        uint16
	DRFHomePortAlgorithm                          [4]uint8
	DRFHomeOperAggregatorKey                      uint16
	DRFHomeOperPartnerAggregatorKey               uint16
	DRFHomeState                                  StateVectorInfo
	DRFNeighborAdminConversationGatewayListDigest Md5Digest
	DRFNeighborAdminConversationPortListDigest    Md5Digest
	DRFNeighborAdminDRCPState                     layers.DRCPState
	DRFNeighborAdminGatewayAlgorithm              [4]uint8
	DRFNeighborAdminPortAlgorithm                 [4]uint8
	// range 1..3
	DRFPortalSystemNumber uint8
	DRFHomeOperDRCPState  layers.DRCPState

	// 9.3.3.2
	DrniPortalSystemGatewayConversation [MAX_CONVERSATION_IDS]bool
	DrniPortalSystemPortConversation    [MAX_CONVERSATION_IDS]bool
}

// DrFindByPortalAddr each portal address is unique within the system
func DrFindByPortalAddr(portaladdr string, dr **DistributedRelay) bool {
	for _, d := range DistributedRelayDBList {
		if d.DrniPortalAddr.String() == portaladdr {
			*dr = d
			return true
		}
	}
	return false
}

// DrFindByName will find the DR based on the DRNI name
func DrFindByName(DrniName string, dr **DistributedRelay) bool {
	for _, d := range DistributedRelayDBList {
		if d.DrniName == DrniName {
			*dr = d
			return true
		}
	}
	return false
}

func DrGetDrcpNext(dr **DistributedRelay) bool {
	returnNext := false
	for _, d := range DistributedRelayDBList {
		if *dr == nil {
			// first agg
			*dr = d
			return true
		} else if (*dr).DrniName == d.DrniName {
			// found agg
			returnNext = true
		} else if returnNext {
			// next agg
			*dr = d
			return true
		}
	}

	*dr = nil
	return false
}

// DrFindByAggregator will find the DR based on the Aggregator that it is
// associated with
func DrFindByAggregator(DrniAggregator int32, dr **DistributedRelay) bool {
	for _, d := range DistributedRelayDBList {
		if d.DrniAggregator == DrniAggregator {
			*dr = d
			return true
		}
	}
	return false
}

// isPortInConversation will check of the provided portList intersected with
// the aggregator port list is greater than zero
func (dr *DistributedRelay) isAggPortInConverstaion(portList []int32) bool {
	a := dr.a

	if a != nil &&
		a.PortNumList != nil {
		for _, ifindex := range a.PortNumList {
			for _, pifindex := range portList {
				if int32(ifindex) == pifindex {
					return true
				}
			}
		}
	}
	return false
}

// setTimeSharingGatwewayDigest, when the port and gateway algorithm
// is set to time sharing then it should be noted that the gateway
// and port algorithm digest
// currently we only support Vlan based
// to start each
// algorithm is as follows:
// Conversations are not bound to a lag link but rather a portal system,
// thus all down traffic will either go to the local aggregator ports
// or IPL if the destination is a remote portal network port (which is not
// an aggregator port).  All up traffic is only destined to another
// aggregator or other network links either in hte local system or accross
// the IPL to the neighbor system.
// If all local aggregator ports are down then the neighbor system must
// forward frames out the aggregator as well as any network links to
// which the frame is destined for
func (dr *DistributedRelay) SetTimeSharingPortAndGatwewayDigest() {
	// algorithm assumes 2P system only
	if dr.DrniGatewayAlgorithm == GATEWAY_ALGORITHM_CVID {
		if !dr.DrniThreeSystemPortal {
			dr.setAdminConvGatewayAndNeighborGatewayListDigest()
			dr.setAdminConvPortAndNeighborPortListDigest()
		}
	}
}

// setAdminConvGatewayAndNeighborGatewayListDigest will set the predetermined
// algorithm as the gateway.  Every even vlan will have its gateway in system
// 2 and every odd vlan will have its gateway in system 1
func (dr *DistributedRelay) setAdminConvGatewayAndNeighborGatewayListDigest() {
	isNewConversation := false
	ghash := md5.New()
	for cid, conv := range ConversationIdMap {
		if conv.Valid && dr.isAggPortInConverstaion(conv.PortList) {

			// mark this call as new so that we can update the state machines
			if dr.DrniConvAdminGateway[cid] == nil {
				dr.DrniConvAdminGateway[cid] = make([]uint8, 0)
				isNewConversation = true
				// Fixed algorithm for 2P system
				// Because we only support sharing by time we don't really care which
				// system is the "gateway" of the conversation because all conversations
				// are free to be delivered on both systems based on bridging rules.
				// Annex G:
				//  A frame received over the IPL shall never be forwarded over the Aggregator Port.
				//  A frame received over the IPL with a DA that was learned from the Aggregator Port shall be discarded.
				//
				// NOTE when other sharing methods are supported then this algorithm will
				// need to be changed
				if math.Mod(float64(conv.Cvlan), 2) == 0 {
					dr.DrniConvAdminGateway[cid] = append(dr.DrniConvAdminGateway[cid], 2)
					dr.DrniConvAdminGateway[cid] = append(dr.DrniConvAdminGateway[cid], 1)
				} else {
					dr.DrniConvAdminGateway[cid] = append(dr.DrniConvAdminGateway[cid], 1)
					dr.DrniConvAdminGateway[cid] = append(dr.DrniConvAdminGateway[cid], 2)
				}
				dr.LaDrLog(fmt.Sprintf("Adding New Gateway Conversation %d portallist[%+v]", cid, dr.DrniConvAdminGateway[cid]))
			}
			buf := new(bytes.Buffer)
			//dr.LaDrLog(fmt.Sprintf("Adding to Gateway Digest:", conv.Cvlan, math.Mod(float64(conv.Cvlan), 2), []uint8{dr.DrniConvAdminGateway[cid][0], dr.DrniConvAdminGateway[cid][1], uint8(cid >> 8 & 0xff), uint8(cid & 0xff)}))
			// network byte order
			binary.Write(buf, binary.BigEndian, []uint8{dr.DrniConvAdminGateway[cid][0], dr.DrniConvAdminGateway[cid][1], uint8(cid >> 8 & 0xff), uint8(cid & 0xff)})
			ghash.Write(buf.Bytes())
		} else {
			buf := new(bytes.Buffer)
			// network byte order
			binary.Write(buf, binary.BigEndian, []uint16{uint16(cid)})
			ghash.Write(buf.Bytes())

			if dr.DrniConvAdminGateway[cid] != nil {
				dr.LaDrLog(fmt.Sprintf("Clearing Gateway Conversation %d portallist[%+v]", cid, dr.DrniConvAdminGateway[cid]))
				isNewConversation = true
			}

			dr.DrniConvAdminGateway[cid] = nil
		}
	}
	for i, val := range ghash.Sum(nil) {
		dr.DrniNeighborAdminConvGatewayListDigest[i] = val
		dr.DRFNeighborAdminConversationGatewayListDigest[i] = val
		dr.DRFHomeConversationGatewayListDigest[i] = val
	}

	// always send regardless of state because all states expect this event
	if isNewConversation &&
		dr.PsMachineFsm != nil {
		dr.ChangePortal = true
		dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
			E:   PsmEventChangePortal,
			Src: DRCPConfigModuleStr,
		}
	}
}

// setAdminConvGatewayAndNeighborGatewayListDigest will set the predetermined
// algorithm as the gateway.  Port Digest is not used as the port conversation
// is determined by hw hashing algorithm, thus setting no priority port list
// against the digest.
func (dr *DistributedRelay) setAdminConvPortAndNeighborPortListDigest() {
	phash := md5.New()
	for cid, _ := range ConversationIdMap {
		buf := new(bytes.Buffer)
		// network byte order
		binary.Write(buf, binary.BigEndian, []uint16{uint16(cid)})
		phash.Write(buf.Bytes())
	}

	for i, val := range phash.Sum(nil) {
		dr.DrniNeighborAdminConvPortListDigest[i] = val
		dr.DRFNeighborAdminConversationPortListDigest[i] = val
		dr.DRFHomeConversationPortListDigest[i] = val
	}
}

// NewDistributedRelay create a new instance of Distributed Relay and
// the associated objects for the IPP ports
func NewDistributedRelay(cfg *DistributedRelayConfig) *DistributedRelay {

	dr := &DistributedRelay{
		DrniId:                      uint32(cfg.DrniPortalSystemNumber),
		DrniName:                    cfg.DrniName,
		DrniPortalPriority:          cfg.DrniPortalPriority,
		DrniThreeSystemPortal:       cfg.DrniThreePortalSystem,
		DrniPortalSystemNumber:      cfg.DrniPortalSystemNumber,
		DrniIntraPortalLinkList:     cfg.DrniIntraPortalLinkList,
		DrniAggregator:              int32(cfg.DrniAggregator),
		DrniPortConversationControl: cfg.DrniPortConversationControl,
		drEvtResponseChan:           make(chan string),
		DrniIPLEncapMap:             make(map[uint32]uint32),
		DrniNetEncapMap:             make(map[uint32]uint32),
		DistributedRelayFunction: DistributedRelayFunction{
			DRFHomeState: StateVectorInfo{mutex: &sync.Mutex{}},
		},
		DrniPSI: true, // by default this is true until the neighbor pkt is received
	}

	neighborPortalSystemNumber := uint32(2)
	if cfg.DrniPortalSystemNumber == 1 {
		neighborPortalSystemNumber = 1

	}
	// Only support two portal system so we need to adjust
	// the ipp port id.  This should ideally come from the user
	// but lets make provisioning as simple as possible
	for i, ippPortId := range cfg.DrniIntraPortalLinkList {
		if ippPortId>>16&0x3 == 0 {
			dr.DrniIntraPortalLinkList[i] = ippPortId | (neighborPortalSystemNumber << 16)
		}
	}

	for i, _ := range dr.DrniPortalSystemState {
		dr.DrniPortalSystemState[i].mutex = &sync.Mutex{}
	}

	/*
		Not allowing user to set we are goign to fill this in via
		setTimeSharingPortAndGatwewayDigest
		for cid, data := range cfg.DrniConvAdminGateway {
			if data != [3]uint8{} {
				dr.DrniConvAdminGateway[cid] = make([]uint8, 0)
				for _, sysnum := range data {
					if sysnum != 0 {
						dr.DrniConvAdminGateway[cid] = append(dr.DrniConvAdminGateway[cid], sysnum)
					}
				}
			}
		}
	*/
	dr.DrniPortalAddr, _ = net.ParseMAC(cfg.DrniPortalAddress)
	for i, macbyte := range dr.DrniPortalAddr {
		dr.DrniAggregatorId[i] = macbyte
	}

	// string format in bits "00000000"
	for i, j := 0, uint32(7); i < 8; i, j = i+1, j-1 {
		val, _ := strconv.Atoi(cfg.DrniNeighborAdminDRCPState[i : i+1])
		dr.DrniNeighborAdminDRCPState |= uint8(val << j)
		dr.DRFNeighborAdminDRCPState |= layers.DRCPState(val << j)
	}

	/*
		Not allowing user to set we are goign to fill this in via
		setTimeSharingPortAndGatwewayDigest
		for i := 0; i < 16; i++ {
			dr.DrniNeighborAdminConvPortListDigest[i] = cfg.DrniNeighborAdminConvPortListDigest[i]
		}
	*/

	// format "00:00:00:00" or "00-00-00-00"
	encapmethod := strings.Split(cfg.DrniEncapMethod, ":")
	if strings.Contains(cfg.DrniEncapMethod, "-") {
		encapmethod = strings.Split(cfg.DrniEncapMethod, "-")
	}
	gatewayalgorithm := strings.Split(cfg.DrniGatewayAlgorithm, ":")
	if strings.Contains(cfg.DrniGatewayAlgorithm, "-") {
		gatewayalgorithm = strings.Split(cfg.DrniGatewayAlgorithm, "-")
	}

	neighborgatewayalgorithm := strings.Split(cfg.DrniNeighborAdminGatewayAlgorithm, ":")
	if strings.Contains(cfg.DrniNeighborAdminGatewayAlgorithm, "-") {
		neighborgatewayalgorithm = strings.Split(cfg.DrniNeighborAdminGatewayAlgorithm, "-")
	}
	//neighborportalgorithm := strings.Split(cfg.DrniNeighborAdminPortAlgorithm, ":")
	var val1, val2, val3, val4 int64
	val1, _ = strconv.ParseInt(encapmethod[0], 16, 16)
	val2, _ = strconv.ParseInt(encapmethod[1], 16, 16)
	val3, _ = strconv.ParseInt(encapmethod[2], 16, 16)
	val4, _ = strconv.ParseInt(encapmethod[3], 16, 16)
	dr.DrniEncapMethod = EncapMethod{uint8(val1), uint8(val2), uint8(val3), uint8(val4)}
	val1, _ = strconv.ParseInt(gatewayalgorithm[0], 16, 16)
	val2, _ = strconv.ParseInt(gatewayalgorithm[1], 16, 16)
	val3, _ = strconv.ParseInt(gatewayalgorithm[2], 16, 16)
	val4, _ = strconv.ParseInt(gatewayalgorithm[3], 16, 16)
	dr.DrniGatewayAlgorithm = [4]uint8{uint8(val1), uint8(val2), uint8(val3), uint8(val4)}
	val1, _ = strconv.ParseInt(neighborgatewayalgorithm[0], 16, 16)
	val2, _ = strconv.ParseInt(neighborgatewayalgorithm[1], 16, 16)
	val3, _ = strconv.ParseInt(neighborgatewayalgorithm[2], 16, 16)
	val4, _ = strconv.ParseInt(neighborgatewayalgorithm[3], 16, 16)
	dr.DrniNeighborAdminGatewayAlgorithm = [4]uint8{uint8(val1), uint8(val2), uint8(val3), uint8(val4)}
	dr.DRFNeighborAdminGatewayAlgorithm = [4]uint8{uint8(val1), uint8(val2), uint8(val3), uint8(val4)}

	for i, data := range cfg.DrniIPLEncapMap {
		dr.DrniIPLEncapMap[uint32(i)] = data
	}
	for i, data := range cfg.DrniNetEncapMap {
		dr.DrniNetEncapMap[uint32(i)] = data
	}

	netMac, _ := net.ParseMAC(cfg.DrniIntraPortalPortProtocolDA)
	dr.DrniPortalPortProtocolIDA = netMac

	// add to the global db's
	DistributedRelayDB[dr.DrniName] = dr
	DistributedRelayDBList = append(DistributedRelayDBList, dr)

	dr.LaDrLog(fmt.Sprintf("Created Distributed Relay %+v\n", dr))
	dr.LaDrLog(fmt.Sprintf("Created Distributed Relay portal %d\n", dr.DrniPortalSystemNumber))

	for _, ippid := range dr.DrniIntraPortalLinkList {
		portid := ippid & 0xffff
		if portid > 0 {
			ipp := NewDRCPIpp(ippid, dr)
			// disabled until an aggregator has been attached
			ipp.DRCPEnabled = false
			dr.Ipplinks = append(dr.Ipplinks, ipp)
		}
	}

	// register for port and lag port updates for this dr
	dr.RegisterForLacpPortUpdates()

	return dr
}

// DeleteDistriutedRelay will delete the distributed relay along with
// the associated IPP links and de-associate from the Aggregator
func (dr *DistributedRelay) DeleteDistributedRelay() {

	// detach was not called externally, so lets call it
	if dr.a != nil {
		dr.DetachAggregatorFromDistributedRelay(dr.DrniAggregator)
	}

	for _, ipp := range dr.Ipplinks {
		ipp.DeleteDRCPIpp()
	}

	// cleanup the tables hosting the dr data
	// cleanup the tables
	if _, ok := DistributedRelayDB[dr.DrniName]; ok {
		delete(DistributedRelayDB, dr.DrniName)
		for i, deldr := range DistributedRelayDBList {
			if deldr == dr {
				DistributedRelayDBList = append(DistributedRelayDBList[:i], DistributedRelayDBList[i+1:]...)
			}
		}
	}
}

// BEGIN will start/build all the Distributed Relay State Machines and
// send the begin event
func (dr *DistributedRelay) BEGIN(restart bool) {

	mEvtChan := make([]chan utils.MachineEvent, 0)
	evt := make([]utils.MachineEvent, 0)

	// there is a case in which we have only called
	// restart and called main functions outside
	// of this scope (TEST for example)
	//prevBegin := p.begin

	// System in being initalized
	//p.begin = true

	if !restart {
		// Portal System Machine
		dr.DrcpPsMachineMain()
		// Gateway Machine
		dr.DrcpGMachineMain()
		// Aggregator Machine
		dr.DrcpAMachineMain()
	}

	// wait group used when stopping all the
	// State mahines associated with this port.
	// want to ensure that all routines are stopped
	// before proceeding with cleanup thus why not
	// create the wg as part of a BEGIN process
	// 1) Portal System Machine
	// 2) Gateway Machine
	// 3) Aggregator Machine
	// Psm
	if dr.PsMachineFsm != nil {
		mEvtChan = append(mEvtChan, dr.PsMachineFsm.PsmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   PsmEventBegin,
			Src: DRCPConfigModuleStr})
	}

	// Gm
	if dr.GMachineFsm != nil {
		mEvtChan = append(mEvtChan, dr.GMachineFsm.GmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   GmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// Am
	if dr.AMachineFsm != nil {
		mEvtChan = append(mEvtChan, dr.AMachineFsm.AmEvents)
		evt = append(evt, utils.MachineEvent{
			E:   AmEventBegin,
			Src: DRCPConfigModuleStr})
	}
	// call the begin event for each
	// distribute the port disable event to various machines
	dr.DistributeMachineEvents(mEvtChan, evt, true)
}

func (dr *DistributedRelay) waitgroupadd(m string) {
	//fmt.Println("Calling wait group add", m)
	dr.wg.Add(1)
}

func (dr *DistributedRelay) waitgroupstop(m string) {
	//fmt.Println("Calling wait group stop", m)
	dr.wg.Done()
}

func (dr *DistributedRelay) Stop() {

	// Psm
	if dr.PsMachineFsm != nil {
		dr.PsMachineFsm.Stop()
		dr.PsMachineFsm = nil
	}
	// Gm
	if dr.GMachineFsm != nil {
		dr.GMachineFsm.Stop()
		dr.GMachineFsm = nil
	}
	// Am
	if dr.AMachineFsm != nil {
		dr.AMachineFsm.Stop()
		dr.AMachineFsm = nil
	}
	dr.wg.Wait()

	close(dr.drEvtResponseChan)
}

// DistributeMachineEvents will distribute the events in parrallel
// to each machine
func (dr *DistributedRelay) DistributeMachineEvents(mec []chan utils.MachineEvent, e []utils.MachineEvent, waitForResponse bool) {

	length := len(mec)
	if len(mec) != len(e) {
		dr.LaDrLog("LADR: Distributing of events failed")
		return
	}

	// send all begin events to each machine in parrallel
	for j := 0; j < length; j++ {
		go func(d *DistributedRelay, w bool, idx int, machineEventChannel []chan utils.MachineEvent, event []utils.MachineEvent) {
			if w {
				event[idx].ResponseChan = d.drEvtResponseChan
			}
			event[idx].Src = DRCPConfigModuleStr
			machineEventChannel[idx] <- event[idx]
		}(dr, waitForResponse, j, mec, e)
	}

	if waitForResponse {
		i := 0
		// lets wait for all the machines to respond
		for {
			select {
			case mStr := <-dr.drEvtResponseChan:
				i++
				dr.LaDrLog(strings.Join([]string{"LADR:", mStr, "response received"}, " "))
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

// 802.1ax-2014 9.3.4.4
func extractGatewayConversationID() {

}

// 802.1ax-2014 9.3.4.4
func extractPortConversationID() {

}

// updatePortalState This function updates the Drni_Portal_System_State[] as follows
func (dr *DistributedRelay) updatePortalState(src string) {

	// update the local portal info
	dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].mutex.Lock()
	dr.DRFHomeState.mutex.Lock()

	dr.LaDrLog(fmt.Sprintf("updatePortalState (%s): DrniPortalSystemState[%d] from DRFHomeState OpState %t updating vector sequence %d portList %v",
		src,
		dr.DrniPortalSystemNumber,
		dr.DRFHomeState.OpState,
		dr.DRFHomeState.GatewayVector[0].Sequence,
		dr.DRFHomeState.PortIdList))

	dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].OpState = dr.DRFHomeState.OpState
	dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].updateGatewayVector(dr.DRFHomeState.GatewayVector[0].Sequence, dr.DRFHomeState.GatewayVector[0].Vector)
	dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].PortIdList = dr.DRFHomeState.PortIdList
	dr.DRFHomeState.mutex.Unlock()
	dr.DrniPortalSystemState[dr.DrniPortalSystemNumber].mutex.Unlock()

	if len(dr.Ipplinks) > 1 {
		// TODO need for the following case when more than a single IPL is supported
		// if any of the other Portal System’s state information is available from two IPPs in this Portal
		// System, then....
	} else if len(dr.Ipplinks) == 1 {
		// single IPP case
		for _, ipp := range dr.Ipplinks {
			dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].mutex.Lock()
			ipp.DRFNeighborState.mutex.Lock()
			ipp.DRFOtherNeighborState.mutex.Lock()
			if ipp.DRFNeighborState.OpState {
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].OpState = ipp.DRFNeighborState.OpState
				seqvector := ipp.DRFNeighborState.GatewayVector[0]
				dr.LaDrLog(fmt.Sprintf("updatePortalState (%s): DrniPortalSystemState[%d] from DRFNeighborState OpState %t updating vector sequence %d portList %v",
					src,
					ipp.DRFNeighborPortalSystemNumber,
					ipp.DRFNeighborState.OpState,
					seqvector.Sequence,
					ipp.DRFNeighborState.PortIdList))
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].updateGatewayVector(seqvector.Sequence, seqvector.Vector)
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].PortIdList = ipp.DRFNeighborState.PortIdList

			}
			if ipp.DRFOtherNeighborState.OpState {
				ipp.DrniNeighborONN = true
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].OpState = ipp.DRFOtherNeighborState.OpState
				seqvector := ipp.DRFOtherNeighborState.GatewayVector[0]
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].updateGatewayVector(seqvector.Sequence, seqvector.Vector)
				dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].PortIdList = ipp.DRFOtherNeighborState.PortIdList
			}
			ipp.DRFOtherNeighborState.mutex.Unlock()
			ipp.DRFNeighborState.mutex.Unlock()
			dr.DrniPortalSystemState[ipp.DRFNeighborPortalSystemNumber].mutex.Unlock()
		}
	}
	// clear unset portals (ignore first index)
	for i, stateinfo := range dr.DrniPortalSystemState {
		if i != 0 && !stateinfo.OpState {
			dr.DrniPortalSystemState[i].mutex.Lock()
			dr.DrniPortalSystemState[i].GatewayVector = nil
			dr.DrniPortalSystemState[i].PortIdList = nil
			dr.DrniPortalSystemState[i].mutex.Unlock()
		}
	}
	/* This is not right because it becomes circular on the setting above we set the portalSystemState based
	   on the ipp.DRFNeighborState, and now we set DRFNeighborSTate based on PortalSystemState
	// TODO revisit logic as this may be incorrect
	// should only be one ipp
	for _, ipp := range dr.Ipplinks {
		for i := uint8(1); i <= MAX_PORTAL_SYSTEM_IDS; i++ {
			ipp.DRFNeighborState.mutex.Lock()
			if i != dr.DRFPortalSystemNumber {
				if !ipp.DRFNeighborState.OpState {
					dr.DrniPortalSystemState[i].mutex.Lock()
					ipp.DRFNeighborState.OpState = dr.DrniPortalSystemState[i].OpState
					for _, seqvector := range dr.DrniPortalSystemState[i].GatewayVector {
						ipp.DRFNeighborState.updateGatewayVector(seqvector.Sequence, seqvector.Vector)
					}
					ipp.DRFNeighborState.PortIdList = dr.DrniPortalSystemState[i].PortIdList
					dr.DrniPortalSystemState[i].mutex.Unlock()
				}
			}
			ipp.DRFNeighborState.mutex.Unlock()
		}
	}
	*/

	// update ipp_portal_system_state
	// TODO If any other Portal System’s state information is available from two IPPs, then
	if len(dr.Ipplinks) > 1 {

	} else if len(dr.Ipplinks) == 1 {
		// single ipl
		// Ipp portal state contains the neighbor state as first entry and if there
		// are any other portals received on this IPP they will follow
		for _, ipp := range dr.Ipplinks {
			ipp.DRFNeighborState.mutex.Lock()
			if ipp.DRFNeighborState.OpState {
				ipp.IppPortalSystemState[ipp.DRFNeighborPortalSystemNumber].mutex.Lock()
				ipp.IppPortalSystemState[ipp.DRFNeighborPortalSystemNumber].OpState = true
				ipp.IppPortalSystemState[ipp.DRFNeighborPortalSystemNumber].updateGatewayVector(ipp.DRFNeighborState.GatewayVector[0].Sequence, ipp.DRFNeighborState.GatewayVector[0].Vector)
				ipp.IppPortalSystemState[ipp.DRFNeighborPortalSystemNumber].PortIdList = ipp.DRFNeighborState.PortIdList
				ipp.IppPortalSystemState[ipp.DRFNeighborPortalSystemNumber].mutex.Unlock()
				// TODO add the any other portal state info from received DRCP here
				// Not being done today because should only be one other system in 2P config
			}
			ipp.DRFNeighborState.mutex.Unlock()
		}
	}
	for _, ipp := range dr.Ipplinks {
		// clear the port sync as the neighbor should not know about this update
		defer ipp.NotifyNTTDRCPUDChange(PsMachineModuleStr, ipp.NTTDRCPDU, true)
		ipp.NTTDRCPDU = true
	}
}

// RegisterForLacpPortUpdates: way to bridge between packages of lacp and drcp
func (dr *DistributedRelay) RegisterForLacpPortUpdates() {
	lacp.RegisterLaPortCreateCb(dr.DrniName, dr.NotifyAggPortCreate)
	lacp.RegisterLaPortDeleteCb(dr.DrniName, dr.NotifyAggPortDelete)
	lacp.RegisterLaPortUpCb(dr.DrniName, dr.NotifyAggPortUp)
	lacp.RegisterLaPortDownCb(dr.DrniName, dr.NotifyAggPortDown)
	lacp.RegisterLaAggCreateCb(dr.DrniName, dr.AttachAggregatorToDistributedRelay)
	lacp.RegisterLaAggDeleteCb(dr.DrniName, dr.DetachAggregatorFromDistributedRelay)
}

// NotifyAggPortCreate called by lacp when Aggregator is created
func (dr *DistributedRelay) NotifyAggPortCreate(ifindex int32) {
	a := dr.a
	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(ifindex), &p) &&
		dr.DrniName == p.DrniName {
		if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME &&
			a != nil &&
			len(dr.a.PortNumList) == 1 {
			dr.SetTimeSharingPortAndGatwewayDigest()
		}
	}
}

// NotifyAggPortDelete called by lacp when Aggregator is deleted
func (dr *DistributedRelay) NotifyAggPortDelete(ifindex int32) {
	a := dr.a
	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(ifindex), &p) &&
		dr.DrniName == p.DrniName {
		if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME &&
			a != nil &&
			len(dr.a.PortNumList) == 0 {
			dr.SetTimeSharingPortAndGatwewayDigest()
		}
	}
}

// NotifyAggPortUp when aggregator port is up
func (dr *DistributedRelay) NotifyAggPortUp(ifindex int32) {

	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(ifindex), &p) &&
		dr.DrniName == p.DrniName {
		dr.LaDrLog(fmt.Sprintf("Agg Port is Up %d", ifindex))
		if dr.DRAggregatorDistributedList == nil {
			dr.DRAggregatorDistributedList = make([]int32, 0)
		}

		dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList, ifindex)

		dr.ChangeDRFPorts = true
		if dr.PsMachineFsm != nil {
			dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
				E:   PsmEventChangeDRFPorts,
				Src: DRCPConfigModuleStr,
			}
		}
	}
}

// NotifyAggPortUp when aggregator port is down
func (dr *DistributedRelay) NotifyAggPortDown(ifindex int32) {

	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(ifindex), &p) &&
		dr.DrniName == p.DrniName {
		dr.LaDrLog(fmt.Sprintf("Agg Port is Down %d", ifindex))
		if dr.DRAggregatorDistributedList != nil {

			for i, ifidx := range dr.DRAggregatorDistributedList {
				if ifidx == ifindex {
					dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList[:i], dr.DRAggregatorDistributedList[i+1:]...)
					break
				}
			}

			dr.ChangeDRFPorts = true
			if dr.PsMachineFsm != nil {
				dr.PsMachineFsm.PsmEvents <- utils.MachineEvent{
					E:   PsmEventChangeDRFPorts,
					Src: DRCPConfigModuleStr,
				}
			}
		}
	}
}

// AttachAggregatorToDistributedRelay: will attach the aggregator and start the Distributed
// relay protocol for the given dr if this agg is associated with a DR
func (dr *DistributedRelay) AttachAggregatorToDistributedRelay(aggId int32) {

	if dr.DrniAggregator == aggId &&
		dr.a == nil {
		var a *lacp.LaAggregator
		if lacp.LaFindAggById(int(aggId), &a) {
			dr.a = a
			a.DrniName = dr.DrniName
			a.ActorOperKey = uint16(dr.DRFHomeOperAggregatorKey)
			a.PartnerOperKey = a.ActorOperKey

			dr.LaDrLog(fmt.Sprintf("Attaching Agg %s %d to DR %s", a.AggName, a.AggId, dr.DrniName))

			// These values should be the same as the admin
			dr.PrevAggregatorId = a.AggMacAddr
			dr.PrevAggregatorPriority = a.AggPriority
			dr.LaDrLog(fmt.Sprintf("Saving Orig SystemId %+v Priority %d", a.AggMacAddr, a.AggPriority))

			// set the aggregator Id for the aggregator as this is used in
			// setDefaultPortalSystemParameters
			a.AggMacAddr = dr.DrniAggregatorId

			// only need to set this once the key has been negotiated.
			if dr.PsMachineFsm != nil &&
				dr.PsMachineFsm.Machine.Curr.CurrentState() == PsmStatePortalSystemUpdate {
				// lets update the aggregator parameters
				// configured ports
				for _, aggport := range a.PortNumList {
					var p *lacp.LaAggPort
					if lacp.LaFindPortById(aggport, &p) {

						dr.LaDrLog(fmt.Sprintf("Aggregator found updating system parameters moving to unselected until DR is synced"))
						if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
							for _, client := range utils.GetAsicDPluginList() {
								for _, ippid := range dr.DrniIntraPortalLinkList {
									inport := ippid & 0xffff
									if inport > 0 {
										dr.LaDrLog(fmt.Sprintf("AttachAgg: Blocking IPP %d to AggPort %d", inport, aggport))
										/* TMP - changes made to progress the compilation. Need to add actual fpPort */
										err := client.IppIngressEgressDrop("fpPort1", "fpPort2")
										if err != nil {
											dr.LaDrLog(fmt.Sprintf("ERROR (AttachAgg) setting Block from %s tolag port %s", utils.GetNameFromIfIndex(int32(inport)), int32(aggport)))
										}
									}
								}
							}
						}
						// assign the new values to the aggregator
						lacp.SetLaAggPortSystemInfoFromDistributedRelay(
							uint16(aggport),
							fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
								dr.DrniPortalAddr[0],
								dr.DrniPortalAddr[1],
								dr.DrniPortalAddr[2],
								dr.DrniPortalAddr[3],
								dr.DrniPortalAddr[4],
								dr.DrniPortalAddr[5]),
							dr.DrniPortalPriority,
							dr.DRFHomeOperAggregatorKey,
							dr.DrniName,
							false)

					} else {
						dr.LaDrLog(fmt.Sprintf("ERROR unable update system info on port %d not found", aggport))
					}
				}
			}
			if len(a.PortNumList) == 0 {
				dr.LaDrLog(fmt.Sprintf("Aggregator found but port list is empty thus not updating system parameters"))
			}

			// add the port to the local distributed list so that the digests can be
			// calculated
			dr.DRAggregatorDistributedList = make([]int32, 0)
			for _, disport := range a.DistributedPortNumList {
				var aggp *lacp.LaAggPort
				foundPort := false
				for lacp.LaGetPortNext(&aggp) && !foundPort {
					if aggp.IntfNum == disport {
						dr.DRAggregatorDistributedList = append(dr.DRAggregatorDistributedList, int32(aggp.PortNum))
						dr.LaDrLog(fmt.Sprintf("Aggregator port in Distributing State", aggp.PortNum))
					}
				}
			}

			// set port and gateway info and digest
			//dr.SetTimeSharingPortAndGatwewayDigest()
			// TODO can we get away with just setting
			if len(a.PortNumList) > 0 {
				// set this to allow for portal system machine to fall through
				// after initialization
				dr.ChangeDRFPorts = true
				dr.ChangePortal = true
			}

			dr.BEGIN(false)
			// start the IPP links
			for _, ipp := range dr.Ipplinks {
				dr.LaDrLog(fmt.Sprintf("Starting Ipp %s", ipp.Name))
				ipp.DRCPEnabled = true
				ipp.BEGIN(false)
			}
		}
	}
}

// DetachCreatedAggregatorFromDistributedRelay: will detach the aggregator and stop the Distributed
// relay protocol for the given dr if since this aggregator is no longer attached
func (dr *DistributedRelay) DetachAggregatorFromDistributedRelay(aggId int32) {
	if dr.DrniAggregator == aggId &&
		dr.a != nil {
		var a *lacp.LaAggregator
		if lacp.LaFindAggById(int(aggId), &a) {
			// lets update the aggregator parameters
			// configured ports
			for _, aggport := range a.PortNumList {
				lacp.SetLaAggPortSystemInfo(
					uint16(aggport),
					fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
						dr.PrevAggregatorId[0],
						dr.PrevAggregatorId[1],
						dr.PrevAggregatorId[2],
						dr.PrevAggregatorId[3],
						dr.PrevAggregatorId[4],
						dr.PrevAggregatorId[5]),
					dr.PrevAggregatorPriority)

				if dr.DrniEncapMethod == ENCAP_METHOD_SHARING_BY_TIME {
					for _, client := range utils.GetAsicDPluginList() {
						for _, ippid := range dr.DrniIntraPortalLinkList {
							inport := ippid & 0xffff
							if inport > 0 {
								dr.LaDrLog(fmt.Sprintf("UnBlocking IPP %d to AggPort %d", inport, aggport))
								/* TEMP - add actual port names */
								client.IppIngressEgressPass("fpPort1", "fpPort2")
							}
						}
					}
				}
			}
			// reset aggregator values
			a.DrniName = ""
			a.AggMacAddr = dr.PrevAggregatorId
			a.AggPriority = dr.PrevAggregatorPriority
			a.ActorOperKey = a.ActorAdminKey
		}
		lacp.DeRegisterLaAggCbAll(dr.DrniName)
		for _, ipp := range dr.Ipplinks {
			ipp.Stop()
		}
		dr.Stop()
		dr.a = nil
	}
}
