package n2xsdk

type Emulation struct {
 Handler string
}

func(np *Emulation) SetGlobalStatisticsUpdateInterval () error {
 //parameters: UpdateInterval
 //AgtEmulation SetGlobalStatisticsUpdateInterval
 return nil
}

func(np *Emulation) GetGlobalStatisticsUpdateInterval ()(string, error) {
 //parameters: 
 //AgtEmulation GetGlobalStatisticsUpdateInterval
 return "", nil
}

func(np *Emulation) Enable () error {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation Enable
 return nil
}

func(np *Emulation) Disable () error {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation Disable
 return nil
}

func(np *Emulation) DisableErrored () error {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation DisableErrored
 return nil
}

func(np *Emulation) GetStateCount ()(string, error) {
 //parameters: DeviceHandle EmulationType State
 //AgtEmulation GetStateCount
 return "", nil
}

func(np *Emulation) GetStateCounts ()(string, error) {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation GetStateCounts
 return "", nil
}

func(np *Emulation) ListCustomStates ()(string, error) {
 //parameters: EmulationType
 //AgtEmulation ListCustomStates
 return "", nil
}

func(np *Emulation) GetCustomStateCount ()(string, error) {
 //parameters: DeviceHandle EmulationType CustomState
 //AgtEmulation GetCustomStateCount
 return "", nil
}

func(np *Emulation) GetPortLimit ()(string, error) {
 //parameters: PortHandle EmulationType
 //AgtEmulation GetPortLimit
 return "", nil
}

func(np *Emulation) ResetStatistics () error {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation ResetStatistics
 return nil
}

func(np *Emulation) Reset () error {
 //parameters: DeviceHandle EmulationType
 //AgtEmulation Reset
 return nil
}

func(np *Emulation) ResetInCustomState () error {
 //parameters: DeviceHandle EmulationType CustomState
 //AgtEmulation ResetInCustomState
 return nil
}

