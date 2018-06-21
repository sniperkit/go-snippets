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
	_ "fmt"
	"l2/lldp/config"
	"l2/lldp/plugin"
	"l2/lldp/utils"
	_ "models/objects"
	"os"
	"os/signal"
	_ "runtime/pprof"
	"strconv"
	"syscall"
	"time"
	"utils/dbutils"
)

/* Create lldp server object for the main handler..
 */
func LLDPNewServer(aPlugin plugin.AsicIntf, lPlugin plugin.ConfigIntf, sPlugin plugin.SystemIntf,
	dbHdl *dbutils.DBUtil) *LLDPServer {
	lldpServerInfo := &LLDPServer{
		asicPlugin: aPlugin,
		CfgPlugin:  lPlugin,
		SysPlugin:  sPlugin,
		lldpDbHdl:  dbHdl,
	}
	// Allocate memory to all the Data Structures
	lldpServerInfo.InitGlobalDS()
	/*
		// Profiling code for lldp
		prof, err := os.Create(LLDP_CPU_PROFILE_FILE)
		if err == nil {
			pprof.StartCPUProfile(prof)
		}
	*/
	return lldpServerInfo
}

/* Allocate memory to all the object which are being used by LLDP server
 */
func (svr *LLDPServer) InitGlobalDS() {
	svr.lldpGblInfo = make(map[int32]LLDPGlobalInfo, LLDP_INITIAL_GLOBAL_INFO_CAPACITY)
	svr.lldpIntfRef2IfIndexMap = make(map[string]int32, LLDP_INITIAL_GLOBAL_INFO_CAPACITY)
	svr.lldpRxPktCh = make(chan InPktChannel, LLDP_RX_PKT_CHANNEL_SIZE)
	svr.lldpTxPktCh = make(chan SendPktChannel, LLDP_TX_PKT_CHANNEL_SIZE)
	svr.lldpSnapshotLen = 1024
	svr.lldpPromiscuous = true
	// LLDP Notifications are atleast 5 seconds apart with default being
	// 30 seconds. So, we can have the leavrage the pcap timeout (read from
	// buffer) to be 1 second.
	svr.lldpTimeout = 500 * time.Millisecond
	svr.GblCfgCh = make(chan *config.Global)
	svr.IntfCfgCh = make(chan *config.IntfConfig, LLDP_PORT_CONFIG_CHANNEL_SIZE)
	svr.IfStateCh = make(chan *config.PortState, LLDP_PORT_STATE_CHANGE_CHANNEL_SIZE)
	svr.UpdateCacheCh = make(chan *config.SystemInfo, 1)
	svr.EventCh = make(chan config.EventInfo, 10)
	svr.PortAttrCh = make(chan *config.PortAttrMsg)
	svr.counter.Send = 0
	svr.counter.Rcvd = 0
	// All Plugin Info
}

/* De-Allocate memory to all the object which are being used by LLDP server
 */
func (svr *LLDPServer) DeInitGlobalDS() {
	// close rx packet channel
	close(svr.lldpRxPktCh)
	close(svr.lldpTxPktCh)
	svr.lldpRxPktCh = nil
	svr.lldpTxPktCh = nil
	svr.lldpGblInfo = nil
}

/* On de-init we will be closing all the pcap handlers that are opened up
 * We will also free up all the pointers from the gblInfo. Otherwise that will
 * lead to memory leak
 */
func (svr *LLDPServer) CloseAllPktHandlers() {

	// close pcap, stop cache timer and free any allocated memory
	for i := 0; i < len(svr.lldpIntfStateSlice); i++ {
		key := svr.lldpIntfStateSlice[i]
		intf, exists := svr.lldpGblInfo[key]
		if !exists {
			continue
		}
		intf.DeInitRuntimeInfo()
		svr.lldpGblInfo[key] = intf
	}
	debug.Logger.Info("closed everything")
}

/* Create global run time information for l2 port and then start rx/tx for that port if state is up
 */
func (svr *LLDPServer) InitL2PortInfo(portInfo *config.PortInfo) {
	intf, _ := svr.lldpGblInfo[portInfo.IfIndex]
	intf.InitRuntimeInfo(portInfo)
	// on fresh start it will not exists but on restart it might
	// default is set to true but LLDP Object is auto-discover and hence we will enable it manually
	// we will overwrite the value based on dbReead but default should always be true
	intf.Enable()
	// set default value for RX/TX mode so that it can be updated while db read
	intf.rxtxMode = config.TXRX
	svr.lldpGblInfo[portInfo.IfIndex] = intf
	svr.lldpIntfStateSlice = append(svr.lldpIntfStateSlice, intf.Port.IfIndex)
	debug.Logger.Debug("Updating IntfRef Map (key, value):(", portInfo.Name, ",", portInfo.IfIndex, ")")
	svr.lldpIntfRef2IfIndexMap[portInfo.Name] = portInfo.IfIndex
}

/*  lldp server: 1) Connect to all the clients
 *		 2) Initialize DB
 *		 3) Read from DB and close DB
 *		 4) Call AsicPlugin for port information
 *		 5) go routine to handle all the channels within lldp server
 */
func (svr *LLDPServer) LLDPStartServer(paramsDir string) {
	// OS Signal channel listener thread
	svr.OSSignalHandle()

	svr.paramsDir = paramsDir
	// Initialize DB
	err := svr.InitDB()
	if err != nil {
		debug.Logger.Err("DB init failed")
	}

	// Start asicd plugin before you do get bulk
	svr.asicPlugin.Start()
	svr.SysPlugin.Start()

	// Get Port Information from Asic, only after reading from DB
	portsInfo := svr.asicPlugin.GetPortsInfo()
	for _, port := range portsInfo {
		svr.InitL2PortInfo(port)
	}
	// Get System Information from Sysd, before we start anything
	svr.SysInfo = svr.SysPlugin.GetSystemInfo(svr.lldpDbHdl)

	// Populate Gbl Configs
	svr.ReadDB()

	// after everything is started then Do Rx/Tx Init
	svr.RunGlobalConfig()
	go svr.ChannelHandler()
}

/*  Create os signal handler channel and initiate go routine for that
 */
func (svr *LLDPServer) OSSignalHandle() {
	sigChannel := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChannel, signalList...)
	go svr.SignalHandler(sigChannel)
}

/* OS signal handler.
 *      If the process get a sighup signal then close all the pcap handlers.
 *      After that delete all the memory which was used during init process
 */
func (svr *LLDPServer) SignalHandler(sigChannel <-chan os.Signal) {
	signal := <-sigChannel
	switch signal {
	case syscall.SIGHUP:
		debug.Logger.Alert("Received SIGHUP Signal")
		svr.CloseAllPktHandlers()
		svr.DeInitGlobalDS()
		svr.CloseDB()
		//pprof.StopCPUProfile()
		debug.Logger.Alert("Exiting!!!!!")
		os.Exit(0)
	default:
		debug.Logger.Info("Unhandled Signal:", signal)
	}
}

/* Create l2 port pcap handler and then start rx and tx on that pcap
 *	Filter is LLDP_BPF_FILTER = "ether proto 0x88cc"
 * Note: API should only and only do
 *  1) pcap create
 *  2) start go routine for Rx/Tx Frames Packet Handler
 *  3) Add the port to UP List
 */
func (svr *LLDPServer) StartRxTx(ifIndex int32, rxtxMode uint8) {
	intf, exists := svr.lldpGblInfo[ifIndex]
	if !exists {
		debug.Logger.Err("No entry for ifindex", ifIndex)
		return
	}
	// if the port is disabled or lldp globally is disabled then no need to start rx/tx...
	if svr.Global.Enable == false || intf.Port.OperState != LLDP_PORT_STATE_UP {
		debug.Logger.Info("Port is down and hence not starting pcap handler yet for", intf.Port.Name)
		return
	}
	if intf.PcapHandle == nil {
		err := intf.CreatePcapHandler(svr.lldpSnapshotLen, svr.lldpPromiscuous, svr.lldpTimeout)
		if err != nil {
			debug.Logger.Alert("Creating Pcap Handler for", intf.Port.Name,
				"failed and hence we will not start LLDP on the port")
			return
		}
		debug.Logger.Info("Pcap created for lldp frames for port:", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
	}
	svr.AddPortToUpState(intf.Port.IfIndex)
	// Everything set up, so now lets start with receiving frames and transmitting frames go routine...
	// If RX routine not running start it
	if rxtxMode != config.TX_ONLY {
		debug.Logger.Debug("mode is not tx only")
		if !intf.RxInfo.RxRunning {
			debug.Logger.Debug("start ReceiveFrames go routine for:", intf.Port.Name)
			go intf.ReceiveFrames(svr.lldpRxPktCh)
			intf.RxInfo.RxRunning = true
		}
	} else {
		// RX go routine could have been spawned due to earlier txrx or rx only modes
		if intf.RxInfo.RxRunning {
			debug.Logger.Debug("RX go routine could have been spawned due to earlier txrx or rx only modes, and hence killing rx routine for:", intf.Port.Name)
			intf.RxKill <- true
			intf.RxInfo.RxRunning = false
			<-intf.RxKill
			intf.counter.Rcvd = 0
		}
	}
	// If TX routine not running start it
	if rxtxMode != config.RX_ONLY {
		debug.Logger.Debug("mode is not rx only")
		if intf.TxInfo.TxTimer == nil {
			debug.Logger.Debug("Start tx timer for:", intf.Port.Name)
			intf.StartTxTimer(svr.lldpTxPktCh)
		}
	} else {
		// TX go routine could have been spawned due to earlier txrx or tx only modes
		debug.Logger.Debug("TX go routine could have been spawned due to earlier txrx or tx only mode stopping tx timer for:", intf.Port.Name)
		intf.TxInfo.StopTxTimer()
		intf.counter.Send = 0
	}
	debug.Logger.Info("started rx tx for port:", intf.Port.Name)
	svr.lldpGblInfo[ifIndex] = intf
	return
}

/*  Send Signal for stopping rx/tx go routine and timers as the pcap handler for
 *  the port is deleted
 */
func (svr *LLDPServer) StopRxTx(ifIndex int32) {
	intf, exists := svr.lldpGblInfo[ifIndex]
	if !exists {
		debug.Logger.Err("No entry for ifIndex", ifIndex)
		return
	}

	// stop the timer
	debug.Logger.Debug("stop tx timer for:", intf.Port.Name)
	intf.TxInfo.StopTxTimer()
	// Delete Pcap Handler
	debug.Logger.Debug("delete pcap handler for:", intf.Port.Name)
	intf.DeletePcapHandler()
	// invalid the cache information
	debug.Logger.Debug("invalidate cache information for:", intf.Port.Name)
	intf.TxInfo.DeleteCacheFrame()
	intf.counter.Rcvd = 0
	intf.counter.Send = 0
	svr.lldpGblInfo[ifIndex] = intf
	debug.Logger.Info("Stop lldp frames rx/tx for port:", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
	svr.DeletePortFromUpState(ifIndex)
}

/*  helper function to inform whether rx channel is closed or open...
 *  Go routine can be exited using this information
 */
func (svr *LLDPServer) ServerRxChClose() bool {
	if svr.lldpRxPktCh == nil {
		return true
	}
	return false
}

/*  delete ifindex from lldpUpIntfStateSlice on port down... we can use this
 *  if user decides to disable lldp on a port
 */
func (svr *LLDPServer) DeletePortFromUpState(ifIndex int32) {
	for idx, _ := range svr.lldpUpIntfStateSlice {
		if svr.lldpUpIntfStateSlice[idx] == ifIndex {
			svr.lldpUpIntfStateSlice = append(svr.lldpUpIntfStateSlice[:idx],
				svr.lldpUpIntfStateSlice[idx+1:]...)
			break
		}
	}
}

/*
 *  Add ifIndex to lldpUpIntfStateSlice on start rx/tx only if it doesn't exists
 */
func (svr *LLDPServer) AddPortToUpState(ifIndex int32) {
	for idx, _ := range svr.lldpUpIntfStateSlice {
		if svr.lldpUpIntfStateSlice[idx] == ifIndex {
			debug.Logger.Alert("Duplicate ADD request for ifIndex:", ifIndex)
			return
		}
	}
	svr.lldpUpIntfStateSlice = append(svr.lldpUpIntfStateSlice, ifIndex)
}

/*  handle l2 state up/down notifications..
 */
func (svr *LLDPServer) UpdateL2IntfStateChange(ifIndex int32, state string) {
	intf, found := svr.lldpGblInfo[ifIndex]
	if !found {
		return
	}
	switch state {
	case "UP":
		debug.Logger.Debug("State UP notification for " + intf.Port.Name + " ifIndex: " + strconv.Itoa(int(intf.Port.IfIndex)))
		intf.Port.OperState = LLDP_PORT_STATE_UP
		svr.lldpGblInfo[ifIndex] = intf
		if intf.isEnabled() {
			// Create Pcap Handler and start rx/tx packets
			svr.StartRxTx(ifIndex, intf.rxtxMode)
		}
	case "DOWN":
		debug.Logger.Debug("State DOWN notification for " + intf.Port.Name + " ifIndex: " + strconv.Itoa(int(intf.Port.IfIndex)))
		intf.Port.OperState = LLDP_PORT_STATE_DOWN
		svr.lldpGblInfo[ifIndex] = intf
		if intf.isEnabled() {
			// Delete Pcap Handler and stop rx/tx packets
			svr.StopRxTx(ifIndex)
		}
	}
}

/*  handle global lldp enable/disable, which will enable/disable lldp for all the ports
 */
func (svr *LLDPServer) handleGlobalConfig() {
	if svr.Global == nil {
		return
	}
	if len(svr.lldpIntfStateSlice) == 0 {
		debug.Logger.Err("No ports on the system")
		return
	}
	debug.Logger.Debug("Doing global init for all the ports in up state", svr.lldpIntfStateSlice,
		"global Info:", *svr.Global)
	// iterate over all the entries in the gblInfo and change the state accordingly
	for _, ifIndex := range svr.lldpIntfStateSlice {
		intf, found := svr.lldpGblInfo[ifIndex]
		if !found {
			debug.Logger.Err("No entry for ifIndex", ifIndex, "in runtime information")
			continue
		}
		debug.Logger.Debug("Init for intf:", intf.Port.Name, "and intf information is",
			intf.isDisabled(), intf.Port.OperState)
		// faster operation
		if intf.isDisabled() || intf.Port.OperState != LLDP_PORT_STATE_UP {
			debug.Logger.Debug("Cannot start LLDP rx/tx for port", intf.Port.Name,
				"as its state is", intf.Port.OperState, "enable is", intf.isDisabled())
			continue
		}
		switch svr.Global.Enable {
		case true:
			debug.Logger.Debug("Global Config Enabled, enabling port rx tx for port:", intf.Port.Name,
				"ifIndex", ifIndex)
			svr.StartRxTx(ifIndex, svr.Global.TxRxMode)
		case false:
			debug.Logger.Debug("Global Config Disabled, disabling port rx tx for port:", intf.Port.Name,
				"ifIndex", ifIndex)
			// do not update the configuration enable/disable state...just stop packet handling
			svr.StopRxTx(ifIndex)
		}
	}

	if svr.Global.Enable == false {
		svr.counter.Rcvd = 0
		svr.counter.Send = 0
	}
}

/*  handle configuration coming from user, which will enable/disable lldp per port
 */
func (svr *LLDPServer) handleIntfConfig(ifIndex int32, enable bool, rxtxMode uint8) {
	intf, found := svr.lldpGblInfo[ifIndex]
	if !found {
		debug.Logger.Err("No entry for ifIndex", ifIndex, "in runtime information")
		return
	}
	intf.rxtxMode = rxtxMode
	switch enable {
	case true:
		debug.Logger.Debug("Config Enable for", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
		intf.Enable()
		svr.lldpGblInfo[ifIndex] = intf
		svr.StartRxTx(ifIndex, rxtxMode)
	case false:
		debug.Logger.Debug("Config Disable for", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
		if intf.isEnabled() { // If Enabled then only do stop rx/tx
			intf.Disable()
			svr.lldpGblInfo[ifIndex] = intf
			svr.StopRxTx(ifIndex)
		}
	}
}

/*  API to send a frame when tx timer expires per port
 */
func (svr *LLDPServer) SendFrame(ifIndex int32) {
	intf, exists := svr.lldpGblInfo[ifIndex]
	// extra check for pcap handle
	if exists && intf.PcapHandle != nil {
		rv := intf.WritePacket(intf.TxInfo.Frame(intf.Port, svr.SysInfo))
		if rv == false {
			intf.TxInfo.SetCache(rv)
		}
	}
	debug.Logger.Debug("Frame send from port:", intf.Port.Name)
	intf.StartTxTimer(svr.lldpTxPktCh)
	intf.counter.Send++
	svr.counter.Send++
	svr.lldpGblInfo[ifIndex] = intf
}

func (svr *LLDPServer) ProcessRcvdPkt(rcvdInfo InPktChannel) {
	intf, exists := svr.lldpGblInfo[rcvdInfo.ifIndex]
	if !exists {
		return
	}
	debug.Logger.Debug("Process Packet Received on port:", intf.Port.Name)
	var err error
	eventInfo := config.EventInfo{}
	intf.RxLock.Lock()
	eventInfo.EventType, err = intf.RxInfo.Process(intf.RxInfo, rcvdInfo.pkt)
	if err != nil {
		intf.RxLock.Unlock()
		debug.Logger.Err("err", err, "while processing rx frame on port",
			intf.Port.Name)
		return
	}
	intf.pktRcvdTime = time.Now()
	intf.counter.Rcvd++
	svr.counter.Rcvd++
	intf.RxLock.Unlock()
	// reset/start timer for recipient information
	intf.RxInfo.CheckPeerEntry(intf.Port.Name, svr.EventCh, rcvdInfo.ifIndex)
	svr.lldpGblInfo[rcvdInfo.ifIndex] = intf
	eventInfo.IfIndex = rcvdInfo.ifIndex

	if eventInfo.EventType != config.NoOp {
		svr.SysPlugin.PublishEvent(eventInfo)
	}
	debug.Logger.Debug("Done Processing Packet for port:", intf.Port.Name)
}

func (svr *LLDPServer) handlePortAttrChange(portAttrMsg *config.PortAttrMsg) {
	debug.Logger.Info("handle port attribute notification msg:", *portAttrMsg)
	intf, exists := svr.lldpGblInfo[portAttrMsg.IfIndex]
	if !exists {
		return
	}
	if portAttrMsg.Description != "" {
		debug.Logger.Debug("Updating Description for Port:", intf.Port.Name)
		intf.Port.Description = portAttrMsg.Description
	}
	if portAttrMsg.Pvid != config.INVALID_PVID {
		debug.Logger.Debug("Updating Pvid for Port:", intf.Port.Name)
		intf.Port.Pvid = portAttrMsg.Pvid
	}
	intf.TxInfo.SetCache(false)
	svr.lldpGblInfo[portAttrMsg.IfIndex] = intf
}

/* To handle all the channels in lldp server... For detail look at the
 * LLDPInitGlobalDS api to see which all channels are getting initialized
 */
func (svr *LLDPServer) ChannelHandler() {

	for {
		select {
		case rcvdInfo, ok := <-svr.lldpRxPktCh:
			if !ok {
				continue // rx channel should be closed only on exit
			}
			svr.ProcessRcvdPkt(rcvdInfo)
		case info, ok := <-svr.lldpTxPktCh:
			if !ok {
				continue
			}
			svr.SendFrame(info.ifIndex)
		case gbl, ok := <-svr.GblCfgCh: // Change in global config
			if !ok {
				debug.Logger.Err("Invalid Value Received on Global Config Channel")
				continue
			}
			debug.Logger.Info("Server Received Global Config", *gbl)
			if svr.Global == nil {
				debug.Logger.Info("Doing Global Config during auto-create")
				svr.Global = &config.Global{}
			}
			svr.Global.Enable = gbl.Enable
			svr.Global.Vrf = gbl.Vrf
			svr.Global.TranmitInterval = gbl.TranmitInterval
			svr.Global.TxRxMode = gbl.TxRxMode
			svr.Global.SnoopAndDrop = gbl.SnoopAndDrop
			// start all interface rx/tx in go routine only
			// @TODO: jgheewala fixme for update in transmit interval
			svr.handleGlobalConfig()
		case intf, ok := <-svr.IntfCfgCh: // Change in interface config
			if !ok {
				continue
			}
			debug.Logger.Info("Server received Intf Config", intf)
			svr.handleIntfConfig(intf.IfIndex, intf.Enable, intf.TxRxMode)
		case ifState, ok := <-svr.IfStateCh: // Change in Port State..
			if !ok {
				continue
			}
			debug.Logger.Info("Server received L2 Intf State Changes for ifIndex:", ifState.IfIndex,
				"state:", ifState.IfState)
			svr.UpdateL2IntfStateChange(ifState.IfIndex, ifState.IfState)
		case sysInfo, ok := <-svr.UpdateCacheCh:
			if !ok {
				continue
			}
			svr.UpdateCache(sysInfo)
		case eventInfo, ok := <-svr.EventCh: //used only for delete
			if !ok {
				continue
			}
			svr.SysPlugin.PublishEvent(eventInfo)
		case portAttrMsg, ok := <-svr.PortAttrCh:
			if !ok {
				continue
			}
			svr.handlePortAttrChange(portAttrMsg)
		}
	}
}

func (svr *LLDPServer) RunGlobalConfig() {
	// Only start rx/tx if, Globally LLDP is enabled, Interface LLDP is enabled and port is in UP state
	// move RX/TX to Channel Handler
	// The below step is really important for us.
	// On Re-Start if lldp global is enable then we will start rx/tx for those ports which are in up state
	// and at the same time we will start the loop for signal handler
	// if fresh start then svr.Global is nil as no global config is done and hence it will be no-op
	// however on re-start lets say you have 100 ports that have lldp running on it in that case your writer
	// channel will create a deadlock as the reader is not yet started... To avoid this we spawn go-routine
	// for handling Global Config before Channel Handler is started
	if svr.Global != nil {
		svr.handleGlobalConfig()
	}
}
