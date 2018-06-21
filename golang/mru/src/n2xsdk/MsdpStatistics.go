package n2xsdk

type MsdpStatistics struct {
 Handler string
}

func(np *MsdpStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtMsdpStatistics ListSelectedSessions
 return "", nil
}

func(np *MsdpStatistics) SelectSessions () error {
 //parameters: Count psaSessionHandles
 //AgtMsdpStatistics SelectSessions
 return nil
}

func(np *MsdpStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtMsdpStatistics DeselectSession
 return nil
}

func(np *MsdpStatistics) GetSessionAccumulatedValue ()(string, error) {
 //parameters: SessionHandle StatisticsType
 //AgtMsdpStatistics GetSessionAccumulatedValue
 return "", nil
}

func(np *MsdpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtMsdpStatistics GetSessionAccumulatedValues
 return "", nil
}

