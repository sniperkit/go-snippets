package n2xsdk

type LdpStatistics struct {
 Handler string
}

func(np *LdpStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtLdpStatistics ListSelectedSessions
 return "", nil
}

func(np *LdpStatistics) SelectSessions () error {
 //parameters: Count psaSessions
 //AgtLdpStatistics SelectSessions
 return nil
}

func(np *LdpStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtLdpStatistics DeselectSession
 return nil
}

func(np *LdpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpStatistics GetSessionAccumulatedValues
 return "", nil
}

