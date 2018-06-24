# Virtual Router Redundancy Protocol

### Introduction
This module implement Virtual Router Redundancy Protocol RFC 5798

### Architecture

                             +---------------+
                             |               |
                             |  User Config  |
                             |               |
                             +---------------+
                                     |         
                                     |
                                     V
                             +---------------+
                  +--------->|               |<-------------+
                  |          |  Initialize   |              |
                  |   +------|               |----------+   |
                  |   |      +---------------+          |   |
                  |   |                                 |   |
                  |   V                                 V   |
          +---------------+                       +---------------+
          |               |---------------------->|               |
          |    Master     |                       |    Backup     |
          |               |<----------------------|               |
          +---------------+                       +---------------+

### Interfaces
 - Create/Delete Virtual Router
 - Change timers for VRRP packet, for e.g: Advertisement Timer

### Configuration
 - VRRP configuration is based of https://tools.ietf.org/html/rfc5798#section-5.2
 - Unless specified each instance of Virtual Router will use the default values specified in the RFC
