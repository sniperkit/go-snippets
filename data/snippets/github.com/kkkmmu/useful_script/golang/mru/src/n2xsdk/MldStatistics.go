package n2xsdk

type MldStatistics struct {
 Handler string
}

func(np *MldStatistics) GetDeviceAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle StatisticsList
 //AgtMldStatistics GetDeviceAccumulatedValues
 return "", nil
}

func(np *MldStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionList
 //AgtMldStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *MldStatistics) GetSessionAccumulatedSpecifiedValues ()(string, error) {
 //parameters: SessionList StatisticsList
 //AgtMldStatistics GetSessionAccumulatedSpecifiedValues
 return "", nil
}

func(np *MldStatistics) ClearStatistics () error {
 //parameters: SessionList
 //AgtMldStatistics ClearStatistics
 return nil
}

