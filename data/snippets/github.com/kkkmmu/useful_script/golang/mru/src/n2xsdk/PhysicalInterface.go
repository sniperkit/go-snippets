package n2xsdk

type PhysicalInterface struct {
 Handler string
}

func(np *PhysicalInterface) ListAvailablePortModes ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface ListAvailablePortModes
 return "", nil
}

func(np *PhysicalInterface) ListAvailablePortModesByMediaType ()(string, error) {
 //parameters: PortHandle MediaType
 //AgtPhysicalInterface ListAvailablePortModesByMediaType
 return "", nil
}

func(np *PhysicalInterface) SetPortMode () error {
 //parameters: PortHandle PortMode
 //AgtPhysicalInterface SetPortMode
 return nil
}

func(np *PhysicalInterface) GetPortMode ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetPortMode
 return "", nil
}

func(np *PhysicalInterface) ListAvailableClockSources ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface ListAvailableClockSources
 return "", nil
}

func(np *PhysicalInterface) SetClockSource () error {
 //parameters: PortHandle ClockSource
 //AgtPhysicalInterface SetClockSource
 return nil
}

func(np *PhysicalInterface) GetClockSource ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetClockSource
 return "", nil
}

func(np *PhysicalInterface) GetClockState ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetClockState
 return "", nil
}

func(np *PhysicalInterface) IsPluginMediaSupported () error {
 //parameters: PortHandle
 //AgtPhysicalInterface IsPluginMediaSupported
 return nil
}

func(np *PhysicalInterface) IsPluginMediaInserted () error {
 //parameters: PortHandle
 //AgtPhysicalInterface IsPluginMediaInserted
 return nil
}

func(np *PhysicalInterface) GetPluginMediaType ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetPluginMediaType
 return "", nil
}

func(np *PhysicalInterface) ListAvailableMediaTypes ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface ListAvailableMediaTypes
 return "", nil
}

func(np *PhysicalInterface) SetMediaType () error {
 //parameters: PortHandle MediaType
 //AgtPhysicalInterface SetMediaType
 return nil
}

func(np *PhysicalInterface) GetMediaType ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetMediaType
 return "", nil
}

func(np *PhysicalInterface) IsMediaTypeSelectionSupported () error {
 //parameters: PortHandle MediaType
 //AgtPhysicalInterface IsMediaTypeSelectionSupported
 return nil
}

func(np *PhysicalInterface) ListAvailableExternalClockRates ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface ListAvailableExternalClockRates
 return "", nil
}

func(np *PhysicalInterface) SetExternalClockRate () error {
 //parameters: PortHandle Rate
 //AgtPhysicalInterface SetExternalClockRate
 return nil
}

func(np *PhysicalInterface) GetExternalClockRate ()(string, error) {
 //parameters: PortHandle
 //AgtPhysicalInterface GetExternalClockRate
 return "", nil
}

