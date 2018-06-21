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
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"log/syslog"
	"testing"
	"utils/asicdClient/mock"
	"utils/logging"
)

var bfdTestServer *BFDServer
var bfdTestSession *BfdSession
var bfdTestSessionParam *BfdSessionParam
var bfdTestControlPacket *BfdControlPacket

func BfdTestNewLogger() *logging.Writer {
	logger := new(logging.Writer)
	logger.SysLogger, _ = syslog.New(syslog.LOG_DEBUG|syslog.LOG_DAEMON, "BFDTEST")
	logger.MyLogLevel = sysdCommonDefs.DEBUG
	return logger
}

func initTestServer() {
	var paramFile string
	fmt.Println("Initializing BFD UT params")
	logger := BfdTestNewLogger()
	bfdTestServer = NewBFDServer(logger)
	bfdTestServer.InitServer(paramFile, &mockasicdclientplugin.MockAsicdClientMgr{})
	initSessionHandlingChans()
	return
}

func initSessionHandlingChans() {
	bfdTestServer.CreateSessionCh = make(chan BfdSessionMgmt)
	bfdTestServer.DeleteSessionCh = make(chan BfdSessionMgmt)
	bfdTestServer.AdminUpSessionCh = make(chan BfdSessionMgmt)
	bfdTestServer.AdminDownSessionCh = make(chan BfdSessionMgmt)
	bfdTestServer.CreatedSessionCh = make(chan int32, MAX_NUM_SESSIONS)
	bfdTestServer.FailedSessionClientCh = make(chan int32, MAX_NUM_SESSIONS)
	bfdTestServer.tobeCreatedSessions = make(map[string]BfdSessionMgmt)
}

func startTestServerChans() {
	for {
		select {
		case <-bfdTestServer.ServerStartedCh:
		case <-bfdTestServer.GlobalConfigCh:
		//case <-bfdTestServer.asicdSubSocketCh:
		//case <-bfdTestServer.asicdSubSocketErrCh:
		case <-bfdTestServer.ribdSubSocketCh:
		case <-bfdTestServer.ribdSubSocketErrCh:
		case <-bfdTestServer.CreateSessionCh:
		case <-bfdTestServer.DeleteSessionCh:
		case <-bfdTestServer.AdminUpSessionCh:
		case <-bfdTestServer.AdminDownSessionCh:
		case <-bfdTestServer.SessionConfigCh:
		case <-bfdTestServer.CreatedSessionCh:
		case <-bfdTestServer.notificationCh:
		case <-bfdTestServer.FailedSessionClientCh:
		case <-bfdTestServer.BfdPacketRecvCh:
		case <-bfdTestServer.SessionParamConfigCh:
		case <-bfdTestServer.SessionParamDeleteCh:
		}
	}
}

func TestCreateBfdServer(t *testing.T) {
	initTestServer()
	if bfdTestServer == nil {
		t.Fatal("Failed to initialize BFD server instance")
	} else {
		t.Log("Successfully initialize BFD server instance")
	}
	go startTestServerChans()
	t.Log("Started go routine to accept messages from all channels supported by server")
}

func TestBuildPortPropertyMap(t *testing.T) {
	err := bfdTestServer.BuildPortPropertyMap()
	if err == nil {
		t.Log("Successfully built port property map: ", len(bfdTestServer.portPropertyMap), " ports")
	} else {
		t.Fatal("Failed to build port property map")
	}
}

func TestCreateASICdSubscriber(t *testing.T) {
	//go bfdTestServer.CreateASICdSubscriber()
	t.Log("Created asicd subscriber go routine")
}

func TestCreateRIBdSubscriber(t *testing.T) {
	go bfdTestServer.CreateRIBdSubscriber()
	t.Log("Created ribd subscriber go routine")
}

func TestNewNormalBfdSession(t *testing.T) {
	bfdTestServer.createDefaultSessionParam()
	fmt.Println("Creating BFD session to 10.1.1.1")
	bfdTestSession = bfdTestServer.NewNormalBfdSession("", "", "10.1.1.1", "default", false, 2)
	if bfdTestSession != nil {
		t.Log("Created BFD session to ", bfdTestSession.state.IpAddr, " session id ", bfdTestSession.state.SessionId)
		if bfdTestSession.state.SessionState != STATE_DOWN {
			t.Fatal("Session created in ", bfdTestSession.state.SessionState, " state")
		}
	} else {
		t.Fatal("Failed to create session")
	}
}

func TestStartSessionServer(t *testing.T) {
	go bfdTestSession.StartSessionServer()
	t.Log("Stated session server for ", bfdTestSession.state.SessionId)
}

func TestStartSessionClient(t *testing.T) {
	go bfdTestSession.StartSessionClient(bfdTestServer)
	t.Log("Stated session client for ", bfdTestSession.state.SessionId)
}

func TestFindBfdSession(t *testing.T) {
	sessionId, found := bfdTestServer.FindBfdSession("10.1.1.1")
	if found {
		t.Log("Found session: ", sessionId)
	} else {
		t.Fatal("Failed to find session to 10.1.1.1")
	}
}

func TestEventHandler(t *testing.T) {
	t.Log("Session state before REMOTE_DOWN event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(REMOTE_DOWN)
	t.Log("Session state after REMOTE_DOWN event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(REMOTE_INIT)
	t.Log("Session state after REMOTE_INIT event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(TIMEOUT)
	t.Log("Session state after TIMEOUT event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(REMOTE_ADMIN_DOWN)
	t.Log("Session state after REMOTE_ADMIN_DOWN event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(ADMIN_UP)
	t.Log("Session state after ADMIN_UP event is ", bfdTestSession.state.SessionState)
	bfdTestSession.EventHandler(REMOTE_UP)
	t.Log("Session state after REMOTE_UP event is ", bfdTestSession.state.SessionState)
}

func TestUpdateBfdSessionControlPacket(t *testing.T) {
	bfdTestSession.UpdateBfdSessionControlPacket()
	t.Log("Updated control packet for session to ", bfdTestSession.state.IpAddr)
	t.Log("BFD control packet is - ", bfdTestSession.bfdPacket)
}

func TestCheckIfAnyProtocolRegistered(t *testing.T) {
	owner := bfdTestSession.CheckIfAnyProtocolRegistered()
	t.Log("Registered protocols for session to ", bfdTestSession.state.IpAddr, " is ", bfdTestSession.state.RegisteredProtocols)
	if owner != true {
		t.Fatal("Expecting USER as owner registered with session to ", bfdTestSession.state.IpAddr)
	}
}

func TestAdminDownBfdSession(t *testing.T) {
	sessionMgmt := BfdSessionMgmt{
		DestIp:   "10.1.1.1",
		Protocol: 2,
	}
	bfdTestServer.AdminDownBfdSession(sessionMgmt)
	t.Log("Session state changed to - ", bfdTestSession.state.SessionState)
}

/*
func TestAdminUpBfdSession(t *testing.T) {
	sessionMgmt := BfdSessionMgmt{
		DestIp:   "10.1.1.1",
		Protocol: 2,
	}
	bfdTestServer.AdminUpBfdSession(sessionMgmt)
	if bfdTestSession.state.SessionState != STATE_DOWN {
		t.Fatal("Failed to change session state to ADMIN_DOWN")
	}
}
*/

func TestSendBfdNotification(t *testing.T) {
	bfdTestSession.SendBfdNotification()
	t.Log("Sent BFD state notification")
}

/*
func TestSendPeriodicControlPackets(t *testing.T) {
	bfdTestSession.SendPeriodicControlPackets()
	t.Log("Sent BFD packet for session to ", bfdTestSession.state.IpAddr)
}

func TestHandleSessionTimeout(t *testing.T) {
	bfdTestSession.HandleSessionTimeout()
	if bfdTestSession.state.SessionState == STATE_DOWN {
		t.Log("Session ", bfdTestSession.state.SessionId, " went to down state due to timeout")
	} else {
		t.Fatal("Session timeout failed")
	}
}

func TestProcessBfdPacket(t *testing.T) {
	bfdTestSession.ProcessBfdPacket(bfdTestSession.bfdPacket)
	t.Log("Processing BFD packet - ", bfdTestSession.bfdPacket)
}
*/

func TestInitiatePollSequence(t *testing.T) {
	bfdTestSession.InitiatePollSequence()
	t.Log("Session ", bfdTestSession.state.SessionId, " initiated poll sequence")
}

func TestDecodeBfdControlPacket(t *testing.T) {
	bfdPacketBuf, _ := bfdTestSession.bfdPacket.CreateBfdControlPacket()
	DecodeBfdControlPacket(bfdPacketBuf)
}
