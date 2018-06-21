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

// bridge.go
package stp

import (
	"fmt"
	"net"
	"sync"
)

const BridgeConfigModuleStr = "BRG CFG"

var BridgeMapTable map[BridgeKey]*Bridge
var BridgeListTable []*Bridge
var StpBridgeMac [6]uint8

type BridgeId [8]uint8
type BridgeKey struct {
	Vlan uint16
}

type Bridge struct {

	// 17.18.1
	Begin bool
	// 17.18.2
	BridgeIdentifier BridgeId
	// 17.18.3
	// Root/Designated equal to Bridge Identifier
	BridgePriority PriorityVector
	// 17.18.4
	BridgeTimes Times
	// 17.18.6
	RootPortId int32
	// 17.18.7
	RootTimes Times
	// Stp IfIndex
	StpPorts []int32

	ForceVersion int32
	TxHoldCount  uint64

	// Vlan
	Vlan uint16

	PrsMachineFsm *PrsMachine

	// store the previous bridge id
	OldRootBridgeIdentifier BridgeId

	// IntfRef
	IntfRef string

	// bridge ifIndex
	BrgIfIndex int32
	// hw stgId
	StgId int32

	// a way to sync all machines
	wg sync.WaitGroup

	DebugLevel int
}

type PriorityVector struct {
	RootBridgeId       BridgeId
	RootPathCost       uint32
	DesignatedBridgeId BridgeId
	DesignatedPortId   uint16
	BridgePortId       uint16
}

type Times struct {
	ForwardingDelay uint16
	HelloTime       uint16
	MaxAge          uint16
	MessageAge      uint16
}

func SaveSwitchMac(switchMac string) {
	netAddr, _ := net.ParseMAC(switchMac)
	StpBridgeMac = [6]uint8{netAddr[0], netAddr[1], netAddr[2], netAddr[3], netAddr[4], netAddr[5]}
}

func NewStpBridge(c *StpBridgeConfig) *Bridge {

	vlan := c.Vlan
	bridgeId := CreateBridgeId(StpBridgeMac, c.Priority, vlan)
	if vlan == DEFAULT_STP_BRIDGE_VLAN {
		bridgeId = CreateBridgeId(StpBridgeMac, c.Priority, 0)
	}

	b := &Bridge{
		Begin:            true,
		ForceVersion:     2,
		BridgeIdentifier: bridgeId,
		BridgePriority: PriorityVector{
			RootBridgeId:       bridgeId,
			RootPathCost:       0,
			DesignatedBridgeId: bridgeId,
			DesignatedPortId:   0,
			BridgePortId:       0,
		},
		BridgeTimes: Times{
			ForwardingDelay: c.ForwardDelay,
			HelloTime:       c.HelloTime,
			MaxAge:          c.MaxAge,
			MessageAge:      0,
		},
		RootPortId: 0, // this will be set once a port is set as root
		RootTimes: Times{ForwardingDelay: c.ForwardDelay,
			HelloTime:  c.HelloTime,
			MaxAge:     c.MaxAge,
			MessageAge: 0}, // this will be set once a port is set as root
		TxHoldCount: uint64(c.TxHoldCount),
		Vlan:        vlan,
		DebugLevel:  c.DebugLevel,
	}

	key := BridgeKey{
		Vlan: b.Vlan,
	}

	BridgeMapTable[key] = b

	if len(BridgeListTable) == 0 {
		BridgeListTable = make([]*Bridge, 0)
	}
	BridgeListTable = append(BridgeListTable, b)

	// TODO lets get the linux bridge
	if c.Vlan == 0 {
		// default vlan
		b.BrgIfIndex = DEFAULT_STP_BRIDGE_VLAN
	} else {
		b.BrgIfIndex = int32(c.Vlan)
	}

	// lets create the stg group
	for _, client := range GetAsicDPluginList() {
		b.StgId = client.CreateStgBridge([]uint16{b.Vlan})
	}
	StpLogger("DEBUG", fmt.Sprintf("NEW BRIDGE: %#v\n", b))
	return b
}

func DelStpBridge(b *Bridge, force bool) {
	if b == nil {
		return
	}
	if force {
		var p *StpPort
		for _, pId := range b.StpPorts {
			if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
				DelStpPort(p)
			}
		}
	} else {
		if len(b.StpPorts) > 0 {
			StpLogger("DEBUG", "ERROR BRIDGE STILL HAS PORTS ASSOCIATED")
			return
		}
	}
	b.Stop()

	key := BridgeKey{
		Vlan: b.Vlan,
	}

	delete(BridgeMapTable, key)
	for i, delBrg := range BridgeListTable {
		if delBrg.BridgeIdentifier == b.BridgeIdentifier {
			if len(BridgeListTable) == 1 {
				BridgeListTable = nil
			} else {
				BridgeListTable = append(BridgeListTable[:i], BridgeListTable[i+1:]...)
				for _, client := range GetAsicDPluginList() {
					client.DeleteStgBridge(b.StgId, []uint16{b.Vlan})
				}
			}
		}
	}
}

func (b *Bridge) Stop() {

	if b.PrsMachineFsm != nil {
		b.PrsMachineFsm.Stop()
		b.PrsMachineFsm = nil
	}

	// lets wait for the machines to close
	b.wg.Wait()

}

func (b *Bridge) BEGIN(restart bool) {
	bridgeResponse := make(chan string)
	if !restart {
		// start all the State machines
		// Port Role Selection State Machine (one instance per bridge)
		b.PrsMachineMain()

	}

	// Prsm
	if b.PrsMachineFsm != nil {
		b.PrsMachineFsm.PrsEvents <- MachineEvent{e: PrsEventBegin,
			src:          BridgeConfigModuleStr,
			responseChan: bridgeResponse}
	}

	<-bridgeResponse
	b.Begin = false
}

func StpFindBridgeById(key BridgeKey, b **Bridge) bool {
	var ok bool
	if *b, ok = BridgeMapTable[key]; ok {
		return true
	}
	return false
}

func StpFindBridgeByIfIndex(brgIfIndex int32, brg **Bridge) bool {
	for _, b := range BridgeMapTable {
		if brgIfIndex == b.BrgIfIndex {
			*brg = b
			return true
		}
	}
	return false
}

func CreateBridgeId(bridgeAddress [6]uint8, bridgePriority uint16, vlan uint16) BridgeId {
	return BridgeId{uint8(bridgePriority&0xf000>>8) | uint8(vlan&0xf00>>8),
		uint8(vlan & 0xff),
		bridgeAddress[0],
		bridgeAddress[1],
		bridgeAddress[2],
		bridgeAddress[3],
		bridgeAddress[4],
		bridgeAddress[5]}

}

func CreateBridgeIdStr(bId BridgeId) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
		bId[0],
		bId[1],
		bId[2],
		bId[3],
		bId[4],
		bId[5],
		bId[6],
		bId[7])
}

func GetBridgeAddrFromBridgeId(b BridgeId) [6]uint8 {
	return [6]uint8{
		b[2],
		b[3],
		b[4],
		b[5],
		b[6],
		b[7],
	}
}

func GetBridgeVlanFromBridgeId(b BridgeId) (vlan uint16) {
	return ((uint16(b[0]) & 0xF) << 8) | uint16(b[1])
}

func GetBridgePriorityFromBridgeId(b BridgeId) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

// Compare BridgeId
// 0 == equal
// > 1 == greater than
// < 1 == less than
func CompareBridgeId(b1 BridgeId, b2 BridgeId) int {
	if b1 == b2 {
		//StpLogger("DEBUG", fmt.Sprintf("CompareBridgeId returns B1 SAME, b1[%#v] b2[%#v]", b1, b2))
		return 0
	} else if (b1[0] < b2[0]) ||
		((b1[0] == b2[0]) && (b1[1] < b2[1])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] < b2[2])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] < b2[3])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] < b2[4])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] == b2[4]) && (b1[5] < b2[5])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] == b2[4]) && (b1[5] == b2[5]) && (b1[6] < b2[6])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] == b2[4]) && (b1[5] == b2[5]) && (b1[6] == b2[6]) && (b1[7] < b2[7])) {
		//StpLogger("DEBUG", fmt.Sprintf("CompareBridgeId returns B1 SUPERIOR, b1[%#v] b2[%#v]", b1, b2))
		return -1
	} else {
		//StpLogger("DEBUG", fmt.Sprintf("CompareBridgeId returns B1 INFERIOR, b1[%#v] b2[%#v]", b1, b2))
		return 1
	}
}

func CompareBridgeAddr(b1 [6]uint8, b2 [6]uint8) int {
	if (b1[0] < b2[0]) ||
		((b1[0] == b2[0]) && (b1[1] < b2[1])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] < b2[2])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] < b2[3])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] < b2[4])) ||
		((b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2]) && (b1[3] == b2[3]) && (b1[4] == b2[4]) && (b1[5] < b2[5])) {
		return -1
	} else if b1 == b2 {
		return 0
	}
	return 1

}

// 17.6 Priority vector calculations
func IsMsgPriorityVectorSuperiorThanPortPriorityVector(msg *PriorityVector, port *PriorityVector) bool {
	/*
		return (CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) < 0) ||
			((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) && (msg.RootPathCost < port.RootPathCost)) ||
			((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) && (msg.RootPathCost == port.RootPathCost) && (CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) < 0)) ||
			((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) && (msg.RootPathCost == port.RootPathCost) && (CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) == 0) && (msg.DesignatedPortId < port.DesignatedPortId)) ||
			(CompareBridgeAddr(GetBridgeAddrFromBridgeId(msg.DesignatedBridgeId),
				GetBridgeAddrFromBridgeId(port.DesignatedBridgeId)) == 0 &&
				(msg.DesignatedPortId == port.DesignatedPortId))
	*/
	if CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) < 0 {
		//StpLogger("DEBUG", "b1 root bridge id superior to b2 root bridge id")
		return true
	} else if (CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) &&
		(msg.RootPathCost < port.RootPathCost) {
		//StpLogger("DEBUG", "b1 root bridge id equal b1 root bridge id and b1 root path superior to b2 root path cost")
		return true
	} else if (CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) &&
		(msg.RootPathCost == port.RootPathCost) &&
		(CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) < 0) {
		//StpLogger("DEBUG", "b1 root bridge id equal b1 root bridge id and b1 root path equal to b2 root path cost, desgn bridge id superior to b1 desgn bridge id")
		return true
	} else if (CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) &&
		(msg.RootPathCost == port.RootPathCost) &&
		(CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) == 0) &&
		(msg.DesignatedPortId < port.DesignatedPortId) {
		//StpLogger("DEBUG", "b1 root bridge id equal b1 root bridge id and b1 root path equal to b2 root path cost, desgn bridge id equal to b1 desgn bridge id, b1 desgn portid superior to b2 desgn portid")
		return true
	} else if CompareBridgeAddr(GetBridgeAddrFromBridgeId(msg.DesignatedBridgeId),
		GetBridgeAddrFromBridgeId(port.DesignatedBridgeId)) == 0 &&
		(msg.DesignatedPortId == port.DesignatedPortId) {
		//StpLogger("DEBUG", "b1 desgn brg addr equal b2 desgn brg addr and b1 desgn portid equal b2 desgn portid")
		return true
	}

	return false

}

func IsMsgPriorityVectorSuperiorOrSameThanPortPriorityVector(msg *PriorityVector, port *PriorityVector) bool {

	return IsMsgPriorityVectorSuperiorThanPortPriorityVector(msg, port)

}

func IsMsgPriorityVectorTheSameOrWorseThanPortPriorityVector(msg *PriorityVector, port *PriorityVector) bool {
	return (msg == port) ||
		IsMsgPriorityVectorSuperiorThanPortPriorityVector(port, msg)
}

func IsMsgPriorityVectorWorseThanPortPriorityVector(msg *PriorityVector, port *PriorityVector) bool {
	return (CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) > 0) ||
		((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0 && (msg.RootPathCost > port.RootPathCost)) ||
			((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) && (msg.RootPathCost == port.RootPathCost) && (CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) > 0))) ||
		((CompareBridgeId(msg.RootBridgeId, port.RootBridgeId) == 0) && (msg.RootPathCost == port.RootPathCost) && (CompareBridgeId(msg.DesignatedBridgeId, port.DesignatedBridgeId) == 0) && (msg.DesignatedPortId > port.DesignatedPortId))

}

func IsDesignatedPriorytVectorNotHigherThanPortPriorityVector(d *PriorityVector, p *PriorityVector) bool {
	// re-use function
	return IsMsgPriorityVectorTheSameOrWorseThanPortPriorityVector(d, p)
}

func (b *Bridge) AllSynced() bool {

	var p *StpPort
	allsynced := false
	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			if p.Selected &&
				p.Role == p.SelectedRole &&
				(p.Synced || p.SelectedRole == PortRoleRootPort) {
				allsynced = true
			} else {
				return false
			}
		}
	}
	return allsynced
}

func (b *Bridge) ReRooted(p *StpPort) bool {
	rerooted := true
	if p.RrWhileTimer.count == 0 {
		rerooted = false
	} else {
		var op *StpPort
		for _, pId := range b.StpPorts {
			if pId == p.IfIndex {
				continue
			}
			if StpFindPortByIfIndex(pId, b.BrgIfIndex, &op) {
				if op.RrWhileTimer.count != 0 {
					rerooted = false
				}
			}
		}
	}
	return rerooted
}
