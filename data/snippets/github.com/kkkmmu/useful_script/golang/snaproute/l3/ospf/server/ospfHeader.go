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
	"encoding/binary"
)

type OSPFHeader struct {
	ver      uint8
	pktType  uint8
	pktlen   uint16
	routerId []byte
	areaId   []byte
	chksum   uint16
	authType uint16
	authKey  []byte
}

func NewOSPFHeader() *OSPFHeader {
	return &OSPFHeader{}
}

func encodeOspfHdr(ospfHdr OSPFHeader) []byte {
	pkt := make([]byte, OSPF_HEADER_SIZE)
	pkt[0] = ospfHdr.ver
	pkt[1] = ospfHdr.pktType
	binary.BigEndian.PutUint16(pkt[2:4], ospfHdr.pktlen)
	copy(pkt[4:8], ospfHdr.routerId)
	copy(pkt[8:12], ospfHdr.areaId)
	//binary.BigEndian.PutUint16(pkt[12:14], ospfHdr.chksum)
	binary.BigEndian.PutUint16(pkt[14:16], ospfHdr.authType)
	//copy(pkt[16:24], ospfHdr.authKey)

	return pkt
}

func decodeOspfHdr(ospfPkt []byte, ospfHdr *OSPFHeader) {
	ospfHdr.ver = uint8(ospfPkt[0])
	ospfHdr.pktType = uint8(ospfPkt[1])
	ospfHdr.pktlen = binary.BigEndian.Uint16(ospfPkt[2:4])
	ospfHdr.routerId = ospfPkt[4:8]
	ospfHdr.areaId = ospfPkt[8:12]
	ospfHdr.chksum = binary.BigEndian.Uint16(ospfPkt[12:14])
	ospfHdr.authType = binary.BigEndian.Uint16(ospfPkt[14:16])
	ospfHdr.authKey = ospfPkt[16:24]
}
