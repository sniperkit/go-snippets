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

// interfaces.go
package utils

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"utils/logging"
)

type IPInfo struct {
	IpAddr          net.IP
	IPv6Addr        net.IP
	IpMask          net.IPMask
	IPv6Mask        net.IPMask
	LinklocalIpAddr string
}

func newIPInfo() *IPInfo {
	return &IPInfo{
		IpAddr:          nil,
		IpMask:          nil,
		IPv6Addr:        nil,
		IPv6Mask:        nil,
		LinklocalIpAddr: "",
	}
}

func (i *IPInfo) isEmpty() bool {
	if i.IpAddr == nil && i.IPv6Addr == nil && i.IpMask == nil && i.IPv6Mask == nil && i.LinklocalIpAddr == "" {
		return true
	}

	return false
}

type InterfaceMgr struct {
	logger      *logging.Writer
	rwMutex     *sync.RWMutex
	ifIndexToIP map[int32]*IPInfo //string
	ipToIfIndex map[string]int32
}

var ifaceMgr *InterfaceMgr

func NewInterfaceMgr(logger *logging.Writer) *InterfaceMgr {
	if ifaceMgr != nil {
		logger.Info("NewInterfaceMgr: Return the existing interface manager", ifaceMgr)
		return ifaceMgr
	}

	ifaceMgr = &InterfaceMgr{
		logger:      logger,
		rwMutex:     &sync.RWMutex{},
		ifIndexToIP: make(map[int32]*IPInfo),
		ipToIfIndex: make(map[string]int32),
	}
	logger.Info("NewInterfaceMgr: Creating new interface manager", ifaceMgr)
	return ifaceMgr
}

func (i *InterfaceMgr) IsIPConfigured(ip string) bool {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	i.logger.Info("IsIPConfigured: ip", ip, "ipToIfIndex", i.ipToIfIndex)
	_, ok := i.ipToIfIndex[ip]
	return ok
}

func (i *InterfaceMgr) GetIfaceIP(ifIndex int32) (ipInfo *IPInfo, err error) {
	var ok bool
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	i.logger.Info("GetIfaceIP: ifIndex", ifIndex, "ifIndexToIP", i.ifIndexToIP)
	if ipInfo, ok = i.ifIndexToIP[ifIndex]; !ok {
		err = errors.New(fmt.Sprintf("Iface %d is not configured", ifIndex))
	}

	return ipInfo, err
}

func (i *InterfaceMgr) GetIfaceIfIdx(ipAddr string) (idx int32, err error) {
	var ok bool
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	i.logger.Info("GetIfaceIdx: ipAddr", ipAddr, "ipAddrToIdx", i.ipToIfIndex)
	if idx, ok = i.ipToIfIndex[ipAddr]; !ok {
		err = errors.New(fmt.Sprintf("Iface %s is not configured", ipAddr))
	}

	return idx, err
}

func (i *InterfaceMgr) AddIface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("AddIface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex", i.ipToIfIndex)

	ip, ipMask, err := net.ParseCIDR(addr)
	if err != nil {
		i.logger.Err("AddIface: ParseCIDR failed for addr", addr, "with error", err)
		return
	}

	var ipInfo *IPInfo
	var ok bool

	if ipInfo, ok = i.ifIndexToIP[ifIndex]; !ok {
		ipInfo = newIPInfo()
		i.ifIndexToIP[ifIndex] = ipInfo
	}

	ipInfo.IpAddr = ip
	ipInfo.IpMask = ipMask.Mask

	i.ipToIfIndex[ip.String()] = ifIndex
}

func (i *InterfaceMgr) AddV6Iface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("AddV6Iface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex", i.ipToIfIndex)

	ip, ipMask, err := net.ParseCIDR(addr)
	if err != nil {
		i.logger.Err("AddV6Iface: ParseCIDR failed for addr", addr, "with error", err)
		return
	}

	var ipInfo *IPInfo
	var ok bool

	if ipInfo, ok = i.ifIndexToIP[ifIndex]; !ok {
		ipInfo = newIPInfo()
		i.ifIndexToIP[ifIndex] = ipInfo
	}

	ipInfo.IPv6Addr = ip
	ipInfo.IPv6Mask = ipMask.Mask

	i.ipToIfIndex[ip.String()] = ifIndex
}

func (i *InterfaceMgr) AddLinkLocalIface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("AddLinkLocalIface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex", i.ipToIfIndex)

	netIP := net.ParseIP(addr)
	if netIP == nil || !netIP.IsLinkLocalUnicast() {
		i.logger.Err("AddLinkLocalIface: ifIndex", ifIndex, "ip", addr, "is not a link local unicast address")
		return
	}

	var ipInfo *IPInfo
	var ok bool

	if ipInfo, ok = i.ifIndexToIP[ifIndex]; !ok {
		ipInfo = newIPInfo()
		i.ifIndexToIP[ifIndex] = ipInfo
	}

	if ipInfo.LinklocalIpAddr != "" {
		i.logger.Err("AddLinkLocalIface: ifIndex", ifIndex, "ip", addr, "link local ip", ipInfo.LinklocalIpAddr,
			"is already set on the interface")
		return
	}

	ipInfo.LinklocalIpAddr = addr
	i.ipToIfIndex[addr] = ifIndex
}

func (i *InterfaceMgr) RemoveIface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("RemoveIface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex",
		i.ipToIfIndex)

	if ipInfo, ok := i.ifIndexToIP[ifIndex]; ok {
		delete(i.ipToIfIndex, ipInfo.IpAddr.String())
		ipInfo.IpAddr = nil
		ipInfo.IpMask = nil

		if ipInfo.isEmpty() {
			delete(i.ifIndexToIP, ifIndex)
		}
	}
}

func (i *InterfaceMgr) RemoveV6Iface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("RemoveV6Iface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex",
		i.ipToIfIndex)

	if ipInfo, ok := i.ifIndexToIP[ifIndex]; ok {
		delete(i.ipToIfIndex, ipInfo.IPv6Addr.String())
		ipInfo.IPv6Addr = nil
		ipInfo.IPv6Mask = nil

		if ipInfo.isEmpty() {
			delete(i.ifIndexToIP, ifIndex)
		}
	}
}

func (i *InterfaceMgr) RemoveLinkLocalIface(ifIndex int32, addr string) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.logger.Info("RemoveLinkLocalIface: ifIndex", ifIndex, "ip", addr, "ifIndexToIP", i.ifIndexToIP, "ipToIfIndex",
		i.ipToIfIndex)

	netIP := net.ParseIP(addr)
	if netIP == nil || !netIP.IsLinkLocalUnicast() {
		i.logger.Err("RemoveLinkLocalIface: ifIndex", ifIndex, "ip", addr, "is not a link local unicast address")
		return
	}

	if ipInfo, ok := i.ifIndexToIP[ifIndex]; ok {
		delete(i.ipToIfIndex, ipInfo.LinklocalIpAddr)
		ipInfo.LinklocalIpAddr = ""

		if ipInfo.isEmpty() {
			delete(i.ifIndexToIP, ifIndex)
		}
	}
}
