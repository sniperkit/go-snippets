package n2xsdk

type FrameRelayInterface struct {
 Handler string
}

func(np *FrameRelayInterface) ScramblerOn () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface ScramblerOn
 return nil
}

func(np *FrameRelayInterface) ScramblerOff () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface ScramblerOff
 return nil
}

func(np *FrameRelayInterface) IsScramblerOn () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface IsScramblerOn
 return nil
}

func(np *FrameRelayInterface) DescramblerOn () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface DescramblerOn
 return nil
}

func(np *FrameRelayInterface) DescramblerOff () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface DescramblerOff
 return nil
}

func(np *FrameRelayInterface) IsDescramblerOn () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface IsDescramblerOn
 return nil
}

func(np *FrameRelayInterface) SetFcsLength () error {
 //parameters: PortHandle FcsLength
 //AgtFrameRelayInterface SetFcsLength
 return nil
}

func(np *FrameRelayInterface) GetFcsLength ()(string, error) {
 //parameters: PortHandle
 //AgtFrameRelayInterface GetFcsLength
 return "", nil
}

func(np *FrameRelayInterface) SetMinimumInterFrameGap () error {
 //parameters: PortHandle InterFrameGap
 //AgtFrameRelayInterface SetMinimumInterFrameGap
 return nil
}

func(np *FrameRelayInterface) GetMinimumInterFrameGap ()(string, error) {
 //parameters: PortHandle
 //AgtFrameRelayInterface GetMinimumInterFrameGap
 return "", nil
}

func(np *FrameRelayInterface) GetSupportedMinimumInterFrameGaps ()(string, error) {
 //parameters: PortHandle
 //AgtFrameRelayInterface GetSupportedMinimumInterFrameGaps
 return "", nil
}

func(np *FrameRelayInterface) SetMtuSize () error {
 //parameters: PortHandle MtuSize
 //AgtFrameRelayInterface SetMtuSize
 return nil
}

func(np *FrameRelayInterface) GetMtuSize ()(string, error) {
 //parameters: PortHandle
 //AgtFrameRelayInterface GetMtuSize
 return "", nil
}

func(np *FrameRelayInterface) SetMplsEncapsulation () error {
 //parameters: PortHandle Encapsulation
 //AgtFrameRelayInterface SetMplsEncapsulation
 return nil
}

func(np *FrameRelayInterface) GetMplsEncapsulation ()(string, error) {
 //parameters: PortHandle
 //AgtFrameRelayInterface GetMplsEncapsulation
 return "", nil
}

func(np *FrameRelayInterface) IsFrameRelaySupported () error {
 //parameters: PortHandle
 //AgtFrameRelayInterface IsFrameRelaySupported
 return nil
}

