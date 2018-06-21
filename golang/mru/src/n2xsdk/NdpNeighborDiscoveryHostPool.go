package n2xsdk

type NdpNeighborDiscoveryHostPool struct {
 Handler string
}

func(np *NdpNeighborDiscoveryHostPool) EnableAutoSendNeighborSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool EnableAutoSendNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) DisableAutoSendNeighborSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool DisableAutoSendNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) IsAutoSendNeighborSolicitationsEnabled () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool IsAutoSendNeighborSolicitationsEnabled
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetCustomState ()(string, error) {
 //parameters: DevicePoolInstance
 //AgtNdpNeighborDiscoveryHostPool GetCustomState
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SendNeighborSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool SendNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) Reset () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool Reset
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) CancelNeighborSolicitations () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool CancelNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetRemoteMacAddress ()(string, error) {
 //parameters: DevicePoolInstance
 //AgtNdpNeighborDiscoveryHostPool GetRemoteMacAddress
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) EnableStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool EnableStatelessTransmit
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) DisableStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool DisableStatelessTransmit
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) IsStatelessTransmitEnabled () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool IsStatelessTransmitEnabled
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitPacketType () error {
 //parameters: DevicePoolHandle PacketType
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitPacketType
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitPacketType ()(string, error) {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitPacketType
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitPacket () error {
 //parameters: DeviceHandle Packet
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitPacket
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitPacket ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitPacket
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) StartStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool StartStatelessTransmit
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) StopStatelessTransmit () error {
 //parameters: DevicePoolHandle
 //AgtNdpNeighborDiscoveryHostPool StopStatelessTransmit
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitStartEventState () error {
 //parameters: DeviceHandle StartEventState
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitStartEventState
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitStartEventState ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitStartEventState
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SendNeighborSolicitationsByPort () error {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool SendNeighborSolicitationsByPort
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) ResetByPort () error {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool ResetByPort
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) CancelNeighborSolicitationsByPort () error {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool CancelNeighborSolicitationsByPort
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) SendAllNeighborSolicitations () error {
 //parameters: 
 //AgtNdpNeighborDiscoveryHostPool SendAllNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) ResetAll () error {
 //parameters: 
 //AgtNdpNeighborDiscoveryHostPool ResetAll
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) CancelAllNeighborSolicitations () error {
 //parameters: 
 //AgtNdpNeighborDiscoveryHostPool CancelAllNeighborSolicitations
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) SetRateControl () error {
 //parameters: PortHandle RateControlMode
 //AgtNdpNeighborDiscoveryHostPool SetRateControl
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetRateControl ()(string, error) {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool GetRateControl
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetFixedRateProfile () error {
 //parameters: PortHandle FixedRate MaxBurstSize
 //AgtNdpNeighborDiscoveryHostPool SetFixedRateProfile
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetFixedRateProfile ()(string, error) {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool GetFixedRateProfile
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetSelfPacedProfile () error {
 //parameters: PortHandle MaxPackets MaxRetries RetryTimeout
 //AgtNdpNeighborDiscoveryHostPool SetSelfPacedProfile
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetSelfPacedProfile ()(string, error) {
 //parameters: PortHandle
 //AgtNdpNeighborDiscoveryHostPool GetSelfPacedProfile
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitStartEvent () error {
 //parameters: DeviceHandle StartEvent
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitStartEvent
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitStartEvent ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitStartEvent
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitStopEvent () error {
 //parameters: DeviceHandle StopEvent
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitStopEvent
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitStopEvent ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitStopEvent
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitStopEventCount () error {
 //parameters: DeviceHandle StopEventCount
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitStopEventCount
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitStopEventCount ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitStopEventCount
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitStopEventTime () error {
 //parameters: DeviceHandle TimeUnits Time
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitStopEventTime
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitStopEventTime ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitStopEventTime
 return "", nil
}

func(np *NdpNeighborDiscoveryHostPool) SetStatelessTransmitRate () error {
 //parameters: DeviceHandle BurstSize PeriodUnits Period
 //AgtNdpNeighborDiscoveryHostPool SetStatelessTransmitRate
 return nil
}

func(np *NdpNeighborDiscoveryHostPool) GetStatelessTransmitRate ()(string, error) {
 //parameters: DeviceHandle
 //AgtNdpNeighborDiscoveryHostPool GetStatelessTransmitRate
 return "", nil
}

