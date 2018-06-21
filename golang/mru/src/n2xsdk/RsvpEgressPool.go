package n2xsdk

type RsvpEgressPool struct {
 Handler string
}

func(np *RsvpEgressPool) IsEgress () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsEgress
 return nil
}

func(np *RsvpEgressPool) ConfigureAsSingleEndPoint () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress NumberTunnels NumberLsps
 //AgtRsvpEgressPool ConfigureAsSingleEndPoint
 return nil
}

func(np *RsvpEgressPool) GetSingleEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetSingleEndPointDetails
 return "", nil
}

func(np *RsvpEgressPool) ConfigureAsMultipleEndPoint () error {
 //parameters: PoolHandle SourceIpAddressFirst SourceIpAddressIncrement SourceIpAddressNumber DestinationIpAddressFirst DestinationIpAddressIncrement DestinationIpAddressNumber PartialMesh
 //AgtRsvpEgressPool ConfigureAsMultipleEndPoint
 return nil
}

func(np *RsvpEgressPool) GetMultipleEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetMultipleEndPointDetails
 return "", nil
}

func(np *RsvpEgressPool) ConfigureAsGridEndPoint () error {
 //parameters: PoolHandle SourceIpAddressFirst SourceIpAddressNumberOfRows SourceIpAddressNumberOfColumns DestinationIpAddressFirst DestinationIpAddressNumberOfRows DestinationIpAddressNumberOfColumns
 //AgtRsvpEgressPool ConfigureAsGridEndPoint
 return nil
}

func(np *RsvpEgressPool) GetGridEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetGridEndPointDetails
 return "", nil
}

func(np *RsvpEgressPool) GetEndPointType ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetEndPointType
 return "", nil
}

func(np *RsvpEgressPool) SetFirstTunnelId () error {
 //parameters: PoolHandle FirstTunnelId
 //AgtRsvpEgressPool SetFirstTunnelId
 return nil
}

func(np *RsvpEgressPool) GetFirstTunnelId ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetFirstTunnelId
 return "", nil
}

func(np *RsvpEgressPool) SetFirstLspId () error {
 //parameters: PoolHandle FirstLspId
 //AgtRsvpEgressPool SetFirstLspId
 return nil
}

func(np *RsvpEgressPool) GetFirstLspId ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetFirstLspId
 return "", nil
}

func(np *RsvpEgressPool) SetBurstInterval () error {
 //parameters: PoolHandle BurstInterval
 //AgtRsvpEgressPool SetBurstInterval
 return nil
}

func(np *RsvpEgressPool) GetBurstInterval ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetBurstInterval
 return "", nil
}

func(np *RsvpEgressPool) SetBurstSize () error {
 //parameters: PoolHandle BurstSize
 //AgtRsvpEgressPool SetBurstSize
 return nil
}

func(np *RsvpEgressPool) GetBurstSize ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetBurstSize
 return "", nil
}

func(np *RsvpEgressPool) Close () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool Close
 return nil
}

func(np *RsvpEgressPool) CloseGracefully () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool CloseGracefully
 return nil
}

func(np *RsvpEgressPool) CloseRemotely () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool CloseRemotely
 return nil
}

func(np *RsvpEgressPool) CloseLsp () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpEgressPool CloseLsp
 return nil
}

func(np *RsvpEgressPool) CloseLspGracefully () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpEgressPool CloseLspGracefully
 return nil
}

func(np *RsvpEgressPool) CloseLspRemotely () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpEgressPool CloseLspRemotely
 return nil
}

func(np *RsvpEgressPool) CloseTunnel () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpEgressPool CloseTunnel
 return nil
}

func(np *RsvpEgressPool) CloseTunnelGracefully () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpEgressPool CloseTunnelGracefully
 return nil
}

func(np *RsvpEgressPool) CloseTunnelRemotely () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpEgressPool CloseTunnelRemotely
 return nil
}

func(np *RsvpEgressPool) GetState ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetState
 return "", nil
}

func(np *RsvpEgressPool) UploadTunnelDetails () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool UploadTunnelDetails
 return nil
}

func(np *RsvpEgressPool) CancelUpload () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool CancelUpload
 return nil
}

func(np *RsvpEgressPool) IsUploadComplete () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsUploadComplete
 return nil
}

func(np *RsvpEgressPool) SetLspsPerUploadPage () error {
 //parameters: PoolHandle LspsPerPage
 //AgtRsvpEgressPool SetLspsPerUploadPage
 return nil
}

func(np *RsvpEgressPool) GetLspsPerUploadPage ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetLspsPerUploadPage
 return "", nil
}

func(np *RsvpEgressPool) GetNumberOfUploadedPages ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetNumberOfUploadedPages
 return "", nil
}

func(np *RsvpEgressPool) GetTunnelLabels ()(string, error) {
 //parameters: PoolHandle PageNumber
 //AgtRsvpEgressPool GetTunnelLabels
 return "", nil
}

func(np *RsvpEgressPool) GetRequiredTunnelLabels ()(string, error) {
 //parameters: PoolHandle IpSourceAddressSize pRequiredSources IpDestinationAddressSize pRequiredDestinations TunnelIdSize pRequiredTunnelIds
 //AgtRsvpEgressPool GetRequiredTunnelLabels
 return "", nil
}

func(np *RsvpEgressPool) GetRequiredLspLabels ()(string, error) {
 //parameters: PoolHandle RequiredSourceAddress RequiredDestinationAddress RequiredTunnelId RequiredLspId
 //AgtRsvpEgressPool GetRequiredLspLabels
 return "", nil
}

func(np *RsvpEgressPool) GetTunnelLabelInformation ()(string, error) {
 //parameters: PoolHandle PageNumber
 //AgtRsvpEgressPool GetTunnelLabelInformation
 return "", nil
}

func(np *RsvpEgressPool) GetRequiredTunnelLabelInformation ()(string, error) {
 //parameters: PoolHandle IpSourceAddressSize pRequiredSources IpDestinationAddressSize pRequiredDestinations TunnelIdSize pRequiredTunnelIds
 //AgtRsvpEgressPool GetRequiredTunnelLabelInformation
 return "", nil
}

func(np *RsvpEgressPool) GetRequiredLspLabelInformation ()(string, error) {
 //parameters: PoolHandle RequiredSourceAddress RequiredDestinationAddress RequiredTunnelId RequiredLspId
 //AgtRsvpEgressPool GetRequiredLspLabelInformation
 return "", nil
}

func(np *RsvpEgressPool) SetBandwidthReservation () error {
 //parameters: PoolHandle Bandwidth
 //AgtRsvpEgressPool SetBandwidthReservation
 return nil
}

func(np *RsvpEgressPool) GetBandwidthReservation ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetBandwidthReservation
 return "", nil
}

func(np *RsvpEgressPool) SetPriority () error {
 //parameters: PoolHandle SetupPriority HoldPriority
 //AgtRsvpEgressPool SetPriority
 return nil
}

func(np *RsvpEgressPool) GetPriority ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetPriority
 return "", nil
}

func(np *RsvpEgressPool) SetAttributesFlag () error {
 //parameters: PoolHandle Flag
 //AgtRsvpEgressPool SetAttributesFlag
 return nil
}

func(np *RsvpEgressPool) GetAttributesFlag ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetAttributesFlag
 return "", nil
}

func(np *RsvpEgressPool) SetResourceAffinities () error {
 //parameters: PoolHandle ExcludeAny IncludeAny IncludeAll
 //AgtRsvpEgressPool SetResourceAffinities
 return nil
}

func(np *RsvpEgressPool) GetResourceAffinities ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetResourceAffinities
 return "", nil
}

func(np *RsvpEgressPool) EnableResourceAffinities () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableResourceAffinities
 return nil
}

func(np *RsvpEgressPool) DisableResourceAffinities () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableResourceAffinities
 return nil
}

func(np *RsvpEgressPool) IsResourceAffinitiesEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsResourceAffinitiesEnabled
 return nil
}

func(np *RsvpEgressPool) SetSessionName () error {
 //parameters: PoolHandle Name
 //AgtRsvpEgressPool SetSessionName
 return nil
}

func(np *RsvpEgressPool) GetSessionName ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetSessionName
 return "", nil
}

func(np *RsvpEgressPool) EnableMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) DisableMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) IsMakeBeforeBreakEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsMakeBeforeBreakEnabled
 return nil
}

func(np *RsvpEgressPool) EnableAutoMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableAutoMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) DisableAutoMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableAutoMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) IsAutoMakeBeforeBreakEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsAutoMakeBeforeBreakEnabled
 return nil
}

func(np *RsvpEgressPool) LspPoolMakeBeforeBreak () error {
 //parameters: PoolHandle LspType
 //AgtRsvpEgressPool LspPoolMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) LspMakeBeforeBreak () error {
 //parameters: PoolHandle AddressPoolIndex LspType
 //AgtRsvpEgressPool LspMakeBeforeBreak
 return nil
}

func(np *RsvpEgressPool) RemoveRsvpEroSubobject () error {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool RemoveRsvpEroSubobject
 return nil
}

func(np *RsvpEgressPool) RemoveAllRsvpEroSubobjects () error {
 //parameters: PoolHandle LspType
 //AgtRsvpEgressPool RemoveAllRsvpEroSubobjects
 return nil
}

func(np *RsvpEgressPool) SetRsvpEroSubobject () error {
 //parameters: PoolHandle LspType Index EroSubobjectType IsLoose Count psaContents
 //AgtRsvpEgressPool SetRsvpEroSubobject
 return nil
}

func(np *RsvpEgressPool) SetRsvpEroIpv4Subobject () error {
 //parameters: PoolHandle LspType Index IsLoose Address Prefix
 //AgtRsvpEgressPool SetRsvpEroIpv4Subobject
 return nil
}

func(np *RsvpEgressPool) SetRsvpEroLabelSubobject () error {
 //parameters: PoolHandle LspType Index Upstream CType IsLoose Label
 //AgtRsvpEgressPool SetRsvpEroLabelSubobject
 return nil
}

func(np *RsvpEgressPool) SetRsvpEroAsNumberSubobject () error {
 //parameters: PoolHandle LspType Index IsLoose AsNumber
 //AgtRsvpEgressPool SetRsvpEroAsNumberSubobject
 return nil
}

func(np *RsvpEgressPool) SetRsvpEroUnnumberedSubobject () error {
 //parameters: PoolHandle LspType Index IsLoose RouterId InterfaceId
 //AgtRsvpEgressPool SetRsvpEroUnnumberedSubobject
 return nil
}

func(np *RsvpEgressPool) GetRsvpEroSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroSubobject
 return "", nil
}

func(np *RsvpEgressPool) GetRsvpEroIpv4Subobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroIpv4Subobject
 return "", nil
}

func(np *RsvpEgressPool) GetRsvpEroLabelSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroLabelSubobject
 return "", nil
}

func(np *RsvpEgressPool) GetRsvpEroAsNumberSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroAsNumberSubobject
 return "", nil
}

func(np *RsvpEgressPool) GetRsvpEroUnnumberedSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroUnnumberedSubobject
 return "", nil
}

func(np *RsvpEgressPool) ListRsvpEroSubobjects ()(string, error) {
 //parameters: PoolHandle LspType
 //AgtRsvpEgressPool ListRsvpEroSubobjects
 return "", nil
}

func(np *RsvpEgressPool) GetNumberOfRsvpEroSubobjects ()(string, error) {
 //parameters: PoolHandle LspType
 //AgtRsvpEgressPool GetNumberOfRsvpEroSubobjects
 return "", nil
}

func(np *RsvpEgressPool) GetRsvpEroSubobjectType ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpEgressPool GetRsvpEroSubobjectType
 return "", nil
}

func(np *RsvpEgressPool) SetTSpecIeeeParameter () error {
 //parameters: PoolHandle TSpecParameter FloatValue
 //AgtRsvpEgressPool SetTSpecIeeeParameter
 return nil
}

func(np *RsvpEgressPool) GetTSpecIeeeParameter ()(string, error) {
 //parameters: PoolHandle TSpecParameter
 //AgtRsvpEgressPool GetTSpecIeeeParameter
 return "", nil
}

func(np *RsvpEgressPool) SetTSpecParameter () error {
 //parameters: PoolHandle TSpecParameter Value
 //AgtRsvpEgressPool SetTSpecParameter
 return nil
}

func(np *RsvpEgressPool) GetTSpecParameter ()(string, error) {
 //parameters: PoolHandle TSpecParameter
 //AgtRsvpEgressPool GetTSpecParameter
 return "", nil
}

func(np *RsvpEgressPool) SetAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter Value
 //AgtRsvpEgressPool SetAdSpecParameter
 return nil
}

func(np *RsvpEgressPool) GetAdSpecParameter ()(string, error) {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpEgressPool GetAdSpecParameter
 return "", nil
}

func(np *RsvpEgressPool) SetProtectionParameter () error {
 //parameters: PoolHandle ProtectionParameter Value
 //AgtRsvpEgressPool SetProtectionParameter
 return nil
}

func(np *RsvpEgressPool) GetProtectionParameter ()(string, error) {
 //parameters: PoolHandle ProtectionParameter
 //AgtRsvpEgressPool GetProtectionParameter
 return "", nil
}

func(np *RsvpEgressPool) SetGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter Value
 //AgtRsvpEgressPool SetGeneralizedUniParameter
 return nil
}

func(np *RsvpEgressPool) GetGeneralizedUniParameter ()(string, error) {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpEgressPool GetGeneralizedUniParameter
 return "", nil
}

func(np *RsvpEgressPool) EnableProtectionObject () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableProtectionObject
 return nil
}

func(np *RsvpEgressPool) DisableProtectionObject () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableProtectionObject
 return nil
}

func(np *RsvpEgressPool) IsProtectionObjectEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsProtectionObjectEnabled
 return nil
}

func(np *RsvpEgressPool) SelectAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpEgressPool SelectAdSpecParameter
 return nil
}

func(np *RsvpEgressPool) DeselectAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpEgressPool DeselectAdSpecParameter
 return nil
}

func(np *RsvpEgressPool) IsAdSpecParameterSelected () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpEgressPool IsAdSpecParameterSelected
 return nil
}

func(np *RsvpEgressPool) SelectGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpEgressPool SelectGeneralizedUniParameter
 return nil
}

func(np *RsvpEgressPool) DeselectGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpEgressPool DeselectGeneralizedUniParameter
 return nil
}

func(np *RsvpEgressPool) IsGeneralizedUniParameterSelected () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpEgressPool IsGeneralizedUniParameterSelected
 return nil
}

func(np *RsvpEgressPool) SetDiffServObjectType () error {
 //parameters: PoolHandle DiffServObjectType
 //AgtRsvpEgressPool SetDiffServObjectType
 return nil
}

func(np *RsvpEgressPool) GetDiffServObjectType ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetDiffServObjectType
 return "", nil
}

func(np *RsvpEgressPool) EnableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableDiffServXBit
 return nil
}

func(np *RsvpEgressPool) DisableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableDiffServXBit
 return nil
}

func(np *RsvpEgressPool) IsDiffServXBitEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsDiffServXBitEnabled
 return nil
}

func(np *RsvpEgressPool) SetDiffServPhbEncodingStyle () error {
 //parameters: PoolHandle EncodingStyle
 //AgtRsvpEgressPool SetDiffServPhbEncodingStyle
 return nil
}

func(np *RsvpEgressPool) GetDiffServPhbEncodingStyle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetDiffServPhbEncodingStyle
 return "", nil
}

func(np *RsvpEgressPool) AddExpPhbIdMap () error {
 //parameters: PoolHandle Exp PhbId
 //AgtRsvpEgressPool AddExpPhbIdMap
 return nil
}

func(np *RsvpEgressPool) RemoveExpPhbIdMap () error {
 //parameters: PoolHandle Exp PhbId
 //AgtRsvpEgressPool RemoveExpPhbIdMap
 return nil
}

func(np *RsvpEgressPool) ListExpPhbIdMap ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool ListExpPhbIdMap
 return "", nil
}

func(np *RsvpEgressPool) SetPhbSchedulingClass () error {
 //parameters: PoolHandle Psc
 //AgtRsvpEgressPool SetPhbSchedulingClass
 return nil
}

func(np *RsvpEgressPool) GetPhbSchedulingClass ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetPhbSchedulingClass
 return "", nil
}

func(np *RsvpEgressPool) SetPathObjectOrder () error {
 //parameters: PoolHandle Count psaObjects
 //AgtRsvpEgressPool SetPathObjectOrder
 return nil
}

func(np *RsvpEgressPool) GetPathObjectOrder ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetPathObjectOrder
 return "", nil
}

func(np *RsvpEgressPool) SetDefaultPathObjectOrder () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool SetDefaultPathObjectOrder
 return nil
}

func(np *RsvpEgressPool) IsDefaultPathObjectOrder () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsDefaultPathObjectOrder
 return nil
}

func(np *RsvpEgressPool) EnableUpstreamLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool EnableUpstreamLabelOverride
 return nil
}

func(np *RsvpEgressPool) DisableUpstreamLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool DisableUpstreamLabelOverride
 return nil
}

func(np *RsvpEgressPool) IsUpstreamLabelOverrideEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsUpstreamLabelOverrideEnabled
 return nil
}

func(np *RsvpEgressPool) SetUnnumberedUpstreamLabel () error {
 //parameters: PoolHandle UpstreamLabel LocalInterfaceId
 //AgtRsvpEgressPool SetUnnumberedUpstreamLabel
 return nil
}

func(np *RsvpEgressPool) GetUnnumberedUpstreamLabel ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetUnnumberedUpstreamLabel
 return "", nil
}

func(np *RsvpEgressPool) SetNumberedUpstreamLabel () error {
 //parameters: PoolHandle UpstreamLabel LocalInterfaceId
 //AgtRsvpEgressPool SetNumberedUpstreamLabel
 return nil
}

func(np *RsvpEgressPool) GetNumberedUpstreamLabel ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetNumberedUpstreamLabel
 return "", nil
}

func(np *RsvpEgressPool) GetPortHandle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetPortHandle
 return "", nil
}

func(np *RsvpEgressPool) GetSessionHandle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpEgressPool GetSessionHandle
 return "", nil
}

func(np *RsvpEgressPool) Enable () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool Enable
 return nil
}

func(np *RsvpEgressPool) Disable () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool Disable
 return nil
}

func(np *RsvpEgressPool) IsEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsEnabled
 return nil
}

func(np *RsvpEgressPool) IsDefault () error {
 //parameters: PoolHandle
 //AgtRsvpEgressPool IsDefault
 return nil
}

