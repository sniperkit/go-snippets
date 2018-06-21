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

// fsmState.go
package config

const BGPConnectRetryTime uint32 = 120 // seconds
const BGPHoldTimeDefault uint32 = 180  // 180 seconds

type BGPFSMState int

const (
	BGPFSMNone BGPFSMState = iota
	BGPFSMIdle
	BGPFSMConnect
	BGPFSMActive
	BGPFSMOpensent
	BGPFSMOpenconfirm
	BGPFSMEstablished
)

var BGPStateToStr = map[BGPFSMState]string{
	BGPFSMNone:        "INIT",
	BGPFSMIdle:        "IDLE",
	BGPFSMConnect:     "CONNECT",
	BGPFSMActive:      "ACTIVE",
	BGPFSMOpensent:    "OPENSENT",
	BGPFSMOpenconfirm: "OPENCONFIRM",
	BGPFSMEstablished: "ESTABLISHED",
}

func GetBGPStateToStr(stateId BGPFSMState) string {
	return BGPStateToStr[stateId]
}
