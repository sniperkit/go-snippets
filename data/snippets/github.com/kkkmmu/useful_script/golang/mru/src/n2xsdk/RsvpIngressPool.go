package n2xsdk

type RsvpIngressPool struct {
 Handler string
}

func(np *RsvpIngressPool) IsEgress () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsEgress
 return nil
}

func(np *RsvpIngressPool) ConfigureAsSingleEndPoint () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress NumberTunnels NumberLsps
 //AgtRsvpIngressPool ConfigureAsSingleEndPoint
 return nil
}

func(np *RsvpIngressPool) GetSingleEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetSingleEndPointDetails
 return "", nil
}

func(np *RsvpIngressPool) ConfigureAsMultipleEndPoint () error {
 //parameters: PoolHandle SourceIpAddressFirst SourceIpAddressIncrement SourceIpAddressNumber DestinationIpAddressFirst DestinationIpAddressIncrement DestinationIpAddressNumber PartialMesh
 //AgtRsvpIngressPool ConfigureAsMultipleEndPoint
 return nil
}

func(np *RsvpIngressPool) GetMultipleEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetMultipleEndPointDetails
 return "", nil
}

func(np *RsvpIngressPool) ConfigureAsGridEndPoint () error {
 //parameters: PoolHandle SourceIpAddressFirst SourceIpAddressNumberOfRows SourceIpAddressNumberOfColumns DestinationIpAddressFirst DestinationIpAddressNumberOfRows DestinationIpAddressNumberOfColumns
 //AgtRsvpIngressPool ConfigureAsGridEndPoint
 return nil
}

func(np *RsvpIngressPool) GetGridEndPointDetails ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetGridEndPointDetails
 return "", nil
}

func(np *RsvpIngressPool) GetEndPointType ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetEndPointType
 return "", nil
}

func(np *RsvpIngressPool) SetFirstTunnelId () error {
 //parameters: PoolHandle FirstTunnelId
 //AgtRsvpIngressPool SetFirstTunnelId
 return nil
}

func(np *RsvpIngressPool) GetFirstTunnelId ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetFirstTunnelId
 return "", nil
}

func(np *RsvpIngressPool) SetFirstLspId () error {
 //parameters: PoolHandle FirstLspId
 //AgtRsvpIngressPool SetFirstLspId
 return nil
}

func(np *RsvpIngressPool) GetFirstLspId ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetFirstLspId
 return "", nil
}

func(np *RsvpIngressPool) SetBurstInterval () error {
 //parameters: PoolHandle BurstInterval
 //AgtRsvpIngressPool SetBurstInterval
 return nil
}

func(np *RsvpIngressPool) GetBurstInterval ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetBurstInterval
 return "", nil
}

func(np *RsvpIngressPool) SetBurstSize () error {
 //parameters: PoolHandle BurstSize
 //AgtRsvpIngressPool SetBurstSize
 return nil
}

func(np *RsvpIngressPool) GetBurstSize ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetBurstSize
 return "", nil
}

func(np *RsvpIngressPool) Close () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool Close
 return nil
}

func(np *RsvpIngressPool) CloseGracefully () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool CloseGracefully
 return nil
}

func(np *RsvpIngressPool) CloseRemotely () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool CloseRemotely
 return nil
}

func(np *RsvpIngressPool) CloseLsp () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpIngressPool CloseLsp
 return nil
}

func(np *RsvpIngressPool) CloseLspGracefully () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpIngressPool CloseLspGracefully
 return nil
}

func(np *RsvpIngressPool) CloseLspRemotely () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpIngressPool CloseLspRemotely
 return nil
}

func(np *RsvpIngressPool) CloseTunnel () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpIngressPool CloseTunnel
 return nil
}

func(np *RsvpIngressPool) CloseTunnelGracefully () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpIngressPool CloseTunnelGracefully
 return nil
}

func(np *RsvpIngressPool) CloseTunnelRemotely () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpIngressPool CloseTunnelRemotely
 return nil
}

func(np *RsvpIngressPool) GetState ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetState
 return "", nil
}

func(np *RsvpIngressPool) UploadTunnelDetails () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool UploadTunnelDetails
 return nil
}

func(np *RsvpIngressPool) CancelUpload () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool CancelUpload
 return nil
}

func(np *RsvpIngressPool) IsUploadComplete () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsUploadComplete
 return nil
}

func(np *RsvpIngressPool) SetLspsPerUploadPage () error {
 //parameters: PoolHandle LspsPerPage
 //AgtRsvpIngressPool SetLspsPerUploadPage
 return nil
}

func(np *RsvpIngressPool) GetLspsPerUploadPage ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetLspsPerUploadPage
 return "", nil
}

func(np *RsvpIngressPool) GetNumberOfUploadedPages ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetNumberOfUploadedPages
 return "", nil
}

func(np *RsvpIngressPool) GetTunnelLabels ()(string, error) {
 //parameters: PoolHandle PageNumber
 //AgtRsvpIngressPool GetTunnelLabels
 return "", nil
}

func(np *RsvpIngressPool) GetRequiredTunnelLabels ()(string, error) {
 //parameters: PoolHandle IpSourceAddressSize pRequiredSources IpDestinationAddressSize pRequiredDestinations TunnelIdSize pRequiredTunnelIds
 //AgtRsvpIngressPool GetRequiredTunnelLabels
 return "", nil
}

func(np *RsvpIngressPool) GetRequiredLspLabels ()(string, error) {
 //parameters: PoolHandle RequiredSourceAddress RequiredDestinationAddress RequiredTunnelId RequiredLspId
 //AgtRsvpIngressPool GetRequiredLspLabels
 return "", nil
}

func(np *RsvpIngressPool) GetTunnelLabelInformation ()(string, error) {
 //parameters: PoolHandle PageNumber
 //AgtRsvpIngressPool GetTunnelLabelInformation
 return "", nil
}

func(np *RsvpIngressPool) GetRequiredTunnelLabelInformation ()(string, error) {
 //parameters: PoolHandle IpSourceAddressSize pRequiredSources IpDestinationAddressSize pRequiredDestinations TunnelIdSize pRequiredTunnelIds
 //AgtRsvpIngressPool GetRequiredTunnelLabelInformation
 return "", nil
}

func(np *RsvpIngressPool) GetRequiredLspLabelInformation ()(string, error) {
 //parameters: PoolHandle RequiredSourceAddress RequiredDestinationAddress RequiredTunnelId RequiredLspId
 //AgtRsvpIngressPool GetRequiredLspLabelInformation
 return "", nil
}

func(np *RsvpIngressPool) SetBandwidthReservation () error {
 //parameters: PoolHandle Bandwidth
 //AgtRsvpIngressPool SetBandwidthReservation
 return nil
}

func(np *RsvpIngressPool) GetBandwidthReservation ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetBandwidthReservation
 return "", nil
}

func(np *RsvpIngressPool) SetPriority () error {
 //parameters: PoolHandle SetupPriority HoldPriority
 //AgtRsvpIngressPool SetPriority
 return nil
}

func(np *RsvpIngressPool) GetPriority ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetPriority
 return "", nil
}

func(np *RsvpIngressPool) SetAttributesFlag () error {
 //parameters: PoolHandle Flag
 //AgtRsvpIngressPool SetAttributesFlag
 return nil
}

func(np *RsvpIngressPool) GetAttributesFlag ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetAttributesFlag
 return "", nil
}

func(np *RsvpIngressPool) SetResourceAffinities () error {
 //parameters: PoolHandle ExcludeAny IncludeAny IncludeAll
 //AgtRsvpIngressPool SetResourceAffinities
 return nil
}

func(np *RsvpIngressPool) GetResourceAffinities ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetResourceAffinities
 return "", nil
}

func(np *RsvpIngressPool) EnableResourceAffinities () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableResourceAffinities
 return nil
}

func(np *RsvpIngressPool) DisableResourceAffinities () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableResourceAffinities
 return nil
}

func(np *RsvpIngressPool) IsResourceAffinitiesEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsResourceAffinitiesEnabled
 return nil
}

func(np *RsvpIngressPool) SetSessionName () error {
 //parameters: PoolHandle Name
 //AgtRsvpIngressPool SetSessionName
 return nil
}

func(np *RsvpIngressPool) GetSessionName ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetSessionName
 return "", nil
}

func(np *RsvpIngressPool) EnableMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) DisableMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) IsMakeBeforeBreakEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsMakeBeforeBreakEnabled
 return nil
}

func(np *RsvpIngressPool) EnableAutoMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableAutoMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) DisableAutoMakeBeforeBreak () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableAutoMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) IsAutoMakeBeforeBreakEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsAutoMakeBeforeBreakEnabled
 return nil
}

func(np *RsvpIngressPool) LspPoolMakeBeforeBreak () error {
 //parameters: PoolHandle LspType
 //AgtRsvpIngressPool LspPoolMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) LspMakeBeforeBreak () error {
 //parameters: PoolHandle AddressPoolIndex LspType
 //AgtRsvpIngressPool LspMakeBeforeBreak
 return nil
}

func(np *RsvpIngressPool) RemoveRsvpEroSubobject () error {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool RemoveRsvpEroSubobject
 return nil
}

func(np *RsvpIngressPool) RemoveAllRsvpEroSubobjects () error {
 //parameters: PoolHandle LspType
 //AgtRsvpIngressPool RemoveAllRsvpEroSubobjects
 return nil
}

func(np *RsvpIngressPool) SetRsvpEroSubobject () error {
 //parameters: PoolHandle LspType Index EroSubobjectType IsLoose Count psaContents
 //AgtRsvpIngressPool SetRsvpEroSubobject
 return nil
}

func(np *RsvpIngressPool) SetRsvpEroIpv4Subobject () error {
 //parameters: PoolHandle LspType Index IsLoose Address Prefix
 //AgtRsvpIngressPool SetRsvpEroIpv4Subobject
 return nil
}

func(np *RsvpIngressPool) SetRsvpEroLabelSubobject () error {
 //parameters: PoolHandle LspType Index Upstream CType IsLoose Label
 //AgtRsvpIngressPool SetRsvpEroLabelSubobject
 return nil
}

func(np *RsvpIngressPool) SetRsvpEroAsNumberSubobject () error {
 //parameters: PoolHandle LspType Index IsLoose AsNumber
 //AgtRsvpIngressPool SetRsvpEroAsNumberSubobject
 return nil
}

func(np *RsvpIngressPool) SetRsvpEroUnnumberedSubobject () error {
 //parameters: PoolHandle LspType Index IsLoose RouterId InterfaceId
 //AgtRsvpIngressPool SetRsvpEroUnnumberedSubobject
 return nil
}

func(np *RsvpIngressPool) GetRsvpEroSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroSubobject
 return "", nil
}

func(np *RsvpIngressPool) GetRsvpEroIpv4Subobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroIpv4Subobject
 return "", nil
}

func(np *RsvpIngressPool) GetRsvpEroLabelSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroLabelSubobject
 return "", nil
}

func(np *RsvpIngressPool) GetRsvpEroAsNumberSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroAsNumberSubobject
 return "", nil
}

func(np *RsvpIngressPool) GetRsvpEroUnnumberedSubobject ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroUnnumberedSubobject
 return "", nil
}

func(np *RsvpIngressPool) ListRsvpEroSubobjects ()(string, error) {
 //parameters: PoolHandle LspType
 //AgtRsvpIngressPool ListRsvpEroSubobjects
 return "", nil
}

func(np *RsvpIngressPool) GetNumberOfRsvpEroSubobjects ()(string, error) {
 //parameters: PoolHandle LspType
 //AgtRsvpIngressPool GetNumberOfRsvpEroSubobjects
 return "", nil
}

func(np *RsvpIngressPool) GetRsvpEroSubobjectType ()(string, error) {
 //parameters: PoolHandle LspType Index
 //AgtRsvpIngressPool GetRsvpEroSubobjectType
 return "", nil
}

func(np *RsvpIngressPool) SetTSpecIeeeParameter () error {
 //parameters: PoolHandle TSpecParameter FloatValue
 //AgtRsvpIngressPool SetTSpecIeeeParameter
 return nil
}

func(np *RsvpIngressPool) GetTSpecIeeeParameter ()(string, error) {
 //parameters: PoolHandle TSpecParameter
 //AgtRsvpIngressPool GetTSpecIeeeParameter
 return "", nil
}

func(np *RsvpIngressPool) SetTSpecParameter () error {
 //parameters: PoolHandle TSpecParameter Value
 //AgtRsvpIngressPool SetTSpecParameter
 return nil
}

func(np *RsvpIngressPool) GetTSpecParameter ()(string, error) {
 //parameters: PoolHandle TSpecParameter
 //AgtRsvpIngressPool GetTSpecParameter
 return "", nil
}

func(np *RsvpIngressPool) SetAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter Value
 //AgtRsvpIngressPool SetAdSpecParameter
 return nil
}

func(np *RsvpIngressPool) GetAdSpecParameter ()(string, error) {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpIngressPool GetAdSpecParameter
 return "", nil
}

func(np *RsvpIngressPool) SetProtectionParameter () error {
 //parameters: PoolHandle ProtectionParameter Value
 //AgtRsvpIngressPool SetProtectionParameter
 return nil
}

func(np *RsvpIngressPool) GetProtectionParameter ()(string, error) {
 //parameters: PoolHandle ProtectionParameter
 //AgtRsvpIngressPool GetProtectionParameter
 return "", nil
}

func(np *RsvpIngressPool) SetGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter Value
 //AgtRsvpIngressPool SetGeneralizedUniParameter
 return nil
}

func(np *RsvpIngressPool) GetGeneralizedUniParameter ()(string, error) {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpIngressPool GetGeneralizedUniParameter
 return "", nil
}

func(np *RsvpIngressPool) EnableProtectionObject () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableProtectionObject
 return nil
}

func(np *RsvpIngressPool) DisableProtectionObject () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableProtectionObject
 return nil
}

func(np *RsvpIngressPool) IsProtectionObjectEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsProtectionObjectEnabled
 return nil
}

func(np *RsvpIngressPool) SelectAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpIngressPool SelectAdSpecParameter
 return nil
}

func(np *RsvpIngressPool) DeselectAdSpecParameter () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpIngressPool DeselectAdSpecParameter
 return nil
}

func(np *RsvpIngressPool) IsAdSpecParameterSelected () error {
 //parameters: PoolHandle AdSpecParameter
 //AgtRsvpIngressPool IsAdSpecParameterSelected
 return nil
}

func(np *RsvpIngressPool) SelectGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpIngressPool SelectGeneralizedUniParameter
 return nil
}

func(np *RsvpIngressPool) DeselectGeneralizedUniParameter () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpIngressPool DeselectGeneralizedUniParameter
 return nil
}

func(np *RsvpIngressPool) IsGeneralizedUniParameterSelected () error {
 //parameters: PoolHandle GeneralizedUniParameter
 //AgtRsvpIngressPool IsGeneralizedUniParameterSelected
 return nil
}

func(np *RsvpIngressPool) SetDiffServObjectType () error {
 //parameters: PoolHandle DiffServObjectType
 //AgtRsvpIngressPool SetDiffServObjectType
 return nil
}

func(np *RsvpIngressPool) GetDiffServObjectType ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetDiffServObjectType
 return "", nil
}

func(np *RsvpIngressPool) EnableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableDiffServXBit
 return nil
}

func(np *RsvpIngressPool) DisableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableDiffServXBit
 return nil
}

func(np *RsvpIngressPool) IsDiffServXBitEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsDiffServXBitEnabled
 return nil
}

func(np *RsvpIngressPool) SetDiffServPhbEncodingStyle () error {
 //parameters: PoolHandle EncodingStyle
 //AgtRsvpIngressPool SetDiffServPhbEncodingStyle
 return nil
}

func(np *RsvpIngressPool) GetDiffServPhbEncodingStyle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetDiffServPhbEncodingStyle
 return "", nil
}

func(np *RsvpIngressPool) AddExpPhbIdMap () error {
 //parameters: PoolHandle Exp PhbId
 //AgtRsvpIngressPool AddExpPhbIdMap
 return nil
}

func(np *RsvpIngressPool) RemoveExpPhbIdMap () error {
 //parameters: PoolHandle Exp PhbId
 //AgtRsvpIngressPool RemoveExpPhbIdMap
 return nil
}

func(np *RsvpIngressPool) ListExpPhbIdMap ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool ListExpPhbIdMap
 return "", nil
}

func(np *RsvpIngressPool) SetPhbSchedulingClass () error {
 //parameters: PoolHandle Psc
 //AgtRsvpIngressPool SetPhbSchedulingClass
 return nil
}

func(np *RsvpIngressPool) GetPhbSchedulingClass ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetPhbSchedulingClass
 return "", nil
}

func(np *RsvpIngressPool) SetPathObjectOrder () error {
 //parameters: PoolHandle Count psaObjects
 //AgtRsvpIngressPool SetPathObjectOrder
 return nil
}

func(np *RsvpIngressPool) GetPathObjectOrder ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetPathObjectOrder
 return "", nil
}

func(np *RsvpIngressPool) SetDefaultPathObjectOrder () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool SetDefaultPathObjectOrder
 return nil
}

func(np *RsvpIngressPool) IsDefaultPathObjectOrder () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsDefaultPathObjectOrder
 return nil
}

func(np *RsvpIngressPool) EnableUpstreamLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableUpstreamLabelOverride
 return nil
}

func(np *RsvpIngressPool) DisableUpstreamLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableUpstreamLabelOverride
 return nil
}

func(np *RsvpIngressPool) IsUpstreamLabelOverrideEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsUpstreamLabelOverrideEnabled
 return nil
}

func(np *RsvpIngressPool) SetUnnumberedUpstreamLabel () error {
 //parameters: PoolHandle UpstreamLabel LocalInterfaceId
 //AgtRsvpIngressPool SetUnnumberedUpstreamLabel
 return nil
}

func(np *RsvpIngressPool) GetUnnumberedUpstreamLabel ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetUnnumberedUpstreamLabel
 return "", nil
}

func(np *RsvpIngressPool) SetNumberedUpstreamLabel () error {
 //parameters: PoolHandle UpstreamLabel LocalInterfaceId
 //AgtRsvpIngressPool SetNumberedUpstreamLabel
 return nil
}

func(np *RsvpIngressPool) GetNumberedUpstreamLabel ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetNumberedUpstreamLabel
 return "", nil
}

func(np *RsvpIngressPool) GetPortHandle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetPortHandle
 return "", nil
}

func(np *RsvpIngressPool) GetSessionHandle ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetSessionHandle
 return "", nil
}

func(np *RsvpIngressPool) Open () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool Open
 return nil
}

func(np *RsvpIngressPool) OpenTunnel () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId
 //AgtRsvpIngressPool OpenTunnel
 return nil
}

func(np *RsvpIngressPool) OpenLsp () error {
 //parameters: PoolHandle IpSourceAddress IpDestinationAddress TunnelId LspId
 //AgtRsvpIngressPool OpenLsp
 return nil
}

func(np *RsvpIngressPool) GetTimeToOpen ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetTimeToOpen
 return "", nil
}

func(np *RsvpIngressPool) EnableEgressLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool EnableEgressLabelOverride
 return nil
}

func(np *RsvpIngressPool) DisableEgressLabelOverride () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool DisableEgressLabelOverride
 return nil
}

func(np *RsvpIngressPool) IsEgressLabelOverrideEnabled () error {
 //parameters: PoolHandle
 //AgtRsvpIngressPool IsEgressLabelOverrideEnabled
 return nil
}

func(np *RsvpIngressPool) SetEgressLabelOverrideStart () error {
 //parameters: PoolHandle LabelOverrideMode LabelOverrideStart
 //AgtRsvpIngressPool SetEgressLabelOverrideStart
 return nil
}

func(np *RsvpIngressPool) GetEgressLabelOverrideStart ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetEgressLabelOverrideStart
 return "", nil
}

func(np *RsvpIngressPool) GetEgressLabelOverrideMode ()(string, error) {
 //parameters: PoolHandle
 //AgtRsvpIngressPool GetEgressLabelOverrideMode
 return "", nil
}

func(np *RsvpIngressPool) OpenPools () error {
 //parameters: Count psaPoolHandles
 //AgtRsvpIngressPool OpenPools
 return nil
}

func(np *RsvpIngressPool) OpenAll () error {
 //parameters: 
 //AgtRsvpIngressPool OpenAll
 return nil
}

func(np *RsvpIngressPool) StopOpenPools () error {
 //parameters: 
 //AgtRsvpIngressPool StopOpenPools
 return nil
}

func(np *RsvpIngressPool) IsOpenPoolsInProgress () error {
 //parameters: 
 //AgtRsvpIngressPool IsOpenPoolsInProgress
 return nil
}

