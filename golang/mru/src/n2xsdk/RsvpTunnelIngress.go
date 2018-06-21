package n2xsdk

type RsvpTunnelIngress struct {
 Handler string
}

func(np *RsvpTunnelIngress) SetTunnelId () error {
 //parameters: TunnelHandle TunnelId
 //AgtRsvpTunnelIngress SetTunnelId
 return nil
}

func(np *RsvpTunnelIngress) GetTunnelId ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetTunnelId
 return "", nil
}

func(np *RsvpTunnelIngress) SetSourceIpAddress () error {
 //parameters: TunnelHandle IpAddress
 //AgtRsvpTunnelIngress SetSourceIpAddress
 return nil
}

func(np *RsvpTunnelIngress) GetSourceIpAddress ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetSourceIpAddress
 return "", nil
}

func(np *RsvpTunnelIngress) SetDestinationIpAddress () error {
 //parameters: TunnelHandle IpAddress
 //AgtRsvpTunnelIngress SetDestinationIpAddress
 return nil
}

func(np *RsvpTunnelIngress) GetDestinationIpAddress ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetDestinationIpAddress
 return "", nil
}

func(np *RsvpTunnelIngress) SetBandwidthReservation () error {
 //parameters: TunnelHandle Bandwidth
 //AgtRsvpTunnelIngress SetBandwidthReservation
 return nil
}

func(np *RsvpTunnelIngress) GetBandwidthReservation ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetBandwidthReservation
 return "", nil
}

func(np *RsvpTunnelIngress) SetBandwidthReservationIeee () error {
 //parameters: TunnelHandle IeeeBandwidth
 //AgtRsvpTunnelIngress SetBandwidthReservationIeee
 return nil
}

func(np *RsvpTunnelIngress) GetBandwidthReservationIeee ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetBandwidthReservationIeee
 return "", nil
}

func(np *RsvpTunnelIngress) SetPriority () error {
 //parameters: TunnelHandle SetupPriority HoldPriority
 //AgtRsvpTunnelIngress SetPriority
 return nil
}

func(np *RsvpTunnelIngress) GetPriority ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetPriority
 return "", nil
}

func(np *RsvpTunnelIngress) SetAttributesFlag () error {
 //parameters: TunnelHandle Flag
 //AgtRsvpTunnelIngress SetAttributesFlag
 return nil
}

func(np *RsvpTunnelIngress) GetAttributesFlag ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetAttributesFlag
 return "", nil
}

func(np *RsvpTunnelIngress) SetResourceAffinities () error {
 //parameters: TunnelHandle ExcludeAny IncludeAny IncludeAll
 //AgtRsvpTunnelIngress SetResourceAffinities
 return nil
}

func(np *RsvpTunnelIngress) GetResourceAffinities ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetResourceAffinities
 return "", nil
}

func(np *RsvpTunnelIngress) EnableResourceAffinities () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress EnableResourceAffinities
 return nil
}

func(np *RsvpTunnelIngress) DisableResourceAffinities () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress DisableResourceAffinities
 return nil
}

func(np *RsvpTunnelIngress) IsResourceAffinitiesEnabled () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress IsResourceAffinitiesEnabled
 return nil
}

func(np *RsvpTunnelIngress) SetSessionName () error {
 //parameters: TunnelHandle Name
 //AgtRsvpTunnelIngress SetSessionName
 return nil
}

func(np *RsvpTunnelIngress) GetSessionName ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetSessionName
 return "", nil
}

func(np *RsvpTunnelIngress) SetExplicitRouteObject () error {
 //parameters: TunnelHandle Count psaRouteList
 //AgtRsvpTunnelIngress SetExplicitRouteObject
 return nil
}

func(np *RsvpTunnelIngress) GetExplicitRouteObject ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetExplicitRouteObject
 return "", nil
}

func(np *RsvpTunnelIngress) AddEroSubobject () error {
 //parameters: TunnelHandle Type
 //AgtRsvpTunnelIngress AddEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) InsertEroSubobject () error {
 //parameters: TunnelHandle Index Type
 //AgtRsvpTunnelIngress InsertEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) RemoveEroSubobject () error {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress RemoveEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) RemoveAllEroSubobjects () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress RemoveAllEroSubobjects
 return nil
}

func(np *RsvpTunnelIngress) SetEroSubobject () error {
 //parameters: TunnelHandle Index Type IsLoose Count psaContents
 //AgtRsvpTunnelIngress SetEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetEroIpv4Subobject () error {
 //parameters: TunnelHandle Index IsLoose Address Prefix
 //AgtRsvpTunnelIngress SetEroIpv4Subobject
 return nil
}

func(np *RsvpTunnelIngress) SetEroLabelSubobject () error {
 //parameters: TunnelHandle Index Upstream CType IsLoose Label
 //AgtRsvpTunnelIngress SetEroLabelSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetEroAsNumberSubobject () error {
 //parameters: TunnelHandle Index IsLoose AsNumber
 //AgtRsvpTunnelIngress SetEroAsNumberSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetEroUnnumberedSubobject () error {
 //parameters: TunnelHandle Index IsLoose RouterId InterfaceId
 //AgtRsvpTunnelIngress SetEroUnnumberedSubobject
 return nil
}

func(np *RsvpTunnelIngress) GetEroSubobject ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetEroIpv4Subobject ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroIpv4Subobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetEroLabelSubobject ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroLabelSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetEroAsNumberSubobject ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroAsNumberSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetEroUnnumberedSubobject ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroUnnumberedSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) ListEroSubobjects ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress ListEroSubobjects
 return "", nil
}

func(np *RsvpTunnelIngress) GetNumberOfEroSubobjects ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress GetNumberOfEroSubobjects
 return "", nil
}

func(np *RsvpTunnelIngress) GetEroSubobjectType ()(string, error) {
 //parameters: TunnelHandle Index
 //AgtRsvpTunnelIngress GetEroSubobjectType
 return "", nil
}

func(np *RsvpTunnelIngress) AddRsvpEroSubobject () error {
 //parameters: TunnelHandle LspType Type
 //AgtRsvpTunnelIngress AddRsvpEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) InsertRsvpEroSubobject () error {
 //parameters: TunnelHandle LspType Index Type
 //AgtRsvpTunnelIngress InsertRsvpEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) RemoveRsvpEroSubobject () error {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress RemoveRsvpEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) RemoveAllRsvpEroSubobjects () error {
 //parameters: TunnelHandle LspType
 //AgtRsvpTunnelIngress RemoveAllRsvpEroSubobjects
 return nil
}

func(np *RsvpTunnelIngress) SetRsvpEroSubobject () error {
 //parameters: TunnelHandle LspType Index Type IsLoose Count psaContents
 //AgtRsvpTunnelIngress SetRsvpEroSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetRsvpEroIpv4Subobject () error {
 //parameters: TunnelHandle LspType Index IsLoose Address Prefix
 //AgtRsvpTunnelIngress SetRsvpEroIpv4Subobject
 return nil
}

func(np *RsvpTunnelIngress) SetRsvpEroLabelSubobject () error {
 //parameters: TunnelHandle LspType Index Upstream CType IsLoose Label
 //AgtRsvpTunnelIngress SetRsvpEroLabelSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetRsvpEroAsNumberSubobject () error {
 //parameters: TunnelHandle LspType Index IsLoose AsNumber
 //AgtRsvpTunnelIngress SetRsvpEroAsNumberSubobject
 return nil
}

func(np *RsvpTunnelIngress) SetRsvpEroUnnumberedSubobject () error {
 //parameters: TunnelHandle LspType Index IsLoose RouterId InterfaceId
 //AgtRsvpTunnelIngress SetRsvpEroUnnumberedSubobject
 return nil
}

func(np *RsvpTunnelIngress) GetRsvpEroSubobject ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetRsvpEroIpv4Subobject ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroIpv4Subobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetRsvpEroLabelSubobject ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroLabelSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetRsvpEroAsNumberSubobject ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroAsNumberSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) GetRsvpEroUnnumberedSubobject ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroUnnumberedSubobject
 return "", nil
}

func(np *RsvpTunnelIngress) ListRsvpEroSubobjects ()(string, error) {
 //parameters: TunnelHandle LspType
 //AgtRsvpTunnelIngress ListRsvpEroSubobjects
 return "", nil
}

func(np *RsvpTunnelIngress) GetNumberOfRsvpEroSubobjects ()(string, error) {
 //parameters: TunnelHandle LspType
 //AgtRsvpTunnelIngress GetNumberOfRsvpEroSubobjects
 return "", nil
}

func(np *RsvpTunnelIngress) GetRsvpEroSubobjectType ()(string, error) {
 //parameters: TunnelHandle LspType Index
 //AgtRsvpTunnelIngress GetRsvpEroSubobjectType
 return "", nil
}

func(np *RsvpTunnelIngress) SetTSpecParameter () error {
 //parameters: TunnelHandle Parameter Value
 //AgtRsvpTunnelIngress SetTSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) GetTSpecParameter ()(string, error) {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress GetTSpecParameter
 return "", nil
}

func(np *RsvpTunnelIngress) SetAdSpecParameter () error {
 //parameters: TunnelHandle Parameter Value
 //AgtRsvpTunnelIngress SetAdSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) GetAdSpecParameter ()(string, error) {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress GetAdSpecParameter
 return "", nil
}

func(np *RsvpTunnelIngress) SetProtectionParameter () error {
 //parameters: TunnelHandle Parameter Value
 //AgtRsvpTunnelIngress SetProtectionParameter
 return nil
}

func(np *RsvpTunnelIngress) GetProtectionParameter ()(string, error) {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress GetProtectionParameter
 return "", nil
}

func(np *RsvpTunnelIngress) SetGeneralizedUniParameter () error {
 //parameters: TunnelHandle Parameter Value
 //AgtRsvpTunnelIngress SetGeneralizedUniParameter
 return nil
}

func(np *RsvpTunnelIngress) GetGeneralizedUniParameter ()(string, error) {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress GetGeneralizedUniParameter
 return "", nil
}

func(np *RsvpTunnelIngress) SetGeneralizedUniTnaParameter () error {
 //parameters: TunnelHandle Parameter Address
 //AgtRsvpTunnelIngress SetGeneralizedUniTnaParameter
 return nil
}

func(np *RsvpTunnelIngress) GetGeneralizedUniTnaParameter ()(string, error) {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress GetGeneralizedUniTnaParameter
 return "", nil
}

func(np *RsvpTunnelIngress) EnableProtectionObject () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress EnableProtectionObject
 return nil
}

func(np *RsvpTunnelIngress) DisableProtectionObject () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress DisableProtectionObject
 return nil
}

func(np *RsvpTunnelIngress) IsProtectionObjectEnabled () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress IsProtectionObjectEnabled
 return nil
}

func(np *RsvpTunnelIngress) SelectTSpecParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress SelectTSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) DeselectTSpecParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress DeselectTSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) IsTSpecParameterSelected () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress IsTSpecParameterSelected
 return nil
}

func(np *RsvpTunnelIngress) SelectAdSpecParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress SelectAdSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) DeselectAdSpecParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress DeselectAdSpecParameter
 return nil
}

func(np *RsvpTunnelIngress) IsAdSpecParameterSelected () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress IsAdSpecParameterSelected
 return nil
}

func(np *RsvpTunnelIngress) SelectGeneralizedUniParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress SelectGeneralizedUniParameter
 return nil
}

func(np *RsvpTunnelIngress) DeselectGeneralizedUniParameter () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress DeselectGeneralizedUniParameter
 return nil
}

func(np *RsvpTunnelIngress) IsGeneralizedUniParameterSelected () error {
 //parameters: TunnelHandle Parameter
 //AgtRsvpTunnelIngress IsGeneralizedUniParameterSelected
 return nil
}

func(np *RsvpTunnelIngress) ListIngressLsps ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress ListIngressLsps
 return "", nil
}

func(np *RsvpTunnelIngress) EnableTraffic () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress EnableTraffic
 return nil
}

func(np *RsvpTunnelIngress) DisableTraffic () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress DisableTraffic
 return nil
}

func(np *RsvpTunnelIngress) IsTrafficEnabled () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress IsTrafficEnabled
 return nil
}

func(np *RsvpTunnelIngress) OpenInstance () error {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress OpenInstance
 return nil
}

func(np *RsvpTunnelIngress) CloseInstance () error {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress CloseInstance
 return nil
}

func(np *RsvpTunnelIngress) CloseInstanceGracefully () error {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress CloseInstanceGracefully
 return nil
}

func(np *RsvpTunnelIngress) CloseInstanceRemotely () error {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress CloseInstanceRemotely
 return nil
}

func(np *RsvpTunnelIngress) CloseTunnel () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress CloseTunnel
 return nil
}

func(np *RsvpTunnelIngress) CloseTunnelGracefully () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress CloseTunnelGracefully
 return nil
}

func(np *RsvpTunnelIngress) CloseTunnelRemotely () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelIngress CloseTunnelRemotely
 return nil
}

func(np *RsvpTunnelIngress) GetState ()(string, error) {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress GetState
 return "", nil
}

func(np *RsvpTunnelIngress) GetAssignedLabel ()(string, error) {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelIngress GetAssignedLabel
 return "", nil
}

