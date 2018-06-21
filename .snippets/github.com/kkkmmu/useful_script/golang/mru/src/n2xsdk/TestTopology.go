package n2xsdk

type TestTopology struct {
 Handler string
}

func(np *TestTopology) EnableSession () error {
 //parameters: SessionHandle
 //AgtTestTopology EnableSession
 return nil
}

func(np *TestTopology) DisableSession () error {
 //parameters: SessionHandle
 //AgtTestTopology DisableSession
 return nil
}

func(np *TestTopology) ResetSession () error {
 //parameters: SessionHandle
 //AgtTestTopology ResetSession
 return nil
}

func(np *TestTopology) ListSessionTypes ()(string, error) {
 //parameters: 
 //AgtTestTopology ListSessionTypes
 return "", nil
}

func(np *TestTopology) AddSession () error {
 //parameters: HandleLen pInterfaceHandle SessionType
 //AgtTestTopology AddSession
 return nil
}

func(np *TestTopology) RemoveSession () error {
 //parameters: SessionHandle
 //AgtTestTopology RemoveSession
 return nil
}

func(np *TestTopology) IsRoutingEngineRequired () error {
 //parameters: SessionType
 //AgtTestTopology IsRoutingEngineRequired
 return nil
}

func(np *TestTopology) SetSessionPoolSize () error {
 //parameters: SessionHandle PoolSize
 //AgtTestTopology SetSessionPoolSize
 return nil
}

func(np *TestTopology) GetSessionPoolSize ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology GetSessionPoolSize
 return "", nil
}

func(np *TestTopology) SetSessionName () error {
 //parameters: SessionHandle SessionName
 //AgtTestTopology SetSessionName
 return nil
}

func(np *TestTopology) GetSessionName ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology GetSessionName
 return "", nil
}

func(np *TestTopology) ListSessionEmulations ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology ListSessionEmulations
 return "", nil
}

func(np *TestTopology) GetEmulationState ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology GetEmulationState
 return "", nil
}

func(np *TestTopology) GetEmulationStateCount ()(string, error) {
 //parameters: SessionHandle DeviceState
 //AgtTestTopology GetEmulationStateCount
 return "", nil
}

func(np *TestTopology) ListEmulationCustomStates ()(string, error) {
 //parameters: EmulationName
 //AgtTestTopology ListEmulationCustomStates
 return "", nil
}

func(np *TestTopology) GetCustomStateCount ()(string, error) {
 //parameters: SessionHandle EmulationName CustomState
 //AgtTestTopology GetCustomStateCount
 return "", nil
}

func(np *TestTopology) GetSessionLimit ()(string, error) {
 //parameters: PortHandle SessionType
 //AgtTestTopology GetSessionLimit
 return "", nil
}

func(np *TestTopology) ListSessions ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology ListSessions
 return "", nil
}

func(np *TestTopology) ListSessionsByType ()(string, error) {
 //parameters: HandleLen pInterfaceHandle SessionType
 //AgtTestTopology ListSessionsByType
 return "", nil
}

func(np *TestTopology) GetSessionType ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology GetSessionType
 return "", nil
}

func(np *TestTopology) ListSutIpAddresses ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology ListSutIpAddresses
 return "", nil
}

func(np *TestTopology) GetSutIpAddress ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology GetSutIpAddress
 return "", nil
}

func(np *TestTopology) GetTesterIpAddress ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology GetTesterIpAddress
 return "", nil
}

func(np *TestTopology) ListSutIpv6Addresses ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology ListSutIpv6Addresses
 return "", nil
}

func(np *TestTopology) ListTesterIpv6Addresses ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology ListTesterIpv6Addresses
 return "", nil
}

func(np *TestTopology) GetSutIpv6Address ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology GetSutIpv6Address
 return "", nil
}

func(np *TestTopology) GetTesterIpv6Address ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology GetTesterIpv6Address
 return "", nil
}

func(np *TestTopology) GetLinkLocalIpv6AddressForGlobalIpv6Address ()(string, error) {
 //parameters: PortHandle aGlobalIpv6Address
 //AgtTestTopology GetLinkLocalIpv6AddressForGlobalIpv6Address
 return "", nil
}

func(np *TestTopology) GetSessionInterface ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopology GetSessionInterface
 return "", nil
}

func(np *TestTopology) SetSessionInterface () error {
 //parameters: SessionHandle HandleLen pInterfaceHandle
 //AgtTestTopology SetSessionInterface
 return nil
}

func(np *TestTopology) GetSessionPoolInterfaces ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtTestTopology GetSessionPoolInterfaces
 return "", nil
}

func(np *TestTopology) SetSessionPoolInterfaces () error {
 //parameters: SessionPoolHandle HandleLen pInterfaceHandleSet
 //AgtTestTopology SetSessionPoolInterfaces
 return nil
}

func(np *TestTopology) IsIpv4VirtualInterfaceSupported () error {
 //parameters: SessionType
 //AgtTestTopology IsIpv4VirtualInterfaceSupported
 return nil
}

func(np *TestTopology) IsIpv6VirtualInterfaceSupported () error {
 //parameters: SessionType
 //AgtTestTopology IsIpv6VirtualInterfaceSupported
 return nil
}

func(np *TestTopology) ListAllVirtualInterfacesOnPort ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology ListAllVirtualInterfacesOnPort
 return "", nil
}

func(np *TestTopology) ListIpv4VirtualInterfacesOnPort ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology ListIpv4VirtualInterfacesOnPort
 return "", nil
}

func(np *TestTopology) ListIpv6VirtualInterfacesOnPort ()(string, error) {
 //parameters: PortHandle
 //AgtTestTopology ListIpv6VirtualInterfacesOnPort
 return "", nil
}

func(np *TestTopology) GetInterfaceDescription ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology GetInterfaceDescription
 return "", nil
}

func(np *TestTopology) GetVirtualInterfaceDescription ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology GetVirtualInterfaceDescription
 return "", nil
}

func(np *TestTopology) GetNumberOfInterfacesInPool ()(string, error) {
 //parameters: HandleLen pInterfaceHandle
 //AgtTestTopology GetNumberOfInterfacesInPool
 return "", nil
}

func(np *TestTopology) GetVirtualInterfaceSutIpAddress ()(string, error) {
 //parameters: HandleLen psaInterfaceHandle
 //AgtTestTopology GetVirtualInterfaceSutIpAddress
 return "", nil
}

func(np *TestTopology) GetVirtualInterfaceSutIpv6Address ()(string, error) {
 //parameters: HandleLen psaInterfaceHandle
 //AgtTestTopology GetVirtualInterfaceSutIpv6Address
 return "", nil
}

func(np *TestTopology) ResetSessionStatistics () error {
 //parameters: SessionHandle
 //AgtTestTopology ResetSessionStatistics
 return nil
}

func(np *TestTopology) ResetEmulationStatistics () error {
 //parameters: SessionHandle EmulationName CustomState
 //AgtTestTopology ResetEmulationStatistics
 return nil
}

func(np *TestTopology) ResetEmulation () error {
 //parameters: SessionHandle EmulationName CustomState
 //AgtTestTopology ResetEmulation
 return nil
}

