Use BCM diag shell fp command to create a rule that counts number of packets that ingress from a physical port

There are SIX steps in creating a FP rule.

 

1. Create a qset that describe the qualifiers

2. Create a fp group with this qset

3. Create a fp entry

4. Add qualifier data to this entry, value/mask

5. Add action to this entry, a fp stat in this case

6. Install the entry

 

# step 1 : create qset with qualifier InPort

fp qset clear

fp qset add inport

 

# step 2 : create group 0 with priority 0

fp group create 0 0

 

# step 3 : create fp entry 0 in fp group 0

fp entry create 0 0

 

# step 4 : qualify inport to be xe0 and mask 0xff

fp qual 0 inport xe0 0xff

 

# step 5 : add action to this entry which add a stat to count bytes and packets

fp stat create group=0 type0=packets type1=bytes

fp stat attach entry=0

 

# step 6 : install fp entry

fp entry install 0
