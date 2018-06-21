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

// port.go
package stp

import (
	"asicd/asicdCommonDefs"
	"asicd/pluginManager/pluginCommon"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	//"github.com/vishvananda/netlink"
	"net"
	"strconv"
	"strings"
	"sync"
	//"syscall"
	"time"
)

var PortMapTable map[PortMapKey]*StpPort
var PortListTable []*StpPort
var PortConfigMap map[int32]portConfig

const PortConfigModuleStr = "PORT CFG"

type PortMapKey struct {
	IfIndex    int32
	BrgIfIndex int32
}

type portConfig struct {
	Name         string
	HardwareAddr net.HardwareAddr
	Speed        int32
	IfIndex      int32
}

type StpPort struct {
	IfIndex        int32
	ProtocolPortId uint16

	// 17.19
	AgeingTime                  int32 // TODO STP functionality
	Agree                       bool
	Agreed                      bool
	AdminEdge                   bool
	AutoEdgePort                bool // optional
	AdminPathCost               int32
	BpduGuard                   bool
	BpduGuardInterval           int32
	BridgeAssurance             bool
	BridgeAssuranceInconsistant bool
	Disputed                    bool
	FdbFlush                    bool
	Forward                     bool
	Forwarding                  bool
	InfoIs                      PortInfoState
	Learn                       bool
	Learning                    bool
	Mcheck                      bool
	MsgPriority                 PriorityVector
	MsgTimes                    Times
	NewInfo                     bool
	OperEdge                    bool
	PortEnabled                 bool
	PortId                      uint16
	PortPathCost                uint32
	PortPriority                PriorityVector
	PortTimes                   Times
	Priority                    uint16
	Proposed                    bool
	Proposing                   bool
	RcvdBPDU                    bool
	RcvdInfo                    PortDesignatedRcvInfo
	RcvdMsg                     bool
	RcvdRSTP                    bool
	RcvdSTP                     bool
	RcvdTc                      bool
	RcvdTcAck                   bool
	RcvdTcn                     bool
	RstpVersion                 bool
	ReRoot                      bool
	Reselect                    bool
	Role                        PortRole
	Selected                    bool
	SelectedRole                PortRole
	SendRSTP                    bool
	Sync                        bool
	Synced                      bool
	TcAck                       bool
	TcProp                      bool
	Tick                        bool
	TxCount                     uint64
	UpdtInfo                    bool
	// 6.4.3
	OperPointToPointMAC  bool
	AdminPointToPointMAC PointToPointMac

	// link operational state
	AdminPortEnabled bool

	// Associated Bridge Id
	BridgeId   BridgeId
	BrgIfIndex int32
	b          *Bridge

	// statistics
	BpduRx  uint64
	BpduTx  uint64
	StpRx   uint64
	StpTx   uint64
	TcRx    uint64
	TcTx    uint64
	TcAckRx uint64
	TcAckTx uint64
	RstpRx  uint64
	RstpTx  uint64
	PvstRx  uint64
	PvstTx  uint64

	ForwardingTransitions uint64

	// 17.17
	EdgeDelayWhileTimer PortTimer
	FdWhileTimer        PortTimer
	HelloWhenTimer      PortTimer
	MdelayWhiletimer    PortTimer
	RbWhileTimer        PortTimer
	RcvdInfoWhiletimer  PortTimer
	RrWhileTimer        PortTimer
	TcWhileTimer        PortTimer
	BAWhileTimer        PortTimer
	BPDUGuardTimer      PortTimer

	PrxmMachineFsm *PrxmMachine
	PtmMachineFsm  *PtmMachine
	PpmmMachineFsm *PpmmMachine
	PtxmMachineFsm *PtxmMachine
	PimMachineFsm  *PimMachine
	BdmMachineFsm  *BdmMachine
	PrtMachineFsm  *PrtMachine
	TcMachineFsm   *TcMachine
	PstMachineFsm  *PstMachine

	begin bool

	// handle used to tx packets to linux if
	handle *pcap.Handle

	// a way to sync all machines
	wg sync.WaitGroup
	// chanel to send response messages
	portChan chan string

	// used to poll linux interface status.  Useful for SIM/TEST
	PollingRoutine bool
	PollingTimer   *time.Timer
}

type PortTimer struct {
	count int32
}

func (t *PortTimer) GetCount() int32 {
	return t.count
}

func NewStpPort(c *StpPortConfig) *StpPort {
	var b *Bridge
	/*
		MigrateTimeDefault = 3
		BridgeHelloTimeDefault = 2
		BridgeMaxAgeDefault = 20
		BridgeForwardDelayDefault = 15
		TransmitHoldCountDefault = 6
	*/
	enabled := c.Enable
	if enabled {
		for i, client := range GetAsicDPluginList() {
			if i == 0 {
				enabled = client.GetPortLinkStatus(c.IfIndex)
			} else {
				tmpena := client.GetPortLinkStatus(c.IfIndex)
				if tmpena != enabled {
					StpLogger("ERROR", fmt.Sprintf("plugin link status do not match %t"))
				}
			}
		}
	} else {
		// in the case of tests we may not find the actual link so lets force
		// enabled to configured value
		StpLogger("INFO", fmt.Sprintf("Did not find port, forcing enabled to %t", c.Enable))
		enabled = c.Enable
	}

	var RootTimes Times
	if StpFindBridgeByIfIndex(c.BrgIfIndex, &b) {
		RootTimes = b.RootTimes
	}
	p := &StpPort{
		IfIndex:              c.IfIndex,
		AutoEdgePort:         false, // default and not configurable
		AdminPathCost:        c.AdminPathCost,
		AdminPointToPointMAC: PointToPointMac(c.AdminPointToPoint),
		// protocol portId
		PortId:              uint16(pluginCommon.GetIdFromIfIndex(c.IfIndex)),
		Priority:            c.Priority, // default usually 0x80
		AdminPortEnabled:    c.Enable,
		PortEnabled:         enabled,
		PortPathCost:        uint32(c.PathCost),
		Role:                PortRoleDisabledPort,
		SelectedRole:        PortRoleDisabledPort,
		PortTimes:           RootTimes,
		SendRSTP:            b.ForceVersion >= 2, // default
		RcvdRSTP:            b.ForceVersion >= 2, // default
		RstpVersion:         b.ForceVersion >= 2,
		Mcheck:              b.ForceVersion >= 2,
		EdgeDelayWhileTimer: PortTimer{count: MigrateTimeDefault},
		FdWhileTimer:        PortTimer{count: int32(b.RootTimes.ForwardingDelay)}, // TODO same as ForwardingDelay above
		HelloWhenTimer:      PortTimer{count: int32(b.RootTimes.HelloTime)},
		MdelayWhiletimer:    PortTimer{count: MigrateTimeDefault},
		RbWhileTimer:        PortTimer{count: int32(b.RootTimes.HelloTime * 2)},
		RcvdInfoWhiletimer:  PortTimer{count: int32(b.RootTimes.HelloTime * 3)},
		RrWhileTimer:        PortTimer{count: int32(b.RootTimes.MaxAge)},
		TcWhileTimer:        PortTimer{count: int32(b.RootTimes.HelloTime)}, // should be updated by newTcWhile func
		BAWhileTimer:        PortTimer{count: int32(b.RootTimes.HelloTime * 3)},
		portChan:            make(chan string),
		BrgIfIndex:          c.BrgIfIndex,
		AdminEdge:           c.AdminEdgePort,
		PortPriority: PriorityVector{
			RootBridgeId:       b.BridgeIdentifier,
			RootPathCost:       0,
			DesignatedBridgeId: b.BridgeIdentifier,
			DesignatedPortId:   uint16(uint16(pluginCommon.GetIdFromIfIndex(c.IfIndex)) | c.Priority<<8),
		},
		BridgeAssurance:   c.BridgeAssurance,
		BpduGuard:         c.BpduGuard,
		BpduGuardInterval: c.BpduGuardInterval,
		b:                 b, // reference to brige
	}

	if c.AdminPathCost == 0 {
		// TODO need to get speed of port to automatically the port path cost
		// Table 17-3
		AutoPathCostDefaultMap := map[int32]uint32{
			10:         200000000,
			100:        20000000,
			1000:       2000000,
			10000:      200000,
			100000:     20000,
			1000000:    2000,
			10000000:   200,
			100000000:  20,
			1000000000: 2,
		}
		speed := PortConfigMap[p.IfIndex].Speed

		p.PortPathCost = AutoPathCostDefaultMap[speed]
		StpLogger("INFO", fmt.Sprintf("Auto Port Path Cost for port %d speed %d = %d", p.IfIndex, speed, p.PortPathCost))
	}

	key := PortMapKey{
		IfIndex:    p.IfIndex,
		BrgIfIndex: p.b.BrgIfIndex,
	}
	portDbMutex.Lock()
	PortMapTable[key] = p
	if len(PortListTable) == 0 {
		PortListTable = make([]*StpPort, 0)
	}
	PortListTable = append(PortListTable, p)
	portDbMutex.Unlock()
	ifName, _ := PortConfigMap[p.IfIndex]
	StpLogger("DEBUG", fmt.Sprintf("NEW PORT: ifname %s %#v\n", ifName.Name, p))

	if p.PortEnabled {
		p.CreateRxTx()
	}

	/*
		if strings.Contains(ifName.Name, "eth") ||
			strings.Contains(ifName.Name, "lo") {
			p.PollLinuxLinkStatus()
		}
	*/
	return p

}

/* NOT NEEDED
func (p *StpPort) PollLinuxLinkStatus() {

	p.PollingTimer = time.NewTimer(time.Second * 1)
	p.PollingRoutine = true

	go func(p *StpPort) {
		StpMachineLogger("DEBUG", "LINUX POLLING", p.IfIndex, p.BrgIfIndex, "Start")
		for {
			if !p.PollingRoutine {
				fmt.Println("Stopping Link Polling Routine")
				return
			}
			select {
			case <-p.PollingTimer.C:
				netif, _ := netlink.LinkByName(PortConfigMap[p.IfIndex].Name)
				netifattr := netif.Attrs()
				//StpLogger("DEBUG", fmt.Sprintf("Polling link flags%#v, running=0x%x up=0x%x check1 %t check2 %t", netifattr.Flags, syscall.IFF_RUNNING, syscall.IFF_UP, ((netifattr.Flags>>6)&0x1) == 1, (netifattr.Flags&1) == 1))
				//if (((netifattr.Flags >> 6) & 0x1) == 1) && (netifattr.Flags&1) == 1 {
				if (netifattr.Flags & 1) == 1 {
					//StpLogger("DEBUG", "LINUX LINK UP")
					prevPortEnabled := p.PortEnabled
					p.PortEnabled = true
					p.NotifyPortEnabled("LINUX LINK STATUS", prevPortEnabled, true)

				} else {
					//StpLogger("DEBUG", "LINUX LINK DOWN")
					prevPortEnabled := p.PortEnabled
					p.PortEnabled = false
					p.NotifyPortEnabled("LINUX LINK STATUS", prevPortEnabled, false)

				}
				p.PollingTimer.Reset(time.Second * 1)
			}
		}
	}(p)
}
*/
func DelStpPort(p *StpPort) {
	p.Stop()
	key := PortMapKey{
		IfIndex:    p.IfIndex,
		BrgIfIndex: p.b.BrgIfIndex,
	}
	portDbMutex.Lock()
	defer portDbMutex.Unlock()

	// remove from global port table
	delete(PortMapTable, key)

	for i, delPort := range PortListTable {
		if delPort.IfIndex == p.IfIndex &&
			delPort.BrgIfIndex == p.BrgIfIndex {
			if len(PortListTable) == 1 {
				PortListTable = nil
			} else {
				PortListTable = append(PortListTable[:i], PortListTable[i+1:]...)
			}
		}
	}
}

func StpFindPortByIfIndex(pId int32, brgId int32, p **StpPort) bool {
	var ok bool
	key := PortMapKey{
		IfIndex:    pId,
		BrgIfIndex: brgId,
	}
	portDbMutex.Lock()
	defer portDbMutex.Unlock()

	//fmt.Println("looking for key in map", key, PortMapTable)
	if *p, ok = PortMapTable[key]; ok {
		return true
	}
	return false
}

func (p *StpPort) Stop() {

	if p.PollingTimer != nil {
		p.PollingTimer.Stop()
		// used to stop the go routine
		p.PollingRoutine = false
		p.PollingTimer = nil
	}

	p.DeleteRxTx()

	if p.PimMachineFsm != nil {
		p.PimMachineFsm.Stop()
		p.PimMachineFsm = nil
	}

	if p.PrxmMachineFsm != nil {
		p.PrxmMachineFsm.Stop()
		p.PrxmMachineFsm = nil
	}

	if p.PtmMachineFsm != nil {
		p.PtmMachineFsm.Stop()
		p.PtmMachineFsm = nil
	}

	if p.PpmmMachineFsm != nil {
		p.PpmmMachineFsm.Stop()
		p.PpmmMachineFsm = nil
	}

	if p.PtxmMachineFsm != nil {
		p.PtxmMachineFsm.Stop()
		p.PtxmMachineFsm = nil
	}

	if p.BdmMachineFsm != nil {
		p.BdmMachineFsm.Stop()
		p.BdmMachineFsm = nil
	}

	if p.PrtMachineFsm != nil {
		p.PrtMachineFsm.Stop()
		p.PrtMachineFsm = nil
	}

	if p.TcMachineFsm != nil {
		p.TcMachineFsm.Stop()
		p.TcMachineFsm = nil
	}

	if p.PstMachineFsm != nil {
		p.PstMachineFsm.Stop()
		p.PstMachineFsm = nil
	}

	// lets wait for the machines to close
	p.wg.Wait()
	close(p.portChan)

}

func (p *StpPort) BEGIN(restart bool) {
	mEvtChan := make([]chan MachineEvent, 0)
	evt := make([]MachineEvent, 0)

	//p.begin = true

	if !restart {
		// start all the State machines
		// Port Timer State Machine
		p.PtmMachineMain()
		// Port Receive State Machine
		p.PrxmMachineMain()
		// Port Protocol Migration State Machine
		p.PpmmMachineMain()
		// Port Transmit State Machine
		p.PtxmMachineMain()
		// Port Information State Machine
		p.PimMachineMain()
		// Bridge Detection State Machine
		p.BdmMachineMain()
		// Port Role Selection State Machine (one instance per bridge)
		//p.PrsMachineMain()
		// Port Role Transitions State Machine
		p.PrtMachineMain()
		// Topology Change State Machine
		p.TcMachineMain()
		// Port State Transition State Machine
		p.PstMachineMain()
	}

	// Prxm
	if p.PrxmMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PrxmMachineFsm.PrxmEvents)
		evt = append(evt, MachineEvent{e: PrxmEventBegin,
			src: PortConfigModuleStr})

		if p.handle != nil {
			p.CreateRxTx()

		}
	}

	// Ptm
	if p.PtmMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PtmMachineFsm.PtmEvents)
		evt = append(evt, MachineEvent{e: PtmEventBegin,
			src: PortConfigModuleStr})
	}

	// Ppm
	if p.PpmmMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PpmmMachineFsm.PpmmEvents)
		evt = append(evt, MachineEvent{e: PpmmEventBegin,
			src: PortConfigModuleStr})
	}

	// Ptxm
	if p.PtxmMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PtxmMachineFsm.PtxmEvents)
		evt = append(evt, MachineEvent{e: PtxmEventBegin,
			src: PortConfigModuleStr})
	}

	// Pim
	if p.PimMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PimMachineFsm.PimEvents)
		evt = append(evt, MachineEvent{e: PimEventBegin,
			src: PortConfigModuleStr})
	}

	// Bdm
	if p.BdmMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.BdmMachineFsm.BdmEvents)
		if p.AdminEdge {
			evt = append(evt, MachineEvent{e: BdmEventBeginAdminEdge,
				src: PortConfigModuleStr})
		} else {
			evt = append(evt, MachineEvent{e: BdmEventBeginNotAdminEdge,
				src: PortConfigModuleStr})
		}
	}

	// Prtm
	if p.PrtMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PrtMachineFsm.PrtEvents)
		evt = append(evt, MachineEvent{e: PrtEventBegin,
			src: PortConfigModuleStr})
	}

	// Tcm
	if p.TcMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.TcMachineFsm.TcEvents)
		evt = append(evt, MachineEvent{e: TcEventBegin,
			src: PortConfigModuleStr})
	}

	// Pstm
	if p.PstMachineFsm != nil {
		mEvtChan = append(mEvtChan, p.PstMachineFsm.PstEvents)
		evt = append(evt, MachineEvent{e: PstEventBegin,
			src: PortConfigModuleStr})
	}

	// call the begin event for each
	// distribute the port disable event to various machines
	p.DistributeMachineEvents(mEvtChan, evt, true)

	// lets start the tick timer
	if p.PtmMachineFsm != nil {
		p.PtmMachineFsm.TickTimerStart()
	}
	StpMachineLogger("INFO", PortConfigModuleStr, p.IfIndex, p.BrgIfIndex, "BEGIN complete")
	p.begin = false
}

func (p *StpPort) CreateRxTx() {

	if p.handle == nil {
		// lets setup the port receive/transmit handle
		ifName, _ := PortConfigMap[p.IfIndex]
		handle, err := pcap.OpenLive(ifName.Name, 65536, true, 50*time.Millisecond)
		if err != nil {
			// failure here may be ok as this may be SIM
			if !strings.Contains(ifName.Name, "SIM") {
				StpLogger("ERROR", fmt.Sprintf("Error creating pcap OpenLive handle for port %d %s %s\n", p.IfIndex, ifName.Name, err))
			}
			return
		}

		filter := fmt.Sprintf("ether dst 01:80:C2:00:00:00 or 01:00:0C:CC:CC:CD")
		err = handle.SetBPFFilter(filter)
		if err != nil {
			StpLogger("ERROR", fmt.Sprintln("Unable to set bpf filter to pcap handler", p.IfIndex, ifName.Name, err))
			return
		}

		StpLogger("INFO", fmt.Sprintf("Creating STP Listener for intf %d %s\n", p.IfIndex, ifName.Name))
		//p.LaPortLog(fmt.Sprintf("Creating Listener for intf", p.IntfNum))
		p.handle = handle

		// start rx routine
		src := gopacket.NewPacketSource(p.handle, layers.LayerTypeEthernet)
		in := src.Packets()
		BpduRxMain(p.IfIndex, p.b.BrgIfIndex, in)
	}
}

func (p *StpPort) DeleteRxTx() {
	if p.handle != nil {
		p.handle.Close()
		p.handle = nil
		StpLogger("INFO", fmt.Sprintf("RX/TX handle closed for port %d\n", p.IfIndex))
	}
}

// DistributeMachineEvents will distribute the events in parrallel
// to each machine
func (p *StpPort) DistributeMachineEvents(mec []chan MachineEvent, e []MachineEvent, waitForResponse bool) {

	length := len(mec)
	if len(mec) != len(e) {
		StpLogger("ERROR", "STPPORT: Distributing of events failed")
		return
	}

	// send all begin events to each machine in parrallel
	for j := 0; j < length; j++ {
		go func(port *StpPort, w bool, idx int, machineEventChannel []chan MachineEvent, event []MachineEvent) {
			if w {
				event[idx].responseChan = p.portChan
			}
			event[idx].src = PortConfigModuleStr
			//fmt.Println("distribute events", machineEventChannel[idx], event[idx])
			machineEventChannel[idx] <- event[idx]
		}(p, waitForResponse, j, mec, e)
	}

	if waitForResponse && length > 0 {
		i := 0
		// lets wait for all the machines to respond
		for {
			select {
			case mStr := <-p.portChan:
				i++
				//fmt.Println("distribute events response received", mStr)
				StpLogger("DEBUG", strings.Join([]string{"STPPORT:", mStr, "response received"}, " "))
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

// This function can be used to know whether or not to apply configuration parameters
// to one port or all ports, as of 3/21/16 applying to all ports
func (p *StpPort) GetPortListToApplyConfigTo() (newlist []*StpPort) {
	portDbMutex.Lock()
	defer portDbMutex.Unlock()
	for _, port := range PortListTable {
		if p.IfIndex == port.IfIndex {
			newlist = append(newlist, port)
		}
	}
	return
}

func (p *StpPort) SetRxPortCounters(ptype BPDURxType) {
	p.BpduRx++
	switch ptype {
	case BPDURxTypeSTP:
		p.StpRx++
	case BPDURxTypeRSTP:
		p.RstpRx++
	case BPDURxTypeTopo:
		p.TcRx++
	case BPDURxTypeTopoAck:
		p.TcAckRx++
	case BPDURxTypePVST:
		p.PvstRx++
	}
}

func (p *StpPort) SetTxPortCounters(ptype BPDURxType) {
	p.BpduTx++
	switch ptype {
	case BPDURxTypeSTP:
		p.StpTx++
	case BPDURxTypeRSTP:
		p.RstpTx++
	case BPDURxTypeTopo:
		p.TcTx++
	case BPDURxTypeTopoAck:
		p.TcAckTx++
	case BPDURxTypePVST:
		p.PvstTx++
	}
}

/*
BPDU info
ProtocolId        uint16
	ProtocolVersionId byte
	BPDUType          byte
	Flags             byte
	RootId            [8]byte
	RootCostPath      uint32
	BridgeId          [8]byte
	PortId            uint16
	MsgAge            uint16
	MaxAge            uint16
	HelloTime         uint16
	FwdDelay          uint16
*/

func (p *StpPort) SaveMsgRcvInfo(data interface{}) {

	bpdumsg := data.(RxBpduPdu)
	bpduLayer := bpdumsg.pdu

	switch bpduLayer.(type) {
	case *layers.STP:
		stp := bpduLayer.(*layers.STP)
		// TODO revisit what the BridgePortId should be
		p.MsgPriority.BridgePortId = stp.PortId
		p.MsgPriority.DesignatedBridgeId = stp.BridgeId
		p.MsgPriority.DesignatedPortId = stp.PortId
		p.MsgPriority.RootBridgeId = stp.RootId
		p.MsgPriority.RootPathCost = stp.RootPathCost

		p.MsgTimes.ForwardingDelay = stp.FwdDelay >> 8
		p.MsgTimes.HelloTime = stp.HelloTime >> 8
		p.MsgTimes.MaxAge = stp.MaxAge >> 8
		p.MsgTimes.MessageAge = stp.MsgAge >> 8

	case *layers.RSTP:
		rstp := bpduLayer.(*layers.RSTP)
		// TODO revisit what the BridgePortId should be
		p.MsgPriority.BridgePortId = rstp.PortId
		p.MsgPriority.DesignatedBridgeId = rstp.BridgeId
		p.MsgPriority.DesignatedPortId = rstp.PortId
		p.MsgPriority.RootBridgeId = rstp.RootId
		p.MsgPriority.RootPathCost = rstp.RootPathCost

		p.MsgTimes.ForwardingDelay = rstp.FwdDelay >> 8
		p.MsgTimes.HelloTime = rstp.HelloTime >> 8
		p.MsgTimes.MaxAge = rstp.MaxAge >> 8
		p.MsgTimes.MessageAge = rstp.MsgAge >> 8

	case *layers.PVST:
		pvst := bpduLayer.(*layers.PVST)
		// TODO revisit what the BridgePortId should be
		p.MsgPriority.BridgePortId = pvst.PortId
		p.MsgPriority.DesignatedBridgeId = pvst.BridgeId
		p.MsgPriority.DesignatedPortId = pvst.PortId
		p.MsgPriority.RootBridgeId = pvst.RootId
		p.MsgPriority.RootPathCost = pvst.RootPathCost

		p.MsgTimes.ForwardingDelay = pvst.FwdDelay >> 8
		p.MsgTimes.HelloTime = pvst.HelloTime >> 8
		p.MsgTimes.MaxAge = pvst.MaxAge >> 8
		p.MsgTimes.MessageAge = pvst.MsgAge >> 8

	}
}

func (p *StpPort) BridgeProtocolVersionGet() uint8 {
	// TODO get the protocol version from the bridge
	// Below is the default
	return layers.RSTPProtocolVersion
}

func (p *StpPort) NotifyPortEnabled(src string, oldportenabled bool, newportenabled bool) {
	// The following Machines need to know about
	// changes in PortEnable State
	// 1) Port Receive
	// 2) Port Protocol Migration
	// 3) Port Information
	// 4) Bridge Detection
	if oldportenabled != newportenabled {
		StpMachineLogger("INFO", src, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("NotifyPortEnabled: %t", newportenabled))
		mEvtChan := make([]chan MachineEvent, 0)
		evt := make([]MachineEvent, 0)

		// notify the state machines
		if !newportenabled {

			if p.EdgeDelayWhileTimer.count != MigrateTimeDefault {
				if p.PrxmMachineFsm != nil {
					mEvtChan = append(mEvtChan, p.PrxmMachineFsm.PrxmEvents)
					evt = append(evt, MachineEvent{e: PrxmEventEdgeDelayWhileNotEqualMigrateTimeAndNotPortEnabled,
						src: src})
				}
			}

			if p.PpmmMachineFsm != nil {
				if p.PpmmMachineFsm.Machine.Curr.CurrentState() == PpmmStateCheckingRSTP {
					if p.MdelayWhiletimer.count != MigrateTimeDefault {
						mEvtChan = append(mEvtChan, p.PpmmMachineFsm.PpmmEvents)
						evt = append(evt, MachineEvent{e: PpmmEventMdelayNotEqualMigrateTimeAndNotPortEnabled,
							src: src})
					}
				} else {
					mEvtChan = append(mEvtChan, p.PpmmMachineFsm.PpmmEvents)
					evt = append(evt, MachineEvent{e: PpmmEventNotPortEnabled,
						src: src})
				}
			}

			if p.BdmMachineFsm != nil {
				if !p.AdminEdge {
					if p.BdmMachineFsm.Machine.Curr.CurrentState() == BdmStateEdge {
						//BdEventNotPortEnabledAndNotAdminEdge
						mEvtChan = append(mEvtChan, p.BdmMachineFsm.BdmEvents)
						evt = append(evt, MachineEvent{e: BdmEventNotPortEnabledAndNotAdminEdge,
							src: src})
					}
				} else {
					if p.BdmMachineFsm.Machine.Curr.CurrentState() == BdmStateNotEdge {
						//BdmEventNotPortEnabledAndAdminEdge
						mEvtChan = append(mEvtChan, p.BdmMachineFsm.BdmEvents)
						evt = append(evt, MachineEvent{e: BdmEventNotPortEnabledAndAdminEdge,
							src: src})
					}
				}
			}
			if p.PimMachineFsm != nil {
				if p.InfoIs != PortInfoStateDisabled {
					mEvtChan = append(mEvtChan, p.PimMachineFsm.PimEvents)
					evt = append(evt, MachineEvent{e: PimEventNotPortEnabledInfoIsNotEqualDisabled,
						src: src})
				}
			}

		} else {

			// reset this timer to allow for packet to be received
			if p.BridgeAssurance {
				p.BAWhileTimer.count = int32(p.b.RootTimes.HelloTime * 3)
			}

			/*
				This should only be triggered from RcvdBpdu being set becuase
				we need packet setn
				if p.PrxmMachineFsm.Machine.Curr.CurrentState() == PrxmStateDiscard {
					if p.RcvdBPDU {
						mEvtChan = append(mEvtChan, p.PrxmMachineFsm.PrxmEvents)
						evt = append(evt, MachineEvent{e: PrxmEventRcvdBpduAndPortEnabled,
							src: src})
					}
				} else if p.PrxmMachineFsm.Machine.Curr.CurrentState() == PrxmStateReceive {
					if p.RcvdBPDU &&
						!p.RcvdMsg {
						mEvtChan = append(mEvtChan, p.PrxmMachineFsm.PrxmEvents)
						evt = append(evt, MachineEvent{e: PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg,
							src: src})
					}
				}
			*/

			if p.PimMachineFsm != nil {
				if p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateDisabled {
					mEvtChan = append(mEvtChan, p.PimMachineFsm.PimEvents)
					evt = append(evt, MachineEvent{e: PimEventPortEnabled,
						src: src})
				}
			}
		}
		if len(mEvtChan) > 0 {
			p.DistributeMachineEvents(mEvtChan, evt, true)
		}
	}
}

func (p *StpPort) NotifyRcvdMsgChanged(src string, oldrcvdmsg bool, newrcvdmsg bool, data interface{}) {
	// The following machines need to know about
	// changed in rcvdMsg state
	// 1) Port Receive
	// 2) Port Information
	//fmt.Println("NotifyRcvdMsgChanged old/new", oldrcvdmsg, newrcvdmsg, p.RcvdMsg, p.PimMachineFsm.Machine.Curr.CurrentState())
	if oldrcvdmsg != newrcvdmsg {
		/*
			NOT a valid transition RcvdMsg is only PRX -> PIM
			if src != PrxmMachineModuleStr {
				if p.PrxmMachineFsm.Machine.Curr.CurrentState() == PrxmStateReceive &&
					p.RcvdBPDU &&
					p.PortEnabled &&
					!p.RcvdMsg {
					p.PrxmMachineFsm.PrxmEvents <- MachineEvent{
						e:   PrxmEventRcvdBpduAndPortEnabledAndNotRcvdMsg,
						src: src,
					}
				}
			}
		*/
		if src != PimMachineModuleStr {
			bpdumsg := data.(RxBpduPdu)
			bpduLayer := bpdumsg.pdu

			if p.PimMachineFsm != nil {
				if p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateDisabled {
					if p.RcvdMsg {
						p.PimMachineFsm.PimEvents <- MachineEvent{
							e:    PimEventRcvdMsg,
							src:  src,
							data: bpduLayer,
						}
					}
				} else if p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateCurrent {
					if p.RcvdMsg &&
						!p.UpdtInfo {
						p.PimMachineFsm.PimEvents <- MachineEvent{
							e:    PimEventRcvdMsgAndNotUpdtInfo,
							src:  src,
							data: bpduLayer,
						}
					} else if p.InfoIs == PortInfoStateReceived &&
						p.RcvdInfoWhiletimer.count == 0 &&
						!p.UpdtInfo &&
						!p.RcvdMsg {
						p.PimMachineFsm.PimEvents <- MachineEvent{
							e:    PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg,
							src:  src,
							data: bpduLayer,
						}
					}
				}
			}
		}
	}
}

func (p *StpPort) ForceVersionSet(src string, val int) {
	// The following machines need to know about
	// changes in forceVersion State
	// 1) Port Protocol Migration
	// 2) Port Role Transitions
	// 3) Port Transmit
}

func (p *StpPort) NotifyUpdtInfoChanged(src string, oldupdtinfo bool, newupdtinfo bool) {
	// The following machines need to know about
	// changes in UpdtInfo State
	// 1) Port Information
	// 2) Port Role Transitions
	// 3) Port Transmit
	if oldupdtinfo != newupdtinfo {
		//StpMachineLogger("DEBUG", src, p.IfIndex, fmt.Sprintf("updateinfo changed %d", newupdtinfo))
		// PI
		if p.UpdtInfo {
			if src != PimMachineModuleStr &&
				p.PimMachineFsm != nil &&
				(p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateAged ||
					p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateCurrent) {
				if p.Selected {
					p.PimMachineFsm.PimEvents <- MachineEvent{
						e:   PimEventSelectedAndUpdtInfo,
						src: src,
					}
				}
			}
		} else {
			if p.PimMachineFsm != nil {
				if src != PimMachineModuleStr &&
					p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateCurrent {
					if p.RcvdMsg {
						p.PimMachineFsm.PimEvents <- MachineEvent{
							e:   PimEventRcvdMsgAndNotUpdtInfo,
							src: src,
						}
					} else if p.InfoIs == PortInfoStateReceived &&
						p.RcvdInfoWhiletimer.count == 0 &&
						!p.RcvdMsg {
						p.PimMachineFsm.PimEvents <- MachineEvent{
							e:   PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg,
							src: src,
						}
					}
				}
			}
			/*
				if src != PtxmMachineModuleStr &&
					p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle &&
					p.Selected {
					if p.HelloWhenTimer.count == 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.SendRSTP &&
						p.NewInfo &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleRootPort &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleDesignatedPort &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}

				}
				if src != PrtMachineModuleStr &&
					p.Selected {
					if p.SelectedRole == PortRoleDisabledPort &&
						p.Role != p.SelectedRole {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.SelectedRole == PortRoleRootPort &&
						p.Role != p.SelectedRole {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.SelectedRole == PortRoleDesignatedPort &&
						p.Role != p.SelectedRole {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.SelectedRole == PortRoleAlternatePort &&
						p.Role != p.SelectedRole {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.SelectedRole == PortRoleBackupPort &&
						p.Role != p.SelectedRole {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSelectedRoleEqualBackupPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisablePort &&
						!p.Learning &&
						!p.Forwarding {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
						if p.FdWhileTimer.count != int32(p.PortTimes.MaxAge) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
						if p.Proposed &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.AllSynced() &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Proposed &&
							p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Forward &&
							!p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.SelectedRole == PortRoleRootPort &&
							p.Role != p.SelectedRole {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RrWhileTimer.count != int32(p.PortTimes.ForwardingDelay) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RstpVersion &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RbWhileTimer.count == 0 &&
							p.RstpVersion &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RstpVersion &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RbWhileTimer.count == 0 &&
							p.RstpVersion &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
						if !p.Forward &&
							!p.Agreed &&
							!p.Proposing &&
							!p.OperEdge {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Learning &&
							!p.Forwarding &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RrWhileTimer.count == 0 &&
							p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							!p.Synced &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							!p.Synced &&
							!p.OperEdge &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RrWhileTimer.count != 0 &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RrWhileTimer.count != 0 &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Disputed &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Disputed &&
							!p.OperEdge &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
						if p.Proposed &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.AllSynced() &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Proposed &&
							p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count != int32(p.PortTimes.ForwardingDelay) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RbWhileTimer.count != int32(2*p.PortTimes.HelloTime) &&
							p.Role == PortRoleBackupPort {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBlockPort {
						if !p.Learning &&
							!p.Forwarding {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}
				}*/
		}
	}
}

func (p *StpPort) NotifySyncedChanged(src string, oldsynced bool, newsynced bool) {

	if oldsynced != newsynced {
		if src != PrtMachineModuleStr {
			if p.PrtMachineFsm != nil {
				if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
					if !p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
				} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
					if p.b.AllSynced() &&
						!p.Agree &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
				} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
					if !p.Learning &&
						!p.Forwarding &&
						!p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.Agreed &&
						!p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.OperEdge &&
						!p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.Sync &&
						p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.Sync &&
						!p.Synced &&
						!p.OperEdge &&
						p.Learn &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.Sync &&
						!p.Synced &&
						!p.OperEdge &&
						p.Forward &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
				} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
					if p.b.AllSynced() &&
						!p.Agree &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if !p.Synced &&
						p.Selected &&
						!p.UpdtInfo {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
				}
			}
		}
	}
}

/*
No need to when these parameters change as selected/updtInfo will
be used as the trigger for update
func (p *StpPort) DesignatedPrioritySet(src string, val *PriorityVector) {
	// The following machines need to know about
	// changes in DesignatedPriority State
	// 1) Port Information
	// 2) Port Transmit
	p.DesignatedPriority = *val
}

func (p *StpPort) DesignatedTimesSet(src string, val *Times) {
	// The following machines need to know about
	// changes in DesignatedTimes State
	// 1) Port Information
	// 2) Port Transmit
	p.DesignatedTimes = *val
}
*/
func (p *StpPort) NotifySelectedChanged(src string, oldselected bool, newselected bool) {
	// The following machines need to know about
	// changes in Selected State
	// 1) Port Information
	// 2) Port Role Transitions
	// 3) Port Transmit
	if oldselected != newselected {
		StpMachineLogger("DEBUG", src, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("NotifySelectedChanged Role[%d] SelectedRole[%d] Forwarding[%t] Learning[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Agreed[%t] Agree[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
			p.Role, p.SelectedRole, p.Forwarding, p.Learning, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))

		// PI
		if p.Selected {
			/*
				if src != PtxmMachineModuleStr &&
					p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle &&
					!p.UpdtInfo {
					if p.HelloWhenTimer.count == 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if p.SendRSTP &&
						p.NewInfo &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleRootPort &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleDesignatedPort &&
						p.TxCount < p.b.TxHoldCount &&
						p.HelloWhenTimer.count != 0 {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}

				}
			*/
			if src == PrsMachineModuleStr &&
				p.PrtMachineFsm != nil {
				if !p.UpdtInfo {
					// PRSM -> PRTM
					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisablePort &&
						!p.Learning &&
						!p.Forwarding {
						p.PrtMachineFsm.PrtEvents <- MachineEvent{
							e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
							src: src,
						}
					}
					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
						if p.FdWhileTimer.count != int32(p.PortTimes.MaxAge) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					} else {
						if p.SelectedRole == PortRoleDisabledPort &&
							p.Role != p.SelectedRole {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}
					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
						if p.Proposed &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.AllSynced() &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Proposed &&
							p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Forward &&
							!p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RrWhileTimer.count != int32(p.PortTimes.ForwardingDelay) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RstpVersion &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.ReRooted(p) &&
							p.RbWhileTimer.count == 0 &&
							p.RstpVersion &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RstpVersion &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.ReRooted(p) &&
							p.RbWhileTimer.count == 0 &&
							p.RstpVersion &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}

					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
						if !p.Forward &&
							!p.Agreed &&
							!p.Proposing &&
							!p.OperEdge {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Learning &&
							!p.Forwarding &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RrWhileTimer.count == 0 &&
							p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							!p.Synced &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync &&
							!p.Synced &&
							!p.OperEdge &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RrWhileTimer.count != 0 &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot &&
							p.RrWhileTimer.count != 0 &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Disputed &&
							!p.OperEdge &&
							p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Disputed &&
							!p.OperEdge &&
							p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.ReRoot &&
							!p.Sync &&
							!p.Learn {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count == 0 &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Agreed &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							p.RrWhileTimer.count == 0 &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.OperEdge &&
							!p.ReRoot &&
							!p.Sync &&
							p.Learn &&
							!p.Forward {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}

					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
						if p.Proposed &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.b.AllSynced() &&
							!p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Proposed &&
							p.Agree {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.FdWhileTimer.count != int32(p.PortTimes.ForwardingDelay) {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.Sync {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.ReRoot {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if !p.Synced {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotSyncedAndSelectedAndNotUpdtInfo,
								src: src,
							}
						} else if p.RbWhileTimer.count != int32(2*p.PortTimes.HelloTime) &&
							p.Role == PortRoleBackupPort {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}

					if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBlockPort {
						if !p.Learning &&
							!p.Forwarding {
							p.PrtMachineFsm.PrtEvents <- MachineEvent{
								e:   PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo,
								src: src,
							}
						}
					}
				} else {
					// PRSM -> PIM
					if p.PimMachineFsm != nil {
						if p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateAged ||
							p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateCurrent {
							if p.UpdtInfo {
								p.PimMachineFsm.PimEvents <- MachineEvent{
									e:   PimEventSelectedAndUpdtInfo,
									src: src,
								}
							}
						}
					}
				}
			}
		}
	}
}

func (p *StpPort) NotifyOperEdgeChanged(src string, oldoperedge bool, newoperedge bool) {
	// The following machines need to know about
	// changes in OperEdge State
	// 1) Port Role Transitions
	// 2) Bridge Detection
	if oldoperedge != newoperedge &&
		p.PrtMachineFsm != nil {
		// Prt update 17.29.3
		if p.PrtMachineFsm != nil &&
			p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort &&
			src != PrtMachineModuleStr {
			if !p.Forward &&
				!p.Agreed &&
				!p.Proposing &&
				!p.OperEdge &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.OperEdge &&
				!p.Synced &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.Sync &&
				!p.Synced &&
				!p.OperEdge &&
				p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.Sync &&
				!p.Synced &&
				!p.OperEdge &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.ReRoot &&
				p.RrWhileTimer.count != 0 &&
				!p.OperEdge &&
				p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.ReRoot &&
				p.RrWhileTimer.count != 0 &&
				!p.OperEdge &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.Disputed &&
				!p.OperEdge &&
				p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.Disputed &&
				!p.OperEdge &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.OperEdge &&
				p.RrWhileTimer.count == 0 &&
				!p.Sync &&
				!p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.OperEdge &&
				!p.ReRoot &&
				!p.Sync &&
				!p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.OperEdge &&
				p.RrWhileTimer.count == 0 &&
				!p.Sync &&
				p.Learn &&
				!p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					src: src,
				}
			} else if p.OperEdge &&
				!p.ReRoot &&
				!p.Sync &&
				p.Learn &&
				!p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					src: src,
				}
			}
		}
		// Bdm
		if p.BdmMachineFsm != nil {
			if p.BdmMachineFsm.Machine.Curr.CurrentState() == BdmStateEdge &&
				src != BdmMachineModuleStr {
				if !p.OperEdge {
					p.BdmMachineFsm.BdmEvents <- MachineEvent{
						e:   BdmEventNotOperEdge,
						src: src,
					}
				}
			}
		}
	}
}

func (p *StpPort) NotifySelectedRoleChanged(src string, oldselectedrole PortRole, newselectedrole PortRole) {

	if oldselectedrole != newselectedrole &&
		p.PrtMachineFsm != nil {
		StpMachineLogger("DEBUG", src, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("NotifySelectedRoleChange: role[%d] selectedRole[%d]", p.Role, p.SelectedRole))
		/*if p.Role != p.SelectedRole {*/
		if newselectedrole == PortRoleDisabledPort {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
				src: src,
			}
		} else if newselectedrole == PortRoleRootPort {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
				src: src,
			}
		} else if newselectedrole == PortRoleDesignatedPort {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
				src: src,
			}
		} else if newselectedrole == PortRoleAlternatePort {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
				src: src,
			}
		} else if newselectedrole == PortRoleBackupPort {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventSelectedRoleEqualBackupPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo,
				src: src,
			}
		}
		/*}*/
	}
}

func (p *StpPort) NotifyProposingChanged(src string, oldproposing bool, newproposing bool) {
	if oldproposing != newproposing {
		if src != BdmMachineModuleStr &&
			p.BdmMachineFsm != nil {
			if p.BdmMachineFsm.Machine.Curr.CurrentState() == BdmStateNotEdge {
				if p.EdgeDelayWhileTimer.count == 0 &&
					p.AutoEdgePort &&
					p.SendRSTP &&
					p.Proposing {
					p.BdmMachineFsm.BdmEvents <- MachineEvent{
						e:   BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
		if src != PrsMachineModuleStr && p.PrtMachineFsm != nil {
			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
				if !p.Forward &&
					!p.Agreed &&
					!p.Proposing &&
					!p.OperEdge &&
					p.Selected &&
					!p.UpdtInfo {
					p.PrtMachineFsm.PrtEvents <- MachineEvent{
						e:   PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo,
						src: PrtMachineModuleStr,
					}
				}
			}
		}

	}
}
func (p *StpPort) NotifyRcvdTcRcvdTcnRcvdTcAck(oldrcvdtc bool, oldrcvdtcn bool, oldrcvdtcack bool, newrcvdtc bool, newrcvdtcn bool, newrcvdtcack bool) {

	// only care if there was a change
	//if oldrcvdtc != newrcvdtc ||
	//	oldrcvdtcn != newrcvdtcn ||
	//	oldrcvdtcack != newrcvdtcack {
	//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("TC state[%s] tcn[%t] tcack[%t] tcn[%t]",
	//	TcStateStrMap[p.TcMachineFsm.Machine.Curr.CurrentState()], p.RcvdTc, p.RcvdTcAck, p.RcvdTcn))
	if p.TcMachineFsm != nil {
		if p.RcvdTc &&
			(p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning ||
				p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateActive) {
			p.TcMachineFsm.TcEvents <- MachineEvent{
				e:   TcEventRcvdTc,
				src: RxModuleStr,
			}
		}
		if p.RcvdTcn &&
			(p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning ||
				p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateActive) {

			p.TcMachineFsm.TcEvents <- MachineEvent{
				e:   TcEventRcvdTcn,
				src: RxModuleStr,
			}
		}
		if p.RcvdTcAck &&
			(p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning ||
				p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateActive) {

			p.TcMachineFsm.TcEvents <- MachineEvent{
				e:   TcEventRcvdTcAck,
				src: RxModuleStr,
			}
		}
		if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning &&
			p.Role != PortRoleRootPort &&
			p.Role != PortRoleDesignatedPort &&
			!(p.Learn || p.Learning) &&
			!(p.RcvdTc || p.RcvdTcn || p.RcvdTcAck || p.TcProp) {

			p.TcMachineFsm.TcEvents <- MachineEvent{
				e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
				src: RxModuleStr,
			}
		}
	}
}

func (p *StpPort) EdgeDelay() uint16 {
	if p.OperPointToPointMAC {
		return MigrateTimeDefault
	} else {
		return p.b.RootTimes.MaxAge
	}
}

// check if any other bridge port is adminEdge
func (p *StpPort) IsAdminEdgePort() bool {
	portDbMutex.Lock()
	defer portDbMutex.Unlock()

	for _, ptmp := range PortListTable {
		if p != ptmp &&
			ptmp.IfIndex == p.IfIndex {
			if p.AdminEdge || p.OperEdge {
				return true
			}
		}
	}
	return false
}

func ConstructPortConfigMap() {
	currMarker := int(asicdCommonDefs.MIN_SYS_PORTS)
	count := 100
	for _, client := range GetAsicDPluginList() {
		StpLogger("DEBUG", "Calling asicd for port config")
		for {
			bulkInfo, err := client.GetBulkPortState(currMarker, count)
			if err != nil {
				StpLogger("ERROR", fmt.Sprintf("GetBulkPortState Error: %s", err))
				return
			}
			StpLogger("DEBUG", fmt.Sprintf("Length of GetBulkPortState: %d", bulkInfo.Count))

			bulkCfgInfo, err := client.GetBulkPort(currMarker, count)
			if err != nil {
				StpLogger("ERROR", fmt.Sprintf("Error: %s", err))
				return
			}

			StpLogger("DEBUG", fmt.Sprintf("Length of GetBulkPortConfig: %d", bulkCfgInfo.Count))
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifindex := bulkInfo.PortStateList[i].IfIndex
				ent := PortConfigMap[ifindex]
				ent.IfIndex = ifindex
				ent.Name = bulkInfo.PortStateList[i].Name
				ent.HardwareAddr, _ = net.ParseMAC(bulkCfgInfo.PortList[i].MacAddr)
				PortConfigMap[ifindex] = ent
				StpLogger("INIT", fmt.Sprintf("Found Port IfIndex %d Name %s\n", ent.IfIndex, ent.Name))
			}
			if more == false {
				return
			}
		}
	}
}

func GetPortNameFromIfIndex(ifindex int32) string {
	if p, ok := PortConfigMap[ifindex]; ok {
		return p.Name
	}
	return ""
}

// IntfRef can be string number or fpPort
func GetIfIndexFromIntfRef(intfref string) int32 {

	for _, p := range PortConfigMap {
		if p.Name == intfref {
			return p.IfIndex
		} else if s, err := strconv.Atoi(intfref); err == nil {
			if int32(s) == p.IfIndex {
				return p.IfIndex
			}
		}
	}
	return 0
}
