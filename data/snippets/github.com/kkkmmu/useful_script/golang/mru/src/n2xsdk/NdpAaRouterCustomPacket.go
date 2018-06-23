package n2xsdk

type NdpAaRouterCustomPacket struct {
 Handler string
}

func(np *NdpAaRouterCustomPacke) AddPdu () error {
 //parameters: DeviceHandle PduHandleList
 //AgtNdpAaRouterCustomPacket AddPdu
 return nil
}

func(np *NdpAaRouterCustomPacke) RemovePdu () error {
 //parameters: DeviceHandle
 //AgtNdpAaRouterCustomPacket RemovePdu
 return nil
}

