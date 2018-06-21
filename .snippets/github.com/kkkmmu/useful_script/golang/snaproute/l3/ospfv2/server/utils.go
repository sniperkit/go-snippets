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
	"errors"
	"net"
	"strconv"
)

func convertDotNotationToUint32(str string) (uint32, error) {
	var val uint32
	ip := net.ParseIP(str)
	if ip == nil {
		return 0, errors.New("Invalid string format")
	}
	ipBytes := ip.To4()
	val = val + uint32(ipBytes[0])
	val = (val << 8) + uint32(ipBytes[1])
	val = (val << 8) + uint32(ipBytes[2])
	val = (val << 8) + uint32(ipBytes[3])
	return val, nil
}

func convertMaskToUint32(mask net.IPMask) uint32 {
	var val uint32

	val = val + uint32(mask[0])
	val = (val << 8) + uint32(mask[1])
	val = (val << 8) + uint32(mask[2])
	val = (val << 8) + uint32(mask[3])
	return val
}

func ParseCIDRToUint32(IpAddr string) (ip uint32, mask uint32, err error) {
	ipAddr, ipNet, err := net.ParseCIDR(IpAddr)
	if err != nil {
		return 0, 0, errors.New("Invalid IP Address")
	}
	ip, _ = convertDotNotationToUint32(ipAddr.String())
	mask = convertMaskToUint32(ipNet.Mask)
	return ip, mask, nil
}
func convertUint32ToDotNotation(val uint32) string {
	p0 := int(val & 0xFF)
	p1 := int((val >> 8) & 0xFF)
	p2 := int((val >> 16) & 0xFF)
	p3 := int((val >> 24) & 0xFF)
	str := strconv.Itoa(p3) + "." + strconv.Itoa(p2) + "." +
		strconv.Itoa(p1) + "." + strconv.Itoa(p0)

	return str
}

func computeCheckSum(pkt []byte) uint16 {
	var csum uint32

	for i := 0; i < len(pkt); i += 2 {
		csum += uint32(pkt[i]) << 8
		csum += uint32(pkt[i+1])
	}
	chkSum := ^uint16((csum >> 16) + csum)
	return chkSum
}

const (
	MODX int = 4102
)

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func computeFletcherChecksum(data []byte, offset uint16) uint16 {
	checksum := 0
	if offset != FLETCHER_CHECKSUM_VALIDATE {
		binary.BigEndian.PutUint16(data[offset:], 0)
	}
	left := len(data)
	c0 := 0
	c1 := 0
	j := 0
	for left != 0 {
		pLen := min(left, MODX)
		for i := 0; i < pLen; i++ {
			c0 = c0 + int(data[j])
			j = j + 1
			c1 = c1 + c0
		}
		c0 = c0 % 255
		c1 = c1 % 255
		left = left - pLen
	}
	x := int((len(data)-int(offset)-1)*c0-c1) % 255
	if x <= 0 {
		x = x + 255
	}
	y := 510 - c0 - x
	if y > 255 {
		y = y - 255
	}

	if offset == FLETCHER_CHECKSUM_VALIDATE {
		checksum = (c1 << 8) + c0
	} else {
		checksum = (x << 8) | (y & 0xff)
	}

	return uint16(checksum)
}

func convertByteToOctetString(data []byte) string {
	var str string
	for i := 0; i < len(data)-1; i++ {
		str = str + strconv.Itoa(int(data[i])) + ":"
	}
	str = str + strconv.Itoa(int(data[len(data)-1]))
	return str
}
