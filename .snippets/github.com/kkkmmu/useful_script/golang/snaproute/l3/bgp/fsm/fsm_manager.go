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

// peer.go
package fsm

import (
	_ "fmt"
	"l3/bgp/baseobjects"
	"l3/bgp/config"
	"l3/bgp/packet"
	"net"
	"sync"
	"utils/logging"
)

type PeerFSMConn struct {
	PeerIP      string
	Established bool
	Conn        *net.Conn
}

type PeerFSMState struct {
	PeerIP string
	State  config.BGPFSMState
}

type PeerAttrs struct {
	PeerIP        string
	BGPId         net.IP
	ASSize        uint8
	HoldTime      uint32
	KeepaliveTime uint32
	AddPathFamily map[packet.AFI]map[packet.SAFI]uint8
}

const (
	BGPCmdReasonNone int = iota
	BGPCmdReasonMaxPrefixExceeded
)

type PeerFSMCommand struct {
	Command int
	Reason  int
}

type FSMManager struct {
	logger         *logging.Writer
	neighborConf   *base.NeighborConf
	gConf          *config.GlobalConfig
	pConf          *config.NeighborConfig
	fsmConnCh      chan PeerFSMConn
	fsmStateCh     chan PeerFSMState
	peerAttrsCh    chan PeerAttrs
	bgpPktSrcCh    chan *packet.BGPPktSrc
	reachabilityCh chan config.ReachabilityInfo
	fsms           map[uint8]*FSM
	AcceptCh       chan net.Conn
	tcpConnFailCh  chan uint8
	CloseCh        chan bool
	StopFSMCh      chan string
	acceptConn     bool
	CommandCh      chan PeerFSMCommand
	BfdStatusCh    chan bool
	activeFSM      uint8
	newConnCh      chan PeerFSMConnState
	fsmMutex       sync.RWMutex
}

func NewFSMManager(logger *logging.Writer, neighborConf *base.NeighborConf, bgpPktSrcCh chan *packet.BGPPktSrc,
	fsmConnCh chan PeerFSMConn, reachabilityCh chan config.ReachabilityInfo) *FSMManager {
	mgr := FSMManager{
		logger:         logger,
		neighborConf:   neighborConf,
		gConf:          neighborConf.Global,
		pConf:          &neighborConf.RunningConf,
		fsmConnCh:      fsmConnCh,
		bgpPktSrcCh:    bgpPktSrcCh,
		reachabilityCh: reachabilityCh,
	}
	mgr.fsms = make(map[uint8]*FSM)
	mgr.AcceptCh = make(chan net.Conn)
	mgr.tcpConnFailCh = make(chan uint8, 2)
	mgr.acceptConn = false
	mgr.CloseCh = make(chan bool)
	mgr.StopFSMCh = make(chan string)
	mgr.CommandCh = make(chan PeerFSMCommand, 5)
	mgr.BfdStatusCh = make(chan bool, 4)
	mgr.activeFSM = uint8(config.ConnDirInvalid)
	mgr.newConnCh = make(chan PeerFSMConnState, 2)
	mgr.fsmMutex = sync.RWMutex{}
	return &mgr
}

func (mgr *FSMManager) Init() {
	fsmId := uint8(config.ConnDirOut)
	fsm := NewFSM(mgr, fsmId, mgr.neighborConf)
	fsm.Init(NewIdleState(fsm))
	go fsm.StartFSM()
	mgr.fsms[fsmId] = fsm
	fsm.passiveTcpEstCh <- true

	for {
		select {
		case inConn := <-mgr.AcceptCh:
			mgr.logger.Infof("Neighbor %s: Received a connection OPEN from far end", mgr.pConf.NeighborAddress)
			if !mgr.acceptConn {
				mgr.logger.Info("Can't accept connection from ", mgr.pConf.NeighborAddress, "yet.")
				inConn.Close()
			} else {
				foundInConn := false
				for _, fsm = range mgr.fsms {
					if fsm != nil && fsm.peerConn != nil && fsm.peerConn.dir == config.ConnDirIn {
						mgr.logger.Info("A FSM is already created for a incoming connection")
						foundInConn = true
						inConn.Close()
						break
					}
				}
				if !foundInConn {
					for fsmId, fsm = range mgr.fsms {
						if fsm != nil {
							mgr.logger.Infof("Neighbor %s: Send inConn message to FSM %d", mgr.pConf.NeighborAddress,
								fsmId)
							fsm.inConnCh <- inConn
							break
						}
					}
				}
			}

		case fsmId := <-mgr.tcpConnFailCh:
			mgr.logger.Infof("FSMManager: Neighbor %s: Received a TCP conn failed from FSM %d",
				mgr.pConf.NeighborAddress, fsmId)
			mgr.fsmTcpConnFailed(fsmId)

		case newConn := <-mgr.newConnCh:
			mgr.logger.Infof("FSMManager: Neighbor %s FSM %d Handle another connection", mgr.pConf.NeighborAddress,
				newConn.id)
			newId := mgr.getNewId(newConn.id)
			mgr.handleAnotherConnection(newId, newConn.connDir, newConn.conn)

		case stopMsg := <-mgr.StopFSMCh:
			mgr.StopFSM(stopMsg)

		case <-mgr.CloseCh:
			mgr.Cleanup()
			return

		case fsmCommand := <-mgr.CommandCh:
			event := BGPFSMEvent(fsmCommand.Command)
			mgr.logger.Infof("FSMManager: Neighbor %s: Received FSM command %d", mgr.pConf.NeighborAddress, event)
			if (event == BGPEventManualStart) || (event == BGPEventManualStop) || (event == BGPEventAutoStop) ||
				(event == BGPEventManualStartPassTcpEst) {
				for id, fsm := range mgr.fsms {
					if fsm != nil {
						mgr.logger.Infof("FSMManager: Neighbor %s: FSM %d Send command %d", mgr.pConf.NeighborAddress,
							id, event)
						fsm.eventRxCh <- PeerFSMEvent{event, fsmCommand.Reason}
					}
				}
			}

		case bfdStatus := <-mgr.BfdStatusCh:
			mgr.handleBfdStatusChange(bfdStatus)
		}
	}
}

func (mgr *FSMManager) handleBfdStatusChange(status bool) {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	for id, fsm := range mgr.fsms {
		if fsm != nil {
			mgr.logger.Infof("FSMManager: Neighbor %s: FSM %d Bfd status %d", mgr.pConf.NeighborAddress, id, status)
			fsm.bfdStatusCh <- status
		}
	}
}

func (mgr *FSMManager) AcceptPeerConn() {
	mgr.acceptConn = true
}

func (mgr *FSMManager) RejectPeerConn() {
	mgr.acceptConn = false
}

func (mgr *FSMManager) fsmTcpConnFailed(id uint8) {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	mgr.logger.Infof("FSMManager: Peer %s FSM %d TCP conn failed", mgr.pConf.NeighborAddress.String(), id)
	if len(mgr.fsms) != 1 && mgr.activeFSM != id {
		mgr.fsmClose(id)
	}
}

func (mgr *FSMManager) fsmClose(id uint8) {
	if closeFSM, ok := mgr.fsms[id]; ok {
		mgr.logger.Infof("FSMManager: Peer %s, close FSM %d", mgr.pConf.NeighborAddress.String(), id)
		closeFSM.closeCh <- true
		mgr.fsmBroken(id, false)
		mgr.fsms[id] = nil
		delete(mgr.fsms, id)
		mgr.logger.Infof("FSMManager: Peer %s, closed FSM %d", mgr.pConf.NeighborAddress.String(), id)
	} else {
		mgr.logger.Infof("FSMManager: Peer %s, FSM %d to close is not found in map %v",
			mgr.pConf.NeighborAddress.String(), id, mgr.fsms)
	}
}

func (mgr *FSMManager) fsmEstablished(id uint8, conn *net.Conn) {
	mgr.logger.Infof("FSMManager: Peer %s FSM %d connection established", mgr.pConf.NeighborAddress.String(), id)
	if _, ok := mgr.fsms[id]; ok {
		mgr.activeFSM = id
		mgr.fsmConnCh <- PeerFSMConn{mgr.neighborConf.Neighbor.NeighborAddress.String(), true, conn}
	} else {
		mgr.logger.Infof("FSMManager: Peer %s FSM %d not found in fsms dict %v", mgr.pConf.NeighborAddress.String(),
			id, mgr.fsms)
	}
	//mgr.Peer.PeerConnEstablished(conn)
}

func (mgr *FSMManager) fsmBroken(id uint8, fsmDelete bool) {
	mgr.logger.Infof("FSMManager: Peer %s FSM %d connection broken", mgr.pConf.NeighborAddress.String(), id)
	if mgr.activeFSM == id {
		mgr.activeFSM = uint8(config.ConnDirInvalid)
		mgr.fsmConnCh <- PeerFSMConn{mgr.neighborConf.Neighbor.NeighborAddress.String(), false, nil}
		//mgr.Peer.PeerConnBroken(fsmDelete)
	}
}

func (mgr *FSMManager) fsmStateChange(id uint8, state config.BGPFSMState) {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	if mgr.activeFSM == id || mgr.activeFSM == uint8(config.ConnDirInvalid) {
		mgr.neighborConf.FSMStateChange(uint32(state))
	}
}

func (mgr *FSMManager) SendUpdateMsg(bgpMsg *packet.BGPMessage) {
	mgr.fsmMutex.RLock()
	defer mgr.fsmMutex.RUnlock()

	if mgr.activeFSM == uint8(config.ConnDirInvalid) {
		mgr.logger.Infof("FSMManager: Neighbor %s FSM is not in ESTABLISHED state", mgr.pConf.NeighborAddress)
		return
	}
	mgr.logger.Infof("FSMManager: Neighbor %s FSM %d - send update", mgr.pConf.NeighborAddress, mgr.activeFSM)
	mgr.fsms[mgr.activeFSM].pktTxCh <- bgpMsg
}

func (mgr *FSMManager) Cleanup() {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	for id, fsm := range mgr.fsms {
		if fsm != nil {
			mgr.logger.Infof("FSMManager: Neighbor %s FSM %d - cleanup FSM", mgr.pConf.NeighborAddress, id)
			fsm.closeCh <- true
			fsm = nil
			mgr.fsmBroken(id, true)
			mgr.fsmStateChange(id, config.BGPFSMIdle)
			mgr.fsms[id] = nil
			delete(mgr.fsms, id)
		}
	}
}

func (mgr *FSMManager) StopFSM(stopMsg string) {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	for id, fsm := range mgr.fsms {
		if fsm != nil {
			mgr.logger.Infof("FSMManager: Neighbor %s FSM %d - Stop FSM", mgr.pConf.NeighborAddress, id)
			fsm.eventRxCh <- PeerFSMEvent{BGPEventTcpConnFails, BGPCmdReasonNone}
			mgr.fsmBroken(id, false)
		}
	}
}

func (mgr *FSMManager) getNewId(id uint8) uint8 {
	return uint8((id + 1) % 2)
}

func (mgr *FSMManager) createFSMForNewConnection(id uint8, connDir config.ConnDir) (*FSM, chan net.Conn) {
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	var state BaseStateIface

	if mgr.fsms[id] != nil {
		mgr.logger.Errf("FSMManager: Neighbor %s - FSM with id %d already exists", mgr.pConf.NeighborAddress, id)
		return nil, nil
	}

	mgr.logger.Infof("FSMManager: Neighbor %s Creating new FSM with id %d", mgr.pConf.NeighborAddress, id)
	fsm := NewFSM(mgr, id, mgr.neighborConf)

	state = NewActiveState(fsm)
	connCh := fsm.inConnCh
	if connDir == config.ConnDirOut {
		state = NewConnectState(fsm)
		connCh = fsm.outConnCh
	}
	fsm.Init(state)
	mgr.fsms[id] = fsm
	return fsm, connCh
}

func (mgr *FSMManager) handleAnotherConnection(id uint8, connDir config.ConnDir, conn *net.Conn) {
	fsm, connCh := mgr.createFSMForNewConnection(id, connDir)
	if fsm != nil {
		go fsm.StartFSM()
		fsm.passiveTcpEstCh <- true
		connCh <- *conn
	}
}

func (mgr *FSMManager) getFSMIdByDir(connDir config.ConnDir) uint8 {
	for id, fsm := range mgr.fsms {
		if fsm != nil && fsm.peerConn != nil && fsm.peerConn.dir == connDir {
			return id
		}
	}

	return uint8(config.ConnDirInvalid)
}

func (mgr *FSMManager) receivedBGPOpenMessage(id uint8, connDir config.ConnDir, openMsg *packet.BGPOpen) bool {
	var closeConnDir config.ConnDir = config.ConnDirInvalid
	var closeFSMId uint8 = uint8(config.ConnDirInvalid)
	mgr.fsmMutex.Lock()
	defer mgr.fsmMutex.Unlock()

	mgr.logger.Infof("FSMManager - Neighbor %s: FSM %d rx OPEN message", mgr.pConf.NeighborAddress, id)

	localBGPId := packet.ConvertIPBytesToUint(mgr.gConf.RouterId.To4())
	bgpIdInt := packet.ConvertIPBytesToUint(openMsg.BGPId.To4())
	for fsmId, fsm := range mgr.fsms {
		if fsmId != id && fsm != nil && fsm.State.state() >= config.BGPFSMOpensent {
			if fsm.State.state() == config.BGPFSMEstablished {
				closeConnDir = connDir
			} else if localBGPId > bgpIdInt {
				closeConnDir = config.ConnDirIn
			} else {
				closeConnDir = config.ConnDirOut
			}
			closeFSMId := mgr.getFSMIdByDir(closeConnDir)
			mgr.logger.Infof("FSMManager - Neighbor %s: Close FSM id %d", mgr.pConf.NeighborAddress, closeFSMId)
			mgr.fsmClose(closeFSMId)
		}
	}

	if closeFSMId == uint8(config.ConnDirInvalid) || closeFSMId != id {
		asSize := packet.GetASSize(openMsg)
		addPathFamily := packet.GetAddPathFamily(openMsg)
		if mgr.fsms[id] != nil {
			mgr.logger.Infof("FSMManager - Neighbor %s: FSM %d set peer attr", mgr.pConf.NeighborAddress, id)
			mgr.neighborConf.SetPeerAttrs(openMsg.BGPId, asSize, mgr.fsms[id].holdTime, mgr.fsms[id].keepAliveTime,
				addPathFamily)
		}
	}

	if closeConnDir == connDir {
		mgr.logger.Infof("FSMManager - Neighbor %s: FSM %d Closing FSM... return false",
			mgr.pConf.NeighborAddress.String(), id)
		return false
	} else {
		return true
	}
}
