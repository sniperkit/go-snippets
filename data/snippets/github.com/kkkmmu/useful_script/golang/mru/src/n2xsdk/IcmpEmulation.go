package n2xsdk

type IcmpEmulation struct {
 Handler string
}

func(np *IcmpEmulation) Enable () error {
 //parameters: PortHandle
 //AgtIcmpEmulation Enable
 return nil
}

func(np *IcmpEmulation) Disable () error {
 //parameters: PortHandle
 //AgtIcmpEmulation Disable
 return nil
}

func(np *IcmpEmulation) IsEnabled () error {
 //parameters: PortHandle
 //AgtIcmpEmulation IsEnabled
 return nil
}

func(np *IcmpEmulation) IsIcmpSupported () error {
 //parameters: PortHandle
 //AgtIcmpEmulation IsIcmpSupported
 return nil
}

