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

package rpc

import (
	"dhcpd"
	"errors"
	"fmt"
	"l3/dhcp/server"
)

func (h *DHCPHandler) SendSetDhcpGlobalConfig(conf *dhcpd.DhcpGlobalConfig) error {
	if conf.DefaultLeaseTime > conf.MaxLeaseTime {
		err := errors.New("Invalid Config: Default Lease Time cannot be more than Max Lease Time")
		return err
	}
	dhcpGlobalConf := server.DhcpGlobalConfig{
		Enable:           conf.Enable,
		DefaultLeaseTime: uint32(conf.DefaultLeaseTime),
		MaxLeaseTime:     uint32(conf.MaxLeaseTime),
	}
	h.server.DhcpGlobalConfCh <- dhcpGlobalConf
	return nil
}

func (h *DHCPHandler) SendSetDhcpIntfConfig(conf *dhcpd.DhcpIntfConfig) error {
	h.logger.Info(fmt.Sprintln("conf:", conf))
	subnet, ret := convertIPStrToUint32(conf.Subnet)
	if !ret {
		err := errors.New("Invalid Subnet")
		return err
	}
	h.logger.Info(fmt.Sprintln("subnet:", subnet))
	subnetMask, ret := convertIPStrToUint32(conf.SubnetMask)
	if !ret {
		err := errors.New("Invalid Subnet Mask")
		return err
	}
	h.logger.Info(fmt.Sprintln("subnetMask:", subnetMask))

	lowerIPBound, higherIPBound, ret := parseIPRangeStr(conf.IPAddrRange)
	if !ret {
		err := errors.New("Invalid IP Address Range")
		return err
	}

	if lowerIPBound&subnetMask != subnet ||
		higherIPBound&subnetMask != subnet {
		err := errors.New("Invalid IP Address Range")
		return err
	}
	h.logger.Info(fmt.Sprintln("lowerIPBound:", lowerIPBound, "higherIPBound:", higherIPBound))

	bCastAddr, ret := convertIPStrToUint32(conf.BroadcastAddr)
	if !ret {
		err := errors.New("Invalid Broadcast IP Address")
		return err
	}
	if bCastAddr&subnetMask != subnet {
		err := errors.New("Invalid Broadcast Address because it is not in same subnet.")
		return err
	}
	h.logger.Info(fmt.Sprintln("bCastAddr:", bCastAddr))

	rtrAddr, ret := convertIPStrToUint32(conf.RouterAddr)
	if !ret {
		err := errors.New("Invalid Router IP Address")
		return err
	}
	if rtrAddr&subnetMask != subnet {
		err := errors.New("Invalid Router Address because it is not in same subnet.")
		return err
	}

	h.logger.Info(fmt.Sprintln("rtrAddr:", rtrAddr))
	dnsAddr, ret := convertIPStrToUint32(conf.DNSServerAddr)
	if !ret {
		err := errors.New("Invalid Domain Name Server Address")
		return err
	}
	h.logger.Info(fmt.Sprintln("DNSServerAddr:", dnsAddr))

	dhcpIntfConf := server.DhcpIntfConfig{
		IntfRef:       conf.IntfRef,
		Subnet:        subnet,
		SubnetMask:    subnetMask,
		LowerIPBound:  lowerIPBound,
		HigherIPBound: higherIPBound,
		BCastAddr:     bCastAddr,
		RtrAddr:       rtrAddr,
		DnsAddr:       dnsAddr,
		DomainName:    conf.DomainName,
		Enable:        conf.Enable,
	}
	h.server.DhcpIntfConfCh <- dhcpIntfConf
	retMsg := <-h.server.DhcpIntfConfRetCh
	return retMsg
}

func (h *DHCPHandler) CreateDhcpGlobalConfig(conf *dhcpd.DhcpGlobalConfig) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received CreateDhcpGlobalConfig:", conf))
	err := h.SendSetDhcpGlobalConfig(conf)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *DHCPHandler) CreateDhcpIntfConfig(conf *dhcpd.DhcpIntfConfig) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received CreateDhcpIntfConfig:", conf))
	err := h.SendSetDhcpIntfConfig(conf)
	if err != nil {
		return false, err
	}
	return true, nil
}
