package n2xsdk

type Ipv6Interface struct {
 Handler string
}

func(np *Ipv6Interface) IsAvailable () error {
 //parameters: hPort
 //AgtIpv6Interface IsAvailable
 return nil
}

func(np *Ipv6Interface) ListAvailablePorts ()(string, error) {
 //parameters: 
 //AgtIpv6Interface ListAvailablePorts
 return "", nil
}

