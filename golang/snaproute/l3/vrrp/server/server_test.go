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
	arpFS "l3/arp/clientMgr/flexswitch"
	ndpClient "l3/ndp/lib"
	ndpFS "l3/ndp/lib/flexswitch"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"log/syslog"
	"net"
	"reflect"
	"syscall"
	"testing"
	"time"
	asicdClnt "utils/asicdClient/mock"
	"utils/commonDefs"
	"utils/logging"
)

var testSvr *VrrpServer

func TestServerInit(t *testing.T) {
	var err error
	logger := new(logging.Writer)
	logger.MyComponentName = "VRRPD"
	logger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "VRRPFSMTEST")
	if err != nil {
		t.Error("failed to initialize syslog:", err)
		return
	} else {
		logger.MyLogLevel = 9 // trace level
		debug.SetLogger(logger)
	}
	testAsicdClnt := &asicdClnt.MockAsicdClientMgr{}
	testNdpInst := ndpFS.NdpdClientStruct{}
	testArpInst := arpFS.ArpdClientStruct{}
	testArpClnt := arpClient.NewArpdClient(commonDefs.MOCK_PLUGIN, "", make([]commonDefs.ClientJson, 0), testArpInst)
	testNdpClnt := ndpClient.NewNdpClient(commonDefs.MOCK_PLUGIN, "", make([]commonDefs.ClientJson, 0), testNdpInst)
	testSvr = VrrpNewServer(testAsicdClnt, testArpClnt, testNdpClnt, nil)
	if testSvr == nil {
		t.Error("Initializing Vrrp Server Failed")
		return
	}
	testSvr.VrrpStartServer()

	if testSvr.V4 == nil || testSvr.V6 == nil || testSvr.GlobalConfig == nil {
		t.Error("initializing vrrp server ds failed")
		return
	}
}

func TestServerDeInit(t *testing.T) {
	if testSvr == nil {
		return
	}
	testSvr.DeAllocateMemory()
	if testSvr.V4 != nil {
		t.Error("failed to de-init V4 interfaces")
		return
	}
	if testSvr.V6 != nil {
		t.Error("failed to de-init V6 interfaces")
		return
	}
	if testSvr.VirtualIpCh != nil {
		t.Error("failed to de-init VirtualIpCh")
		return
	}

	if testSvr.dmnBase != nil {
		t.Error("failed to de-init dmnBase")
		return
	}

	if testSvr.UpdateTxCh != nil {
		t.Error("failed to de-init UpdateTxCh")
		return
	}

	if testSvr.UpdateRxCh != nil {
		t.Error("failed to de-init UpdateRxCh")
		return
	}

	if testSvr.L3IntfNotifyCh != nil {
		t.Error("failed to de-init L3IntfNotifyCh")
		return
	}

	if testSvr.Intf != nil {
		t.Error("failed to de-init Interface information")
		return
	}

	if testSvr.GlobalConfig != nil {
		t.Error("failed to de-init GlobalConfig")
		return
	}

	if testSvr.GblCfgCh != nil {
		t.Error("failed to de-init GblCfgCh")
		return
	}

	if testSvr.CfgCh != nil {
		t.Error("failed to de-init CfgCh")
		return
	}
	testSvr = nil
}

func goToSleep() {
	time.Sleep(75 * time.Millisecond)
}

func TestGlobalConfig(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	gblCfg := &common.GlobalConfig{
		Vrf:       "default",
		Enable:    false,
		Operation: common.CREATE,
	}
	testSvr.GblCfgCh <- gblCfg
	goToSleep()
	wantGblState := &common.GlobalState{
		Vrf:           gblCfg.Vrf,
		Status:        gblCfg.Enable,
		V4Intfs:       0,
		V6Intfs:       0,
		TotalRxFrames: 0,
		TotalTxFrames: 0,
	}
	state := testSvr.GetGlobalState(gblCfg.Vrf)
	if !reflect.DeepEqual(state, wantGblState) {
		t.Error("Vrrp Global State Mis-Match During Global Create")
		t.Error("	    wantGblState:", *wantGblState)
		t.Error("	    rcvdGblState:", *state)
		return
	}
	gblCfg.Enable = true
	testSvr.GblCfgCh <- gblCfg
	goToSleep()
	wantGblState = &common.GlobalState{
		Vrf:           gblCfg.Vrf,
		Status:        gblCfg.Enable,
		V4Intfs:       0,
		V6Intfs:       0,
		TotalRxFrames: 0,
		TotalTxFrames: 0,
	}
	state = testSvr.GetGlobalState(gblCfg.Vrf)
	if !reflect.DeepEqual(state, wantGblState) {
		t.Error("Vrrp Global State Mis-Match During Global Create")
		t.Error("	    wantGblState:", *wantGblState)
		t.Error("	    rcvdGblState:", *state)
		return
	}
	TestServerDeInit(t)
}

var testIfIndex = int32(100)
var ipv4Intf = &common.BaseIpInfo{
	IfIndex:   testIfIndex,
	IntfRef:   "lo",
	IpAddr:    "172.18.0.3/24",
	OperState: common.STATE_DOWN,
	IpType:    syscall.AF_INET,
}

func TestIPv4Notifications(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	ipIntf := ipv4Intf
	ipIntf.MsgType = common.IP_MSG_CREATE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V4) != 1 {
		t.Error("failed to handle ipv4 interface create notification by create a new v4 entry")
		t.Error("	    ", testSvr.V4)
		return
	}
	ipIntf.MsgType = ""
	v4Entry := testSvr.V4[testIfIndex]
	if !reflect.DeepEqual(v4Entry.Cfg.Info, *ipIntf) {
		t.Error("failed to copy base ipv4 interface")
		t.Error("	    wantedIpInfo:", *ipIntf)
		t.Error("	    gotIpInfo:", v4Entry.Cfg.Info)
		return
	}

	ipIntf.MsgType = common.IP_MSG_STATE_CHANGE
	ipIntf.OperState = common.STATE_UP
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V4) != 1 {
		t.Error("failed to handle ipv4 interface update notification by updating v4 entry")
		t.Error("	    ", testSvr.V4)
		return
	}
	ipIntf.MsgType = ""
	if !reflect.DeepEqual(v4Entry.Cfg.Info, *ipIntf) {
		t.Error("failed to copy base ipv4 interface")
		t.Error("	    wantedIpInfo:", *ipIntf)
		t.Error("	    gotIpInfo:", v4Entry.Cfg.Info)
		return
	}

	ipIntf.MsgType = common.IP_MSG_DELETE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V4) == 1 {
		t.Error("failed to handle ipv4 interface delete notification by deleting existing v4 entry")
		t.Error("	    ", *testSvr.V4[testIfIndex])
		return
	}
	TestServerDeInit(t)
}

var ipv6Intf = &common.BaseIpInfo{
	IfIndex:   testIfIndex,
	IntfRef:   "lo",
	IpAddr:    "fe80::d898:ddff:fe7b:975a/64",
	OperState: common.STATE_DOWN,
	IpType:    syscall.AF_INET6,
}

func TestIPv6Notifications(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	ipIntf := ipv6Intf
	ipIntf.MsgType = common.IP_MSG_CREATE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V6) != 1 {
		t.Error("failed to handle ipv6 interface create notification by create a new v6 entry")
		t.Error("	    ", testSvr.V4)
		return
	}
	ipIntf.MsgType = ""
	v6Entry := testSvr.V6[testIfIndex]
	if !reflect.DeepEqual(v6Entry.Cfg.Info, *ipIntf) {
		t.Error("failed to copy base ipv6 interface")
		t.Error("	    wantedIpInfo:", *ipIntf)
		t.Error("	    gotIpInfo:", v6Entry.Cfg.Info)
		return
	}

	ipIntf.MsgType = common.IP_MSG_STATE_CHANGE
	ipIntf.OperState = common.STATE_UP
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V6) != 1 {
		t.Error("failed to handle ipv6 interface update notification by updating v6 entry")
		t.Error("	    ", testSvr.V6)
		return
	}
	ipIntf.MsgType = ""
	if !reflect.DeepEqual(v6Entry.Cfg.Info, *ipIntf) {
		t.Error("failed to update ipv6 interface")
		t.Error("	    wantedIpInfo:", *ipIntf)
		t.Error("	    gotIpInfo:", v6Entry.Cfg.Info)
		return
	}
	ipIntf.IpAddr = "3000::1/64"
	ipIntf.MsgType = common.IP_MSG_CREATE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V6) != 1 {
		t.Error("failed to handle ipv6 interface update notification by updating v6 entry")
		t.Error("	    ", testSvr.V6)
		return
	}
	if v6Entry.Cfg.GlobalScopeIp != ipIntf.IpAddr {
		t.Error("failed to update ipv6 global scope ip Addr")
		return
	}

	ipIntf.MsgType = common.IP_MSG_DELETE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V6) != 1 {
		t.Error("failed to handle deleting one ipv6 address from the interface")
		t.Error("	    ", testSvr.V6)
		return
	}
	ipIntf.IpAddr = "fe80::d898:ddff:fe7b:975a/64"
	ipIntf.MsgType = common.IP_MSG_DELETE
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	if len(testSvr.V6) == 1 {
		t.Error("failed to handle ipv6 interface delete notification by deleting existing v6 entry")
		t.Error("	    ", *testSvr.V6[testIfIndex])
		return
	}
	TestServerDeInit(t)
}

var testIntfBaseCfg = common.IntfCfg{
	IfIndex:               testIfIndex,
	IntfRef:               "lo",
	VRID:                  1,
	Priority:              101,
	VirtualIPAddr:         "172.18.0.1/32",
	AdvertisementInterval: 1,
	PreemptMode:           true,
	AcceptMode:            false,
}

func testEnableGlobalConfig() {
	testSvr.GlobalConfig.Vrf = "default"
	testSvr.GlobalConfig.Enable = false
	gblCfg := &common.GlobalConfig{
		Vrf:       "default",
		Enable:    true,
		Operation: common.UPDATE,
	}
	testSvr.GblCfgCh <- gblCfg
	goToSleep()
}

func testDisableGlobalConfig() {
	gblCfg := &common.GlobalConfig{
		Vrf:       "default",
		Enable:    false,
		Operation: common.UPDATE,
	}
	testSvr.GblCfgCh <- gblCfg
	goToSleep()
}

func TestVrrpV4IntfConfig(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	cfg := testIntfBaseCfg
	cfg.AdminState = false
	cfg.Version = common.VERSION2
	cfg.Operation = common.CREATE
	cfg.IpType = syscall.AF_INET
	_, err := testSvr.ValidConfiguration(&cfg)
	if err == nil {
		t.Error("vrrp configuration should be rejected as there is no l3 intf that got created")
		return
	}
	l3Intf := ipv4Intf
	l3Intf.OperState = common.STATE_UP
	l3Intf.MsgType = common.IP_MSG_CREATE
	testSvr.L3IntfNotifyCh <- l3Intf
	goToSleep()
	if len(testSvr.V4) != 1 {
		t.Error("failed to handle ipv4 interface create notification by create a new v4 entry")
		t.Error("	    ", testSvr.V4)
		return
	}
	testSvr.CfgCh <- &cfg
	goToSleep()
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	vrrpIntf, exists := testSvr.Intf[key]
	if !exists {
		t.Error("failed to init fsm for vrrp v4 version configuration:", cfg)
		t.Error("	    key used for searching vrrp interface:", key)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		return
	}
	if vrrpIntf.Fsm.IsRunning() {
		t.Error("fsm should not be started as VRRP Global is not yet enabled")
		return
	}

	v4State := testSvr.GetEntry(key)
	if v4State == nil {
		t.Error("get vrrp v4 interface by intfRef & VRID failed")
		return
	}
	ip, _, _ := net.ParseCIDR(l3Intf.IpAddr)
	wantStateInfo := common.State{
		IntfRef:                 cfg.IntfRef,
		Vrid:                    cfg.VRID,
		OperState:               common.STATE_DOWN,
		CurrentFsmState:         "Initialize",
		AdverRx:                 0,
		AdverTx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  ip.String(),
		VirtualIp:               cfg.VirtualIPAddr,
		VirtualRouterMACAddress: vrrpIntf.GetVMac(),
		AdvertisementInterval:   cfg.AdvertisementInterval,
	}
	// hacking time stamp to be empty
	v4State.LastAdverRx = ""
	v4State.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, *v4State) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v4State)
		return
	}

	cfg.AdminState = true
	cfg.Operation = common.UPDATE
	testSvr.CfgCh <- &cfg
	goToSleep()
	if !reflect.DeepEqual(cfg, *vrrpIntf.Fsm.Config) {
		t.Error("vrrp update config didn't update fsm configuration")
		t.Error("	    want cfg:", cfg)
		t.Error("	    fsm config:", *vrrpIntf.Fsm.Config)
		return
	}

	if !reflect.DeepEqual(cfg, *vrrpIntf.Config) {
		t.Error("vrrp update config didn't update interface configuration")
		t.Error("	    want cfg:", cfg)
		t.Error("	    intf config:", *vrrpIntf.Config)
		return
	}
	vrrpIntf = testSvr.Intf[key]
	if vrrpIntf.Fsm.IsRunning() {
		t.Error("after updating adming state without vrrp globally enabled fsm should not start")
		t.Error("	    vrrp global information:", *testSvr.GlobalConfig)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		t.Error("	    vrrp fsm information:", *vrrpIntf.Fsm.Config)
		return
	}

	testEnableGlobalConfig()
	vrrpIntf = testSvr.Intf[key]
	if vrrpIntf.Fsm.IsRunning() == false {
		t.Error("after global config enable fsm needs to start for vrrp interfaces")
		t.Error("	    vrrp global information:", *testSvr.GlobalConfig)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		return
	}
	wantStateInfo.OperState = common.STATE_UP
	wantStateInfo.MasterDownTimer = 3
	wantStateInfo.CurrentFsmState = "Backup"
	goToSleep()
	v4State = testSvr.GetEntry(key)
	if v4State == nil {
		t.Error("get vrrp v4 interface by intfRef & VRID failed")
		return
	}
	// hacking time stamp to be empty
	v4State.LastAdverRx = ""
	v4State.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, *v4State) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v4State)
		return
	}
	_, _, v4Bulk := testSvr.GetV4Intfs(0, 100)
	if len(v4Bulk) == 0 && len(v4Bulk) != 1 {
		t.Error("failed to get bulk vrrp v4 entries")
		return
	}

	v4Entry := v4Bulk[0]
	// hacking time stamp to be empty
	v4Entry.LastAdverRx = ""
	v4Entry.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, v4Entry) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v4Entry)
		return
	}

	testDisableGlobalConfig()
	goToSleep()
	goToSleep()
	goToSleep()
	if vrrpIntf.Fsm.IsRunning() {
		t.Error("after global config disable fsm need to stop for vrrp interfaces")
		t.Error("	    vrrp global information:", *testSvr.GlobalConfig)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		t.Error("	    vrrp fsm information:", vrrpIntf.Fsm)
		return
	}

	cfg.Operation = common.DELETE
	testSvr.CfgCh <- &cfg
	goToSleep()

	if len(testSvr.Intf) != 0 {
		t.Error("failed to delete vrrp v4 interface configuration")
		t.Error("	    vrrp global information:", *testSvr.GlobalConfig)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		return
	}
	_, _, v4Bulk = testSvr.GetV4Intfs(0, 100)
	if len(v4Bulk) != 0 {
		t.Error("failed to delete v4 entry after vrrp v4 is deleted")
		return
	}
	TestServerDeInit(t)
}

func TestVrrpV6IntfConfig(t *testing.T) {
	t.Log("=====TestVrrpV6IntfConfig======")
	TestServerInit(t)
	goToSleep()
	ipIntf := ipv6Intf
	ipIntf.MsgType = common.IP_MSG_CREATE
	ipIntf.OperState = common.STATE_UP
	testSvr.L3IntfNotifyCh <- ipIntf
	ip, _, _ := net.ParseCIDR(ipIntf.IpAddr)
	goToSleep()
	if len(testSvr.V6) != 1 {
		t.Error("failed to handle ipv6 interface create notification by create a new v6 entry")
		t.Error("	    ", testSvr.V4)
		return
	}
	cfg := testIntfBaseCfg
	cfg.AdminState = true
	cfg.Version = common.VERSION3
	cfg.Operation = common.CREATE
	cfg.IpType = syscall.AF_INET6
	_, err := testSvr.ValidConfiguration(&cfg)
	if err != nil {
		t.Error("vrrp configuration should not be rejected as there is l3 intf that got created")
		return
	}
	testEnableGlobalConfig()
	testSvr.CfgCh <- &cfg
	goToSleep()
	wantGblState := &common.GlobalState{
		Vrf:           testSvr.GlobalConfig.Vrf,
		Status:        testSvr.GlobalConfig.Enable,
		V4Intfs:       0,
		V6Intfs:       1,
		TotalRxFrames: 0,
		TotalTxFrames: 0,
	}
	state := testSvr.GetGlobalState(testSvr.GlobalConfig.Vrf)
	if !reflect.DeepEqual(state, wantGblState) {
		t.Error("Vrrp Global State Mis-Match after creating vrrp v6 interface")
		t.Error("	    wantGblState:", *wantGblState)
		t.Error("	    rcvdGblState:", *state)
		return
	}

	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	vrrpIntf, exists := testSvr.Intf[key]
	if !exists {
		t.Error("failed to init fsm for vrrp v6 version configuration:", cfg)
		t.Error("	    key used for searching vrrp interface:", key)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		return
	}
	if !vrrpIntf.Fsm.IsRunning() {
		t.Error("fsm should be started as VRRP Global is enabled & admin state is True")
		return
	}

	if !reflect.DeepEqual(cfg, *vrrpIntf.Fsm.Config) {
		t.Error("vrrp config mis-match with fsm configuration")
		t.Error("	    want cfg:", cfg)
		t.Error("	    fsm config:", *vrrpIntf.Fsm.Config)
		return
	}

	if !reflect.DeepEqual(cfg, *vrrpIntf.Config) {
		t.Error("vrrp config mis-match with interface configuration")
		t.Error("	    want cfg:", cfg)
		t.Error("	    intf config:", *vrrpIntf.Config)
		return
	}

	wantStateInfo := common.State{
		IntfRef:                 cfg.IntfRef,
		Vrid:                    cfg.VRID,
		OperState:               common.STATE_UP,
		CurrentFsmState:         "Initialize",
		AdverRx:                 0,
		AdverTx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  ip.String(),
		VirtualIp:               cfg.VirtualIPAddr,
		MasterDownTimer:         3,
		VirtualRouterMACAddress: vrrpIntf.GetVMac(),
		AdvertisementInterval:   cfg.AdvertisementInterval,
	}
	_, _, v6Bulk := testSvr.GetV6Intfs(0, 100)
	if len(v6Bulk) == 0 && len(v6Bulk) != 1 {
		t.Error("failed to get bulk vrrp v4 entries")
		return
	}

	v6Entry := v6Bulk[0]
	// hacking time stamp to be empty
	v6Entry.LastAdverRx = ""
	v6Entry.LastAdverTx = ""
	wantStateInfo.CurrentFsmState = "Backup"
	if !reflect.DeepEqual(wantStateInfo, v6Entry) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v6Entry)
		return
	}
	ipIntf.MsgType = common.IP_MSG_STATE_CHANGE
	ipIntf.OperState = common.STATE_DOWN
	testSvr.L3IntfNotifyCh <- ipIntf
	goToSleep()
	v6State := testSvr.GetEntry(key)
	if v6State == nil {
		t.Error("get vrrp v6 interface by intfRef & VRID failed")
		return
	}
	// hacking time stamp to be empty
	v6State.LastAdverRx = ""
	v6State.LastAdverTx = ""
	wantStateInfo.CurrentFsmState = "Initialize"
	if !reflect.DeepEqual(wantStateInfo, *v6State) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v6State)
		return
	}

	ipIntf.MsgType = common.IP_MSG_STATE_CHANGE
	ipIntf.OperState = common.STATE_UP
	testSvr.L3IntfNotifyCh <- ipIntf
	v6State = testSvr.GetEntry(key)
	if v6State == nil {
		t.Error("get vrrp v6 interface by intfRef & VRID failed")
		return
	}
	// hacking time stamp to be empty
	v6State.LastAdverRx = ""
	v6State.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, *v6State) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", v6State)
		return
	}

	cfg.Operation = common.DELETE
	testSvr.CfgCh <- &cfg
	goToSleep()
	goToSleep()
	goToSleep()

	if len(testSvr.Intf) != 0 {
		t.Error("failed to delete vrrp v6 interface configuration")
		t.Error("	    vrrp global information:", *testSvr.GlobalConfig)
		t.Error("	    vrrp interface information:", testSvr.Intf)
		return
	}
	wantGblState = &common.GlobalState{
		Vrf:           testSvr.GlobalConfig.Vrf,
		Status:        testSvr.GlobalConfig.Enable,
		V4Intfs:       0,
		V6Intfs:       0,
		TotalRxFrames: 0,
		TotalTxFrames: 0,
	}
	state = testSvr.GetGlobalState(testSvr.GlobalConfig.Vrf)
	if !reflect.DeepEqual(state, wantGblState) {
		t.Error("Vrrp Global State Mis-Match During Global Create")
		t.Error("	    wantGblState:", *wantGblState)
		t.Error("	    rcvdGblState:", *state)
		return
	}
	_, _, v6Bulk = testSvr.GetV6Intfs(0, 100)
	if len(v6Bulk) != 0 {
		t.Error("failed to delete v6 entry after vrrp v6 is deleted")
		return
	}
	TestServerDeInit(t)
}

func TestV4V6Intf(t *testing.T) {
	var ipIntf IPIntf
	v4Obj := &V4Intf{}
	key := constructIntfKey("vlan149", 1, syscall.AF_INET)
	ipIntf = v4Obj
	ipIntf.SetVrrpIntfKey(key)
	dbKey := ipIntf.GetVrrpIntfKey()
	if !reflect.DeepEqual(key, *dbKey) {
		t.Error("failed to get vrrp interface from v4 object")
		t.Error("	want key:", key)
		t.Error("	got key:", *dbKey)
		return
	}
	v6Obj := &V6Intf{}
	ipIntf = v6Obj
	ipIntf.SetVrrpIntfKey(key)
	dbKey = ipIntf.GetVrrpIntfKey()
	if !reflect.DeepEqual(key, *dbKey) {
		t.Error("failed to get vrrp interface from v6 object")
		t.Error("	want key:", key)
		t.Error("	got key:", *dbKey)
		return
	}
}

func TestGlobalRxTxCount(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	var empty struct{}
	testSvr.UpdateRxCh <- empty
	testSvr.UpdateTxCh <- empty
	goToSleep()
	if testSvr.globalState.TotalRxFrames != 1 {
		t.Error("failed to update rx count")
		return
	}
	if testSvr.globalState.TotalTxFrames != 1 {
		t.Error("failed to update tx count")
		return
	}
	TestServerDeInit(t)
}

func TestDummyDbRead(t *testing.T) {
	TestServerInit(t)
	goToSleep()
	testSvr.readVrrpGblCfg()
	testSvr.readVrrpV4IntfCfg()
	testSvr.readVrrpV6IntfCfg()
	if testSvr.GlobalConfig.Vrf == "default" {
		t.Error("if no db handle then Global Config should not be done")
		return
	}
	if len(testSvr.V4) != 0 {
		t.Error("no entries are stored in db and hence V4 interfaces count should be 0")
		return
	}
	if len(testSvr.V6) != 0 {
		t.Error("no entries are stored in db and hence V6 interfaces count should be 0")
		return
	}
	TestServerDeInit(t)
}
