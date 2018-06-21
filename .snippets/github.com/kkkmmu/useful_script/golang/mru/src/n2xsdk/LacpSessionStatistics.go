package n2xsdk

type LacpSessionStatistics struct {
 Handler string
}

func(np *LacpSessionStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSessionStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *LacpSessionStatistics) GetSessionConvergenceTimes ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSessionStatistics GetSessionConvergenceTimes
 return "", nil
}

func(np *LacpSessionStatistics) ClearSessionStatistics () error {
 //parameters: SessionHandle
 //AgtLacpSessionStatistics ClearSessionStatistics
 return nil
}

