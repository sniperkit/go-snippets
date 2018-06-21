package n2xsdk

type LdpSession struct {
 Handler string
}

func(np *LdpSession) SetSutIpAddress () error {
 //parameters: SessionHandle RemoteIpAddress
 //AgtLdpSession SetSutIpAddress
 return nil
}

func(np *LdpSession) GetSutIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetSutIpAddress
 return "", nil
}

func(np *LdpSession) SetTesterIpAddress () error {
 //parameters: SessionHandle LocalIpAddress
 //AgtLdpSession SetTesterIpAddress
 return nil
}

func(np *LdpSession) GetTesterIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetTesterIpAddress
 return "", nil
}

func(np *LdpSession) SetSubnetMask () error {
 //parameters: SessionHandle Mask
 //AgtLdpSession SetSubnetMask
 return nil
}

func(np *LdpSession) GetSubnetMask ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetSubnetMask
 return "", nil
}

func(np *LdpSession) SetPeerType () error {
 //parameters: SessionHandle PeerType
 //AgtLdpSession SetPeerType
 return nil
}

func(np *LdpSession) GetPeerType ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetPeerType
 return "", nil
}

func(np *LdpSession) SetHoldTime () error {
 //parameters: SessionHandle HoldTimeSec
 //AgtLdpSession SetHoldTime
 return nil
}

func(np *LdpSession) GetHoldTime ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetHoldTime
 return "", nil
}

func(np *LdpSession) SetLabelSpaceId () error {
 //parameters: SessionHandle LabelSpaceId
 //AgtLdpSession SetLabelSpaceId
 return nil
}

func(np *LdpSession) GetLabelSpaceId ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetLabelSpaceId
 return "", nil
}

func(np *LdpSession) SetKeepAliveTime () error {
 //parameters: SessionHandle KeepAliveTimeSec
 //AgtLdpSession SetKeepAliveTime
 return nil
}

func(np *LdpSession) GetKeepAliveTime ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetKeepAliveTime
 return "", nil
}

func(np *LdpSession) SetRouterId () error {
 //parameters: SessionHandle RouterId
 //AgtLdpSession SetRouterId
 return nil
}

func(np *LdpSession) GetRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetRouterId
 return "", nil
}

func(np *LdpSession) EnableAutoTesterRouterId () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableAutoTesterRouterId
 return nil
}

func(np *LdpSession) DisableAutoTesterRouterId () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableAutoTesterRouterId
 return nil
}

func(np *LdpSession) IsAutoTesterRouterIdEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsAutoTesterRouterIdEnabled
 return nil
}

func(np *LdpSession) SetTesterRouterId () error {
 //parameters: SessionHandle RouterId
 //AgtLdpSession SetTesterRouterId
 return nil
}

func(np *LdpSession) GetTesterRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetTesterRouterId
 return "", nil
}

func(np *LdpSession) SetSutRouterId () error {
 //parameters: SessionHandle SutRouterId
 //AgtLdpSession SetSutRouterId
 return nil
}

func(np *LdpSession) GetSutRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetSutRouterId
 return "", nil
}

func(np *LdpSession) EnableTransportAddress () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableTransportAddress
 return nil
}

func(np *LdpSession) DisableTransportAddress () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableTransportAddress
 return nil
}

func(np *LdpSession) IsTransportAddressEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsTransportAddressEnabled
 return nil
}

func(np *LdpSession) SetTransportAddress () error {
 //parameters: SessionHandle TransportAddress
 //AgtLdpSession SetTransportAddress
 return nil
}

func(np *LdpSession) GetTransportAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetTransportAddress
 return "", nil
}

func(np *LdpSession) SetVcDirectionality () error {
 //parameters: SessionHandle VcDir
 //AgtLdpSession SetVcDirectionality
 return nil
}

func(np *LdpSession) GetVcDirectionality ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetVcDirectionality
 return "", nil
}

func(np *LdpSession) SetMergeCapability () error {
 //parameters: SessionHandle MergeCap
 //AgtLdpSession SetMergeCapability
 return nil
}

func(np *LdpSession) GetMergeCapability ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetMergeCapability
 return "", nil
}

func(np *LdpSession) SetLabelAdvertisementMode () error {
 //parameters: SessionHandle LAMode
 //AgtLdpSession SetLabelAdvertisementMode
 return nil
}

func(np *LdpSession) GetLabelAdvertisementMode ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetLabelAdvertisementMode
 return "", nil
}

func(np *LdpSession) SetMaxPduLength () error {
 //parameters: SessionHandle MaxPduLengthBytes
 //AgtLdpSession SetMaxPduLength
 return nil
}

func(np *LdpSession) GetMaxPduLength ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetMaxPduLength
 return "", nil
}

func(np *LdpSession) Enable () error {
 //parameters: SessionHandle
 //AgtLdpSession Enable
 return nil
}

func(np *LdpSession) Disable () error {
 //parameters: SessionHandle
 //AgtLdpSession Disable
 return nil
}

func(np *LdpSession) GetEnableFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetEnableFlag
 return "", nil
}

func(np *LdpSession) SetSessionAggregation () error {
 //parameters: SessionHandle
 //AgtLdpSession SetSessionAggregation
 return nil
}

func(np *LdpSession) UnSetSessionAggregation () error {
 //parameters: SessionHandle
 //AgtLdpSession UnSetSessionAggregation
 return nil
}

func(np *LdpSession) IsSessionAggregationSet () error {
 //parameters: SessionHandle
 //AgtLdpSession IsSessionAggregationSet
 return nil
}

func(np *LdpSession) EnableLoopDetection () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableLoopDetection
 return nil
}

func(np *LdpSession) DisableLoopDetection () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableLoopDetection
 return nil
}

func(np *LdpSession) IsLoopDetectionEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsLoopDetectionEnabled
 return nil
}

func(np *LdpSession) EnableLspDiscardMode () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableLspDiscardMode
 return nil
}

func(np *LdpSession) DisableLspDiscardMode () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableLspDiscardMode
 return nil
}

func(np *LdpSession) IsLspDiscardModeEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsLspDiscardModeEnabled
 return nil
}

func(np *LdpSession) EnableNotificationOnUnsupportedFec () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableNotificationOnUnsupportedFec
 return nil
}

func(np *LdpSession) DisableNotificationOnUnsupportedFec () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableNotificationOnUnsupportedFec
 return nil
}

func(np *LdpSession) IsNotificationOnUnsupportedFecEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsNotificationOnUnsupportedFecEnabled
 return nil
}

func(np *LdpSession) EnableNotificationOnUnsupportedTlv () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableNotificationOnUnsupportedTlv
 return nil
}

func(np *LdpSession) DisableNotificationOnUnsupportedTlv () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableNotificationOnUnsupportedTlv
 return nil
}

func(np *LdpSession) IsNotificationOnUnsupportedTlvEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsNotificationOnUnsupportedTlvEnabled
 return nil
}

func(np *LdpSession) EnableFtSessionTlvFlag () error {
 //parameters: SessionHandle Flag
 //AgtLdpSession EnableFtSessionTlvFlag
 return nil
}

func(np *LdpSession) DisableFtSessionTlvFlag () error {
 //parameters: SessionHandle Flag
 //AgtLdpSession DisableFtSessionTlvFlag
 return nil
}

func(np *LdpSession) IsFtSessionTlvFlagEnabled () error {
 //parameters: SessionHandle Flag
 //AgtLdpSession IsFtSessionTlvFlagEnabled
 return nil
}

func(np *LdpSession) EnableMd5Authentication () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableMd5Authentication
 return nil
}

func(np *LdpSession) DisableMd5Authentication () error {
 //parameters: SessionHandle
 //AgtLdpSession DisableMd5Authentication
 return nil
}

func(np *LdpSession) IsMd5AuthenticationEnabled () error {
 //parameters: SessionHandle
 //AgtLdpSession IsMd5AuthenticationEnabled
 return nil
}

func(np *LdpSession) SetMd5AuthenticationKey () error {
 //parameters: SessionHandle Md5Key
 //AgtLdpSession SetMd5AuthenticationKey
 return nil
}

func(np *LdpSession) GetMd5AuthenticationKey ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetMd5AuthenticationKey
 return "", nil
}

func(np *LdpSession) EnableTLV () error {
 //parameters: SessionHandle TlvType
 //AgtLdpSession EnableTLV
 return nil
}

func(np *LdpSession) DisableTLV () error {
 //parameters: SessionHandle TlvType
 //AgtLdpSession DisableTLV
 return nil
}

func(np *LdpSession) IsTlvEnabled () error {
 //parameters: SessionHandle TlvType
 //AgtLdpSession IsTlvEnabled
 return nil
}

func(np *LdpSession) SetFtSessionTlvTimer () error {
 //parameters: SessionHandle Timer TimerValue
 //AgtLdpSession SetFtSessionTlvTimer
 return nil
}

func(np *LdpSession) GetFtSessionTlvTimer ()(string, error) {
 //parameters: SessionHandle Timer
 //AgtLdpSession GetFtSessionTlvTimer
 return "", nil
}

func(np *LdpSession) EnableConformance () error {
 //parameters: SessionHandle
 //AgtLdpSession EnableConformance
 return nil
}

func(np *LdpSession) GetLastError ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetLastError
 return "", nil
}

func(np *LdpSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetState
 return "", nil
}

func(np *LdpSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtLdpSession GetStateSummary
 return "", nil
}

func(np *LdpSession) GetLspStateSummaryBySession ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpSession GetLspStateSummaryBySession
 return "", nil
}

func(np *LdpSession) GetLspStateSummary ()(string, error) {
 //parameters: 
 //AgtLdpSession GetLspStateSummary
 return "", nil
}

func(np *LdpSession) SimulateRestart () error {
 //parameters: SessionHandle CloseType DownTime
 //AgtLdpSession SimulateRestart
 return nil
}

func(np *LdpSession) EnableAllSessions () error {
 //parameters: 
 //AgtLdpSession EnableAllSessions
 return nil
}

func(np *LdpSession) DisableAllSessions () error {
 //parameters: 
 //AgtLdpSession DisableAllSessions
 return nil
}

