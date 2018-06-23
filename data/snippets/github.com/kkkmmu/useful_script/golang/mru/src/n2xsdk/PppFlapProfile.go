package n2xsdk

type PppFlapProfile struct {
 Handler string
}

func(np *PppFlapProfile) SetPortParameters () error {
 //parameters: PortHandle FlapOrder FlapCycles FlapDelay
 //AgtPppFlapProfile SetPortParameters
 return nil
}

func(np *PppFlapProfile) GetPortParameters ()(string, error) {
 //parameters: PortHandle
 //AgtPppFlapProfile GetPortParameters
 return "", nil
}

func(np *PppFlapProfile) SetSessionPoolParameters () error {
 //parameters: SessionPoolHandle FlapCycles FlapDelay
 //AgtPppFlapProfile SetSessionPoolParameters
 return nil
}

func(np *PppFlapProfile) GetSessionPoolParameters ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppFlapProfile GetSessionPoolParameters
 return "", nil
}

func(np *PppFlapProfile) SetObjects () error {
 //parameters: FlapObjectType Count psaObjectHandles
 //AgtPppFlapProfile SetObjects
 return nil
}

func(np *PppFlapProfile) GetObjects ()(string, error) {
 //parameters: 
 //AgtPppFlapProfile GetObjects
 return "", nil
}

func(np *PppFlapProfile) StartFlap () error {
 //parameters: 
 //AgtPppFlapProfile StartFlap
 return nil
}

func(np *PppFlapProfile) StopFlap () error {
 //parameters: 
 //AgtPppFlapProfile StopFlap
 return nil
}

func(np *PppFlapProfile) GetPortFlapState ()(string, error) {
 //parameters: PortHandle
 //AgtPppFlapProfile GetPortFlapState
 return "", nil
}

func(np *PppFlapProfile) GetSessionPoolFlapState ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppFlapProfile GetSessionPoolFlapState
 return "", nil
}

