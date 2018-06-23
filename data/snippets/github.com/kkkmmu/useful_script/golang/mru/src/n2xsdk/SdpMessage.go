package n2xsdk

type SdpMessage struct {
 Handler string
}

func(np *SdpMessage) IsSdpInterfaceEnabled () error {
 //parameters: DeviceHandle
 //AgtSdpMessage IsSdpInterfaceEnabled
 return nil
}

func(np *SdpMessage) EnableSdpInterface () error {
 //parameters: DeviceHandle
 //AgtSdpMessage EnableSdpInterface
 return nil
}

func(np *SdpMessage) DisableSdpInterface () error {
 //parameters: DeviceHandle
 //AgtSdpMessage DisableSdpInterface
 return nil
}

func(np *SdpMessage) IsUserDefinedMessageEnabled () error {
 //parameters: DeviceHandle
 //AgtSdpMessage IsUserDefinedMessageEnabled
 return nil
}

func(np *SdpMessage) EnableUserDefinedMessage () error {
 //parameters: DeviceHandle
 //AgtSdpMessage EnableUserDefinedMessage
 return nil
}

func(np *SdpMessage) DisableUserDefinedMessage () error {
 //parameters: DeviceHandle
 //AgtSdpMessage DisableUserDefinedMessage
 return nil
}

func(np *SdpMessage) SetUserDefinedMessage () error {
 //parameters: DeviceHandle EveryNthInstance UserDefinedMessage
 //AgtSdpMessage SetUserDefinedMessage
 return nil
}

func(np *SdpMessage) GetUserDefinedMessage ()(string, error) {
 //parameters: DeviceHandle
 //AgtSdpMessage GetUserDefinedMessage
 return "", nil
}

func(np *SdpMessage) AddMediaLine () error {
 //parameters: DeviceHandle
 //AgtSdpMessage AddMediaLine
 return nil
}

func(np *SdpMessage) RemoveMediaLine () error {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage RemoveMediaLine
 return nil
}

func(np *SdpMessage) ListMediaLines ()(string, error) {
 //parameters: DeviceHandle
 //AgtSdpMessage ListMediaLines
 return "", nil
}

func(np *SdpMessage) IsLocalMediaIpAddressEnabled () error {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage IsLocalMediaIpAddressEnabled
 return nil
}

func(np *SdpMessage) EnableLocalMediaIpAddress () error {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage EnableLocalMediaIpAddress
 return nil
}

func(np *SdpMessage) DisableLocalMediaIpAddress () error {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage DisableLocalMediaIpAddress
 return nil
}

func(np *SdpMessage) SetLocalMediaIpv4Address () error {
 //parameters: DeviceHandle MediaRowId Ipv4Address Increment Repeat
 //AgtSdpMessage SetLocalMediaIpv4Address
 return nil
}

func(np *SdpMessage) GetLocalMediaIpv4Address ()(string, error) {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage GetLocalMediaIpv4Address
 return "", nil
}

func(np *SdpMessage) SetLocalMediaIpv6Address () error {
 //parameters: DeviceHandle MediaRowId Ipv6Address Increment Repeat
 //AgtSdpMessage SetLocalMediaIpv6Address
 return nil
}

func(np *SdpMessage) GetLocalMediaIpv6Address ()(string, error) {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage GetLocalMediaIpv6Address
 return "", nil
}

func(np *SdpMessage) SetMediaType () error {
 //parameters: DeviceHandle MediaRowId MediaType
 //AgtSdpMessage SetMediaType
 return nil
}

func(np *SdpMessage) GetMediaType ()(string, error) {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage GetMediaType
 return "", nil
}

func(np *SdpMessage) AddMediaPayload () error {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage AddMediaPayload
 return nil
}

func(np *SdpMessage) RemoveMediaPayload () error {
 //parameters: DeviceHandle MediaLineRowId MediaPayloadLineRowId
 //AgtSdpMessage RemoveMediaPayload
 return nil
}

func(np *SdpMessage) ListMediaPayloadRows ()(string, error) {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage ListMediaPayloadRows
 return "", nil
}

func(np *SdpMessage) ListMediaPayloads ()(string, error) {
 //parameters: DeviceHandle MediaLineRowId
 //AgtSdpMessage ListMediaPayloads
 return "", nil
}

func(np *SdpMessage) SetAudioPayload () error {
 //parameters: DeviceHandle MediaLineRowId MediaPayloadLineRowId AudioPayload
 //AgtSdpMessage SetAudioPayload
 return nil
}

func(np *SdpMessage) GetAudioPayload ()(string, error) {
 //parameters: DeviceHandle MediaLineRowId MediaPayloadLineRowId
 //AgtSdpMessage GetAudioPayload
 return "", nil
}

func(np *SdpMessage) SetMediaPort () error {
 //parameters: DeviceHandle MediaRowId MediaPort Increment Repeat
 //AgtSdpMessage SetMediaPort
 return nil
}

func(np *SdpMessage) GetMediaPort ()(string, error) {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage GetMediaPort
 return "", nil
}

func(np *SdpMessage) SetMediaAttribute () error {
 //parameters: DeviceHandle MediaRowId MediaAttribute
 //AgtSdpMessage SetMediaAttribute
 return nil
}

func(np *SdpMessage) GetMediaAttribute ()(string, error) {
 //parameters: DeviceHandle MediaRowId
 //AgtSdpMessage GetMediaAttribute
 return "", nil
}

