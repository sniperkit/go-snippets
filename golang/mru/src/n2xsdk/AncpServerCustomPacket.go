package n2xsdk

type AncpServerCustomPacket struct {
 Handler string
}

func(np *ncpServerCustomPacke) AddPdu () error {
 //parameters: PeerPoolHandle PduHandleList
 //AgtAncpServerCustomPacket AddPdu
 return nil
}

func(np *ncpServerCustomPacke) RemovePdu () error {
 //parameters: PeerPoolHandle
 //AgtAncpServerCustomPacket RemovePdu
 return nil
}

func(np *ncpServerCustomPacke) SendCustomPacket () error {
 //parameters: PeerPoolHandle
 //AgtAncpServerCustomPacket SendCustomPacket
 return nil
}

