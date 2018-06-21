vxlan package can be run in one of two modes

Config:
"vxlanmode" - "linux" or "proxy"
1) "linux" - Using linux vxlan interface (NOT SUPPORTED)
2) "proxy" - Using vxland vxlan proxy

"standalone" - true or false
1) true - run without asicd, all linux configurations will be done within vxland
2) false - run with asicd

Case 2 will perform the necessary encap and decap of the vxlan packets.
Users must specify in the vxlan.conf file which version they would like to run.

Both version will create a linux interface called vtep<vtepid>
Case 1 - Vxlan interface
Case 2 - vEth interface, where the peer interface will be vtep<vtepid>Int only to be used
         by vxland.  However no restrictions are made to prevent users from doing anything
		with this interface. 
		
This will be added to the associated VXLAN (VNI) bridge.
=======
# Virtual Extensible LAN (VxLan)

### Introduction
Brief Introduction about the module

### Architecture
Pictorial Representation of module architecture. Flow diagram

### Interfaces
Exposed Interfaces

### Configuration
Location of configuration and expected entries in configuration file
