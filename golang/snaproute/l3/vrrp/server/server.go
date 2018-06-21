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
	"l3/arp/clientMgr"
	ndpClient "l3/ndp/lib"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"os"
	"os/signal"
	"syscall"
	"utils/asicdClient"
	"utils/dmnBase"
)

type VrrpServer struct {
	// All System Related Information
	SwitchPlugin       asicdClient.AsicdClientIntf
	ArpClient          arpClient.ArpdClientIntf
	NdpClient          ndpClient.NdpdClientIntf
	dmnBase            *dmnBase.FSBaseDmn
	GlobalConfig       *common.GlobalConfig
	globalState        common.GlobalState
	V4                 map[int32]*V4Intf
	V6                 map[int32]*V6Intf
	Intf               map[KeyInfo]VrrpInterface // key is struct IntfRef, VRID, Version which is KeyInfo
	v4Intfs            []KeyInfo                 // list of v4 vrrp interfaces that got created
	v6Intfs            []KeyInfo                 // list of v6 vrrp interfaces that got created
	V4IntfRefToIfIndex map[string]int32          // key is intfRef and value is ifIndex
	V6IntfRefToIfIndex map[string]int32          // key is intfRef and valud if ifIndex
	CfgCh              chan *common.IntfCfg      // Starting from hereAll Channels Used during Events
	GblCfgCh           chan *common.GlobalConfig
	L3IntfNotifyCh     chan *common.BaseIpInfo
	VirtualIpCh        chan *common.VirtualIpInfo // used for updating virtual ip state in hardware/linux
	UpdateTxCh         chan struct{}
	UpdateRxCh         chan struct{}
}

type VrrpTxChannelInfo struct {
	key      string
	priority uint16 // any value > 255 means ignore it
}

type KeyInfo struct {
	IntfRef string
	VRID    int32
	IpType  int
}

func (svr *VrrpServer) EventListener() {
	// Start receviing in rpc values in the channell
	for {
		select {

		case gCfg, ok := <-svr.GblCfgCh:
			if ok {
				svr.HandleGlobalConfig(gCfg)
			}
		case cfg, ok := <-svr.CfgCh:
			if ok {
				svr.HandleVrrpIntfConfig(cfg)
			}
		case l3NotifyInfo, ok := <-svr.L3IntfNotifyCh:
			if ok {
				svr.HandleIpNotification(l3NotifyInfo)
			}
		case vipUpdateInfo, ok := <-svr.VirtualIpCh:
			if ok {
				svr.UpdateVirtualIntf(vipUpdateInfo)
			}
		case _, ok := <-svr.UpdateRxCh:
			if ok {
				svr.updateRxCount()
			}
		case _, ok := <-svr.UpdateTxCh:
			if ok {
				svr.updateTxCount()
			}
		}
	}
}

func (svr *VrrpServer) GetSystemInfo() {
	// Get IP Information
	svr.GetIPIntfs()
}

func (svr *VrrpServer) InitGlobalDS() {
	svr.GlobalConfig = new(common.GlobalConfig)
	svr.V6 = make(map[int32]*V6Intf, VRRP_GLOBAL_INFO_DEFAULT_SIZE)
	svr.V4 = make(map[int32]*V4Intf, VRRP_GLOBAL_INFO_DEFAULT_SIZE)
	svr.Intf = make(map[KeyInfo]VrrpInterface, VRRP_GLOBAL_INFO_DEFAULT_SIZE)
	svr.V4IntfRefToIfIndex = make(map[string]int32, VRRP_GLOBAL_INFO_DEFAULT_SIZE)
	svr.V6IntfRefToIfIndex = make(map[string]int32, VRRP_GLOBAL_INFO_DEFAULT_SIZE)
	svr.GblCfgCh = make(chan *common.GlobalConfig)
	svr.CfgCh = make(chan *common.IntfCfg, VRRP_INTF_CONFIG_CH_SIZE)
	svr.L3IntfNotifyCh = make(chan *common.BaseIpInfo)
	svr.VirtualIpCh = make(chan *common.VirtualIpInfo)
	svr.UpdateTxCh = make(chan struct{})
	svr.UpdateRxCh = make(chan struct{})
}

func (svr *VrrpServer) DeAllocateMemory() {
	svr.V4 = nil
	svr.V6 = nil
	for _, intf := range svr.Intf {
		intf.DeInitVrrpIntf()
	}
	svr.Intf = nil
	svr.GlobalConfig = nil
	svr.V4IntfRefToIfIndex = nil
	svr.V6IntfRefToIfIndex = nil
	svr.GblCfgCh = nil
	svr.CfgCh = nil
	svr.L3IntfNotifyCh = nil
	svr.VirtualIpCh = nil
	svr.UpdateRxCh = nil
	svr.UpdateTxCh = nil
}

func (svr *VrrpServer) signalHandler(sigChannel <-chan os.Signal) {
	signal := <-sigChannel
	switch signal {
	case syscall.SIGHUP:
		debug.Logger.Alert("Received SIGHUP Signal")
		svr.DeAllocateMemory()
		debug.Logger.Alert("Exiting!!!!!!")
		os.Exit(0)
	default:
		debug.Logger.Info("Unhandled Signal:", signal)
	}
}

func (svr *VrrpServer) OSSignalHandle() {
	sigChannel := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChannel, signalList...)
	go svr.signalHandler(sigChannel)
}

func (svr *VrrpServer) VrrpStartServer() {
	svr.OSSignalHandle()
	svr.InitGlobalDS()
	svr.GetSystemInfo()
	svr.ReadDB()
	go svr.EventListener()
}

func VrrpNewServer(sPlugin asicdClient.AsicdClientIntf, arpClient arpClient.ArpdClientIntf, ndpClient ndpClient.NdpdClientIntf,
	dmnBase *dmnBase.FSBaseDmn) *VrrpServer {
	vrrpServer := &VrrpServer{}
	vrrpServer.SwitchPlugin = sPlugin
	vrrpServer.dmnBase = dmnBase
	vrrpServer.ArpClient = arpClient
	vrrpServer.NdpClient = ndpClient
	return vrrpServer
}

const (

	// Error Message
	VRRP_INVALID_VRID                   = "VRID is invalid"
	VRRP_CLIENT_CONNECTION_NOT_REQUIRED = "Connection to Client is not required"
	VRRP_SAME_OWNER                     = "Local Router should not be same as the VRRP Ip Address"
	VRRP_MISSING_VRID_CONFIG            = "VRID is not configured on interface"
	VRRP_CHECKSUM_ERR                   = "VRRP checksum failure"
	VRRP_INVALID_PCAP                   = "Invalid Pcap Handler"
	VRRP_VLAN_NOT_CREATED               = "Create Vlan before configuring VRRP"
	VRRP_IPV4_INTF_NOT_CREATED          = "Create IPv4 interface before configuring VRRP"
	VRRP_DATABASE_LOCKED                = "database is locked"

	// Default Size
	VRRP_GLOBAL_INFO_DEFAULT_SIZE         = 50
	VRRP_VLAN_MAPPING_DEFAULT_SIZE        = 5
	VRRP_INTF_STATE_SLICE_DEFAULT_SIZE    = 5
	VRRP_LINUX_INTF_MAPPING_DEFAULT_SIZE  = 5
	VRRP_INTF_IPADDR_MAPPING_DEFAULT_SIZE = 5
	VRRP_RX_BUF_CHANNEL_SIZE              = 5
	VRRP_TX_BUF_CHANNEL_SIZE              = 1
	VRRP_FSM_CHANNEL_SIZE                 = 1
	VRRP_INTF_CONFIG_CH_SIZE              = 1
	VRRP_TOTAL_INTF_CONFIG_ELEMENTS       = 7
)
