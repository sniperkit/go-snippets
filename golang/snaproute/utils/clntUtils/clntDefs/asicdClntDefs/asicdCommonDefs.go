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

package asicdClntDefs

import ()

const (
	IfTypePort = iota
	IfTypeLag
	IfTypeVlan
	IfTypeP2P
	IfTypeBcast
	IfTypeLoopback
	IfTypeSecondary
	IfTypeVirtual
	IfTypeVtep
	IfTypeNull
)

const (
	// state values copied from asicd
	INTF_STATE_DOWN = 0
	INTF_STATE_UP   = 1
)

const (
	// this needs to match asicd server
	PORT_ATTR_PHY_INTF_TYPE = 0x00000001
	PORT_ATTR_ADMIN_STATE   = 0x00000002
	PORT_ATTR_MAC_ADDR      = 0x00000004
	PORT_ATTR_SPEED         = 0x00000008
	PORT_ATTR_DUPLEX        = 0x00000010
	PORT_ATTR_AUTONEG       = 0x00000020
	PORT_ATTR_MEDIA_TYPE    = 0x00000040
	PORT_ATTR_MTU           = 0x00000080
	PORT_ATTR_BREAKOUT_MODE = 0x00000100
	PORT_ATTR_LOOPBACK_MODE = 0x00000200
	PORT_ATTR_ENABLE_FEC    = 0x00000400
	PORT_ATTR_TX_PRBS_EN    = 0x00000800
	PORT_ATTR_RX_PRBS_EN    = 0x00001000
	PORT_ATTR_PRBS_POLY     = 0x00002000
	PORT_ATTR_DESCRIPTION   = 0x00004000
	PORT_ATTR_PVID          = 0x00008000
)

const (
	PORT_UPDATE_ATTR_INTF_REF        = 0x1
	PORT_UPDATE_ATTR_IF_INDEX        = 0x2
	PORT_UPDATE_ATTR_DESCRIPTION     = 0x4
	PORT_UPDATE_ATTR_PHY_INTF_TYPE   = 0x8
	PORT_UPDATE_ATTR_ADMIN_STATE     = 0x10
	PORT_UPDATE_ATTR_MAC_ADDR        = 0x20
	PORT_UPDATE_ATTR_SPEED           = 0x40
	PORT_UPDATE_ATTR_DUPLEX          = 0x80
	PORT_UPDATE_ATTR_AUTONEG         = 0x100
	PORT_UPDATE_ATTR_MEDIA_TYPE      = 0x200
	PORT_UPDATE_ATTR_MTU             = 0x400
	PORT_UPDATE_ATTR_BREAK_OUT_MODE  = 0x800
	PORT_UPDATE_ATTR_LOOPBACK_MODE   = 0x1000
	PORT_UPDATE_ATTR_ENABLE_FEC      = 0x2000
	PORT_UPDATE_ATTR_PRBS_TX_ENABLE  = 0x4000
	PORT_UPDATE_ATTR_PRBS_RX_ENABLE  = 0x8000
	PORT_UPDATE_ATTR_PRBS_POLYNOMIAL = 0x10000
)
