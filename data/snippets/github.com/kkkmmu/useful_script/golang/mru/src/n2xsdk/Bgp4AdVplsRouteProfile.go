package n2xsdk

type Bgp4AdVplsRouteProfile struct {
 Handler string
}

func(np *Bgp4AdVplsRouteProfile) EnableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4AdVplsRouteProfile EnableAutoPathAttribute
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4AdVplsRouteProfile DisableAutoPathAttribute
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsAutoPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4AdVplsRouteProfile IsAutoPathAttributeEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetExportVplsIdIncrementingRange () error {
 //parameters: RouteProfileHandle VplsIdType FirstVplsId VplsIdIncrement PercentOverlap
 //AgtBgp4AdVplsRouteProfile SetExportVplsIdIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetExportVplsIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle VplsIdType
 //AgtBgp4AdVplsRouteProfile GetExportVplsIdIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstNextHop PrefixLength NextHopIncrement NextHopRepeat
 //AgtBgp4AdVplsRouteProfile SetNextHopIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle NextHopList
 //AgtBgp4AdVplsRouteProfile SetNextHopList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetNextHopList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetNextHop () error {
 //parameters: RouteProfileHandle NextHop
 //AgtBgp4AdVplsRouteProfile SetNextHop
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetNextHop
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetNextHopAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetNextHopAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetExportRouteTargetIncrementingRange () error {
 //parameters: RouteProfileHandle RouteTargetType FirstRouteTarget RouteTargetIncrement PercentOverlap
 //AgtBgp4AdVplsRouteProfile SetExportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetExportRouteTargetIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4AdVplsRouteProfile GetExportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetRouteDistinguisherIncrementingRange () error {
 //parameters: RouteProfileHandle RouteDistinguisherType FirstRouteDistinguisher RouteDistinguisherIncrement PercentOverlap
 //AgtBgp4AdVplsRouteProfile SetRouteDistinguisherIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetRouteDistinguisherIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4AdVplsRouteProfile GetRouteDistinguisherIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetRouteDistinguisherAuto ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4AdVplsRouteProfile GetRouteDistinguisherAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4AdVplsRouteProfile GetState
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4AdVplsRouteProfile IsEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile EnableAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile DisableAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4AdVplsRouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4AdVplsRouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4AdVplsRouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4AdVplsRouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4AdVplsRouteProfile SetAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4AdVplsRouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4AdVplsRouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4AdVplsRouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4AdVplsRouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4AdVplsRouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4AdVplsRouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4AdVplsRouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4AdVplsRouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4AdVplsRouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4AdVplsRouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4AdVplsRouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4AdVplsRouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4AdVplsRouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4AdVplsRouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4AdVplsRouteProfile SetOrigin
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4AdVplsRouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4AdVplsRouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4AdVplsRouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4AdVplsRouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4AdVplsRouteProfile SetAggregator
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4AdVplsRouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4AdVplsRouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4AdVplsRouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4AdVplsRouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4AdVplsRouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4AdVplsRouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4AdVplsRouteProfile SetCommunities
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4AdVplsRouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4AdVplsRouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4AdVplsRouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4AdVplsRouteProfile Advertise
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4AdVplsRouteProfile Withdraw
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Enable
 return nil
}

func(np *Bgp4AdVplsRouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile Disable
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetExportRouteTargetType () error {
 //parameters: RouteProfileHandle RouteTargetTypeValue
 //AgtBgp4AdVplsRouteProfile SetExportRouteTargetType
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetExportRouteTargetType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetExportRouteTargetType
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) SetRouteDistinguisherType () error {
 //parameters: RouteProfileHandle RouteDistinguisherTypeValue
 //AgtBgp4AdVplsRouteProfile SetRouteDistinguisherType
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetRouteDistinguisherType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetRouteDistinguisherType
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) EnableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile EnableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4AdVplsRouteProfile) DisableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile DisableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4AdVplsRouteProfile) IsAutoRouteDistinguisherEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile IsAutoRouteDistinguisherEnabled
 return nil
}

func(np *Bgp4AdVplsRouteProfile) SetVpnsPerPeer () error {
 //parameters: RouteProfileHandle NumVpnsPerPeer
 //AgtBgp4AdVplsRouteProfile SetVpnsPerPeer
 return nil
}

func(np *Bgp4AdVplsRouteProfile) GetVpnsPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4AdVplsRouteProfile) GetTotalNumberOfVpns ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4AdVplsRouteProfile GetTotalNumberOfVpns
 return "", nil
}

