package n2xsdk

type LinkOamDte struct {
 Handler string
}

func(np *LinkOamDte) EnableAutomaticDiscovery () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte EnableAutomaticDiscovery
 return nil
}

func(np *LinkOamDte) DisableAutomaticDiscovery () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte DisableAutomaticDiscovery
 return nil
}

func(np *LinkOamDte) IsAutomaticDiscoveryEnabled () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte IsAutomaticDiscoveryEnabled
 return nil
}

func(np *LinkOamDte) EnableOrganizationSpecificInformation () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte EnableOrganizationSpecificInformation
 return nil
}

func(np *LinkOamDte) DisableOrganizationSpecificInformation () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte DisableOrganizationSpecificInformation
 return nil
}

func(np *LinkOamDte) IsOrganizationSpecificInformationEnabled () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte IsOrganizationSpecificInformationEnabled
 return nil
}

func(np *LinkOamDte) EnableCriticalLinkEvents () error {
 //parameters: SessionPoolHandle CriticalLinkEventType
 //AgtLinkOamDte EnableCriticalLinkEvents
 return nil
}

func(np *LinkOamDte) DisableCriticalLinkEvents () error {
 //parameters: SessionPoolHandle CriticalLinkEventType
 //AgtLinkOamDte DisableCriticalLinkEvents
 return nil
}

func(np *LinkOamDte) IsCriticalLinkEventsEnabled () error {
 //parameters: SessionPoolHandle CriticalLinkEventType
 //AgtLinkOamDte IsCriticalLinkEventsEnabled
 return nil
}

func(np *LinkOamDte) EnableLinkEventNotification () error {
 //parameters: SessionPoolHandle LinkNotificationType
 //AgtLinkOamDte EnableLinkEventNotification
 return nil
}

func(np *LinkOamDte) DisableLinkEventNotification () error {
 //parameters: SessionPoolHandle LinkNotificationType
 //AgtLinkOamDte DisableLinkEventNotification
 return nil
}

func(np *LinkOamDte) IsLinkEventNotificationEnabled () error {
 //parameters: SessionPoolHandle LinkNotificationType
 //AgtLinkOamDte IsLinkEventNotificationEnabled
 return nil
}

func(np *LinkOamDte) EnableOrganizationSpecificEvent () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte EnableOrganizationSpecificEvent
 return nil
}

func(np *LinkOamDte) DisableOrganizationSpecificEvent () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte DisableOrganizationSpecificEvent
 return nil
}

func(np *LinkOamDte) IsOrganizationSpecificEventEnabled () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte IsOrganizationSpecificEventEnabled
 return nil
}

func(np *LinkOamDte) EnableDiscoveryCapability () error {
 //parameters: SessionPoolHandle DiscoveryCapabilityType
 //AgtLinkOamDte EnableDiscoveryCapability
 return nil
}

func(np *LinkOamDte) DisableDiscoveryCapability () error {
 //parameters: SessionPoolHandle DiscoveryCapabilityType
 //AgtLinkOamDte DisableDiscoveryCapability
 return nil
}

func(np *LinkOamDte) IsDiscoveryCapabilityEnabled () error {
 //parameters: SessionPoolHandle DiscoveryCapabilityType
 //AgtLinkOamDte IsDiscoveryCapabilityEnabled
 return nil
}

func(np *LinkOamDte) SetOamMode () error {
 //parameters: SessionPoolHandle OamMode
 //AgtLinkOamDte SetOamMode
 return nil
}

func(np *LinkOamDte) GetOamMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetOamMode
 return "", nil
}

func(np *LinkOamDte) SetDestinationMacType () error {
 //parameters: SessionPoolHandle DestinationMacType
 //AgtLinkOamDte SetDestinationMacType
 return nil
}

func(np *LinkOamDte) GetDestinationMacType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetDestinationMacType
 return "", nil
}

func(np *LinkOamDte) SetUnicastMac () error {
 //parameters: SessionPoolHandle UnicastMac
 //AgtLinkOamDte SetUnicastMac
 return nil
}

func(np *LinkOamDte) GetUnicastMac ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetUnicastMac
 return "", nil
}

func(np *LinkOamDte) GetOui ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetOui
 return "", nil
}

func(np *LinkOamDte) SetOui () error {
 //parameters: SessionPoolHandle Oui
 //AgtLinkOamDte SetOui
 return nil
}

func(np *LinkOamDte) GetVendorSpecificInformation ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetVendorSpecificInformation
 return "", nil
}

func(np *LinkOamDte) SetVendorSpecificInformation () error {
 //parameters: SessionPoolHandle VendorSpecificInformation
 //AgtLinkOamDte SetVendorSpecificInformation
 return nil
}

func(np *LinkOamDte) GetMaxPduLength ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetMaxPduLength
 return "", nil
}

func(np *LinkOamDte) SetMaxPduLength () error {
 //parameters: SessionPoolHandle MaxOamPduLength
 //AgtLinkOamDte SetMaxPduLength
 return nil
}

func(np *LinkOamDte) GetNumberOfPdus ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetNumberOfPdus
 return "", nil
}

func(np *LinkOamDte) SetNumberOfPdus () error {
 //parameters: SessionPoolHandle PduCount
 //AgtLinkOamDte SetNumberOfPdus
 return nil
}

func(np *LinkOamDte) SetLinkEventNotificationFlowControl () error {
 //parameters: SessionPoolHandle EventsPerInterval Interval NumberOfIntervals
 //AgtLinkOamDte SetLinkEventNotificationFlowControl
 return nil
}

func(np *LinkOamDte) GetLinkEventNotificationFlowControl ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetLinkEventNotificationFlowControl
 return "", nil
}

func(np *LinkOamDte) GetOrganizationSpecificEventValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetOrganizationSpecificEventValue
 return "", nil
}

func(np *LinkOamDte) SetOrganizationSpecificEventValue () error {
 //parameters: SessionPoolHandle OrganizationSpecificEventValue
 //AgtLinkOamDte SetOrganizationSpecificEventValue
 return nil
}

func(np *LinkOamDte) GetOrganizationSpecificInformationValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte GetOrganizationSpecificInformationValue
 return "", nil
}

func(np *LinkOamDte) SetOrganizationSpecificInformationValue () error {
 //parameters: SessionPoolHandle OrganizationSpecificInformationValue
 //AgtLinkOamDte SetOrganizationSpecificInformationValue
 return nil
}

func(np *LinkOamDte) StartDiscovery () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte StartDiscovery
 return nil
}

func(np *LinkOamDte) SendEnableRemoteLoopback () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte SendEnableRemoteLoopback
 return nil
}

func(np *LinkOamDte) SendDisableRemoteLoopback () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte SendDisableRemoteLoopback
 return nil
}

func(np *LinkOamDte) StartLinkEventNotification () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte StartLinkEventNotification
 return nil
}

func(np *LinkOamDte) StopLinkEventNotification () error {
 //parameters: SessionPoolHandle
 //AgtLinkOamDte StopLinkEventNotification
 return nil
}

