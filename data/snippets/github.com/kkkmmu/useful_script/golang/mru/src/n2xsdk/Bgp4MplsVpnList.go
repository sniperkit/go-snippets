package n2xsdk

type Bgp4MplsVpnList struct {
 Handler string
}

func(np *Bgp4MplsVpnLis) AddVpns () error {
 //parameters: TypeCount psaTypes AdminCount psaAdministrators NumberCount psaAssignedNumbers
 //AgtBgp4MplsVpnList AddVpns
 return nil
}

func(np *Bgp4MplsVpnLis) AddVpnRange () error {
 //parameters: NumVpns Type FirstAdministrator AdministratorModifier FirstAssignedNumber AssignedNumberModifier
 //AgtBgp4MplsVpnList AddVpnRange
 return nil
}

func(np *Bgp4MplsVpnLis) ListVpns ()(string, error) {
 //parameters: 
 //AgtBgp4MplsVpnList ListVpns
 return "", nil
}

func(np *Bgp4MplsVpnLis) RemoveVpns () error {
 //parameters: Count psaVpnHandles
 //AgtBgp4MplsVpnList RemoveVpns
 return nil
}

func(np *Bgp4MplsVpnLis) ClearVpns () error {
 //parameters: 
 //AgtBgp4MplsVpnList ClearVpns
 return nil
}

