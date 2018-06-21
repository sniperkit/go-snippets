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
	"errors"
	"fmt"
	"l3/vrrp/api"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"strings"
	"syscall"
	"vrrpd"
)

func (h *ConfigHandler) CreateVrrpGlobal(cfg *vrrpd.VrrpGlobal) (r bool, err error) {
	debug.Logger.Info("Thrift request for creating vrrp global object:", *cfg)
	gblCfg := &common.GlobalConfig{cfg.Vrf, cfg.Enable, common.CREATE}
	api.CreateVrrpGbl(gblCfg)
	debug.Logger.Info("Thrift returning for creating vrrp global object true, nil")
	return true, nil
}

func (h *ConfigHandler) UpdateVrrpGlobal(ocfg *vrrpd.VrrpGlobal, cfg *vrrpd.VrrpGlobal, attrset []bool, op []*vrrpd.PatchOpInfo) (r bool, err error) {
	debug.Logger.Info("Thrift request for updating vrrp global object:", *cfg)
	gblCfg := &common.GlobalConfig{cfg.Vrf, cfg.Enable, common.UPDATE}
	api.UpdateVrrpGbl(gblCfg)
	debug.Logger.Info("Thrift returning for updating vrrp global object true, nil")
	return true, nil
}

func (h *ConfigHandler) DeleteVrrpGlobal(cfg *vrrpd.VrrpGlobal) (r bool, err error) {
	debug.Logger.Info("Thrift request for deleting vrrp global object:", *cfg)
	err = errors.New("Deleting Vrrp Global Object is not Supported")
	r = false
	debug.Logger.Info("Thrift returning for deleting vrrp global object:", r, err)
	return r, err
}

func (h *ConfigHandler) CreateVrrpV4Intf(cfg *vrrpd.VrrpV4Intf) (r bool, err error) {
	debug.Logger.Info("Thrift request received for creating vrrp v4 interface common for:", *cfg)
	if !strings.Contains(cfg.Address, common.NETMASK_DELIMITER) {
		cfg.Address += common.NETMASK_DELIMITER + common.SLASH_32
	}
	v4Cfg := &common.IntfCfg{
		IntfRef:               cfg.IntfRef,
		VRID:                  cfg.VRID,
		Priority:              cfg.Priority,
		VirtualIPAddr:         cfg.Address,
		AdvertisementInterval: cfg.AdvertisementInterval,
		PreemptMode:           cfg.PreemptMode,
		AcceptMode:            cfg.AcceptMode,
		Operation:             common.CREATE,
		IpType:                syscall.AF_INET,
	}
	if cfg.AdminState == common.STATE_UP {
		v4Cfg.AdminState = true
	}
	if cfg.Version == common.VERSION2_STR {
		v4Cfg.Version = common.VERSION2
	} else if cfg.Version == common.VERSION3_STR {
		v4Cfg.Version = common.VERSION3
	}
	debug.Logger.Info("Push create cfg:", *v4Cfg, "to api layer")
	r, err = api.VrrpIntfConfig(v4Cfg)
	debug.Logger.Info("Thrift request returning for creating vrrp v4 interface config returning:", r, err)
	return r, err
}
func (h *ConfigHandler) UpdateVrrpV4Intf(origconfig *vrrpd.VrrpV4Intf, newconfig *vrrpd.VrrpV4Intf, attrset []bool, op []*vrrpd.PatchOpInfo) (r bool, err error) {
	debug.Logger.Info("Thrift request received for updating vrrp v4 interface config for:", *origconfig, "to new:", *newconfig)
	if !strings.Contains(newconfig.Address, common.NETMASK_DELIMITER) {
		newconfig.Address += common.NETMASK_DELIMITER + common.SLASH_32
	}
	if origconfig.Version != newconfig.Version {
		return false, errors.New(fmt.Sprintln("Changing Version from:", origconfig.Version, "to", newconfig.Version, "is not supported"))
	}
	v4Cfg := &common.IntfCfg{
		IntfRef:               newconfig.IntfRef,
		VRID:                  newconfig.VRID,
		Priority:              newconfig.Priority,
		VirtualIPAddr:         newconfig.Address,
		AdvertisementInterval: newconfig.AdvertisementInterval,
		PreemptMode:           newconfig.PreemptMode,
		AcceptMode:            newconfig.AcceptMode,
		Operation:             common.UPDATE,
		IpType:                syscall.AF_INET,
	}
	if newconfig.AdminState == common.STATE_UP {
		v4Cfg.AdminState = true
	}
	if newconfig.Version == common.VERSION2_STR {
		v4Cfg.Version = common.VERSION2
	} else if newconfig.Version == common.VERSION3_STR {
		v4Cfg.Version = common.VERSION3
	}
	debug.Logger.Info("Push update cfg:", *v4Cfg, "to api layer")
	r, err = api.VrrpIntfConfig(v4Cfg)
	debug.Logger.Info("Thrift request returning for updating vrrp v4 interface config for:", *origconfig, "to new:", *newconfig, "returning:", r, err)
	return true, nil
}

func (h *ConfigHandler) DeleteVrrpV4Intf(cfg *vrrpd.VrrpV4Intf) (r bool, err error) {
	debug.Logger.Info("Thrift request received for deleting vrrp v4 interface cfg for:", *cfg)
	if !strings.Contains(cfg.Address, common.NETMASK_DELIMITER) {
		cfg.Address += common.NETMASK_DELIMITER + common.SLASH_32
	}
	v4Cfg := &common.IntfCfg{
		IntfRef:               cfg.IntfRef,
		VRID:                  cfg.VRID,
		Priority:              cfg.Priority,
		VirtualIPAddr:         cfg.Address,
		AdvertisementInterval: cfg.AdvertisementInterval,
		PreemptMode:           cfg.PreemptMode,
		AcceptMode:            cfg.AcceptMode,
		Operation:             common.DELETE,
		IpType:                syscall.AF_INET,
	}
	if cfg.AdminState == common.STATE_UP {
		v4Cfg.AdminState = true
	}
	if cfg.Version == common.VERSION2_STR {
		v4Cfg.Version = common.VERSION2
	} else if cfg.Version == common.VERSION3_STR {
		v4Cfg.Version = common.VERSION3
	}
	debug.Logger.Info("Push delete cfg:", *v4Cfg, "to api layer")
	r, err = api.VrrpIntfConfig(v4Cfg)
	debug.Logger.Info("Thrift request returning for deleting vrrp v4 interface config returning:", r, err)
	return r, err
}

func (h *ConfigHandler) CreateVrrpV6Intf(cfg *vrrpd.VrrpV6Intf) (r bool, err error) {
	debug.Logger.Info("Thrift request received for creating vrrp v6 interface cfg for:", *cfg)
	if !strings.Contains(cfg.Address, common.NETMASK_DELIMITER) {
		cfg.Address += common.NETMASK_DELIMITER + common.SLASH_64
	}
	v6Cfg := &common.IntfCfg{
		IntfRef:               cfg.IntfRef,
		VRID:                  cfg.VRID,
		Priority:              cfg.Priority,
		VirtualIPAddr:         cfg.Address,
		AdvertisementInterval: cfg.AdvertisementInterval,
		PreemptMode:           cfg.PreemptMode,
		AcceptMode:            cfg.AcceptMode,
		Version:               common.VERSION3,
		Operation:             common.CREATE,
		IpType:                syscall.AF_INET6,
	}
	if cfg.AdminState == common.STATE_UP {
		v6Cfg.AdminState = true
	}
	r, err = api.VrrpIntfConfig(v6Cfg)
	debug.Logger.Info("Thrift request returning for creating vrrp v6 interface cfg returning:", r, err)
	return r, err
}

func (h *ConfigHandler) UpdateVrrpV6Intf(origconfig *vrrpd.VrrpV6Intf, newconfig *vrrpd.VrrpV6Intf, attrset []bool, op []*vrrpd.PatchOpInfo) (r bool, err error) {
	debug.Logger.Info("Thrift request received for updating vrrp v6 interface config for:", *origconfig, "to new:", *newconfig)
	if !strings.Contains(newconfig.Address, common.NETMASK_DELIMITER) {
		newconfig.Address += common.NETMASK_DELIMITER + common.SLASH_64
	}
	v6Cfg := &common.IntfCfg{
		IntfRef:               newconfig.IntfRef,
		VRID:                  newconfig.VRID,
		Priority:              newconfig.Priority,
		VirtualIPAddr:         newconfig.Address,
		AdvertisementInterval: newconfig.AdvertisementInterval,
		PreemptMode:           newconfig.PreemptMode,
		AcceptMode:            newconfig.AcceptMode,
		Version:               common.VERSION3,
		Operation:             common.UPDATE,
		IpType:                syscall.AF_INET6,
	}
	if newconfig.AdminState == common.STATE_UP {
		v6Cfg.AdminState = true
	}
	r, err = api.VrrpIntfConfig(v6Cfg)
	debug.Logger.Info("Thrift request returning for updating vrrp v6 interface config for:", *origconfig, "to new:", *newconfig, "returning:", r, err)
	return true, nil
}

func (h *ConfigHandler) DeleteVrrpV6Intf(cfg *vrrpd.VrrpV6Intf) (r bool, err error) {
	debug.Logger.Info("Thrift request received for deleting vrrp v6 interface cfg for:", *cfg)
	if !strings.Contains(cfg.Address, common.NETMASK_DELIMITER) {
		cfg.Address += common.NETMASK_DELIMITER + common.SLASH_64
	}
	v6Cfg := &common.IntfCfg{
		IntfRef:               cfg.IntfRef,
		VRID:                  cfg.VRID,
		Priority:              cfg.Priority,
		VirtualIPAddr:         cfg.Address,
		AdvertisementInterval: cfg.AdvertisementInterval,
		PreemptMode:           cfg.PreemptMode,
		AcceptMode:            cfg.AcceptMode,
		Version:               common.VERSION3,
		Operation:             common.DELETE,
		IpType:                syscall.AF_INET6,
	}
	if cfg.AdminState == common.STATE_UP {
		v6Cfg.AdminState = true
	}
	r, err = api.VrrpIntfConfig(v6Cfg)
	debug.Logger.Info("Thrift request returning for deleting vrrp v6 interface cfg returning:", r, err)
	return r, err
}
