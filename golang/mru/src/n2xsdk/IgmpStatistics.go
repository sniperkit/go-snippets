package n2xsdk

type IgmpStatistics struct {
 Handler string
}

func(np *IgmpStatistics) GetDeviceAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle StatisticsList
 //AgtIgmpStatistics GetDeviceAccumulatedValues
 return "", nil
}

func(np *IgmpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionList
 //AgtIgmpStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *IgmpStatistics) GetSessionAccumulatedSpecifiedValues ()(string, error) {
 //parameters: SessionList StatisticsList
 //AgtIgmpStatistics GetSessionAccumulatedSpecifiedValues
 return "", nil
}

func(np *IgmpStatistics) GetSessionAccumulatedValue ()(string, error) {
 //parameters: SessionList StatisticsList
 //AgtIgmpStatistics GetSessionAccumulatedValue
 return "", nil
}

func(np *IgmpStatistics) ClearStatistics () error {
 //parameters: SessionList
 //AgtIgmpStatistics ClearStatistics
 return nil
}

func(np *IgmpStatistics) ClearSessionStatistics () error {
 //parameters: SessionList
 //AgtIgmpStatistics ClearSessionStatistics
 return nil
}

func(np *IgmpStatistics) ListSelectedSessions ()(string, error) {
 //parameters: 
 //AgtIgmpStatistics ListSelectedSessions
 return "", nil
}

func(np *IgmpStatistics) SelectSessions () error {
 //parameters: Count psaSessionHandles
 //AgtIgmpStatistics SelectSessions
 return nil
}

func(np *IgmpStatistics) DeselectSession () error {
 //parameters: SessionHandle
 //AgtIgmpStatistics DeselectSession
 return nil
}

func(np *IgmpStatistics) GetSessionAccumulatedVersionValues ()(string, error) {
 //parameters: SessionHandle IgmpVersion
 //AgtIgmpStatistics GetSessionAccumulatedVersionValues
 return "", nil
}

