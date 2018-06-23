How to Configure the LPM Table for IPv4 and IPv6 Combinations

This article describes on how to configure the LPM table for IPv4 and IPv6 combinations with different scenarios.

Based on configuration property variables, the LPM Table scaling varies conditionally.

The example below illustrates and explains the different combinations and scenarios (IPv4 and IPv6) using the with BCM56640 (Triumph 3) device. 

 

Step #1: Create VLAN, L3 Interface, and L3 egress object using the following Broadcom Diag shell commands:

sc l3egressmode=0x1

vlan create 2 pbm=ge

l3 intf add vlan=2 mac=00:00:00:00:00:02

l3 egress add mac=0:0:2:2:2:2 intf=0 port=ge0 mod=0

 

Default Case #1 (8k=>IPv4-32b, 2k=>IPv6-128b)

    IPv4 32-bit:

    for i=1,8192 'l3 defip add ip=$i mask=255.255.255.255 intf=100002'

    l3 defip add ip=8193 mask=255.255.255.255 intf=100002

    for i=1,15360 'l3 defip add ip=$i mask=255.255.255.255 intf=100002'

     

    IPv6 64-bit:

    for i=1,7680 'l3 ip6route add vrf=0 ip=3001:0:$i:0:0:0:0:0 masklen=64 intf=100002'

    for i=1,7168 'l3 ip6route add vrf=0 ip=3001:0:$i:0:0:0:0:0 masklen=64 intf=100002'

     

    IPv6b 128-bit:

    for i=1,2048 'l3 ip6route add vrf=0 ip=3001:0:0:0:0:$i:0:0 masklen=128 intf=100002'

    l3 ip6route add vrf=0 ip=3001:0:0:0:0:2049:0:0 masklen=128 intf=100002

     

    User Config Case #2: 

    Scenario A

    config add lpm_scaling_enable=1

    config add ipv6_lpm_128b_enable=1

    config add lpm_ipv6_128b_reserved=0

    config add num_ipv6_lpm_128b_entries=256

     

    sc l3egressmode=0x1

    vlan create 2 pbm=ge

    l3 intf add vlan=2 mac=00:00:00:00:00:02

    l3 egress add mac=0:0:2:2:2:2 intf=0 port=ge0 mod=0 

    for i=1,256 'l3 ip6route add vrf=0 ip=3001:0:0:0:0:$i:0:0 masklen=128 intf=100002'

    for i=1,7168 'l3 ip6route add vrf=0 ip=3001:0:$i:0:0:0:0:0 masklen=64 intf=100002'

    l3 ip6route add vrf=0 ip=3001:0:7169:0:0:0:0:0 masklen=64 intf=100002

     

    Scenario B

    config add lpm_scaling_enable=1

    config add ipv6_lpm_128b_enable=1

    config add lpm_ipv6_128b_reserved=1

    config add num_ipv6_lpm_128b_entries=256

     

    sc l3egressmode=0x1

    vlan create 2 pbm=ge

    l3 intf add vlan=2 mac=00:00:00:00:00:02

    l3 egress add mac=0:0:2:2:2:2 intf=0 port=ge0 mod=0

     

    IPv4 and IPv6 Entries and Combination Calculations:

    => (8K-2*ipv6128)*2 => IPv4 32 bit entries or (8k-2*ipv6128) => IPV6-64 bit

    -> (8000-512) * 2 => 15360 IPv4, (8000-512) => 7680 IPv6 64 bit

     

    Test Case #A IPv6 128-bit and IPv4 32-bit: 256 (Partition1) and 15360 (Partition2, Partition1)

    for i=1,256 'l3 ip6route add vrf=0 ip=3001:0:0:0:0:$i:0:0 masklen=128 intf=100002'

    for i=1,15360 'l3 defip add ip=$i mask=255.255.255.255 intf=100002'

    l3 defip add ip=15361 mask=255.255.255.255 intf=100002

     

    BCM.0> l3 defip add ip=15361 mask=255.255.255.255 intf=100002

    L3: Error adding route table: Table full 

     

    Test Case #B IPv6 128-bit and IPv6 64-bit test: 256 (Partition1) and  6912 (Partition2) 

    As per programmer’s guide (56640-PG105-RDS.pdf), Table 15 illustrates the different

    supported configurations and the usage of Partition1 and Partition2 as per the configuration.

     

    Actual Partition1 size = 1024

    Configured entries =  256

    Remaining Partition1 size = 1024-256 = 768

     

    Actual Partition2 size = 1024*6 = 6144

    Partition2 is used for IPv6 64-bit entries and Partition1 is meant for IPv6 128-bit size entries. However, we can use this space for IPv6 64-b-t entries but at a cost of entries being unused.

     

    So total IPv6 64-bit entries = 6144 + 768 = 6912 

    for i=1,256 'l3 ip6route add vrf=0 ip=3001:0:0:0:0:$i:0:0 masklen=128 intf=100002'

    for i=1,6912 'l3 ip6route add vrf=0 ip=3001:0:$i:0:0:0:0:0 masklen=64 intf=100002'

    l3 ip6route add vrf=0 ip=3001:0:6913:0:0:0:0:0 masklen=64 intf=100002
