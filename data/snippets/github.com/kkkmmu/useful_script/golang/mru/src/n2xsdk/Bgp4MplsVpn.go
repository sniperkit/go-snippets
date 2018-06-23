package n2xsdk

type Bgp4MplsVpn struct {
 Handler string
}

func(np *Bgp4MplsVpn) SetRouteTarget () error {
 //parameters: hVpn Type Administrator AssignedNumber
 //AgtBgp4MplsVpn SetRouteTarget
 return nil
}

func(np *Bgp4MplsVpn) GetRouteTarget ()(string, error) {
 //parameters: hVpn
 //AgtBgp4MplsVpn GetRouteTarget
 return "", nil
}

func(np *Bgp4MplsVpn) SetName () error {
 //parameters: hVpn Name
 //AgtBgp4MplsVpn SetName
 return nil
}

func(np *Bgp4MplsVpn) GetName ()(string, error) {
 //parameters: hVpn
 //AgtBgp4MplsVpn GetName
 return "", nil
}

