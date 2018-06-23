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

package asicdCommonDefs

import (
	"asicd/pluginManager/pluginCommon"
	"asicdInt"
)

const (
	PUB_SOCKET_ADDR                      = "ipc:///tmp/asicd_all.ipc"
	PUB_SOCKET_RIBD_CLIENT_ADDR          = "ipc:///tmp/asicd_ribd.ipc"
	SYS_RSVD_VLAN                        = pluginCommon.SYS_RSVD_VLAN
	MIN_SYS_PORTS                        = pluginCommon.MIN_SYS_PORTS
	MAX_SYS_PORTS                        = pluginCommon.MAX_SYS_PORTS
	HASH_SEL_SRCDSTMAC                   = pluginCommon.HASHTYPE_SRCMAC_DSTMAC
	HASH_SEL_SRCDSTIP                    = pluginCommon.HASHTYPE_SRCIP_DSTIP
	NOTIFY_IPV4_ROUTE_CREATE_FAILURE     = pluginCommon.NOTIFY_IPV4_ROUTE_CREATE_FAILURE
	NOTIFY_IPV4_ROUTE_DELETE_FAILURE     = pluginCommon.NOTIFY_IPV4_ROUTE_DELETE_FAILURE
	NOTIFY_IPV6_ROUTE_CREATE_FAILURE     = pluginCommon.NOTIFY_IPV6_ROUTE_CREATE_FAILURE
	NOTIFY_IPV6_ROUTE_DELETE_FAILURE     = pluginCommon.NOTIFY_IPV6_ROUTE_DELETE_FAILURE
	NOTIFY_L2INTF_STATE_CHANGE           = pluginCommon.NOTIFY_L2INTF_STATE_CHANGE
	NOTIFY_IPV4_L3INTF_STATE_CHANGE      = pluginCommon.NOTIFY_IPV4_L3INTF_STATE_CHANGE
	NOTIFY_IPV6_L3INTF_STATE_CHANGE      = pluginCommon.NOTIFY_IPV6_L3INTF_STATE_CHANGE
	NOTIFY_VLAN_CREATE                   = pluginCommon.NOTIFY_VLAN_CREATE
	NOTIFY_VLAN_DELETE                   = pluginCommon.NOTIFY_VLAN_DELETE
	NOTIFY_VLAN_UPDATE                   = pluginCommon.NOTIFY_VLAN_UPDATE
	NOTIFY_IPV4INTF_CREATE               = pluginCommon.NOTIFY_IPV4INTF_CREATE
	NOTIFY_IPV4INTF_DELETE               = pluginCommon.NOTIFY_IPV4INTF_DELETE
	NOTIFY_IPV6INTF_CREATE               = pluginCommon.NOTIFY_IPV6INTF_CREATE
	NOTIFY_IPV6INTF_DELETE               = pluginCommon.NOTIFY_IPV6INTF_DELETE
	NOTIFY_LAG_CREATE                    = pluginCommon.NOTIFY_LAG_CREATE
	NOTIFY_LAG_DELETE                    = pluginCommon.NOTIFY_LAG_DELETE
	NOTIFY_LAG_UPDATE                    = pluginCommon.NOTIFY_LAG_UPDATE
	NOTIFY_LOGICAL_INTF_CREATE           = pluginCommon.NOTIFY_LOGICAL_INTF_CREATE
	NOTIFY_LOGICAL_INTF_DELETE           = pluginCommon.NOTIFY_LOGICAL_INTF_DELETE
	NOTIFY_LOGICAL_INTF_UPDATE           = pluginCommon.NOTIFY_LOGICAL_INTF_UPDATE
	NOTIFY_IPV4NBR_MAC_MOVE              = pluginCommon.NOTIFY_IPV4NBR_MAC_MOVE
	NOTIFY_IPV6NBR_MAC_MOVE              = pluginCommon.NOTIFY_IPV6NBR_MAC_MOVE
	NOTIFY_VTEP_CREATE                   = pluginCommon.NOTIFY_VTEP_CREATE
	NOTIFY_VTEP_DELETE                   = pluginCommon.NOTIFY_VTEP_DELETE
	NOTIFY_PORT_CONFIG_MODE_CHANGE       = pluginCommon.NOTIFY_PORT_CONFIG_MODE_CHANGE
	NOTIFY_PORT_ATTR_CHANGE              = pluginCommon.NOTIFY_PORT_ATTR_CHANGE
	INTF_STATE_UP                        = pluginCommon.INTF_STATE_UP
	INTF_STATE_DOWN                      = pluginCommon.INTF_STATE_DOWN
	INTF_TYPE_MASK                       = pluginCommon.INTF_TYPE_MASK
	INTF_TYPE_SHIFT                      = pluginCommon.INTF_TYPE_SHIFT
	INTF_ID_MASK                         = pluginCommon.INTF_ID_MASK
	INTF_ID_SHIFT                        = pluginCommon.INTF_ID_SHIFT
	STP_PORT_STATE_BLOCKING              = pluginCommon.STP_PORT_STATE_BLOCKING
	STP_PORT_STATE_LEARNING              = pluginCommon.STP_PORT_STATE_LEARNING
	STP_PORT_STATE_FORWARDING            = pluginCommon.STP_PORT_STATE_FORWARDING
	IP_TYPE_IPV6                         = pluginCommon.IP_TYPE_IPV6
	IP_TYPE_IPV4                         = pluginCommon.IP_TYPE_IPV4
	PORT_MODE_L2                         = pluginCommon.PORT_MODE_L2
	PORT_MODE_L3                         = pluginCommon.PORT_MODE_L3
	PORT_MODE_UNCONFIGURED               = pluginCommon.PORT_MODE_UNCONFIGURED
	PORT_MODE_INTERNAL                   = pluginCommon.PORT_MODE_INTERNAL
	NOTIFY_IPV4VIRTUAL_INTF_CREATE       = pluginCommon.NOTIFY_IPV4VIRTUAL_INTF_CREATE
	NOTIFY_IPV4VIRTUAL_INTF_DELETE       = pluginCommon.NOTIFY_IPV4VIRTUAL_INTF_DELETE
	NOTIFY_IPV6VIRTUAL_INTF_CREATE       = pluginCommon.NOTIFY_IPV6VIRTUAL_INTF_CREATE
	NOTIFY_IPV6VIRTUAL_INTF_DELETE       = pluginCommon.NOTIFY_IPV6VIRTUAL_INTF_DELETE
	NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE = pluginCommon.NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE
	NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE = pluginCommon.NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE
)

var GetIntfIdFromIfIndex pluginCommon.GetId = pluginCommon.GetIdFromIfIndex
var GetIntfTypeFromIfIndex pluginCommon.GetType = pluginCommon.GetTypeFromIfIndex
var GetIfIndexFromIntfIdAndIntfType pluginCommon.GetIfIndex = pluginCommon.GetIfIndexFromIdType

type AsicdNotification pluginCommon.AsicdNotification
type L2IntfStateNotifyMsg pluginCommon.L2IntfStateNotifyMsg
type IPv4L3IntfStateNotifyMsg pluginCommon.IPv4L3IntfStateNotifyMsg
type IPv6L3IntfStateNotifyMsg pluginCommon.IPv6L3IntfStateNotifyMsg
type VlanNotifyMsg pluginCommon.VlanNotifyMsg
type LogicalIntfNotifyMsg pluginCommon.LogicalIntfNotifyMsg
type LagNotifyMsg pluginCommon.LagNotifyMsg
type IPv4IntfNotifyMsg pluginCommon.IPv4IntfNotifyMsg
type IPv6IntfNotifyMsg pluginCommon.IPv6IntfNotifyMsg
type IPv4NbrMacMoveNotifyMsg pluginCommon.IPv4NbrMacMoveNotifyMsg
type IPv6NbrMacMoveNotifyMsg pluginCommon.IPv6NbrMacMoveNotifyMsg
type IPv4RouteAddDelFailNotifyMsg struct {
	routeList []asicdInt.IPv4Route
}
type VtepNotifyMsg pluginCommon.VtepNotifyMsg
type PortConfigModeChgNotifyMsg pluginCommon.PortConfigModeChgNotifyMsg
type PortAttrChangeNotifyMsg pluginCommon.PortAttrChangeNotifyMsg
type IPv4VirtualIntfNotifyMsg pluginCommon.IPv4VirtualIntfNotifyMsg
type IPv6VirtualIntfNotifyMsg pluginCommon.IPv6VirtualIntfNotifyMsg
type IPv4VirtualIntfStateNotifyMsg pluginCommon.IPv4VirtualIntfStateNotifyMsg
type IPv6VirtualIntfStateNotifyMsg pluginCommon.IPv6VirtualIntfStateNotifyMsg
type SubIntfPluginObj pluginCommon.SubIntfPluginObj
