package n2xsdk

type FrameTransmitter struct {
 Handler string
}

func(np *FrameTransmitter) AddPdu () error {
 //parameters: Count ProtocolList PduLength
 //AgtFrameTransmitter AddPdu
 return nil
}

func(np *FrameTransmitter) AddPduWithPacketType () error {
 //parameters: OuterProtocol PacketType PduLength
 //AgtFrameTransmitter AddPduWithPacketType
 return nil
}

func(np *FrameTransmitter) RemovePdu () error {
 //parameters: PduHandle
 //AgtFrameTransmitter RemovePdu
 return nil
}

func(np *FrameTransmitter) GetPduLength ()(string, error) {
 //parameters: PduHandle
 //AgtFrameTransmitter GetPduLength
 return "", nil
}

func(np *FrameTransmitter) SetPduLength () error {
 //parameters: PduHandle PduLength
 //AgtFrameTransmitter SetPduLength
 return nil
}

func(np *FrameTransmitter) TransmitPdu () error {
 //parameters: PortHandle PduHandle TransmitOptions
 //AgtFrameTransmitter TransmitPdu
 return nil
}

func(np *FrameTransmitter) TransmitFrame () error {
 //parameters: PortHandle FrameLength psaFrameBytes TransmitOptions
 //AgtFrameTransmitter TransmitFrame
 return nil
}

func(np *FrameTransmitter) GetLastFrameDetails ()(string, error) {
 //parameters: PortHandle
 //AgtFrameTransmitter GetLastFrameDetails
 return "", nil
}

