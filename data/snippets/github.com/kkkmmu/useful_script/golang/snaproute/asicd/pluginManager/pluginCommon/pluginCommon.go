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

package pluginCommon

import (
	"net"
	"utils/commonDefs"
)

/*
#include "pluginCommon.h"
*/
import "C"

var SwitchMacAddr net.HardwareAddr

//Consolidated list of constants used by underlying plugins
const (
	//Generic consts
	DEFAULT_SWITCH_MAC_ADDR = "00:11:22:33:44:55"
	MAC_ADDR_LEN            = C.MAC_ADDR_LEN
	INVALID_OBJECT_ID       = 0xFFFFFFFFFFFFFFFF
	INVALID_IFINDEX         = -1

	//System related consts
	BOOT_MODE_COLDBOOT = C.BOOT_MODE_COLDBOOT
	BOOT_MODE_WARMBOOT = C.BOOT_MODE_WARMBOOT
	INTF_STATE_UP      = C.INTF_STATE_UP
	INTF_STATE_DOWN    = C.INTF_STATE_DOWN
	MIN_SYS_PORTS      = 0
	MAX_SYS_PORTS      = 256
	ASICD_CONFIG_FILE  = "asicdConf.json"
	MAX_VENDOR_ID_LEN  = C.MAX_VENDOR_ID_LEN
	MAX_PART_NUM_LEN   = C.MAX_PART_NUM_LEN
	MAX_REV_ID_LEN     = C.MAX_REV_ID_LEN

	//FDB relate consts
	MAC_ENTRY_LEARNED = C.MAC_ENTRY_LEARNED
	MAC_ENTRY_AGED    = C.MAC_ENTRY_AGED

	//Vlan related consts
	SVI_PREFIX        = "vlan"
	SYS_RSVD_VLAN     = -1
	MAX_VLAN_ID       = C.MAX_VLAN_ID
	DEFAULT_VLAN_ID   = C.DEFAULT_VLAN_ID
	SYS_RSVD_VLAN_MIN = 3835
	SYS_RSVD_VLAN_MAX = 4090

	//Port related consts
	//when adding any new attribute and updating values here please update
	//utils/commonDefs/asicdClient.go also
	MAX_PORT_STAT_TYPES                   = C.portStatTypesMax
	MAX_BUFFER_STAT_TYPES                 = C.bufferPortStatTypesMax
	MAX_BUFFER_GLOBAL_STAT_TYPES          = C.bufferGlobalStatTypesMax
	PORT_BREAKOUT_MODE_UNSUPPORTED_STRING = "unsupported"
	PORT_BREAKOUT_MODE_UNSUPPORTED        = C.PORT_BREAKOUT_MODE_UNSUPPORTED
	PORT_BREAKOUT_MODE_1x40               = C.PORT_BREAKOUT_MODE_1x40
	PORT_BREAKOUT_MODE_4x10               = C.PORT_BREAKOUT_MODE_4x10
	PORT_BREAKOUT_MODE_4x25               = C.PORT_BREAKOUT_MODE_4x25
	PORT_BREAKOUT_MODE_2x50               = C.PORT_BREAKOUT_MODE_2x50
	PORT_BREAKOUT_MODE_1x100              = C.PORT_BREAKOUT_MODE_1x100
	PORT_ATTR_PHY_INTF_TYPE               = C.PORT_ATTR_PHY_INTF_TYPE
	PORT_ATTR_ADMIN_STATE                 = C.PORT_ATTR_ADMIN_STATE
	PORT_ATTR_MAC_ADDR                    = C.PORT_ATTR_MAC_ADDR
	PORT_ATTR_SPEED                       = C.PORT_ATTR_SPEED
	PORT_ATTR_DUPLEX                      = C.PORT_ATTR_DUPLEX
	PORT_ATTR_AUTONEG                     = C.PORT_ATTR_AUTONEG
	PORT_ATTR_MEDIA_TYPE                  = C.PORT_ATTR_MEDIA_TYPE
	PORT_ATTR_MTU                         = C.PORT_ATTR_MTU
	PORT_ATTR_BREAKOUT_MODE               = C.PORT_ATTR_BREAKOUT_MODE
	PORT_ATTR_LOOPBACK_MODE               = C.PORT_ATTR_LOOPBACK_MODE
	PORT_ATTR_ENABLE_FEC                  = C.PORT_ATTR_ENABLE_FEC
	PORT_ATTR_TX_PRBS_EN                  = C.PORT_ATTR_TX_PRBS_EN
	PORT_ATTR_RX_PRBS_EN                  = C.PORT_ATTR_RX_PRBS_EN
	PORT_ATTR_PRBS_POLY                   = C.PORT_ATTR_PRBS_POLY
	PORT_ATTR_DESCRIPTION                 = C.PORT_ATTR_DESCRIPTION
	PORT_ATTR_PVID                        = C.PORT_ATTR_PVID
	PORT_MODE_UNCONFIGURED                = "Unconfigured"
	PORT_MODE_L2                          = "L2"
	PORT_MODE_L3                          = "L3"
	PORT_MODE_INTERNAL                    = "Internal"
	PORT_LOOPBACK_MODE_NONE               = C.PORT_LOOPBACK_MODE_NONE
	PORT_LOOPBACK_MODE_MAC                = C.PORT_LOOPBACK_MODE_MAC
	PORT_LOOPBACK_MODE_PHY                = C.PORT_LOOPBACK_MODE_PHY

	//Intf related constants
	IP_TYPE_IPV4 = C.IP_TYPE_IPV4
	IP_TYPE_IPV6 = C.IP_TYPE_IPV6

	//STP related consts
	STP_PORT_STATE_BLOCKING   = C.StpPortStateBlocking
	STP_PORT_STATE_LEARNING   = C.StpPortStateLearning
	STP_PORT_STATE_FORWARDING = C.StpPortStateForwarding

	//Lag related consts
	LAG_PREFIX             = "LAG"
	HASHTYPE_SRCMAC_DSTMAC = C.HASHTYPE_SRCMAC_DSTMAC
	HASHTYPE_SRCIP_DSTIP   = C.HASHTYPE_SRCIP_DSTIP

	//Next hop related consts
	NEIGHBOR_TYPE_COPY_TO_CPU       = C.NEIGHBOR_TYPE_COPY_TO_CPU
	NEIGHBOR_TYPE_BLACKHOLE         = C.NEIGHBOR_TYPE_BLACKHOLE
	NEIGHBOR_TYPE_FULL_SPEC_NEXTHOP = C.NEIGHBOR_TYPE_FULL_SPEC_NEXTHOP
	NEIGHBOR_L2_ACCESS_TYPE_PORT    = C.NEIGHBOR_L2_ACCESS_TYPE_PORT
	NEIGHBOR_L2_ACCESS_TYPE_LAG     = C.NEIGHBOR_L2_ACCESS_TYPE_LAG

	//Next hop related consts
	NEXTHOP_TYPE_COPY_TO_CPU       = C.NEXTHOP_TYPE_COPY_TO_CPU
	NEXTHOP_TYPE_BLACKHOLE         = C.NEXTHOP_TYPE_BLACKHOLE
	NEXTHOP_TYPE_FULL_SPEC_NEXTHOP = C.NEXTHOP_TYPE_FULL_SPEC_NEXTHOP
	NEXTHOP_L2_ACCESS_TYPE_PORT    = C.NEXTHOP_L2_ACCESS_TYPE_PORT
	NEXTHOP_L2_ACCESS_TYPE_LAG     = C.NEXTHOP_L2_ACCESS_TYPE_LAG
	INVALID_NEXTHOP_ID             = 0xFFFFFFFFFFFFFFFF

	// vxlan related consts
	VXLAN_VTEP_PREFIX  = "Vtep"
	VXLAN_VXLAN_PREFIX = "Vxlan"

	// SUB interface related consts
	SUB_INTF_SECONDARY_TYPE = "secondary"
	SUB_INTF_VIRTUAL_TYPE   = "virtual"
	SECONDARY_TYPE          = C.SECONDARY_INTF
	VIRTUAL_TYPE            = C.VIRTUAL_INTF

	//Route related consts
	MAX_NEXTHOPS_PER_GROUP      = C.MAX_NEXTHOPS_PER_GROUP
	ROUTE_OPERATION_TYPE_UPDATE = C.ROUTE_OPERATION_TYPE_UPDATE
	ROUTE_TYPE_CONNECTED        = C.ROUTE_TYPE_CONNECTED
	ROUTE_TYPE_MULTIPATH        = C.ROUTE_TYPE_MULTIPATH
	ROUTE_TYPE_SINGLEPATH       = C.ROUTE_TYPE_SINGLEPATH
	ROUTE_TYPE_V6               = C.ROUTE_TYPE_V6
	ROUTE_TYPE_NULL             = C.ROUTE_TYPE_NULL

	//ACL and COPP related consts
	MAX_COPP_CLASS_COUNT = C.CoppClassCount
	ACL_L3_ROUTABLE      = C.ACL_L3_ROUTABLE
	ACL_IPV4_ENABLED     = C.ACL_IPV4_ENABLED
	ACL_IPV6_ENABLED     = C.ACL_IPV6_ENABLED
	ACL_MAC_ENABLED      = C.ACL_MAC_ENABLED
	ACL_VLAN_ENABLED     = C.ACL_VLAN_ENABLED
	ACL_PORT_ENABLED     = C.ACL_PORT_ENABLED
	ACL_HITLESS_UPDATE   = C.ACL_HITLESS_UPDATE
)

// Interface id/type mgmt
const (
	INTF_TYPE_MASK  = C.INTF_TYPE_MASK
	INTF_TYPE_SHIFT = C.INTF_TYPE_SHIFT
	INTF_ID_MASK    = C.INTF_ID_MASK
	INTF_ID_SHIFT   = C.INTF_ID_SHIFT
)

const (
	PORT_PROTOCOL_ARP       = 0x1
	PORT_PROTOCOL_DHCP      = 0x2
	PORT_PROTOCOL_DHCPRELAY = 0x4
	PORT_PROTOCOL_BGP       = 0x8
	PORT_PROTOCOL_OSPF      = 0x10
	PORT_PROTOCOL_VXLAN     = 0x20
	PORT_PROTOCOL_MPLS      = 0x40
	PORT_PROTOCOL_BFD       = 0x40
)

const (
	ACL_DIRECTION_IN  = "IN"
	ACL_DIRECTION_OUT = "OUT"
	ACL_TYPE_IP       = "IPv4"
	ACL_TYPE_MAC      = "MAC"
	ACL_TYPE_IP6      = "IPv6"
	ACL_TYPE_MLAG     = "MLAG"
)

const (
	PLATFORM_ID_INVALID     = C.PLATFORM_ID_INVALID
	PLATFORM_CEL_REDSTONE   = "x86-64-cel-redstone-xp-r0"
	PLATFORM_LNKD_OPEN19    = "x86-64-qwave-hawk-r0"
	PLATFORM_CEL_SEASTONE   = "x86-64-cel-seastone-r0"
	PLATFORM_ACC_AS5712     = "x86-64-accton-as5712-54x-r0"
	PLATFORM_ACC_AS5812_54T = "x86-64-accton-as5812-54t-r0"
	PLATFORM_ACC_WEDGE40    = "x86-64-accton-wedge-16x-r0"
	PLATFORM_ACC_WEDGE100   = "x86-64-facebook-wedge100-r0"
	PLATFORM_ACC_WEDGE100_2 = "x86-64-accton-wedge100-32x-r0"
	PLATFORM_ACC_AS7712     = "x86-64-accton-as7712-32x-r0"
	PLATFORM_ACC_AS6712     = "x86-64-accton-as6712-32x-r0"
	PLATFORM_ACC_AS6812     = "x86-64-accton-as6812-32x-r0"
	PLATFORM_CEL_VOYAGER    = "x86-64-cel-voyager"
	//Verify platform name
	PLATFORM_ACC_AS5912 = "x86-64-accton-as5912-32x-r0"
	PLATFORM_DELL_S6000 = "x86-64-dell-s6000-s1220-r0"
)

var PlatformIdMap map[string]int32 = map[string]int32{
	PLATFORM_CEL_REDSTONE:   C.PLATFORM_CEL_REDSTONE,
	PLATFORM_LNKD_OPEN19:    C.PLATFORM_LNKD_OPEN19,
	PLATFORM_CEL_SEASTONE:   C.PLATFORM_CEL_SEASTONE,
	PLATFORM_ACC_AS5712:     C.PLATFORM_ACC_AS5712,
	PLATFORM_ACC_AS5812_54T: C.PLATFORM_ACC_AS5812_54T,
	PLATFORM_ACC_WEDGE40:    C.PLATFORM_ACC_WEDGE40,
	PLATFORM_ACC_WEDGE100:   C.PLATFORM_ACC_WEDGE100,
	PLATFORM_ACC_WEDGE100_2: C.PLATFORM_ACC_WEDGE100,
	PLATFORM_ACC_AS7712:     C.PLATFORM_ACC_AS7712,
	PLATFORM_ACC_AS6712:     C.PLATFORM_ACC_AS6712,
	PLATFORM_ACC_AS6812:     C.PLATFORM_ACC_AS6812,
	PLATFORM_CEL_VOYAGER:    C.PLATFORM_CEL_VOYAGER,
	//Verify platform name
	PLATFORM_ACC_AS5912: C.PLATFORM_ACC_AS5912,
	PLATFORM_DELL_S6000: C.PLATFORM_DELL_S6000,
}

// Func types for intf id mgmt
type GetId func(int32) int
type GetType func(int32) int
type GetIfIndex func(int, int) int32

func GetTypeFromIfIndex(ifIndex int32) int {
	return int((ifIndex & INTF_TYPE_MASK) >> INTF_TYPE_SHIFT)
}
func GetIdFromIfIndex(ifIndex int32) int {
	return int((ifIndex & INTF_ID_MASK) >> INTF_ID_SHIFT)
}
func GetIfIndexFromIdType(ifId, ifType int) int32 {
	return int32((ifId & INTF_ID_MASK) | ((ifType << INTF_TYPE_SHIFT) & INTF_TYPE_MASK))
}
func IsZeroIP(ip []uint8, ipType int) bool {
	/*    if ipType == syscall.AF_INET {
	    if bytes.Equal(ip,net.IPv4zero) {
		    return true
	    }
	} else {
	    if bytes.Equal(ip,net.IPv6zero) {
		    return true
	    }
	}
	return false*/
	for i := 0; i < len(ip); i++ {
		if ip[i] != 0 {
			return false
		}
	}
	return true
}

const (
	LB_MODE_NONE = "NONE"
	LB_MODE_MAC  = "MAC"
	LB_MODE_PHY  = "PHY"
	LB_MODE_RMT  = "RMT"
)

var LBModeEnumToStr map[uint8]string = map[uint8]string{
	uint8(C.PORT_LOOPBACK_MODE_NONE): LB_MODE_NONE,
	uint8(C.PORT_LOOPBACK_MODE_MAC):  LB_MODE_MAC,
	uint8(C.PORT_LOOPBACK_MODE_PHY):  LB_MODE_PHY,
	uint8(C.PORT_LOOPBACK_MODE_RMT):  LB_MODE_RMT,
}

var LBModeStrToEnum map[string]uint8 = map[string]uint8{
	LB_MODE_NONE: uint8(C.PORT_LOOPBACK_MODE_NONE),
	LB_MODE_MAC:  uint8(C.PORT_LOOPBACK_MODE_MAC),
	LB_MODE_PHY:  uint8(C.PORT_LOOPBACK_MODE_PHY),
	LB_MODE_RMT:  uint8(C.PORT_LOOPBACK_MODE_RMT),
}

const (
	PRBS_POLY_2POW7  = "2^7"
	PRBS_POLY_2POW23 = "2^23"
	PRBS_POLY_2POW31 = "2^31"
)

var PrbsStrToEnum map[string]uint8 = map[string]uint8{
	PRBS_POLY_2POW7:  uint8(C.PRBS_POLY_2POW7),
	PRBS_POLY_2POW23: uint8(C.PRBS_POLY_2POW23),
	PRBS_POLY_2POW31: uint8(C.PRBS_POLY_2POW31),
}

var PrbsEnumToStr map[uint8]string = map[uint8]string{
	uint8(C.PRBS_POLY_2POW7):  PRBS_POLY_2POW7,
	uint8(C.PRBS_POLY_2POW23): PRBS_POLY_2POW23,
	uint8(C.PRBS_POLY_2POW31): PRBS_POLY_2POW31,
}

var OnOffState map[int]string = map[int]string{0: "OFF", 1: "ON"}

const (
	STATE_UP   = "UP"
	STATE_DOWN = "DOWN"
)

var UpDownState map[int]string = map[int]string{0: STATE_DOWN, 1: STATE_UP}
var UpDownStateStr map[string]int = map[string]int{STATE_DOWN: 0, STATE_UP: 1}
var IfType map[int]string = map[int]string{
	int(C.PortIfTypeMII):    "MII",
	int(C.PortIfTypeGMII):   "GMII",
	int(C.PortIfTypeSGMII):  "SGMII",
	int(C.PortIfTypeQSGMII): "QSGMII",
	int(C.PortIfTypeSFI):    "SFI",
	int(C.PortIfTypeXFI):    "XFI",
	int(C.PortIfTypeXAUI):   "XAUI",
	int(C.PortIfTypeXLAUI):  "XLAUI",
	int(C.PortIfTypeRXAUI):  "RXAUI",
	int(C.PortIfTypeCR):     "CR",
	int(C.PortIfTypeCR2):    "CR2",
	int(C.PortIfTypeCR4):    "CR4",
	int(C.PortIfTypeKR):     "KR",
	int(C.PortIfTypeKR2):    "KR2",
	int(C.PortIfTypeKR4):    "KR4",
	int(C.PortIfTypeSR):     "SR",
	int(C.PortIfTypeSR2):    "SR2",
	int(C.PortIfTypeSR4):    "SR4",
	int(C.PortIfTypeSR10):   "SR10",
	int(C.PortIfTypeLR):     "LR",
	int(C.PortIfTypeLR4):    "LR4",
}
var PluginIfType map[string]int = map[string]int{
	"MII":    int(C.PortIfTypeMII),
	"GMII":   int(C.PortIfTypeGMII),
	"SGMII":  int(C.PortIfTypeSGMII),
	"QSGMII": int(C.PortIfTypeQSGMII),
	"SFI":    int(C.PortIfTypeSFI),
	"XFI":    int(C.PortIfTypeXFI),
	"XAUI":   int(C.PortIfTypeXAUI),
	"XLAUI":  int(C.PortIfTypeXLAUI),
	"RXAUI":  int(C.PortIfTypeRXAUI),
	"CR":     int(C.PortIfTypeCR),
	"CR2":    int(C.PortIfTypeCR2),
	"CR4":    int(C.PortIfTypeCR4),
	"KR":     int(C.PortIfTypeKR),
	"KR2":    int(C.PortIfTypeKR2),
	"KR4":    int(C.PortIfTypeKR4),
	"SR":     int(C.PortIfTypeSR),
	"SR2":    int(C.PortIfTypeSR2),
	"SR4":    int(C.PortIfTypeSR4),
	"SR10":   int(C.PortIfTypeSR10),
	"LR":     int(C.PortIfTypeLR),
	"LR4":    int(C.PortIfTypeLR4),
}
var PlatformSfpTypeMap map[int]int = map[int]int{
	commonDefs.SFF_MODULE_TYPE_100G_AOC:      int(C.PortIfTypeSR4),
	commonDefs.SFF_MODULE_TYPE_100G_BASE_CR4: int(C.PortIfTypeCR4),
	commonDefs.SFF_MODULE_TYPE_100G_BASE_SR4: int(C.PortIfTypeSR4),
	commonDefs.SFF_MODULE_TYPE_100G_BASE_LR4: int(C.PortIfTypeLR4),
	commonDefs.SFF_MODULE_TYPE_40G_BASE_CR4:  int(C.PortIfTypeCR4),
	commonDefs.SFF_MODULE_TYPE_40G_BASE_SR4:  int(C.PortIfTypeSR4),
	commonDefs.SFF_MODULE_TYPE_40G_BASE_LR4:  int(C.PortIfTypeLR4),
	commonDefs.SFF_MODULE_TYPE_40G_BASE_CR:   int(C.PortIfTypeCR),
	commonDefs.SFF_MODULE_TYPE_40G_BASE_SR2:  int(C.PortIfTypeSR2),
	commonDefs.SFF_MODULE_TYPE_10G_BASE_SR:   int(C.PortIfTypeSR),
	commonDefs.SFF_MODULE_TYPE_10G_BASE_LR:   int(C.PortIfTypeLR),
	commonDefs.SFF_MODULE_TYPE_10G_BASE_CR:   int(C.PortIfTypeCR),
}

const (
	FULL_DUPLEX = "Full_Duplex"
	HALF_DUPLEX = "Half_Duplex"
)
const (
	PluginOp_NA = iota
	PluginOp_Add
	PluginOp_Del
)

var DuplexType map[int]string = map[int]string{
	int(C.HalfDuplex): HALF_DUPLEX,
	int(C.FullDuplex): FULL_DUPLEX,
}

const (
	//Asicd notification msgs new notifications should always be added at the bottom of the list
	//also when this list is update make sure utils/commonDefs/asicdClient.go is updated or someone will
	//hunt you down ;)
	NOTIFY_L2INTF_STATE_CHANGE = iota
	NOTIFY_IPV4_L3INTF_STATE_CHANGE
	NOTIFY_IPV6_L3INTF_STATE_CHANGE
	NOTIFY_VLAN_CREATE
	NOTIFY_VLAN_DELETE
	NOTIFY_VLAN_UPDATE
	NOTIFY_LOGICAL_INTF_CREATE
	NOTIFY_LOGICAL_INTF_DELETE
	NOTIFY_LOGICAL_INTF_UPDATE
	NOTIFY_IPV4INTF_CREATE
	NOTIFY_IPV4INTF_DELETE
	NOTIFY_IPV6INTF_CREATE
	NOTIFY_IPV6INTF_DELETE
	NOTIFY_LAG_CREATE
	NOTIFY_LAG_DELETE
	NOTIFY_LAG_UPDATE
	NOTIFY_IPV4NBR_MAC_MOVE
	NOTIFY_IPV6NBR_MAC_MOVE
	NOTIFY_IPV4_ROUTE_CREATE_FAILURE
	NOTIFY_IPV4_ROUTE_DELETE_FAILURE
	NOTIFY_IPV6_ROUTE_CREATE_FAILURE
	NOTIFY_IPV6_ROUTE_DELETE_FAILURE
	NOTIFY_VTEP_CREATE
	NOTIFY_VTEP_DELETE
	NOTIFY_MPLSINTF_STATE_CHANGE
	NOTIFY_MPLSINTF_CREATE
	NOTIFY_MPLSINTF_DELETE
	NOTIFY_PORT_CONFIG_MODE_CHANGE
	NOTIFY_PORT_ATTR_CHANGE
	NOTIFY_IPV4VIRTUAL_INTF_CREATE
	NOTIFY_IPV4VIRTUAL_INTF_DELETE
	NOTIFY_IPV6VIRTUAL_INTF_CREATE
	NOTIFY_IPV6VIRTUAL_INTF_DELETE
	NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE
	NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE
)

// Format of asicd's published messages
type AsicdNotification struct {
	MsgType uint8
	Msg     []byte
}

// Following sections contains formats for individual message types
type L2IntfStateNotifyMsg struct {
	IfIndex int32
	IfState uint8
}
type MplsIntfStateNotifyMsg struct {
	IfIndex int32
	IfState uint8
}
type IPv4L3IntfStateNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	IfState uint8
}
type IPv6L3IntfStateNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	IfState uint8
}
type VlanNotifyMsg struct {
	VlanId      uint16
	VlanIfIndex int32
	VlanName    string
	TagPorts    []int32
	UntagPorts  []int32
}
type LogicalIntfNotifyMsg struct {
	IfIndex         int32
	LogicalIntfName string
}
type LagNotifyMsg struct {
	LagName     string
	IfIndex     int32
	IfIndexList []int32
}
type MplsIntfNotifyMsg struct {
	IfIndex int32
}
type IPv4IntfNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	IntfRef string
}
type IPv4NbrMacMoveNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	VlanId  int32
}
type IPv6NbrMacMoveNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	VlanId  int32
}
type IPv6IntfNotifyMsg struct {
	IpAddr  string
	IfIndex int32
	IntfRef string
}
type VtepNotifyMsg struct {
	VtepName   string
	IfIndex    int32
	Vni        int32
	SrcIfIndex int32
	SrcIfName  string
}
type PortConfigModeChgNotifyMsg struct {
	IfIndex int32
	OldMode string
	NewMode string
}
type PortAttrChangeNotifyMsg struct {
	IfIndex     int32
	Mtu         int32
	Description string
	Pvid        int32
	AttrMask    int32
}

type IPv4VirtualIntfNotifyMsg struct {
	IfIndex       int32
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	IfName        string
}

type IPv6VirtualIntfNotifyMsg struct {
	IfIndex       int32
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	IfName        string
}

type IPv4VirtualIntfStateNotifyMsg struct {
	IfIndex int32
	IpAddr  string
	IfState uint8
}

type IPv6VirtualIntfStateNotifyMsg struct {
	IfIndex int32
	IpAddr  string
	IfState uint8
}

// Struct containing info required for mapping asic ports to linux interfaces
type IfMapInfo struct {
	IfName string
	Port   int
}

// Struct used for configuring ports
type PortConfig struct {
	PortNum           int32
	PortName          string
	Description       string
	PhyIntfType       string
	AdminState        string
	MacAddr           string
	Speed             int32
	Duplex            string
	Autoneg           string
	MediaType         string
	Mtu               int32
	LogicalPortInfo   bool
	MappedToHw        bool
	BreakOutMode      int32
	BreakOutSupported bool
	LoopbackMode      uint8
	EnableFEC         bool
	PRBSTxEnable      bool
	PRBSRxEnable      bool
	PRBSPolynomial    uint8
	Dirty             bool
	SfpId             uint8
}

// Struct used for retrieving port state
type PortState struct {
	PortNum                     int32
	IfIndex                     int32
	Name                        string
	OperState                   string
	NumUpEvents                 int32
	LastUpEventTime             string
	NumDownEvents               int32
	LastDownEventTime           string
	Pvid                        int32
	IfInOctets                  int64
	IfInUcastPkts               int64
	IfInDiscards                int64
	IfInErrors                  int64
	IfInUnknownProtos           int64
	IfOutOctets                 int64
	IfOutUcastPkts              int64
	IfOutDiscards               int64
	IfOutErrors                 int64
	IfEtherUnderSizePktCnt      int64
	IfEtherOverSizePktCnt       int64
	IfEtherFragments            int64
	IfEtherCRCAlignError        int64
	IfEtherJabber               int64
	IfEtherPkts                 int64
	IfEtherMCPkts               int64
	IfEtherBcastPkts            int64
	IfEtherPkts64OrLessOctets   int64
	IfEtherPkts65To127Octets    int64
	IfEtherPkts128To255Octets   int64
	IfEtherPkts256To511Octets   int64
	IfEtherPkts512To1023Octets  int64
	IfEtherPkts1024To1518Octets int64
	ErrDisableReason            string
	PRBSRxErrCnt                int64
}

// Buffer State for retrieving buffer state
type BufferPortState struct {
	IntfRef            string
	IfIndex            int32
	EgressPort         uint64
	PortBufferPortStat uint64
}

// Buffer state for global buffer stats
type BufferGlobalState struct {
	DeviceId          int32
	BufferStat        uint64
	EgressBufferStat  uint64
	IngressBufferStat uint64
}

// Struct used for ECMP/WCMP group creation
type NextHopGroupMemberInfo struct {
	IpAddr    interface{} //uint32
	NextHopId uint64
	Weight    int
	RifId     int32
}

var CoppStringToInt map[string]int = map[string]int{
	"ArpUC":    int(C.CoppArpUC),
	"ArpMC":    int(C.CoppArpMC),
	"BGP":      int(C.CoppBgp),
	"ICMPv4UC": int(C.CoppIcmpV4UC),
	"ICMPv4BC": int(C.CoppIcmpV4BC),
	"STP":      int(C.CoppStp),
	"LACP":     int(C.CoppLacp),
	"BFD":      int(C.CoppBfd),
	"ICMPv6":   int(C.CoppIcmpv6),
	"LLDP":     int(C.CoppLldp),
}

var CoppType map[int]string = map[int]string{
	int(C.CoppArpUC):    "ArpUC",
	int(C.CoppArpMC):    "ArpMC",
	int(C.CoppBgp):      "BGP",
	int(C.CoppIcmpV4UC): "ICMPv4UC",
	int(C.CoppIcmpV4BC): "ICMPv4BC",
	int(C.CoppStp):      "STP",
	int(C.CoppLacp):     "LACP",
	int(C.CoppBfd):      "BFD",
	int(C.CoppIcmpv6):   "ICMPv6",
	int(C.CoppLldp):     "LLDP",
	//	int(C.CoppSsh):      "SSH",
	//	int(C.CoppHttp):     "HTTP",
}

type CoppStatState struct {
	CoppClass    string
	PeakRate     int32
	BurstRate    int32
	GreenPackets int64
	RedPackets   int64
}

//struct used for ACL creation
//Acl proto type to IP number map
var AclProtoType map[string]int = map[string]int{
	"TCP":    C.COPP_IP_PROTOCOL_IP_NUMBER_TCP,
	"UDP":    C.COPP_IP_PROTOCOL_IP_NUMBER_UDP,
	"ICMPV4": C.COPP_IP_PROTOCOL_IPV4_NUMBER_ICMP,
	"ICMPV6": C.COPP_IP_PROTOCOL_IPV6_NUMBER_ICMP,
	"OSPFV2": C.COPP_IP_PROTOCOL_IP_NUMBER_OSPFV2,
}

const (
	AclEq    = 0
	AclNeq   = 1
	AclRange = 2
)

type Acl struct {
	AclName   string
	AclType   string
	Flags     int32
	VlanList  []int32
	PortList  []int32
	Direction string
	Priority  int
	Action    string
	Filter    interface{}
}

type AclIpv4Filter struct {
	Proto       int
	SourceIp    []uint8
	DestIp      []uint8
	SourceMask  []uint8
	DestMask    []uint8
	SrcIntf     int32
	DstIntf     int32
	L4SrcPort   int32
	L4DstPort   int32
	L4PortMatch int
	L4MinPort   int32
	L4MaxPort   int32
}

type AclMacFilter struct {
	SourceMac  net.HardwareAddr
	DestMac    net.HardwareAddr
	SourceMask net.HardwareAddr
	DestMask   net.HardwareAddr
}

type AclIpv6Filter struct {
	SourceIpv6   []uint8
	DestIpv6     []uint8
	SourceMaskv6 []uint8
	DestMaskv6   []uint8
	L4SrcPort    int32
	L4DstPort    int32
	L4PortMatch  int
	L4MinPort    int32
	L4MaxPort    int32
	Proto        int
}

/* HW object for Acl */
type AclConfig struct {
	AclName string
	HwIndex int   //Hw entry index
	Flags   int32 //Acl flags for entrys
	StatId  int   // stat object
}
type AclState struct {
	AclName string
	Stat    uint64
}

type Copp struct {
	ProtoId          int
	PolicerPeakRate  int
	PolicerBurstRate int
	CpuQueue         int
}

//Format of callback functions for updating DBs in individual resource managers
type ProcessLinkStateChangeCB func(int32, int32, string, string)
type InitPortConfigDBCB func(*PortConfig)
type InitPortStateDBCB func(int32, string)
type UpdatePortStateDBCB func(*PortState)
type UpdateLagDBCB func(int, int, []int32)
type UpdateIPNeighborDBCB func(*PluginIPNeighborInfo)
type UpdateIPv4NeighborDBCB func(uint32, net.HardwareAddr, int32, uint64, int)
type UpdateIPv4RouteDBCB func(uint32, uint32, uint32)
type UpdateIPv4NextHopDBCB func()
type UpdateIPv4NextHopGroupDBCB func()
type MacEntryTableCB func(int, int32, int32, net.HardwareAddr)
type InitBufferPortStateDBCB func(int32, string)
type UpdateBufferPortStateDBCB func(*BufferPortState)
type InitBufferGlobalStateDBCB func(int) int
type UpdateBufferGlobalStateDBCB func(*BufferGlobalState)
type UpdateCoppStatStateDBCB func(*CoppStatState)
type UpdateAclConfigDBCB func(*AclConfig)
type UpdateAclStateDBCB func(*AclState)

func ComputeStrSetDifference(a, b []string) (aDiffb []string) {
	var bMap map[string]bool = make(map[string]bool, 0)
	if len(a) == 0 {
		return a
	}
	if len(b) == 0 {
		return a
	}
	for _, elem := range b {
		bMap[elem] = true
	}
	for _, elem := range a {
		if _, ok := bMap[elem]; !ok {
			aDiffb = append(aDiffb, elem)
		}
	}
	return aDiffb
}
func ComputeSetDifference(a, b []int32) (aDiffb []int32) {
	var bMap map[int32]bool = make(map[int32]bool, 0)
	if len(a) == 0 {
		return a
	}
	if len(b) == 0 {
		return a
	}
	for _, elem := range b {
		bMap[elem] = true
	}
	for _, elem := range a {
		if _, ok := bMap[elem]; !ok {
			aDiffb = append(aDiffb, elem)
		}
	}
	return aDiffb
}

type PluginIPInfo struct {
	IfIndex      int32
	IpAddr       []uint32
	MaskLen      int
	IpType       int
	VlanId       int
	IfName       string
	Address      string // if a plugin uses string then it will be copied in here... for e.g. 192.168.1.1/32
	RefCount     int
	InstallRoute int
	IPv6Type     int
	LinkLocalIp  string
	AdminState   int
}

// Struct used for configuring sub interface on all plugins
type SubIntfPluginObj struct {
	IntfRef        string
	IfIndex        int32
	SubIntfIfIndex int32
	IpAddr         []uint32
	MaskLen        int
	IpType         int
	VlanId         int
	Address        string
	Type           string
	StateUp        bool
	MacAddr        net.HardwareAddr
	IpNet          *net.IPNet
	IPv6Type       int
	RefCount       int
}

type PluginIPRouteInfo struct {
	PrefixIp      []uint8
	NwAddr        string
	IpType        int //destination ip type (v4/v6)
	Mask          []uint8
	NextHopIp     []uint8
	NextHopIpType int //next hop ip type (v4/v6)
	NextHopIfType int
	Weight        int
	//info to use to communicate with plugins
	Op            uint8
	RouteFlags    uint32
	NextHopId     uint64
	RouterIfIndex int32
	IfName        string
	RifId         int32
	RouteDBKeys   interface{}
	DoneChannel   chan int
	RouteManager  interface{}
	EcmpRoute     bool
	RouteEntity   interface{}
	RetHdlrFunc   func(ipInfo *PluginIPRouteInfo, rMgr interface{}, plugin interface{}, ret int)
}

type PluginIPNeighborInfo struct {
	IfIndex       int32
	L2IfIndex     int32
	IpAddr        []uint32
	NeighborFlags uint32
	NextHopFlags  uint32
	NextHopId     uint64
	VlanId        int
	MacAddr       net.HardwareAddr
	IpType        int
	Address       string
	IfName        string
	OperationType int
}

type PluginLagInfo struct {
	IfName     string
	HwId       *int
	HashType   int
	MemberList []int32
}

type PluginUpdateLagInfo struct {
	IfName        string
	HwId          *int
	HashType      int
	OldMemberList []int32
	MemberList    []int32
}

type Inventory struct {
	VendorId   string
	PartNumber string
	RevisionId string
}

type PluginVlanInfo struct {
	VlanId          int
	RsvdVlan        bool
	PortList        []int32
	UntagPortList   []int32
	IfNameList      []string
	UntagIfNameList []string
}

type PluginVlanUpdateInfo struct {
	VlanId                int
	AddPortList           []int32
	DeletePortList        []int32
	AddUntagPortList      []int32
	DeleteUntagPortList   []int32
	AddIfNameList         []string
	DeleteIfNameList      []string
	AddUntagIfNameList    []string
	DeleteUntagIfNameList []string
}
