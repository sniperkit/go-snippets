package n2xsdk

type PingManager struct {
 Handler string
}

func(np *PingManager) SetSourcePort () error {
 //parameters: PortHandle
 //AgtPingManager SetSourcePort
 return nil
}

func(np *PingManager) GetSourcePort ()(string, error) {
 //parameters: 
 //AgtPingManager GetSourcePort
 return "", nil
}

func(np *PingManager) SetDestinationPorts () error {
 //parameters: Count psaPortHandles
 //AgtPingManager SetDestinationPorts
 return nil
}

func(np *PingManager) ListDestinationPorts ()(string, error) {
 //parameters: 
 //AgtPingManager ListDestinationPorts
 return "", nil
}

func(np *PingManager) SetDestinationAddresses () error {
 //parameters: Count psaIpAddresses
 //AgtPingManager SetDestinationAddresses
 return nil
}

func(np *PingManager) ListDestinationAddresses ()(string, error) {
 //parameters: 
 //AgtPingManager ListDestinationAddresses
 return "", nil
}

func(np *PingManager) SetParameter () error {
 //parameters: Parameter Value
 //AgtPingManager SetParameter
 return nil
}

func(np *PingManager) GetParameter ()(string, error) {
 //parameters: Parameter
 //AgtPingManager GetParameter
 return "", nil
}

func(np *PingManager) SetPingMode () error {
 //parameters: Mode
 //AgtPingManager SetPingMode
 return nil
}

func(np *PingManager) StartEchoRequests () error {
 //parameters: 
 //AgtPingManager StartEchoRequests
 return nil
}

func(np *PingManager) StopEchoRequests () error {
 //parameters: 
 //AgtPingManager StopEchoRequests
 return nil
}

func(np *PingManager) IsRunning () error {
 //parameters: 
 //AgtPingManager IsRunning
 return nil
}

func(np *PingManager) GetResults ()(string, error) {
 //parameters: 
 //AgtPingManager GetResults
 return "", nil
}

