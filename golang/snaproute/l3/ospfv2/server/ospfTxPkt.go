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
	"github.com/google/gopacket/pcap"
)

func (server *OSPFV2Server) SendOspfPkt(key IntfConfKey, ospfPkt []byte) error {
	entry, _ := server.IntfConfMap[key]
	handle := entry.txHdl.SendPcapHdl
	if handle == nil {
		server.logger.Err("Invalid pcap handle")
		err := errors.New("Invalid pcap handle")
		return err
	}
	entry.txHdl.SendMutex.Lock()
	err := handle.WritePacketData(ospfPkt)
	entry.txHdl.SendMutex.Unlock()
	return err
}

func (server *OSPFV2Server) InitTxPkt(intfKey IntfConfKey) error {
	intfEnt, _ := server.IntfConfMap[intfKey]
	ifName := intfEnt.IfName
	sendHdl, err := pcap.OpenLive(ifName, snapshotLen, promiscuous, pcapTimeout)
	if err != nil {
		server.logger.Err("Error opening pcap handler on ", ifName)
		return err
	}
	intfEnt.txHdl.SendPcapHdl = sendHdl
	server.IntfConfMap[intfKey] = intfEnt
	return nil
}

func (server *OSPFV2Server) DeinitTxPkt(intfKey IntfConfKey) {
	intfEnt, _ := server.IntfConfMap[intfKey]
	intfEnt.txHdl.SendPcapHdl.Close()
	intfEnt.txHdl.SendPcapHdl = nil
	server.IntfConfMap[intfKey] = intfEnt
}
