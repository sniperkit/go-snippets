package n2xsdk

type NdpAddressAutoconfigHostPool struct {
 Handler string
}

func(np *NdpAddressAutoconfigHostPool) EnableAutoSendRouterSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool EnableAutoSendRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) DisableAutoSendRouterSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool DisableAutoSendRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) IsAutoSendRouterSolicitationsEnabled () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool IsAutoSendRouterSolicitationsEnabled
 return nil
}

func(np *NdpAddressAutoconfigHostPool) SendRouterSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool SendRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) Reset () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool Reset
 return nil
}

func(np *NdpAddressAutoconfigHostPool) CancelRouterSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool CancelRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetCustomState ()(string, error) {
 //parameters: DevicePoolInstance
 //AgtNdpAddressAutoconfigHostPool GetCustomState
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) EnableStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool EnableStatelessTransmit
 return nil
}

func(np *NdpAddressAutoconfigHostPool) DisableStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool DisableStatelessTransmit
 return nil
}

func(np *NdpAddressAutoconfigHostPool) IsStatelessTransmitEnabled () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool IsStatelessTransmitEnabled
 return nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitPacketType () error {
 //parameters: DevicePoolHandle PacketType
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitPacketType
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitPacketType ()(string, error) {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitPacketType
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitPacket () error {
 //parameters: DeviceHandle Packet
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitPacket
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitPacket ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitPacket
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) StartStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool StartStatelessTransmit
 return nil
}

func(np *NdpAddressAutoconfigHostPool) StopStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostPool StopStatelessTransmit
 return nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitStartEventState () error {
 //parameters: DeviceHandle StartEventState
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitStartEventState
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitStartEventState ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitStartEventState
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SendRouterSolicitationsByPort () error {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool SendRouterSolicitationsByPort
 return nil
}

func(np *NdpAddressAutoconfigHostPool) ResetByPort () error {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool ResetByPort
 return nil
}

func(np *NdpAddressAutoconfigHostPool) CancelRouterSolicitationsByPort () error {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool CancelRouterSolicitationsByPort
 return nil
}

func(np *NdpAddressAutoconfigHostPool) SendAllRouterSolicitations () error {
 //parameters: 
 //AgtNdpAddressAutoconfigHostPool SendAllRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) ResetAll () error {
 //parameters: 
 //AgtNdpAddressAutoconfigHostPool ResetAll
 return nil
}

func(np *NdpAddressAutoconfigHostPool) CancelAllRouterSolicitations () error {
 //parameters: 
 //AgtNdpAddressAutoconfigHostPool CancelAllRouterSolicitations
 return nil
}

func(np *NdpAddressAutoconfigHostPool) SetRateControl () error {
 //parameters: PortHandle RateControlMode
 //AgtNdpAddressAutoconfigHostPool SetRateControl
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetRateControl ()(string, error) {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool GetRateControl
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetFixedRateProfile () error {
 //parameters: PortHandle FixedRate MaxBurstSize
 //AgtNdpAddressAutoconfigHostPool SetFixedRateProfile
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetFixedRateProfile ()(string, error) {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool GetFixedRateProfile
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetSelfPacedProfile () error {
 //parameters: PortHandle MaxPackets MaxRetries RetryTimeout
 //AgtNdpAddressAutoconfigHostPool SetSelfPacedProfile
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetSelfPacedProfile ()(string, error) {
 //parameters: PortHandle
 //AgtNdpAddressAutoconfigHostPool GetSelfPacedProfile
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitStartEvent () error {
 //parameters: DeviceHandle StartEvent
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitStartEvent
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitStartEvent ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitStartEvent
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitStopEvent () error {
 //parameters: DeviceHandle StopEvent
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitStopEvent
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitStopEvent ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitStopEvent
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitStopEventCount () error {
 //parameters: DeviceHandle StopEventCount
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitStopEventCount
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitStopEventCount ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitStopEventCount
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitStopEventTime () error {
 //parameters: DeviceHandle TimeUnits Time
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitStopEventTime
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitStopEventTime ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitStopEventTime
 return "", nil
}

func(np *NdpAddressAutoconfigHostPool) SetStatelessTransmitRate () error {
 //parameters: DeviceHandle BurstSize PeriodUnits Period
 //AgtNdpAddressAutoconfigHostPool SetStatelessTransmitRate
 return nil
}

func(np *NdpAddressAutoconfigHostPool) GetStatelessTransmitRate ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpAddressAutoconfigHostPool GetStatelessTransmitRate
 return "", nil
}

