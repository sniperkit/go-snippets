1. Basic: 
    There are multiple of tables within the device. These tables are used as part of the switching criteria. For example, L2 tables are used to lookup packet L2 header. The VLAN_Translation table is used to lookup incoming VLAN and translate to another. How the device searches the entry within each table is depending on the table lookup algorithm. There are three types of search algorithms:
    1. Hash Search Lookup.
    2. Direct Index Lookup.
    3. TCAM Search Lookup.

    Hash Search Lookup:
        The Hash Search Engine Lookup uses the input key as a searching criteria into the tables. This include the L2 MAC table, L3 host table, MAC VLAN, IPMC searchs(s, g) and (*, g), and ECMP table for Reverse Path Forwarding(RPF) checks. Ingress VLAN translation, and egress VLAN translation.
        The hasing key can be selected from the HASH_CONTROL register. For the L2 or L3 table, users can choose hashing algorithm such as upper or lower 16-bit or 32-bit CRC. When the hashing is computed, it returens an index into the bucket. Each bucket has eight entries. In the read operation, the hashed value gets indexed into the bucket. The eight entries within that bucket are then read and compared. In the write operation, the first available entry in that bucket is written. The process continues until all entries are full. Hash tables are most efficient and cost effective on a given size. Howere, because of the nature of hashing, not all entries in a table can be 100% utilized. For example, two different addresses may hashed into the same bucket index. A bucket index may happen to have space for one last address entry, so leaving the second address to be dropped.
        The devies supports a dual-hashing algorithm. This algorithm increases the efficiency of the hashing table. In dual-hashing, the L2/L3 table is divided into two buckets. Each bucket has four entires, instead of eight. Dual hashing for L2_ENTRY table is enabled via L2_AUX_HASH_CONTROL register. For L3 lookup table, it is the L3_AUX_HASH_CONTROL register. When dual-hashing is enabled, uses can configure the hash mode for each bucket. For example, the first bucket can use CRC32_UPPER hashing algorithm, when the second bucket can use CRC32_LOWER hashing algorithm. The result is two indices into two buckets.
        To effectively use dual-hasing mode, software must be able to keep track with the hashing values and rearrange address entries when needed. For example, if the hasing result of address X belongs to buckets A and D, and address Y belongs to buckets A and C, but buckets A and C are full, then software must determine if entries in bucket D is Available. If entry in bucket D is available, address X is then moved to bucket D, so that bucket A is now available for address Y.
    
    Direct Index Lookup.
        The Direct Index Lookup uses a direct search method. Based on the search index, it directly goes to the entry within that table to perform a read/write operation. Tables that support direct index lookup include the VLAN and PORT tables.

    TCAM Search Lookup.
        The CAM search lookup examines every entry of a table simultanenously and results in the first entry that mathes the lookup key. The search bits on a TCAM entry can be masked ofr wild search. Some of the tables that are TCAM-based are VCAP, ICAP, and ECAP. A TCAM look-up can be triggered by any one of the following when using the ContenAware engine:
            1. VRF/default routing/Policy routing/Ingress ACL
            2. Egress ACL
            3. HiGig ACL
        While the CAM search lookup performs a search, A TCAM is used to parse the contents of the packet and set its acitons when a match is made. The rules are organized into slices. Each slice has a field selector, a 256-entry TCAM key or rule, a 256-entry policy table that contains the action fields corresponding to each of the CAM rules, and a 256  meter and traffic counter.

    这里要注意：拿到一个表的名字以后要先确定该表的查找方法（HASH/Direct Index/TCAM）。HASH表由于其本质，冲突是不可避免的，因此所有硬件表项不可避免的不能得到充分利用。
    测试中经常看到L2 table和L3_ENTRY硬件上其实还有空表项但是SDK调用返回表已满，这就是由于其都是HASH表引起的。
    

2. BCM_XGS3_L3_DEFIP_IP4_CNT(unit)
   BCM_XGS3_L3_DEFIP_IP6_CNT(unit)
   BCM_DEFIP_PAIR128_USED_COUNT(unit)
   BCM_XGS3_L3_IP4_CNT(unit)
   BCM_XGS3_L3_IP4_IPMC_CNT(unit)
   BCM_XGS3_L3_IP6_CNT(unit)
   BCM_XGS3_L3_IP6_IPMC_CNT(unit)
   BCM_XGS3_L3_ECMP_MAX_PATHS(unit)
   BCM_XGS3_L3_IF_COUNT(unit)
   BCM_XGS3_L3_DEFIP_TBL_SIZE(unit)
   BCM_XGS3_L3_ECMP_MAX_GROUPS(unit)
   BCM_XGS3_L3_TBL_SIZE(unit)
   BCM_XGS3_L3_ECMP_IN_USE(unit)

3. L3_ENTRY table supports two types of view for providing the associated data:
    Traditional view
        If the address is found, the entry outputs the NEXT_HOP_INDEX to point to INITIAL_ING_L3_NEXT_HOP, ING_L3_NEXT_HOP and EGR_L3_NEXT_HOP tables to provide the associated data.
    Extended view
        If the address is found, the entry outputs the associated data, such as the next-hop MAC address, destination port, L3_INTF_NUM and so on, for routing purpose directly.
