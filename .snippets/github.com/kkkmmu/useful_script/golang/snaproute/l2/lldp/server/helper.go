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
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"l2/lldp/config"
	"l2/lldp/packet"
	"l2/lldp/utils"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

/*  Init l2 port information for global runtime information
 */
func (intf *LLDPGlobalInfo) InitRuntimeInfo(portConf *config.PortInfo) {
	intf.Port = *portConf
	intf.RxInfo = packet.RxInit()
	intf.TxInfo = packet.TxInit(LLDP_DEFAULT_TX_INTERVAL, LLDP_DEFAULT_TX_HOLD_MULTIPLIER)
	intf.RxKill = make(chan bool)
	intf.TxDone = make(chan bool)
	intf.RxLock = &sync.RWMutex{}
}

/*  De-Init l2 port information
 */
func (intf *LLDPGlobalInfo) DeInitRuntimeInfo() {
	intf.StopCacheTimer()
	intf.DeletePcapHandler()
	intf.FreeDynamicMemory()
}

/*  Delete l2 port pcap handler
 */
func (intf *LLDPGlobalInfo) DeletePcapHandler() {
	if intf.PcapHandle != nil {
		// Send go routine kill signal right away before even we do anything else
		if intf.RxInfo.RxRunning {
			debug.Logger.Debug("sending go routine kill signal for:", intf.Port.Name)
			intf.RxKill <- true
			intf.RxInfo.RxRunning = false
			<-intf.RxKill
		}
		// @FIXME: some bug in close handling that causes 5 mins delay
		intf.PcapHandle.Close()
		intf.PcapHandle = nil
		debug.Logger.Debug("pcap closed for:", intf.Port.Name)
	}
}

/*  Based on configuration we will enable disable lldp per port
 */
func (intf *LLDPGlobalInfo) Enable() {
	intf.enable = true
}

/*  Based on configuration we will enable disable lldp per port
 */
func (intf *LLDPGlobalInfo) Disable() {
	intf.enable = false
}

/*  Check LLDP is disabled or not
 */
func (intf *LLDPGlobalInfo) isDisabled() bool {
	return !intf.enable
}

/*  Check LLDP is enabled or not
 */
func (intf *LLDPGlobalInfo) isEnabled() bool {
	return intf.enable
}

/*  Stop RX cache timer
 */
func (intf *LLDPGlobalInfo) StopCacheTimer() {
	intf.RxLock.Lock()
	defer intf.RxLock.Unlock()
	if intf.RxInfo.ClearCacheTimer == nil {
		return
	}
	intf.RxInfo.ClearCacheTimer.Stop()
}

/*  Return back all the memory which was allocated using new
 */
func (intf *LLDPGlobalInfo) FreeDynamicMemory() {
	intf.RxLock.Lock()
	defer intf.RxLock.Unlock()
	intf.RxInfo.RxFrame = nil
	intf.RxInfo.RxLinkInfo = nil
}

/*  Create Pcap Handler
 */
func (intf *LLDPGlobalInfo) CreatePcapHandler(lldpSnapshotLen int32, lldpPromiscuous bool, lldpTimeout time.Duration) error {
	debug.Logger.Debug("Creating Pcap for port:", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
	pcapHdl, err := pcap.OpenLive(intf.Port.Name, lldpSnapshotLen, lldpPromiscuous, lldpTimeout)
	if err != nil {
		debug.Logger.Err(fmt.Sprintln("Creating Pcap Handler failed for", intf.Port.Name, "Error:", err))
		return errors.New("Creating Pcap Failed")
	}
	err = pcapHdl.SetBPFFilter(LLDP_BPF_FILTER)
	if err != nil {
		debug.Logger.Err(fmt.Sprintln("setting filter:", LLDP_BPF_FILTER, "for", intf.Port.Name,
			"failed with error:", err))
		return errors.New("Setting BPF Filter Failed")
	}
	intf.PcapHandle = pcapHdl
	debug.Logger.Debug("Pcap Created for port:", intf.Port.Name, "ifIndex:", intf.Port.IfIndex)
	return nil
}

/*  Get Chassis Id info
 *	 Based on SubType Return the string, mac address then form string using
 *	 net package
 */
func (intf *LLDPGlobalInfo) GetChassisIdInfo() string {

	retVal := ""
	switch intf.RxInfo.RxFrame.ChassisID.Subtype {
	case layers.LLDPChassisIDSubTypeReserved:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubTypeChassisComp:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubtypeIfaceAlias:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubTypePortComp:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubTypeMACAddr:
		var mac net.HardwareAddr
		mac = intf.RxInfo.RxFrame.ChassisID.ID
		return mac.String()
	case layers.LLDPChassisIDSubTypeNetworkAddr:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubtypeIfaceName:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPChassisIDSubTypeLocal:
		debug.Logger.Debug("Need to handle this case")
	default:
		return retVal

	}
	return retVal
}

/*  Get Port Id info
 *	 Based on SubType Return the string, mac address then form string using
 *	 net package
 */
func (intf *LLDPGlobalInfo) GetPortIdInfo() string {

	retVal := ""
	switch intf.RxInfo.RxFrame.PortID.Subtype {
	case layers.LLDPPortIDSubtypeReserved:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPPortIDSubtypeIfaceAlias:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPPortIDSubtypePortComp:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPPortIDSubtypeMACAddr:
		var mac net.HardwareAddr
		mac = intf.RxInfo.RxFrame.ChassisID.ID
		return mac.String()
	case layers.LLDPPortIDSubtypeNetworkAddr:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPPortIDSubtypeIfaceName:
		return string(intf.RxInfo.RxFrame.PortID.ID)
	case layers.LLDPPortIDSubtypeAgentCircuitID:
		debug.Logger.Debug("Need to handle this case")
	case layers.LLDPPortIDSubtypeLocal:
		debug.Logger.Debug("Need to handle this case")
	default:
		return retVal

	}
	return retVal
}

/*  Get System Capability info
 *	 Based on booleans value Return the string, which states what system capabilities are enabled
 */
func (intf *LLDPGlobalInfo) GetSystemCap() string {
	retVal := ""
	systemCap := intf.RxInfo.RxLinkInfo.SysCapabilities.SystemCap
	if systemCap.Other {
		retVal += "Other, "
	}
	if systemCap.Repeater {
		retVal += "Repeater, "
	}
	if systemCap.Bridge {
		retVal += "Bridge, "
	}
	if systemCap.WLANAP {
		retVal += "WlanAP, "
	}
	if systemCap.Router {
		retVal += "Router, "
	}
	if systemCap.Phone {
		retVal += "Phone, "
	}
	if systemCap.DocSis {
		retVal += "DocSis, "
	}
	if systemCap.StationOnly {
		retVal += "StationOnly, "
	}
	if systemCap.CVLAN {
		retVal += "CVlan, "
	}
	if systemCap.SVLAN {
		retVal += "SVlan, "
	}
	if systemCap.TMPR {
		retVal += "TMPR, "
	}

	return strings.TrimSuffix(retVal, ", ")
}

/*  Get Enabled Capability info
 *	 Based on booleans value Return the string, which states what enabled capabilities are enabled
 */
func (intf *LLDPGlobalInfo) GetEnabledCap() string {
	retVal := ""
	enabledCap := intf.RxInfo.RxLinkInfo.SysCapabilities.EnabledCap
	if enabledCap.Other {
		retVal += "Other, "
	}
	if enabledCap.Repeater {
		retVal += "Repeater, "
	}
	if enabledCap.Bridge {
		retVal += "Bridge, "
	}
	if enabledCap.WLANAP {
		retVal += "WlanAP, "
	}
	if enabledCap.Router {
		retVal += "Router, "
	}
	if enabledCap.Phone {
		retVal += "Phone, "
	}
	if enabledCap.DocSis {
		retVal += "DocSis, "
	}
	if enabledCap.StationOnly {
		retVal += "StationOnly, "
	}
	if enabledCap.CVLAN {
		retVal += "CVlan, "
	}
	if enabledCap.SVLAN {
		retVal += "SVlan, "
	}
	if enabledCap.TMPR {
		retVal += "TMPR, "
	}

	return strings.TrimSuffix(retVal, ", ")
}

/*  Get Peer Host Name information
 *
 */
func (intf *LLDPGlobalInfo) GetPeerHostName() string {
	return intf.RxInfo.RxLinkInfo.SysName
}

/*  Get Peer Host Name information
 *
 */
func (intf *LLDPGlobalInfo) GetSystemDescription() string {
	return intf.RxInfo.RxLinkInfo.SysDescription
}

/*  Get Peer Description Information
 *
 */
func (intf *LLDPGlobalInfo) GetPortDescription() string {
	return intf.RxInfo.RxLinkInfo.PortDescription
}

/*  Get PVID Information
 *
 */
func (intf *LLDPGlobalInfo) GetPvid() int32 {
	return int32(intf.RxInfo.OrgTLV802_1.PVID)
}

/*  dump received lldp frame and other TX information
 */
func (intf LLDPGlobalInfo) DumpFrame() {
	intf.RxLock.RLock()
	defer intf.RxLock.RUnlock()
	debug.Logger.Debug("L2 Port:", intf.Port.IfIndex, "Port IfIndex:", intf.Port.IfIndex)
	debug.Logger.Debug("SrcMAC:", intf.RxInfo.SrcMAC.String(), "DstMAC:", intf.RxInfo.DstMAC.String())
	debug.Logger.Debug("ChassisID info is", intf.RxInfo.RxFrame.ChassisID)
	debug.Logger.Debug("PortID info is", intf.RxInfo.RxFrame.PortID)
	debug.Logger.Debug("TTL info is", intf.RxInfo.RxFrame.TTL)
	debug.Logger.Debug("Optional Values is", intf.RxInfo.RxLinkInfo)
}

/*  Api used to get entry.. This is mainly used by LLDP Server API Layer when it get config from
 *  North Bound Plugin...
 */
func (svr *LLDPServer) EntryExist(intfRef string) (int32, bool) {
	// first check whether the input is all numbers
	if ifIndex, err := strconv.Atoi(intfRef); err == nil {
		_, exists := svr.lldpGblInfo[int32(ifIndex)]
		if exists {
			return int32(ifIndex), exists
		}
	} else {
		// this is proper interface reference lets check the xRef
		ifIndex, exists := svr.lldpIntfRef2IfIndexMap[intfRef]
		if exists {
			return ifIndex, exists
		}
	}
	return -1, false
}

/*  Api to update system cache on next send frame
 */
func (svr *LLDPServer) UpdateCache(sysInfo *config.SystemInfo) {
	svr.SysInfo = sysInfo
	debug.Logger.Debug("Updated system information:", *svr.SysInfo)
	for _, ifIndex := range svr.lldpUpIntfStateSlice {
		intf, exists := svr.lldpGblInfo[ifIndex]
		if !exists {
			continue
		}
		intf.TxInfo.SetCache(false)
	}
}

/*  Api to get mode in uint8 format given a string
 */
func getMode(mode string) (rv uint8) {
	switch mode {
	case config.TX_RX_MODE_TxRx:
		rv = config.TXRX
	case config.TX_RX_MODE_TxOnly:
		rv = config.TX_ONLY
	case config.TX_RX_MODE_RxOnly:
		rv = config.RX_ONLY
	default:
		rv = config.INVALID_MODE
	}
	return rv
}

/*  Api to get mode in string format given a uint8
 */
func getModeString(mode uint8) (rv string) {
	switch mode {
	case config.TXRX:
		rv = config.TX_RX_MODE_TxRx
	case config.TX_ONLY:
		rv = config.TX_RX_MODE_TxOnly
	case config.RX_ONLY:
		rv = config.TX_RX_MODE_RxOnly
	default:
		rv = "invalid"
	}
	return rv
}
