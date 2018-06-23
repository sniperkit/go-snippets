package n2xsdk

type StatisticsControl struct {
 Handler string
}

func(np *StatisticsControl) SetOctetCountMode () error {
 //parameters: OctetCountMode
 //AgtStatisticsControl SetOctetCountMode
 return nil
}

func(np *StatisticsControl) GetOctetCountMode ()(string, error) {
 //parameters: 
 //AgtStatisticsControl GetOctetCountMode
 return "", nil
}

func(np *StatisticsControl) SetExpectedPayloadProtection () error {
 //parameters: PortHandle PayloadProtectionType CustomByteOffset
 //AgtStatisticsControl SetExpectedPayloadProtection
 return nil
}

func(np *StatisticsControl) GetExpectedPayloadProtection ()(string, error) {
 //parameters: PortHandle
 //AgtStatisticsControl GetExpectedPayloadProtection
 return "", nil
}

func(np *StatisticsControl) SetServiceDisruptionGuardTime () error {
 //parameters: PortHandle GuardTime
 //AgtStatisticsControl SetServiceDisruptionGuardTime
 return nil
}

func(np *StatisticsControl) GetServiceDisruptionGuardTime ()(string, error) {
 //parameters: PortHandle
 //AgtStatisticsControl GetServiceDisruptionGuardTime
 return "", nil
}

func(np *StatisticsControl) SetErroredFrameFilter () error {
 //parameters: PortHandle ErroredFrameFilter
 //AgtStatisticsControl SetErroredFrameFilter
 return nil
}

func(np *StatisticsControl) GetErroredFrameFilter ()(string, error) {
 //parameters: PortHandle
 //AgtStatisticsControl GetErroredFrameFilter
 return "", nil
}

func(np *StatisticsControl) GetStreamStatisticsMode ()(string, error) {
 //parameters: 
 //AgtStatisticsControl GetStreamStatisticsMode
 return "", nil
}

func(np *StatisticsControl) ListFrameMatchersForStatistics ()(string, error) {
 //parameters: FrameMatcherStatistics NumPorts psaPortHandles
 //AgtStatisticsControl ListFrameMatchersForStatistics
 return "", nil
}

func(np *StatisticsControl) ListStatisticsForFrameMatcher ()(string, error) {
 //parameters: FrameMatcherHandle
 //AgtStatisticsControl ListStatisticsForFrameMatcher
 return "", nil
}

func(np *StatisticsControl) GetMaxStreamsOnPort ()(string, error) {
 //parameters: PortHandle
 //AgtStatisticsControl GetMaxStreamsOnPort
 return "", nil
}

func(np *StatisticsControl) GetMaxStreamsPerSession ()(string, error) {
 //parameters: 
 //AgtStatisticsControl GetMaxStreamsPerSession
 return "", nil
}

func(np *StatisticsControl) GetMaxStatisticsSlots ()(string, error) {
 //parameters: 
 //AgtStatisticsControl GetMaxStatisticsSlots
 return "", nil
}

func(np *StatisticsControl) GetUsedStatisticsSlots ()(string, error) {
 //parameters: 
 //AgtStatisticsControl GetUsedStatisticsSlots
 return "", nil
}

func(np *StatisticsControl) GetNumSlotsForStatistics ()(string, error) {
 //parameters: NumStatistics psaStatistics
 //AgtStatisticsControl GetNumSlotsForStatistics
 return "", nil
}

func(np *StatisticsControl) SetIndexType () error {
 //parameters: NumPorts pPortHandles IndexType
 //AgtStatisticsControl SetIndexType
 return nil
}

func(np *StatisticsControl) GetIndexType ()(string, error) {
 //parameters: NumPorts pPortHandles
 //AgtStatisticsControl GetIndexType
 return "", nil
}

func(np *StatisticsControl) ListSupportedIndexTypes ()(string, error) {
 //parameters: PortHandle
 //AgtStatisticsControl ListSupportedIndexTypes
 return "", nil
}

