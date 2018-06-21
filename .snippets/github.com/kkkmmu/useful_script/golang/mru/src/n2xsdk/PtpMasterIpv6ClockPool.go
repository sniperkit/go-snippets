package n2xsdk

type PtpMasterIpv6ClockPool struct {
 Handler string
}

func(np *PtpMasterIpv6ClockPool) SetClockQuality () error {
 //parameters: SessionPoolHandle ClockClass ClockAccuracy OffsetScaledLogVariance Priority1 Priority2 AlternateMasterFlag StepsRemoved TimeSource
 //AgtPtpMasterIpv6ClockPool SetClockQuality
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetClockQuality ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetClockQuality
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnableAnnounce
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisableAnnounce
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledAnnounce
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetAnnounceAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval ReceiptTimeout
 //AgtPtpMasterIpv6ClockPool SetAnnounceAttributes
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetAnnounceAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetAnnounceAttributes
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnableSync
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisableSync
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledSync
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetSyncAttributes () error {
 //parameters: SessionPoolHandle TwoStepFlagMode LogMessageInterval
 //AgtPtpMasterIpv6ClockPool SetSyncAttributes
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetSyncAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetSyncAttributes
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnableFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisableFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetFollowUpAttributes () error {
 //parameters: SessionPoolHandle FollowUpDelay
 //AgtPtpMasterIpv6ClockPool SetFollowUpAttributes
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetFollowUpAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetFollowUpAttributes
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnableDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisableDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetDelayResponseAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval
 //AgtPtpMasterIpv6ClockPool SetDelayResponseAttributes
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetDelayResponseAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDelayResponseAttributes
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetDomainNumber () error {
 //parameters: SessionPoolHandle DomainNumber
 //AgtPtpMasterIpv6ClockPool SetDomainNumber
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetDomainNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDomainNumber
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpMasterIpv6ClockPool SetPortNumber
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPortNumber
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpMasterIpv6ClockPool SetClockId
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetClockId
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetClockIdIncrementingRange () error {
 //parameters: SessionPoolHandle FirstClockId Repeat Increment
 //AgtPtpMasterIpv6ClockPool SetClockIdIncrementingRange
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetClockIdIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetClockIdIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle ClockIdList
 //AgtPtpMasterIpv6ClockPool SetClockIdList
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) GetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetClockIdList
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetPathDelayMode () error {
 //parameters: SessionPoolHandle PathDelayMode
 //AgtPtpMasterIpv6ClockPool SetPathDelayMode
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPathDelayMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPathDelayMode
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnablePeerDelayRequest
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisablePeerDelayRequest
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledPeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledPeerDelayRequest
 return nil
}

func(np *PtpMasterIpv6ClockPool) EnablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnablePeerDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisablePeerDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledPeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledPeerDelayResponse
 return nil
}

func(np *PtpMasterIpv6ClockPool) EnablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledPeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledPeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetPeerToPeerPathDelayAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval TwoStepFlagMode FollowUpDelay
 //AgtPtpMasterIpv6ClockPool SetPeerToPeerPathDelayAttributes
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPeerToPeerPathDelayAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPeerToPeerPathDelayAttributes
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetTimestampCalibration () error {
 //parameters: SessionPoolHandle TxOffset RxOffset
 //AgtPtpMasterIpv6ClockPool SetTimestampCalibration
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetTimestampCalibration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetTimestampCalibration
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtPtpMasterIpv6ClockPool GetVlanPriority
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtPtpMasterIpv6ClockPool SetVlanPriority
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetDestinationAddressType () error {
 //parameters: SessionPoolHandle AddressType
 //AgtPtpMasterIpv6ClockPool SetDestinationAddressType
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetDestinationAddressType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDestinationAddressType
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) IsIpv6PriorityNoCodePointFieldSelected () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsIpv6PriorityNoCodePointFieldSelected
 return nil
}

func(np *PtpMasterIpv6ClockPool) SelectIpv6PriorityNoCodePointField () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool SelectIpv6PriorityNoCodePointField
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetIpv6Priority ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetIpv6Priority
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetIpv6Priority () error {
 //parameters: SessionPoolHandle Ipv6Priority
 //AgtPtpMasterIpv6ClockPool SetIpv6Priority
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsIpv6PriorityDiffServFieldSelected () error {
 //parameters: SessionPoolHandle DiffServPerHopBehaviorGroup
 //AgtPtpMasterIpv6ClockPool IsIpv6PriorityDiffServFieldSelected
 return nil
}

func(np *PtpMasterIpv6ClockPool) SelectIpv6PriorityDiffServField () error {
 //parameters: SessionPoolHandle DiffServPerHopBehaviorGroup
 //AgtPtpMasterIpv6ClockPool SelectIpv6PriorityDiffServField
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetIpv6PriorityDiffServValue ()(string, error) {
 //parameters: SessionPoolHandle DiffServCodePointField
 //AgtPtpMasterIpv6ClockPool GetIpv6PriorityDiffServValue
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetIpv6PriorityDiffServValue () error {
 //parameters: SessionPoolHandle DiffServCodePointConfigurableField Value
 //AgtPtpMasterIpv6ClockPool SetIpv6PriorityDiffServValue
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetDestinationAddress () error {
 //parameters: SessionPoolHandle FirstAddress
 //AgtPtpMasterIpv6ClockPool SetDestinationAddress
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetDestinationAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDestinationAddress
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetDestinationAddressIncrementingRange () error {
 //parameters: SessionPoolHandle FirstAddress PrefixLength Increment Repeat PercentageOverlap
 //AgtPtpMasterIpv6ClockPool SetDestinationAddressIncrementingRange
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetDestinationAddressIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDestinationAddressIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetDestinationAddressList ()(string, error) {
 //parameters: SessionPoolHandle IpAddressList
 //AgtPtpMasterIpv6ClockPool SetDestinationAddressList
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) GetDestinationAddressList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetDestinationAddressList
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetPeerDestinationAddress () error {
 //parameters: SessionPoolHandle PeerAddress
 //AgtPtpMasterIpv6ClockPool SetPeerDestinationAddress
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPeerDestinationAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPeerDestinationAddress
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) EnableUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool EnableUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv6ClockPool) DisableUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool DisableUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv6ClockPool) IsEnabledUnicastNegotiation () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool IsEnabledUnicastNegotiation
 return nil
}

func(np *PtpMasterIpv6ClockPool) SetUnicastGrantDuration () error {
 //parameters: SessionPoolHandle UnicastGrantDuration
 //AgtPtpMasterIpv6ClockPool SetUnicastGrantDuration
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetUnicastGrantDuration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetUnicastGrantDuration
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetNumberOfUnicastDestinations () error {
 //parameters: SessionPoolHandle UnicastDestinationCount
 //AgtPtpMasterIpv6ClockPool SetNumberOfUnicastDestinations
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetNumberOfUnicastDestinations ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetNumberOfUnicastDestinations
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetPeerPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpMasterIpv6ClockPool SetPeerPortNumber
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPeerPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPeerPortNumber
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetPeerClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpMasterIpv6ClockPool SetPeerClockId
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetPeerClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetPeerClockId
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetUdpSourcePort () error {
 //parameters: SessionPoolHandle FirstPort
 //AgtPtpMasterIpv6ClockPool SetUdpSourcePort
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetUdpSourcePort ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetUdpSourcePort
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetUdpSourcePortIncrementingRange () error {
 //parameters: SessionPoolHandle FirstPort Increment Repeat
 //AgtPtpMasterIpv6ClockPool SetUdpSourcePortIncrementingRange
 return nil
}

func(np *PtpMasterIpv6ClockPool) GetUdpSourcePortIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetUdpSourcePortIncrementingRange
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) SetUdpSourcePortList ()(string, error) {
 //parameters: SessionPoolHandle UdpPortList
 //AgtPtpMasterIpv6ClockPool SetUdpSourcePortList
 return "", nil
}

func(np *PtpMasterIpv6ClockPool) GetUdpSourcePortList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterIpv6ClockPool GetUdpSourcePortList
 return "", nil
}

