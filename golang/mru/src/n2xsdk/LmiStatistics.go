package n2xsdk

type LmiStatistics struct {
 Handler string
}

func(np *LmiStatistics) SelectPorts () error {
 //parameters: Count psaPorts
 //AgtLmiStatistics SelectPorts
 return nil
}

func(np *LmiStatistics) ListSelectedPorts ()(string, error) {
 //parameters: 
 //AgtLmiStatistics ListSelectedPorts
 return "", nil
}

func(np *LmiStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: PortHandle
 //AgtLmiStatistics GetAccumulatedValues
 return "", nil
}

func(np *LmiStatistics) ClearPortStatistics () error {
 //parameters: PortHandle
 //AgtLmiStatistics ClearPortStatistics
 return nil
}

