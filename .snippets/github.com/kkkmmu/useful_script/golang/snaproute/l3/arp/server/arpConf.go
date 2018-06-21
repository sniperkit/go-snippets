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
)

func (server *ARPServer) processResolveIPv4(conf ResolveIPv4) {
	server.logger.Debug(fmt.Sprintln("Received ResolveIPv4 call for TargetIP:", conf.TargetIP, "ifIndex:", conf.IfId))
	if conf.TargetIP == "0.0.0.0" {
		return
	}
	IfIndex := conf.IfId
	_, exist := server.l3IntfPropMap[IfIndex]
	if !exist {
		server.logger.Err("No such L3 Interface exist", IfIndex)
		return
	}
	server.updateL3ArpEntry(conf.TargetIP, true, IfIndex)
	server.sendArpReq(conf.TargetIP, IfIndex)
}

func (server *ARPServer) processDeleteResolvedIPv4(ipAddr string) {
	server.logger.Info("Delete Resolved IPv4 for ipAddr:", ipAddr)
	server.arpDeleteArpEntryFromRibCh <- ipAddr
}

func (server *ARPServer) processArpConf(conf ArpConf) (int, error) {
	server.logger.Debug("Received ARP Timeout Value via Configuration:", conf.RefTimeout)
	if conf.RefTimeout < server.MinRefreshTimeout {
		server.logger.Err("Refresh Timeout is below minimum allowed refresh timeout value of:", server.MinRefreshTimeout)
		err := errors.New("Invalid Timeout Value")
		return 0, err
	} else if conf.RefTimeout == server.ConfRefreshTimeout {
		server.logger.Err("Arp is already configured with Refresh Timeout Value of:", server.ConfRefreshTimeout, "(seconds)")
		return 0, nil
	}

	server.timeoutCounter = conf.RefTimeout / server.timerGranularity
	server.arpEntryCntUpdateCh <- server.timeoutCounter
	return 0, nil
}

func (server *ARPServer) processArpAction(msg ArpActionMsg) {
	server.logger.Info("Processing Arp Action msg", msg)
	server.arpActionProcessCh <- msg
}

func (server *ARPServer) processDeleteArpEntryInt(entry *DeleteArpEntry) {
	server.logger.Info("Delete arp for ipAddr:", entry.IpAddr, "received from protocol running on flexswitch")
	server.arpDeleteArpEntryIntCh <- entry.IpAddr
}

func (server *ARPServer) processGarp(info *GarpEntry) {
	go server.SendGarp(info.IfName, info.MacAddr, info.IpAddr)
}
