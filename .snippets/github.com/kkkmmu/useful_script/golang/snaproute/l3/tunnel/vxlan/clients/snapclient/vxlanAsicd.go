package snapclient

import (
	"asicd/asicdCommonDefs"
	"asicdInt"
	"asicdServices"
	"encoding/json"
	"fmt"
	vxlan "l3/tunnel/vxlan/protocol"
	"net"
	"strconv"
	"strings"
	"utils/commonDefs"

	nanomsg "github.com/op/go-nanomsg"
)

type AsicdClient struct {
	VXLANClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

type portVlanValue struct {
	ifIndex string
	refCnt  int
}

var asicdclnt AsicdClient
var PortVlanDb map[uint16][]*portVlanValue

func ConvertVxlanConfigToVxlanAsicdConfig(c *vxlan.VxlanDbEntry) *asicdInt.Vxlan {
	v := &asicdInt.Vxlan{
		Vni: int32(c.VNI),
	}
	for _, untagvlan := range c.UntaggedVlanId {
		v.UntaggedVlanId = append(v.UntaggedVlanId, int16(untagvlan))
	}
	for _, vlan := range c.VlanId {
		v.VlanId = append(v.VlanId, int16(vlan))
	}
	return v
}

func ConvertVtepToVxlanAsicdConfig(vtep *vxlan.VtepDbEntry) *asicdInt.Vtep {

	return &asicdInt.Vtep{
		Vni:            int32(vtep.Vni),
		IfName:         vtep.VtepName,
		SrcIfIndex:     vtep.SrcIfIndex,
		SrcIfName:      vtep.SrcIfName,
		UDP:            int16(vtep.UDP),
		TTL:            int16(vtep.TTL),
		DSCP:           int16(vtep.TOS),
		MTU:            int32(vtep.MTU),
		SrcIp:          vtep.SrcIp.String(),
		DstIp:          vtep.DstIp.String(),
		VlanId:         int16(vtep.VlanId),
		SrcMac:         vtep.SrcMac.String(),
		NextHopIfIndex: vtep.NextHop.IfIndex,
		NextHopIfName:  vtep.NextHop.IfName,
		NextHopIp:      vtep.NextHop.Ip.String(),
	}
}

// ConstructPortConfigMap:
// Let caller know what ports are valid in system
func (intf VXLANSnapClient) ConstructPortConfigMap() {
	currMarker := asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS)
	if asicdclnt.ClientHdl != nil {
		count := asicdServices.Int(asicdCommonDefs.MAX_SYS_PORTS)
		intf.thriftmutex.Lock()
		for {
			bulkInfo, err := asicdclnt.ClientHdl.GetBulkPortState(currMarker, count)
			if err != nil {
				intf.thriftmutex.Unlock()
				return
			}

			bulkCfgInfo, err := asicdclnt.ClientHdl.GetBulkPort(currMarker, count)
			if err != nil {
				intf.thriftmutex.Unlock()
				return
			}

			objCount := int(bulkInfo.Count)
			logger.Info("Creating Ports objcnt", objCount)
			more := bool(bulkInfo.More)
			currMarker = asicdServices.Int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				if bulkInfo.PortStateList[i].IntfRef != bulkCfgInfo.PortList[i].IntfRef {
					logger.Err(fmt.Sprintln("Error IntfRef differ at index", bulkInfo.PortStateList[i].IntfRef, bulkCfgInfo.PortList[i].IntfRef))
				}
				ifindex := bulkInfo.PortStateList[i].IfIndex
				netMac, _ := net.ParseMAC(bulkCfgInfo.PortList[i].MacAddr)
				PortNum, _ := strconv.Atoi(bulkInfo.PortStateList[i].IntfRef)
				config := vxlan.PortConfig{
					PortNum:      int32(PortNum),
					IfIndex:      ifindex,
					Name:         bulkInfo.PortStateList[i].Name,
					HardwareAddr: netMac,
				}
				logger.Info("Creating Port", config)
				go func() {
					serverchannels.VxlanPortCreate <- config
				}()
				logger.Info("Creating Port done")
			}
			if more == false {
				intf.thriftmutex.Unlock()
				return
			}
		}
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) deleteASICdSubscriber() error {
	logger.Info(fmt.Sprintf("Disconnecting from ASICd publisher"))
	if intf.asicdSubSocket != nil {
		intf.asicdSubSocket.Unsubscribe("")
		intf.asicdSubSocket.Close()
	}
	//intf.asicdSubSocket = nil
	return nil
}

// createASICdSubscriber
//
func (intf VXLANSnapClient) createASICdSubscriber() error {
	var err error
	address := asicdCommonDefs.PUB_SOCKET_ADDR
	if intf.asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		logger.Err(fmt.Sprintln("Failed to create ASICd subscribe socket, error:", err))
		return err
	}

	if err = intf.asicdSubSocket.Subscribe(""); err != nil {
		logger.Err(fmt.Sprintln("Failed to subscribe to \"\" on ASICd subscribe socket, error:", err))
		return err
	}

	if _, err = intf.asicdSubSocket.Connect(address); err != nil {
		logger.Err(fmt.Sprintln("Failed to connect to ASICd publisher socket, address:", address, "error:", err))
		return err
	}

	logger.Debug(fmt.Sprintln("Connected to ASICd publisher at address:", address, intf))
	if err = intf.asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		logger.Err(fmt.Sprintln("Failed to set the buffer size for ASICd publisher socket, error:", err))
		return err
	}

	go intf.listenAsicdEvents()

	return nil
}

func (intf VXLANSnapClient) listenAsicdEvents() error {
	logger.Info("Started Listener for Asicd events")
	for {
		//logger.Info(fmt.Sprintln("Read on ASICd subscriber socket...", intf.asicdSubSocket, intf))
		asicdrxBuf, err := intf.asicdSubSocket.Recv(0)
		if err != nil {
			logger.Err(fmt.Sprintln("Recv on ASICd subscriber socket failed with error:", err))
			intf.asicdSubSocketErrCh <- err
			return nil
		}
		//logger.Info(fmt.Sprintln("ASIC subscriber recv returned:", asicdrxBuf))
		intf.asicdSubSocketCh <- asicdrxBuf
	}

	return nil
}

func (intf VXLANSnapClient) processLinkDownEvent(ifindex int) {
	if downcblst, ok := intf.portDownEventList[int32(ifindex)]; ok {
		for _, downcb := range downcblst {
			downcb(int32(ifindex))
		}
	}
}

func (intf VXLANSnapClient) processLinkUpEvent(ifindex int) {
	if upcblst, ok := intf.portUpEventList[int32(ifindex)]; ok {
		for _, upcb := range upcblst {
			upcb(int32(ifindex))
		}
	}
}

func (intf VXLANSnapClient) RegisterLinkUpDownEvents(ifindex int32, upcb vxlan.PortEvtCb, downcb vxlan.PortEvtCb) {
	//if _, ok := intf.portDownEventList[ifindex]; !ok {
	//	intf.portDownEventList[ifindex] = make([]vxlan.PortEvtCb, 0)
	//	intf.portUpEventList[ifindex] = make([]vxlan.PortEvtCb, 0)
	//}
	logger.Debug("Registering for link up/down event cb for", ifindex)
	intf.portDownEventList[ifindex] = append(intf.portDownEventList[ifindex], downcb)
	intf.portUpEventList[ifindex] = append(intf.portUpEventList[ifindex], upcb)
}

func (intf VXLANSnapClient) DeRegisterLinkUpDownEvents(ifindex int32, upcb vxlan.PortEvtCb, downcb vxlan.PortEvtCb) {
	if cbList, ok := intf.portDownEventList[ifindex]; ok {
		for idx, cb := range cbList {
			if &downcb == &cb {
				intf.portDownEventList[ifindex] = append(intf.portDownEventList[ifindex][:idx], intf.portDownEventList[ifindex][idx+1:]...)
				intf.portUpEventList[ifindex] = append(intf.portUpEventList[ifindex][:idx], intf.portUpEventList[ifindex][idx+1:]...)
			}
		}
	}
}

func (intf VXLANSnapClient) GetLinkState(ifname string) string {
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		bulkInfo, err := asicdclnt.ClientHdl.GetBulkPortState(asicdServices.Int(asicdCommonDefs.MIN_SYS_PORTS), asicdServices.Int(asicdCommonDefs.MAX_SYS_PORTS))
		intf.thriftmutex.Unlock()
		if err == nil && bulkInfo.Count != 0 {
			objCount := int64(bulkInfo.Count)
			for i := int64(0); i < objCount; i++ {
				if bulkInfo.PortStateList[i].IntfRef == ifname {
					return bulkInfo.PortStateList[i].OperState
				}
			}
		}
		logger.Err(fmt.Sprintf("GetLinkState: could not get status for port %d, failure in get method\n", ifname))
	}
	return "DOWN"
}

// processAsicdNotification:
// Proceses the various ASICD notifications for the following events
// 1) Vlan membership updates - to create Virtual Access Ports
// 2) IPv4 interface updates - to inform when VTEP src interface has an update
//    for its IP address
func (intf VXLANSnapClient) processAsicdNotification(asicdrxBuf []byte) {

	var msg asicdCommonDefs.AsicdNotification
	err := json.Unmarshal(asicdrxBuf, &msg)
	if err != nil {
		logger.Err("Unable to unmarshal asicdrxBuf:", asicdrxBuf)
		return
	}
	if msg.MsgType == commonDefs.NOTIFY_IPV4INTF_CREATE ||
		msg.MsgType == commonDefs.NOTIFY_IPV4INTF_DELETE {
		var v4Msg commonDefs.IPv4IntfNotifyMsg
		err = json.Unmarshal(msg.Msg, &v4Msg)
		if err != nil {
			logger.Err("Unable to unmashal ipv4intf msg:", msg.Msg)
			return
		}
		intf.updateIpIntf(v4Msg.IpAddr, v4Msg.IfIndex, v4Msg.MsgType)

	} else if msg.MsgType == commonDefs.NOTIFY_IPV6INTF_CREATE ||
		msg.MsgType == commonDefs.NOTIFY_IPV6INTF_DELETE {
		var v6Msg commonDefs.IPv6IntfNotifyMsg
		err = json.Unmarshal(msg.Msg, &v6Msg)
		if err != nil {
			logger.Err("Unable to unmashal ipv4intf msg:", msg.Msg)
			return
		}
		intf.updateIpIntf(v6Msg.IpAddr, v6Msg.IfIndex, v6Msg.MsgType)

	} else if msg.MsgType == commonDefs.NOTIFY_L2INTF_STATE_CHANGE {
		var l2Msg commonDefs.L2IntfStateNotifyMsg
		err = json.Unmarshal(msg.Msg, &l2Msg)
		if err != nil {
			logger.Err("Unable to unmashal ipv4intf msg:", msg.Msg)
			return
		}
		logger.Debug(fmt.Sprintf("Msg L2 linkstatus = %d msg port = %d\n", l2Msg.IfState, l2Msg.IfIndex))
		if l2Msg.IfState == asicdCommonDefs.INTF_STATE_DOWN {
			intf.processLinkDownEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex)) //asicd always sends out link State events for PHY ports
		} else {
			intf.processLinkUpEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex))
		}

	}
	/*
		else if msg.MsgType == commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE {
			var l3v4Msg commonDefs.IPv4L3IntfStateNotifyMsg
			err = json.Unmarshal(msg.Msg, &l3v4Msg)
			if err != nil {
				logger.Err("Unable to unmashal ipv4intf msg:", msg.Msg)
				return
			}
			logger.Info("Msg L3v4 linkstatus = %d msg port = %d\n", l3v4Msg.IfState, l3v4Msg.IfIndex)
			if l3v4Msg.IfState == asicdCommonDefs.INTF_STATE_DOWN {
				intf.processLinkDownEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l3v4Msg.IfIndex)) //asicd always sends out link State events for PHY ports
			} else {
				intf.processLinkUpEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l3v4Msg.IfIndex))
			}
		} else if msg.MsgType == commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE {
			var l3v6Msg commonDefs.IPv6L3IntfStateNotifyMsg
			err = json.Unmarshal(msg.Msg, &l3v6Msg)
			if err != nil {
				logger.Err("Unable to unmashal ipv4intf msg:", msg.Msg)
				return
			}
			logger.Info("Msg L3v6 linkstatus = %d msg port = %d\n", l3v6Msg.IfState, l3v6Msg.IfIndex)
			if l3v6Msg.IfState == asicdCommonDefs.INTF_STATE_DOWN {
				intf.processLinkDownEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l3v6Msg.IfIndex)) //asicd always sends out link State events for PHY ports
			} else {
				intf.processLinkUpEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l3v6Msg.IfIndex))
			}
		}
	*/
	/*
		else if rxMsg.MsgType == asicdCommonDefs.NOTIFY_VLAN_UPDATE {
				//Vlan Create Msg
				logger.Info("Recvd VLAN notification")
				var vlanMsg asicdCommonDefs.VlanNotifyMsg
				err = json.Unmarshal(rxMsg.Msg, &vlanMsg)
				if err != nil {
					logger.Err(fmt.Sprintln("Unable to unmashal vlanNotifyMsg:", rxMsg.Msg))
					return
				}
				intf.updateVlanAccessPorts(vlanMsg, rxMsg.MsgType)
			} */
	/*
		} else if rxMsg.MsgType == asicdCommonDefs.NOTIFY_L3INTF_STATE_CHANGE {
			//L3_INTF_STATE_CHANGE
			logger.Info("Recvd INTF_STATE_CHANGE notification")
			var l3IntfMsg asicdCommonDefs.L3IntfStateNotifyMsg
			err = json.Unmarshal(rxMsg.Msg, &l3IntfMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("Unable to unmashal l3IntfStateNotifyMsg:", rxMsg.Msg))
				return
			}
			server.processL3StateChange(l3IntfMsg)
		} else if rxMsg.MsgType == asicdCommonDefs.NOTIFY_L2INTF_STATE_CHANGE {
			//L2_INTF_STATE_CHANGE
			logger.Info("Recvd INTF_STATE_CHANGE notification")
			var l2IntfMsg asicdCommonDefs.L2IntfStateNotifyMsg
			err = json.Unmarshal(rxMsg.Msg, &l2IntfMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("Unable to unmashal l2IntfStateNotifyMsg:", rxMsg.Msg))
				return
			}
			//server.processL2StateChange(l2IntfMsg)
		} else if rxMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_CREATE ||
			rxMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_UPDATE ||
			rxMsg.MsgType == asicdCommonDefs.NOTIFY_LAG_DELETE {
			logger.Info("Recvd NOTIFY_LAG notification")
			var lagMsg asicdCommonDefs.LagNotifyMsg
			err = json.Unmarshal(rxMsg.Msg, &lagMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("Unable to unmashal lagNotifyMsg:", rxMsg.Msg))
				return
			}
			server.updateLagInfra(lagMsg, rxMsg.MsgType)
		} else if rxMsg.MsgType == asicdCommonDefs.NOTIFY_IPV4NBR_MAC_MOVE {
			//IPv4 Neighbor mac move
			logger.Info("Recvd IPv4NBR_MAC_MOVE notification")
			var macMoveMsg asicdCommonDefs.IPv4NbrMacMoveNotifyMsg
			err = json.Unmarshal(rxMsg.Msg, &macMoveMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("Unable to unmashal macMoveNotifyMsg:", rxMsg.Msg))
				return
			}
			server.processIPv4NbrMacMove(macMoveMsg)
		}
	*/
}

// GetAccessPorts:
// Discovered the ports which have been provisioned with membership to the
// vxlan vlan.   This method will be called after Vxlan has been provisioned
/*
func (intf VXLANSnapClient) GetAccessPortVlan(vlan uint16) {
	curMark := 0
	logger.Info("Calling Asicd for getting Vlan Property")
	count := 100
	more := true
	if asicdclnt.ClientHdl != nil {
		for more {
			bulkVlanInfo, _ := asicdclnt.ClientHdl.GetBulkVlan(asicdInt.Int(curMark), asicdInt.Int(count))
			if bulkVlanInfo == nil {
				break
			}
			// Get bulk on vlan state can re-use curMark and count used by get bulk vlan, as there is a 1:1 mapping in terms of cfg/state objs
			bulkVlanStateInfo, _ := asicdclnt.ClientHdl.GetBulkVlanState(asicdServices.Int(curMark), asicdServices.Int(count))
			if bulkVlanStateInfo == nil {
				break
			}
			objCnt := int(bulkVlanInfo.Count)
			more = bool(bulkVlanInfo.More)
			curMark = int(bulkVlanInfo.EndIdx)
			for i := 0; i < objCnt; i++ {
				ifindex := int(bulkVlanStateInfo.VlanStateList[i].IfIndex)
				config := vxlan.VxlanAccessPortVlan{
					Command:  vxlan.VxlanCommandCreate,
					VlanId:   uint16(asicdCommonDefs.GetIfIndexFromIntfIdAndIntfType(ifindex, commonDefs.IfTypeVlan)),
					IntfList: bulkVlanInfo.VlanList[i].IfIndexList,
				}
				// lets send the config back to the server
				serverchannels.VxlanAccessPortVlanUpdate <- config
			}
		}
	}
	//logger.Info(fmt.Sprintln("Vlan Property Map:", server.vlanPropMap))
}

//
func (intf VXLANSnapClient) updateVlanAccessPorts(msg asicdCommonDefs.VlanNotifyMsg, msgType uint8) {
	vlanId := int(msg.VlanId)
	portList := msg.UntagPorts
	if msgType == asicdCommonDefs.NOTIFY_VLAN_UPDATE { //VLAN UPDATE
		logger.Info(fmt.Sprintln("Received Vlan Update Notification Vlan:", vlanId, "PortList:", portList))
		config := vxlan.VxlanAccessPortVlan{
			Command:  vxlan.VxlanCommandUpdate,
			VlanId:   uint16(vlanId),
			IntfList: portList,
		}
		// lets send the config back to the server
		serverchannels.VxlanAccessPortVlanUpdate <- config

	} else if msgType == asicdCommonDefs.NOTIFY_VLAN_DELETE { // VLAN DELETE
		logger.Info(fmt.Sprintln("Received Vlan Delete Notification Vlan:", vlanId, "PortList:", portList))
		config := vxlan.VxlanAccessPortVlan{
			Command:  vxlan.VxlanCommandDelete,
			VlanId:   uint16(vlanId),
			IntfList: portList,
		}
		// lets send the config back to the server
		serverchannels.VxlanAccessPortVlanUpdate <- config
	}
}
*/
func (intf VXLANSnapClient) updateIpIntf(IpAddr string, ifindex int32, msgType uint8) {
	ipAddr := net.ParseIP(IpAddr)
	IfIndex := ifindex

	nextindex := 0
	count := 1024
	var IfName string
	foundIntf := false
	if asicdclnt.ClientHdl != nil {
		exitnotfound := true
		intf.thriftmutex.Lock()

		// TODO when api is available should just call GetIntf...
		for exitnotfound && !foundIntf {
			bulkIntf, _ := asicdclnt.ClientHdl.GetBulkIntf(asicdInt.Int(nextindex), asicdInt.Int(count))
			//logger.Info(fmt.Sprintln("Return from GetBulkIntf", bulkIntf))
			if bulkIntf == nil {
				exitnotfound = false
			} else {
				for _, intf := range bulkIntf.IntfList {
					if intf.IfName == IfName {
						IfIndex = intf.IfIndex
						foundIntf = true
						//logger.Info(fmt.Sprintln("Found IfName", IfName, bulkIntf))
						break
					}
				}

				if !foundIntf {
					if bulkIntf.More == true {
						nextindex += count
					} else {
						exitnotfound = false
					}
				} else {
					exitnotfound = false
				}
			}
		}
		if foundIntf {
			logicalIntfState, _ := asicdclnt.ClientHdl.GetLogicalIntfState(IfName)
			if logicalIntfState != nil {
				mac, _ := net.ParseMAC(logicalIntfState.SrcMac)
				config := vxlan.VxlanIntfInfo{
					Command:  vxlan.VxlanCommandCreate,
					Ip:       ipAddr,
					Mac:      mac,
					IfIndex:  IfIndex,
					IntfName: IfName,
				}

				if msgType == asicdCommonDefs.NOTIFY_VLAN_CREATE {
					config.Command = vxlan.VxlanCommandCreate
					serverchannels.Vxlanintfinfo <- config

				} else if msgType == asicdCommonDefs.NOTIFY_VLAN_DELETE {
					config.Command = vxlan.VxlanCommandDelete
					serverchannels.Vxlanintfinfo <- config
				}
			}
		}
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) CreateVxlan(vxlan *vxlan.VxlanDbEntry) {
	// convert a vxland config to hw config
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		asicdclnt.ClientHdl.CreateVxlan(ConvertVxlanConfigToVxlanAsicdConfig(vxlan))
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) DeleteVxlan(vxlan *vxlan.VxlanDbEntry) {
	// convert a vxland config to hw config
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		asicdclnt.ClientHdl.DeleteVxlan(ConvertVxlanConfigToVxlanAsicdConfig(vxlan))
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) UpdateVxlan(vni uint32, addvlanlist, delvlanlist, addUntaggedVlan, delUntaggedVlan []uint16) {
	// convert a vxland config to hw config
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()

		convertadd := make([]int16, 0)
		if len(addvlanlist) > 0 {
			for _, vlan := range addvlanlist {
				convertadd = append(convertadd, int16(vlan))
			}
		}
		convertdel := make([]int16, 0)
		if len(delvlanlist) > 0 {
			for _, vlan := range delvlanlist {
				convertdel = append(convertdel, int16(vlan))
			}
		}
		convertadduntag := make([]int16, 0)
		if len(addvlanlist) > 0 {
			for _, vlan := range addUntaggedVlan {
				convertadd = append(convertadd, int16(vlan))
			}
		}
		convertdeluntag := make([]int16, 0)
		if len(delvlanlist) > 0 {
			for _, vlan := range delUntaggedVlan {
				convertdel = append(convertdel, int16(vlan))
			}
		}

		asicdclnt.ClientHdl.UpdateVxlan(int32(vni), convertadd, convertdel, convertadduntag, convertdeluntag)
		intf.thriftmutex.Unlock()
	}
}

/*
func (intf VXLANSnapClient) AddHostToVxlan(vni int32, intfreflist, untagintfreflist []string) {
	if asicdclnt.ClientHdl != nil {
		asicdclnt.ClientHdl.AddHostInterfaceToVxlan(vni, intfreflist, untagintfreflist)
	}

}

func (intf VXLANSnapClient) DelHostFromVxlan(vni int32, intfreflist, untagintfreflist []string) {
	if asicdclnt.ClientHdl != nil {
		asicdclnt.ClientHdl.DelHostInterfaceFromVxlan(vni, intfreflist, untagintfreflist)
	}

}
*/
// CreateVtep:
// Creates a VTEP interface with the ASICD.  Should create an interface within
// the HW as well as within Linux stack.   AsicD also requires that vlan membership is
// provisioned separately from VTEP.  The vlan in question is the VLAN found
// within the VXLAN header.
func (intf VXLANSnapClient) CreateVtep(vtep *vxlan.VtepDbEntry, vteplistener chan<- vxlan.MachineEvent) {
	// convert a vxland config to hw config
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		if vtep.VlanId != 0 {
			// need to create a vlan membership of the vtep vlan Id
			if _, ok := PortVlanDb[vtep.VlanId]; !ok {
				v := &portVlanValue{
					ifIndex: vtep.SrcIfName,
					refCnt:  1,
				}
				pbmp := []string{}
				PortVlanDb[vtep.VlanId] = append(PortVlanDb[vtep.VlanId], v)
				pbmp = append(pbmp, fmt.Sprintf("%d", vtep.NextHop.IfIndex))

				asicdVlan := &asicdServices.Vlan{
					VlanId:   int32(vtep.VlanId),
					IntfList: pbmp,
				}
				asicdclnt.ClientHdl.CreateVlan(asicdVlan)
			} else {
				portExists := -1
				for i, p := range PortVlanDb[vtep.VlanId] {
					if p.ifIndex == vtep.SrcIfName {
						portExists = i
						break
					}
				}
				if portExists == -1 {
					oldpbmp := []string{}
					for _, p := range PortVlanDb[vtep.VlanId] {
						oldpbmp = append(oldpbmp, fmt.Sprintf("%s", p.ifIndex))
					}
					v := &portVlanValue{
						ifIndex: vtep.SrcIfName,
						refCnt:  1,
					}
					PortVlanDb[vtep.VlanId] = append(PortVlanDb[vtep.VlanId], v)
					newpbmp := []string{}
					for _, p := range PortVlanDb[vtep.VlanId] {
						newpbmp = append(newpbmp, fmt.Sprintf("%s", p.ifIndex))
					}

					oldAsicdVlan := &asicdServices.Vlan{
						VlanId:   int32(vtep.VlanId),
						IntfList: oldpbmp,
					}
					newAsicdVlan := &asicdServices.Vlan{
						VlanId:   int32(vtep.VlanId),
						IntfList: newpbmp,
					}
					// note if the thrift attribute id's change then
					// this attr may need to be updated
					attrset := []bool{false, true, false}
					op := []*asicdServices.PatchOpInfo{}
					asicdclnt.ClientHdl.UpdateVlan(oldAsicdVlan, newAsicdVlan, attrset, op)
				} else {
					v := PortVlanDb[vtep.VlanId][portExists]
					v.refCnt++
					PortVlanDb[vtep.VlanId][portExists] = v
				}
			}
		}
		// create the vtep
		id, _ := asicdclnt.ClientHdl.CreateVxlanVtep(ConvertVtepToVxlanAsicdConfig(vtep))
		intf.thriftmutex.Unlock()

		data := vxlan.VtepCreateCfgData{
			VtepName: vtep.VtepHandleName,
			IfIndex:  id,
		}

		event := vxlan.MachineEvent{
			E:    vxlan.VxlanVtepEventHwConfigComplete,
			Src:  vxlan.VXLANSnapClientStr,
			Data: data,
		}

		vteplistener <- event
	}
}

// DeleteVtep:
// Delete a VTEP interface with the ASICD.  Should create an interface within
// the HW as well as within Linux stack. AsicD also requires that vlan membership is
// provisioned separately from VTEP.  The vlan in question is the VLAN found
// within the VXLAN header.
func (intf VXLANSnapClient) DeleteVtep(vtep *vxlan.VtepDbEntry) {
	// convert a vxland config to hw config
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		// delete the vtep
		asicdclnt.ClientHdl.DeleteVxlanVtep(ConvertVtepToVxlanAsicdConfig(vtep))

		if vtep.VlanId != 0 {
			// update the vlan the vtep was using
			if _, ok := PortVlanDb[vtep.VlanId]; ok {
				portExists := -1
				for i, p := range PortVlanDb[vtep.VlanId] {
					if p.ifIndex == vtep.SrcIfName {
						portExists = i
						break
					}
				}
				if portExists != -1 {
					v := PortVlanDb[vtep.VlanId][portExists]
					v.refCnt--
					PortVlanDb[vtep.VlanId][portExists] = v

					// lets remove this port from the vlan
					if v.refCnt == 0 {
						oldpbmp := []string{}
						for _, p := range PortVlanDb[vtep.VlanId] {
							oldpbmp = append(oldpbmp, fmt.Sprintf("%s", p.ifIndex))
						}
						// remove from local list
						PortVlanDb[vtep.VlanId] = append(PortVlanDb[vtep.VlanId][:portExists], PortVlanDb[vtep.VlanId][portExists+1:]...)
						newpbmp := []string{}
						for _, p := range PortVlanDb[vtep.VlanId] {
							newpbmp = append(newpbmp, fmt.Sprintf("%s", p.ifIndex))
						}

						oldAsicdVlan := &asicdServices.Vlan{
							VlanId:   int32(vtep.VlanId),
							IntfList: oldpbmp,
						}
						newAsicdVlan := &asicdServices.Vlan{
							VlanId:   int32(vtep.VlanId),
							IntfList: newpbmp,
						}
						// note if the thrift attribute id's change then
						// this attr may need to be updated
						attrset := []bool{false, true, false}
						op := []*asicdServices.PatchOpInfo{}
						asicdclnt.ClientHdl.UpdateVlan(oldAsicdVlan, newAsicdVlan, attrset, op)
					}
					// lets remove the vlan
					if len(PortVlanDb[vtep.VlanId]) == 0 {

						asicdVlan := &asicdServices.Vlan{
							VlanId: int32(vtep.VlanId),
						}
						asicdclnt.ClientHdl.DeleteVlan(asicdVlan)
						delete(PortVlanDb, vtep.VlanId)

					}
				}
			}
		}
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) UpdateVtepAttr(vtepName string, vni uint32, tos, ttl uint8, mtu uint16) {
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		asicdclnt.ClientHdl.UpdateVxlanVtepAttr(vtepName, int32(vni), int16(tos), int16(ttl), int32(mtu))
		intf.thriftmutex.Unlock()
	}
}

func (intf VXLANSnapClient) GetAllVlans() []uint16 {
	vlanlist := make([]uint16, 0)

	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		vlanInfo, _ := asicdclnt.ClientHdl.GetBulkVlan(1, 4094)
		intf.thriftmutex.Unlock()

		for _, vlandata := range vlanInfo.VlanList {
			vlanlist = append(vlanlist, uint16(vlandata.VlanId))

		}
	}

	return vlanlist
}

func (intf VXLANSnapClient) GetIntfInfo(IfName string, intfchan chan<- vxlan.MachineEvent) {
	// TODO
	nextindex := 0
	count := 1024
	var IfIndex int32
	var mac net.HardwareAddr
	foundIntf := false
	if asicdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		exitnotfound := true
		// TODO when api is available should just call GetIntf...
		for exitnotfound && !foundIntf {
			bulkIntf, _ := asicdclnt.ClientHdl.GetBulkIntf(asicdInt.Int(nextindex), asicdInt.Int(count))
			if bulkIntf == nil {
				exitnotfound = false
			} else {
				logger.Info(fmt.Sprintln("Return from GetBulkIntf"))
				for _, intf := range bulkIntf.IntfList {
					if intf.IfName == IfName {
						IfIndex = intf.IfIndex
						foundIntf = true
						logger.Debug(fmt.Sprintln("Found IfName", IfName))
						break
					}
				}

				if !foundIntf {
					if bulkIntf.More == true {
						nextindex += count
					} else {
						exitnotfound = false
					}
				} else {
					exitnotfound = false
				}
			}
		}
		// if we found the interface all other objects at least the logical interface should exist
		if foundIntf {

			// interface exists now what type of interface is it
			intfType := asicdCommonDefs.GetIntfTypeFromIfIndex(IfIndex)
			switch intfType {
			case commonDefs.IfTypePort:
				portNum := asicdCommonDefs.GetIntfIdFromIfIndex(int32(IfIndex))
				phyIntfState, _ := asicdclnt.ClientHdl.GetPort(fmt.Sprintf("%d", portNum))
				logger.Debug(fmt.Sprintln("Return from GetPort", phyIntfState))
				if phyIntfState != nil {
					mac, _ = net.ParseMAC(phyIntfState.MacAddr)
				} else {
					// did not find src mac exiting
					intf.thriftmutex.Unlock()
					return
				}
			//case IfTypeLag:
			//case IfTypeVlan:
			//case IfTypeP2P:
			//case IfTypeBcast:
			case commonDefs.IfTypeLoopback:
				logicalIntfState, _ := asicdclnt.ClientHdl.GetLogicalIntfState(IfName)
				logger.Debug(fmt.Sprintln("Return from GetLogicalIntfState", logicalIntfState))
				if logicalIntfState != nil {
					mac, _ = net.ParseMAC(logicalIntfState.SrcMac)
				} else {
					// did not find src mac exiting
					intf.thriftmutex.Unlock()
					return
				}

				//case IfTypeSecondary:
				//case IfTypeVirtual:
				//case IfTypeVtep:
			}

			ipV4, err := asicdclnt.ClientHdl.GetIPv4IntfState(IfName)
			logger.Debug(fmt.Sprintln("Return from GetIPv4IntfState ", IfName, ipV4))
			if err == nil {
				ipaddrstr := strings.Split(ipV4.IpAddr, "/")[0]
				ip := net.ParseIP(ipaddrstr)
				config := vxlan.VxlanIntfInfo{
					Command:  vxlan.VxlanCommandCreate,
					Ip:       ip,
					Mac:      mac,
					IfIndex:  IfIndex,
					IntfName: IfName,
				}
				logger.Debug(fmt.Sprintln("Intf info complete, sending config to channel", config))
				event := vxlan.MachineEvent{
					E:    vxlan.VxlanVtepEventSrcInterfaceResolved,
					Src:  vxlan.VXLANSnapClientStr,
					Data: config,
				}
				intfchan <- event
			}
		} else {
			logger.Err(fmt.Sprintln("ERROR Did not find IfName", IfName))
		}
		intf.thriftmutex.Unlock()
	}
}
