package n2xsdk

type LagLacp struct {
 Handler string
}

func(np *LagLacp) SetPrimaryPort () error {
 //parameters: LagHandle PortHandle
 //AgtLagLacp SetPrimaryPort
 return nil
}

func(np *LagLacp) GetPrimaryPort ()(string, error) {
 //parameters: LagHandle
 //AgtLagLacp GetPrimaryPort
 return "", nil
}

func(np *LagLacp) SetTransmitPort () error {
 //parameters: LagHandle PortHandle
 //AgtLagLacp SetTransmitPort
 return nil
}

func(np *LagLacp) GetTransmitPort ()(string, error) {
 //parameters: LagHandle
 //AgtLagLacp GetTransmitPort
 return "", nil
}

func(np *LagLacp) SetLinkedPorts () error {
 //parameters: LagHandle PortHandleList
 //AgtLagLacp SetLinkedPorts
 return nil
}

func(np *LagLacp) ListLinkedPorts ()(string, error) {
 //parameters: LagHandle
 //AgtLagLacp ListLinkedPorts
 return "", nil
}

func(np *LagLacp) StartAll () error {
 //parameters: LagHandle
 //AgtLagLacp StartAll
 return nil
}

func(np *LagLacp) StopAll () error {
 //parameters: LagHandle
 //AgtLagLacp StopAll
 return nil
}

func(np *LagLacp) StartPort () error {
 //parameters: LagHandle PortHandle
 //AgtLagLacp StartPort
 return nil
}

func(np *LagLacp) StopPort () error {
 //parameters: LagHandle PortHandle
 //AgtLagLacp StopPort
 return nil
}

func(np *LagLacp) ClearAllStatistics () error {
 //parameters: LagHandle
 //AgtLagLacp ClearAllStatistics
 return nil
}

func(np *LagLacp) ClearPortStatistics () error {
 //parameters: LagHandle PortHandle
 //AgtLagLacp ClearPortStatistics
 return nil
}

func(np *LagLacp) GetState ()(string, error) {
 //parameters: LagHandle PortHandle
 //AgtLagLacp GetState
 return "", nil
}

func(np *LagLacp) GetRoute ()(string, error) {
 //parameters: LagHandle PortHandle
 //AgtLagLacp GetRoute
 return "", nil
}

func(np *LagLacp) GetRedirectedStatistics ()(string, error) {
 //parameters: LagHandle
 //AgtLagLacp GetRedirectedStatistics
 return "", nil
}

func(np *LagLacp) GetPortStatistics ()(string, error) {
 //parameters: LagHandle PortHandle
 //AgtLagLacp GetPortStatistics
 return "", nil
}

