package n2xsdk

type Bgp4MVpnIpv6RouteProfile struct {
 Handler string
}

func(np *Bgp4MVpnIpv6RouteProfile) SetVpnsPerPeer () error {
 //parameters: RouteProfileHandle NumVpnsPerPeer
 //AgtBgp4MVpnIpv6RouteProfile SetVpnsPerPeer
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetVpnsPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetRouteType () error {
 //parameters: RouteProfileHandle RouteType
 //AgtBgp4MVpnIpv6RouteProfile SetRouteType
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetRouteType
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoOriginatorIp () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoOriginatorIp
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoOriginatorIp () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoOriginatorIp
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoOriginatorIpEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsAutoOriginatorIpEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorIpAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorIpAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetOriginatorIpIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorIp PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetOriginatorIpIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorIpIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorIpIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetDistributionMode () error {
 //parameters: RouteProfileHandle DistributionMode
 //AgtBgp4MVpnIpv6RouteProfile SetDistributionMode
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetDistributionMode ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetDistributionMode
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSgPairsPerVpn () error {
 //parameters: RouteProfileHandle SgPairsPerVpn
 //AgtBgp4MVpnIpv6RouteProfile SetSgPairsPerVpn
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSgPairsPerVpn ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetSgPairsPerVpn
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetGroupsPerVpn () error {
 //parameters: RouteProfileHandle GroupsPerVpn
 //AgtBgp4MVpnIpv6RouteProfile SetGroupsPerVpn
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetGroupsPerVpn ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetGroupsPerVpn
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSourcesPerVpn () error {
 //parameters: RouteProfileHandle SourcesPerVpn
 //AgtBgp4MVpnIpv6RouteProfile SetSourcesPerVpn
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSourcesPerVpn ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetSourcesPerVpn
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetGroupAddressSubRangeOffsets () error {
 //parameters: RouteProfileHandle GroupAddress MsbOffsetList
 //AgtBgp4MVpnIpv6RouteProfile SetGroupAddressSubRangeOffsets
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetGroupAddressSubRangeOffsets ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetGroupAddressSubRangeOffsets
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetGroupAddressSubRange () error {
 //parameters: RouteProfileHandle SubRangeInstance Increment Count Repeat
 //AgtBgp4MVpnIpv6RouteProfile SetGroupAddressSubRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetGroupAddressSubRange ()(string, error) {
 //parameters: RouteProfileHandle SubRangeInstance
 //AgtBgp4MVpnIpv6RouteProfile GetGroupAddressSubRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetGroupAddressIncrementingRange () error {
 //parameters: RouteProfileHandle GroupAddress PrefixLength Increment Count
 //AgtBgp4MVpnIpv6RouteProfile SetGroupAddressIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetGroupAddressIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetGroupAddressIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetGroupAddress () error {
 //parameters: RouteProfileHandle GroupAddress
 //AgtBgp4MVpnIpv6RouteProfile SetGroupAddress
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetGroupAddress ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetGroupAddress
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSourceAddressSubRangeOffsets () error {
 //parameters: RouteProfileHandle SourceAddress MsbOffsetList
 //AgtBgp4MVpnIpv6RouteProfile SetSourceAddressSubRangeOffsets
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSourceAddressSubRangeOffsets ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetSourceAddressSubRangeOffsets
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSourceAddressSubRange () error {
 //parameters: RouteProfileHandle SubRangeInstance Increment Count Repeat
 //AgtBgp4MVpnIpv6RouteProfile SetSourceAddressSubRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSourceAddressSubRange ()(string, error) {
 //parameters: RouteProfileHandle SubRangeInstance
 //AgtBgp4MVpnIpv6RouteProfile GetSourceAddressSubRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSourceAddressIncrementingRange () error {
 //parameters: RouteProfileHandle SourceAddress PrefixLength Increment Count
 //AgtBgp4MVpnIpv6RouteProfile SetSourceAddressIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSourceAddressIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetSourceAddressIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetSourceAddress () error {
 //parameters: RouteProfileHandle SourceAddress
 //AgtBgp4MVpnIpv6RouteProfile SetSourceAddress
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetSourceAddress ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetSourceAddress
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnablePmsiTunnelAttribute () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnablePmsiTunnelAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisablePmsiTunnelAttribute () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisablePmsiTunnelAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsPmsiTunnelAttributeEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsPmsiTunnelAttributeEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnablePmsiTunnelLeafFlag () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnablePmsiTunnelLeafFlag
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisablePmsiTunnelLeafFlag () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisablePmsiTunnelLeafFlag
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsPmsiTunnelLeafFlagEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsPmsiTunnelLeafFlagEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetPmsiTunnelType () error {
 //parameters: RouteProfileHandle PmsiTunnelType
 //AgtBgp4MVpnIpv6RouteProfile SetPmsiTunnelType
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetPmsiTunnelType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetPmsiTunnelType
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoLabelMode () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoLabelMode
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoLabelMode () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoLabelMode
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoLabelModeEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsAutoLabelModeEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetUserLabelIncrementingRange () error {
 //parameters: RouteProfileHandle StartLabel LabelIncrement LabelRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetUserLabelIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetUserLabelIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetUserLabelIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoUnicastTunnelEndpointMode () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoUnicastTunnelEndpointMode
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoUnicastTunnelEndpointMode () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoUnicastTunnelEndpointMode
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoUnicastTunnelEndpointModeEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsAutoUnicastTunnelEndpointModeEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetUnicastTunnelEndpointAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetUnicastTunnelEndpointAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetUnicastTunnelEndpointIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorIp PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetUnicastTunnelEndpointIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetUnicastTunnelEndpointIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetUnicastTunnelEndpointIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoPathAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoPathAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4MVpnIpv6RouteProfile IsAutoPathAttributeEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstNextHop PrefixLength NextHopIncrement NextHopRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetNextHopIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle NextHopList
 //AgtBgp4MVpnIpv6RouteProfile SetNextHopList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetNextHopList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetNextHop () error {
 //parameters: RouteProfileHandle NextHop
 //AgtBgp4MVpnIpv6RouteProfile SetNextHop
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetNextHop
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetNextHopAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetNextHopAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetLinkLocalNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstLinkLocalNextHop PrefixLength LinkLocalNextHopIncrement LinkLocalNextHopRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetLinkLocalNextHopIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetLinkLocalNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetLinkLocalNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle LinkLocalNextHopList
 //AgtBgp4MVpnIpv6RouteProfile SetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetLinkLocalNextHop () error {
 //parameters: RouteProfileHandle LinkLocalNextHop
 //AgtBgp4MVpnIpv6RouteProfile SetLinkLocalNextHop
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetLinkLocalNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetLinkLocalNextHop
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4MVpnIpv6RouteProfile GetState
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4MVpnIpv6RouteProfile IsEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile EnableAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile DisableAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4MVpnIpv6RouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4MVpnIpv6RouteProfile SetAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4MVpnIpv6RouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4MVpnIpv6RouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4MVpnIpv6RouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4MVpnIpv6RouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4MVpnIpv6RouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4MVpnIpv6RouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4MVpnIpv6RouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4MVpnIpv6RouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4MVpnIpv6RouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4MVpnIpv6RouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4MVpnIpv6RouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4MVpnIpv6RouteProfile SetOrigin
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4MVpnIpv6RouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4MVpnIpv6RouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4MVpnIpv6RouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4MVpnIpv6RouteProfile SetAggregator
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4MVpnIpv6RouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4MVpnIpv6RouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4MVpnIpv6RouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4MVpnIpv6RouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4MVpnIpv6RouteProfile SetCommunities
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4MVpnIpv6RouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4MVpnIpv6RouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4MVpnIpv6RouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4MVpnIpv6RouteProfile Advertise
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4MVpnIpv6RouteProfile Withdraw
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Enable
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile Disable
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetExportRouteTargetIncrementingRange () error {
 //parameters: RouteProfileHandle RouteTargetType FirstRouteTarget RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4MVpnIpv6RouteProfile SetExportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetExportRouteTargetIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4MVpnIpv6RouteProfile GetExportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetExportRouteTargetList ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType RouteTargetList
 //AgtBgp4MVpnIpv6RouteProfile SetExportRouteTargetList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetExportRouteTargetList ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4MVpnIpv6RouteProfile GetExportRouteTargetList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetExportRouteTarget () error {
 //parameters: RouteProfileHandle RouteTargetType RouteTarget
 //AgtBgp4MVpnIpv6RouteProfile SetExportRouteTarget
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetExportRouteTarget ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4MVpnIpv6RouteProfile GetExportRouteTarget
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetRouteDistinguisherIncrementingRange () error {
 //parameters: RouteProfileHandle RouteDistinguisherType FirstRouteDistinguisher RouteDistinguisherIncrement RouteDistinguisherRepeat PercentOverlap
 //AgtBgp4MVpnIpv6RouteProfile SetRouteDistinguisherIncrementingRange
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteDistinguisherIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4MVpnIpv6RouteProfile GetRouteDistinguisherIncrementingRange
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteDistinguisherAuto ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4MVpnIpv6RouteProfile GetRouteDistinguisherAuto
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType RouteDistinguisherList
 //AgtBgp4MVpnIpv6RouteProfile SetRouteDistinguisherList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteDistinguisherList ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4MVpnIpv6RouteProfile GetRouteDistinguisherList
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetRouteDistinguisher () error {
 //parameters: RouteProfileHandle RouteDistinguisherType RouteDistinguisher
 //AgtBgp4MVpnIpv6RouteProfile SetRouteDistinguisher
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteDistinguisher ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4MVpnIpv6RouteProfile GetRouteDistinguisher
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetExportRouteTargetType () error {
 //parameters: RouteProfileHandle RouteTargetTypeValue
 //AgtBgp4MVpnIpv6RouteProfile SetExportRouteTargetType
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetExportRouteTargetType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetExportRouteTargetType
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) SetRouteDistinguisherType () error {
 //parameters: RouteProfileHandle RouteDistinguisherTypeValue
 //AgtBgp4MVpnIpv6RouteProfile SetRouteDistinguisherType
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetRouteDistinguisherType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetRouteDistinguisherType
 return "", nil
}

func(np *Bgp4MVpnIpv6RouteProfile) EnableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile EnableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) DisableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile DisableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) IsAutoRouteDistinguisherEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile IsAutoRouteDistinguisherEnabled
 return nil
}

func(np *Bgp4MVpnIpv6RouteProfile) GetTotalNumberOfVpns ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4MVpnIpv6RouteProfile GetTotalNumberOfVpns
 return "", nil
}

