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
	"fmt"
	"net"
	"strconv"
)

const big = 0xFFFFFF

func dtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}

func parseIPFmt(s string) []byte {
	var p [4]byte
	i := 0
	for j := 0; j < 4; j++ {
		if i >= len(s) {
			// Missing octets.
			return nil
		}
		if j > 0 {
			if s[i] != '.' {
				return nil
			}
			i++
		}
		var (
			n  int
			ok bool
		)
		n, i, ok = dtoi(s, i)
		if !ok || n > 0xFF {
			return nil
		}
		p[j] = byte(n)
	}
	if i != len(s) {
		return nil
	}
	return []byte{p[0], p[1], p[2], p[3]}
}

func parseIntFmt(str string) []byte {
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return nil
	}
	var p [4]byte
	fmt.Println(i)
	p[0] = byte(i & 0xFF)
	p[1] = byte((i & 0xFF00) >> 8)
	p[2] = byte((i & 0xFF0000) >> 16)
	p[3] = byte((i & 0xFF000000) >> 24)
	return []byte{p[0], p[1], p[2], p[3]}
}

func convertAreaOrRouterId(str string) []byte {
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			return parseIPFmt(str)
		}
	}
	return parseIntFmt(str)
}

func convertAreaOrRouterIdUint32(str string) uint32 {
	return convertIPv4ToUint32(convertAreaOrRouterId(str))
}

func convertAuthKey(s string) []byte {
	var p [8]byte
	i := 0
	for j := 0; j < 8; j++ {
		if i >= len(s) {
			// Missing octets.
			return nil
		}
		if j > 0 {
			if s[i] != '.' {
				return nil
			}
			i++
		}
		var (
			n  int
			ok bool
		)
		n, i, ok = dtoi(s, i)
		if !ok || n > 0xFF {
			return nil
		}
		p[j] = byte(n)
	}
	if i != len(s) {
		return nil
	}
	return []byte{p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7]}
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

func bytesEqual(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i, b := range x {
		if y[i] != b {
			return false
		}
	}
	return true
}

func isInSubnet(ifIpAddr net.IP, srcIp net.IP, netMask net.IPMask) bool {
	net1 := ifIpAddr.Mask(netMask)
	net2 := srcIp.Mask(netMask)
	if net1.Equal(net2) {
		return true
	}
	return false
}

func convertIPv4ToUint32(ip []byte) uint32 {
	var val uint32 = 0

	val = val + uint32(ip[0])
	val = (val << 8) + uint32(ip[1])
	val = (val << 8) + uint32(ip[2])
	val = (val << 8) + uint32(ip[3])

	return val
}

func convertUint32ToIPv4(val uint32) string {
	p0 := int(val & 0xFF)
	p1 := int((val >> 8) & 0xFF)
	p2 := int((val >> 16) & 0xFF)
	p3 := int((val >> 24) & 0xFF)
	str := strconv.Itoa(p3) + "." + strconv.Itoa(p2) + "." +
		strconv.Itoa(p1) + "." + strconv.Itoa(p0)

	return str
}

func convertIPInByteToString(ip []byte) string {
	return convertUint32ToIPv4(convertIPv4ToUint32(ip))
}

const (
	MODX                       int    = 4102
	FLETCHER_CHECKSUM_VALIDATE uint16 = 0xffff
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
