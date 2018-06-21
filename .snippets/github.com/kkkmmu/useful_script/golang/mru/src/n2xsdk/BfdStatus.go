package n2xsdk

type BfdStatus struct {
 Handler string
}

func(np *BfdStatus) GetStatus ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle Status
 //AgtBfdStatus GetStatus
 return "", nil
}

func(np *BfdStatus) GetRemoteDiscriminator ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetRemoteDiscriminator
 return "", nil
}

func(np *BfdStatus) GetLastDownTimestamp ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastDownTimestamp
 return "", nil
}

func(np *BfdStatus) GetLastUpTimestamp ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastUpTimestamp
 return "", nil
}

func(np *BfdStatus) GetLastRxDiagnostic ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastRxDiagnostic
 return "", nil
}

func(np *BfdStatus) GetLastTxDiagnostic ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastTxDiagnostic
 return "", nil
}

func(np *BfdStatus) GetLastReceivedMultiplier ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastReceivedMultiplier
 return "", nil
}

func(np *BfdStatus) GetLastReceivedMinTxInterval ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastReceivedMinTxInterval
 return "", nil
}

func(np *BfdStatus) GetLastReceivedMinRxInterval ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastReceivedMinRxInterval
 return "", nil
}

func(np *BfdStatus) GetLastReceivedMinEchoRxInterval ()(string, error) {
 //parameters: SessionPoolHandle InstanceHandle
 //AgtBfdStatus GetLastReceivedMinEchoRxInterval
 return "", nil
}

