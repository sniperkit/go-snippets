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
	"l3/ndp/api"
	"l3/ndp/config"
	"l3/ndp/debug"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/commonDefs"
)

var switchInst asicdClntIntfs.AsicdClntIntf

func GetSwitchInst(fileName string) asicdClntIntfs.AsicdClntIntf {
	notifyHdl := &AsicNotificationHdl{}
	asicdClntInitParams, err := clntIntfs.NewBaseClntInitParams(ASICD_DMN, debug.Logger, notifyHdl, fileName)
	switchInst, err = asicdClntIntfs.NewAsicdClntInit(asicdClntInitParams)
	if err != nil {
		debug.Logger.Err("Error Initializing new asicd client")
		panic(err)
	}
	return switchInst
}

func (notifyHdl *AsicNotificationHdl) ProcessNotification(msg clntIntfs.NotifyMsg) {
	if !api.InitComplete() {
		return
	}
	switch msg.(type) {
	case asicdClntDefs.IPv6IntfNotifyMsg:
		// create/delete ipv6 interface notification case
		ipv6Msg := msg.(asicdClntDefs.IPv6IntfNotifyMsg)
		if switchInst.GetIntfTypeFromIfIndex(ipv6Msg.IfIndex) != commonDefs.IfTypeLoopback {
			if ipv6Msg.MsgType == asicdClntDefs.NOTIFY_IPV6INTF_CREATE {
				debug.Logger.Debug("Received Asicd IPV6 INTF Notfication CREATE:", ipv6Msg)
				api.SendIPIntfNotfication(ipv6Msg.IfIndex, ipv6Msg.IpAddr, ipv6Msg.IntfRef, config.CONFIG_CREATE)
			} else {
				debug.Logger.Debug("Received Asicd IPV6 INTF Notfication DELETE:", ipv6Msg)
				api.SendIPIntfNotfication(ipv6Msg.IfIndex, ipv6Msg.IpAddr, ipv6Msg.IntfRef, config.CONFIG_DELETE)
			}
		}
	case asicdClntDefs.IPv6L3IntfStateNotifyMsg:
		// state up/down for ipv6 interface case
		l3Msg := msg.(asicdClntDefs.IPv6L3IntfStateNotifyMsg)
		// only get state notification if ip type is v6 && not loopback
		if switchInst.GetIntfTypeFromIfIndex(l3Msg.IfIndex) != commonDefs.IfTypeLoopback {
			if l3Msg.IfState == asicdClntDefs.INTF_STATE_UP {
				debug.Logger.Debug("Received Asicd L3 State Notfication UP:", l3Msg)
				api.SendL3PortNotification(l3Msg.IfIndex, config.STATE_UP, l3Msg.IpAddr)
			} else {
				debug.Logger.Debug("Received Asicd L3 State Notfication DOWN:", l3Msg)
				api.SendL3PortNotification(l3Msg.IfIndex, config.STATE_DOWN, l3Msg.IpAddr)
			}
		}
	case asicdClntDefs.L2IntfStateNotifyMsg:
		l2Msg := msg.(asicdClntDefs.L2IntfStateNotifyMsg)
		if l2Msg.IfState == asicdClntDefs.INTF_STATE_UP {
			debug.Logger.Debug("Received Asicd L2 Port Notfication UP:", l2Msg)
			api.SendL3PortNotification(l2Msg.IfIndex, config.STATE_UP, config.L2_NOTIFICATION)
		} else {
			debug.Logger.Debug("Received Asicd L2 Port Notfication DOWN:", l2Msg)
			api.SendL3PortNotification(l2Msg.IfIndex, config.STATE_DOWN, config.L2_NOTIFICATION)
		}
	case asicdClntDefs.VlanNotifyMsg:
		vlanMsg := msg.(asicdClntDefs.VlanNotifyMsg)
		debug.Logger.Debug("Received Asicd Vlan Notfication:", vlanMsg)
		oper := ""
		switch vlanMsg.MsgType {
		case asicdClntDefs.NOTIFY_VLAN_CREATE:
			debug.Logger.Debug("Received Asicd VLAN CREATE")
			oper = config.CONFIG_CREATE
		case asicdClntDefs.NOTIFY_VLAN_DELETE:
			debug.Logger.Debug("Received Asicd VLAN DELETE")
			oper = config.CONFIG_DELETE
		case asicdClntDefs.NOTIFY_VLAN_UPDATE:
			debug.Logger.Debug("Received Asicd VLAN UPDATE")
			oper = config.CONFIG_UPDATE
		}
		api.SendVlanNotification(oper, int32(vlanMsg.VlanId), vlanMsg.VlanIfIndex, vlanMsg.VlanName, vlanMsg.UntagPorts, vlanMsg.TagPorts)
	case asicdClntDefs.IPv6NbrMacMoveNotifyMsg:
		macMoveMsg := msg.(asicdClntDefs.IPv6NbrMacMoveNotifyMsg)
		debug.Logger.Debug("Received Asicd IPv6 Neighbor Mac Move Notification:", macMoveMsg)
		api.SendMacMoveNotification(macMoveMsg.IpAddr, macMoveMsg.IfIndex, macMoveMsg.VlanId)
	case asicdClntDefs.IPv6VirtualIntfNotifyMsg:
		virIpMsg := msg.(asicdClntDefs.IPv6VirtualIntfNotifyMsg)
		if virIpMsg.MsgType == asicdClntDefs.NOTIFY_IPV6VIRTUAL_INTF_CREATE {
			debug.Logger.Debug("Received Asicd IPV6 Virtual INTF Notfication CREATE:", virIpMsg)
			api.SendVirtualIpNotification(virIpMsg.IfIndex, virIpMsg.ParentIfIndex, virIpMsg.IpAddr, virIpMsg.MacAddr, virIpMsg.IfName,
				config.VIRTUAL_CREATE)
		} else {
			debug.Logger.Debug("Received Asicd IPV6 Virtual INTF Notfication DELETE:", virIpMsg)
			api.SendVirtualIpNotification(virIpMsg.IfIndex, virIpMsg.ParentIfIndex, virIpMsg.IpAddr, virIpMsg.MacAddr, virIpMsg.IfName,
				config.VIRTUAL_DELETE)
		}

	case asicdClntDefs.IPv6VirtualIntfStateNotifyMsg:
		virIpMsg := msg.(asicdClntDefs.IPv6VirtualIntfStateNotifyMsg)
		if virIpMsg.IfState == asicdClntDefs.INTF_STATE_UP {
			debug.Logger.Debug("Received Asicd IPV6 Virtual Intf State Up:", virIpMsg)
			api.SendVirtualIpStateMsg(virIpMsg.IfIndex, virIpMsg.IpAddr, config.STATE_UP)
		} else {
			debug.Logger.Debug("Received Asicd IPV6 Virtual Intf State Down:", virIpMsg)
			api.SendVirtualIpStateMsg(virIpMsg.IfIndex, virIpMsg.IpAddr, config.STATE_DOWN)
		}
	}
}
