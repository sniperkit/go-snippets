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
	"models/objects"
	"net"
	"testing"
	"utils/logging"
)

type DBClient struct {
	t *testing.T
}

func (d *DBClient) Init() error {
	return nil
}

func (d *DBClient) AddObject(obj objects.ConfigObj) error {
	d.t.Log("StateDBClient:AddObject")
	return nil
}

func (d *DBClient) DeleteObject(obj objects.ConfigObj) error {
	d.t.Log("StateDBClient:DeleteObject")
	return nil
}

func (d *DBClient) UpdateObject(obj objects.ConfigObj) error {
	d.t.Log("StateDBClient:UpdateObject")
	return nil
}

func (d *DBClient) DeleteAllObjects(obj objects.ConfigObj) error {
	d.t.Log("StateDBClient:DeleteAllObjects")
	return nil
}

func constructRib(t *testing.T, logger *logging.Writer, gConf *config.GlobalConfig) *LocRib {
	routeMgr := &RouteMgr{t}
	dbClient := &DBClient{t}
	locRib := NewLocRib(logger, routeMgr, dbClient, gConf)
	return locRib
}

func TestLocRib(t *testing.T) {
	logger := getLogger(t)
	gConf, _ := getConfObjects("192.168.0.100", uint32(1234), uint32(4321))
	locRib := NewLocRib(logger, nil, nil, gConf)
	if locRib != nil {
		t.Log("LocRib successfully created")
	}
}

func TestGetReachability(t *testing.T) {
	logger := getLogger(t)
	gConf, _ := getConfObjects("192.168.0.100", uint32(1234), uint32(4321))
	locRib := constructRib(t, logger, gConf)
	network := "60.1.1.0"
	reachInfo := locRib.GetReachabilityInfo(network)
	if locRib.reachabilityMap[network] != reachInfo {
		t.Fatal("LocRib:GetReachabilityInfo failed for network", network)
	}
}

func TestGetDest(t *testing.T) {
	logger := getLogger(t)
	gConf, _ := getConfObjects("192.168.0.100", uint32(1234), uint32(4321))
	locRib := constructRib(t, logger, gConf)
	network := "60.1.1.0"
	nlri := packet.NewIPPrefix(net.ParseIP(network), 24)
	protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)

	dest, exists := locRib.GetDest(nlri, protoFamily, false)
	if exists {
		t.Fatal("LocRib:GetDest failed to create Destination for network", network, "protocol family", protoFamily)
	}

	dest, exists = locRib.GetDest(nlri, protoFamily, true)
	if exists {
		t.Fatal("LocRib:GetDest failed to create Destination for network", network, "protocol family", protoFamily)
	}
	if dest == nil {
		t.Fatal("LocRib:GetDest Destination is nil for network", network, "protocol family", protoFamily)
	}

	sameDest, exists := locRib.GetDest(nlri, protoFamily, false)
	if !exists {
		t.Fatal("LocRib:GetDest failed to create Destination for network", network, "protocol family", protoFamily)
	}
	if dest != sameDest {
		t.Fatal("LocRib:GetDest did not return the same destination for", network, "protocol family", protoFamily)
	}

	anotherDest, exists := locRib.GetDest(nlri, protoFamily, true)
	if !exists {
		t.Fatal("LocRib:GetDest failed to create Destination for network", network, "protocol family", protoFamily)
	}
	if dest != anotherDest {
		t.Fatal("LocRib:GetDest did not return the same destination for", network, "protocol family", protoFamily)
	}

	network2 := "60.1.1.192"
	nlri2 := packet.NewIPPrefix(net.ParseIP(network2), 24)
	dest2, exists := locRib.GetDest(nlri2, protoFamily, true)
	if exists {
		t.Fatal("LocRib:GetDest failed to create Destination for network", network, "protocol family", protoFamily)
	}
	if dest2 == nil {
		t.Fatal("LocRib:GetDest Destination is nil for network", network, "protocol family", protoFamily)
	}
}

func constructIPPrefix(t *testing.T, ips ...string) []packet.NLRI {
	nlri := make([]packet.NLRI, 0)
	for _, cidrIP := range ips {
		ip, ipNet, err := net.ParseCIDR(cidrIP)
		if err != nil {
			t.Fatal("ParseCIDR for ip", ip, "failed with error:", err)
		}
		ones, _ := ipNet.Mask.Size()
		prefix := packet.NewIPPrefix(ip, uint8(ones))
		nlri = append(nlri, prefix)
	}
	return nlri
}

func TestProcessUpdate(t *testing.T) {
	logger := getLogger(t)
	neighbor := "192.168.0.100"
	localAS := uint32(1234)
	peerAS := uint32(4321)
	gConf, pConf := getConfObjects(neighbor, localAS, peerAS)
	nConf := base.NewNeighborConf(logger, gConf, nil, *pConf)
	locRib := constructRib(t, logger, gConf)
	pathAttrs := constructPathAttrs(net.ParseIP(neighbor), peerAS, peerAS+3, peerAS+6)
	nlri := constructIPPrefix(t, "30.1.10.0/24", "40.1.0.0/16")
	msg := packet.NewBGPUpdateMessage(nil, pathAttrs, nlri)
	bgpPktSrc := packet.BGPPktSrc{Src: neighbor, Msg: msg}

	body := bgpPktSrc.Msg.Body.(*packet.BGPUpdate)
	path := NewPath(locRib, nConf, body.PathAttributes, nil, RouteTypeEGP)
	updated := make(map[uint32]map[*Path][]*Destination)
	withdrawn := make([]*Destination, 0)
	updatedAddPaths := make([]*Destination, 0)
	addedAllPrefixes := true
	protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)

	updated, withdrawn, updatedAddPaths, addedAllPrefixes = locRib.ProcessUpdate(nConf, path, nlri, nil, protoFamily,
		0, updated, withdrawn, updatedAddPaths)
	if len(updated[protoFamily]) != 1 {
		t.Fatal("LocRib:ProcessUpdate - Found more path in protocol family", protoFamily)
	}
	for path, destinations := range updated[protoFamily] {
		if len(destinations) != 2 {
			t.Fatalf("LocRib:ProcessUpdate - Did not find 2 destinations %+v for path %+v", destinations, path)
		}
	}
	if len(withdrawn) != 0 {
		t.Fatal("LocRib:ProcessUpdate - Found withdrawn paths, withdrawn=", withdrawn)
	}
	if len(updatedAddPaths) > 0 {
		t.Fatal("LocRib:ProcessUpdate - Found add paths, updatedAddPaths=", updatedAddPaths)
	}
	if !addedAllPrefixes {
		t.Fatal("LocRib:ProcessUpdate - Did not add all prefixes to RIB")
	}

	pathAttrs = constructPathAttrs(net.ParseIP(neighbor), peerAS, peerAS+4, peerAS+8)
	removeNLRI := constructIPPrefix(t, "40.1.0.0/16")
	nlri = constructIPPrefix(t, "30.1.10.0/24", "60.1.0.0/16")
	msg = packet.NewBGPUpdateMessage(removeNLRI, pathAttrs, nlri)
	bgpPktSrc = packet.BGPPktSrc{Src: neighbor, Msg: msg}

	body = bgpPktSrc.Msg.Body.(*packet.BGPUpdate)
	path = NewPath(locRib, nConf, body.PathAttributes, nil, RouteTypeEGP)
	updated = make(map[uint32]map[*Path][]*Destination)
	withdrawn = make([]*Destination, 0)
	updatedAddPaths = make([]*Destination, 0)
	addedAllPrefixes = true

	updated, withdrawn, updatedAddPaths, addedAllPrefixes = locRib.ProcessUpdate(nConf, path, nlri, removeNLRI,
		protoFamily, 0, updated, withdrawn, updatedAddPaths)
	if len(updated[protoFamily]) != 1 {
		t.Fatal("LocRib:ProcessUpdate - Found more path in protocol family", protoFamily)
	}

	for path, destinations := range updated[protoFamily] {
		if len(destinations) != 2 {
			t.Fatalf("LocRib:ProcessUpdate - Did not find 2 destinations %+v for path %+v", destinations, path)
		}
	}
	if len(withdrawn) != 1 {
		t.Fatal("LocRib:ProcessUpdate - Found withdrawn paths, withdrawn=", withdrawn)
	}
	if len(updatedAddPaths) > 0 {
		t.Fatal("LocRib:ProcessUpdate - Found add paths, updatedAddPaths=", updatedAddPaths)
	}
	if !addedAllPrefixes {
		t.Fatal("LocRib:ProcessUpdate - Did not add all prefixes to RIB")
	}
}
