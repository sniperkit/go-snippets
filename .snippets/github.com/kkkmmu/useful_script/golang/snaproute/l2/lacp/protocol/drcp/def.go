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

// def.go
package drcp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	CONVERSATION_ID_TYPE_PORT    = 0
	CONVERSATION_ID_TYPE_GATEWAY = 1

	MAX_IPP_LINKS         = 3
	MAX_CONVERSATION_IDS  = 4096
	MAX_PORTAL_SYSTEM_IDS = 3

	DrniFastPeriodicTime time.Duration = time.Second * 1
	DrniSlowPeriodictime time.Duration = time.Second * 30
	DrniShortTimeoutTime time.Duration = 3 * DrniFastPeriodicTime
	DrniLongTimeoutTime  time.Duration = 3 * DrniSlowPeriodictime
)

var GATEWAY_ALGORITHM_NULL GatewayAlgorithm = GatewayAlgorithm{}
var GATEWAY_ALGORITHM_RESERVED GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x00}
var GATEWAY_ALGORITHM_CVID GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x01}
var GATEWAY_ALGORITHM_SVID GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x02}
var GATEWAY_ALGORITHM_ISID GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x03}
var GATEWAY_ALGORITHM_TE_SID GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x04}
var GATEWAY_ALGORITHM_ECMP_FLOW_HASH GatewayAlgorithm = GatewayAlgorithm{0x00, 0x80, 0xC2, 0x05}

var ENCAP_METHOD_SHARING_NULL [4]uint8 = [4]uint8{0x00, 0x00, 0x00, 0x00}
var ENCAP_METHOD_SEPARATE_LINKS [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x00}
var ENCAP_METHOD_SHARING_BY_TIME [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x01}
var ENCAP_METHOD_SHARING_BY_TAG [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x02}
var ENCAP_METHOD_SHARING_BY_ITAG [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x03}
var ENCAP_METHOD_SHARING_BY_BTAG [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x04}
var ENCAP_METHOD_SHARING_BY_PSEUDOWIRE [4]uint8 = [4]uint8{0x00, 0x80, 0xC2, 0x05}

type GatewayAlgorithm [4]uint8
type EncapMethod [4]uint8
type Md5Digest [16]uint8

func (g *GatewayAlgorithm) String() string {
	return fmt.Sprintf("%02x-%02x-%02x-%02x", g[0], g[1], g[2], g[3])
}

func (g *EncapMethod) String() string {
	return fmt.Sprintf("%02x-%02x-%02x-%02x", g[0], g[1], g[2], g[3])
}

func (d Md5Digest) get16Bytes() [16]uint8 {
	return [16]uint8{
		d[0], d[1], d[2], d[3],
		d[4], d[5], d[6], d[7],
		d[8], d[9], d[10], d[11],
		d[12], d[13], d[14], d[15],
	}
}

func (d Md5Digest) calculatePortDigest(portList [][]uint16) Md5Digest {
	hash := md5.New()
	if portList != nil {
		for i, ports := range portList {
			buf := new(bytes.Buffer)
			data := ports
			data = append(data, uint16(i))
			// network byte order
			binary.Write(buf, binary.BigEndian, data)
			hash.Write(buf.Bytes())
		}
	}

	digest := hash.Sum(nil)
	for i, _ := range digest {
		d[i] = digest[i]
	}
	return d
}

func (d Md5Digest) calculateGatewayDigest(gatewayList [][]uint8) Md5Digest {
	hash := md5.New()
	if gatewayList != nil {
		for _, sysnums := range gatewayList {
			buf := new(bytes.Buffer)
			// network byte order
			binary.Write(buf, binary.BigEndian, sysnums)
			hash.Write(buf.Bytes())
		}
	}

	digest := hash.Sum(nil)
	for i, _ := range digest {
		d[i] = digest[i]
	}
	return d
}
