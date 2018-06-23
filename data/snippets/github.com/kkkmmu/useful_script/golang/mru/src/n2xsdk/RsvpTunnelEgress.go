package n2xsdk

type RsvpTunnelEgress struct {
 Handler string
}

func(np *RsvpTunnelEgress) CloseInstance () error {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelEgress CloseInstance
 return nil
}

func(np *RsvpTunnelEgress) CloseTunnel () error {
 //parameters: TunnelHandle
 //AgtRsvpTunnelEgress CloseTunnel
 return nil
}

func(np *RsvpTunnelEgress) ListEgressLsps ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelEgress ListEgressLsps
 return "", nil
}

func(np *RsvpTunnelEgress) GetDestinationIpAddress ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelEgress GetDestinationIpAddress
 return "", nil
}

func(np *RsvpTunnelEgress) GetSourceIpAddress ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelEgress GetSourceIpAddress
 return "", nil
}

func(np *RsvpTunnelEgress) GetTunnelId ()(string, error) {
 //parameters: TunnelHandle
 //AgtRsvpTunnelEgress GetTunnelId
 return "", nil
}

func(np *RsvpTunnelEgress) GetState ()(string, error) {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelEgress GetState
 return "", nil
}

func(np *RsvpTunnelEgress) GetAssignedLabel ()(string, error) {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelEgress GetAssignedLabel
 return "", nil
}

func(np *RsvpTunnelEgress) GetAssignedLabelInformation ()(string, error) {
 //parameters: TunnelHandle LspId
 //AgtRsvpTunnelEgress GetAssignedLabelInformation
 return "", nil
}

