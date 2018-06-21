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

package flexswitch

import (
	"asicd/asicdCommonDefs"
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"utils/commonDefs"
	"utils/logging"
)

var asicdSubSocket *nanomsg.SubSocket

type processMsg func(uint8, []byte, *logging.Writer) (commonDefs.AsicdNotifyMsg, error)

var AsicdMsgMap map[uint8]processMsg = map[uint8]processMsg{
	asicdCommonDefs.NOTIFY_L2INTF_STATE_CHANGE:           processL2IntdStateNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE:      processIPv4L3IntfStateNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE:      processIPv6L3IntfStateNotifyMsg,
	asicdCommonDefs.NOTIFY_VLAN_CREATE:                   processVlanNotifyMsg,
	asicdCommonDefs.NOTIFY_VLAN_DELETE:                   processVlanNotifyMsg,
	asicdCommonDefs.NOTIFY_VLAN_UPDATE:                   processVlanNotifyMsg,
	asicdCommonDefs.NOTIFY_LOGICAL_INTF_CREATE:           processLogicalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_LOGICAL_INTF_DELETE:           processLogicalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_LOGICAL_INTF_UPDATE:           processLogicalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4INTF_CREATE:               processIPv4IntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4INTF_DELETE:               processIPv4IntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6INTF_CREATE:               processIPv6IntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6INTF_DELETE:               processIPv6IntfNotifyMsg,
	asicdCommonDefs.NOTIFY_LAG_CREATE:                    processLagNotifyMsg,
	asicdCommonDefs.NOTIFY_LAG_DELETE:                    processLagNotifyMsg,
	asicdCommonDefs.NOTIFY_LAG_UPDATE:                    processLagNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4NBR_MAC_MOVE:              processIPv4NbrMacMoveNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6NBR_MAC_MOVE:              processIPv6NbrMacMoveNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4_ROUTE_CREATE_FAILURE:     processIPv4RouteAddDelNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4_ROUTE_DELETE_FAILURE:     processIPv4RouteAddDelNotifyMsg,
	asicdCommonDefs.NOTIFY_PORT_CONFIG_MODE_CHANGE:       processPortConfigModeChgNotifyMsg,
	asicdCommonDefs.NOTIFY_PORT_ATTR_CHANGE:              processPortAttrChangeNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4VIRTUAL_INTF_CREATE:       processIPv4VirutalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4VIRTUAL_INTF_DELETE:       processIPv4VirutalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6VIRTUAL_INTF_CREATE:       processIPv6VirutalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6VIRTUAL_INTF_DELETE:       processIPv6VirutalIntfNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE: processIPv4VirtualIntfStateNotifyMsg,
	asicdCommonDefs.NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE: processIPv6VirtualIntfStateNotifyMsg,
}

func processL2IntdStateNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var l2Msg asicdCommonDefs.L2IntfStateNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &l2Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal L2IntfStateNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.L2IntfStateNotifyMsg{
		MsgType: rxMsgType,
		IfIndex: l2Msg.IfIndex,
		IfState: l2Msg.IfState,
	}
	return msg, nil
}

func processIPv4L3IntfStateNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var l3Msg asicdCommonDefs.IPv4L3IntfStateNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &l3Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv4L3IntfStateNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv4L3IntfStateNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  l3Msg.IpAddr,
		IfIndex: l3Msg.IfIndex,
		IfState: l3Msg.IfState,
	}
	return msg, nil
}

func processIPv6L3IntfStateNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var l3Msg asicdCommonDefs.IPv6L3IntfStateNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &l3Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv6L3IntfStateNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv6L3IntfStateNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  l3Msg.IpAddr,
		IfIndex: l3Msg.IfIndex,
		IfState: l3Msg.IfState,
	}
	return msg, nil
}

func processIPv4VirtualIntfStateNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var l3Msg asicdCommonDefs.IPv4VirtualIntfStateNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &l3Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv4VirtualIntfStateNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv4VirtualIntfStateNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  l3Msg.IpAddr,
		IfIndex: l3Msg.IfIndex,
		IfState: l3Msg.IfState,
	}
	return msg, nil
}

func processIPv6VirtualIntfStateNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var l3Msg asicdCommonDefs.IPv6VirtualIntfStateNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &l3Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv6VirtualIntfStateNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv6VirtualIntfStateNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  l3Msg.IpAddr,
		IfIndex: l3Msg.IfIndex,
		IfState: l3Msg.IfState,
	}
	return msg, nil
}

func processVlanNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var vlanMsg asicdCommonDefs.VlanNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &vlanMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal vlanCreate:", rxMsg))
		return msg, err
	}
	msg = commonDefs.VlanNotifyMsg{
		MsgType:     rxMsgType,
		VlanId:      vlanMsg.VlanId,
		VlanIfIndex: vlanMsg.VlanIfIndex,
		VlanName:    vlanMsg.VlanName,
		TagPorts:    vlanMsg.TagPorts,
		UntagPorts:  vlanMsg.UntagPorts,
	}
	return msg, nil
}

func processLogicalIntfNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var logicalMsg asicdCommonDefs.LogicalIntfNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &logicalMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal logical Interface:", rxMsg))
		return msg, err
	}
	msg = commonDefs.LogicalIntfNotifyMsg{
		MsgType:         rxMsgType,
		IfIndex:         logicalMsg.IfIndex,
		LogicalIntfName: logicalMsg.LogicalIntfName,
	}
	return msg, nil

}

func processIPv4IntfNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var ipv4Msg asicdCommonDefs.IPv4IntfNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &ipv4Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv4 Intf:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv4IntfNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  ipv4Msg.IpAddr,
		IfIndex: ipv4Msg.IfIndex,
		IntfRef: ipv4Msg.IntfRef,
	}

	return msg, nil
}

func processIPv6IntfNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var ipv6Msg asicdCommonDefs.IPv6IntfNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &ipv6Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv6 Intf:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv6IntfNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  ipv6Msg.IpAddr,
		IfIndex: ipv6Msg.IfIndex,
		IntfRef: ipv6Msg.IntfRef,
	}

	return msg, nil
}

func processIPv4VirutalIntfNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var ipv4Msg asicdCommonDefs.IPv4VirtualIntfNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &ipv4Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv4 Virtual Intf:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv4VirtualIntfNotifyMsg{
		MsgType:       rxMsgType,
		IfIndex:       ipv4Msg.IfIndex,
		ParentIfIndex: ipv4Msg.ParentIfIndex,
		IpAddr:        ipv4Msg.IpAddr,
		MacAddr:       ipv4Msg.MacAddr,
		IfName:        ipv4Msg.IfName,
	}

	return msg, nil
}

func processIPv6VirutalIntfNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var ipv4Msg asicdCommonDefs.IPv6VirtualIntfNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &ipv4Msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal IPv6 Virtual Intf:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv6VirtualIntfNotifyMsg{
		MsgType:       rxMsgType,
		IfIndex:       ipv4Msg.IfIndex,
		ParentIfIndex: ipv4Msg.ParentIfIndex,
		IpAddr:        ipv4Msg.IpAddr,
		MacAddr:       ipv4Msg.MacAddr,
		IfName:        ipv4Msg.IfName,
	}

	return msg, nil
}

func processLagNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var lagMsg asicdCommonDefs.LagNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &lagMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal Lag Create:", rxMsg))
		return msg, err
	}
	msg = commonDefs.LagNotifyMsg{
		MsgType:     rxMsgType,
		LagName:     lagMsg.LagName,
		IfIndex:     lagMsg.IfIndex,
		IfIndexList: lagMsg.IfIndexList,
	}

	return msg, nil
}

func processIPv4NbrMacMoveNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var macMoveMsg asicdCommonDefs.IPv4NbrMacMoveNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &macMoveMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal Mac Move:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv4NbrMacMoveNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  macMoveMsg.IpAddr,
		IfIndex: macMoveMsg.IfIndex,
		VlanId:  macMoveMsg.VlanId,
	}
	return msg, err
}

func processIPv6NbrMacMoveNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var macMoveMsg asicdCommonDefs.IPv6NbrMacMoveNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &macMoveMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal Mac Move:", rxMsg))
		return msg, err
	}
	msg = commonDefs.IPv6NbrMacMoveNotifyMsg{
		MsgType: rxMsgType,
		IpAddr:  macMoveMsg.IpAddr,
		IfIndex: macMoveMsg.IfIndex,
		VlanId:  macMoveMsg.VlanId,
	}
	return msg, err
}

func processIPv4RouteAddDelNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var msg commonDefs.AsicdNotifyMsg
	return msg, nil
}

func processPortConfigModeChgNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var portCfgChgMsg asicdCommonDefs.PortConfigModeChgNotifyMsg
	var msg commonDefs.AsicdNotifyMsg
	err := json.Unmarshal(rxMsg, &portCfgChgMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal PortConfigModeChgNotifyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.PortConfigModeChgNotifyMsg{
		IfIndex: portCfgChgMsg.IfIndex,
		OldMode: portCfgChgMsg.OldMode,
		NewMode: portCfgChgMsg.NewMode,
	}
	return msg, nil
}

func processPortAttrChangeNotifyMsg(rxMsgType uint8, rxMsg []byte, logger *logging.Writer) (commonDefs.AsicdNotifyMsg, error) {
	var attrChgMsg asicdCommonDefs.PortAttrChangeNotifyMsg
	var msg commonDefs.AsicdNotifyMsg

	err := json.Unmarshal(rxMsg, &attrChgMsg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmashal  PortConfigMtuChgNotigyMsg:", rxMsg))
		return msg, err
	}
	msg = commonDefs.PortAttrChangeNotifyMsg{
		IfIndex:     attrChgMsg.IfIndex,
		Mtu:         attrChgMsg.Mtu,
		Description: attrChgMsg.Description,
		Pvid:        attrChgMsg.Pvid,
		AttrMask:    attrChgMsg.AttrMask,
	}
	return msg, nil
}

func listenForASICdUpdates(address string, logger *logging.Writer) (err error) {
	if asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		logger.Err(fmt.Sprintln("Failed to create ASICd subscribe socket, error:", err))
		return err
	}

	if err = asicdSubSocket.Subscribe(""); err != nil {
		logger.Err(fmt.Sprintln("Failed to subscribe to \"\" on ASICd subscribe socket, error:", err))
		return err
	}

	if _, err = asicdSubSocket.Connect(address); err != nil {
		logger.Err(fmt.Sprintln("Failed to connect to ASICd publisher socket, address:", address, "error:", err))
		return err
	}

	logger.Debug(fmt.Sprintln("Connected to ASICd publisher at address:", address))
	if err = asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		logger.Err(fmt.Sprintln("Failed to set the buffer size for ASICd publisher socket, error:", err))
		return err
	}
	return nil
}

func InitFSAsicdSubscriber(nHdl commonDefs.AsicdClientStruct) error {
	err := listenForASICdUpdates(asicdCommonDefs.PUB_SOCKET_ADDR, nHdl.Logger)
	if err != nil {
		nHdl.Logger.Err(fmt.Sprintln("Unable to open FS ASICd Subscriber", err))
		return err
	}

	go createASICdSubscriber(nHdl)
	return nil
}

func createASICdSubscriber(nHdl commonDefs.AsicdClientStruct) {
	for {
		nHdl.Logger.Debug("Read on ASICd subscriber socket...")
		asicdrxBuf, err := asicdSubSocket.Recv(0)
		if err != nil {
			nHdl.Logger.Err(fmt.Sprintln("Recv on ASICd subscriber socket failed with error:", err))
			//intfSubClientIntf.ProcessIntfNotification(asicdrxBuf)
			continue
		}
		processFSAsicdNotification(asicdrxBuf, nHdl)
		//evtCh.HALSubSocketCh <- asicdrxBuf
	}
}

func processNotification(rxMsg asicdCommonDefs.AsicdNotification, clientHdl commonDefs.AsicdClientStruct) {
	if clientHdl.NMap[rxMsg.MsgType] {
		processNotifyMsg := AsicdMsgMap[rxMsg.MsgType]
		msg, err := processNotifyMsg(rxMsg.MsgType, rxMsg.Msg, clientHdl.Logger)
		if err == nil {
			clientHdl.NHdl.ProcessNotification(msg)
		}
	}
}

func processFSAsicdNotification(asicdrxBuf []byte, nHdl commonDefs.AsicdClientStruct) {
	var rxMsg asicdCommonDefs.AsicdNotification
	err := json.Unmarshal(asicdrxBuf, &rxMsg)
	if err != nil {
		nHdl.Logger.Err(fmt.Sprintln("Unable to unmarshal asicdrxBuf:", asicdrxBuf))
		return
	}

	processNotification(rxMsg, nHdl)
}
