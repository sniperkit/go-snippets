package n2xsdk

type Statistics struct {
 Handler string
}

func(np *Statistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtStatistics GetType
 return "", nil
}

func(np *Statistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtStatistics GetName
 return "", nil
}

func(np *Statistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtStatistics GetLockCount
 return "", nil
}

func(np *Statistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts psaDestinationPortHandles
 //AgtStatistics SelectPorts
 return nil
}

func(np *Statistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedPorts
 return "", nil
}

func(np *Statistics) ListPortsWithStreamsSelected ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListPortsWithStreamsSelected
 return "", nil
}

func(np *Statistics) SelectStatisticsSet () error {
 //parameters: StatisticsHandle StatisticsSet
 //AgtStatistics SelectStatisticsSet
 return nil
}

func(np *Statistics) GetSelectedStatisticsSet ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics GetSelectedStatisticsSet
 return "", nil
}

func(np *Statistics) SelectStreamGroups () error {
 //parameters: StatisticsHandle Count psaStreamGroupHandles
 //AgtStatistics SelectStreamGroups
 return nil
}

func(np *Statistics) ListSelectedStreamGroups ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedStreamGroups
 return "", nil
}

func(np *Statistics) SelectStreams () error {
 //parameters: StatisticsHandle NumStreamGroupHandles psaStreamGroupHandles NumStreamIndexes psaStreamIndexes NumDestinationPortHandles psaDestinationPortHandles
 //AgtStatistics SelectStreams
 return nil
}

func(np *Statistics) ListSelectedStreams ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedStreams
 return "", nil
}

func(np *Statistics) ClearSelectedStreams () error {
 //parameters: StatisticsHandle
 //AgtStatistics ClearSelectedStreams
 return nil
}

func(np *Statistics) SelectConnectionGroups () error {
 //parameters: StatisticsHandle Count psaConnectionGroupHandles
 //AgtStatistics SelectConnectionGroups
 return nil
}

func(np *Statistics) SelectConnectionGroupsByPorts () error {
 //parameters: StatisticsHandle NumPortHandles psaPortHandles
 //AgtStatistics SelectConnectionGroupsByPorts
 return nil
}

func(np *Statistics) SelectConnectionGroupsByCategories () error {
 //parameters: StatisticsHandle NumCategoryClasses psaCategoryClasses NumCategoryData psaCategoryData
 //AgtStatistics SelectConnectionGroupsByCategories
 return nil
}

func(np *Statistics) ClearSelectedConnectionGroups () error {
 //parameters: StatisticsHandle
 //AgtStatistics ClearSelectedConnectionGroups
 return nil
}

func(np *Statistics) ListSelectedConnectionGroups ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedConnectionGroups
 return "", nil
}

func(np *Statistics) ListSelectedConnectionGroupPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedConnectionGroupPorts
 return "", nil
}

func(np *Statistics) ListAvailableStatisticsCategories ()(string, error) {
 //parameters: 
 //AgtStatistics ListAvailableStatisticsCategories
 return "", nil
}

func(np *Statistics) ListSelectedStatisticsCategories ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedStatisticsCategories
 return "", nil
}

func(np *Statistics) GetConnectionGroupStatistics ()(string, error) {
 //parameters: StatisticsHandle NumConnectionGroups psaConnectionGroupHandles
 //AgtStatistics GetConnectionGroupStatistics
 return "", nil
}

func(np *Statistics) GetStatistics ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics GetStatistics
 return "", nil
}

func(np *Statistics) GetStreamGroupStatistics ()(string, error) {
 //parameters: StatisticsHandle StreamGroupHandle
 //AgtStatistics GetStreamGroupStatistics
 return "", nil
}

func(np *Statistics) GetStreamStatistics ()(string, error) {
 //parameters: StatisticsHandle StreamGroupHandle StreamIndex DestinationPortHandles
 //AgtStatistics GetStreamStatistics
 return "", nil
}

func(np *Statistics) GetCurrentInterval ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics GetCurrentInterval
 return "", nil
}

func(np *Statistics) GetAvailableStreamCount ()(string, error) {
 //parameters: DestinationPortHandle
 //AgtStatistics GetAvailableStreamCount
 return "", nil
}

func(np *Statistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtStatistics ListAvailableStatistics
 return "", nil
}

func(np *Statistics) SelectStatistics () error {
 //parameters: StatisticsHandle NumStatistics psaStatistics
 //AgtStatistics SelectStatistics
 return nil
}

func(np *Statistics) ListSelectedStatistics ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedStatistics
 return "", nil
}

func(np *Statistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics GetAccumulatedValues
 return "", nil
}

func(np *Statistics) IsStatisticsSupported () error {
 //parameters: DestinationPortHandle
 //AgtStatistics IsStatisticsSupported
 return nil
}

func(np *Statistics) SelectFieldValues () error {
 //parameters: StatisticsHandle NumFieldInitialValues psaFieldInitialValues NumFieldValueCounts psaFieldValueCounts NumFieldValueIncrements psaFieldValueIncrements NumDestinationPortHandles psaDestinationPortHandles
 //AgtStatistics SelectFieldValues
 return nil
}

func(np *Statistics) ListSelectedFieldValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtStatistics ListSelectedFieldValues
 return "", nil
}

func(np *Statistics) GetAvailableStreamCountOnSession ()(string, error) {
 //parameters: 
 //AgtStatistics GetAvailableStreamCountOnSession
 return "", nil
}

