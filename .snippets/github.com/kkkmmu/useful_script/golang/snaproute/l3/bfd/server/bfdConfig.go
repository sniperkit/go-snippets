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
	"l3/bfd/bfddCommonDefs"
	"net"
	"time"
)

type GlobalConfig struct {
	Vrf    string
	Enable bool
}

type GlobalState struct {
	Vrf                  string
	Enable               bool
	NumSessions          uint32
	NumUpSessions        uint32
	NumDownSessions      uint32
	NumAdminDownSessions uint32
}

type SessionConfig struct {
	DestIp    string
	ParamName string
	Interface string
	PerLink   bool
	Protocol  bfddCommonDefs.BfdSessionOwner
	Operation bfddCommonDefs.BfdSessionOperation
}

type SessionState struct {
	IpAddr                    string
	LocalAddr                 string
	SessionId                 int32
	ParamName                 string
	Interface                 string
	InterfaceSpecific         bool
	PerLinkSession            bool
	LocalMacAddr              net.HardwareAddr
	RemoteMacAddr             net.HardwareAddr
	RegisteredProtocols       []bool
	SessionState              BfdSessionState
	RemoteSessionState        BfdSessionState
	LocalDiscriminator        uint32
	RemoteDiscriminator       uint32
	LocalDiagType             BfdDiagnostic
	DesiredMinTxInterval      int32
	RequiredMinRxInterval     int32
	RemoteMinRxInterval       int32
	DetectionMultiplier       int32
	RemoteDetectionMultiplier int32
	DemandMode                bool
	RemoteDemandMode          bool
	AuthType                  AuthenticationType
	AuthSeqKnown              bool
	ReceivedAuthSeq           uint32
	SentAuthSeq               uint32
	NumTxPackets              uint32
	NumRxPackets              uint32
	ToDownCount               uint32
	ToUpCount                 uint32
	UpTime                    time.Time
}

type SessionParamConfig struct {
	Name                      string
	LocalMultiplier           int32
	DesiredMinTxInterval      int32
	RequiredMinRxInterval     int32
	RequiredMinEchoRxInterval int32
	DemandEnabled             bool
	AuthenticationEnabled     bool
	AuthenticationType        AuthenticationType
	AuthenticationKeyId       int32
	AuthenticationData        string
}

type SessionParamState struct {
	Name                      string
	NumSessions               int32
	LocalMultiplier           int32
	DesiredMinTxInterval      int32
	RequiredMinRxInterval     int32
	RequiredMinEchoRxInterval int32
	DemandEnabled             bool
	AuthenticationEnabled     bool
	AuthenticationType        AuthenticationType
	AuthenticationKeyId       int32
	AuthenticationData        string
}
