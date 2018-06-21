package n2xsdk

type RsvpSession struct {
 Handler string
}

func(np *RsvpSession) SetSutIpAddress () error {
 //parameters: SessionHandle RemoteIpAddress
 //AgtRsvpSession SetSutIpAddress
 return nil
}

func(np *RsvpSession) GetSutIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetSutIpAddress
 return "", nil
}

func(np *RsvpSession) SetTesterIpAddress () error {
 //parameters: SessionHandle LocalIpAddress
 //AgtRsvpSession SetTesterIpAddress
 return nil
}

func(np *RsvpSession) GetTesterIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetTesterIpAddress
 return "", nil
}

func(np *RsvpSession) SetRefreshPeriod () error {
 //parameters: SessionHandle PeriodMs
 //AgtRsvpSession SetRefreshPeriod
 return nil
}

func(np *RsvpSession) GetRefreshPeriod ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRefreshPeriod
 return "", nil
}

func(np *RsvpSession) SetLTimerK () error {
 //parameters: SessionHandle LTimerK
 //AgtRsvpSession SetLTimerK
 return nil
}

func(np *RsvpSession) GetLTimerK ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetLTimerK
 return "", nil
}

func(np *RsvpSession) EnableHelloSignalling () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableHelloSignalling
 return nil
}

func(np *RsvpSession) DisableHelloSignalling () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableHelloSignalling
 return nil
}

func(np *RsvpSession) IsHelloSignallingEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsHelloSignallingEnabled
 return nil
}

func(np *RsvpSession) SetHelloTTL () error {
 //parameters: SessionHandle HelloTTL
 //AgtRsvpSession SetHelloTTL
 return nil
}

func(np *RsvpSession) GetHelloTTL ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetHelloTTL
 return "", nil
}

func(np *RsvpSession) SetHelloInterval () error {
 //parameters: SessionHandle HelloInterval
 //AgtRsvpSession SetHelloInterval
 return nil
}

func(np *RsvpSession) GetHelloInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetHelloInterval
 return "", nil
}

func(np *RsvpSession) SetHelloTimerK () error {
 //parameters: SessionHandle HelloTimerK
 //AgtRsvpSession SetHelloTimerK
 return nil
}

func(np *RsvpSession) GetHelloTimerK ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetHelloTimerK
 return "", nil
}

func(np *RsvpSession) SetHelloIncludeRestartCapOption () error {
 //parameters: SessionHandle RestartCapOption
 //AgtRsvpSession SetHelloIncludeRestartCapOption
 return nil
}

func(np *RsvpSession) GetHelloIncludeRestartCapOption ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetHelloIncludeRestartCapOption
 return "", nil
}

func(np *RsvpSession) EnableRefreshReduction () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableRefreshReduction
 return nil
}

func(np *RsvpSession) DisableRefreshReduction () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableRefreshReduction
 return nil
}

func(np *RsvpSession) IsRefreshReductionEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsRefreshReductionEnabled
 return nil
}

func(np *RsvpSession) SetRefreshReductionInitTimeout () error {
 //parameters: SessionHandle RefreshReductionInitTimeout
 //AgtRsvpSession SetRefreshReductionInitTimeout
 return nil
}

func(np *RsvpSession) GetRefreshReductionInitTimeout ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRefreshReductionInitTimeout
 return "", nil
}

func(np *RsvpSession) SetRefreshReductionExpBackoff () error {
 //parameters: SessionHandle RefreshReductionExpBackoff
 //AgtRsvpSession SetRefreshReductionExpBackoff
 return nil
}

func(np *RsvpSession) GetRefreshReductionExpBackoff ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRefreshReductionExpBackoff
 return "", nil
}

func(np *RsvpSession) SetRefreshReductionRetryLimit () error {
 //parameters: SessionHandle RefreshReductionRetryLimit
 //AgtRsvpSession SetRefreshReductionRetryLimit
 return nil
}

func(np *RsvpSession) GetRefreshReductionRetryLimit ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRefreshReductionRetryLimit
 return "", nil
}

func(np *RsvpSession) SetRefreshReductionMaxAckDelay () error {
 //parameters: SessionHandle RefreshReductionMaxAckDelay
 //AgtRsvpSession SetRefreshReductionMaxAckDelay
 return nil
}

func(np *RsvpSession) GetRefreshReductionMaxAckDelay ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRefreshReductionMaxAckDelay
 return "", nil
}

func(np *RsvpSession) EnableSrefresh () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableSrefresh
 return nil
}

func(np *RsvpSession) DisableSrefresh () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableSrefresh
 return nil
}

func(np *RsvpSession) IsSrefreshEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsSrefreshEnabled
 return nil
}

func(np *RsvpSession) EnableBundle () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableBundle
 return nil
}

func(np *RsvpSession) DisableBundle () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableBundle
 return nil
}

func(np *RsvpSession) IsBundleEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsBundleEnabled
 return nil
}

func(np *RsvpSession) EnableEgressRroInsertion () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableEgressRroInsertion
 return nil
}

func(np *RsvpSession) DisableEgressRroInsertion () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableEgressRroInsertion
 return nil
}

func(np *RsvpSession) IsEgressRroInsertionEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsEgressRroInsertionEnabled
 return nil
}

func(np *RsvpSession) EnableEgressRroLabelInsertion () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableEgressRroLabelInsertion
 return nil
}

func(np *RsvpSession) DisableEgressRroLabelInsertion () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableEgressRroLabelInsertion
 return nil
}

func(np *RsvpSession) IsEgressRroLabelInsertionEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsEgressRroLabelInsertionEnabled
 return nil
}

func(np *RsvpSession) SetEgressRroLabelStart () error {
 //parameters: SessionHandle EgressRroLabelStart
 //AgtRsvpSession SetEgressRroLabelStart
 return nil
}

func(np *RsvpSession) GetEgressRroLabelStart ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetEgressRroLabelStart
 return "", nil
}

func(np *RsvpSession) SetEgressRroLabelIncrement () error {
 //parameters: SessionHandle EgressRroLabelIncrement
 //AgtRsvpSession SetEgressRroLabelIncrement
 return nil
}

func(np *RsvpSession) GetEgressRroLabelIncrement ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetEgressRroLabelIncrement
 return "", nil
}

func(np *RsvpSession) EnableEgressRroNodeIdFlag () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableEgressRroNodeIdFlag
 return nil
}

func(np *RsvpSession) DisableEgressRroNodeIdFlag () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableEgressRroNodeIdFlag
 return nil
}

func(np *RsvpSession) IsEgressRroNodeIdFlagEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsEgressRroNodeIdFlagEnabled
 return nil
}

func(np *RsvpSession) SetEgressRroNodeIdPrefix () error {
 //parameters: SessionHandle EgressRroNodeIdPrefixAddress EgressRroNodeIdPrefixLength
 //AgtRsvpSession SetEgressRroNodeIdPrefix
 return nil
}

func(np *RsvpSession) GetEgressRroNodeIdPrefix ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetEgressRroNodeIdPrefix
 return "", nil
}

func(np *RsvpSession) EnableGreTunnel () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableGreTunnel
 return nil
}

func(np *RsvpSession) DisableGreTunnel () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableGreTunnel
 return nil
}

func(np *RsvpSession) IsGreTunnelEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsGreTunnelEnabled
 return nil
}

func(np *RsvpSession) EnableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableGreTunnelChecksumField
 return nil
}

func(np *RsvpSession) DisableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableGreTunnelChecksumField
 return nil
}

func(np *RsvpSession) IsGreTunnelChecksumFieldEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsGreTunnelChecksumFieldEnabled
 return nil
}

func(np *RsvpSession) SetGreTunnelLocalAddress () error {
 //parameters: SessionHandle LocalAddr
 //AgtRsvpSession SetGreTunnelLocalAddress
 return nil
}

func(np *RsvpSession) GetGreTunnelLocalAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetGreTunnelLocalAddress
 return "", nil
}

func(np *RsvpSession) SetGreTunnelRemoteAddress () error {
 //parameters: SessionHandle RemoteAddr
 //AgtRsvpSession SetGreTunnelRemoteAddress
 return nil
}

func(np *RsvpSession) GetGreTunnelRemoteAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetGreTunnelRemoteAddress
 return "", nil
}

func(np *RsvpSession) GetNullLabelOption ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetNullLabelOption
 return "", nil
}

func(np *RsvpSession) Enable () error {
 //parameters: SessionHandle
 //AgtRsvpSession Enable
 return nil
}

func(np *RsvpSession) Disable () error {
 //parameters: SessionHandle
 //AgtRsvpSession Disable
 return nil
}

func(np *RsvpSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetState
 return "", nil
}

func(np *RsvpSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtRsvpSession GetStateSummary
 return "", nil
}

func(np *RsvpSession) GetTunnelStateSummary ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetTunnelStateSummary
 return "", nil
}

func(np *RsvpSession) GetLastError ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetLastError
 return "", nil
}

func(np *RsvpSession) GetEnableFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetEnableFlag
 return "", nil
}

func(np *RsvpSession) GetDetailedState ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetDetailedState
 return "", nil
}

func(np *RsvpSession) AddIngressPool () error {
 //parameters: SessionHandle
 //AgtRsvpSession AddIngressPool
 return nil
}

func(np *RsvpSession) ReplicatePool () error {
 //parameters: SessionHandle PoolHandle
 //AgtRsvpSession ReplicatePool
 return nil
}

func(np *RsvpSession) RemovePool () error {
 //parameters: SessionHandle PoolHandle
 //AgtRsvpSession RemovePool
 return nil
}

func(np *RsvpSession) ListPools ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession ListPools
 return "", nil
}

func(np *RsvpSession) ListIngressPools ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession ListIngressPools
 return "", nil
}

func(np *RsvpSession) ClearPools () error {
 //parameters: SessionHandle ClearIngressPools ClearEgressPools
 //AgtRsvpSession ClearPools
 return nil
}

func(np *RsvpSession) EnableEgress () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableEgress
 return nil
}

func(np *RsvpSession) DisableEgress () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableEgress
 return nil
}

func(np *RsvpSession) IsEgressEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsEgressEnabled
 return nil
}

func(np *RsvpSession) GetDefaultEgressPoolHandle ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetDefaultEgressPoolHandle
 return "", nil
}

func(np *RsvpSession) EnableGracefulRestart () error {
 //parameters: SessionHandle
 //AgtRsvpSession EnableGracefulRestart
 return nil
}

func(np *RsvpSession) DisableGracefulRestart () error {
 //parameters: SessionHandle
 //AgtRsvpSession DisableGracefulRestart
 return nil
}

func(np *RsvpSession) IsGracefulRestartEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpSession IsGracefulRestartEnabled
 return nil
}

func(np *RsvpSession) GracefulRestart () error {
 //parameters: SessionHandle
 //AgtRsvpSession GracefulRestart
 return nil
}

func(np *RsvpSession) ForceGracefulRestart () error {
 //parameters: SessionHandle GracefulRestartType
 //AgtRsvpSession ForceGracefulRestart
 return nil
}

func(np *RsvpSession) SetRestartCapRestartTime () error {
 //parameters: SessionHandle RestartTime
 //AgtRsvpSession SetRestartCapRestartTime
 return nil
}

func(np *RsvpSession) GetRestartCapRestartTime ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRestartCapRestartTime
 return "", nil
}

func(np *RsvpSession) SetRestartCapRecoveryTime () error {
 //parameters: SessionHandle RecoveryTime
 //AgtRsvpSession SetRestartCapRecoveryTime
 return nil
}

func(np *RsvpSession) GetRestartCapRecoveryTime ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetRestartCapRecoveryTime
 return "", nil
}

func(np *RsvpSession) SetGracefulRestartType () error {
 //parameters: SessionHandle GracefulRestartType
 //AgtRsvpSession SetGracefulRestartType
 return nil
}

func(np *RsvpSession) GetGracefulRestartType ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpSession GetGracefulRestartType
 return "", nil
}

func(np *RsvpSession) SetEgressLabelOverrideRange () error {
 //parameters: MinimumLabel MaximumLabel
 //AgtRsvpSession SetEgressLabelOverrideRange
 return nil
}

func(np *RsvpSession) GetEgressLabelOverrideRange ()(string, error) {
 //parameters: 
 //AgtRsvpSession GetEgressLabelOverrideRange
 return "", nil
}

func(np *RsvpSession) EnableAllSessions () error {
 //parameters: 
 //AgtRsvpSession EnableAllSessions
 return nil
}

func(np *RsvpSession) DisableAllSessions () error {
 //parameters: 
 //AgtRsvpSession DisableAllSessions
 return nil
}

func(np *RsvpSession) SetNullLabelOption () error {
 //parameters: SessionHandle
 //AgtRsvpSession SetNullLabelOption
 return nil
}

