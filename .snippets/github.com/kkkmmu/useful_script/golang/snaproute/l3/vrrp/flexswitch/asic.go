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

package flexswitch

import (
	"l3/vrrp/api"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"sync"
	"syscall"
	"utils/asicdClient"
	"utils/commonDefs"
)

var switchInst *commonDefs.AsicdClientStruct = nil
var fsPlugin asicdClient.AsicdClientIntf
var once sync.Once

func initAsicdNotification() commonDefs.AsicdNotification {
	nMap := make(commonDefs.AsicdNotification)
	nMap = commonDefs.AsicdNotification{
		commonDefs.NOTIFY_L2INTF_STATE_CHANGE:       false, // 0
		commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE:  true,  // 1
		commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE:  true,  // 2
		commonDefs.NOTIFY_VLAN_CREATE:               false, // 3
		commonDefs.NOTIFY_VLAN_DELETE:               false, // 4
		commonDefs.NOTIFY_VLAN_UPDATE:               false, // 5
		commonDefs.NOTIFY_LOGICAL_INTF_CREATE:       false, // 6
		commonDefs.NOTIFY_LOGICAL_INTF_DELETE:       false, // 7
		commonDefs.NOTIFY_LOGICAL_INTF_UPDATE:       false, // 8
		commonDefs.NOTIFY_IPV4INTF_CREATE:           true,  // 9
		commonDefs.NOTIFY_IPV4INTF_DELETE:           true,  // 10
		commonDefs.NOTIFY_IPV6INTF_CREATE:           true,  // 11
		commonDefs.NOTIFY_IPV6INTF_DELETE:           true,  // 12
		commonDefs.NOTIFY_LAG_CREATE:                false, // 13
		commonDefs.NOTIFY_LAG_DELETE:                false, // 14
		commonDefs.NOTIFY_LAG_UPDATE:                false, // 15
		commonDefs.NOTIFY_IPV4NBR_MAC_MOVE:          true,  // 16
		commonDefs.NOTIFY_IPV6NBR_MAC_MOVE:          true,  // 17
		commonDefs.NOTIFY_IPV4_ROUTE_CREATE_FAILURE: false, // 17
		commonDefs.NOTIFY_IPV4_ROUTE_DELETE_FAILURE: false, // 18
		commonDefs.NOTIFY_IPV6_ROUTE_CREATE_FAILURE: false, // 19
		commonDefs.NOTIFY_IPV6_ROUTE_DELETE_FAILURE: false, // 20
		commonDefs.NOTIFY_VTEP_CREATE:               false, // 21
		commonDefs.NOTIFY_VTEP_DELETE:               false, // 22
		commonDefs.NOTIFY_MPLSINTF_STATE_CHANGE:     false, // 23
		commonDefs.NOTIFY_MPLSINTF_CREATE:           false, // 24
		commonDefs.NOTIFY_MPLSINTF_DELETE:           false, // 25
		commonDefs.NOTIFY_PORT_CONFIG_MODE_CHANGE:   false, // 26
	}
	return nMap
}

type AsicNotificationHdl struct {
}

func GetSwitchInst() *commonDefs.AsicdClientStruct {
	once.Do(func() {
		notifyMap := initAsicdNotification()
		notifyHdl := &AsicNotificationHdl{}
		switchInst = &commonDefs.AsicdClientStruct{
			NHdl: notifyHdl,
			NMap: notifyMap,
		}
	})
	return switchInst
}

func (notifyHdl *AsicNotificationHdl) ProcessNotification(msg commonDefs.AsicdNotifyMsg) {
	if !api.InitComplete() {
		return
	}
	switch msg.(type) {
	case commonDefs.IPv6IntfNotifyMsg:
		// create/delete ipv6 interface notification case
		ipv6Msg := msg.(commonDefs.IPv6IntfNotifyMsg)
		if !fsPlugin.IsLoopbackType(ipv6Msg.IfIndex) {
			ipInfo := &common.BaseIpInfo{
				IfIndex:   ipv6Msg.IfIndex,
				IntfRef:   ipv6Msg.IntfRef,
				IpAddr:    ipv6Msg.IpAddr,
				IpType:    syscall.AF_INET6,
				OperState: common.STATE_DOWN, // during create lets hardcode the oper state to be down
			}
			if ipv6Msg.MsgType == commonDefs.NOTIFY_IPV6INTF_CREATE {
				debug.Logger.Debug("Received Asicd IPV6 INTF Notfication CREATE:", ipv6Msg)
				ipInfo.MsgType = common.IP_MSG_CREATE
				api.SendIpIntfNotification(ipInfo)
			} else {
				debug.Logger.Debug("Received Asicd IPV6 INTF Notfication DELETE:", ipv6Msg)
				ipInfo.MsgType = common.IP_MSG_DELETE
				api.SendIpIntfNotification(ipInfo)
			}
		}

	case commonDefs.IPv6L3IntfStateNotifyMsg:
		// state up/down for ipv6 interface case
		ipv6Msg := msg.(commonDefs.IPv6L3IntfStateNotifyMsg)
		// only get state notification if ip type is v6 && not loopback
		if !fsPlugin.IsLoopbackType(ipv6Msg.IfIndex) {
			ipInfo := &common.BaseIpInfo{
				IfIndex: ipv6Msg.IfIndex,
				IpAddr:  ipv6Msg.IpAddr,
				IpType:  syscall.AF_INET6,
				MsgType: common.IP_MSG_STATE_CHANGE,
			}
			if ipv6Msg.IfState == commonDefs.INTF_STATE_UP {
				debug.Logger.Debug("Received Asicd L3 State Notfication UP:", ipv6Msg)
				ipInfo.OperState = common.STATE_UP
				api.SendIpIntfNotification(ipInfo)
			} else {
				debug.Logger.Debug("Received Asicd L3 State Notfication DOWN:", ipv6Msg)
				ipInfo.OperState = common.STATE_DOWN
				api.SendIpIntfNotification(ipInfo)
			}
		}
	case commonDefs.IPv4IntfNotifyMsg:
		ipv4Msg := msg.(commonDefs.IPv4IntfNotifyMsg)
		if !fsPlugin.IsLoopbackType(ipv4Msg.IfIndex) {
			ipInfo := &common.BaseIpInfo{
				IfIndex:   ipv4Msg.IfIndex,
				IntfRef:   ipv4Msg.IntfRef,
				IpAddr:    ipv4Msg.IpAddr,
				IpType:    syscall.AF_INET,
				OperState: common.STATE_DOWN, // during create lets hardcode the oper state to be down
			}
			if ipv4Msg.MsgType == commonDefs.NOTIFY_IPV4INTF_CREATE {
				debug.Logger.Debug("Received Asicd IPV4 INTF Notfication CREATE:", ipv4Msg)
				ipInfo.MsgType = common.IP_MSG_CREATE
				api.SendIpIntfNotification(ipInfo)
			} else {
				debug.Logger.Debug("Received Asicd IPV4 INTF Notfication DELETE:", ipv4Msg)
				ipInfo.MsgType = common.IP_MSG_DELETE
				api.SendIpIntfNotification(ipInfo)
			}
		}
	case commonDefs.IPv4L3IntfStateNotifyMsg:
		// state up/down for ipv6 interface case
		ipv4Msg := msg.(commonDefs.IPv4L3IntfStateNotifyMsg)
		// only get state notification if ip type is v6 && not loopback
		if !fsPlugin.IsLoopbackType(ipv4Msg.IfIndex) {
			ipInfo := &common.BaseIpInfo{
				IfIndex: ipv4Msg.IfIndex,
				IpAddr:  ipv4Msg.IpAddr,
				IpType:  syscall.AF_INET,
				MsgType: common.IP_MSG_STATE_CHANGE,
			}
			if ipv4Msg.IfState == commonDefs.INTF_STATE_UP {
				debug.Logger.Debug("Received Asicd L3 State Notfication UP:", ipv4Msg)
				ipInfo.OperState = common.STATE_UP
				api.SendIpIntfNotification(ipInfo)
			} else {
				debug.Logger.Debug("Received Asicd L3 State Notfication DOWN:", ipv4Msg)
				ipInfo.OperState = common.STATE_DOWN
				api.SendIpIntfNotification(ipInfo)
			}
		}
	}
}
