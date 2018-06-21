package n2xsdk

type Device struct {
 Handler string
}

func(np *Device) Enable () error {
 //parameters: DeviceHandle
 //AgtDevice Enable
 return nil
}

func(np *Device) Disable () error {
 //parameters: DeviceHandle
 //AgtDevice Disable
 return nil
}

func(np *Device) SetPoolSize () error {
 //parameters: DeviceHandle PoolSize
 //AgtDevice SetPoolSize
 return nil
}

func(np *Device) GetPoolSize ()(string, error) {
 //parameters: DeviceHandle
 //AgtDevice GetPoolSize
 return "", nil
}

func(np *Device) ResetStatistics () error {
 //parameters: DeviceHandle
 //AgtDevice ResetStatistics
 return nil
}

func(np *Device) Reset () error {
 //parameters: DeviceHandle
 //AgtDevice Reset
 return nil
}

func(np *Device) GetPort ()(string, error) {
 //parameters: DeviceHandle
 //AgtDevice GetPort
 return "", nil
}

func(np *Device) ListStreamGroups ()(string, error) {
 //parameters: DeviceHandle
 //AgtDevice ListStreamGroups
 return "", nil
}

func(np *Device) RemoveAllStreamGroups () error {
 //parameters: DeviceHandle
 //AgtDevice RemoveAllStreamGroups
 return nil
}

func(np *Device) ListMeshes ()(string, error) {
 //parameters: DeviceHandle
 //AgtDevice ListMeshes
 return "", nil
}

func(np *Device) RemoveAllMeshes () error {
 //parameters: DeviceHandle
 //AgtDevice RemoveAllMeshes
 return nil
}

func(np *Device) AddEmulation () error {
 //parameters: DeviceHandle EmulationType
 //AgtDevice AddEmulation
 return nil
}

func(np *Device) RemoveEmulation () error {
 //parameters: DeviceHandle EmulationType
 //AgtDevice RemoveEmulation
 return nil
}

func(np *Device) RemoveAllEmulations () error {
 //parameters: DeviceHandle
 //AgtDevice RemoveAllEmulations
 return nil
}

func(np *Device) ListEmulations ()(string, error) {
 //parameters: DeviceHandle
 //AgtDevice ListEmulations
 return "", nil
}

func(np *Device) ListAvailableEmulations ()(string, error) {
 //parameters: 
 //AgtDevice ListAvailableEmulations
 return "", nil
}

