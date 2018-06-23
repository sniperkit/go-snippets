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

// config.go
package stp

import (
	"errors"
	"fmt"
	"math"
)

// StpBridgeConfig config data
type StpBridgeConfig struct {
	IfIndex      int32
	Address      string
	Priority     uint16
	MaxAge       uint16
	HelloTime    uint16
	ForwardDelay uint16
	ForceVersion int32
	TxHoldCount  int32
	Vlan         uint16
	DebugLevel   int
}

// StpPortConfig config data
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

// store the port config for each port
var StpPortConfigMap map[int32]StpPortConfig

// store the bridge config for each bridge
var StpBridgeConfigMap map[int32]StpBridgeConfig

// StpPortConfigGet gets the current port config
func StpPortConfigGet(pId int32) *StpPortConfig {
	c, ok := StpPortConfigMap[pId]
	if ok {
		return &c
	}
	return nil
}

// StpBrgConfigGet gets the current bridge config
func StpBrgConfigGet(bId int32) *StpBridgeConfig {
	c, ok := StpBridgeConfigMap[bId]
	if ok {
		return &c
	}
	return nil
}

func StpBrgConfigDelete(bId int32) error {
	if StpBrgConfigGet(bId) != nil {
		delete(StpBridgeConfigMap, bId)
		return nil
	}
	return errors.New(fmt.Sprintf("Error Trying to Delete Bridge %d Config that does not exist", bId))
}

// StpPortConfigSave Save the last config given by user this is a validation
// check as well so that all port contain the same config
func StpPortConfigSave(c *StpPortConfig, update bool) error {
	brgIfIndex := c.BrgIfIndex
	c.BrgIfIndex = 0
	if _, ok := StpPortConfigMap[c.IfIndex]; !ok {
		//fmt.Println("Saving Port Config", c.IfIndex)
		StpPortConfigMap[c.IfIndex] = *c
	} else {
		if !update && *c != StpPortConfigMap[c.IfIndex] {
			c.BrgIfIndex = brgIfIndex
			// TODO failing for now will need to add code to update all other bridges that use
			// this physical port
			//return errors.New(fmt.Sprintf("Error Port %d Provisioning does not agree with previously created bridge port prev[%#v] new[%#v]",
			//	c.IfIndex, StpPortConfigMap[c.IfIndex], *c))
		} else if update {
			StpPortConfigMap[c.IfIndex] = *c
		}
	}
	c.BrgIfIndex = brgIfIndex
	return nil
}

// StpPortConfigSave Save the last config given by user this is a validation
// check as well so that all port contain the same config
func StpBrgConfigSave(c *StpBridgeConfig) error {
	//fmt.Println("Saving Bridge Config", c.Vlan)
	StpBridgeConfigMap[int32(c.Vlan)] = *c
	return nil
}

// Delete the saved port configuration
func StpPortConfigDelete(Id int32) {
	if _, ok := StpPortConfigMap[Id]; ok {
		delete(StpPortConfigMap, Id)
	}
}

// StpBrgConfigParamCheck will validate the bridge config paramaters
func StpBrgConfigParamCheck(c *StpBridgeConfig, create bool) error {
	var b *Bridge

	if create {
		// Bridges are unique
		_, ok := StpBridgeConfigMap[int32(c.Vlan)]
		if StpFindBridgeByIfIndex(int32(c.Vlan), &b) || ok {
			fmt.Println(BridgeMapTable, StpFindBridgeByIfIndex(int32(c.Vlan), &b), StpBrgConfigGet(int32(c.Vlan)))
			if c.Vlan == DEFAULT_STP_BRIDGE_VLAN {
				errors.New(fmt.Sprintf("Invalid Config, Default Bridge %d already exists", c.Vlan))
			} else {
				return errors.New(fmt.Sprintf("Invalid Config, Bridge %d already exists", c.Vlan))
			}
		}
	}

	// Table 17-2 says the values can be 0-61140 in increments of 4096
	if math.Mod(float64(c.Priority), 4096) != 0 || c.Priority > 61440 {
		return errors.New(fmt.Sprintf("Invalid Bridge Priority %d valid values 0-61440 increments of 4096", c.Priority))
	}

	// valid values according to Table 17-1
	if c.MaxAge < 6 ||
		c.MaxAge > 40 {
		return errors.New(fmt.Sprintf("Invalid Bridge Max Age %d valid range 6.0 - 40.0", c.MaxAge))
	}

	if c.HelloTime < 1 ||
		c.HelloTime > 2 {
		return errors.New(fmt.Sprintf("Invalid Bridge Hello Time %d valid range 1.0 - 2.0", c.HelloTime))
	}

	if c.ForwardDelay < 3 ||
		c.ForwardDelay > 30 {
		return errors.New(fmt.Sprintf("Invalid Bridge Hello Time %d valid range 3.0 - 30.0", c.ForwardDelay))
	}

	// 1 == STP
	// 2 == RSTP
	// 3 == MSTP currently not support
	if c.ForceVersion != 1 &&
		c.ForceVersion != 2 {
		return errors.New(fmt.Sprintf("Invalid Bridge Force Version %d valid 1 (STP) 2 (RSTP)", c.ForceVersion))
	}

	if c.TxHoldCount < 1 ||
		c.TxHoldCount > 10 {
		return errors.New(fmt.Sprintf("Invalid Bridge Tx Hold Count %d valid range 1 - 10", c.TxHoldCount))
	}

	// if zero is used then we will convert this to use default
	if c.Vlan != 0 {
		if c.Vlan > 4095 {
			return errors.New(fmt.Sprintf("Invalid Bridge Vlan %d valid range 1 - 4095", c.Vlan))
		}
	}

	// lets store the configuration
	return StpBrgConfigSave(c)
}

// StpPortConfigParamCheck will validate the config paramater for a bridge port
func StpPortConfigParamCheck(c *StpPortConfig, update bool, create bool) error {
	var b *Bridge

	// bridge must be valid for a bridge port to be created
	if !StpFindBridgeByIfIndex(c.BrgIfIndex, &b) && StpBrgConfigGet(c.BrgIfIndex) == nil {
		return errors.New(fmt.Sprintf("Invalid BrgIfIndex %d, port %d must be associated with valid bridge", c.BrgIfIndex, c.IfIndex))
	}

	var p *StpPort
	if create &&
		StpFindPortByIfIndex(c.IfIndex, c.BrgIfIndex, &p) {
		return errors.New(fmt.Sprintf("ERROR Stp Port already created vlan %d port %d", c.BrgIfIndex, c.IfIndex))
	}

	// Table 17-2
	if math.Mod(float64(c.Priority), 16) != 0 || c.Priority > 240 {
		return errors.New(fmt.Sprintf("Invalid Port %d Priority %d valid values 0-240 increments of 16", c.IfIndex, c.Priority))
	}

	// Table 17-3
	if c.AdminPathCost < 1 || c.AdminPathCost > 200000000 {
		return errors.New(fmt.Sprintf("Invalid Port %d Path Cost %d valid values 0 (AUTO) or 1 - 200,000,000", c.IfIndex, c.AdminPathCost))
	}

	if (!c.AdminEdgePort) &&
		c.BpduGuard {
		return errors.New(fmt.Sprintf("Invalid Port %d Bpdu Guard only available on Edge Ports", c.IfIndex))
	}

	if (c.AdminEdgePort) &&
		c.BridgeAssurance {
		return errors.New(fmt.Sprintf("Invalid Port %d Bridge Assurance only available on non Edge Ports", c.IfIndex))
	}

	// all bridge port configurations are applied against all bridge ports applied to a given
	// port, updates are applied to all bridge ports
	// 9/20/16 relaxing this restriction as users will not know this
	/*
		if !update {
			brgifindex := c.BrgIfIndex
			c.BrgIfIndex = 0
			if _, ok := StpPortConfigMap[c.IfIndex]; ok {
				if *c != StpPortConfigMap[c.IfIndex] {
					return errors.New(fmt.Sprintf("Error Port %d Provisioning does not agree with previously created bridge port prev[%#v] new[%#v]",
						c.IfIndex, StpPortConfigMap[c.IfIndex], *c))
				}
			}
			c.BrgIfIndex = brgifindex
		}
	*/
	return StpPortConfigSave(c, update)
}

func StpBridgeCreate(c *StpBridgeConfig) error {
	var b *Bridge
	tmpaddr := c.Address
	if tmpaddr == "" {
		tmpaddr = "00:AA:AA:BB:BB:DD"
	}

	key := BridgeKey{
		Vlan: c.Vlan,
	}

	if !StpFindBridgeById(key, &b) {
		b = NewStpBridge(c)
		b.BEGIN(false)

	} else {
		return errors.New(fmt.Sprintf("Invalid config, bridge vlan %d already exists", c.Vlan))
	}
	return nil
}

func StpBridgeDelete(c *StpBridgeConfig) error {
	var b *Bridge

	key := BridgeKey{
		Vlan: c.Vlan,
	}
	if StpFindBridgeById(key, &b) {
		DelStpBridge(b, true)
		for _, btmp := range StpBridgeConfigMap {
			if btmp.Vlan == c.Vlan {
				StpBrgConfigDelete(int32(c.Vlan))
			}
		}
	} else {
		return errors.New(fmt.Sprintf("Invalid config, bridge vlan %d does not exists", c.Vlan))
	}
	return nil
}

func StpPortCreate(c *StpPortConfig) error {
	var p *StpPort
	var b *Bridge
	if !StpFindPortByIfIndex(c.IfIndex, c.BrgIfIndex, &p) {
		// lets store the configuration
		err := StpPortConfigSave(c, false)
		if err != nil {
			return err
		}

		// nothing should happen until a birdge is assigned to the port
		if StpFindBridgeByIfIndex(c.BrgIfIndex, &b) {
			p := NewStpPort(c)
			StpPortAddToBridge(p.IfIndex, p.BrgIfIndex)
		}
	} else {
		return errors.New(fmt.Sprintf("Invalid config, port %d bridge %d already exists", c.IfIndex, c.BrgIfIndex))
	}
	return nil
}

func StpPortDelete(c *StpPortConfig) error {
	var p *StpPort
	var b *Bridge
	if StpFindPortByIfIndex(c.IfIndex, c.BrgIfIndex, &p) {
		if StpFindBridgeByIfIndex(p.BrgIfIndex, &b) {
			StpPortDelFromBridge(c.IfIndex, p.BrgIfIndex)
		}
		DelStpPort(p)
		delete(StpPortConfigMap, c.IfIndex)
	} else {
		return errors.New(fmt.Sprintf("Invalid config, port %d bridge %d does not exists", c.IfIndex, c.BrgIfIndex))
	}
	return nil
}

func StpPortAddToBridge(pId int32, brgifindex int32) {
	var p *StpPort
	var b *Bridge
	if StpFindPortByIfIndex(pId, brgifindex, &p) && StpFindBridgeByIfIndex(brgifindex, &b) {
		p.BridgeId = b.BridgeIdentifier
		b.StpPorts = append(b.StpPorts, pId)
		p.BEGIN(false)

		if p.BdmMachineFsm != nil {
			// check all other bridge ports to see if any are AdminEdge
			isOtherBrgPortOperEdge := p.IsAdminEdgePort()
			if !p.AdminEdge &&
				isOtherBrgPortOperEdge {
				p.BdmMachineFsm.BdmEvents <- MachineEvent{
					e:   BdmEventBeginAdminEdge,
					src: "CONFIG: AdminEgeSet",
				}
			} else if p.AdminEdge && !isOtherBrgPortOperEdge {
				portDbMutex.Lock()
				for _, ptmp := range PortListTable {
					if p != ptmp {
						p.BdmMachineFsm.BdmEvents <- MachineEvent{
							e:   BdmEventBeginAdminEdge,
							src: "CONFIG: AdminEgeSet",
						}
					}
				}
				portDbMutex.Unlock()
			}
		}

	} else {
		StpLogger("ERROR", fmt.Sprintf("ERROR did not find bridge[%#v] or port[%d]", brgifindex, pId))
	}
}

func StpPortDelFromBridge(pId int32, brgifindex int32) {
	var p *StpPort
	var b *Bridge
	if StpFindPortByIfIndex(pId, brgifindex, &p) && StpFindBridgeByIfIndex(brgifindex, &b) {
		// lets disable the port before we remove it so that way
		// other ports can trigger tc event
		p.NotifyPortEnabled("CONFIG DEL", p.PortEnabled, false)
		p.PortEnabled = false
		// detach the port from the bridge stp port list
		for idx, ifindex := range b.StpPorts {
			if ifindex == p.IfIndex {
				b.StpPorts = append(b.StpPorts[:idx], b.StpPorts[idx+1:]...)
				break
			}
		}
	} else {
		StpLogger("ERROR", fmt.Sprintf("ERROR did not find bridge[%#v] or port[%d]", brgifindex, pId))
	}
}

func StpPortEnable(pId int32, bId int32, enable bool) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		if p.AdminPortEnabled != enable {
			if p.AdminPortEnabled {
				if p.PortEnabled {
					p.NotifyPortEnabled("CONFIG: ", p.PortEnabled, false)
					p.PortEnabled = false
				}
			} else {
				for _, client := range GetAsicDPluginList() {
					if client.GetPortLinkStatus(pId) {
						defer p.NotifyPortEnabled("CONFIG: ", p.PortEnabled, true)
						p.PortEnabled = true
					}
				}
			}
			p.AdminPortEnabled = enable
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Port Enable", pId, bId))
}

func StpPortLinkUp(pId int32) {
	portDbMutex.Lock()
	defer portDbMutex.Unlock()

	for _, p := range PortListTable {
		if p.IfIndex == pId {
			p.CreateRxTx()
			if p.AdminPortEnabled {
				defer p.NotifyPortEnabled("LINK EVENT", p.PortEnabled, true)
				p.PortEnabled = true
			}
		}
	}
}

func StpPortLinkDown(pId int32) {
	portDbMutex.Lock()
	defer portDbMutex.Unlock()

	for _, p := range PortListTable {
		if p.IfIndex == pId {
			p.DeleteRxTx()
			defer p.NotifyPortEnabled("LINK EVENT", p.PortEnabled, false)
			p.PortEnabled = false
		}
	}
}

func StpBrgPrioritySet(bId int32, priority uint16) error {
	// get bridge
	var b *Bridge
	var p *StpPort
	if StpFindBridgeByIfIndex(bId, &b) {
		prio := GetBridgePriorityFromBridgeId(b.BridgeIdentifier)
		if prio != priority {
			c := StpBrgConfigGet(bId)
			c.Priority = priority
			err := StpBrgConfigParamCheck(c, false)
			if err == nil {
				addr := GetBridgeAddrFromBridgeId(b.BridgeIdentifier)
				vlan := GetBridgeVlanFromBridgeId(b.BridgeIdentifier)
				b.BridgeIdentifier = CreateBridgeId(addr, priority, vlan)
				b.BridgePriority.DesignatedBridgeId = b.BridgeIdentifier

				for _, pId := range b.StpPorts {
					if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
						p.Selected = false
						p.Reselect = true
					}
				}
				if b.PrsMachineFsm != nil {
					b.PrsMachineFsm.PrsEvents <- MachineEvent{
						e:   PrsEventReselect,
						src: "CONFIG: BrgPrioritySet",
					}
				}
			}
			return err
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid bridge %d supplied for setting Priority", bId))
}

func StpBrgForceVersion(bId int32, version int32) error {

	var b *Bridge
	var p *StpPort
	if StpFindBridgeByIfIndex(bId, &b) {
		// version 1 STP
		// version 2 RSTP
		if b.ForceVersion != version {
			c := StpBrgConfigGet(bId)
			c.ForceVersion = version
			err := StpBrgConfigParamCheck(c, false)
			if err == nil {
				b.ForceVersion = version
				for _, pId := range b.StpPorts {
					if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
						if b.ForceVersion == 1 {
							p.RstpVersion = false
						} else {
							p.RstpVersion = true
						}
						p.BEGIN(true)
					}
				}
			}
			return err
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid bridge %d supplied for setting Force Version", bId))
}

// StpPortPrioritySet will set the priority on all bridge ports
func StpPortPrioritySet(pId int32, bId int32, priority uint16) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		if p.Priority != priority {
			/*
				c := StpPortConfigGet(pId)
				c.Priority = priority
				c.BrgIfIndex = bId
				err := StpPortConfigParamCheck(c, true)
				if err == nil {
					StpPortConfigSave(c, true)
					// apply to all bridge ports
					for _, port := range p.GetPortListToApplyConfigTo() {

						port.Priority = priority
						port.Selected = false
						port.Reselect = true

						if port.b.PrsMachineFsm != nil {
							port.b.PrsMachineFsm.PrsEvents <- MachineEvent{
								e:   PrsEventReselect,
								src: "CONFIG: PortPrioritySet",
							}
						}
					}
				}
			*/
			p.Priority = priority
			p.Selected = false
			p.Reselect = true
			if p.b.PrsMachineFsm != nil {
				p.b.PrsMachineFsm.PrsEvents <- MachineEvent{
					e:   PrsEventReselect,
					src: "CONFIG: PortPrioritySet",
				}
			}

			//return err
			return nil
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Port Priority", pId, bId))
}

// StpPortPortPathCostSet N/A
func StpPortPortPathCostSet(pId int32, bId int32, pathcost uint32) error {
	// TODO
	return nil
}

// StpPortAdminEdgeSet will set all bridge port as admin edge ports
func StpPortAdminEdgeSet(pId int32, bId int32, adminedge bool) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		p.AdminEdge = adminedge
		if p.OperEdge != adminedge {
			/*
				c := StpPortConfigGet(pId)
				c.AdminEdgePort = adminedge
				c.BrgIfIndex = bId
				err := StpPortConfigParamCheck(c, true)
				if err == nil {
					StpPortConfigSave(c, true)
					p.AdminEdge = adminedge
					isOtherBrgPortOperEdge := p.IsAdminEdgePort()
					// if we transition from Admin Edge to non-Admin edge
					if !p.AdminEdge && !isOtherBrgPortOperEdge {
						p.BdmMachineFsm.BdmEvents <- MachineEvent{
							e:   BdmEventBeginNotAdminEdge,
							src: "CONFIG: AdminEgeSet",
						}
						for _, ptmp := range PortListTable {
							if p != ptmp &&
								p.IfIndex == ptmp.IfIndex {
								p.BdmMachineFsm.BdmEvents <- MachineEvent{
									e:   BdmEventBeginNotAdminEdge,
									src: "CONFIG: AdminEgeSet",
								}
							}
						}

					} else if p.AdminEdge && !isOtherBrgPortOperEdge {
						p.BdmMachineFsm.BdmEvents <- MachineEvent{
							e:   BdmEventBeginAdminEdge,
							src: "CONFIG: AdminEgeSet",
						}

						for _, ptmp := range PortListTable {
							if p != ptmp &&
								p.IfIndex == ptmp.IfIndex {
								p.BdmMachineFsm.BdmEvents <- MachineEvent{
									e:   BdmEventBeginAdminEdge,
									src: "CONFIG: AdminEgeSet",
								}
							}
						}
					}
				}
			*/
			p.AdminEdge = adminedge
			if !p.AdminEdge {
				p.BdmMachineFsm.BdmEvents <- MachineEvent{
					e:   BdmEventBeginNotAdminEdge,
					src: "CONFIG: AdminEgeSet",
				}
			} else {
				p.BdmMachineFsm.BdmEvents <- MachineEvent{
					e:   BdmEventBeginAdminEdge,
					src: "CONFIG: AdminEgeSet",
				}
			}

			//return err
			return nil
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Port Admin Edge", pId, bId))
}

func StpBrgForwardDelaySet(bId int32, fwddelay uint16) error {
	var b *Bridge
	var p *StpPort
	if StpFindBridgeByIfIndex(bId, &b) {
		c := StpBrgConfigGet(bId)
		c.ForwardDelay = fwddelay
		err := StpBrgConfigParamCheck(c, false)
		if err == nil {
			b.BridgeTimes.ForwardingDelay = fwddelay
			// if we are root lets update the port times
			if b.RootPortId == 0 {
				b.RootTimes.ForwardingDelay = fwddelay
				for _, pId := range b.StpPorts {
					if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
						p.PortTimes.ForwardingDelay = b.RootTimes.ForwardingDelay
					}
				}
			}
		}
		return err
	}
	return errors.New(fmt.Sprintf("Invalid bridge %d supplied for setting Forwarding Delay", bId))
}

func StpBrgHelloTimeSet(bId int32, hellotime uint16) error {
	var b *Bridge
	var p *StpPort
	if StpFindBridgeByIfIndex(bId, &b) {
		c := StpBrgConfigGet(bId)
		c.HelloTime = hellotime
		err := StpBrgConfigParamCheck(c, false)
		if err == nil {
			b.BridgeTimes.HelloTime = hellotime
			// if we are root lets update the port times
			if b.RootPortId == 0 {
				b.RootTimes.HelloTime = hellotime
				for _, pId := range b.StpPorts {
					if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
						p.PortTimes.HelloTime = b.RootTimes.HelloTime
					}
				}
			}
		}
		return err
	}
	return errors.New(fmt.Sprintf("Invalid bridge %d supplied for setting Hello Time", bId))
}

func StpBrgMaxAgeSet(bId int32, maxage uint16) error {
	var b *Bridge
	var p *StpPort
	if StpFindBridgeByIfIndex(bId, &b) {
		c := StpBrgConfigGet(bId)
		c.MaxAge = maxage
		err := StpBrgConfigParamCheck(c, false)
		if err == nil {
			b.BridgeTimes.MaxAge = maxage
			// if we are root lets update the port times
			if b.RootPortId == 0 {
				b.RootTimes.MaxAge = maxage
				for _, pId := range b.StpPorts {
					if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
						p.PortTimes.MaxAge = b.RootTimes.MaxAge
					}
				}
			}
		}
		return err
	}
	return errors.New(fmt.Sprintf("Invalid bridge %d supplied for setting Max Age", bId))
}

func StpBrgTxHoldCountSet(bId int32, txholdcount uint16) error {
	var b *Bridge
	if StpFindBridgeByIfIndex(bId, &b) {
		c := StpBrgConfigGet(bId)
		c.TxHoldCount = int32(txholdcount)
		err := StpBrgConfigParamCheck(c, false)
		if err == nil {
			b.TxHoldCount = uint64(txholdcount)
		}
		return err
	}
	return nil
}

func StpPortProtocolMigrationSet(pId int32, bId int32, protocolmigration bool) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		if p.Mcheck != protocolmigration {
			/*
				c := StpPortConfigGet(pId)
				if protocolmigration {
					c.ProtocolMigration = int32(1)
				} else {
					c.ProtocolMigration = int32(0)
				}

				c.BrgIfIndex = bId
				err := StpPortConfigParamCheck(c, true)
				if err == nil {
					StpPortConfigSave(c, true)

					// apply to all bridge ports
					for _, port := range p.GetPortListToApplyConfigTo() {

						if protocolmigration {
							port.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMcheck,
								src: "CONFIG: ProtocolMigrationSet",
							}
						}
						port.Mcheck = protocolmigration
					}
				}
			*/
			if protocolmigration {
				p.PpmmMachineFsm.PpmmEvents <- MachineEvent{e: PpmmEventMcheck,
					src: "CONFIG: ProtocolMigrationSet",
				}
			}
			p.Mcheck = protocolmigration

			//return err
			return nil
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Protcol Migration", pId, bId))
}

func StpPortBpduGuardSet(pId int32, bId int32, bpduguard bool) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		if p.BpduGuard != bpduguard {
			/*
				c := StpPortConfigGet(pId)
				c.BpduGuard = bpduguard
				c.BrgIfIndex = bId
				err := StpPortConfigParamCheck(c, true)
				if err == nil {
					StpPortConfigSave(c, true)

					// apply to all bridge ports
					for _, port := range p.GetPortListToApplyConfigTo() {
						if bpduguard {
							StpMachineLogger("INFO", "CONFIG", port.IfIndex, port.BrgIfIndex, "Setting BPDU Guard")
						} else {
							StpMachineLogger("INFO", "CONFIG", port.IfIndex, port.BrgIfIndex, "Clearing BPDU Guard")
						}
						port.BpduGuard = bpduguard
					}
				}
			*/
			if bpduguard {
				StpMachineLogger("INFO", "CONFIG", p.IfIndex, p.BrgIfIndex, "Setting BPDU Guard")
			} else {
				StpMachineLogger("INFO", "CONFIG", p.IfIndex, p.BrgIfIndex, "Clearing BPDU Guard")
			}
			p.BpduGuard = bpduguard

			//return err
			return nil
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Bpdu Guard", pId, bId))
}

func StpPortBridgeAssuranceSet(pId int32, bId int32, bridgeassurance bool) error {
	var p *StpPort
	if StpFindPortByIfIndex(pId, bId, &p) {
		if p.BridgeAssurance != bridgeassurance &&
			!p.OperEdge {
			/*
				c := StpPortConfigGet(pId)
				c.BridgeAssurance = bridgeassurance
				c.BrgIfIndex = bId
				err := StpPortConfigParamCheck(c, true)
				if err == nil {
					StpPortConfigSave(c, true)
					// apply to all bridge ports
					for _, port := range p.GetPortListToApplyConfigTo() {
						if bridgeassurance {
							StpMachineLogger("INFO", "CONFIG", port.IfIndex, port.BrgIfIndex, "Setting Bridge Assurance")
						} else {
							StpMachineLogger("INFO", "CONFIG", port.IfIndex, port.BrgIfIndex, "Clearing Bridge Assurance")
						}
						port.BridgeAssurance = bridgeassurance
						port.BridgeAssuranceInconsistant = false
						port.BAWhileTimer.count = int32(p.b.RootTimes.HelloTime * 3)
					}
				}
			*/
			p.BridgeAssurance = bridgeassurance
			p.BridgeAssuranceInconsistant = false
			p.BAWhileTimer.count = int32(p.b.RootTimes.HelloTime * 3)

			//return err
			return nil
		} else {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid port %d or bridge %d supplied for setting Bridge Assurance", pId, bId))
}
