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

package stp

import (
	"net"
	"testing"
	"time"
)

func StpBridgeConfigSetup() *StpBridgeConfig {
	brg := &StpBridgeConfig{
		Address:      "00:11:22:33:44:55",
		Priority:     32768,
		MaxAge:       10,
		HelloTime:    1,
		ForwardDelay: 6,
		ForceVersion: 2, // RSTP
		TxHoldCount:  2,
		Vlan:         1,
	}
	return brg
}

func MemoryCheck(t *testing.T) {
	if len(PortListTable) != 0 {
		t.Error("Error cleaning of StpPort DB")
	}

	if len(PortMapTable) != 0 {
		t.Error("Error cleaning of StpPort MAP DB")
	}

	if len(BridgeMapTable) != 0 {
		t.Error("Error cleaning of Bridge MAP table DB")
	}

	if len(BridgeListTable) != 0 {
		t.Error("Error cleaning of Bridge List table DB")
	}

}

func StpPortConfigSetup(createbridge, update bool) (*StpPortConfig, *StpBridgeConfig) {

	var brg *StpBridgeConfig
	if createbridge {
		brg = StpBridgeConfigSetup()
		// bridge must exist
		StpBridgeCreate(brg)
	}

	p := &StpPortConfig{
		IfIndex:           1,
		Priority:          128,
		Enable:            true,
		PathCost:          200000,
		ProtocolMigration: 1,
		AdminPointToPoint: 0,
		AdminEdgePort:     false,
		AdminPathCost:     200000,
		BrgIfIndex:        1,
		BridgeAssurance:   false,
		BpduGuard:         false,
		BpduGuardInterval: 0,
	}

	PortConfigMap[p.IfIndex] = portConfig{Name: "lo",
		HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x11, 0x22, 0x22, 0x33},
	}

	// set dummy mac
	SaveSwitchMac("00:11:22:33:44:55")

	return p, brg
}

func TestStpBridgeCreationDeletion(t *testing.T) {
	defer MemoryCheck(t)
	// when creating a bridge it is recommended that the following calls are made

	// this is specific to test but config object should be filled in
	brgcfg := StpBridgeConfigSetup()
	// verify the paramaters are correct
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err != nil {
		t.Error("ERROR valid config failed", err)
	}
	// create the bridge
	err = StpBridgeCreate(brgcfg)
	if err != nil {
		t.Error("ERROR valid creation failed", err)
	}
	// delete the bridge
	// config cleanup is done as part of bridge delete
	err = StpBridgeDelete(brgcfg)
	if err != nil {
		t.Error("ERROR valid deletion failed", err)
	}
}

func TestStpPortCreationDeletion(t *testing.T) {
	defer MemoryCheck(t)
	// when creating a bridge it is recommended that the following calls are made

	// this is specific to test but config object should be filled in
	brgcfg := StpBridgeConfigSetup()
	// verify the paramaters are correct
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err != nil {
		t.Error("ERROR valid brg config failed", err)
	}
	// create the bridge
	err = StpBridgeCreate(brgcfg)
	if err != nil {
		t.Error("ERROR valid brg creation failed", err)
	}

	// configure valid stp port
	pcfg, _ := StpPortConfigSetup(false, false)
	pcfg.BrgIfIndex = 1
	pcfg.BridgeAssurance = true

	// valid config
	err = StpPortConfigParamCheck(pcfg, false, true)
	if err != nil {
		t.Error("ERROR valid stp port config failed", err)
	}

	// create the stp port
	err = StpPortCreate(pcfg)
	if err != nil {
		t.Error("ERROR valid stp port creation failed", err)
	}

	var p *StpPort
	if !StpFindPortByIfIndex(pcfg.IfIndex, pcfg.BrgIfIndex, &p) {
		t.Errorf("ERROR unable to find bridge port that was just created", err)

	}
	waitChan := make(chan bool)
	go func() {
		for i := 0; i < 30; i++ {
			if p.Forwarding == false || p.Learning == false {
				time.Sleep(time.Second * 1)
			}
		}
		waitChan <- true
	}()
	<-waitChan

	if p.Forwarding == false || p.Learning == false {
		t.Error("ERROR Bridge Port did not come up in a defaulted learning/forwarding state")
	}

	// delete the stp port
	err = StpPortDelete(pcfg)
	if err != nil {
		t.Error("ERROR valid stp port deletion failed", err)
	}

	if StpFindPortByIfIndex(pcfg.IfIndex, pcfg.BrgIfIndex, &p) {
		t.Errorf("ERROR found bridge port that was just deleted", err)

	}

	// delete the bridge
	// config cleanup is done as part of bridge delete
	err = StpBridgeDelete(brgcfg)
	if err != nil {
		t.Error("ERROR valid deletion failed", err)
	}
}

func TestStpPortAdminEdgeCreationDeletion(t *testing.T) {
	defer MemoryCheck(t)
	// when creating a bridge it is recommended that the following calls are made

	// this is specific to test but config object should be filled in
	brgcfg := StpBridgeConfigSetup()
	// verify the paramaters are correct
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err != nil {
		t.Error("ERROR valid brg config failed", err)
	}
	// create the bridge
	err = StpBridgeCreate(brgcfg)
	if err != nil {
		t.Error("ERROR valid brg creation failed", err)
	}

	// configure valid stp port
	pcfg, _ := StpPortConfigSetup(false, false)
	pcfg.BrgIfIndex = 1
	pcfg.AdminEdgePort = true
	pcfg.BpduGuard = true

	// valid config
	err = StpPortConfigParamCheck(pcfg, false, true)
	if err != nil {
		t.Error("ERROR valid stp port config failed", err)
	}

	// create the stp port
	err = StpPortCreate(pcfg)
	if err != nil {
		t.Error("ERROR valid stp port creation failed", err)
	}

	var p *StpPort
	if !StpFindPortByIfIndex(pcfg.IfIndex, pcfg.BrgIfIndex, &p) {
		t.Errorf("ERROR unable to find bridge port that was just created", err)

	}
	waitChan := make(chan bool)
	go func() {
		for i := 0; i < 30; i++ {
			if p.Forwarding == false || p.Learning == false {
				time.Sleep(time.Second * 1)
			}
		}
		waitChan <- true
	}()
	<-waitChan

	if p.Forwarding == false || p.Learning == false {
		t.Error("ERROR Bridge Port did not come up in a defaulted learning/forwarding state")
	}

	// delete the stp port
	err = StpPortDelete(pcfg)
	if err != nil {
		t.Error("ERROR valid stp port deletion failed", err)
	}

	if StpFindPortByIfIndex(pcfg.IfIndex, pcfg.BrgIfIndex, &p) {
		t.Errorf("ERROR found bridge port that was just deleted", err)

	}

	// delete the bridge
	// config cleanup is done as part of bridge delete
	err = StpBridgeDelete(brgcfg)
	if err != nil {
		t.Error("ERROR valid deletion failed", err)
	}
}

func TestStpBridgeParamCheckPriority(t *testing.T) {
	defer MemoryCheck(t)
	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.Priority = 11111
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid priority was set should have errored", brgcfg.Priority)
	}

	// now lets send a good value according to table 802.1D 17-2
	// 0 - 61440 in increments of 4096
	for i := uint16(0); i <= 61440/4096; i++ {
		brgcfg.Priority = 4096 * i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid priority was set should not have errored", brgcfg.Priority, err)
		}
	}

	// lets create the bridge
	StpBridgeCreate(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// lets update the bridge priority attribute
	err = StpBrgPrioritySet(b.BrgIfIndex, 4096)
	if err != nil {
		t.Error("ERRROR Setting bridge priority to a valid value", err)
	}
	prio := GetBridgePriorityFromBridgeId(b.BridgeIdentifier)
	if prio != (4096 | brgcfg.Vlan) {
		t.Error("ERROR Bridge Priority not set properly in packet", prio, b.BridgeIdentifier)
	}

	// lets update the bridge priority attribute to an invalid value
	err = StpBrgPrioritySet(b.BrgIfIndex, 400)
	if err == nil {
		t.Error("ERRROR Setting bridge priority to an invalid value", err)
	}

	StpBridgeDelete(brgcfg)
}

func TestStpBridgeParamCheckMaxAge(t *testing.T) {
	defer MemoryCheck(t)
	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.MaxAge = 200
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid max age was set should have errored", brgcfg.MaxAge)
	}
	// set bad value
	brgcfg.MaxAge = 5
	err = StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid max age was set should have errored", brgcfg.MaxAge)
	}

	// now lets send all possible values according to table 802.1D 17-1
	for i := uint16(6); i <= 40; i++ {
		brgcfg.MaxAge = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid max age was set should not have errored", brgcfg.MaxAge)
		}
	}

	// create the bridge
	StpBridgeCreate(brgcfg)
	defer StpBridgeDelete(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// let change the max age
	err = StpBrgMaxAgeSet(b.BrgIfIndex, 15)
	if err != nil {
		t.Error("ERROR: Valid hello time set failed", err)
	}

	// let change the Max age to an invalid value
	err = StpBrgMaxAgeSet(b.BrgIfIndex, 100)
	if err == nil {
		t.Error("ERROR: invalid hello time passed", err)
	}
}

func TestStpBridgeParamCheckHelloTime(t *testing.T) {
	defer MemoryCheck(t)
	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.HelloTime = 0
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid hello time was set should have errored", brgcfg.HelloTime)
	}
	// set bad value
	brgcfg.HelloTime = 5
	err = StpBrgConfigParamCheck(brgcfg, false)
	if err == nil {
		t.Error("ERROR an invalid hello time was set should have errored", brgcfg.HelloTime)
	}

	// now lets send all possible values according to table 802.1D 17-1
	for i := uint16(1); i <= 2; i++ {
		brgcfg.HelloTime = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid hello time was set should not have errored", brgcfg.HelloTime)
		}
	}

	// create the bridge
	StpBridgeCreate(brgcfg)
	defer StpBridgeDelete(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// let change the forwarding delay
	err = StpBrgHelloTimeSet(b.BrgIfIndex, 1)
	if err != nil {
		t.Error("ERROR: Valid hello time set failed", err)
	}

	// let change the forwarding delay to an invalid value
	err = StpBrgHelloTimeSet(b.BrgIfIndex, 5)
	if err == nil {
		t.Error("ERROR: invalid hello time passed", err)
	}
}

func TestStpBridgeParamCheckFowardingDelay(t *testing.T) {
	defer MemoryCheck(t)
	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.ForwardDelay = 0
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid forwarding delay was set should have errored", brgcfg.ForwardDelay)
	}
	// set bad value
	brgcfg.ForwardDelay = 50
	err = StpBrgConfigParamCheck(brgcfg, false)
	if err == nil {
		t.Error("ERROR an invalid forwardng delay was set should have errored", brgcfg.ForwardDelay)
	}

	// now lets send all possible values according to table 802.1D 17-1
	for i := uint16(4); i <= 30; i++ {
		brgcfg.ForwardDelay = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid forwarding delay was set should not have errored", brgcfg.ForwardDelay)
		}
	}

	// create the bridge
	StpBridgeCreate(brgcfg)
	defer StpBridgeDelete(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// let change the forwarding delay
	err = StpBrgForwardDelaySet(b.BrgIfIndex, 6)
	if err != nil {
		t.Error("ERROR: Valid version set failed", err)
	}

	// let change the forwarding delay to an invalid value
	err = StpBrgForwardDelaySet(b.BrgIfIndex, 50)
	if err == nil {
		t.Error("ERROR: invalid version set passed", err)
	}
}

func TestStpBridgeParamCheckTxHoldCount(t *testing.T) {
	defer MemoryCheck(t)
	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.TxHoldCount = 0
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid tx hold count was set should have errored", brgcfg.TxHoldCount)
	}
	// set bad value
	brgcfg.TxHoldCount = 50
	err = StpBrgConfigParamCheck(brgcfg, false)
	if err == nil {
		t.Error("ERROR an invalid tx hold count was set should have errored", brgcfg.TxHoldCount)
	}

	// now lets send all possible values according to table 802.1D 17-1
	for i := int32(1); i <= 10; i++ {
		brgcfg.TxHoldCount = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid tx hold count was set should not have errored", brgcfg.TxHoldCount)
		}
	}

	// create the bridge
	StpBridgeCreate(brgcfg)
	defer StpBridgeDelete(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// let change the tx hold count
	err = StpBrgTxHoldCountSet(b.BrgIfIndex, 4)
	if err != nil {
		t.Error("ERROR: Valid tx hold count set failed", err)
	}

	// let change the tx hold count to an invalid value
	err = StpBrgTxHoldCountSet(b.BrgIfIndex, 100)
	if err == nil {
		t.Error("ERROR: invalid tx hold count passed", err)
	}

}

func TestStpBridgeParamCheckVlan(t *testing.T) {
	defer MemoryCheck(t)

	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.Vlan = 4096
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid vlan was set should have errored", brgcfg.Vlan)
	}

	// now lets send all possible values according to table 802.1Q
	for i := uint16(1); i <= 4094; i++ {
		brgcfg.Vlan = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid vlan was set should not have errored", brgcfg.Vlan)
		}
		StpBrgConfigDelete(int32(brgcfg.Vlan))
	}
}

func TestStpBridgeParamCheckForceVersion(t *testing.T) {
	defer MemoryCheck(t)

	// setup
	brgcfg := StpBridgeConfigSetup()

	// set bad value
	brgcfg.ForceVersion = 0
	err := StpBrgConfigParamCheck(brgcfg, true)
	if err == nil {
		t.Error("ERROR an invalid force version was set should have errored", brgcfg.ForceVersion)
	}

	brgcfg.ForceVersion = 3
	err = StpBrgConfigParamCheck(brgcfg, false)
	if err == nil {
		t.Error("ERROR an invalid force version was set should have errored", brgcfg.ForceVersion)
	}

	// now lets send all possible values according to 802.1D
	for i := int32(1); i <= 2; i++ {
		brgcfg.ForceVersion = i
		err = StpBrgConfigParamCheck(brgcfg, false)
		if err != nil {
			t.Error("ERROR valid force version was set should not have errored", brgcfg.ForceVersion)
		}
	}

	// create the bridge
	StpBridgeCreate(brgcfg)
	defer StpBridgeDelete(brgcfg)

	var b *Bridge
	key := BridgeKey{
		Vlan: brgcfg.Vlan,
	}
	if !StpFindBridgeById(key, &b) {
		t.Error("ERROR: unable to find bridge")
	}

	// let change the version
	err = StpBrgForceVersion(b.BrgIfIndex, 1)
	if err != nil {
		t.Error("ERROR: Valid version set failed", err)
	}

	// let change the version
	err = StpBrgForceVersion(b.BrgIfIndex, 2)
	if err != nil {
		t.Error("ERROR: Valid version set failed", err)
	}

	// let change the version
	err = StpBrgForceVersion(b.BrgIfIndex, 10)
	if err == nil {
		t.Error("ERROR: invalid version set passed", err)
	}
}

/*
type StpPortConfig struct {
	IfIndex           int32
	Priority          uint16
	Enable            bool
	PathCost          int32
	ProtocolMigration int32
	AdminPointToPoint int32
	AdminEdgePort     bool
	AdminPathCost     int32
	BrgIfIndex        int32
	BridgeAssurance   bool
	BpduGuard         bool
	BpduGuardInterval int32
}
*/

func TestStpPortParamBrgIfIndex(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(false, false)
	defer StpPortConfigDelete(p.IfIndex)

	// port created without bridge being created
	err := StpPortConfigParamCheck(p, true, true)
	if err == nil {
		t.Error("ERROR: an invalid brgIfndex was set should have errored", p.BrgIfIndex, err)
	}

	// bridge created, and port created with invalid bridge reference
	p, b = StpPortConfigSetup(true, false)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	p.BrgIfIndex = 1000
	err = StpPortConfigParamCheck(p, true, true)
	if err == nil {
		t.Error("ERROR: an invalid brgIfndex was set should have errored", p.BrgIfIndex, err)
	}

	p, b = StpPortConfigSetup(false, false)
	err = StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: an valid brgIfndex was set should not have errored", p.BrgIfIndex, err)
	}
}

func TestStpPortParamPriority(t *testing.T) {
	defer MemoryCheck(t)

	p, b := StpPortConfigSetup(true, false)

	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	// lets create the bridge port so that we can try and update it later
	StpPortCreate(p)
	defer StpPortDelete(p)

	// invalid value
	p.Priority = 241
	err := StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid priority was set should have errored", p.Priority)
	}
	p.Priority = 1
	err = StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid priority was set should have errored", p.Priority, err)
	}

	// valid values according to 802.1D table 17-2
	for i := uint16(0); i <= 240/16; i++ {
		p.Priority = 16 * i
		err := StpPortConfigParamCheck(p, true, false)
		if err != nil {
			t.Error("ERROR: valid priority was set should not have errored", p.Priority, err)
		}
	}

	//brgifindex := p.BrgIfIndex
	// lets pretend another bridge port is being created and Priority is different
	brg := StpBridgeConfigSetup()
	brg.Vlan = 100
	// bridge must exist
	StpBridgeCreate(brg)
	defer StpBridgeDelete(brg)
	/*
		Invalid test as we are no longer bound to a port based provisioning
		p.BrgIfIndex = 100
		p.Priority = 16
		err = StpPortConfigParamCheck(p, false)
		if err == nil {
			t.Error("ERROR: an invalid port config change priority was set should have errored", p.Priority, err)
		}
		p.BrgIfIndex = brgifindex
	*/
	// lets change the port priority on the fly
	err = StpPortPrioritySet(p.IfIndex, p.BrgIfIndex, 32)
	p.Priority = 32
	err = StpPortConfigParamCheck(p, false, false)
	if err != nil {
		t.Error("ERROR: set a valid port priority 32 should not have failed ", err)
	}

	// set an invalid port priority
	err = StpPortPrioritySet(p.IfIndex, p.BrgIfIndex, 50)
	p.Priority = 50
	err = StpPortConfigParamCheck(p, false, false)
	if err == nil {
		t.Error("ERROR: set an ivalid port priority 50 should have failed", err)
	}

	// give test time to complete
	time.Sleep(time.Millisecond * 10)

}

func TestStpPortParamAdminPathCost(t *testing.T) {
	defer MemoryCheck(t)

	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	StpPortCreate(p)
	defer StpPortDelete(p)

	// invalid value
	p.AdminPathCost = 200000001
	err := StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid admin path cost was set should have errored", p.AdminPathCost, err)
	}
	p.AdminPathCost = 0
	err = StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid admin path cost was set should have errored", p.AdminPathCost, err)
	}

	// valid recommended values according to 802.1D table 17-3
	// technically the range is 1-2000000000
	for i := int32(1); i <= 9; i++ {
		p.AdminPathCost = 2 * (i * 10)
		err := StpPortConfigParamCheck(p, true, false)
		if err != nil {
			t.Error("ERROR: valid admin path cost was set should not have errored", p.AdminPathCost, err)
		}
	}
	StpPortConfigSave(p, false)

	// lets pretend another bridge port is being created and Priority is different
	brg := StpBridgeConfigSetup()
	brg.Vlan = 100
	// bridge must exist
	StpBridgeCreate(brg)
	defer StpBridgeDelete(brg)
	/*
		Invalid test as we are no longer bound to port based provisioning
		p.BrgIfIndex = 100
		p.AdminPathCost = 200
		err = StpPortConfigParamCheck(p, false)
		if err == nil {
			t.Error("ERROR: an invalid port config change admin path cost was set should have errored", p.AdminPathCost, err)
		}
	*/
}

func TestStpPortParamBridgeAssurance(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	// lets create the bridge port so that we can try and update it later
	StpPortCreate(p)
	defer StpPortDelete(p)

	// bridge assurance is only valid on an non edge port
	p.AdminEdgePort = true
	p.BridgeAssurance = true
	err := StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid port config Admin Edge and Bridge Assurance set should have errored", p.AdminEdgePort, p.BridgeAssurance, err)
	}

	// bridge assurance is only valid on an non edge port
	p.AdminEdgePort = false
	p.BridgeAssurance = true
	err = StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: valid port config Bridge Assurance set should not have errored", p.AdminEdgePort, p.BridgeAssurance, err)
	}

	// lets save the config
	StpPortConfigSave(p, false)

	ifIndex := p.IfIndex
	brgIfIndex := p.BrgIfIndex
	// lets update Admin Edge when Bridge Assurance is already enabled
	p = StpPortConfigGet(ifIndex)
	if p == nil {
		t.Error("ERROR: could not find port config")
	}
	p.IfIndex = ifIndex
	p.BrgIfIndex = brgIfIndex
	p.AdminEdgePort = true
	err = StpPortConfigParamCheck(p, true, true)
	if err == nil {
		t.Error("ERROR: invalid port config Bridge Assurance set should have errored", p.AdminEdgePort, p.BridgeAssurance, err)
	}

	// set admin edge to false so that we can set bridge assurance
	err = StpPortAdminEdgeSet(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: failed to set port as an admin edge port", err)
	}

	// lets set hte bridge assurance on a non-admin edge port
	err = StpPortBridgeAssuranceSet(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: failed to set Bridge Assurance on port", err)
	}

	// lets set hte bridge assurance on a non-admin edge port
	err = StpPortBridgeAssuranceSet(p.IfIndex, p.BrgIfIndex, true)
	if err != nil {
		t.Error("ERROR: failed to set Bridge Assurance on port", err)
	}

	// lets set hte bridge assurance on a non-admin edge port again to ensure that
	// nothing changed
	err = StpPortBridgeAssuranceSet(p.IfIndex, p.BrgIfIndex, true)
	if err != nil {
		t.Error("ERROR: failed to set Bridge Assurance on port", err)
	}

	var port *StpPort
	if !StpFindPortByIfIndex(p.IfIndex, p.BrgIfIndex, &port) {
		t.Error("ERROR: did not find bridge port")
	}

	if !port.BridgeAssurance {
		t.Error("ERROR: Why was bridge assurance not set in db record")
	}

	// set admin edge to true while bridge assurance is enabled, should fail
	err = StpPortAdminEdgeSet(p.IfIndex, p.BrgIfIndex, true)
	p.BridgeAssurance = true
	err = StpPortConfigParamCheck(p, false, false)
	if err == nil {
		t.Error("ERROR: failed to set port as an admin edge port because Bridge Assurance is enabled", err)
	}

	// lets disable bridge assurance
	err = StpPortBridgeAssuranceSet(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: failed to clear Bridge Assurance on port", err)
	}

	if !StpFindPortByIfIndex(p.IfIndex, p.BrgIfIndex, &port) {
		t.Error("ERROR: did not find bridge port")
	}

	if port.BridgeAssurance {
		t.Error("ERROR: Why was bridge assurance set in db record")
	}

	// set admin edge to true while bridge assurance is disabled
	err = StpPortAdminEdgeSet(p.IfIndex, p.BrgIfIndex, true)
	if err != nil {
		t.Error("ERROR: failed to set port as an admin edge port because Bridge Assurance is enabled", err)
	}

	// give test time to complete
	time.Sleep(time.Millisecond * 10)
}

func TestStpPortParamBpduGuard(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}
	StpPortCreate(p)
	defer StpPortDelete(p)

	// bpdu guard is only valid on an edge port
	p.AdminEdgePort = false
	p.BpduGuard = true
	err := StpPortConfigParamCheck(p, true, false)
	if err == nil {
		t.Error("ERROR: an invalid port config NOT Admin Edge and Bpdu Guard set should have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// bpdu guard only valid on edge port
	p.AdminEdgePort = true
	p.BpduGuard = true
	err = StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: valid port config bpdu guard and admin state set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// lets save the config and create the port so we can play around with update
	StpPortConfigSave(p, false)

	brgIfIndex := p.BrgIfIndex
	// lets update Admin Edge when bpdu guard is already enabled
	p = StpPortConfigGet(p.IfIndex)
	if p == nil {
		t.Error("ERROR: could not find port config")
	}
	p.BrgIfIndex = brgIfIndex
	// disable admin edge, which should fail
	p.AdminEdgePort = false
	err = StpPortAdminEdgeSet(p.IfIndex, p.BrgIfIndex, p.AdminEdgePort)
	err = StpPortConfigParamCheck(p, false, false)
	if err == nil {
		t.Error("ERROR: invalid port config bpdu Guard is enabled set should have errored", p.AdminEdgePort, p.BpduGuard, err)
	}
	// reset to good state
	p.AdminEdgePort = true
	p.BpduGuard = false
	// lets disable bpdu guard
	err = StpPortBpduGuardSet(p.IfIndex, p.BrgIfIndex, p.BpduGuard)
	if err != nil {
		t.Error("ERROR: valid port config bdpu guard being unset set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}
	// reset to good state
	p.BpduGuard = true
	// re-enabled bpdu guard
	err = StpPortBpduGuardSet(p.IfIndex, p.BrgIfIndex, p.BpduGuard)
	if err != nil {
		t.Error("ERROR: valid port config bdpu guard being set set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// set again to ensure that nothing changed (code coverage call)
	p.BpduGuard = true
	// re-enabled bpdu guard
	err = StpPortBpduGuardSet(p.IfIndex, p.BrgIfIndex, p.BpduGuard)
	if err != nil {
		t.Error("ERROR: valid port config bdpu guard being set set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// disable bpdu guard
	p.BpduGuard = false
	err = StpPortBpduGuardSet(p.IfIndex, p.BrgIfIndex, p.BpduGuard)
	if err != nil {
		t.Error("ERROR: valid port config bdpu guard being unset set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// disable admin edge
	p.AdminEdgePort = false
	err = StpPortAdminEdgeSet(p.IfIndex, p.BrgIfIndex, p.AdminEdgePort)
	if err != nil {
		t.Error("ERROR: valid port config bpdu Guard is no longer enabled set should not have errored", p.AdminEdgePort, p.BpduGuard, err)
	}

	// give test time to complete
	time.Sleep(time.Millisecond * 30)

}

func TestStpPortParamProtocolMigration(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	// protocol migration is a state machine check attribute
	p.ProtocolMigration = 1
	err := StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: invalid port config protocol migration not set", err)
	}

	// protocol migration is a state machine check attribute
	p.ProtocolMigration = 0
	err = StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: invalid port config protocol migration not set", err)
	}

	// lets save the config and create the port so we can play around with update
	StpPortConfigSave(p, false)
	StpPortCreate(p)
	defer StpPortDelete(p)

	brgIfIndex := p.BrgIfIndex
	// lets update Admin Edge when bpdu guard is already enabled
	p = StpPortConfigGet(p.IfIndex)
	if p == nil {
		t.Error("ERROR: could not find port config")
	}
	p.BrgIfIndex = brgIfIndex
	// enable protocol migration
	p.ProtocolMigration = 1
	err = StpPortProtocolMigrationSet(p.IfIndex, p.BrgIfIndex, true)
	if err != nil {
		t.Error("ERROR: invalid port config protocol migration not set", err)
	}
	// set again to ensure that nothing changed (code coverage call)
	p.ProtocolMigration = 0
	// re-enabled bpdu guard
	err = StpPortProtocolMigrationSet(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: invalid port config protocol migration not set", err)
	}

	// disable protocol migration
	p.ProtocolMigration = 0
	// lets disable bpdu guard
	err = StpPortProtocolMigrationSet(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: invalid port config protocol migration not set", err)
	}
}

func TestStpPortPortEnable(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	// port enable
	p.Enable = true
	err := StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: invalid port config Enable not set", err)
	}
	StpPortConfigDelete(p.IfIndex)

	// port disable
	p.Enable = false
	err = StpPortConfigParamCheck(p, true, false)
	if err != nil {
		t.Error("ERROR: invalid port config Enable not cleared", err)
	}
	StpPortConfigDelete(p.IfIndex)

	p.Enable = true
	// lets save the config and create the port so we can play around with update
	StpPortConfigSave(p, false)
	err = StpPortCreate(p)
	if err != nil {
		t.Error("ERROR: port creation failed", err)
	}
	defer StpPortDelete(p)

	brgIfIndex := p.BrgIfIndex
	// lets update Admin Edge when bpdu guard is already enabled
	p = StpPortConfigGet(p.IfIndex)
	if p == nil {
		t.Error("ERROR: could not find port config")
	}
	p.BrgIfIndex = brgIfIndex
	// port disable
	p.Enable = false
	err = StpPortEnable(p.IfIndex, p.BrgIfIndex, p.Enable)
	if err != nil {
		t.Error("ERROR: invalid port config port enable not disabled", err)
	}

	// port enable
	p.Enable = true
	err = StpPortEnable(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: invalid port config port enable not enabled", err)
	}

	// enable again (code coverage)
	p.Enable = true
	err = StpPortEnable(p.IfIndex, p.BrgIfIndex, false)
	if err != nil {
		t.Error("ERROR: invalid port config port enable not enabled", err)
	}

	// lets give the test some time to complete
	// otherwise a crash may occur due to an event
	// being processed when the port id being deleted.
	// This is unlikely to happen from external user
	// however should be noted that a true solution
	// is needed
	time.Sleep(time.Millisecond * 10)

}

func TestStpPortLinkUpDown(t *testing.T) {
	defer MemoryCheck(t)
	p, b := StpPortConfigSetup(true, false)
	defer StpPortConfigDelete(p.IfIndex)
	if b != nil {
		defer StpBridgeDelete(b)
	}

	StpPortCreate(p)
	defer StpPortDelete(p)

	// Calls below made for purposes of code coverage
	// Could verify the state machines states but this
	// should be done as part of state machine testing

	// link up
	StpPortLinkUp(p.IfIndex)

	// link down
	StpPortLinkDown(p.IfIndex)

	// link up
	StpPortLinkUp(p.IfIndex)

	// link down
	StpPortLinkDown(p.IfIndex)

	// lets give the test some time to complete
	// otherwise a crash may occur due to an event
	// being processed when the port id being deleted.
	// This is unlikely to happen from external user
	// however should be noted that a true solution
	// is needed
	time.Sleep(time.Millisecond * 10)

}
