#ifndef PLUGIN_MGR_H
#define PLUGIN_MGR_H
#include <stdint.h>
#define LENOF(st) ((sizeof(st))/(sizeof(st[0])))
#define MAX_NUM_PORTS 256
#define MAX_CMD_SIZE 128
#define MAC_ADDR_LEN 6
#define DEFAULT_VLAN_ID 4095
#define CPU_VLAN_ID 4095
#define ALL_INTERFACES -1
#define NETIF_NAME_LEN_MAX 16
#define BOOT_MODE_COLDBOOT 0
#define BOOT_MODE_WARMBOOT 1
#define MAX_PORT_NAME_LEN 16
#define MAX_PORT_DESC_LEN 128
#define MAX_NUM_SUB_PORTS 4
#define MAX_VLAN_ID 4096
#define MAX_IF_TYPE_ENUM 128
#define MAX_MEDIA_TYPE_ENUM 64
#define MAX_DUPLEX_TYPE_ENUM 3
#define NEIGHBOR_TYPE_COPY_TO_CPU 0x1
#define NEIGHBOR_TYPE_BLACKHOLE 0x2
#define NEIGHBOR_TYPE_FULL_SPEC_NEXTHOP 0x4
#define NEIGHBOR_L2_ACCESS_TYPE_PORT 0x8
#define NEIGHBOR_L2_ACCESS_TYPE_LAG 0x10
#define NEIGHBOR_TYPE_VXLAN 0x20
#define NEXTHOP_TYPE_COPY_TO_CPU 0x1
#define NEXTHOP_TYPE_BLACKHOLE 0x2
#define NEXTHOP_TYPE_FULL_SPEC_NEXTHOP 0x4
#define NEXTHOP_L2_ACCESS_TYPE_PORT 0x8
#define NEXTHOP_L2_ACCESS_TYPE_LAG 0x10
#define INTF_STATE_DOWN 0
#define INTF_STATE_UP 1
#define ALL_ZEROS_MAC_ADDR {0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
#define LACP_MAC_ADDR {0x01, 0x80, 0xc2, 0x00, 0x00, 0x02}
#define INVALID_OBJECT_ID 0xFFFFFFFFFFFFFFFF
#define MAC_ENTRY_LEARNED 0x1
#define MAC_ENTRY_AGED 0x2
#define ROUTE_TYPE_CONNECTED             0x00000001
#define ROUTE_TYPE_SINGLEPATH            0x00000002
#define ROUTE_TYPE_MULTIPATH             0x00000004
#define ROUTE_OPERATION_TYPE_UPDATE      0x00000008
#define ROUTE_TYPE_V6                    0x000000010
#define ROUTE_TYPE_NULL                  0x000000020
#define MAX_NEXTHOPS_PER_GROUP 32
#define INTF_TYPE_MASK 0x7f000000
#define INTF_TYPE_SHIFT 24
#define INTF_ID_MASK 0xffffff
#define INTF_ID_SHIFT 0
#define INTF_ID_FROM_IFINDEX(ifIndex) ((ifIndex & INTF_ID_MASK) >> INTF_ID_SHIFT)
#define INTF_TYPE_FROM_IFINDEX(ifIndex) ((ifIndex & INTF_TYPE_MASK) >> INTF_TYPE_SHIFT)
#define IFINDEX_FROM_IF_ID_TYPE(id, type) ((id & INTF_ID_MASK)|((type << INTF_TYPE_SHIFT) & INTF_TYPE_MASK))
#define MAX_NEXTHOPS_PER_ECMP_ROUTE 32
#define VXLAN_UDP_DEST_PORT 4789

//Platform Id Strings
#define PLATFORM_ID_INVALID -1
enum PlatformId{
    PLATFORM_CEL_REDSTONE = 0,
    PLATFORM_LNKD_HAWK,
    PLATFORM_LNKD_OPEN19,
    PLATFORM_CEL_SEASTONE,
    PLATFORM_ACC_AS5712,
    PLATFORM_ACC_AS5812_54T,
    PLATFORM_ACC_WEDGE40,
    PLATFORM_ACC_WEDGE100,
    PLATFORM_ACC_AS7712,
    PLATFORM_ACC_AS6712,
    PLATFORM_ACC_AS6812,
    PLATFORM_CEL_VOYAGER,
    PLATFORM_ACC_AS5912,
    PLATFORM_DELL_S6000
}; 

//Global consts
#define MAX_VENDOR_ID_LEN 32
#define MAX_PART_NUM_LEN 8
#define MAX_REV_ID_LEN 8

#define	PORT_PROTOCOL_ARP 		0x1
#define	PORT_PROTOCOL_DHCP 		0x2
#define	PORT_PROTOCOL_DHCP_RELAY 	0x4
#define	PORT_PROTOCOL_BGP 		0x8
#define	PORT_PROTOCOL_OSPF 		0x10
#define	PORT_PROTOCOL_VXLAN 		0x20
#define	PORT_PROTOCOL_MPLS		0x40
#define	PORT_PROTOCOL_BFD 		0x80

//Port attribute update flags
#define PORT_ATTR_PHY_INTF_TYPE 0x00000001
#define PORT_ATTR_ADMIN_STATE   0x00000002
#define PORT_ATTR_MAC_ADDR      0x00000004
#define PORT_ATTR_SPEED         0x00000008
#define PORT_ATTR_DUPLEX        0x00000010
#define PORT_ATTR_AUTONEG       0x00000020
#define PORT_ATTR_MEDIA_TYPE    0x00000040
#define PORT_ATTR_MTU           0x00000080
#define PORT_ATTR_BREAKOUT_MODE 0x00000100
#define PORT_ATTR_LOOPBACK_MODE 0x00000200
#define PORT_ATTR_ENABLE_FEC    0x00000400
#define PORT_ATTR_TX_PRBS_EN    0x00000800
#define PORT_ATTR_RX_PRBS_EN    0x00001000
#define PORT_ATTR_PRBS_POLY     0x00002000
#define PORT_ATTR_DESCRIPTION   0x00004000
#define PORT_ATTR_PVID          0x00008000

//Port breakout modes
#define PORT_BREAKOUT_MODE_UNSUPPORTED 0
#define PORT_BREAKOUT_MODE_1x40  0x00000001
#define PORT_BREAKOUT_MODE_4x10  0x00000002
#define PORT_BREAKOUT_MODE_1x100 0x00000004
#define PORT_BREAKOUT_MODE_4x25  0x00000008
#define PORT_BREAKOUT_MODE_2x50  0x00000010

//Port loopback modes
#define PORT_LOOPBACK_MODE_NONE 0x1
#define PORT_LOOPBACK_MODE_MAC  0x2
#define PORT_LOOPBACK_MODE_PHY  0x4
#define PORT_LOOPBACK_MODE_RMT  0x8

//Port PRBS polynomial type
#define PRBS_POLY_2POW7   0x1
#define PRBS_POLY_2POW23  0x2
#define PRBS_POLY_2POW31  0x4

// IP values copied from golang syscall package
#define IP_TYPE_IPV4 0x2
#define IP_TYPE_IPV6 0xa

/*Policer rate definitions */
#define COPP_POLICER_HIGH_BW     8000
#define COPP_POLICER_MED_BW      5000
#define COPP_POLICER_LOW_BW      1000
#define COPP_POLICER_HIGH_BURST  5000
#define COPP_POLICER_MED_BURST   1000
#define COPP_POLICER_LOW_BURST   500
#define COPP_MAX_STAT_TYPE       2

/* CoPP constants */
#define COPP_L4_PORT_HTTP                   80
#define COPP_L4_PORT_SSH                    22
#define COPP_L4_PORT_BGP                    179
#define COPP_L4_PORT_BFD                    3784
#define COPP_L4_PORT_MASK                   0xFFFFFFFF

#define COPP_ARP_ETHER_TYPE                 0x0806
#define COPP_LLDP_ETHER_TYPE                0x88cc
#define COPP_STP_ETHER_TYPE                 0x010B
#define COPP_LACP_ETHER_TYPE                0x8809 
#define COPP_L4_PORT_INVALID                0
#define COPP_IP_PROTOCOL_IP_NUMBER_TCP      0x06
#define COPP_IP_PROTOCOL_IP_NUMBER_UDP      0x11
#define COPP_IP_PROTOCOL_IPV4_NUMBER_ICMP   0x01
#define COPP_IP_PROTOCOL_IPV6_NUMBER_ICMP   0x3A
#define COPP_IP_PROTOCOL_IP_NUMBER_OSPFV2   0x59
#define COPP_IP_PROTOCOL_IP_NUMBER_MASK     0xFF
#define COPP_IP_PROTOCOL_INVALID            -1

#define COPP_DST_IP_LOCAL_DATA              0x01
#define COPP_DST_IP_LOCAL_MASK              0x0F

#define COPP_L2_BROADCAST_DEST              {0xFF, 0xFF, 0xFF, \
                                                 0xFF, 0xFF, 0xFF}
#define COPP_L2_ADDR_MASK                   {0xFF, 0xFF, 0xFF, \
                                                 0xFF, 0xFF, 0xFF}
#define COPP_L2_ETHER_TYPE_MASK             0xFFFF

#define COPP_L3_IPV4_BROADCAST_ADDR         "255.255.255.255"
#define COPP_L3_IPV4_MCAST_ADDR             "224.0.0.0"
#define COPP_L3_IPV4_ADDR_MASK              COPP_L3_IPV4_BROADCAST_ADDR
#define COPP_L3_IPV4_MCAST_ADDR_MASK        "240.0.0.0"
#define COPP_L3_IPV6_MCAST_ADDR             {0xFF, 0x00, 0x00, 0x00,\
                                                 0x00, 0x00, 0x00, 0x00,\
                                                 0x00, 0x00, 0x00, 0x00,\
                                                 0x00, 0x00, 0x00, 0x00}
#define COPP_L3_IPV6_MCAST_ADDR_MASK        COPP_L3_IPV6_MCAST_ADDR

#define         ACL_NAME_MAX_LEN 255
/* Acl flags */
#define         ACL_L3_ROUTABLE   0x00000001
#define         ACL_IPV6_ENABLED  0x00000002
#define         ACL_IPV4_ENABLED  0x00000004
#define         ACL_MAC_ENABLED   0x00000008
#define         ACL_VLAN_ENABLED  0x00000010
#define         ACL_PORT_ENABLED   0x00000012
#define         ACL_HITLESS_UPDATE  0x00000014

/*Sub interfaces TYPE*/
#define SECONDARY_INTF 1
#define VIRTUAL_INTF 2

/* STP STATE definitions */
enum stpPortStates {
    StpPortStateBlocking = 0,
    StpPortStateLearning,
    StpPortStateForwarding,
    StpPortStateCount
};

enum hashTypes {
    HASHTYPE_SRCMAC_DSTMAC = 0,
    HASHTYPE_SRCIP_DSTIP = 6,
    HASHTYPE_END
};
enum mediaType {
    MediaTypeCount
};
enum duplexType {
    HalfDuplex = 0,
    FullDuplex,
    DuplexCount
};
enum portIfType {
    PortIfTypeMII,
    PortIfTypeGMII,
    PortIfTypeSGMII,
    PortIfTypeQSGMII,
    PortIfTypeSFI,
    PortIfTypeXFI,
    PortIfTypeXAUI,
    PortIfTypeXLAUI,
    PortIfTypeRXAUI,
    PortIfTypeCR,
    PortIfTypeCR2,
    PortIfTypeCR4,
    PortIfTypeKR,
    PortIfTypeKR2,
    PortIfTypeKR4,
    PortIfTypeSR,
    PortIfTypeSR2,
    PortIfTypeSR4,
    PortIfTypeSR10,
    PortIfTypeLR,
    PortIfTypeLR4,
    PortIfTypeCount
};
enum portStatTypes {
    IfInOctets = 0,
    IfInUcastPkts,
    IfInDiscards,
    IfInErrors,
    IfInUnknownProtos,
    IfOutOctets,
    IfOutUcastPkts,
    IfOutDiscards,
    IfOutErrors,
	IfEtherUnderSizePktCnt,
	IfEtherOverSizePktCnt,
	IfEtherFragments,
	IfEtherCRCAlignError,
	IfEtherJabber,
	IfEtherPkts,
	IfEtherMCPkts,
	IfEtherBcastPkts,
	IfEtherPkts64OrLessOctets,
	IfEtherPkts65To127Octets,
	IfEtherPkts128To255Octets,
	IfEtherPkts256To511Octets,
	IfEtherPkts512To1023Octets,
	IfEtherPkts1024To1518Octets,
    portStatTypesMax
};

enum BufferPortStatTpes {
EgressPort = 0,
PortBufferStat,
bufferPortStatTypesMax
};

enum BufferGlobalStatTypes{
BufferStat = 0,
EgressBufferStat,
IngressBufferStat,
bufferGlobalStatTypesMax
};

typedef struct ifMapInfo_s {
    char *ifName;
    int port;
} ifMapInfo_t;

typedef struct portConfig_s {
    uint8_t portNum;
    enum portIfType ifType;
    uint8_t adminState;
    uint8_t macAddr[MAC_ADDR_LEN];
    uint32_t portSpeed;
    enum duplexType duplex;
    uint8_t autoneg;
    enum mediaType media;
    int mtu;
    uint32_t breakOutMode;
    uint8_t isBreakOutSupported;
    uint8_t mappedToHw;
    uint8_t logicalPortInfo;
    uint8_t loopbackMode;
    uint8_t enableFEC;
    uint8_t txPrbsEn;
    uint8_t rxPrbsEn;
    uint8_t prbsPoly;
    uint8_t dirty;
    uint8_t sfpId;
} portConfig;
typedef struct portState_s {
    uint8_t portNum;
    uint8_t operState;
    char portName[MAX_PORT_NAME_LEN];
    uint64_t stats[portStatTypesMax];
    uint64_t prbsErrCnt;
} portState;
typedef struct bufferPortState_s {
    uint8_t portNum;
    char portName[MAX_PORT_NAME_LEN];
    uint64_t stats[bufferPortStatTypesMax]; 
} bufferPortState;
typedef struct bufferGlobalState_s {
	uint32_t deviceId;
	uint64_t stats[bufferGlobalStatTypesMax];
}bufferGlobalState;
typedef struct portBreakOutInfo_s {
    uint8_t logPortNum;
    uint8_t parentPort;
    uint32_t breakOutMode;
    uint8_t totalNumSubPorts;
} portBreakOutInfo_t;

enum coppClassTypes{
	CoppArpUC = 0,
	CoppArpMC,
	CoppBgp,
	CoppIcmpV4UC,
	CoppIcmpV4BC,
	CoppStp,
	CoppLacp,
	CoppBfd,
	CoppIcmpv6,  
	CoppLldp,
	//CoppSsh,
	//CoppHttp,
	CoppClassCount
};

typedef struct coppStatState_s {
	int     coppClass;
	int     peakRate;
	int 	burstRate;
	uint64_t greenPackets;
	uint64_t redPackets;
} coppStatState_t;

enum aclAction{
	AclDeny = 0,
	AclAllow
};
enum aclDirection{
	AclIn = 0,
	AclOut
};

enum aclType{
	AclIp = 0,
	AclMac,
	AclSVI,
	AclMlag // used by internal daemon
};

enum aclL4PortMatch {
	AclEq = 0,
	AclNeq,
        AclRange
};

typedef struct aclIPv4Filter_t {
        uint8_t * sourceIp;
        uint8_t* destIp;
        uint8_t* sourceMask;
        uint8_t* destMask;
        int   proto;
        int   srcintf;
        int   dstintf;
        int   l4SrcPort;
        int   l4DstPort;
        int   l4PortMatch;
        int   l4MinPort;
        int   l4MaxPort;
}aclIPv4Filter;

typedef struct aclMACFilter_t {
        uint8_t* sourceMac;
        uint8_t* destMac;
        uint8_t* sourceMask;
        uint8_t* destMask;
}aclMACFilter;

typedef struct aclIPv6Filter_t {
        uint8_t * sourceIpv6;
        uint8_t* destIpv6;
        uint8_t* sourceMaskv6;
        uint8_t* destMaskv6;
	int   proto;
        int   l4SrcPort;
        int   l4DstPort;
        int   l4PortMatch;
        int   l4MinPort;
        int   l4MaxPort;
}aclIPv6Filter;

typedef struct filter_t {
	aclIPv4Filter filterIPv4;
	aclMACFilter filterMac;
 	aclIPv6Filter filterIPv6;
}filter;

typedef struct acl_t {
        char *aclName;
	int  aclLen;
	int aclType;
        int priority;
	int   action;
	int direction;
	int *portList;
	int portListLen;
	int *vlanList;
	int vlanListLen;
	filter filterConf;
}acl;

typedef struct aclConfig {
	char *aclName;
	uint64_t hwIndex;
	int statId;
}aclConfig;


typedef struct aclState_s {
	char* aclName;
	uint64_t stat;
} aclState;

typedef struct copp_s {
	int protoId;
	int policerPeakRate;
	int policerBurstRate;
	int cpuQueue;
}copp;

#endif //PLUGIN_MGR_H
