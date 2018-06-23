package n2xsdk

type AncpClientCustomPacket struct {
 Handler string
}

func(np *ncpClientCustomPacke) AddPdu () error {
 //parameters: PeerPoolHandle PduHandleList
 //AgtAncpClientCustomPacket AddPdu
 return nil
}

func(np *ncpClientCustomPacke) RemovePdu () error {
 //parameters: PeerPoolHandle
 //AgtAncpClientCustomPacket RemovePdu
 return nil
}

func(np *ncpClientCustomPacke) SendCustomPacket () error {
 //parameters: PeerPoolHandle
 //AgtAncpClientCustomPacket SendCustomPacket
 return nil
}

