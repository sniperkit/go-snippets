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

package server

import (
	"asicd/asicdCommonDefs"
	"errors"
	"l3/bfd/bfddCommonDefs"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils/commonDefs"
)

const (
	PACKET_QUEUE_SIZE = 128
	MAX_NUM_SESSIONS  = 1024
)

func (server *BFDServer) StartSessionHandler() error {
	server.CreateSessionCh = make(chan BfdSessionMgmt)
	server.DeleteSessionCh = make(chan BfdSessionMgmt)
	server.AdminUpSessionCh = make(chan BfdSessionMgmt)
	server.AdminDownSessionCh = make(chan BfdSessionMgmt)
	server.ResetSessionCh = make(chan int32)
	server.CreatedSessionCh = make(chan int32, MAX_NUM_SESSIONS)
	server.FailedSessionClientCh = make(chan int32, MAX_NUM_SESSIONS)
	server.tobeCreatedSessions = make(map[string]BfdSessionMgmt)
	go server.StartBfdSesionServer()
	go server.StartBfdSessionRxTx()
	go server.StartSessionRetryHandler()
	for {
		select {
		case sessionMgmt := <-server.CreateSessionCh:
			server.CreateBfdSession(sessionMgmt)
		case sessionMgmt := <-server.DeleteSessionCh:
			server.DeleteBfdSession(sessionMgmt)
		case sessionMgmt := <-server.AdminUpSessionCh:
			server.AdminUpBfdSession(sessionMgmt)
		case sessionMgmt := <-server.AdminDownSessionCh:
			server.AdminDownBfdSession(sessionMgmt)
		case sessionId := <-server.ResetSessionCh:
			server.ResetBfdSession(sessionId)

		}
	}
	return nil
}

func (server *BFDServer) DispatchReceivedBfdPacket(ipAddr string, bfdPacket *BfdControlPacket) error {
	sessionId := int32(bfdPacket.YourDiscriminator)
	session, exist := server.bfdGlobal.Sessions[sessionId]
	if !exist {
		session, exist = server.bfdGlobal.SessionsByIp[ipAddr]
	}
	if exist && session != nil {
		session.sessionLock.Lock()
		if session.IsSessionActive() && session.rxInterval != 0 {
			if session.sessionTimer != nil {
				session.sessionTimer.Reset(time.Duration(session.rxInterval) * time.Millisecond)
				session.ReceivedPacketCh <- bfdPacket
			}
		}
		session.sessionLock.Unlock()
	} else {
		/*
			// Create a session as discovered. This can be enabled for active mode of bfd.
			server.logger.Info("No session found for ", ipAddr, " creating a session")
			sessionMgmt := BfdSessionMgmt{
				DestIp:   ipAddr,
				Protocol: bfddCommonDefs.DISC,
				PerLink:  false,
			}
			server.CreateSessionCh <- sessionMgmt
		*/
	}
	return nil
}

func (server *BFDServer) StartBfdSesionServer() error {
	var ipAddr string
	var bfdPacket *BfdControlPacket
	var err error
	destAddr := net.JoinHostPort("", strconv.Itoa(DEST_PORT))
	ServerAddr, err := net.ResolveUDPAddr("udp", destAddr)
	if err != nil {
		server.logger.Info("Failed ResolveUDPAddr ", destAddr, err)
		return err
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		server.logger.Info("Failed ListenUDP ", err)
		return err
	}
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	server.logger.Info("Started BFD session server on ", destAddr)
	for {
		length, udpAddr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			server.logger.Info("Failed to read from ", ServerAddr)
		} else {
			if length >= DEFAULT_CONTROL_PACKET_LEN {
				bfdPacket, err = DecodeBfdControlPacket(buf[0:length])
				if err != nil {
					server.logger.Info("Failed to decode packet - ", err)
					continue
				} else {
					ipAddr, _, err = net.SplitHostPort(udpAddr.String())
					if err == nil {
						err = server.DispatchReceivedBfdPacket(ipAddr, bfdPacket)
						if err != nil {
							server.logger.Info("Failed to dispatch received packet")
						}
					}
				}
			}
		}
	}
	return nil
}

func (server *BFDServer) StartBfdSessionRxTx() error {
	for {
		select {
		case createdSessionId := <-server.CreatedSessionCh:
			session := server.bfdGlobal.Sessions[createdSessionId]
			if session != nil {
				session.SessionStopClientCh = make(chan bool)
				session.SessionStopServerCh = make(chan bool)
				session.ReceivedPacketCh = make(chan *BfdControlPacket, PACKET_QUEUE_SIZE)
				session.isClientActive = false
				if session.state.PerLinkSession {
					server.logger.Info("Starting PerLink server for session ", createdSessionId)
					go session.StartPerLinkSessionServer(server)
					server.logger.Info("Starting PerLink client for session ", createdSessionId)
					go session.StartPerLinkSessionClient(server)
				} else {
					server.logger.Info("Starting server for session ", createdSessionId)
					go session.StartSessionServer()
					server.logger.Info("Starting client for session ", createdSessionId)
					go session.StartSessionClient(server)
				}
			} else {
				server.logger.Info("Bfd session could not be initiated for ", createdSessionId)
			}
		case failedClientSessionId := <-server.FailedSessionClientCh:
			session := server.bfdGlobal.Sessions[failedClientSessionId]
			if session != nil {
				server.bfdGlobal.InactiveSessionsIdSlice = append(server.bfdGlobal.InactiveSessionsIdSlice, failedClientSessionId)
			}
		}
	}
	return nil
}

func (server *BFDServer) StartSessionRetryHandler() error {
	server.logger.Info("Starting session retry handler")
	retryTimer := time.NewTicker(time.Second * 5)
	for t := range retryTimer.C {
		_ = t
		for i := 0; i < len(server.bfdGlobal.InactiveSessionsIdSlice); i++ {
			sessionId := server.bfdGlobal.InactiveSessionsIdSlice[i]
			session := server.bfdGlobal.Sessions[sessionId]
			if session != nil {
				server.logger.Info("Session retry handler restarting session", sessionId, "active", session.isClientActive)
				if session.isClientActive == false {
					if session.state.PerLinkSession {
						server.logger.Info("Starting PerLink client for inactive session ", sessionId)
						go session.StartPerLinkSessionClient(server)
					} else {
						server.logger.Info("Starting client for inactive session ", sessionId)
						go session.StartSessionClient(server)
					}
				}
				server.bfdGlobal.InactiveSessionsIdSlice = append(server.bfdGlobal.InactiveSessionsIdSlice[:i], server.bfdGlobal.InactiveSessionsIdSlice[i+1:]...)
			}
		}
	}
	server.logger.Info("Session retry handler exiting ...")
	return nil
}

func (server *BFDServer) processSessionConfig(sessionConfig SessionConfig) error {
	sessionMgmt := BfdSessionMgmt{
		DestIp:    sessionConfig.DestIp,
		ParamName: sessionConfig.ParamName,
		Interface: sessionConfig.Interface,
		Protocol:  sessionConfig.Protocol,
		PerLink:   sessionConfig.PerLink,
	}
	switch sessionConfig.Operation {
	case bfddCommonDefs.CREATE:
		server.CreateSessionCh <- sessionMgmt
	case bfddCommonDefs.DELETE:
		server.DeleteSessionCh <- sessionMgmt
	case bfddCommonDefs.ADMINUP:
		server.AdminUpSessionCh <- sessionMgmt
	case bfddCommonDefs.ADMINDOWN:
		server.AdminDownSessionCh <- sessionMgmt
	}
	return nil
}

func (server *BFDServer) SendAdminUpToAllNeighbors() error {
	for _, session := range server.bfdGlobal.Sessions {
		session.StartBfdSession()
	}
	return nil
}

func (server *BFDServer) SendAdminDownToAllNeighbors() error {
	for _, session := range server.bfdGlobal.Sessions {
		session.StopBfdSession()
	}
	return nil
}

func (server *BFDServer) SendDeleteToAllSessions() error {
	for _, session := range server.bfdGlobal.Sessions {
		session.SessionStopClientCh <- true
		session.SessionStopServerCh <- true
	}
	return nil
}

func (server *BFDServer) GetNewSessionId() int32 {
	var sessionIdUsed bool
	var sessionId int32
	sessionId = 0
	if server.bfdGlobal.NumSessions < MAX_NUM_SESSIONS {
		sessionIdUsed = true //By default assume the sessionId is already used.
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		for sessionIdUsed {
			sessionId = r1.Int31n(MAX_NUM_SESSIONS)
			if _, exist := server.bfdGlobal.Sessions[sessionId]; exist {
				server.logger.Info("GetNewSessionId: sessionId ", sessionId, " is in use, Generating a new one")
			} else {
				if sessionId != 0 {
					sessionIdUsed = false
				}
			}
		}
	}
	return sessionId
}

func (server *BFDServer) GetIfIndexFromDestIp(DestIp string) (int32, string, error) {
	var ifIndex int32
	var localAddr string
	server.ribdClient.ClientHdl.TrackReachabilityStatus(DestIp, "BFD", "add")
	reachabilityInfo, err := server.ribdClient.ClientHdl.GetRouteReachabilityInfo(DestIp, -1)
	server.logger.Info("Reachability info ", reachabilityInfo)
	if err != nil || !reachabilityInfo.IsReachable {
		err = errors.New(DestIp + " is not reachable")
	} else {
		ifIndex = int32(reachabilityInfo.NextHopIfIndex)
		localAddr = reachabilityInfo.Ipaddr
		server.logger.Info("GetIfIndexFromDestIp: DestIp: ", DestIp, "IfIndex: ", ifIndex, "LocalIP: ", localAddr)
	}
	return ifIndex, localAddr, err
}

func (server *BFDServer) GetTxJitter() int32 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	jitter := r1.Int31n(TX_JITTER)
	return jitter
}

func (server *BFDServer) NewNormalBfdSession(Interface string, LocalIp string, DestIp string, ParamName string, PerLink bool, Protocol bfddCommonDefs.BfdSessionOwner) *BfdSession {
	bfdSession := &BfdSession{}
	sessionId := server.GetNewSessionId()
	if sessionId == 0 {
		server.logger.Info("Failed to get sessionId")
		return nil
	}
	bfdSession.state.SessionId = sessionId
	bfdSession.state.IpAddr = DestIp
	bfdSession.state.LocalAddr = LocalIp
	bfdSession.state.Interface = Interface
	bfdSession.state.PerLinkSession = PerLink
	if PerLink {
		bfdSession.state.LocalMacAddr, _ = server.getMacAddrFromIntfName(Interface)
		bfdSession.state.RemoteMacAddr, _ = net.ParseMAC(bfdDedicatedMac)
		bfdSession.useDedicatedMac = true
	}
	bfdSession.state.RegisteredProtocols = make([]bool, bfddCommonDefs.MAX_APPS)
	bfdSession.state.RegisteredProtocols[Protocol] = true
	bfdSession.state.SessionState = STATE_DOWN
	bfdSession.state.RemoteSessionState = STATE_DOWN
	bfdSession.state.LocalDiscriminator = uint32(bfdSession.state.SessionId)
	bfdSession.state.LocalDiagType = DIAG_NONE
	bfdSession.txInterval = STARTUP_TX_INTERVAL / 1000
	bfdSession.txJitter = server.GetTxJitter()
	sessionParam, exist := server.bfdGlobal.SessionParams[ParamName]
	if exist {
		bfdSession.state.ParamName = ParamName
	} else {
		bfdSession.state.ParamName = "default"
		sessionParam, _ = server.bfdGlobal.SessionParams["default"]
	}
	sessionParam.state.NumSessions++
	bfdSession.rxInterval = (STARTUP_RX_INTERVAL * sessionParam.state.LocalMultiplier) / 1000
	bfdSession.state.DesiredMinTxInterval = sessionParam.state.DesiredMinTxInterval
	bfdSession.state.RequiredMinRxInterval = sessionParam.state.RequiredMinRxInterval
	bfdSession.state.DetectionMultiplier = sessionParam.state.LocalMultiplier
	bfdSession.state.DemandMode = sessionParam.state.DemandEnabled
	bfdSession.authEnabled = sessionParam.state.AuthenticationEnabled
	bfdSession.authType = AuthenticationType(sessionParam.state.AuthenticationType)
	bfdSession.authSeqNum = 1
	bfdSession.authKeyId = uint32(sessionParam.state.AuthenticationKeyId)
	bfdSession.authData = sessionParam.state.AuthenticationData
	bfdSession.paramChanged = true
	bfdSession.bfdPacket = NewBfdControlPacketDefault()
	bfdSession.server = server
	bfdSession.sessionLock = sync.RWMutex{}
	server.bfdGlobal.Sessions[sessionId] = bfdSession
	server.bfdGlobal.SessionsByIp[DestIp] = bfdSession
	server.bfdGlobal.NumSessions++
	server.bfdGlobal.SessionsIdSlice = append(server.bfdGlobal.SessionsIdSlice, sessionId)
	server.logger.Info("New session : ", sessionId, " created on : ", Interface)
	server.CreatedSessionCh <- sessionId
	return bfdSession
}

func (server *BFDServer) NewPerLinkBfdSessions(IfIndex int32, LocalIp string, DestIp string, ParamName string, Protocol bfddCommonDefs.BfdSessionOwner) error {
	lag, exist := server.lagPropertyMap[IfIndex]
	if exist {
		for _, link := range lag.Links {
			IfName, _ := server.getLinuxIntfName(IfIndex)
			bfdSession := server.NewNormalBfdSession(IfName, LocalIp, DestIp, ParamName, true, Protocol)
			if bfdSession == nil {
				server.logger.Info("Failed to create perlink session on ", link)
			}
		}
	} else {
		server.logger.Info("Unknown lag ", IfIndex, " can not create perlink sessions")
	}
	return nil
}

func (server *BFDServer) NewBfdSession(DestIp string, ParamName string, Interface string, Protocol bfddCommonDefs.BfdSessionOwner, PerLink bool) *BfdSession {
	var IfType int
	var interfaceSpecific bool
	if Interface != "None" && Interface != "" {
		interfaceSpecific = true
	}
	IfIndex, localIp, err := server.GetIfIndexFromDestIp(DestIp)
	if err != nil {
		server.logger.Err("Session creation failed " + err.Error())
		return nil
	} else {
		IfType = asicdCommonDefs.GetIntfTypeFromIfIndex(IfIndex)
		IfName, err := server.getLinuxIntfName(IfIndex)
		if err == nil {
			if interfaceSpecific && IfName != Interface {
				server.logger.Info("Bfd session to ", DestIp, " cannot be created on interface ", Interface)
				return nil
			}
		}
	}
	if IfType == commonDefs.IfTypeLag && PerLink {
		server.NewPerLinkBfdSessions(IfIndex, localIp, DestIp, ParamName, Protocol)
	} else {
		bfdSession := server.NewNormalBfdSession(Interface, localIp, DestIp, ParamName, false, Protocol)
		bfdSession.state.InterfaceSpecific = interfaceSpecific
		return bfdSession
	}
	return nil
}

func (server *BFDServer) UpdateBfdSessionsUsingParam(paramName string) error {
	sessionParam, paramExist := server.bfdGlobal.SessionParams[paramName]
	for _, session := range server.bfdGlobal.Sessions {
		if session.state.ParamName == paramName {
			if paramExist {
				session.state.DesiredMinTxInterval = sessionParam.state.DesiredMinTxInterval
				session.state.RequiredMinRxInterval = sessionParam.state.RequiredMinRxInterval
				session.state.DetectionMultiplier = sessionParam.state.LocalMultiplier
				session.state.DemandMode = sessionParam.state.DemandEnabled
				session.authEnabled = sessionParam.state.AuthenticationEnabled
				session.authType = AuthenticationType(sessionParam.state.AuthenticationType)
				session.authKeyId = uint32(sessionParam.state.AuthenticationKeyId)
				session.authData = sessionParam.state.AuthenticationData
			} else {
				session.state.DesiredMinTxInterval = DEFAULT_DESIRED_MIN_TX_INTERVAL
				session.state.RequiredMinRxInterval = DEFAULT_REQUIRED_MIN_RX_INTERVAL
				session.state.DetectionMultiplier = DEFAULT_DETECT_MULTI
				session.state.DemandMode = false
				session.authEnabled = false
			}
			session.paramChanged = true
			session.InitiatePollSequence()
		}
	}
	return nil
}

func (server *BFDServer) FindBfdSession(DestIp string) (sessionId int32, found bool) {
	found = false
	for sessionId, session := range server.bfdGlobal.Sessions {
		if session.state.IpAddr == DestIp {
			return sessionId, true
		}
	}
	return sessionId, found
}

func (server *BFDServer) FindBfdSessionContainingAddr(DestIp string) (sessionId int32, found bool) {
	found = false
	for sessionId, session := range server.bfdGlobal.Sessions {
		if strings.HasPrefix(session.state.IpAddr, DestIp) {
			return sessionId, true
		}
	}
	return sessionId, found
}

func NewBfdControlPacketDefault() *BfdControlPacket {
	bfdControlPacket := &BfdControlPacket{
		Version:    DEFAULT_BFD_VERSION,
		Diagnostic: DIAG_NONE,
		State:      STATE_DOWN,
		Poll:       false,
		Final:      false,
		ControlPlaneIndependent:   false,
		AuthPresent:               false,
		Demand:                    false,
		Multipoint:                false,
		DetectMult:                DEFAULT_DETECT_MULTI,
		MyDiscriminator:           0,
		YourDiscriminator:         0,
		DesiredMinTxInterval:      DEFAULT_DESIRED_MIN_TX_INTERVAL,
		RequiredMinRxInterval:     DEFAULT_REQUIRED_MIN_RX_INTERVAL,
		RequiredMinEchoRxInterval: DEFAULT_REQUIRED_MIN_ECHO_RX_INTERVAL,
		AuthHeader:                nil,
	}
	return bfdControlPacket
}

// CreateBfdSession initializes a session and starts cpntrol packets exchange.
// This function is called when a protocol registers with BFD to monitor a destination IP.
func (server *BFDServer) CreateBfdSession(sessionMgmt BfdSessionMgmt) (*BfdSession, error) {
	var bfdSession *BfdSession
	var err error
	DestIp := sessionMgmt.DestIp
	ParamName := sessionMgmt.ParamName
	Interface := sessionMgmt.Interface
	Protocol := sessionMgmt.Protocol
	PerLink := sessionMgmt.PerLink
	sessionIp := DestIp
	if Interface != "" {
		ipAddr := net.ParseIP(DestIp)
		if ipAddr.To4() == nil {
			if ipAddr.IsLinkLocalUnicast() {
				sessionIp = sessionIp + "%" + Interface
			}
		}
	}
	sessionId, found := server.FindBfdSession(sessionIp)
	if !found {
		server.logger.Info("CreateSession ", sessionIp, ParamName, Interface, Protocol, PerLink)
		bfdSession = server.NewBfdSession(sessionIp, ParamName, Interface, Protocol, PerLink)
		if bfdSession != nil {
			server.logger.Info("Bfd session created ", bfdSession.state.SessionId, bfdSession.state.IpAddr)
		} else {
			server.logger.Info("CreateSession failed for ", sessionIp, Protocol)
			// Store the session config in tobeCrearedSessions map.
			if _, exist := server.tobeCreatedSessions[sessionIp]; !exist {
				server.tobeCreatedSessions[sessionIp] = sessionMgmt
				server.logger.Info("Stored session config for ", sessionIp, " waiting for reachability")
			}
			err = errors.New("Failed to create session to " + sessionIp)
		}
	} else {
		server.logger.Info("Bfd session already exists ", sessionIp, Protocol, sessionId)
		bfdSession = server.bfdGlobal.Sessions[sessionId]
		if !bfdSession.state.RegisteredProtocols[Protocol] {
			bfdSession.state.RegisteredProtocols[Protocol] = true
		}
	}
	return bfdSession, err
}

func (server *BFDServer) SessionDeleteHandler(session *BfdSession, Protocol bfddCommonDefs.BfdSessionOwner, ForceDel bool) error {
	var i int
	sessionId := session.state.SessionId
	session.state.RegisteredProtocols[Protocol] = false
	if ForceDel || session.CheckIfAnyProtocolRegistered() == false {
		if session.IsSessionActive() {
			session.SessionStopClientCh <- true
		}
		session.SessionStopServerCh <- true
		server.bfdGlobal.SessionParams[session.state.ParamName].state.NumSessions--
		server.bfdGlobal.NumSessions--
		delete(server.bfdGlobal.Sessions, sessionId)
		delete(server.bfdGlobal.SessionsByIp, session.state.IpAddr)
		for i = 0; i < len(server.bfdGlobal.SessionsIdSlice); i++ {
			if server.bfdGlobal.SessionsIdSlice[i] == sessionId {
				break
			}
		}
		server.bfdGlobal.SessionsIdSlice = append(server.bfdGlobal.SessionsIdSlice[:i], server.bfdGlobal.SessionsIdSlice[i+1:]...)
		server.logger.Info("Deleted session ", sessionId)
	}
	return nil
}

func (server *BFDServer) DeletePerLinkSessions(DestIp string, Protocol bfddCommonDefs.BfdSessionOwner, ForceDel bool) error {
	for _, session := range server.bfdGlobal.Sessions {
		if session.state.IpAddr == DestIp {
			server.SessionDeleteHandler(session, Protocol, ForceDel)
		}
	}
	return nil
}

// DeleteBfdSession ceases the session.
// A session down control packet is sent to BFD neighbor before deleting the session.
// This function is called when a protocol decides to stop monitoring the destination IP.
func (server *BFDServer) DeleteBfdSession(sessionMgmt BfdSessionMgmt) error {
	DestIp := sessionMgmt.DestIp
	Interface := sessionMgmt.Interface
	Protocol := sessionMgmt.Protocol
	ForceDel := sessionMgmt.ForceDel
	sessionIp := DestIp
	if Interface != "" {
		ipAddr := net.ParseIP(DestIp)
		if ipAddr.To4() == nil {
			if ipAddr.IsLinkLocalUnicast() {
				sessionIp = sessionIp + "%" + Interface
			}
		}
	}
	server.logger.Info("DeleteSession ", sessionIp, Protocol)
	sessionId, found := server.FindBfdSession(sessionIp)
	if found {
		session := server.bfdGlobal.Sessions[sessionId]
		if session.state.PerLinkSession {
			server.DeletePerLinkSessions(sessionIp, Protocol, ForceDel)
		} else {
			server.SessionDeleteHandler(session, Protocol, ForceDel)
		}
		server.ribdClient.ClientHdl.TrackReachabilityStatus(DestIp, "BFD", "del")
	} else {
		server.logger.Info("Bfd session not found ", sessionId)
	}
	return nil
}

func (server *BFDServer) ResetBfdSession(sessionId int32) error {
	server.logger.Info("ResetSession: SessionId", sessionId)
	session := server.bfdGlobal.Sessions[sessionId]
	if session.IsSessionActive() {
		session.SessionStopClientCh <- true
	}
	session.SessionStopServerCh <- true
	session.ResetLocalSessionParams()
	session.isClientActive = false
	server.logger.Info("Starting server for session ", sessionId)
	go session.StartSessionServer()
	server.logger.Info("Starting client for session ", sessionId)
	go session.StartSessionClient(server)
	return nil
}

func (server *BFDServer) AdminUpPerLinkBfdSessions(DestIp string) error {
	for _, session := range server.bfdGlobal.Sessions {
		if session.state.IpAddr == DestIp {
			session.StartBfdSession()
		}
	}
	return nil
}

// AdminUpBfdSession ceases the session.
func (server *BFDServer) AdminUpBfdSession(sessionMgmt BfdSessionMgmt) error {
	DestIp := sessionMgmt.DestIp
	Protocol := sessionMgmt.Protocol
	server.logger.Info("AdminDownSession ", DestIp, Protocol)
	sessionId, found := server.FindBfdSession(DestIp)
	if found {
		session := server.bfdGlobal.Sessions[sessionId]
		if session.state.PerLinkSession {
			server.AdminUpPerLinkBfdSessions(DestIp)
		} else {
			server.bfdGlobal.Sessions[sessionId].StartBfdSession()
		}
	} else {
		server.logger.Info("Bfd session not found ", sessionId)
	}
	return nil
}

func (server *BFDServer) AdminDownPerLinkBfdSessions(DestIp string) error {
	for _, session := range server.bfdGlobal.Sessions {
		if session.state.IpAddr == DestIp {
			session.StopBfdSession()
		}
	}
	return nil
}

// AdminDownBfdSession ceases the session.
func (server *BFDServer) AdminDownBfdSession(sessionMgmt BfdSessionMgmt) error {
	DestIp := sessionMgmt.DestIp
	Protocol := sessionMgmt.Protocol
	server.logger.Info("AdminDownSession ", DestIp, Protocol)
	sessionId, found := server.FindBfdSession(DestIp)
	if found {
		session := server.bfdGlobal.Sessions[sessionId]
		if session.state.PerLinkSession {
			server.AdminDownPerLinkBfdSessions(DestIp)
		} else {
			server.bfdGlobal.Sessions[sessionId].StopBfdSession()
		}
	} else {
		server.logger.Info("Bfd session not found ", sessionId)
	}
	return nil
}

// This function handles NextHop change from RIB.
// A Poll control packet will be sent to BFD neighbor and expect a Final control packet.
func (server *BFDServer) HandleNextHopChange(DestIp string, IfIndex int32, Reachable bool) error {
	server.logger.Info("HandleNextHopChange - ", DestIp, IfIndex, Reachable)
	if Reachable {
		// Go through the list of tobeCreatedSessions and try to recreate.
		for _, sessionMgmt := range server.tobeCreatedSessions {
			_, err := server.CreateBfdSession(sessionMgmt)
			if err == nil {
				delete(server.tobeCreatedSessions, sessionMgmt.DestIp)
			}
		}

		// TODO: Go through all the sessions that are InterfaceSpecific and match the ifIndex.
		// If reachability to DestIp is through a different interface then bring down the session.
	}
	return nil
}
