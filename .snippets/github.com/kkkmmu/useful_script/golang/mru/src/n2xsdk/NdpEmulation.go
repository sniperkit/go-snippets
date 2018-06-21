package n2xsdk

type NdpEmulation struct {
 Handler string
}

func(np *NdpEmulation) GetState ()(string, error) {
 //parameters: PortHandle
 //AgtNdpEmulation GetState
 return "", nil
}

func(np *NdpEmulation) SendNeighborSolicitationsByAddressPool () error {
 //parameters: PortHandle SutIpv6Address AddressPoolHandle
 //AgtNdpEmulation SendNeighborSolicitationsByAddressPool
 return nil
}

func(np *NdpEmulation) SendAllNeighborSolicitationsByAddressPool () error {
 //parameters: PortHandle AddressPoolHandle
 //AgtNdpEmulation SendAllNeighborSolicitationsByAddressPool
 return nil
}

func(np *NdpEmulation) EnableSendAllNdpRequestsAtTestStart () error {
 //parameters: 
 //AgtNdpEmulation EnableSendAllNdpRequestsAtTestStart
 return nil
}

func(np *NdpEmulation) DisableSendAllNdpRequestsAtTestStart () error {
 //parameters: 
 //AgtNdpEmulation DisableSendAllNdpRequestsAtTestStart
 return nil
}

func(np *NdpEmulation) IsSendAllNdpRequestsAtTestStartEnabled () error {
 //parameters: 
 //AgtNdpEmulation IsSendAllNdpRequestsAtTestStartEnabled
 return nil
}

func(np *NdpEmulation) GetSelfPacedState ()(string, error) {
 //parameters: PortHandle
 //AgtNdpEmulation GetSelfPacedState
 return "", nil
}

func(np *NdpEmulation) SendNeighborSolicitations () error {
 //parameters: PortHandle VlanIds SutIpv6Address
 //AgtNdpEmulation SendNeighborSolicitations
 return nil
}

func(np *NdpEmulation) SendAllNeighborSolicitations () error {
 //parameters: PortHandle
 //AgtNdpEmulation SendAllNeighborSolicitations
 return nil
}

func(np *NdpEmulation) SendNeighborSolicitationsBySession () error {
 //parameters: PortHandle SutIpv6Address SessionPoolHandle
 //AgtNdpEmulation SendNeighborSolicitationsBySession
 return nil
}

func(np *NdpEmulation) SendAllNeighborSolicitationsBySession () error {
 //parameters: PortHandle SessionPoolHandle
 //AgtNdpEmulation SendAllNeighborSolicitationsBySession
 return nil
}

func(np *NdpEmulation) ClearMacTable () error {
 //parameters: PortHandle
 //AgtNdpEmulation ClearMacTable
 return nil
}

func(np *NdpEmulation) Enable () error {
 //parameters: PortHandle
 //AgtNdpEmulation Enable
 return nil
}

func(np *NdpEmulation) Disable () error {
 //parameters: PortHandle
 //AgtNdpEmulation Disable
 return nil
}

func(np *NdpEmulation) EnableSelfPacedResolution () error {
 //parameters: PortHandle
 //AgtNdpEmulation EnableSelfPacedResolution
 return nil
}

func(np *NdpEmulation) DisableSelfPacedResolution () error {
 //parameters: PortHandle
 //AgtNdpEmulation DisableSelfPacedResolution
 return nil
}

