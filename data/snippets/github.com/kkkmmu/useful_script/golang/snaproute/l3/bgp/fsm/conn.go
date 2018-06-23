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

// conn.go
package fsm

import (
	"errors"
	"fmt"
	"l3/bgp/config"
	"l3/bgp/packet"
	"l3/bgp/utils"
	"math/rand"
	"net"
	"strings"
	"time"
	"utils/logging"
	"utils/netUtils"

	"golang.org/x/net/ipv4"
)

type OutTCPConn struct {
	fsm          *FSM
	logger       *logging.Writer
	ifaceMgr     *utils.InterfaceMgr
	fsmConnCh    chan net.Conn
	fsmConnErrCh chan PeerConnErr
	StopConnCh   chan bool
	id           uint32
}

func NewOutTCPConn(fsm *FSM, fsmConnCh chan net.Conn, fsmConnErrCh chan PeerConnErr) *OutTCPConn {
	outConn := OutTCPConn{
		fsm:          fsm,
		logger:       fsm.logger,
		ifaceMgr:     utils.NewInterfaceMgr(fsm.logger),
		fsmConnCh:    fsmConnCh,
		fsmConnErrCh: fsmConnErrCh,
		StopConnCh:   make(chan bool, 2),
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	outConn.id = r.Uint32()
	fsm.logger.Info("Neighbor:", fsm.pConf.NeighborAddress, "FSM", fsm.id, "Creating new out TCP conn with id",
		outConn.id)
	return &outConn
}

func (o *OutTCPConn) Connect(seconds uint32, remote, local string, connCh chan net.Conn, errCh chan error) {
	o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "Connect start, local IP:", local,
		"remote IP:", remote)
	remoteIP, _, err := net.SplitHostPort(remote)
	if err != nil {
		errCh <- err
		return
	}

	reachableCh := make(chan config.ReachabilityResult)
	reachabilityInfo := config.ReachabilityInfo{
		IP:          o.fsm.pConf.NeighborAddress.String(),
		ReachableCh: reachableCh,
		IfIndex:     o.fsm.pConf.IfIndex,
	}
	o.fsm.Manager.reachabilityCh <- reachabilityInfo
	reachable := <-reachableCh
	if reachable.Err != nil {
		duration := uint32(3)
		if (duration * 2) < seconds {
			for {
				select {
				case <-time.After(time.Duration(duration) * time.Second):
					o.fsm.Manager.reachabilityCh <- reachabilityInfo
					reachable = <-reachableCh
				}
				seconds -= duration
				if reachable.Err == nil || seconds <= (duration*2) {
					break
				}
			}
		}
		if reachable.Err != nil {
			errCh <- config.AddressNotResolvedError{"Neighbor is not reachable"}
			return
		}
	}

	if local == "" && reachable.NextHopInfo != nil {
		nextHopIP := net.ParseIP(strings.TrimSpace(reachable.NextHopInfo.NextHopIp))
		if nextHopIP != nil && (nextHopIP.Equal(net.IPv4zero) || nextHopIP.Equal(net.IPv6zero)) {
			o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "Next hop ip", nextHopIP,
				"is 0, Set source ip for the TCP connection to", reachable.NextHopInfo.IPAddr)
			local = net.JoinHostPort(strings.TrimSpace(reachable.NextHopInfo.IPAddr), "0")
		}
	}

	if local != "" {
		o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "local IP is set to", local,
			"local IP len", len(local))
		localIP, _, err := net.SplitHostPort(local)
		if err != nil {
			o.logger.Err("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "SplitHostPort for local IP",
				local, "failed with error", err)
			errCh <- err
			return
		}

		if strings.TrimSpace(localIP) != "" {
			if !o.ifaceMgr.IsIPConfigured(strings.TrimSpace(localIP)) {
				errCh <- errors.New(fmt.Sprintf("Local IP %s is not configured on the switch", localIP))
				return
			}
		} else {
			o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
				"local IP is empty, set the local address to empty")
			local = ""
		}
	}

	o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
		"Connect called... calling DialTimeout with", seconds, "second timeout", "OutTCPCOnn id", o.id)
	socket, err := netUtils.ConnectSocket("tcp", remote, local)
	defer netUtils.CloseSocket(socket)
	if err != nil {
		o.logger.Err("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "ConnectSocket failed with error", err)
		errCh <- err
		return
	}

	if o.fsm.pConf.AuthPassword != "" {
		o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "Set MD5 option on the socket:",
			socket, "password:", o.fsm.pConf.AuthPassword)
		err = netUtils.SetSockoptTCPMD5(socket, remoteIP, o.fsm.pConf.AuthPassword)
		if err != nil {
			o.logger.Err("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
				"Set MD5 option on the socket failed with error", err)
			errCh <- err
			return
		}
	}

	duration := uint32(10)
	if duration < seconds {
		duration = seconds
	}

	err = netUtils.Connect(socket, "tcp", remote, local, time.Duration(duration)*time.Second)
	if err != nil {
		o.logger.Err("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "Connect failed with error", err)
		errCh <- err
		return
	}

	conn, err := netUtils.ConvertFdToConn(socket)
	if err != nil {
		errCh <- err
	} else {
		packetConn := ipv4.NewConn(conn)
		ttl := 1
		if o.fsm.pConf.MultiHopEnable {
			ttl = int(o.fsm.pConf.MultiHopTTL)
		}
		if err = packetConn.SetTTL(ttl); err != nil {
			conn.Close()
			errCh <- err
			return
		}
		connCh <- conn
	}
}

func (o *OutTCPConn) ConnectToPeer(seconds uint32, remote, local string) {
	var stopConn bool = false
	connCh := make(chan net.Conn)
	errCh := make(chan error)

	o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id, "ConnectToPeer called", "OutTCPCOnn id",
		o.id)
	connTime := seconds - 3
	if connTime <= 0 {
		connTime = seconds
	}

	done := false
	go o.Connect(seconds, remote, local, connCh, errCh)

	for {
		select {
		case conn := <-connCh:
			o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
				"ConnectToPeer: Connected to peer", remote, "OutTCPCOnn id", o.id)
			if stopConn {
				conn.Close()
				return
			}

			done = true
			o.fsmConnCh <- conn

		case err := <-errCh:
			o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
				"ConnectToPeer: Failed to connect to peer", remote, "with error:", err, "OutTCPCOnn id", o.id)
			if stopConn {
				return
			}

			done = true
			o.fsmConnErrCh <- PeerConnErr{0, err}

		case <-o.StopConnCh:
			o.logger.Info("Neighbor:", o.fsm.pConf.NeighborAddress, "FSM", o.fsm.id,
				"ConnectToPeer: Recieved stop connecting to peer", remote, "OutTCPCOnn id", o.id)
			if done {
				return
			}
			stopConn = true
		}
	}
}

type PeerConn struct {
	fsm       *FSM
	logger    *logging.Writer
	dir       config.ConnDir
	conn      *net.Conn
	id        uint32
	peerAttrs packet.BGPPeerAttrs

	readCh chan bool
	stopCh chan bool
	exitCh chan bool
}

func NewPeerConn(fsm *FSM, dir config.ConnDir, conn *net.Conn, id uint32) *PeerConn {
	peerConn := PeerConn{
		fsm:    fsm,
		logger: fsm.logger,
		dir:    dir,
		conn:   conn,
		id:     id,
		peerAttrs: packet.BGPPeerAttrs{
			ASSize:           2,
			AddPathsRxActual: false,
		},
		readCh: make(chan bool),
		stopCh: make(chan bool),
		exitCh: make(chan bool),
	}

	return &peerConn
}

func (p *PeerConn) StartReading() {
	stopReading := false
	readError := false
	doneReadingCh := make(chan bool)
	stopReadingCh := make(chan bool)
	exitCh := make(chan bool)

	p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "conn:StartReading called")
	go p.ReadPkt(doneReadingCh, stopReadingCh, exitCh)
	p.readCh <- true

	for {
		select {
		case <-p.stopCh:
			stopReading = true
			if readError {
				p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
					"readError is true, send stopReadingCh")
				stopReadingCh <- true
			}

		case <-exitCh:
			p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "conn: exit channel")
			(*p.conn).Close()
			p.exitCh <- true
			return

		case readOk := <-doneReadingCh:
			if stopReading {
				p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
					"stopReading is true, send stopReadingCh")
				stopReadingCh <- true
			} else {
				if readOk {
					p.readCh <- true
				} else {
					p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
						"read failed, set readError to true")
					readError = true
				}
			}
		}
	}
}

func (p *PeerConn) StopReading(exitCh chan bool) {
	p.exitCh = exitCh
	p.stopCh <- true
}

func (p *PeerConn) readPartialPkt(length int) ([]byte, error) {
	buf := make([]byte, length)
	var totalRead int = 0
	var read int = 0
	var err error
	for totalRead < length {
		read, err = (*p.conn).Read(buf[totalRead:])
		if err != nil {
			return buf, err
		}
		totalRead += read
		//p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "conn:readPartialPkt -",
		//	"read", read, "bytes, total read", totalRead, "bytes, lenght =", length))
	}
	return buf, err
}

func (p *PeerConn) DecodeMessage(header *packet.BGPHeader, buf []byte) (*packet.BGPMessage, *packet.BGPMessageError,
	bool) {
	var msgErr *packet.BGPMessageError
	msg := packet.NewBGPMessage()
	err := msg.Decode(header, buf, p.peerAttrs)
	//bgpPktInfo := packet.NewBGPPktInfo(msg, nil)
	msgOk := true
	if header.Type == packet.BGPMsgTypeNotification {
		msgOk = false
	}

	if err != nil {
		p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "BGP packet body decode failed, err:",
			err)
		//bgpPktInfo = packet.NewBGPPktInfo(msg, err.(*packet.BGPMessageError))
		bgpErr := err.(packet.BGPMessageError)
		msgErr = &bgpErr
		msgOk = false
	} else if header.Type == packet.BGPMsgTypeOpen {
		peerAS := packet.GetPeerAS(msg.Body.(*packet.BGPOpen))
		if peerAS != p.fsm.pConf.PeerAS {
			msgOk = false
			msgErr = &packet.BGPMessageError{packet.BGPOpenMsgError, packet.BGPBadPeerAS, nil, "Bad peer AS"}
		} else {
			p.peerAttrs.ASSize = packet.GetASSize(msg.Body.(*packet.BGPOpen))
			p.peerAttrs.AddPathFamily = packet.GetAddPathFamily(msg.Body.(*packet.BGPOpen))
			addPathsTxFarEnd := packet.IsAddPathsTxEnabledForIPv4(p.peerAttrs.AddPathFamily)
			p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "Far end can send add paths")
			if addPathsTxFarEnd && p.fsm.pConf.AddPathsRx {
				p.peerAttrs.AddPathsRxActual = true
				p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "negotiated to recieve add paths from far end")
			}
		}
	}

	return msg, msgErr, msgOk
}

func (p *PeerConn) ReadPkt(doneCh chan bool, stopCh chan bool, exitCh chan bool) {
	p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "conn:ReadPkt called")
	var t time.Time
	var header *packet.BGPHeader
	for {
		select {
		case <-p.readCh:
			header = nil
			(*p.conn).SetReadDeadline(time.Now().Add(time.Duration(3) * time.Second))
			buf, err := p.readPartialPkt(int(packet.BGPMsgHeaderLen))
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					doneCh <- true
					continue
				} else {
					p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
						"readPartialPkt DID NOT time out, returned err:", err, "nerr:", nerr)
					p.fsm.outConnErrCh <- PeerConnErr{p.id, err}
					doneCh <- false
					break
				}
			}

			header = packet.NewBGPHeader()
			err = header.Decode(buf)
			if err != nil {
				p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
					"BGP packet header decode failed")
				//bgpPktInfo := packet.NewBGPPktInfo(nil, err.(*packet.BGPMessageError))
				bgpErr := err.(packet.BGPMessageError)
				p.fsm.pktRxCh <- packet.NewBGPPktInfo(nil, &bgpErr)
				doneCh <- false
				continue
			}

			if header.Type != packet.BGPMsgTypeKeepAlive {
				p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id,
					"Recieved BGP packet type=", header.Type, "len=", header.Len())
			}

			(*p.conn).SetReadDeadline(t)
			if header.Len() > packet.BGPMsgHeaderLen {
				buf, err = p.readPartialPkt(int(header.Len() - packet.BGPMsgHeaderLen))
				if err != nil {
					p.fsm.outConnErrCh <- PeerConnErr{p.id, err}
					doneCh <- false
					break
				}
			} else {
				buf = make([]byte, 0)
			}

			if header.Type != packet.BGPMsgTypeKeepAlive {
				p.logger.Infof("Neighbor:%s FSM %d Received BGP packet %x", p.fsm.pConf.NeighborAddress, p.fsm.id, buf)
			}

			msg, msgErr, msgOk := p.DecodeMessage(header, buf)
			p.fsm.pktRxCh <- packet.NewBGPPktInfo(msg, msgErr)
			doneCh <- msgOk

		case <-stopCh:
			p.logger.Info("Neighbor:", p.fsm.pConf.NeighborAddress, "FSM", p.fsm.id, "Closing the peer connection")
			if p.conn != nil {
				(*p.conn).Close()
			}
			exitCh <- true
			return
		}
	}
}
