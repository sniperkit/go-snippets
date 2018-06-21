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
	"l3/ndp/config"
	"l3/ndp/debug"
	"l3/ndp/packet"
	"l3/ndp/publisher"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dmnBase"
)

func NDPNewServer(sPlugin asicdClntIntfs.AsicdClntIntf, dmnBase *dmnBase.FSBaseDmn) *NDPServer {
	svr := &NDPServer{}
	svr.SwitchPlugin = sPlugin
	svr.dmnBase = dmnBase
	// Profiling code for lldp
	prof, err := os.Create(NDP_CPU_PROFILE_FILE)
	if err == nil {
		pprof.StartCPUProfile(prof)
	}
	return svr
}

/* OS signal handler.
 *      If the process get a sighup signal then close all the pcap handlers.
 *      After that delete all the memory which was used during init process
 */
func (svr *NDPServer) SignalHandler(sigChannel <-chan os.Signal) {
	signal := <-sigChannel
	switch signal {
	case syscall.SIGHUP:
		debug.Logger.Alert("Received SIGHUP Signal")
		svr.DeInitGlobalDS()
		pprof.StopCPUProfile()
		debug.Logger.Alert("Exiting!!!!!")
		os.Exit(0)
	default:
		debug.Logger.Info("Unhandled Signal:", signal)
	}
}

/*  Create os signal handler channel and initiate go routine for that
 */
func (svr *NDPServer) OSSignalHandle() {
	sigChannel := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChannel, signalList...)
	go svr.SignalHandler(sigChannel)
}

func (svr *NDPServer) InitGlobalDS() {
	svr.L2Port = make(map[int32]PhyPort, NDP_SERVER_MAP_INITIAL_CAP)
	svr.SwitchMacMapEntries = make(map[string]struct{}, NDP_SERVER_MAP_INITIAL_CAP)
	svr.L3Port = make(map[int32]Interface, NDP_SERVER_MAP_INITIAL_CAP)
	svr.VlanInfo = make(map[int32]config.VlanInfo, NDP_SERVER_MAP_INITIAL_CAP)
	svr.virtualIp = make(map[int32]*config.VirtualIpInfo, NDP_SERVER_MAP_INITIAL_CAP)
	svr.VlanIfIdxVlanIdMap = make(map[string]int32, NDP_SERVER_MAP_INITIAL_CAP)
	svr.NeighborInfo = make(map[string]config.NeighborConfig, NDP_SERVER_MAP_INITIAL_CAP)
	svr.L3IfIntfRefToIfIndex = make(map[string]int32, NDP_SERVER_MAP_INITIAL_CAP)
	svr.PhyPortToL3PortMap = make(map[int32]L3Info)
	svr.Dot1QToVlanIfIndex = make(map[int32]int32)
	svr.IpIntfCh = make(chan *config.IPIntfNotification, NDP_SERVER_ASICD_NOTIFICATION_CH_SIZE)
	svr.VlanCh = make(chan *config.VlanNotification)
	svr.MacMoveCh = make(chan *config.MacMoveNotification)
	svr.RxPktCh = make(chan *RxPktInfo, NDP_SERVER_INITIAL_CHANNEL_SIZE)
	svr.PktDataCh = make(chan config.PacketData, NDP_SERVER_INITIAL_CHANNEL_SIZE)
	svr.ActionCh = make(chan *config.ActionData)
	svr.VirtualIpCh = make(chan *config.VirtualIpInfo)
	svr.SnapShotLen = 1024
	svr.Promiscuous = false
	svr.Timeout = 1 * time.Second
	svr.NeigborEntryLock = &sync.RWMutex{}
	svr.Packet = packet.Init()

	//configuration channels
	svr.GlobalCfg = make(chan NdpConfig)

	// init publisher
	pub := publisher.NewPublisher()
	pub.InitPublisher()
	svr.notifyChan = pub.PubChan.All
}

func (svr *NDPServer) DeInitGlobalDS() {
	svr.L2Port = nil
	svr.L3Port = nil
	svr.IpIntfCh = nil
	svr.VlanCh = nil
	svr.RxPktCh = nil
}

/*
 * API: it will collect all ipv6 interface ports from the system... If needed we can collect port information
 *      also from the system.
 *	After the information is collected, if the oper state is up then we will start rx/tx
 */
func (svr *NDPServer) InitSystem() {
	// Get ports information
	svr.GetPorts()

	// Get vlans information
	svr.GetVlans()

	// Get IP Information
	svr.GetIPIntf()

	// Get Virtual IP Information
	svr.getVirtualIpIntf()

	// Check status of IP Interface and then start RX/TX for that ip interface
	for _, ipIntf := range svr.L3Port {
		if ipIntf.OperState == config.STATE_UP {
			svr.StartRxTx(ipIntf.IfIndex)
		}
	}
}

func (svr *NDPServer) UpdateInterfaceTimers() {
	for key, intf := range svr.L3Port {
		intf.UpdateTimer(svr.NdpConfig)
		svr.L3Port[key] = intf
	}
}

func (svr *NDPServer) EventsListener() {
	for {
		select {
		// global configuration channel
		case globalCfg, ok := <-svr.GlobalCfg:
			if !ok {
				continue
			}
			update := svr.NdpConfig.Create(globalCfg)
			if update {
				svr.UpdateInterfaceTimers()
			}
		case vlanInfo, ok := <-svr.VlanCh:
			if !ok {
				continue
			}
			svr.HandleVlanNotification(vlanInfo)
		// ipv6 interface create/delete state up/down notification channel
		case ipIntfNotify := <-svr.IpIntfCh:
			switch ipIntfNotify.Operation {
			case config.CONFIG_CREATE, config.CONFIG_DELETE:
				svr.HandleIPIntfCreateDelete(ipIntfNotify)
			case config.STATE_UP, config.STATE_DOWN:
				// we need to received l2, l3 state up notification via one channel only
				// by doing so we will maintain the order in which the state notifications are
				// coming
				switch ipIntfNotify.IpAddr {
				case config.L2_NOTIFICATION:
					phyPortStateCh := &config.PortState{
						IfIndex: ipIntfNotify.IfIndex,
						IfState: ipIntfNotify.Operation,
					}
					svr.HandlePhyPortStateNotification(phyPortStateCh)
				default:
					svr.HandleStateNotification(ipIntfNotify)
				}
			}
		// packet rx channel
		case rxChInfo, ok := <-svr.RxPktCh:
			if !ok {
				continue
			}
			svr.counter.Rcvd++
			svr.ProcessRxPkt(rxChInfo.ifIndex, rxChInfo.pkt)
		// packet tx channel on timer expiry
		case pktData, ok := <-svr.PktDataCh:
			if !ok {
				continue
			}
			svr.ProcessTimerExpiry(pktData)
		// mac move notification channel
		case macMoveInfo, ok := <-svr.MacMoveCh:
			if !ok {
				continue
			}
			svr.SoftwareUpdateNbrEntry(macMoveInfo)
		// action notification
		case actionData, ok := <-svr.ActionCh:
			if !ok {
				continue
			}
			svr.HandleAction(actionData)
		case vipInfo, ok := <-svr.VirtualIpCh:
			if ok {
				svr.HandleVirtualIpNotification(vipInfo)
			}
		}
	}
}

/*  ndp server:
 * 1) OS Signal Handler
 * 2) Read from DB and close DB
 * 3) Connect to all the clients
 * 4) Call AsicPlugin for port information
 * 5) go routine to handle all the channels within lldp server
 */

func (svr *NDPServer) NDPStartServer() {
	svr.OSSignalHandle()
	svr.ReadDB()
	svr.InitGlobalDS()
	svr.InitSystem()
	go svr.EventsListener()
}
