package n2xsdk

type VsrStatistics struct {
 Handler string
}

func(np *VsrStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtVsrStatistics GetType
 return "", nil
}

func(np *VsrStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtVsrStatistics GetName
 return "", nil
}

func(np *VsrStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtVsrStatistics GetLockCount
 return "", nil
}

func(np *VsrStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtVsrStatistics ListSelectedPorts
 return "", nil
}

func(np *VsrStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtVsrStatistics SelectPorts
 return nil
}

func(np *VsrStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtVsrStatistics GetAccumulatedValues
 return "", nil
}

func(np *VsrStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtVsrStatistics LockItem
 return nil
}

func(np *VsrStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtVsrStatistics UnlockItem
 return nil
}

func(np *VsrStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtVsrStatistics ListAvailableStatistics
 return "", nil
}

func(np *VsrStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtVsrStatistics ListSelectedStatistics
 return "", nil
}

func(np *VsrStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtVsrStatistics SelectStatistics
 return nil
}

