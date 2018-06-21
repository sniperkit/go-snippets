package n2xsdk

type MultiplayStatisticsControl struct {
 Handler string
}

func(np *MultiplayStatisticsControl) Enable () error {
 //parameters: NumPorts pPortHandles
 //AgtMultiplayStatisticsControl Enable
 return nil
}

func(np *MultiplayStatisticsControl) Disable () error {
 //parameters: NumPorts pPortHandles
 //AgtMultiplayStatisticsControl Disable
 return nil
}

func(np *MultiplayStatisticsControl) IsEnabled () error {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl IsEnabled
 return nil
}

func(np *MultiplayStatisticsControl) IsMultiplayStatisticsSupported () error {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl IsMultiplayStatisticsSupported
 return nil
}

func(np *MultiplayStatisticsControl) EnableDualMode () error {
 //parameters: NumPorts pPortHandles
 //AgtMultiplayStatisticsControl EnableDualMode
 return nil
}

func(np *MultiplayStatisticsControl) DisableDualMode () error {
 //parameters: NumPorts pPortHandles
 //AgtMultiplayStatisticsControl DisableDualMode
 return nil
}

func(np *MultiplayStatisticsControl) IsDualModeEnabled () error {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl IsDualModeEnabled
 return nil
}

func(np *MultiplayStatisticsControl) SetFieldConfiguration () error {
 //parameters: NumPorts pPortHandles FieldConfig
 //AgtMultiplayStatisticsControl SetFieldConfiguration
 return nil
}

func(np *MultiplayStatisticsControl) GetFieldConfiguration ()(string, error) {
 //parameters: NumPorts pPortHandles
 //AgtMultiplayStatisticsControl GetFieldConfiguration
 return "", nil
}

func(np *MultiplayStatisticsControl) ListAvailableFields ()(string, error) {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl ListAvailableFields
 return "", nil
}

func(np *MultiplayStatisticsControl) GetMaximumNumberOfGroups ()(string, error) {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl GetMaximumNumberOfGroups
 return "", nil
}

func(np *MultiplayStatisticsControl) SetGroupBaseAddress () error {
 //parameters: PortHandle AddressField GroupAddress
 //AgtMultiplayStatisticsControl SetGroupBaseAddress
 return nil
}

func(np *MultiplayStatisticsControl) GetGroupAddressRange ()(string, error) {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl GetGroupAddressRange
 return "", nil
}

func(np *MultiplayStatisticsControl) SetInitialVlanValue () error {
 //parameters: PortHandle VlanField VlanId
 //AgtMultiplayStatisticsControl SetInitialVlanValue
 return nil
}

func(np *MultiplayStatisticsControl) GetVlanRange ()(string, error) {
 //parameters: PortHandle
 //AgtMultiplayStatisticsControl GetVlanRange
 return "", nil
}

