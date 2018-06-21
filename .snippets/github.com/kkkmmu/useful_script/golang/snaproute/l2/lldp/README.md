Package lldpServer implement IEEE 802.1AB Link Layer Discovery Protocol.
## Architecture
![alt text](https://github.com/SnapRoute/l2/blob/master/lldp/docs/LLDP_Design.png "Architecture")

## Support
 - Enable/Disable LLDP per interface/port
 - Chassis Id TLV
 - Port Id TLV
 - TTL Tlv
 - System Description TLV
 - Hostname TLV
 - Managment Address (subtype IPv4 Address) TLV
 - Marshalling/Un-Marshalling of all above TLV's

##Future Work
 - User based configuration for Optional TLV's.
 - Statistics for packet rx/tx
 - Chassis Id TLV
 - Port Id TLV
 - TTL Tlv
 - Marshalling/Un-Marshalling of Mandatory TLV's

