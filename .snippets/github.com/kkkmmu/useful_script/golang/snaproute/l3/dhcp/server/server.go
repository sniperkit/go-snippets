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
	"asicdServices"
	"encoding/json"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	//"github.com/google/gopacket/pcap"
	nanomsg "github.com/op/go-nanomsg"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"utils/ipcutils"
	"utils/logging"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type DhcpClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
}

type DhcpGlobalConfig struct {
	Enable           bool
	DefaultLeaseTime uint32
	MaxLeaseTime     uint32
}

type DhcpIntfConfig struct {
	Enable        bool
	IntfRef       string
	Subnet        uint32
	SubnetMask    uint32
	LowerIPBound  uint32
	HigherIPBound uint32
	BCastAddr     uint32
	RtrAddr       uint32
	DnsAddr       uint32
	DomainName    string
}

type DhcpIntfKey struct {
	subnet     uint32
	subnetMask uint32
}

type DhcpOfferedData struct {
	LeaseTime     uint32
	MacAddr       string
	TransactionId uint32
	RefreshTimer  *time.Timer
	StaleTimer    *time.Timer
	State         uint8
}

type DhcpIntfData struct {
	enable        bool
	l3IfIdx       int32
	lowerIPBound  uint32
	higherIPBound uint32
	rtrAddr       uint32
	dnsAddr       uint32
	domainName    string
	usedIpPool    map[uint32]DhcpOfferedData
	usedIpToMac   map[string]uint32
	dhcpMsg       []byte
}

type DHCPServer struct {
	logger              *logging.Writer
	DhcpGlobalConf      DhcpGlobalConfig
	DhcpGlobalConfCh    chan DhcpGlobalConfig
	DhcpIntfConfCh      chan DhcpIntfConfig
	DhcpIntfConfRetCh   chan error
	DhcpIntfConfMap     map[DhcpIntfKey]DhcpIntfData
	asicdSubSocket      *nanomsg.SubSocket
	asicdSubSocketCh    chan []byte
	asicdSubSocketErrCh chan error
	//l3PropertyMap       map[DhcpIntfKey]int32
	l3IntfPropMap   map[int32]L3Property
	portPropertyMap map[int32]PortProperty
	vlanPropertyMap map[int32]VlanProperty
	lagPropertyMap  map[int32]LagProperty
	asicdClient     AsicdClient
	InitDone        chan bool
	pcapTimeout     time.Duration
	promiscuous     bool
	snapshotLen     int32
}

func NewDHCPServer(logger *logging.Writer) *DHCPServer {
	dhcpServer := &DHCPServer{}
	dhcpServer.logger = logger
	dhcpServer.DhcpGlobalConfCh = make(chan DhcpGlobalConfig)
	dhcpServer.DhcpIntfConfCh = make(chan DhcpIntfConfig)
	dhcpServer.DhcpIntfConfRetCh = make(chan error)
	dhcpServer.DhcpIntfConfMap = make(map[DhcpIntfKey]DhcpIntfData)
	dhcpServer.asicdSubSocketCh = make(chan []byte)
	dhcpServer.asicdSubSocketErrCh = make(chan error)
	//dhcpServer.l3PropertyMap = make(map[DhcpIntfKey]int32)
	dhcpServer.l3IntfPropMap = make(map[int32]L3Property)
	dhcpServer.portPropertyMap = make(map[int32]PortProperty)
	dhcpServer.vlanPropertyMap = make(map[int32]VlanProperty)
	dhcpServer.lagPropertyMap = make(map[int32]LagProperty)
	dhcpServer.InitDone = make(chan bool)
	return dhcpServer
}

func (server *DHCPServer) initDhcpParams() {
	server.logger.Debug("Calling initParams...")
	server.snapshotLen = 65549
	server.promiscuous = false
	server.pcapTimeout = time.Duration(1) * time.Second
}

func (server *DHCPServer) connectToServers(paramsFile string) {
	server.logger.Debug(fmt.Sprintln("Inside connectToClients...paramsFile", paramsFile))
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		server.logger.Err("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		server.logger.Err("Error in Unmarshalling Json")
		return
	}

	for _, client := range clientsList {
		if client.Name == "asicd" {
			server.logger.Debug(fmt.Sprintln("found asicd at port", client.Port))
			server.asicdClient.Address = "localhost:" + strconv.Itoa(client.Port)
			server.asicdClient.Transport, server.asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.asicdClient.Address)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Failed to connect to Asicd, retrying until connection is successful"))
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					server.asicdClient.Transport, server.asicdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(server.asicdClient.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						server.logger.Err("Still can't connect to Asicd, retrying..")
					}
				}

			}
			server.logger.Info("Dhcpd is connected to Asicd")
			server.asicdClient.ClientHdl = asicdServices.NewASICDServicesClientFactory(server.asicdClient.Transport, server.asicdClient.PtrProtocolFactory)
		}
	}
}

func (server *DHCPServer) sigHandler(sigChan <-chan os.Signal) {
	server.logger.Debug("Inside sigHandler....")
	signal := <-sigChan
	switch signal {
	case syscall.SIGHUP:
		server.logger.Debug("Received SIGHUP signal")
		os.Exit(0)
	default:
		server.logger.Err(fmt.Sprintln("Unhandled signal : ", signal))
	}
}

func (server *DHCPServer) InitServer(paramDir string) {
	server.logger.Debug("Starting Dhcp Server")
	server.initDhcpParams()
	fileName := paramDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}
	fileName = fileName + "clients.json"
	server.connectToServers(fileName)
	server.buildDhcpInfra()
	server.logger.Debug("Listen for ASICd updates")
	server.listenForASICdUpdates(asicdCommonDefs.PUB_SOCKET_ADDR)
	go server.createASICdSubscriber()

	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)
	go server.sigHandler(sigChan)
	server.processDhcpInfra()
}

func (server *DHCPServer) StartServer(paramDir string) {
	server.logger.Debug(fmt.Sprintln("Inside Start Server...", paramDir))
	server.InitServer(paramDir)
	server.InitDone <- true
	for {
		select {
		case dhcpGlobalConf := <-server.DhcpGlobalConfCh:
			server.processDhcpGlobalConf(dhcpGlobalConf)
		case dhcpIntfConf := <-server.DhcpIntfConfCh:
			err, l3IfIdx := server.processDhcpIntfConf(dhcpIntfConf)
			server.DhcpIntfConfRetCh <- err
			if err == nil {
				server.handleDhcpIntfConf(l3IfIdx)
			}
		case asicdrxBuf := <-server.asicdSubSocketCh:
			server.processAsicdNotification(asicdrxBuf)
		case <-server.asicdSubSocketErrCh:
		}
	}
}
