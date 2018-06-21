package n2xsdk

type NdpHostPool struct {
 Handler string
}

func(np *NdpHostPool) EnableAutoSendRouterSolicitation () error {
 //parameters: DevicePoolHandle
 //AgtNdpHostPool EnableAutoSendRouterSolicitation
 return nil
}

func(np *NdpHostPool) DisableAutoSendRouterSolicitation () error {
 //parameters: DevicePoolHandle
 //AgtNdpHostPool DisableAutoSendRouterSolicitation
 return nil
}

func(np *NdpHostPool) IsAutoSendRouterSolicitationEnabled () error {
 //parameters: DevicePoolHandle
 //AgtNdpHostPool IsAutoSendRouterSolicitationEnabled
 return nil
}

func(np *NdpHostPool) SendRouterSolicitation () error {
 //parameters: DeviceIdentifiers
 //AgtNdpHostPool SendRouterSolicitation
 return nil
}

func(np *NdpHostPool) ResetAddressAutoconfig () error {
 //parameters: DeviceIdentifiers
 //AgtNdpHostPool ResetAddressAutoconfig
 return nil
}

func(np *NdpHostPool) GetAddressAutoconfigState ()(string, error) {
 //parameters: DeviceIdentifiers
 //AgtNdpHostPool GetAddressAutoconfigState
 return "", nil
}

