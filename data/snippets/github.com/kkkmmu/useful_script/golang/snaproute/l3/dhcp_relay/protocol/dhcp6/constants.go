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

package dhcp6

// DHCP Packet global constants
const (
	DHCP_SERVER_PORT     = 547
	DHCP_CLIENT_PORT     = 546
	DHCP_BROADCAST_IP    = "255.255.255.255"
	DHCP_NO_IP           = "0.0.0.0"
	DHCP_REDIS_DB_PORT   = ":6379"
	HOP_COUNT_LIMIT      = 32
	ALL_DRA_SERVERS_ADDR = "FF02::1:2"
	ALL_SERVERS_ADDR     = "FF05::1:3"
)

// DHCP Message Types
const (
	SOLICIT     MsgType = 1
	ADVERTISE   MsgType = 2
	REQUEST     MsgType = 3
	CONFIRM     MsgType = 4
	RENEW       MsgType = 5
	REBIND      MsgType = 6
	REPLY       MsgType = 7
	RELEASE     MsgType = 8
	DECLINE     MsgType = 9
	RECONFIGURE MsgType = 10
	INFO_REQ    MsgType = 11
	RELAY_FORW  MsgType = 12
	RELAY_REPL  MsgType = 13
)

// DHCP Option Types
const (
	OPTION_CLIENTID  OptionType = 1
	OPTION_IA_NA     OptionType = 3
	OPTION_IA_TA     OptionType = 4
	OPTION_IAADDR    OptionType = 5
	OPTION_RELAY_MSG OptionType = 9
)
