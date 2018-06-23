How to Program L3 host to use L3 Extended Host Entry
     Newer devices such as Trident2 and Triumph3 support L3 extended host entries where the next hop information such as the GLP, MAC address, and interface number, are in the host entry itself, instead of in the ING/EGR next hop tables. In the software, this feature is also referred to as embedded nexthops.

     Programming the L3 extended host entry requires that L3 Egress Mode, e.g. sc L3EgressMode=1, is enabled and the bcm_l3_host_tinformation has:

      

     l3a_intf set to the L3 interface id

     l3a_port_tgid set to the GLP to reach the nexthop

     l3a_nexthop_mac must have the nexthop MAC

      

     In comparison, with non-extended L3 host, the bcm_l3_host_t l3a_intf is set to the L3 egress object id.  Please notice that l3a_intf is of type bcm_if_t. The same bcm_if_t can refer to:

      

     - a L3 interface id, or

     - a L3 egress object id, or

     - a L3 egress multipath object id

      
Examples
      1. Using the BCM diag shell, to program the L3 host without extended L3 host
       

      l3 intf add vlan=3 mac=0x4 intf=3                               <<<< L3 interface: SA, VLAN

      l3 egress add mac=0x3 port=modport(0,2) intf=3                  <<<< L3 egress object: DA, GLP

      l3 l3table add ip=192.168.1.1 intf=$egr_object_id               <<<< L3 host: DIP, L3 egress object id


      To examine the programmed entry

      dump chg l3_entry_ipv4_unicast


      2. Using the BCM diag shell, to program the L3 host with extended L3 host
       

      l3 intf add vlan=3 mac=0x4 intf=3                                <<<< L3 interface: SA, VLAN

      l3 l3table add ip=192.168.1.1 mac=0x3 port=modport(0,2) intf=3   <<<< L3 host: DIP, L3 interface, DA, GLP


      To examine the programmed entry


      dump chg l3_entry_ipv4_multicast

     which shows the IPV4UC_EXT view of the L3_ENTRY table.

        
