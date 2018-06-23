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

// 802.1D-2004 17.28 Port Role Selection State Machine
//The Port Role Selection state machine shall implement the function specified by the state diagram in Figure
//17-19, the definitions in 17.13, 17.16, 17.20, and 17.21, and the variable declarations in 17.17, 17.18, and
//17.19. It selects roles for all Bridge Ports.
//On initialization all Bridge Ports are assigned the Disabled Port Role. Whenever any Bridge Port’s reselect
//variable (17.19.34) is set by the Port Information state machine (17.27), spanning tree information including
//the designatedPriority (17.19.4) and designatedTimes (17.19.5) for each Port is recomputed and its Port
//Role (selectedRole, 17.19.37) updated by the updtRolesTree() procedure (17.21.25). The reselect variables
//are cleared before computation starts so that recomputation will take place if new information becomes
//available while the computation is in progress.
package stp

import (
	"fmt"
	//"time"
	"utils/fsm"
)

const PrsMachineModuleStr = "PRSM"

const (
	PrsStateNone = iota + 1
	PrsStateInitBridge
	PrsStateRoleSelection
)

var PrsStateStrMap map[fsm.State]string

func PrsMachineStrStateMapInit() {
	PrsStateStrMap = make(map[fsm.State]string)
	PrsStateStrMap[PrsStateNone] = "None"
	PrsStateStrMap[PrsStateInitBridge] = "Init Bridge"
	PrsStateStrMap[PrsStateRoleSelection] = "Role Selection"
}

const (
	PrsEventBegin = iota + 1
	PrsEventUnconditionallFallThrough
	PrsEventReselect
)

// PrsMachine holds FSM and current State
// and event channels for State transitions
type PrsMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	b *Bridge

	// debug level
	debugLevel int

	// machine specific events
	PrsEvents chan MachineEvent
	// enable logging
	PrsLogEnableEvent chan bool
}

func (m *PrsMachine) GetCurrStateStr() string {
	return PrsStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PrsMachine) GetPrevStateStr() string {
	return PrsStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewStpPimMachine will create a new instance of the LacpRxMachine
func NewStpPrsMachine(b *Bridge) *PrsMachine {
	prsm := &PrsMachine{
		b:                 b,
		debugLevel:        b.DebugLevel,
		PrsEvents:         make(chan MachineEvent, 50),
		PrsLogEnableEvent: make(chan bool)}

	b.PrsMachineFsm = prsm

	return prsm
}

func (prsm *PrsMachine) PrsLogger(s string) {
	StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, prsm.b.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (prsm *PrsMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if prsm.Machine == nil {
		prsm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	prsm.Machine.Rules = r
	prsm.Machine.Curr = &StpStateEvent{
		strStateMap: PrsStateStrMap,
		logEna:      true,
		logger:      prsm.PrsLogger,
		owner:       PrsMachineModuleStr,
		ps:          PrsStateNone,
		s:           PrsStateNone,
	}

	return prsm.Machine
}

// Stop should clean up all resources
func (prsm *PrsMachine) Stop() {

	close(prsm.PrsEvents)
	close(prsm.PrsLogEnableEvent)
}

// PrsMachineInitBridge
func (prsm *PrsMachine) PrsMachineInitBridge(m fsm.Machine, data interface{}) fsm.State {
	prsm.updtRoleDisabledTree()
	return PrsStateInitBridge
}

// PrsMachineRoleSelection
func (prsm *PrsMachine) PrsMachineRoleSelection(m fsm.Machine, data interface{}) fsm.State {
	prsm.clearReselectTree()
	prsm.updtRolesTree()
	prsm.setSelectedTree()

	return PrsStateRoleSelection
}

func PrsMachineFSMBuild(b *Bridge) *PrsMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrxmMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	prsm := NewStpPrsMachine(b)

	// BEGIN -> INIT_BRIDGE
	rules.AddRule(PrsStateNone, PrsEventBegin, prsm.PrsMachineInitBridge)

	// UNINTENTIONAL FALL THROUGH -> ROLE SELECTION
	rules.AddRule(PrsStateInitBridge, PrsEventUnconditionallFallThrough, prsm.PrsMachineRoleSelection)

	// RESLECT -> ROLE SELECTION
	rules.AddRule(PrsStateRoleSelection, PrsEventReselect, prsm.PrsMachineRoleSelection)

	// Create a new FSM and apply the rules
	prsm.Apply(&rules)

	return prsm
}

// PrsMachineMain:
func (b *Bridge) PrsMachineMain() {

	// Build the State machine for STP Bridge Detection State Machine according to
	// 802.1d Section 17.25
	prsm := PrsMachineFSMBuild(b)
	b.wg.Add(1)

	// set the inital State
	prsm.Machine.Start(prsm.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PrsMachine) {
		StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, b.BrgIfIndex, "Machine Start")
		defer m.b.wg.Done()
		for {
			select {
			case event, ok := <-m.PrsEvents:
				//fmt.Println("Event Rx", event.src, event.e)
				if ok {
					rv := m.Machine.ProcessEvent(event.src, event.e, nil)
					if rv != nil {
						StpMachineLogger("ERROR", PrsMachineModuleStr, -1, b.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, event.e, PrsStateStrMap[m.Machine.Curr.CurrentState()]))
					} else {
						if m.Machine.Curr.CurrentState() == PrsStateInitBridge {
							rv := m.Machine.ProcessEvent(PrsMachineModuleStr, PrsEventUnconditionallFallThrough, nil)
							if rv != nil {
								StpMachineLogger("ERROR", PrsMachineModuleStr, -1, b.BrgIfIndex, fmt.Sprintf("%s event[%d] currState[%s]\n", rv, event.e, PrsStateStrMap[m.Machine.Curr.CurrentState()]))
							}
						}
					}

					if event.responseChan != nil {
						SendResponse(PrsMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, b.BrgIfIndex, "Machine End")
					return
				}
			case ena := <-m.PrsLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(prsm)
}

// clearReselectTree: 17.21.2
func (prsm *PrsMachine) clearReselectTree() {
	var p *StpPort
	b := prsm.b

	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			//if p.PortEnabled {
			p.Reselect = false
			//}
		}
	}
}

func (prsm *PrsMachine) updtRoleDisabledTree() {
	var p *StpPort
	b := prsm.b

	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDisabledPort)
			p.SelectedRole = PortRoleDisabledPort
		}
	}
}

// updtRolesTree: 17.21.25
func (prsm *PrsMachine) updtRolesTree() {

	b := prsm.b

	var p *StpPort
	var rootPortId int32
	rootPathVector := PriorityVector{
		RootBridgeId:       b.BridgePriority.DesignatedBridgeId,
		DesignatedBridgeId: b.BridgePriority.DesignatedBridgeId,
	}

	// 17.21.25 (c)(1)
	rootTimes := Times{
		ForwardingDelay: b.BridgeTimes.ForwardingDelay,
		HelloTime:       b.BridgeTimes.HelloTime,
		MaxAge:          b.BridgeTimes.MaxAge,
		MessageAge:      b.BridgeTimes.MessageAge,
	}

	tmpVector := rootPathVector

	// lets consider each port a root to begin with
	myBridgeId := rootPathVector.RootBridgeId

	// lets find the root port
	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: InfoIs %d", p.InfoIs))
			// 17.21.25 (a)
			if p.InfoIs == PortInfoStateReceived {

				/*if CompareBridgeAddr(GetBridgeAddrFromBridgeId(myBridgeId),
					GetBridgeAddrFromBridgeId(p.PortPriority.DesignatedBridgeId)) == 0 {
					continue
				}*/
				StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: port root bridge %#v  tmpRootBridge %#v", p.PortPriority.RootBridgeId, tmpVector.RootBridgeId))
				compare := CompareBridgeId(p.PortPriority.RootBridgeId, tmpVector.RootBridgeId)
				switch compare {
				// 17.21.25 (b) bridge is superior
				case -1:

					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: Root Bridge Received is SUPERIOR port Priority %#v", p.PortPriority))
					}
					tmpVector.RootBridgeId = p.PortPriority.RootBridgeId
					tmpVector.RootPathCost = p.PortPriority.RootPathCost + p.PortPathCost
					tmpVector.DesignatedBridgeId = p.PortPriority.DesignatedBridgeId
					tmpVector.DesignatedPortId = p.PortPriority.DesignatedPortId
					rootPortId = int32(p.Priority<<8 | p.PortId)
					// 17.21.25 (c)(2)
					rootTimes = p.PortTimes
				case 0:
					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: Root Bridge Received by port SAME")
					}
					// 17.21.25 (b) path cost or port determines root
					tmpCost := p.PortPriority.RootPathCost + p.PortPathCost
					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: rx+txCost[%d] bridgeCost[%d]", tmpCost, tmpVector.RootPathCost))
					}
					if tmpCost < tmpVector.RootPathCost {
						if prsm.debugLevel > 1 {
							StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: DesignatedBridgeId received by port is SUPERIOR")
						}
						tmpVector.RootPathCost = tmpCost
						tmpVector.DesignatedBridgeId = p.PortPriority.DesignatedBridgeId
						tmpVector.DesignatedPortId = p.PortPriority.DesignatedPortId
						rootPortId = int32(p.Priority<<8 | p.PortId)
					} else if tmpCost == tmpVector.RootPathCost {
						if prsm.debugLevel > 1 {
							StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: DesignatedBridgeId received by port is SAME")
						}
						if p.PortPriority.DesignatedPortId <
							tmpVector.DesignatedPortId {
							if prsm.debugLevel > 1 {
								StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: DesignatedPortId received by port is SUPPERIOR")
							}
							tmpVector.DesignatedPortId = p.PortPriority.DesignatedPortId
							rootPortId = int32(p.Priority<<8 | p.PortId)
						} else if p.PortPriority.DesignatedPortId ==
							tmpVector.DesignatedPortId {
							var rp *StpPort
							var localPortId int32
							if StpFindPortByIfIndex(rootPortId, b.BrgIfIndex, &rp) {
								rootPortId = int32((rp.Priority << 8) | p.PortId)
								localPortId = int32((p.Priority << 8) | p.PortId)
								if localPortId < rootPortId {
									if prsm.debugLevel > 1 {
										StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: received portId is SUPPERIOR")
									}
									rootPortId = int32(p.Priority<<8 | p.PortId)
								}
							}
						}
					}
				}
			}
		}
	}

	// lets copy over the tmpVector over to the rootPathVector
	if rootPortId != 0 {
		if prsm.debugLevel > 1 {
			StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, b.BrgIfIndex, fmt.Sprintf("updtRolesTree: Port %d selected as the root port", rootPortId))
		}
		compare := CompareBridgeAddr(GetBridgeAddrFromBridgeId(b.BridgePriority.RootBridgeId),
			GetBridgeAddrFromBridgeId(tmpVector.RootBridgeId))
		if compare != 0 {
			b.OldRootBridgeIdentifier = b.BridgePriority.RootBridgeId
		}

		b.BridgePriority.RootBridgeId = tmpVector.RootBridgeId
		b.BridgePriority.RootPathCost = tmpVector.RootPathCost
		b.RootTimes = rootTimes
		b.RootTimes.MessageAge += 1
		b.RootPortId = rootPortId
	} else {
		if prsm.debugLevel > 1 {
			StpMachineLogger("DEBUG", PrsMachineModuleStr, 0, b.BrgIfIndex, "updtRolesTree: This bridge is the root bridge")
		}
		compare := CompareBridgeAddr(GetBridgeAddrFromBridgeId(b.BridgeIdentifier),
			GetBridgeAddrFromBridgeId(tmpVector.RootBridgeId))
		if compare != 0 {
			b.OldRootBridgeIdentifier = b.BridgePriority.RootBridgeId
		}

		b.BridgePriority.RootBridgeId = tmpVector.RootBridgeId
		b.BridgePriority.RootPathCost = tmpVector.RootPathCost
		b.RootTimes = rootTimes
		b.RootPortId = 0
	}
	if prsm.debugLevel > 1 {
		StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, b.BrgIfIndex, fmt.Sprintf("BridgePriority: %#v  BridgeTimes: %#v", b.BridgePriority, b.RootTimes))
	}
	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {

			// 17.21.25 (e)
			p.PortTimes = b.RootTimes
			p.PortTimes.HelloTime = b.BridgeTimes.HelloTime

			desgPortId := p.PortPriority.DesignatedPortId
			brgPortId := p.PortPriority.BridgePortId
			localPortId := p.Priority<<8 | p.PortId

			p.b.BridgePriority.DesignatedPortId = 0
			p.PortPriority.DesignatedPortId = 0
			p.PortPriority.BridgePortId = 0

			if prsm.debugLevel > 1 {
				StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, b.BrgIfIndex, fmt.Sprintf("updtRolesTree: portEnabled %t, infoIs %d\n", p.PortEnabled, p.InfoIs))
			}

			if p.BridgeAssurance &&
				!p.OperEdge &&
				p.BridgeAssuranceInconsistant {
				defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, false)
				p.UpdtInfo = false
				defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleAlternatePort)
				p.SelectedRole = PortRoleAlternatePort
				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: Bridge Assurance port role selected ALTERNATE")
				}
			} else if !p.PortEnabled || p.InfoIs == PortInfoStateDisabled {
				// 17.21.25 (f) if port is disabled
				defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDisabledPort)
				p.SelectedRole = PortRoleDisabledPort

				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree:1 port role selected DISABLED")
				}
			} else if p.InfoIs == PortInfoStateAged {
				// 17.21.25 (g)
				defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, true)
				p.UpdtInfo = true
				defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDesignatedPort)
				p.SelectedRole = PortRoleDesignatedPort
				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree:1 port role selected DESIGNATED")
				}
			} else if p.InfoIs == PortInfoStateMine {
				// 17.21.25 (h)
				defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDesignatedPort)
				p.SelectedRole = PortRoleDesignatedPort

				if p.b.BridgePriority == p.PortPriority &&
					desgPortId == localPortId {
					if p.PortTimes != b.RootTimes {
						if prsm.debugLevel > 1 {
							StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: port times[%#v] != root times[%#v]", p.PortTimes, p.b.RootTimes))
						}
						defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, true)
						p.UpdtInfo = true
					}
				} else {
					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: BridgePriority[%#v] != PortPriority[%#v]", p.b.BridgePriority, p.PortPriority))
					}
					defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, true)
					p.UpdtInfo = true
				}
				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree:2 port role selected DESIGNATED")
				}
			} else if p.InfoIs == PortInfoStateReceived {
				// 17.21.25 (i)
				if rootPortId == int32(localPortId) {
					defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleRootPort)
					p.SelectedRole = PortRoleRootPort
					// this will allow for packets to tx the interface
					defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, false)
					p.UpdtInfo = false
					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: port role selected ROOT")
					}
				} else {
					// 17.21.25 (j), (k), (l)
					// designated not higher than port priority
					p.b.BridgePriority.DesignatedPortId = localPortId
					p.PortPriority.DesignatedPortId = desgPortId
					if prsm.debugLevel > 1 {
						StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: check not better BridgePriority[%#v] PortPriority[%#v]", p.b.BridgePriority, p.PortPriority))
					}
					if IsDesignatedPriorytVectorNotHigherThanPortPriorityVector(&p.b.BridgePriority, &p.PortPriority) {
						if prsm.debugLevel > 1 {
							StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("updtRolesTree: check addr not same myBridge[%#v] portBridge[%#v]", GetBridgeAddrFromBridgeId(myBridgeId), GetBridgeAddrFromBridgeId(p.PortPriority.DesignatedBridgeId)))
						}
						if CompareBridgeAddr(GetBridgeAddrFromBridgeId(p.PortPriority.DesignatedBridgeId),
							GetBridgeAddrFromBridgeId(myBridgeId)) != 0 {
							defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleAlternatePort)
							p.SelectedRole = PortRoleAlternatePort
							// this will allow for packets to tx the interface
							defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, false)
							p.UpdtInfo = false
							if prsm.debugLevel > 1 {
								StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: port role selected ALTERNATE")
							}
						} else {

							if (p.Priority<<8 | p.PortId) != p.PortPriority.DesignatedPortId {
								defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleBackupPort)
								p.SelectedRole = PortRoleBackupPort
								// this will allow for packets to tx the interface
								defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, false)
								p.UpdtInfo = false
								if prsm.debugLevel > 1 {
									StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree: port role selected BACKUP")
								}
							} else {
								defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDesignatedPort)
								p.SelectedRole = PortRoleDesignatedPort
								defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, true)
								p.UpdtInfo = true
								if prsm.debugLevel > 1 {
									StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree:3 port role selected DESIGNATED")
								}
							}
						}
					} else {
						//if p.SelectedRole != PortRoleDesignatedPort {
						defer p.NotifySelectedRoleChanged(PrsMachineModuleStr, p.SelectedRole, PortRoleDesignatedPort)
						p.SelectedRole = PortRoleDesignatedPort
						defer p.NotifyUpdtInfoChanged(PrsMachineModuleStr, p.UpdtInfo, true)
						p.UpdtInfo = true
						//}
						if prsm.debugLevel > 1 {
							StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "updtRolesTree:4 port role selected DESIGNATED")
						}
					}
				}
			}
			p.PortPriority.DesignatedPortId = desgPortId
			p.PortPriority.BridgePortId = brgPortId

		}
	}
}

// setSelectedTree: 17.21.16
func (prsm *PrsMachine) setSelectedTree() {
	var p *StpPort
	b := prsm.b
	setAllSelectedTrue := true

	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			if p.Reselect {
				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, "setSelectedTree: is in reselet mode")
				}
				setAllSelectedTrue = false
				break
			}
		}
	}
	if setAllSelectedTrue {
		if prsm.debugLevel > 1 {
			StpMachineLogger("DEBUG", PrsMachineModuleStr, -1, b.BrgIfIndex, "setSelectedTree: setting all ports as selected")
		}
		for _, pId := range b.StpPorts {
			if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
				if prsm.debugLevel > 1 {
					StpMachineLogger("DEBUG", PrsMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("setSelectedTree: setting selected prev selected state %t", p.Selected))
				}
				defer p.NotifySelectedChanged(PrsMachineModuleStr, p.Selected, true)
				p.Selected = true
			}
		}
	}
}
