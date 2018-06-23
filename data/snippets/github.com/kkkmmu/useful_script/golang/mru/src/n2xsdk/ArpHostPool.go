package n2xsdk

type ArpHostPool struct {
 Handler string
}

func(np *rpHostPool) EnableAutoSendRequests () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool EnableAutoSendRequests
 return nil
}

func(np *rpHostPool) DisableAutoSendRequests () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool DisableAutoSendRequests
 return nil
}

func(np *rpHostPool) IsAutoSendRequestsEnabled () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool IsAutoSendRequestsEnabled
 return nil
}

func(np *rpHostPool) GetCustomState ()(string, error) {
 //parameters: DevicePoolInstance
 //AgtArpHostPool GetCustomState
 return "", nil
}

func(np *rpHostPool) SendRequests () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool SendRequests
 return nil
}

func(np *rpHostPool) Reset () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool Reset
 return nil
}

func(np *rpHostPool) CancelRequests () error {
 //parameters: DevicePoolHandle
 //AgtArpHostPool CancelRequests
 return nil
}

func(np *rpHostPool) GetRemoteMacAddress ()(string, error) {
 //parameters: DevicePoolInstance
 //AgtArpHostPool GetRemoteMacAddress
 return "", nil
}

func(np *rpHostPool) SendRequestsByPort () error {
 //parameters: PortHandle
 //AgtArpHostPool SendRequestsByPort
 return nil
}

func(np *rpHostPool) ResetByPort () error {
 //parameters: PortHandle
 //AgtArpHostPool ResetByPort
 return nil
}

func(np *rpHostPool) CancelRequestsByPort () error {
 //parameters: PortHandle
 //AgtArpHostPool CancelRequestsByPort
 return nil
}

func(np *rpHostPool) SendAllRequests () error {
 //parameters: 
 //AgtArpHostPool SendAllRequests
 return nil
}

func(np *rpHostPool) ResetAll () error {
 //parameters: 
 //AgtArpHostPool ResetAll
 return nil
}

func(np *rpHostPool) CancelAllRequests () error {
 //parameters: 
 //AgtArpHostPool CancelAllRequests
 return nil
}

func(np *rpHostPool) SetRateControl () error {
 //parameters: PortHandle RateControlMode
 //AgtArpHostPool SetRateControl
 return nil
}

func(np *rpHostPool) GetRateControl ()(string, error) {
 //parameters: PortHandle
 //AgtArpHostPool GetRateControl
 return "", nil
}

func(np *rpHostPool) SetFixedRateProfile () error {
 //parameters: PortHandle FixedRate MaxBurstSize
 //AgtArpHostPool SetFixedRateProfile
 return nil
}

func(np *rpHostPool) GetFixedRateProfile ()(string, error) {
 //parameters: PortHandle
 //AgtArpHostPool GetFixedRateProfile
 return "", nil
}

func(np *rpHostPool) SetSelfPacedProfile () error {
 //parameters: PortHandle MaxPackets MaxRetries RetryTimeout
 //AgtArpHostPool SetSelfPacedProfile
 return nil
}

func(np *rpHostPool) GetSelfPacedProfile ()(string, error) {
 //parameters: PortHandle
 //AgtArpHostPool GetSelfPacedProfile
 return "", nil
}

func(np *rpHostPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtArpHostPool GetVlanPriority
 return "", nil
}

func(np *rpHostPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtArpHostPool SetVlanPriority
 return nil
}

