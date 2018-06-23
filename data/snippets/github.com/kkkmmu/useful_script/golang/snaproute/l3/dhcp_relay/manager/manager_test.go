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

package manager

import (
	"dhcprelayd"
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"l3/dhcp_relay/infra"
	"log/syslog"
	"models/events"
	"models/objects"
	"net"
	"reflect"
	"testing"
	"time"
	asicdmock "utils/asicdClient/mock"
	"utils/commonDefs"
	"utils/logging"
)

type Processor4Fake struct {
	EnabledFlag bool

	ClientStateSlice     []*dhcprelayd.DHCPRelayClientState
	ClientStateMap       map[string]*dhcprelayd.DHCPRelayClientState
	IntfStateSlice       []*dhcprelayd.DHCPRelayIntfState
	IntfStateMap         map[string]*dhcprelayd.DHCPRelayIntfState
	IntfServerStateSlice []*dhcprelayd.DHCPRelayIntfServerState
	IntfServerStateMap   map[string]*dhcprelayd.DHCPRelayIntfServerState
}

type Processor6Fake struct {
	EnabledFlag bool

	ClientStateSlice     []*dhcprelayd.DHCPv6RelayClientState
	ClientStateMap       map[string]*dhcprelayd.DHCPv6RelayClientState
	IntfStateSlice       []*dhcprelayd.DHCPv6RelayIntfState
	IntfStateMap         map[string]*dhcprelayd.DHCPv6RelayIntfState
	IntfServerStateSlice []*dhcprelayd.DHCPv6RelayIntfServerState
	IntfServerStateMap   map[string]*dhcprelayd.DHCPv6RelayIntfServerState
}

func (pProc *Processor4Fake) GetEnabledFlag() bool {
	return pProc.EnabledFlag
}

func (pProc *Processor4Fake) SetEnabledFlag() {
	pProc.EnabledFlag = true
}

func (pProc *Processor4Fake) StartRxTx() {
	pProc.EnabledFlag = true
}

func (pProc *Processor4Fake) StopRxTx() {
	pProc.EnabledFlag = false
}

func (pProc *Processor4Fake) GetClientStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayClientState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.ClientStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPRelayClientState, actualCount)
	copy(result, pProc.ClientStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor4Fake) GetClientState(
	macAddr string) (*dhcprelayd.DHCPRelayClientState, bool) {

	if clientState, ok := pProc.ClientStateMap[macAddr]; ok {
		return clientState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor4Fake) GetIntfStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPRelayIntfState, actualCount)
	copy(result, pProc.IntfStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor4Fake) GetIntfState(
	ifName string) (*dhcprelayd.DHCPRelayIntfState, bool) {

	if intfState, ok := pProc.IntfStateMap[ifName]; ok {
		return intfState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor4Fake) GetIntfServerStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfServerState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfServerStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPRelayIntfServerState, actualCount)
	copy(result, pProc.IntfServerStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor4Fake) GetIntfServerState(
	ifName string, serverAddr string) (*dhcprelayd.DHCPRelayIntfServerState, bool) {

	intfServerStateKey := ifName + "_" + serverAddr
	if intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]; ok {
		return intfServerState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor4Fake) ProcessCreateDRAIntf(ifName string) {
	return
}

func (pProc *Processor4Fake) ProcessDeleteDRAIntf(ifName string) {
	return
}

func (pProc *Processor4Fake) ProcessActiveDRAIntf(ifIdx int) {
	return
}

func (pProc *Processor4Fake) ProcessInactiveDRAIntf(ifIdx int) {
	return
}

func (pProc *Processor6Fake) GetEnabledFlag() bool {
	return pProc.EnabledFlag
}

func (pProc *Processor6Fake) SetEnabledFlag() {
	pProc.EnabledFlag = true
}

func (pProc *Processor6Fake) StartRxTx() {
	pProc.EnabledFlag = true
}

func (pProc *Processor6Fake) StopRxTx() {
	pProc.EnabledFlag = false
}

func (pProc *Processor6Fake) GetClientStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayClientState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.ClientStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayClientState, actualCount)
	copy(result, pProc.ClientStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor6Fake) GetClientState(
	macAddr string) (*dhcprelayd.DHCPv6RelayClientState, bool) {

	if clientState, ok := pProc.ClientStateMap[macAddr]; ok {
		return clientState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor6Fake) GetIntfStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayIntfState, actualCount)
	copy(result, pProc.IntfStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor6Fake) GetIntfState(
	ifName string) (*dhcprelayd.DHCPv6RelayIntfState, bool) {

	if intfState, ok := pProc.IntfStateMap[ifName]; ok {
		return intfState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor6Fake) GetIntfServerStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfServerState) {

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfServerStateSlice)
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayIntfServerState, actualCount)
	copy(result, pProc.IntfServerStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor6Fake) GetIntfServerState(
	ifName string, serverAddr string) (*dhcprelayd.DHCPv6RelayIntfServerState, bool) {

	intfServerStateKey := ifName + "_" + serverAddr
	if intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]; ok {
		return intfServerState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor6Fake) ProcessCreateDRAIntf(ifName string) {
	return
}

func (pProc *Processor6Fake) ProcessDeleteDRAIntf(ifName string) {
	return
}

func (pProc *Processor6Fake) ProcessActiveDRAIntf(ifIdx int) {
	return
}

func (pProc *Processor6Fake) ProcessInactiveDRAIntf(ifIdx int) {
	return
}

type DbFake struct {
}

func (dbHdl DbFake) Connect() error {
	return nil
}
func (dbHdl DbFake) Disconnect() {}
func (dbHdl DbFake) StoreObjectInDb(objects.ConfigObj) error {
	return nil
}
func (dbHdl DbFake) DeleteObjectFromDb(objects.ConfigObj) error {
	return nil
}
func (dbHdl DbFake) GetObjectFromDb(objects.ConfigObj, string) (objects.ConfigObj, error) {
	return nil, nil
}
func (dbHdl DbFake) GetKey(objects.ConfigObj) string {
	return ""
}
func (dbHdl DbFake) GetAllObjFromDb(obj objects.ConfigObj) ([]objects.ConfigObj, error) {
	resultObjs := []objects.ConfigObj{}
	switch reflect.TypeOf(obj) {
	case reflect.TypeOf(objects.DHCPRelayIntf{}):
		resultObjs = append(resultObjs, objects.DHCPRelayIntf{
			IntfRef:  "eth0",
			Enable:   true,
			ServerIp: []string{"11.0.0.1"},
		})
		resultObjs = append(resultObjs, objects.DHCPRelayIntf{
			IntfRef:  "eth1",
			Enable:   true,
			ServerIp: []string{"12.0.0.1"},
		})
		resultObjs = append(resultObjs, objects.DHCPRelayIntf{
			IntfRef:  "eth2",
			Enable:   true,
			ServerIp: []string{"13.0.0.1"},
		})
	case reflect.TypeOf(objects.DHCPRelayGlobal{}):
		resultObjs = append(resultObjs, objects.DHCPRelayGlobal{
			Vrf:           "default",
			Enable:        true,
			HopCountLimit: 32,
		})
	case reflect.TypeOf(objects.DHCPv6RelayIntf{}):
		resultObjs = append(resultObjs, objects.DHCPv6RelayIntf{
			IntfRef:       "eth0",
			Enable:        true,
			ServerIp:      []string{"2001::1"},
			UpstreamIntfs: []string{},
		})
		resultObjs = append(resultObjs, objects.DHCPv6RelayIntf{
			IntfRef:       "eth1",
			Enable:        true,
			ServerIp:      []string{"2001::2"},
			UpstreamIntfs: []string{},
		})
		resultObjs = append(resultObjs, objects.DHCPv6RelayIntf{
			IntfRef:       "eth2",
			Enable:        true,
			ServerIp:      []string{"2001::3"},
			UpstreamIntfs: []string{},
		})
	case reflect.TypeOf(objects.DHCPv6RelayGlobal{}):
		resultObjs = append(resultObjs, objects.DHCPv6RelayGlobal{
			Vrf:           "default",
			Enable:        true,
			HopCountLimit: 32,
		})
	}
	return resultObjs, nil
}
func (dbHdl DbFake) CompareObjectsAndDiff(objects.ConfigObj, map[string]bool, objects.ConfigObj) ([]bool, error) {
	return []bool{}, nil
}

func (dbHdl DbFake) CompareObjectDefaultAndDiff(objects.ConfigObj, objects.ConfigObj) ([]bool, error) {
	return nil, nil
}

func (dbHdl DbFake) StoreObjectDefaultInDb(objects.ConfigObj) error {
	return nil
}

func (dbHdl DbFake) UpdateObjectInDb(objects.ConfigObj, objects.ConfigObj, []bool) error {
	return nil
}
func (dbHdl DbFake) MergeDbAndConfigObj(objects.ConfigObj, objects.ConfigObj, []bool) (objects.ConfigObj, error) {
	return nil, nil
}
func (dbHdl DbFake) GetBulkObjFromDb(obj objects.ConfigObj, startIndex, count int64) (error, int64, int64, bool, []objects.ConfigObj) {
	return nil, 0, 0, true, nil
}
func (dbHdl DbFake) Publish(string, interface{}, interface{}) {}
func (dbHdl DbFake) StoreValInDb(interface{}, interface{}, interface{}) error {
	return nil
}
func (dbHdl DbFake) DeleteValFromDb(interface{}) error {
	return nil
}
func (dbHdl DbFake) GetAllKeys(interface{}) (interface{}, error) {
	return nil, nil
}
func (dbHdl DbFake) GetValFromDB(key interface{}, field interface{}) (val interface{}, err error) {
	return nil, nil
}
func (dbHdl DbFake) StoreEventObjectInDb(events.EventObj) error {
	return nil
}
func (dbHdl DbFake) GetEventObjectFromDb(events.EventObj, string) (events.EventObj, error) {
	return nil, nil
}
func (dbHdl DbFake) GetAllEventObjFromDb(events.EventObj) ([]events.EventObj, error) {
	return []events.EventObj{}, nil
}
func (dbHdl DbFake) MergeDbAndConfigObjForPatchUpdate(objects.ConfigObj, objects.ConfigObj, []objects.PatchOpInfo) (objects.ConfigObj, []bool, error) {
	return nil, []bool{}, nil
}
func (dbHdl DbFake) StoreUUIDToObjKeyMap(objKey string) (string, error) {
	return "", nil
}
func (dbHdl DbFake) DeleteUUIDToObjKeyMap(uuid, objKey string) error {
	return nil
}
func (dbHdl DbFake) GetUUIDFromObjKey(objKey string) (string, error) {
	return "", nil
}
func (dbHdl DbFake) GetObjKeyFromUUID(uuid string) (string, error) {
	return "", nil
}
func (dbHdl DbFake) MergeDbObjKeys(obj, dbObj objects.ConfigObj) (objects.ConfigObj, error) {
	return nil, nil
}

func NewLogger(name string, tag string, listenToConfig bool) (*logging.Writer, error) {
	var err error
	srLogger := new(logging.Writer)
	srLogger.MyComponentName = name

	srLogger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}

	// srLogger.GlobalLogging = true
	srLogger.MyLogLevel = sysdCommonDefs.INFO
	return srLogger, err
}

func InitTestDRAMgr() (*DRAMgr, error) {
	lgr, err := NewLogger("dhcprelayd", "dhcprelayd", true)
	if err != nil {
		return nil, err
	}
	iMgr := infra.NewInfraMgr(lgr, &asicdmock.MockAsicdClientMgr{})
	pProc4F := &Processor4Fake{}
	pProc6F := &Processor6Fake{}
	draMgr := &DRAMgr{
		DbHdl:  DbFake{},
		Logger: lgr,
		IMgr:   iMgr,
		PProc4: pProc4F,
		PProc6: pProc6F,
	}
	return draMgr, nil
}

// drav4
func TestCreateDRAv4Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.2", "10.0.0.3"},
	}
	_, err = draMgr.CreateDRAv4Interface(drav4IntfCfg)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	if !draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	t.Log("PASS: CreateDRAv4Interface")
}

func TestCreateDRAv4Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	ipv4Intf := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	_, err = draMgr.CreateDRAv4Interface(drav4IntfCfg)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	t.Log("PASS: CreateDRAv4Interface")
}

func TestCreateDRAv4Intf3(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	ipv4Intf := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.1.1.3.8"},
	}
	_, err = draMgr.CreateDRAv4Interface(drav4IntfCfg)
	if err == nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	t.Log("PASS: CreateDRAv4Interface")
}

func TestCreateDRAv4Intf4(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1"},
	}
	_, err = draMgr.CreateDRAv4Interface(drav4IntfCfg)
	if err == nil {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Interface")
		return
	}
	t.Log("PASS: CreateDRAv4Interface")
}

func TestUpdateDRAv4Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	ipv4Intf := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[1] = draIntfCfgPre
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc4.SetEnabledFlag()
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   false,
		ServerIp: []string{"10.0.0.1"},
	}
	_, err = draMgr.UpdateDRAv4Interface(
		draIntfCfgPre, drav4IntfCfg, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	draIntfCfgPost := draMgr.IMgr.DRAv4Intfs["eth0"]
	expectedIntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   false,
		ServerIp: []string{"10.0.0.1"},
	}
	if draIntfCfgPost.Enable != expectedIntfCfg.Enable {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	if len(draIntfCfgPost.ServerIp) != len(expectedIntfCfg.ServerIp) {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	for i, ipAddr := range draIntfCfgPost.ServerIp {
		if ipAddr != expectedIntfCfg.ServerIp[i] {
			t.Errorf("FAIL: UpdateDRAv4Interface")
			return
		}
	}
	t.Log("PASS: UpdateDRAv4Interface")
}

func TestUpdateDRAv4Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	ipv4Intf := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav4IntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{},
	}
	draMgr.PProc4.SetEnabledFlag()
	_, err = draMgr.UpdateDRAv4Interface(
		draIntfCfgPre, drav4IntfCfg, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	if !draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	draIntfCfgPost := draMgr.IMgr.DRAv4Intfs["eth0"]
	expectedIntfCfg := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{},
	}
	if draIntfCfgPost.Enable != expectedIntfCfg.Enable {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	if len(draIntfCfgPost.ServerIp) != len(expectedIntfCfg.ServerIp) {
		t.Errorf("FAIL: UpdateDRAv4Interface")
		return
	}
	for i, ipAddr := range draIntfCfgPost.ServerIp {
		if ipAddr != expectedIntfCfg.ServerIp[i] {
			t.Errorf("FAIL: UpdateDRAv4Interface")
			return
		}
	}
	t.Log("PASS: UpdateDRAv4Interface")
}

func TestDeleteDRAv4Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	ipv4Intf := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[1] = draIntfCfgPre
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc4.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv4Interface("eth0")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	if _, ok := draMgr.IMgr.DRAv4Intfs["eth0"]; ok {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	t.Log("PASS: DeleteDRAv4Interface")
}

func TestDeleteDRAv4Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv4Intf2 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IntfProps[2] = ipv4Intf2
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv4Intfs["eth1"] = draIntf2CfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[2] = draIntf2CfgPre
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc4.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv4Interface("eth0")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	if !draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	if _, ok := draMgr.IMgr.DRAv4Intfs["eth0"]; ok {
		t.Errorf("FAIL: DeleteDRAv4Interface")
		return
	}
	t.Log("PASS: DeleteDRAv4Interface")
}

func TestCreateDRAv4Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draIntf1CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntf1CfgPre
	drav4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.CreateDRAv4Global(drav4Global)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	if !draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	if draMgr.IMgr.DRAv4Global == nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	t.Log("PASS: CreateDRAv4Global")
}

func TestCreateDRAv4Global2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	drav4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.CreateDRAv4Global(drav4Global)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	if draMgr.IMgr.DRAv4Global == nil {
		t.Errorf("FAIL: CreateDRAv4Global")
		return
	}
	t.Log("PASS: CreateDRAv4Global")
}

func TestUpdateDRAv4Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv4Intf2 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.2",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IntfProps[2] = ipv4Intf2
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv4Intfs["eth1"] = draIntf2CfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[2] = draIntf2CfgPre
	oldDraV4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.IMgr.DRAv4Global = oldDraV4Global
	draMgr.PProc4.SetEnabledFlag()
	newDraV4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        false,
		HopCountLimit: 31,
	}
	_, err = draMgr.UpdateDRAv4Global(
		oldDraV4Global, newDraV4Global, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if draMgr.IMgr.GetActiveDRAv4IntfCount() > 0 {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	curDraVgGlobal := draMgr.IMgr.DRAv4Global
	if curDraVgGlobal.Enable != newDraV4Global.Enable {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if curDraVgGlobal.HopCountLimit != newDraV4Global.HopCountLimit {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	t.Log("PASS: UpdateDRAv4Global")
}

func TestUpdateDRAv4Global2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv4Intf2 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.2",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IntfProps[2] = ipv4Intf2
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv4Intfs["eth1"] = draIntf2CfgPre
	oldDraV4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        false,
		HopCountLimit: 32,
	}
	draMgr.IMgr.DRAv4Global = oldDraV4Global
	draMgr.PProc4.SetEnabledFlag()
	newDraV4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.UpdateDRAv4Global(
		oldDraV4Global, newDraV4Global, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if !draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if draMgr.IMgr.GetActiveDRAv4IntfCount() <= 0 {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	curDraVgGlobal := draMgr.IMgr.DRAv4Global
	if curDraVgGlobal.Enable != newDraV4Global.Enable {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	if curDraVgGlobal.HopCountLimit != newDraV4Global.HopCountLimit {
		t.Errorf("FAIL: UpdateDRAv4Global")
		return
	}
	t.Log("PASS: UpdateDRAv4Global")
}

func TestDeleteDRAv4Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Global")
		return
	}
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv4IntfProps[1] = ipv4Intf1
	draMgr.IMgr.IPv4IfRefToIfIndex["eth0"] = 1
	draIntf1CfgPre := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1", "10.0.0.2"},
	}
	draMgr.IMgr.DRAv4Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv4Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc4.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv4Global("default")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv4Global")
		return
	}
	if draMgr.PProc4.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv4Global")
		return
	}
	if draMgr.IMgr.DRAv4Global != nil {
		t.Errorf("FAIL: DeleteDRAv4Global")
		return
	}
	t.Log("PASS: DeleteDRAv4Global")
}

//States
func GetV4DummyStates() (
	[]*dhcprelayd.DHCPRelayClientState,
	map[string]*dhcprelayd.DHCPRelayClientState,
	[]*dhcprelayd.DHCPRelayIntfState,
	map[string]*dhcprelayd.DHCPRelayIntfState,
	[]*dhcprelayd.DHCPRelayIntfServerState,
	map[string]*dhcprelayd.DHCPRelayIntfServerState) {

	clientStates := []*dhcprelayd.DHCPRelayClientState{}
	clientStates = append(clientStates, &dhcprelayd.DHCPRelayClientState{
		MacAddr:         "aa:bb:cc:dd:ee:ff",
		ServerIp:        "10.0.0.1",
		OfferedIp:       "10.3.4.5",
		GatewayIp:       "11.0.0.1",
		AcceptedIp:      "10.3.4.5",
		RequestedIp:     "10.3.4.5",
		ClientDiscover:  time.Now().String(),
		ClientRequest:   time.Now().String(),
		ClientRequests:  1124,
		ClientResponses: 1345,
		ServerOffer:     time.Now().String(),
		ServerAck:       time.Now().String(),
		ServerRequests:  932,
		ServerResponses: 878,
	},
	)
	clientStates = append(clientStates, &dhcprelayd.DHCPRelayClientState{
		MacAddr:         "bb:cc:dd:ee:ff:aa",
		ServerIp:        "20.0.0.1",
		OfferedIp:       "10.4.5.6",
		GatewayIp:       "12.0.0.1",
		AcceptedIp:      "10.4.5.6",
		RequestedIp:     "10.4.5.6",
		ClientDiscover:  time.Now().String(),
		ClientRequest:   time.Now().String(),
		ClientRequests:  1224,
		ClientResponses: 1445,
		ServerOffer:     time.Now().String(),
		ServerAck:       time.Now().String(),
		ServerRequests:  942,
		ServerResponses: 888,
	},
	)
	clientStates = append(clientStates, &dhcprelayd.DHCPRelayClientState{
		MacAddr:         "cc:dd:ee:ff:aa:bb",
		ServerIp:        "30.0.0.1",
		OfferedIp:       "10.5.6.7",
		GatewayIp:       "13.0.0.1",
		AcceptedIp:      "10.5.6.7",
		RequestedIp:     "10.5.6.7",
		ClientDiscover:  time.Now().String(),
		ClientRequest:   time.Now().String(),
		ClientRequests:  1324,
		ClientResponses: 1545,
		ServerOffer:     time.Now().String(),
		ServerAck:       time.Now().String(),
		ServerRequests:  952,
		ServerResponses: 898,
	},
	)
	intfStates := []*dhcprelayd.DHCPRelayIntfState{}
	intfStates = append(intfStates, &dhcprelayd.DHCPRelayIntfState{
		IntfRef:           "eth0",
		TotalDrops:        1123,
		TotalDhcpClientRx: 34567,
		TotalDhcpClientTx: 33457,
		TotalDhcpServerRx: 32567,
		TotalDhcpServerTx: 34222,
	},
	)
	intfStates = append(intfStates, &dhcprelayd.DHCPRelayIntfState{
		IntfRef:           "eth1",
		TotalDrops:        1223,
		TotalDhcpClientRx: 35567,
		TotalDhcpClientTx: 34457,
		TotalDhcpServerRx: 33567,
		TotalDhcpServerTx: 35222,
	},
	)
	intfStates = append(intfStates, &dhcprelayd.DHCPRelayIntfState{
		IntfRef:           "eth2",
		TotalDrops:        1323,
		TotalDhcpClientRx: 36567,
		TotalDhcpClientTx: 35457,
		TotalDhcpServerRx: 34567,
		TotalDhcpServerTx: 36222,
	},
	)
	intfServerStates := []*dhcprelayd.DHCPRelayIntfServerState{}
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPRelayIntfServerState{
		IntfRef:   "eth0",
		ServerIp:  "10.0.0.1",
		Request:   34567,
		Responses: 33457,
	},
	)
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPRelayIntfServerState{
		IntfRef:   "eth1",
		ServerIp:  "20.0.0.1",
		Request:   35567,
		Responses: 34457,
	},
	)
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPRelayIntfServerState{
		IntfRef:   "eth2",
		ServerIp:  "30.0.0.1",
		Request:   36567,
		Responses: 35457,
	},
	)
	intfRefMap := make(map[string]int)
	intfRefMap["eth0"] = 1
	intfRefMap["eth1"] = 2
	intfRefMap["eth2"] = 3

	clientStateMap := make(map[string]*dhcprelayd.DHCPRelayClientState)
	intfStateMap := make(map[string]*dhcprelayd.DHCPRelayIntfState)
	intfServerStateMap := make(map[string]*dhcprelayd.DHCPRelayIntfServerState)
	for _, clientState := range clientStates {
		clientStateMap[clientState.MacAddr] = clientState
	}
	for _, intfState := range intfStates {
		intfStateMap[intfState.IntfRef] = intfState
	}
	for _, intfServerState := range intfServerStates {
		intfServerStateKey := intfServerState.IntfRef + "_" + intfServerState.ServerIp
		intfServerStateMap[intfServerStateKey] = intfServerState
	}

	return clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap
}

func GetV4DummyStateProcessor(
	clientStates []*dhcprelayd.DHCPRelayClientState,
	clientStateMap map[string]*dhcprelayd.DHCPRelayClientState,
	intfStates []*dhcprelayd.DHCPRelayIntfState,
	intfStateMap map[string]*dhcprelayd.DHCPRelayIntfState,
	intfServerStates []*dhcprelayd.DHCPRelayIntfServerState,
	intfServerStateMap map[string]*dhcprelayd.DHCPRelayIntfServerState) *Processor4Fake {

	pProcF := &Processor4Fake{}
	pProcF.ClientStateSlice = clientStates
	pProcF.ClientStateMap = clientStateMap
	pProcF.IntfStateSlice = intfStates
	pProcF.IntfStateMap = intfStateMap
	pProcF.IntfServerStateSlice = intfServerStates
	pProcF.IntfServerStateMap = intfServerStateMap
	return pProcF
}

func PopulateV4DummyInfra(iMgr *infra.InfraMgr) {
	iMgr.IPv4IfRefToIfIndex = make(map[string]int)
	iMgr.IPv4IntfProps = make(map[int]*infra.IPv4IntfProperty)
	iMgr.DRAv4Intfs = make(map[string]*dhcprelayd.DHCPRelayIntf)
	iMgr.ActiveDRAv4Intfs = make(map[int]*dhcprelayd.DHCPRelayIntf)

	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "10.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv4Intf2 := &infra.IPv4IntfProperty{
		IpAddr:  "20.0.0.2",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	ipv4Intf3 := &infra.IPv4IntfProperty{
		IpAddr:  "30.0.0.3",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 3,
		IfRef:   "eth2",
		State:   false,
	}
	iMgr.IPv4IntfProps[1] = ipv4Intf1
	iMgr.IPv4IntfProps[2] = ipv4Intf2
	iMgr.IPv4IntfProps[3] = ipv4Intf3
	iMgr.IPv4IfRefToIfIndex["eth0"] = 1
	iMgr.IPv4IfRefToIfIndex["eth1"] = 2
	iMgr.IPv4IfRefToIfIndex["eth2"] = 3

	drav4Global := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav4IntfCfg1 := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.0.0.1"},
	}
	drav4IntfCfg2 := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"20.0.0.1"},
	}
	drav4IntfCfg3 := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  "eth2",
		Enable:   true,
		ServerIp: []string{"30.0.0.1"},
	}
	iMgr.DRAv4Global = drav4Global
	iMgr.DRAv4Intfs["eth0"] = drav4IntfCfg1
	iMgr.DRAv4Intfs["eth1"] = drav4IntfCfg2
	iMgr.DRAv4Intfs["eth2"] = drav4IntfCfg3
	iMgr.ActiveDRAv4Intfs[1] = drav4IntfCfg1
	iMgr.ActiveDRAv4Intfs[2] = drav4IntfCfg2
}

func TestGetDRAv4ClientState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4ClientState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv4ClientState("aa:bb:cc:dd:ee:ff")
	if err != nil {
		t.Errorf("FAIL: GetDRAv4ClientState")
		return
	}
	t.Log("PASS: GetDRAv4ClientState")
}

func TestGetDRAv4ClientStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4ClientStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	resultClientStates := []*dhcprelayd.DHCPRelayClientState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv4ClientState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultClientStates = append(resultClientStates, bulkInfo.DHCPRelayClientStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultClientStates) != 3 {
		t.Errorf("FAIL: GetDRAv4ClientStateSlice")
		return
	}
	t.Log("PASS: GetDRAv4ClientStateSlice")
}

func TestGetDRAv4IntfState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv4IntfState("eth0")
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfState")
		return
	}
	t.Log("PASS: GetDRAv4IntfState")
}

func TestGetDRAv4IntfStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	resultIntfStates := []*dhcprelayd.DHCPRelayIntfState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv4IntfState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultIntfStates = append(resultIntfStates, bulkInfo.DHCPRelayIntfStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultIntfStates) != 3 {
		t.Errorf("FAIL: GetDRAv4IntfStateSlice")
		return
	}
	t.Log("PASS: GetDRAv4IntfStateSlice")
}

func TestGetDRAv4IntfServerState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfServerState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv4IntfServerState("eth0", "10.0.0.1")
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfServerState")
		return
	}
	_, err = draMgr.GetDRAv4IntfServerState("eth0", "20.0.0.1")
	if err == nil {
		t.Errorf("FAIL: GetDRAv4IntfServerState")
		return
	}
	t.Log("PASS: GetDRAv4IntfServerState")
}

func TestGetDRAv4IntfServerStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv4IntfServerStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV4DummyStates()
	pProcF := GetV4DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc4 = pProcF
	PopulateV4DummyInfra(draMgr.IMgr)
	resultIntfServerStates := []*dhcprelayd.DHCPRelayIntfServerState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv4IntfServerState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultIntfServerStates = append(resultIntfServerStates, bulkInfo.DHCPRelayIntfServerStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultIntfServerStates) != 3 {
		t.Errorf("FAIL: GetDRAv4IntfServerStateSlice")
		return
	}
	t.Log("PASS: GetDRAv4IntfServerStateSlice")
}

// IPv4 Create Notify
func TestV4ProcessAsicdNotification1(t *testing.T) {
	//	draMgr, err := InitTestDRAMgr()
	//	if err != nil {
	//		t.Errorf("FAIL: IPv4 Create Notify")
	//		return
	//	}
	//	PopulateV4DummyInfra(draMgr.IMgr)
	//	drav4IntfCfg4 := &dhcprelayd.DHCPRelayIntf{
	//		IntfRef:  "eth3",
	//		Enable:   true,
	//		ServerIp: []string{"40.0.0.1"},
	//	}
	//	draMgr.IMgr.DRAv4Intfs["eth3"] = drav4IntfCfg4
	//	notf := commonDefs.IPv4IntfNotifyMsg{
	//		MsgType: commonDefs.NOTIFY_IPV4INTF_CREATE,
	//		IpAddr:  "40.0.0.1",
	//		IfIndex: 4,
	//		IntfRef: "eth3",
	//	}
	//	draMgr.ProcessAsicdNotification(notf)
	//	fmt.Println(draMgr.IMgr.GetActiveDRAv4IntfCount())
	//	if draMgr.IMgr.GetActiveDRAv4IntfCount() != 4 {
	//		t.Errorf("FAIL: IPv4 Create Notify")
	//		return
	//	}
	//	t.Log("PASS: IPv4 Create Notify")
}

// IPv4 Delete Notify
func TestV4ProcessAsicdNotification2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv4 Delete Notify")
		return
	}
	PopulateV4DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv4IntfNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV4INTF_DELETE,
		IpAddr:  "11.0.0.1",
		IfIndex: 1,
		IntfRef: "eth0",
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() != 1 {
		t.Errorf("FAIL: IPv4 Delete Notify")
		return
	}
	t.Log("PASS: IPv4 Delete Notify")
}

// IPv4 State Up Notify
func TestV4ProcessAsicdNotification3(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv4 State Up Notify")
		return
	}
	PopulateV4DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv4L3IntfStateNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE,
		IpAddr:  "13.0.0.1",
		IfIndex: 3,
		IfState: 1,
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() != 3 {
		t.Errorf("FAIL: IPv4 State Up Notify")
		return
	}
	t.Log("PASS: IPv4 State Up Notify")
}

// IPv4 State Down Notify
func TestV4ProcessAsicdNotification4(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv4 State Down Notify")
		return
	}
	PopulateV4DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv4L3IntfStateNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE,
		IpAddr:  "12.0.0.1",
		IfIndex: 2,
		IfState: 0,
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv4IntfCount() != 1 {
		t.Errorf("FAIL: IPv4 State Down Notify")
		return
	}
	t.Log("PASS: IPv4 State Down Notify")
}

func PopulateDummyIPv4Infra(iMgr *infra.InfraMgr) {
	ipv4Intf1 := &infra.IPv4IntfProperty{
		IpAddr:  "11.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv4Intf2 := &infra.IPv4IntfProperty{
		IpAddr:  "12.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	ipv4Intf3 := &infra.IPv4IntfProperty{
		IpAddr:  "13.0.0.1",
		Netmask: net.CIDRMask(24, 32),
		IfIndex: 3,
		IfRef:   "eth2",
		State:   false,
	}
	iMgr.IPv4IntfProps[1] = ipv4Intf1
	iMgr.IPv4IntfProps[2] = ipv4Intf2
	iMgr.IPv4IntfProps[3] = ipv4Intf3
	iMgr.IPv4IfRefToIfIndex["eth0"] = 1
	iMgr.IPv4IfRefToIfIndex["eth1"] = 2
	iMgr.IPv4IfRefToIfIndex["eth2"] = 3
	//	drav4IntfCfg1 := &dhcprelayd.DHCPRelayIntf{
	//		IntfRef:  "eth0",
	//		Enable:   true,
	//		ServerIp: []string{"10.0.0.1"},
	//	}
	//	iMgr.DRAv4Intfs["eth0"] = drav4IntfCfg1
	//	iMgr.ActiveDRAv4Intfs[1] = drav4IntfCfg1
	//	iMgr.DRAv4Global = &dhcprelayd.DHCPRelayGlobal{
	//		Vrf:           "default",
	//		Enable:        true,
	//		HopCountLimit: 32,
	//	}
}

// Daemon reload test (Read DB etc.)
func TestV4DaemonReload(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DRAv4 Daemon Reload")
		return
	}
	PopulateDummyIPv4Infra(draMgr.IMgr)
	draMgr.InitDRAMgr()
	if draMgr.IMgr.GetActiveDRAv4IntfCount() != 2 {
		t.Errorf("FAIL: DRAv4 Daemon Reload")
		return
	}
	t.Log("PASS: DRAv4 Daemon Reload")
}

// drav6
func TestCreateDRAv6Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	_, err = draMgr.CreateDRAv6Interface(drav6IntfCfg)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	if !draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	t.Log("PASS: CreateDRAv6Interface")
}

func TestCreateDRAv6Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	ipv6Intf := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	_, err = draMgr.CreateDRAv6Interface(drav6IntfCfg)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	t.Log("PASS: CreateDRAv6Interface")
}

func TestCreateDRAv6Intf3(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	ipv6Intf := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"10.1.1.3"},
	}
	_, err = draMgr.CreateDRAv6Interface(drav6IntfCfg)
	if err == nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	t.Log("PASS: CreateDRAv6Interface")
}

func TestCreateDRAv6Intf4(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1"},
	}
	_, err = draMgr.CreateDRAv6Interface(drav6IntfCfg)
	if err == nil {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Interface")
		return
	}
	t.Log("PASS: CreateDRAv6Interface")
}

func TestUpdateDRAv6Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	ipv6Intf := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[1] = draIntfCfgPre
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc6.SetEnabledFlag()
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   false,
		ServerIp: []string{"2456:db8::1"},
	}
	_, err = draMgr.UpdateDRAv6Interface(
		draIntfCfgPre, drav6IntfCfg, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	draIntfCfgPost := draMgr.IMgr.DRAv6Intfs["eth0"]
	expectedIntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   false,
		ServerIp: []string{"2456:db8::1"},
	}
	if draIntfCfgPost.Enable != expectedIntfCfg.Enable {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	if len(draIntfCfgPost.ServerIp) != len(expectedIntfCfg.ServerIp) {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	for i, ipAddr := range draIntfCfgPost.ServerIp {
		if ipAddr != expectedIntfCfg.ServerIp[i] {
			t.Errorf("FAIL: UpdateDRAv6Interface")
			return
		}
	}
	t.Log("PASS: UpdateDRAv6Interface")
}

func TestUpdateDRAv6Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	ipv6Intf := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav6IntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{},
	}
	draMgr.PProc6.SetEnabledFlag()
	_, err = draMgr.UpdateDRAv6Interface(
		draIntfCfgPre, drav6IntfCfg, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	if !draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	draIntfCfgPost := draMgr.IMgr.DRAv6Intfs["eth0"]
	expectedIntfCfg := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{},
	}
	if draIntfCfgPost.Enable != expectedIntfCfg.Enable {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	if len(draIntfCfgPost.ServerIp) != len(expectedIntfCfg.ServerIp) {
		t.Errorf("FAIL: UpdateDRAv6Interface")
		return
	}
	for i, ipAddr := range draIntfCfgPost.ServerIp {
		if ipAddr != expectedIntfCfg.ServerIp[i] {
			t.Errorf("FAIL: UpdateDRAv6Interface")
			return
		}
	}
	t.Log("PASS: UpdateDRAv6Interface")
}

func TestDeleteDRAv6Intf1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	ipv6Intf := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draIntfCfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntfCfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[1] = draIntfCfgPre
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc6.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv6Interface("eth0")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	if _, ok := draMgr.IMgr.DRAv6Intfs["eth0"]; ok {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	t.Log("PASS: DeleteDRAv6Interface")
}

func TestDeleteDRAv6Intf2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv6Intf2 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IntfProps[2] = ipv6Intf2
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"2466:db8::1", "2466:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv6Intfs["eth1"] = draIntf2CfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[2] = draIntf2CfgPre
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc6.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv6Interface("eth0")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	if !draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	if _, ok := draMgr.IMgr.DRAv6Intfs["eth0"]; ok {
		t.Errorf("FAIL: DeleteDRAv6Interface")
		return
	}
	t.Log("PASS: DeleteDRAv6Interface")
}

func TestCreateDRAv6Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draIntf1CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntf1CfgPre
	drav6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.CreateDRAv6Global(drav6Global)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	if !draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	if draMgr.IMgr.DRAv6Global == nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	t.Log("PASS: CreateDRAv6Global")
}

func TestCreateDRAv6Global2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	drav6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.CreateDRAv6Global(drav6Global)
	if err != nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	if draMgr.IMgr.DRAv6Global == nil {
		t.Errorf("FAIL: CreateDRAv6Global")
		return
	}
	t.Log("PASS: CreateDRAv6Global")
}

func TestUpdateDRAv6Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv6Intf2 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IntfProps[2] = ipv6Intf2
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"2466:db8::1", "2466:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv6Intfs["eth1"] = draIntf2CfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[2] = draIntf2CfgPre
	oldDraV6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.IMgr.DRAv6Global = oldDraV6Global
	draMgr.PProc6.SetEnabledFlag()
	newDraV6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        false,
		HopCountLimit: 31,
	}
	_, err = draMgr.UpdateDRAv6Global(
		oldDraV6Global, newDraV6Global, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if draMgr.IMgr.GetActiveDRAv6IntfCount() > 0 {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	curDraVgGlobal := draMgr.IMgr.DRAv6Global
	if curDraVgGlobal.Enable != newDraV6Global.Enable {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if curDraVgGlobal.HopCountLimit != newDraV6Global.HopCountLimit {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	t.Log("PASS: UpdateDRAv6Global")
}

func TestUpdateDRAv6Global2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv6Intf2 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IntfProps[2] = ipv6Intf2
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth1"] = 2
	draIntf1CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draIntf2CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"2466:db8::1", "2466:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.DRAv6Intfs["eth1"] = draIntf2CfgPre
	oldDraV6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        false,
		HopCountLimit: 32,
	}
	draMgr.IMgr.DRAv6Global = oldDraV6Global
	draMgr.PProc6.SetEnabledFlag()
	newDraV6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	_, err = draMgr.UpdateDRAv6Global(
		oldDraV6Global, newDraV6Global, []bool{false, true, true, false}, nil)
	if err != nil {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if !draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if draMgr.IMgr.GetActiveDRAv6IntfCount() <= 0 {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	curDraVgGlobal := draMgr.IMgr.DRAv6Global
	if curDraVgGlobal.Enable != newDraV6Global.Enable {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	if curDraVgGlobal.HopCountLimit != newDraV6Global.HopCountLimit {
		t.Errorf("FAIL: UpdateDRAv6Global")
		return
	}
	t.Log("PASS: UpdateDRAv6Global")
}

func TestDeleteDRAv6Global1(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Global")
		return
	}
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2031:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	draMgr.IMgr.IPv6IntfProps[1] = ipv6Intf1
	draMgr.IMgr.IPv6IfRefToIfIndex["eth0"] = 1
	draIntf1CfgPre := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2456:db8::1", "2456:db8::2"},
	}
	draMgr.IMgr.DRAv6Intfs["eth0"] = draIntf1CfgPre
	draMgr.IMgr.ActiveDRAv6Intfs[1] = draIntf1CfgPre
	draMgr.IMgr.DRAv6Global = &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	draMgr.PProc6.SetEnabledFlag()
	_, err = draMgr.DeleteDRAv6Global("default")
	if err != nil {
		t.Errorf("FAIL: DeleteDRAv6Global")
		return
	}
	if draMgr.PProc6.GetEnabledFlag() {
		t.Errorf("FAIL: DeleteDRAv6Global")
		return
	}
	if draMgr.IMgr.DRAv6Global != nil {
		t.Errorf("FAIL: DeleteDRAv6Global")
		return
	}
	t.Log("PASS: DeleteDRAv6Global")
}

//States
func GetV6DummyStates() (
	[]*dhcprelayd.DHCPv6RelayClientState,
	map[string]*dhcprelayd.DHCPv6RelayClientState,
	[]*dhcprelayd.DHCPv6RelayIntfState,
	map[string]*dhcprelayd.DHCPv6RelayIntfState,
	[]*dhcprelayd.DHCPv6RelayIntfServerState,
	map[string]*dhcprelayd.DHCPv6RelayIntfServerState) {

	clientStates := []*dhcprelayd.DHCPv6RelayClientState{}
	clientStates = append(clientStates, &dhcprelayd.DHCPv6RelayClientState{
		MacAddr:           "fe80::1",
		ClientSolicit:     time.Now().String(),
		ClientRequest:     time.Now().String(),
		ClientConfirm:     time.Now().String(),
		ClientRenew:       time.Now().String(),
		ClientRebind:      time.Now().String(),
		ClientRelease:     time.Now().String(),
		ClientDecline:     time.Now().String(),
		ClientInfoRequest: time.Now().String(),
		ClientRequests:    1124,
		ClientResponses:   1345,
		ServerAdvertise:   time.Now().String(),
		ServerReply:       time.Now().String(),
		ServerReconfigure: time.Now().String(),
		ServerRequests:    932,
		ServerResponses:   878,
	},
	)
	clientStates = append(clientStates, &dhcprelayd.DHCPv6RelayClientState{
		MacAddr:           "fe80::2",
		ClientSolicit:     time.Now().String(),
		ClientRequest:     time.Now().String(),
		ClientConfirm:     time.Now().String(),
		ClientRenew:       time.Now().String(),
		ClientRebind:      time.Now().String(),
		ClientRelease:     time.Now().String(),
		ClientDecline:     time.Now().String(),
		ClientInfoRequest: time.Now().String(),
		ClientRequests:    1224,
		ClientResponses:   1445,
		ServerAdvertise:   time.Now().String(),
		ServerReply:       time.Now().String(),
		ServerReconfigure: time.Now().String(),
		ServerRequests:    942,
		ServerResponses:   888,
	},
	)
	clientStates = append(clientStates, &dhcprelayd.DHCPv6RelayClientState{
		MacAddr:           "fe80::3",
		ClientSolicit:     time.Now().String(),
		ClientRequest:     time.Now().String(),
		ClientConfirm:     time.Now().String(),
		ClientRenew:       time.Now().String(),
		ClientRebind:      time.Now().String(),
		ClientRelease:     time.Now().String(),
		ClientDecline:     time.Now().String(),
		ClientInfoRequest: time.Now().String(),
		ClientRequests:    1324,
		ClientResponses:   1545,
		ServerAdvertise:   time.Now().String(),
		ServerReply:       time.Now().String(),
		ServerReconfigure: time.Now().String(),
		ServerRequests:    952,
		ServerResponses:   898,
	},
	)

	intfStates := []*dhcprelayd.DHCPv6RelayIntfState{}
	intfStates = append(intfStates, &dhcprelayd.DHCPv6RelayIntfState{
		IntfRef:           "eth0",
		TotalDrops:        1123,
		TotalDhcpClientRx: 34567,
		TotalDhcpClientTx: 33457,
		TotalDhcpServerRx: 32567,
		TotalDhcpServerTx: 34222,
	},
	)
	intfStates = append(intfStates, &dhcprelayd.DHCPv6RelayIntfState{
		IntfRef:           "eth1",
		TotalDrops:        1223,
		TotalDhcpClientRx: 35567,
		TotalDhcpClientTx: 34457,
		TotalDhcpServerRx: 33567,
		TotalDhcpServerTx: 35222,
	},
	)
	intfStates = append(intfStates, &dhcprelayd.DHCPv6RelayIntfState{
		IntfRef:           "eth2",
		TotalDrops:        1323,
		TotalDhcpClientRx: 36567,
		TotalDhcpClientTx: 35457,
		TotalDhcpServerRx: 34567,
		TotalDhcpServerTx: 36222,
	},
	)

	intfServerStates := []*dhcprelayd.DHCPv6RelayIntfServerState{}
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPv6RelayIntfServerState{
		IntfRef:   "eth0",
		ServerIp:  "2001::1",
		Request:   34567,
		Responses: 33457,
	},
	)
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPv6RelayIntfServerState{
		IntfRef:   "eth1",
		ServerIp:  "2001::2",
		Request:   35567,
		Responses: 34457,
	},
	)
	intfServerStates = append(intfServerStates, &dhcprelayd.DHCPv6RelayIntfServerState{
		IntfRef:   "eth2",
		ServerIp:  "2001::3",
		Request:   36567,
		Responses: 35457,
	},
	)

	intfRefMap := make(map[string]int)
	intfRefMap["eth0"] = 1
	intfRefMap["eth1"] = 2
	intfRefMap["eth2"] = 3

	clientStateMap := make(map[string]*dhcprelayd.DHCPv6RelayClientState)
	intfStateMap := make(map[string]*dhcprelayd.DHCPv6RelayIntfState)
	intfServerStateMap := make(map[string]*dhcprelayd.DHCPv6RelayIntfServerState)
	for _, clientState := range clientStates {
		clientStateMap[clientState.MacAddr] = clientState
	}
	for _, intfState := range intfStates {
		intfStateMap[intfState.IntfRef] = intfState
	}
	for _, intfServerState := range intfServerStates {
		intfServerStateKey := intfServerState.IntfRef + "_" + intfServerState.ServerIp
		intfServerStateMap[intfServerStateKey] = intfServerState
	}

	return clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap
}

func GetV6DummyStateProcessor(
	clientStates []*dhcprelayd.DHCPv6RelayClientState,
	clientStateMap map[string]*dhcprelayd.DHCPv6RelayClientState,
	intfStates []*dhcprelayd.DHCPv6RelayIntfState,
	intfStateMap map[string]*dhcprelayd.DHCPv6RelayIntfState,
	intfServerStates []*dhcprelayd.DHCPv6RelayIntfServerState,
	intfServerStateMap map[string]*dhcprelayd.DHCPv6RelayIntfServerState) *Processor6Fake {

	pProcF := &Processor6Fake{}
	pProcF.ClientStateSlice = clientStates
	pProcF.ClientStateMap = clientStateMap
	pProcF.IntfStateSlice = intfStates
	pProcF.IntfStateMap = intfStateMap
	pProcF.IntfServerStateSlice = intfServerStates
	pProcF.IntfServerStateMap = intfServerStateMap
	return pProcF
}

func PopulateV6DummyInfra(iMgr *infra.InfraMgr) {
	iMgr.IPv6IfRefToIfIndex = make(map[string]int)
	iMgr.IPv6IntfProps = make(map[int]*infra.IPv6IntfProperty)
	iMgr.DRAv6Intfs = make(map[string]*dhcprelayd.DHCPv6RelayIntf)
	iMgr.ActiveDRAv6Intfs = make(map[int]*dhcprelayd.DHCPv6RelayIntf)

	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv6Intf2 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	ipv6Intf3 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::3",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 3,
		IfRef:   "eth2",
		State:   false,
	}
	iMgr.IPv6IntfProps[1] = ipv6Intf1
	iMgr.IPv6IntfProps[2] = ipv6Intf2
	iMgr.IPv6IntfProps[3] = ipv6Intf3
	iMgr.IPv6IfRefToIfIndex["eth0"] = 1
	iMgr.IPv6IfRefToIfIndex["eth1"] = 2
	iMgr.IPv6IfRefToIfIndex["eth2"] = 3

	drav6Global := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           "default",
		Enable:        true,
		HopCountLimit: 32,
	}
	drav6IntfCfg1 := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth0",
		Enable:   true,
		ServerIp: []string{"2001::1"},
	}
	drav6IntfCfg2 := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth1",
		Enable:   true,
		ServerIp: []string{"2001::2"},
	}
	drav6IntfCfg3 := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  "eth2",
		Enable:   true,
		ServerIp: []string{"2001::3"},
	}
	iMgr.DRAv6Global = drav6Global
	iMgr.DRAv6Intfs["eth0"] = drav6IntfCfg1
	iMgr.DRAv6Intfs["eth1"] = drav6IntfCfg2
	iMgr.DRAv6Intfs["eth2"] = drav6IntfCfg3
	iMgr.ActiveDRAv6Intfs[1] = drav6IntfCfg1
	iMgr.ActiveDRAv6Intfs[2] = drav6IntfCfg2
}

func TestGetDRAv6ClientState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6ClientState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv6ClientState("fe80::1")
	if err != nil {
		t.Errorf("FAIL: GetDRAv6ClientState")
		return
	}
	t.Log("PASS: GetDRAv6ClientState")
}

func TestGetDRAv6ClientStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6ClientStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	resultClientStates := []*dhcprelayd.DHCPv6RelayClientState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv6ClientState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultClientStates = append(resultClientStates, bulkInfo.DHCPv6RelayClientStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultClientStates) != 3 {
		t.Errorf("FAIL: GetDRAv6ClientStateSlice")
		return
	}
	t.Log("PASS: GetDRAv6ClientStateSlice")
}

func TestGetDRAv6IntfState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv6IntfState("eth0")
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfState")
		return
	}
	t.Log("PASS: GetDRAv6IntfState")
}

func TestGetDRAv6IntfStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	resultIntfStates := []*dhcprelayd.DHCPv6RelayIntfState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv6IntfState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultIntfStates = append(resultIntfStates, bulkInfo.DHCPv6RelayIntfStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultIntfStates) != 3 {
		t.Errorf("FAIL: GetDRAv6IntfStateSlice")
		return
	}
	t.Log("PASS: GetDRAv6IntfStateSlice")
}

func TestGetDRAv6IntfServerState(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfServerState")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	_, err = draMgr.GetDRAv6IntfServerState("eth1", "2001::2")
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfServerState")
		return
	}
	_, err = draMgr.GetDRAv6IntfServerState("eth1", "2001::1")
	if err == nil {
		t.Errorf("FAIL: GetDRAv6IntfServerState")
		return
	}
	t.Log("PASS: GetDRAv6IntfServerState")
}

func TestGetDRAv6IntfServerStateSlice(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: GetDRAv6IntfServerStateSlice")
		return
	}
	clientStates, clientStateMap, intfStates, intfStateMap, intfServerStates, intfServerStateMap :=
		GetV6DummyStates()
	pProcF := GetV6DummyStateProcessor(
		clientStates,
		clientStateMap,
		intfStates,
		intfStateMap,
		intfServerStates,
		intfServerStateMap,
	)
	draMgr.PProc6 = pProcF
	PopulateV6DummyInfra(draMgr.IMgr)
	resultIntfServerStates := []*dhcprelayd.DHCPv6RelayIntfServerState{}
	curIdx := 0
	count := 2
	for {
		bulkInfo := draMgr.GetBulkDRAv6IntfServerState(curIdx, count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			resultIntfServerStates = append(resultIntfServerStates, bulkInfo.DHCPv6RelayIntfServerStateList[i])
		}
		if !bulkInfo.More {
			break
		}
		curIdx = int(bulkInfo.EndIdx)
	}
	if len(resultIntfServerStates) != 3 {
		t.Errorf("FAIL: GetDRAv6IntfServerStateSlice")
		return
	}
	t.Log("PASS: GetDRAv6IntfServerStateSlice")
}

// IPv6 Create Notify
func TestV6ProcessAsicdNotification1(t *testing.T) {
	//	draMgr, err := InitTestDRAMgr()
	//	if err != nil {
	//		t.Errorf("FAIL: IPv6 Create Notify")
	//		return
	//	}
	//	PopulateV6DummyInfra(draMgr.IMgr)
	//	drav6IntfCfg4 := &dhcprelayd.DHCPv6RelayIntf{
	//		IntfRef:  "eth3",
	//		Enable:   true,
	//		ServerIp: []string{"2001::4"},
	//	}
	//	draMgr.IMgr.DRAv6Intfs["eth3"] = drav6IntfCfg4
	//	notf := commonDefs.IPv6IntfNotifyMsg{
	//		MsgType: commonDefs.NOTIFY_IPV6INTF_CREATE,
	//		IpAddr:  "2091:db8::4",
	//		IfIndex: 4,
	//		IntfRef: "eth3",
	//	}
	//	draMgr.ProcessAsicdNotification(notf)
	//	fmt.Println(draMgr.IMgr.GetActiveDRAv6IntfCount())
	//	if draMgr.IMgr.GetActiveDRAv6IntfCount() != 4 {
	//		t.Errorf("FAIL: IPv6 Create Notify")
	//		return
	//	}
	//	t.Log("PASS: IPv6 Create Notify")
}

// IPv6 Delete Notify
func TestV6ProcessAsicdNotification2(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv6 Delete Notify")
		return
	}
	PopulateV6DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv6IntfNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV6INTF_DELETE,
		IpAddr:  "2091:db8::1",
		IfIndex: 1,
		IntfRef: "eth0",
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() != 1 {
		t.Errorf("FAIL: IPv6 Delete Notify")
		return
	}
	t.Log("PASS: IPv6 Delete Notify")
}

// IPv6 State Up Notify
func TestV6ProcessAsicdNotification3(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv6 State Up Notify")
		return
	}
	PopulateV6DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv6L3IntfStateNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE,
		IpAddr:  "2091:db8::3",
		IfIndex: 3,
		IfState: 1,
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() != 3 {
		t.Errorf("FAIL: IPv6 State Up Notify")
		return
	}
	t.Log("PASS: IPv6 State Up Notify")
}

// IPv6 State Down Notify
func TestV6ProcessAsicdNotification4(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: IPv6 State Down Notify")
		return
	}
	PopulateV6DummyInfra(draMgr.IMgr)
	notf := commonDefs.IPv6L3IntfStateNotifyMsg{
		MsgType: commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE,
		IpAddr:  "2091:db8::2",
		IfIndex: 2,
		IfState: 0,
	}
	draMgr.ProcessAsicdNotification(notf)
	if draMgr.IMgr.GetActiveDRAv6IntfCount() != 1 {
		t.Errorf("FAIL: IPv6 State Down Notify")
		return
	}
	t.Log("PASS: IPv6 State Down Notify")
}

func PopulateDummyIPv6Infra(iMgr *infra.InfraMgr) {
	ipv6Intf1 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::1",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 1,
		IfRef:   "eth0",
		State:   true,
	}
	ipv6Intf2 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::2",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 2,
		IfRef:   "eth1",
		State:   true,
	}
	ipv6Intf3 := &infra.IPv6IntfProperty{
		IpAddr:  "2091:db8::3",
		Netmask: net.CIDRMask(64, 128),
		IfIndex: 3,
		IfRef:   "eth2",
		State:   false,
	}
	iMgr.IPv6IntfProps[1] = ipv6Intf1
	iMgr.IPv6IntfProps[2] = ipv6Intf2
	iMgr.IPv6IntfProps[3] = ipv6Intf3
	iMgr.IPv6IfRefToIfIndex["eth0"] = 1
	iMgr.IPv6IfRefToIfIndex["eth1"] = 2
	iMgr.IPv6IfRefToIfIndex["eth2"] = 3
}

// Daemon reload test (Read DB etc.)
func TestV6DaemonReload(t *testing.T) {
	draMgr, err := InitTestDRAMgr()
	if err != nil {
		t.Errorf("FAIL: DRAv6 Daemon Reload")
		return
	}
	PopulateDummyIPv6Infra(draMgr.IMgr)
	draMgr.InitDRAMgr()
	if draMgr.IMgr.GetActiveDRAv6IntfCount() != 2 {
		t.Errorf("FAIL: DRAv6 Daemon Reload")
		return
	}
	t.Log("PASS: DRAv6 Daemon Reload")
}
