package n2xsdk

type StatisticsLog struct {
 Handler string
}

func(np *StatisticsLo) EnableLogging () error {
 //parameters: 
 //AgtStatisticsLog EnableLogging
 return nil
}

func(np *StatisticsLo) DisableLogging () error {
 //parameters: 
 //AgtStatisticsLog DisableLogging
 return nil
}

func(np *StatisticsLo) IsLoggingEnabled () error {
 //parameters: 
 //AgtStatisticsLog IsLoggingEnabled
 return nil
}

func(np *StatisticsLo) SetLogFile () error {
 //parameters: LogFile
 //AgtStatisticsLog SetLogFile
 return nil
}

func(np *StatisticsLo) GetLogFile ()(string, error) {
 //parameters: 
 //AgtStatisticsLog GetLogFile
 return "", nil
}

func(np *StatisticsLo) SetLoggingInterval () error {
 //parameters: Multiple
 //AgtStatisticsLog SetLoggingInterval
 return nil
}

func(np *StatisticsLo) GetLoggingInterval ()(string, error) {
 //parameters: 
 //AgtStatisticsLog GetLoggingInterval
 return "", nil
}

func(np *StatisticsLo) SelectPorts () error {
 //parameters: NumPorts psaPorts
 //AgtStatisticsLog SelectPorts
 return nil
}

func(np *StatisticsLo) ListSelectedPorts ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListSelectedPorts
 return "", nil
}

func(np *StatisticsLo) ClearSelectedStreams () error {
 //parameters: 
 //AgtStatisticsLog ClearSelectedStreams
 return nil
}

func(np *StatisticsLo) ListPortsWithStreamsSelected ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListPortsWithStreamsSelected
 return "", nil
}

func(np *StatisticsLo) SelectStreamGroups () error {
 //parameters: Count psaStreamGroupHandles
 //AgtStatisticsLog SelectStreamGroups
 return nil
}

func(np *StatisticsLo) ListSelectedStreamGroups ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListSelectedStreamGroups
 return "", nil
}

func(np *StatisticsLo) SelectStreams () error {
 //parameters: NumStreamGroupHandles psaStreamGroupHandles NumStreamIndexes psaStreamIndexes NumDestinationPortHandles psaDestinationPorts
 //AgtStatisticsLog SelectStreams
 return nil
}

func(np *StatisticsLo) ListSelectedStreams ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListSelectedStreams
 return "", nil
}

func(np *StatisticsLo) SelectFieldValues () error {
 //parameters: NumFieldInitialValues psaFieldInitialValues NumFieldValueCounts psaFieldValueCounts NumFieldValueIncrements psaFieldValueIncrements NumDestinationPortHandles psaDestinationPortHandles
 //AgtStatisticsLog SelectFieldValues
 return nil
}

func(np *StatisticsLo) ListSelectedFieldValues ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListSelectedFieldValues
 return "", nil
}

func(np *StatisticsLo) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListAvailableStatistics
 return "", nil
}

func(np *StatisticsLo) SelectStatistics () error {
 //parameters: NumStatistics psaStatistics
 //AgtStatisticsLog SelectStatistics
 return nil
}

func(np *StatisticsLo) ListSelectedStatistics ()(string, error) {
 //parameters: 
 //AgtStatisticsLog ListSelectedStatistics
 return "", nil
}

