package n2xsdk

type EthernetOamLoopback struct {
 Handler string
}

func(np *EthernetOamLoopback) SetStandard () error {
 //parameters: SessionPoolHandle EthOamStandard
 //AgtEthernetOamLoopback SetStandard
 return nil
}

func(np *EthernetOamLoopback) GetStandard ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetStandard
 return "", nil
}

func(np *EthernetOamLoopback) SetDestinationAddressType () error {
 //parameters: SessionPoolHandle DestinationAddressType
 //AgtEthernetOamLoopback SetDestinationAddressType
 return nil
}

func(np *EthernetOamLoopback) GetDestinationAddressType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetDestinationAddressType
 return "", nil
}

func(np *EthernetOamLoopback) SetDestinationMacAddress () error {
 //parameters: SessionPoolHandle MacAddress
 //AgtEthernetOamLoopback SetDestinationMacAddress
 return nil
}

func(np *EthernetOamLoopback) GetDestinationMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetDestinationMacAddress
 return "", nil
}

func(np *EthernetOamLoopback) SetMdLevel () error {
 //parameters: SessionPoolHandle Level
 //AgtEthernetOamLoopback SetMdLevel
 return nil
}

func(np *EthernetOamLoopback) GetMdLevel ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetMdLevel
 return "", nil
}

func(np *EthernetOamLoopback) SetMegLevel () error {
 //parameters: SessionPoolHandle Level
 //AgtEthernetOamLoopback SetMegLevel
 return nil
}

func(np *EthernetOamLoopback) GetMegLevel ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetMegLevel
 return "", nil
}

func(np *EthernetOamLoopback) SetTransactionIdentifier () error {
 //parameters: SessionPoolHandle TransId
 //AgtEthernetOamLoopback SetTransactionIdentifier
 return nil
}

func(np *EthernetOamLoopback) GetTransactionIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetTransactionIdentifier
 return "", nil
}

func(np *EthernetOamLoopback) EnableSenderIdTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback EnableSenderIdTlv
 return nil
}

func(np *EthernetOamLoopback) DisableSenderIdTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback DisableSenderIdTlv
 return nil
}

func(np *EthernetOamLoopback) IsSenderIdTlvEnabled () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback IsSenderIdTlvEnabled
 return nil
}

func(np *EthernetOamLoopback) SetSenderIdTlv () error {
 //parameters: SessionPoolHandle ChassisIdLength ChassisIdSubType ChassisID MgmtAddressDomain MgmtAddressDomainLength MgmtAddress MgmtAddresslength
 //AgtEthernetOamLoopback SetSenderIdTlv
 return nil
}

func(np *EthernetOamLoopback) GetSenderIdTlv ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetSenderIdTlv
 return "", nil
}

func(np *EthernetOamLoopback) EnableOrganizationSpecificTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback EnableOrganizationSpecificTlv
 return nil
}

func(np *EthernetOamLoopback) DisableOrganizationSpecificTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback DisableOrganizationSpecificTlv
 return nil
}

func(np *EthernetOamLoopback) IsOrganizationSpecificTlvEnabled () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback IsOrganizationSpecificTlvEnabled
 return nil
}

func(np *EthernetOamLoopback) SetOrganizationSpecificTlv () error {
 //parameters: SessionPoolHandle OrgSpecTlvLen OrgUniqId OrgSpecSubType OrgSpecTlvValue
 //AgtEthernetOamLoopback SetOrganizationSpecificTlv
 return nil
}

func(np *EthernetOamLoopback) GetOrganizationSpecificTlv ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetOrganizationSpecificTlv
 return "", nil
}

func(np *EthernetOamLoopback) EnableDataTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback EnableDataTlv
 return nil
}

func(np *EthernetOamLoopback) DisableDataTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback DisableDataTlv
 return nil
}

func(np *EthernetOamLoopback) IsDataTlvEnabled () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback IsDataTlvEnabled
 return nil
}

func(np *EthernetOamLoopback) SetDataTlv () error {
 //parameters: SessionPoolHandle DataTlvLength DataTlvValue
 //AgtEthernetOamLoopback SetDataTlv
 return nil
}

func(np *EthernetOamLoopback) GetDataTlv ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetDataTlv
 return "", nil
}

func(np *EthernetOamLoopback) EnableTestTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback EnableTestTlv
 return nil
}

func(np *EthernetOamLoopback) DisableTestTlv () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback DisableTestTlv
 return nil
}

func(np *EthernetOamLoopback) IsTestTlvEnabled () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback IsTestTlvEnabled
 return nil
}

func(np *EthernetOamLoopback) SetTestTlv () error {
 //parameters: SessionPoolHandle TestTlvLength TestTlvPatternType
 //AgtEthernetOamLoopback SetTestTlv
 return nil
}

func(np *EthernetOamLoopback) GetTestTlv ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetTestTlv
 return "", nil
}

func(np *EthernetOamLoopback) SetBurstSize () error {
 //parameters: SessionPoolHandle BurstSize
 //AgtEthernetOamLoopback SetBurstSize
 return nil
}

func(np *EthernetOamLoopback) GetBurstSize ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetBurstSize
 return "", nil
}

func(np *EthernetOamLoopback) SetBurstDelay () error {
 //parameters: SessionPoolHandle BurstDelay
 //AgtEthernetOamLoopback SetBurstDelay
 return nil
}

func(np *EthernetOamLoopback) GetBurstDelay ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetBurstDelay
 return "", nil
}

func(np *EthernetOamLoopback) EnableRenew () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback EnableRenew
 return nil
}

func(np *EthernetOamLoopback) DisableRenew () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback DisableRenew
 return nil
}

func(np *EthernetOamLoopback) IsRenewEnabled () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback IsRenewEnabled
 return nil
}

func(np *EthernetOamLoopback) SetRenewTimeout () error {
 //parameters: SessionPoolHandle RenewTimeout
 //AgtEthernetOamLoopback SetRenewTimeout
 return nil
}

func(np *EthernetOamLoopback) GetRenewTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetRenewTimeout
 return "", nil
}

func(np *EthernetOamLoopback) SetResponseTimeout () error {
 //parameters: SessionPoolHandle ResponseTimeout
 //AgtEthernetOamLoopback SetResponseTimeout
 return nil
}

func(np *EthernetOamLoopback) GetResponseTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback GetResponseTimeout
 return "", nil
}

func(np *EthernetOamLoopback) Start () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback Start
 return nil
}

func(np *EthernetOamLoopback) ForceRenew () error {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopback ForceRenew
 return nil
}

func(np *EthernetOamLoopback) SetIncrementingParameter () error {
 //parameters: SessionPoolHandle IncrementingParameter Value Repeat Increment
 //AgtEthernetOamLoopback SetIncrementingParameter
 return nil
}

func(np *EthernetOamLoopback) GetIncrementingParameter ()(string, error) {
 //parameters: SessionPoolHandle IncrementingParameter
 //AgtEthernetOamLoopback GetIncrementingParameter
 return "", nil
}

func(np *EthernetOamLoopback) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtEthernetOamLoopback GetVlanPriority
 return "", nil
}

func(np *EthernetOamLoopback) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtEthernetOamLoopback SetVlanPriority
 return nil
}

