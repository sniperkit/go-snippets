package n2xsdk

type NdpHostStatistics struct {
 Handler string
}

func(np *NdpHostStatistics) GetStatistics ()(string, error) {
 //parameters: DeviceIdentifiers StatisticsTypes
 //AgtNdpHostStatistics GetStatistics
 return "", nil
}

func(np *NdpHostStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtNdpHostStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *NdpHostStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: DevicePoolHandle
 //AgtNdpHostStatistics GetStatisticsUpdateCount
 return "", nil
}

func(np *NdpHostStatistics) ClearStatistics () error {
 //parameters: DevicePoolHandle
 //AgtNdpHostStatistics ClearStatistics
 return nil
}

