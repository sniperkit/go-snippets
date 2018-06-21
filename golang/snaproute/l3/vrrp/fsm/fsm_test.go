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

package fsm

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"log/syslog"
	"reflect"
	"syscall"
	"testing"
	"time"
	"utils/logging"
)

var ipv6Packet = []byte{
	0x33, 0x33, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x5e, 0x00, 0x02, 0x01, 0x86, 0xdd, 0x60, 0x00,
	0x00, 0x00, 0x00, 0x18, 0x70, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x1d,
	0xfc, 0xff, 0xfe, 0xcf, 0x15, 0xfc, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x31, 0x01, 0x64, 0x01, 0x00, 0x64, 0xd1, 0x9d, 0xfe, 0x80,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
}
var testFsm *FSM
var testv6IntfCfg = &common.IntfCfg{
	IntfRef:               "lo",
	IfIndex:               0,
	VRID:                  1,
	Priority:              150,
	VirtualIPAddr:         "fe80::172:16:0:1/64",
	AdvertisementInterval: 3,
	PreemptMode:           true,
	AcceptMode:            false,
	AdminState:            true,
	Version:               common.VERSION3,
	Operation:             1,
	IpType:                syscall.AF_INET6,
}
var testL3Info = &common.BaseIpInfo{33554581, "lo", "fe80::c837:abff:febe:ad4", "UP", "", syscall.AF_INET6}
var testVipCh = make(chan *common.VirtualIpInfo)
var testRxCh = make(chan struct{})
var testTxCh = make(chan struct{})
var testExitMimicServerCh = make(chan struct{})
var testV6Vmac = "00-00-5E-00-02-01"
var testV4Vmac = "00-00-5E-00-01-01"
var wantStateInfo common.State

func mimicServer(t *testing.T) {
	for {
		select {
		case _, ok := <-testVipCh:
			if ok {
				//t.Log("Server received virtual ip update information:", *vipInfo)
			}
		case _, ok := <-testRxCh:
			if ok {
				//t.Log("Update rx count")
			}
		case _, ok := <-testTxCh:
			if ok {
				wantStateInfo.AdverTx++
			}
		case _, ok := <-testExitMimicServerCh:
			if ok {
				return
			}
		}
	}
}

func TestV6FsmInit(t *testing.T) {
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

	testFsm = InitFsm(testv6IntfCfg, testL3Info, testVipCh, testRxCh, testTxCh)
	if testFsm == nil {
		t.Error("Initializing FSM Failed")
		return
	}

	go mimicServer(t)
	time.Sleep(50 * time.Millisecond)
}

func testFsmDeInit(t *testing.T) {
	testExitMimicServerCh <- testFsm.empty
	time.Sleep(50 * time.Millisecond)
	testFsm.DeInitFsm()
	if testFsm.Config != nil {
		t.Error("failed to delete interface config")
		return
	}
	if testFsm.AdverTimer != nil {
		t.Error("failed to stop & delete advertisement timer")
		return
	}

	if testFsm.vipCh == nil {
		t.Error("return virtual ip interface channel for server shouldn't be deleted")
		return
	}

	if testFsm.stateInfo != nil {
		t.Error("Failed to delete FSM state information")
		return
	}

	if testFsm.pktCh != nil {
		t.Error("failed to delete packet channel")
		return
	}

	if testFsm.pHandle != nil {
		t.Error("failed to delte pcap handler for fsm")
	}

	if testFsm.PktInfo != nil {
		t.Error("failed to delete object for packet encoding & decoding called packetInfo")
		return
	}

	if testFsm.MasterDownTimer != nil {
		t.Error("failed to stop & delete master down timer")
		return
	}

	if testFsm.IntfEventCh != nil {
		t.Error("failed to delete interface event channel handler")
		return
	}
}

func TestV6UpdateIntfConfig(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testv6IntfCfg.Priority = 125
	testFsm.UpdateConfig(testv6IntfCfg)
	if !reflect.DeepEqual(testv6IntfCfg, testFsm.Config) {
		t.Error("Failed to update interface config")
		t.Error("	    Wanted Config:", *testv6IntfCfg)
		t.Error("	    Received Config:", *testFsm.Config)
		return
	}
	testFsmDeInit(t)
}

func TestFsmIsRunning(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	go testFsm.StartFsm()
	time.Sleep(50 * time.Millisecond)
	if testFsm.IsRunning() == false {
		t.Error("FSM started but is running flag is set to false")
		return
	}
	testFsmDeInit(t)
}

func TestGetStateInfo(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	wantStateInfo := common.State{
		IntfRef:                 testv6IntfCfg.IntfRef,
		Vrid:                    testv6IntfCfg.VRID,
		OperState:               common.STATE_DOWN,
		CurrentFsmState:         VRRP_INITIALIZE_STATE_STRING,
		MasterIp:                "",
		AdverRx:                 0,
		AdverTx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  testL3Info.IpAddr,
		VirtualIp:               testv6IntfCfg.VirtualIPAddr,
		VirtualRouterMACAddress: testV6Vmac,
		AdvertisementInterval:   testv6IntfCfg.AdvertisementInterval,
		MasterDownTimer:         0,
	}
	state := common.State{}
	testFsm.GetStateInfo(&state)
	if !reflect.DeepEqual(wantStateInfo, state) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", state)
		return
	}
	testFsmDeInit(t)
}

func TestFsmUpGetStateInfo(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	go testFsm.StartFsm()
	time.Sleep(10 * time.Millisecond)
	wantStateInfo := common.State{
		IntfRef:                 testv6IntfCfg.IntfRef,
		Vrid:                    testv6IntfCfg.VRID,
		OperState:               common.STATE_UP,
		CurrentFsmState:         VRRP_INITIALIZE_STATE_STRING,
		MasterIp:                "",
		AdverRx:                 0,
		AdverTx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  testL3Info.IpAddr,
		VirtualIp:               testv6IntfCfg.VirtualIPAddr,
		VirtualRouterMACAddress: testV6Vmac,
		AdvertisementInterval:   testv6IntfCfg.AdvertisementInterval,
		MasterDownTimer:         0,
	}
	state := common.State{}
	testFsm.GetStateInfo(&state)
	if !reflect.DeepEqual(wantStateInfo, state) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", state)
		return
	}
	testFsmDeInit(t)
}

func TestV6VMacCreate(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testFsm.createVirtualMac()
	if testFsm.VirtualMACAddress != testV6Vmac {
		t.Error("Failed creating virtual mac for v6 interface")
		t.Error("	    wanted vmac:", testV6Vmac)
		t.Error("	    received vmac:", testFsm.VirtualMACAddress)
		return
	}
	testFsmDeInit(t)
}

func TestV4VMacCreate(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testFsm.ipType = syscall.AF_INET
	testFsm.createVirtualMac()
	if testFsm.VirtualMACAddress != testV4Vmac {
		t.Error("Failed creating virtual mac for v4 interface")
		t.Error("	    wanted vmac:", testV4Vmac)
		t.Error("	    received vmac:", testFsm.VirtualMACAddress)
		return
	}
	testFsmDeInit(t)
}

func TestInitPktListener(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testFsm.initPktListener()
	if testFsm.pHandle == nil {
		t.Error("Failed creating pkt listener")
		return
	}
	testFsmDeInit(t)
}

func TestProcessRcvdPktStateToBackup(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	go testFsm.StartFsm()
	time.Sleep(500 * time.Millisecond)
	p := gopacket.NewPacket(ipv6Packet, layers.LinkTypeEthernet, gopacket.Default)
	pktInfo := &PktChannelInfo{p}
	testFsm.processRcvdPkt(pktInfo)

	wantStateInfo := common.State{
		IntfRef:                 testv6IntfCfg.IntfRef,
		Vrid:                    testv6IntfCfg.VRID,
		OperState:               common.STATE_UP,
		CurrentFsmState:         VRRP_BACKUP_STATE_STRING,
		MasterIp:                "fe80::8a1d:fcff:fecf:15fc",
		AdverRx:                 1,
		AdverTx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  testL3Info.IpAddr,
		VirtualIp:               testv6IntfCfg.VirtualIPAddr,
		VirtualRouterMACAddress: testV6Vmac,
		AdvertisementInterval:   testv6IntfCfg.AdvertisementInterval,
		MasterDownTimer:         10,
	}
	state := common.State{}
	testFsm.GetStateInfo(&state)
	// hacking time stamp to be empty
	state.LastAdverRx = ""
	state.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, state) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", state)
		return
	}
	testFsmDeInit(t)
}

func TestProcessRcvdPktStateToMaster(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testv6IntfCfg.Priority = VRRP_MASTER_PRIORITY
	testv6IntfCfg.AdvertisementInterval = 100
	testv6IntfCfg.AdminState = true
	testFsm.UpdateConfig(testv6IntfCfg)
	if !reflect.DeepEqual(testv6IntfCfg, testFsm.Config) {
		t.Error("Failed to update interface config")
		t.Error("	    Wanted Config:", *testv6IntfCfg)
		t.Error("	    Received Config:", *testFsm.Config)
		return
	}
	p := gopacket.NewPacket(ipv6Packet, layers.LinkTypeEthernet, gopacket.Default)
	pktInfo := &PktChannelInfo{p}
	testFsm.processRcvdPkt(pktInfo)
	wantStateInfo = common.State{
		IntfRef:                 testv6IntfCfg.IntfRef,
		Vrid:                    testv6IntfCfg.VRID,
		OperState:               common.STATE_DOWN, // it should not be up
		CurrentFsmState:         VRRP_MASTER_STATE_STRING,
		MasterIp:                testL3Info.IpAddr, // master state and hence my own ip
		AdverRx:                 0,
		LastAdverRx:             "",
		LastAdverTx:             "",
		IpAddr:                  testL3Info.IpAddr,
		VirtualIp:               testv6IntfCfg.VirtualIPAddr,
		VirtualRouterMACAddress: testV6Vmac,
		AdvertisementInterval:   testv6IntfCfg.AdvertisementInterval,
		MasterDownTimer:         0,
	}
	time.Sleep(50 * time.Millisecond)
	state := common.State{}
	testFsm.GetStateInfo(&state)
	// hacking time stamp to be empty
	state.LastAdverRx = ""
	state.LastAdverTx = ""
	if !reflect.DeepEqual(wantStateInfo, state) {
		t.Error("Failure getting state information from fsm")
		t.Error("	    want state info:", wantStateInfo)
		t.Error("	    got state info:", state)
		return
	}
	testFsmDeInit(t)
}

func TestMasterEqualPriorityState(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testFsm.transitionToMaster()
	testv6IntfCfg.Priority = 100
	testv6IntfCfg.AdvertisementInterval = 1
	testFsm.UpdateConfig(testv6IntfCfg)
	p := gopacket.NewPacket(ipv6Packet, layers.LinkTypeEthernet, gopacket.Default)
	pktInfo := testFsm.PktInfo.Decode(p)
	decodeInfo := &DecodedInfo{PktInfo: pktInfo}
	testFsm.master(decodeInfo)

	time.Sleep(50 * time.Millisecond)
	state := common.State{}
	testFsm.GetStateInfo(&state)
	if state.CurrentFsmState != VRRP_MASTER_STATE_STRING {
		t.Error("For equal priority and sender ip < local ip, state should be master")
		return
	}
	// changing ip address to make sure that we transition to backup state
	testFsm.ipAddr = "fe80::8a1d:fcff:fecf:15fb"
	testFsm.master(decodeInfo)
	time.Sleep(50 * time.Millisecond)
	testFsm.GetStateInfo(&state)
	if state.CurrentFsmState != VRRP_BACKUP_STATE_STRING {
		t.Error("For equal priority and sender ip > local ip, state should be backup")
		return
	}
	testFsmDeInit(t)
}

func TestBackupState(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	testFsm.transitionToBackup(testv6IntfCfg.AdvertisementInterval)
	p := gopacket.NewPacket(ipv6Packet, layers.LinkTypeEthernet, gopacket.Default)
	pktInfo := testFsm.PktInfo.Decode(p)
	decodeInfo := &DecodedInfo{PktInfo: pktInfo}
	testFsm.backup(decodeInfo)
	state := common.State{}
	testFsm.GetStateInfo(&state)
	if state.CurrentFsmState != VRRP_BACKUP_STATE_STRING {
		t.Error("local priority is higher than the priority in pkt received and hence fsm should be in backup state")
		return
	}
	testFsmDeInit(t)
}

func TestMasterDownTimerAndMasterAdverTimer(t *testing.T) {
	TestV6FsmInit(t)
	if testFsm == nil {
		return
	}
	// handle master down time to 1 second
	testFsm.MasterDownValue = 1
	testFsm.handleMasterDownTimer()
	if testFsm.MasterDownTimer == nil {
		t.Error("Failed to start master down timer")
		return
	}
	// sleep for down timer and then get state object
	time.Sleep(time.Duration(testFsm.MasterDownValue) * time.Second)
	time.Sleep(50 * time.Millisecond)
	state := common.State{}
	testFsm.GetStateInfo(&state)
	if state.CurrentFsmState != VRRP_MASTER_STATE_STRING {
		t.Error("on master down timer expiry we should have transition to master")
		return
	}
	if testFsm.AdverTimer == nil {
		t.Error("failed to start AdvertisementInterval after transitioning to master")
		return
	}
	testFsmDeInit(t)
}
