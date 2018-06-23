package n2xsdk

type FastFindEndOfTestControl struct {
 Handler string
}

func(np *FastFindEndOfTestControl) RetrieveData () error {
 //parameters: 
 //AgtFastFindEndOfTestControl RetrieveData
 return nil
}

func(np *FastFindEndOfTestControl) CancelDataRetrieval () error {
 //parameters: 
 //AgtFastFindEndOfTestControl CancelDataRetrieval
 return nil
}

func(np *FastFindEndOfTestControl) GetDataRetrievalStatus ()(string, error) {
 //parameters: 
 //AgtFastFindEndOfTestControl GetDataRetrievalStatus
 return "", nil
}

func(np *FastFindEndOfTestControl) EnableDataRetrievalOnAllPorts () error {
 //parameters: EnableOnAllPorts
 //AgtFastFindEndOfTestControl EnableDataRetrievalOnAllPorts
 return nil
}

func(np *FastFindEndOfTestControl) IsDataRetrievalSetOnAllPortsEnabled () error {
 //parameters: 
 //AgtFastFindEndOfTestControl IsDataRetrievalSetOnAllPortsEnabled
 return nil
}

func(np *FastFindEndOfTestControl) SelectPortsForDataRetrieval () error {
 //parameters: PortHandles
 //AgtFastFindEndOfTestControl SelectPortsForDataRetrieval
 return nil
}

func(np *FastFindEndOfTestControl) ListSelectedPortsForDataRetrieval ()(string, error) {
 //parameters: 
 //AgtFastFindEndOfTestControl ListSelectedPortsForDataRetrieval
 return "", nil
}

func(np *FastFindEndOfTestControl) IsAutoRetrieveDataEnabled () error {
 //parameters: 
 //AgtFastFindEndOfTestControl IsAutoRetrieveDataEnabled
 return nil
}

func(np *FastFindEndOfTestControl) EnableAutoRetrieveData () error {
 //parameters: 
 //AgtFastFindEndOfTestControl EnableAutoRetrieveData
 return nil
}

func(np *FastFindEndOfTestControl) DisableAutoRetrieveData () error {
 //parameters: 
 //AgtFastFindEndOfTestControl DisableAutoRetrieveData
 return nil
}

