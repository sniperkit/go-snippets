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

package manager

import (
	"dhcprelayd"
)

type IPv4ProcessorIntf interface {
	StartRxTx()
	StopRxTx()
	GetEnabledFlag() bool
	SetEnabledFlag()
	ProcessCreateDRAIntf(string)
	ProcessDeleteDRAIntf(string)
	ProcessActiveDRAIntf(int)
	ProcessInactiveDRAIntf(int)
	GetClientState(string) (*dhcprelayd.DHCPRelayClientState, bool)
	GetClientStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPRelayClientState)
	GetIntfState(string) (*dhcprelayd.DHCPRelayIntfState, bool)
	GetIntfStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfState)
	GetIntfServerState(string, string) (*dhcprelayd.DHCPRelayIntfServerState, bool)
	GetIntfServerStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfServerState)
}

type IPv6ProcessorIntf interface {
	StartRxTx()
	StopRxTx()
	GetEnabledFlag() bool
	SetEnabledFlag()
	ProcessCreateDRAIntf(string)
	ProcessDeleteDRAIntf(string)
	ProcessActiveDRAIntf(int)
	ProcessInactiveDRAIntf(int)
	GetClientState(string) (*dhcprelayd.DHCPv6RelayClientState, bool)
	GetClientStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPv6RelayClientState)
	GetIntfState(string) (*dhcprelayd.DHCPv6RelayIntfState, bool)
	GetIntfStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfState)
	GetIntfServerState(string, string) (*dhcprelayd.DHCPv6RelayIntfServerState, bool)
	GetIntfServerStateSlice(int, int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfServerState)
}
