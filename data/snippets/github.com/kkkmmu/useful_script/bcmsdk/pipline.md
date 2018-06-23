1. Tunnel Termination Lookup.
    The purpose of this block is to determine what kinds of forwarding lookups and processing should be performed for a given packet. The lookup table is TCAM-based table, where the key presented to the table is composed of the destination MAC address, the VLAN ID/VFI derived in the previous stages, and the incoming source interface(port/virtual port). When the lookup in the TCAM results in a hit, additional bits are retrived that indicate what kind of lookup and furture processing the packet is eligible for. This includes lookups and processing for:
        1. IPv4/IPv6 unicast flows.
        2. IPv4/IPv6 multicast flows.
        3. MPLS/VPLS/VPWS/ flows.
        4. MAC-in-MAC flows.
        5. VXLAN/L2GRE/NSM.
        6. FCoE flows.
    The type of processing allowed is also qualifed by the packet's EtherType to ensure the packet is of the correct format for the indicated type of processing. For example, even if the termination hits and indicates a packet is eligible for MPLS procesing, the logi does not allow this unless the packet's EtherType is 0x8847. In the event that the TCAM does not have an entry that hits, or the entry that his does not indicate that any kind of additional lookup are allowed, the packet is only eligible for standard Layer-2 lookups and forwarding.
2. Tunnel Adaption
    Once the packet is qualified for further processing by the Tunnel termination lookups, lookups are performed to derive the forwarding attributes such as SVP, L3_IIF, or VFI.

3. VLAN Membership Check.
    Once this block in the pipeline is reached, the packet is now associated with a VLAN/VFI based on the results from the previous blocks in the pipeline. The purpose of this block is to ensure that the packet is associated with a valid VLAN/VFI and that the ingress port/SVP is a member of the VLAN/VFI.
    This block is alos used to provide a bitmap of the ports associated with the LVAN in the event that the packet must be flooded to the VLAN. Finally, for the ports the packet should be forwarded to, a spanning tree state check is done to ensure that the port is a forwarding state.
    There are two configurations that can potentially allow a packet to be switched, even though the ingress port was not a member of the VLAN/VFI.
        1. The first is a per port/virtual port configuation bit that can be used to disable VLAN membership checks in the ingress pipeline. The allow a packet to arrive on a port that is not associated with the VLAN, yet still be switched to a port that is a member of the VLAN.
        2. The other way to bypass the VLAN membership check is via the VCAP or VLAN translation stages; where there is an explict control indication that membership checks should not be performed. The VLAN membership bitmap is obtained through the VLAN_ATTRS_2 table, which is indexed by the VLAN ID assignment in the previous stages. If the membership check is enabled, this bitmap is used to ensure that the ingress-port is a member of the VLAN. This bitmap is used in the event that the packet cannot be resovled to a single switching destination and must be flooded to entire VLAN.


4. Forwarding Inforamtion Database.
    The forwarding database lookup block is used to resovle the packet to a know forwarding destination. The type of resolution attempted is based on the outputs from the adaption lookups. Based on this, one of the following database lookup is performed:
        1. Layer 2 database.
        2. Host/multicast database.
        3. Route lookup.
    In the case of a layer 2 lookup, a key is formed using the destination MAC address as well as the VLAN ID/VFI. A lookup is then done using a hashed version of the key to determin if the destination MAC address is known. If it is, the lookup returns a hit along with the logical destination. The destination is considered logical, because the destination may be a link aggregation group, a virtual port, or a multicast group. At the same time that the destination MAC lookup is done, a source MAC lookup is done to determin if the adddress is known. If the source MAC address is unknown, the system can be configured to react in a number of ways on the new unknown MAC address. 
    When the previous block indicates the packet is eligible for IPv4/IPv6 unicast processing, a key is formed using the destination IP address, as well as a Virtual route forwarind(VRF) identifier that can be used to segreage the L3 database into different virtual databases. Host lookup is then done using a hashed version of the key to determin if there is a match. If the lookup results in a miss, the a subsequent route lookup is done for a default route. If the route lookup is also a miss, the packet can be dropped or copied to the CPU. If the packet hist a route or host lookup, the the associated entry returns either a nexthop index(if there is only a single path throught the network to the host or router) or an equal-cost multipath(ECMP) group nubmer(if there are multiple paths to the host router). The next hop index is ued in a later stage to retrieves attributes related to the destination host. The ECMP group is used to retrieve a list of next hop indices that are resoved to a single nexthop index via a hashing function.
