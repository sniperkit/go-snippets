//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
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
	"net"
)

type EthHdrMetadata struct {
	SrcMAC net.HardwareAddr
}

func NewEthHdrMetadata() *EthHdrMetadata {
	return &EthHdrMetadata{}
}

type IpHdrMetadata struct {
	SrcIP     uint32
	DstIP     uint32
	DstIPType DstIPType
}

func NewIpHdrMetadata() *IpHdrMetadata {
	return &IpHdrMetadata{}
}

type OSPFHeader struct {
	Ver      uint8
	PktType  uint8
	Pktlen   uint16
	RouterId uint32
	AreaId   uint32
	Chksum   uint16
	AuthType uint16
	AuthKey  []byte
}

func NewOSPFHeader() *OSPFHeader {
	return &OSPFHeader{}
}

func encodeOspfHdr(ospfHdr OSPFHeader) []byte {
	pkt := make([]byte, OSPF_HEADER_SIZE)
	pkt[0] = ospfHdr.Ver
	pkt[1] = ospfHdr.PktType
	binary.BigEndian.PutUint16(pkt[2:4], ospfHdr.Pktlen)
	binary.BigEndian.PutUint32(pkt[4:8], ospfHdr.RouterId)
	binary.BigEndian.PutUint32(pkt[8:12], ospfHdr.AreaId)
	//binary.BigEndian.PutUint16(pkt[12:14], ospfHdr.chksum)
	binary.BigEndian.PutUint16(pkt[14:16], ospfHdr.AuthType)
	//copy(pkt[16:24], ospfHdr.authKey)

	return pkt
}

func decodeOspfHdr(ospfPkt []byte, ospfHdr *OSPFHeader) {
	ospfHdr.Ver = uint8(ospfPkt[0])
	ospfHdr.PktType = uint8(ospfPkt[1])
	ospfHdr.Pktlen = binary.BigEndian.Uint16(ospfPkt[2:4])
	ospfHdr.RouterId = binary.BigEndian.Uint32(ospfPkt[4:8])
	ospfHdr.AreaId = binary.BigEndian.Uint32(ospfPkt[8:12])
	ospfHdr.Chksum = binary.BigEndian.Uint16(ospfPkt[12:14])
	ospfHdr.AuthType = binary.BigEndian.Uint16(ospfPkt[14:16])
	ospfHdr.AuthKey = ospfPkt[16:24]
}

type OspfHdrMetadata struct {
	PktType  uint8
	Pktlen   uint16
	Backbone bool
	RouterId uint32
	AreaId   uint32
}

func NewOspfHdrMetadata() *OspfHdrMetadata {
	return &OspfHdrMetadata{}
}

type OspfPktStruct struct {
	Data      []byte
	EthHdrMd  *EthHdrMetadata
	IpHdrMd   *IpHdrMetadata
	OspfHdrMd *OspfHdrMetadata
}

func NewOspfPktStruct() *OspfPktStruct {
	return &OspfPktStruct{}
}

type OSPFHelloData struct {
	Netmask         uint32
	HelloInterval   uint16
	Options         uint8
	RtrPrio         uint8
	RtrDeadInterval uint32
	DRtrIpAddr      uint32
	BDRtrIpAddr     uint32
	NbrList         []uint32
}

func NewOSPFHelloData() *OSPFHelloData {
	return &OSPFHelloData{}
}

func encodeOspfHelloData(helloData OSPFHelloData, nbrList []uint32) []byte {
	pkt := make([]byte, OSPF_HELLO_MIN_SIZE+len(nbrList)*4)
	binary.BigEndian.PutUint32(pkt[0:4], helloData.Netmask)
	binary.BigEndian.PutUint16(pkt[4:6], helloData.HelloInterval)
	pkt[6] = helloData.Options
	pkt[7] = helloData.RtrPrio
	binary.BigEndian.PutUint32(pkt[8:12], helloData.RtrDeadInterval)
	binary.BigEndian.PutUint32(pkt[12:16], helloData.DRtrIpAddr)
	binary.BigEndian.PutUint32(pkt[16:20], helloData.BDRtrIpAddr)
	start := OSPF_HELLO_MIN_SIZE
	end := start + 4
	for _, nbr := range nbrList {
		binary.BigEndian.PutUint32(pkt[start:end], nbr)
		start = end
		end = start + 4
	}
	return pkt
}

func decodeOspfHelloData(data []byte, ospfHelloData *OSPFHelloData) {
	ospfHelloData.Netmask = binary.BigEndian.Uint32(data[0:4])
	ospfHelloData.HelloInterval = binary.BigEndian.Uint16(data[4:6])
	ospfHelloData.Options = data[6]
	ospfHelloData.RtrPrio = data[7]
	ospfHelloData.RtrDeadInterval = binary.BigEndian.Uint32(data[8:12])
	ospfHelloData.DRtrIpAddr = binary.BigEndian.Uint32(data[12:16])
	ospfHelloData.BDRtrIpAddr = binary.BigEndian.Uint32(data[16:20])
	numNbrs := (len(data) - OSPF_HELLO_MIN_SIZE) / 4
	start := OSPF_HELLO_MIN_SIZE
	end := start + 4
	for idx := 0; idx < numNbrs; idx++ {
		nbr := binary.BigEndian.Uint32(data[start:end])
		ospfHelloData.NbrList = append(ospfHelloData.NbrList, nbr)
		start = end
		end = start + 4
	}
}
