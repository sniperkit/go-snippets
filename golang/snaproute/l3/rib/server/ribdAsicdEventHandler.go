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

// ribdEventHandler.go
package server

import (
	defs "l3/rib/ribdCommonDefs"
	"net"
	"ribd"
	"strconv"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/clntUtils/clntIntfs"
)

func (ribdServiceHandler *RIBDServer) processAsicdNotification(msg clntIntfs.NotifyMsg) {
	switch msg.(type) {
	case asicdClntDefs.LogicalIntfNotifyMsg:
		logicalIntfNotifyMsg := msg.(asicdClntDefs.LogicalIntfNotifyMsg)
		if logicalIntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_LOGICAL_INTF_CREATE {
			ribdServiceHandler.Logger.Info("NOTIFY_LOGICAL_INTF_CREATE received")
			ribdServiceHandler.ProcessLogicalIntfCreateEvent(logicalIntfNotifyMsg)
		}
		break
	case asicdClntDefs.LagNotifyMsg:
		lagIntfNotifyMsg := msg.(asicdClntDefs.LagNotifyMsg)
		if lagIntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_LAG_CREATE {
			ribdServiceHandler.Logger.Info(".NOTIFY_LAG_CREATE received")
			ribdServiceHandler.ProcessLagIntfCreateEvent(lagIntfNotifyMsg)
		}
		break
	case asicdClntDefs.VlanNotifyMsg:
		vlanNotifyMsg := msg.(asicdClntDefs.VlanNotifyMsg)
		if vlanNotifyMsg.MsgType == asicdClntDefs.NOTIFY_VLAN_CREATE {
			ribdServiceHandler.Logger.Info("asicdCommonDefs.NOTIFY_VLAN_CREATE")
			ribdServiceHandler.ProcessVlanCreateEvent(vlanNotifyMsg)
		}
		break
	case asicdClntDefs.IPv4L3IntfStateNotifyMsg:
		ipv4L3IntfStateNotifyMsg := msg.(asicdClntDefs.IPv4L3IntfStateNotifyMsg)
		if ipv4L3IntfStateNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV4_L3INTF_STATE_CHANGE event")
			ribdServiceHandler.Logger.Info("Msg linkstatus = ", ipv4L3IntfStateNotifyMsg.IfState, " msg  ifId ", ipv4L3IntfStateNotifyMsg.IfIndex)
			if ipv4L3IntfStateNotifyMsg.IfState == 0 {
				ribdServiceHandler.ProcessIPv4IntfDownEvent(ipv4L3IntfStateNotifyMsg.IpAddr, ipv4L3IntfStateNotifyMsg.IfIndex)
			} else {
				ribdServiceHandler.ProcessIPv4IntfUpEvent(ipv4L3IntfStateNotifyMsg.IpAddr, ipv4L3IntfStateNotifyMsg.IfIndex)
			}
		}
		break
	case asicdClntDefs.IPv6L3IntfStateNotifyMsg:
		ipv6L3IntfStateNotifyMsg := msg.(asicdClntDefs.IPv6L3IntfStateNotifyMsg)
		if ipv6L3IntfStateNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV6_L3INTF_STATE_CHANGE event")
			ribdServiceHandler.Logger.Info("Msg linkstatus = ", ipv6L3IntfStateNotifyMsg.IfState, " msg  ifId ", ipv6L3IntfStateNotifyMsg.IfIndex)
			if ipv6L3IntfStateNotifyMsg.IfState == 0 {
				ribdServiceHandler.ProcessIPv6IntfDownEvent(ipv6L3IntfStateNotifyMsg.IpAddr, ipv6L3IntfStateNotifyMsg.IfIndex)
			} else {
				ribdServiceHandler.ProcessIPv6IntfUpEvent(ipv6L3IntfStateNotifyMsg.IpAddr, ipv6L3IntfStateNotifyMsg.IfIndex)
			}
		}
		break
	case asicdClntDefs.IPv4IntfNotifyMsg:
		ipv4IntfNotifyMsg := msg.(asicdClntDefs.IPv4IntfNotifyMsg)
		if ipv4IntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV4INTF_CREATE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV4INTF_CREATE event")
			ribdServiceHandler.Logger.Info("Received NOTIFY_IPV4INTF_CREATE ipAddr ", ipv4IntfNotifyMsg.IpAddr, " ifIndex = ", ipv4IntfNotifyMsg.IfIndex)
			ribdServiceHandler.ProcessIPv4IntfCreateEvent(ipv4IntfNotifyMsg)
		} else if ipv4IntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV4INTF_DELETE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV4INTF_DELETE  event")
			ribdServiceHandler.Logger.Info("Received ipv4 intf delete with ipAddr ", ipv4IntfNotifyMsg.IpAddr, " ifIndex = ", ipv4IntfNotifyMsg.IfIndex)
			ribdServiceHandler.ProcessIPv4IntfDeleteEvent(ipv4IntfNotifyMsg)
		}
		break
	case asicdClntDefs.IPv6IntfNotifyMsg:
		ipv6IntfNotifyMsg := msg.(asicdClntDefs.IPv6IntfNotifyMsg)
		if ipv6IntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV6INTF_CREATE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV6INTF_CREATE event")
			ribdServiceHandler.Logger.Info("Received NOTIFY_IPV6INTF_CREATE ipAddr ", ipv6IntfNotifyMsg.IpAddr, " ifIndex = ", ipv6IntfNotifyMsg.IfIndex)
			ribdServiceHandler.ProcessIPv6IntfCreateEvent(ipv6IntfNotifyMsg)
		} else if ipv6IntfNotifyMsg.MsgType == asicdClntDefs.NOTIFY_IPV6INTF_DELETE {
			ribdServiceHandler.Logger.Info("NOTIFY_IPV6INTF_DELETE event")
			ribdServiceHandler.Logger.Info("Received ipv6 intf delete with ipAddr ", ipv6IntfNotifyMsg.IpAddr, " ifIndex = ", ipv6IntfNotifyMsg.IfIndex)
			ribdServiceHandler.ProcessIPv6IntfDeleteEvent(ipv6IntfNotifyMsg)
		}
		break
	}
}

func (ribdServiceHandler *RIBDServer) ProcessLogicalIntfCreateEvent(logicalIntfNotifyMsg asicdClntDefs.LogicalIntfNotifyMsg) {
	ifId := logicalIntfNotifyMsg.IfIndex
	if IntfIdNameMap == nil {
		IntfIdNameMap = make(map[int32]IntfEntry)
	}
	intfEntry := IntfEntry{name: logicalIntfNotifyMsg.LogicalIntfName}
	ribdServiceHandler.Logger.Info("ProcessLogicalIntfCreateEvent:Updating IntfIdMap at index ", ifId, " with name ", logicalIntfNotifyMsg.LogicalIntfName)
	IntfIdNameMap[int32(ifId)] = intfEntry
	if IfNameToIfIndex == nil {
		IfNameToIfIndex = make(map[string]int32)
	}
	IfNameToIfIndex[logicalIntfNotifyMsg.LogicalIntfName] = ifId

}
func (ribdServiceHandler *RIBDServer) ProcessLagIntfCreateEvent(lagIntfNotifyMsg asicdClntDefs.LagNotifyMsg) {
	ifId := lagIntfNotifyMsg.IfIndex
	if IntfIdNameMap == nil {
		IntfIdNameMap = make(map[int32]IntfEntry)
	}
	intfEntry := IntfEntry{name: lagIntfNotifyMsg.LagName}
	ribdServiceHandler.Logger.Info("ProcessLagIntfCreateEvent:Updating IntfIdMap at index ", ifId, " with name ", lagIntfNotifyMsg.LagName)
	IntfIdNameMap[int32(ifId)] = intfEntry
	if IfNameToIfIndex == nil {
		IfNameToIfIndex = make(map[string]int32)
	}
	IfNameToIfIndex[lagIntfNotifyMsg.LagName] = ifId

}
func (ribdServiceHandler *RIBDServer) ProcessVlanCreateEvent(vlanNotifyMsg asicdClntDefs.VlanNotifyMsg) {
	ifId := vlanNotifyMsg.VlanIfIndex //asicdClntIntfs.GetIfIndexFromIntfIdAndIntfType(int(vlanNotifyMsg.VlanId), commonDefs.IfTypeVlan)
	ribdServiceHandler.Logger.Info("vlanId ", vlanNotifyMsg.VlanId, " ifId:", ifId)
	if IntfIdNameMap == nil {
		IntfIdNameMap = make(map[int32]IntfEntry)
	}
	intfEntry := IntfEntry{name: vlanNotifyMsg.VlanName}
	IntfIdNameMap[int32(ifId)] = intfEntry
	if IfNameToIfIndex == nil {
		IfNameToIfIndex = make(map[string]int32)
	}
	IfNameToIfIndex[vlanNotifyMsg.VlanName] = ifId
}
func (ribdServiceHandler *RIBDServer) ProcessIPv4IntfCreateEvent(msg asicdClntDefs.IPv4IntfNotifyMsg) {

	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(msg.IpAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	ribdServiceHandler.Logger.Info("Calling createv4Route with ipaddr ", ipAddrStr, " mask ", ipMaskStr, " nextHopIntRef: ", strconv.Itoa(int(msg.IfIndex)))
	//fmt.Println("Calling createv4Route with ipaddr ", ipAddrStr, " mask ", ipMaskStr, " nextHopIntRef: ", strconv.Itoa(int(msg.IfIndex)))
	cfg := ribd.IPv4Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "0.0.0.0",
		NextHopIntRef: strconv.Itoa(int(msg.IfIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)

	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.Add,
	}
}
func (ribdServiceHandler *RIBDServer) ProcessIPv6IntfCreateEvent(msg asicdClntDefs.IPv6IntfNotifyMsg) {
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(msg.IpAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	ribdServiceHandler.Logger.Info("Calling createRoute with ipaddr ", ipAddrStr, " mask ", ipMaskStr, " nextHopIntRef: ", strconv.Itoa(int(msg.IfIndex)))
	cfg := ribd.IPv6Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "::",
		NextHopIntRef: strconv.Itoa(int(msg.IfIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)

	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.Addv6,
	}
}
func (ribdServiceHandler *RIBDServer) ProcessIPv4IntfDeleteEvent(msg asicdClntDefs.IPv4IntfNotifyMsg) {
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(msg.IpAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	ribdServiceHandler.Logger.Info("Calling deletev4Route with ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	cfg := ribd.IPv4Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "0.0.0.0",
		NextHopIntRef: strconv.Itoa(int(msg.IfIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)
	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.Del,
	}

}
func (ribdServiceHandler *RIBDServer) ProcessIPv6IntfDeleteEvent(msg asicdClntDefs.IPv6IntfNotifyMsg) {
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(msg.IpAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	ribdServiceHandler.Logger.Info("Calling deleteRoute with ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	cfg := ribd.IPv6Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "::",
		NextHopIntRef: strconv.Itoa(int(msg.IfIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)
	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.Delv6,
	}
}
