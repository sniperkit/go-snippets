package n2xsdk

type EmulationPduList struct {
 Handler string
}

func(np *EmulationPduLis) AddPdu () error {
 //parameters: Count ProtocolList
 //AgtEmulationPduList AddPdu
 return nil
}

func(np *EmulationPduLis) AddPduWithPacketType () error {
 //parameters: OuterProtocol PacketType
 //AgtEmulationPduList AddPduWithPacketType
 return nil
}

func(np *EmulationPduLis) RemovePdu () error {
 //parameters: EmulationPduHandle
 //AgtEmulationPduList RemovePdu
 return nil
}

func(np *EmulationPduLis) SetName () error {
 //parameters: EmulationPduHandle Name
 //AgtEmulationPduList SetName
 return nil
}

func(np *EmulationPduLis) GetName ()(string, error) {
 //parameters: EmulationPduHandle
 //AgtEmulationPduList GetName
 return "", nil
}

func(np *EmulationPduLis) GetHandle ()(string, error) {
 //parameters: Name
 //AgtEmulationPduList GetHandle
 return "", nil
}

func(np *EmulationPduLis) ListHandles ()(string, error) {
 //parameters: 
 //AgtEmulationPduList ListHandles
 return "", nil
}

func(np *EmulationPduLis) ListHandlesByOuterProtocol ()(string, error) {
 //parameters: OuterProtocol
 //AgtEmulationPduList ListHandlesByOuterProtocol
 return "", nil
}

