package n2xsdk

type NdpAddressAutoconfigHostStatistics struct {
 Handler string
}

func(np *NdpAddressAutoconfigHostStatistics) GetStatistics ()(string, error) {
 //parameters: DeviceIdentifiers StatisticsTypes
 //AgtNdpAddressAutoconfigHostStatistics GetStatistics
 return "", nil
}

func(np *NdpAddressAutoconfigHostStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtNdpAddressAutoconfigHostStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *NdpAddressAutoconfigHostStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostStatistics GetStatisticsUpdateCount
 return "", nil
}

func(np *NdpAddressAutoconfigHostStatistics) ClearStatistics () error {
 //parameters: DevicePoolHandle
 //AgtNdpAddressAutoconfigHostStatistics ClearStatistics
 return nil
}

