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
	"errors"
	"fmt"
	"l3/ndp/debug"
)

func (cfg *NdpConfig) Validate(vrf string, retransmit uint32, reachableTime uint32, raTime uint8) (bool, error) {
	if retransmit == 0 {
		return false, errors.New(fmt.Sprintln("Invalid Re-Transmit Time", retransmit))
	}
	if reachableTime == 0 {
		return false, errors.New(fmt.Sprintln("Invalid ReachableTime", reachableTime))
	}

	if raTime == 0 {
		return false, errors.New(fmt.Sprintln("Invalid Router Advertisement Interval", raTime))
	}

	if vrf != "default" {
		return false, errors.New(fmt.Sprintln("Global Config is only supported for default VRF", vrf))
	}

	return true, nil
}

func (cfg *NdpConfig) Create(gCfg NdpConfig) bool {
	update := false
	if cfg.Vrf == "" {
		debug.Logger.Debug("Received Global Config Create for NDP:", gCfg)
	} else {
		debug.Logger.Debug("Received Global Config Update for NDP:", gCfg)
		update = true
	}
	cfg.Vrf = gCfg.Vrf
	cfg.RetransTime = gCfg.RetransTime
	cfg.RaRestransmitTime = gCfg.RaRestransmitTime
	cfg.ReachableTime = gCfg.ReachableTime
	return update
}
