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

// vxlan_linux.go
// NOTE: this is meant for testing, it should eventually live in asicd
package test_linux

import (
	"fmt"
	"net"
	"time"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	//"os/exec"
	"utils/logging"
)

// options "proxy", "linux"
// TODO eventually read this from config file
var VxlanConfigMode string = "proxy"
var VxlanDB map[uint32]VxlanDbEntry
var macDb map[int32][]*VxlanMacDbEntry

type VxlanMacDbEntry struct {
	vtepName string
	mac      string
}

type VxlanDbEntry struct {
	VNI    uint32
	VlanId uint16 // used to tag inner ethernet frame when egressing
	Group  net.IP // multicast group IP
	MTU    uint32 // MTU size for each VTEP
	Brg    *netlink.Bridge
	Links  []*netlink.Link
}

type VxlanLinux struct {
	logger *logging.Writer
}

// bridge for the VNI
type VxlanConfig struct {
	VNI    uint32
	VlanId uint16 // used to tag inner ethernet frame when egressing
	Group  net.IP // multicast group IP
	MTU    uint32 // MTU size for each VTEP
}

// tunnel endpoint for the VxLAN
type VtepConfig struct {
	Vni                   uint32           `SNAPROUTE: KEY` //VxLAN ID.
	VtepName              string           //VTEP instance name.
	SrcIfName             string           //Source interface ifIndex.
	UDP                   uint16           //vxlan udp port.  Deafult is the iana default udp port
	TTL                   uint16           //TTL of the Vxlan tunnel
	TOS                   uint16           //Type of Service
	InnerVlanHandlingMode int32            //The inner vlan tag handling mode.
	Learning              bool             //specifies if unknown source link layer  addresses and IP addresses are entered into the VXLAN  device forwarding database.
	Rsc                   bool             //specifies if route short circuit is turned on.
	L2miss                bool             //specifies if netlink LLADDR miss notifications are generated.
	L3miss                bool             //specifies if netlink IP ADDR miss notifications are generated.
	TunnelSrcIp           net.IP           //Source IP address for the static VxLAN tunnel
	TunnelDstIp           net.IP           //Destination IP address for the static VxLAN tunnel
	VlanId                uint16           //Vlan Id to encapsulate with the vtep tunnel ethernet header
	TunnelSrcMac          net.HardwareAddr //Src Mac assigned to the VTEP within this VxLAN. If an address is not assigned the the local switch address will be used.
	TunnelDstMac          net.HardwareAddr
}

func NewVxlanLinux(logger *logging.Writer) *VxlanLinux {
	initVxlanDB()
	return &VxlanLinux{
		logger: logger,
	}

}

func initVxlanDB() {
	if VxlanDB == nil {
		VxlanDB = make(map[uint32]VxlanDbEntry)
	}
}

// createVxLAN is the equivalent to creating a bridge in the linux
// The VNI is actually associated with the VTEP so lets just create a bridge
// if necessary
func (v *VxlanLinux) CreateVxLAN(c *VxlanDbEntry) {

	if _, ok := VxlanDB[c.VNI]; !ok {
		VxlanDB[c.VNI] = VxlanDbEntry{
			VNI:    c.VNI,
			VlanId: c.VlanId,
			Group:  c.Group,
			MTU:    c.MTU,
			Links:  make([]*netlink.Link, 0),
		}
		// lets create a bridge if it does not exists
		// bridge should be based on the VLAN used by a
		// customer.
		brname := fmt.Sprintf("br%d", c.VNI)
		bridge := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: brname,
				MTU:  int(c.MTU),
			},
		}

		if err := netlink.LinkAdd(bridge); err != nil {
			v.logger.Err(err.Error())
		}

		link, err := netlink.LinkByName(bridge.Attrs().Name)
		if err != nil {
			v.logger.Err(err.Error())
		}

		vxlanDbEntry := VxlanDB[c.VNI]
		vxlanDbEntry.Brg = link.(*netlink.Bridge)
		VxlanDB[c.VNI] = vxlanDbEntry
		// lets set the vtep interface to up
		if err := netlink.LinkSetUp(bridge); err != nil {
			v.logger.Err(err.Error())
		}
	}
}

func (v *VxlanLinux) DeleteVxLAN(c *VxlanDbEntry) {

	if vxlan, ok := VxlanDB[c.VNI]; ok {
		for i, link := range vxlan.Links {
			// lets set the vtep interface to up
			if err := netlink.LinkSetDown(*link); err != nil {
				v.logger.Err(err.Error())
			}
			if err := netlink.LinkDel(*link); err != nil {
				v.logger.Err(err.Error())
			}

			vxlanDbEntry := VxlanDB[c.VNI]
			vxlanDbEntry.Links = append(vxlanDbEntry.Links[:i], vxlanDbEntry.Links[i+1:]...)
			VxlanDB[c.VNI] = vxlanDbEntry
		}

		link, err := netlink.LinkByName(vxlan.Brg.Name)
		if err != nil {
			v.logger.Err(err.Error())
		}

		// lets set the vtep interface to up
		if err := netlink.LinkSetDown(link); err != nil {
			v.logger.Err(err.Error())
		}
		if err := netlink.LinkDel(link); err != nil {
			v.logger.Err(err.Error())
		}

		delete(VxlanDB, c.VNI)
	}
}

func (v *VxlanLinux) CreateVtep(c *VtepConfig) {

	link, err := netlink.LinkByName(c.SrcIfName)
	if err != nil {
		v.logger.Err(fmt.Sprintf("Error finding link %s: %s", c.SrcIfName, err.Error()))
		return
	}

	fmt.Println("Config Mode", VxlanConfigMode)
	if VxlanConfigMode == "linux" {
		/* 4/6/16 DID Not work, packets were never received on VTEP */
		vtep := &netlink.Vxlan{
			LinkAttrs: netlink.LinkAttrs{
				Name: c.VtepName,
				//MasterIndex: VxlanDB[c.Vni].Brg.Attrs().Index,
				HardwareAddr: c.TunnelSrcMac,
				//MTU:          VxlanDB[c.Vni].Brg.Attrs().MTU,
				MTU: 1550,
			},
			VxlanId:      int(c.Vni),
			VtepDevIndex: link.Attrs().Index,
			SrcAddr:      c.TunnelSrcIp,
			Group:        c.TunnelDstIp,
			TTL:          int(c.TTL),
			TOS:          int(c.TOS),
			Learning:     c.Learning,
			Proxy:        false,
			RSC:          c.Rsc,
			L2miss:       false,
			L3miss:       false,
			UDPCSum:      true,
			NoAge:        false,
			GBP:          false,
			Age:          300,
			Port:         int(nl.Swap16(c.UDP)),
			PortLow:      int(c.UDP),
			PortHigh:     int(c.UDP),
		}
		//equivalent to linux command:
		// ip link add DEVICE type vxlan id ID [ dev PHYS_DEV  ] [ { group
		//         | remote } IPADDR ] [ local IPADDR ] [ ttl TTL ] [ tos TOS ] [
		//          port MIN MAX ] [ [no]learning ] [ [no]proxy ] [ [no]rsc ] [
		//          [no]l2miss ] [ [no]l3miss ]
		if err := netlink.LinkAdd(vtep); err != nil {
			v.logger.Err(err.Error())
		}

	} else {

		// Veth will create two interfaces
		// VtepName and VtepName + Int
		// the VtepNam + Int interface will be used by Vxland to rx packets
		// from other daemons and to send packets received from physical port
		// to the daemons
		//
		//  physical port <--> vxland (if vxlan packet) <--> vtepName Int <-->
		//  vtepName <--> Other Daemons listening
		//  on this vtepName interface
		vtep := &netlink.Veth{
			LinkAttrs: netlink.LinkAttrs{
				Name:         c.VtepName,
				MasterIndex:  VxlanDB[c.Vni].Brg.Attrs().Index,
				HardwareAddr: c.TunnelSrcMac,
				MTU:          VxlanDB[c.Vni].Brg.Attrs().MTU,
			},
			PeerName: c.VtepName + "Int",
		}
		//equivalent to linux command:
		// ip link add DEVICE type vxlan id ID [ dev PHYS_DEV  ] [ { group
		//         | remote } IPADDR ] [ local IPADDR ] [ ttl TTL ] [ tos TOS ] [
		//          port MIN MAX ] [ [no]learning ] [ [no]proxy ] [ [no]rsc ] [
		//          [no]l2miss ] [ [no]l3miss ]
		if err := netlink.LinkAdd(vtep); err != nil {
			v.logger.Err(err.Error())
		}

	}

	link, err = netlink.LinkByName(c.VtepName)
	if err != nil {
		v.logger.Err(fmt.Sprintf("Link by Name vtep:", err.Error()))
	}

	// found that hte mac we are trying to set fails lets try and add it again
	if err := netlink.LinkSetHardwareAddr(link, c.TunnelSrcMac); err != err {
		v.logger.Err(fmt.Sprintf("LinkSetHardwareAddr vtep:", err.Error()))
	}

	// equivalent to linux command:
	/* bridge fdb add - add a new fdb entry
	       This command creates a new fdb entry.

	       LLADDR the Ethernet MAC address.

	       dev DEV
	              the interface to which this address is associated.

	              self - the address is associated with a software fdb (default)

	              embedded - the address is associated with an offloaded fdb

	              router - the destination address is associated with a router.
	              Valid if the referenced device is a VXLAN type device and has
	              route shortcircuit enabled.

	      The next command line parameters apply only when the specified device
	      DEV is of type VXLAN.

	       dst IPADDR
	              the IP address of the destination VXLAN tunnel endpoint where
	              the Ethernet MAC ADDRESS resides.

	       vni VNI
	              the VXLAN VNI Network Identifier (or VXLAN Segment ID) to use to
	              connect to the remote VXLAN tunnel endpoint.  If omitted the
	              value specified at vxlan device creation will be used.

	       port PORT
	              the UDP destination PORT number to use to connect to the remote
	              VXLAN tunnel endpoint.  If omitted the default value is used.

	       via DEVICE
	              device name of the outgoing interface for the VXLAN device
	              driver to reach the remote VXLAN tunnel endpoint.


			// values taken from linux/neighbour.h
	*/
	if VxlanConfigMode == "linux" {
		if c.TunnelDstIp != nil &&
			c.TunnelDstMac != nil {
			// adds arp entry related to tunnel dst
			neigh := &netlink.Neigh{
				LinkIndex:    link.Attrs().Index,
				Family:       netlink.NDA_VNI,                           // NDA_VNI
				State:        netlink.NUD_NOARP | netlink.NUD_PERMANENT, // NUD_NOARP (0x40) | NUD_PERMANENT (0x80)
				Type:         1,
				Flags:        netlink.NTF_SELF, // NTF_SELF
				IP:           c.TunnelDstIp,
				HardwareAddr: c.TunnelDstMac,
			}
			v.logger.Info(fmt.Sprintf("neighbor: %#v", neigh))
			if err := netlink.NeighSet(neigh); err != nil {
				v.logger.Err(fmt.Sprintf("NeighSet:", err.Error()))
			}
		retry_neighbor_set:
			neighList, err := netlink.NeighList(neigh.LinkIndex, neigh.Family)
			if err == nil {

				for _, n := range neighList {
					foundNeighbor := false
					if len(neigh.IP) == len(n.IP) {
						for i, _ := range neigh.IP {
							if neigh.IP[i] == n.IP[i] {
								foundNeighbor = true
							} else {
								foundNeighbor = false
							}
						}
					}
					if foundNeighbor {
						v.logger.Info("Found Neighbor ip")
						if n.State == netlink.NUD_FAILED {
							v.logger.Info(fmt.Sprintf("retry neighbor: %#v", neigh))
							if err := netlink.NeighSet(neigh); err != nil {
								v.logger.Err(fmt.Sprintf("NeighSet:", err.Error()))
								goto retry_neighbor_set
							}
						}
					}
				}
			}
		} else {
			v.logger.Info(fmt.Sprintf("neighbor: not configured dstIp %#v dstmac %#v", c.TunnelDstIp, c.TunnelDstMac))
		}
	}
	vxlanDbEntry := VxlanDB[uint32(c.Vni)]
	vxlanDbEntry.Links = append(vxlanDbEntry.Links, &link)
	VxlanDB[uint32(c.Vni)] = vxlanDbEntry

	if err := netlink.LinkSetMaster(link, vxlanDbEntry.Brg); err != nil {
		v.logger.Err(err.Error())
	}

	/* ON RECREATE - Link up is failing with reason:
	   transport endpoint is not connected lets delay
	   till it is connected */
	// lets set the vtep interface to up
	for i := 0; i < 10; i++ {
		err := netlink.LinkSetUp(link)
		if err != nil && i < 10 {
			v.logger.Info(fmt.Sprintf("createVtep: %s link not connected yet waiting 5ms", c.VtepName))
			time.Sleep(time.Millisecond * 5)
		} else if err != nil {
			v.logger.Err(err.Error())
		} else {
			break
		}
	}
}

func (v *VxlanLinux) DeleteVtep(c *VtepConfig) {

	foundEntry := false
	if vxlanentry, ok := VxlanDB[c.Vni]; ok {
		for i, link := range vxlanentry.Links {
			var linkName string
			if VxlanConfigMode == "linux" {
				linkName = (*link).(*netlink.Vxlan).Attrs().Name
			} else {
				linkName = (*link).(*netlink.Veth).Attrs().Name
			}
			if linkName == c.VtepName {
				v.logger.Info(fmt.Sprintf("deleteVtep: link found %s looking for %s", linkName, c.VtepName))
				foundEntry = true
				vxlanDbEntry := VxlanDB[c.Vni]
				vxlanDbEntry.Links = append(vxlanDbEntry.Links[:i], vxlanDbEntry.Links[i+1:]...)
				VxlanDB[c.Vni] = vxlanDbEntry
				break
			}
		}
	}

	if foundEntry {
		link, err := netlink.LinkByName(c.VtepName)
		if err != nil {
			v.logger.Err(err.Error())
		}
		if err := netlink.LinkSetDown(link); err != nil {
			v.logger.Err(err.Error())
		}

		if err := netlink.LinkDel(link); err != nil {
			v.logger.Err(err.Error())
		}
	} else {
		v.logger.Err("Unable to find vtep in vxlan db")
	}
}

func (v *VxlanLinux) LearnFdbVtep(mac string, vtepname string, ifindex int32) {

	if macDb == nil {
		macDb = make(map[int32][]*VxlanMacDbEntry, 0)
	}

	if macList, ok := macDb[ifindex]; ok {
		for _, macentry := range macList {
			if macentry.mac == mac {
				return
			}
		}
		macDb[ifindex] = append(macDb[ifindex], &VxlanMacDbEntry{mac: mac,
			vtepName: vtepname})
		link, _ := netlink.LinkByName(vtepname)
		if link != nil {
			netmac, _ := net.ParseMAC(mac)
			neigh := &netlink.Neigh{
				LinkIndex: link.Attrs().Index,
				//Family:       netlink.NDA_VNI,                           // NDA_VNI
				State:        netlink.NUD_NOARP, // NUD_NOARP (0x40) | NUD_PERMANENT (0x80)
				Type:         1,
				Flags:        netlink.NTF_SELF, // NTF_SELF
				HardwareAddr: netmac,
			}
			if err := netlink.NeighAppend(neigh); err != nil {
				v.logger.Err(fmt.Sprintf("NeighSet:", err.Error()))
			}
		}
	}
}
