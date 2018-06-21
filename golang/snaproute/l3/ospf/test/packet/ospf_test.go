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

/*
Ospf_packet.go
*/
package packettest

import (
	//  "encoding/hex"
	"fmt"
	"l3/ospf/config"
	"l3/ospf/server"
	"net"
	"testing"
	"utils/logging"
)

func _TestOSPFDBDecode(t *testing.T) {
	logger, _ := logging.NewLogger("ospfd", "OSPF", true)
	ospfServer := server.NewOSPFServer(logger)
	ospfHeader := server.NewOspfHdrMetadata()
	ipHdrMd := server.NewIpHdrMetadata()
	ifkey := server.IntfConfKey{
		IPAddr:  config.IpAddress("10.1.1.1"),
		IntfIdx: 1,
	}

	srcmac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	data := []byte{0xc2, 0x03, 0x56, 0xcc, 0x00, 0x01, 0xc2, 0x00, 0x56, 0xab, 0x00, 0x00, 0x08, 0x00, 0x45, 0xc0, 0x00, 0x40, 0x00, 0x7b, 0x00, 0x00, 0x01, 0x59, 0x7c, 0x28, 0x0a, 0x00, 0x14, 0x02, 0x0a, 0x00, 0x14, 0x01, 0x02, 0x02, 0x00, 0x20, 0x05, 0x05, 0x05, 0x05, 0x00, 0x00, 0x00, 0x14, 0x87, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0xdc, 0x52, 0x07, 0x00, 0x00, 0x14, 0x92, 0xff, 0xf6, 0x00, 0x03, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01}
	fmt.Printf("Decode DB: Start")
	ospfServer.ProcessRxDbdPkt(data, ospfHeader, ipHdrMd, ifkey, srcmac)
	fmt.Printf("Decode DB: Success")
}
