package n2xsdk

type LmpSession struct {
 Handler string
}

func(np *LmpSession) NegotiateHelloConfig () error {
 //parameters: SessionHandle
 //AgtLmpSession NegotiateHelloConfig
 return nil
}

func(np *LmpSession) SetLocalNodeId () error {
 //parameters: SessionHandle NodeId
 //AgtLmpSession SetLocalNodeId
 return nil
}

func(np *LmpSession) GetLocalNodeId ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetLocalNodeId
 return "", nil
}

func(np *LmpSession) SetRemoteNodeId () error {
 //parameters: SessionHandle NodeId
 //AgtLmpSession SetRemoteNodeId
 return nil
}

func(np *LmpSession) GetRemoteNodeId ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetRemoteNodeId
 return "", nil
}

func(np *LmpSession) SetTesterIpccAddress () error {
 //parameters: SessionHandle LocalIpccAddress
 //AgtLmpSession SetTesterIpccAddress
 return nil
}

func(np *LmpSession) GetTesterIpccAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetTesterIpccAddress
 return "", nil
}

func(np *LmpSession) SetSutIpccAddress () error {
 //parameters: SessionHandle RemoteIpccAddress
 //AgtLmpSession SetSutIpccAddress
 return nil
}

func(np *LmpSession) GetSutIpccAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetSutIpccAddress
 return "", nil
}

func(np *LmpSession) SetCcId () error {
 //parameters: SessionHandle CcId
 //AgtLmpSession SetCcId
 return nil
}

func(np *LmpSession) GetCcId ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetCcId
 return "", nil
}

func(np *LmpSession) Enable () error {
 //parameters: SessionHandle
 //AgtLmpSession Enable
 return nil
}

func(np *LmpSession) Disable () error {
 //parameters: SessionHandle
 //AgtLmpSession Disable
 return nil
}

func(np *LmpSession) EnableGreTunnel () error {
 //parameters: SessionHandle
 //AgtLmpSession EnableGreTunnel
 return nil
}

func(np *LmpSession) DisableGreTunnel () error {
 //parameters: SessionHandle
 //AgtLmpSession DisableGreTunnel
 return nil
}

func(np *LmpSession) IsGreTunnelEnabled () error {
 //parameters: SessionHandle
 //AgtLmpSession IsGreTunnelEnabled
 return nil
}

func(np *LmpSession) SetGreTunnelLocalAddress () error {
 //parameters: SessionHandle LocalAddr
 //AgtLmpSession SetGreTunnelLocalAddress
 return nil
}

func(np *LmpSession) GetGreTunnelLocalAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetGreTunnelLocalAddress
 return "", nil
}

func(np *LmpSession) SetGreTunnelRemoteAddress () error {
 //parameters: SessionHandle RemoteAddr
 //AgtLmpSession SetGreTunnelRemoteAddress
 return nil
}

func(np *LmpSession) GetGreTunnelRemoteAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetGreTunnelRemoteAddress
 return "", nil
}

func(np *LmpSession) EnableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtLmpSession EnableGreTunnelChecksumField
 return nil
}

func(np *LmpSession) DisableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtLmpSession DisableGreTunnelChecksumField
 return nil
}

func(np *LmpSession) IsGreTunnelChecksumFieldEnabled () error {
 //parameters: SessionHandle
 //AgtLmpSession IsGreTunnelChecksumFieldEnabled
 return nil
}

func(np *LmpSession) SetConfigRetransmitInterval () error {
 //parameters: SessionHandle Interval
 //AgtLmpSession SetConfigRetransmitInterval
 return nil
}

func(np *LmpSession) GetConfigRetransmitInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetConfigRetransmitInterval
 return "", nil
}

func(np *LmpSession) SetHelloInterval () error {
 //parameters: SessionHandle HelloInterval
 //AgtLmpSession SetHelloInterval
 return nil
}

func(np *LmpSession) GetHelloInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetHelloInterval
 return "", nil
}

func(np *LmpSession) SetHelloDeadInterval () error {
 //parameters: SessionHandle HelloDeadInterval
 //AgtLmpSession SetHelloDeadInterval
 return nil
}

func(np *LmpSession) GetHelloDeadInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetHelloDeadInterval
 return "", nil
}

func(np *LmpSession) EnableHelloConfigNegotiable () error {
 //parameters: SessionHandle
 //AgtLmpSession EnableHelloConfigNegotiable
 return nil
}

func(np *LmpSession) DisableHelloConfigNegotiable () error {
 //parameters: SessionHandle
 //AgtLmpSession DisableHelloConfigNegotiable
 return nil
}

func(np *LmpSession) IsHelloConfigNegotiable () error {
 //parameters: SessionHandle
 //AgtLmpSession IsHelloConfigNegotiable
 return nil
}

func(np *LmpSession) SetActiveConfig () error {
 //parameters: SessionHandle
 //AgtLmpSession SetActiveConfig
 return nil
}

func(np *LmpSession) SetPassiveConfig () error {
 //parameters: SessionHandle
 //AgtLmpSession SetPassiveConfig
 return nil
}

func(np *LmpSession) IsActiveConfig () error {
 //parameters: SessionHandle
 //AgtLmpSession IsActiveConfig
 return nil
}

func(np *LmpSession) Open () error {
 //parameters: SessionHandle
 //AgtLmpSession Open
 return nil
}

func(np *LmpSession) OpenAllSessions () error {
 //parameters: 
 //AgtLmpSession OpenAllSessions
 return nil
}

func(np *LmpSession) Close () error {
 //parameters: SessionHandle
 //AgtLmpSession Close
 return nil
}

func(np *LmpSession) CloseAllSessions () error {
 //parameters: 
 //AgtLmpSession CloseAllSessions
 return nil
}

func(np *LmpSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetState
 return "", nil
}

func(np *LmpSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtLmpSession GetStateSummary
 return "", nil
}

func(np *LmpSession) SetTeLinkCorrelationMode () error {
 //parameters: SessionHandle TeLinkCorrelationMode
 //AgtLmpSession SetTeLinkCorrelationMode
 return nil
}

func(np *LmpSession) GetTeLinkCorrelationMode ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession GetTeLinkCorrelationMode
 return "", nil
}

func(np *LmpSession) AddTeLink () error {
 //parameters: SessionHandle EmulationPduHandle
 //AgtLmpSession AddTeLink
 return nil
}

func(np *LmpSession) RemoveTeLink () error {
 //parameters: SessionHandle TeLinkHandle
 //AgtLmpSession RemoveTeLink
 return nil
}

func(np *LmpSession) ListTeLinks ()(string, error) {
 //parameters: SessionHandle
 //AgtLmpSession ListTeLinks
 return "", nil
}

func(np *LmpSession) EnableAllSessions () error {
 //parameters: 
 //AgtLmpSession EnableAllSessions
 return nil
}

func(np *LmpSession) DisableAllSessions () error {
 //parameters: 
 //AgtLmpSession DisableAllSessions
 return nil
}

