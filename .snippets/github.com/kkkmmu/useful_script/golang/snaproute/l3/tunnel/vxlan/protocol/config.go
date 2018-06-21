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

// config.go
// Config entry is based on thrift data structures.
package vxlan

import (
	"fmt"
	"net"
	//"strings"
	"errors"
	"vxland"
)

const (
	VxlanCommandCreate = iota + 1
	VxlanCommandDelete
	VxlanCommandUpdate
)

type cfgFileJson struct {
	SwitchMac        string            `json:"SwitchMac"`
	PluginList       []string          `json:"PluginList"`
	IfNameMap        map[string]string `json:"IfNameMap"`
	IfNamePrefix     map[string]string `json:"IfNamePrefix"`
	SysRsvdVlanRange string            `json:"SysRsvdVlanRange"`
}

type VxLanConfigChannels struct {
	Vxlancreate               chan VxlanConfig
	Vxlandelete               chan VxlanConfig
	Vxlanupdate               chan VxlanUpdate
	Vtepcreate                chan VtepConfig
	Vtepdelete                chan VtepConfig
	Vtepupdate                chan VtepUpdate
	VxlanAccessPortVlanUpdate chan VxlanAccessPortVlan
	VxlanNextHopUpdate        chan VxlanNextHopIp
	VxlanPortCreate           chan PortConfig
	Vxlanintfinfo             chan VxlanIntfInfo
}

type VxlanIntfInfo struct {
	Command  int
	IntfName string
	IfIndex  int32
	Mac      net.HardwareAddr
	Ip       net.IP
}

type VxlanNextHopIp struct {
	Command   int
	Ip        net.IP
	Intf      int32
	IntfName  string
	NextHopIp net.IP
}

type VxlanAccessPortVlan struct {
	Command  int
	VlanId   uint16
	IntfList []int32
}

type VxlanUpdate struct {
	Oldconfig VxlanConfig
	Newconfig VxlanConfig
	Attr      []string
}

type VtepUpdate struct {
	Oldconfig VtepConfig
	Newconfig VtepConfig
	Attr      []string
}

// bridges attached to the VNI, mapping table to know
// what vlan maps to what VNI used for filtering packets on RX
type VxlanConfig struct {
	VNI          uint32
	UntaggedVlan []uint16
	VlanId       []uint16 // used to tag inner ethernet frame when egressing
	Enable       bool
}

type PortConfig struct {
	Name         string
	HardwareAddr net.HardwareAddr
	Speed        int32
	PortNum      int32
	IfIndex      int32
}

// tunnel endpoint for the VxLAN
type VtepConfig struct {
	Vni                   uint32           `SNAPROUTE: KEY` //VxLAN ID.
	VtepName              string           //VTEP instance name.
	SrcIfName             string           //Source interface ifIndex.
	UDP                   uint16           //vxlan udp port.  Deafult is the iana default udp port
	TTL                   uint16           //TTL of the Vxlan tunnel
	TOS                   uint16           //Type of Service
	MTU                   uint32           //Maximum Transmission Unit
	InnerVlanHandlingMode int32            //The inner vlan tag handling mode.
	Learning              bool             //specifies if unknown source link layer  addresses and IP addresses are entered into the VXLAN  device forwarding database.
	Rsc                   bool             //specifies if route short circuit is turned on.
	L2miss                bool             //specifies if netlink LLADDR miss notifications are generated.
	L3miss                bool             //specifies if netlink IP ADDR miss notifications are generated.
	TunnelSrcIp           net.IP           //Source IP address for the static VxLAN tunnel
	TunnelDstIp           net.IP           //Destination IP address for the static VxLAN tunnel
	VlanId                uint16           //Vlan Id to encapsulate with the vtep tunnel ethernet header
	TunnelSrcMac          net.HardwareAddr //Src Mac assigned to the VTEP within this VxLAN. If an address is not assigned the the local switch address will be used.
	TunnelDstMac          net.HardwareAddr // Optional - may be looked up based on TunnelNextHopIp
	TunnelNextHopIP       net.IP           // NextHopIP is used to find the DMAC for the tunnel within Asicd
	Enable                bool
}

func ConvertInt32ToBool(val int32) bool {
	if val == 0 {
		return false
	}
	return true
}

// VxlanConfigCheck
// Validate the VXLAN provisioning
func VxlanConfigCheck(c *VxlanConfig) error {
	if GetVxlanDBEntry(c.VNI) != nil {
		return errors.New(fmt.Sprintln("Error VxlanInstance Exists vni is not unique", c))
	}
	if len(c.VlanId) > 1 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Only support one VNI->VLAN map", c))
	}
	if len(c.UntaggedVlan) > 1 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Only support one VNI->VLAN map", c))
	}
	if len(c.UntaggedVlan) > 0 && len(c.VlanId) > 0 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Only support one VNI->VLAN map", c))
	}

	// could optimize by making this call once on startup and listen for vlan create
	// updates from asicd but for now this is the simple call (will only affect speed
	// of provisioning or if asicd is not available)
	var vlanList []uint16
	for _, client := range ClientIntf {
		vlanList = client.GetAllVlans()
	}

	unprovvlanlist := make([]uint16, 0)
	for _, vlan := range c.VlanId {
		foundVlan := false

		for _, provvlan := range vlanList {
			if vlan == provvlan {
				foundVlan = true
			}
		}
		if !foundVlan {
			unprovvlanlist = append(unprovvlanlist, vlan)
		}
	}
	if len(unprovvlanlist) > 0 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Prov failed because following vlans were not provisioned %#v", unprovvlanlist))
	}

	return nil
}

func VxlanConfigUpdateCheck(oc *VxlanConfig, nc *VxlanConfig) error {
	if oc.VNI != nc.VNI {
		return errors.New(fmt.Sprintln("Error Unsupported Attribute VNI Update, must delete then create"))
	}
	vxlan := GetVxlanDBEntry(nc.VNI)
	if vxlan == nil {
		return errors.New(fmt.Sprintln("Error Error VxlanInstance Does not exists"))
	}

	if len(nc.VlanId) > 1 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Update Only support one VNI->VLAN map", nc))
	}
	if len(nc.UntaggedVlan) > 1 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Update Only support one VNI->VLAN map", nc))
	}
	if len(nc.UntaggedVlan) > 0 && len(nc.VlanId) > 0 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Update Only support one VNI->VLAN map", nc))
	}

	// could optimize by making this call once on startup and listen for vlan create
	// updates from asicd but for now this is the simple call (will only affect speed
	// of provisioning or if asicd is not available)
	var vlanList []uint16
	for _, client := range ClientIntf {
		vlanList = client.GetAllVlans()
	}

	unprovvlanlist := make([]uint16, 0)
	for _, vlan := range nc.VlanId {
		foundVlan := false
		for _, provvlan := range vlanList {
			if vlan == provvlan {
				foundVlan = true
				break
			}
		}
		if !foundVlan {
			unprovvlanlist = append(unprovvlanlist, vlan)
		}
	}
	if len(unprovvlanlist) > 0 {
		return errors.New(fmt.Sprintln("Error VxlanInstance Prov failed because following vlans were not provisioned %#v", unprovvlanlist))
	}

	return nil
}

// VtepConfigCheck
// Validate the VTEP provisioning
func VtepConfigCheck(c *VtepConfig, create bool) error {
	key := VtepDbKey{
		Name:  c.VtepName,
		Vni:   c.Vni,
		DstIp: c.TunnelDstIp.String(),
	}
	if GetVtepDBEntry(&key) != nil && create {
		return errors.New(fmt.Sprintln("Error VtepInstance Exists name is not unique", c))
	}

	if c.MTU < 64 || c.MTU > 65535 {
		return errors.New("Error VtepInstance MTU must be between 64-65535")
	}

	if c.TTL > 255 {
		return errors.New("Error VtepInstance TTL must be between 0-255")
	}

	if c.TOS > 255 {
		return errors.New("Error VtepInstance TOS field must be between 0-255")
	}

	return nil
}

// ConvertVxlanInstanceToVxlanConfig:
// Convert thrift struct to vxlan config and check that db entry exists already
func ConvertVxlanInstanceToVxlanConfig(c *vxland.VxlanInstance, create bool) (*VxlanConfig, error) {

	if !create &&
		GetVxlanDBEntry(uint32(c.Vni)) == nil {
		return nil, errors.New(fmt.Sprintln("Error VxlanInstance does not Exists", c))

	}

	if c.AdminState != "UP" &&
		c.AdminState != "DOWN" {
		return nil, errors.New(fmt.Sprintln("Error VxlanInstance Error unsupported Admin State String must be UP/DOWN"))
	}

	vxlan := &VxlanConfig{

		VNI: uint32(c.Vni),
	}
	for _, untagvlan := range c.UntaggedVlanId {
		vxlan.UntaggedVlan = append(vxlan.UntaggedVlan, uint16(untagvlan))
	}

	for _, vlan := range c.VlanId {
		vxlan.VlanId = append(vxlan.VlanId, uint16(vlan))
	}
	vxlan.Enable = false
	if c.AdminState == "UP" {
		vxlan.Enable = true
	}

	return vxlan, nil
}

func getVtepName(intf string) string {
	vtepName := intf
	//if !strings.Contains("vtep", intf) {
	//	vtepName = "vtep" + intf
	//}
	return vtepName
}

// ConvertVxlanVtepInstanceToVtepConfig:
// Convert thrift struct to vxlan config
func ConvertVxlanVtepInstanceToVtepConfig(c *vxland.VxlanVtepInstance) ([]*VtepConfig, error) {

	var mac net.HardwareAddr
	var ip net.IP
	var name string
	//var ok bool
	vtepName := getVtepName(c.Intf)
	name = c.IntfRef
	ip = net.ParseIP(c.SrcIp)

	if c.SrcIp != "" && ip == nil {
		return nil, errors.New(fmt.Sprintln("Error VxlanVtepInstance unsupported SrcIp Format"))
	}

	for _, dstip := range c.DstIp {
		dip := net.ParseIP(dstip)
		if dip == nil {
			return nil, errors.New(fmt.Sprintln("Error VxlanVtepInstance destination Ip invalid", dip))
		}
	}

	if c.SrcIp == "" && c.IntfRef == "" {
		return nil, errors.New(fmt.Sprintln("Error VxlanVtepInstance SrcIp or IntfRef required"))
	}

	/* TODO need to create a generic way to get an interface name, mac, ip
	if c.SrcIp == "0.0.0.0" && c.IntfRef != "" {
		// need to get the appropriate IntfRef type
		ok, name, mac, ip = snapclient.asicDGetLoopbackInfo()
		if !ok {
			errorstr := "VTEP: Src Tunnel Info not provisioned yet, loopback intf needed"
			logger.Info(errorstr)
			return &VtepConfig{}, errors.New(errorstr)
		}
		fmt.Println("loopback info:", name, mac, ip)
		if c.SrcIp != "0.0.0.0" {
			ip = net.ParseIP(c.SrcIp)
		}
		logger.Info(fmt.Sprintf("Forcing Vtep %s to use Lb %s SrcMac %s Ip %s", vtepName, name, mac, ip))
	}
	*/
	if c.AdminState != "UP" &&
		c.AdminState != "DOWN" {
		return nil, errors.New(fmt.Sprintln("Error VxlanInstance Error unsupported Admin State String must be UP/DOWN"))
	}
	enable := false
	if c.AdminState == "UP" {
		enable = true
	}

	vteps := make([]*VtepConfig, 0)

	for _, dstip := range c.DstIp {
		dip := net.ParseIP(dstip)
		vtep := &VtepConfig{
			Enable:    enable,
			Vni:       uint32(c.Vni),
			VtepName:  vtepName,
			SrcIfName: name,
			MTU:       uint32(c.Mtu),
			UDP:       uint16(c.DstUDP),
			TTL:       uint16(c.TTL),
			TOS:       uint16(c.TOS),
			InnerVlanHandlingMode: c.InnerVlanHandlingMode,
			TunnelSrcIp:           ip,
			TunnelDstIp:           dip,
			VlanId:                uint16(c.VlanId),
			TunnelSrcMac:          mac,
		}
		vteps = append(vteps, vtep)
	}

	return vteps, nil
}

func UpdateThriftVxLAN(c *VxlanUpdate) {
	// important to note that the attrset starts at index 0 which is the BaseObj
	// which is not the first element on the thrift obj, thus we need to skip
	// this attribute
	updateCfg := false
	disableCfg := false
	enableCfg := false
	for _, objName := range c.Attr {
		if objName == "VlanId" ||
			objName == "UntaggedVlanId" {
			logger.Info(fmt.Sprintf("Updating vlan list tag: %#v,", c.Newconfig.VlanId))
			updateCfg = true
		} else if objName == "AdminState" {
			logger.Info(fmt.Sprintf("Updating admin state %t", c.Newconfig.Enable))
			if c.Newconfig.Enable {
				enableCfg = true
			} else {
				disableCfg = true
			}
		}
	}

	if disableCfg {
		DeleteVxLAN(&c.Oldconfig, true)
		saveVxLanConfigData(&c.Newconfig)
	} else if enableCfg {
		saveVxLanConfigData(&c.Newconfig)
		CreateVxLAN(&c.Newconfig)
	} else if updateCfg {
		newvlanlist := make([]uint16, 0)
		delvlanlist := make([]uint16, 0)
		newuntagvlanlist := make([]uint16, 0)
		deluntagvlanlist := make([]uint16, 0)

		vxlan, ok := GetVxlanDB()[c.Newconfig.VNI]
		if ok {
			for idx, nvlan := range c.Newconfig.VlanId {
				foundvlan := false
				for _, ovlan := range c.Oldconfig.VlanId {
					if nvlan == ovlan {
						foundvlan = true
						break
					}
				}
				if !foundvlan {
					newvlanlist = append(newvlanlist, c.Newconfig.VlanId[idx])
				}
			}
			for idx, nvlan := range c.Oldconfig.VlanId {
				foundvlan := false
				for _, ovlan := range c.Newconfig.VlanId {
					if nvlan == ovlan {
						foundvlan = true
						break
					}
				}
				if !foundvlan {
					delvlanlist = append(delvlanlist, c.Oldconfig.VlanId[idx])
				}
			}

			for idx, nvlan := range c.Newconfig.UntaggedVlan {
				foundvlan := false
				for _, ovlan := range c.Oldconfig.UntaggedVlan {
					if nvlan == ovlan {
						foundvlan = true
						break
					}
				}
				if !foundvlan {
					newuntagvlanlist = append(newuntagvlanlist, c.Newconfig.UntaggedVlan[idx])
				}
			}
			for idx, nvlan := range c.Oldconfig.UntaggedVlan {
				foundvlan := false
				for _, ovlan := range c.Newconfig.UntaggedVlan {
					if nvlan == ovlan {
						foundvlan = true
						break
					}
				}
				if !foundvlan {
					deluntagvlanlist = append(deluntagvlanlist, c.Oldconfig.UntaggedVlan[idx])
				}
			}

			// TODO add lock around db as it is used as a filtering
			if len(delvlanlist) > 0 {
				for _, dv := range delvlanlist {
					for idx, vlan := range vxlan.VlanId {
						if dv == vlan {
							vxlan.VlanId = append(vxlan.VlanId[:idx], vxlan.VlanId[idx+1:]...)
							break
						}
					}
				}
			}

			if len(newvlanlist) > 0 {
				for _, nv := range newvlanlist {
					vxlan.VlanId = append(vxlan.VlanId, nv)
				}
			}

			for _, client := range ClientIntf {
				client.UpdateVxlan(vxlan.VNI, newvlanlist, delvlanlist, newuntagvlanlist, deluntagvlanlist)
			}
		}
	}
}

func UpdateThriftVtep(c *VtepUpdate) {

	updateobj := false
	enableobj := false
	disableobj := false
	recreateobj := false
	for _, objName := range c.Attr {
		if objName == "TOS" ||
			objName == "Mtu" ||
			objName == "TTL" {
			logger.Debug("UpdateVtep ", objName)
			updateobj = true
		} else if objName == "AdminState" {
			if c.Newconfig.Enable {
				logger.Debug("UpdateVtep ", objName, "Enabled")
				enableobj = true
			} else {
				logger.Debug("UpdateVtep", objName, "Disabled")
				disableobj = true
			}
		} else {
			logger.Debug("UpdateVtep %s Service Affecting Change", objName)
			recreateobj = true
		}
	}

	if disableobj {
		key := &VtepDbKey{
			Name:  c.Newconfig.VtepName,
			Vni:   c.Newconfig.Vni,
			DstIp: c.Newconfig.TunnelDstIp.String(),
		}
		vtep := GetVtepDBEntry(key)
		if vtep != nil {
			vtep.VxlanVtepMachineFsm.VxlanVtepEvents <- MachineEvent{
				E:   VxlanVtepEventDisable,
				Src: VxlanVtepMachineModuleStr,
			}
			saveVtepConfigData(&(c.Newconfig))
		}
	} else if enableobj {
		key := &VtepDbKey{
			Name:  c.Newconfig.VtepName,
			Vni:   c.Newconfig.Vni,
			DstIp: c.Newconfig.TunnelDstIp.String(),
		}
		vtep := GetVtepDBEntry(key)
		if vtep != nil {
			saveVtepConfigData(&(c.Newconfig))
			ReProvisionVtep(vtep)
		}
	} else if recreateobj {
		DeleteVtep(&(c.Oldconfig))
		CreateVtep(&(c.Newconfig))
	} else if updateobj {
		for _, attr := range c.Attr {
			if attr == "TOS" {
				UpdateVtepTOS(c.Newconfig.VtepName, uint32(c.Newconfig.Vni), c.Newconfig.TunnelDstIp.String(), uint8(c.Newconfig.TOS))
			} else if attr == "Mtu" {
				UpdateVtepMTU(c.Newconfig.VtepName, uint32(c.Newconfig.Vni), c.Newconfig.TunnelDstIp.String(), uint16(c.Newconfig.MTU))
			} else if attr == "TTL" {
				UpdateVtepTTL(c.Newconfig.VtepName, uint32(c.Newconfig.Vni), c.Newconfig.TunnelDstIp.String(), uint8(c.Newconfig.TTL))
			}
		}
	}
}

func UpdateVtepTOS(vtepName string, vni uint32, dstip string, tos uint8) {
	key := &VtepDbKey{
		Name:  vtepName,
		Vni:   vni,
		DstIp: dstip,
	}
	vtep := GetVtepDBEntry(key)
	if vtep != nil {
		// only need to set the tos attribute as vxland constructs the tos
		vtep.TOS = tos
		vtepDB[*key] = vtep
		for idx, v := range vtepDbList {
			if vtep.VtepConfigName == v.VtepConfigName &&
				vtep.Vni == v.Vni &&
				vtep.DstIp.String() == v.DstIp.String() {
				vtepDbList = append(vtepDbList[:idx], vtepDbList[idx+1:]...)
				break
			}
		}
		vtepDbList = append(vtepDbList, vtep)
		for _, client := range ClientIntf {
			client.UpdateVtepAttr(vtep.VtepName, vtep.Vni, tos, vtep.TOS, vtep.MTU)
		}
	}
}

func UpdateVtepTTL(vtepName string, vni uint32, dstip string, ttl uint8) {
	key := &VtepDbKey{
		Name:  vtepName,
		Vni:   vni,
		DstIp: dstip,
	}
	vtep := GetVtepDBEntry(key)
	if vtep != nil {
		// only need to set the tos attribute as vxland constructs the tos
		vtep.TTL = ttl
		vtepDB[*key] = vtep
		for idx, v := range vtepDbList {
			if vtep.VtepConfigName == v.VtepConfigName &&
				vtep.Vni == v.Vni &&
				vtep.DstIp.String() == v.DstIp.String() {
				vtepDbList = append(vtepDbList[:idx], vtepDbList[idx+1:]...)
				break
			}
		}
		vtepDbList = append(vtepDbList, vtep)
		for _, client := range ClientIntf {
			client.UpdateVtepAttr(vtep.VtepName, vtep.Vni, vtep.TOS, ttl, vtep.MTU)
		}
	}
}

func UpdateVtepMTU(vtepName string, vni uint32, dstip string, mtu uint16) {
	key := &VtepDbKey{
		Name:  vtepName,
		Vni:   vni,
		DstIp: dstip,
	}
	vtep := GetVtepDBEntry(key)
	if vtep != nil {
		// only need to set the tos attribute as vxland constructs the tos
		vtep.MTU = mtu
		vtepDB[*key] = vtep
		for idx, v := range vtepDbList {
			if vtep.VtepConfigName == v.VtepConfigName &&
				vtep.Vni == v.Vni &&
				vtep.DstIp.String() == v.DstIp.String() {
				vtepDbList = append(vtepDbList[:idx], vtepDbList[idx+1:]...)
				break
			}
		}
		vtepDbList = append(vtepDbList, vtep)
		for _, client := range ClientIntf {
			client.UpdateVtepAttr(vtep.VtepName, vtep.Vni, vtep.TOS, vtep.TTL, mtu)
		}
	}
}

//HandleNextHopChange:
// Handle notifications from RIB that the next hop reachabilty has changed
func HandleNextHopChange(dip net.IP, nexthopip net.IP, nexthopIfIndex int32, nexthopIfName string, reachable bool) {
	// TOOD do some work to find all VTEP's and deprovision the entries
	for _, vtep := range GetVtepDB() {
		if reachable &&
			vtep.DstIp.String() == dip.String() {
			if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() == VxlanVtepStateInterface {

				nexthopinfo := VtepNextHopInfo{
					Ip:      nexthopip,
					IfIndex: nexthopIfIndex,
					IfName:  nexthopIfName,
				}
				event := MachineEvent{
					E:    VxlanVtepEventNextHopInfoResolved,
					Src:  VxlanVtepMachineModuleStr,
					Data: nexthopinfo,
				}

				vtep.VxlanVtepMachineFsm.VxlanVtepEvents <- event

			}
		} else if !reachable &&
			vtep.DstIp.String() == dip.String() {
			// set state
			// tearing down the connection appropriately
			if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() >= VxlanVtepStateNextHopInfo {
				// deprovision the vtep
				DeProvisionVtep(vtep, false)
			}
		}
	}
}
