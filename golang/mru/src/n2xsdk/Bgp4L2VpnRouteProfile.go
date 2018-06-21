package n2xsdk

type Bgp4L2VpnRouteProfile struct {
 Handler string
}

func(np *Bgp4L2VpnRouteProfile) EnableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4L2VpnRouteProfile EnableAutoPathAttribute
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableAutoPathAttribute () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4L2VpnRouteProfile DisableAutoPathAttribute
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsAutoPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttributeAuto
 //AgtBgp4L2VpnRouteProfile IsAutoPathAttributeEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetNextHopIncrementingRange () error {
 //parameters: RouteProfileHandle FirstNextHop PrefixLength NextHopIncrement NextHopRepeat
 //AgtBgp4L2VpnRouteProfile SetNextHopIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetNextHopIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetNextHopIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle NextHopList
 //AgtBgp4L2VpnRouteProfile SetNextHopList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetNextHopList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetNextHopList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetNextHop () error {
 //parameters: RouteProfileHandle NextHop
 //AgtBgp4L2VpnRouteProfile SetNextHop
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetNextHop ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetNextHop
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetNextHopAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetNextHopAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetLabelBaseMode () error {
 //parameters: RouteProfileHandle LabelBaseMode
 //AgtBgp4L2VpnRouteProfile SetLabelBaseMode
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetLabelBaseMode ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetLabelBaseMode
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetUserLabelBase () error {
 //parameters: RouteProfileHandle UserLabelBase
 //AgtBgp4L2VpnRouteProfile SetUserLabelBase
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetUserLabelBase ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetUserLabelBase
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetLastLabelAllocated ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetLastLabelAllocated
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVeBlockSize () error {
 //parameters: RouteProfileHandle VeBlockSize
 //AgtBgp4L2VpnRouteProfile SetVeBlockSize
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVeBlockSize ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVeBlockSize
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVeBlocksPerVeId () error {
 //parameters: RouteProfileHandle NumVeBlocksPerVeId
 //AgtBgp4L2VpnRouteProfile SetVeBlocksPerVeId
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVeBlocksPerVeId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVeBlocksPerVeId
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVeIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstVeId VeIdIncrement VeIdRepeat
 //AgtBgp4L2VpnRouteProfile SetVeIdIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVeIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVeIdIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVeIdList ()(string, error) {
 //parameters: RouteProfileHandle VeIdList
 //AgtBgp4L2VpnRouteProfile SetVeIdList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetVeIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVeIdList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVeId () error {
 //parameters: RouteProfileHandle VeId
 //AgtBgp4L2VpnRouteProfile SetVeId
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVeId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVeId
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) EnableVcControlWord () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile EnableVcControlWord
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableVcControlWord () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile DisableVcControlWord
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsVcControlWordEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile IsVcControlWordEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetVcInterfaceMtu () error {
 //parameters: RouteProfileHandle Mtu
 //AgtBgp4L2VpnRouteProfile SetVcInterfaceMtu
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVcInterfaceMtu ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVcInterfaceMtu
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetVcEncapsulationType () error {
 //parameters: RouteProfileHandle EncapsulationType
 //AgtBgp4L2VpnRouteProfile SetVcEncapsulationType
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVcEncapsulationType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVcEncapsulationType
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetExportRouteTargetIncrementingRange () error {
 //parameters: RouteProfileHandle RouteTargetType FirstRouteTarget RouteTargetIncrement PercentOverlap
 //AgtBgp4L2VpnRouteProfile SetExportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetExportRouteTargetIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteTargetType
 //AgtBgp4L2VpnRouteProfile GetExportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetRouteDistinguisherIncrementingRange () error {
 //parameters: RouteProfileHandle RouteDistinguisherType FirstRouteDistinguisher RouteDistinguisherIncrement PercentOverlap
 //AgtBgp4L2VpnRouteProfile SetRouteDistinguisherIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetRouteDistinguisherIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4L2VpnRouteProfile GetRouteDistinguisherIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetRouteDistinguisherAuto ()(string, error) {
 //parameters: RouteProfileHandle RouteDistinguisherType
 //AgtBgp4L2VpnRouteProfile GetRouteDistinguisherAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetPeerPoolHandle ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetPeerPoolHandle
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetState ()(string, error) {
 //parameters: RouteProfileInstance
 //AgtBgp4L2VpnRouteProfile GetState
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) IsEnabled () error {
 //parameters: RouteProfileInstance
 //AgtBgp4L2VpnRouteProfile IsEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) EnableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile EnableAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Enable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Enable4ByteAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile DisableAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Disable4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Disable4ByteAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile IsAsPathEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Is4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Is4ByteAsPathEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) EnableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile EnableAutoAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) EnableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile EnableAuto4ByteAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableAutoAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile DisableAutoAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableAuto4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile DisableAuto4ByteAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsAutoAsPathEnabled () error {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile IsAutoAsPathEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsAuto4ByteAsPathEnabled () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile IsAuto4ByteAsPathEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle AsPathAttributeType FirstAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4L2VpnRouteProfile SetAsPathIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsPathIncrementingRange () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType First4ByteAsNumber AsNumberIncrement AsNumberRepeat
 //AgtBgp4L2VpnRouteProfile Set4ByteAsPathIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile GetAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsPathIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsPathIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4L2VpnRouteProfile SetAsPathList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4L2VpnRouteProfile Set4ByteAsPathList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetAsPathList ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile GetAsPathList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsPathList ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsPathList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetAsPath () error {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumber
 //AgtBgp4L2VpnRouteProfile SetAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsPath () error {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumber
 //AgtBgp4L2VpnRouteProfile Set4ByteAsPath
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetAsPath ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile GetAsPath
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsPath ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsPath
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) AddAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsNumberList
 //AgtBgp4L2VpnRouteProfile AddAsPathToListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Add4ByteAsPathToListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsNumberList
 //AgtBgp4L2VpnRouteProfile Add4ByteAsPathToListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) RemoveAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4L2VpnRouteProfile RemoveAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Remove4ByteAsPathFromListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4L2VpnRouteProfile Remove4ByteAsPathFromListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) ClearAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType
 //AgtBgp4L2VpnRouteProfile ClearAsPathListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Clear4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType
 //AgtBgp4L2VpnRouteProfile Clear4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) ListAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeType AsPathListHandle
 //AgtBgp4L2VpnRouteProfile ListAsPathListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) List4ByteAsPathListOfLists ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeType AsPathListHandle
 //AgtBgp4L2VpnRouteProfile List4ByteAsPathListOfLists
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle AsPathAttributeTypeAuto
 //AgtBgp4L2VpnRouteProfile GetAsPathAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsPathAuto ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto
 //AgtBgp4L2VpnRouteProfile Get4ByteAsPathAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsPathAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle 4ByteAsPathAttributeTypeAuto FourByteFormatType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsPathAutoWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) EnablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4L2VpnRouteProfile EnablePathAttribute
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisablePathAttribute () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4L2VpnRouteProfile DisablePathAttribute
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsPathAttributeEnabled () error {
 //parameters: RouteProfileHandle PathAttribute
 //AgtBgp4L2VpnRouteProfile IsPathAttributeEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetOrigin () error {
 //parameters: RouteProfileHandle Origin
 //AgtBgp4L2VpnRouteProfile SetOrigin
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetOrigin ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetOrigin
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetLocalPreference () error {
 //parameters: RouteProfileHandle LocalPreference
 //AgtBgp4L2VpnRouteProfile SetLocalPreference
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetLocalPreference ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetLocalPreference
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetMultiExitDiscriminator () error {
 //parameters: RouteProfileHandle MultiExitDiscriminator
 //AgtBgp4L2VpnRouteProfile SetMultiExitDiscriminator
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetMultiExitDiscriminator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetMultiExitDiscriminator
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle FirstAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4L2VpnRouteProfile SetAggregatorIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregatorIncrementingRange () error {
 //parameters: RouteProfileHandle First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregatorIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregatorIncrementingRangeWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType First4ByteAsNumber AsNumberIncrement AsNumberRepeat FirstIpAddress PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregatorIncrementingRangeWithFormat
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregatorIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregatorIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregatorIncrementingRangeWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregatorIncrementingRangeWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4L2VpnRouteProfile SetAggregatorList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle AsNumberList IpAddressList
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType AsNumberList IpAddressList
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetAggregatorList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregatorList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregatorList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregatorListWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregatorListWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4L2VpnRouteProfile SetAggregator
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregator () error {
 //parameters: RouteProfileHandle AsNumber IpAddress
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregator
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Set4ByteAsAggregatorWithFormat () error {
 //parameters: RouteProfileHandle FourByteFormatType AsNumber IpAddress
 //AgtBgp4L2VpnRouteProfile Set4ByteAsAggregatorWithFormat
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetAggregator
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregator ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregator
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAsAggregatorWithFormat ()(string, error) {
 //parameters: RouteProfileHandle FourByteFormatType
 //AgtBgp4L2VpnRouteProfile Get4ByteAsAggregatorWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetAggregatorAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAggregatorAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Get4ByteAggregatorAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Get4ByteAggregatorAutoWithFormat ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Get4ByteAggregatorAutoWithFormat
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetOriginatorIdIncrementingRange () error {
 //parameters: RouteProfileHandle FirstOriginatorId PrefixLength AddressIncrement AddressRepeat
 //AgtBgp4L2VpnRouteProfile SetOriginatorIdIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetOriginatorIdIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetOriginatorIdIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle OriginatorIdList
 //AgtBgp4L2VpnRouteProfile SetOriginatorIdList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetOriginatorIdList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetOriginatorIdList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetOriginatorId () error {
 //parameters: RouteProfileHandle OriginatorId
 //AgtBgp4L2VpnRouteProfile SetOriginatorId
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetOriginatorId ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetOriginatorId
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetOriginatorIdAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetOriginatorIdAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetClusterList ()(string, error) {
 //parameters: RouteProfileHandle ClusterList
 //AgtBgp4L2VpnRouteProfile SetClusterList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetClusterList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetClusterList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetClusterListAuto ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetClusterListAuto
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetCommunities () error {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4L2VpnRouteProfile SetCommunities
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetCommunities ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetCommunities
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetCommunitiesIncrementingRange () error {
 //parameters: RouteProfileHandle CommunitiesList CommunityIncrement CommunityRepeat
 //AgtBgp4L2VpnRouteProfile SetCommunitiesIncrementingRange
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetCommunitiesIncrementingRange ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetCommunitiesIncrementingRange
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle CommunitiesList
 //AgtBgp4L2VpnRouteProfile SetCommunitiesList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetCommunitiesList ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetCommunitiesList
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) EnableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile EnableTrafficDestinations
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableTrafficDestinations () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile DisableTrafficDestinations
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsTrafficDestinationsEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile IsTrafficDestinationsEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetFlapDistribution () error {
 //parameters: RouteProfileHandle FlapPercentage FlapOffset
 //AgtBgp4L2VpnRouteProfile SetFlapDistribution
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetFlapDistribution ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetFlapDistribution
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) Advertise () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4L2VpnRouteProfile Advertise
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Withdraw () error {
 //parameters: RouteProfileIdentifiers
 //AgtBgp4L2VpnRouteProfile Withdraw
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Enable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Enable
 return nil
}

func(np *Bgp4L2VpnRouteProfile) Disable () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile Disable
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetExportRouteTargetType () error {
 //parameters: RouteProfileHandle RouteTargetTypeValue
 //AgtBgp4L2VpnRouteProfile SetExportRouteTargetType
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetExportRouteTargetType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetExportRouteTargetType
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) SetRouteDistinguisherType () error {
 //parameters: RouteProfileHandle RouteDistinguisherTypeValue
 //AgtBgp4L2VpnRouteProfile SetRouteDistinguisherType
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetRouteDistinguisherType ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetRouteDistinguisherType
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) EnableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile EnableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4L2VpnRouteProfile) DisableAutoRouteDistinguisher () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile DisableAutoRouteDistinguisher
 return nil
}

func(np *Bgp4L2VpnRouteProfile) IsAutoRouteDistinguisherEnabled () error {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile IsAutoRouteDistinguisherEnabled
 return nil
}

func(np *Bgp4L2VpnRouteProfile) SetVpnsPerPeer () error {
 //parameters: RouteProfileHandle NumVpnsPerPeer
 //AgtBgp4L2VpnRouteProfile SetVpnsPerPeer
 return nil
}

func(np *Bgp4L2VpnRouteProfile) GetVpnsPerPeer ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4L2VpnRouteProfile) GetTotalNumberOfVpns ()(string, error) {
 //parameters: RouteProfileHandle
 //AgtBgp4L2VpnRouteProfile GetTotalNumberOfVpns
 return "", nil
}

