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
	"reflect"
	"testing"
)

const (
	testInvalidVrf   = "xasdfas"
	testInvalidTimer = 0
)

func testGlobalConfigNdpOperations(gblCfg NdpConfig, t *testing.T) {
	update := testNdpServer.NdpConfig.Create(gblCfg)
	if update {
		t.Error("Global Create should not be treated as update")
		return
	}

	if !reflect.DeepEqual(testNdpServer.NdpConfig, gblCfg) {
		t.Error("Failed to create ndp global config")
	}

	rv, err := testNdpServer.NdpConfig.Validate(testInvalidVrf, gblCfg.RetransTime, gblCfg.ReachableTime, gblCfg.RaRestransmitTime)
	if err == nil || rv == true {
		t.Error("Creating", testInvalidVrf, "is not supported and hence should have failed")
		return
	}
	rv, err = testNdpServer.NdpConfig.Validate(gblCfg.Vrf, 0, gblCfg.ReachableTime, gblCfg.RaRestransmitTime)
	if err == nil || rv == true {
		t.Error("Assigning 0 as re-transmit timer is not supported and hence should have failed")
		return
	}
	rv, err = testNdpServer.NdpConfig.Validate(gblCfg.Vrf, gblCfg.RetransTime, 0, gblCfg.RaRestransmitTime)
	if err == nil || rv == true {
		t.Error("Assigning 0 as reachable timer is not supported and hence should have failed")
		return
	}
	rv, err = testNdpServer.NdpConfig.Validate(gblCfg.Vrf, gblCfg.RetransTime, gblCfg.ReachableTime, 0)
	if err == nil || rv == true {
		t.Error("Assigning 0 as ra timer is not supported and hence should have failed")
		return
	}
}

func TestNdpGlobal(t *testing.T) {
	InitNDPTestServer()
	gblCfg := NdpConfig{
		Vrf:               "default",
		ReachableTime:     30000,
		RetransTime:       1,
		RaRestransmitTime: 5,
	}
	testGlobalConfigNdpOperations(gblCfg, t)
}
