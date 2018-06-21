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
	"l3/vrrp/debug"
	"l3/vrrp/packet"
	"time"
)

func (f *FSM) transitionToBackup(advInt int32) {
	debug.Logger.Debug(FSM_PREFIX, "advertisement timer to be used in backup state for",
		"calculating master down timer is ", f.Config.AdvertisementInterval)
	// Bring Down Sub-Interface only if previous state is same as current state
	if f.previousState == f.State {
		debug.Logger.Debug(FSM_PREFIX, "Disabling VirtualIp as interface:", f.Config.IntfRef, "is acting as backup")
		f.updateVirtualIP(false /*enable*/)
	}
	// Re-Calculate Down timer value
	f.calculateDownValue(advInt)
	// Set/Reset Master Down Timer
	f.handleMasterDownTimer()
	//(165) + Transition to the {Backup} state
	f.State = VRRP_BACKUP_STATE
	f.updateCurrentState()
}

func (f *FSM) backup(decodeInfo *DecodedInfo) {
	pktInfo := decodeInfo.PktInfo
	hdr := pktInfo.Hdr
	/* @TODO:
	   (305) - If the protected IPvX address is an IPv4 address, then:
	   (310) + MUST NOT respond to ARP requests for the IPv4 address(es) associated with the virtual router.
	   (315) - else // protected addr is IPv6
	   (320) + MUST NOT respond to ND Neighbor Solicitation messages for the IPv6 address(es) associated with the virtual router.
	   (325) + MUST NOT send ND Router Advertisement messages for the virtual router.
	   (330) -endif // was protected addr IPv4?
	*/
	// Check dmac address from the inPacket and if it is same discard the packet
	if pktInfo.DstMac == f.VirtualMACAddress {
		debug.Logger.Err("DMAC is equal to VMac and hence discarding the packet")
		return
	}
	// MUST NOT accept packets addressed to the IPvX address(es)
	// associated with the virtual router. @TODO: check with Hari
	if pktInfo.DstIp == f.ipAddr {
		debug.Logger.Err("dst ip is equal to interface ip, dropping the packet")
		return
	}
	//(420) - If an ADVERTISEMENT is received, then:
	if hdr.Type == packet.VRRP_PKT_TYPE_ADVERTISEMENT {
		f.updateRxStInfo(pktInfo)
		// (425) + If the Priority in the ADVERTISEMENT is zero, then:
		if hdr.Priority == 0 {
			//(430) * Set the Master_Down_Timer to Skew_Time
			f.MasterDownValue = f.SkewTime
			f.handleMasterDownTimer()
		} else { // (440) priority non-zero
			/*
			 *	(445) * If Preempt_Mode is False, or if the Priority in the
			 *	ADVERTISEMENT is greater than or equal to the local
			 *	Priority, then:
			 */
			if f.Config.PreemptMode == false || int32(hdr.Priority) >= f.Config.Priority {
				/*
				 * (450) @ Set Master_Adver_Interval to Adver Interval contained in the ADVERTISEMENT
				 * (460) @ Reset the Master_Down_Timer to Master_Down_Interval
				 * (455) @ Recompute the Master_Down_Interval
				 *
				 * api used will be transitionToBackup() which will do the exact
				 * things mentioned above, sorry if you think the naming doesn't
				 * sound correct
				 */
				f.transitionToBackup(int32(hdr.MaxAdverInt))
			} else { //     (465) * else // preempt was true or priority was less
				//          (470) @ Discard the ADVERTISEMENT
			} // endif preempt test
		} // endif was priority zero
	} // endif was advertisement received
	// end BACKUP STATE
}

func (f *FSM) handleMasterDownTimer() {
	if f.MasterDownTimer != nil {
		f.MasterDownTimer.Reset(time.Duration(f.MasterDownValue) * time.Second)
	} else {
		var MasterDownTimer_func func()
		// On Timer expiration we will transition to master
		MasterDownTimer_func = func() {
			debug.Logger.Info(FSM_PREFIX, "master down timer expired..transition to Master")
			f.previousState = f.State
			f.transitionToMaster()
		}
		debug.Logger.Info(FSM_PREFIX, "setting down timer to", f.MasterDownValue)
		// Set Timer expire func...
		f.MasterDownTimer = time.AfterFunc(time.Duration(f.MasterDownValue)*time.Second, MasterDownTimer_func)
	}
}
