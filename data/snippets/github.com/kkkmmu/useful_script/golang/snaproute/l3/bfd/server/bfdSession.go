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
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"l3/bfd/bfddCommonDefs"
	"net"
	"strconv"
	"time"
)

func (session *BfdSession) StartSessionServer() error {
	session.server.logger.Info("Started session server for ", session.state.SessionId)
	for {
		select {
		case bfdPacket := <-session.ReceivedPacketCh:
			session.state.NumRxPackets++
			session.ProcessBfdPacket(bfdPacket)
		case <-session.SessionStopServerCh:
			session.server.logger.Info("Exiting session server ", session.state.SessionId)
			return nil
		}
	}

	return nil
}

func (session *BfdSession) StartSessionClient(server *BFDServer) error {
	var err error
	server.logger.Info("Starting session client for ", session.state.SessionId)
	destAddr := net.JoinHostPort(session.state.IpAddr, strconv.Itoa(DEST_PORT))
	ServerAddr, err := net.ResolveUDPAddr("udp", destAddr)
	if err != nil {
		server.logger.Info("Failed ResolveUDPAddr ", destAddr, err)
		server.FailedSessionClientCh <- session.state.SessionId
		return err
	}
	localAddr := net.JoinHostPort(session.state.LocalAddr, strconv.Itoa(int(SRC_PORT+session.state.SessionId)))
	ClientAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		server.logger.Info("Failed ResolveUDPAddr ", localAddr, err)
		server.FailedSessionClientCh <- session.state.SessionId
		return err
	}
	Conn, err := net.DialUDP("udp", ClientAddr, ServerAddr)
	if err != nil {
		server.logger.Info("Failed DialUDP ", ClientAddr, ServerAddr, err)
		server.FailedSessionClientCh <- session.state.SessionId
		return err
	}
	session.sessionLock.Lock()
	session.txConn = Conn
	server.logger.Info("Started session client for ", destAddr, localAddr)
	defer session.txConn.Close()
	session.txTimer = time.AfterFunc(time.Duration(session.txInterval)*time.Millisecond, func() { session.SendPeriodicControlPackets() })
	session.sessionTimer = time.AfterFunc(time.Duration(session.rxInterval)*time.Millisecond, func() { session.HandleSessionTimeout() })
	defer session.txTimer.Stop()
	defer session.sessionTimer.Stop()
	session.isClientActive = true
	session.sessionLock.Unlock()
	for {
		select {
		case <-session.SessionStopClientCh:
			server.logger.Info("Exiting session client ", session.state.SessionId)
			return nil
		}
	}
}

/* State Machine
                             +--+
                             |  | UP, TIMER
                             |  V
                     DOWN  +------+  INIT
              +------------|      |------------+
              |            | DOWN |            |
              |  +-------->|      |<--------+  |
              |  |         +------+         |  |
              |  |                          |  |
              |  |                          |  |
              |  |                     DOWN,|  |
              |  |TIMER                TIMER|  |
              V  |                          |  V
            +------+                      +------+
       +----|      |                      |      |----+
   DOWN|    | INIT |--------------------->|  UP  |    |INIT, UP
       +--->|      | INIT, UP             |      |<---+
            +------+                      +------+
*/
// EventHandler is called after receiving a BFD packet from remote.
func (session *BfdSession) EventHandler(event BfdSessionEvent) error {
	var err error
	if session.IsSessionActive() == false {
		session.server.logger.Info("Cannot process event ", event, " Session ", session.state.SessionId, " not active")
		err = errors.New("Session is not active. No event can be processed.")
		return err
	}
	switch session.state.SessionState {
	case STATE_ADMIN_DOWN:
		session.server.logger.Info("Received ", event, " event for an admindown session")
	case STATE_DOWN:
		switch event {
		case REMOTE_DOWN:
			session.MoveToInitState()
		case REMOTE_INIT:
			session.MoveToUpState()
		case ADMIN_UP:
			session.MoveToDownState()
		case ADMIN_DOWN:
			session.LocalAdminDown()
		case REMOTE_ADMIN_DOWN:
			session.RemoteAdminDown()
		case TIMEOUT, REMOTE_UP:
		}
	case STATE_INIT:
		switch event {
		case REMOTE_INIT, REMOTE_UP:
			session.MoveToUpState()
		case TIMEOUT:
			session.MoveToDownState()
		case ADMIN_DOWN:
			session.LocalAdminDown()
		case REMOTE_ADMIN_DOWN:
			session.RemoteAdminDown()
		case REMOTE_DOWN, ADMIN_UP:
		}
	case STATE_UP:
		switch event {
		case REMOTE_DOWN, TIMEOUT:
			session.MoveToDownState()
		case ADMIN_DOWN:
			session.LocalAdminDown()
		case REMOTE_ADMIN_DOWN:
			session.RemoteAdminDown()
		case REMOTE_INIT, REMOTE_UP, ADMIN_UP:
		}
	}
	return err
}

func (session *BfdSession) CanProcessBfdControlPacket(bfdPacket *BfdControlPacket) bool {
	var canProcess bool
	canProcess = true
	if bfdPacket.Version != DEFAULT_BFD_VERSION {
		canProcess = false
		session.server.logger.Info("Can't process version mismatch ", bfdPacket.Version, DEFAULT_BFD_VERSION)
	}
	if bfdPacket.DetectMult == 0 {
		canProcess = false
		session.server.logger.Info("Can't process detect multi ", bfdPacket.DetectMult)
	}
	if bfdPacket.Multipoint {
		canProcess = false
		session.server.logger.Info("Can't process Multipoint ", bfdPacket.Multipoint)
	}
	if bfdPacket.MyDiscriminator == 0 {
		canProcess = false
		session.server.logger.Info("Can't process remote discriminator ", bfdPacket.MyDiscriminator)
	}
	/*
		if bfdPacket.YourDiscriminator == 0 {
			if session.state.SessionState == STATE_UP {
				canProcess = false
				session.server.logger.Info("Can't process packet with my discriminator ", bfdPacket.YourDiscriminator, " in up state")
			}
		}
	*/
	/*
		if bfdPacket.YourDiscriminator == 0 {
			canProcess = false
			session.server.logger.Info("Can't process local discriminator ", bfdPacket.YourDiscriminator)
		} else {
			sessionId := bfdPacket.YourDiscriminator
			session := server.bfdGlobal.Sessions[int32(sessionId)]
			if session != nil {
				if session.state.SessionState == STATE_ADMIN_DOWN {
					canProcess = false
				}
			}
		}
	*/
	return canProcess
}

func (session *BfdSession) AuthenticateReceivedControlPacket(bfdPacket *BfdControlPacket) bool {
	var authenticated bool
	if !bfdPacket.AuthPresent {
		authenticated = true
	} else {
		copiedPacket := &BfdControlPacket{}
		*copiedPacket = *bfdPacket
		authType := bfdPacket.AuthHeader.Type
		keyId := uint32(bfdPacket.AuthHeader.AuthKeyID)
		authData := bfdPacket.AuthHeader.AuthData
		seqNum := bfdPacket.AuthHeader.SequenceNumber
		if authType == session.authType {
			if authType == BFD_AUTH_TYPE_SIMPLE {
				session.server.logger.Info("Authentication type simple: keyId, authData ", keyId, string(authData))
				if keyId == session.authKeyId && string(authData) == session.authData {
					authenticated = true
				}
			} else {
				if seqNum >= session.state.ReceivedAuthSeq && keyId == session.authKeyId {
					var binBuf bytes.Buffer
					copiedPacket.AuthHeader.AuthData = []byte(session.authData)
					binary.Write(&binBuf, binary.BigEndian, copiedPacket)
					switch authType {
					case BFD_AUTH_TYPE_KEYED_MD5, BFD_AUTH_TYPE_METICULOUS_MD5:
						var authDataSum [16]byte
						authDataSum = md5.Sum(binBuf.Bytes())
						if bytes.Equal(authData[:], authDataSum[:]) {
							authenticated = true
						} else {
							session.server.logger.Info("Authentication data did't match for type: ", authType)
						}
					case BFD_AUTH_TYPE_KEYED_SHA1, BFD_AUTH_TYPE_METICULOUS_SHA1:
						var authDataSum [20]byte
						authDataSum = sha1.Sum(binBuf.Bytes())
						if bytes.Equal(authData[:], authDataSum[:]) {
							authenticated = true
						} else {
							session.server.logger.Info("Authentication data did't match for type: ", authType)
						}
					}
				} else {
					session.server.logger.Info("Sequence number and key id check failed: ", seqNum, session.state.ReceivedAuthSeq, keyId, session.authKeyId)
				}
			}
		} else {
			session.server.logger.Info("Authentication type did't match: ", authType, session.authType)
		}
	}
	return authenticated
}

func (session *BfdSession) ProcessBfdPacket(bfdPacket *BfdControlPacket) error {
	var event BfdSessionEvent
	authenticated := session.AuthenticateReceivedControlPacket(bfdPacket)
	if authenticated == false {
		session.server.logger.Info("Can't authenticatereceived bfd packet for session ", session.state.SessionId)
		return nil
	}
	canProcess := session.CanProcessBfdControlPacket(bfdPacket)
	if canProcess == false {
		session.server.logger.Info("Can't process received bfd packet for session ", session.state.SessionId)
		return nil
	}
	if session.state.SessionState == STATE_UP && session.state.RemoteSessionState == STATE_UP {
		session.rxInterval = (int32(bfdPacket.DesiredMinTxInterval) * int32(bfdPacket.DetectMult)) / 1000
	} else {
		session.rxInterval = (STARTUP_RX_INTERVAL * int32(bfdPacket.DetectMult)) / 1000
	}
	session.CheckAnyRemoteParamChanged(bfdPacket)
	session.RemoteChangedDemandMode(bfdPacket)
	session.ProcessPollSequence(bfdPacket)
	if session.rxInterval == 0 ||
		session.state.SessionState == STATE_ADMIN_DOWN ||
		session.state.RemoteSessionState == STATE_ADMIN_DOWN {
		session.sessionTimer.Stop()
	}
	switch session.state.RemoteSessionState {
	case STATE_DOWN:
		event = REMOTE_DOWN
		if session.state.SessionState != STATE_DOWN {
			session.state.LocalDiagType = DIAG_NEIGHBOR_SIGNAL_DOWN
		}
	case STATE_INIT:
		event = REMOTE_INIT
	case STATE_UP:
		event = REMOTE_UP
		if session.state.SessionState == STATE_UP {
			if session.txInterval != session.state.DesiredMinTxInterval/1000 {
				session.txInterval = session.state.DesiredMinTxInterval / 1000
				session.txTimer.Reset(0)
			}
		}
	case STATE_ADMIN_DOWN:
		event = REMOTE_ADMIN_DOWN
	}
	session.EventHandler(event)
	return nil
}

func (session *BfdSession) CheckAnyRemoteParamChanged(bfdPacket *BfdControlPacket) error {
	if session.state.RemoteSessionState != bfdPacket.State ||
		session.state.RemoteDiscriminator != bfdPacket.MyDiscriminator ||
		session.state.RemoteMinRxInterval != int32(bfdPacket.RequiredMinRxInterval) {
		session.remoteParamChanged = true
	}
	session.state.RemoteSessionState = bfdPacket.State
	session.state.RemoteDiscriminator = bfdPacket.MyDiscriminator
	session.state.RemoteMinRxInterval = int32(bfdPacket.RequiredMinRxInterval)
	session.state.RemoteDetectionMultiplier = int32(bfdPacket.DetectMult)
	return nil
}

func (session *BfdSession) UpdateBfdSessionControlPacket() error {
	session.bfdPacket.Diagnostic = session.state.LocalDiagType
	session.bfdPacket.State = session.state.SessionState
	session.bfdPacket.DetectMult = uint8(session.state.DetectionMultiplier)
	session.bfdPacket.MyDiscriminator = session.state.LocalDiscriminator
	session.bfdPacket.YourDiscriminator = session.state.RemoteDiscriminator
	if session.state.SessionState == STATE_UP && session.state.RemoteSessionState == STATE_UP {
		if session.bfdPacket.DesiredMinTxInterval == time.Duration(STARTUP_TX_INTERVAL) ||
			session.bfdPacket.RequiredMinRxInterval == time.Duration(STARTUP_RX_INTERVAL) {
			session.bfdPacket.DesiredMinTxInterval = time.Duration(session.state.DesiredMinTxInterval)
			session.bfdPacket.RequiredMinRxInterval = time.Duration(session.state.RequiredMinRxInterval)
			session.InitiatePollSequence()
		}

		wasDemand := session.bfdPacket.Demand
		session.bfdPacket.Demand = session.state.DemandMode
		isDemand := session.bfdPacket.Demand
		if !wasDemand && isDemand {
			session.server.logger.Info("Enabled demand for session ", session.state.SessionId)
			session.sessionTimer.Stop()
		}
		if wasDemand && !isDemand {
			session.server.logger.Info("Disabled demand for session ", session.state.SessionId)
			session.sessionTimer.Reset(time.Duration(session.rxInterval) * time.Millisecond)
		}
	} else {
		session.bfdPacket.DesiredMinTxInterval = time.Duration(STARTUP_TX_INTERVAL)
		session.bfdPacket.RequiredMinRxInterval = time.Duration(STARTUP_RX_INTERVAL)
	}
	session.bfdPacket.Poll = session.pollSequence
	session.pollSequence = false
	session.bfdPacket.Final = session.pollSequenceFinal
	session.pollSequenceFinal = false
	if session.authEnabled {
		session.bfdPacket.AuthPresent = true
		session.bfdPacket.AuthHeader.Type = session.authType
		if session.authType != BFD_AUTH_TYPE_SIMPLE {
			session.bfdPacket.AuthHeader.SequenceNumber = session.authSeqNum
		}
		if session.authType == BFD_AUTH_TYPE_METICULOUS_MD5 || session.authType == BFD_AUTH_TYPE_METICULOUS_SHA1 {
			session.authSeqNum++
		}
		session.bfdPacket.AuthHeader.AuthKeyID = uint8(session.authKeyId)
		session.bfdPacket.AuthHeader.AuthData = []byte(session.authData)
	} else {
		session.bfdPacket.AuthPresent = false
	}
	session.pollChanged = false
	session.paramChanged = false
	session.stateChanged = false
	session.remoteParamChanged = false
	return nil
}

func (session *BfdSession) CheckIfAnyProtocolRegistered() bool {
	for i := bfddCommonDefs.BfdSessionOwner(1); i < bfddCommonDefs.MAX_APPS; i++ {
		if session.state.RegisteredProtocols[i] == true {
			return true
		}
	}
	return false
}

// Stop session as Bfd is disabled globally. Do not delete
func (session *BfdSession) StopBfdSession() error {
	session.EventHandler(ADMIN_DOWN)
	session.state.LocalDiagType = DIAG_ADMIN_DOWN
	return nil
}

func (session *BfdSession) GetBfdSessionNotification() bool {
	var bfdState bool
	bfdState = false
	if session.state.SessionState == STATE_UP ||
		session.state.SessionState == STATE_ADMIN_DOWN ||
		session.state.RemoteSessionState == STATE_ADMIN_DOWN {
		bfdState = true
	}
	return bfdState
}

func (session *BfdSession) SendBfdNotification() error {
	bfdState := session.GetBfdSessionNotification()
	if bfdState != session.notifiedState {
		session.notifiedState = bfdState
		bfdNotification := bfddCommonDefs.BfddNotifyMsg{
			DestIp: session.state.IpAddr,
			State:  bfdState,
		}
		bfdNotificationBuf, err := json.Marshal(bfdNotification)
		if err != nil {
			session.server.logger.Err("Failed to marshal BfdSessionNotification message for session ", session.state.SessionId)
		}
		session.server.notificationCh <- bfdNotificationBuf
	}
	return nil
}

// Restart session that was stopped earlier due to global Bfd disable.
func (session *BfdSession) StartBfdSession() error {
	session.sessionTimer.Reset(time.Duration(session.rxInterval) * time.Millisecond)
	txInterval := session.ApplyTxJitter()
	session.txTimer.Reset(time.Duration(txInterval) * time.Millisecond)
	session.state.SessionState = STATE_DOWN
	session.EventHandler(ADMIN_UP)
	return nil
}

func (session *BfdSession) IsSessionActive() bool {
	if session.isClientActive {
		return true
	} else {
		return false
	}
}

func (session *BfdSession) ResetLocalSessionParams() error {
	session.state.SessionState = STATE_DOWN
	session.state.NumRxPackets = 0
	session.state.NumTxPackets = 0
	session.state.ToUpCount = 0
	session.state.ToDownCount = 0
	return nil
}

func (session *BfdSession) ResetRemoteSessionParams() error {
	session.state.RemoteDiscriminator = 0
	session.state.RemoteSessionState = STATE_DOWN
	session.remoteParamChanged = true
	return nil
}

func (session *BfdSession) LocalAdminDown() error {
	session.state.SessionState = STATE_ADMIN_DOWN
	session.state.LocalDiagType = DIAG_ADMIN_DOWN
	session.stateChanged = true
	session.SendBfdNotification()
	session.txInterval = STARTUP_TX_INTERVAL / 1000
	session.txTimer.Reset(0)
	session.rxInterval = (STARTUP_RX_INTERVAL * session.state.DetectionMultiplier) / 1000
	session.sessionTimer.Stop()
	return nil
}

func (session *BfdSession) RemoteAdminDown() error {
	session.state.RemoteSessionState = STATE_ADMIN_DOWN
	session.state.LocalDiagType = DIAG_NEIGHBOR_SIGNAL_DOWN
	session.SendBfdNotification()
	session.txInterval = STARTUP_TX_INTERVAL / 1000
	session.txTimer.Reset(0)
	session.rxInterval = (STARTUP_RX_INTERVAL * session.state.DetectionMultiplier) / 1000
	session.sessionTimer.Stop()
	return nil
}

func (session *BfdSession) MoveToDownState() error {
	session.state.SessionState = STATE_DOWN
	session.state.ToDownCount++
	session.server.logger.Info("Session ", session.state.SessionId, " moved to down state at ", time.Now().String())
	session.movedToDownState = true
	session.useDedicatedMac = true
	session.stateChanged = true
	if session.authType == BFD_AUTH_TYPE_KEYED_MD5 || session.authType == BFD_AUTH_TYPE_KEYED_SHA1 {
		session.authSeqNum++
	}
	session.SendBfdNotification()
	session.txInterval = STARTUP_TX_INTERVAL / 1000
	session.txTimer.Reset(time.Duration(session.txInterval) * time.Millisecond)
	session.rxInterval = (STARTUP_RX_INTERVAL * session.state.DetectionMultiplier) / 1000
	session.sessionTimer.Reset(time.Duration(session.rxInterval) * time.Millisecond)
	return nil
}

func (session *BfdSession) MoveToInitState() error {
	session.state.SessionState = STATE_INIT
	session.stateChanged = true
	session.useDedicatedMac = true
	session.movedToDownState = false
	return nil
}

func (session *BfdSession) MoveToUpState() error {
	session.state.SessionState = STATE_UP
	session.state.UpTime = time.Now()
	session.state.ToUpCount++
	session.stateChanged = true
	session.movedToDownState = false
	session.state.LocalDiagType = DIAG_NONE
	session.SendBfdNotification()
	return nil
}

func (session *BfdSession) ApplyTxJitter() int32 {
	return (int32(float32(session.txInterval) * (1 - float32(session.txJitter)/100)))
}

func (session *BfdSession) NeedBfdPacketUpdate() bool {
	if session.paramChanged || session.remoteParamChanged ||
		session.pollChanged || session.stateChanged {
		return true
	}
	return false
}

func (session *BfdSession) SendPeriodicControlPackets() {
	var err error
	var packetUpdated bool
	if session.NeedBfdPacketUpdate() {
		packetUpdated = true
		session.UpdateBfdSessionControlPacket()
		session.bfdPacketBuf, err = session.bfdPacket.CreateBfdControlPacket()
		if err != nil {
			session.server.logger.Info("Failed to create control packet for session ", session.state.SessionId)
		}
	}
	_, err = session.txConn.Write(session.bfdPacketBuf)
	if err != nil {
		session.server.logger.Info("failed to send control packet for session ", session.state.SessionId)
	} else {
		session.state.NumTxPackets++
	}
	if packetUpdated {
		// Re-compute the packet to clear any flag set in the previously sent packet
		session.UpdateBfdSessionControlPacket()
		session.bfdPacketBuf, err = session.bfdPacket.CreateBfdControlPacket()
		if err != nil {
			session.server.logger.Info("Failed to create control packet for session ", session.state.SessionId)
		}
		packetUpdated = false
	}
	if session.state.SessionState == STATE_UP || session.state.RemoteSessionState == STATE_UP {
		session.txInterval = session.state.DesiredMinTxInterval / 1000
	}
	txTimer := session.ApplyTxJitter()
	session.txTimer.Reset(time.Duration(txTimer) * time.Millisecond)
}

func (session *BfdSession) HandleSessionTimeout() {
	if session.state.SessionState != STATE_DOWN &&
		session.state.SessionState != STATE_ADMIN_DOWN {
		session.server.logger.Info("Timer expired for: ", session.state.IpAddr, " session id ", session.state.SessionId, " prev state ", session.server.ConvertBfdSessionStateValToStr(session.state.SessionState), " at ", time.Now().String())
	}
	if session.movedToDownState {
		session.server.logger.Info("Resetting remote params of session ", session.state.SessionId, " at ", time.Now().String())
		session.ResetRemoteSessionParams()
		session.movedToDownState = false
	}
	session.state.LocalDiagType = DIAG_TIME_EXPIRED
	session.EventHandler(TIMEOUT)
	session.sessionTimer.Reset(time.Duration(session.rxInterval) * time.Millisecond)
}

func (session *BfdSession) RemoteChangedDemandMode(bfdPacket *BfdControlPacket) error {
	var wasDemandMode, isDemandMode bool
	wasDemandMode = session.state.RemoteDemandMode
	session.state.RemoteDemandMode = bfdPacket.Demand
	if session.state.RemoteDemandMode {
		isDemandMode = true
		session.txTimer.Stop()
	}
	if wasDemandMode && !isDemandMode {
		txInterval := session.ApplyTxJitter()
		session.txTimer.Reset(time.Duration(txInterval) * time.Millisecond)
	}
	return nil
}

func (session *BfdSession) InitiatePollSequence() error {
	if !session.pollSequence {
		session.server.logger.Info("Starting poll sequence for session ", session.state.SessionId)
		session.pollSequence = true
		session.pollChanged = true
	}
	return nil
}

func (session *BfdSession) ProcessPollSequence(bfdPacket *BfdControlPacket) error {
	if session.state.SessionState != STATE_ADMIN_DOWN {
		if bfdPacket.Poll {
			session.server.logger.Info("Received packet with poll bit for session ", session.state.SessionId)
			session.pollSequenceFinal = true
			session.pollChanged = true
		}
		if bfdPacket.Final {
			session.server.logger.Info("Received packet with final bit for session ", session.state.SessionId)
			session.pollSequence = false
			session.pollChanged = true
		}
	}
	return nil
}
