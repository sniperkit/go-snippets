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

package server

import (
	"errors"
	"fmt"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"l3/vrrp/packet"
	"net"
	"syscall"
	"utils/netUtils"
)

func (svr *VrrpServer) getIPv4Intfs() {
	debug.Logger.Info("Get all ipv4 interfaces from asicd")
	ipv4Info, err := svr.SwitchPlugin.GetAllIPv4IntfState()
	if err != nil {
		debug.Logger.Err("Failed to get all IPv4 interfaces, err:", err)
		return
	}

	for _, obj := range ipv4Info {
		// do not care for loopback interface
		if svr.SwitchPlugin.IsLoopbackType(obj.IfIndex) {
			continue
		}
		v4Obj := &V4Intf{}
		ipInfo := &common.BaseIpInfo{
			IntfRef:   obj.IntfRef,
			IfIndex:   obj.IfIndex,
			OperState: obj.OperState,
			IpAddr:    obj.IpAddr,
		}
		v4Obj.Init(ipInfo)
		svr.V4[obj.IfIndex] = v4Obj
		svr.V4IntfRefToIfIndex[obj.IntfRef] = obj.IfIndex
	}
}

func (svr *VrrpServer) getIPv6Intfs() {
	debug.Logger.Info("Get all ipv6 interfaces from asicd")
	ipv6Info, err := svr.SwitchPlugin.GetAllIPv6IntfState()
	if err != nil {
		debug.Logger.Err("Failed to get all IPv6 interfaces, err:", err)
		return
	}
	for _, obj := range ipv6Info {
		// store only ipv6 link local as vrrp is point to point
		if !svr.SwitchPlugin.IsLoopbackType(obj.IfIndex) && netUtils.IsIpv6LinkLocal(obj.IpAddr) {
			v6Obj := &V6Intf{}
			ipInfo := &common.BaseIpInfo{
				IntfRef:   obj.IntfRef,
				IfIndex:   obj.IfIndex,
				OperState: obj.OperState,
				IpAddr:    obj.IpAddr,
			}
			v6Obj.Init(ipInfo)
			svr.V6[obj.IfIndex] = v6Obj
			svr.V6IntfRefToIfIndex[obj.IntfRef] = obj.IfIndex
		}
	}
}

func (svr *VrrpServer) GetIPIntfs() {
	svr.getIPv4Intfs()
	svr.getIPv6Intfs()
}

func constructIntfKey(intfRef string, vrid int32, ipType int) KeyInfo {
	return KeyInfo{intfRef, vrid, ipType}
}

func (svr *VrrpServer) ValidateCreateConfig(cfg *common.IntfCfg) (bool, error) {
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	if _, exists := svr.Intf[key]; exists {
		return false, errors.New(fmt.Sprintln("Vrrp Interface already created for config:", cfg,
			"only update is allowed"))
	}
	if cfg.IpType == syscall.AF_INET {
		// check if ipv4 address is configured on the intfRef
		_, v4exists := svr.V4IntfRefToIfIndex[cfg.IntfRef]
		if !v4exists {
			return false, errors.New(fmt.Sprintln("Vrrp V4 cannot be configured as no l3 Interface found for:", cfg.IntfRef))
		}
	} else if cfg.IpType == syscall.AF_INET6 {
		// check if ipv6 address is configured on the intRef
		_, v6exists := svr.V6IntfRefToIfIndex[cfg.IntfRef]

		if !v6exists {
			return false, errors.New(fmt.Sprintln("Vrrp V6 cannot be configured as no l3 Interface found for:", cfg.IntfRef))
		}
	} else {
		return false, errors.New("Invalid ip type")
	}
	debug.Logger.Info("Validation of create config:", *cfg, "is success")
	return true, nil
}

func (svr *VrrpServer) ValidateUpdateConfig(cfg *common.IntfCfg) (bool, error) {
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	intf, exists := svr.Intf[key]
	if !exists {
		return false, errors.New(fmt.Sprintln("Vrrp Interface doesn't exists for key:", key, "config:", cfg,
			"please do create before updating entry"))
	}
	if intf.Config.VRID != cfg.VRID {
		return false, errors.New("Updating VRID is not allowed")
	}
	return true, nil
}

func (svr *VrrpServer) ValidateDeleteConfig(cfg *common.IntfCfg) (bool, error) {
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	if _, exists := svr.Intf[key]; !exists {
		return false, errors.New(fmt.Sprintln("Vrrp Interface was not created for config:", cfg))
	}
	return true, nil
}

func (svr *VrrpServer) ValidConfiguration(cfg *common.IntfCfg) (bool, error) {
	if cfg.VRID == 0 {
		return false, errors.New(fmt.Sprintln(VRRP_INVALID_VRID, cfg.VRID))
	}
	switch cfg.Operation {
	case common.CREATE:
		return svr.ValidateCreateConfig(cfg)
	case common.UPDATE:
		return svr.ValidateUpdateConfig(cfg)
	case common.DELETE:
		return svr.ValidateDeleteConfig(cfg)
	}
	return false, errors.New("Invalid Operation received for Vrrp Interface Config")
}

/* Update Intf List which can be used during state
 */
func (svr *VrrpServer) updateIntfList(key KeyInfo, ipType int, insert bool) {
	debug.Logger.Debug("updating state object intf key list for:", key, ipType, insert)
	switch insert {
	case true:
		// new vrrp configured insert the entry into lists
		if ipType == syscall.AF_INET {
			svr.v4Intfs = append(svr.v4Intfs, key)
		} else {
			svr.v6Intfs = append(svr.v6Intfs, key)
		}
	case false:
		if ipType == syscall.AF_INET {
			for idx, _ := range svr.v4Intfs {
				if svr.v4Intfs[idx] == key {
					svr.v4Intfs = append(svr.v4Intfs[:idx], svr.v4Intfs[idx+1:]...)
					break
				}
			}
		} else {
			for idx, _ := range svr.v6Intfs {
				if svr.v6Intfs[idx] == key {
					svr.v6Intfs = append(svr.v6Intfs[:idx], svr.v6Intfs[idx+1:]...)
					break
				}
			}
		}
	}
	debug.Logger.Debug("after updating len of v4intfs is:", len(svr.v4Intfs), "len of v6Intf is:", len(svr.v6Intfs))
}

/* During Create of Virtual Interface Enable should always be set to false... when
 * vrrp interface becomes master it will request for the interface to be in up state
 * Input: (vrrp interface config, virtual mac)
 */
func (svr *VrrpServer) CreateVirtualIntf(cfg *common.IntfCfg, vMac string) {
	debug.Logger.Info("Vrrp Creating Virtual Interface for:", cfg.IntfRef, cfg.VirtualIPAddr, vMac)
	switch cfg.IpType {
	case syscall.AF_INET:
		svr.SwitchPlugin.CreateVirtualIPv4Intf(cfg.IntfRef, cfg.VirtualIPAddr, vMac, false /*enable*/)
	case syscall.AF_INET6:
		svr.SwitchPlugin.CreateVirtualIPv6Intf(cfg.IntfRef, cfg.VirtualIPAddr, vMac, false /*enable*/)
	}
}

/* Update Virtual Interface by changing the state as requested
 * Input: (intfRef, virtual ip address, macAddress, enable)
 */
func (svr *VrrpServer) UpdateVirtualIntf(virtualIpInfo *common.VirtualIpInfo) {
	switch virtualIpInfo.IpType {
	case syscall.AF_INET:
		ip, _, _ := net.ParseCIDR(virtualIpInfo.IpAddr)
		ip = ip.To4()
		if virtualIpInfo.Enable {
			// Call arp to delete the nextHop Entry
			debug.Logger.Info("Calling Arp to delete nexthop entry:", ip.String())
			svr.ArpClient.DeleteArpEntry(ip.String())
		}
		err := svr.SwitchPlugin.UpdateVirtualIPv4Intf(virtualIpInfo.IntfRef, virtualIpInfo.IpAddr, virtualIpInfo.MacAddr, virtualIpInfo.Enable)
		if err != nil {
			debug.Logger.Err("Failed to update virtual ip in asicd")
		}
	case syscall.AF_INET6:
		ip, _, _ := net.ParseCIDR(virtualIpInfo.IpAddr)
		ip = ip.To16()
		if virtualIpInfo.Enable {
			debug.Logger.Info("Calling Ndp to delete nexthop entry:", ip.String())
			svr.NdpClient.DeleteNdpEntry(ip.String())
		}
		svr.SwitchPlugin.UpdateVirtualIPv6Intf(virtualIpInfo.IntfRef, virtualIpInfo.IpAddr, virtualIpInfo.MacAddr, virtualIpInfo.Enable)
	}
}

/* During Delete of Virtual Interface Enable should always be set to false
 * Input: (vrrp interface config, virtual mac)
 */
func (svr *VrrpServer) DeleteVirtualIntf(cfg *common.IntfCfg, vMac string) {
	debug.Logger.Info("Vrrp Deleting Virtual Interface for:", cfg.IntfRef, cfg.VirtualIPAddr, vMac)
	switch cfg.IpType {
	case syscall.AF_INET:
		svr.SwitchPlugin.DeleteVirtualIPv4Intf(cfg.IntfRef, cfg.VirtualIPAddr, vMac, false /*enable*/)
	case syscall.AF_INET6:
		svr.SwitchPlugin.DeleteVirtualIPv6Intf(cfg.IntfRef, cfg.VirtualIPAddr, vMac, false /*enable*/)
	}
}

/*
 *  Handling Vrrp Interface Configuration
 */
func (svr *VrrpServer) HandlerVrrpIntfCreateConfig(cfg *common.IntfCfg) {
	debug.Logger.Info("Received vrrp interface create common:", *cfg)
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	intf, exists := svr.Intf[key]
	if exists {
		debug.Logger.Err("During Create we should not have any entry in the DB")
		return
	}
	debug.Logger.Debug("Constructed Key for vrrp interface is:", key)
	l3Info := &common.BaseIpInfo{}
	l3Info.IntfRef = cfg.IntfRef
	l3Info.IpType = cfg.IpType
	var ipIntf IPIntf
	var ifIndex int32
	// Get DB based on config version
	switch cfg.IpType {
	case syscall.AF_INET:
		ifIndex, exists = svr.V4IntfRefToIfIndex[cfg.IntfRef]
		debug.Logger.Debug("v4 ifIndex found in reverse map for:", cfg.IntfRef, "is:", ifIndex, "exists:", exists)
		if exists {
			ipIntf, exists = svr.V4[ifIndex]
		}
	// if cross reference exists then only set l3Info else just pass go defaults and it will updated
	// later once we have configured ipv4 or ipv6 interface
	case syscall.AF_INET6:
		ifIndex, exists = svr.V6IntfRefToIfIndex[cfg.IntfRef]
		debug.Logger.Debug("v6 ifIndex found in reverse map for:", cfg.IntfRef, "is:", ifIndex, "exists:", exists)
		if exists {
			ipIntf, exists = svr.V6[ifIndex]
		}
	}
	// if entry exists then only you should get information from DB otherwise it should be nothing
	// Information collected from DB is L3 interface ip address and operation state
	if exists {
		debug.Logger.Debug("ip interface exists and hence get information from DB")
		l3Info.IfIndex = ifIndex
		ipIntf.GetObjFromDb(l3Info)
	} else {
		debug.Logger.Err("cannot create config as no l3 interface is found")
		return
	}
	intf.InitVrrpIntf(cfg, l3Info, svr.VirtualIpCh, svr.UpdateRxCh, svr.UpdateTxCh)
	// if l3 interface was created before vrrp interface then there might be a chance that interface is already
	// up... if that's the case then lets start fsm right away
	if l3Info.OperState == common.STATE_UP && svr.GlobalConfig.Enable && cfg.AdminState {
		// during create always call start fsm
		intf.StartFsm()
		debug.Logger.Info("Fsm is initialized for the interface, now calling create virtual interface")
		svr.CreateVirtualIntf(cfg, intf.GetVMac())
	}
	debug.Logger.Debug("storing interface at location:", key)
	svr.Intf[key] = intf
	ipIntf.SetVrrpIntfKey(key)
	svr.updateIntfList(key, cfg.IpType, true /*insert*/)
}

func (svr *VrrpServer) HandleVrrpIntfUpdateConfig(cfg *common.IntfCfg) {
	debug.Logger.Info("Received interface update config:", *cfg)
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	intf, exists := svr.Intf[key]
	if !exists {
		debug.Logger.Err("Cannot perform update as no interface found in db for key:", key)
		return
	}
	oper := intf.UpdateConfig(cfg)
	if oper == common.CREATE {
		svr.CreateVirtualIntf(cfg, intf.GetVMac())
	} else if oper == common.DELETE {
		svr.DeleteVirtualIntf(cfg, intf.GetVMac())
	}
	svr.Intf[key] = intf
}

func (svr *VrrpServer) HandleVrrpIntfDeleteConfig(cfg *common.IntfCfg) {
	debug.Logger.Info("Received vrrp interface delete config:", *cfg)
	key := constructIntfKey(cfg.IntfRef, cfg.VRID, cfg.IpType)
	intf, exists := svr.Intf[key]
	if !exists {
		// this should never happen as Validate should have taken care of this
		debug.Logger.Err("no vrrp interface found for:", key)
		return
	}
	mac := intf.GetVMac()
	intf.DeInitVrrpIntf()
	if svr.GlobalConfig.Enable {
		svr.DeleteVirtualIntf(cfg, mac)
	}
	delete(svr.Intf, key)
	svr.updateIntfList(key, cfg.IpType, false /*delete*/)
}

func (svr *VrrpServer) HandleVrrpIntfConfig(cfg *common.IntfCfg) {
	debug.Logger.Info("svr handling vrrp interface configuration:", *cfg)
	switch cfg.Operation {
	case common.CREATE:
		svr.HandlerVrrpIntfCreateConfig(cfg)
		svr.addConfiguredIntfCount(cfg.IpType)
	case common.UPDATE:
		svr.HandleVrrpIntfUpdateConfig(cfg)
	case common.DELETE:
		svr.HandleVrrpIntfDeleteConfig(cfg)
		svr.removeConfiguredIntfCount(cfg.IpType)
	}

}

func (svr *VrrpServer) HandleProtocolMacEntry(add bool) {
	switch add {
	case true:
		svr.SwitchPlugin.EnablePacketReception(packet.VRRP_PROTOCOL_MAC, -1, 1)
		svr.SwitchPlugin.EnablePacketReception(packet.VRRP_V6_PROTOCOL_MAC, -1, 1)
	case false:
		svr.SwitchPlugin.DisablePacketReception(packet.VRRP_PROTOCOL_MAC, -1, 1)
		svr.SwitchPlugin.EnablePacketReception(packet.VRRP_V6_PROTOCOL_MAC, -1, 1)
	}
}

/*
 * We can get vrrp interface configurations before even vrrp is enabled....let's handle that scenario here
 * by starting fsm if vrrp enabled
 * by stopping fsm if vrrp disabled
 */
func (svr *VrrpServer) HandleVrrpEnableDisable(enable bool) {
	debug.Logger.Info("vrrp globally:", enable, "handling it")
	for key, intf := range svr.Intf {
		if enable {
			intf.UpdateIpState()
			svr.CreateVirtualIntf(intf.Config, intf.GetVMac())
		} else {
			intf.StopFsm()
			svr.DeleteVirtualIntf(intf.Config, intf.GetVMac())
		}
		svr.Intf[key] = intf
	}
}

func (svr *VrrpServer) HandleGlobalConfig(gCfg *common.GlobalConfig) {
	debug.Logger.Info("Handling Global Config for:", *gCfg)
	svr.GlobalConfig.Enable = gCfg.Enable
	switch gCfg.Operation {
	case common.CREATE:
		debug.Logger.Info("Vrrp Global Object Created")
		svr.GlobalConfig.Vrf = gCfg.Vrf
	case common.UPDATE:
		debug.Logger.Info("Vrrp Global Updated:", *svr.GlobalConfig)
		if gCfg.Enable {
			debug.Logger.Info("Vrrp Enabled, configuring Protocol Mac")
			svr.HandleProtocolMacEntry(true /*Enable*/)
		} else {
			debug.Logger.Info("Vrrp Disabled, deleting Protocol Mac")
			svr.HandleProtocolMacEntry(false /*Enable*/)
		}
		svr.HandleVrrpEnableDisable(gCfg.Enable)
	case common.DB_UPDATE:
		debug.Logger.Info("Vrrp Global DB Updated:", *svr.GlobalConfig)
		svr.GlobalConfig.Vrf = gCfg.Vrf
		if gCfg.Enable {
			debug.Logger.Info("Vrrp Enabled, configuring Protocol Mac during restart")
			svr.HandleProtocolMacEntry(true /*Enable*/)
		}
	}
	svr.updateGlobalStatus()
}

func (svr *VrrpServer) HandleIpStateChange(msg *common.BaseIpInfo) {
	var ipIntf IPIntf
	var exists bool

	switch msg.IpType {
	case syscall.AF_INET:
		ipIntf, exists = svr.V4[msg.IfIndex]
	case syscall.AF_INET6:
		ipIntf, exists = svr.V6[msg.IfIndex]
	}

	if !exists {
		debug.Logger.Err("No Entry found for:", *msg, "during state:", msg.OperState, "notification")
		return
	}
	// update sw state for ip interface with new information
	ipIntf.Update(msg)

	// check if vrrp is enabled or not
	if svr.GlobalConfig.Enable == false {
		debug.Logger.Info("Vrrp is not enabled and hence just updating ip information")
		return
	}
	// get the vrrp interface key
	key := ipIntf.GetVrrpIntfKey()
	if key == nil {
		// if no key then it means that no vrrp interface is created
		debug.Logger.Warning("No vrrp interface attached to ip interface:", ipIntf.GetIntfRef())
		return
	}
	intf, exists := svr.Intf[*key]
	if !exists {
		debug.Logger.Warning("No Vrrp Interface configured and hence nothing to do")
		return
	}
	intf.UpdateOperState(msg.OperState)
	intf.UpdateIpState()
	svr.Intf[*key] = intf
}

func (svr *VrrpServer) HandleIpNotification(msg *common.BaseIpInfo) {
	debug.Logger.Info("handling ip notification:", *msg)
	switch msg.MsgType {
	case common.IP_MSG_CREATE:
		switch msg.IpType {
		case syscall.AF_INET:
			v4, exists := svr.V4[msg.IfIndex]
			if !exists {
				v4 = &V4Intf{}
				v4.Init(msg)
				svr.V4[msg.IfIndex] = v4
				svr.V4IntfRefToIfIndex[msg.IntfRef] = msg.IfIndex
				debug.Logger.Info("Reverse v4 ip intf to ifIndex cached for:", msg.IntfRef, "-------->", msg.IfIndex)
			}
		case syscall.AF_INET6:
			v6, exists := svr.V6[msg.IfIndex]
			if !exists {
				v6 = &V6Intf{}
				v6.Init(msg)
				svr.V6IntfRefToIfIndex[msg.IntfRef] = msg.IfIndex
				debug.Logger.Info("Reverse v6 ip intf to ifIndex cached for:", msg.IntfRef, "-------->", msg.IfIndex)
			} else {
				v6.updateIp(msg.IpAddr)
			}
			svr.V6[msg.IfIndex] = v6
		}
	case common.IP_MSG_DELETE:
		switch msg.IpType {
		case syscall.AF_INET:
			v4, exists := svr.V4[msg.IfIndex]
			if exists {
				v4.DeInit(msg)
				delete(svr.V4, msg.IfIndex)
			}
		case syscall.AF_INET6:
			// most likely we will get two delete one for linkscope and other for global-scope
			v6, exists := svr.V6[msg.IfIndex]
			if exists {
				v6.DeInit(msg)
				if !netUtils.IsIpv6LinkLocal(msg.IpAddr) {
					v6.Cfg.GlobalScopeIp = ""
				} else {
					v6.Cfg.Info.IpAddr = ""
				}
				if v6.Cfg.GlobalScopeIp == "" && v6.Cfg.Info.IpAddr == "" {
					delete(svr.V6, msg.IfIndex)
				}
			}
		}

	case common.IP_MSG_STATE_CHANGE:
		svr.HandleIpStateChange(msg)
	}
}

func (svr *VrrpServer) addConfiguredIntfCount(ipType int) {
	if ipType == syscall.AF_INET {
		svr.globalState.V4Intfs++
	} else if ipType == syscall.AF_INET6 {
		svr.globalState.V6Intfs++
	}
}

func (svr *VrrpServer) removeConfiguredIntfCount(ipType int) {
	if ipType == syscall.AF_INET {
		svr.globalState.V4Intfs--
	} else if ipType == syscall.AF_INET6 {
		svr.globalState.V6Intfs--
	}
}

func (svr *VrrpServer) updateRxCount() {
	svr.globalState.TotalRxFrames++
}

func (svr *VrrpServer) updateTxCount() {
	svr.globalState.TotalTxFrames++
}

func (svr *VrrpServer) updateGlobalStatus() {
	svr.globalState.Status = svr.GlobalConfig.Enable
	svr.globalState.Vrf = svr.GlobalConfig.Vrf
}
