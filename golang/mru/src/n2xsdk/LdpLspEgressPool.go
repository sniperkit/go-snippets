package n2xsdk

type LdpLspEgressPool struct {
 Handler string
}

func(np *LdpLspEgressPool) SetNumberOfLsps () error {
 //parameters: PoolHandle NumLsps Modifier
 //AgtLdpLspEgressPool SetNumberOfLsps
 return nil
}

func(np *LdpLspEgressPool) GetNumberOfLsps ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetNumberOfLsps
 return "", nil
}

func(np *LdpLspEgressPool) GetFecType ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetFecType
 return "", nil
}

func(np *LdpLspEgressPool) SetPrefixFec () error {
 //parameters: PoolHandle FirstAddress Incr PfxLength IncrMask
 //AgtLdpLspEgressPool SetPrefixFec
 return nil
}

func(np *LdpLspEgressPool) GetPrefixFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetPrefixFec
 return "", nil
}

func(np *LdpLspEgressPool) SetHostFec () error {
 //parameters: PoolHandle IpAddress Incr IncrMask
 //AgtLdpLspEgressPool SetHostFec
 return nil
}

func(np *LdpLspEgressPool) GetHostFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetHostFec
 return "", nil
}

func(np *LdpLspEgressPool) SetCrLdpFec () error {
 //parameters: PoolHandle ActionFlag LspId IpAddress PfxLength IncrLspId
 //AgtLdpLspEgressPool SetCrLdpFec
 return nil
}

func(np *LdpLspEgressPool) GetCrLdpFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetCrLdpFec
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsFec () error {
 //parameters: PoolHandle CBit VcType GroupId VcId IncrVcId
 //AgtLdpLspEgressPool SetL2MplsFec
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsFec
 return "", nil
}

func(np *LdpLspEgressPool) SetPathVectorTLV () error {
 //parameters: PoolHandle PathCount psaContentList
 //AgtLdpLspEgressPool SetPathVectorTLV
 return nil
}

func(np *LdpLspEgressPool) GetPathVectorTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetPathVectorTLV
 return "", nil
}

func(np *LdpLspEgressPool) SetHopCountTLV () error {
 //parameters: PoolHandle HopCount
 //AgtLdpLspEgressPool SetHopCountTLV
 return nil
}

func(np *LdpLspEgressPool) GetHopCountTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetHopCountTLV
 return "", nil
}

func(np *LdpLspEgressPool) EnableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool EnableDiffServXBit
 return nil
}

func(np *LdpLspEgressPool) DisableDiffServXBit () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool DisableDiffServXBit
 return nil
}

func(np *LdpLspEgressPool) IsDiffServXBitEnabled () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool IsDiffServXBitEnabled
 return nil
}

func(np *LdpLspEgressPool) SetDiffServPhbEncodingStyle () error {
 //parameters: PoolHandle EncodingStyle
 //AgtLdpLspEgressPool SetDiffServPhbEncodingStyle
 return nil
}

func(np *LdpLspEgressPool) GetDiffServPhbEncodingStyle ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetDiffServPhbEncodingStyle
 return "", nil
}

func(np *LdpLspEgressPool) EnableTLV () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspEgressPool EnableTLV
 return nil
}

func(np *LdpLspEgressPool) DisableTLV () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspEgressPool DisableTLV
 return nil
}

func(np *LdpLspEgressPool) IsTlvEnabled () error {
 //parameters: PoolHandle TlvType
 //AgtLdpLspEgressPool IsTlvEnabled
 return nil
}

func(np *LdpLspEgressPool) Open () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool Open
 return nil
}

func(np *LdpLspEgressPool) Close () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool Close
 return nil
}

func(np *LdpLspEgressPool) IsEnabled () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool IsEnabled
 return nil
}

func(np *LdpLspEgressPool) Cancel () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool Cancel
 return nil
}

func(np *LdpLspEgressPool) GetPoolState ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetPoolState
 return "", nil
}

func(np *LdpLspEgressPool) GetNumberOfOpenedLsps ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetNumberOfOpenedLsps
 return "", nil
}

func(np *LdpLspEgressPool) GetLspInfo ()(string, error) {
 //parameters: PoolHandle LspCount
 //AgtLdpLspEgressPool GetLspInfo
 return "", nil
}

func(np *LdpLspEgressPool) SetChunkCount () error {
 //parameters: PoolHandle ChunkCount
 //AgtLdpLspEgressPool SetChunkCount
 return nil
}

func(np *LdpLspEgressPool) GetChunkCount ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetChunkCount
 return "", nil
}

func(np *LdpLspEgressPool) SetChunkDelay () error {
 //parameters: PoolHandle ChunkDelay
 //AgtLdpLspEgressPool SetChunkDelay
 return nil
}

func(np *LdpLspEgressPool) GetChunkDelay ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetChunkDelay
 return "", nil
}

func(np *LdpLspEgressPool) SetFilterFecType () error {
 //parameters: PoolHandle FecType
 //AgtLdpLspEgressPool SetFilterFecType
 return nil
}

func(np *LdpLspEgressPool) SetFilterState () error {
 //parameters: PoolHandle State
 //AgtLdpLspEgressPool SetFilterState
 return nil
}

func(np *LdpLspEgressPool) SetFilterLabel () error {
 //parameters: PoolHandle Label
 //AgtLdpLspEgressPool SetFilterLabel
 return nil
}

func(np *LdpLspEgressPool) SetFilterDiffServType () error {
 //parameters: PoolHandle DiffServType
 //AgtLdpLspEgressPool SetFilterDiffServType
 return nil
}

func(np *LdpLspEgressPool) SetFilterLspId () error {
 //parameters: PoolHandle LspId
 //AgtLdpLspEgressPool SetFilterLspId
 return nil
}

func(np *LdpLspEgressPool) SetFilterSrcAddress () error {
 //parameters: PoolHandle SrcAddress
 //AgtLdpLspEgressPool SetFilterSrcAddress
 return nil
}

func(np *LdpLspEgressPool) SetFilterLen () error {
 //parameters: PoolHandle Len
 //AgtLdpLspEgressPool SetFilterLen
 return nil
}

func(np *LdpLspEgressPool) SetFilterRouterId () error {
 //parameters: PoolHandle RouterId
 //AgtLdpLspEgressPool SetFilterRouterId
 return nil
}

func(np *LdpLspEgressPool) SetFilterActFlag () error {
 //parameters: PoolHandle ActFlag
 //AgtLdpLspEgressPool SetFilterActFlag
 return nil
}

func(np *LdpLspEgressPool) SetFilterCbit () error {
 //parameters: PoolHandle CBit
 //AgtLdpLspEgressPool SetFilterCbit
 return nil
}

func(np *LdpLspEgressPool) SetFilterVcType () error {
 //parameters: PoolHandle VcType
 //AgtLdpLspEgressPool SetFilterVcType
 return nil
}

func(np *LdpLspEgressPool) SetFilterGroupId () error {
 //parameters: PoolHandle GroupId
 //AgtLdpLspEgressPool SetFilterGroupId
 return nil
}

func(np *LdpLspEgressPool) SetFilterVcId () error {
 //parameters: PoolHandle VcId
 //AgtLdpLspEgressPool SetFilterVcId
 return nil
}

func(np *LdpLspEgressPool) SetFilterL2MplsMtu () error {
 //parameters: PoolHandle MtuValue
 //AgtLdpLspEgressPool SetFilterL2MplsMtu
 return nil
}

func(np *LdpLspEgressPool) SetFilterL2MplsConcatAtmCell () error {
 //parameters: PoolHandle ConcatAtmCell
 //AgtLdpLspEgressPool SetFilterL2MplsConcatAtmCell
 return nil
}

func(np *LdpLspEgressPool) SetFilterL2MplsPayload () error {
 //parameters: PoolHandle CemPayload
 //AgtLdpLspEgressPool SetFilterL2MplsPayload
 return nil
}

func(np *LdpLspEgressPool) SetFilterL2MplsInterfaceString () error {
 //parameters: PoolHandle InterfaceString
 //AgtLdpLspEgressPool SetFilterL2MplsInterfaceString
 return nil
}

func(np *LdpLspEgressPool) SetFilterL2MplsOptions () error {
 //parameters: PoolHandle CemOptions
 //AgtLdpLspEgressPool SetFilterL2MplsOptions
 return nil
}

func(np *LdpLspEgressPool) SetL2MplsGenPwIdFec () error {
 //parameters: PoolHandle CBit VcType
 //AgtLdpLspEgressPool SetL2MplsGenPwIdFec
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsGenPwIdFec ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsGenPwIdFec
 return "", nil
}

func(np *LdpLspEgressPool) SetNullLabelType () error {
 //parameters: PoolHandle NullLabelType
 //AgtLdpLspEgressPool SetNullLabelType
 return nil
}

func(np *LdpLspEgressPool) GetNullLabelType ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetNullLabelType
 return "", nil
}

func(np *LdpLspEgressPool) SetEgressDiffServEtypeTLV () error {
 //parameters: PoolHandle ExpCount psaExp PhbldCount psaPhbId
 //AgtLdpLspEgressPool SetEgressDiffServEtypeTLV
 return nil
}

func(np *LdpLspEgressPool) GetEgressDiffServEtypeTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetEgressDiffServEtypeTLV
 return "", nil
}

func(np *LdpLspEgressPool) SetEgressDiffServLtypeTLV () error {
 //parameters: PoolHandle Psc
 //AgtLdpLspEgressPool SetEgressDiffServLtypeTLV
 return nil
}

func(np *LdpLspEgressPool) GetEgressDiffServLtypeTLV ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetEgressDiffServLtypeTLV
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsMtu () error {
 //parameters: PoolHandle MtuValue
 //AgtLdpLspEgressPool SetL2MplsMtu
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsMtu ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsMtu
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsConcatAtmCell () error {
 //parameters: PoolHandle ConcatAtmCell
 //AgtLdpLspEgressPool SetL2MplsConcatAtmCell
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsConcatAtmCell ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsConcatAtmCell
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsPayload () error {
 //parameters: PoolHandle CemPayload
 //AgtLdpLspEgressPool SetL2MplsPayload
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsPayload ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsPayload
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsInterfaceString () error {
 //parameters: PoolHandle InterfaceString
 //AgtLdpLspEgressPool SetL2MplsInterfaceString
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsInterfaceString ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsInterfaceString
 return "", nil
}

func(np *LdpLspEgressPool) SetL2MplsOptions () error {
 //parameters: PoolHandle CemOptions
 //AgtLdpLspEgressPool SetL2MplsOptions
 return nil
}

func(np *LdpLspEgressPool) GetL2MplsOptions ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetL2MplsOptions
 return "", nil
}

func(np *LdpLspEgressPool) EnableL2MplsInterfaceParameter () error {
 //parameters: PoolHandle InterfaceType
 //AgtLdpLspEgressPool EnableL2MplsInterfaceParameter
 return nil
}

func(np *LdpLspEgressPool) DisableL2MplsInterfaceParameter () error {
 //parameters: PoolHandle InterfaceType
 //AgtLdpLspEgressPool DisableL2MplsInterfaceParameter
 return nil
}

func(np *LdpLspEgressPool) IsL2MplsInterfaceParameterEnabled () error {
 //parameters: PoolHandle InterfaceType
 //AgtLdpLspEgressPool IsL2MplsInterfaceParameterEnabled
 return nil
}

func(np *LdpLspEgressPool) SetRequestedVlanId () error {
 //parameters: PoolHandle RequestedVlanId
 //AgtLdpLspEgressPool SetRequestedVlanId
 return nil
}

func(np *LdpLspEgressPool) GetRequestedVlanId ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetRequestedVlanId
 return "", nil
}

func(np *LdpLspEgressPool) SetCemBitRate () error {
 //parameters: PoolHandle CemBitRate
 //AgtLdpLspEgressPool SetCemBitRate
 return nil
}

func(np *LdpLspEgressPool) GetCemBitRate ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetCemBitRate
 return "", nil
}

func(np *LdpLspEgressPool) SetFrDlciLength () error {
 //parameters: PoolHandle FrDlciLength
 //AgtLdpLspEgressPool SetFrDlciLength
 return nil
}

func(np *LdpLspEgressPool) GetFrDlciLength ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetFrDlciLength
 return "", nil
}

func(np *LdpLspEgressPool) SetFcsRetentionIndicator () error {
 //parameters: PoolHandle FcsRetentionIndicator
 //AgtLdpLspEgressPool SetFcsRetentionIndicator
 return nil
}

func(np *LdpLspEgressPool) GetFcsRetentionIndicator ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetFcsRetentionIndicator
 return "", nil
}

func(np *LdpLspEgressPool) SetFragmentationIndicator () error {
 //parameters: PoolHandle FragmentationIndicator
 //AgtLdpLspEgressPool SetFragmentationIndicator
 return nil
}

func(np *LdpLspEgressPool) GetFragmentationIndicator ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetFragmentationIndicator
 return "", nil
}

func(np *LdpLspEgressPool) SetTdmOption () error {
 //parameters: PoolHandle TdmOption
 //AgtLdpLspEgressPool SetTdmOption
 return nil
}

func(np *LdpLspEgressPool) GetTdmOption ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetTdmOption
 return "", nil
}

func(np *LdpLspEgressPool) SetVccvParameter () error {
 //parameters: PoolHandle VccvParameter
 //AgtLdpLspEgressPool SetVccvParameter
 return nil
}

func(np *LdpLspEgressPool) GetVccvParameter ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetVccvParameter
 return "", nil
}

func(np *LdpLspEgressPool) EnableL2MplsPwStatus () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool EnableL2MplsPwStatus
 return nil
}

func(np *LdpLspEgressPool) DisableL2MplsPwStatus () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool DisableL2MplsPwStatus
 return nil
}

func(np *LdpLspEgressPool) IsL2MplsPwStatusEnabled () error {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool IsL2MplsPwStatusEnabled
 return nil
}

func(np *LdpLspEgressPool) SetPwStatus () error {
 //parameters: PoolHandle PWStatusCode
 //AgtLdpLspEgressPool SetPwStatus
 return nil
}

func(np *LdpLspEgressPool) GetPwStatus ()(string, error) {
 //parameters: PoolHandle
 //AgtLdpLspEgressPool GetPwStatus
 return "", nil
}

