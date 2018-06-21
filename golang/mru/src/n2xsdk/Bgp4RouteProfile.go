package n2xsdk

type Bgp4RouteProfile struct {
 Handler string
}

func(np *Bgp4RouteProfile) SetRoutesPerPeer () error {
 //parameters: RouteProfileHandle RoutesPerPeer
 //AgtBgp4RouteProfile SetRoutesPerPeer
 return nil
}

func(np *Bgp4RouteProfile) GetRoutesPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetRoutesPerPeer
 return "", nil
}

func(np *Bgp4RouteProfile) GetTotalRoutesInProfile ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetTotalRoutesInProfile
 return "", nil
}

func(np *Bgp4RouteProfile) GetNthRouteInProfile ()(string, error) {
 //parameters: RouteProfileHandle RouteProfileIndex
 //AgtBgp4RouteProfile GetNthRouteInProfile
 return "", nil
}

func(np *Bgp4RouteProfile) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4RouteProfile) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4RouteProfile) EnableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile EnableLabeling
 return nil
}

func(np *Bgp4RouteProfile) DisableLabeling () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile DisableLabeling
 return nil
}

func(np *Bgp4RouteProfile) IsLabelingEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile IsLabelingEnabled
 return nil
}

func(np *Bgp4RouteProfile) SetLabelingMode () error {
 //parameters: RouteProfileHandle LabelingMode
 //AgtBgp4RouteProfile SetLabelingMode
 return nil
}

func(np *Bgp4RouteProfile) GetLabelingMode ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetLabelingMode
 return "", nil
}

func(np *Bgp4RouteProfile) SetUserLabelType () error {
 //parameters: RouteProfileHandle UserLabelType
 //AgtBgp4RouteProfile SetUserLabelType
 return nil
}

func(np *Bgp4RouteProfile) GetUserLabelType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetUserLabelType
 return "", nil
}

func(np *Bgp4RouteProfile) SetUserLabelIncrementingRange () error {
 //parameters: RouteProfileHandle StartLabel LabelIncrement LabelRepeat
 //AgtBgp4RouteProfile SetUserLabelIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) GetUserLabelIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetUserLabelIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) SetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle LabelList
 //AgtBgp4RouteProfile SetUserLabelList
 return "", nil
}

func(np *Bgp4RouteProfile) GetUserLabelList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetUserLabelList
 return "", nil
}

func(np *Bgp4RouteProfile) SetUserLabel () error {
 //parameters: RouteProfileHandle Label
 //AgtBgp4RouteProfile SetUserLabel
 return nil
}

func(np *Bgp4RouteProfile) GetUserLabel ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetUserLabel
 return "", nil
}

func(np *Bgp4RouteProfile) GetLabels ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile GetLabels
 return "", nil
}

func(np *Bgp4RouteProfile) GetLabelPoolState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile GetLabelPoolState
 return "", nil
}

func(np *Bgp4RouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4RouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile GetState
 return "", nil
}

func(np *Bgp4RouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4RouteProfile IsEnabled
 return nil
}

func(np *Bgp4RouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile EnableAsPath
 return nil
}

func(np *Bgp4RouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4RouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile DisableAsPath
 return nil
}

func(np *Bgp4RouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4RouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4RouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4RouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4RouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4RouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4RouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4RouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4RouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4RouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4RouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4RouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4RouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4RouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4RouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4RouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4RouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4RouteProfile SetAsPath
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4RouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4RouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4RouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4RouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4RouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4RouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4RouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4RouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4RouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4RouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4RouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4RouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4RouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4RouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4RouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4RouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4RouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4RouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4RouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4RouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4RouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4RouteProfile SetOrigin
 return nil
}

func(np *Bgp4RouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4RouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4RouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4RouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4RouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4RouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4RouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4RouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4RouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4RouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4RouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4RouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4RouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4RouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4RouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4RouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4RouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4RouteProfile SetAggregator
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4RouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4RouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4RouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4RouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4RouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4RouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4RouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4RouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4RouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4RouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4RouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4RouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4RouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4RouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4RouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4RouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4RouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4RouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4RouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4RouteProfile SetCommunities
 return nil
}

func(np *Bgp4RouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4RouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4RouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4RouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4RouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4RouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4RouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4RouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4RouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4RouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4RouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4RouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4RouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4RouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4RouteProfile Advertise
 return nil
}

func(np *Bgp4RouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4RouteProfile Withdraw
 return nil
}

func(np *Bgp4RouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Enable
 return nil
}

func(np *Bgp4RouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4RouteProfile Disable
 return nil
}

