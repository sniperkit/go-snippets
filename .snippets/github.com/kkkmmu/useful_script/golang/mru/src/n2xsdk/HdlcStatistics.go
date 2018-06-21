package n2xsdk

type HdlcStatistics struct {
 Handler string
}

func(np *HdlcStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtHdlcStatistics GetType
 return "", nil
}

func(np *HdlcStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtHdlcStatistics GetName
 return "", nil
}

func(np *HdlcStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtHdlcStatistics GetLockCount
 return "", nil
}

func(np *HdlcStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtHdlcStatistics ListSelectedPorts
 return "", nil
}

func(np *HdlcStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtHdlcStatistics SelectPorts
 return nil
}

func(np *HdlcStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtHdlcStatistics GetAccumulatedValues
 return "", nil
}

func(np *HdlcStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtHdlcStatistics LockItem
 return nil
}

func(np *HdlcStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtHdlcStatistics UnlockItem
 return nil
}

func(np *HdlcStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtHdlcStatistics ListAvailableStatistics
 return "", nil
}

func(np *HdlcStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtHdlcStatistics ListSelectedStatistics
 return "", nil
}

func(np *HdlcStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtHdlcStatistics SelectStatistics
 return nil
}

