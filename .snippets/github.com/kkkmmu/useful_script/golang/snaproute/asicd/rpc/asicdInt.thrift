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
// _______   __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __  
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  | 
// |  |__   |  |     |  |__   \  V  /     |   (----  \   \/    \/   /  |  |  ---|  |---- |  ,---- |  |__|  | 
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   | 
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  | 
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__| 
//                                                                                                           

namespace go asicdInt
typedef i32 int
typedef i16 uint16

// Struct for Configuring Reserved Mac Addr on Chip
// so that packets are punted to CPU when enabled
struct RsvdProtocolMacConfig {
        1: string MacAddr
        2: string MacAddrMask
        3: i32 VlanId
}
struct Lag {
    1:i32 LagIfIndex
    2:i32 HashType
    3:list<i32> IfIndexList
    4:string LagName
}
struct LagGetInfo {
    1: int StartIdx
    2: int EndIdx
    3: int Count
    4: bool More
    5: list<Lag> LagList
}

struct Vtep {
	1 : i32 IfIndex
	2 : string IfName
	3 : i32 Vni
	4 : i32 SrcIfIndex
	5 : string SrcIfName
	6 : string SrcMac
	7 : string DstIp
	8 : string SrcIp
	9 : i16 VlanId
	10 : i16 UDP
	11 : i16 TTL
	12 : i16 DSCP
	13 : i32 NextHopIfIndex
	14 : string NextHopIfName
	15 : string NextHopIp 
	16 : bool Learning
	17 : i32 MTU
}
struct Vxlan {
	1: i32 Vni
	2: list<i16> UntaggedVlanId
	3: list<i16> VlanId
}

struct IPv4NextHop {
    1: string NextHopIp
    2: i32 Weight
    3: i32 NextHopIfType
}
struct IPv6NextHop {
    1: string NextHopIp
    2: i32 Weight
    3: i32 NextHopIfType
}
struct IPv4Route {
    1: string destinationNw
    2: string networkMask
    3: list<IPv4NextHop> NextHopList
}
struct IPv6Route {
    1: string destinationNw
    2: string networkMask
    3: list<IPv6NextHop> NextHopList
}
struct Vlan {
	1 : i32 VlanId
	2 : list<i32> IfIndexList
	3 : list<i32> UntagIfIndexList
    4:string VlanName
}
struct VlanGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Vlan> VlanList
}
struct Intf {
    1: string IfName
    2: i32 IfIndex
}
struct IntfGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Intf> IntfList
}
struct IPRouteHwState {
	1 : string DestinationNw
	2 : string NextHopIps
	3 : string RouteCreatedTime
	4 : string RouteUpdatedTime
}
struct IPRouteHwStateGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<IPRouteHwState> IPRouteHwStateList
}

struct Port {
	1 : string IntfRef
	2 : i32 IfIndex
	3 : string Description
	4 : string PhyIntfType
	5 : string AdminState
	6 : string MacAddr
	7 : i32 Speed
	8 : string Duplex
	9 : string Autoneg
	10 : string MediaType
	11 : i32 Mtu
	12 : string BreakOutMode
	13 : string LoopbackMode
	14 : bool EnableFEC
	15 : bool PRBSTxEnable
	16 : bool PRBSRxEnable
	17 : string PRBSPolynomial
}

struct SFlowIntfBulkInfo {
	1: list<SFlowIntfInfo> SFlowIntfInfoList
}

struct SFlowIntfInfo {
	1: string IntfRef
    2: string NetDevName
	3: i32 PortIfIdx
}

service ASICDINTServices {
    // All services listed here are utilities provided to other daemons. These are hand-written and not auto genereated.

    //Port
    list<Port> GetAllPortsWithDirtyCache()

    //Vlan
    VlanGetInfo GetBulkVlan(1: int fromIndex, 2: int count);

    //Intf
    IntfGetInfo GetBulkIntf(1: int fromIndex, 2: int count);

    //STP
    i32 CreateStg(1:list<i32> vlanList);
    bool DeleteStg(1:i32 stgId);
    bool SetPortStpState(1:i32 stgId, 2:i32 port, 3:i32 stpState);
    i32 GetPortStpState(1:i32 stgId, 2:i32 port);
    bool UpdateStgVlanList(1:i32 stgId, 2:list<i32> vlanList);
    oneway void FlushFdbStgGroup(1:i32 stgId, 2:i32 port);

    //LAG
    i32 CreateLag(1:string ifName, 2:i32 hashType, 3:string ifIndexList);
    i32 DeleteLag(1:i32 ifIndex);
    i32 UpdateLag(1:i32 ifIndex, 2:i32 hashType, 3:string ifIndexList);
    LagGetInfo GetBulkLag(1:int fromIndex, 2:int count);
    bool CreateLagCfgIntfList(1:string ifName, 2:list<i32> ifIndexList);
    bool UpdateLagCfgIntfList(1:string ifName, 2:list<i32> ifIndexList);
    bool DeleteLagCfgIntfList(1:string ifName, 2:list<i32> ifIndexList);

    //IPv4 neighbors
    i32 CreateIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 UpdateIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 DeleteIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);

    //IPv4 routes
    oneway void OnewayCreateIPv4Route(1:list<IPv4Route> ipv4RouteList);
    oneway void OnewayDeleteIPv4Route(1:list<IPv4Route> ipv4RouteList);

    //IPv6 neighbors
    i32 CreateIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 UpdateIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 DeleteIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    
    //IPv6 routes
    oneway void OnewayCreateIPv6Route(1:list<IPv6Route> ipv6RouteList);
    oneway void OnewayDeleteIPv6Route(1:list<IPv6Route> ipv6RouteList);
	IPRouteHwStateGetInfo GetBulkIPRouteHwState(1: int fromIndex, 2: int count);
	IPRouteHwState GetIPRouteHwState(1: string DestinationNw);

    //Protocol Mac Addr
    bool EnablePacketReception(1:RsvdProtocolMacConfig config);
    bool DisablePacketReception(1:RsvdProtocolMacConfig config);
	
    //Err-disable	
    bool ErrorDisablePort(1: i32 ifIndex, 2:string AdminState, 3:string ErrDisableReason)
	
    // auto generated for vxland, and copied here as config is proxied
    // so if interface changes must update here as well.  Only difference is
    // the asicd tends to return id values and the function names changed 
    i32 CreateVxlanVtep(1: Vtep config);
    bool DeleteVxlanVtep(1: Vtep config);
    bool UpdateVxlanVtepAttr(1: string vtepName, 2: i32 vni, 3: i16 tos, 4: i16 ttl, 5: i32 mtu);
    bool CreateVxlan(1: Vxlan config);
    bool DeleteVxlan(1: Vxlan config);
    bool UpdateVxlan(1: i32 vni, 2:list<i16>addvlanlist, 3:list<i16>delvlanlist, 4:list<i16>adduntaggedVlan, 5:list<i16>deluntaggedVlan);	
    bool LearnFdbVtep(1:string mac, 2:string vtep, 3:i32 ifindex);

    // Update Virtual IP
    bool UpdateVirtualIPv4Intf(1: string intRef, 2: string Type, 3:string IpAddr, 4:string MacAddr, 5: bool Enable);
    bool UpdateVirtualIPv6Intf(1: string intRef, 2: string Type, 3:string IpAddr, 4:string MacAddr, 5: bool Enable);

   //SFlow
   bool EnableSFlowSampling(1:i32 PortIfIdx, 2:bool Enable)
   bool SetSFlowSamplingRate(1:i32 PortIfIdx, 2:i32 Rate)
   SFlowIntfBulkInfo GetBulkSFlowIntfInfo()
   bool IsLinuxOnlyPlugin()
}
