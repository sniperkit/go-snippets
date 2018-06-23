package n2xsdk

type LdpLspIngressPool struct {
 Handler string
}

func(np *LdpLspIngressPool) SetNumberOfLsps () error {
 //parameters: PoolHandle NumLsps Modifier
 //AgtLdpLspIngressPool SetNumberOfLsps
 return nil
}

func(np *LdpLspIngressPool) GetNumberOfLsps ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetNumberOfLsps
 return "", nil
}

func(np *LdpLspIngressPool) GetFecType ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetFecType
 return "", nil
}

func(np *LdpLspIngressPool) SetPrefixFec () error {
 //parameters: PoolHandle FirstAddress Incr PfxLength IncrMask
 //AgtLdpLspIngressPool SetPrefixFec
 return nil
}

func(np *LdpLspIngressPool) GetPrefixFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetPrefixFec
 return "", nil
}

func(np *LdpLspIngressPool) SetHostFec () error {
 //parameters: PoolHandle IpAddress Incr IncrMask
 //AgtLdpLspIngressPool SetHostFec
 return nil
}

func(np *LdpLspIngressPool) GetHostFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetHostFec
 return "", nil
}

func(np *LdpLspIngressPool) SetCrLdpFec () error {
 //parameters: PoolHandle ActionFlag LspId IpAddress PfxLength IncrLspId
 //AgtLdpLspIngressPool SetCrLdpFec
 return nil
}

func(np *LdpLspIngressPool) GetCrLdpFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetCrLdpFec
 return "", nil
}

func(np *LdpLspIngressPool) SetL2MplsFec () error {
 //parameters: PoolHandle CBit VcType GroupId VcId IncrVcId
 //AgtLdpLspIngressPool SetL2MplsFec
 return nil
}

func(np *LdpLspIngressPool) GetL2MplsFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetL2MplsFec
 return "", nil
}

func(np *LdpLspIngressPool) SetPathVectorTLV () error {
 //parameters: PoolHandle PathCount psaContentList
 //AgtLdpLspIngressPool SetPathVectorTLV
 return nil
}

func(np *LdpLspIngressPool) GetPathVectorTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetPathVectorTLV
 return "", nil
}

func(np *LdpLspIngressPool) SetHopCountTLV () error {
 //parameters: PoolHandle HopCount
 //AgtLdpLspIngressPool SetHopCountTLV
 return nil
}

func(np *LdpLspIngressPool) GetHopCountTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetHopCountTLV
 return "", nil
}

func(np *LdpLspIngressPool) EnableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool EnableDiffServXBit
 return nil
}

func(np *LdpLspIngressPool) DisableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool DisableDiffServXBit
 return nil
}

func(np *LdpLspIngressPool) IsDiffServXBitEnabled () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool IsDiffServXBitEnabled
 return nil
}

func(np *LdpLspIngressPool) SetDiffServPhbEncodingStyle () error {
 //parameters: PoolHandle EncodingStyle
 //AgtLdpLspIngressPool SetDiffServPhbEncodingStyle
 return nil
}

func(np *LdpLspIngressPool) GetDiffServPhbEncodingStyle ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetDiffServPhbEncodingStyle
 return "", nil
}

func(np *LdpLspIngressPool) EnableTLV () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspIngressPool EnableTLV
 return nil
}

func(np *LdpLspIngressPool) DisableTLV () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspIngressPool DisableTLV
 return nil
}

func(np *LdpLspIngressPool) IsTlvEnabled () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspIngressPool IsTlvEnabled
 return nil
}

func(np *LdpLspIngressPool) Open () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool Open
 return nil
}

func(np *LdpLspIngressPool) Close () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool Close
 return nil
}

func(np *LdpLspIngressPool) IsEnabled () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool IsEnabled
 return nil
}

func(np *LdpLspIngressPool) Cancel () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool Cancel
 return nil
}

func(np *LdpLspIngressPool) GetPoolState ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetPoolState
 return "", nil
}

func(np *LdpLspIngressPool) GetNumberOfOpenedLsps ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetNumberOfOpenedLsps
 return "", nil
}

func(np *LdpLspIngressPool) GetLspInfo ()(string, error) {
 //parameters: PoolHandle LspCount
 //AgtLdpLspIngressPool GetLspInfo
 return "", nil
}

func(np *LdpLspIngressPool) SetChunkCount () error {
 //parameters: PoolHandle ChunkCount
 //AgtLdpLspIngressPool SetChunkCount
 return nil
}

func(np *LdpLspIngressPool) GetChunkCount ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetChunkCount
 return "", nil
}

func(np *LdpLspIngressPool) SetChunkDelay () error {
 //parameters: PoolHandle ChunkDelay
 //AgtLdpLspIngressPool SetChunkDelay
 return nil
}

func(np *LdpLspIngressPool) GetChunkDelay ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetChunkDelay
 return "", nil
}

func(np *LdpLspIngressPool) SetFilterFecType () error {
 //parameters: PoolHandle FecType
 //AgtLdpLspIngressPool SetFilterFecType
 return nil
}

func(np *LdpLspIngressPool) SetFilterState () error {
 //parameters: PoolHandle State
 //AgtLdpLspIngressPool SetFilterState
 return nil
}

func(np *LdpLspIngressPool) SetFilterLabel () error {
 //parameters: PoolHandle Label
 //AgtLdpLspIngressPool SetFilterLabel
 return nil
}

func(np *LdpLspIngressPool) SetFilterDiffServType () error {
 //parameters: PoolHandle DiffServType
 //AgtLdpLspIngressPool SetFilterDiffServType
 return nil
}

func(np *LdpLspIngressPool) SetFilterLspId () error {
 //parameters: PoolHandle LspId
 //AgtLdpLspIngressPool SetFilterLspId
 return nil
}

func(np *LdpLspIngressPool) SetFilterSrcAddress () error {
 //parameters: PoolHandle SrcAddress
 //AgtLdpLspIngressPool SetFilterSrcAddress
 return nil
}

func(np *LdpLspIngressPool) SetFilterLen () error {
 //parameters: PoolHandle Len
 //AgtLdpLspIngressPool SetFilterLen
 return nil
}

func(np *LdpLspIngressPool) SetFilterRouterId () error {
 //parameters: PoolHandle RouterId
 //AgtLdpLspIngressPool SetFilterRouterId
 return nil
}

func(np *LdpLspIngressPool) SetFilterActFlag () error {
 //parameters: PoolHandle ActFlag
 //AgtLdpLspIngressPool SetFilterActFlag
 return nil
}

func(np *LdpLspIngressPool) SetFilterCbit () error {
 //parameters: PoolHandle CBit
 //AgtLdpLspIngressPool SetFilterCbit
 return nil
}

func(np *LdpLspIngressPool) SetFilterVcType () error {
 //parameters: PoolHandle VcType
 //AgtLdpLspIngressPool SetFilterVcType
 return nil
}

func(np *LdpLspIngressPool) SetFilterGroupId () error {
 //parameters: PoolHandle GroupId
 //AgtLdpLspIngressPool SetFilterGroupId
 return nil
}

func(np *LdpLspIngressPool) SetFilterVcId () error {
 //parameters: PoolHandle VcId
 //AgtLdpLspIngressPool SetFilterVcId
 return nil
}

func(np *LdpLspIngressPool) SetFilterL2MplsMtu () error {
 //parameters: PoolHandle MtuValue
 //AgtLdpLspIngressPool SetFilterL2MplsMtu
 return nil
}

func(np *LdpLspIngressPool) SetFilterL2MplsConcatAtmCell () error {
 //parameters: PoolHandle ConcatAtmCell
 //AgtLdpLspIngressPool SetFilterL2MplsConcatAtmCell
 return nil
}

func(np *LdpLspIngressPool) SetFilterL2MplsPayload () error {
 //parameters: PoolHandle CemPayload
 //AgtLdpLspIngressPool SetFilterL2MplsPayload
 return nil
}

func(np *LdpLspIngressPool) SetFilterL2MplsInterfaceString () error {
 //parameters: PoolHandle InterfaceString
 //AgtLdpLspIngressPool SetFilterL2MplsInterfaceString
 return nil
}

func(np *LdpLspIngressPool) SetFilterL2MplsOptions () error {
 //parameters: PoolHandle CemOptions
 //AgtLdpLspIngressPool SetFilterL2MplsOptions
 return nil
}

func(np *LdpLspIngressPool) SetL2MplsGenPwIdFec () error {
 //parameters: PoolHandle CBit VcType
 //AgtLdpLspIngressPool SetL2MplsGenPwIdFec
 return nil
}

func(np *LdpLspIngressPool) GetL2MplsGenPwIdFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetL2MplsGenPwIdFec
 return "", nil
}

func(np *LdpLspIngressPool) SetTrafficTLV () error {
 //parameters: PoolHandle Frequency Weight PeakDataRate PeakBurstSize CommittedDataRate CommittedBurstSize ExcessBurstSize NegotiableFlag
 //AgtLdpLspIngressPool SetTrafficTLV
 return nil
}

func(np *LdpLspIngressPool) GetTrafficTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetTrafficTLV
 return "", nil
}

func(np *LdpLspIngressPool) AddErSubobject () error {
 //parameters: PoolHandle Type
 //AgtLdpLspIngressPool AddErSubobject
 return nil
}

func(np *LdpLspIngressPool) RemoveErSubobject () error {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool RemoveErSubobject
 return nil
}

func(np *LdpLspIngressPool) RemoveAllErSubobjects () error {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool RemoveAllErSubobjects
 return nil
}

func(np *LdpLspIngressPool) SetErTLV () error {
 //parameters: PoolHandle Index Type IsLoose Count psaContents
 //AgtLdpLspIngressPool SetErTLV
 return nil
}

func(np *LdpLspIngressPool) SetErIpv4Subobject () error {
 //parameters: PoolHandle Index IsLoose Address Prefix
 //AgtLdpLspIngressPool SetErIpv4Subobject
 return nil
}

func(np *LdpLspIngressPool) SetErAsNumberSubobject () error {
 //parameters: PoolHandle Index IsLoose AsNumber
 //AgtLdpLspIngressPool SetErAsNumberSubobject
 return nil
}

func(np *LdpLspIngressPool) SetErLspidSubobject () error {
 //parameters: PoolHandle Index RouterId LspId
 //AgtLdpLspIngressPool SetErLspidSubobject
 return nil
}

func(np *LdpLspIngressPool) GetErTLV ()(string, error) {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool GetErTLV
 return "", nil
}

func(np *LdpLspIngressPool) GetErIpv4Subobject ()(string, error) {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool GetErIpv4Subobject
 return "", nil
}

func(np *LdpLspIngressPool) GetErAsNumberSubobject ()(string, error) {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool GetErAsNumberSubobject
 return "", nil
}

func(np *LdpLspIngressPool) GetErLspidSubobject ()(string, error) {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool GetErLspidSubobject
 return "", nil
}

func(np *LdpLspIngressPool) ListErSubobjects ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool ListErSubobjects
 return "", nil
}

func(np *LdpLspIngressPool) GetNumberOfErSubobjects ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetNumberOfErSubobjects
 return "", nil
}

func(np *LdpLspIngressPool) GetErSubobjectType ()(string, error) {
 //parameters: PoolHandle Index
 //AgtLdpLspIngressPool GetErSubobjectType
 return "", nil
}

func(np *LdpLspIngressPool) SetPreemptTLV () error {
 //parameters: PoolHandle SetPriority HoldPriority
 //AgtLdpLspIngressPool SetPreemptTLV
 return nil
}

func(np *LdpLspIngressPool) GetPreemptTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetPreemptTLV
 return "", nil
}

func(np *LdpLspIngressPool) GetPinningTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetPinningTLV
 return "", nil
}

func(np *LdpLspIngressPool) SetPinningTLV () error {
 //parameters: PoolHandle PBit
 //AgtLdpLspIngressPool SetPinningTLV
 return nil
}

func(np *LdpLspIngressPool) GetResourceBitMaskTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetResourceBitMaskTLV
 return "", nil
}

func(np *LdpLspIngressPool) SetResourceBitMaskTLV () error {
 //parameters: PoolHandle ResourceBitMask
 //AgtLdpLspIngressPool SetResourceBitMaskTLV
 return nil
}

func(np *LdpLspIngressPool) GetIngressDiffServEtypeTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetIngressDiffServEtypeTLV
 return "", nil
}

func(np *LdpLspIngressPool) SetIngressDiffServEtypeTLV () error {
 //parameters: PoolHandle ExpCount psaExp PhbldCount psaPhbId
 //AgtLdpLspIngressPool SetIngressDiffServEtypeTLV
 return nil
}

func(np *LdpLspIngressPool) SetIngressDiffServLtypeTLV () error {
 //parameters: PoolHandle Psc
 //AgtLdpLspIngressPool SetIngressDiffServLtypeTLV
 return nil
}

func(np *LdpLspIngressPool) GetIngressDiffServLtypeTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspIngressPool GetIngressDiffServLtypeTLV
 return "", nil
}

