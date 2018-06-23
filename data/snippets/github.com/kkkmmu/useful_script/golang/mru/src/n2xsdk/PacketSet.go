package n2xsdk

type PacketSet struct {
 Handler string
}

func(np *PacketSe) GetPacketSetStats ()(string, error) {
 //parameters: PacketSetHandle
 //AgtPacketSet GetPacketSetStats
 return "", nil
}

func(np *PacketSe) GetPacketSource ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketSource
 return "", nil
}

func(np *PacketSe) GetPacketNumber ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketNumber
 return "", nil
}

func(np *PacketSe) GetPacketLength ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketLength
 return "", nil
}

func(np *PacketSe) GetPacketDirection ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketDirection
 return "", nil
}

func(np *PacketSe) GetPacketSourceIndex ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketSourceIndex
 return "", nil
}

func(np *PacketSe) GetLayer1Mode ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetLayer1Mode
 return "", nil
}

func(np *PacketSe) GetLayer2Offset ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetLayer2Offset
 return "", nil
}

func(np *PacketSe) GetLayer2Encapsulation ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetLayer2Encapsulation
 return "", nil
}

func(np *PacketSe) GetLayer2HeaderLength ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetLayer2HeaderLength
 return "", nil
}

func(np *PacketSe) GetLayer3PayloadOffset ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetLayer3PayloadOffset
 return "", nil
}

func(np *PacketSe) GetReceiveTimestamp ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetReceiveTimestamp
 return "", nil
}

func(np *PacketSe) GetPacketOctets ()(string, error) {
 //parameters: PacketSetHandle PacketIndex OctetOffset
 //AgtPacketSet GetPacketOctets
 return "", nil
}

func(np *PacketSe) GetTestPayloadOffset ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetTestPayloadOffset
 return "", nil
}

func(np *PacketSe) GetTestPayloadField ()(string, error) {
 //parameters: PacketSetHandle PacketIndex Field
 //AgtPacketSet GetTestPayloadField
 return "", nil
}

func(np *PacketSe) GetPacketLatency ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetPacketLatency
 return "", nil
}

func(np *PacketSe) GetStatusFlags ()(string, error) {
 //parameters: PacketSetHandle PacketIndex
 //AgtPacketSet GetStatusFlags
 return "", nil
}

func(np *PacketSe) GetProtocolField ()(string, error) {
 //parameters: PacketSetHandle PacketIndex Protocol ProtocolOffset ProtocolField
 //AgtPacketSet GetProtocolField
 return "", nil
}

func(np *PacketSe) GetProtocolField32 ()(string, error) {
 //parameters: PacketSetHandle PacketIndex Protocol ProtocolOffset ProtocolField
 //AgtPacketSet GetProtocolField32
 return "", nil
}

func(np *PacketSe) ReloadProtocolFile () error {
 //parameters: 
 //AgtPacketSet ReloadProtocolFile
 return nil
}

func(np *PacketSe) GetEntirePacketWithOverhead ()(string, error) {
 //parameters: PacketSetHandle PacketIndex NumRequestedTestFields psaRequestedFields
 //AgtPacketSet GetEntirePacketWithOverhead
 return "", nil
}

