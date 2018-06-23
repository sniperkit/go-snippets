//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
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
	"os"
	"os/signal"
	"syscall"
	"time"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type ArpEntry struct {
	MacAddr   string
	VlanId    int
	IfName    string
	L3IfIdx   int
	Counter   int
	TimeStamp time.Time
	PortNum   int
	Type      bool //True : RIB False: RX
}

type ArpState struct {
	IpAddr         string
	MacAddr        string
	VlanId         int
	Intf           string
	ExpiryTimeLeft string
}

type ResolveIPv4 struct {
	TargetIP string
	IfId     int
}

type ArpConf struct {
	RefTimeout int
}

type GarpEntry struct {
	IfName  string
	MacAddr string
	IpAddr  string
}

type ARPServer struct {
	IsLinuxOnly             bool
	logger                  logging.LoggerIntf
	arpCache                map[string]ArpEntry //Key: Dest IpAddr
	AsicdSubSocketCh        chan clntIntfs.NotifyMsg
	dbHdl                   *dbutils.DBUtil
	eventDbHdl              *dbutils.DBUtil
	snapshotLen             int32
	pcapTimeout             time.Duration
	promiscuous             bool
	ConfRefreshTimeout      int
	MinRefreshTimeout       int
	timerGranularity        int
	timeout                 time.Duration
	timeoutCounter          int
	minCnt                  int
	retryCnt                int
	probeWait               int
	probeNum                int
	probeMax                int
	probeMin                int
	arpSliceRefreshTimer    *time.Timer
	arpSliceRefreshDuration time.Duration
	usrConfDbName           string
	l3IntfPropMap           map[int]L3IntfProperty        //Key: IfIndex
	portPropMap             map[int]PortProperty          //Key: IfIndex
	vlanPropMap             map[int]VlanProperty          //Key: IfIndex
	lagPropMap              map[int]LagProperty           //Key: IfIndex
	virtualIntfPropMap      map[int32]VirtualIntfProperty //key: IfIndex
	arpSlice                []string
	arpEntryUpdateCh        chan UpdateArpEntryMsg
	arpEntryDeleteCh        chan DeleteArpEntryMsg
	arpEntryMacMoveCh       chan asicdClntDefs.IPv4NbrMacMoveNotifyMsg
	arpEntryCntUpdateCh     chan int
	arpSliceRefreshStartCh  chan bool
	arpSliceRefreshDoneCh   chan bool
	arpCounterUpdateCh      chan bool
	arpActionProcessCh      chan ArpActionMsg
	ResolveIPv4Ch           chan ResolveIPv4
	DeleteResolvedIPv4Ch    chan DeleteResolvedIPv4
	DeleteArpEntryCh        chan *DeleteArpEntry
	GarpEntryCh             chan *GarpEntry
	ArpConfCh               chan ArpConf
	dumpArpTable            bool
	InitDone                chan bool

	ArpActionCh                chan ArpActionMsg
	arpDeleteArpEntryFromRibCh chan string
	arpDeleteArpEntryIntCh     chan string
	SysRsvdVlan                int

	AsicdPlugin asicdClntIntfs.AsicdClntIntf
}

func NewARPServer(logger logging.LoggerIntf) *ARPServer {
	arpServer := &ARPServer{}
	arpServer.logger = logger
	arpServer.arpCache = make(map[string]ArpEntry)
	arpServer.AsicdSubSocketCh = make(chan clntIntfs.NotifyMsg)
	arpServer.l3IntfPropMap = make(map[int]L3IntfProperty)
	arpServer.lagPropMap = make(map[int]LagProperty)
	arpServer.vlanPropMap = make(map[int]VlanProperty)
	arpServer.portPropMap = make(map[int]PortProperty)
	arpServer.virtualIntfPropMap = make(map[int32]VirtualIntfProperty)
	arpServer.arpSlice = make([]string, 0)
	arpServer.arpEntryUpdateCh = make(chan UpdateArpEntryMsg)
	arpServer.arpEntryDeleteCh = make(chan DeleteArpEntryMsg)
	arpServer.arpEntryCntUpdateCh = make(chan int)
	arpServer.arpSliceRefreshStartCh = make(chan bool)
	arpServer.arpSliceRefreshDoneCh = make(chan bool)
	arpServer.arpCounterUpdateCh = make(chan bool)
	arpServer.arpActionProcessCh = make(chan ArpActionMsg)
	arpServer.arpDeleteArpEntryFromRibCh = make(chan string)
	arpServer.arpDeleteArpEntryIntCh = make(chan string)
	arpServer.ResolveIPv4Ch = make(chan ResolveIPv4)
	arpServer.DeleteResolvedIPv4Ch = make(chan DeleteResolvedIPv4)
	arpServer.DeleteArpEntryCh = make(chan *DeleteArpEntry)
	arpServer.GarpEntryCh = make(chan *GarpEntry)
	arpServer.ArpConfCh = make(chan ArpConf)
	arpServer.InitDone = make(chan bool)
	arpServer.ArpActionCh = make(chan ArpActionMsg)
	arpServer.arpEntryMacMoveCh = make(chan asicdClntDefs.IPv4NbrMacMoveNotifyMsg)
	return arpServer
}

func (server *ARPServer) initArpParams() {
	server.logger.Debug("Calling initParams...")
	server.snapshotLen = 65549
	server.promiscuous = false
	server.pcapTimeout = time.Duration(1) * time.Second
	server.timerGranularity = 1
	server.ConfRefreshTimeout = 600 / server.timerGranularity
	server.MinRefreshTimeout = 300 / server.timerGranularity
	server.timeout = time.Duration(server.timerGranularity) * time.Second
	server.timeoutCounter = 600
	server.retryCnt = 5
	server.minCnt = 1
	server.probeWait = 5
	server.probeNum = 5
	server.probeMax = 20
	server.probeMin = 10
	server.arpSliceRefreshDuration = time.Duration(10) * time.Minute
	server.dumpArpTable = false
}

func (server *ARPServer) sigHandler(sigChan <-chan os.Signal) {
	server.logger.Debug("Inside sigHandler....")
	signal := <-sigChan
	switch signal {
	case syscall.SIGHUP:
		server.logger.Debug("Received SIGHUP signal")
		server.printArpEntries()
		server.logger.Debug("Closing DB handler")
		if server.dbHdl != nil {
			server.dbHdl.Disconnect()
		}
		os.Exit(0)
	default:
		server.logger.Err(fmt.Sprintln("Unhandled signal : ", signal))
	}
}

func (server *ARPServer) initializeEvents() error {
	server.eventDbHdl = dbutils.NewDBUtil(server.logger)
	err := server.eventDbHdl.Connect()
	if err != nil {
		server.logger.Err("Failed to create the DB handle")
		return err
	}

	return eventUtils.InitEvents("ARPD", server.eventDbHdl, server.eventDbHdl, server.logger, 1000)
}

func (server *ARPServer) InitServer() {
	server.logger.Debug("Starting Arp Server")
	server.SysRsvdVlan = server.AsicdPlugin.GetSysRsvdVlan()
	server.logger.Debug("Listen for ASICd updates")
	var err error
	server.IsLinuxOnly, err = server.AsicdPlugin.IsLinuxOnlyPlugin()
	if err != nil {
		server.logger.Err("Error getting Plugin Info from Asicd", err)
	}
	if server.IsLinuxOnly == false {
		server.initArpParams()
		server.buildArpInfra()

		err = server.initiateDB()
		if err != nil {
			server.logger.Err(fmt.Sprintln("DB Initialization failure...", err))
		} else {
			server.logger.Debug("ArpCache DB has been initiated successfully...")
			server.updateArpCacheFromDB()
		}

		go server.updateArpCache()
		if server.dbHdl != nil {
			server.getArpGlobalConfig()
		}

		// Initialize Events
		err = server.initializeEvents()
		if err != nil {
			server.logger.Err(fmt.Sprintln("Unable to initialize events", err))
		}

		sigChan := make(chan os.Signal, 1)
		signalList := []os.Signal{syscall.SIGHUP}
		signal.Notify(sigChan, signalList...)
		go server.sigHandler(sigChan)
		go server.refreshArpSlice()
		go server.arpCacheTimeout()
	}
	server.FlushLinuxArpCache()
}

func (server *ARPServer) StartServer() {
	server.InitServer()
	server.InitDone <- true
	for {
		select {
		case arpConf := <-server.ArpConfCh:
			server.processArpConf(arpConf)
		case rConf := <-server.ResolveIPv4Ch:
			server.processResolveIPv4(rConf)
		case rConf := <-server.DeleteResolvedIPv4Ch:
			server.processDeleteResolvedIPv4(rConf.IpAddr)
		case arpEntryInfo, ok := <-server.DeleteArpEntryCh:
			if ok {
				server.processDeleteArpEntryInt(arpEntryInfo)
			}
			/*
				case garpInfo, ok := <-server.GarpEntryCh:
					if ok {
						server.processGarp(garpInfo)
					}
			*/
		case arpActionMsg := <-server.ArpActionCh:
			server.processArpAction(arpActionMsg)
		case msg := <-server.AsicdSubSocketCh:
			if server.IsLinuxOnly == false {
				server.processAsicdNotification(msg)
			}
		}
	}
}
