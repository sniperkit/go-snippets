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

package server

import (
	"errors"
	"fmt"
	"time"
)

func (server *ARPServer) GetArpEntry(ipAddr string) (arpState ArpState, err error) {
	arpEnt, exist := server.arpCache[ipAddr]
	if !exist {
		err = errors.New(fmt.Sprintln("Unable to find Arp entry for given IP:", ipAddr))
		return arpState, err
	}

	arpState.IpAddr = ipAddr
	if arpEnt.MacAddr != "incomplete" {
		arpState.MacAddr = arpEnt.MacAddr
		arpState.Intf = arpEnt.IfName
		arpState.VlanId = arpEnt.VlanId
		curTime := time.Now()
		expiryTime := time.Duration(server.timerGranularity*server.timeoutCounter) * time.Second
		timeElapsed := curTime.Sub(arpEnt.TimeStamp)
		timeLeft := expiryTime - timeElapsed
		arpState.ExpiryTimeLeft = timeLeft.String()
	} else {
		arpState.MacAddr = arpEnt.MacAddr
		arpState.Intf = "N/A"
		arpState.VlanId = -1
		arpState.ExpiryTimeLeft = "N/A"
	}
	return arpState, nil
}

func (server *ARPServer) GetLinuxArpEntry(ipAddr string) (arpState ArpLinuxState, err error) {
	arpCache := GetLinuxArpCache()

	for _, arpEnt := range arpCache {
		if arpEnt.IpAddr == ipAddr {
			arpState.IpAddr = ipAddr
			if arpEnt.Flags == "0x0" {
				arpState.MacAddr = "incomplete"
				arpState.HWType = "N/A"
			} else {
				arpState.MacAddr = arpEnt.MacAddr
				if arpEnt.HWType == "0x1" {
					arpState.HWType = "ether"
				} else {
					arpState.HWType = "non-ether"
				}
			}
			arpState.IfName = arpEnt.IfName
			return arpState, nil
		}
	}

	err = errors.New(fmt.Sprintln("Unable to find Arp entry for given IP:", ipAddr))
	return arpState, err
}
