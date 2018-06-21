package n2xsdk

type DeviceList struct {
 Handler string
}

func(np *DeviceLis) Add () error {
 //parameters: PortHandle
 //AgtDeviceList Add
 return nil
}

func(np *DeviceLis) Remove () error {
 //parameters: DeviceHandle
 //AgtDeviceList Remove
 return nil
}

func(np *DeviceLis) ListHandles ()(string, error) {
 //parameters: PortHandle
 //AgtDeviceList ListHandles
 return "", nil
}

func(np *DeviceLis) ListHandlesWithEmulation ()(string, error) {
 //parameters: PortHandle EmulationType
 //AgtDeviceList ListHandlesWithEmulation
 return "", nil
}

func(np *DeviceLis) SetName () error {
 //parameters: DeviceHandle DeviceName
 //AgtDeviceList SetName
 return nil
}

func(np *DeviceLis) GetName ()(string, error) {
 //parameters: DeviceHandle
 //AgtDeviceList GetName
 return "", nil
}

func(np *DeviceLis) SetLinkDefaults () error {
 //parameters: LinkDefaults
 //AgtDeviceList SetLinkDefaults
 return nil
}

func(np *DeviceLis) GetLinkDefaults ()(string, error) {
 //parameters: 
 //AgtDeviceList GetLinkDefaults
 return "", nil
}

