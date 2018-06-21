package n2xsdk

type RsvpTunnelEgressList struct {
 Handler string
}

func(np *RsvpTunnelEgressLis) EnableEgress () error {
 //parameters: SessionHandle
 //AgtRsvpTunnelEgressList EnableEgress
 return nil
}

func(np *RsvpTunnelEgressLis) EnableAllEgress () error {
 //parameters: 
 //AgtRsvpTunnelEgressList EnableAllEgress
 return nil
}

func(np *RsvpTunnelEgressLis) DisableEgress () error {
 //parameters: SessionHandle
 //AgtRsvpTunnelEgressList DisableEgress
 return nil
}

func(np *RsvpTunnelEgressLis) DisableAllEgress () error {
 //parameters: 
 //AgtRsvpTunnelEgressList DisableAllEgress
 return nil
}

func(np *RsvpTunnelEgressLis) IsEgressEnabled () error {
 //parameters: SessionHandle
 //AgtRsvpTunnelEgressList IsEgressEnabled
 return nil
}

func(np *RsvpTunnelEgressLis) ListTunnels ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpTunnelEgressList ListTunnels
 return "", nil
}

func(np *RsvpTunnelEgressLis) GetEnableFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtRsvpTunnelEgressList GetEnableFlag
 return "", nil
}

