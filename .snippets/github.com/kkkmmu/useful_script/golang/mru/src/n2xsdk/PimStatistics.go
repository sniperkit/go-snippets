package n2xsdk

type PimStatistics struct {
 Handler string
}

func(np *PimStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtPimStatistics ListSelectedSessions
 return "", nil
}

func(np *PimStatistics) SelectSessions () error {
 //parameters: Count psaSessionHandles
 //AgtPimStatistics SelectSessions
 return nil
}

func(np *PimStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtPimStatistics DeselectSession
 return nil
}

func(np *PimStatistics) GetSessionAccumulatedValue ()(string, error) {
 //parameters: SessionHandle StatisticsType
 //AgtPimStatistics GetSessionAccumulatedValue
 return "", nil
}

func(np *PimStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtPimStatistics GetSessionAccumulatedValues
 return "", nil
}

