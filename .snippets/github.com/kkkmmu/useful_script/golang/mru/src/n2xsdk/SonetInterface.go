package n2xsdk

type SonetInterface struct {
 Handler string
}

func(np *SonetInterface) ScramblerOn () error {
 //parameters: PortHandle
 //AgtSonetInterface ScramblerOn
 return nil
}

func(np *SonetInterface) ScramblerOff () error {
 //parameters: PortHandle
 //AgtSonetInterface ScramblerOff
 return nil
}

func(np *SonetInterface) IsScramblerOn () error {
 //parameters: PortHandle
 //AgtSonetInterface IsScramblerOn
 return nil
}

func(np *SonetInterface) DescramblerOn () error {
 //parameters: PortHandle
 //AgtSonetInterface DescramblerOn
 return nil
}

func(np *SonetInterface) DescramblerOff () error {
 //parameters: PortHandle
 //AgtSonetInterface DescramblerOff
 return nil
}

func(np *SonetInterface) IsDescramblerOn () error {
 //parameters: PortHandle
 //AgtSonetInterface IsDescramblerOn
 return nil
}

func(np *SonetInterface) SetFramingMode () error {
 //parameters: PortHandle FramingMode
 //AgtSonetInterface SetFramingMode
 return nil
}

func(np *SonetInterface) GetFramingMode ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface GetFramingMode
 return "", nil
}

func(np *SonetInterface) IsSonetSupported () error {
 //parameters: PortHandle
 //AgtSonetInterface IsSonetSupported
 return nil
}

func(np *SonetInterface) SetRate () error {
 //parameters: PortHandle Rate
 //AgtSonetInterface SetRate
 return nil
}

func(np *SonetInterface) GetRate ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface GetRate
 return "", nil
}

func(np *SonetInterface) ListAvailableRates ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface ListAvailableRates
 return "", nil
}

func(np *SonetInterface) GetClockOffsetRange ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface GetClockOffsetRange
 return "", nil
}

func(np *SonetInterface) GetClockOffsetStep ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface GetClockOffsetStep
 return "", nil
}

func(np *SonetInterface) SetClockOffset () error {
 //parameters: PortHandle ClockOffset
 //AgtSonetInterface SetClockOffset
 return nil
}

func(np *SonetInterface) GetClockOffset ()(string, error) {
 //parameters: PortHandle
 //AgtSonetInterface GetClockOffset
 return "", nil
}

