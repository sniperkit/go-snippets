package n2xsdk

type Bgp4VpnIpv6RouteProfile struct {
 Handler string
}

func(np *Bgp4VpnIpv6RouteProfile) EnableSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsSourceAsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsSourceAsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetSourceAsType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetSourceAsType
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableAutoSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAuto4ByteSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableAuto4ByteSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableAutoSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAuto4ByteSourceAs () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableAuto4ByteSourceAs
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoSourceAsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsAutoSourceAsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetSourceAsAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetSourceAsAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAuto4ByteSourceAsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsAuto4ByteSourceAsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteSourceAsAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteSourceAsAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetSourceAsIncrementingRange () error {
 //parameters: RouteProfileHandle First2ByteSourceAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4VpnIpv6RouteProfile SetSourceAsIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteSourceAsIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteSourceAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4VpnIpv6RouteProfile Set4ByteSourceAsIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetSourceAsIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetSourceAsIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteSourceAsIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteSourceAsIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableVrfImportRt () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableVrfImportRt
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableVrfImportRt () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableVrfImportRt
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsVrfImportRtEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsVrfImportRtEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoVrfRtPeAddress () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableAutoVrfRtPeAddress
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoVrfRtPeAddress () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableAutoVrfRtPeAddress
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoVrfRtPeAddressEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsAutoVrfRtPeAddressEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetVrfRtPeAddressAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetVrfRtPeAddressAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetVrfRtPeAddressIncrementingRange () error {
 //parameters: RouteProfileHandle FirstRouteTarget PrefixLength RouteTargetIncrement RouteTargetRepeat
 //AgtBgp4VpnIpv6RouteProfile SetVrfRtPeAddressIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetVrfRtPeAddressIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetVrfRtPeAddressIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetVrfRtIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstRouteTargetId RouteTargetIncrement PercentOverlap
 //AgtBgp4VpnIpv6RouteProfile SetVrfRtIdIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetVrfRtIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetVrfRtIdIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesPerVpn () error {
 //parameters: RouteProfileHandle NumRoutesPerVpn
 //AgtBgp4VpnIpv6RouteProfile SetRoutesPerVpn
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesPerVpn ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRoutesPerVpn
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) CreateMatchingVpnVrfPool () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile CreateMatchingVpnVrfPool
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableConnector () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableConnector
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableConnector () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableConnector
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsConnectorEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsConnectorEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoConnectorIpAddress () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableAutoConnectorIpAddress
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoConnectorIpAddress () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableAutoConnectorIpAddress
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoConnectorIpAddressEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsAutoConnectorIpAddressEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorIpAddressAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorIpAddressAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorRouteDistinguisherIncrementingRange () error {
 //parameters: RouteProfileHandle FirstRouteDistinguisher RouteDistinguisherIncrement RouteDistinguisherRepeat PercentOverlap
 //AgtBgp4VpnIpv6RouteProfile SetConnectorRouteDistinguisherIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorRouteDistinguisherIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorRouteDistinguisherIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorIpAddressIncrementingRange () error {
 //parameters: RouteProfileHandle FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4VpnIpv6RouteProfile SetConnectorIpAddressIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorIpAddressIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorIpAddressIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherList
 //AgtBgp4VpnIpv6RouteProfile SetConnectorRouteDistinguisherList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorRouteDistinguisherList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorIpAddressList ()(string, error) {
 //parameters: RouteProfileHandle IpAddressList
 //AgtBgp4VpnIpv6RouteProfile SetConnectorIpAddressList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorIpAddressList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorIpAddressList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorRouteDistinguisher () error {
 //parameters: RouteProfileHandle RouteDistinguisher
 //AgtBgp4VpnIpv6RouteProfile SetConnectorRouteDistinguisher
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorRouteDistinguisher ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorRouteDistinguisher
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetConnectorIpAddress () error {
 //parameters: RouteProfileHandle IpAddress
 //AgtBgp4VpnIpv6RouteProfile SetConnectorIpAddress
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetConnectorIpAddress ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetConnectorIpAddress
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetExportRouteTargetType () error {
 //parameters: RouteProfileHandle RouteTargetTypeValue
 //AgtBgp4VpnIpv6RouteProfile SetExportRouteTargetType
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetExportRouteTargetType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetExportRouteTargetType
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRouteDistinguisherType () error {
 //parameters: RouteProfileHandle RouteDistinguisherTypeValue
 //AgtBgp4VpnIpv6RouteProfile SetRouteDistinguisherType
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRouteDistinguisherType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRouteDistinguisherType
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoRouteDistinguisherEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsAutoRouteDistinguisherEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetVpnsPerPeer () error {
 //parameters: RouteProfileHandle NumVpnsPerPeer
 //AgtBgp4VpnIpv6RouteProfile SetVpnsPerPeer
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetVpnsPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetTotalNumberOfVpns ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetTotalNumberOfVpns
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetExportRouteTargetIncrementingRange () error {
 //parameters: RouteProfileHandle RouteTargetType FirstRouteTarget RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4VpnIpv6RouteProfile SetExportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetExportRouteTargetIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4VpnIpv6RouteProfile GetExportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetExportRouteTargetList ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType RouteTargetList
 //AgtBgp4VpnIpv6RouteProfile SetExportRouteTargetList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetExportRouteTargetList ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4VpnIpv6RouteProfile GetExportRouteTargetList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetExportRouteTarget () error {
 //parameters: RouteProfileHandle RouteTargetType RouteTarget
 //AgtBgp4VpnIpv6RouteProfile SetExportRouteTarget
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetExportRouteTarget ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4VpnIpv6RouteProfile GetExportRouteTarget
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRouteDistinguisherIncrementingRange () error {
 //parameters: RouteProfileHandle RouteDistinguisherType FirstRouteDistinguisher RouteDistinguisherIncrement RouteDistinguisherRepeat PercentOverlap
 //AgtBgp4VpnIpv6RouteProfile SetRouteDistinguisherIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRouteDistinguisherIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4VpnIpv6RouteProfile GetRouteDistinguisherIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRouteDistinguisherAuto ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4VpnIpv6RouteProfile GetRouteDistinguisherAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType RouteDistinguisherList
 //AgtBgp4VpnIpv6RouteProfile SetRouteDistinguisherList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4VpnIpv6RouteProfile GetRouteDistinguisherList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRouteDistinguisher () error {
 //parameters: RouteProfileHandle RouteDistinguisherType RouteDistinguisher
 //AgtBgp4VpnIpv6RouteProfile SetRouteDistinguisher
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRouteDistinguisher ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4VpnIpv6RouteProfile GetRouteDistinguisher
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesIncrementingRange () error {
 //parameters: RouteProfileHandle FirstRoute PrefixLengthIpv6 RouteIncrement PercentOverlap
 //AgtBgp4VpnIpv6RouteProfile SetRoutesIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRoutesIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesPrefixLength () error {
 //parameters: RouteProfileHandle PrefixLengthIpv6
 //AgtBgp4VpnIpv6RouteProfile SetRoutesPrefixLength
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesPrefixLength ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRoutesPrefixLength
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesSubRangeOffsets () error {
 //parameters: RouteProfileHandle FirstRoute MsbOffsetList
 //AgtBgp4VpnIpv6RouteProfile SetRoutesSubRangeOffsets
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesSubRangeOffsets ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRoutesSubRangeOffsets
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesSubRange () error {
 //parameters: RouteProfileHandle SubRangeInstance RouteIncrement Count Repeat
 //AgtBgp4VpnIpv6RouteProfile SetRoutesSubRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesSubRange ()(string, error) {
 //parameters: RouteProfileHandle SubRangeInstance
 //AgtBgp4VpnIpv6RouteProfile GetRoutesSubRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4VpnIpv6RouteProfile EnableAutoPathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4VpnIpv6RouteProfile DisableAutoPathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4VpnIpv6RouteProfile IsAutoPathAttributeEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstNextHop PrefixLength NextHopIncrement NextHopRepeat
 //AgtBgp4VpnIpv6RouteProfile SetNextHopIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle NextHopList
 //AgtBgp4VpnIpv6RouteProfile SetNextHopList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetNextHopList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetNextHop () error {
 //parameters: RouteProfileHandle NextHop
 //AgtBgp4VpnIpv6RouteProfile SetNextHop
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetNextHop
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNextHopAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetNextHopAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetLinkLocalNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstLinkLocalNextHop PrefixLength LinkLocalNextHopIncrement LinkLocalNextHopRepeat
 //AgtBgp4VpnIpv6RouteProfile SetLinkLocalNextHopIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLinkLocalNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetLinkLocalNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle LinkLocalNextHopList
 //AgtBgp4VpnIpv6RouteProfile SetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetLinkLocalNextHop () error {
 //parameters: RouteProfileHandle LinkLocalNextHop
 //AgtBgp4VpnIpv6RouteProfile SetLinkLocalNextHop
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLinkLocalNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetLinkLocalNextHop
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetRoutesPerPeer () error {
 //parameters: RouteProfileHandle RoutesPerPeer
 //AgtBgp4VpnIpv6RouteProfile SetRoutesPerPeer
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetRoutesPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetRoutesPerPeer
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetTotalRoutesInProfile ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetTotalRoutesInProfile
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNthRouteInProfile ()(string, error) {
 //parameters: RouteProfileHandle RouteProfileIndex
 //AgtBgp4VpnIpv6RouteProfile GetNthRouteInProfile
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableLabeling
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableLabeling
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsLabelingEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsLabelingEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetLabelingMode () error {
 //parameters: RouteProfileHandle LabelingMode
 //AgtBgp4VpnIpv6RouteProfile SetLabelingMode
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLabelingMode ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetLabelingMode
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetUserLabelType () error {
 //parameters: RouteProfileHandle UserLabelType
 //AgtBgp4VpnIpv6RouteProfile SetUserLabelType
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetUserLabelType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetUserLabelType
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetUserLabelIncrementingRange () error {
 //parameters: RouteProfileHandle StartLabel LabelIncrement LabelRepeat
 //AgtBgp4VpnIpv6RouteProfile SetUserLabelIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetUserLabelIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetUserLabelIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle LabelList
 //AgtBgp4VpnIpv6RouteProfile SetUserLabelList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetUserLabelList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetUserLabel () error {
 //parameters: RouteProfileHandle Label
 //AgtBgp4VpnIpv6RouteProfile SetUserLabel
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetUserLabel ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetUserLabel
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLabels ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile GetLabels
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLabelPoolState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile GetLabelPoolState
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile GetState
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4VpnIpv6RouteProfile IsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile EnableAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile DisableAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4VpnIpv6RouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4VpnIpv6RouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4VpnIpv6RouteProfile SetAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4VpnIpv6RouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4VpnIpv6RouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4VpnIpv6RouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4VpnIpv6RouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4VpnIpv6RouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4VpnIpv6RouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4VpnIpv6RouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4VpnIpv6RouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4VpnIpv6RouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4VpnIpv6RouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4VpnIpv6RouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4VpnIpv6RouteProfile SetOrigin
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4VpnIpv6RouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4VpnIpv6RouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4VpnIpv6RouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4VpnIpv6RouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4VpnIpv6RouteProfile SetAggregator
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4VpnIpv6RouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4VpnIpv6RouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4VpnIpv6RouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4VpnIpv6RouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4VpnIpv6RouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4VpnIpv6RouteProfile SetCommunities
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4VpnIpv6RouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4VpnIpv6RouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4VpnIpv6RouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4VpnIpv6RouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4VpnIpv6RouteProfile Advertise
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4VpnIpv6RouteProfile Withdraw
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Enable
 return nil
}

func(np *Bgp4VpnIpv6RouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4VpnIpv6RouteProfile Disable
 return nil
}

