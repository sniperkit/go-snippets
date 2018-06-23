package n2xsdk

type PtpMasterIpv4ClockPool struct {
 Handler string
}

func(np *PtpMasterIpv4ClockPool) SetClockQuality () error {
 //parameters: SessionPoolHandle ClockClass ClockAccuracy OffsetScaledLogVariance Priority1 Priority2 AlternateMasterFlag StepsRemoved TimeSource
 //AgtPtpMasterIpv4ClockPool SetClockQuality
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetClockQuality ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetClockQuality
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnableAnnounce
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisableAnnounce
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledAnnounce
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetAnnounceAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval ReceiptTimeout
 //AgtPtpMasterIpv4ClockPool SetAnnounceAttributes
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetAnnounceAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetAnnounceAttributes
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnableSync
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisableSync
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledSync
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetSyncAttributes () error {
 //parameters: SessionPoolHandle TwoStepFlagMode LogMessageInterval
 //AgtPtpMasterIpv4ClockPool SetSyncAttributes
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetSyncAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetSyncAttributes
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnableFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisableFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetFollowUpAttributes () error {
 //parameters: SessionPoolHandle FollowUpDelay
 //AgtPtpMasterIpv4ClockPool SetFollowUpAttributes
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetFollowUpAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetFollowUpAttributes
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnableDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisableDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetDelayResponseAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval
 //AgtPtpMasterIpv4ClockPool SetDelayResponseAttributes
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetDelayResponseAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDelayResponseAttributes
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetDomainNumber () error {
 //parameters: SessionPoolHandle DomainNumber
 //AgtPtpMasterIpv4ClockPool SetDomainNumber
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetDomainNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDomainNumber
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpMasterIpv4ClockPool SetPortNumber
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPortNumber
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpMasterIpv4ClockPool SetClockId
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetClockId
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetClockIdIncrementingRange () error {
 //parameters: SessionPoolHandle FirstClockId Repeat Increment
 //AgtPtpMasterIpv4ClockPool SetClockIdIncrementingRange
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetClockIdIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetClockIdIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle ClockIdList
 //AgtPtpMasterIpv4ClockPool SetClockIdList
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) GetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetClockIdList
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetPathDelayMode () error {
 //parameters: SessionPoolHandle PathDelayMode
 //AgtPtpMasterIpv4ClockPool SetPathDelayMode
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPathDelayMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPathDelayMode
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnablePeerDelayRequest
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisablePeerDelayRequest
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledPeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledPeerDelayRequest
 return nil
}

func(np *PtpMasterIpv4ClockPool) EnablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnablePeerDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisablePeerDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledPeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledPeerDelayResponse
 return nil
}

func(np *PtpMasterIpv4ClockPool) EnablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledPeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledPeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetPeerToPeerPathDelayAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval TwoStepFlagMode FollowUpDelay
 //AgtPtpMasterIpv4ClockPool SetPeerToPeerPathDelayAttributes
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPeerToPeerPathDelayAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPeerToPeerPathDelayAttributes
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetTimestampCalibration () error {
 //parameters: SessionPoolHandle TxOffset RxOffset
 //AgtPtpMasterIpv4ClockPool SetTimestampCalibration
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetTimestampCalibration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetTimestampCalibration
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtPtpMasterIpv4ClockPool GetVlanPriority
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtPtpMasterIpv4ClockPool SetVlanPriority
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetDestinationAddressType () error {
 //parameters: SessionPoolHandle AddressType
 //AgtPtpMasterIpv4ClockPool SetDestinationAddressType
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetDestinationAddressType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDestinationAddressType
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) IsIpv4PriorityNoCodePointFieldSelected () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsIpv4PriorityNoCodePointFieldSelected
 return nil
}

func(np *PtpMasterIpv4ClockPool) SelectIpv4PriorityNoCodePointField () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool SelectIpv4PriorityNoCodePointField
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetIpv4Priority ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetIpv4Priority
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetIpv4Priority () error {
 //parameters: SessionPoolHandle Ipv4Priority
 //AgtPtpMasterIpv4ClockPool SetIpv4Priority
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsIpv4PriorityTypeOfServiceFieldSelected () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsIpv4PriorityTypeOfServiceFieldSelected
 return nil
}

func(np *PtpMasterIpv4ClockPool) SelectIpv4PriorityTypeOfServiceField () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool SelectIpv4PriorityTypeOfServiceField
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetIpv4PriorityTypeOfServiceValue ()(string, error) {
 //parameters: SessionPoolHandle TosCodePointField
 //AgtPtpMasterIpv4ClockPool GetIpv4PriorityTypeOfServiceValue
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetIpv4PriorityTypeOfServiceValue () error {
 //parameters: SessionPoolHandle TosCodePointField Value
 //AgtPtpMasterIpv4ClockPool SetIpv4PriorityTypeOfServiceValue
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsIpv4PriorityDiffServFieldSelected () error {
 //parameters: SessionPoolHandle DiffServPerHopBehaviorGroup
 //AgtPtpMasterIpv4ClockPool IsIpv4PriorityDiffServFieldSelected
 return nil
}

func(np *PtpMasterIpv4ClockPool) SelectIpv4PriorityDiffServField () error {
 //parameters: SessionPoolHandle DiffServPerHopBehaviorGroup
 //AgtPtpMasterIpv4ClockPool SelectIpv4PriorityDiffServField
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetIpv4PriorityDiffServValue ()(string, error) {
 //parameters: SessionPoolHandle DiffServCodePointField
 //AgtPtpMasterIpv4ClockPool GetIpv4PriorityDiffServValue
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetIpv4PriorityDiffServValue () error {
 //parameters: SessionPoolHandle DiffServCodePointConfigurableField Value
 //AgtPtpMasterIpv4ClockPool SetIpv4PriorityDiffServValue
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetDestinationAddress () error {
 //parameters: SessionPoolHandle FirstAddress
 //AgtPtpMasterIpv4ClockPool SetDestinationAddress
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetDestinationAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDestinationAddress
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetDestinationAddressIncrementingRange () error {
 //parameters: SessionPoolHandle FirstAddress PrefixLength Increment Repeat PercentageOverlap
 //AgtPtpMasterIpv4ClockPool SetDestinationAddressIncrementingRange
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetDestinationAddressIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDestinationAddressIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetDestinationAddressList ()(string, error) {
 //parameters: SessionPoolHandle IpAddressList
 //AgtPtpMasterIpv4ClockPool SetDestinationAddressList
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) GetDestinationAddressList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetDestinationAddressList
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetPeerDestinationAddress () error {
 //parameters: SessionPoolHandle PeerAddress
 //AgtPtpMasterIpv4ClockPool SetPeerDestinationAddress
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPeerDestinationAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPeerDestinationAddress
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) EnableUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool EnableUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv4ClockPool) DisableUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool DisableUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv4ClockPool) IsEnabledUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool IsEnabledUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv4ClockPool) SetUnicastGrantDuration () error {
 //parameters: SessionPoolHandle UnicastGrantDuration
 //AgtPtpMasterIpv4ClockPool SetUnicastGrantDuration
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetUnicastGrantDuration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetUnicastGrantDuration
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetNumberOfUnicastDestinations () error {
 //parameters: SessionPoolHandle UnicastDestinationCount
 //AgtPtpMasterIpv4ClockPool SetNumberOfUnicastDestinations
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetNumberOfUnicastDestinations ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetNumberOfUnicastDestinations
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetPeerPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpMasterIpv4ClockPool SetPeerPortNumber
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPeerPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPeerPortNumber
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetPeerClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpMasterIpv4ClockPool SetPeerClockId
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetPeerClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetPeerClockId
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetUdpSourcePort () error {
 //parameters: SessionPoolHandle FirstPort
 //AgtPtpMasterIpv4ClockPool SetUdpSourcePort
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetUdpSourcePort ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetUdpSourcePort
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetUdpSourcePortIncrementingRange () error {
 //parameters: SessionPoolHandle FirstPort Increment Repeat
 //AgtPtpMasterIpv4ClockPool SetUdpSourcePortIncrementingRange
 return nil
}

func(np *PtpMasterIpv4ClockPool) GetUdpSourcePortIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetUdpSourcePortIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) SetUdpSourcePortList ()(string, error) {
 //parameters: SessionPoolHandle UdpPortList
 //AgtPtpMasterIpv4ClockPool SetUdpSourcePortList
 return "", nil
}

func(np *PtpMasterIpv4ClockPool) GetUdpSourcePortList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv4ClockPool GetUdpSourcePortList
 return "", nil
}

