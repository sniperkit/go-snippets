//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package dhcp4

import (
	"time"
)

// DHCP Packet global constants
const (
        DHCP_PACKET_MIN_SIZE    = 272
        DHCP_PACKET_HEADER_SIZE = 16
        DHCP_PACKET_MIN_BYTES   = 240
        DHCP_SERVER_PORT        = 67
        DHCP_CLIENT_PORT        = 68
        DHCP_BROADCAST_IP       = "255.255.255.255"
        DHCP_NO_IP              = "0.0.0.0"
        DHCP_REDIS_DB_PORT      = ":6379"
)

const (
	snapshot_len   int32         = 1024
	promiscuous    bool          = false
	timeout        time.Duration = 30 * time.Second
)

// DHCP Available Options enum type.... This will cover most of the options type
const (
        End                          DhcpOptionCode = 255
        Pad                          DhcpOptionCode = 0
        OptionSubnetMask             DhcpOptionCode = 1
        OptionTimeOffset             DhcpOptionCode = 2
        OptionRouter                 DhcpOptionCode = 3
        OptionTimeServer             DhcpOptionCode = 4
        OptionNameServer             DhcpOptionCode = 5
        OptionDomainNameServer       DhcpOptionCode = 6
        OptionLogServer              DhcpOptionCode = 7
        OptionCookieServer           DhcpOptionCode = 8
        OptionLPRServer              DhcpOptionCode = 9
        OptionImpressServer          DhcpOptionCode = 10
        OptionResourceLocationServer DhcpOptionCode = 11
        OptionHostName               DhcpOptionCode = 12
        OptionBootFileSize           DhcpOptionCode = 13
        OptionMeritDumpFile          DhcpOptionCode = 14
        OptionDomainName             DhcpOptionCode = 15
        OptionSwapServer             DhcpOptionCode = 16
        OptionRootPath               DhcpOptionCode = 17
        OptionExtensionsPath         DhcpOptionCode = 18

        // IP Layer Parameters per Host
        OptionIPForwardingEnableDisable          DhcpOptionCode = 19
        OptionNonLocalSourceRoutingEnableDisable DhcpOptionCode = 20
        OptionPolicyFilter                       DhcpOptionCode = 21
        OptionMaximumDatagramReassemblySize      DhcpOptionCode = 22
        OptionDefaultIPTimeToLive                DhcpOptionCode = 23
        OptionPathMTUAgingTimeout                DhcpOptionCode = 24
        OptionPathMTUPlateauTable                DhcpOptionCode = 25

        // IP Layer Parameters per Interface
        OptionInterfaceMTU              DhcpOptionCode = 26
        OptionAllSubnetsAreLocal        DhcpOptionCode = 27
        OptionBroadcastAddress          DhcpOptionCode = 28
        OptionPerformMaskDiscovery      DhcpOptionCode = 29
        OptionMaskSupplier              DhcpOptionCode = 30
        OptionPerformRouterDiscovery    DhcpOptionCode = 31
        OptionRouterSolicitationAddress DhcpOptionCode = 32
        OptionStaticRoute               DhcpOptionCode = 33

        // Link Layer Parameters per Interface
        OptionTrailerEncapsulation  DhcpOptionCode = 34
        OptionARPCacheTimeout       DhcpOptionCode = 35
        OptionEthernetEncapsulation DhcpOptionCode = 36

        // TCP Parameters
        OptionTCPDefaultTTL        DhcpOptionCode = 37
        OptionTCPKeepaliveInterval DhcpOptionCode = 38
        OptionTCPKeepaliveGarbage  DhcpOptionCode = 39

        // Application and Service Parameters
        OptionNetworkInformationServiceDomain            DhcpOptionCode = 40
        OptionNetworkInformationServers                  DhcpOptionCode = 41
        OptionNetworkTimeProtocolServers                 DhcpOptionCode = 42
        OptionVendorSpecificInformation                  DhcpOptionCode = 43
        OptionNetBIOSOverTCPIPNameServer                 DhcpOptionCode = 44
        OptionNetBIOSOverTCPIPDatagramDistributionServer DhcpOptionCode = 45
        OptionNetBIOSOverTCPIPNodeType                   DhcpOptionCode = 46
        OptionNetBIOSOverTCPIPScope                      DhcpOptionCode = 47
        OptionXWindowSystemFontServer                    DhcpOptionCode = 48
        OptionXWindowSystemDisplayManager                DhcpOptionCode = 49
        OptionNetworkInformationServicePlusDomain        DhcpOptionCode = 64
        OptionNetworkInformationServicePlusServers       DhcpOptionCode = 65
        OptionMobileIPHomeAgent                          DhcpOptionCode = 68
        OptionSimpleMailTransportProtocol                DhcpOptionCode = 69
        OptionPostOfficeProtocolServer                   DhcpOptionCode = 70
        OptionNetworkNewsTransportProtocol               DhcpOptionCode = 71
        OptionDefaultWorldWideWebServer                  DhcpOptionCode = 72
        OptionDefaultFingerServer                        DhcpOptionCode = 73
        OptionDefaultInternetRelayChatServer             DhcpOptionCode = 74
        OptionStreetTalkServer                           DhcpOptionCode = 75
        OptionStreetTalkDirectoryAssistance              DhcpOptionCode = 76

        OptionRelayAgentInformation DhcpOptionCode = 82

        // DHCP Extensions
        OptionRequestedIPAddress     DhcpOptionCode = 50
        OptionIPAddressLeaseTime     DhcpOptionCode = 51
        OptionOverload               DhcpOptionCode = 52
        OptionDHCPMessageType        DhcpOptionCode = 53
        OptionServerIdentifier       DhcpOptionCode = 54
        OptionParameterRequestList   DhcpOptionCode = 55
        OptionMessage                DhcpOptionCode = 56
        OptionMaximumDHCPMessageSize DhcpOptionCode = 57
        OptionRenewalTimeValue       DhcpOptionCode = 58
        OptionRebindingTimeValue     DhcpOptionCode = 59
        OptionVendorClassIdentifier  DhcpOptionCode = 60
        OptionClientIdentifier       DhcpOptionCode = 61

        OptionTFTPServerName DhcpOptionCode = 66
        OptionBootFileName   DhcpOptionCode = 67

        OptionUserClass DhcpOptionCode = 77

        OptionClientArchitecture DhcpOptionCode = 93

        OptionTZPOSIXString    DhcpOptionCode = 100
        OptionTZDatabaseString DhcpOptionCode = 101

        OptionClasslessRouteFormat DhcpOptionCode = 121
)

// Dhcp OpCodes Types
const (
        Request OpCode = 1 // From Client
        Reply   OpCode = 2 // From Server
)

// DHCP Client/Server Message Type 53
const (
        DhcpDiscover MessageType = 1 // From Client
        DhcpOffer    MessageType = 2 // From Server
        DhcpRequest  MessageType = 3 // From Client
        DhcpDecline  MessageType = 4 // From Client
        DhcpACK      MessageType = 5 // From Server
        DhcpNAK      MessageType = 6 // From Server
        DhcpRelease  MessageType = 7 // From Client
        DhcpInform   MessageType = 8 // From Client
)

const (
        CLIENT_CONNECTION_NOT_REQUIRED = "Connection to client is not required"
)
