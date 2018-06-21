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
	"encoding/json"
	"github.com/google/gopacket/pcap"
	nanomsg "github.com/op/go-nanomsg"
	"io/ioutil"
	"l3/bfd/bfddCommonDefs"
	"net"
	"os"
	"os/signal"
	"ribd"
	"strconv"
	"sync"
	"syscall"
	"time"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dbutils"
	"utils/ipcutils"
	"utils/logging"
)

var (
	bfdSnapshotLen   int32  = 65549                   // packet capture length
	bfdPromiscuous   bool   = false                   // mode
	bfdDedicatedMac  string = "01:00:5E:90:00:01"     // Dest MAC perlink packets till neighbor's MAC is learned
	bfdPcapFilter    string = "udp and dst port 3784" // packet capture filter
	bfdPcapFilterLag string = "udp and dst port 6784" // packet capture filter
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type RibdClient struct {
	ipcutils.IPCClientBase
	ClientHdl *ribd.RIBDServicesClient
}

type BfdSessionMgmt struct {
	DestIp    string
	ParamName string
	Interface string
	Protocol  bfddCommonDefs.BfdSessionOwner
	PerLink   bool
	ForceDel  bool
}

type BfdSession struct {
	state               SessionState
	rxInterval          int32
	sessionTimer        *time.Timer
	txInterval          int32
	txTimer             *time.Timer
	TxTimeoutCh         chan int32
	txJitter            int32
	SessionTimeoutCh    chan int32
	bfdPacket           *BfdControlPacket
	bfdPacketBuf        []byte
	ReceivedPacketCh    chan *BfdControlPacket
	SessionStopClientCh chan bool
	SessionStopServerCh chan bool
	pollSequence        bool
	pollSequenceFinal   bool
	pollChanged         bool
	authEnabled         bool
	authType            AuthenticationType
	authSeqNum          uint32
	authKeyId           uint32
	authData            string
	txConn              net.Conn
	sendPcapHandle      *pcap.Handle
	recvPcapHandle      *pcap.Handle
	useDedicatedMac     bool
	paramChanged        bool
	remoteParamChanged  bool
	stateChanged        bool
	isClientActive      bool
	movedToDownState    bool
	notifiedState       bool
	server              *BFDServer
	sessionLock         sync.RWMutex
}

type BfdSessionParam struct {
	state SessionParamState
}

type BfdGlobal struct {
	Vrf                     string
	Enabled                 bool
	NumSessions             uint32
	Sessions                map[int32]*BfdSession
	SessionsIdSlice         []int32
	SessionsByIp            map[string]*BfdSession
	InactiveSessionsIdSlice []int32
	NumSessionParams        uint32
	SessionParams           map[string]*BfdSessionParam
	NumUpSessions           uint32
	NumDownSessions         uint32
	NumAdminDownSessions    uint32
}

type RecvedBfdPacket struct {
	IpAddr    string
	Len       int32
	PacketBuf []byte
}

type BFDServer struct {
	logger                *logging.Writer
	ServerStartedCh       chan bool
	ribdClient            RibdClient
	AsicdPlugin           asicdClntIntfs.AsicdClntIntf
	GlobalConfigCh        chan GlobalConfig
	AsicdSubSocketCh      chan clntIntfs.NotifyMsg
	ribdSubSocket         *nanomsg.SubSocket
	ribdSubSocketCh       chan []byte
	ribdSubSocketErrCh    chan error
	portPropertyMap       map[int32]PortProperty
	vlanPropertyMap       map[int32]VlanProperty
	CreateSessionCh       chan BfdSessionMgmt
	DeleteSessionCh       chan BfdSessionMgmt
	AdminUpSessionCh      chan BfdSessionMgmt
	AdminDownSessionCh    chan BfdSessionMgmt
	ResetSessionCh        chan int32
	SessionConfigCh       chan SessionConfig
	CreatedSessionCh      chan int32
	bfddPubSocket         *nanomsg.PubSocket
	lagPropertyMap        map[int32]LagProperty
	notificationCh        chan []byte
	FailedSessionClientCh chan int32
	BfdPacketRecvCh       chan RecvedBfdPacket
	SessionParamConfigCh  chan SessionParamConfig
	SessionParamDeleteCh  chan string
	tobeCreatedSessions   map[string]BfdSessionMgmt
	bfdGlobal             BfdGlobal
}

func NewBFDServer(logger *logging.Writer) *BFDServer {
	bfdServer := &BFDServer{}
	bfdServer.logger = logger
	bfdServer.ServerStartedCh = make(chan bool)
	bfdServer.GlobalConfigCh = make(chan GlobalConfig)
	bfdServer.AsicdSubSocketCh = make(chan clntIntfs.NotifyMsg)
	bfdServer.ribdSubSocketCh = make(chan []byte)
	bfdServer.ribdSubSocketErrCh = make(chan error)
	bfdServer.portPropertyMap = make(map[int32]PortProperty)
	bfdServer.vlanPropertyMap = make(map[int32]VlanProperty)
	bfdServer.lagPropertyMap = make(map[int32]LagProperty)
	bfdServer.SessionConfigCh = make(chan SessionConfig)
	bfdServer.notificationCh = make(chan []byte)
	bfdServer.SessionParamConfigCh = make(chan SessionParamConfig)
	bfdServer.SessionParamDeleteCh = make(chan string)
	bfdServer.bfdGlobal.Enabled = false
	bfdServer.bfdGlobal.NumSessions = 0
	bfdServer.bfdGlobal.Sessions = make(map[int32]*BfdSession)
	bfdServer.bfdGlobal.SessionsIdSlice = []int32{}
	bfdServer.bfdGlobal.SessionsByIp = make(map[string]*BfdSession)
	bfdServer.bfdGlobal.InactiveSessionsIdSlice = []int32{}
	bfdServer.bfdGlobal.NumSessionParams = 0
	bfdServer.bfdGlobal.SessionParams = make(map[string]*BfdSessionParam)
	bfdServer.bfdGlobal.NumUpSessions = 0
	bfdServer.bfdGlobal.NumDownSessions = 0
	bfdServer.bfdGlobal.NumAdminDownSessions = 0
	return bfdServer
}

func (server *BFDServer) SigHandler(dbHdl *dbutils.DBUtil) {
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP:
				server.SendDeleteToAllSessions()
				time.Sleep(500 * time.Millisecond)
				server.logger.Info("Stopped all sessions")
				dbHdl.Disconnect()
				server.logger.Info("Exting!!!")
				os.Exit(0)
			default:
			}
		}
	}
}

func (server *BFDServer) ConnectToServers(paramsFile string) {
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		server.logger.Info("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		server.logger.Info("Error in Unmarshalling Json")
		return
	}

	for _, client := range clientsList {
		if client.Name == "ribd" {
			server.logger.Info("found ribd at port", client.Port)
			server.ribdClient.Address = "localhost:" + strconv.Itoa(client.Port)
			server.ribdClient.TTransport, server.ribdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.ribdClient.Address)
			if err != nil {
				server.logger.Info("Failed to connect to Ribd, retrying until connection is successful")
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					server.ribdClient.TTransport, server.ribdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.ribdClient.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						server.logger.Info("Still can't connect to Ribd, retrying...")
					}
				}
			}
			if server.ribdClient.TTransport != nil && server.ribdClient.PtrProtocolFactory != nil {
				server.ribdClient.ClientHdl = ribd.NewRIBDServicesClientFactory(server.ribdClient.TTransport, server.ribdClient.PtrProtocolFactory)
				server.ribdClient.IsConnected = true
				server.logger.Info("Bfdd is connected to Ribd")
			}
		}
	}
}

func (server *BFDServer) InitPublisher(pub_str string) (pub *nanomsg.PubSocket) {
	server.logger.Info("Setting up ", pub_str, "publisher")
	pub, err := nanomsg.NewPubSocket()
	if err != nil {
		server.logger.Info("Failed to open pub socket")
		return nil
	}
	ep, err := pub.Bind(pub_str)
	if err != nil {
		server.logger.Info("Failed to bind pub socket - ", ep)
		return nil
	}
	err = pub.SetSendBuffer(1024)
	if err != nil {
		server.logger.Info("Failed to set send buffer size")
		return nil
	}
	return pub
}

func (server *BFDServer) PublishSessionNotifications() {
	server.bfddPubSocket = server.InitPublisher(bfddCommonDefs.PUB_SOCKET_ADDR)
	for {
		select {
		case event := <-server.notificationCh:
			_, err := server.bfddPubSocket.Send(event, nanomsg.DontWait)
			if err == syscall.EAGAIN {
				server.logger.Err("Failed to publish event")
			}
		}
	}
}

func (server *BFDServer) InitServer(paramFile string) {
	server.logger.Info("Starting Bfd Server")
	server.ConnectToServers(paramFile)
	server.initBfdGlobalConfDefault()
	server.BuildPortPropertyMap()
	server.BuildLagPropertyMap()
	server.createDefaultSessionParam()
}

func (server *BFDServer) StartServer(paramFile string, dbHdl *dbutils.DBUtil) {
	// Initialize BFD server from params file
	server.InitServer(paramFile)
	// Start subcriber for RIBd events
	go server.CreateRIBdSubscriber()
	// Start session management handler
	go server.StartSessionHandler()
	// Initialize and run notification publisher
	go server.PublishSessionNotifications()

	server.ServerStartedCh <- true

	// Now, wait on below channels to process
	for {
		select {
		case gConf := <-server.GlobalConfigCh:
			server.logger.Info("Received call for performing Global Configuration", gConf)
			server.processGlobalConfig(gConf)
		case asicdrxBuf := <-server.AsicdSubSocketCh:
			server.processAsicdNotification(asicdrxBuf)
		case ribdrxBuf := <-server.ribdSubSocketCh:
			server.processRibdNotification(ribdrxBuf)
		case <-server.ribdSubSocketErrCh:
		case sessionConfig := <-server.SessionConfigCh:
			server.logger.Info("Received call for performing Session Configuration", sessionConfig)
			server.processSessionConfig(sessionConfig)
		case sessionParamConfig := <-server.SessionParamConfigCh:
			server.logger.Info("Received call for performing Session Param Configuration", sessionParamConfig)
			server.processSessionParamConfig(sessionParamConfig)
		case paramName := <-server.SessionParamDeleteCh:
			server.logger.Info("Received call for performing Session Param Delete", paramName)
			server.processSessionParamDelete(paramName)
		}
	}
}
