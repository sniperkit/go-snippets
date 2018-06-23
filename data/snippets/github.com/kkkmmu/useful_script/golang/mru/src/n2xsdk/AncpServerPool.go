package n2xsdk

type AncpServerPool struct {
 Handler string
}

func(np *ncpServerPool) SetAncpStandard () error {
 //parameters: SessionPoolHandle AncpStandard
 //AgtAncpServerPool SetAncpStandard
 return nil
}

func(np *ncpServerPool) GetAncpStandard ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetAncpStandard
 return "", nil
}

func(np *ncpServerPool) SetGsmpv3Standard () error {
 //parameters: SessionPoolHandle Gsmpv3Standard
 //AgtAncpServerPool SetGsmpv3Standard
 return nil
}

func(np *ncpServerPool) GetGsmpv3Standard ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetGsmpv3Standard
 return "", nil
}

func(np *ncpServerPool) SetPartitionId () error {
 //parameters: SessionPoolHandle PartitionId Repeat Increment
 //AgtAncpServerPool SetPartitionId
 return nil
}

func(np *ncpServerPool) GetPartitionId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetPartitionId
 return "", nil
}

func(np *ncpServerPool) SetMaxConcurrentConnections () error {
 //parameters: SessionPoolHandle MaxConcurrentConnections
 //AgtAncpServerPool SetMaxConcurrentConnections
 return nil
}

func(np *ncpServerPool) GetMaxConcurrentConnections ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetMaxConcurrentConnections
 return "", nil
}

func(np *ncpServerPool) SetAdjacencyKeepAliveTime () error {
 //parameters: SessionPoolHandle KeepAliveTime
 //AgtAncpServerPool SetAdjacencyKeepAliveTime
 return nil
}

func(np *ncpServerPool) GetAdjacencyKeepAliveTime ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetAdjacencyKeepAliveTime
 return "", nil
}

func(np *ncpServerPool) SetPortManagementMessageRate () error {
 //parameters: SessionPoolHandle MessagesPerInterval Interval
 //AgtAncpServerPool SetPortManagementMessageRate
 return nil
}

func(np *ncpServerPool) GetPortManagementMessageRate ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerPool GetPortManagementMessageRate
 return "", nil
}

func(np *ncpServerPool) EnableCapability () error {
 //parameters: SessionPoolHandle AncpCapability
 //AgtAncpServerPool EnableCapability
 return nil
}

func(np *ncpServerPool) DisableCapability () error {
 //parameters: SessionPoolHandle AncpCapability
 //AgtAncpServerPool DisableCapability
 return nil
}

func(np *ncpServerPool) IsCapabilityEnabled () error {
 //parameters: SessionPoolHandle AncpCapability
 //AgtAncpServerPool IsCapabilityEnabled
 return nil
}

func(np *ncpServerPool) SendReset () error {
 //parameters: SessionIdentifiers
 //AgtAncpServerPool SendReset
 return nil
}

func(np *ncpServerPool) SendPortManagementLineConfig () error {
 //parameters: SessionIdentifiers CircuitId ProfileName
 //AgtAncpServerPool SendPortManagementLineConfig
 return nil
}

func(np *ncpServerPool) SendPortManagementOam () error {
 //parameters: SessionIdentifiers CircuitId
 //AgtAncpServerPool SendPortManagementOam
 return nil
}

func(np *ncpServerPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtAncpServerPool GetVlanPriority
 return "", nil
}

func(np *ncpServerPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtAncpServerPool SetVlanPriority
 return nil
}

func(np *ncpServerPool) IsIpv4PriorityNoCodePointFieldSelected () error {
 //parameters: DeviceHandle
 //AgtAncpServerPool IsIpv4PriorityNoCodePointFieldSelected
 return nil
}

func(np *ncpServerPool) SelectIpv4PriorityNoCodePointField () error {
 //parameters: DeviceHandle
 //AgtAncpServerPool SelectIpv4PriorityNoCodePointField
 return nil
}

func(np *ncpServerPool) GetIpv4Priority ()(string, error) {
 //parameters: DeviceHandle
 //AgtAncpServerPool GetIpv4Priority
 return "", nil
}

func(np *ncpServerPool) SetIpv4Priority () error {
 //parameters: DeviceHandle Ipv4Priority
 //AgtAncpServerPool SetIpv4Priority
 return nil
}

func(np *ncpServerPool) IsIpv4PriorityTypeOfServiceFieldSelected () error {
 //parameters: DeviceHandle
 //AgtAncpServerPool IsIpv4PriorityTypeOfServiceFieldSelected
 return nil
}

func(np *ncpServerPool) SelectIpv4PriorityTypeOfServiceField () error {
 //parameters: DeviceHandle
 //AgtAncpServerPool SelectIpv4PriorityTypeOfServiceField
 return nil
}

func(np *ncpServerPool) GetIpv4PriorityTypeOfServiceValue ()(string, error) {
 //parameters: DeviceHandle TosCodePointField
 //AgtAncpServerPool GetIpv4PriorityTypeOfServiceValue
 return "", nil
}

func(np *ncpServerPool) SetIpv4PriorityTypeOfServiceValue () error {
 //parameters: DeviceHandle TosCodePointField Value
 //AgtAncpServerPool SetIpv4PriorityTypeOfServiceValue
 return nil
}

func(np *ncpServerPool) IsIpv4PriorityDiffServFieldSelected () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtAncpServerPool IsIpv4PriorityDiffServFieldSelected
 return nil
}

func(np *ncpServerPool) SelectIpv4PriorityDiffServField () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtAncpServerPool SelectIpv4PriorityDiffServField
 return nil
}

func(np *ncpServerPool) GetIpv4PriorityDiffServValue ()(string, error) {
 //parameters: DeviceHandle DiffServCodePointField
 //AgtAncpServerPool GetIpv4PriorityDiffServValue
 return "", nil
}

func(np *ncpServerPool) SetIpv4PriorityDiffServValue () error {
 //parameters: DeviceHandle DiffServCodePointConfigurableField Value
 //AgtAncpServerPool SetIpv4PriorityDiffServValue
 return nil
}

