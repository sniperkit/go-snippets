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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ArpLinuxState struct {
	IpAddr  string
	HWType  string
	MacAddr string
	IfName  string
}

type ArpLinuxEntry struct {
	IpAddr  string
	HWType  string
	Flags   string
	MacAddr string
	Mask    string
	IfName  string
}

const (
	f_IPAddr int = iota
	f_HWType
	f_Flags
	f_HWAddr
	f_HWMask
	f_IfName
)

func GetLinuxArpCache() []ArpLinuxEntry {
	fp, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil
	}

	defer fp.Close()

	s := bufio.NewScanner(fp)
	s.Scan() // Skip the field description
	var arpLinuxEntry = make([]ArpLinuxEntry, 0)
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		arpEnt := ArpLinuxEntry{
			IpAddr:  fields[f_IPAddr],
			HWType:  fields[f_HWType],
			Flags:   fields[f_Flags],
			MacAddr: fields[f_HWAddr],
			Mask:    fields[f_HWMask],
			IfName:  fields[f_IfName],
		}
		arpLinuxEntry = append(arpLinuxEntry, arpEnt)
	}
	return arpLinuxEntry
}

func (server *ARPServer) FlushLinuxArpCache() {
	server.logger.Info("Flushing linux arp")
	server.logger.Info(fmt.Sprintln("L3 Property:", server.l3IntfPropMap))
	arpEntry := GetLinuxArpCache()
	for _, arpEnt := range arpEntry {
		if arpEnt.Flags != "0x0" &&
			arpEnt.HWType == "0x1" {
			for _, ent := range server.l3IntfPropMap {
				if arpEnt.IfName == ent.IfName {
					server.deleteLinuxArp(arpEnt.IpAddr)
					break
				}

			}
		}
	}
}

func (server *ARPServer) deleteLinuxArp(ipAddr string) {
	cmd := exec.Command("arp", "-d", ipAddr)
	if err := cmd.Run(); err != nil {
		server.logger.Err(fmt.Sprintln("Error deleting linux arp entry for ", ipAddr))
	}

}
