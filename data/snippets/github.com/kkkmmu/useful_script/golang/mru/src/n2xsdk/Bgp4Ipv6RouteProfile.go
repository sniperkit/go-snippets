package n2xsdk

type Bgp4Ipv6RouteProfile struct {
 Handler string
}

func(np *Bgp4Ipv6RouteProfile) SetRoutesIncrementingRange () error {
 //parameters: RouteProfileHandle FirstRoute PrefixLengthIpv6 RouteIncrement PercentOverlap
 //AgtBgp4Ipv6RouteProfile SetRoutesIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetRoutesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetRoutesIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetRoutesPrefixLength () error {
 //parameters: RouteProfileHandle PrefixLengthIpv6
 //AgtBgp4Ipv6RouteProfile SetRoutesPrefixLength
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetRoutesPrefixLength ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetRoutesPrefixLength
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetRoutesSubRangeOffsets () error {
 //parameters: RouteProfileHandle FirstRoute MsbOffsetList
 //AgtBgp4Ipv6RouteProfile SetRoutesSubRangeOffsets
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetRoutesSubRangeOffsets ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetRoutesSubRangeOffsets
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetRoutesSubRange () error {
 //parameters: RouteProfileHandle SubRangeInstance RouteIncrement Count Repeat
 //AgtBgp4Ipv6RouteProfile SetRoutesSubRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetRoutesSubRange ()(string, error) {
 //parameters: RouteProfileHandle SubRangeInstance
 //AgtBgp4Ipv6RouteProfile GetRoutesSubRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) EnableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4Ipv6RouteProfile EnableAutoPathAttribute
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4Ipv6RouteProfile DisableAutoPathAttribute
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsAutoPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4Ipv6RouteProfile IsAutoPathAttributeEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) SetNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstNextHop PrefixLength NextHopIncrement NextHopRepeat
 //AgtBgp4Ipv6RouteProfile SetNextHopIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle NextHopList
 //AgtBgp4Ipv6RouteProfile SetNextHopList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetNextHopList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetNextHop () error {
 //parameters: RouteProfileHandle NextHop
 //AgtBgp4Ipv6RouteProfile SetNextHop
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetNextHop
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetNextHopAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetNextHopAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetLinkLocalNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstLinkLocalNextHop PrefixLength LinkLocalNextHopIncrement LinkLocalNextHopRepeat
 //AgtBgp4Ipv6RouteProfile SetLinkLocalNextHopIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetLinkLocalNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetLinkLocalNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle LinkLocalNextHopList
 //AgtBgp4Ipv6RouteProfile SetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetLinkLocalNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetLinkLocalNextHopList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetLinkLocalNextHop () error {
 //parameters: RouteProfileHandle LinkLocalNextHop
 //AgtBgp4Ipv6RouteProfile SetLinkLocalNextHop
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetLinkLocalNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetLinkLocalNextHop
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetRoutesPerPeer () error {
 //parameters: RouteProfileHandle RoutesPerPeer
 //AgtBgp4Ipv6RouteProfile SetRoutesPerPeer
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetRoutesPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetRoutesPerPeer
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetTotalRoutesInProfile ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetTotalRoutesInProfile
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetNthRouteInProfile ()(string, error) {
 //parameters: RouteProfileHandle RouteProfileIndex
 //AgtBgp4Ipv6RouteProfile GetNthRouteInProfile
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) EnableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile EnableLabeling
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile DisableLabeling
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsLabelingEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile IsLabelingEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) SetLabelingMode () error {
 //parameters: RouteProfileHandle LabelingMode
 //AgtBgp4Ipv6RouteProfile SetLabelingMode
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetLabelingMode ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetLabelingMode
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetUserLabelType () error {
 //parameters: RouteProfileHandle UserLabelType
 //AgtBgp4Ipv6RouteProfile SetUserLabelType
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetUserLabelType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetUserLabelType
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetUserLabelIncrementingRange () error {
 //parameters: RouteProfileHandle StartLabel LabelIncrement LabelRepeat
 //AgtBgp4Ipv6RouteProfile SetUserLabelIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetUserLabelIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetUserLabelIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle LabelList
 //AgtBgp4Ipv6RouteProfile SetUserLabelList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetUserLabelList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetUserLabel () error {
 //parameters: RouteProfileHandle Label
 //AgtBgp4Ipv6RouteProfile SetUserLabel
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetUserLabel ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetUserLabel
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetLabels ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile GetLabels
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetLabelPoolState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile GetLabelPoolState
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile GetState
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4Ipv6RouteProfile IsEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile EnableAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile DisableAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4Ipv6RouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4Ipv6RouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4Ipv6RouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4Ipv6RouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4Ipv6RouteProfile SetAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4Ipv6RouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4Ipv6RouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4Ipv6RouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4Ipv6RouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4Ipv6RouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4Ipv6RouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4Ipv6RouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4Ipv6RouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4Ipv6RouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4Ipv6RouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4Ipv6RouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4Ipv6RouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4Ipv6RouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4Ipv6RouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4Ipv6RouteProfile SetOrigin
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4Ipv6RouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4Ipv6RouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4Ipv6RouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4Ipv6RouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4Ipv6RouteProfile SetAggregator
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4Ipv6RouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4Ipv6RouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4Ipv6RouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4Ipv6RouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4Ipv6RouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4Ipv6RouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4Ipv6RouteProfile SetCommunities
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4Ipv6RouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4Ipv6RouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4Ipv6RouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4Ipv6RouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4Ipv6RouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4Ipv6RouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4Ipv6RouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4Ipv6RouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4Ipv6RouteProfile Advertise
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4Ipv6RouteProfile Withdraw
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Enable
 return nil
}

func(np *Bgp4Ipv6RouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4Ipv6RouteProfile Disable
 return nil
}

