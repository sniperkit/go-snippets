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

// timers
package stp

import (
	//"fmt"
	"time"
)

type TimerType int

const (
	TimerTypeEdgeDelayWhile TimerType = iota
	TimerTypeFdWhile
	TimerTypeHelloWhen
	TimerTypeMdelayWhile
	TimerTypeRbWhile
	TimerTypeRcvdInfoWhile
	TimerTypeRrWhile
	TimerTypeTcWhile
	TimerTypeBAWhile
)

var TimerTypeStrMap map[TimerType]string

func TimerTypeStrStateMapInit() {
	TimerTypeStrMap = make(map[TimerType]string)
	TimerTypeStrMap[TimerTypeEdgeDelayWhile] = "Edge Delay Timer"
	TimerTypeStrMap[TimerTypeFdWhile] = "Fd While Timer"
	TimerTypeStrMap[TimerTypeHelloWhen] = "Hello When Timer"
	TimerTypeStrMap[TimerTypeMdelayWhile] = "Migration Delay Timer"
	TimerTypeStrMap[TimerTypeRbWhile] = "Recent Backup Timer"
	TimerTypeStrMap[TimerTypeRcvdInfoWhile] = "Received Info Timer"
	TimerTypeStrMap[TimerTypeRrWhile] = "Recent Root Timer"
	TimerTypeStrMap[TimerTypeTcWhile] = "topology Change Timer"
	TimerTypeStrMap[TimerTypeBAWhile] = "Bridge Assurance While Timer"
}

// TickTimerStart: Port Timers Tick timer
func (m *PtmMachine) TickTimerStart() {

	if m.TickTimer == nil {
		m.TickTimer = time.NewTimer(time.Second * 1)
	} else {
		m.TickTimer.Reset(time.Second * 1)
	}
}

// TickTimerStop
// Stop the running timer
func (m *PtmMachine) TickTimerStop() {

	if m.TickTimer != nil {
		m.TickTimer.Stop()
	}
}

func (m *PtmMachine) TickTimerDestroy() {
	m.TickTimerStop()
	m.TickTimer = nil
}

func (p *StpPort) ResetTimerCounters(counterType TimerType) {
	// TODO
}

func (p *StpPort) DecrementTimerCounters() {
	/*StpMachineLogger("DEBUG", "PTIM", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("EdgeDelayWhile[%d] FdWhileTimer[%d] HelloWhen[%d] MdelayWhile[%d] RbWhile[%d] RcvdInfoWhile[%d] RrWhile[%d] TcWhile[%d] txcount[%d]",
	p.EdgeDelayWhileTimer.count,
	p.FdWhileTimer.count,
	p.HelloWhenTimer.count,
	p.MdelayWhiletimer.count,
	p.RbWhileTimer.count,
	p.RcvdInfoWhiletimer.count,
	p.RrWhileTimer.count,
	p.TcWhileTimer.count,
	p.TxCount))
	*/
	// 17.19.44
	if p.TxCount > 0 {
		p.TxCount--
	}

	// ed owner
	if p.EdgeDelayWhileTimer.count > 0 {
		p.EdgeDelayWhileTimer.count--

		if p.EdgeDelayWhileTimer.count == 0 {
			defer p.NotifyEdgeDelayWhileTimerExpired()
		} else {
			if p.PrxmMachineFsm != nil &&
				p.EdgeDelayWhileTimer.count != MigrateTimeDefault &&
				!p.PortEnabled {
				p.PrxmMachineFsm.PrxmEvents <- MachineEvent{
					e:   PrxmEventEdgeDelayWhileNotEqualMigrateTimeAndNotPortEnabled,
					src: PrtMachineModuleStr,
				}
			}
		}
	}
	// Prt owner
	if p.FdWhileTimer.count > 0 {
		if p.PrtMachineFsm != nil {

			if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDisabledPort &&
				uint16(p.FdWhileTimer.count) != p.b.BridgeTimes.MaxAge &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileNotEqualMaxAgeAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort &&
				uint16(p.FdWhileTimer.count) != p.b.BridgeTimes.ForwardingDelay &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileNotEqualForwardDelayAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		}
		p.FdWhileTimer.count--

		if p.FdWhileTimer.count == 0 {
			defer p.NotifyFdWhileTimerExpired()
		}
	}
	// ptx owner
	if p.HelloWhenTimer.count > 0 {
		// if hellowhen never expires then packets should not get transmitted
		if p.PortEnabled {
			p.HelloWhenTimer.count--

			if p.HelloWhenTimer.count == 0 {
				defer p.NotifyHelloWhenTimerExpired()
			} else {
				if p.PtxmMachineFsm != nil &&
					p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle {
					if p.SendRSTP &&
						p.NewInfo &&
						(p.TxCount < p.b.TxHoldCount) &&
						p.Selected &&
						!p.UpdtInfo {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventSendRSTPAndNewInfoAndTxCountLessThanTxHoldCoundAndHelloWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: PtxmMachineModuleStr,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleRootPort &&
						(p.TxCount < p.b.TxHoldCount) &&
						p.Selected &&
						!p.UpdtInfo {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndRootPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: PtxmMachineModuleStr,
						}
					} else if !p.SendRSTP &&
						p.NewInfo &&
						p.Role == PortRoleDesignatedPort &&
						(p.TxCount < p.b.TxHoldCount) &&
						p.Selected &&
						!p.UpdtInfo {
						p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
							e:   PtxmEventNotSendRSTPAndNewInfoAndDesignatedPortAndTxCountLessThanTxHoldCountAndHellWhenNotEqualZeroAndSelectedAndNotUpdtInfo,
							src: PtxmMachineModuleStr,
						}
					}
				}
			}
		} else {
			p.HelloWhenTimer.count = int32(p.PortTimes.HelloTime)
		}

	}

	// ppm owner
	if p.MdelayWhiletimer.count > 0 {
		p.MdelayWhiletimer.count--

		if p.MdelayWhiletimer.count == 0 {
			defer p.NotifyMdelayWhileTimerExpired()
		}
	}
	// prt owner
	if p.RbWhileTimer.count > 0 {
		// this case should reset the rbwhiletimer
		if p.PrtMachineFsm != nil &&
			p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateAlternatePort &&
			p.Role == PortRoleBackupPort &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventRbWhileNotEqualTwoTimesHelloTimeAndRoleEqualsBackupPortAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		}
		p.RbWhileTimer.count--

		if p.RbWhileTimer.count == 0 {
			defer p.NotifyRbWhileTimerExpired()
		}
	}
	// pi owner
	if p.RcvdInfoWhiletimer.count > 0 {
		p.RcvdInfoWhiletimer.count--

		if p.RcvdInfoWhiletimer.count == 0 {
			defer p.NotifyRcvdInfoWhileTimerExpired()
		}
	}
	// prt owner
	if p.RrWhileTimer.count > 0 {
		p.RrWhileTimer.count--
		if p.PrtMachineFsm != nil &&
			p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {

			if p.RrWhileTimer.count != int32(p.b.RootTimes.ForwardingDelay) &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventRrWhileNotEqualFwdDelayAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
			// lets just reset the rrwhile count which is normally done based on
			// transition to root port state, but in order to not have
			// root port states constantly transition to root we will just
			// set this here
			//p.RrWhileTimer.count = int32(p.PortTimes.ForwardingDelay)
		} else {
			if p.RrWhileTimer.count != 0 &&
				p.PrtMachineFsm != nil &&
				p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
				if p.ReRoot &&
					!p.OperEdge &&
					p.Learn &&
					p.Selected &&
					!p.UpdtInfo {
					p.PrtMachineFsm.PrtEvents <- MachineEvent{
						e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndLearnAndSelectedAndNotUpdtInfo,
						src: PrtMachineModuleStr,
					}
				} else if p.ReRoot &&
					!p.OperEdge &&
					p.Forward &&
					p.Selected &&
					!p.UpdtInfo {
					p.PrtMachineFsm.PrtEvents <- MachineEvent{
						e:   PrtEventReRootAndRrWhileNotEqualZeroAndNotOperEdgeAndForwardAndSelectedAndNotUpdtInfo,
						src: PrtMachineModuleStr,
					}
				}
			}
		}

		if p.RrWhileTimer.count == 0 {
			defer p.NotifyRrWhileTimerExpired()
		}
	}
	// tc owner
	if p.TcWhileTimer.count > 0 {
		p.TcWhileTimer.count--

		if p.TcWhileTimer.count == 0 {
			defer p.NotifyTcWhileTimerExpired()
		}
	}
	// Bridge Assurance
	if p.BridgeAssurance &&
		!p.OperEdge &&
		p.PortEnabled &&
		p.BAWhileTimer.count > 0 &&
		p.RcvdBPDU {
		p.BAWhileTimer.count--

		if p.BAWhileTimer.count == 0 {
			p.BridgeAssuranceInconsistant = true
			p.NotifySelectedRoleChanged("BAM", p.SelectedRole, PortRoleDisabledPort)
			p.SelectedRole = PortRoleDisabledPort
		}
	}

	if p.BpduGuard &&
		p.OperEdge &&
		p.BPDUGuardTimer.count > 0 {
		p.BPDUGuardTimer.count--
		// condition has not been detected lets clear the
		// Detection
		if p.BPDUGuardTimer.count == 0 {
			for _, client := range GetAsicDPluginList() {
				client.BPDUGuardDetected(p.IfIndex, false)
			}
		}
	}
}

func (p *StpPort) NotifyEdgeDelayWhileTimerExpired() {

	if p.BdmMachineFsm != nil &&
		p.BdmMachineFsm.Machine.Curr.CurrentState() == BdmStateNotEdge &&
		p.AutoEdgePort &&
		p.SendRSTP &&
		p.Proposing {
		p.BdmMachineFsm.BdmEvents <- MachineEvent{
			e:   BdmEventEdgeDelayWhileEqualZeroAndAutoEdgeAndSendRSTPAndProposing,
			src: BdmMachineModuleStr,
		}

	}
}

func (p *StpPort) NotifyFdWhileTimerExpired() {
	//StpMachineLogger("DEBUG", "TIMER", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("FdWhileTimerExpired: PrtState[%s] RstpVersion[%t] Learn[%t] Forward[%t] Sync[%t] reroot[%t] rrwhile[%d] selected[%t] updtInfo[%t]",
	//	PrtStateStrMap[p.PrtMachineFsm.Machine.Curr.CurrentState()], p.RstpVersion, p.Learn, p.Forward, p.Sync, p.ReRoot, p.RrWhileTimer.count, p.Selected, p.UpdtInfo))
	if p.PrtMachineFsm != nil {
		if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
			// events from Figure 17-21
			if p.RstpVersion &&
				!p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
					src: PtmMachineModuleStr,
				}
			} else if p.RstpVersion &&
				p.Learn &&
				!p.Forward {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
					src: PtmMachineModuleStr,
				}
			}
		} else if p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
			// events from Figure 17-22
			if p.RrWhileTimer.count == 0 &&
				!p.Sync &&
				!p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}

			} else if !p.ReRoot &&
				!p.Sync &&
				!p.Learn &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndNotLearnSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}

			} else if p.RrWhileTimer.count == 0 &&
				!p.Sync &&
				p.Learn &&
				!p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			} else if !p.ReRoot &&
				!p.Sync &&
				p.Learn &&
				!p.Forward &&
				p.Selected &&
				!p.UpdtInfo {
				p.PrtMachineFsm.PrtEvents <- MachineEvent{
					e:   PrtEventFdWhileEqualZeroAndNotReRootAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
					src: PrtMachineModuleStr,
				}
			}
		}
	}
}

func (p *StpPort) NotifyHelloWhenTimerExpired() {

	if p.PtxmMachineFsm != nil &&
		p.PtxmMachineFsm.Machine.Curr.CurrentState() == PtxmStateIdle {
		if p.Selected &&
			!p.UpdtInfo {
			p.PtxmMachineFsm.PtxmEvents <- MachineEvent{
				e:   PtxmEventHelloWhenEqualsZeroAndSelectedAndNotUpdtInfo,
				src: PtxmMachineModuleStr,
			}
		}
	}
}

func (p *StpPort) NotifyMdelayWhileTimerExpired() {

	if p.PpmmMachineFsm != nil &&
		p.PpmmMachineFsm.Machine.Curr.CurrentState() == PpmmStateCheckingRSTP ||
		p.PpmmMachineFsm.Machine.Curr.CurrentState() == PpmmStateSelectingSTP {
		p.PpmmMachineFsm.PpmmEvents <- MachineEvent{
			e:   PpmmEventMdelayWhileEqualZero,
			src: PpmmMachineModuleStr,
		}
	}
}

func (p *StpPort) NotifyRbWhileTimerExpired() {

	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateRootPort {
		if p.b.ReRooted(p) &&
			p.RstpVersion &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndNotLearnAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.b.ReRooted(p) &&
			p.RstpVersion &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventReRootedAndRbWhileEqualZeroAndRstpVersionAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		}
	}
}

func (p *StpPort) NotifyRcvdInfoWhileTimerExpired() {
	if p.PimMachineFsm != nil &&
		p.PimMachineFsm.Machine.Curr.CurrentState() == PimStateCurrent {
		if p.InfoIs == PortInfoStateReceived &&
			!p.UpdtInfo &&
			!p.RcvdMsg {
			p.PimMachineFsm.PimEvents <- MachineEvent{
				e:   PimEventInflsEqualReceivedAndRcvdInfoWhileEqualZeroAndNotUpdtInfoAndNotRcvdMsg,
				src: PimMachineModuleStr,
			}
		}
	}
}

func (p *StpPort) NotifyRrWhileTimerExpired() {
	if p.PrtMachineFsm != nil &&
		p.PrtMachineFsm.Machine.Curr.CurrentState() == PrtStateDesignatedPort {
		if p.ReRoot &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventRrWhileEqualZeroAndReRootAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.FdWhileTimer.count == 0 &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}

		} else if p.Agreed &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.OperEdge &&
			!p.Sync &&
			!p.Learn &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndNotLearnAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.FdWhileTimer.count == 0 &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventFdWhileEqualZeroAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.Agreed &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventAgreedAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		} else if p.OperEdge &&
			!p.Sync &&
			p.Learn &&
			!p.Forward &&
			p.Selected &&
			!p.UpdtInfo {
			p.PrtMachineFsm.PrtEvents <- MachineEvent{
				e:   PrtEventOperEdgeAndRrWhileEqualZeroAndNotSyncAndLearnAndNotForwardAndSelectedAndNotUpdtInfo,
				src: PrtMachineModuleStr,
			}
		}
	}

}
func (p *StpPort) NotifyTcWhileTimerExpired() {

}
