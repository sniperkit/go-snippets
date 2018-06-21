package n2xsdk

type FrameRelayPvc struct {
 Handler string
}

func(np *FrameRelayPvc) SetHeaderFields () error {
 //parameters: PvcHandle Dlci Fecn Becn CR DE
 //AgtFrameRelayPvc SetHeaderFields
 return nil
}

func(np *FrameRelayPvc) GetHeaderFields ()(string, error) {
 //parameters: PvcHandle
 //AgtFrameRelayPvc GetHeaderFields
 return "", nil
}

func(np *FrameRelayPvc) SetEncapsulation () error {
 //parameters: PvcHandle Encapsulation
 //AgtFrameRelayPvc SetEncapsulation
 return nil
}

func(np *FrameRelayPvc) GetEncapsulation ()(string, error) {
 //parameters: PvcHandle
 //AgtFrameRelayPvc GetEncapsulation
 return "", nil
}

func(np *FrameRelayPvc) SetIpAddresses () error {
 //parameters: PvcHandle TesterIpAddress SutIpAddress PrefixLength
 //AgtFrameRelayPvc SetIpAddresses
 return nil
}

func(np *FrameRelayPvc) GetIpAddresses ()(string, error) {
 //parameters: PvcHandle
 //AgtFrameRelayPvc GetIpAddresses
 return "", nil
}

func(np *FrameRelayPvc) GetPvcParameters ()(string, error) {
 //parameters: PvcHandle
 //AgtFrameRelayPvc GetPvcParameters
 return "", nil
}

