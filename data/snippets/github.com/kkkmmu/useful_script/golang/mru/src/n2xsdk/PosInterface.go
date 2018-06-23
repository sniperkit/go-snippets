package n2xsdk

type PosInterface struct {
 Handler string
}

func(np *PosInterface) SetFramingMode () error {
 //parameters: PortHandle Framing
 //AgtPosInterface SetFramingMode
 return nil
}

func(np *PosInterface) GetFramingMode ()(string, error) {
 //parameters: PortHandle
 //AgtPosInterface GetFramingMode
 return "", nil
}

func(np *PosInterface) GetSupportedFramingModes ()(string, error) {
 //parameters: PortHandle
 //AgtPosInterface GetSupportedFramingModes
 return "", nil
}

func(np *PosInterface) ScramblerOn () error {
 //parameters: PortHandle
 //AgtPosInterface ScramblerOn
 return nil
}

func(np *PosInterface) ScramblerOff () error {
 //parameters: PortHandle
 //AgtPosInterface ScramblerOff
 return nil
}

func(np *PosInterface) IsScramblerOn () error {
 //parameters: PortHandle
 //AgtPosInterface IsScramblerOn
 return nil
}

func(np *PosInterface) DescramblerOn () error {
 //parameters: PortHandle
 //AgtPosInterface DescramblerOn
 return nil
}

func(np *PosInterface) DescramblerOff () error {
 //parameters: PortHandle
 //AgtPosInterface DescramblerOff
 return nil
}

func(np *PosInterface) IsDescramblerOn () error {
 //parameters: PortHandle
 //AgtPosInterface IsDescramblerOn
 return nil
}

func(np *PosInterface) IsPosSupported () error {
 //parameters: PortHandle
 //AgtPosInterface IsPosSupported
 return nil
}

func(np *PosInterface) SetMinimumInterFrameGap () error {
 //parameters: PortHandle InterFrameGap
 //AgtPosInterface SetMinimumInterFrameGap
 return nil
}

func(np *PosInterface) GetMinimumInterFrameGap ()(string, error) {
 //parameters: PortHandle
 //AgtPosInterface GetMinimumInterFrameGap
 return "", nil
}

func(np *PosInterface) GetSupportedMinimumInterFrameGaps ()(string, error) {
 //parameters: PortHandle
 //AgtPosInterface GetSupportedMinimumInterFrameGaps
 return "", nil
}

