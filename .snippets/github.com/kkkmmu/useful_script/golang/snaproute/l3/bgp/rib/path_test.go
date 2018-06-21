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

// path_test.go
package rib

import (
	"l3/bgp/baseobjects"
	"l3/bgp/config"
	"l3/bgp/packet"
	"l3/bgp/utils"
	"net"
	"testing"
	"utils/logging"
)

func getLogger(t *testing.T) *logging.Writer {
	logger, err := logging.NewLogger("bgpd", "BGP", true)
	if err != nil {
		t.Fatal("Failed to start the logger. Exiting!!")
		return nil
	}
	utils.SetLogger(logger)
	return logger
}

func getConfObjects(neighbor string, localAS, peerAS uint32) (*config.GlobalConfig, *config.NeighborConfig) {
	gConf := &config.GlobalConfig{}
	gConf.AS = localAS
	gConf.RouterId = net.ParseIP("10.1.10.100")
	pConf := &config.NeighborConfig{}
	pConf.NeighborAddress = net.ParseIP(neighbor)
	pConf.PeerAS = peerAS
	return gConf, pConf
}

func constructPathAttrs(nh net.IP, asList ...uint32) []packet.BGPPathAttr {
	pathAttrs := make([]packet.BGPPathAttr, 0)

	origin := packet.NewBGPPathAttrOrigin(packet.BGPPathAttrOriginIncomplete)
	pathAttrs = append(pathAttrs, origin)

	asSeg := packet.NewBGPAS4PathSegmentSeq()
	for _, as := range asList {
		asSeg.AppendAS(as)
	}
	asPath := packet.NewBGPPathAttrASPath()
	asPath.ASSize = 4
	asPath.AppendASPathSegment(asSeg)
	pathAttrs = append(pathAttrs, asPath)

	nextHop := packet.NewBGPPathAttrNextHop()
	nextHop.Value = nh
	pathAttrs = append(pathAttrs, nextHop)

	return pathAttrs
}

func TestPath(t *testing.T) {
	logger := getLogger(t)
	gConf, pConf := getConfObjects("192.168.0.100", uint32(1234), uint32(4321))
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	pathAttrs := constructPathAttrs(pConf.NeighborAddress, pConf.PeerAS, pConf.PeerAS+1)
	locRib := NewLocRib(logger, nil, nil, gConf)
	path := NewPath(locRib, nConf, pathAttrs, nil, RouteTypeEGP)
	if path != nil {
		t.Log("Path successfully created")
	}
}
