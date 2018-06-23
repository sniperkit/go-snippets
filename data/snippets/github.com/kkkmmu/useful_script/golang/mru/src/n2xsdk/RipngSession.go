package n2xsdk

type RipngSession struct {
 Handler string
}

func(np *RipngSession) SetInterfaceIpAddress () error {
 //parameters: SessionHandle LocalIpAddress
 //AgtRipngSession SetInterfaceIpAddress
 return nil
}

func(np *RipngSession) SetSutIpAddress () error {
 //parameters: SessionHandle RemoteIpAddress
 //AgtRipngSession SetSutIpAddress
 return nil
}

func(np *RipngSession) SetRouterId () error {
 //parameters: SessionHandle RouterId
 //AgtRipngSession SetRouterId
 return nil
}

func(np *RipngSession) SetSubnetMask () error {
 //parameters: SessionHandle SubnetMask
 //AgtRipngSession SetSubnetMask
 return nil
}

func(np *RipngSession) SetMetric () error {
 //parameters: SessionHandle Metric
 //AgtRipngSession SetMetric
 return nil
}

func(np *RipngSession) SetUpdateInterval () error {
 //parameters: SessionHandle UpdateInterval
 //AgtRipngSession SetUpdateInterval
 return nil
}

func(np *RipngSession) SetTriggeredInterval () error {
 //parameters: SessionHandle TriggeredInterval
 //AgtRipngSession SetTriggeredInterval
 return nil
}

func(np *RipngSession) SetUpdateControl () error {
 //parameters: SessionHandle UpdateControl
 //AgtRipngSession SetUpdateControl
 return nil
}

func(np *RipngSession) SetExpirationInterval () error {
 //parameters: SessionHandle ExpirationInterval
 //AgtRipngSession SetExpirationInterval
 return nil
}

func(np *RipngSession) SetGarbageInterval () error {
 //parameters: SessionHandle GarbageInterval
 //AgtRipngSession SetGarbageInterval
 return nil
}

func(np *RipngSession) Enable () error {
 //parameters: SessionHandle
 //AgtRipngSession Enable
 return nil
}

func(np *RipngSession) Disable () error {
 //parameters: SessionHandle
 //AgtRipngSession Disable
 return nil
}

func(np *RipngSession) EnableAllSessions () error {
 //parameters: 
 //AgtRipngSession EnableAllSessions
 return nil
}

func(np *RipngSession) DisableAllSessions () error {
 //parameters: 
 //AgtRipngSession DisableAllSessions
 return nil
}

func(np *RipngSession) GetInterfaceIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetInterfaceIpAddress
 return "", nil
}

func(np *RipngSession) GetSutIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetSutIpAddress
 return "", nil
}

func(np *RipngSession) GetRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetRouterId
 return "", nil
}

func(np *RipngSession) GetSubnetMask ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetSubnetMask
 return "", nil
}

func(np *RipngSession) GetMetric ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetMetric
 return "", nil
}

func(np *RipngSession) GetUpdateInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetUpdateInterval
 return "", nil
}

func(np *RipngSession) GetTriggeredInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetTriggeredInterval
 return "", nil
}

func(np *RipngSession) GetUpdateControl ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetUpdateControl
 return "", nil
}

func(np *RipngSession) GetExpirationInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetExpirationInterval
 return "", nil
}

func(np *RipngSession) GetGarbageInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetGarbageInterval
 return "", nil
}

func(np *RipngSession) AdvertiseAllRoutePools () error {
 //parameters: SessionHandle
 //AgtRipngSession AdvertiseAllRoutePools
 return nil
}

func(np *RipngSession) WithdrawAllRoutePools () error {
 //parameters: SessionHandle
 //AgtRipngSession WithdrawAllRoutePools
 return nil
}

func(np *RipngSession) StartFlapAll () error {
 //parameters: 
 //AgtRipngSession StartFlapAll
 return nil
}

func(np *RipngSession) StopFlapAll () error {
 //parameters: 
 //AgtRipngSession StopFlapAll
 return nil
}

func(np *RipngSession) EnableAllRoutePools () error {
 //parameters: 
 //AgtRipngSession EnableAllRoutePools
 return nil
}

func(np *RipngSession) DisableAllRoutePools () error {
 //parameters: 
 //AgtRipngSession DisableAllRoutePools
 return nil
}

func(np *RipngSession) IsEnabled () error {
 //parameters: SessionHandle
 //AgtRipngSession IsEnabled
 return nil
}

func(np *RipngSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtRipngSession GetStateSummary
 return "", nil
}

func(np *RipngSession) GetLastError ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetLastError
 return "", nil
}

func(np *RipngSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngSession GetState
 return "", nil
}

