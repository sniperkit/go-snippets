package n2xsdk

type LmiEmulation struct {
 Handler string
}

func(np *LmiEmulation) Enable () error {
 //parameters: PortHandle
 //AgtLmiEmulation Enable
 return nil
}

func(np *LmiEmulation) Disable () error {
 //parameters: PortHandle
 //AgtLmiEmulation Disable
 return nil
}

func(np *LmiEmulation) IsEnabled () error {
 //parameters: PortHandle
 //AgtLmiEmulation IsEnabled
 return nil
}

func(np *LmiEmulation) SetVersion () error {
 //parameters: PortHandle Version
 //AgtLmiEmulation SetVersion
 return nil
}

func(np *LmiEmulation) GetVersion ()(string, error) {
 //parameters: PortHandle
 //AgtLmiEmulation GetVersion
 return "", nil
}

func(np *LmiEmulation) SetEmulationMode () error {
 //parameters: PortHandle Mode
 //AgtLmiEmulation SetEmulationMode
 return nil
}

func(np *LmiEmulation) GetEmulationMode ()(string, error) {
 //parameters: PortHandle
 //AgtLmiEmulation GetEmulationMode
 return "", nil
}

