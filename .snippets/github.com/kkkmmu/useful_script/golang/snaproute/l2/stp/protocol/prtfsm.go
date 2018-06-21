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

// 802.1D-2004 17.29 Port Role Selection State Machine
//The Port Role Selection state machine shall implement the function specified by the state diagram in Figure
//17-19, the definitions in 17.13, 17.16, 17.20, and 17.21, and the variable declarations in 17.17, 17.18, and
//17.19. It selects roles for all Bridge Ports.
//On initialization all Bridge Ports are assigned the Disabled Port Role. Whenever any Bridge Port’s reselect
//variable (17.19.34) is set by the Port Information state machine (17.27), spanning tree information including
//the designatedPriority (17.19.4) and designatedTimes (17.19.5) for each Port is recomputed and its Port
//Role (selectedRole, 17.19.37) updated by the updtRolesTree() procedure (17.21.25). The reselect variables
//are cleared before computation starts so that recomputation will take place if new information becomes
//available while the computation is in progress.
//
package stp

import (
	"fmt"
	//"time"
	"utils/fsm"
)

const PrtMachineModuleStr = "PRTM"

const (
	PrtStateNone = iota + 1
	// Role: Disabled
	PrtStateInitPort
	PrtStateDisablePort
	PrtStateDisabledPort
	// Role Root
	PrtStateRootPort
	PrtStateReRoot
	PrtStateRootAgreed
	PrtStateRootProposed
	PrtStateRootForward
	PrtStateRootLearn
	PrtStateReRooted
	// Role Designated
	PrtStateDesignatedPort
	PrtStateDesignatedRetired
	PrtStateDesignatedSynced
	PrtStateDesignatedPropose
	PrtStateDesignatedForward
	PrtStateDesignatedLearn
	PrtStateDesignatedDiscard
	// Role Alternate Backup
	PrtStateAlternatePort
	PrtStateAlternateAgreed
	PrtStateAlternateProposed
	PrtStateBlockPort
	PrtStateBackupPort
)

var PrtStateStrMap map[fsm.State]string

func PrtMachineStrStateMapInit() {
	PrtStateStrMap = make(map[fsm.State]string)
	PrtStateStrMap[PrtStateNone] = "None"
	PrtStateStrMap[PrtStateInitPort] = "Init Port"
	PrtStateStrMap[PrtStateDisablePort] = "Disable Port"
	PrtStateStrMap[PrtStateDisabledPort] = "Disabled Port"
	PrtStateStrMap[PrtStateRootPort] = "Root Port"
	PrtStateStrMap[PrtStateReRoot] = "Re-Root"
	PrtStateStrMap[PrtStateRootAgreed] = "Root Agreed"
	PrtStateStrMap[PrtStateRootProposed] = "Root Proposed"
	PrtStateStrMap[PrtStateRootForward] = "Root Forward"
	PrtStateStrMap[PrtStateRootLearn] = "Root Learn"
	PrtStateStrMap[PrtStateReRooted] = "Re-Rooted"
	PrtStateStrMap[PrtStateDesignatedPort] = "Designated Port"
	PrtStateStrMap[PrtStateDesignatedRetired] = "Designated Retired"
	PrtStateStrMap[PrtStateDesignatedSynced] = "Designated Synced"
	PrtStateStrMap[PrtStateDesignatedPropose] = "Designated Propose"
	PrtStateStrMap[PrtStateDesignatedForward] = "Designated Forward"
	PrtStateStrMap[PrtStateDesignatedLearn] = "Designated Learn"
	PrtStateStrMap[PrtStateDesignatedDiscard] = "Designated Discard"
	PrtStateStrMap[PrtStateAlternatePort] = "Alternate Port"
	PrtStateStrMap[PrtStateAlternateAgreed] = "Alternate Agreed"
	PrtStateStrMap[PrtStateAlternateProposed] = "Alternate Proposed"
	PrtStateStrMap[PrtStateBlockPort] = "Block Port"
	PrtStateStrMap[PrtStateBackupPort] = "Backup Port"
}

const (
	PrtEventBegin                     = iota + 1
	PrtEventUnconditionallFallThrough // 2
	// events taken from Figure 17.20 Disabled Port role transitions
	PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo // 3
	PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo                              // 4
	PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo                                    // 5
	PrtEventSyncAndSelectedAndNotUpdtInfo                                                     // 6 also applies to Alternate and Backup Port role
	PrtEventReRootAndSelectedAndNotUpdtInfo                                                   // 7 also applies to Alternate and Backup Port role
	PrtEventNotSyncedAndSelectedAndNotUpdtInfo                                                // 8 also applies to Alternate and Backup Port role
	// events taken from Figure 17.21 Root Port role transitions
	PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo           // 9
	PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo                                            // 10 also applies to Alternate and Backup Port role
	PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo                                           // 11 also applies to Alternate and Backup Port role
	PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo                                               // 12 also applies to Alternate and Backup Port role
	PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo                                         // 13
	PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo                                        // 14
	PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo                                               // 15
	PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo                      // 16
	PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo           // 17
	PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo            // 18
	PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo // 19
	// events take from Figure 17-22 Designated port role transitions
	PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo      // 20
	PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo             // 21
	PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo                         // 22
	PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo                                              // 23
	PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo                                            // 24
	PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo                                                   // 25
	PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo                                       // 26
	PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo                          // 27
	PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo                        // 28
	PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo              // 29
	PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo            // 30
	PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo                                  // 31
	PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo                                // 32
	PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo        // 33
	PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo                  // 34
	PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo                  // 35
	PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo                         // 36
	PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo                // 37
	PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo                       // 38
	PrtEventNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo                                  // 39
	PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo // 40
	PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo        // 41
	PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo        // 42
	PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo               // 43
	PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo      // 44
	PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo             // 45
	PrtEventNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo                        // 46
	// events taken from Figure 17-23 Alternate and Backup Port role transitions
	PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo   // 47
	PrtEventSelectedRoleEqualBackupPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo  // 48
	PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo // 49
	PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo                             // 50
)

// PrtMachine holds FSM and current State
// and event channels for State transitions
type PrtMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PrtEvents chan MachineEvent
	// enable logging
	PrtLogEnableEvent chan bool
}

func (m *PrtMachine) GetCurrStateStr() string {
	return PrtStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *PrtMachine) GetPrevStateStr() string {
	return PrtStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewStpPrtMachine will create a new instance of the LacpRxMachine
func NewStpPrtMachine(p *StpPort) *PrtMachine {
	prtm := &PrtMachine{
		p:                 p,
		PrtEvents:         make(chan MachineEvent, 50),
		PrtLogEnableEvent: make(chan bool)}

	p.PrtMachineFsm = prtm

	return prtm
}

func (prtm *PrtMachine) PrtLogger(s string) {
	StpMachineLogger("DEBUG", PrtMachineModuleStr, prtm.p.IfIndex, prtm.p.BrgIfIndex, s)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (prtm *PrtMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if prtm.Machine == nil {
		prtm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	prtm.Machine.Rules = r
	prtm.Machine.Curr = &StpStateEvent{
		strStateMap: PrtStateStrMap,
		logEna:      true,
		logger:      prtm.PrtLogger,
		owner:       PrtMachineModuleStr,
		ps:          PrtStateNone,
		s:           PrtStateNone,
	}

	return prtm.Machine
}

// Stop should clean up all resources
func (prtm *PrtMachine) Stop() {

	close(prtm.PrtEvents)
	close(prtm.PrtLogEnableEvent)
}

// PrtMachineInitPort
func (prtm *PrtMachine) PrtMachineInitPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyRoleChanged(p.Role, PortRoleDisabledPort)
	p.Role = PortRoleDisabledPort
	defer prtm.NotifyLearnChanged(p.Learn, false)
	p.Learn = false
	defer prtm.NotifyForwardChanged(p.Forward, false)
	p.Forward = false
	p.Synced = false
	p.Sync = true
	p.ReRoot = true
	p.RrWhileTimer.count = int32(p.b.BridgeTimes.ForwardingDelay)
	p.FdWhileTimer.count = int32(p.b.BridgeTimes.MaxAge)
	p.RbWhileTimer.count = 0
	return PrtStateInitPort
}

//PrtMachineDisablePort
func (prtm *PrtMachine) PrtMachineDisablePort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyRoleChanged(p.Role, p.SelectedRole)
	p.Role = p.SelectedRole
	defer prtm.NotifyLearnChanged(p.Learn, false)
	p.Learn = false
	defer prtm.NotifyForwardChanged(p.Forward, false)
	p.Forward = false
	return PrtStateDisablePort
}

//PrtMachineDisablePort
func (prtm *PrtMachine) PrtMachineDisabledPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.FdWhileTimer.count = int32(p.b.BridgeTimes.MaxAge)
	p.Synced = true
	p.RrWhileTimer.count = 0
	p.Sync = false
	p.ReRoot = false
	return PrtStateDisabledPort
}

//PrtMachineRootProposed
func (prtm *PrtMachine) PrtMachineRootProposed(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	prtm.setSyncTree(p.IfIndex)
	p.Proposed = false
	return PrtStateRootProposed
}

//PrtMachineRootAgreed
func (prtm *PrtMachine) PrtMachineRootAgreed(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.Proposed = false
	p.Sync = false
	p.Agree = true
	defer prtm.NotifyNewInfoChanged(p.NewInfo, true)
	p.NewInfo = true
	return PrtStateRootAgreed
}

//PrtMachineReroot
func (prtm *PrtMachine) PrtMachineReRoot(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	prtm.setReRootTree(p.IfIndex)
	return PrtStateReRoot
}

//PrtMachineRootForward
func (prtm *PrtMachine) PrtMachineRootForward(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.FdWhileTimer.count = 0
	defer prtm.NotifyForwardChanged(p.Forward, true)
	p.Forward = true
	return PrtStateRootForward
}

//PrtMachineRootLearn
func (prtm *PrtMachine) PrtMachineRootLearn(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.FdWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
	defer prtm.NotifyLearnChanged(p.Learn, true)
	p.Learn = true
	return PrtStateRootLearn
}

//PrtMachineReRooted
func (prtm *PrtMachine) PrtMachineReRooted(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.ReRoot = false
	return PrtStateReRooted
}

//PrtMachineRootPort
func (prtm *PrtMachine) PrtMachineRootPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyRoleChanged(p.Role, PortRoleRootPort)
	p.Role = PortRoleRootPort
	p.RrWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
	return PrtStateRootPort
}

//PrtMachineDesignatedPropose
func (prtm *PrtMachine) PrtMachineDesignatedPropose(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer p.NotifyProposingChanged(PrtMachineModuleStr, p.Proposing, true)
	p.Proposing = true
	p.EdgeDelayWhileTimer.count = int32(p.EdgeDelay())
	defer prtm.NotifyNewInfoChanged(p.NewInfo, true)
	p.NewInfo = true
	return PrtStateDesignatedPropose
}

//PrtMachineDesignatedSynced
func (prtm *PrtMachine) PrtMachineDesignatedSynced(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.RrWhileTimer.count = 0
	p.Synced = true
	p.Sync = false
	return PrtStateDesignatedSynced
}

//PrtMachineDesignatedRetired
func (prtm *PrtMachine) PrtMachineDesignatedRetired(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.ReRoot = false
	return PrtStateDesignatedRetired
}

//PrtMachineDesignatedForward
func (prtm *PrtMachine) PrtMachineDesignatedForward(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyForwardChanged(p.Forward, true)
	p.Forward = true
	p.FdWhileTimer.count = 0
	p.Agreed = p.SendRSTP
	return PrtStateDesignatedForward
}

//PrtMachineDesignatedLearn
func (prtm *PrtMachine) PrtMachineDesignatedLearn(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyLearnChanged(p.Learn, true)
	p.Learn = true
	p.FdWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
	return PrtStateDesignatedLearn
}

//PrtMachineDesignatedDiscard
func (prtm *PrtMachine) PrtMachineDesignatedDiscard(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyLearnChanged(p.Learn, false)
	p.Learn = false
	defer prtm.NotifyForwardChanged(p.Forward, false)
	p.Forward = false
	p.Disputed = false
	p.FdWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
	return PrtStateDesignatedDiscard
}

//PrtMachineDesignatedPort
func (prtm *PrtMachine) PrtMachineDesignatedPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyRoleChanged(p.Role, PortRoleDesignatedPort)
	p.Role = PortRoleDesignatedPort
	return PrtStateDesignatedPort
}

//PrtMachineAlternateProposed
func (prtm *PrtMachine) PrtMachineAlternateProposed(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	prtm.setSyncTree(p.IfIndex)
	p.Proposed = false
	return PrtStateAlternateProposed
}

//PrtMachineAlternateAgreed
func (prtm *PrtMachine) PrtMachineAlternateAgreed(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.Proposed = false
	p.Agree = true
	defer prtm.NotifyNewInfoChanged(p.NewInfo, true)
	p.NewInfo = true
	return PrtStateAlternateAgreed
}

//PrtMachineBlockPort
func (prtm *PrtMachine) PrtMachineBlockPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	defer prtm.NotifyRoleChanged(p.Role, p.SelectedRole)
	p.Role = p.SelectedRole
	defer prtm.NotifyLearnChanged(p.Learn, false)
	p.Learn = false
	defer prtm.NotifyForwardChanged(p.Forward, false)
	p.Forward = false
	return PrtStateBlockPort
}

//PrtMachineBackupPort
func (prtm *PrtMachine) PrtMachineBackupPort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.RbWhileTimer.count = int32(2 * p.PortTimes.HelloTime)
	return PrtStateBackupPort
}

//PrtMachineAlternatePort
func (prtm *PrtMachine) PrtMachineAlternatePort(m fsm.Machine, data interface{}) fsm.State {
	p := prtm.p
	p.FdWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
	p.Synced = true
	p.RrWhileTimer.count = 0
	p.Sync = false
	p.ReRoot = false
	return PrtStateAlternatePort
}

func PrtMachineFSMBuild(p *StpPort) *PrtMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new PrtMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the DISCARD State
	prtm := NewStpPrtMachine(p)

	// BEGIN -> INIT PORT
	rules.AddRule(PrtStateNone, PrtEventBegin, prtm.PrtMachineInitPort)
	// Disabled
	rules.AddRule(PrtStateInitPort, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDisablePort, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventBegin, prtm.PrtMachineInitPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateReRoot, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootForward, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootLearn, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateReRooted, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootPort, PrtEventBegin, prtm.PrtMachineInitPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventBegin, prtm.PrtMachineInitPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateBlockPort, PrtEventBegin, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateBackupPort, PrtEventBegin, prtm.PrtMachineInitPort)

	// UNCONDITIONALFALLTHROUGH -> DISABLEPORT/ROOTPORT/DESIGNATEDPORT/ALTERNATEPORT
	// Disabled
	rules.AddRule(PrtStateInitPort, PrtEventUnconditionallFallThrough, prtm.PrtMachineDisablePort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateReRoot, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootForward, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootLearn, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateReRooted, PrtEventUnconditionallFallThrough, prtm.PrtMachineRootPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventUnconditionallFallThrough, prtm.PrtMachineDesignatedPort)
	// Alternate Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventUnconditionallFallThrough, prtm.PrtMachineAlternatePort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventUnconditionallFallThrough, prtm.PrtMachineAlternatePort)
	rules.AddRule(PrtStateBlockPort, PrtEventUnconditionallFallThrough, prtm.PrtMachineAlternatePort)
	rules.AddRule(PrtStateBackupPort, PrtEventUnconditionallFallThrough, prtm.PrtMachineAlternatePort)

	// NOTLEARNING and NOTFORWARDING and SELECTED and NOT UPDTINFO -> DISABLED PORT
	rules.AddRule(PrtStateDisablePort, PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo, prtm.PrtMachineDisabledPort)
	rules.AddRule(PrtStateBlockPort, PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// FDWHILE NOT EQUAL MAXAGE and SELECTED and NOT UPDTINFO -> DISABLED PORT
	rules.AddRule(PrtStateDisabledPort, PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo, prtm.PrtMachineDisabledPort)

	// SYNC and SELECTED and NOT UPDTINFO  -> DISABLED PORT
	rules.AddRule(PrtStateDisabledPort, PrtEventSyncAndSelectedAndNotUpdtInfo, prtm.PrtMachineDisabledPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSyncAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// REROOT and SELECTED and NOT UPDTINFO  -> DISABLED PORT
	rules.AddRule(PrtStateDisabledPort, PrtEventReRootAndSelectedAndNotUpdtInfo, prtm.PrtMachineDisabledPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventReRootAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// NOTSYNCEDand SELECTED and NOT UPDTINFO  -> DISABLED PORT
	rules.AddRule(PrtStateDisabledPort, PrtEventNotSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineDisabledPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventNotSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// SELECTEDROLE EQUALS DISABLEDPORT and ROLE NOTEQUAL SELECTED ROLE and SELECTED and NOT UPDTINFO  -> DISABLE PORT
	rules.AddRule(PrtStateInitPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDisablePort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateReRoot, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootForward, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootLearn, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateReRooted, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateRootPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateBlockPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)
	rules.AddRule(PrtStateBackupPort, PrtEventSelectedRoleEqualDisabledPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineInitPort)

	// PROPOSED and NOTAGREE and SELECTED and NOT UPDTINFO  -> ROOTPROPOSED
	rules.AddRule(PrtStateRootPort, PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootProposed)

	// ALLSYNCED AND NOTAGREE and SELECTED and NOTUPDTINFO and SELECTED and NOT UPDTINFO -> ROOTAGREED
	rules.AddRule(PrtStateRootPort, PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootAgreed)

	// PROPOSED and AGREE and SELECTED and NOT UPDTINFO -> ROOTAGREED
	rules.AddRule(PrtStateRootPort, PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootAgreed)

	// NOT FORWARD and NOT REROOT and SELECTED and NOT UPDTINFO -> REROOT
	rules.AddRule(PrtStateRootPort, PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo, prtm.PrtMachineReRoot)

	// RRWHILE NOT EQUAL FWDDELAY and SELECTED and NOT UPDTINFO -> ROOT PORT
	rules.AddRule(PrtStateRootPort, PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)

	// REROOT and FORWARDand SELECTED and NOT UPDTINFO  -> REROOTED
	rules.AddRule(PrtStateRootPort, PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineReRooted)

	// FDWHILE EQUAL ZERO and RSTPVERSION and NOTLEARN and SELECTED and NOT UPDTINFO -> ROOTLEARN
	rules.AddRule(PrtStateRootPort, PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootLearn)

	// REROOTED and RBWHILE EQUAL ZERO and RSTPVERSION and NOTLEARN and SELECTED and NOT UPDTINFO -> ROOTLEARN
	rules.AddRule(PrtStateRootPort, PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootLearn)

	// FDWHILE EQUAL ZERO and RSTPVERSION and LEARN and NOT FORWARD and SELECTED and NOT UPDTINFO -> ROOTLEARN
	rules.AddRule(PrtStateRootPort, PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootForward)

	// REROOTED and RBWHILE EQUAL ZERO and RSTPVERSION and LEARN and NOT FORWARD and SELECTED and NOT UPDTINFO -> ROOTLEARN
	rules.AddRule(PrtStateRootPort, PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootForward)

	// SELECTEDROLE EQUALS ROOTPORT and ROLE NOTEQUAL SELECTEDROLE and SELECTED and NOT UPDTINFO -> ROOTPORT
	rules.AddRule(PrtStateInitPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDisablePort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateReRoot, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootForward, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootLearn, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateReRooted, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateRootPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateBlockPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)
	rules.AddRule(PrtStateBackupPort, PrtEventSelectedRoleEqualRootPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineRootPort)

	// NOTFORWARD and NOTAGREED and NOTPROPOSING and NOTOPEREDGEand SELECTED and NOT UPDTINFO -> DESIGNATEDPROPOSE
	rules.AddRule(PrtStateDesignatedPort, PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPropose)

	// NOTLEARNING and NOTFORWARDING and NOTSYNCED and SELECTED and NOT UPDTINFO -> DESIGNATEDSYNCED
	rules.AddRule(PrtStateDesignatedPort, PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedSynced)

	// AGREED and NOTSYNCED and SELECTED and NOT UPDTINFO -> DESIGNATEDSYNCED
	rules.AddRule(PrtStateDesignatedPort, PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedSynced)

	// OPEREDGE and NOTSYNCED and SELECTED and NOT UPDTINFO -> DESIGNATEDSYNCED
	rules.AddRule(PrtStateDesignatedPort, PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedSynced)

	// SYNC and SYNCED and SELECTED and NOT UPDTINFO -> DESIGNATEDSYNCED
	rules.AddRule(PrtStateDesignatedPort, PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedSynced)

	// RRWHILE EQUAL ZERO and REROOT and SELECTED and NOT UPDTINFO  -> DESIGNATEDRETIRED
	rules.AddRule(PrtStateDesignatedPort, PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedRetired)

	// SELECTEDROLE EQUALS DESIGNATEDPORT and ROLE NOT EQUAL SELECTEDROLE and SELECTED and NOT UPDTINFO -> DESIGNATEDPORT
	rules.AddRule(PrtStateInitPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDisablePort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateReRoot, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateRootForward, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateRootLearn, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateReRooted, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateRootPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateBlockPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)
	rules.AddRule(PrtStateBackupPort, PrtEventSelectedRoleEqualDesignatedPortAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedPort)

	// SYNC and NOTSYNCED and NOTOPEREDGE and LEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// SYNC AND NOTSYNCED and NOTOPEREDGE and FORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// REROOT and RRWHILE NOTEQUAL ZERO and NOTOPEREDGE and LEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// REROOT and RRWHILE NOTEQUAL ZERO and NOTOPEREDGE and FORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// DISPUTED and NOTOPEREDGE and LEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// DISPUTED and NOTOPEREDGE and FORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDDISCARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedDiscard)

	// FDWHILE EQUALS ZERO and RRWHILE EQUALS ZERO and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// FDWHILE EQUALS ZERO and NOTREROOT and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// AGREED and RRWHILE EQUALS ZERO and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// AGREED and NOTREROOT and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// OPEREDGE and RRWHILE EQUALS ZERO and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// OPEREDGE and NOTREROOT and NOTSYNC and NOTLEARN and SELECTED and NOT UPDTINFO -> DESIGNATEDLEARN
	rules.AddRule(PrtStateDesignatedPort, PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedLearn)

	// FDWHILE EQUALS ZERO and RRWHILE EQUALS ZERO and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// FDWHILE EQUALS ZERO and NOTREROOT and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// AGREED and RRWHILE EQUALS ZERO and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// AGREED and NOTREROOT and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// OPEREDGE and RRWHILE EQUALS ZERO and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// OPEREDGE and NOTREROOT and NOTSYNC and LEARN and NOTFORWARD and SELECTED and NOT UPDTINFO -> DESIGNATEDFORWARD
	rules.AddRule(PrtStateDesignatedPort, PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, prtm.PrtMachineDesignatedForward)

	// SELECTEDROLE EQUALS ALTERNATE and ROLE NOT EQUAL SELECTEDROLE and SELECTED and NOT UPDTINFO -> ALTERNATEPORT
	rules.AddRule(PrtStateInitPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDisablePort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateReRoot, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootForward, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootLearn, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateReRooted, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateBlockPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateBackupPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)

	// SELECTEDROLE EQUALS BACKUP and ROLE NOT EQUAL SELECTEDROLE and SELECTED and NOT UPDTINFO -> ALTERNATEPORT
	rules.AddRule(PrtStateInitPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDisablePort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDisabledPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Root
	rules.AddRule(PrtStateRootProposed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootAgreed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateReRoot, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootForward, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootLearn, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateReRooted, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateRootPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Designated
	rules.AddRule(PrtStateDesignatedPropose, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedSynced, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedRetired, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedForward, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedLearn, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedDiscard, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateDesignatedPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	// Alternate/Backup
	rules.AddRule(PrtStateAlternateProposed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateAlternateAgreed, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateAlternatePort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateBlockPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)
	rules.AddRule(PrtStateBackupPort, PrtEventSelectedRoleEqualAlternateAndRoleNotEqualSelectedRoleAndSelectedAndNotUpdtInfo, prtm.PrtMachineBlockPort)

	// PROPOSED and NOTAGREE and SELECTED and NOTUPDTINFO > ALTERNATEPROPOSED
	rules.AddRule(PrtStateAlternatePort, PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternateProposed)

	// ALLSYNCED and NOTAGREE and SELECTED and NOTUPDTINFO > ALTERNATEPROPOSED
	rules.AddRule(PrtStateAlternatePort, PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternateAgreed)

	// PROPOSED and AGREE and SELECTED and NOTUPDTINFO > ALTERNATEPROPOSED
	rules.AddRule(PrtStateAlternatePort, PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternateAgreed)

	// FDWHILE NOTEQUAL FORWARDDELAY and SELECTED and NOTUPDTINFO > ALTERNATEPROPOSED
	rules.AddRule(PrtStateAlternatePort, PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// FDWHILE NOTEQUAL FORWARDDELAY and SELECTED and NOTUPDTINFO > ALTERNATEPROPOSED
	rules.AddRule(PrtStateAlternatePort, PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo, prtm.PrtMachineAlternatePort)

	// RBWHILE NOTEQUAL 2*HELLOTIME and ROLE EQUAL BACKUPPORT and SELECTED and NOTUPDTINFO > BACKUPPORT
	rules.AddRule(PrtStateAlternatePort, PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo, prtm.PrtMachineBackupPort)

	// Create a new FSM and apply the rules
	prtm.Apply(&rules)

	return prtm
}

// PimMachineMain:
func (p *StpPort) PrtMachineMain() {

	// Build the State machine for STP Port Role Transitions State Machine according to
	// 802.1d Section 17.29
	prtm := PrtMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	prtm.Machine.Start(prtm.Machine.Curr.PreviousState())

	// lets create a go routing which will wait for the specific events
	go func(m *PrtMachine) {
		StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine Start")
		defer m.p.wg.Done()
		for {
			select {

			case event, ok := <-m.PrtEvents:

				if ok {
					//StpMachineLogger("DEBUG", PrtMachineModuleStr, m.p.IfIndex, m.p.BrgIfIndex, fmt.Sprintf("Event Rx", event.src, event.e))
					if m.Machine.Curr.CurrentState() == PrtStateNone && event.e != PrtEventBegin {
						m.PrtEvents <- event
						break
					}

					rv := m.Machine.ProcessEvent(event.src, event.e, nil)
					if rv != nil {
						StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s src[%s]state[%s]event[%d]\n", rv, event.src, PrtStateStrMap[m.Machine.Curr.CurrentState()], event.e))
					} else {
						// for faster state transitions
						m.ProcessPostStateProcessing()
					}

					if event.responseChan != nil {
						SendResponse(PrtMachineModuleStr, event.responseChan)
					}
				} else {
					StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, "Machine End")
					return
				}

			case ena := <-m.PrtLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(prtm)
}

func (prtm *PrtMachine) NotifyRoleChanged(oldrole PortRole, newrole PortRole) {
	// The following machines need to know about
	// changes in Role State
	// 1) Port State Transitions
	// 2) Topology Change
	p := prtm.p
	if oldrole != newrole {
		if p.TcMachineFsm != nil {
			if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning {
				if p.Role != PortRoleRootPort &&
					p.Role != PortRoleDesignatedPort &&
					!p.Learn &&
					p.Learning &&
					!p.RcvdTc &&
					!p.RcvdTcn &&
					!p.RcvdTcAck &&
					!p.TcProp {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
						src: PrtMachineModuleStr,
					}
				} else if p.Role == PortRoleRootPort &&
					p.Forward &&
					!p.OperEdge {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleEqualRootPortAndForwardAndNotOperEdge,
						src: PrtMachineModuleStr,
					}
				}
			} else if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateActive {
				if p.Role != PortRoleRootPort &&
					p.Role != PortRoleDesignatedPort {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPort,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
	}
}

func (prtm *PrtMachine) NotifyForwardChanged(oldforward bool, newforward bool) {
	p := prtm.p
	if oldforward != newforward {
		// Pst
		if p.PstMachineFsm != nil {
			if p.PstMachineFsm.Machine.Curr.CurrentState() == PstStateLearning {
				if p.Forward {
					p.PstMachineFsm.PstEvents <- MachineEvent{
						e:   PstEventForward,
						src: PrtMachineModuleStr,
					}
				}

			} else if p.PstMachineFsm.Machine.Curr.CurrentState() == PstStateForwarding {
				if !p.Forward {
					p.PstMachineFsm.PstEvents <- MachineEvent{
						e:   PstEventNotForward,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
		// Tc
		if p.TcMachineFsm != nil {
			if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning {
				if p.Role == PortRoleRootPort &&
					p.Forward &&
					!p.OperEdge {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleEqualRootPortAndForwardAndNotOperEdge,
						src: PrtMachineModuleStr,
					}
				} else if p.Role == PortRoleDesignatedPort &&
					p.Forward &&
					!p.OperEdge {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleEqualDesignatedPortAndForwardAndNotOperEdge,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
	}
}

func (prtm *PrtMachine) NotifyLearnChanged(oldlearn bool, newlearn bool) {
	p := prtm.p
	if oldlearn != newlearn {
		// Pst
		if p.PstMachineFsm != nil {
			if p.PstMachineFsm.Machine.Curr.CurrentState() == PstStateDiscarding {
				if p.Learn {
					p.PstMachineFsm.PstEvents <- MachineEvent{
						e:   PstEventLearn,
						src: PrtMachineModuleStr,
					}
				}

			} else if p.PstMachineFsm.Machine.Curr.CurrentState() == PstStateLearning {
				if !p.Learn {
					p.PstMachineFsm.PstEvents <- MachineEvent{
						e:   PstEventNotLearn,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
		// Tc
		if p.TcMachineFsm != nil {
			if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateLearning {
				if p.Role != PortRoleRootPort &&
					p.Role != PortRoleDesignatedPort &&
					!p.Learn &&
					!p.Learning &&
					!p.RcvdTc &&
					!p.RcvdTcn &&
					!p.RcvdTcAck &&
					!p.TcProp {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventRoleNotEqualRootPortAndRoleNotEqualDesignatedPortAndNotLearnAndNotLearningAndNotRcvdTcAndNotRcvdTcnAndNotRcvdTcAckAndNotTcProp,
						src: PrtMachineModuleStr,
					}
				}
			} else if p.TcMachineFsm.Machine.Curr.CurrentState() == TcStateInactive {
				if p.Learn &&
					!p.FdbFlush {
					p.TcMachineFsm.TcEvents <- MachineEvent{
						e:   TcEventLearnAndNotFdbFlush,
						src: PrtMachineModuleStr,
					}
				}
			}
		}
	}
}

func (prtm *PrtMachine) NotifyNewInfoChanged(oldnewinfo bool, newnewinfo bool) {
	p := prtm.p
	if oldnewinfo != newnewinfo {
		if p.PtxmMachineFsm != nil {
			if p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle {
				if p.SendRSTP &&
					p.NewInfo &&
					p.TxCount < p.b.TxHoldCount &&
					p.HelloWhenTimer.count != 0 &&
					p.Selected &&
					!p.UpdtInfo {
					rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo, nil)
					if rv != nil {
						StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo))
					}
				}
			}
		}
	}
}

func (prtm *PrtMachine) ProcessPostStateInitPort() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateInitPort {
		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, fmt.Sprintf("PrtStateInitPort (post) Forwarding[%t] Learning[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Agreed[%t] Agree[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
		//	p.Forwarding, p.Learning, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateRootProposed() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootProposed {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessingPostStateRootAgreed() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootAgreed {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessingPostStateReRoot() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateReRoot {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessingPostStateRootForward() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootForward {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessingPostStateRootLearn() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootLearn {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessingPostStateReRooted() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateReRooted {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}

}

func (prtm *PrtMachine) ProcessPostStateRootPort() {
	p := prtm.p
	b := p.b
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("PrtStateRootPort (post) Forwarding[%t] Forward[%t] Learning[%t] Learn[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Agreed[%t] Agree[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
		//	p.Forwarding, p.Forward, p.Learning, p.Learn, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))
		if p.Proposed &&
			!p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.b.AllSynced() &&
			!p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Proposed &&
			p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if !p.Forward &&
			!p.ReRoot &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotForwardAndNotReRootAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.RrWhileTimer.count != int32(p.PortTimes.ForwardingDelay) &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.ReRoot &&
			p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			p.RstpVersion &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if b.ReRooted(p) &&
			p.RbWhileTimer.count == 0 &&
			p.RstpVersion &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			p.RstpVersion &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if b.ReRooted(p) &&
			p.RbWhileTimer.count == 0 &&
			p.RstpVersion &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateDesignatedPort() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, fmt.Sprintf("PrtStateDesignatedPort (post) Forwarding[%t] Forward[%t] Learning[%t] Learn[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Synced[%t] Sync[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
		//	p.Forwarding, p.Forward, p.Learning, p.Learn, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))
		if !p.Forwarding &&
			!p.Agreed &&
			!p.Proposing &&
			!p.OperEdge &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotForwardAndNotAgreedAndNotProposingAndNotOperEdgeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if !p.Learning &&
			!p.Forwarding &&
			!p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotLearningAndNotForwardingAndNotSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Agreed &&
			!p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAgreedAndNotSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.OperEdge &&
			!p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventOperEdgeAndNotSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Sync &&
			p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.RrWhileTimer.count == 0 &&
			p.ReRoot &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Sync &&
			!p.Synced &&
			!p.OperEdge &&
			p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Sync &&
			!p.Synced &&
			!p.OperEdge &&
			p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.ReRoot &&
			p.RrWhileTimer.count != 0 &&
			!p.OperEdge &&
			p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.ReRoot &&
			p.RrWhileTimer.count != 0 &&
			!p.OperEdge &&
			p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Disputed &&
			!p.OperEdge &&
			p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventDisputedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Disputed &&
			!p.OperEdge &&
			p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventDisputedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			!p.ReRoot &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Agreed &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Agreed &&
			!p.ReRoot &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAgreedAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.OperEdge &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.OperEdge &&
			!p.ReRoot &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventOperEdgeAndNotReRootAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count == 0 &&
			!p.ReRoot &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Agreed &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Agreed &&
			!p.ReRoot &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAgreedAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.OperEdge &&
			p.RrWhileTimer.count == 0 &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.OperEdge &&
			!p.ReRoot &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventOperEdgeAndNotReRootAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateDesignatedPropose() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPropose {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}
func (prtm *PrtMachine) ProcessingPostStateDesignatedSynced() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedSynced {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}
func (prtm *PrtMachine) ProcessingPostStateDesignatedRetired() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedRetired {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}
func (prtm *PrtMachine) ProcessingPostStateDesignatedForward() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedForward {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}
func (prtm *PrtMachine) ProcessingPostStateDesignatedLearn() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedLearn {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}
func (prtm *PrtMachine) ProcessingPostStateDesignatedDiscard() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedDiscard {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateBlockedPort() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBlockPort {
		if !p.Learning &&
			!p.Forwarding &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateAlternateProposed() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternateProposed {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateAlternateAgreed() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternateAgreed {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateBackupPort() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateBackupPort {
		rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventUnconditionallFallThrough, nil)
		if rv != nil {
			StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventUnconditionallFallThrough))
		} else {
			prtm.ProcessPostStateProcessing()
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateAlternatePort() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
		if p.Proposed &&
			!p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventProposedAndNotAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.b.AllSynced() &&
			!p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventAllSyncedAndNotAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Proposed &&
			p.Agree &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventProposedAndAgreeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.FdWhileTimer.count != int32(p.PortTimes.ForwardingDelay) &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Sync &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventSyncAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventSyncAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.ReRoot &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if !p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.RbWhileTimer.count != int32(2*p.PortTimes.HelloTime) &&
			p.Role == PortRoleBackupPort &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateDisable() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		prtm.Machine.Curr.CurrentState() == PrtStateDisablePort {
		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, fmt.Sprintf("PrtStateDisablePort (post) Forwarding[%t] Learning[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Agreed[%t] Agree[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
		//	p.Forwarding, p.Learning, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))

		if !p.Learning &&
			!p.Forwarding &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotLearningAndNotForwardingAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessingPostStateDisabled() {
	p := prtm.p
	if p.PrtMachineFsm != nil &&
		prtm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
		//StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, fmt.Sprintf("PrtStateDisabledPort (post) Forwarding[%t] Learning[%t] Agreed[%t] Agree[%t]\nProposing[%t] OperEdge[%t] Agreed[%t] Agree[%t]\nReRoot[%t] Selected[%t], UpdtInfo[%t] Fdwhile[%d] rrWhile[%d]\n",
		//	p.Forwarding, p.Learning, p.Agreed, p.Agree, p.Proposing, p.OperEdge, p.Synced, p.Sync, p.ReRoot, p.Selected, p.UpdtInfo, p.FdWhileTimer.count, p.RrWhileTimer.count))
		if p.FdWhileTimer.count != int32(p.PortTimes.MaxAge) &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.Sync &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventSyncAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventSyncAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if p.ReRoot &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventReRootAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventReRootAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		} else if !p.Synced &&
			p.Selected &&
			!p.UpdtInfo {
			rv := prtm.Machine.ProcessEvent(PrtMachineModuleStr, PrtEventNotSyncedAndSelectedAndNotUpdtInfo, nil)
			if rv != nil {
				StpMachineLogger("ERROR", PrtMachineModuleStr, p.IfIndex, p.BrgIfIndex, fmt.Sprintf("%s post state[%s]event[%d]\n", rv, PrtStateStrMap[prtm.Machine.Curr.CurrentState()], PrtEventNotSyncedAndSelectedAndNotUpdtInfo))
			} else {
				prtm.ProcessPostStateProcessing()
			}
		}
	}
}

func (prtm *PrtMachine) ProcessPostStateProcessing() {
	// Disabled states
	prtm.ProcessPostStateInitPort()
	prtm.ProcessingPostStateDisable()
	prtm.ProcessingPostStateDisabled()
	// Root states
	prtm.ProcessPostStateRootPort()
	prtm.ProcessingPostStateRootProposed()
	prtm.ProcessingPostStateRootAgreed()
	prtm.ProcessingPostStateReRoot()
	prtm.ProcessingPostStateRootForward()
	prtm.ProcessingPostStateRootLearn()
	prtm.ProcessingPostStateReRooted()
	// Designated states
	prtm.ProcessingPostStateDesignatedPort()
	prtm.ProcessingPostStateDesignatedPropose()
	prtm.ProcessingPostStateDesignatedSynced()
	prtm.ProcessingPostStateDesignatedRetired()
	prtm.ProcessingPostStateDesignatedForward()
	prtm.ProcessingPostStateDesignatedLearn()
	prtm.ProcessingPostStateDesignatedDiscard()
	// Alternate and Backup states
	prtm.ProcessingPostStateBlockedPort()
	prtm.ProcessingPostStateAlternatePort()
	prtm.ProcessingPostStateAlternateProposed()
	prtm.ProcessingPostStateAlternateAgreed()
	prtm.ProcessingPostStateBackupPort()
}

func (prtm *PrtMachine) NotifyReRootChanged(oldreroot bool, newreroot bool) {
	p := prtm.p
	// only need to handle reroot == true cases
	// because this is triggered via setReRootTree which needs to notify
	// all other ports
	if oldreroot != newreroot &&
		p.PrtMachineFsm != nil {
		if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
			if p.ReRoot &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
			if p.ReRoot &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndForwardAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
			if p.RrWhileTimer.count == 0 &&
				p.ReRoot &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if p.ReRoot &&
				p.RrWhileTimer.count != 0 &&
				!p.OperEdge &&
				p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if p.ReRoot &&
				p.RrWhileTimer.count != 0 &&
				!p.OperEdge &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
			if p.ReRoot &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventReRootAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		}
	}

}

func (prtm *PrtMachine) NotifySyncChanged(oldsync bool, newsync bool) {
	p := prtm.p
	// only need to handle sync == true cases
	// because this is triggered via setSyncTree which needs to notify
	// all other ports
	if oldsync != newsync &&
		p.PrtMachineFsm != nil {
		/*StpMachineLogger("DEBUG", PrtMachineModuleStr, p.IfIndex, fmt.Sprintf("notifySyncChanged: state[%s] synced[%t] operedge[%t] learn[%t] forward[%t] selected[%t] updtInfo[%t]",
		PrtStateStrMap[p.PrtMachineFsm.Machine.Curr.CurrentState()],
		p.Synced,
		p.OperEdge,
		p.Learn,
		p.Forward,
		p.Selected,
		p.UpdtInfo))*/
		if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort {
			if p.Sync &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
			if p.Sync &&
				p.Synced &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndSyncedAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if p.Sync &&
				!p.Synced &&
				!p.OperEdge &&
				p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if p.Sync &&
				!p.Synced &&
				!p.OperEdge &&
				p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndNotSyncedAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort {
			if p.Sync &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventSyncAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		}
	}
}

func (prtm *PrtMachine) setSyncTree(ifindex int32) {
	b := prtm.p.b
	var p *StpPort
	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			// skip calling ifindex because post state processing will handle
			// its change.
			if ifindex != pId {
				defer p.PrtMachineFsm.NotifySyncChanged(p.Sync, true)
			}
			p.Sync = true
		}
	}
}

func (prtm *PrtMachine) setReRootTree(ifindex int32) {
	b := prtm.p.b
	var p *StpPort
	for _, pId := range b.StpPorts {
		if StpFindPortByIfIndex(pId, b.BrgIfIndex, &p) {
			// skip calling ifindex because post state processing will handle
			// its change.
			if ifindex != pId {
				defer p.PrtMachineFsm.NotifyReRootChanged(p.ReRoot, true)
			}
			p.ReRoot = true
		}
	}
}
