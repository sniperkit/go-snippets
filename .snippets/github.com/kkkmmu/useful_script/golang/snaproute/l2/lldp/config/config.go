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

package config

const (
	TX_RX_MODE_TxRx   = "TxRx"
	TXRX              = 1
	TX_RX_MODE_TxOnly = "TxOnly"
	TX_ONLY           = 2
	TX_RX_MODE_RxOnly = "RxOnly"
	RX_ONLY           = 3
	INVALID_MODE      = 255
)

const (
	INVALID_PVID = -1
)

type Global struct {
	Vrf             string
	Enable          bool
	TranmitInterval int32
	TxRxMode        uint8
	SnoopAndDrop    bool
}

// this is used for auto-discovery
type Intf struct {
	IntfRef  string
	Enable   bool
	TxRxMode string
}

// this is used to update configuration request coming from client to server
type IntfConfig struct {
	IfIndex  int32
	Enable   bool
	TxRxMode uint8
}

type PortInfo struct {
	IfIndex     int32
	Name        string
	OperState   string
	MacAddr     string
	Description string
	Pvid        int32
}

type PortState struct {
	IfIndex int32
	IfState string
}

type IntfState struct {
	IntfRef             string
	IfIndex             int32
	Enable              bool
	SendFrames          int32
	ReceivedFrames      int32
	LocalPort           string
	PeerMac             string
	PeerPort            string
	PeerHostName        string
	PortDescription     string
	VlanId              int32
	HoldTime            string
	SystemDescription   string
	SystemCapabilities  string
	EnabledCapabilities string
}

type GlobalState struct {
	Vrf             string
	Enable          bool
	TranmitInterval int32
	Neighbors       int32
	TotalTxFrames   int32
	TotalRxFrames   int32
}

type EventInfo struct {
	IfIndex   int32
	EventType int
}

const (
	_ = iota
	Learned
	Updated
	Removed
	NoOp
)

type SystemInfo struct {
	Vrf         string
	MgmtIp      string
	Hostname    string
	SwitchMac   string
	SwVersion   string
	Description string
}

type PortAttrMsg struct {
	IfIndex     int32
	Description string
	Pvid        int32
}
