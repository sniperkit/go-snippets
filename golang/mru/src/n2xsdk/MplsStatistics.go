package n2xsdk

type MplsStatistics struct {
 Handler string
}

func(np *MplsStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtMplsStatistics GetType
 return "", nil
}

func(np *MplsStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtMplsStatistics GetName
 return "", nil
}

func(np *MplsStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtMplsStatistics GetLockCount
 return "", nil
}

func(np *MplsStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtMplsStatistics ListSelectedPorts
 return "", nil
}

func(np *MplsStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtMplsStatistics SelectPorts
 return nil
}

func(np *MplsStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtMplsStatistics GetAccumulatedValues
 return "", nil
}

func(np *MplsStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtMplsStatistics LockItem
 return nil
}

func(np *MplsStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtMplsStatistics UnlockItem
 return nil
}

func(np *MplsStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtMplsStatistics ListAvailableStatistics
 return "", nil
}

func(np *MplsStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtMplsStatistics ListSelectedStatistics
 return "", nil
}

func(np *MplsStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtMplsStatistics SelectStatistics
 return nil
}

