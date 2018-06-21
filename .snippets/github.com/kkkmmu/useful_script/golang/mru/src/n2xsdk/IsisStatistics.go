package n2xsdk

type IsisStatistics struct {
 Handler string
}

func(np *IsisStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtIsisStatistics ListSelectedSessions
 return "", nil
}

func(np *IsisStatistics) SelectSessions () error {
 //parameters: Count psaSessions
 //AgtIsisStatistics SelectSessions
 return nil
}

func(np *IsisStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtIsisStatistics DeselectSession
 return nil
}

func(np *IsisStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisStatistics GetSessionAccumulatedValues
 return "", nil
}

