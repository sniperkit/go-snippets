//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
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
	"dhcprelayd"
	// "models/objects"
)

// Holds message Opcode
type ServerOpId int

const ( // Request message codes
	CREATE_DHCPRELAY_GLOBAL = iota
	UPDATE_DHCPRELAY_GLOBAL
	DELETE_DHCPRELAY_GLOBAL

	CREATE_DHCPRELAY_INTF
	UPDATE_DHCPRELAY_INTF
	DELETE_DHCPRELAY_INTF

	GET_DHCPRELAY_CLIENT_STATE
	GETBLK_DHCPRELAY_CLIENT_STATE

	GET_DHCPRELAY_INTF_STATE
	GETBLK_DHCPRELAY_INTF_STATE

	GET_DHCPRELAY_INTFSERVER_STATE
	GETBLK_DHCPRELAY_INTFSERVER_STATE

	CREATE_DHCPV6RELAY_GLOBAL
	UPDATE_DHCPV6RELAY_GLOBAL
	DELETE_DHCPV6RELAY_GLOBAL

	CREATE_DHCPV6RELAY_INTF
	UPDATE_DHCPV6RELAY_INTF
	DELETE_DHCPV6RELAY_INTF

	GET_DHCPV6RELAY_CLIENT_STATE
	GETBLK_DHCPV6RELAY_CLIENT_STATE

	GET_DHCPV6RELAY_INTF_STATE
	GETBLK_DHCPV6RELAY_INTF_STATE

	GET_DHCPV6RELAY_INTFSERVER_STATE
	GETBLK_DHCPV6RELAY_INTFSERVER_STATE
)

type ServerRequest struct {
	Op   ServerOpId
	Data interface{}
}

type ServerReply struct {
	Obj   interface{}
	Error error
}

type GetBulkInArgs struct {
	FromIdx int
	Count   int
}

// DHCP Relay Global Config
type CreateDHCPRelayGlobalInArgs struct {
	DHCPRelayGlobal *dhcprelayd.DHCPRelayGlobal
}

type UpdateDHCPRelayGlobalInArgs struct {
	DHCPRelayGlobalOld *dhcprelayd.DHCPRelayGlobal
	DHCPRelayGlobalNew *dhcprelayd.DHCPRelayGlobal
	AttrSet            []bool
}

type DeleteDHCPRelayGlobalInArgs struct {
	Vrf string
}

// DHCP Interface Config
type CreateDHCPRelayIntfInArgs struct {
	DHCPRelayIntf *dhcprelayd.DHCPRelayIntf
}

type UpdateDHCPRelayIntfInArgs struct {
	DHCPRelayIntfOld *dhcprelayd.DHCPRelayIntf
	DHCPRelayIntfNew *dhcprelayd.DHCPRelayIntf
	AttrSet          []bool
}

type DeleteDHCPRelayIntfInArgs struct {
	IntfRef string
}

// DHCP Client State
type GetDHCPRelayClientStateInArgs struct {
	MacAddr string
}

type GetDHCPRelayClientStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayClientState
	Err error
}

type GetBulkDHCPRelayClientStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayClientStateGetInfo
	Err error
}

// Interface State
type GetDHCPRelayIntfStateInArgs struct {
	IntfRef string
}

type GetDHCPRelayIntfStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayIntfState
	Err error
}

type GetBulkDHCPRelayIntfStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayIntfStateGetInfo
	Err error
}

// Interface Server State
type GetDHCPRelayIntfServerStateInArgs struct {
	IntfRef  string
	ServerIp string
}

type GetDHCPRelayIntfServerStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayIntfServerState
	Err error
}

type GetBulkDHCPRelayIntfServerStateOutArgs struct {
	Obj *dhcprelayd.DHCPRelayIntfServerStateGetInfo
	Err error
}

// DHCPv6 Relay Global Config
type CreateDHCPv6RelayGlobalInArgs struct {
	DHCPv6RelayGlobal *dhcprelayd.DHCPv6RelayGlobal
}

type UpdateDHCPv6RelayGlobalInArgs struct {
	DHCPv6RelayGlobalOld *dhcprelayd.DHCPv6RelayGlobal
	DHCPv6RelayGlobalNew *dhcprelayd.DHCPv6RelayGlobal
	AttrSet              []bool
}

type DeleteDHCPv6RelayGlobalInArgs struct {
	Vrf string
}

// DHCPv6 Interface Config
type CreateDHCPv6RelayIntfInArgs struct {
	DHCPv6RelayIntf *dhcprelayd.DHCPv6RelayIntf
}

type UpdateDHCPv6RelayIntfInArgs struct {
	DHCPv6RelayIntfOld *dhcprelayd.DHCPv6RelayIntf
	DHCPv6RelayIntfNew *dhcprelayd.DHCPv6RelayIntf
	AttrSet            []bool
}

type DeleteDHCPv6RelayIntfInArgs struct {
	IntfRef string
}

// DHCPv6 Client State
type GetDHCPv6RelayClientStateInArgs struct {
	MacAddr string
}

type GetDHCPv6RelayClientStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayClientState
	Err error
}

type GetBulkDHCPv6RelayClientStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayClientStateGetInfo
	Err error
}

// Interface State
type GetDHCPv6RelayIntfStateInArgs struct {
	IntfRef string
}

type GetDHCPv6RelayIntfStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayIntfState
	Err error
}

type GetBulkDHCPv6RelayIntfStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayIntfStateGetInfo
	Err error
}

// Interface Server State
type GetDHCPv6RelayIntfServerStateInArgs struct {
	IntfRef  string
	ServerIp string
}

type GetDHCPv6RelayIntfServerStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayIntfServerState
	Err error
}

type GetBulkDHCPv6RelayIntfServerStateOutArgs struct {
	Obj *dhcprelayd.DHCPv6RelayIntfServerStateGetInfo
	Err error
}
