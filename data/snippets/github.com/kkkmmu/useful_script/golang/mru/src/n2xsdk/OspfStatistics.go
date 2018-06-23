package n2xsdk

type OspfStatistics struct {
 Handler string
}

func(np *OspfStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtOspfStatistics ListSelectedSessions
 return "", nil
}

func(np *OspfStatistics) ListSelectedOspfv2Sessions ()(string, error) {
 //parameters: 
 //AgtOspfStatistics ListSelectedOspfv2Sessions
 return "", nil
}

func(np *OspfStatistics) ListSelectedOspfv3Sessions ()(string, error) {
 //parameters: 
 //AgtOspfStatistics ListSelectedOspfv3Sessions
 return "", nil
}

func(np *OspfStatistics) SelectSessions () error {
 //parameters: Count psaSessions
 //AgtOspfStatistics SelectSessions
 return nil
}

func(np *OspfStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtOspfStatistics DeselectSession
 return nil
}

func(np *OspfStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionHandle
 //AgtOspfStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *OspfStatistics) GetOspfv2SessionAccumulatedValue ()(string, error) {
 //parameters: SessionHandle StatisticsType
 //AgtOspfStatistics GetOspfv2SessionAccumulatedValue
 return "", nil
}

func(np *OspfStatistics) GetOspfv3SessionAccumulatedValue ()(string, error) {
 //parameters: SessionHandle StatisticsTypeV3
 //AgtOspfStatistics GetOspfv3SessionAccumulatedValue
 return "", nil
}

func(np *OspfStatistics) ClearSessionStatistics () error {
 //parameters: SessionHandle
 //AgtOspfStatistics ClearSessionStatistics
 return nil
}

