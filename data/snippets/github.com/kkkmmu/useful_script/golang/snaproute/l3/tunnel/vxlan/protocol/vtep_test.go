// _vtep_test.go
package vxlan

import (
	"asicd/pluginManager/pluginCommon"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"testing"
	"time"
	"utils/commonDefs"
	"utils/logging"
)

var vtepcreatedone chan bool
var vxlancreatedone chan bool
var vtepdeletedone chan bool
var vxlandeletedone chan bool

type mockintf struct {
	BaseClientIntf
	// triggers for test to fail various interface functions to test behavior
	failCreateVtep           bool
	failCreateVxlan          bool
	failDeleteVtep           bool
	failDeleteVxlan          bool
	failGetIntfInfo          bool
	failCreateAccessPortVlan bool
	failDeleteAccessPortVlan bool
	failGetAccessPorts       bool
	failGetNextHop           bool
	failResolveNexHop        bool
	faillink                 bool
}

func (b mockintf) IsClientIntfType(client VXLANClientIntf, clientStr string) bool {
	switch client.(type) {
	case *BaseClientIntf:
		if clientStr == "SnapMockTestClient" {
			return true
		}
	}
	return false
}

func (b mockintf) GetIntfInfo(name string, intfchan chan<- MachineEvent) {

	logger.Info("MOCK: Calling GetIntfInfo")
	if !b.failGetIntfInfo {
		mac, _ := net.ParseMAC("00:01:02:03:04:05")
		ip := net.ParseIP("100.1.1.1")
		intfinfo := VxlanIntfInfo{
			Command:  VxlanCommandCreate,
			IntfName: "eth0",
			IfIndex:  pluginCommon.GetIfIndexFromIdType(1, commonDefs.IfTypePort),
			Mac:      mac,
			Ip:       ip,
		}
		intfchan <- MachineEvent{
			E:    VxlanVtepEventSrcInterfaceResolved,
			Src:  "TEST",
			Data: intfinfo,
		}
	}
}
func (b mockintf) CreateVtep(vtep *VtepDbEntry, vtepname chan<- MachineEvent) {
	if !b.failCreateVtep {
		logger.Info(fmt.Sprintf("Create vtep %#v", vtep))
		event := MachineEvent{
			E:   VxlanVtepEventHwConfigComplete,
			Src: VXLANSnapClientStr,
			Data: VtepCreateCfgData{
				VtepName: "Vtep0Int",
				IfIndex:  12345,
			},
		}
		vtepname <- event

		vtepcreatedone <- true
	}
}
func (b mockintf) DeleteVtep(vtep *VtepDbEntry) {
	if !b.failDeleteVtep {
		vtepdeletedone <- true
	}
}
func (b mockintf) CreateVxlan(vxlan *VxlanDbEntry) {
	if !b.failCreateVxlan {
		logger.Info("MOCK: Calling Vxlan Create done")
		vxlancreatedone <- true
	}
}
func (b mockintf) DeleteVxlan(vxlan *VxlanDbEntry) {
	if !b.failDeleteVxlan {
		logger.Info("MOCK: calling Delete Vxlan")
		vxlandeletedone <- true
	}
}
func (b mockintf) GetAccessPorts(vlan uint16) {
	if !b.failGetAccessPorts {

	}
}
func (b mockintf) UpdateAccessPorts() {

}
func (b mockintf) CreateAccessPortVlan(vlan uint16, intfList []int) {
	if !b.failCreateAccessPortVlan {

	}
}
func (b mockintf) DeleteAccessPortVlan(vlan uint16, intfList []int) {
	if !b.failDeleteAccessPortVlan {

	}
}
func (b mockintf) GetNextHopInfo(ip net.IP, nexthopchan chan<- MachineEvent) bool {
	logger.Info("MOCK: Calling GetNextHopInfo")
	nexthopip := net.ParseIP("100.1.1.2")
	if !b.failGetNextHop {

		nexthop := VxlanNextHopIp{
			Command:   VxlanCommandCreate,
			Ip:        ip,
			Intf:      1,
			IntfName:  "eth0",
			NextHopIp: nexthopip,
		}

		nexthopchan <- MachineEvent{
			E:    VxlanVtepEventNextHopInfoResolved,
			Src:  "TEST",
			Data: nexthop,
		}
		return true
	} else {
		logger.Info("MOCK: force fail")
	}
	return false
}
func (b mockintf) ResolveNextHopMac(nextHopIp net.IP, nextHopIfName string, nexthopmacchan chan<- MachineEvent) {
	logger.Info("MOCK: Calling ResolveNextHopMac")
	mac, _ := net.ParseMAC("00:55:44:33:22:11")
	if !b.failResolveNexHop {
		nexthopmac := mac
		nexthopmacchan <- MachineEvent{
			E:    VxlanVtepEventNextHopInfoMacResolved,
			Src:  "TEST",
			Data: nexthopmac,
		}
	} else {
		logger.Info("MOCK: force fail")
	}
}

func (b mockintf) GetLinkState(ifname string) string {
	if b.faillink {
		return "DOWN"
	}
	return "UP"
}

func MockFuncRxTx(vtep *VtepDbEntry) {
	logger.Info(fmt.Sprintf("MOCK: going to listen on interface %s", vtep.VtepName))
}

func setup() {
	VxlanGlobalStateSet(VXLAN_GLOBAL_ENABLE)
	setVxlanTestLogger()
	vtepcreatedone = make(chan bool, 1)
	vtepdeletedone = make(chan bool, 1)
	vxlancreatedone = make(chan bool, 1)
	vxlandeletedone = make(chan bool, 1)
}

func teardown(t *testing.T) {
	if len(vxlanDB) > 0 {
		t.Error("Error vxlandb not cleaned up")
	}
	if len(vxlanDbList) > 0 {
		t.Error("Error vxlandblist not cleaned up")
	}
	if len(vtepDB) > 0 {
		t.Error("Error vtepdb not cleaned up")
	}
	if len(vtepDbList) > 0 {
		t.Error("Error vtepdblist not cleaned up")
	}
	logger.Close()
	logger = nil
	close(vtepcreatedone)
	close(vxlancreatedone)
	close(vtepdeletedone)
	close(vxlandeletedone)
	exec.Command("/bin/rm", "UsrConfDb.db")
	DeRegisterClients()
	SetLogger(nil)
}

func setVxlanTestLogger() {

	logger, _ := logging.NewLogger("vxland", "TEST", true)
	SetLogger(logger)
}

func TimerTest(v *VtepDbEntry, exitchan chan<- bool) {
	cnt := 0
	for {
		time.Sleep(time.Millisecond * 10)
		if v.ticksTillConfig > 5 {
			exitchan <- true
			return
		}
		if cnt > 10 {
			exitchan <- true
			return
		}
		cnt++
	}
}

// TestVtepNameGeneration
func TestVtepNameGeneration(t *testing.T) {

	//setVxlanTestLogger()

	name := GenInternalVtepName()
	if name != "Vtep1" {
		t.Errorf("Generated Name is not correct expected Vtep1 got", name)
	}

	name = GenInternalVtepName()
	if name != "Vtep2" {
		t.Errorf("Generated Name is not correct expected Vtep1 got", name)
	}

	name = GenInternalVtepName()
	if name != "Vtep3" {
		t.Errorf("Generated Name is not correct expected Vtep1 got", name)
	}

	if len(vtepNameIdList) != 0 {
		t.Errorf("VtepNameList is not empty")
	}

	FreeGenInternalVtepName(name)
	if len(vtepNameIdList) != 1 {
		t.Errorf("VtepNameList does not contain the entry that was just freed", vtepNameIdList)
	}

	if vtepNameIdList[0] != 3 {
		t.Errorf("VtepNameList entry id was not freed properly", vtepNameIdList)
	}

	name = GenInternalVtepName()
	if name != "Vtep3" {
		t.Errorf("Generated Name is not correct expected Vtep1 got", name)
	}
	if len(vtepNameIdList) != 0 {
		t.Errorf("VtepNameList is not empty", vtepNameIdList)
	}

	if vtepNameIdCnt != 4 {
		t.Errorf("VtepNameIdCnt is not correct", vtepNameIdCnt)
	}

	FreeGenInternalVtepName("Vtep1")
	FreeGenInternalVtepName("Vtep2")
	FreeGenInternalVtepName("Vtep3")

	// try and delete an entry that does not exist
	FreeGenInternalVtepName("Vtep4")

	for i := 0; i < 100; i++ {
		GenInternalVtepName()
	}

	FreeGenInternalVtepName("Vtep34")

	if vtepNameIdList[0] != 34 {
		t.Errorf("VtepNameList did not delete previous entry Vtep34")
	}

	//SetLogger(nil)
}

// TestFSMValidVxlanVtepCreate:
// Test creation of vxlan before vtep
func TestFSMValidVxlanVtepCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{})

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      0,
	}

	// create the vxlan
	CreateVxLAN(vxlanConfig)

	<-vxlancreatedone

	CreateVtep(vtepConfig)

	// need to wait for test to hwcreate to be called
	<-vtepcreatedone

	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}

	vtep := GetVtepDBEntry(key)
	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() < VxlanVtepStateHwConfig {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateHwConfig, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)

	<-vtepdeletedone

	DeleteVxLAN(vxlanConfig, false)

	<-vxlandeletedone

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}

	VxlanVtepRxTx = oldRxTx
}

// TestFSMValidVtepVxlanCreate:
// Test creation of vtep before vxlan
func TestFSMValidVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{})

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	CreateVtep(vtepConfig)

	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}

	vtep := GetVtepDBEntry(key)

	// delay to allow for time for FSM to come up
	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	<-exitchan

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateDetached {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateDetached, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	// need to wait for test to hwcreate to be called
	CreateVxLAN(vxlanConfig)

	<-vxlancreatedone
	<-vtepcreatedone

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() < VxlanVtepStateHwConfig {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateHwConfig, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)

	<-vtepdeletedone

	DeleteVxLAN(vxlanConfig, false)

	<-vxlandeletedone

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}

}

// TestFSMInValidNumVlanInVxlanMapVtepVxlanCreate:
// check validation, this is a temporary test till we can support multiple vlans in the
// vni map
func TestFSMInValidNumVlanInVxlanMapVtepVxlanCreate(t *testing.T) {

	// setup common test info
	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200, 300},
	}

	err := VxlanConfigCheck(vxlanConfig)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}
	vxlanConfig = &VxlanConfig{
		Enable:       true,
		VNI:          100,
		VlanId:       []uint16{200},
		UntaggedVlan: []uint16{300},
	}
	err = VxlanConfigCheck(vxlanConfig)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}
	vxlanConfig = &VxlanConfig{
		Enable:       true,
		VNI:          100,
		UntaggedVlan: []uint16{200, 300},
	}
	err = VxlanConfigCheck(vxlanConfig)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}

	// check update
	oc := &VxlanConfig{
		Enable:       true,
		VNI:          100,
		UntaggedVlan: []uint16{200},
	}
	nc := &VxlanConfig{
		Enable:       true,
		VNI:          100,
		UntaggedVlan: []uint16{200, 300},
	}
	err = VxlanConfigUpdateCheck(oc, nc)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}
	oc = &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}
	nc = &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200, 300},
	}
	err = VxlanConfigUpdateCheck(oc, nc)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}
	oc = &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}
	nc = &VxlanConfig{
		Enable:       true,
		VNI:          100,
		VlanId:       []uint16{200},
		UntaggedVlan: []uint16{300},
	}
	err = VxlanConfigUpdateCheck(oc, nc)
	if err == nil {
		t.Errorf("Error don't support more than one vlan in create")
	}
}

// TestFSMCreateVtepNoVxlan
// Test that FSM is not running when vxlan is not configured
func TestFSMCreateVtepNoVxlan(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{})

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}

	vtep := GetVtepDBEntry(key)
	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	<-exitchan

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateDetached {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateInit, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)
}

// TestFSMIntfFailVtepVxlanCreate:
// Test creation of vtep and the src interface does not exist
// basically provisioning incomplete
func TestFSMIntfFailVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{
		failGetIntfInfo: true,
	})

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case done := <-vtepcreatedone:
				if done {
					t.Errorf("Vtep should not have been created the interface did not exist")
				}
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateInit {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateInit, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)

	DeleteVxLAN(vxlanConfig, false)

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}

// TestFSMIntfFailThenSendIntfSuccessVtepVxlanCreate:
// Test creation of vtep and the src interface does not exist
// basically provisioning incomplete, then once interface exists
// notify state machine
func TestFSMIntfFailThenSendIntfSuccessVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	x := mockintf{
		failGetIntfInfo: true,
	}
	// all apis are set to not fail
	RegisterClients(x)

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case done := <-vtepcreatedone:
				if done {
					t.Errorf("Vtep should not have been created the interface did not exist")
				}
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateInit {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateInit, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	// allow this to proceed in order to process
	x.failGetIntfInfo = false

	// lets get the info necessary info, this is "sort of" simulating
	// a notification from the client to the server
	mac, _ := net.ParseMAC("00:01:02:03:04:05")
	ip := net.ParseIP("100.1.1.1")
	intfinfo := VxlanIntfInfo{
		IntfName: "eth0",
		IfIndex:  pluginCommon.GetIfIndexFromIdType(1, commonDefs.IfTypePort),
		Mac:      mac,
		Ip:       ip,
	}

	vtep.VxlanVtepMachineFsm.VxlanVtepEvents <- MachineEvent{
		E:    VxlanVtepEventSrcInterfaceResolved,
		Src:  "TEST",
		Data: intfinfo,
	}

	<-vtepcreatedone

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() < VxlanVtepStateHwConfig {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateHwConfig, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)

	<-vtepdeletedone

	DeleteVxLAN(vxlanConfig, false)

	<-vxlandeletedone

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}

// TestFSMNextHopFailVtepVxlanCreate:
// Test next hop ip has not been discovered yet thus verify state
func TestFSMNextHopFailVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{
		failGetNextHop: true,
	})

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case done := <-vtepcreatedone:
				if done {
					t.Errorf("Vtep should not have been created the interface did not exist")
				}
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateInterface {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateInterface, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	if vtep.NextHop.Ip.String() != "0.0.0.0" &&
		vtep.NextHop.Ip.String() != "" &&
		vtep.NextHop.Ip != nil {
		t.Errorf("Why was the next hop IP address[%s] found for interface [%s]", vtep.NextHop.Ip, vtep.SrcIfName)
	}

	DeleteVtep(vtepConfig)

	DeleteVxLAN(vxlanConfig, false)

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}

// TestFSMNextHopFailThenSucceedVtepVxlanCreate:
// Test next hop ip has not been discovered yet thus verify state
// then notify state machine when next hop has been found
func TestFSMNextHopFailThenSucceedVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	x := mockintf{
		failGetNextHop: true,
	}
	// all apis are set to not fail
	RegisterClients(x)

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case done := <-vtepcreatedone:
				if done {
					t.Errorf("Vtep should not have been created the interface did not exist")
				}
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateInterface {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateInterface, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	if vtep.NextHop.Ip.String() != "0.0.0.0" &&
		vtep.NextHop.Ip.String() != "" &&
		vtep.NextHop.Ip != nil {
		t.Errorf("Why was the next hop IP address[%s] found for interface [%s]", vtep.NextHop.Ip, vtep.SrcIfName)
	}

	// notify that next hop found
	x.failGetNextHop = false

	// notify next hop has been found
	nexthopip := net.ParseIP("100.1.1.100")
	nexthopinfo := VtepNextHopInfo{
		Ip:      nexthopip,
		IfIndex: 1,
		IfName:  "eth0",
	}

	vtep.VxlanVtepMachineFsm.VxlanVtepEvents <- MachineEvent{
		E:    VxlanVtepEventNextHopInfoResolved,
		Src:  "TEST",
		Data: nexthopinfo,
	}

	<-vtepcreatedone

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() < VxlanVtepStateHwConfig {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateHwConfig, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	DeleteVtep(vtepConfig)

	DeleteVxLAN(vxlanConfig, false)

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}

// TestFSMResolveNextHopFailVtepVxlanCreate:
// Test that next hop ip mac address does not exist yet
func TestFSMResolveNextHopMacFailVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	// all apis are set to not fail
	RegisterClients(mockintf{
		failResolveNexHop: true,
	})

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case done := <-vtepcreatedone:
				if done {
					t.Errorf("Vtep should not have been created the interface did not exist")
				}
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateNextHopInfo {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateNextHopInfo, vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState())
	}

	if vtep.SrcIp.String() == "0.0.0.0" ||
		vtep.SrcIp.String() == "" ||
		vtep.SrcIp == nil {
		t.Errorf("Why was the IP address[%s] not found for interface [%s]", vtep.SrcIp, vtep.SrcIfName)
	}

	if vtep.DstMac.String() != "00:00:00:00:00:00" &&
		vtep.DstMac != nil {
		t.Errorf("Why was the Dst MAC address[%s] found for interface [%s]", vtep.DstMac, vtep.SrcIfName)
	}

	DeleteVtep(vtepConfig)

	DeleteVxLAN(vxlanConfig, false)

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}

// TestFSMResolveNextHopFailVtepVxlanCreate:
// Test that next hop ip mac address does not exist yet
func xTestFSMlinkDownFailCausingRibReachabilityVtepVxlanCreate(t *testing.T) {

	// setup common test info
	setup()
	defer teardown(t)

	oldRxTx := VxlanVtepRxTx
	VxlanVtepRxTx = MockFuncRxTx

	defer func() {
		VxlanVtepRxTx = oldRxTx
	}()

	mymock := mockintf{}
	// all apis are set to not fail
	RegisterClients(mymock)

	vxlanConfig := &VxlanConfig{
		Enable: true,
		VNI:    100,
		VlanId: []uint16{200},
	}

	vtepConfig := &VtepConfig{
		Enable:      true,
		Vni:         100,
		VtepName:    "vtep100",
		SrcIfName:   "eth0",
		TunnelDstIp: net.ParseIP("100.1.1.2"),
		VlanId:      200,
	}

	// create vxlan
	CreateVxLAN(vxlanConfig)
	<-vxlancreatedone

	// create vtep
	CreateVtep(vtepConfig)

	// should only be one entry
	key := &VtepDbKey{
		Name:  vtepConfig.VtepName,
		Vni:   vtepConfig.Vni,
		DstIp: vtepConfig.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)

	exitchan := make(chan bool, 1)
	go TimerTest(vtep, exitchan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-vtepcreatedone:
				return
			case <-exitchan:
				return
			}
		}
	}()

	wg.Wait()

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateStart {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateStrMap[VxlanVtepStateStart], VxlanVtepStateStrMap[vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState()])
	}

	DeRegisterClients()
	mymock.failGetNextHop = true
	RegisterClients(mymock)

	// simulate link down which should cause RIB to signal reachability message
	ifindex := GetIfIndexFromIfName("eth0")
	ProcessLinkDownCb(ifindex)

	dip := net.ParseIP("100.1.1.2")
	nexthopip := net.ParseIP("100.1.1.2")
	nexthopIfIndex := ifindex
	nexthopIfName := "eth0"
	reachable := false

	// now simulate the not reachabile message
	HandleNextHopChange(dip, nexthopip, nexthopIfIndex, nexthopIfName, reachable)
	<-vtepdeletedone
	fmt.Println("Return from HandleNextHopeChange")
	waitevt := make(chan bool)

	go func() {
		var x = 0
		for x = 0; x < 5 && vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateNextHopInfo; x++ {
			fmt.Println("waiting	", VxlanVtepStateStrMap[vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState()])
			time.Sleep(time.Millisecond * 10)
		}
		waitevt <- true
	}()
	time.Sleep(time.Second * 4)
	<-waitevt

	if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() != VxlanVtepStateNextHopInfo {
		t.Errorf("State not as expected expected[%s] actual[%s]", VxlanVtepStateStrMap[VxlanVtepStateNextHopInfo], VxlanVtepStateStrMap[vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState()])
	}

	DeleteVtep(vtepConfig)

	<-vtepdeletedone

	DeleteVxLAN(vxlanConfig, false)

	<-vxlandeletedone

	if len(GetVxlanDB()) != 0 {
		t.Errorf("Vxlan db not empty as expected")
	}

	if len(GetVtepDB()) != 0 {
		t.Errorf("Vtep db not empty as expected")
	}
}
